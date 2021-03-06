package recurring

import (
	"time"

	"github.com/lloydmeta/tasques/internal/domain/task"
)

type Schedule interface {
	Next(t time.Time) time.Time
}

type ScheduleParser interface {
	Parse(spec string) (Schedule, error)
}

// A thin wrapper interface around Go Cron
type Scheduler interface {
	ScheduleParser

	// Schedules the given recurring Task to be inserted at intervals according
	//to the Task's ScheduleExpression.
	//
	// Note that it takes a full Task (not a ref) because the scheduler likely
	// does scheduling asynchronously, and an address is dangerous.
	//
	// This function will act like an upsert: if there already exists a Task
	// that has the same RecurringTaskId, the existing one is first unscheduled, then the new
	// one scheduled.
	//
	// Since we assume the ScheduleExpression is _valid_, there should be no
	// errors, but what the hey 🤷🏻‍♂️
	Schedule(task Task) error

	// Stops the given recurring Task from being inserted at intervals
	//
	// Returns true if it was previously scheduled, false otherwise
	Unschedule(taskId task.RecurringTaskId) bool

	// Starts the Scheduler in its own Go routing
	Start()

	// Stops the Scheduler
	Stop()
}
