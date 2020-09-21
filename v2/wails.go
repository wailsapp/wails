// Package wails is the main package of the Wails project.
// It is used by client applications.
package wails

import (
	"github.com/wailsapp/wails/v2/internal/app"
	"github.com/wailsapp/wails/v2/internal/appoptions"
	"github.com/wailsapp/wails/v2/internal/runtime/goruntime"
)

// Runtime is an alias for the goruntime.Runtime struct
type Runtime = goruntime.Runtime

// Options is an alias for the app.Options struct
type Options = appoptions.Options

// MacOptions are Mac specific options
type MacOptions = appoptions.MacOptions

// CreateAppWithOptions creates an application based on the given config
func CreateAppWithOptions(options *Options) *app.App {
	return app.CreateApp(options)
}

// CreateApp creates an application based on the given title, width and height
func CreateApp(title string, width int, height int) *app.App {

	options := &Options{
		Title:  title,
		Width:  width,
		Height: height,
	}

	return app.CreateApp(options)
}
