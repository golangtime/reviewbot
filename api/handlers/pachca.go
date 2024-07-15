package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type FindUser struct {
	Per int `json:"per"`
}

func (h *Handler) FindPachcaUser(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")

	httpClient := &http.Client{}

	u, _ := url.Parse("https://api.pachca.com/api/shared/v1/users")

	u.Query().Add("per", "100")
	u.Query().Add("query", query)

	u.RawQuery = r.URL.Query().Encode()

	request, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		log.Println("pachca request error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", h.pachcaToken))

	resp, err := httpClient.Do(request)
	if err != nil {
		log.Println("pachca error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("pachca: != 200", resp.StatusCode)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	w.Header().Add("Content-Type", "application/json")

	w.Write(body)
}

func (h *Handler) FindChat(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")

	httpClient := &http.Client{}

	u, _ := url.Parse("https://api.pachca.com/api/shared/v1/chats")

	u.Query().Add("per", "100")
	u.Query().Add("page", page)

	u.RawQuery = r.URL.Query().Encode()

	request, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		log.Println("pachca request error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", h.pachcaToken))

	resp, err := httpClient.Do(request)
	if err != nil {
		log.Println("pachca error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("pachca: != 200, was = %d\n", resp.StatusCode)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	w.Header().Add("Content-Type", "application/json")

	w.Write(body)
}
