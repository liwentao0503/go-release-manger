package releaseManage

import (
	"fmt"
	"time"
)

// Scheduler stores the internal task list and provides an interface for task management.
type Scheduler struct {
	// tasks is the internal task list used to store tasks that are currently scheduled.
	tasks []*Task
}

// New will create a new scheduler instance that allows users to create and manage tasks.
func New() *Scheduler {
	s := &Scheduler{}
	s.tasks = make([]*Task, 0)
	return s
}

// Add will add a task to the task list and schedule it. Once added, tasks will wait the defined time interval and then
// execute. This means a task with a 15 second interval will be triggered 15 seconds after Add is complete. Not before
// or after (excluding typical machine time jitter).
//
//	// Add a task
//	id, err := scheduler.Add(&tasks.Task{
//		Interval: time.Duration(30 * time.Second),
//		TaskFunc: func() error {
//			// Put your logic here
//		}(),
//		ErrFunc: func(err error) {
//			// Put custom error handling here
//		}(),
//	})
//	if err != nil {
//		// Do stuff
//	}
func (schd *Scheduler) Add(tasks ...*Task) error {
	for _, t := range tasks {
		// Check if TaskFunc is nil before doing anything
		if t.TaskFunc == nil {
			return fmt.Errorf("task function cannot be nil")
		}

		// Ensure Interval is never 0, this would cause Timer to panic
		if t.Interval <= time.Duration(0) {
			return fmt.Errorf("task interval must be defined")
		}

		if t.AfterFunc == nil {
			t.AfterFunc = func() {}
		}

		t.done = make(chan struct{})

		// Ensure MaxRetry is less 1
		if t.MaxRetry < 1 {
			t.MaxRetry = 1
		}

		// Add task to schedule
		schd.tasks = append(schd.tasks, t)
	}

	return nil
}

// The returned task should be treated as read-only, and not modified outside of this package. Doing so, may cause
// panics.
func (schd *Scheduler) Tasks() []*Task {
	return schd.tasks
}

// scheduleTask creates the underlying scheduled task. If StartAfter is set, this routine will wait until the
// time specified.
func (schd *Scheduler) scheduleTask(t *Task) {
	time.Sleep(t.DelayTime)
	t.timer = time.AfterFunc(t.Interval, func() {
		select {
		case <-t.Ctx.Done():
			return
		default:
		}
		schd.execTask(t)
	})
}

// execTask is the underlying scheduler, it is used to trigger and execute tasks.
func (schd *Scheduler) execTask(t *Task) {
	var err error

	if err = t.TaskFunc(); err != nil {
		t.MaxRetry--
	} else {
		t.AfterFunc()
		// if success return
		t.done <- struct{}{}
		return
	}

	if t.MaxRetry == 0 || t.RunOnce {
		t.timer.Stop()
		if t.ErrFunc != nil && err != nil {
			t.ErrFunc(err)
		}
		t.done <- struct{}{}
		return
	}

	t.timer.Reset(t.Interval)
}

// StartTask start tasks in queue order
func (schd *Scheduler) StartTask() {
	for _, v := range schd.tasks {
		schd.scheduleTask(v)
		<-v.done
	}
}
