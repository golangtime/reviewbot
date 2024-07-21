package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/golangtime/reviewbot/api"
)

func errorResponse(w http.ResponseWriter, logger *slog.Logger, err error, statusCode int) {
	logger.Error("request error", "error", err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := api.AddRepoResponse{
		Success: false,
	}

	_ = json.NewEncoder(w).Encode(&resp)
}
