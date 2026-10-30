//go:build darwin || linux

package commands

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/atterpac/refresh/engine"
	"github.com/atterpac/refresh/process"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A closing terminal sends SIGHUP. The watcher must stop and kill the primary
// process and its children, which run in their own process group and never
// see the hangup (#6328).
func TestWatcherStopsProcessesOnSIGHUP(t *testing.T) {
	dir := t.TempDir()
	pidFile := filepath.Join(dir, "primary.pid")
	childFile := filepath.Join(dir, "child.pid")

	watcherEngine, err := engine.NewEngineFromConfig(engine.Config{
		RootPath: dir,
		LogLevel: "error",
		ExecStruct: []process.Execute{{
			// Like `wails3 task run` starting the app: the primary spawns a child.
			Cmd:        `echo $$ > "$PID_FILE"; sleep 300 & echo $! > "$CHILD_FILE"; wait`,
			Env:        map[string]string{"PID_FILE": pidFile, "CHILD_FILE": childFile},
			Type:       process.Primary,
			ExitPolicy: process.ExitPolicyShutdown,
		}},
	})
	require.NoError(t, err)

	done := make(chan error, 1)
	go func() { done <- runWatcherEngine(watcherEngine, false) }()

	var pid, child int
	require.Eventually(t, func() bool {
		var err error
		if pid, err = readPID(pidFile); err != nil {
			return false
		}
		child, err = readPID(childFile)
		return err == nil
	}, 10*time.Second, 20*time.Millisecond, "primary process did not start")

	require.NoError(t, syscall.Kill(os.Getpid(), syscall.SIGHUP))

	select {
	case err := <-done:
		assert.NoError(t, err)
	case <-time.After(10 * time.Second):
		t.Fatal("watcher did not stop after SIGHUP")
	}
	assert.ErrorIs(t, syscall.Kill(pid, 0), syscall.ESRCH, "primary process outlived the watcher")
	assert.Eventually(t, func() bool { return processExited(child) },
		10*time.Second, 20*time.Millisecond, "child process outlived the watcher")
}

func readPID(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(data)))
}

// processExited reports whether pid is gone. An orphaned child can linger as a
// zombie until its new parent reaps it, which also counts as exited.
func processExited(pid int) bool {
	err := syscall.Kill(pid, 0)
	if errors.Is(err, syscall.ESRCH) {
		return true
	}
	if err != nil {
		return false
	}
	state, err := exec.Command("ps", "-o", "stat=", "-p", strconv.Itoa(pid)).Output()
	return err == nil && strings.HasPrefix(strings.TrimSpace(string(state)), "Z")
}
