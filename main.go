package main

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	_ "github.com/mattn/go-sqlite3"

	"github.com/golangtime/reviewbot/api"
	"github.com/golangtime/reviewbot/bot"
	"github.com/golangtime/reviewbot/client/github"
	"github.com/golangtime/reviewbot/db"
	"github.com/golangtime/reviewbot/job"
	"github.com/golangtime/reviewbot/scheduler"
	"github.com/golangtime/reviewbot/scheduler/cron"
)

func StartAPI(logger *slog.Logger, database *sql.DB) error {
	repo := db.Repository{}

	ctrl := api.NewAPIV1(database, repo)

	http.HandleFunc("/repos/add", func(w http.ResponseWriter, r *http.Request) {
		var req api.AddRepoRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			logger.Error("request error", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = ctrl.AddRepo(req.Owner, req.Name, req.MinApprovals)
		if err != nil {
			logger.Error("add repo error", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			resp := api.AddRepoResponse{
				Success: false,
			}
			json.NewEncoder(w).Encode(&resp)
			return
		}

		resp := api.AddRepoResponse{
			Success: true,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&resp)
	})

	http.HandleFunc("/repos/list", func(w http.ResponseWriter, r *http.Request) {
		var req api.ListReposRequest

		json.NewDecoder(r.Body).Decode(&req)

		repos, err := ctrl.ListRepo(req.Owner)
		if err != nil {
			logger.Error("list repo error", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			resp := api.AddRepoResponse{
				Success: false,
			}
			json.NewEncoder(w).Encode(&resp)
			return
		}

		var repoResponse []*api.Repo
		for _, r := range repos {
			repoResponse = append(repoResponse, &api.Repo{
				Name:         r.Name,
				MinApprovals: r.MinApprovals,
			})
		}

		resp := api.ListReposResponse{
			Repos: repoResponse,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&resp)
	})

	err := http.ListenAndServe(":8000", nil)
	return err
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	dbConn, err := sql.Open("sqlite3", "bot.db")
	if err != nil {
		panic(err)
	}

	db.MigrateDB(dbConn)

	go StartAPI(logger, dbConn)

	gitClient := github.New()

	jobFunc := job.NewJob(dbConn, logger, gitClient)

	botScheduler, err := cron.NewCron(jobFunc)
	if err != nil {
		logger.Error("scheduler error", "error", err)
		os.Exit(1)
	}
	defer botScheduler.Stop()

	bot := bot.New(botScheduler, logger)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	logger.Info("starting bot")
	bot.Start(scheduler.Schedule{
		Records: []scheduler.ScheduleRecord{{0, 0, 0}},
	})

	select {
	case sig := <-sigCh:
		logger.Info("signal received", "signal", sig.String())
	}

	logger.Info("bot stopped")
}
