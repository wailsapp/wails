package cmd

import (
	"bytes"
	"os"
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
	cmd.Env = append(os.Environ(), "GO111MODULE=on")
	var stdo, stde bytes.Buffer
	cmd.Stdout = &stdo
	cmd.Stderr = &stde
	err = cmd.Run()
	stdout = string(stdo.Bytes())
	stderr = string(stde.Bytes())
	return
}

// RunInDirectory runs the given command in the given directory
func (sh *ShellHelper) RunInDirectory(dir string, command string, vars ...string) (stdout, stderr string, err error) {
	cmd := exec.Command(command, vars...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GO111MODULE=on")
	var stdo, stde bytes.Buffer
	cmd.Stdout = &stdo
	cmd.Stderr = &stde
	err = cmd.Run()
	stdout = string(stdo.Bytes())
	stderr = string(stde.Bytes())
	return
}
