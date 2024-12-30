package main

import (
	"embed"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/index.html
var assets embed.FS

type App struct{}

func (a *App) GetCurrentInstanceInfo() map[string]interface{} {
	return map[string]interface{}{
		"args":       os.Args,
		"workingDir": getCurrentWorkingDir(),
	}
}

func main() {

	var window *application.WebviewWindow
	app := application.New(application.Options{
		Name:        "Single Instance Example",
		LogLevel:    slog.LevelDebug,
		Description: "An example of single instance functionality in Wails v3",
		Services: []application.Service{
			application.NewService(&App{}),
		},
		SingleInstance: &application.SingleInstanceOptions{
			UniqueID: "com.wails.example.single-instance",
			OnSecondInstanceLaunch: func(data application.SecondInstanceData) {
				if window != nil {
					window.EmitEvent("secondInstanceLaunched", data)
					window.Restore()
					window.Focus()
				}
				log.Printf("Second instance launched with args: %v\n", data.Args)
				log.Printf("Working directory: %s\n", data.WorkingDir)
				if data.AdditionalData != nil {
					log.Printf("Additional data: %v\n", data.AdditionalData)
				}
			},
			AdditionalData: map[string]string{
				"launchtime": time.Now().Local().String(),
			},
		},
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
	})

	window = app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:  "Single Instance Demo",
		Width:  800,
		Height: 700,
		URL:    "/",
	})

	app.Run()
}

func getCurrentWorkingDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}
