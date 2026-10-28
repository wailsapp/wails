package exec

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
)

func TestConcurrentTaskCacheSavesKeepAllEntries(t *testing.T) {
	base := t.TempDir()
	cache, err := LoadTaskCache(base)
	if err != nil {
		t.Fatal(err)
	}
	const workerCount, tasksPerWorker = 64, 4
	start := make(chan struct{})
	errors := make(chan error, workerCount)
	var workers sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		workers.Go(func() {
			<-start
			for j := 0; j < tasksPerWorker; j++ {
				task := &ast.Task{Name: fmt.Sprintf("worker-%d-task-%d", i, j)}
				if err := cache.RecordTask(task, base); err != nil {
					errors <- err
					return
				}
			}
		})
	}
	close(start)
	workers.Wait()
	close(errors)
	for err := range errors {
		t.Error(err)
	}
	data, err := os.ReadFile(filepath.Join(base, cacheDirName, cacheFileName))
	if err != nil {
		t.Fatal(err)
	}
	if !json.Valid(data) {
		t.Fatal("concurrent saves produced invalid JSON")
	}
	reloaded, err := LoadTaskCache(base)
	if err != nil {
		t.Fatal(err)
	}
	if got := len(reloaded.Entries); got != workerCount*tasksPerWorker {
		t.Errorf("reloaded %d task entries, want %d", got, workerCount*tasksPerWorker)
	}
}
