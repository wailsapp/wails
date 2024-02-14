package main

import (
	"embed"
	"os"
	"plugin_demo/plugins/hashes"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/plugins/experimental/server"
	"github.com/wailsapp/wails/v3/plugins/kvstore"
	"github.com/wailsapp/wails/v3/plugins/log"
	"github.com/wailsapp/wails/v3/plugins/single_instance"
	"github.com/wailsapp/wails/v3/plugins/sqlite"
	"github.com/wailsapp/wails/v3/plugins/start_at_login"
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
			"log":    log.NewPlugin(),
			"sqlite": sqlite.NewPlugin(&sqlite.Config{
				DBFile: "test.db",
			}),
			"kvstore": kvstore.NewPlugin(&kvstore.Config{
				Filename: "store.json",
				AutoSave: true,
			}),
			"server": server.NewPlugin(&server.Config{
				Enabled: true,
				Port:    34115,
			}),
			"single_instance": single_instance.NewPlugin(&single_instance.Config{
				// When true, the original app will be activated when a second instance is launched
				ActivateAppOnSubsequentLaunch: true,
			}),
			"start_at_login": start_at_login.NewPlugin(start_at_login.Config{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
	})

	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		DevToolsEnabled:        true,
		OpenInspectorOnStartup: true,
	})

	err := app.Run()

	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
