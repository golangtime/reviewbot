package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/golangtime/reviewbot/api"
	"github.com/golangtime/reviewbot/api/dto"
	"github.com/golangtime/reviewbot/client"
	"github.com/golangtime/reviewbot/db"
)

type GitClients struct {
	Github    client.GitClient
	Bitbucket client.GitClient
}

type Handler struct {
	db          *sql.DB
	repo        db.Repo
	logger      *slog.Logger
	gitClients  GitClients
	pachcaToken string
}

func NewHandler(
	dbConn *sql.DB,
	repo db.Repo,
	logger *slog.Logger,
	gitClients GitClients,
	pachcaToken string,
) *Handler {
	return &Handler{
		gitClients:  gitClients,
		db:          dbConn,
		repo:        repo,
		logger:      logger,
		pachcaToken: pachcaToken,
	}
}

func (h *Handler) ListPullRequests(w http.ResponseWriter, r *http.Request) {
	logger := h.logger

	var req dto.ListPullRequests

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Error("request error", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var gitClient client.GitClient

	switch req.Provider {
	case "github":
		gitClient = h.gitClients.Github
	case "bitbucket":
		gitClient = h.gitClients.Bitbucket
	}

	pullRequests, err := gitClient.ListPullRequests(req.Owner, req.Repo)
	if err != nil {
		logger.Error("list pull requests error", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		resp := api.AddRepoResponse{
			Success: false,
		}
		json.NewEncoder(w).Encode(&resp)
		return
	}

	var result []dto.PullRequest
	for _, pr := range pullRequests {
		var reviewers []dto.PullRequestReviewer
		for _, reviewer := range pr.Reviewers {
			reviewers = append(reviewers, dto.PullRequestReviewer{
				ID:    reviewer.ID,
				Email: reviewer.Email,
			})
		}

		result = append(result, dto.PullRequest{
			ID:        pr.ExternalID,
			Link:      pr.Link,
			Reviewers: reviewers,
		})
	}

	resp := dto.ListPullRequestsResponse{
		Result: result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&resp)
}

func (h *Handler) UpdateNotificationRule(w http.ResponseWriter, r *http.Request) {
	logger := h.logger

	var req dto.UpdateNotificationRuleRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Error("request error", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.repo.UpdateNotificationRule(h.db, req.UserID, req.NotificationType, req.ProviderID, req.Priority)
	if err != nil {
		logger.Error("update notification rule error", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		resp := api.AddRepoResponse{
			Success: false,
		}
		json.NewEncoder(w).Encode(&resp)
		return
	}

	resp := dto.UpdateNotificationRuleResponse{
		Success: true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&resp)
}
