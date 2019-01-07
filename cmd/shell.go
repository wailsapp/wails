package cmd

import (
	"bytes"
	"os/exec"
)

// ShellHelper helps with Shell commands
type ShellHelper struct {
}

// NewShellHelper creates a new ShellHelper!
func NewShellHelper() *ShellHelper {
	return &ShellHelper{}
}

// Run the given command
func (sh *ShellHelper) Run(command string, vars ...string) (stdout, stderr string, err error) {
	cmd := exec.Command(command, vars...)
	var stdo, stde bytes.Buffer
	cmd.Stdout = &stdo
	cmd.Stderr = &stde
	err = cmd.Run()
	stdout = string(stdo.Bytes())
	stderr = string(stde.Bytes())
	return
}
