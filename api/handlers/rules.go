package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/golangtime/reviewbot/api"
)

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
