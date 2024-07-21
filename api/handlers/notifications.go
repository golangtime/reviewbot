package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/golangtime/reviewbot/api"
)

func (v *Handler) ListPendingNotifications(w http.ResponseWriter, r *http.Request) {
	logger := v.logger

	var req api.ListNotificationRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Error("list notifications request decode error", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		resp := api.AddRepoResponse{
			Success: false,
		}
		json.NewEncoder(w).Encode(&resp)
		return
	}

	records, err := v.repo.ListPendingNotifications(v.db)
	if err != nil {
		logger.Error("list notifications error", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		resp := api.AddRepoResponse{
			Success: false,
		}
		json.NewEncoder(w).Encode(&resp)
		return
	}

	response := make([]*api.Notification, 0, len(records))

	for _, r := range records {
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
}
