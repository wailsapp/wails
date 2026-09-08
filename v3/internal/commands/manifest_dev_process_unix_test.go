//go:build !windows

package commands

import (
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

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
