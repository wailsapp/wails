//go:build !windows

package updater

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

// crossDeviceErr builds the synthetic EXDEV failure os.Rename returns for a
// cross-filesystem move. Tests inject it via renameFunc so the fallback is
// exercised deterministically without needing two real filesystems.
func crossDeviceErr(src, dst string) error {
	return &os.LinkError{Op: "rename", Old: src, New: dst, Err: errCrossDevice}
}

// withRenameFunc swaps renameFunc for the duration of a test.
func withRenameFunc(t *testing.T, fn func(oldpath, newpath string) error) {
	t.Helper()
	prev := renameFunc
	renameFunc = fn
	t.Cleanup(func() { renameFunc = prev })
}

// TestIsCrossDevice verifies EXDEV classification: only cross-device link
// errors trigger the copy fallback; nil, plain, and other errno failures do
// not.
func TestIsCrossDevice(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil error", nil, false},
		{"exdev link error", crossDeviceErr("a", "b"), true},
		{"non-exdev link error", &os.LinkError{Op: "rename", Old: "a", New: "b", Err: syscall.EACCES}, false},
		{"plain error", errors.New("boom"), false},
		{"wrapped exdev", &os.LinkError{Op: "rename", Old: "a", New: "b", Err: syscall.EXDEV}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isCrossDevice(tt.err); got != tt.want {
				t.Errorf("isCrossDevice(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

// TestBothDirs verifies the directory-pair detection that selects the
// clear-slot-then-retry replace path for .app bundle swaps.
func TestBothDirs(t *testing.T) {
	dir := t.TempDir()
	fileA := filepath.Join(dir, "a.bin")
	fileB := filepath.Join(dir, "b.bin")
	writeFile(t, fileA, []byte("A"))
	writeFile(t, fileB, []byte("B"))
	dirA := filepath.Join(dir, "a.app")
	dirB := filepath.Join(dir, "b.app")
	if err := os.MkdirAll(dirA, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dirB, 0o755); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		target string
		new    string
		want   bool
	}{
		{"both dirs", dirA, dirB, true},
		{"target file new dir", fileA, dirA, false},
		{"target dir new file", dirA, fileA, false},
		{"both files", fileA, fileB, false},
		{"missing target", filepath.Join(dir, "nope"), dirA, false},
		{"missing new", dirA, filepath.Join(dir, "nope"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bothDirs(tt.target, tt.new); got != tt.want {
				t.Errorf("bothDirs() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestRenameOrCopy covers the move strategy: direct rename on the same
// filesystem, copy-and-delete fallback on EXDEV, and as-is return of other
// errors with both sides left untouched.
func TestRenameOrCopy(t *testing.T) {
	t.Run("same filesystem renames", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "src.bin")
		dst := filepath.Join(dir, "dst.bin")
		writeFile(t, src, []byte("NEW"))
		writeFile(t, dst, []byte("OLD"))

		if err := renameOrCopy(src, dst); err != nil {
			t.Fatalf("renameOrCopy: %v", err)
		}
		if got := readFile(t, dst); string(got) != "NEW" {
			t.Errorf("dst contents: %q (want NEW)", got)
		}
		if _, err := os.Stat(src); !os.IsNotExist(err) {
			t.Errorf("src should have been renamed away: %v", err)
		}
	})

	t.Run("exdev falls back to copy and removes src", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "src.bin")
		dst := filepath.Join(dir, "dst.bin")
		writeFile(t, src, []byte("NEW"))
		writeFile(t, dst, []byte("OLD"))
		withRenameFunc(t, func(oldpath, newpath string) error {
			return crossDeviceErr(oldpath, newpath)
		})

		if err := renameOrCopy(src, dst); err != nil {
			t.Fatalf("renameOrCopy: %v", err)
		}
		if got := readFile(t, dst); string(got) != "NEW" {
			t.Errorf("dst contents: %q (want NEW)", got)
		}
		if _, err := os.Stat(src); !os.IsNotExist(err) {
			t.Errorf("src should have been removed after copy: %v", err)
		}
	})

	t.Run("non-exdev error returns as-is without touching dst", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "src.bin")
		dst := filepath.Join(dir, "dst.bin")
		writeFile(t, src, []byte("NEW"))
		writeFile(t, dst, []byte("OLD"))
		sentinel := &os.LinkError{Op: "rename", Old: src, New: dst, Err: syscall.EACCES}
		withRenameFunc(t, func(_, _ string) error { return sentinel })

		err := renameOrCopy(src, dst)
		if !errors.Is(err, syscall.EACCES) {
			t.Fatalf("err = %v (want EACCES)", err)
		}
		if got := readFile(t, dst); string(got) != "OLD" {
			t.Errorf("dst must be untouched, got %q", got)
		}
		if got := readFile(t, src); string(got) != "NEW" {
			t.Errorf("src must be untouched, got %q", got)
		}
	})

	t.Run("exdev copy failure is wrapped", func(t *testing.T) {
		dir := t.TempDir()
		missing := filepath.Join(dir, "missing.bin")
		dst := filepath.Join(dir, "dst.bin")
		writeFile(t, dst, []byte("OLD"))
		withRenameFunc(t, func(oldpath, newpath string) error {
			return crossDeviceErr(oldpath, newpath)
		})

		if err := renameOrCopy(missing, dst); err == nil {
			t.Fatal("expected error for missing src")
		}
	})
}

// A failed swap must never leave the install slot without a runnable
// binary: non-EXDEV failures return the error and preserve both sides.
func TestReplaceTarget_FailurePreservesOriginal(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "app.bin")
	newPath := filepath.Join(dir, "app.bin.new")
	writeFile(t, target, []byte("OLD"))
	writeFile(t, newPath, []byte("NEW"))
	withRenameFunc(t, func(oldpath, newpath string) error {
		return &os.LinkError{Op: "rename", Old: oldpath, New: newpath, Err: syscall.EACCES}
	})

	if err := replaceTarget(target, newPath); err == nil {
		t.Fatal("expected error")
	}
	if got := readFile(t, target); string(got) != "OLD" {
		t.Errorf("target must be untouched, got %q", got)
	}
	if got := readFile(t, newPath); string(got) != "NEW" {
		t.Errorf("newPath must be untouched, got %q", got)
	}
}

// Directory replaces must land clean rather than merging stale files from
// the old bundle into the new one — including across the EXDEV fallback.
func TestReplaceTarget_DirReplace(t *testing.T) {
	tests := []struct {
		name      string
		forceCopy bool
	}{
		{"same filesystem rename", false},
		{"cross-device copy", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			target := filepath.Join(dir, "App.app")
			newPath := filepath.Join(dir, "App.app.new")
			makeAppBundle(t, target, "old-bin")
			// Stale file present only in the old bundle.
			writeFile(t, filepath.Join(target, "Contents", "stale.txt"), []byte("stale"))
			makeAppBundle(t, newPath, "new-bin")
			if tt.forceCopy {
				withRenameFunc(t, func(oldpath, newpath string) error {
					return crossDeviceErr(oldpath, newpath)
				})
			}

			if err := replaceTarget(target, newPath); err != nil {
				t.Fatalf("replaceTarget: %v", err)
			}
			if got := readFile(t, filepath.Join(target, "Contents", "MacOS", "exe")); string(got) != "new-bin" {
				t.Errorf("bundle exe: %q (want new-bin)", got)
			}
			if _, err := os.Stat(filepath.Join(target, "Contents", "stale.txt")); !os.IsNotExist(err) {
				t.Errorf("stale file must not survive a clean replace: %v", err)
			}
			if _, err := os.Stat(newPath); !os.IsNotExist(err) {
				t.Errorf("newPath should have been moved away: %v", err)
			}
		})
	}
}

// End-to-end: the /tmp-on-tmpfs versus $HOME layout from the bug report.
// With rename forced to EXDEV, the helper must still swap via copy, restore
// the executable bit, launch once, and clean the backup.
func TestRunHelperSwap_CrossDeviceFile(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "app.bin")
	newPath := filepath.Join(dir, "app.bin.new")
	writeFile(t, target, []byte("OLD"))
	if err := os.Chmod(target, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, newPath, []byte("NEW")) // 0o644 — non-executable staging artifact
	withRenameFunc(t, func(oldpath, newpath string) error {
		return crossDeviceErr(oldpath, newpath)
	})

	l := &fakeLauncher{}
	code := runHelperSwap(target, newPath, 0, filepath.Join(dir, "log"), instantWaiter, l)
	if code != 0 {
		t.Fatalf("code: %d (want 0)", code)
	}
	if got := readFile(t, target); string(got) != "NEW" {
		t.Errorf("target contents: %q (want NEW)", got)
	}
	info, err := os.Stat(target)
	if err != nil {
		t.Fatal(err)
	}
	if mode := info.Mode().Perm(); mode != 0o755 {
		t.Errorf("post-swap mode: got %o, want 0755", mode)
	}
	if _, err := os.Stat(target + ".bak"); !os.IsNotExist(err) {
		t.Errorf("backup should be cleaned up: %v", err)
	}
	if len(l.calls) != 1 || l.calls[0] != target {
		t.Errorf("launcher calls: %+v", l.calls)
	}
}

// dirNames returns the base names of dir's entries. Staging tests use it to
// assert no temporary sibling survives a staged swap.
func dirNames(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return names
}

// TestRenameOrCopy_StagedFileSwap verifies the EXDEV file path lands via a
// temporary sibling: dst ends complete, src is removed, and no staging file
// is left behind.
func TestRenameOrCopy_StagedFileSwap(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.bin")
	dst := filepath.Join(dir, "dst.bin")
	writeFile(t, src, []byte("NEW"))
	writeFile(t, dst, []byte("OLD"))
	withRenameFunc(t, func(oldpath, newpath string) error {
		return crossDeviceErr(oldpath, newpath)
	})

	if err := renameOrCopy(src, dst); err != nil {
		t.Fatalf("renameOrCopy: %v", err)
	}
	if got := readFile(t, dst); string(got) != "NEW" {
		t.Errorf("dst contents: %q (want NEW)", got)
	}
	if got := dirNames(t, dir); len(got) != 1 || got[0] != "dst.bin" {
		t.Errorf("dir entries after swap: %q (want only dst.bin)", got)
	}
}

// TestRenameOrCopy_StagingFailurePreservesDst forces the staged copy to fail
// (unreadable src) and verifies dst keeps its complete original contents,
// src is left alone, and the staging file is cleaned up.
func TestRenameOrCopy_StagingFailurePreservesDst(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.bin")
	dst := filepath.Join(dir, "dst.bin")
	writeFile(t, src, []byte("NEW"))
	if err := os.Chmod(src, 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(src, 0o644) })
	writeFile(t, dst, []byte("OLD"))
	withRenameFunc(t, func(oldpath, newpath string) error {
		return crossDeviceErr(oldpath, newpath)
	})

	if err := renameOrCopy(src, dst); err == nil {
		t.Fatal("expected error for unreadable src")
	}
	if got := readFile(t, dst); string(got) != "OLD" {
		t.Errorf("dst must keep its complete original, got %q", got)
	}
	if got := dirNames(t, dir); len(got) != 2 {
		t.Errorf("dir entries after failed swap: %q (want src+dst only)", got)
	}
}

// End-to-end recovery for a failed directory swap: after bothDirs clears the
// target slot, a non-EXDEV rename failure must fall through to the backup
// restore — original bundle back at target, relaunched, no helper env.
func TestRunHelperSwap_DirReplaceFails_RestoresBackup(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "App.app")
	newPath := filepath.Join(dir, "App.app.new")
	makeAppBundle(t, target, "old-bin")
	makeAppBundle(t, newPath, "new-bin")
	for _, key := range []string{envHelperMode, envHelperTarget, envHelperNew, envHelperPID, envHelperLog} {
		t.Setenv(key, "1")
	}
	withRenameFunc(t, func(oldpath, newpath string) error {
		return &os.LinkError{Op: "rename", Old: oldpath, New: newpath, Err: syscall.EACCES}
	})

	calls := 0
	l := &funcLauncher{fn: func(path string) error {
		calls++
		if path != target {
			t.Errorf("launch path: %q (want %q)", path, target)
		}
		if got := readFile(t, filepath.Join(path, "Contents", "MacOS", "exe")); string(got) != "old-bin" {
			t.Errorf("restored bundle exe: %q (want old-bin)", got)
		}
		for _, key := range []string{envHelperMode, envHelperTarget, envHelperNew, envHelperPID, envHelperLog} {
			if value := os.Getenv(key); value != "" {
				t.Errorf("%s leaked to restored application: %q", key, value)
			}
		}
		return nil
	}}

	code := runHelperSwap(target, newPath, 0, filepath.Join(dir, "log"), instantWaiter, l)
	if code != 13 {
		t.Fatalf("code: %d (want 13)", code)
	}
	if calls != 1 {
		t.Fatalf("launches: %d (want 1)", calls)
	}
	if got := readFile(t, filepath.Join(target, "Contents", "MacOS", "exe")); string(got) != "old-bin" {
		t.Errorf("after restore bundle exe: %q (want old-bin)", got)
	}
	if got := readFile(t, filepath.Join(newPath, "Contents", "MacOS", "exe")); string(got) != "new-bin" {
		t.Errorf("staged bundle must be untouched, got %q", got)
	}
}
