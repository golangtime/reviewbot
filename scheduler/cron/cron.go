package cron

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/golangtime/reviewbot/scheduler"
)

type CronScheduler struct {
	TestMode bool
	sch      gocron.Scheduler
	f        func()
}

func NewCron(f func()) (*CronScheduler, error) {
	sch, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}
	return &CronScheduler{
		TestMode: false,
		sch:      sch,
		f:        f,
	}, nil
}

func (s *CronScheduler) runTest() error {
	s.f()

	_, err := s.sch.NewJob(
		gocron.DurationJob(
			3600*time.Second,
		),
		gocron.NewTask(s.f),
	)
	s.sch.Start()
	return err
}

func (s *CronScheduler) runDaily(schedule scheduler.Schedule) error {
	cronSchedule := []gocron.AtTime{}

	for _, r := range schedule.Records {
		cronSchedule = append(cronSchedule, gocron.NewAtTime(
			uint(r.Hour), uint(r.Minute), uint(r.Second),
		))
	}

	_, err := s.sch.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(cronSchedule[0], cronSchedule[1:]...),
		),
		gocron.NewTask(s.f),
		gocron.WithStartAt(gocron.WithStartImmediately()),
	)

	s.sch.Start()
	return err
}

func (s *CronScheduler) Start(schedule scheduler.Schedule) error {
	if s.TestMode {
		s.runTest()
	} else {
		s.runDaily(schedule)
	}

	return nil
}

func (s *CronScheduler) Stop() {
	s.sch.Shutdown()
}
