package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/golangtime/reviewbot/api"
)

func (h *Handler) AddNotificationRule(w http.ResponseWriter, r *http.Request) {
	logger := h.logger
	if r.Method == "PUT" {
		h.UpdateNotificationRule(w, r)
		return
	}

	var req api.AddNotificationRuleRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errorResponse(w, logger, err, http.StatusBadRequest)
		return
	}

	err = h.repo.AddNotificationRule(h.db, req.UserID, req.NotificationType, req.ProviderID, req.ChatID, req.Priority)
	if err != nil {
		errorResponse(w, logger, err, http.StatusInternalServerError)
		return
	}

	resp := api.AddNotificationRuleResponse{
		Success: true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&resp)
}

func (h *Handler) RemoveNotificationRule(w http.ResponseWriter, r *http.Request) {
	logger := h.logger

	var req api.RemoveNotitifcationRuleRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Error("request error", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.repo.RemoveNotificationRule(h.db, req.ID)
	if err != nil {
		logger.Error("remove notification rule error", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		resp := api.RemoveNotitifcationRuleResponse{
			Success: false,
		}
		json.NewEncoder(w).Encode(&resp)
		return
	}

	resp := api.RemoveNotitifcationRuleResponse{
		Success: true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&resp)
}

func (h *Handler) ListNotificationRules(w http.ResponseWriter, r *http.Request) {
	logger := h.logger
	var req api.ListNotificationRulesRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errorResponse(w, logger, err, http.StatusBadRequest)
		return
	}

	result, err := h.repo.ListNotificationRules(h.db)
	if err != nil {
		errorResponse(w, logger, err, http.StatusInternalServerError)
		return
	}

	list := make([]*api.NotificationRule, 0, len(result))
	for _, r := range result {
		list = append(list, &api.NotificationRule{
			ID:               r.ID,
			UserID:           r.UserID,
			NotificationType: r.NotificationType,
			ProviderID:       r.ProviderID,
			ChatID:           r.ChatID,
			Priority:         r.Priority,
		})
	}

	resp := api.ListNotificationRulesResponse{
		Result: list,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&resp)
}
