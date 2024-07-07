package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/golangtime/reviewbot/api"
	"github.com/golangtime/reviewbot/api/dto"
	"github.com/golangtime/reviewbot/db"
)

type Handler struct {
	db     *sql.DB
	repo   db.Repo
	logger *slog.Logger
}

func NewHandler(dbConn *sql.DB, repo db.Repo, logger *slog.Logger) *Handler {
	return &Handler{
		db:     dbConn,
		repo:   repo,
		logger: logger,
	}
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
		logger.Error("add notification rule error", "error", err)
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
