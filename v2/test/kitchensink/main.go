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
		MaxWidth:  1400,
		MaxHeight: 800,
		Mac: &mac.Options{
			WebviewIsTransparent:          true,
			WindowBackgroundIsTranslucent: true,
			TitleBar:                      mac.TitleBarHiddenInset(),
		},
		LogLevel: logger.INFO,
	})

	app.Bind(&Logger{})

	app.Run()
}
