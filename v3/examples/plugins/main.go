package main

import (
	"embed"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/plugins/browser"
	"github.com/wailsapp/wails/v3/plugins/kvstore"
	"github.com/wailsapp/wails/v3/plugins/log"
	"github.com/wailsapp/wails/v3/plugins/single_instance"
	"github.com/wailsapp/wails/v3/plugins/sqlite"
	"github.com/wailsapp/wails/v3/plugins/start_at_login"
	"os"
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
			"hashes":  hashes.NewPlugin(),
			"browser": browser.NewPlugin(),
			"log":     log.NewPlugin(),
			"sqlite": sqlite.NewPlugin(&sqlite.Config{
				DBFile: "test.db",
			}),
			"kvstore": kvstore.NewPlugin(&kvstore.Config{
				Filename: "store.json",
				AutoSave: true,
			}),
			"single_instance": single_instance.NewPlugin(&single_instance.Config{
				// When true, the original app will be activated when a second instance is launched
				ActivateAppOnSubsequentLaunch: true,
			}),
			"start_at_login": start_at_login.NewPlugin(),
		},
		Assets: application.AssetOptions{
			FS: assets,
		},
	})

	window := app.NewWebviewWindow()
	window.ToggleDevTools()

	err := app.Run()

	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
