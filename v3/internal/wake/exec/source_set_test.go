package exec

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
)

func TestGoCacheSourceSetChanges(t *testing.T) {
	for _, scenario := range []string{"add", "delete"} {
		t.Run(scenario, func(t *testing.T) {
			root := t.TempDir()
			candidate := filepath.Join(root, "changed.go")
			if scenario == "delete" {
				writeOldGoSource(t, candidate)
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
				t.Fatal("unchanged input set should allow a cache hit")
			}
			if scenario == "add" {
				writeOldGoSource(t, candidate)
			} else if err := os.Remove(candidate); err != nil {
				t.Fatal(err)
			}
			sources, output = executor.goCmdInputs(goCmdBuild, task, "go build", root)
			if cache.ShouldSkipGoCmd(task.Name, "go build", sources, output) {
				t.Errorf("source %s did not invalidate the cache", scenario)
			}
		})
	}
}

func writeOldGoSource(t *testing.T, path string) {
	t.Helper()
	if err := os.WriteFile(path, []byte("package fixture\n"), 0644); err != nil {
		t.Fatal(err)
	}
	old := time.Now().Add(-time.Hour)
	if err := os.Chtimes(path, old, old); err != nil {
		t.Fatal(err)
	}
}

func TestGoCacheSourceSetOrderAndPersistence(t *testing.T) {
	root := t.TempDir()
	a, b := filepath.Join(root, "a.go"), filepath.Join(root, "b.go")
	writeOldGoSource(t, a)
	writeOldGoSource(t, b)
	cache, err := LoadTaskCache(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := cache.RecordGoCmd("build", "go build", []string{a, b}); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadTaskCache(root)
	if err != nil {
		t.Fatal(err)
	}
	if !loaded.ShouldSkipGoCmd("build", "go build", []string{b, a, b}, "") {
		t.Error("reordered or repeated input paths should retain the persisted cache hit")
	}
}

func TestLegacyGoCacheEntryRebuildsOnce(t *testing.T) {
	cache := &TaskCache{
		Dir: t.TempDir(),
		Entries: map[string]TaskCacheEntry{
			"build" + goCacheSuffix: {Hash: goCmdHash("go build"), LastRun: time.Now()},
		},
	}
	if cache.ShouldSkipGoCmd("build", "go build", nil, "") {
		t.Error("legacy command-only cache entry should not be trusted")
	}
	if err := cache.RecordGoCmd("build", "go build", nil); err != nil {
		t.Fatal(err)
	}
	if !cache.ShouldSkipGoCmd("build", "go build", nil, "") {
		t.Error("refreshed entry should allow a cache hit")
	}
}
