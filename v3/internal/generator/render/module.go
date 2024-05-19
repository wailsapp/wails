package render

import (
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/generator/collect"
	"golang.org/x/tools/go/types/typeutil"
)

// module gathers data that is used when rendering a single JS/TS module.
type module struct {
	*Renderer
	*flags.GenerateBindingsOptions

	Imports *collect.ImportMap

	postponedCreates typeutil.Map
}

// Runtime returns the import path for the Wails JS runtime module.
func (m *module) Runtime() string {
	if m.UseBundledRuntime {
		return "/wails/runtime.js"
	} else {
		return "@wailsio/runtime"
	}
}
