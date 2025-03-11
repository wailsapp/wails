package main

import (
	"embed"
	"log/slog"
	"os"

	"gin-service/services"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/*
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "Gin Service Demo",
		Description: "A demo of using Gin in Wails services",
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		LogLevel: slog.LevelDebug,
		Services: []application.Service{
			application.NewServiceWithOptions(services.NewGinService(), application.ServiceOptions{
				Route: "/api",
			}),
		},
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
	})

	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:  "Gin Service Demo",
		Width:  1024,
		Height: 768,
	})

	err := app.Run()

	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
