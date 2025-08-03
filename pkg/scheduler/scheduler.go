package scheduler

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
)

type Scheduler struct {
	scheduler gocron.Scheduler
}

func New() (*Scheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	s.Start()

	return &Scheduler{scheduler: s}, nil
}
func (s *Scheduler) Stop() error {
	return s.scheduler.Shutdown()
}

func (s *Scheduler) NewJob(cron string, fn any, params ...any) error {
	_, err := s.scheduler.NewJob(
		gocron.CronJob(cron, true),
		gocron.NewTask(fn, params...),
	)
	return err
}

// run every second
func (s *Scheduler) EverySecond(fn any, params ...any) error {
	return s.NewJob("* * * * * *", fn, params...)
}

// run every minute
func (s *Scheduler) EveryMinute(fn any, params ...any) error {
	return s.NewJob("0 * * * * *", fn, params...)
}

// run every hour
func (s *Scheduler) EveryHour(fn any, params ...any) error {
	return s.NewJob("0 0 * * * *", fn, params...)
}

// run every day
func (s *Scheduler) EveryDay(fn any, params ...any) error {
	return s.NewJob("0 0 0 * * *", fn, params...)
}

// run every week
func (s *Scheduler) EveryWeek(fn any, params ...any) error {
	return s.NewJob("0 0 0 * * 0", fn, params...)
}

// run every month
func (s *Scheduler) EveryMonth(fn any, params ...any) error {
	return s.NewJob("0 0 0 1 * *", fn, params...)
}

// run every year
func (s *Scheduler) EveryYear(fn any, params ...any) error {
	return s.NewJob("0 0 0 1 1 *", fn, params...)
}

// run every day at specific time
func (s *Scheduler) DailyAt(hour, minute int, fn any, params ...any) error {
	return s.NewJob(fmt.Sprintf("%d %d * * *", minute, hour), fn, params...)
}

// run every week at specific time
func (s *Scheduler) WeeklyAt(day int, hour, minute int, fn any, params ...any) error {
	return s.NewJob(fmt.Sprintf("%d %d * * %d", minute, hour, day), fn, params...)
}

// run every month at specific time
func (s *Scheduler) MonthlyAt(day, hour, minute int, fn any, params ...any) error {
	return s.NewJob(fmt.Sprintf("%d %d %d * *", minute, hour, day), fn, params...)
}

// run every year at specific time
func (s *Scheduler) YearlyAt(month, day, hour, minute int, fn any, params ...any) error {
	return s.NewJob(fmt.Sprintf("%d %d %d %d *", minute, hour, day, month), fn, params...)
}

// run after specific duration
func (s *Scheduler) Duration(duration time.Duration, fn any, params ...any) (gocron.Job, error) {
	return s.scheduler.NewJob(
		gocron.DurationJob(
			duration,
		),
		gocron.NewTask(fn, params...),
	)
}

// run after specific duration
func (s *Scheduler) After(duration time.Duration, fn any, params ...any) (gocron.Job, error) {
	return s.scheduler.NewJob(
		gocron.DurationJob(
			duration,
		),
		gocron.NewTask(fn, params...),
	)
}

// run after seconds
func (s *Scheduler) AfterSecound(fn any, params ...any) (gocron.Job, error) {
	return s.scheduler.NewJob(
		gocron.DurationJob(
			1*time.Second,
		),
		gocron.NewTask(fn, params...),
	)
}

// run after minute
func (s *Scheduler) AfterMinute(fn any, params ...any) (gocron.Job, error) {
	return s.scheduler.NewJob(
		gocron.DurationJob(
			1*time.Minute,
		),
		gocron.NewTask(fn, params...),
	)
}

// run after hour
func (s *Scheduler) AfterHour(fn any, params ...any) (gocron.Job, error) {
	return s.scheduler.NewJob(
		gocron.DurationJob(
			1*time.Hour,
		),
		gocron.NewTask(fn, params...),
	)
}

// run after day
func (s *Scheduler) AfterDay(fn any, params ...any) (gocron.Job, error) {
	return s.scheduler.NewJob(
		gocron.DurationJob(
			1*time.Hour*24,
		),
		gocron.NewTask(fn, params...),
	)
}

// run after week
func (s *Scheduler) AfterWeek(fn any, params ...any) (gocron.Job, error) {
	return s.scheduler.NewJob(
		gocron.DurationJob(
			1*time.Hour*24*7,
		),
		gocron.NewTask(fn, params...),
	)
}

// run after month
func (s *Scheduler) AfterMonth(fn any, params ...any) (gocron.Job, error) {
	return s.scheduler.NewJob(
		gocron.DurationJob(
			1*time.Hour*24*30,
		),
		gocron.NewTask(fn, params...),
	)
}

// run after year
func (s *Scheduler) AfterYear(fn any, params ...any) (gocron.Job, error) {
	return s.scheduler.NewJob(
		gocron.DurationJob(
			1*time.Hour*24*365,
		),
		gocron.NewTask(fn, params...),
	)
}
