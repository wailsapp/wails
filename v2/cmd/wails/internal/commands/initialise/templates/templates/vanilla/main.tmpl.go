package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed frontend/src
var assets embed.FS

func main() {

	// Create application with options
	app := NewApp()

	err := wails.Run(&options.App{
		Title:             "{{.ProjectName}}",
		Width:             800,
		Height:            600,
		MinWidth:          400,
		MinHeight:         400,
		MaxWidth:          1280,
		MaxHeight:         1024,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		RGBA:              &options.RGBA{0, 0, 0, 255},
		Assets:            assets,
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
		LogLevel:   logger.DEBUG,
		OnStartup:  app.startup,
		OnShutdown: app.shutdown,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
