package commands

import (
	"os"

	"github.com/pterm/pterm"

	"github.com/wailsapp/wails/v3/internal/flags"
)

func Build(_ *flags.Build) error {
	pterm.Info.Println("`wails build` is an alias for `wails task build`. Use `wails task` for much better control over your builds.")
	os.Args = []string{"wails", "task", "build"}
	return RunTask(&RunTaskOptions{}, []string{})
}
