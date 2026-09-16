//go:build !windows

package commands

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartAndroidAVDCleansUpANewProcessWhenStartupIsCancelled(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	tools := t.TempDir()
	pidFile := filepath.Join(root, "emulator.pid")
	emulator := "#!/bin/sh\nprintf '%s' \"$$\" > \"$EMULATOR_PID_FILE\"\nexec sleep 30\n"
	require.NoError(t, os.WriteFile(filepath.Join(tools, "emulator"), []byte(emulator), 0o755))
	t.Setenv("PATH", tools+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("EMULATOR_PID_FILE", pidFile)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errCh := make(chan error, 1)
	go func() {
		_, err := startAndroidAVD(ctx, "Pixel_8",
			func(context.Context) ([]androidDevice, error) { return nil, nil },
			func(context.Context, ...string) (string, error) { return "", nil },
		)
		errCh <- err
	}()
	require.Eventually(t, func() bool {
		_, err := os.Stat(pidFile)
		return err == nil
	}, 3*time.Second, 10*time.Millisecond, "emulator process did not start")
	cancel()
	err := <-errCh
	require.Error(t, err)
	pidText, readErr := os.ReadFile(pidFile)
	require.NoError(t, readErr)
	pid, parseErr := strconv.Atoi(strings.TrimSpace(string(pidText)))
	require.NoError(t, parseErr)
	t.Cleanup(func() { _ = syscall.Kill(pid, syscall.SIGKILL) })
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if signalErr := syscall.Kill(pid, 0); errors.Is(signalErr, syscall.ESRCH) {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("new emulator process %d survived cancelled startup", pid)
}

func TestStartAndroidAVDPreservesANewProcessAfterReadiness(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	tools := t.TempDir()
	pidFile := filepath.Join(root, "emulator.pid")
	emulator := "#!/bin/sh\nprintf '%s' \"$$\" > \"$EMULATOR_PID_FILE\"\nexec sleep 30\n"
	require.NoError(t, os.WriteFile(filepath.Join(tools, "emulator"), []byte(emulator), 0o755))
	t.Setenv("PATH", tools+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("EMULATOR_PID_FILE", pidFile)
	device := androidDevice{Serial: "emulator-5554", State: "device", Emulator: true}
	discoveries := 0
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	actual, err := startAndroidAVD(ctx, "Pixel_8",
		func(context.Context) ([]androidDevice, error) {
			discoveries++
			if discoveries == 1 {
				return nil, nil
			}
			if _, err := os.Stat(pidFile); err != nil {
				return nil, nil
			}
			return []androidDevice{device}, nil
		},
		func(_ context.Context, arguments ...string) (string, error) {
			switch arguments[len(arguments)-1] {
			case "name":
				return "Pixel_8\nOK\n", nil
			case "sys.boot_completed":
				return "1\n", nil
			case "android":
				return "package:/system/framework/framework-res.apk\n", nil
			default:
				return "", errors.New("unexpected adb call")
			}
		},
	)
	require.NoError(t, err)
	assert.Equal(t, device, actual)
	pidText, readErr := os.ReadFile(pidFile)
	require.NoError(t, readErr)
	pid, parseErr := strconv.Atoi(strings.TrimSpace(string(pidText)))
	require.NoError(t, parseErr)
	t.Cleanup(func() { _ = syscall.Kill(pid, syscall.SIGKILL) })
	require.NoError(t, syscall.Kill(pid, 0), "ready emulator must remain running")
}

func TestStartAndroidAVDReportsANewProcessThatExitsBeforeReadiness(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	tools := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(tools, "emulator"), []byte("#!/bin/sh\nexit 7\n"), 0o755))
	t.Setenv("PATH", tools+string(os.PathListSeparator)+os.Getenv("PATH"))
	_, err := startAndroidAVD(context.Background(), "Pixel_8",
		func(context.Context) ([]androidDevice, error) { return nil, nil },
		func(context.Context, ...string) (string, error) { return "", nil },
	)
	require.ErrorContains(t, err, "exited before becoming ready")
}
