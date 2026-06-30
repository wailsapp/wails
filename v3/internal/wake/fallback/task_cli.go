// Package fallback drives the embedded Task runtime as a backstop when wake
// encounters a Taskfile feature it doesn't implement itself (dotenv, requires,
// interval, run modes other than "always", short, defer, non-interleaved
// output). The Task runtime is compiled into wails3 — there is no external
// `task` binary to look for — so this is an in-process call, not a subprocess.
package fallback

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/wailsapp/task/v3"
	"github.com/wailsapp/task/v3/taskfile/ast"
)

// ErrUnsupported is the sentinel a wake caller may wrap when bailing out to the
// Task runtime. Kept for backwards compatibility with the previous package
// surface; nothing in the new fallback path returns it directly.
type ErrUnsupported struct {
	Feature string
}

func (e *ErrUnsupported) Error() string {
	return fmt.Sprintf("unsupported feature: %s", e.Feature)
}

// TaskCLI runs the named task using the in-process Task runtime, applying any
// `KEY=VALUE` strings in env as Taskfile-level CLI variables (matching the
// semantics of `task NAME KEY=VALUE`). dir is the working directory; the
// empty string falls through to the Taskfile's own root_path / discovery.
//
// Naming carries the historical "CLI" suffix because the function used to
// shell out to an external `task` binary. It hasn't done that since this
// commit — the Task runtime is embedded.
func TaskCLI(name, dir string, env []string) error {
	e := task.Executor{
		Dir:                 dir,
		Color:               true,
		DisableVersionCheck: true,
		Stdin:               os.Stdin,
		Stdout:              os.Stdout,
		Stderr:              os.Stderr,
	}
	if err := e.Setup(); err != nil {
		return fmt.Errorf("wake fallback: task setup: %w", err)
	}

	call := &ast.Call{
		Task: name,
		Vars: &ast.Vars{},
	}

	// Pre-existing wake call sites pass env as a slice of "KEY=VALUE" strings.
	// In Task's model those land as Taskfile-level vars and as call-level vars
	// — set both so commands that read either path see the same values.
	if e.Taskfile != nil && e.Taskfile.Vars == nil {
		e.Taskfile.Vars = &ast.Vars{}
	}
	for _, kv := range env {
		k, v, ok := strings.Cut(kv, "=")
		if !ok {
			continue
		}
		val := ast.Var{Value: v}
		call.Vars.Set(k, val)
		if e.Taskfile != nil {
			e.Taskfile.Vars.Set(k, val)
		}
	}

	// Match commands.RunTask's "no-watch" defaults.
	e.Watch = false
	e.Interval = 0 * time.Second

	return e.RunTask(context.Background(), call)
}

// Available used to report whether an external `task` binary was on PATH. The
// embedded runtime is always available; keeping the function so existing
// callers compile without churn.
func Available() bool { return true }
