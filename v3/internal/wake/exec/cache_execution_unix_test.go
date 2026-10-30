//go:build !windows

package exec

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
)

func TestCacheMissExecutesChangedCommand(t *testing.T) {
	base := t.TempDir()
	cache, err := LoadTaskCache(base)
	if err != nil {
		t.Fatal(err)
	}
	task := &ast.Task{
		Name: "build", Dir: base, Generates: []string{"output.txt"},
		Cmds: []*ast.Cmd{{Cmd: "echo old > output.txt"}},
	}
	taskfile := &ast.Taskfile{Tasks: map[string]*ast.Task{"build": task}}
	run := func() {
		t.Helper()
		executor := &Executor{Taskfile: taskfile, Dir: base, Cache: cache}
		if err := executor.Execute(context.Background(), "build"); err != nil {
			t.Fatal(err)
		}
	}
	output := filepath.Join(base, "output.txt")
	run()
	task.Cmds[0].Cmd = "echo new > output.txt"
	run()
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "new\n" {
		t.Errorf("changed command produced %q, want new output", data)
	}
	if err := os.WriteFile(output, []byte("cache hit control"), 0644); err != nil {
		t.Fatal(err)
	}
	run()
	data, err = os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "cache hit control" {
		t.Errorf("unchanged cached task unexpectedly executed: %q", data)
	}
}
