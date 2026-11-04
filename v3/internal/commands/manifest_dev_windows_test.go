//go:build windows

package commands

import (
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func init() {
	if os.Getenv("WAILS_WINDOWS_RESTART_HELPER") != "1" {
		return
	}
	listener, err := net.Listen("tcp", "127.0.0.1:"+os.Getenv(wailsVitePort))
	if err != nil {
		os.Exit(23)
	}
	defer listener.Close()
	for {
		connection, err := listener.Accept()
		if err != nil {
			os.Exit(24)
		}
		connection.Close()
	}
}

func TestManifestWindowsBackendStartupFailureRestoresPreviousImage(t *testing.T) {
	root := t.TempDir()
	image, err := os.ReadFile(os.Args[0])
	require.NoError(t, err)
	binary := filepath.Join(root, "backend.exe")
	require.NoError(t, os.WriteFile(binary, image, 0755))
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := listener.Addr().(*net.TCPAddr).Port
	address := listener.Addr().String()
	require.NoError(t, listener.Close())
	t.Setenv("WAILS_WINDOWS_RESTART_HELPER", "1")
	args := []string{"--config-path", `C:\Users\Example User\testing.yaml`, "--help", ""}
	old, err := startManifestAppWithArguments(root, binary, "", port, nil, args)
	require.NoError(t, err)
	defer old.stop(time.Second)
	require.NoError(t, waitForProcessTCP(t.Context(), old, address, 5*time.Second))
	require.Equal(t, image, old.restartImage)
	old.stop(time.Second)
	require.NoError(t, os.WriteFile(binary, []byte("invalid replacement executable"), 0755))
	failed, err := startManifestApp(root, binary, "", port)
	require.Error(t, err)
	require.Nil(t, failed)
	restored, err := restoreManifestWindowsApp(t.Context(), root, old, "", port)
	require.NoError(t, err)
	defer restored.stop(time.Second)
	require.Equal(t, args, restored.cmd.Args[1:], "rollback must retain the previous arguments")
	require.NotEqual(t, old.cmd.Process.Pid, restored.cmd.Process.Pid)
	require.NoError(t, waitForProcessTCP(t.Context(), restored, address, 5*time.Second))
	restoredImage, err := os.ReadFile(binary)
	require.NoError(t, err)
	require.Equal(t, image, restoredImage)
}
