package releaseManage

import (
	"fmt"
	"time"
)

// Scheduler stores the internal task list and provides an interface for task management.
type Scheduler struct {
	abnormalEnd chan struct{}
	// tasks is the internal task list used to store tasks that are currently scheduled.
	steps []*Step
}

// New will create a new scheduler instance that allows users to create and manage tasks.
func New() *Scheduler {
	s := &Scheduler{}
	s.steps = make([]*Step, 0)
	return s
}

func (schd *Scheduler) Add(steps ...*Step) error {
	for _, s := range steps {
		// Check if TaskFunc is nil before doing anything
		if s.TaskFunc == nil {
			return fmt.Errorf("task function cannot be nil")
		}

		// Ensure Interval is never 0, this would cause Timer to panic
		if s.Interval <= time.Duration(0) {
			return fmt.Errorf("task interval must be defined")
		}

		if s.AfterFunc == nil {
			s.AfterFunc = func() {}
		}

		if s.ErrFunc == nil {
			s.ErrFunc = func(error) {}
		}

		s.done = make(chan struct{})

		// Ensure MaxRetry is less 1
		if s.MaxRetry < 1 {
			s.MaxRetry = 1
		}

		// Add task to schedule
		schd.steps = append(schd.steps, s)
	}

	return nil
}

// The returned task should be treated as read-only, and not modified outside of this package. Doing so, may cause
// panics.
func (schd *Scheduler) GetStepsExecutionStatus() []StepStatus {
	stepStatus := make([]StepStatus, 0, len(schd.steps))
	for _, v := range schd.steps {
		stepStatus = append(stepStatus, v.StepStatus)
	}
	return stepStatus
}

// scheduleTask creates the underlying scheduled task. If StartAfter is set, this routine will wait until the
// time specified.
func (schd *Scheduler) scheduleTask(s *Step) {
	time.Sleep(s.DelayTime)
	s.StartTime = time.Now()
	s.timer = time.AfterFunc(s.Interval, func() {
		select {
		case <-s.Ctx.Done():
			fmt.Println("main ctx has canceled")
			return
		default:
		}
		schd.execTask(s)
	})
}

// execTask is the underlying scheduler, it is used to trigger and execute tasks.
func (schd *Scheduler) execTask(s *Step) {
	var err error
	if err = s.TaskFunc(); err == nil {
		s.stepDone()
		return
	}

	s.MaxRetry--
	if s.MaxRetry == 0 {
		s.stepFailed(err)
		if s.GlobalAbnormalEnd {
			schd.StopReleaseManage()
		}
		return
	}

	s.timer.Reset(s.Interval)
}

// StopReleaseManage stop release manage by chan
func (schd *Scheduler) StopReleaseManage() {
	schd.abnormalEnd <- struct{}{}
}

// StartStep start tasks in queue order
func (schd *Scheduler) ReleaseManage(start int) {
	for i := start; i < len(schd.steps); i++ {
		schd.scheduleTask(schd.steps[i])
		select {
		case <-schd.abnormalEnd:
			return
		case <-schd.steps[i].done:
		}
	}
}
