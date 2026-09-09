package resolve

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
)

func TestBuildDAG(t *testing.T) {
	tf := &ast.Taskfile{
		Tasks: map[string]*ast.Task{
			"build": {
				Name: "build",
				Deps: []*ast.Dep{{Task: "deps:one"}, {Task: "deps:two"}},
			},
			"deps:one": {
				Name: "deps:one",
			},
			"deps:two": {
				Name: "deps:two",
			},
		},
	}

	dag, err := BuildDAG(tf, "build")
	require.NoError(t, err)
	require.Len(t, dag.Order, 3)

	oneIdx := -1
	twoIdx := -1
	buildIdx := -1
	for i, name := range dag.Order {
		switch name {
		case "deps:one":
			oneIdx = i
		case "deps:two":
			twoIdx = i
		case "build":
			buildIdx = i
		}
	}

	require.Greater(t, buildIdx, oneIdx)
	require.Greater(t, buildIdx, twoIdx)
}

func TestBuildDAGNotFound(t *testing.T) {
	tf := &ast.Taskfile{
		Tasks: map[string]*ast.Task{
			"build": {Name: "build"},
		},
	}

	_, err := BuildDAG(tf, "nonexistent")
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found")
}

func TestDetectCycle(t *testing.T) {
	tf := &ast.Taskfile{
		Tasks: map[string]*ast.Task{
			"a": {
				Name: "a",
				Deps: []*ast.Dep{{Task: "b"}},
			},
			"b": {
				Name: "b",
				Deps: []*ast.Dep{{Task: "a"}},
			},
		},
	}

	_, err := BuildDAG(tf, "a")
	require.Error(t, err)
	require.Contains(t, err.Error(), "cycle")
}

func TestMergeTask(t *testing.T) {
	base := &ast.Task{
		Name:    "build",
		Dir:     "/base",
		Summary: "base summary",
		Cmds: []*ast.Cmd{
			{Cmd: "echo base"},
		},
		Vars: map[string]*ast.Var{
			"A": {Static: "base"},
		},
		Env: map[string]string{
			"FOO": "base",
		},
	}

	override := &ast.Task{
		Dir: "/override",
		Cmds: []*ast.Cmd{
			{Cmd: "echo override"},
		},
		Vars: map[string]*ast.Var{
			"A": {Static: "override"},
			"B": {Static: "new"},
		},
		Env: map[string]string{
			"FOO": "override",
			"BAR": "new",
		},
	}

	result := MergeTask(base, override)

	require.Equal(t, "/override", result.Dir)
	// Local-wins: override cmds replace base cmds (they do not append).
	require.Len(t, result.Cmds, 1)
	require.Equal(t, "echo override", result.Cmds[0].Cmd)
	require.Equal(t, "override", result.Vars["A"].Static)
	require.Equal(t, "new", result.Vars["B"].Static)
	require.Equal(t, "override", result.Env["FOO"])
	require.Equal(t, "new", result.Env["BAR"])
}
