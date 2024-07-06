package scheduler

type Schedule struct {
	Records []ScheduleRecord
}

type ScheduleRecord struct {
	Hour   int
	Minute int
	Second int
}

type Scheduler interface {
	Start(Schedule) error
	Stop()
}
