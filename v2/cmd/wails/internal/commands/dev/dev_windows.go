//go:build windows
// +build windows

package dev

import (
	"os"
	"os/exec"
	"strconv"
)

func setParentGID(_ *exec.Cmd) {}

func killProc(cmd *exec.Cmd, devCommand string) {
	// Credit: https://stackoverflow.com/a/44551450
	// For whatever reason, killing an npm script on windows just doesn't exit properly with cancel
	if cmd != nil && cmd.Process != nil {
		kill := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(cmd.Process.Pid))
		kill.Stderr = os.Stderr
		kill.Stdout = os.Stdout
		err := kill.Run()
		if err != nil {
			if err.Error() != "exit status 1" {
				LogRed("Error from '%s': %s", devCommand, err.Error())
			}
		}
	}
}
