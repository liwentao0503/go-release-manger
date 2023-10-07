# go-release-manger

[![Go Report Card](https://goreportcard.com/badge/github.com/liwentao0503/go-release-manger)](https://goreportcard.com/report/github.com/liwentao0503/go-release-manger) 
[![PkgGoDev](https://pkg.go.dev/badge/github.com/liwentao0503/go-release-manger)](https://pkg.go.dev/github.com/liwentao0503/go-release-manger)

Package tasks is an easy to use in-process scheduler for recurring tasks in Go. Similar to a publishing system that publishes in order, 
it supports single-step retry, single-step delay, single-step completion hook, single-step error handling, and context graceful exit.

Tasks focus on the sequence of task execution

For simplicity this task scheduler uses the time.Duration type to specify intervals. This allows for a simple interface 
and flexible control over when tasks are executed.

## Key Features

- **Sequential Execution**: For example, the publishing system needs to publish strictly in the order of dependencies.
- **Optimized Goroutine Scheduling**: Tasks leverages Go's `time.AfterFunc()` function to reduce sleeping goroutines and optimize CPU scheduling.
- **Flexible Task Intervals**: Tasks uses the `time.Duration` type to specify intervals, offering a simple interface and flexible control over task execution timing.
- **Delayed Task Start**: Schedule tasks to start at a later time by specifying a delay time, allowing for greater control over task execution.
- **Max Retry Tasks**: Schedule the task up to xx times by setting the `MaxRetry` flag, which is ideal for one-time tasks or task execution failure timeout retries.
- **Custom Error Handling**: Define a custom error handling function to handle errors returned by tasks, enabling tailored error handling logic.

