package commands

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// syncDirs makes dst's contents identical to src's, then removes src.
//
// Unlike a RemoveAll+Rename swap, dst itself is never deleted or renamed:
// on Windows a file watcher (e.g. Vite's chokidar) holds an open handle on
// watched directories, which makes deleting or renaming over them fail with
// "Access is denied" (#5515). Updating individual files also avoids the
// chokidar rename-event loop caused by deleting and recreating the watched
// directory (#3976), and leaves unchanged files untouched so the dev server
// only reloads what actually changed.
func syncDirs(src, dst string) error {
	info, err := os.Lstat(dst)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		// Fast path: no destination, move the whole tree at once.
		if renameErr := withRetry(func() error { return os.Rename(src, dst) }); renameErr == nil {
			return nil
		}
		// Fall through to the per-file sync, which recreates dst below.
	case err != nil:
		return err
	case !info.IsDir():
		// The destination exists but is not a directory: clear it, then
		// move the tree into place.
		if err := withRetry(func() error { return os.RemoveAll(dst) }); err != nil {
			return err
		}
		if renameErr := withRetry(func() error { return os.Rename(src, dst) }); renameErr == nil {
			return nil
		}
	}

	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}

	// Copy phase: bring every entry of src into dst, skipping identical files.
	keep := make(map[string]bool)
	err = filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		keep[rel] = true
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			if info, err := os.Lstat(target); err == nil && !info.IsDir() {
				if err := withRetry(func() error { return os.Remove(target) }); err != nil {
					return err
				}
			}
			return os.MkdirAll(target, 0o755)
		}
		if sameFileContent(path, target) {
			return nil
		}
		return replaceFile(path, target)
	})
	if err != nil {
		return err
	}

	// Delete phase: drop anything in dst that the generator did not produce.
	var stale []string
	err = filepath.WalkDir(dst, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(dst, path)
		if err != nil {
			return err
		}
		if rel == "." || keep[rel] {
			return nil
		}
		stale = append(stale, path)
		if d.IsDir() {
			return fs.SkipDir
		}
		return nil
	})
	if err != nil {
		return err
	}
	for _, path := range stale {
		if err := withRetry(func() error { return os.RemoveAll(path) }); err != nil {
			return err
		}
	}

	return withRetry(func() error { return os.RemoveAll(src) })
}

// replaceFile moves src over dst, clearing dst first if it is not a regular
// file (e.g. a directory now replaced by a file of the same name).
func replaceFile(src, dst string) error {
	if info, err := os.Lstat(dst); err == nil && !info.Mode().IsRegular() {
		if err := withRetry(func() error { return os.RemoveAll(dst) }); err != nil {
			return err
		}
	}
	err := withRetry(func() error { return os.Rename(src, dst) })
	if err == nil {
		return nil
	}
	// os.Rename replaces existing files on Windows (MOVEFILE_REPLACE_EXISTING),
	// but not when dst has the read-only attribute set. Clear dst and retry.
	if removeErr := withRetry(func() error { return os.Remove(dst) }); removeErr == nil {
		err = withRetry(func() error { return os.Rename(src, dst) })
	}
	return err
}

// sameFileContent reports whether a and b are regular files with equal
// content. Generated bindings are small, so a full read is fine.
func sameFileContent(a, b string) bool {
	infoA, errA := os.Stat(a)
	infoB, errB := os.Stat(b)
	if errA != nil || errB != nil || !infoA.Mode().IsRegular() || !infoB.Mode().IsRegular() || infoA.Size() != infoB.Size() {
		return false
	}
	dataA, err := os.ReadFile(a)
	if err != nil {
		return false
	}
	dataB, err := os.ReadFile(b)
	if err != nil {
		return false
	}
	return bytes.Equal(dataA, dataB)
}

// withRetry retries op with backoff on Windows, where file watchers and
// antivirus scanners briefly lock files and directories. Other platforms get
// a single attempt.
func withRetry(op func() error) error {
	var err error
	delay := 8 * time.Millisecond
	for attempt := 0; ; attempt++ {
		if err = op(); err == nil || runtime.GOOS != "windows" || attempt == 5 {
			return err
		}
		time.Sleep(delay)
		delay *= 2
	}
}
