package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func TestCleanRemovesOnlyGeneratedWorkspace(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, manifest.WriteMinimal(root, manifest.Project{Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0"}))
	require.NoError(t, os.MkdirAll(filepath.Join(root, ".wails", "build"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, ".wails", "build", "generated"), []byte("state"), 0o644))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "bin"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "bin", "user-artifact"), []byte("keep"), 0o644))
	previous, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(root))
	t.Cleanup(func() { _ = os.Chdir(previous) })
	require.NoError(t, Clean(nil))
	assert.NoDirExists(t, filepath.Join(root, ".wails"))
	assert.FileExists(t, filepath.Join(root, "bin", "user-artifact"))
}
