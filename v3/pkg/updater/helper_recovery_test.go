package updater

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRunHelperSwap_StagedPayloadVanishes_RelaunchesWithoutHelperEnv(t *testing.T) {
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
	// Remove the staged file while the helper waits for the parent, so the
	// payload is gone by the time the swap would run.
	//
	// This is caught by the post-wait validation, before the backup is taken
	// and before the target is removed — hence 18 rather than the 13 the
	// replace loop would report. Detecting it earlier is the point: nothing is
	// disturbed, so recovery is a relaunch rather than a restore that could
	// itself fail.
	wait := func(int, time.Duration) error { return os.Remove(newPath) }
	code := runHelperSwap(target, newPath, 1234, filepath.Join(dir, "log"), wait, l)
	if code != 18 || calls != 1 {
		t.Fatalf("code=%d launches=%d, want 18 and 1", code, calls)
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
