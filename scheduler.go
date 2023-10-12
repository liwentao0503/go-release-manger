package releaseManage

import (
	"context"
	"fmt"
	"os"
	"time"
)

// Scheduler stores the internal step list and provides an interface for step management.
type Scheduler struct {
	// HostName Scheduler's name
	Name   string
	ctx    context.Context
	cancel context.CancelFunc
	// steps is the internal step list used to store steps that are currently scheduled.
	steps []*Step
}

// New will create a new scheduler instance that allows users to create and manage steps.
func New() *Scheduler {
	s := &Scheduler{}
	s.steps = make([]*Step, 0)
	s.Name, _ = os.Hostname()
	s.ctx, s.cancel = context.WithCancel(context.Background())
	return s
}

func (schd *Scheduler) Add(steps ...*Step) error {
	for _, s := range steps {
		if err := s.check(); err != nil {
			return err
		}

		// Add step to schedule
		schd.steps = append(schd.steps, s)
	}

	return nil
}

// The returned step should be treated as read-only, and not modified outside of this package. Doing so, may cause
// panics.
func (schd *Scheduler) GetStepsExecutionStatus() []StepStatus {
	stepStatus := make([]StepStatus, 0, len(schd.steps))
	for _, v := range schd.steps {
		stepStatus = append(stepStatus, v.StepStatus)
	}
	return stepStatus
}

// scheduleStep creates the underlying scheduled step. If StartAfter is set, this routine will wait until the
// time specified.
func (schd *Scheduler) scheduleStep(s *Step) {
	time.Sleep(s.DelayTime)
	s.StartTime = time.Now()
	s.timer = time.AfterFunc(s.Interval, func() {
		select {
		case <-s.Ctx.Done():
			fmt.Println("main ctx has canceled")
			return
		default:
		}
		schd.execStep(s)
	})
}

// execStep is the underlying scheduler, it is used to trigger and execute steps.
func (schd *Scheduler) execStep(s *Step) {
	var err error
	if err = s.StepFunc(); err == nil {
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
	schd.cancel()
}

// ReleaseManage start steps in queue order
func (schd *Scheduler) ReleaseManage(ctx context.Context, reportBeat func(), start int) {
	go schd.reportBeat(ctx, reportBeat)
	for i := start; i < len(schd.steps); i++ {
		schd.scheduleStep(schd.steps[i])
		select {
		case <-schd.ctx.Done():
			return
		case <-schd.steps[i].done:
		}
	}
}

func (schd *Scheduler) reportBeat(ctx context.Context, doFunc func()) {
	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			doFunc()
		case <-ctx.Done():
			fmt.Println("reportBeat main ctx has canceled")
			return
		case <-schd.ctx.Done():
			return
		}
	}
}
