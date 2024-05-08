package render

import (
	"go/types"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/parser/collect"
	"golang.org/x/tools/go/types/typeutil"
)

// module gathers data that is used when rendering a single JS/TS module.
type module struct {
	*Renderer
	*flags.GenerateBindingsOptions

	Imports *collect.ImportMap

	postponedCreates typeMapWorkaround
}

// Runtime returns the import path for the Wails JS runtime module.
func (m *module) Runtime() string {
	if m.UseBundledRuntime {
		return "/wails/runtime.js"
	} else {
		return "@wailsio/runtime"
	}
}

// typeutil.Map has a bug where struct fields with alias types result in panics.
// typeMapWorkaround is a drop-in workaround that offers the interface we need.
// It should be replaced by typeutil.Map as soon as the bug is fixed.
type typeMapWorkaround struct {
	main     typeutil.Map
	fallback map[types.Type]any
}

func (cm *typeMapWorkaround) Len() int {
	return cm.main.Len() + len(cm.fallback)
}

func (cm *typeMapWorkaround) At(typ types.Type) (result any) {
	defer func() {
		if recover() != nil {
			result = cm.fallback[typ]
		}
	}()

	return cm.main.At(typ)
}

func (cm *typeMapWorkaround) Set(typ types.Type, value any) {
	defer func() {
		if recover() != nil {
			if cm.fallback == nil {
				cm.fallback = make(map[types.Type]any)
			}
			cm.fallback[typ] = value
		}
	}()

	cm.main.Set(typ, value)
}

func (cm *typeMapWorkaround) Iterate(yield func(key types.Type, value any)) {
	cm.main.Iterate(yield)
	for key, value := range cm.fallback {
		yield(key, value)
	}
}
