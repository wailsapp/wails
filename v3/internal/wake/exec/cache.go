package exec

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
	"github.com/wailsapp/wails/v3/internal/wake/platform"
)

const cacheDirName = ".wake"
const cacheFileName = "cache.json"

type TaskCacheEntry struct {
	Hash    string    `json:"hash"`
	LastRun time.Time `json:"last_run"`
}

// TaskCache is shared across parallel workers: ShouldSkip / ShouldSkipGoCmd
// read Entries while RecordTask / RecordGoCmd write to it, and Save serialises
// the whole map. All access goes through mu.
type TaskCache struct {
	mu      sync.RWMutex
	Dir     string                    `json:"-"`
	Entries map[string]TaskCacheEntry `json:"entries"`
}

func LoadTaskCache(baseDir string) (*TaskCache, error) {
	cachePath := filepath.Join(baseDir, cacheDirName, cacheFileName)
	tc := &TaskCache{
		Dir:     baseDir,
		Entries: make(map[string]TaskCacheEntry),
	}

	data, err := os.ReadFile(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return tc, nil
		}
		return nil, fmt.Errorf("wake: read cache: %w", err)
	}

	if err := json.Unmarshal(data, tc); err != nil {
		return tc, nil
	}

	return tc, nil
}

func (tc *TaskCache) Save() error {
	tc.mu.RLock()
	dir := filepath.Join(tc.Dir, cacheDirName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		tc.mu.RUnlock()
		return fmt.Errorf("wake: create cache dir: %w", err)
	}
	path := filepath.Join(dir, cacheFileName)
	data, err := json.MarshalIndent(tc, "", "  ")
	tc.mu.RUnlock()
	if err != nil {
		return fmt.Errorf("wake: marshal cache: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

func ComputeTaskHash(task *ast.Task, baseDir string) (string, error) {
	h := sha256.New()

	h.Write([]byte(task.Name))
	h.Write([]byte("\n"))

	for _, cmd := range task.Cmds {
		h.Write([]byte(cmd.Cmd))
		h.Write([]byte(cmd.Task))
		h.Write([]byte("\n"))
	}

	sortedKeys := make([]string, 0, len(task.Env))
	for k := range task.Env {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	for _, k := range sortedKeys {
		h.Write([]byte(k + "=" + task.Env[k] + "\n"))
	}

	sortedSources := make([]string, len(task.Sources))
	copy(sortedSources, task.Sources)
	sort.Strings(sortedSources)
	for _, pattern := range sortedSources {
		h.Write([]byte(pattern + "\n"))
	}

	sortedGenerates := make([]string, len(task.Generates))
	copy(sortedGenerates, task.Generates)
	sort.Strings(sortedGenerates)
	for _, pattern := range sortedGenerates {
		h.Write([]byte(pattern + "\n"))
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func (tc *TaskCache) ShouldSkip(task *ast.Task, baseDir string) bool {
	if len(task.Sources) == 0 && len(task.Generates) == 0 && len(task.Status) == 0 {
		return false
	}

	taskDir := baseDir
	if task.Dir != "" {
		if filepath.IsAbs(task.Dir) {
			taskDir = task.Dir
		} else {
			taskDir = filepath.Join(baseDir, task.Dir)
		}
	}

	if len(task.Status) > 0 {
		for _, cmd := range task.Status {
			// platform.ShellCommand picks `sh -c` on Unix and `cmd /C` on
			// Windows. The previous unconditional `sh -c` here broke status-
			// based skipping on stock Windows installs.
			c := platform.ShellCommand(cmd)
			if err := c.Run(); err != nil {
				return false
			}
		}
		return true
	}

	for _, pattern := range task.Generates {
		if !globExists(taskDir, pattern) {
			return false
		}
	}

	hash, err := ComputeTaskHash(task, baseDir)
	if err != nil {
		return false
	}

	tc.mu.RLock()
	entry, ok := tc.Entries[task.Name]
	tc.mu.RUnlock()
	if !ok {
		return false
	}

	if entry.Hash != hash {
		return false
	}

	sources, excludes := splitSources(task.Sources)
	for _, pattern := range sources {
		files := globMatches(taskDir, pattern)
		files = applyExcludes(files, taskDir, excludes)
		for _, file := range files {
			info, err := os.Stat(file)
			if err != nil {
				continue
			}
			if info.ModTime().After(entry.LastRun) {
				return false
			}
		}
	}

	return true
}

func splitSources(sources []string) ([]string, []string) {
	var includes, excludes []string
	for _, s := range sources {
		if strings.HasPrefix(s, "exclude:") || strings.HasPrefix(s, "exclude ") {
			pattern := strings.TrimSpace(strings.TrimPrefix(s, "exclude:"))
			pattern = strings.TrimSpace(strings.TrimPrefix(pattern, "exclude "))
			if pattern != "" {
				excludes = append(excludes, pattern)
			}
		} else {
			includes = append(includes, s)
		}
	}
	return includes, excludes
}

func applyExcludes(files []string, baseDir string, excludes []string) []string {
	if len(excludes) == 0 {
		return files
	}

	var result []string
	for _, file := range files {
		rel, err := filepath.Rel(baseDir, file)
		if err != nil {
			rel = file
		}
		// Taskfile glob patterns are slash-separated. On Windows
		// filepath.Rel produces backslash paths, which then fail to
		// match a slash-based pattern. Normalise both sides to slashes
		// before matching so exclude rules behave the same on every OS.
		rel = filepath.ToSlash(rel)
		excluded := false
		for _, ex := range excludes {
			if matchesGlob(rel, filepath.ToSlash(ex)) {
				excluded = true
				break
			}
		}
		if !excluded {
			result = append(result, file)
		}
	}
	return result
}

func matchesGlob(path, pattern string) bool {
	if strings.Contains(pattern, "**") {
		return recursiveMatch(path, pattern)
	}
	matched, err := filepath.Match(pattern, path)
	if err != nil {
		return false
	}
	return matched
}

func recursiveMatch(path, pattern string) bool {
	parts := strings.SplitN(pattern, "**", 2)
	prefix := parts[0]
	suffix := ""
	if len(parts) > 1 {
		suffix = strings.TrimLeft(parts[1], "/")
	}

	// The portion before `**` is a literal path prefix the file must live under.
	if prefix != "" && !strings.HasPrefix(path, prefix) {
		return false
	}

	if suffix == "" || suffix == "*" {
		return true
	}

	if matched, _ := filepath.Match(suffix, filepath.Base(path)); matched {
		return true
	}
	return strings.HasSuffix(path, suffix)
}

func (tc *TaskCache) RecordTask(task *ast.Task, baseDir string) error {
	hash, err := ComputeTaskHash(task, baseDir)
	if err != nil {
		return err
	}

	tc.mu.Lock()
	tc.Entries[task.Name] = TaskCacheEntry{
		Hash:    hash,
		LastRun: time.Now(),
	}
	tc.mu.Unlock()

	return tc.Save()
}

func globMatches(baseDir, pattern string) []string {
	if strings.Contains(pattern, "**") {
		return recursiveGlob(baseDir, pattern)
	}
	matches, err := filepath.Glob(filepath.Join(baseDir, pattern))
	if err != nil {
		return nil
	}
	return matches
}

func recursiveGlob(baseDir, pattern string) []string {
	var results []string
	parts := strings.SplitN(pattern, "**", 2)
	prefix := parts[0]
	suffix := ""
	if len(parts) > 1 {
		suffix = strings.TrimLeft(parts[1], "/")
	}

	// Only walk the subtree named by the literal prefix (e.g. `frontend/dist`),
	// not the whole baseDir. Without this, `frontend/dist/**/*` matched every
	// file under baseDir, including unrelated paths like `.wake/cache.json`.
	root := baseDir
	if prefix != "" {
		root = filepath.Join(baseDir, prefix)
	}

	filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		if suffix == "" || suffix == "*" {
			results = append(results, path)
			return nil
		}
		if matchesPattern(filepath.Base(path), suffix) || strings.HasSuffix(path, suffix) {
			results = append(results, path)
		}
		return nil
	})

	return results
}

func matchesPattern(path, pattern string) bool {
	matched, err := filepath.Match(pattern, filepath.Base(path))
	if err != nil {
		return false
	}
	return matched
}
