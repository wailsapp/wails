//go:build !windows

package exec

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
)

func TestExecutorDefaultsToBaseDirectory(t *testing.T) {
	for _, directory := range []string{"", "nested"} {
		t.Run("dir="+directory, func(t *testing.T) {
			base := t.TempDir()
			caller := t.TempDir()
			wantDir := filepath.Join(base, directory)
			if err := os.MkdirAll(wantDir, 0755); err != nil {
				t.Fatal(err)
			}
			t.Chdir(caller)
			task := &ast.Task{Name: "write", Dir: directory, Cmds: []*ast.Cmd{{Cmd: "echo expected > marker.txt"}}}
			executor := &Executor{
				Taskfile: &ast.Taskfile{Tasks: map[string]*ast.Task{"write": task}},
				Dir:      base,
			}
			if err := executor.Execute(context.Background(), "write"); err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(filepath.Join(wantDir, "marker.txt"))
			if err != nil || string(data) != "expected\n" {
				t.Errorf("task output = %q, %v; want expected output in %s", data, err, wantDir)
			}
			if _, err := os.Stat(filepath.Join(caller, "marker.txt")); !os.IsNotExist(err) {
				t.Errorf("task wrote to the caller directory: %v", err)
			}
		})
	}
}
