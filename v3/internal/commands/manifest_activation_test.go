package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/migration"
)

func TestManifestActivationUsesPrivateMigrationReport(t *testing.T) {
	writeManifest := func(t *testing.T, root string) {
		t.Helper()
		require.NoError(t, manifest.WriteMinimal(root, manifest.Project{Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0"}))
	}
	writeTaskfile := func(t *testing.T, root string) {
		t.Helper()
		require.NoError(t, os.WriteFile(filepath.Join(root, "Taskfile.yml"), []byte("version: '3'\n"), 0o644))
	}

	t.Run("native manifest", func(t *testing.T) {
		root := t.TempDir()
		writeManifest(t, root)
		active, err := activeManifestProjectAt(root)
		require.NoError(t, err)
		assert.True(t, active)
	})

	t.Run("pending migration", func(t *testing.T) {
		root := t.TempDir()
		writeManifest(t, root)
		writeTaskfile(t, root)
		require.NoError(t, migration.Write(root, migration.Report{Version: 1, Complete: false, Sources: map[string]string{}}))
		active, err := activeManifestProjectAt(root)
		require.NoError(t, err)
		assert.False(t, active)
	})

	t.Run("completed migration", func(t *testing.T) {
		root := t.TempDir()
		writeManifest(t, root)
		writeTaskfile(t, root)
		require.NoError(t, migration.Write(root, migration.Report{Version: 1, Complete: true, Sources: map[string]string{}}))
		active, err := activeManifestProjectAt(root)
		require.NoError(t, err)
		assert.True(t, active)
	})

	t.Run("ambiguous without report", func(t *testing.T) {
		root := t.TempDir()
		writeManifest(t, root)
		writeTaskfile(t, root)
		_, err := activeManifestProjectAt(root)
		require.ErrorContains(t, err, "both wails.toml and a legacy Taskfile exist")
	})
}
