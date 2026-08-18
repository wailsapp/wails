package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func TestManifestPresenceIsTheCutoverFlag(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, manifest.WriteMinimal(root, manifest.Project{Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0"}))
	require.NoError(t, os.WriteFile(filepath.Join(root, "Taskfile.yml"), []byte("version: '3'\n"), 0o644))
	active, err := activeManifestProjectAt(root)
	require.NoError(t, err)
	assert.True(t, active, "Taskfiles must be ignored once wails.hcl exists")
}

func TestMigrationDraftDoesNotActivateNativeRouting(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.MigratedFilename), manifest.Minimal(manifest.Project{Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0"}), 0o644))
	active, err := activeManifestProjectAt(root)
	require.NoError(t, err)
	assert.False(t, active)
}
