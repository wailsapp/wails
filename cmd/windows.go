// +build windows

package cmd

import (
	"os"

	"golang.org/x/sys/windows"
)

// Credit: https://stackoverflow.com/a/52579002

func init() {
	stdout := windows.Handle(os.Stdout.Fd())
	var originalMode uint32

	_ = windows.GetConsoleMode(stdout, &originalMode)
	_ = windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
}
