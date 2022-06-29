// Package wails is the main package of the Wails project.
// It is used by client applications.
package wails

import (
	app "github.com/wailsapp/wails/v2/internal/appng"
	"github.com/wailsapp/wails/v2/internal/signal"
	"github.com/wailsapp/wails/v2/pkg/options"
)

// Run creates an application based on the given config and executes it
func Run(options *options.App) error {

	if options.RGBA != nil {
		println("---- WARNING ----")
		println("The `RGBA` option has been deprecated. Please use `BackgroundColour`.")

		if options.BackgroundColour == nil {
			options.BackgroundColour = options.RGBA
		}
	}

	// Call an Init method manually
	err := Init()
	if err != nil {
		return err
	}

	mainapp, err := app.CreateApp(options)
	if err != nil {
		return err
	}

	signal.OnShutdown(func() {
		mainapp.Shutdown()
	})

	signal.Start()

	return mainapp.Run()
}
