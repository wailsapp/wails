//go:build !windows

package updater

import (
	"errors"
	"fmt"
	"os"
	"syscall"
)

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

// replaceTarget removes the existing file or directory at target and renames
// newPath into its place. On Unix this is straightforward: open file handles
// remain valid against the unlinked inode, so we can delete a running binary
// and immediately put a new one at its path. macOS .app bundles are
// directories, hence RemoveAll rather than Remove.
func replaceTarget(target, newPath string) error {
	if err := os.RemoveAll(target); err != nil {
		return err
	}
	return renameOrCopy(newPath, target)
}

// rename is os.Rename, swappable so tests can simulate EXDEV.
var rename = os.Rename

// renameOrCopy moves src to dst. When they sit on different filesystems it
// copies to a sibling of dst, fsyncs, and renames that into place, so dst
// never holds a partial copy.
func renameOrCopy(src, dst string) error {
	err := rename(src, dst)
	if !errors.Is(err, syscall.EXDEV) {
		return err
	}
	tmp := dst + ".wails-new"
	_ = os.RemoveAll(tmp)
	if err := copyAny(src, tmp); err != nil {
		_ = os.RemoveAll(tmp)
		return fmt.Errorf("cross-device copy %s -> %s: %w", src, tmp, err)
	}
	if err := rename(tmp, dst); err != nil {
		_ = os.RemoveAll(tmp)
		return err
	}
	_ = os.RemoveAll(src)
	return nil
}
