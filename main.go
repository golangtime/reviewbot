package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	_ "github.com/mattn/go-sqlite3"

	"github.com/golangtime/reviewbot/api"
	"github.com/golangtime/reviewbot/api/handlers"
	"github.com/golangtime/reviewbot/bot"
	"github.com/golangtime/reviewbot/client/bitbucket"
	"github.com/golangtime/reviewbot/client/github"
	"github.com/golangtime/reviewbot/config"
	"github.com/golangtime/reviewbot/db"
	"github.com/golangtime/reviewbot/intergrations/email"
	"github.com/golangtime/reviewbot/job"
	"github.com/golangtime/reviewbot/scheduler"
	"github.com/golangtime/reviewbot/scheduler/cron"
)

func StartAPI(logger *slog.Logger, database *sql.DB) error {
	repo := db.Repository{}

	handler := handlers.NewHandler(database, repo, logger)

	ctrl := api.NewAPIV1(database, repo)

	http.HandleFunc("/repos/add", func(w http.ResponseWriter, r *http.Request) {
		var req api.AddRepoRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			logger.Error("request error", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = ctrl.AddRepo(req.Owner, req.Name, req.Provider, req.MinApprovals)
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
				Provider:     r.Provider,
				MinApprovals: r.MinApprovals,
			})
		}

		resp := api.ListReposResponse{
			Repos: repoResponse,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&resp)
	})

	http.HandleFunc("/notification/pending", func(w http.ResponseWriter, r *http.Request) {
		var req api.ListNotificationRequest

		json.NewDecoder(r.Body).Decode(&req)

		result, err := ctrl.ListPendingNotifications(req.QueueType)
		if err != nil {
			logger.Error("list notifications error", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			resp := api.AddRepoResponse{
				Success: false,
			}
			json.NewEncoder(w).Encode(&resp)
			return
		}

		var response []*api.Notification
		for _, r := range result {
			response = append(response, &api.Notification{
				Recepient:   r.Recepient,
				Link:        r.Link,
				UserID:      r.UserID,
				CreatedAt:   r.CreatedAt,
				ReservedFor: r.ReservedFor,
			})
		}

		resp := api.ListNotificationResponse{
			Notifications: response,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&resp)
	})

	http.HandleFunc("/notification/rules", func(w http.ResponseWriter, r *http.Request) {
		var req api.ListNotificationRulesRequest

		json.NewDecoder(r.Body).Decode(&req)

		result, err := ctrl.ListNotificationRules()
		if err != nil {
			logger.Error("list notification rules error", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			resp := api.AddRepoResponse{
				Success: false,
			}
			json.NewEncoder(w).Encode(&resp)
			return
		}

		var response []*api.NotificationRule
		for _, r := range result {
			response = append(response, &api.NotificationRule{
				UserID:           r.UserID,
				NotificationType: r.NotificationType,
				ProviderID:       r.ProviderID,
				Priority:         r.Priority,
			})
		}

		resp := api.ListNotificationRulesResponse{
			Result: response,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&resp)
	})

	http.HandleFunc("/notification/add_rule", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			handler.UpdateNotificationRule(w, r)
			return
		}

		var req api.AddNotificationRuleRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			logger.Error("request error", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = ctrl.AddNotificationRule(req.UserID, req.NotificationType, req.ProviderID, req.Priority)
		if err != nil {
			logger.Error("add notification rule error", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			resp := api.AddRepoResponse{
				Success: false,
			}
			json.NewEncoder(w).Encode(&resp)
			return
		}

		resp := api.AddNotificationRuleResponse{
			Success: true,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&resp)
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

	go StartAPI(logger, dbConn)

	gitClient := github.New()
	bitbucketClient := bitbucket.New(cfg.Bitbucket.URL, cfg.Bitbucket.User, cfg.Bitbucket.Password)

	jobFunc := job.NewJob(dbConn, logger, job.GitClients{
		Github:    gitClient,
		Bitbucket: bitbucketClient,
	})

	botScheduler, err := cron.NewCron(jobFunc)
	if err != nil {
		logger.Error("scheduler error", "error", err)
		os.Exit(1)
	}
	defer botScheduler.Stop()

	log.Printf("%+v", cfg)

	emailSender := email.NewEmailSender(logger, cfg.Email.From, cfg.Email.User, cfg.Email.Password)

	bot.RegisterNotificationSender("email", emailSender)

	bot := bot.New(dbConn, botScheduler, logger)

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
