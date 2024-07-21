package main

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	_ "github.com/mattn/go-sqlite3"

	"github.com/golangtime/reviewbot/api/handlers"
	"github.com/golangtime/reviewbot/bot"
	"github.com/golangtime/reviewbot/client/bitbucket"
	"github.com/golangtime/reviewbot/client/github"
	"github.com/golangtime/reviewbot/config"
	"github.com/golangtime/reviewbot/db"
	"github.com/golangtime/reviewbot/intergrations/email"
	"github.com/golangtime/reviewbot/intergrations/pachca"
	"github.com/golangtime/reviewbot/job"
	"github.com/golangtime/reviewbot/scheduler"
	"github.com/golangtime/reviewbot/scheduler/cron"
)

func StartAPI(cfg *config.Config, logger *slog.Logger, database *sql.DB, gitClients handlers.GitClients) error {
	repo := db.Repository{}

	handler := handlers.NewHandler(database, repo, logger, gitClients, cfg.Pachca.Token)

	http.HandleFunc("/repos/add", func(w http.ResponseWriter, r *http.Request) {
		handler.AddRepo(w, r)
	})

	http.HandleFunc("/repos/list", func(w http.ResponseWriter, r *http.Request) {
		handler.ListRepos(w, r)
	})

	http.HandleFunc("/repos/remove", func(w http.ResponseWriter, r *http.Request) {
		handler.RemoveRepo(w, r)
	})

	http.HandleFunc("/pull_requests", func(w http.ResponseWriter, r *http.Request) {
		handler.ListPullRequests(w, r)
	})

	http.HandleFunc("/notification/pending", func(w http.ResponseWriter, r *http.Request) {
		handler.ListPendingNotifications(w, r)
	})

	http.HandleFunc("/notification/rules", func(w http.ResponseWriter, r *http.Request) {
		handler.ListNotificationRules(w, r)
	})

	http.HandleFunc("/notification/add_rule", func(w http.ResponseWriter, r *http.Request) {
		handler.AddNotificationRule(w, r)
	})

	http.HandleFunc("/notification/remove_rule", func(w http.ResponseWriter, r *http.Request) {
		handler.RemoveNotificationRule(w, r)
	})

	http.HandleFunc("/pachca/user", func(w http.ResponseWriter, r *http.Request) {
		handler.FindPachcaUser(w, r)
	})

	http.HandleFunc("/pachca/chat", func(w http.ResponseWriter, r *http.Request) {
		handler.FindChat(w, r)
	})

	err := http.ListenAndServe(":8000", nil)
	return err
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	dbConn, err := sql.Open("sqlite3", "bot.db")
	if err != nil {
		panic(err)
	}

	db.MigrateDB(dbConn)

	gitClient := github.New()
	bitbucketClient := bitbucket.New(cfg.Bitbucket.URL, cfg.Bitbucket.User, cfg.Bitbucket.Password)

	go StartAPI(cfg, logger, dbConn, handlers.GitClients{
		Github:    gitClient,
		Bitbucket: bitbucketClient,
	})

	jobFunc := job.NewJob(dbConn, logger, job.GitClients{
		Github:    gitClient,
		Bitbucket: bitbucketClient,
	})

	botScheduler, err := cron.NewCron(cfg.Schedule.TestMode, jobFunc)
	if err != nil {
		logger.Error("scheduler error", "error", err)
		os.Exit(1)
	}
	defer botScheduler.Stop()

	log.Printf("%+v", cfg)

	emailSender := email.NewEmailSender(logger, cfg.Email.From, cfg.Email.User, cfg.Email.Password)
	bot.RegisterNotificationSender("email", emailSender)

	pachcaSender := pachca.NewPachcaSender(logger, cfg.Pachca.Token)
	bot.RegisterNotificationSender("pachca", pachcaSender)

	bot := bot.New(dbConn, botScheduler, logger)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	var schedule []scheduler.ScheduleRecord
	for _, sh := range cfg.Schedule.Default {
		schedule = append(schedule, scheduler.ScheduleRecord{
			Hour: sh[0], Minute: sh[1], Second: sh[2],
		})
	}

	logger.Info("starting bot", "schedule", schedule)
	bot.Start(scheduler.Schedule{
		Records: schedule,
	})

	select {
	case sig := <-sigCh:
		logger.Info("signal received", "signal", sig.String())
	}

	logger.Info("bot stopped")
}
