package exec

import (
	"context"
	"errors"
	"fmt"
	"os"
	osexec "os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/internal/report"
	"github.com/wailsapp/wails/v3/internal/wake/ast"
	"github.com/wailsapp/wails/v3/internal/wake/cmds"
	"github.com/wailsapp/wails/v3/internal/wake/parse"
)

type Executor struct {
	Taskfile   *ast.Taskfile
	Dir        string
	Level      report.Verbosity
	Reporter   report.Reporter
	MaxWorkers int
	Cache      *TaskCache
	executed   map[string]bool
	reported   map[string]bool
	depRuns    map[string]bool
	mu         sync.Mutex
}

// beginStep opens a reported step for task unless it has no real commands or was
// already reported. Returning false means the caller should not emit step
// output: this is the single chokepoint that bounds the [k/N] counter at k <= N,
// even when a cache-pruned subtree and a real execution path reach the same task.
func (e *Executor) beginStep(task *ast.Task, label string) bool {
	if !hasRealCmds(task) {
		return false
	}
	e.mu.Lock()
	if e.reported == nil {
		e.reported = make(map[string]bool)
	}
	if e.reported[task.Name] {
		e.mu.Unlock()
		return false
	}
	e.reported[task.Name] = true
	e.mu.Unlock()
	e.rep().StepStart(task.Name, label)
	return true
}

// reportStep emits a complete, instantaneous step (cache hit, skip).
func (e *Executor) reportStep(task *ast.Task, label string, status report.Status) {
	if e.beginStep(task, label) {
		e.rep().StepEnd(status, 0)
	}
}

// reportPrunedCached reports task and the dependency subtree a cache hit just
// pruned, all as cached. Without it the live counter would never reach N on an
// incremental build, because a cache hit skips deps that the plan counted.
func (e *Executor) reportPrunedCached(task *ast.Task) {
	for _, d := range task.Deps {
		if t := e.Taskfile.Tasks[d.Task]; t != nil {
			e.reportPrunedCached(t)
		}
	}
	for _, c := range task.Cmds {
		if c.Task != "" {
			if t := e.Taskfile.Tasks[c.Task]; t != nil {
				e.reportPrunedCached(t)
			}
		}
	}
	e.reportStep(task, "", report.StatusCached)
}

// rep returns the reporter, defaulting to a no-op so the executor is safe to use
// without a reporter wired in (e.g. in tests).
func (e *Executor) rep() report.Reporter {
	if e.Reporter == nil {
		return report.Nop{}
	}
	return e.Reporter
}

func (e *Executor) Execute(ctx context.Context, target string) error {
	return e.ExecuteWithVars(ctx, target, nil)
}

func (e *Executor) ExecuteWithVars(ctx context.Context, target string, extraVars map[string]*ast.Var) error {
	e.mu.Lock()
	if e.executed == nil {
		e.executed = make(map[string]bool)
	}
	if e.executed[target] {
		e.mu.Unlock()
		e.log("task %q already executed, skipping", target)
		return nil
	}
	e.mu.Unlock()

	task := e.Taskfile.Tasks[target]
	if task == nil {
		return fmt.Errorf("wake: task %q not found", target)
	}

	if task.Prompt != "" {
		if !confirm(task.Prompt) {
			return fmt.Errorf("wake: task %q cancelled by user", target)
		}
	}

	if e.Cache != nil && e.Cache.ShouldSkip(task, e.Dir) {
		e.log("task %q up-to-date (cache hit)", task.Name)
		e.reportPrunedCached(task)
		e.mu.Lock()
		e.executed[target] = true
		e.mu.Unlock()
		return nil
	}

	if err := e.runTask(ctx, task, extraVars); err != nil {
		return err
	}

	e.mu.Lock()
	e.executed[target] = true
	e.mu.Unlock()

	if e.Cache != nil {
		if err := e.Cache.RecordTask(task, e.Dir); err != nil {
			e.log("failed to record cache for task %q: %v", task.Name, err)
		}
	}

	RecordRun(target)
	return nil
}

func (e *Executor) resolveDepVars(depVars map[string]*ast.Var, mergedVars map[string]*ast.Var) map[string]*ast.Var {
	if depVars == nil {
		return nil
	}

	resolved := make(map[string]*ast.Var)
	for k, v := range depVars {
		if v.Ref != "" {
			refName := strings.TrimPrefix(v.Ref, ".")
			if ref, ok := mergedVars[refName]; ok {
				val := ref.Value
				if val == "" {
					val = ref.Static
				}
				expanded := parse.ExpandTemplates(val, mergedVars)
				resolved[k] = &ast.Var{Static: expanded, Value: expanded}
				e.rep().Debug(report.DebugLine{
					Category: "var", Subject: k, Arrow: expanded,
					Fields: []report.DebugField{{Key: "ref", Val: v.Ref}},
				})
			} else {
				e.rep().Debug(report.DebugLine{
					Category: "var", Subject: k,
					Fields: []report.DebugField{{Key: "ref", Val: v.Ref}, {Key: "status", Val: "unresolved"}},
				})
				resolved[k] = &ast.Var{Static: v.Ref, Value: v.Ref}
			}
		} else {
			val := v.Value
			if val == "" {
				val = v.Static
			}
			expanded := parse.ExpandTemplates(val, mergedVars)
			resolved[k] = &ast.Var{Static: expanded, Value: expanded}
		}
	}
	return resolved
}

func (e *Executor) runTask(ctx context.Context, task *ast.Task, depVars map[string]*ast.Var) error {
	if !matchesPlatform(task.Platforms) {
		e.reportStep(task, "", report.StatusSkipped)
		return nil
	}

	if err := checkPreconditions(task); err != nil {
		return err
	}

	if isUpToDate(task, e.Dir) {
		e.log("skipping task %q (up-to-date)", task.Name)
		e.reportPrunedCached(task)
		return nil
	}

	mergedVars := e.mergeVars(task, depVars)

	if e.depRuns == nil {
		e.depRuns = make(map[string]bool)
	}

	for _, dep := range task.Deps {
		expandedDep := parse.ExpandTemplates(dep.Task, mergedVars)
		resolvedDep := e.resolveTaskName(expandedDep, task.Name)

		resolvedDepVars := e.resolveDepVars(dep.Vars, mergedVars)

		if e.Level >= report.Debug {
			fields := make([]report.DebugField, 0, len(resolvedDepVars))
			for k, v := range resolvedDepVars {
				val := v.Value
				if val == "" {
					val = v.Static
				}
				fields = append(fields, report.DebugField{Key: k, Val: val})
			}
			sort.Slice(fields, func(i, j int) bool { return fields[i].Key < fields[j].Key })
			e.rep().Debug(report.DebugLine{Category: "dep", Subject: task.Name, Arrow: resolvedDep, Fields: fields})
		}

		method := task.Method
		if method == "" {
			method = "all"
		}

		switch method {
		case "once":
			if e.depRuns[resolvedDep] {
				continue
			}
			e.depRuns[resolvedDep] = true
		case "none":
			continue
		}

		if err := e.ExecuteWithVars(ctx, resolvedDep, resolvedDepVars); err != nil {
			return err
		}
	}

	var env = os.Environ()
	for k, v := range task.Env {
		expanded := parse.ExpandTemplates(v, mergedVars)
		env = append(env, k+"="+expanded)
	}
	// Tell wails subprocess producers (e.g. `wails3 generate bindings`) to emit
	// live feedback as wire events on stdout instead of printing their own UI;
	// the output capture decodes and routes them to this build's reporter.
	env = append(env, "WAKE_REPORT=1")

	dir := parse.ExpandTemplates(task.Dir, mergedVars)
	if dir != "" && !strings.HasPrefix(dir, "/") {
		dir = e.Dir + "/" + dir
	}

	label := parse.ExpandTemplates(task.Label, mergedVars)

	// Implicit caching of native Go commands (`go build`, `go mod tidy`).
	// The task CLI always re-runs these because the Taskfile declares no
	// sources/generates; wake derives the inputs itself and skips the
	// subprocess when nothing has changed. Only applies to single-command
	// tasks that opt out of explicit caching.
	var goRecord func()
	if e.Cache != nil && len(task.Cmds) == 1 &&
		len(task.Sources) == 0 && len(task.Generates) == 0 && len(task.Status) == 0 {
		goDir := dir
		if goDir == "" {
			goDir = e.Dir
		}
		expanded := expandCmd(task.Cmds[0].Cmd, mergedVars)
		if kind := classifyGoCmd(expanded); kind != goCmdNone {
			sources, output := e.goCmdInputs(kind, task, expanded, goDir)
			if e.Cache.ShouldSkipGoCmd(task.Name, expanded, sources, output) {
				e.log("task %q up-to-date (go-cache hit)", task.Name)
				e.reportStep(task, label, report.StatusCached)
				return nil
			}
			goRecord = func() {
				if err := e.Cache.RecordGoCmd(task.Name, expanded); err != nil {
					e.log("go-cache record failed for %q: %v", task.Name, err)
				}
			}
		}
	}

	// A "step" is a task that runs at least one real command; pure wrappers
	// (only dep/task-ref cmds) emit no step of their own — their work shows up
	// as the steps of the tasks they dispatch to. This keeps the [k/N] counter
	// aligned with the visible per-command work and avoids nested steps.
	// A "step" is a task that runs at least one real command; pure wrappers
	// (only dep/task-ref cmds) emit no step of their own. step reports whether
	// this run should be displayed (false if already reported via a cache-pruned
	// path). Commands run regardless; only display is gated.
	step := e.beginStep(task, label)
	var cw *captureWriter
	var start time.Time
	if hasRealCmds(task) {
		cw = newCaptureWriter(e.rep(), e.Level >= report.Verbose)
		start = time.Now()
	}

	for _, cmd := range task.Cmds {
		if err := e.runCmd(ctx, cmd, dir, env, mergedVars, cw); err != nil {
			if !cmd.IgnoreError {
				if cw != nil {
					cw.flush()
				}
				// A failure always surfaces, even if the step was previously
				// reported (e.g. cached) and so would otherwise be silent.
				if !step {
					e.rep().StepStart(task.Name, label)
				}
				e.rep().StepFailed(e.failure(task, cmd, mergedVars, cw, err))
				return err
			}
		}
	}

	if goRecord != nil {
		goRecord()
	}

	if step {
		cw.flush()
		e.rep().StepEnd(report.StatusOK, time.Since(start))
		e.reportArtifacts(task, dir)
	}

	return nil
}

// reportArtifacts walks the task's `generates:` patterns, resolves each glob
// against the task's working directory, and registers each resulting file as
// a build artifact with the reporter. Stat failures are ignored (the file
// may have been intentionally cleaned up before reporting) — only files we
// can confirm exist are surfaced.
func (e *Executor) reportArtifacts(task *ast.Task, dir string) {
	if len(task.Generates) == 0 {
		return
	}
	for _, pattern := range task.Generates {
		full := pattern
		if !filepath.IsAbs(full) {
			full = filepath.Join(dir, pattern)
		}
		matches, err := filepath.Glob(full)
		if err != nil {
			continue
		}
		for _, m := range matches {
			info, err := os.Stat(m)
			if err != nil || info.IsDir() {
				continue
			}
			// Keep the displayed path relative to the wake root when we can
			// — full absolute paths in the summary are noise.
			display := m
			if rel, err := filepath.Rel(e.Dir, m); err == nil && !strings.HasPrefix(rel, "..") {
				display = rel
			}
			e.rep().Artifact(report.Artifact{
				Path: display,
				Size: info.Size(),
			})
		}
	}
}

// hasRealCmds reports whether task runs any shell/native command (as opposed to
// only dispatching to other tasks via deps or task-ref commands).
func hasRealCmds(task *ast.Task) bool {
	for _, c := range task.Cmds {
		if c.Cmd != "" {
			return true
		}
	}
	return false
}

// failure assembles a render-ready Failure for a command that errored.
func (e *Executor) failure(task *ast.Task, cmd *ast.Cmd, vars map[string]*ast.Var, cw *captureWriter, err error) report.Failure {
	f := report.Failure{
		Task:    task.Name,
		Command: parse.ExpandTemplates(cmd.Cmd, vars),
		Output:  cw.output(),
		Err:     err,
	}
	var ee *osexec.ExitError
	if errors.As(err, &ee) {
		f.ExitCode = ee.ExitCode()
	}
	return f
}

func (e *Executor) resolveTaskName(name, contextTask string) string {
	// Same-namespace resolution wins over a same-named top-level task, matching
	// resolveInNamespace at parse time. See its comment for why.
	if strings.Contains(contextTask, ":") {
		parts := strings.SplitN(contextTask, ":", 2)
		prefix := parts[0]

		candidate := prefix + ":" + name
		if _, ok := e.Taskfile.Tasks[candidate]; ok {
			return candidate
		}

		candidate2 := prefix + ":common:" + strings.TrimPrefix(name, "common:")
		if _, ok := e.Taskfile.Tasks[candidate2]; ok {
			return candidate2
		}
	}

	if _, ok := e.Taskfile.Tasks[name]; ok {
		return name
	}

	for _, task := range e.Taskfile.Tasks {
		for _, alias := range task.Aliases {
			if alias == name {
				return task.Name
			}
		}
	}

	for incName := range e.Taskfile.Includes {
		candidate := incName + ":" + name
		if _, ok := e.Taskfile.Tasks[candidate]; ok {
			return candidate
		}
	}

	return name
}

func (e *Executor) runCmd(ctx context.Context, cmd *ast.Cmd, dir string, env []string, vars map[string]*ast.Var, cw *captureWriter) error {
	if cmd.For != nil {
		return e.runForLoop(ctx, cmd, dir, env, vars)
	}
	if cmd.Task != "" {
		expandedTask := parse.ExpandTemplates(cmd.Task, vars)
		resolvedTask := e.resolveTaskName(expandedTask, "")
		subTask := e.Taskfile.Tasks[resolvedTask]
		if subTask == nil {
			return fmt.Errorf("wake: sub-task %q not found", expandedTask)
		}
		return e.runTask(ctx, subTask, cmd.Vars)
	}
	if cmd.Cmd == "" {
		return nil
	}

	if !matchesPlatform(cmd.Platforms) {
		return nil
	}

	expandedCmd := parse.ExpandTemplates(cmd.Cmd, vars)
	e.rep().StepCommand(expandedCmd)

	executor := cmds.Route(expandedCmd, cmds.RouteOptions{
		Dir: dir,
		Env: env,
	})

	// Output is captured (shown only on failure) and streamed live when verbose;
	// the writer also intercepts wire events from subprocess producers and routes
	// them to the live reporter. cw is non-nil whenever the task is a step.
	if cw != nil {
		if s, ok := executor.(cmds.OutputSetter); ok {
			s.SetOutput(cw, cw)
		}
	}

	return executor.Run()
}

func (e *Executor) runForLoop(ctx context.Context, cmd *ast.Cmd, dir string, env []string, vars map[string]*ast.Var) error {
	fl := cmd.For
	if fl == nil {
		return nil
	}

	var items []string
	if fl.Var != "" {
		vr := vars[fl.Var]
		if vr == nil {
			vr = e.Taskfile.Vars[fl.Var]
		}
		if vr != nil {
			val := vr.Value
			if val == "" {
				val = vr.Static
			}
			items = strings.Fields(val)
		}
	}
	if len(items) == 0 {
		items = fl.Items
	}

	for _, item := range items {
		loopVars := make(map[string]*ast.Var)
		for k, v := range vars {
			loopVars[k] = v
		}
		loopVars["ITEM"] = &ast.Var{Static: item, Value: item}

		if fl.Task != "" {
			expandedTask := parse.ExpandTemplates(fl.Task, loopVars)
			subTask := e.Taskfile.Tasks[expandedTask]
			if subTask == nil {
				return fmt.Errorf("wake: for-loop task %q not found", expandedTask)
			}
			if err := e.runTask(ctx, subTask, fl.Vars); err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *Executor) mergeVars(task *ast.Task, depVars map[string]*ast.Var) map[string]*ast.Var {
	merged := make(map[string]*ast.Var)

	for k, v := range e.Taskfile.Vars {
		merged[k] = v
	}

	if depVars != nil {
		for k, v := range depVars {
			merged[k] = v
		}
	}

	for k, v := range task.Vars {
		merged[k] = v
	}

	for k, v := range task.Env {
		if _, exists := merged[k]; !exists {
			merged[k] = &ast.Var{Static: v}
		}
	}

	// Expand templated vars to a fixed point. A single pass is order-dependent:
	// a var like `{{.OUTPUT | default .DEFAULT_OUTPUT}}` resolves to empty when
	// DEFAULT_OUTPUT happens to be visited later in Go's randomized map order.
	// Iterating until no value changes makes resolution deterministic.
	for iter := 0; iter < 10; iter++ {
		changed := false
		for k, v := range merged {
			needsExpand := v.Static != "" && strings.Contains(v.Static, "{{")
			if !needsExpand && v.Value != "" && strings.Contains(v.Value, "{{") {
				needsExpand = true
			}
			if !needsExpand {
				continue
			}
			src := v.Value
			if src == "" {
				src = v.Static
			}
			expanded := parse.ExpandTemplates(src, merged)
			if expanded != v.Value {
				merged[k] = &ast.Var{Static: v.Static, Value: expanded, Ref: v.Ref, Shell: v.Shell}
				changed = true
			}
		}
		if !changed {
			break
		}
	}

	return merged
}

func ExecuteParallel(ctx context.Context, tasks []*ast.Task, execFn func(context.Context, *ast.Task) error, maxWorkers int) error {
	if maxWorkers <= 0 {
		maxWorkers = 4
	}

	sem := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup
	errCh := make(chan error, len(tasks))

	for _, t := range tasks {
		wg.Add(1)
		go func(t *ast.Task) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			if err := execFn(ctx, t); err != nil {
				errCh <- err
			}
		}(t)
	}

	wg.Wait()
	close(errCh)

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

func (e *Executor) log(format string, args ...interface{}) {
	if e.Level >= report.Debug {
		e.rep().Debug(report.DebugLine{Category: "exec", Subject: fmt.Sprintf(format, args...)})
	}
}

func matchesPlatform(platforms []string) bool {
	if len(platforms) == 0 {
		return true
	}
	goos := os.Getenv("GOOS")
	for _, p := range platforms {
		if p == goos {
			return true
		}
	}
	return false
}
