package releaseManage

import (
	"context"
	"fmt"
	"time"
)

type StepRetry struct {
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

type StepExecutionStatus uint8

const (
	// defauit 0 is running
	StepExecutionSuccess StepExecutionStatus = iota + 1
	StepExecutionSingleFailed
	StepExecutionGlobalFailed
)

func (s StepExecutionStatus) GetResult(err error) string {
	switch s {
	case StepExecutionSuccess:
		return ""
	case StepExecutionSingleFailed:
		return fmt.Sprintf("single step failed err: %v", err)
	case StepExecutionGlobalFailed:
		return fmt.Sprintf("Single-step failure leads to global failure err: %v", err)
	}
	return ""
}

type StepStatus struct {
	Status       StepExecutionStatus
	StartTime    time.Time
	EndTime      time.Time
	DurationTime time.Duration
	Result       string
}

// Step contains the scheduled step details and control mechanisms. This struct is used during the creation of steps.
// It allows users to control how and when steps are executed.
type Step struct {
	// Current running status
	StepStatus

	StepRetry

	// DelayTime is used to specify a delay time for the scheduler. When set, steps will wait for the specified
	// time to start the schedule timer. When not set, the previous step and the next step are executed concurrently.
	DelayTime time.Duration

	// stepFunc is the user defined function to execute as part of this step.
	// Either stepFunc or FuncWithstepContext must be defined. If both are defined, FuncWithstepContext will be used.
	StepFunc func() error

	// AfterFunc is executed after StepFunc is executed correctly. Execute only once
	AfterFunc func()

	// ErrFunc allows users to define a function that is called when steps return an error. If ErrFunc is nil,
	// errors from steps will be ignored.
	ErrFunc func(error)
	// ErrFunc is executed with an error. AfterFunc is executed without an error.

	// If it fails within the number of attempts of MaxRetry, whether to exit the entire Scheduler
	GlobalAbnormalEnd bool

	// ctx is the internal context used to control step cancelation.
	Ctx context.Context

	// timer is the internal step timer. This is stored here to provide control via main scheduler functions.
	timer *time.Timer

	// done is the internal finish signal. Keep in queue order
	done chan struct{}
}

func (s *Step) saveResult(err error) {
	if err != nil {
		if s.GlobalAbnormalEnd {
			s.Status = StepExecutionGlobalFailed
			s.Result = s.Status.GetResult(err)
			return
		}
		s.Status = StepExecutionSingleFailed
		s.Result = s.Status.GetResult(err)
	} else {
		s.Status = StepExecutionSuccess
	}

	s.done <- struct{}{}
}

func (s *Step) saveStepTime() {
	s.EndTime = time.Now()
	s.DurationTime = time.Since(s.StartTime)
}

func (s *Step) stepDone() {
	s.saveStepTime()
	if s.AfterFunc != nil {
		s.AfterFunc()
	}
	s.saveResult(nil)
}

func (s *Step) stepFailed(err error) {
	s.saveStepTime()
	s.timer.Stop()
	if s.ErrFunc != nil {
		s.ErrFunc(err)
	}
	s.saveResult(err)
}

func (s *Step) check() error {
	// Check if StepFunc is nil before doing anything
	if s.StepFunc == nil {
		return fmt.Errorf("step function cannot be nil")
	}

	// Ensure Interval is never 0, this would cause Timer to panic
	if s.Interval <= time.Duration(0) {
		return fmt.Errorf("step interval must be defined")
	}

	s.done = make(chan struct{})

	// Ensure MaxRetry is less 1
	if s.MaxRetry < 1 {
		s.MaxRetry = 1
	}

	return nil
}
