//go:build !windows

package exec

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
)

func TestGoCommandChainsAreNotImplicitlyCached(t *testing.T) {
	for _, tt := range []struct {
		command string
		runs    int
	}{
		{"go build", 1},
		{"go build && echo chained >> side-effect.txt", 2},
	} {
		t.Run(tt.command, func(t *testing.T) {
			base := t.TempDir()
			if err := os.WriteFile(filepath.Join(base, "go"), []byte("#!/bin/sh\necho invoked >> invocations.txt\n"), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(base, "go.mod"), []byte("module example.invalid/cachefixture\ngo 1.25\n"), 0644); err != nil {
				t.Fatal(err)
			}
			t.Setenv("PATH", base+string(os.PathListSeparator)+os.Getenv("PATH"))
			cache, err := LoadTaskCache(base)
			if err != nil {
				t.Fatal(err)
			}
			task := &ast.Task{Name: "build", Cmds: []*ast.Cmd{{Cmd: tt.command}}}
			taskfile := &ast.Taskfile{Tasks: map[string]*ast.Task{"build": task}}
			for i := 0; i < 2; i++ {
				executor := &Executor{Taskfile: taskfile, Dir: base, Cache: cache}
				if err := executor.Execute(context.Background(), "build"); err != nil {
					t.Fatal(err)
				}
			}
			data, err := os.ReadFile(filepath.Join(base, "invocations.txt"))
			if err != nil || strings.Count(string(data), "invoked\n") != tt.runs {
				t.Errorf("Go invocations = %q, %v; want %d runs", data, err, tt.runs)
			}
			if tt.runs == 2 {
				data, err := os.ReadFile(filepath.Join(base, "side-effect.txt"))
				if err != nil || string(data) != "chained\nchained\n" {
					t.Errorf("chain side effects = %q, %v; want two executions", data, err)
				}
			}
		})
	}
}
