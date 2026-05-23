package exec

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
	"github.com/wailsapp/wails/v3/internal/wake/cmds"
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

	if err := e.runTask(ctx, task); err != nil {
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

func (e *Executor) runTask(ctx context.Context, task *ast.Task) error {
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

	var env = os.Environ()
	for k, v := range task.Env {
		env = append(env, k+"="+v)
	}

	dir := task.Dir
	if dir != "" && !strings.HasPrefix(dir, "/") {
		dir = e.Dir + "/" + dir
	}

	for _, cmd := range task.Cmds {
		if err := e.runCmd(ctx, cmd, dir, env); err != nil {
			if !cmd.IgnoreError {
				return err
			}
		}
	}

	return nil
}

func (e *Executor) runCmd(ctx context.Context, cmd *ast.Cmd, dir string, env []string) error {
	if cmd.For != nil {
		return e.runForLoop(ctx, cmd, dir, env)
	}
	if cmd.Task != "" {
		subTask := e.Taskfile.Tasks[cmd.Task]
		if subTask == nil {
			return fmt.Errorf("wake: sub-task %q not found", cmd.Task)
		}
		return e.runTask(ctx, subTask)
	}
	if cmd.Cmd == "" {
		return nil
	}

	if !matchesPlatform(cmd.Platforms) {
		return nil
	}

	silent := cmd.Silent || e.Taskfile.Silent || e.Silent
	if !silent {
		fmt.Printf("[wake] %s\n", cmd.Cmd)
	}

	executor := cmds.Route(cmd.Cmd, cmds.RouteOptions{
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

func (e *Executor) runForLoop(ctx context.Context, cmd *ast.Cmd, dir string, env []string) error {
	fl := cmd.For
	if fl == nil {
		return nil
	}

	var items []string
	if fl.Var != "" {
		vr := e.Taskfile.Vars[fl.Var]
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
		_ = item
		if fl.Task != "" {
			subTask := e.Taskfile.Tasks[fl.Task]
			if subTask == nil {
				return fmt.Errorf("wake: for-loop task %q not found", fl.Task)
			}
			if err := e.runTask(ctx, subTask); err != nil {
				return err
			}
		}
	}

	return nil
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
