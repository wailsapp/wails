package wake

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake/ast"
)

func TestParseSimpleTaskfile(t *testing.T) {
	dir := t.TempDir()
	tfPath := filepath.Join(dir, "Taskfile.yml")
	err := os.WriteFile(tfPath, []byte(`
version: "3"
tasks:
  hello:
    cmds:
      - echo hello
`), 0644)
	require.NoError(t, err)

	tf, err := Parse(tfPath)
	require.NoError(t, err)
	require.NotNil(t, tf)
	require.Equal(t, "3", tf.Version)
	require.Contains(t, tf.Tasks, "hello")
	require.Len(t, tf.Tasks["hello"].Cmds, 1)
	require.Equal(t, "echo hello", tf.Tasks["hello"].Cmds[0].Cmd)
}

func TestParseWithIncludes(t *testing.T) {
	dir := t.TempDir()

	mainTf := filepath.Join(dir, "Taskfile.yml")
	err := os.WriteFile(mainTf, []byte(`
version: "3"
includes:
  common:
    taskfile: ./common/Taskfile.yml
    dir: ./common
tasks:
  build:
    deps:
      - task: common:setup
    cmds:
      - echo build
`), 0644)
	require.NoError(t, err)

	commonDir := filepath.Join(dir, "common")
	err = os.MkdirAll(commonDir, 0755)
	require.NoError(t, err)

	commonTf := filepath.Join(commonDir, "Taskfile.yml")
	err = os.WriteFile(commonTf, []byte(`
version: "3"
tasks:
  setup:
    cmds:
      - echo setup
`), 0644)
	require.NoError(t, err)

	tf, err := Parse(mainTf)
	require.NoError(t, err)
	require.NotNil(t, tf)
	require.Contains(t, tf.Tasks, "build")
	require.Contains(t, tf.Tasks, "common:setup")
}

func TestParseWithVars(t *testing.T) {
	dir := t.TempDir()
	tfPath := filepath.Join(dir, "Taskfile.yml")
	err := os.WriteFile(tfPath, []byte(`
version: "3"
vars:
  GREETING: hello world
tasks:
  greet:
    cmds:
      - echo {{.GREETING}}
`), 0644)
	require.NoError(t, err)

	tf, err := Parse(tfPath)
	require.NoError(t, err)
	require.NotNil(t, tf)
	require.Contains(t, tf.Vars, "GREETING")
	require.Equal(t, "hello world", tf.Vars["GREETING"].Static)
}

func TestParseWithPlatforms(t *testing.T) {
	dir := t.TempDir()
	tfPath := filepath.Join(dir, "Taskfile.yml")
	err := os.WriteFile(tfPath, []byte(`
version: "3"
tasks:
  darwin-only:
    platforms:
      - darwin
    cmds:
      - echo darwin
  linux-only:
    platforms:
      - linux
    cmds:
      - echo linux
`), 0644)
	require.NoError(t, err)

	tf, err := Parse(tfPath)
	require.NoError(t, err)
	require.NotNil(t, tf)
	require.Contains(t, tf.Tasks, "darwin-only")
	require.Contains(t, tf.Tasks, "linux-only")
	require.Equal(t, []string{"darwin"}, tf.Tasks["darwin-only"].Platforms)
}

func TestParseWithDeps(t *testing.T) {
	dir := t.TempDir()
	tfPath := filepath.Join(dir, "Taskfile.yml")
	err := os.WriteFile(tfPath, []byte(`
version: "3"
tasks:
  build:
    deps:
      - task: deps:one
      - task: deps:two
    cmds:
      - echo build
  deps:one:
    cmds:
      - echo one
  deps:two:
    cmds:
      - echo two
`), 0644)
	require.NoError(t, err)

	tf, err := Parse(tfPath)
	require.NoError(t, err)
	require.NotNil(t, tf)
	require.Len(t, tf.Tasks["build"].Deps, 2)
	require.Equal(t, "deps:one", tf.Tasks["build"].Deps[0].Task)
	require.Equal(t, "deps:two", tf.Tasks["build"].Deps[1].Task)
}

func TestParseWithForLoop(t *testing.T) {
	dir := t.TempDir()
	tfPath := filepath.Join(dir, "Taskfile.yml")
	err := os.WriteFile(tfPath, []byte(`
version: "3"
vars:
  ITEMS: a b c
tasks:
  loop:
    cmds:
      - for: { var: ITEMS }
        task: process
  process:
    cmds:
      - echo process
`), 0644)
	require.NoError(t, err)

	tf, err := Parse(tfPath)
	require.NoError(t, err)
	require.NotNil(t, tf)
	require.NotNil(t, tf.Tasks["loop"].Cmds[0].For)
	require.Equal(t, "ITEMS", tf.Tasks["loop"].Cmds[0].For.Var)
	require.Equal(t, "process", tf.Tasks["loop"].Cmds[0].Task)
}

func TestParseWithPreconditions(t *testing.T) {
	dir := t.TempDir()
	tfPath := filepath.Join(dir, "Taskfile.yml")
	err := os.WriteFile(tfPath, []byte(`
version: "3"
tasks:
  check:
    preconditions:
      - sh: test -d .
        msg: current directory must exist
    cmds:
      - echo ok
`), 0644)
	require.NoError(t, err)

	tf, err := Parse(tfPath)
	require.NoError(t, err)
	require.NotNil(t, tf)
	require.Len(t, tf.Tasks["check"].Precondition, 1)
	require.Equal(t, "test -d .", tf.Tasks["check"].Precondition[0].Sh)
	require.Equal(t, "current directory must exist", tf.Tasks["check"].Precondition[0].Msg)
}

func TestParseWithEnv(t *testing.T) {
	dir := t.TempDir()
	tfPath := filepath.Join(dir, "Taskfile.yml")
	err := os.WriteFile(tfPath, []byte(`
version: "3"
tasks:
  with-env:
    env:
      FOO: bar
      BAZ: qux
    cmds:
      - echo $FOO
`), 0644)
	require.NoError(t, err)

	tf, err := Parse(tfPath)
	require.NoError(t, err)
	require.NotNil(t, tf)
	require.Equal(t, "bar", tf.Tasks["with-env"].Env["FOO"])
	require.Equal(t, "qux", tf.Tasks["with-env"].Env["BAZ"])
}

func TestParseRejectsWrongVersion(t *testing.T) {
	dir := t.TempDir()
	tfPath := filepath.Join(dir, "Taskfile.yml")
	err := os.WriteFile(tfPath, []byte(`
version: "2"
tasks:
  hello:
    cmds:
      - echo hello
`), 0644)
	require.NoError(t, err)

	_, err = Parse(tfPath)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported Taskfile version")
}

func TestDiscoverAndParse(t *testing.T) {
	dir := t.TempDir()

	tfPath := filepath.Join(dir, "Taskfile.yml")
	err := os.WriteFile(tfPath, []byte(`
version: "3"
tasks:
  hello:
    cmds:
      - echo hello
`), 0644)
	require.NoError(t, err)

	tf, err := discoverAndParse(dir)
	require.NoError(t, err)
	require.NotNil(t, tf)
	require.Contains(t, tf.Tasks, "hello")
}

func TestResolveInNamespacePrefersLocal(t *testing.T) {
	tf := &ast.Taskfile{
		Tasks: map[string]*ast.Task{
			"build":          {Name: "build"},
			"darwin:build":   {Name: "darwin:build"},
			"darwin:package": {Name: "darwin:package"},
			"darwin:common:go:mod:tidy": {Name: "darwin:common:go:mod:tidy"},
		},
		Includes: map[string]*ast.Include{"common": {}},
	}

	// A bare dep `build` inside the darwin namespace must bind to darwin:build,
	// not the same-named top-level wrapper (the prod-packaging dep-var bug).
	if got, ok := resolveInNamespace(tf, "darwin:package", "build"); !ok || got != "darwin:build" {
		t.Errorf("resolveInNamespace(darwin:package, build) = %q,%v; want darwin:build,true", got, ok)
	}

	// common:-qualified dep inside a namespace resolves to the namespaced copy.
	if got, ok := resolveInNamespace(tf, "darwin:build", "common:go:mod:tidy"); !ok || got != "darwin:common:go:mod:tidy" {
		t.Errorf("resolveInNamespace(darwin:build, common:go:mod:tidy) = %q,%v; want darwin:common:go:mod:tidy,true", got, ok)
	}

	// No namespace context: bare top-level name resolves to itself.
	if got, ok := resolveInNamespace(tf, "run", "build"); !ok || got != "build" {
		t.Errorf("resolveInNamespace(run, build) = %q,%v; want build,true", got, ok)
	}
}

func TestUseWakeEnvVar(t *testing.T) {
	os.Setenv("WAILS_USE_WAKE", "true")
	defer os.Unsetenv("WAILS_USE_WAKE")
	require.True(t, useWake())

	os.Setenv("WAILS_USE_WAKE", "false")
	require.False(t, useWake())

	os.Unsetenv("WAILS_USE_WAKE")
	require.False(t, useWake())
}
