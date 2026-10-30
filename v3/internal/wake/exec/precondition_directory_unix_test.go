//go:build !windows

package exec

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
)

func TestPreconditionsUseTaskDirectory(t *testing.T) {
	for _, directory := range []string{"", "work", "{{.WORK}}"} {
		t.Run("dir="+directory, func(t *testing.T) {
			base := t.TempDir()
			work := base
			if directory != "" {
				work = filepath.Join(base, "work")
			}
			if err := os.MkdirAll(work, 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(work, "ready"), nil, 0644); err != nil {
				t.Fatal(err)
			}
			t.Chdir(t.TempDir())
			task := &ast.Task{
				Name: "guarded", Dir: directory,
				Precondition: []*ast.Precondition{{Sh: "test -f ready", Msg: "ready file missing"}},
				Cmds:         []*ast.Cmd{{Cmd: "echo done > output.txt"}},
			}
			executor := &Executor{
				Taskfile: &ast.Taskfile{
					Tasks: map[string]*ast.Task{"guarded": task},
					Vars:  map[string]*ast.Var{"WORK": {Value: "work"}},
				},
				Dir: base,
			}
			if err := executor.Execute(context.Background(), "guarded"); err != nil {
				t.Fatal(err)
			}
			if data, err := os.ReadFile(filepath.Join(work, "output.txt")); err != nil || string(data) != "done\n" {
				t.Errorf("guarded task output = %q, %v", data, err)
			}
		})
	}
}
