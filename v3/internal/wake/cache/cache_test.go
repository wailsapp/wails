package cache

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSnapshotFilesIncludesExecutableMode(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "hook.sh")
	require.NoError(t, os.WriteFile(path, []byte("#!/bin/sh\n"), 0o644))
	store, err := OpenCache(root)
	require.NoError(t, err)
	plain, err := store.SnapshotFiles("hook", "hook.sh")
	require.NoError(t, err)

	require.NoError(t, os.Chmod(path, 0o755))
	executable, err := store.SnapshotFiles("hook", "hook.sh")
	require.NoError(t, err)
	assert.NotEqual(t, plain, executable)
}
