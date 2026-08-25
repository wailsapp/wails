//go:build windows

package updater

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

// processQueryLimitedInformation grants only enough access to ask whether a
// process is still running — the minimum permissions required to call
// GetExitCodeProcess. Defined directly here rather than depending on
// golang.org/x/sys/windows so the updater stays on the standard library.
const processQueryLimitedInformation = 0x1000

// stillActive is the exit code Windows returns from GetExitCodeProcess while
// the process is still running. (Defined in MSDN as STILL_ACTIVE = 259.)
const stillActive = 259

// errNotSameDevice is the Windows error code returned by os.Rename when src
// and dst reside on different volumes (ERROR_NOT_SAME_DEVICE = 17 / 0x11).
// Defined locally to avoid a dependency on golang.org/x/sys/windows.
const errNotSameDevice = syscall.Errno(17)

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

// renameOrCopy attempts to move a file via os.Rename. If it fails specifically
// because src and dst are on different volumes (ERROR_NOT_SAME_DEVICE / errno 17),
// it transparently falls back to a secure copy-and-delete strategy.
//
// All other os.Rename errors (permission denied, target locked, etc.) are
// returned as-is so the caller sees the real failure reason.
func renameOrCopy(src, dst string) error {
	err := os.Rename(src, dst)
	if err == nil {
		return nil
	}

	var linkErr *os.LinkError
	if !errors.As(err, &linkErr) || !errors.Is(linkErr.Err, errNotSameDevice) {
		return err
	}

	if err := copyFileExec(src, dst); err != nil {
		return fmt.Errorf("cross-volume copy %s -> %s: %w", src, dst, err)
	}

	_ = os.Remove(src)
	return nil
}

// copyFileExec duplicates the executable from src to dst. It preserves the
// original file mode, flushes data to the storage device, and ensures that
// any partially written destination file is deleted if an error occurs.
func copyFileExec(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	info, err := in.Stat()
	if err != nil {
		return err
	}

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}

	// Deferred cleanup to handle failures gracefully. If an error occurs during
	// copying or syncing, the partial file at dst will be removed.
	defer func() {
		_ = out.Close()
		if err != nil {
			_ = os.Remove(dst)
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	// Ensure all dirty pages are flushed to disk before heading back.
	return out.Sync()
}

// replaceTarget puts the file at newPath into target's slot. On Windows the
// kernel keeps an executable's image file held for some time after the
// process that ran it exits — long enough that os.Remove(target) fails with
// "Access is denied" through multiple retry attempts. Discovered against
// wailsapp/updater-demo on Windows 11 amd64 (helper logged 20 consecutive
// unlinkat failures, ~10s total).
//
// Windows does, however, allow renaming a file whose image is still mapped.
// So we rename the target aside (giving the new file a free slot at target)
// and let the stale .old file be cleaned up later — best-effort delete here,
// and a sweep in maybeCleanReplacedAsides on the next helper run takes care
// of any leftovers once the kernel has finally released them.
func replaceTarget(target, newPath string) error {
	// Best-effort first try: if nothing's actually holding it, a normal
	// remove + rename is cleaner (no .old file left behind).
	if err := os.RemoveAll(target); err == nil {
		// Sweep any .old files leftover from prior updates — by now the
		// kernel has released them.
		sweepRenameAsides(target)
		return renameOrCopy(newPath, target)
	}

	aside := fmt.Sprintf("%s.old.%d", target, time.Now().UnixNano())
	if err := os.Rename(target, aside); err != nil {
		return fmt.Errorf("rename-aside %s -> %s: %w", target, aside, err)
	}

	if err := renameOrCopy(newPath, target); err != nil {
		// Revert the original file back to its target slot to prevent half-states.
		_ = os.Rename(aside, target)
		return err
	}

	// The just-created aside is probably still mapped by the kernel; this
	// remove will fail. Sweep grabs older asides whose owning processes are
	// long gone.
	_ = os.Remove(aside)
	sweepRenameAsides(target)
	return nil
}

// sweepRenameAsides best-effort-deletes any "<target>.old.*" siblings.
// Without this, a Windows app updated N times accumulates N stale
// executables in its install directory.
func sweepRenameAsides(target string) {
	matches, err := filepath.Glob(target + ".old.*")
	if err != nil {
		return
	}
	for _, m := range matches {
		_ = os.Remove(m)
	}
}
