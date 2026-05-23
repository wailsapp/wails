package resolve

import (
	"github.com/wailsapp/wails/v3/internal/wake/ast"
)

func ApplyOverrides(tf *ast.Taskfile, named map[string]func(*ast.Task) *ast.Task, all []func(*ast.Task) *ast.Task) {
	for name, task := range tf.Tasks {
		if fn, ok := named[name]; ok {
			tf.Tasks[name] = fn(task)
		}
		for _, fn := range all {
			tf.Tasks[name] = fn(tf.Tasks[name])
		}
	}
}

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
	result.Silent = override.Silent || base.Silent
	result.Internal = override.Internal || base.Internal
	result.Interactive = override.Interactive || base.Interactive
	if override.Prompt != "" {
		result.Prompt = override.Prompt
	}
	if override.Method != "" {
		result.Method = override.Method
	}

	result.Cmds = append(base.Cmds, override.Cmds...)
	result.Deps = append(base.Deps, override.Deps...)
	result.Precondition = append(override.Precondition, base.Precondition...)

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
