package commands

import (
	"context"

	"github.com/go-task/task/v3"
	"github.com/go-task/task/v3/taskfile"
)

type BuildOptions struct {
}

func Build(options *BuildOptions) error {
	e := task.Executor{}
	err := e.Setup()
	if err != nil {
		return err
	}
	build := taskfile.Call{
		Task: "build",
		Vars: nil,
	}
	return e.Run(context.Background(), build)
}
