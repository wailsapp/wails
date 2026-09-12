package updater

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRunHelperSwap_ReplaceFails_RelaunchesWithoutHelperEnv(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "app.bin")
	newPath := filepath.Join(dir, "new.bin")
	writeFile(t, target, []byte("OLD"))
	writeFile(t, newPath, []byte("NEW"))
	for _, key := range []string{envHelperMode, envHelperTarget, envHelperNew, envHelperPID, envHelperLog} {
		t.Setenv(key, "1")
	}
	calls := 0
	l := &funcLauncher{fn: func(path string) error {
		calls++
		if path != target || string(readFile(t, path)) != "OLD" {
			t.Errorf("expected restored application, got %q", path)
		}
		for _, key := range []string{envHelperMode, envHelperTarget, envHelperNew, envHelperPID, envHelperLog} {
			if value := os.Getenv(key); value != "" {
				t.Errorf("%s leaked to restored application: %q", key, value)
			}
		}
		return nil
	}}
	// Remove the staged file after validation to exercise failed replacement
	// and rollback through the full helper flow.
	wait := func(int, time.Duration) error { return os.Remove(newPath) }
	code := runHelperSwap(target, newPath, 1234, filepath.Join(dir, "log"), wait, l)
	if code != 13 || calls != 1 {
		t.Fatalf("code=%d launches=%d, want 13 and 1", code, calls)
	}
}

func TestRunHelperSwap_BackupFails_RelaunchesOriginal(t *testing.T) {
	for _, launchFails := range []bool{false, true} {
		t.Run(map[bool]string{false: "recovered", true: "launch_error"}[launchFails], func(t *testing.T) {
			dir := t.TempDir()
			// The target name fits, but adding .bak exceeds the 255-byte
			// component limit on common filesystems.
			target := filepath.Join(dir, strings.Repeat("a", 252))
			newPath := filepath.Join(dir, "new.bin")
			writeFile(t, target, []byte("OLD"))
			writeFile(t, newPath, []byte("NEW"))
			t.Setenv(envHelperMode, "1")
			waited, calls := false, 0
			wait := func(int, time.Duration) error { waited = true; return nil }
			l := &funcLauncher{fn: func(path string) error {
				calls++
				if !waited || path != target || os.Getenv(envHelperMode) != "" {
					t.Error("recovery must wait for parent and launch original without helper mode")
				}
				if launchFails {
					return errors.New("launch failed")
				}
				return nil
			}}
			code := runHelperSwap(target, newPath, 1234, filepath.Join(dir, "log"), wait, l)
			if code != 12 || calls != 1 || string(readFile(t, target)) != "OLD" {
				t.Fatalf("code=%d launches=%d; original must remain intact", code, calls)
			}
		})
	}
}
