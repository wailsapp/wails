package updater_test

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/wailsapp/wails/v3/pkg/updater"
)

// TestRestartHelperProcess re-executes the real helper entry point after
// simulated application initialization, without running a graphical app.
func TestRestartHelperProcess(t *testing.T) {
	pidFile := os.Getenv("WAILS_TEST_HELPER_PID_FILE")
	if pidFile == "" {
		return
	}
	if err := os.WriteFile(pidFile, []byte(strconv.Itoa(os.Getpid())), 0o600); err != nil {
		t.Fatal(err)
	}
	delay, err := time.ParseDuration(os.Getenv("WAILS_TEST_HELPER_DELAY"))
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(delay)
	updater.HandleHelperMode()
	t.Fatal("helper mode unexpectedly returned")
}

func TestRestart_HelperStartupDeadline(t *testing.T) {
	for _, tt := range []struct {
		name    string
		delay   time.Duration
		timeout time.Duration
		wantErr bool
	}{
		{name: "default permits slow initialization", delay: 6 * time.Second},
		{name: "configured timeout permits startup", delay: 150 * time.Millisecond, timeout: 2 * time.Second},
		{name: "configured timeout stops helper", delay: time.Second, timeout: 100 * time.Millisecond, wantErr: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			pidFile := filepath.Join(t.TempDir(), "helper.pid")
			t.Setenv("WAILS_TEST_HELPER_PID_FILE", pidFile)
			t.Setenv("WAILS_TEST_HELPER_DELAY", tt.delay.String())
			t.Cleanup(updater.SetSelfExecutableForTest(os.Executable))
			t.Cleanup(updater.SetNewDetachedCommandForTest(func(path string) *exec.Cmd {
				return exec.Command(path, "-test.run=^TestRestartHelperProcess$")
			}))

			host := &fakeHost{}
			body := []byte("test payload")
			p := &fakeProvider{name: "p", rel: &updater.Release{
				Version: "2.0.0", Artifact: updater.Artifact{Filename: "app.bin", Size: int64(len(body))},
			}, body: body}
			u := updater.New(host)
			if err := u.Init(updater.Config{CurrentVersion: "1.0.0", Providers: []updater.Provider{p}, HelperReadyTimeout: tt.timeout}); err != nil {
				t.Fatal(err)
			}
			if _, err := u.Check(context.Background()); err != nil {
				t.Fatal(err)
			}
			if err := u.DownloadAndInstall(context.Background()); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = os.RemoveAll(filepath.Dir(u.DownloadedPath())) })
			err := u.Restart(context.Background())
			if err == nil {
				stopRestartHelper(t, pidFile)
			}
			if tt.wantErr {
				if !errors.Is(err, updater.ErrHelperNotReady) || host.quits != 0 {
					t.Fatalf("timeout: error=%v, quits=%d", err, host.quits)
				}
			} else if err != nil || host.quits != 1 {
				t.Fatalf("startup: error=%v, quits=%d", err, host.quits)
			}
		})
	}
}

func TestInit_RejectsNegativeHelperReadyTimeout(t *testing.T) {
	u := updater.New(&fakeHost{})
	err := u.Init(updater.Config{CurrentVersion: "1.0.0", Providers: []updater.Provider{&fakeProvider{}}, HelperReadyTimeout: -time.Second})
	if err == nil {
		t.Fatal("negative helper startup timeout was accepted")
	}
}

// stopRestartHelper prevents a successful helper from swapping the test binary
// when the fake host leaves its parent running.
func stopRestartHelper(t *testing.T, pidFile string) {
	t.Helper()
	// The fake host keeps the parent alive. Stop and reap its waiting
	// helper before the test process exits and it could attempt a swap.
	raw, readErr := os.ReadFile(pidFile)
	if readErr != nil {
		t.Fatal(readErr)
	}
	pid, parseErr := strconv.Atoi(string(raw))
	if parseErr != nil || pid <= 0 {
		t.Fatalf("invalid helper pid %q", raw)
	}
	process, findErr := os.FindProcess(pid)
	if findErr != nil {
		t.Fatal(findErr)
	}
	_ = process.Kill()
	_, _ = process.Wait()
}
