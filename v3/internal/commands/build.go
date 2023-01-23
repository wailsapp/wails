package commands

import (
	"os"

	"github.com/wailsapp/wails/v3/internal/flags"
)

func Build(_ *flags.Build) error {
	os.Args = []string{"wails", "task", "build"}
	return RunTask(&RunTaskOptions{}, []string{})
}
