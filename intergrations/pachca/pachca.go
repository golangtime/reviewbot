package pachca

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
)

type PachcaSender struct {
	token      string
	logger     *slog.Logger
	httpClient *http.Client
}

func NewPachcaSender(logger *slog.Logger, token string) *PachcaSender {
	httpClient := &http.Client{}
	return &PachcaSender{
		token:      token,
		logger:     logger,
		httpClient: httpClient,
	}
}

type PachcaMessage struct {
	Message SendMessage `json:"message"`
}

type SendMessage struct {
	EntityType string `json:"entity_type"`
	EntityID   int64  `json:"entity_id"`
	Content    string `json:"content"`
}

type UserData struct {
	Data struct {
		Nickname string `json:"nickname"`
	} `json:"data"`
}

func (s *PachcaSender) GetUserNickname(userID int) string {
	httpClient := &http.Client{}

	u, _ := url.Parse(fmt.Sprintf("https://api.pachca.com/api/shared/v1/users/%d", userID))

	request, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		log.Println("pachca request error:", err)
		return ""
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token))

	resp, err := httpClient.Do(request)
	if err != nil {
		log.Println("pachca error: ", err)
		return ""
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("pachca: != 200, was = %d\n", resp.StatusCode)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var user UserData
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Println("unmarshal error: ", err)
		return ""
	}

	return user.Data.Nickname
}

func (s *PachcaSender) Send(providerID string, chatID int64, link string) error {
	entityID, err := strconv.Atoi(providerID)
	if err != nil {
		return err
	}

	var msg SendMessage

	if chatID > 0 {
		nickname := s.GetUserNickname(entityID)

		msg = SendMessage{
			EntityType: "discussion",
			EntityID:   int64(chatID),
			Content:    fmt.Sprintf("@%s Hi! You need to review the following pull request: %s", nickname, link),
		}
	} else {
		msg = SendMessage{
			EntityType: "user",
			EntityID:   int64(entityID),
			Content:    fmt.Sprintf("You need to review the following pull request: %s", link),
		}
	}

	log.Printf("send pachca message: %+v\n", msg)

	body, err := json.Marshal(&PachcaMessage{msg})
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, "https://api.pachca.com/api/shared/v1/messages", bytes.NewReader(body))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token))

	resp, err := s.httpClient.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		log.Println(string(respBody))
		return fmt.Errorf("pachca error: !=200, was = %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Println(string(respBody))

	return nil
}
