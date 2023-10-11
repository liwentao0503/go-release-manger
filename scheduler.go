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
		if err := s.check(); err != nil {
			return err
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

// ReleaseManage start tasks in queue order
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
