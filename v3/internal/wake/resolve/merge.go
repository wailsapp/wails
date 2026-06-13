package resolve

import (
	"github.com/wailsapp/wails/v3/internal/wake/ast"
)

// MergeTaskfile layers a local taskfile over a base one with local-wins
// precedence: a task defined in local replaces/overrides the base task of the
// same name (see MergeTask), a task that exists only in local is added, and
// top-level vars defined in local override base vars per key. This is what
// powers Taskfile.local.yml overriding a base Taskfile.yml.
func MergeTaskfile(base, local *ast.Taskfile) {
	if local == nil {
		return
	}
	if base.Tasks == nil {
		base.Tasks = make(map[string]*ast.Task, len(local.Tasks))
	}
	for name, localTask := range local.Tasks {
		if baseTask, ok := base.Tasks[name]; ok {
			base.Tasks[name] = MergeTask(baseTask, localTask)
		} else {
			base.Tasks[name] = localTask
		}
	}
	if len(local.Vars) > 0 {
		base.Vars = mergeVarMaps(base.Vars, local.Vars)
	}
}

// MergeTask overlays override onto base with local-wins precedence:
//   - list fields (cmds, deps, sources, generates, platforms, status,
//     preconditions, aliases) REPLACE the base list whenever override provides
//     a non-empty one, otherwise the base list is kept;
//   - map fields (vars, env) merge per key with override winning;
//   - scalar fields take the override value when it is set.
//
// Replacement (not append) is deliberate: redefining a task in
// Taskfile.local.yml should supersede the base definition, not extend it.
func MergeTask(base, override *ast.Task) *ast.Task {
	result := *base
	if override.Dir != "" {
		result.Dir = override.Dir
	}
	if override.Summary != "" {
		result.Summary = override.Summary
	}
	if override.Desc != "" {
		result.Desc = override.Desc
	}
	if override.Label != "" {
		result.Label = override.Label
	}
	// Boolean fields use OR semantics: if either base or override has the
	// flag set, the merged task carries it. This is a deliberate trade-off,
	// not a bug: ast.Task models these as plain bool (not *bool), so the
	// parser cannot tell "unset" apart from "explicitly false" — an
	// override file declaring `silent: false` is indistinguishable from one
	// that omits the field entirely. OR-merging means the more conservative
	// behaviour wins ("silent if anyone in the chain wants it silent"),
	// which is the safer default for a redefinition layered over a base.
	//
	// If we later want overrides to be able to flip a base's true back to
	// false, this needs to become *bool with a parse-time presence check;
	// none of the wake-routed code paths currently read these fields, so
	// the limitation is undocumented user-visible behaviour rather than a
	// runtime issue.
	result.Silent = override.Silent || base.Silent
	result.Internal = override.Internal || base.Internal
	result.Interactive = override.Interactive || base.Interactive
	if override.Prompt != "" {
		result.Prompt = override.Prompt
	}
	if override.Method != "" {
		result.Method = override.Method
	}

	if len(override.Cmds) > 0 {
		result.Cmds = override.Cmds
	}
	if len(override.Deps) > 0 {
		result.Deps = override.Deps
	}
	if len(override.Sources) > 0 {
		result.Sources = override.Sources
	}
	if len(override.Generates) > 0 {
		result.Generates = override.Generates
	}
	if len(override.Platforms) > 0 {
		result.Platforms = override.Platforms
	}
	if len(override.Status) > 0 {
		result.Status = override.Status
	}
	if len(override.Precondition) > 0 {
		result.Precondition = override.Precondition
	}
	if len(override.Aliases) > 0 {
		result.Aliases = override.Aliases
	}

	result.Vars = mergeVarMaps(base.Vars, override.Vars)
	result.Env = mergeStringMaps(base.Env, override.Env)

	return &result
}

func mergeVarMaps(base, override map[string]*ast.Var) map[string]*ast.Var {
	result := make(map[string]*ast.Var)
	for k, v := range base {
		result[k] = v
	}
	for k, v := range override {
		result[k] = v
	}
	return result
}

func mergeStringMaps(base, override map[string]string) map[string]string {
	result := make(map[string]string)
	for k, v := range base {
		result[k] = v
	}
	for k, v := range override {
		result[k] = v
	}
	return result
}
