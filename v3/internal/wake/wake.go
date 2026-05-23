package wake

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	Vars     map[string]string
	Verbose  bool
	Silent   bool
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

	if err := override.LoadTaskfileOverrides(dir); err != nil {
		return fmt.Errorf("wake: load Taskfile overrides: %w", err)
	}

	resolve.ApplyOverrides(tf, override.Named(), override.All())

	resolve.FilterPlatforms(tf)

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

	ctx := context.Background()
	for _, taskName := range dag.Order {
		if err := ex.Execute(ctx, taskName); err != nil {
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
		task.Label = parse.ExpandTemplates(task.Label, tf.Vars)
		task.Summary = parse.ExpandTemplates(task.Summary, tf.Vars)

		for i, dep := range task.Deps {
			dep.Task = parse.ExpandTemplates(dep.Task, tf.Vars)
			task.Deps[i] = dep
		}

		for _, cmd := range task.Cmds {
			cmd.Cmd = parse.ExpandTemplates(cmd.Cmd, tf.Vars)
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

func resolveDepNamespaces(tf *ast.Taskfile) {
	for _, task := range tf.Tasks {
		for i, dep := range task.Deps {
			if _, ok := tf.Tasks[dep.Task]; ok {
				continue
			}

			if strings.Contains(task.Name, ":") {
				parts := strings.SplitN(task.Name, ":", 2)
				candidate := parts[0] + ":" + dep.Task
				if _, ok := tf.Tasks[candidate]; ok {
					task.Deps[i].Task = candidate
					continue
				}
			}

			for incName := range tf.Includes {
				candidate := incName + ":" + dep.Task
				if _, ok := tf.Tasks[candidate]; ok {
					task.Deps[i].Task = candidate
					break
				}
			}
		}

		for _, cmd := range task.Cmds {
			if cmd.Task == "" {
				continue
			}
			if _, ok := tf.Tasks[cmd.Task]; ok {
				continue
			}

			if strings.Contains(task.Name, ":") {
				parts := strings.SplitN(task.Name, ":", 2)
				candidate := parts[0] + ":" + cmd.Task
				if _, ok := tf.Tasks[candidate]; ok {
					cmd.Task = candidate
					continue
				}
			}

			for incName := range tf.Includes {
				candidate := incName + ":" + cmd.Task
				if _, ok := tf.Tasks[candidate]; ok {
					cmd.Task = candidate
					break
				}
			}
		}
	}
}

func useWake() bool {
	if env := os.Getenv("WAILS_USE_WAKE"); env != "" {
		return env == "true"
	}
	return false
}
