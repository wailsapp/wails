package exec

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
	"github.com/wailsapp/wails/v3/internal/wake/parse"
)

// goCacheSuffix keys implicit native-Go cache entries separately from the
// explicit sources/generates entries written by RecordTask.
const goCacheSuffix = "#go"

type goCmdKind int

const (
	goCmdNone goCmdKind = iota
	goCmdModTidy
	goCmdBuild
)

// classifyGoCmd recognises the native Go commands wake can cache implicitly.
// The string must already be template-expanded.
func classifyGoCmd(cmd string) goCmdKind {
	f := strings.Fields(cmd)
	if len(f) < 2 || f[0] != "go" {
		return goCmdNone
	}
	switch {
	case f[1] == "build":
		return goCmdBuild
	case f[1] == "mod" && len(f) >= 3 && f[2] == "tidy":
		return goCmdModTidy
	}
	return goCmdNone
}

// goCmdInputs derives the implicit source files an output of a recognised Go
// command. For `go build` the inputs are the module's own .go files, go.mod,
// go.sum, plus everything the task's dependencies generate (e.g. the embedded
// frontend/dist), so a change to any embedded asset still forces a rebuild.
func (e *Executor) goCmdInputs(kind goCmdKind, task *ast.Task, expandedCmd, dir string) (sources []string, output string) {
	goMod := filepath.Join(dir, "go.mod")
	goSum := filepath.Join(dir, "go.sum")

	switch kind {
	case goCmdModTidy:
		return []string{goMod, goSum}, ""
	case goCmdBuild:
		output = parseOutputFlag(expandedCmd)
		if output != "" && !filepath.IsAbs(output) {
			output = filepath.Join(dir, output)
		}
		sources = collectGoFiles(dir)
		sources = append(sources, goMod, goSum)
		sources = append(sources, e.depGenerates(task, map[string]bool{})...)
		return sources, output
	}
	return nil, ""
}

// parseOutputFlag extracts the value of the `-o` flag from a `go build` line.
// Both `-o <path>` and `-o=<path>` are accepted; the previous implementation
// only handled the two-token form, so `go build -o=bin/app` slipped through
// with output="" and the cache could skip without the binary existing.
func parseOutputFlag(cmd string) string {
	f := strings.Fields(cmd)
	for i := 0; i < len(f); i++ {
		if f[i] == "-o" && i+1 < len(f) {
			return f[i+1]
		}
		if strings.HasPrefix(f[i], "-o=") {
			return strings.TrimPrefix(f[i], "-o=")
		}
	}
	return ""
}

// collectGoFiles walks root and returns all non-test .go files, skipping
// directories that never contain buildable module sources.
func collectGoFiles(root string) []string {
	var files []string
	filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			switch d.Name() {
			case "node_modules", ".git", ".wake", ".task", "bin", "dist", "vendor":
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			files = append(files, path)
		}
		return nil
	})
	return files
}

// depGenerates returns the concrete files produced by the transitive
// dependency closure of task, resolved through the same name resolution used
// during execution.
func (e *Executor) depGenerates(task *ast.Task, seen map[string]bool) []string {
	var out []string
	for _, dep := range task.Deps {
		name := e.resolveTaskName(dep.Task, task.Name)
		if seen[name] {
			continue
		}
		seen[name] = true

		depTask := e.Taskfile.Tasks[name]
		if depTask == nil {
			continue
		}

		depDir := e.Dir
		if depTask.Dir != "" {
			if filepath.IsAbs(depTask.Dir) {
				depDir = depTask.Dir
			} else {
				depDir = filepath.Join(e.Dir, depTask.Dir)
			}
		}

		for _, pattern := range depTask.Generates {
			out = append(out, globMatches(depDir, pattern)...)
		}

		out = append(out, e.depGenerates(depTask, seen)...)
	}
	return out
}

func goCmdHash(expandedCmd string) string {
	h := sha256.Sum256([]byte(expandedCmd))
	return fmt.Sprintf("%x", h[:])
}

// ShouldSkipGoCmd reports whether a previously-cached native Go command can be
// skipped: the command (flags) is unchanged, the expected output still exists,
// and no input file has been modified since the last successful run.
func (tc *TaskCache) ShouldSkipGoCmd(name, expandedCmd string, sources []string, output string) bool {
	tc.mu.RLock()
	entry, ok := tc.Entries[name+goCacheSuffix]
	tc.mu.RUnlock()
	if !ok {
		return false
	}
	if entry.Hash != goCmdHash(expandedCmd) {
		return false
	}
	if output != "" {
		if _, err := os.Stat(output); err != nil {
			return false
		}
	}
	for _, src := range sources {
		info, err := os.Stat(src)
		if err != nil {
			continue // optional inputs (e.g. missing go.sum) don't force a rebuild
		}
		if info.ModTime().After(entry.LastRun) {
			if os.Getenv("WAKE_DEBUG") != "" {
				fmt.Fprintf(os.Stderr, "[go-cache] %s: newer source %s (%s > %s)\n", name, src, info.ModTime(), entry.LastRun)
			}
			return false
		}
	}
	return true
}

// RecordGoCmd stores the implicit cache entry for a native Go command after a
// successful run.
func (tc *TaskCache) RecordGoCmd(name, expandedCmd string) error {
	tc.mu.Lock()
	tc.Entries[name+goCacheSuffix] = TaskCacheEntry{
		Hash:    goCmdHash(expandedCmd),
		LastRun: time.Now(),
	}
	tc.mu.Unlock()
	return tc.Save()
}

// expandCmd is a thin wrapper so exec.go can expand without importing parse at
// the call site twice.
func expandCmd(cmd string, vars map[string]*ast.Var) string {
	return parse.ExpandTemplates(cmd, vars)
}
