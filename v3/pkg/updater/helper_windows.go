//go:build windows

package updater

import "syscall"

// processQueryLimitedInformation grants only enough access to ask whether a
// process is still running — the minimum permissions required to call
// GetExitCodeProcess. Defined directly here rather than depending on
// golang.org/x/sys/windows so the updater stays on the standard library.
const processQueryLimitedInformation = 0x1000

// stillActive is the exit code Windows returns from GetExitCodeProcess while
// the process is still running. (Defined in MSDN as STILL_ACTIVE = 259.)
const stillActive = 259

// platformIsAlive reports whether pid names a running process. On Windows
// the previous os.Process.Signal(nil) probe always returned an error
// (syscall.EWINDOWS for live processes, ErrProcessDone for dead ones), so
// waitForPID short-circuited and the helper attempted to swap the binary
// while the parent still held it open — causing the swap-retry loop to grind
// against the lock instead of waiting.
//
// Open the process with PROCESS_QUERY_LIMITED_INFORMATION and ask the
// kernel for its exit code directly; STILL_ACTIVE distinguishes a live
// process from one that has already terminated.
func platformIsAlive(pid int) bool {
	h, err := syscall.OpenProcess(processQueryLimitedInformation, false, uint32(pid))
	if err != nil {
		return false
	}
	defer syscall.CloseHandle(h)
	var code uint32
	if err := syscall.GetExitCodeProcess(h, &code); err != nil {
		return false
	}
	return code == stillActive
}
