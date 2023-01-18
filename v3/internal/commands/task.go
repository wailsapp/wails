package commands

import (
	"context"
	"fmt"

	"github.com/pterm/pterm"

	"github.com/go-task/task/v3"
	"github.com/go-task/task/v3/taskfile"
)

type RunTaskOptions struct {
	Name string `name:"n" description:"The name of the task to run"`
}

func RunTask(options *RunTaskOptions) error {
	if options.Name == "" {
		return fmt.Errorf("name of task required")
	}

	e := task.Executor{}
	err := e.Setup()
	if err != nil {
		return err
	}
	build := taskfile.Call{
		Task: options.Name,
		Vars: nil,
	}
	return e.Run(context.Background(), build)
}

type ListTaskOptions struct {
}

func ListTasks(options *ListTaskOptions) error {
	e := task.Executor{}
	if err := e.Setup(); err != nil {
		return err
	}
	tasks := e.GetTaskList()
	if len(tasks) == 0 {
		return fmt.Errorf("no tasks found. Ensure there is a `Taskfile.yml` in your project. You can generate a default takfile by running `wails generate defaults`")
	}
	tableData := [][]string{
		{"Task", "Summary"},
	}
	println()

	for _, thisTask := range tasks {
		if thisTask.Internal {
			continue
		}
		var thisRow = make([]string, 2)
		thisRow[0] = thisTask.Task
		thisRow[1] = thisTask.Summary
		tableData = append(tableData, thisRow)
	}
	err := pterm.DefaultTable.WithHasHeader(true).WithHeaderRowSeparator("-").WithData(tableData).Render()
	if err != nil {
		return err
	}
	println()
	return nil
}
