package main

import (
	"embed"
	"log"
	"log/slog"
	"net/http"
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
					assetsDir, err := filepath.Abs(filepath.Join(dir, "assets"))
					if err != nil {
						next.ServeHTTP(w, r)
						return
					}

					// URL paths always use '/', while filesystem paths are OS-specific.
					// Clean before joining so traversal segments cannot escape assets.
					cleanPath := path.Clean("/" + r.URL.Path)
					relativePath := filepath.FromSlash(strings.TrimPrefix(cleanPath, "/"))
					if filepath.IsAbs(relativePath) || filepath.VolumeName(relativePath) != "" {
						next.ServeHTTP(w, r)
						return
					}

					resolvedPath, err := filepath.Abs(filepath.Join(assetsDir, relativePath))
					if err != nil {
						next.ServeHTTP(w, r)
						return
					}

					rel, err := filepath.Rel(assetsDir, resolvedPath)
					if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
						next.ServeHTTP(w, r)
						return
					}

					// Stat follows symlinks, so resolve the existing target too. This
					// prevents an asset symlink from exposing a file outside assets.
					resolvedRealPath, err := filepath.EvalSymlinks(resolvedPath)
					if err != nil {
						next.ServeHTTP(w, r)
						return
					}
					assetsRealPath, err := filepath.EvalSymlinks(assetsDir)
					if err != nil {
						next.ServeHTTP(w, r)
						return
					}
					realRel, err := filepath.Rel(assetsRealPath, resolvedRealPath)
					if err != nil || realRel == ".." || strings.HasPrefix(realRel, ".."+string(filepath.Separator)) {
						next.ServeHTTP(w, r)
						return
					}

					// Serve through a rooted file server rather than passing a
					// request-derived filesystem path to http.ServeFile. The path
					// above has already been resolved and checked against assets.
					fileRequest := *r
					fileURL := *r.URL
					fileURL.Path = "/" + filepath.ToSlash(realRel)
					fileURL.RawPath = ""
					fileRequest.URL = &fileURL
					http.FileServer(http.Dir(assetsRealPath)).ServeHTTP(w, &fileRequest)
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
