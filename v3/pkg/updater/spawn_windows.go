//go:build windows

package updater

import (
	"errors"
	"fmt"
	"os/exec"
	"syscall"
)

// applyDetachAttrs marks the child as detached and asks Windows to keep it
// outside the parent's Job. File managers and launchers commonly use a
// kill-on-close Job; without CREATE_BREAKAWAY_FROM_JOB the updater helper is
// terminated as soon as the application exits.
func applyDetachAttrs(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	const (
		detachedProcess        = 0x00000008
		createNoWindow         = 0x08000000
		createNewProcGroup     = 0x00000200
		createBreakawayFromJob = 0x01000000
	)
	cmd.SysProcAttr.CreationFlags |= detachedProcess | createNoWindow | createNewProcGroup | createBreakawayFromJob
	cmd.SysProcAttr.HideWindow = true
}

func wrapHelperSpawnError(err error) error {
	if errors.Is(err, syscall.ERROR_ACCESS_DENIED) {
		return fmt.Errorf("%w: %v", ErrJobBreakawayDenied, err)
	}
	return fmt.Errorf("updater: spawn helper: %w", err)
}
