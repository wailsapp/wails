//go:build windows

package updater

import (
	"errors"
	"fmt"
	"os/exec"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
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

// startHelper preserves normal child creation for jobs that neither permit
// breakaway nor kill their processes on close. Querying a null job handle
// reports the immediate job; outer jobs can still impose their own limits.
// If the query fails (including when there is no job), keep requesting breakaway.
func startHelper(cmd *exec.Cmd) error {
	var info windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION
	err := windows.QueryInformationJobObject(0, windows.JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(&info)), uint32(unsafe.Sizeof(info)), nil)
	const needsBreakaway = windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE |
		windows.JOB_OBJECT_LIMIT_BREAKAWAY_OK | windows.JOB_OBJECT_LIMIT_SILENT_BREAKAWAY_OK
	if err == nil && info.BasicLimitInformation.LimitFlags&needsBreakaway == 0 && cmd.SysProcAttr != nil {
		cmd.SysProcAttr.CreationFlags &^= windows.CREATE_BREAKAWAY_FROM_JOB
	}
	return cmd.Start()
}
