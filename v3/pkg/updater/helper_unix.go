//go:build !windows

package updater

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
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
// transparently falls back to a copy-and-delete strategy.
//
// Files are staged through a temporary sibling beside dst that is synced
// and then atomically renamed into place, so a crash mid-copy can never
// leave a partial binary at dst: either the complete old file or the
// complete new file is there. Directories (.app bundles) cannot be
// atomically replaced on Unix, so those still copy straight onto the
// cleared slot; every failure mode where the helper process survives is
// covered by the outer backup/restore in runHelperSwap.
//
// All other rename errors (permission denied, missing source, etc.) are
// returned as-is so the caller sees the real failure reason and —
// critically — the original at dst is left untouched for the outer
// backup/restore instead of being deleted first.
func renameOrCopy(src, dst string) error {
	err := renameFunc(src, dst)
	if err == nil {
		return nil
	}
	if !isCrossDevice(err) {
		return err
	}

	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("cross-device copy %s -> %s: %w", src, dst, err)
	}
	if info.IsDir() {
		if err := copyAny(src, dst); err != nil {
			return fmt.Errorf("cross-device copy %s -> %s: %w", src, dst, err)
		}
	} else if err := stageFileCopy(src, dst, info.Mode()); err != nil {
		return err
	}

	_ = os.RemoveAll(src)
	return nil
}

// stageFileCopy copies the file at src to dst via a temporary sibling in
// dst's directory. The staging file carries src's mode, is synced to stable
// storage, and is then atomically renamed over dst — dst holds either the
// complete old or the complete new file, never a partial write. Any staging
// file left behind by a failure is removed before returning.
func stageFileCopy(src, dst string, mode os.FileMode) error {
	tmp, err := os.CreateTemp(filepath.Dir(dst), ".wails-update-*")
	if err != nil {
		return fmt.Errorf("cross-device copy %s -> %s: %w", src, dst, err)
	}
	tmpName := tmp.Name()
	_ = tmp.Close()
	cleanup := func() { _ = os.Remove(tmpName) }

	// copyFile truncates the already-created staging file without
	// re-applying the mode, so set it explicitly before the swap.
	if err := copyFile(src, tmpName, mode); err != nil {
		cleanup()
		return fmt.Errorf("cross-device copy %s -> %s: %w", src, dst, err)
	}
	if err := os.Chmod(tmpName, mode); err != nil {
		cleanup()
		return fmt.Errorf("cross-device copy %s -> %s: %w", src, dst, err)
	}
	// Same-directory swap: provably the same filesystem, so this cannot
	// EXDEV. Plain os.Rename (not renameFunc) is deliberate — renameFunc
	// exists only to simulate the cross-device src→dst move under test.
	if err := os.Rename(tmpName, dst); err != nil {
		cleanup()
		return fmt.Errorf("cross-device copy %s -> %s: %w", src, dst, err)
	}
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
