//go:build darwin || linux

package commands

import (
	"os"
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
// process, which runs in its own process group and never sees the hangup (#6328).
func TestWatcherStopsProcessesOnSIGHUP(t *testing.T) {
	dir := t.TempDir()
	pidFile := filepath.Join(dir, "primary.pid")

	watcherEngine, err := engine.NewEngineFromConfig(engine.Config{
		RootPath: dir,
		LogLevel: "error",
		ExecStruct: []process.Execute{{
			Cmd:        "echo $$ > " + pidFile + "; exec sleep 300",
			Type:       process.Primary,
			ExitPolicy: process.ExitPolicyShutdown,
		}},
	})
	require.NoError(t, err)

	done := make(chan error, 1)
	go func() { done <- runWatcherEngine(watcherEngine, false) }()

	var pid int
	require.Eventually(t, func() bool {
		data, err := os.ReadFile(pidFile)
		if err != nil {
			return false
		}
		pid, err = strconv.Atoi(strings.TrimSpace(string(data)))
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
}
