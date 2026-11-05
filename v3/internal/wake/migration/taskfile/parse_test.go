package taskfile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseBasic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "Taskfile.yml")
	err := os.WriteFile(path, []byte(`
version: "3"
tasks:
  hello:
    cmds:
      - echo hello
`), 0644)
	require.NoError(t, err)

	tf, err := Parse(path)
	require.NoError(t, err)
	require.Equal(t, "3", tf.Version)
	require.Contains(t, tf.Tasks, "hello")
}

func TestParseNoTasksYieldsNonNilMap(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "Taskfile.yml")
	err := os.WriteFile(path, []byte(`
version: "3"
vars:
  FOO: bar
`), 0644)
	require.NoError(t, err)

	tf, err := Parse(path)
	require.NoError(t, err)
	require.NotNil(t, tf.Tasks, "Tasks must be non-nil even with no tasks: key")
	require.Empty(t, tf.Tasks)
}

func TestParseIncludes(t *testing.T) {
	dir := t.TempDir()

	mainPath := filepath.Join(dir, "Taskfile.yml")
	err := os.WriteFile(mainPath, []byte(`
version: "3"
includes:
  sub:
    taskfile: ./sub/Taskfile.yml
    dir: ./sub
tasks:
  main:
    cmds:
      - echo main
`), 0644)
	require.NoError(t, err)

	subDir := filepath.Join(dir, "sub")
	err = os.MkdirAll(subDir, 0755)
	require.NoError(t, err)

	subPath := filepath.Join(subDir, "Taskfile.yml")
	err = os.WriteFile(subPath, []byte(`
version: "3"
tasks:
  sub:
    cmds:
      - echo sub
`), 0644)
	require.NoError(t, err)

	tf, err := Parse(mainPath)
	require.NoError(t, err)
	require.Contains(t, tf.Includes, "sub")

	require.Nil(t, tf.Includes["sub"].Resolved, "migration inspects includes without executing or resolving them")
}
