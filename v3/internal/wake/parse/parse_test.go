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

func TestResolveIncludesInjectsTopLevelVarsIntoClonedTasks(t *testing.T) {
	// Regression test: included file's top-level vars (e.g. CROSS_IMAGE) must be
	// available in each cloned task's Vars map so preconditions and commands that
	// reference {{.CROSS_IMAGE}} expand correctly at execution time.
	dir := t.TempDir()

	subDir := filepath.Join(dir, "sub")
	require.NoError(t, os.MkdirAll(subDir, 0755))

	require.NoError(t, os.WriteFile(filepath.Join(subDir, "Taskfile.yml"), []byte(`
version: "3"
vars:
  CROSS_IMAGE: wails-cross
  STATIC_VAR: from-include

tasks:
  build:
    cmds:
      - echo build
  build:docker:
    vars:
      OWN_VAR: own-value
    cmds:
      - docker run {{.CROSS_IMAGE}}
`), 0644))

	require.NoError(t, os.WriteFile(filepath.Join(dir, "Taskfile.yml"), []byte(`
version: "3"
includes:
  linux:
    taskfile: ./sub/Taskfile.yml
tasks:
  main:
    cmds:
      - echo main
`), 0644))

	tf, err := Parse(filepath.Join(dir, "Taskfile.yml"))
	require.NoError(t, err)
	require.NoError(t, ResolveIncludes(tf))

	// Both cloned tasks must carry the included file's top-level vars.
	buildTask, ok := tf.Tasks["linux:build"]
	require.True(t, ok, "linux:build must exist after include resolution")
	require.Contains(t, buildTask.Vars, "CROSS_IMAGE", "top-level included var must propagate to linux:build")
	require.Equal(t, "wails-cross", buildTask.Vars["CROSS_IMAGE"].Static)

	dockerTask, ok := tf.Tasks["linux:build:docker"]
	require.True(t, ok, "linux:build:docker must exist after include resolution")
	require.Contains(t, dockerTask.Vars, "CROSS_IMAGE", "top-level included var must propagate to linux:build:docker")
	require.Equal(t, "wails-cross", dockerTask.Vars["CROSS_IMAGE"].Static)

	// Task-level vars must take precedence over included file's top-level vars.
	require.Contains(t, dockerTask.Vars, "OWN_VAR", "task-level vars must be retained")
	require.Equal(t, "own-value", dockerTask.Vars["OWN_VAR"].Static)

	// The injected vars must be deep-copied, not shared by pointer across cloned
	// tasks: later in-place mutation (ResolveVars/ResolveAllVarShells) of one
	// task's vars must not leak into the others that share the included file.
	require.NotSame(t, buildTask.Vars["CROSS_IMAGE"], dockerTask.Vars["CROSS_IMAGE"],
		"injected included vars must not be shared by pointer across cloned tasks")
	buildTask.Vars["CROSS_IMAGE"].Static = "mutated"
	require.Equal(t, "wails-cross", dockerTask.Vars["CROSS_IMAGE"].Static,
		"mutating one cloned task's var must not affect another's")
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
