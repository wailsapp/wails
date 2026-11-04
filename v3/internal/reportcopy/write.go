//go:build !windows

// Package reportcopy writes diagnostic text to the system clipboard.
package reportcopy

import (
	"os/exec"
	"strings"
)

func writeCommand(text, name string, args ...string) error {
	command := exec.Command(name, args...)
	command.Stdin = strings.NewReader(text)
	return command.Run()
}
