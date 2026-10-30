package exec

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
	"github.com/wailsapp/wails/v3/internal/wake/cmds"
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
	switch cmds.Route(cmd, cmds.RouteOptions{}).(type) {
	case *cmds.GoBuildCmd:
		return goCmdBuild
	case *cmds.GoModTidyCmd:
		return goCmdModTidy
	}
	return goCmdNone
}

// goCmdInputs derives implicit sources and output for a recognised Go command.
// For `go build` the inputs are the module's Go/native source files, go.mod,
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
		sources = collectGoSourceFiles(dir)
		sources = append(sources, goMod, goSum)
		sources = append(sources, e.depGenerates(task, map[string]bool{})...)
		return sources, output
	}
	return nil, ""
}

// parseOutputFlag reuses native argument parsing so quoted and escaped paths
// match the output location used by command execution.
func parseOutputFlag(cmd string) string {
	if build, ok := cmds.Route(cmd, cmds.RouteOptions{}).(*cmds.GoBuildCmd); ok {
		return build.Opts.Output
	}
	return ""
}

// collectGoSourceFiles walks root for Go and native compilation inputs, excluding
// Go test files and the existing generated/dependency directories.
func collectGoSourceFiles(root string) []string {
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
		switch filepath.Ext(path) {
		case ".go":
			if !strings.HasSuffix(path, "_test.go") {
				files = append(files, path)
			}
		case ".c", ".cc", ".cpp", ".cxx", ".m", ".h", ".hh", ".hpp", ".hxx",
			".f", ".F", ".for", ".f90", ".s", ".S", ".sx", ".swig", ".swigcxx", ".syso":
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
	var encoded [sha256.Size * 2]byte
	hex.Encode(encoded[:], h[:])
	return string(encoded[:])
}

// goCmdInputsHash includes the discovered file set without depending on its order.
func goCmdInputsHash(expandedCmd string, sources []string) string {
	paths := slices.Clone(sources)
	slices.Sort(paths)
	paths = slices.Compact(paths)
	// Successful command arguments and file names cannot contain NUL bytes.
	return goCmdHash(expandedCmd + "\x00" + strings.Join(paths, "\x00"))
}

// ShouldSkipGoCmd reports whether a previously-cached native Go command can be
// skipped: the command and discovered input-file set are unchanged, the output
// still exists, and no input has been modified since the last successful run.
func (tc *TaskCache) ShouldSkipGoCmd(name, expandedCmd string, sources []string, output string) bool {
	tc.mu.RLock()
	entry, ok := tc.Entries[name+goCacheSuffix]
	tc.mu.RUnlock()
	if !ok {
		return false
	}
	if entry.Hash != goCmdInputsHash(expandedCmd, sources) {
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
func (tc *TaskCache) RecordGoCmd(name, expandedCmd string, sources []string) error {
	tc.mu.Lock()
	tc.Entries[name+goCacheSuffix] = TaskCacheEntry{
		Hash:    goCmdInputsHash(expandedCmd, sources),
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
