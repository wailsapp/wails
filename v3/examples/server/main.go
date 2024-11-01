package main

import (
	"embed"
	"io/fs"
	"os"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/plugins/experimental/server"
)

//go:embed assets/*
var assets embed.FS

func main() {
	staticAssets, _ := fs.Sub(assets, "assets")
	app := application.New(application.Options{
		Name:        "Server Demo",
		Description: "server only demo of the plugins API",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(staticAssets),
		},
		Services: []application.Service{
			application.NewService(&GreetService{}),
			application.NewService(server.NewPlugin(&server.Config{
				Host:    "0.0.0.0",
				Port:    34115,
				Enabled: true,
				Assets:  staticAssets,
			})),
		},
	})

	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			app.EmitEvent("time", now)
			time.Sleep(time.Second)
		}
	}()

	// This is to test desktop vs server
	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		URL:             "/",
		DevToolsEnabled: true,
	})

	err := app.Run()

	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
