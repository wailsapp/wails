package updater

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"
)

// fakeLauncher records each launch attempt without actually exec'ing
// anything. errsByCall lets a test fail only specific invocations (e.g.
// first launch fails, second-launch-during-restore succeeds).
type fakeLauncher struct {
	mu         sync.Mutex
	calls      []string
	errsByCall []error
}

func (f *fakeLauncher) launch(path string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	idx := len(f.calls)
	f.calls = append(f.calls, path)
	if idx < len(f.errsByCall) {
		return f.errsByCall[idx]
	}
	return nil
}

// instantWaiter is a processWaiter that returns immediately. Used when the
// test doesn't want to spend wall-clock time waiting for a (nonexistent) PID.
func instantWaiter(_ int, _ time.Duration) error { return nil }

func TestRunHelperSwap_HappyPath_File(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "app.bin")
	newPath := filepath.Join(dir, "app.bin.new")
	writeFile(t, target, []byte("OLD"))
	writeFile(t, newPath, []byte("NEW"))

	l := &fakeLauncher{}
	code := runHelperSwap(target, newPath, 0, filepath.Join(dir, "log"), instantWaiter, l)
	if code != 0 {
		t.Fatalf("code: %d", code)
	}
	if got := readFile(t, target); string(got) != "NEW" {
		t.Errorf("target contents: %q", got)
	}
	if _, err := os.Stat(newPath); !os.IsNotExist(err) {
		t.Errorf("new path should have been renamed away: %v", err)
	}
	if _, err := os.Stat(target + ".bak"); !os.IsNotExist(err) {
		t.Errorf("backup should be cleaned up: %v", err)
	}
	if len(l.calls) != 1 || l.calls[0] != target {
		t.Errorf("launcher calls: %+v", l.calls)
	}
}

// The downloaded artifact is created via os.Create which masks 0o666 against
// umask — on Unix the executable bit isn't set, so a direct rename would
// produce a non-runnable binary at target. The helper must restore the
// original target's mode after the swap.
func TestRunHelperSwap_PreservesExecutableBit(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("executable bit is not a Unix concept on Windows")
	}
	dir := t.TempDir()
	target := filepath.Join(dir, "app.bin")
	newPath := filepath.Join(dir, "app.bin.new")
	writeFile(t, target, []byte("OLD"))
	if err := os.Chmod(target, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, newPath, []byte("NEW")) // 0o644 — non-executable

	l := &fakeLauncher{}
	code := runHelperSwap(target, newPath, 0, filepath.Join(dir, "log"), instantWaiter, l)
	if code != 0 {
		t.Fatalf("code: %d", code)
	}
	info, err := os.Stat(target)
	if err != nil {
		t.Fatal(err)
	}
	if mode := info.Mode().Perm(); mode != 0o755 {
		t.Errorf("post-swap mode: got %o, want 0755 — exec bit not restored", mode)
	}
}

func TestRunHelperSwap_TargetMissing_FailsEarly(t *testing.T) {
	dir := t.TempDir()
	missing := filepath.Join(dir, "missing")
	newPath := filepath.Join(dir, "new")
	writeFile(t, newPath, []byte("NEW"))

	code := runHelperSwap(missing, newPath, 0, filepath.Join(dir, "log"), instantWaiter, &fakeLauncher{})
	if code != 10 {
		t.Fatalf("code: %d (want 10)", code)
	}
}

func TestRunHelperSwap_NewMissing_FailsEarly(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "app.bin")
	writeFile(t, target, []byte("OLD"))

	code := runHelperSwap(target, filepath.Join(dir, "nope"), 0, filepath.Join(dir, "log"), instantWaiter, &fakeLauncher{})
	if code != 11 {
		t.Fatalf("code: %d (want 11)", code)
	}
}

func TestRunHelperSwap_LaunchFails_RestoresBackup(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "app.bin")
	newPath := filepath.Join(dir, "app.bin.new")
	writeFile(t, target, []byte("OLD"))
	writeFile(t, newPath, []byte("NEW"))

	l := &fakeLauncher{errsByCall: []error{errors.New("nope")}}
	code := runHelperSwap(target, newPath, 0, filepath.Join(dir, "log"), instantWaiter, l)
	if code != 15 {
		t.Fatalf("code: %d (want 15)", code)
	}
	// After restore the target should have its original contents back.
	if got := readFile(t, target); string(got) != "OLD" {
		t.Errorf("after restore target contents: %q (want OLD)", got)
	}
	// Launcher gets called twice: once for the new (fails), once for the restore.
	if len(l.calls) != 2 {
		t.Errorf("expected 2 launch attempts, got %d: %+v", len(l.calls), l.calls)
	}
}

func TestRunHelperSwap_AppBundleDirectory(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "App.app")
	newPath := filepath.Join(dir, "App.app.new")
	makeAppBundle(t, target, "old-bin")
	makeAppBundle(t, newPath, "new-bin")

	l := &fakeLauncher{}
	code := runHelperSwap(target, newPath, 0, filepath.Join(dir, "log"), instantWaiter, l)
	if code != 0 {
		t.Fatalf("code: %d", code)
	}
	if got := readFile(t, filepath.Join(target, "Contents", "MacOS", "exe")); string(got) != "new-bin" {
		t.Errorf("bundle contents after swap: %q", got)
	}
}

func TestRunHelperSwap_AppBundle_RestoreOnLaunchFailure(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "App.app")
	newPath := filepath.Join(dir, "App.app.new")
	makeAppBundle(t, target, "old-bin")
	makeAppBundle(t, newPath, "new-bin")

	l := &fakeLauncher{errsByCall: []error{errors.New("kaboom")}}
	code := runHelperSwap(target, newPath, 0, filepath.Join(dir, "log"), instantWaiter, l)
	if code != 15 {
		t.Fatalf("code: %d (want 15)", code)
	}
	if got := readFile(t, filepath.Join(target, "Contents", "MacOS", "exe")); string(got) != "old-bin" {
		t.Errorf("after restore bundle contents: %q (want old-bin)", got)
	}
}

// When the parent process refuses to exit within the wait timeout, the helper
// must abort before touching the target. On Windows the swap would otherwise
// grind against the file lock; on macOS `open -n` would launch a second
// instance alongside the still-running parent. Either way the user ends up
// with a worse state than before they tried to update.
func TestRunHelperSwap_ParentWaitTimeout_Aborts(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "app.bin")
	newPath := filepath.Join(dir, "app.bin.new")
	writeFile(t, target, []byte("OLD"))
	writeFile(t, newPath, []byte("NEW"))

	timeoutWaiter := func(_ int, _ time.Duration) error {
		return errors.New("parent still alive")
	}
	l := &fakeLauncher{}
	code := runHelperSwap(target, newPath, 1234, filepath.Join(dir, "log"), timeoutWaiter, l)
	if code != 17 {
		t.Fatalf("code: %d (want 17 — parent-timeout abort)", code)
	}
	if got := readFile(t, target); string(got) != "OLD" {
		t.Errorf("target should be untouched, got %q", got)
	}
	if len(l.calls) != 0 {
		t.Errorf("launcher must not be called after abort, got %+v", l.calls)
	}
}

// Regression for "user clicks Restart, app silently exits 11 instead of
// relaunching." Helper-mode env vars were inherited by the launched binary
// because exec.Command defaults cmd.Env to os.Environ(); HandleHelperMode at
// the top of the new process saw the still-set sentinels and ran another
// (doomed) swap. Fix: scrub the env in the helper before launching.
func TestRunHelperSwap_ClearsHelperEnvBeforeLaunch(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "app.bin")
	newPath := filepath.Join(dir, "app.bin.new")
	writeFile(t, target, []byte("OLD"))
	writeFile(t, newPath, []byte("NEW"))

	// Seed the env exactly as Restart() does before spawning the helper.
	t.Setenv(envHelperMode, "1")
	t.Setenv(envHelperTarget, target)
	t.Setenv(envHelperNew, newPath)
	t.Setenv(envHelperPID, "1234")
	t.Setenv(envHelperLog, filepath.Join(dir, "log"))

	// envAtLaunch is captured by the launcher at the moment it would spawn the
	// new binary — that's exactly the snapshot the inherited exec would see.
	var envAtLaunch map[string]string
	envCapturingLauncher := &funcLauncher{fn: func(path string) error {
		envAtLaunch = map[string]string{}
		for _, k := range []string{envHelperMode, envHelperTarget, envHelperNew, envHelperPID, envHelperLog} {
			envAtLaunch[k] = os.Getenv(k)
		}
		return nil
	}}

	code := runHelperSwap(target, newPath, 0, filepath.Join(dir, "log"), instantWaiter, envCapturingLauncher)
	if code != 0 {
		t.Fatalf("swap code: %d", code)
	}
	for k, v := range envAtLaunch {
		if v != "" {
			t.Errorf("env var %s leaked to launched process: %q", k, v)
		}
	}
}

// funcLauncher is a launcher whose launch is provided by the test. Different
// from fakeLauncher in that the callback runs custom inspection logic.
type funcLauncher struct {
	fn func(path string) error
}

func (f *funcLauncher) launch(path string) error { return f.fn(path) }

func TestHandleHelperMode_NoEnv_Returns(t *testing.T) {
	// When the sentinel env var is absent the function must return
	// immediately and NOT touch os.Exit.
	t.Setenv(envHelperMode, "")
	HandleHelperMode()
}

// --- helpers ---

func writeFile(t *testing.T, path string, body []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, body, 0o644); err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// makeAppBundle creates a minimal Contents/MacOS/exe layout under root,
// with `payload` as the executable's contents. This is the structure macOS
// treats as a .app bundle.
func makeAppBundle(t *testing.T, root, payload string) {
	t.Helper()
	exe := filepath.Join(root, "Contents", "MacOS", "exe")
	writeFile(t, exe, []byte(payload))
}
