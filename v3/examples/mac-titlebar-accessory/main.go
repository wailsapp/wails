package main

import (
	"embed"
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

var statusCycle = []string{"idle", "syncing", "synced"}

func main() {
	app := application.New(application.Options{
		Name:        "mac-titlebar-accessory",
		Description: "A demo of AddTitlebarAccessory: independent webviews pinned to the titlebar",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Titlebar Accessory",
		Width:  620,
		Height: 380,
		URL:    "/",
	})

	leading, err := window.AddTitlebarAccessory(application.MacTitlebarAccessoryOptions{
		Name:     "ping",
		URL:      "/leading.html",
		Position: application.AccessoryLeading,
		Width:    92,
		Height:   24,
	})
	if err != nil {
		log.Fatal(err)
	}
	_ = leading

	status, err := window.AddTitlebarAccessory(application.MacTitlebarAccessoryOptions{
		Name:     "status",
		URL:      "/trailing.html",
		Position: application.AccessoryTrailing,
		Width:    110,
		Height:   24,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Proves ExecJS reaches the accessory's own webview independently of the
	// main window: the status pill cycles on a timer driven entirely from
	// here, with nothing in the main window's content aware of it.
	go func() {
		i := 0
		for {
			time.Sleep(2500 * time.Millisecond)
			state := statusCycle[i%len(statusCycle)]
			status.ExecJS(`window.setStatus && window.setStatus("` + state + `")`)
			i++
		}
	}()

	app.Event.On("titlebar:ping", func(e *application.CustomEvent) {
		log.Printf("titlebar:ping received from the leading accessory: %v", e.Data)
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
