//go:build !windows

package updater

import (
	"os/exec"
	"syscall"
)

// applyDetachAttrs puts the child in its own session so it isn't reaped or
// signalled along with the parent.
func applyDetachAttrs(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Setsid = true
}
