//go:build !windows

package exec

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
)

func TestStatusChecksTaskDirectory(t *testing.T) {
	for _, tt := range []struct {
		name      string
		absolute  bool
		taskReady bool
	}{
		{"relative directory", false, true},
		{"absolute directory", true, true},
		{"caller marker does not skip task", false, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			base := t.TempDir()
			work := filepath.Join(base, "work")
			if err := os.Mkdir(work, 0755); err != nil {
				t.Fatal(err)
			}
			markerDir := base
			if tt.taskReady {
				markerDir = work
			}
			if err := os.WriteFile(filepath.Join(markerDir, "ready"), nil, 0644); err != nil {
				t.Fatal(err)
			}
			t.Chdir(base)
			task := &ast.Task{Dir: "work", Status: []string{"test -f ready"}}
			if tt.absolute {
				task.Dir = work
			}
			cache := &TaskCache{}
			if got := cache.ShouldSkip(task, base); got != tt.taskReady {
				t.Errorf("ShouldSkip = %v, want %v", got, tt.taskReady)
			}
		})
	}
}
