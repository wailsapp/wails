package main

import (
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

func main() {

	// Create application with options
	app := NewBasic()

	err := wails.Run(&options.App{
		Title:     "{{.ProjectName}}",
		Width:     800,
		Height:    600,
		MinWidth:  800,
		MinHeight: 600,
		Mac: &mac.Options{
			WebviewIsTransparent:          true,
			WindowBackgroundIsTranslucent: true,
			TitleBar:                      mac.TitleBarHiddenInset(),
			Menu:                          menu.DefaultMacMenu(),
		},
		LogLevel: logger.DEBUG,
		Startup:  app.startup,
		Shutdown: app.shutdown,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
