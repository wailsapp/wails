package exec

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
)

func TestQuotedGoOutputCache(t *testing.T) {
	for _, command := range []string{`go build -o "app with spaces"`, `go build -o="app with spaces"`} {
		t.Run(command, func(t *testing.T) {
			base := t.TempDir()
			wantOutput := filepath.Join(base, "app with spaces")
			if err := os.WriteFile(wantOutput, nil, 0644); err != nil {
				t.Fatal(err)
			}
			cache, err := LoadTaskCache(base)
			if err != nil {
				t.Fatal(err)
			}
			task := &ast.Task{Name: "build"}
			executor := &Executor{Taskfile: &ast.Taskfile{}, Dir: base}
			sources, output := executor.goCmdInputs(goCmdBuild, task, command, base)
			if output != wantOutput {
				t.Errorf("cache output = %q, want %q", output, wantOutput)
			}
			if err := cache.RecordGoCmd(task.Name, command); err != nil {
				t.Fatal(err)
			}
			if !cache.ShouldSkipGoCmd(task.Name, command, sources, output) {
				t.Error("existing quoted output should allow a cache hit")
			}
			if err := os.Remove(wantOutput); err != nil {
				t.Fatal(err)
			}
			if cache.ShouldSkipGoCmd(task.Name, command, sources, output) {
				t.Error("deleted quoted output should force a cache miss")
			}
		})
	}
}
