package main

import (
	"embed"
	"github.com/wailsapp/wails/v3/examples/services/hashes"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/services/fileserver"
	"github.com/wailsapp/wails/v3/pkg/services/kvstore"
	"github.com/wailsapp/wails/v3/pkg/services/log"
	"github.com/wailsapp/wails/v3/pkg/services/sqlite"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
)

//go:embed assets/*
var assets embed.FS

func main() {

	// Get the local directory of this source file
	// This isn't needed when running the example with `go run .`
	// but is needed when running the example from an IDE
	_, thisFile, _, _ := runtime.Caller(0)
	localDir := filepath.Dir(thisFile)

	rootPath := filepath.Join(localDir, "files")
	app := application.New(application.Options{
		Name:        "Services Demo",
		Description: "A demo of the services API",
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		LogLevel: slog.LevelDebug,
		Services: []application.Service{
			application.NewService(hashes.New()),
			application.NewService(sqlite.New(&sqlite.Config{
				DBFile: "test.db",
			})),
			application.NewService(kvstore.New(&kvstore.Config{
				Filename: "store.json",
				AutoSave: true,
			})),
			application.NewService(log.New()),
			application.NewService(fileserver.New(&fileserver.Config{
				RootPath: rootPath,
			}), application.ServiceOptions{
				Route: "/files",
			}),
		},
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
	})

	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Width:  1024,
		Height: 768,
	})

	err := app.Run()

	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
