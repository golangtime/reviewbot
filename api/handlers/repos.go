package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golangtime/reviewbot/api"
)

func (h *Handler) AddRepo(w http.ResponseWriter, r *http.Request) {
	logger := h.logger
	var req api.AddRepoRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errorResponse(w, logger, err, http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		errorResponse(w, logger, fmt.Errorf("empty repository name"), http.StatusBadRequest)
		return
	}

	if req.Owner == "" {
		errorResponse(w, logger, fmt.Errorf("empty owner name"), http.StatusBadRequest)
		return
	}

	err = h.repo.AddRepo(h.db, req.Owner, req.Name, req.MinApprovals, req.Provider)
	if err != nil {
		errorResponse(w, logger, err, http.StatusInternalServerError)
		return
	}

	resp := api.AddRepoResponse{
		Success: true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&resp)
}

func (h *Handler) RemoveRepo(w http.ResponseWriter, r *http.Request) {
	logger := h.logger

	var req api.RemoveRepoRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errorResponse(w, logger, err, http.StatusBadRequest)
		return
	}

	err = h.repo.RemoveRepo(h.db, req.Owner, req.Name)
	if err != nil {
		errorResponse(w, logger, err, http.StatusInternalServerError)
		return
	}

	resp := api.RemoveRepoResponse{
		Success: true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&resp)
}

func (h *Handler) ListRepos(w http.ResponseWriter, r *http.Request) {
	logger := h.logger

	var req api.ListReposRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errorResponse(w, logger, err, http.StatusBadRequest)
		return
	}

	repos, err := h.repo.ListRepos(h.db, req.Owner)
	if err != nil {
		errorResponse(w, logger, err, http.StatusInternalServerError)
		return
	}

	repoResponse := make([]*api.Repo, 0, len(repos))

	for _, r := range repos {
		repoResponse = append(repoResponse, &api.Repo{
			Owner:        r.Owner,
			Name:         r.Name,
			Provider:     r.Provider,
			MinApprovals: r.MinApprovals,
		})
	}

	resp := api.ListReposResponse{
		Repos:   repoResponse,
		Count:   len(repoResponse),
		Success: true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&resp)
}
