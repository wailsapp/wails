package exec

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
	"github.com/wailsapp/wails/v3/internal/wake/parse"
	"github.com/wailsapp/wails/v3/internal/wake/platform"
)

type runCache struct {
	mu       sync.Mutex
	lastRuns map[string]time.Time
}

var cache = &runCache{
	lastRuns: make(map[string]time.Time),
}

func checkPreconditions(task *ast.Task, vars map[string]*ast.Var) error {
	for _, pc := range task.Precondition {
		// Precondition `sh:` strings are Go templates just like cmds. They must
		// be expanded against the task's resolved vars before running; otherwise
		// a guard like `{{if eq .OBFUSCATED "true"}}command -v garble{{else}}true{{end}}`
		// reaches the shell verbatim, fails to parse as a command, and surfaces
		// the precondition's `msg:` as a spurious failure on every build.
		sh := parse.ExpandTemplates(pc.Sh, vars)
		if strings.TrimSpace(sh) == "" {
			continue
		}
		c := platform.ShellCommand(sh)
		if err := c.Run(); err != nil {
			msg := parse.ExpandTemplates(pc.Msg, vars)
			if msg == "" {
				msg = fmt.Sprintf("precondition failed: %q", sh)
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
			c := platform.ShellCommand(cmd)
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

	// Read lastRuns under the mutex — RecordRun writes to this map and runs
	// concurrently when Parallel=true. The map snapshot was previously read
	// raw, which both produced inconsistent up-to-date decisions and could
	// panic on Go's "concurrent map read and map write" detector.
	cache.mu.Lock()
	lastRun := cache.lastRuns[task.Name]
	cache.mu.Unlock()

	for _, pattern := range task.Sources {
		// globMatches understands Taskfile-style "**" recursive globs;
		// stdlib filepath.Glob does not, and silently returning no matches
		// would make recursively-watched directories spuriously up-to-date.
		matches := globMatches(baseDir, pattern)
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

// globExists reports whether at least one file matches pattern within
// baseDir. Uses globMatches so "**" recursive patterns behave consistently
// with the rest of the cache layer.
func globExists(baseDir, pattern string) bool {
	return len(globMatches(baseDir, pattern)) > 0
}

func RecordRun(taskName string) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.lastRuns[taskName] = time.Now()
}
