package resolve

import (
	"github.com/wailsapp/wails/v3/internal/wake/ast"
	"github.com/wailsapp/wails/v3/internal/wake/platform"
)

func FilterPlatforms(tf *ast.Taskfile) {
	for _, task := range tf.Tasks {
		if !platform.Filter(task.Platforms) {
			task.Cmds = nil
			task.Deps = nil
		}
	}
}
