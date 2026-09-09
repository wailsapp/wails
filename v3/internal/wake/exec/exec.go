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
	"sync/atomic"
	"time"

	"github.com/wailsapp/wails/v3/internal/report"
	"github.com/wailsapp/wails/v3/internal/wake/ast"
	"github.com/wailsapp/wails/v3/internal/wake/cmds"
	"github.com/wailsapp/wails/v3/internal/wake/parse"
	wakeplatform "github.com/wailsapp/wails/v3/internal/wake/platform"
)

type Executor struct {
	Taskfile        *ast.Taskfile
	Dir             string
	Level           report.Verbosity
	Reporter        report.Reporter
	MaxWorkers      int
	Cache           *TaskCache
	// Parallel toggles the in-task dep-fanout. When true, each task's
	// `deps:` are executed concurrently (matching Taskfile semantics);
	// when false, sequentially.
	Parallel        bool
	// Force, when true, makes every task re-run regardless of cache state.
	// Wired up from WAKE_FORCE=true (or a `wails3 build --clean` flag) at
	// the wake.Execute boundary. Bypasses both the Taskfile-declared cache
	// check (sources/generates/status) and the implicit native-Go cache.
	Force           bool
	executed        map[string]bool
	reported        map[string]bool
	depRuns         map[string]bool
	failureReported atomic.Bool
	mu              sync.Mutex
}

// FailureReported reports whether the executor has rendered at least one
// StepFailed panel during this run. wake.go uses it to decide whether to
// surface the top-level error itself — when nothing has been rendered, the
// user gets no diagnostic unless wake echoes the error message.
func (e *Executor) FailureReported() bool { return e.failureReported.Load() }

// beginStep opens a reported step for task unless it has no real commands or
// was already reported. Returns the StepID (0 if no step was opened) so the
// caller can thread it through paired step-scoped methods. This is the single
// chokepoint that bounds the [k/N] counter at k <= N even when a cache-pruned
// subtree and a real execution path reach the same task.
func (e *Executor) beginStep(task *ast.Task, label string) report.StepID {
	if !hasRealCmds(task) {
		return 0
	}
	e.mu.Lock()
	if e.reported == nil {
		e.reported = make(map[string]bool)
	}
	if e.reported[task.Name] {
		e.mu.Unlock()
		return 0
	}
	e.reported[task.Name] = true
	e.mu.Unlock()
	return e.rep().StepStart(task.Name, label)
}

// reportStep emits a complete, instantaneous step (cache hit, skip).
func (e *Executor) reportStep(task *ast.Task, label string, status report.Status) {
	if id := e.beginStep(task, label); id != 0 {
		e.rep().StepEnd(id, status, 0)
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

	if !e.Force && e.Cache != nil && e.Cache.ShouldSkip(task, e.Dir) {
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

	// Preconditions are templated against the task's resolved vars (their `sh:`
	// guards reference vars like .OBFUSCATED), so they need the merged map and
	// must run before the up-to-date short-circuit below. Compute the map
	// lazily: tasks without preconditions that turn out to be up-to-date skip
	// the fixed-point expansion entirely.
	var mergedVars map[string]*ast.Var
	if len(task.Precondition) > 0 {
		mergedVars = e.mergeVars(task, depVars)
		if err := checkPreconditions(task, mergedVars); err != nil {
			return err
		}
	}

	if !e.Force && isUpToDate(task, e.Dir) {
		e.log("skipping task %q (up-to-date)", task.Name)
		e.reportPrunedCached(task)
		return nil
	}

	if mergedVars == nil {
		mergedVars = e.mergeVars(task, depVars)
	}

	// depRuns is shared across goroutines once the parallel dep fanout is
	// in flight, so take e.mu for both the lazy-init and every read/write.
	// The previous code mutated it without locks and could panic on Go's
	// "concurrent map read and map write" detector under Parallel=true.
	e.mu.Lock()
	if e.depRuns == nil {
		e.depRuns = make(map[string]bool)
	}
	e.mu.Unlock()

	// Taskfile semantics: deps run in parallel before cmds. In serial mode
	// (Parallel=false) we fall back to a sequential walk; otherwise we spawn
	// one goroutine per dep and wait for all. Each goroutine bottoms out in
	// e.ExecuteWithVars, which is mutex-safe; runaway concurrency from very
	// wide dep fans is bounded by the OS scheduler — we don't add a per-task
	// worker pool because the outer call chain is usually shallow.
	type depPlan struct {
		name string
		vars map[string]*ast.Var
	}
	var plans []depPlan
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
			e.mu.Lock()
			already := e.depRuns[resolvedDep]
			if !already {
				e.depRuns[resolvedDep] = true
			}
			e.mu.Unlock()
			if already {
				continue
			}
		case "none":
			continue
		}

		plans = append(plans, depPlan{name: resolvedDep, vars: resolvedDepVars})
	}

	if !e.Parallel || len(plans) <= 1 {
		for _, p := range plans {
			if err := e.ExecuteWithVars(ctx, p.name, p.vars); err != nil {
				return err
			}
		}
	} else {
		var wg sync.WaitGroup
		errCh := make(chan error, len(plans))
		for _, p := range plans {
			wg.Add(1)
			go func(p depPlan) {
				defer wg.Done()
				if err := e.ExecuteWithVars(ctx, p.name, p.vars); err != nil {
					errCh <- err
				}
			}(p)
		}
		wg.Wait()
		close(errCh)
		for err := range errCh {
			if err != nil {
				return err
			}
		}
	}

	var env = os.Environ()
	// Taskfile-level env first; per-task env follows and wins on collisions
	// (matches the upstream task semantics, and matches the order any sane
	// reader of the Taskfile expects).
	for k, v := range e.Taskfile.Env {
		expanded := parse.ExpandTemplates(v, mergedVars)
		env = append(env, k+"="+expanded)
	}
	for k, v := range task.Env {
		expanded := parse.ExpandTemplates(v, mergedVars)
		env = append(env, k+"="+expanded)
	}
	// Tell wails subprocess producers (e.g. `wails3 generate bindings`) to emit
	// live feedback as wire events on stdout instead of printing their own UI;
	// the output capture decodes and routes them to this build's reporter.
	env = append(env, "WAKE_REPORT=1")

	dir := parse.ExpandTemplates(task.Dir, mergedVars)
	if dir != "" && !filepath.IsAbs(dir) {
		// filepath.IsAbs handles `C:\…` on Windows; the previous
		// `strings.HasPrefix("/")` check misclassified Windows absolute
		// paths as relative and prefixed them with e.Dir, producing nonsense
		// like `C:\proj\C:\foo` for the task's working directory.
		dir = filepath.Join(e.Dir, dir)
	}

	label := parse.ExpandTemplates(task.Label, mergedVars)

	// Implicit caching of native Go commands (`go build`, `go mod tidy`).
	// The task CLI always re-runs these because the Taskfile declares no
	// sources/generates; wake derives the inputs itself and skips the
	// subprocess when nothing has changed. Only applies to single-command
	// tasks that opt out of explicit caching.
	var goRecord func()
	if !e.Force && e.Cache != nil && len(task.Cmds) == 1 &&
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
	// (only dep/task-ref cmds) emit no step of their own. stepID is non-zero
	// when this run should be displayed (zero if already reported via a
	// cache-pruned path). Commands run regardless; only display is gated.
	stepID := e.beginStep(task, label)
	var cw *captureWriter
	var start time.Time
	if hasRealCmds(task) {
		cw = newCaptureWriter(e.rep(), stepID, e.Level >= report.Verbose)
		start = time.Now()
	}

	for _, cmd := range task.Cmds {
		if err := e.runCmd(ctx, cmd, dir, env, mergedVars, stepID, cw); err != nil {
			if !cmd.IgnoreError {
				if cw != nil {
					cw.flush()
				}
				// Only render a failure panel for a real shell command. A
				// task-ref cmd (cmd.Cmd == "") simply propagates the error
				// from a sub-task that has already rendered its own panel —
				// rendering a second one here would duplicate the message and,
				// because the parent task usually has no captureWriter, would
				// crash trying to attach captured output that doesn't exist.
				if cmd.Cmd != "" {
					if stepID == 0 {
						stepID = e.rep().StepStart(task.Name, label)
					}
					e.rep().StepFailed(stepID, e.failure(task, cmd, mergedVars, cw, err))
					e.failureReported.Store(true)
				}
				return err
			}
		}
	}

	if goRecord != nil {
		goRecord()
	}

	if stepID != 0 {
		cw.flush()
		e.rep().StepEnd(stepID, report.StatusOK, time.Since(start))
	}

	return nil
}

// reportGoBuildArtifact extracts a `go build ... -o <path>` output from the
// expanded command string and registers the resulting file as a build
// artifact. We deliberately do NOT walk task `generates:` patterns: those
// produce intermediate outputs (icons, bindings) that aren't what the user
// asks for at the end of `wails3 build` — they want the binary. Catching
// go-build's `-o` argument hits exactly that.
//
// Called after a real-cmd execution succeeds. dir is the task's resolved
// working directory; the -o value is rooted there if relative.
func (e *Executor) reportGoBuildArtifact(expandedCmd, dir string) {
	if !strings.HasPrefix(strings.TrimLeft(expandedCmd, " \t"), "go build") {
		return
	}
	out := extractGoBuildOutput(expandedCmd)
	if out == "" {
		return
	}
	full := out
	if !filepath.IsAbs(full) {
		base := dir
		if base == "" {
			base = e.Dir
		}
		full = filepath.Join(base, out)
	}
	info, err := os.Stat(full)
	if err != nil || info.IsDir() {
		return
	}
	display := full
	if rel, err := filepath.Rel(e.Dir, full); err == nil && !strings.HasPrefix(rel, "..") {
		display = rel
	}
	e.rep().Artifact(report.Artifact{
		Path: display,
		Size: info.Size(),
		Kind: "binary",
	})
}

// extractGoBuildOutput scans a `go build ...` command line for the value
// passed to -o (either `-o X` or `-o=X`). Quotes around X are stripped.
// Returns "" if no -o flag is present.
func extractGoBuildOutput(cmd string) string {
	// shellSplit lives in cmds/ but the algorithm is small; we re-implement
	// the minimal quote-aware split here to avoid a cycle.
	tokens := splitArgsQuoted(cmd)
	for i := 0; i < len(tokens); i++ {
		t := tokens[i]
		if t == "-o" && i+1 < len(tokens) {
			return tokens[i+1]
		}
		if strings.HasPrefix(t, "-o=") {
			return strings.TrimPrefix(t, "-o=")
		}
	}
	return ""
}

func splitArgsQuoted(s string) []string {
	var (
		out    []string
		cur    strings.Builder
		inS    bool
		inD    bool
		active bool
	)
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '\\' && !inS && i+1 < len(s):
			cur.WriteByte(s[i+1])
			i++
			active = true
		case c == '\'' && !inD:
			inS = !inS
			active = true
		case c == '"' && !inS:
			inD = !inD
			active = true
		case (c == ' ' || c == '\t') && !inS && !inD:
			if active {
				out = append(out, cur.String())
				cur.Reset()
				active = false
			}
		default:
			cur.WriteByte(c)
			active = true
		}
	}
	if active {
		out = append(out, cur.String())
	}
	return out
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

func (e *Executor) runCmd(ctx context.Context, cmd *ast.Cmd, dir string, env []string, vars map[string]*ast.Var, stepID report.StepID, cw *captureWriter) error {
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
		// cmd-level vars (the `vars:` block on a `task:` cmd ref) need to
		// be resolved in the *calling* task's scope before being passed to
		// the sub-task. Without this, templates like
		// `SCRIPT: '{{if eq .DEV "true"}}...{{end}}'` reach the sub-task
		// unresolved and end up rendering empty there — same path as the
		// existing resolveDepVars treatment for `task.Deps` vars.
		resolvedVars := e.resolveDepVars(cmd.Vars, vars)
		return e.runTask(ctx, subTask, resolvedVars)
	}
	if cmd.Cmd == "" {
		return nil
	}

	if !matchesPlatform(cmd.Platforms) {
		return nil
	}

	expandedCmd := parse.ExpandTemplates(cmd.Cmd, vars)
	e.rep().StepCommand(stepID, expandedCmd)

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

	err := executor.Run()
	if err == nil {
		// Only surface a `go build -o <path>` output as an Artifact when the
		// command actually succeeded. The previous `defer` ran on every code
		// path, which could promote a stale binary left over from a prior
		// successful build into the summary's Output panel as if this run
		// had produced it.
		e.reportGoBuildArtifact(expandedCmd, dir)
	}
	return err
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
			// Resolve fl.Vars (the per-cmd `vars:` block on the for-loop
			// task ref) against loopVars first, so any template that
			// references {{.ITEM}} is expanded to the current loop item
			// before the sub-task receives it. The previous code passed
			// fl.Vars through unresolved, so per-item templates never
			// saw the current item and rendered empty.
			resolved := e.resolveDepVars(fl.Vars, loopVars)
			if err := e.runTask(ctx, subTask, resolved); err != nil {
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

	// Task-local vars layer on top of depVars, but with a twist: if a task
	// var's template references the same name (the "passthrough" pattern
	// `SCRIPT: "{{.SCRIPT}}"`), evaluate it against the *current* merged map
	// so it picks up the inherited depVars value instead of overwriting it
	// with an immediately-self-referential empty. Without this, every
	// Taskfile that declares accepted vars via passthrough silently loses
	// them at the next task boundary.
	for k, v := range task.Vars {
		if v != nil && v.Value == "" && strings.Contains(v.Static, "{{."+k) {
			expanded := parse.ExpandTemplates(v.Static, merged)
			if !strings.Contains(expanded, "{{") {
				merged[k] = &ast.Var{Static: v.Static, Value: expanded}
				continue
			}
		}
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
	// `GOOS` env is only reliably set during cross-compilation. At normal
	// build/run time it is usually unset, and `os.Getenv` would return ""
	// — making *every* platform-gated task skip silently. Fall back to
	// the runtime OS via the platform package so wake's behaviour matches
	// upstream task and the obvious user expectation.
	goos := os.Getenv("GOOS")
	if goos == "" {
		goos = wakeplatform.OS()
	}
	for _, p := range platforms {
		if p == goos {
			return true
		}
	}
	return false
}
