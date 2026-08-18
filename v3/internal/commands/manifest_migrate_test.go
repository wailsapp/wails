package commands

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func TestMigrationDraftIsInactiveAndNonDestructive(t *testing.T) {
	root := t.TempDir()
	doc := manifest.NewDocument(manifest.Project{Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0"})
	require.NoError(t, manifest.WriteMigrationDraft(root, doc))
	require.FileExists(t, filepath.Join(root, manifest.MigratedFilename))
	_, err := manifest.Load(root, "")
	require.ErrorContains(t, err, "could not find wails.hcl")
}

func TestActivateMigrationRenamesOnlyTheReviewedDraft(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.MigratedFilename), manifest.Minimal(manifest.Project{Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0"}), 0o644))
	// A stock-shaped Taskfile gives analysis something to inspect. It remains in
	// place after activation and is ignored only because wails.hcl now exists.
	require.NoError(t, os.WriteFile(filepath.Join(root, "Taskfile.yml"), []byte("version: '3'\ntasks: {}\n"), 0o644))
	previous, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(root))
	t.Cleanup(func() { _ = os.Chdir(previous) })
	require.NoError(t, Migrate(&MigrateOptions{Activate: true}))
	assert.FileExists(t, manifest.Filename)
	assert.FileExists(t, "Taskfile.yml")
	assert.NoFileExists(t, manifest.MigratedFilename)
}

func TestActivateMigrationNeverOverwritesAnActiveManifest(t *testing.T) {
	root := t.TempDir()
	draft := manifest.Minimal(manifest.Project{Name: "draft", ProductName: "Draft", Identifier: "com.example.draft", Version: "1.0.0"})
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.MigratedFilename), draft, 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "Taskfile.yml"), []byte("version: '3'\ntasks: {}\n"), 0o644))
	active := []byte("existing manifest\n")
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), active, 0o644))

	err := activateMigration(root, &MigrateOptions{Activate: true})
	require.ErrorContains(t, err, "refusing to overwrite")
	actual, readErr := os.ReadFile(filepath.Join(root, manifest.Filename))
	require.NoError(t, readErr)
	assert.Equal(t, active, actual)
	assert.FileExists(t, filepath.Join(root, manifest.MigratedFilename))
}

func TestActivateMigrationDryRunEmitsJSONWhenRequested(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.MigratedFilename), manifest.Minimal(manifest.Project{Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0"}), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "Taskfile.yml"), []byte("version: '3'\ntasks: {}\n"), 0o644))

	output := captureMigrationStdout(t, func() {
		require.NoError(t, activateMigration(root, &MigrateOptions{Activate: true, DryRun: true, JSON: true}))
	})
	assert.Contains(t, output, `"complete": true`)
	assert.NoFileExists(t, filepath.Join(root, manifest.Filename))
	assert.FileExists(t, filepath.Join(root, manifest.MigratedFilename))
}

func captureMigrationStdout(t *testing.T, fn func()) string {
	t.Helper()
	original := os.Stdout
	reader, writer, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = writer
	done := make(chan string, 1)
	go func() {
		var buffer bytes.Buffer
		_, _ = io.Copy(&buffer, reader)
		done <- buffer.String()
	}()
	t.Cleanup(func() {
		os.Stdout = original
		_ = reader.Close()
	})
	fn()
	require.NoError(t, writer.Close())
	os.Stdout = original
	return <-done
}
