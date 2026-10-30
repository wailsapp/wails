package exec

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
)

func TestTaskDirectoryInvalidatesCache(t *testing.T) {
	for _, absolute := range []bool{false, true} {
		name := "relative"
		if absolute {
			name = "absolute"
		}
		t.Run(name, func(t *testing.T) {
			base := t.TempDir()
			directories := []string{"first", "second"}
			for i, directory := range directories {
				path := filepath.Join(base, directory)
				if err := os.Mkdir(path, 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(path, "output.txt"), []byte("existing output"), 0644); err != nil {
					t.Fatal(err)
				}
				if absolute {
					directories[i] = path
				}
			}
			cache, err := LoadTaskCache(base)
			if err != nil {
				t.Fatal(err)
			}
			task := &ast.Task{Name: "build", Dir: directories[0], Generates: []string{"output.txt"}}
			if err := cache.RecordTask(task, base); err != nil {
				t.Fatal(err)
			}
			if !cache.ShouldSkip(task, base) {
				t.Fatal("unchanged task should use its recorded cache entry")
			}
			task.Dir = directories[1]
			if cache.ShouldSkip(task, base) {
				t.Error("changed task directory incorrectly reused the old cache entry")
			}
			if err := cache.RecordTask(task, base); err != nil {
				t.Fatal(err)
			}
			if !cache.ShouldSkip(task, base) {
				t.Error("task should be cached after recording its new directory")
			}
		})
	}
}
