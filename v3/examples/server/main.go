package main

import (
	"embed"
	"os"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/plugins/experimental/server"
	"github.com/wailsapp/wails/v3/plugins/log"
)

//go:embed assets/*
var assets embed.FS

func main() {

	app := application.New(application.Options{
		Name:        "Server Demo",
		Description: "server only demo of the plugins API",
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		Plugins: map[string]application.Plugin{
			"log": log.NewPlugin(),
			"server": server.NewPlugin(&server.Config{
				Host:    "0.0.0.0",
				Port:    34115,
				Enabled: true,
			}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
	})
	go func() {
		for {
			app.Events.Emit(&application.WailsEvent{
				Name: "ping",
				Data: "are you alive?",
			})
			time.Sleep(10 * time.Second)
		}
	}()

	//	window := app.NewWebviewWindow()
	//	window.ToggleDevTools()

	err := app.Run()

	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
