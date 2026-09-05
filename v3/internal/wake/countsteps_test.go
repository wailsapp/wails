package wake

import (
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
)

// countSteps must count only tasks that run a real command, following both deps
// and task-ref commands, deduping shared deps.
func TestCountSteps(t *testing.T) {
	tf := &ast.Taskfile{
		Tasks: map[string]*ast.Task{
			// Pure wrapper: dispatches via a task-ref cmd, no real command.
			"build": {
				Name: "build",
				Cmds: []*ast.Cmd{{Task: "darwin:build"}},
			},
			// Real command + deps.
			"darwin:build": {
				Name: "darwin:build",
				Deps: []*ast.Dep{{Task: "tidy"}, {Task: "frontend"}},
				Cmds: []*ast.Cmd{{Cmd: "go build"}},
			},
			"tidy":     {Name: "tidy", Cmds: []*ast.Cmd{{Cmd: "go mod tidy"}}},
			"frontend": {Name: "frontend", Deps: []*ast.Dep{{Task: "tidy"}}, Cmds: []*ast.Cmd{{Cmd: "npm run build"}}},
			// Reachable only via a task-ref cmd, has a real command.
			"bundle": {Name: "bundle", Cmds: []*ast.Cmd{{Cmd: "zip out"}}},
		},
	}
	// Make darwin:build also dispatch to bundle via a task-ref cmd.
	tf.Tasks["darwin:build"].Cmds = append(tf.Tasks["darwin:build"].Cmds, &ast.Cmd{Task: "bundle"})

	// From "build" (a pure wrapper): darwin:build, tidy, frontend, bundle = 4
	// real-command tasks; "build" itself is not counted.
	if got := countSteps(tf, "build"); got != 4 {
		t.Fatalf("countSteps(build) = %d, want 4", got)
	}

	// tidy is shared by darwin:build and frontend but counted once.
	if got := countSteps(tf, "frontend"); got != 2 {
		t.Fatalf("countSteps(frontend) = %d, want 2 (frontend, tidy)", got)
	}

	if got := countSteps(tf, "tidy"); got != 1 {
		t.Fatalf("countSteps(tidy) = %d, want 1", got)
	}
}
