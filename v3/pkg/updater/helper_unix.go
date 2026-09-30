//go:build !windows

package updater

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
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

// replaceTarget moves the existing file or directory at target aside,
// renames newPath into its place, and then deletes the aside. On Unix a
// running binary can be renamed and unlinked because open handles stay valid
// against the inode. Keeping the aside until the new payload is in place
// means a failed or slow cross-device copy never leaves target missing.
// macOS .app bundles are directories, hence RemoveAll.
func replaceTarget(target, newPath string) error {
	aside := target + ".wails-old"
	if err := os.RemoveAll(aside); err != nil {
		return err
	}
	if err := rename(target, aside); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	if err := renameOrCopy(newPath, target); err != nil {
		_ = rename(aside, target)
		return err
	}
	_ = os.RemoveAll(aside)
	return nil
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
	// copyTree merges into an existing directory, so a stale intermediate
	// must be gone before copying or its extra files would ship.
	if err := os.RemoveAll(tmp); err != nil {
		return fmt.Errorf("clear %s: %w", tmp, err)
	}
	if err := copyAny(src, tmp); err != nil {
		_ = os.RemoveAll(tmp)
		return fmt.Errorf("cross-device copy %s -> %s: %w", src, tmp, err)
	}
	if err := rename(tmp, dst); err != nil {
		_ = os.RemoveAll(tmp)
		return err
	}
	// Make the rename durable before deleting the only other copy.
	if err := syncDir(filepath.Dir(dst)); err != nil {
		return fmt.Errorf("sync %s: %w", filepath.Dir(dst), err)
	}
	_ = os.RemoveAll(src)
	return nil
}

func syncDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	// Some filesystems cannot fsync a directory; there is nothing more to do there.
	if err := d.Sync(); err != nil && !errors.Is(err, syscall.EINVAL) && !errors.Is(err, syscall.ENOTSUP) {
		_ = d.Close()
		return err
	}
	return d.Close()
}
