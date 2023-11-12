package commands

import (
	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/flags"
	"os"
)

func Build(_ *flags.Build) error {
	return wrapTask("build")
}

func Dev(_ *flags.Dev) error {
	return wrapTask("dev")
}

func Package(_ *flags.Package) error {
	return wrapTask("package")
}

func wrapTask(command string) error {
	pterm.Warning.Printf("`wails3 %s` is an alias for `wails3 task %s`. Use `wails task` for better control and more options.\n", command, command)
	os.Args = []string{"wails3", "task", command}
	return RunTask(&RunTaskOptions{}, []string{})
}
