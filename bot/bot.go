package bot

import (
	"log/slog"

	"github.com/golangtime/reviewbot/scheduler"
)

type Bot struct {
	logger    *slog.Logger
	scheduler scheduler.Scheduler
}

func New(scheduler scheduler.Scheduler, logger *slog.Logger) *Bot {
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
