package resolve

import (
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
)

func TestMergeTaskReplacesCmds(t *testing.T) {
	base := &ast.Task{
		Name: "build",
		Cmds: []*ast.Cmd{{Cmd: "echo base"}},
	}
	override := &ast.Task{
		Name: "build",
		Cmds: []*ast.Cmd{{Cmd: "echo local"}},
	}

	got := MergeTask(base, override)
	if len(got.Cmds) != 1 {
		t.Fatalf("cmds = %d, want 1 (replace, not append)", len(got.Cmds))
	}
	if got.Cmds[0].Cmd != "echo local" {
		t.Errorf("cmd = %q, want echo local", got.Cmds[0].Cmd)
	}
}

func TestMergeTaskKeepsBaseListWhenOverrideEmpty(t *testing.T) {
	base := &ast.Task{
		Name: "build",
		Cmds: []*ast.Cmd{{Cmd: "echo base"}},
		Deps: []*ast.Dep{{Task: "setup"}},
	}
	override := &ast.Task{
		Name: "build",
		Env:  map[string]string{"FOO": "bar"},
	}

	got := MergeTask(base, override)
	if len(got.Cmds) != 1 || got.Cmds[0].Cmd != "echo base" {
		t.Errorf("cmds = %v, want base cmds preserved", got.Cmds)
	}
	if len(got.Deps) != 1 || got.Deps[0].Task != "setup" {
		t.Errorf("deps = %v, want base deps preserved", got.Deps)
	}
	if got.Env["FOO"] != "bar" {
		t.Errorf("env FOO = %q, want bar", got.Env["FOO"])
	}
}

func TestMergeTaskMergesEnvPerKey(t *testing.T) {
	base := &ast.Task{
		Name: "build",
		Env:  map[string]string{"KEEP": "yes", "OVERRIDE": "old"},
	}
	override := &ast.Task{
		Name: "build",
		Env:  map[string]string{"OVERRIDE": "new", "ADD": "added"},
	}

	got := MergeTask(base, override)
	if got.Env["KEEP"] != "yes" {
		t.Errorf("KEEP = %q, want yes", got.Env["KEEP"])
	}
	if got.Env["OVERRIDE"] != "new" {
		t.Errorf("OVERRIDE = %q, want new (local wins)", got.Env["OVERRIDE"])
	}
	if got.Env["ADD"] != "added" {
		t.Errorf("ADD = %q, want added", got.Env["ADD"])
	}
}

func TestMergeTaskfileAddsAndOverrides(t *testing.T) {
	base := &ast.Taskfile{
		Vars: map[string]*ast.Var{
			"GREETING": {Static: "hello"},
			"KEEP":     {Static: "base"},
		},
		Tasks: map[string]*ast.Task{
			"build": {Name: "build", Cmds: []*ast.Cmd{{Cmd: "echo base-build"}}},
			"test":  {Name: "test", Cmds: []*ast.Cmd{{Cmd: "echo base-test"}}},
		},
	}
	local := &ast.Taskfile{
		Vars: map[string]*ast.Var{
			"GREETING": {Static: "hola"},
		},
		Tasks: map[string]*ast.Task{
			"build":  {Name: "build", Cmds: []*ast.Cmd{{Cmd: "echo local-build"}}},
			"deploy": {Name: "deploy", Cmds: []*ast.Cmd{{Cmd: "echo deploy"}}},
		},
	}

	MergeTaskfile(base, local)

	// Overridden task: local cmds replace base cmds.
	if base.Tasks["build"].Cmds[0].Cmd != "echo local-build" {
		t.Errorf("build cmd = %q, want echo local-build", base.Tasks["build"].Cmds[0].Cmd)
	}
	// Untouched base task survives.
	if base.Tasks["test"].Cmds[0].Cmd != "echo base-test" {
		t.Errorf("test cmd = %q, want echo base-test", base.Tasks["test"].Cmds[0].Cmd)
	}
	// Local-only task is added.
	if _, ok := base.Tasks["deploy"]; !ok {
		t.Error("deploy task not added from local")
	}
	// Vars: local wins per key, base-only var preserved.
	if base.Vars["GREETING"].Static != "hola" {
		t.Errorf("GREETING = %q, want hola (local wins)", base.Vars["GREETING"].Static)
	}
	if base.Vars["KEEP"].Static != "base" {
		t.Errorf("KEEP = %q, want base (preserved)", base.Vars["KEEP"].Static)
	}
}

func fullBaseTask() *ast.Task {
	return &ast.Task{
		Name:         "build",
		Dir:          "/base",
		Summary:      "base summary",
		Desc:         "base desc",
		Label:        "base label",
		Prompt:       "base prompt",
		Method:       "timestamp",
		Cmds:         []*ast.Cmd{{Cmd: "base cmd"}},
		Deps:         []*ast.Dep{{Task: "base-dep"}},
		Sources:      []string{"base/**/*.go"},
		Generates:    []string{"base/out"},
		Platforms:    []string{"linux"},
		Status:       []string{"test -f base"},
		Precondition: []*ast.Precondition{{Sh: "test -d base"}},
		Aliases:      []string{"b"},
		Vars:         map[string]*ast.Var{"V": {Static: "base"}},
		Env:          map[string]string{"E": "base"},
	}
}

// TestMergeTaskOverrideSetsEveryField exercises the "override provides it" side
// of every branch in MergeTask: each scalar and list field should take the
// override value.
func TestMergeTaskOverrideSetsEveryField(t *testing.T) {
	base := fullBaseTask()
	override := &ast.Task{
		Name:         "build",
		Dir:          "/over",
		Summary:      "over summary",
		Desc:         "over desc",
		Label:        "over label",
		Silent:       true,
		Internal:     true,
		Interactive:  true,
		Prompt:       "over prompt",
		Method:       "checksum",
		Cmds:         []*ast.Cmd{{Cmd: "over cmd"}},
		Deps:         []*ast.Dep{{Task: "over-dep"}},
		Sources:      []string{"over/**/*.go"},
		Generates:    []string{"over/out"},
		Platforms:    []string{"darwin"},
		Status:       []string{"test -f over"},
		Precondition: []*ast.Precondition{{Sh: "test -d over"}},
		Aliases:      []string{"o"},
		Vars:         map[string]*ast.Var{"V": {Static: "over"}},
		Env:          map[string]string{"E": "over"},
	}

	got := MergeTask(base, override)

	if got.Dir != "/over" || got.Summary != "over summary" || got.Desc != "over desc" || got.Label != "over label" {
		t.Errorf("scalar fields not overridden: %+v", got)
	}
	if !got.Silent || !got.Internal || !got.Interactive {
		t.Errorf("bool fields not OR'd true: silent=%v internal=%v interactive=%v", got.Silent, got.Internal, got.Interactive)
	}
	if got.Prompt != "over prompt" || got.Method != "checksum" {
		t.Errorf("prompt/method not overridden: %q %q", got.Prompt, got.Method)
	}
	if got.Cmds[0].Cmd != "over cmd" || got.Deps[0].Task != "over-dep" {
		t.Errorf("cmds/deps not replaced")
	}
	if got.Sources[0] != "over/**/*.go" || got.Generates[0] != "over/out" {
		t.Errorf("sources/generates not replaced")
	}
	if got.Platforms[0] != "darwin" || got.Status[0] != "test -f over" {
		t.Errorf("platforms/status not replaced")
	}
	if got.Precondition[0].Sh != "test -d over" || got.Aliases[0] != "o" {
		t.Errorf("precondition/aliases not replaced")
	}
	if got.Vars["V"].Static != "over" || got.Env["E"] != "over" {
		t.Errorf("vars/env not overridden")
	}
}

// TestMergeTaskEmptyOverrideKeepsBase exercises the "override omits it" side of
// every branch: an override that only names the task must preserve all base
// fields.
func TestMergeTaskEmptyOverrideKeepsBase(t *testing.T) {
	base := fullBaseTask()
	got := MergeTask(base, &ast.Task{Name: "build"})

	if got.Dir != "/base" || got.Summary != "base summary" || got.Desc != "base desc" || got.Label != "base label" {
		t.Errorf("scalar fields not preserved: %+v", got)
	}
	if got.Prompt != "base prompt" || got.Method != "timestamp" {
		t.Errorf("prompt/method not preserved: %q %q", got.Prompt, got.Method)
	}
	if len(got.Cmds) != 1 || got.Cmds[0].Cmd != "base cmd" {
		t.Errorf("cmds not preserved: %v", got.Cmds)
	}
	if len(got.Deps) != 1 || got.Deps[0].Task != "base-dep" {
		t.Errorf("deps not preserved: %v", got.Deps)
	}
	if len(got.Sources) != 1 || got.Sources[0] != "base/**/*.go" {
		t.Errorf("sources not preserved: %v", got.Sources)
	}
	if len(got.Generates) != 1 || got.Generates[0] != "base/out" {
		t.Errorf("generates not preserved: %v", got.Generates)
	}
	if len(got.Platforms) != 1 || got.Platforms[0] != "linux" {
		t.Errorf("platforms not preserved: %v", got.Platforms)
	}
	if len(got.Status) != 1 || got.Status[0] != "test -f base" {
		t.Errorf("status not preserved: %v", got.Status)
	}
	if len(got.Precondition) != 1 || got.Precondition[0].Sh != "test -d base" {
		t.Errorf("precondition not preserved: %v", got.Precondition)
	}
	if len(got.Aliases) != 1 || got.Aliases[0] != "b" {
		t.Errorf("aliases not preserved: %v", got.Aliases)
	}
	if got.Vars["V"].Static != "base" || got.Env["E"] != "base" {
		t.Errorf("vars/env not preserved")
	}
}

func TestMergeTaskfileNilLocalIsNoop(t *testing.T) {
	base := &ast.Taskfile{
		Tasks: map[string]*ast.Task{"build": {Name: "build"}},
	}
	MergeTaskfile(base, nil)
	if len(base.Tasks) != 1 {
		t.Errorf("tasks = %d, want 1", len(base.Tasks))
	}
}

// TestMergeTaskfileNilBaseTasks guards the case where the base declares no
// `tasks:` (nil map): adding a local-only task must not panic.
func TestMergeTaskfileNilBaseTasks(t *testing.T) {
	base := &ast.Taskfile{} // Tasks is nil
	local := &ast.Taskfile{
		Tasks: map[string]*ast.Task{"deploy": {Name: "deploy"}},
	}
	MergeTaskfile(base, local)
	if _, ok := base.Tasks["deploy"]; !ok {
		t.Error("deploy not added to a nil-Tasks base")
	}
}
