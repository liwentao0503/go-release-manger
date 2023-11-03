package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	releaseManage "release-manger"
)

func main() {
	// Create Context used to cancel downstream Goroutines
	mainCtx, mainCancel := context.WithCancel(context.Background())
	scheduler := releaseManage.New()

	retryStep := releaseManage.StepRetry{
		Interval: 1 * time.Second,
		MaxRetry: 3,
	}
	steps := &releaseManage.Step{
		Ctx:       mainCtx,
		StepRetry: retryStep,
		StepFunc: func() error {
			fmt.Println("step1")
			return nil
		},
		AfterFunc: func() {
			fmt.Println("step1 has finished ")
		},
		ErrFunc: func(error) {
			fmt.Println("wwsdasdas2")
		},
	}

	steps2 := &releaseManage.Step{
		Ctx:               mainCtx,
		StepRetry:         retryStep,
		DelayTime:         1 * time.Second,
		GlobalAbnormalEnd: true,
		StepFunc: func() error {
			fmt.Println("step2")
			return fmt.Errorf("step2 err")
		},
		ErrFunc: func(error) {
			fmt.Println("step2 failed use ErrFunc")
		},
		AfterFunc: func() {
			fmt.Println("step2 has finish")
		},
	}

	steps3 := &releaseManage.Step{
		Ctx:       mainCtx,
		StepRetry: retryStep,
		DelayTime: 1 * time.Second,
		StepFunc: func() error {
			fmt.Println("step3")
			return nil
		},
	}

	steps4 := &releaseManage.Step{
		Ctx:       mainCtx,
		StepRetry: retryStep,
		DelayTime: 1 * time.Second,
		StepFunc: func() error {
			fmt.Println("step4")
			return nil
		},
	}

	scheduler.Add(steps, steps2, steps3, steps4)

	go scheduler.ReleaseManage(mainCtx, func() {
		fmt.Printf("%s is doing\n", scheduler.Name)
	}, 0)

	quit := make(chan os.Signal, 1)

	// `signal.Notify` registers the given channel to
	// receive notifications of the specified signals.
	signal.Notify(quit, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	mainCancel()

	fmt.Println(scheduler.GetStepsExecutionStatus())
	time.Sleep(1 * time.Second)
}
