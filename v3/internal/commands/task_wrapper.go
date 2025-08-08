package commands

import (
	"github.com/wailsapp/wails/v3/internal/term"
	"os"

	"github.com/wailsapp/wails/v3/internal/flags"
)

func Build(options *flags.Build) error {
	taskName := "build"
	
	// If static flag is enabled, use the static build task
	if options.Static {
		taskName = "build:static"
	}

	// Set up environment variables for the task
	if options.Static {
		os.Setenv("STATIC", "true")
	}
	if options.Compiler != "" {
		os.Setenv("CC", options.Compiler)
	}

	return wrapTask(taskName)
}

func Package(options *flags.Package) error {
	taskName := "package"
	
	// If static flag is enabled, use the static package task
	if options.Static {
		taskName = "package:static"
	}

	// Set up environment variables for the task
	if options.Static {
		os.Setenv("STATIC", "true")
	}
	if options.Compiler != "" {
		os.Setenv("CC", options.Compiler)
	}

	return wrapTask(taskName)
}

func wrapTask(command string) error {
	term.Warningf("`wails3 %s` is an alias for `wails3 task %s`. Use `wails task` for better control and more options.\n", command, command)
	os.Args = []string{"wails3", "task", command}
	return RunTask(&RunTaskOptions{}, []string{})
}
