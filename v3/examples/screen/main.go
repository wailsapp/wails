package main

import (
	"embed"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path"
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

					// Resolve the assets directory to an absolute, cleaned path.
					assetsDirAbs, err := filepath.Abs(filepath.Clean(assetsDir))
					if err != nil {
						// If we cannot resolve the assets directory safely, fall back to default handler.
						next.ServeHTTP(w, r)
						return
					}

					// Clean the requested URL path using path.Clean (HTTP paths always use forward slashes).
					cleanPath := path.Clean("/" + r.URL.Path)

					// Reject Windows drive-letter (e.g. "C:/...") or UNC-style absolute paths.
					if len(cleanPath) >= 2 && cleanPath[1] == ':' {
						next.ServeHTTP(w, r)
						return
					}

					// Treat the request path as relative by stripping the leading forward slash.
					relativePath := strings.TrimPrefix(cleanPath, "/")
					// Convert to OS-specific path separators for filesystem operations.
					relativePath = filepath.FromSlash(relativePath)

					// Resolve the requested path against the absolute assets directory.
					resolvedPath, err := filepath.Abs(filepath.Join(assetsDirAbs, relativePath))
					if err != nil {
						// If the path cannot be resolved, fall back to default handler.
						next.ServeHTTP(w, r)
						return
					}

					// Ensure the resolved path is still within the assets directory.
					// This check prevents path traversal attacks like "/../../../etc/passwd".
					if resolvedPath != assetsDirAbs && !strings.HasPrefix(resolvedPath, assetsDirAbs+string(filepath.Separator)) {
						// Path traversal attempt detected, fall back to default handler.
						next.ServeHTTP(w, r)
						return
					}

					// Path is validated to be within assetsDirAbs above.
					if _, err := os.Stat(resolvedPath); err == nil { // #nosec G304 // lgtm[go/path-injection] -- path validated above
						// Serve file from disk to make testing easy
						http.ServeFile(w, r, resolvedPath) // #nosec G304 // lgtm[go/path-injection] -- path validated above
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
