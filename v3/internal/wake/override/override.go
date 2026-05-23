package override

import (
	"github.com/wailsapp/wails/v3/internal/wake/ast"
	"github.com/wailsapp/wails/v3/internal/wake/parse"
	"github.com/wailsapp/wails/v3/internal/wake/resolve"
)

var (
	namedOverrides = make(map[string]func(*ast.Task) *ast.Task)
	allOverrides   []func(*ast.Task) *ast.Task
)

func Override(taskName string, fn func(*ast.Task) *ast.Task) {
	namedOverrides[taskName] = fn
}

func OverrideAll(fn func(*ast.Task) *ast.Task) {
	allOverrides = append(allOverrides, fn)
}

func ClearOverrides() {
	namedOverrides = make(map[string]func(*ast.Task) *ast.Task)
	allOverrides = nil
}

func Named() map[string]func(*ast.Task) *ast.Task {
	return namedOverrides
}

func All() []func(*ast.Task) *ast.Task {
	return allOverrides
}

func LoadTaskfileOverride(path string) error {
	ov, err := parse.Parse(path)
	if err != nil {
		return err
	}
	if err := parse.ResolveIncludes(ov); err != nil {
		return err
	}

	for name, ovTask := range ov.Tasks {
		name := name
		ovTask := ovTask
		Override(name, func(base *ast.Task) *ast.Task {
			return resolve.MergeTask(base, ovTask)
		})
	}

	return nil
}

func LoadTaskfileOverrides(dir string, patterns ...string) error {
	if len(patterns) == 0 {
		patterns = []string{
			"Taskfile.local.yml",
			"Taskfile.local.yaml",
			"Taskfile.override.yml",
			"Taskfile.override.yaml",
		}
	}

	for _, pattern := range patterns {
		path := dir + "/" + pattern
		if _, err := parse.Parse(path); err == nil {
			if err := LoadTaskfileOverride(path); err != nil {
				return err
			}
		}
	}

	return nil
}
