package migration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadAbsentReport(t *testing.T) {
	_, exists, err := Read(t.TempDir())
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestWriteReadReport(t *testing.T) {
	root := t.TempDir()
	want := Report{Version: 1, CompletedBy: "v3", Complete: false, Sources: map[string]string{"Taskfile.yml": "digest"}}
	require.NoError(t, Write(root, want))
	got, exists, err := Read(root)
	require.NoError(t, err)
	assert.True(t, exists)
	assert.Equal(t, want, got)
}

func TestReadRejectsMalformedReport(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Dir(Path(root)), 0o755))
	require.NoError(t, os.WriteFile(Path(root), []byte("{"), 0o644))
	_, exists, err := Read(root)
	assert.True(t, exists)
	require.ErrorContains(t, err, "parse .wails/migration-report.json")
}
