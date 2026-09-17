//go:build !windows

package updater

import (
	"errors"
	"fmt"
	"os"
	"syscall"
)

// errCrossDevice is the errno returned by os.Rename when src and dst reside
// on different filesystems (EXDEV). On Linux this is errno 18; the canonical
// layout that hits it is tmpfs /tmp staging versus a persistent $HOME
// install dir. Defined locally so the updater stays on the standard library.
const errCrossDevice = syscall.EXDEV

// renameFunc performs the underlying rename. Held in a package-level var so
// tests can inject a synthetic EXDEV without needing two real filesystems
// (mirrors selfExecutable / newDetachedCommand in spawn.go).
var renameFunc = os.Rename

// platformIsAlive reports whether pid names a running process. Sending the
// no-op signal 0 to a pid either succeeds (process is running and we have
// permission) or fails (process is gone / no permission).
func platformIsAlive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return proc.Signal(syscall.Signal(0)) == nil
}

// replaceTarget puts the file or directory at newPath into target's slot.
// On Unix open file handles remain valid against the unlinked inode, so a
// running binary can be replaced in place. rename(2) atomically replaces a
// file destination, hence no RemoveAll is needed on the common file path —
// the original stays in place unless the move itself succeeds. macOS .app
// bundles are directories: those cannot be renamed onto a non-empty
// directory nor safely merged via copy, so the slot is cleared first and the
// move retried once for a clean replace.
func replaceTarget(target, newPath string) error {
	// Directory-over-directory needs an empty slot so the retry lands clean
	// instead of merging stale files from the old bundle into the new one.
	if bothDirs(target, newPath) {
		if err := os.RemoveAll(target); err != nil {
			return err
		}
		return renameOrCopy(newPath, target)
	}
	return renameOrCopy(newPath, target)
}

// renameOrCopy attempts to move src to dst via rename. If it fails
// specifically because src and dst are on different filesystems (EXDEV), it
// transparently falls back to a copy-and-delete strategy reusing copyAny so
// both files and .app directories are handled with modes and fsync intact.
//
// All other rename errors (permission denied, missing source, etc.) are
// returned as-is so the caller sees the real failure reason and —
// critically — the original at dst is left untouched for the outer
// backup/restore in runHelperSwap instead of being deleted first.
func renameOrCopy(src, dst string) error {
	err := renameFunc(src, dst)
	if err == nil {
		return nil
	}
	if !isCrossDevice(err) {
		return err
	}

	if err := copyAny(src, dst); err != nil {
		return fmt.Errorf("cross-device copy %s -> %s: %w", src, dst, err)
	}

	_ = os.RemoveAll(src)
	return nil
}

// isCrossDevice reports whether err is the EXDEV failure rename returns for
// a cross-filesystem move. Only this narrow case triggers the copy fallback.
func isCrossDevice(err error) bool {
	var linkErr *os.LinkError
	return errors.As(err, &linkErr) && errors.Is(linkErr.Err, errCrossDevice)
}

// bothDirs reports whether target and newPath both exist and are
// directories. Used to select the clear-slot-then-retry path for .app
// bundle replaces.
func bothDirs(target, newPath string) bool {
	ti, err := os.Stat(target)
	if err != nil || !ti.IsDir() {
		return false
	}
	ni, err := os.Stat(newPath)
	if err != nil || !ni.IsDir() {
		return false
	}
	return true
}
