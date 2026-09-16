//go:build unix

package benchmark

import (
	"path/filepath"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSnapshotArtifactRejectsNamedPipe(t *testing.T) {
	root := t.TempDir()
	fifo := filepath.Join(root, "fifo")
	require.NoError(t, syscall.Mkfifo(fifo, 0o600))
	_, err := snapshotArtifact(fifo)
	assert.ErrorContains(t, err, "not a regular file or directory")

	directory := filepath.Join(root, "tree")
	require.NoError(t, syscall.Mkdir(directory, 0o755))
	require.NoError(t, syscall.Mkfifo(filepath.Join(directory, "fifo"), 0o600))
	_, err = snapshotArtifact(directory)
	assert.ErrorContains(t, err, "contains non-regular file")
}
