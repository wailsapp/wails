package main

import (
	"embed"
	"log/slog"
	"os"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/services/keygen"
)

//go:embed embed/*
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "Keygen Integration Demo",
		Description: "A demo application showing Keygen licensing and update features",
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		Windows: application.WindowsOptions{
			WebviewUserDataPath: "keygen-demo",
		},
		Linux: application.LinuxOptions{
			WindowIsTranslucent: false,
		},
		LogLevel: slog.LevelDebug,
		Services: []application.Service{
			// Initialize Keygen service with demo account
			application.NewService(keygen.New(keygen.ServiceOptions{
				AccountID:      "demo",      // Replace with your Keygen account ID
				ProductID:      "prod_demo", // Replace with your product ID
				LicenseKey:     "",          // Will be set via UI
				PublicKey:      "",          // Replace with your Ed25519 public key
				CurrentVersion: "1.0.0",
				AutoCheck:      true,
				UpdateChannel:  "stable",
			})),
		},
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
	})

	// Create the application window
	window := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:  "Keygen Integration Demo",
		Width:  1024,
		Height: 768,
		URL:    "/",
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarDefault,
			InvisibleTitleBarHeight: 40,
		},
		Windows: application.WindowsWindow{
			DisableWindowIcon: false,
		},
		BackgroundColour: application.NewRGBA(255, 255, 255, 255),
	})

	// Add the App binding for frontend access
	window.RegisterHook(application.HookWindowEvent("ready"), func(event *application.WindowEvent) {
		appInstance := &App{
			window: window,
			keygen: keygen.New(keygen.ServiceOptions{
				AccountID:      "demo",
				ProductID:      "prod_demo",
				LicenseKey:     "",
				PublicKey:      "",
				CurrentVersion: "1.0.0",
				AutoCheck:      true,
				UpdateChannel:  "stable",
			}),
		}
		window.SetBindings(appInstance)
	})

	// Run the application
	err := app.Run()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
