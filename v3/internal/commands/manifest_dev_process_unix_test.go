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
	case "parent":
		child := exec.Command(os.Args[0], "-test.run=TestManifestProcessTreeHelper")
		child.Env = mergeEnvironment(os.Environ(), []string{"WAILS_PROCESS_TREE_HELPER=child"})
		if err := child.Start(); err != nil {
			os.Exit(10)
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
