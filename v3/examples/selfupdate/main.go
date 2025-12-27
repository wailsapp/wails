package main

import (
	"embed"
	"log/slog"
	"os"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/services/selfupdate"
)

//go:embed assets/*
var assets embed.FS

// version is set at build time using -ldflags "-X main.version=1.0.0"
var version = "0.0.1"

func main() {
	app := application.New(application.Options{
		Name:        "Selfupdate Demo",
		Description: "A demo of the selfupdate service",
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		LogLevel: slog.LevelDebug,
		Services: []application.Service{
			application.NewService(selfupdate.NewWithConfig(&selfupdate.Config{
				// Set the current version of your application
				CurrentVersion: version,
				// Use GitHub as the update source (default)
				Source: selfupdate.SourceGitHub,
				// The repository to check for updates (owner/repo format)
				// Change this to your own repository
				Repository: "example/myapp",
				// Optional: Allow pre-release versions
				AllowPrerelease: false,
				// Optional: Configure signature verification
				// Signature: &selfupdate.SignatureConfig{
				// 	Type:      selfupdate.SignatureECDSA,
				// 	PublicKey: string(publicKeyPEM),
				// },
				// Optional: Configure UI theming
				UI: &selfupdate.UIConfig{
					Title:            "Update Available",
					BackgroundColor:  "#1a1a2e",
					TextColor:        "#eaeaea",
					AccentColor:      "#0f3460",
					ProgressBarColor: "#e94560",
					ButtonColor:      "#e94560",
					ButtonTextColor:  "#ffffff",
				},
			})),
		},
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Selfupdate Demo - v" + version,
		Width:  900,
		Height: 600,
	})

	err := app.Run()

	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
