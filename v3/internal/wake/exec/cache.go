package exec

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
)

const cacheDirName = ".wake"
const cacheFileName = "cache.json"

type TaskCacheEntry struct {
	Hash    string    `json:"hash"`
	LastRun time.Time `json:"last_run"`
}

type TaskCache struct {
	Dir     string                     `json:"-"`
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
	dir := filepath.Join(tc.Dir, cacheDirName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("wake: create cache dir: %w", err)
	}

	path := filepath.Join(dir, cacheFileName)
	data, err := json.MarshalIndent(tc, "", "  ")
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

	hash, err := ComputeTaskHash(task, baseDir)
	if err != nil {
		return false
	}

	entry, ok := tc.Entries[task.Name]
	if !ok {
		return false
	}

	if entry.Hash != hash {
		return false
	}

	for _, pattern := range task.Sources {
		files := globMatches(baseDir, pattern)
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

func (tc *TaskCache) RecordTask(task *ast.Task, baseDir string) error {
	hash, err := ComputeTaskHash(task, baseDir)
	if err != nil {
		return err
	}

	tc.Entries[task.Name] = TaskCacheEntry{
		Hash:    hash,
		LastRun: time.Now(),
	}

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
	parts := strings.Split(pattern, "**")
	suffix := strings.TrimLeft(parts[len(parts)-1], "/")

	filepath.WalkDir(baseDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(baseDir, path)
		if err != nil {
			return nil
		}
		if suffix == "" || strings.HasSuffix(rel, suffix) || matchesPattern(rel, suffix) {
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
