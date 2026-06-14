package parse

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
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

	err = ResolveIncludes(tf)
	require.NoError(t, err)
	require.NotNil(t, tf.Includes["sub"].Resolved)
	require.Contains(t, tf.Tasks, "sub:sub")
}

func TestResolveVars(t *testing.T) {
	vars := map[string]*ast.Var{
		"A": {Static: "hello"},
		"B": {Ref: ".A"},
	}

	err := ResolveVars(vars)
	require.NoError(t, err)
	require.Equal(t, "hello", vars["A"].Value)
	require.Equal(t, "hello", vars["B"].Value)
}

func TestResolveVarsCycle(t *testing.T) {
	vars := map[string]*ast.Var{
		"A": {Ref: ".B"},
		"B": {Ref: ".A"},
	}

	err := ResolveVars(vars)
	require.Error(t, err)
	require.Contains(t, err.Error(), "cycle")
}

func TestExpandTemplates(t *testing.T) {
	vars := map[string]*ast.Var{
		"NAME": {Static: "world", Value: "world"},
	}

	result := ExpandTemplates("hello {{.NAME}}", vars)
	require.Equal(t, "hello world", result)
}

func TestPopulateBuiltins(t *testing.T) {
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

	PopulateBuiltins(tf)

	require.Contains(t, tf.Vars, "OS")
	require.Contains(t, tf.Vars, "ARCH")
	require.Contains(t, tf.Vars, "OSFAMILY")
	require.Contains(t, tf.Vars, "NUMCPU")
	require.Contains(t, tf.Vars, "ROOT_DIR")
	require.Contains(t, tf.Vars, "TASKFILE")
	require.Contains(t, tf.Vars, "TASKFILE_DIR")
	require.Contains(t, tf.Vars, "exeExt")
}
