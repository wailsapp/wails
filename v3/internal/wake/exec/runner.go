package exec

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
)

type runCache struct {
	mu       sync.Mutex
	lastRuns map[string]time.Time
}

var cache = &runCache{
	lastRuns: make(map[string]time.Time),
}

func checkPreconditions(task *ast.Task) error {
	for _, pc := range task.Precondition {
		if pc.Sh == "" {
			continue
		}
		c := exec.Command("sh", "-c", pc.Sh)
		if err := c.Run(); err != nil {
			msg := pc.Msg
			if msg == "" {
				msg = fmt.Sprintf("precondition failed: %q", pc.Sh)
			}
			return fmt.Errorf("wake: %s", msg)
		}
	}
	return nil
}

func isUpToDate(task *ast.Task, baseDir string) bool {
	if len(task.Sources) == 0 && len(task.Generates) == 0 {
		return false
	}

	if len(task.Status) > 0 {
		for _, cmd := range task.Status {
			c := exec.Command("sh", "-c", cmd)
			if err := c.Run(); err != nil {
				return false
			}
		}
		return true
	}

	for _, pattern := range task.Generates {
		if !globExists(baseDir, pattern) {
			return false
		}
	}

	lastRun := cache.lastRuns[task.Name]
	for _, pattern := range task.Sources {
		matches, err := filepath.Glob(filepath.Join(baseDir, pattern))
		if err != nil {
			continue
		}
		for _, file := range matches {
			info, err := os.Stat(file)
			if err != nil {
				continue
			}
			if info.ModTime().After(lastRun) {
				return false
			}
		}
	}

	return true
}

func globExists(baseDir, pattern string) bool {
	matches, err := filepath.Glob(filepath.Join(baseDir, pattern))
	if err != nil {
		return false
	}
	return len(matches) > 0
}

func RecordRun(taskName string) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.lastRuns[taskName] = time.Now()
}
