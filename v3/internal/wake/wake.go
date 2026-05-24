package wake

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
	"github.com/wailsapp/wails/v3/internal/wake/exec"
	"github.com/wailsapp/wails/v3/internal/wake/fallback"
	"github.com/wailsapp/wails/v3/internal/wake/override"
	"github.com/wailsapp/wails/v3/internal/wake/parse"
	"github.com/wailsapp/wails/v3/internal/wake/resolve"
)

type ExecuteOptions struct {
	Dir       string
	Platform  string
	Arch      string
	Dev       bool
	Vars      map[string]string
	Verbose   bool
	Silent    bool
	Parallel  bool
}

func Parse(path string) (*ast.Taskfile, error) {
	tf, err := parse.Parse(path)
	if err != nil {
		return nil, err
	}
	if err := parse.ResolveIncludes(tf); err != nil {
		return nil, err
	}
	parse.PopulateBuiltins(tf)
	return tf, nil
}

func Execute(name string, opts ExecuteOptions) error {
	if !useWake() {
		if fallback.Available() {
			return fallback.TaskCLI(name, opts.Dir, nil)
		}
		return fmt.Errorf("wake: WAILS_USE_WAKE not set and task CLI not available")
	}

	dir := opts.Dir
	if dir == "" {
		dir = "."
	}

	tf, err := discoverAndParse(dir)
	if err != nil {
		return err
	}

	if reason, supported := checkSupported(tf); !supported {
		if fallback.Available() {
			return fallback.TaskCLI(name, opts.Dir, nil)
		}
		return fmt.Errorf("wake: unsupported feature: %s (task CLI not available for fallback)", reason)
	}

	if err := override.LoadTaskfileOverrides(dir); err != nil {
		return fmt.Errorf("wake: load Taskfile overrides: %w", err)
	}

	resolve.ApplyOverrides(tf, override.Named(), override.All())

	resolve.FilterPlatforms(tf)

	filterTaskNamespaces(tf, name)

	if err := resolveVars(tf); err != nil {
		return err
	}

	expandTemplates(tf)

	dag, err := resolve.BuildDAG(tf, name)
	if err != nil {
		return err
	}

	cache, err := exec.LoadTaskCache(dir)
	if err != nil {
		return fmt.Errorf("wake: load cache: %w", err)
	}

	ex := &exec.Executor{
		Taskfile: tf,
		Dir:      dir,
		Verbose:  opts.Verbose,
		Silent:   opts.Silent,
		Cache:    cache,
	}

	if opts.Verbose {
		fmt.Fprintf(os.Stderr, "[wake-dag] target=%q order=%v\n", name, dag.Order)
	}

	ctx := context.Background()
	if opts.Parallel {
		if err := executeParallel(ctx, ex, name, dag); err != nil {
			return err
		}
	} else {
		if err := ex.Execute(ctx, name); err != nil {
			return err
		}
	}

	return nil
}

func discoverAndParse(dir string) (*ast.Taskfile, error) {
	candidates := []string{
		filepath.Join(dir, "Taskfile.yml"),
		filepath.Join(dir, "Taskfile.yaml"),
		filepath.Join(dir, "build", "Taskfile.yml"),
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return Parse(path)
		}
	}

	return nil, fmt.Errorf("wake: no Taskfile found in %s", dir)
}

func resolveVars(tf *ast.Taskfile) error {
	if err := parse.ResolveAllVarShells(tf.Vars); err != nil {
		return err
	}
	for _, task := range tf.Tasks {
		if err := parse.ResolveAllVarShells(task.Vars); err != nil {
			return err
		}
	}

	if err := parse.ResolveVars(tf.Vars); err != nil {
		return err
	}
	for _, task := range tf.Tasks {
		if err := parse.ResolveVars(task.Vars); err != nil {
			return err
		}
	}
	return nil
}

func expandTemplates(tf *ast.Taskfile) {
	type taskUpdate struct {
		oldName string
		newName string
		task    *ast.Task
	}

	var updates []taskUpdate
	for taskName, task := range tf.Tasks {
		newName := parse.ExpandTemplates(taskName, tf.Vars)
		task.Name = newName
		task.Dir = parse.ExpandTemplates(task.Dir, tf.Vars)
		// Label is NOT pre-expanded here - it's expanded at execution time with mergedVars
		task.Summary = parse.ExpandTemplates(task.Summary, tf.Vars)

		for i, dep := range task.Deps {
			dep.Task = parse.ExpandTemplates(dep.Task, tf.Vars)
			task.Deps[i] = dep
		}

		for _, cmd := range task.Cmds {
			// Commands are NOT pre-expanded here - they're expanded at execution time with mergedVars
			cmd.Task = parse.ExpandTemplates(cmd.Task, tf.Vars)
		}

		for k, v := range task.Env {
			task.Env[k] = parse.ExpandTemplates(v, tf.Vars)
		}

		updates = append(updates, taskUpdate{oldName: taskName, newName: newName, task: task})
	}

	for _, u := range updates {
		delete(tf.Tasks, u.oldName)
		tf.Tasks[u.newName] = u.task
	}

	resolveDepNamespaces(tf)
}

func filterTaskNamespaces(tf *ast.Taskfile, target string) {
	if !strings.Contains(target, ":") {
		return
	}

	prefix := strings.SplitN(target, ":", 2)[0]

	platformPrefixes := map[string]bool{
		"darwin":  true,
		"linux":   true,
		"windows": true,
		"ios":     true,
		"android": true,
	}

	if !platformPrefixes[prefix] {
		return
	}

	for name := range tf.Tasks {
		if strings.HasPrefix(name, "common:") {
			delete(tf.Tasks, name)
			continue
		}

		if strings.Contains(name, ":") {
			parts := strings.SplitN(name, ":", 2)
			taskPrefix := parts[0]
			if platformPrefixes[taskPrefix] && taskPrefix != prefix {
				delete(tf.Tasks, name)
				continue
			}

			if len(parts) >= 3 {
				taskPrefix2 := parts[1]
				if platformPrefixes[taskPrefix2] && taskPrefix2 != prefix {
					delete(tf.Tasks, name)
					continue
				}
			}
		}
	}
}

func resolveDepNamespaces(tf *ast.Taskfile) {
	for _, task := range tf.Tasks {
		for i, dep := range task.Deps {
			if resolved, ok := resolveInNamespace(tf, task.Name, dep.Task); ok {
				task.Deps[i].Task = resolved
			}
		}
		for _, cmd := range task.Cmds {
			if cmd.Task == "" {
				continue
			}
			if resolved, ok := resolveInNamespace(tf, task.Name, cmd.Task); ok {
				cmd.Task = resolved
			}
		}
	}
}

// resolveInNamespace resolves a short task reference made from within
// contextName's namespace. A bare name inside `darwin:foo` resolves to
// `darwin:<name>` *before* any same-named top-level task — this is Taskfile's
// local-namespace-wins rule. Without it, a dep like `build` in `darwin:package`
// binds to the top-level `build` wrapper, which silently drops the caller's
// vars (e.g. PRODUCTION=true), so prod packaging built with dev flags.
func resolveInNamespace(tf *ast.Taskfile, contextName, name string) (string, bool) {
	if strings.Contains(contextName, ":") {
		prefix := strings.SplitN(contextName, ":", 2)[0]
		if _, ok := tf.Tasks[prefix+":"+name]; ok {
			return prefix + ":" + name, true
		}
		candidate := prefix + ":common:" + strings.TrimPrefix(name, "common:")
		if _, ok := tf.Tasks[candidate]; ok {
			return candidate, true
		}
	}
	if _, ok := tf.Tasks[name]; ok {
		return name, true
	}
	for incName := range tf.Includes {
		if _, ok := tf.Tasks[incName+":"+name]; ok {
			return incName + ":" + name, true
		}
	}
	return name, false
}

func useWake() bool {
	if env := os.Getenv("WAILS_USE_WAKE"); env != "" {
		return env == "true"
	}
	return false
}

func checkSupported(tf *ast.Taskfile) (string, bool) {
	if len(tf.Dotenv) > 0 {
		return "dotenv", false
	}
	if tf.Output != "" && tf.Output != "interleaved" {
		return "output: " + tf.Output, false
	}
	if tf.Requires != nil {
		return "requires", false
	}
	if tf.Interval != "" {
		return "interval", false
	}

	for _, task := range tf.Tasks {
		if task.Run != "" && task.Run != "always" {
			return "run: " + task.Run + " in task " + task.Name, false
		}
		if task.Short != "" {
			return "short in task " + task.Name, false
		}
		if len(task.Defer) > 0 {
			return "defer in task " + task.Name, false
		}
		if task.Interval != "" {
			return "interval in task " + task.Name, false
		}
	}

	return "", true
}

func executeParallel(ctx context.Context, ex *exec.Executor, target string, dag *resolve.DAG) error {
	inDeg := make(map[string]int)
	for k, v := range dag.InDegree {
		inDeg[k] = v
	}

	completed := make(map[string]bool)
	var mu sync.Mutex

	for len(completed) < len(dag.Tasks) {
		var ready []*ast.Task
		for _, task := range dag.Tasks {
			if completed[task.Name] {
				continue
			}
			if inDeg[task.Name] == 0 {
				ready = append(ready, task)
			}
		}

		if len(ready) == 0 {
			return fmt.Errorf("wake: deadlock detected in parallel execution")
		}

		var wg sync.WaitGroup
		errCh := make(chan error, len(ready))

		for _, task := range ready {
			wg.Add(1)
			go func(t *ast.Task) {
				defer wg.Done()
				if err := ex.Execute(ctx, t.Name); err != nil {
					errCh <- err
				}
				mu.Lock()
				completed[t.Name] = true
				for _, dependent := range dag.Edges[t.Name] {
					inDeg[dependent]--
				}
				mu.Unlock()
			}(task)
		}

		wg.Wait()
		close(errCh)

		for err := range errCh {
			if err != nil {
				return err
			}
		}
	}

	return nil
}
