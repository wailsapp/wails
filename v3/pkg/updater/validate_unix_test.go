//go:build !windows

package updater

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunHelperSwap_PostSwapValidationRestoresBackup(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "App.app")
	staged := filepath.Join(dir, "staging", "App.app")
	makeAppBundle(t, target, "old-bin")
	makeAppBundle(t, staged, "new-bin")
	originalRename := renameFunc
	t.Cleanup(func() { renameFunc = originalRename })
	renameFunc = func(src, dst string) error {
		if err := os.Rename(src, dst); err != nil {
			return err
		}
		return os.Truncate(filepath.Join(dst, "Contents", "Info.plist"), 0)
	}
	l := &funcLauncher{fn: func(string) error { return nil }}
	code := runHelperSwap(target, staged, 0, filepath.Join(dir, "log"), instantWaiter, l)
	if code != 19 {
		t.Fatalf("code = %d, want 19 (post-swap rollback)", code)
	}
	if got := readFile(t, filepath.Join(target, "Contents", "MacOS", "exe")); string(got) != "old-bin" {
		t.Fatalf("original not restored: %q", got)
	}
}
