package commands

import (
	"github.com/wailsapp/wails/v3/internal/flags"
)

func Build(options *flags.Build) error {
	return RunTask(&RunTaskOptions{
		Name: "build",
	})
}
