package commands

import (
	"context"
	"fmt"

	"github.com/go-task/task/v3"
	"github.com/go-task/task/v3/taskfile"
)

type RunOptions struct {
	Task string `name:"t" description:"The name of the task to run"`
}

func Run(options *RunOptions) error {
	if options.Task == "" {
		return fmt.Errorf("task name is required")
	}
	e := task.Executor{}
	err := e.Setup()
	if err != nil {
		return err
	}
	build := taskfile.Call{
		Task: options.Task,
		Vars: nil,
	}
	return e.Run(context.Background(), build)
}
