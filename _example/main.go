package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	releaseManage "github.com/liwentao0503/go-release-manger"
)

func main() {
	// Create Context used to cancel downstream Goroutines
	mainCtx, mainCancel := context.WithCancel(context.Background())
	scheduler := releaseManage.New()

	retryStep := releaseManage.StepRetry{
		Interval: 1 * time.Second,
		MaxRetry: 3,
	}
	tasks := &releaseManage.Step{
		Ctx:       mainCtx,
		StepRetry: retryStep,
		TaskFunc: func() error {
			fmt.Println("task1")
			return nil
		},
		AfterFunc: func() {
			fmt.Println("task1 has finished ")
		},
		ErrFunc: func(error) {
			fmt.Println("wwsdasdas2")
		},
	}

	tasks2 := &releaseManage.Step{
		Ctx:               mainCtx,
		StepRetry:         retryStep,
		DelayTime:         1 * time.Second,
		GlobalAbnormalEnd: true,
		TaskFunc: func() error {
			fmt.Println("task2")
			return fmt.Errorf("task2 err")
		},
		ErrFunc: func(error) {
			fmt.Println("task2 failed use ErrFunc")
		},
		AfterFunc: func() {
			fmt.Println("task2 has finish")
		},
	}

	tasks3 := &releaseManage.Step{
		Ctx:       mainCtx,
		StepRetry: retryStep,
		DelayTime: 1 * time.Second,
		TaskFunc: func() error {
			fmt.Println("task3")
			return nil
		},
	}

	tasks4 := &releaseManage.Step{
		Ctx:       mainCtx,
		StepRetry: retryStep,
		DelayTime: 1 * time.Second,
		TaskFunc: func() error {
			fmt.Println("task4")
			return nil
		},
	}

	scheduler.Add(tasks, tasks2, tasks3, tasks4)

	go scheduler.StartStep(0)

	quit := make(chan os.Signal, 1)

	// `signal.Notify` registers the given channel to
	// receive notifications of the specified signals.
	signal.Notify(quit, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	mainCancel()

	fmt.Println(scheduler.GetStepsExecutionStatus())
	time.Sleep(1 * time.Second)
}
