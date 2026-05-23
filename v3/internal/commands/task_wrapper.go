package commands

import (
	"os"
	"runtime"
	"strings"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/wake"
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
	goos := os.Getenv("GOOS")
	if goos == "" {
		goos = runtime.GOOS
	}
	goarch := os.Getenv("GOARCH")
	if goarch == "" {
		goarch = runtime.GOARCH
	}

	var remainingArgs []string

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

	taskName := action
	if validPlatforms[goos] {
		taskName = goos + ":" + action
	}

	remainingArgs = append(remainingArgs, "ARCH="+goarch)

	if useWake() {
		return runWakeTask(taskName, goos, goarch, remainingArgs)
	}

	newArgs := []string{"wails3", "task", taskName}
	newArgs = append(newArgs, remainingArgs...)
	os.Args = newArgs
	return runTaskFunc(&RunTaskOptions{Name: taskName}, remainingArgs)
}

func useWake() bool {
	return os.Getenv("WAILS_USE_WAKE") == "true"
}

func runWakeTask(taskName, goos, goarch string, cliVars []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	vars := make(map[string]string)
	for _, v := range cliVars {
		if strings.Contains(v, "=") {
			parts := strings.SplitN(v, "=", 2)
			if len(parts) == 2 {
				vars[parts[0]] = parts[1]
			}
		}
	}

	opts := wake.ExecuteOptions{
		Dir:      dir,
		Platform: goos,
		Arch:     goarch,
		Vars:     vars,
		Verbose:  os.Getenv("WAKE_VERBOSE") != "",
		Silent:   os.Getenv("WAKE_SILENT") != "",
	}

	return wake.Execute(taskName, opts)
}
