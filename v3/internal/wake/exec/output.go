package exec

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
	"github.com/wailsapp/wails/v3/internal/wake/platform"
)

type OutputMode int

const (
	OutputDefault OutputMode = iota
	OutputSilent
	OutputVerbose
)

type Runner struct {
	Dir      string
	Env      []string
	OutMode  OutputMode
	Taskfile *ast.Taskfile
}

func (r *Runner) Run(cmd *ast.Cmd) error {
	if cmd.For != nil {
		return r.runFor(cmd)
	}
	if cmd.Task != "" {
		return r.runSubTask(cmd.Task)
	}
	if cmd.Cmd == "" {
		return nil
	}

	if !platform.Filter(cmd.Platforms) {
		return nil
	}

	silent := cmd.Silent || r.Taskfile.Silent
	if r.OutMode == OutputSilent {
		silent = true
	}

	if !silent {
		fmt.Printf("wake: %s\n", cmd.Cmd)
	}

	c := exec.Command("sh", "-c", cmd.Cmd)
	c.Dir = r.Dir
	c.Env = append(os.Environ(), r.Env...)

	if silent && r.OutMode != OutputVerbose {
		c.Stdout = nil
		c.Stderr = nil
	} else {
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
	}

	return c.Run()
}

func (r *Runner) runFor(cmd *ast.Cmd) error {
	fl := cmd.For
	if fl == nil || fl.Var == "" {
		return nil
	}

	vr := r.Taskfile.Vars[fl.Var]
	if vr == nil {
		return nil
	}

	val := vr.Value
	if val == "" {
		val = vr.Static
	}
	items := strings.Fields(val)

	for _, item := range items {
		_ = item
		if fl.Task != "" {
			sub := r.Taskfile.Tasks[fl.Task]
			if sub == nil {
				return fmt.Errorf("wake: for-loop task %q not found", fl.Task)
			}
			subRunner := &Runner{
				Dir:      r.Dir,
				Env:      r.Env,
				OutMode:  r.OutMode,
				Taskfile: r.Taskfile,
			}
			for _, c := range sub.Cmds {
				if err := subRunner.Run(c); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (r *Runner) runSubTask(name string) error {
	task := r.Taskfile.Tasks[name]
	if task == nil {
		return fmt.Errorf("wake: sub-task %q not found", name)
	}
	for _, c := range task.Cmds {
		if err := r.Run(c); err != nil {
			return err
		}
	}
	return nil
}
