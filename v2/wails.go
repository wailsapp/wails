// Package wails is the main package of the Wails project.
// It is used by client applications.
package wails

import (
	"github.com/wailsapp/wails/v2/internal/app"
	"github.com/wailsapp/wails/v2/internal/runtime"
	"github.com/wailsapp/wails/v2/pkg/options"
)

// Runtime is an alias for the runtime.Runtime struct
type Runtime = runtime.Runtime

// CreateAppWithOptions creates an application based on the given config
func CreateAppWithOptions(options *options.App) *app.App {
	return app.CreateApp(options)
}

// CreateApp creates an application based on the given title, width and height
func CreateApp(title string, width int, height int) *app.App {

	options := &options.App{
		Title:  title,
		Width:  width,
		Height: height,
	}

	return app.CreateApp(options)
}
