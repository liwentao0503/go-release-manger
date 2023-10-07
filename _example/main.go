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

	retryStep := releaseManage.RetryStep{
		Interval: 1 * time.Second,
		MaxRetry: 3,
	}
	tasks := &releaseManage.Task{
		Ctx:       mainCtx,
		RetryStep: retryStep,
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

	tasks2 := &releaseManage.Task{
		Ctx:               mainCtx,
		RetryStep:         retryStep,
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

	tasks3 := &releaseManage.Task{
		Ctx:       mainCtx,
		RetryStep: retryStep,
		DelayTime: 1 * time.Second,
		TaskFunc: func() error {
			fmt.Println("task3")
			return nil
		},
	}

	tasks4 := &releaseManage.Task{
		Ctx:       mainCtx,
		RetryStep: retryStep,
		DelayTime: 1 * time.Second,
		TaskFunc: func() error {
			fmt.Println("task4")
			return nil
		},
	}

	scheduler.Add(tasks, tasks2, tasks3, tasks4)

	go scheduler.StartTask(0)

	quit := make(chan os.Signal, 1)

	// `signal.Notify` registers the given channel to
	// receive notifications of the specified signals.
	signal.Notify(quit, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	mainCancel()
}
