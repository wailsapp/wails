package main

import (
	"embed"
	_ "embed"
	"image/color"
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/services/dock"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend/dist
var assets embed.FS

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {
	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.

	dockService := dock.NewWithOptions(dock.BadgeOptions{
		TextColour:       color.RGBA{255, 255, 204, 255},
		BackgroundColour: color.RGBA{16, 124, 16, 255},
		FontName:         "consolab.ttf",
		FontSize:         20,
		SmallFontSize:    14,
	})

	app := application.New(application.Options{
		Name:        "badge",
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			application.NewService(dockService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Window 1",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})

	app.Event.On("remove:badge", func(event *application.CustomEvent) {
		err := dockService.RemoveBadge()
		if err != nil {
			log.Fatal(err)
		}
	})

	app.Event.On("set:badge", func(event *application.CustomEvent) {
		text := event.Data.(string)
		err := dockService.SetBadge(text)
		if err != nil {
			log.Fatal(err)
		}
	})

	// Create a goroutine that emits an event containing the current time every second.
	// The frontend can listen to this event and update the UI accordingly.
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				now := time.Now().Format(time.RFC1123)
				app.Event.Emit("time", now)
			case <-app.Context().Done():
				return
			}
		}
	}()

	// Run the application. This blocks until the application has been exited.
	err := app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
