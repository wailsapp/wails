package main

import (
	wails "github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

func main() {

	// Create application with options
	app := wails.CreateAppWithOptions(&options.App{
		Title:     "Kitchen Sink",
		Width:     1024,
		Height:    768,
		MinWidth:  800,
		MinHeight: 600,
		Mac: &mac.Options{
			WebviewIsTransparent:          true,
			WindowBackgroundIsTranslucent: true,
			// Comment out line below to see Window.SetTitle() work
			TitleBar: mac.TitleBarHiddenInset(),
			Menu:     createApplicationMenu(),
		},
		LogLevel: logger.TRACE,
	})

	app.Bind(&Events{})
	app.Bind(&Logger{})
	app.Bind(&Browser{})
	app.Bind(&System{})
	app.Bind(&Dialog{})
	app.Bind(&Window{})

	app.Run()
}
