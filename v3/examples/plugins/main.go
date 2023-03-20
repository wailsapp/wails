package main

import (
	"embed"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/plugins/browser"
	"github.com/wailsapp/wails/v3/plugins/kvstore"
	"log"
	"plugin_demo/plugins/hashes"
	//"plugin_demo/plugins/hashes"
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
			"hashes":  hashes.NewPlugin(),
			"browser": browser.NewPlugin(),
			"kvstore": kvstore.NewPlugin(&kvstore.Config{
				Filename: "store.json",
				AutoSave: true,
			}),
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
