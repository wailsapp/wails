package wake

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/internal/report"
	"github.com/wailsapp/wails/v3/internal/report/pulse"
	"github.com/wailsapp/wails/v3/internal/wake/ast"
	"github.com/wailsapp/wails/v3/internal/wake/exec"
	"github.com/wailsapp/wails/v3/internal/wake/fallback"
	"github.com/wailsapp/wails/v3/internal/wake/override"
	"github.com/wailsapp/wails/v3/internal/wake/parse"
	"github.com/wailsapp/wails/v3/internal/wake/resolve"
)

type ExecuteOptions struct {
	Dir      string
	Platform string
	Arch     string
	Dev      bool
	Verb     string // command the user ran (e.g. "build"); shown in the header
	Vars     map[string]string
	Silent   bool // show failures only
	Verbose  bool // show commands + live subprocess output
	Debug    bool // show resolver internals (implies verbose detail)
	Parallel bool
	Force    bool // ignore all task and Go caches; run every task from scratch
}

// verbosity collapses the boolean options into a single level. Debug wins, then
// Verbose, then Silent; the default is Normal.
func (o ExecuteOptions) verbosity() report.Verbosity {
	switch {
	case o.Debug:
		return report.Debug
	case o.Verbose:
		return report.Verbose
	case o.Silent:
		return report.Silent
	default:
		return report.Normal
	}
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

	if !overridesDisabled() {
		locals, err := override.LoadLocal(dir)
		if err != nil {
			return fmt.Errorf("wake: load local Taskfile overrides: %w", err)
		}
		for _, local := range locals {
			resolve.MergeTaskfile(tf, local)
		}
	}

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

	level := opts.verbosity()
	rep := pulse.New(os.Stdout, level)

	// Print a single-line experimental notice once per wake invocation, unless
	// the user has explicitly silenced wake output. Wake is gated behind
	// WAILS_USE_WAKE=true and the default build path remains the Task runtime;
	// the notice exists so users running the opt-in path always know they're
	// on the experimental track and can correlate any anomalies they hit.
	// Suppress it under WAKE_SILENT (which already implies a non-interactive
	// or scripted invocation) and during cache-hit dry-runs / DAG previews.
	if !opts.Silent && os.Getenv("WAKE_NOTICE") != "off" {
		fmt.Fprintln(os.Stderr,
			"  wake (experimental) · set WAILS_USE_WAKE= to disable · WAKE_NOTICE=off to hide this notice")
	}

	// Make the reporter reachable by in-process producers for the duration of
	// the build. Subprocess producers reach it over the wire protocol instead
	// (see internal/report and the executor's output capture).
	report.SetActive(rep)
	defer report.SetActive(nil)

	ex := &exec.Executor{
		Taskfile: tf,
		Dir:      dir,
		Level:    level,
		Reporter: rep,
		Cache:    cache,
		Force:    opts.Force,
		Parallel: opts.Parallel,
	}

	rep.BuildStart(opts.Verb, name, countSteps(tf, name))

	if level >= report.Debug {
		rep.Debug(report.DebugLine{Category: "dag", Subject: name, Arrow: strings.Join(dag.Order, " · ")})
	}

	ctx := context.Background()
	start := time.Now()
	var runErr error
	if opts.Parallel {
		runErr = executeParallel(ctx, ex, name, dag)
	} else {
		runErr = ex.Execute(ctx, name)
	}

	// If the build errored but nothing rendered a per-step failure panel
	// (e.g. an early task-resolution error, a precondition that aborted
	// before any step ran), surface the underlying message to the reporter
	// as a synthetic step failure so the user actually sees *why*. Without
	// this, the build collapses to a bare "build failed after 0ms" verdict.
	if runErr != nil && !ex.FailureReported() {
		id := rep.StepStart(name, "")
		rep.StepFailed(id, report.Failure{
			Task: name,
			Err:  runErr,
		})
	}

	rep.BuildEnd(time.Since(start), runErr == nil)
	if runErr != nil {
		return errReported{runErr}
	}
	return nil
}

// countSteps mirrors execution to predict how many steps will be shown: it walks
// the target's transitive dependencies (deps) and task-ref commands, and counts
// the unique tasks that run at least one real command. This is the N in the
// [k/N] step counter. Names are already namespace-resolved at this point.
func countSteps(tf *ast.Taskfile, target string) int {
	seen := make(map[string]bool)
	steps := make(map[string]bool)

	var visit func(name string)
	visit = func(name string) {
		if seen[name] {
			return
		}
		seen[name] = true
		t := tf.Tasks[name]
		if t == nil {
			return
		}
		for _, d := range t.Deps {
			visit(d.Task)
		}
		real := false
		for _, c := range t.Cmds {
			if c.Cmd != "" {
				real = true
			}
			if c.Task != "" {
				visit(c.Task)
			}
			if c.For != nil && c.For.Task != "" {
				visit(c.For.Task)
			}
		}
		if real {
			steps[name] = true
		}
	}

	visit(target)
	return len(steps)
}

// errReported marks an error whose failure has already been rendered by the
// reporter, so the top-level CLI can avoid printing it a second time.
type errReported struct{ err error }

func (e errReported) Error() string { return e.err.Error() }
func (e errReported) Unwrap() error { return e.err }

// IsReported reports whether err was already rendered to the build UI.
func IsReported(err error) bool {
	var r errReported
	return errors.As(err, &r)
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
	// After settling the simple cases (shell, ref, plain static), evaluate
	// any root vars whose Static is a template — e.g.
	// `{{.PACKAGE_MANAGER | default "npm"}}`. Doing this only at the root
	// level is deliberate: task-local templates frequently reference root
	// vars and must wait for mergeVars at execution time, where the full
	// scope is available.
	parse.ExpandVarTemplates(tf.Vars)

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
			// SplitN(..., 2) capped at 2 elements, so the previous
			// `len(parts) >= 3` nested-namespace branch was dead code.
			// Use the full Split here so a nested platform prefix like
			// `darwin:linux:foo` (rare but possible after an aggressive
			// include merge) is also filtered out on a darwin build.
			parts := strings.Split(name, ":")
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

// overridesDisabled reports whether local override files (Taskfile.local.* /
// Taskfile.override.*) should be ignored. Setting WAILS_NO_OVERRIDES=true skips
// auto-loading entirely, giving CI and security-sensitive builds deterministic
// behaviour from the committed base Taskfile alone.
func overridesDisabled() bool {
	return os.Getenv("WAILS_NO_OVERRIDES") == "true"
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
