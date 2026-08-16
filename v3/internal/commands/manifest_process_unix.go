//go:build !windows

package commands

import (
	"os"
	"os/exec"
	"syscall"
)

func configureManifestProcess(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func signalManifestProcess(process *os.Process, signal os.Signal) error {
	if value, ok := signal.(syscall.Signal); ok {
		return syscall.Kill(-process.Pid, value)
	}
	return process.Signal(signal)
}

func killManifestProcess(process *os.Process) error {
	return syscall.Kill(-process.Pid, syscall.SIGKILL)
}
