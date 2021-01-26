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

// Store is an alias for the Store object
type Store = runtime.Store

// Run creates an application based on the given config and executes it
func Run(options *options.App) error {
	app, err := app.CreateApp(options)
	if err != nil {
		return err
	}

	return app.Run()
}
