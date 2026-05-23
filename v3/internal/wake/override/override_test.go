package override

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
)

func TestOverrideNamed(t *testing.T) {
	ClearOverrides()

	Override("build", func(t *ast.Task) *ast.Task {
		t.Cmds = append(t.Cmds, &ast.Cmd{Cmd: "echo overridden"})
		return t
	})

	require.Contains(t, Named(), "build")

	task := &ast.Task{
		Name: "build",
		Cmds: []*ast.Cmd{{Cmd: "echo original"}},
	}
	result := Named()["build"](task)
	require.Len(t, result.Cmds, 2)
	require.Equal(t, "echo overridden", result.Cmds[1].Cmd)
}

func TestOverrideAll(t *testing.T) {
	ClearOverrides()

	OverrideAll(func(t *ast.Task) *ast.Task {
		t.Env = map[string]string{"GLOBAL": "true"}
		return t
	})

	require.Len(t, All(), 1)

	task := &ast.Task{Name: "any"}
	result := All()[0](task)
	require.Equal(t, "true", result.Env["GLOBAL"])
}

func TestLoadTaskfileOverride(t *testing.T) {
	ClearOverrides()

	dir := t.TempDir()
	ovPath := filepath.Join(dir, "Taskfile.local.yml")
	err := os.WriteFile(ovPath, []byte(`
version: "3"
tasks:
  build:
    cmds:
      - echo post-build
    env:
      CUSTOM: value
`), 0644)
	require.NoError(t, err)

	err = LoadTaskfileOverride(ovPath)
	require.NoError(t, err)
	require.Contains(t, Named(), "build")

	task := &ast.Task{
		Name: "build",
		Cmds: []*ast.Cmd{{Cmd: "echo build"}},
		Env:  map[string]string{"EXISTING": "yes"},
	}
	result := Named()["build"](task)
	require.Len(t, result.Cmds, 2)
	require.Equal(t, "echo build", result.Cmds[0].Cmd)
	require.Equal(t, "echo post-build", result.Cmds[1].Cmd)
	require.Equal(t, "value", result.Env["CUSTOM"])
	require.Equal(t, "yes", result.Env["EXISTING"])
}

func TestLoadTaskfileOverrideNotFound(t *testing.T) {
	ClearOverrides()
	err := LoadTaskfileOverride("/nonexistent/Taskfile.local.yml")
	require.Error(t, err)
}

func TestLoadTaskfileOverrides(t *testing.T) {
	ClearOverrides()

	dir := t.TempDir()

	localPath := filepath.Join(dir, "Taskfile.local.yml")
	err := os.WriteFile(localPath, []byte(`
version: "3"
tasks:
  build:
    cmds:
      - echo local-override
`), 0644)
	require.NoError(t, err)

	err = LoadTaskfileOverrides(dir)
	require.NoError(t, err)
	require.Contains(t, Named(), "build")
}

func TestLoadTaskfileOverridesNoneExist(t *testing.T) {
	ClearOverrides()
	dir := t.TempDir()
	err := LoadTaskfileOverrides(dir)
	require.NoError(t, err)
	require.Empty(t, Named())
}
