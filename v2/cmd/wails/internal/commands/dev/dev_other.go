//go:build darwin || linux
// +build darwin linux

package dev

import (
	"os/exec"
	"syscall"

	"github.com/wailsapp/wails/v2/cmd/wails/internal/logutils"
	"golang.org/x/sys/unix"
)

func setParentGID(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
}

func killProc(cmd *exec.Cmd, devCommand string) {
	if cmd == nil || cmd.Process == nil {
		return
	}

	// Experiencing the same issue on macOS BigSur
	// I'm using Vite, but I would imagine this could be an issue with Node (npm) in general
	// Also, after several edit/rebuild cycles any abnormal shutdown (crash or CTRL-C) may still leave Node running
	// Credit: https://stackoverflow.com/a/29552044/14764450 (same page as the Windows solution above)
	// Not tested on *nix
	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err == nil {
		err := syscall.Kill(-pgid, unix.SIGTERM) // note the minus sign
		if err != nil {
			logutils.LogRed("Error from '%s' when attempting to kill the process: %s", devCommand, err.Error())
		}
	}
}
