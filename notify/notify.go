package notify

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golangtime/reviewbot/db"
)

type Notification struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Link   string `json:"link"`
}

type Sender struct {
	logger *slog.Logger
	db     *sql.DB
	repo   *db.Repository

	// senders
	emailSender  NotificationSender
	pachcaSender NotificationSender
}

type NotificationSender interface {
	Send(providerID string, chatID int64, link string) error
}

type SenderOptions struct {
	EmailSender  NotificationSender
	PachcaSender NotificationSender
}

func NewSender(dbConn *sql.DB, logger *slog.Logger, opts *SenderOptions) *Sender {
	repo := &db.Repository{}

	return &Sender{
		repo:         repo,
		db:           dbConn,
		logger:       logger,
		emailSender:  opts.EmailSender,
		pachcaSender: opts.PachcaSender,
	}
}

func (s *Sender) Start() error {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			notifications, err := s.repo.ListPendingNotifications(s.db)
			if err != nil {
				s.logger.Error("ListPendingNotifications error", "error", err)
				continue
			}

			if len(notifications) == 0 {
				continue
			}

			for _, notification := range notifications {
				err := s.Send(&Notification{
					ID:     notification.ID,
					UserID: notification.UserID,
					Link:   notification.Link,
				})
				if err != nil {
					s.logger.Error("send error", "error", err)
				}
			}
		}
	}
}

func (s *Sender) Send(notification *Notification) error {
	tx, err := s.db.Begin()
	if err != nil {
		return nil
	}

	defer func() {
		if err == nil {
			err = tx.Commit()
			if err != nil {
				s.logger.Error("commit error", "error", err)
			}
			return
		}

		err = tx.Rollback()
		if err != nil {
			s.logger.Error("rollback error", "error", err)
		}
	}()

	// TODO - locking
	_, err = s.db.Exec("SELECT 1 FROM notification_queue WHERE id = $1", notification.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Info("skip notification, already sent", "notification_id", notification.ID)
			return nil
		}
		return err
	}

	s.logger.Debug("send notification", "notification_id", notification.ID)

	notificationRule, err := s.repo.ListNotificationRuleByUser(s.db, notification.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Info("no notification rule for user found", "notification_id", notification.ID, "user_id", notification.UserID)
			return nil
		}
		return fmt.Errorf("get notification rule for user error: %w", err)
	}

	defer func() {
		_, err = s.db.Exec("DELETE FROM notification_queue WHERE id = $1", notification.ID)
		if err != nil {
			s.logger.Error("delete notification error", "error", err)
			return
		}
	}()

	switch notificationRule.NotificationType {
	case "email":
		return s.SendEmail(notificationRule.ProviderID, notification.Link)
	case "pachca":
		return s.SendPachca(notificationRule.ProviderID, notificationRule.ChatID, notification.Link)
	default:
		s.logger.Warn("notification provider not found", "notification_type", notificationRule.NotificationType)
	}

	return nil
}

func (s *Sender) SendEmail(email, link string) error {
	return s.emailSender.Send(email, 0, link)
}

func (s *Sender) SendPachca(providerID string, chatID int64, link string) error {
	return s.pachcaSender.Send(providerID, chatID, link)
}
