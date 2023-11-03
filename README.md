# release-manger

[![Go Report Card](https://goreportcard.com/badge/github.com/wt1i/go-release-manger)](https://goreportcard.com/report/github.com/wt1i/go-release-manger) 
[![PkgGoDev](https://pkg.go.dev/badge/github.com/wt1i/go-release-manger)](https://pkg.go.dev/github.com/wt1i/go-release-manger)

Package releaseManage is an easy to use in-process scheduler for recurring steps in Go. Similar to a publishing system that publishes in order, 
it supports single-step retry, single-step delay, single-step completion hook, single-step error handling, and context graceful exit.

For simplicity this step scheduler uses the time.Duration type to specify intervals. This allows for a simple interface 
and flexible control over when steps are executed.

## Key Features

- **Sequential Execution**: For example, the publishing system needs to publish strictly in the order of dependencies.
- **Optimized Goroutine Scheduling**: Steps leverages Go's `time.AfterFunc()` function to reduce sleeping goroutines and optimize CPU scheduling.
- **Flexible Step Intervals**: Steps uses the `time.Duration` type to specify intervals, offering a simple interface and flexible control over step execution timing.
- **Delayed Step Start**: Schedule steps to start at a later time by specifying a delay time, allowing for greater control over step execution.
- **Max Retry Steps**: Schedule the step up to xx times by setting the `MaxRetry` flag, which is ideal for one-time steps or step execution failure timeout retries.
- **Custom Error Handling**: Define a custom error handling function to handle errors returned by steps, enabling tailored error handling logic.

## Flow Chart
![image](https://github.com/wt1i/go-release-manger/blob/main/img/flow_chart.png)