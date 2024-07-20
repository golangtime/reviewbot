package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/golangtime/reviewbot/api"
)

func (h *Handler) RemoveRepo(w http.ResponseWriter, r *http.Request) {
	logger := h.logger

	var req api.RemoveRepoRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Error("request error", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.repo.RemoveRepo(h.db, req.Owner, req.Name)
	if err != nil {
		logger.Error("remove repository error", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		resp := api.RemoveRepoResponse{
			Success: false,
		}
		json.NewEncoder(w).Encode(&resp)
		return
	}

	resp := api.RemoveRepoResponse{
		Success: true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&resp)
}
