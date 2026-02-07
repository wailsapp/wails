package commands

import (
	"os"
	"runtime"
	"strings"

	"github.com/wailsapp/wails/v3/internal/flags"
)

// runTaskFunc is a variable to allow mocking in tests
var runTaskFunc = RunTask

// validPlatforms for GOOS
var validPlatforms = map[string]bool{
	"windows": true,
	"darwin":  true,
	"linux":   true,
}

func Build(buildFlags *flags.Build, otherArgs []string) error {
	if buildFlags.Tags != "" {
		otherArgs = append(otherArgs, "EXTRA_TAGS="+buildFlags.Tags)
	}
	return wrapTask("build", otherArgs)
}

func Package(_ *flags.Package, otherArgs []string) error {
	return wrapTask("package", otherArgs)
}

func SignWrapper(_ *flags.SignWrapper, otherArgs []string) error {
	return wrapTask("sign", otherArgs)
}

func wrapTask(action string, otherArgs []string) error {
	// Check environment first, then allow args to override
	goos := os.Getenv("GOOS")
	if goos == "" {
		goos = runtime.GOOS
	}
	goarch := os.Getenv("GOARCH")
	if goarch == "" {
		goarch = runtime.GOARCH
	}

	var remainingArgs []string

	// Args override environment
	for _, arg := range otherArgs {
		switch {
		case strings.HasPrefix(arg, "GOOS="):
			goos = strings.TrimPrefix(arg, "GOOS=")
		case strings.HasPrefix(arg, "GOARCH="):
			goarch = strings.TrimPrefix(arg, "GOARCH=")
		default:
			remainingArgs = append(remainingArgs, arg)
		}
	}

	// Determine task name based on GOOS
	taskName := action
	if validPlatforms[goos] {
		taskName = goos + ":" + action
	}

	// Pass ARCH to task (always set, defaults to current architecture)
	remainingArgs = append(remainingArgs, "ARCH="+goarch)

	// Rebuild os.Args to include the command and all additional arguments
	newArgs := []string{"wails3", "task", taskName}
	newArgs = append(newArgs, remainingArgs...)
	os.Args = newArgs
	// Pass the task name via options and remainingArgs as CLI variables
	return runTaskFunc(&RunTaskOptions{Name: taskName}, remainingArgs)
}
