package bot

import (
	"database/sql"
	"log/slog"

	"github.com/golangtime/reviewbot/notify"
	"github.com/golangtime/reviewbot/scheduler"
)

type Bot struct {
	logger    *slog.Logger
	scheduler scheduler.Scheduler
}

var notificationSenderRegistry = map[string]notify.NotificationSender{}

func RegisterNotificationSender(notificationType string, sender notify.NotificationSender) {
	notificationSenderRegistry[notificationType] = sender
}

func New(dbConn *sql.DB, scheduler scheduler.Scheduler, logger *slog.Logger) *Bot {
	backgroundNotificationSender := notify.NewSender(dbConn, logger, &notify.SenderOptions{
		EmailSender:  notificationSenderRegistry["email"],
		PachcaSender: notificationSenderRegistry["pachca"],
	})

	go backgroundNotificationSender.Start()

	return &Bot{
		logger:    logger,
		scheduler: scheduler,
	}
}

func (b *Bot) Start(schedule scheduler.Schedule) {
	b.scheduler.Start(schedule)
}

func (b *Bot) Stop() error {
	return nil
}
