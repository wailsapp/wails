package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:     "My Test App",
		Width:     1024,
		Height:    768,
		MinWidth:  400,
		MinHeight: 300,
		Frameless: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind: []interface{}{
			app,
			&GreetService{},
		},
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId: "com.test.myapp",
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
