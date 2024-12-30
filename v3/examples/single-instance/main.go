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

var encryptionKey = [32]byte{
	0x1e, 0x1f, 0x1c, 0x1d, 0x1a, 0x1b, 0x18, 0x19,
	0x16, 0x17, 0x14, 0x15, 0x12, 0x13, 0x10, 0x11,
	0x0e, 0x0f, 0x0c, 0x0d, 0x0a, 0x0b, 0x08, 0x09,
	0x06, 0x07, 0x04, 0x05, 0x02, 0x03, 0x00, 0x01,
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
			UniqueID:      "com.wails.example.single-instance",
			EncryptionKey: encryptionKey,
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
