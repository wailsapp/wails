package main

import (
	"embed"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/*
var assets embed.FS

func main() {

	app := application.New(application.Options{
		Name:        "Screen Demo",
		Description: "A demo of the Screen API",
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		Windows: application.WindowsOptions{
			WndProcInterceptor:            nil,
			DisableQuitOnLastWindowClosed: false,
			WebviewUserDataPath:           "",
			WebviewBrowserPath:            "",
		},
		Services: []application.Service{
			application.NewService(&ScreenService{}),
		},
		LogLevel: slog.LevelError,
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
			Middleware: func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Disable caching
					w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
					w.Header().Set("Pragma", "no-cache")
					w.Header().Set("Expires", "0")

					_, filename, _, _ := runtime.Caller(0)
					dir := filepath.Dir(filename)
					assetsDir := filepath.Join(dir, "assets")

					// Clean and validate the path to prevent directory traversal
					cleanPath := filepath.Clean(r.URL.Path)
					fullPath := filepath.Join(assetsDir, cleanPath)

					// Ensure the resolved path is still within the assets directory
					if !strings.HasPrefix(fullPath, assetsDir+string(filepath.Separator)) && fullPath != assetsDir {
						// Path traversal attempt detected, fall back to default handler
						next.ServeHTTP(w, r)
						return
					}

					if _, err := os.Stat(fullPath); err == nil {
						// Serve file from disk to make testing easy
						http.ServeFile(w, r, fullPath)
					} else {
						// Passthrough to the default asset handler if file not found on disk
						next.ServeHTTP(w, r)
					}
				})
			},
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Screen Demo",
		Width:  800,
		Height: 600,
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInsetUnified,
			InvisibleTitleBarHeight: 50,
		},
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
