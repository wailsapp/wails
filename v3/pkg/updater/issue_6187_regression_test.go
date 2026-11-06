package updater

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestWorkflowReproHelperTimeoutLeavesOriginalRelaunchable encodes the report's
// exit-17 path: when the parent does not exit within the wait window the helper
// returns 17 without launching anything, even though a working original binary
// is still on disk. The assertion records the actual observable outcome so that
// a run demonstrates (or refutes) that no relaunch happened after the abort.
func TestWorkflowReproHelperTimeoutLeavesOriginalRelaunchable(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TMPDIR", dir)
	t.Setenv("TEMP", dir)
	t.Setenv("TMP", dir)

	target := filepath.Join(dir, "app.bin")
	newPath := filepath.Join(dir, "app.bin.new")
	if err := os.WriteFile(target, []byte("OLD"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(newPath, []byte("NEW"), 0o755); err != nil {
		t.Fatal(err)
	}

	// seed helper env exactly as Restart would, then let runHelperSwap proceed.
	t.Setenv(envHelperMode, "1")
	t.Setenv(envHelperTarget, target)
	t.Setenv(envHelperNew, newPath)

	triggerReached := false
	wait := func(_ int, d time.Duration) error {
		triggerReached = true
		if d <= 0 {
			t.Fatalf("waiter must receive a positive timeout, got %v", d)
		}
		return errors.New("parent still alive past timeout")
	}

	var launches []string
	l := &funcLauncher{fn: func(path string) error {
		launches = append(launches, path)
		if strings.Contains(strings.Join(os.Environ(), "\n"), envHelperMode+"=1") {
			t.Errorf("REPRO_ASSERTION: helper env %s leaked to relaunched process", envHelperMode)
		}
		return nil
	}}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	done := make(chan int, 1)
	go func() {
		done <- runHelperSwap(target, newPath, 4242, filepath.Join(dir, "log"), wait, l)
	}()

	var code int
	select {
	case code = <-done:
	case <-ctx.Done():
		t.Fatal("runHelperSwap did not return within 5s while waiter returned immediately")
	}

	if !triggerReached {
		t.Fatalf("positive control failed: timeout waiter was never invoked (code=%d)", code)
	}

	oldBytes, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("reading target: %v", err)
	}

	t.Logf("REPRO_OBSERVATION: exit_code=%d launches=%d target_unchanged=%t", code, len(launches), string(oldBytes) == "OLD")

	if code != 17 {
		t.Fatalf("setup: expected parent-timeout abort code 17, got %d", code)
	}
	if len(launches) == 0 {
		t.Errorf("REPRO_ASSERTION: parent-timeout abort exit_code=%d launches=%d target_unchanged=%t; original intact but nothing was relaunched", code, len(launches), string(oldBytes) == "OLD")
	}
}
