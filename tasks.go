package releaseManage

import (
	"context"
	"time"
)

type RetryStep struct {
	// Interval is the frequency that the step executes. Defining this at 30 seconds, will result in a step that
	// runs every 30 seconds.
	//
	// The below are common examples to get started with.
	//
	//  // Every 30 seconds
	//  time.Duration(30 * time.Second)
	Interval time.Duration

	// MaxRetry is max retry times
	MaxRetry int
}

// Task contains the scheduled task details and control mechanisms. This struct is used during the creation of tasks.
// It allows users to control how and when tasks are executed.
type Task struct {
	RetryStep

	// DelayTime is used to specify a delay time for the scheduler. When set, tasks will wait for the specified
	// time to start the schedule timer. When not set, the previous task and the next task are executed concurrently.
	DelayTime time.Duration

	// TaskFunc is the user defined function to execute as part of this task.
	// Either TaskFunc or FuncWithTaskContext must be defined. If both are defined, FuncWithTaskContext will be used.
	TaskFunc func() error

	// AfterFunc is executed after TaskFunc is executed correctly. Execute only once
	AfterFunc func()

	// ErrFunc allows users to define a function that is called when tasks return an error. If ErrFunc is nil,
	// errors from tasks will be ignored.
	ErrFunc func(error)
	// ErrFunc is executed with an error. AfterFunc is executed without an error.

	// If it fails within the number of attempts of MaxRetry, whether to exit the entire Scheduler
	GlobalAbnormalEnd bool

	// ctx is the internal context used to control task cancelation.
	Ctx context.Context

	// timer is the internal task timer. This is stored here to provide control via main scheduler functions.
	timer *time.Timer

	// done is the internal finish signal. Keep in queue order
	done chan struct{}
}
