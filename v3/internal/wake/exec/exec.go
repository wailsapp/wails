package exec

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
	"github.com/wailsapp/wails/v3/internal/wake/cmds"
	"github.com/wailsapp/wails/v3/internal/wake/parse"
)

type Executor struct {
	Taskfile   *ast.Taskfile
	Dir        string
	Verbose    bool
	Silent     bool
	MaxWorkers int
	Cache      *TaskCache
	executed   map[string]bool
	mu         sync.Mutex
}

func (e *Executor) Execute(ctx context.Context, target string) error {
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
		fmt.Fprintf(os.Stderr, "[wake-debug] task %q not found in Taskfile.Tasks (total: %d)\n", target, len(e.Taskfile.Tasks))
		return fmt.Errorf("wake: task %q not found", target)
	}

	if task.Prompt != "" {
		if !confirm(task.Prompt) {
			return fmt.Errorf("wake: task %q cancelled by user", target)
		}
	}

	if e.Cache != nil && e.Cache.ShouldSkip(task, e.Dir) {
		e.log("task %q up-to-date (cache hit)", task.Name)
		e.mu.Lock()
		e.executed[target] = true
		e.mu.Unlock()
		return nil
	}

	if err := e.runTask(ctx, task, nil); err != nil {
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

func (e *Executor) runTask(ctx context.Context, task *ast.Task, depVars map[string]*ast.Var) error {
	if !matchesPlatform(task.Platforms) {
		return nil
	}

	if err := checkPreconditions(task); err != nil {
		return err
	}

	if isUpToDate(task, e.Dir) {
		e.log("skipping task %q (up-to-date)", task.Name)
		return nil
	}

	mergedVars := e.mergeVars(task, depVars)

	for _, dep := range task.Deps {
		expandedDep := parse.ExpandTemplates(dep.Task, mergedVars)
		resolvedDep := e.resolveTaskName(expandedDep, task.Name)
		if err := e.Execute(ctx, resolvedDep); err != nil {
			return err
		}
	}

	var env = os.Environ()
	for k, v := range task.Env {
		expanded := parse.ExpandTemplates(v, mergedVars)
		env = append(env, k+"="+expanded)
	}

	dir := parse.ExpandTemplates(task.Dir, mergedVars)
	if dir != "" && !strings.HasPrefix(dir, "/") {
		dir = e.Dir + "/" + dir
	}

	for _, cmd := range task.Cmds {
		if err := e.runCmd(ctx, cmd, dir, env, mergedVars); err != nil {
			if !cmd.IgnoreError {
				return err
			}
		}
	}

	return nil
}

func (e *Executor) resolveTaskName(name, contextTask string) string {
	if _, ok := e.Taskfile.Tasks[name]; ok {
		return name
	}

	if strings.Contains(contextTask, ":") {
		parts := strings.SplitN(contextTask, ":", 2)
		candidate := parts[0] + ":" + name
		if _, ok := e.Taskfile.Tasks[candidate]; ok {
			return candidate
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

func (e *Executor) runCmd(ctx context.Context, cmd *ast.Cmd, dir string, env []string, vars map[string]*ast.Var) error {
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

	silent := cmd.Silent || e.Taskfile.Silent || e.Silent
	if !silent {
		fmt.Printf("[wake] %s\n", expandedCmd)
	}

	executor := cmds.Route(expandedCmd, cmds.RouteOptions{
		Dir: dir,
		Env: env,
	})

	if silent {
		if sc, ok := executor.(*cmds.ShellCmd); ok {
			sc.Stdout = nil
			sc.Stderr = nil
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

	for k, v := range task.Vars {
		merged[k] = v
	}

	if depVars != nil {
		for k, v := range depVars {
			merged[k] = v
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
	if e.Verbose {
		fmt.Printf(format+"\n", args...)
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
