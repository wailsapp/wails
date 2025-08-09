package commands

import (
	"github.com/wailsapp/wails/v3/internal/term"
	"os"

	"github.com/wailsapp/wails/v3/internal/flags"
)

// runTaskFunc is a variable to allow mocking in tests
var runTaskFunc = RunTask

func Build(_ *flags.Build, otherArgs []string) error {
	return wrapTask("build", otherArgs)
}

func Package(_ *flags.Package, otherArgs []string) error {
	return wrapTask("package", otherArgs)
}

func wrapTask(command string, otherArgs []string) error {
	term.Warningf("`wails3 %s` is an alias for `wails3 task %s`. Use `wails task` for better control and more options.\n", command, command)
	// Rebuild os.Args to include the command and all additional arguments
	newArgs := []string{"wails3", "task", command}
	newArgs = append(newArgs, otherArgs...)
	os.Args = newArgs
	// Pass the task name via options and otherArgs as CLI variables
	return runTaskFunc(&RunTaskOptions{Name: command}, otherArgs)
}
