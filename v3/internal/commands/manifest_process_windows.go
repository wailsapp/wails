//go:build windows

package commands

import (
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

func configureManifestProcess(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}
}

func signalManifestProcess(process *os.Process, _ os.Signal) error {
	return exec.Command("taskkill", "/T", "/PID", strconv.Itoa(process.Pid)).Run()
}

func killManifestProcess(process *os.Process) error {
	return exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(process.Pid)).Run()
}
