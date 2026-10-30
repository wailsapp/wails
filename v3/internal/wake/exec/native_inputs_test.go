package exec

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
)

func TestGoBuildCacheTracksNativeSources(t *testing.T) {
	for _, tc := range []struct {
		name     string
		wantSkip bool
	}{
		{"native.c", false}, {"native.cc", false}, {"native.cpp", false}, {"native.cxx", false},
		{"native.m", false}, {"native.h", false}, {"native.hh", false}, {"native.hpp", false}, {"native.hxx", false},
		{"native.f", false}, {"native.F", false}, {"native.for", false}, {"native.f90", false},
		{"native.s", false}, {"native.S", false}, {"native.sx", false},
		{"native.swig", false}, {"native.swigcxx", false}, {"native.syso", false},
		{"notes.txt", true}, {"main_test.go", true}, {"node_modules/ignored.c", true}, {"bin/ignored.s", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			// Only input discovery and timestamps are under test; no compiler is invoked.
			path := filepath.Join(root, tc.name)
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte("fixture"), 0644); err != nil {
				t.Fatal(err)
			}
			cache, err := LoadTaskCache(root)
			if err != nil {
				t.Fatal(err)
			}
			task := &ast.Task{Name: "build"}
			executor := &Executor{Taskfile: &ast.Taskfile{}, Dir: root}
			sources, output := executor.goCmdInputs(goCmdBuild, task, "go build", root)
			if err := cache.RecordGoCmd(task.Name, "go build", sources); err != nil {
				t.Fatal(err)
			}
			if !cache.ShouldSkipGoCmd(task.Name, "go build", sources, output) {
				t.Fatal("unchanged inputs should allow a cache hit")
			}
			changed := time.Now().Add(time.Hour)
			if err := os.Chtimes(path, changed, changed); err != nil {
				t.Fatal(err)
			}
			sources, output = executor.goCmdInputs(goCmdBuild, task, "go build", root)
			if got := cache.ShouldSkipGoCmd(task.Name, "go build", sources, output); got != tc.wantSkip {
				t.Errorf("cache skip after changing %s = %t, want %t", tc.name, got, tc.wantSkip)
			}
		})
	}
}
