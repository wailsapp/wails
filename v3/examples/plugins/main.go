package main

import (
	"embed"
	_ "embed"
	"github.com/wailsapp/wails/v3/pkg/application"
	"log"
	"plugin_demo/plugins/hashes"
)

//go:embed assets/*
var assets embed.FS

func main() {

	app := application.New(application.Options{
		Name:        "Plugin Demo",
		Description: "A demo of the plugins API",
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		Plugins: map[string]application.Plugin{
			"hashes": hashes.NewPlugin(),
		},
		Assets: application.AssetOptions{
			FS: assets,
		},
	})

	window := app.NewWebviewWindow()
	window.ToggleDevTools()

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
