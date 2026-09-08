//go:build !windows

package commands

import (
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func TestManifestProcessStopTerminatesChildProcessGroup(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	address := listener.Addr().String()
	require.NoError(t, listener.Close())
	ready := filepath.Join(t.TempDir(), "ready")
	process, err := startManifestProcess(t.TempDir(), os.Args[0], []string{
		"WAILS_PROCESS_TREE_HELPER=parent",
		"WAILS_PROCESS_TREE_ADDRESS=" + address,
		"WAILS_PROCESS_TREE_READY=" + ready,
	}, "-test.run=TestManifestProcessTreeHelper")
	require.NoError(t, err)
	require.Eventually(t, func() bool {
		_, statErr := os.Stat(ready)
		return statErr == nil
	}, 2*time.Second, 10*time.Millisecond)
	connection, err := net.DialTimeout("tcp", address, 200*time.Millisecond)
	require.NoError(t, err)
	require.NoError(t, connection.Close())

	process.stop(50 * time.Millisecond)
	require.Eventually(t, func() bool {
		connection, dialErr := net.DialTimeout("tcp", address, 20*time.Millisecond)
		if dialErr == nil {
			_ = connection.Close()
			return false
		}
		return true
	}, 2*time.Second, 20*time.Millisecond)
}

func TestManifestProcessTreeHelper(t *testing.T) {
	switch os.Getenv("WAILS_PROCESS_TREE_HELPER") {
	case "parent", "orphan-parent":
		child := exec.Command(os.Args[0], "-test.run=TestManifestProcessTreeHelper")
		child.Env = mergeEnvironment(os.Environ(), []string{"WAILS_PROCESS_TREE_HELPER=child"})
		if err := child.Start(); err != nil {
			os.Exit(10)
		}
		if os.Getenv("WAILS_PROCESS_TREE_HELPER") == "orphan-parent" {
			deadline := time.Now().Add(2 * time.Second)
			for time.Now().Before(deadline) {
				if _, err := os.Stat(os.Getenv("WAILS_PROCESS_TREE_READY")); err == nil {
					return
				}
				time.Sleep(10 * time.Millisecond)
			}
			os.Exit(14)
		}
		if err := child.Wait(); err != nil {
			os.Exit(11)
		}
	case "child":
		listener, err := net.Listen("tcp", os.Getenv("WAILS_PROCESS_TREE_ADDRESS"))
		if err != nil {
			os.Exit(12)
		}
		defer listener.Close()
		if err := os.WriteFile(os.Getenv("WAILS_PROCESS_TREE_READY"), []byte("ready"), 0o600); err != nil {
			os.Exit(13)
		}
		select {}
	}
}

func TestManifestProcessCleansChildrenWhenWrapperExits(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	address := listener.Addr().String()
	require.NoError(t, listener.Close())
	ready := filepath.Join(t.TempDir(), "ready")
	process, err := startManifestProcess(t.TempDir(), os.Args[0], []string{
		"WAILS_PROCESS_TREE_HELPER=orphan-parent", "WAILS_PROCESS_TREE_ADDRESS=" + address, "WAILS_PROCESS_TREE_READY=" + ready,
	}, "-test.run=^TestManifestProcessTreeHelper$")
	require.NoError(t, err)
	t.Cleanup(func() { _ = killManifestProcess(process.cmd.Process) })
	select {
	case <-process.done:
	case <-time.After(3 * time.Second):
		t.Fatal("wrapper did not exit")
	}
	require.NoError(t, process.waitError())
	require.FileExists(t, ready)
	process.stop(10 * time.Millisecond)
	require.Eventually(t, func() bool {
		connection, err := net.DialTimeout("tcp", address, 20*time.Millisecond)
		if err == nil {
			_ = connection.Close()
			return false
		}
		return true
	}, time.Second, 10*time.Millisecond)
}

func TestManifestFrontendDevReceivesDeclaredEnvironment(t *testing.T) {
	root := t.TempDir()
	config := manifest.Config{Frontend: manifest.Frontend{Directory: ".", Dev: []string{"sh", "-c", "printf '%s' \"$WAILS_TEST_FRONTEND_VALUE\" > env-result"}, Environment: map[string]string{"WAILS_TEST_FRONTEND_VALUE": "configured-value"}}}
	process, err := startFrontendDev(root, config, "127.0.0.1", 9245, "http://127.0.0.1:9245")
	require.NoError(t, err)
	t.Cleanup(func() { process.stop(time.Second) })
	select {
	case <-process.done:
	case <-time.After(time.Second):
		t.Fatal("frontend did not finish")
	}
	require.NoError(t, process.waitError())
	value, err := os.ReadFile(filepath.Join(root, "env-result"))
	require.NoError(t, err)
	require.Equal(t, "configured-value", string(value))
	next := config
	next.Frontend.Environment = map[string]string{"WAILS_TEST_FRONTEND_VALUE": "changed"}
	require.True(t, frontendSessionChanged(config, next), "environment-only reload must restart the frontend")
}
