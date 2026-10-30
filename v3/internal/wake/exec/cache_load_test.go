package exec

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
)

func TestLoadedCacheCanRecordTasks(t *testing.T) {
	for _, tt := range []struct {
		name, content string
		wantEntries   int
	}{
		{"null entries", `{"entries":null}`, 1},
		{"missing entries", `{}`, 1},
		{"empty entries", `{"entries":{}}`, 1},
		{"invalid JSON", `{`, 1},
		{"existing entry", `{"entries":{"existing":{"hash":"prior"}}}`, 2},
	} {
		t.Run(tt.name, func(t *testing.T) {
			base := t.TempDir()
			dir := filepath.Join(base, cacheDirName)
			if err := os.Mkdir(dir, 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(dir, cacheFileName), []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}
			cache, err := LoadTaskCache(base)
			if err != nil {
				t.Fatal(err)
			}
			if err := cache.RecordTask(&ast.Task{Name: "build"}, base); err != nil {
				t.Fatal(err)
			}
			reloaded, err := LoadTaskCache(base)
			if err != nil {
				t.Fatal(err)
			}
			if len(reloaded.Entries) != tt.wantEntries || reloaded.Entries["build"].Hash == "" {
				t.Errorf("saved cache entries = %+v, want %d including build", reloaded.Entries, tt.wantEntries)
			}
		})
	}
}
