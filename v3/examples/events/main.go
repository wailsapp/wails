package main

import (
	"embed"
	_ "embed"
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed assets
var assets embed.FS

func main() {

	app := application.New(application.Options{
		Name:        "customEventProcessor Demo",
		Description: "A demo of the customEventProcessor API",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Custom event handling
	app.Event.On("myevent", func(e *application.CustomEvent) {
		app.Logger.Info("[Go] CustomEvent received", "name", e.Name, "data", e.Data, "sender", e.Sender, "cancelled", e.IsCancelled())
	})

	// OS specific application events
	app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(event *application.ApplicationEvent) {
		go func() {
			ticker := time.NewTicker(10 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					// This emits a custom event every 10 seconds
					// As it's sent from the application, the sender will be blank
					app.Event.Emit("myevent", "hello")
				case <-app.Context().Done():
					return
				}
			}
		}()
	})

	app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(event *application.ApplicationEvent) {
		app.Logger.Info("System theme changed!")
		if event.Context().IsDarkMode() {
			app.Logger.Info("System is now using dark mode!")
		} else {
			app.Logger.Info("System is now using light mode!")
		}
	})

	// Platform agnostic events
	app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(event *application.ApplicationEvent) {
		app.Logger.Info("events.Common.ApplicationStarted fired!")
	})

	win1 := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Window 1",
		Name:  "Window 1",
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInsetUnified,
			InvisibleTitleBarHeight: 50,
		},
	})

	var countdown = 3

	win1.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		countdown--
		if countdown == 0 {
			app.Logger.Info("Window 1 Closing!")
			return
		}
		app.Logger.Info("Window 1 Closing? Nope! Not closing!")
		e.Cancel()
	})

	win2 := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Window 2",
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInsetUnified,
			InvisibleTitleBarHeight: 50,
		},
	})

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				win2.EmitEvent("windowevent", "ooooh!")
			case <-app.Context().Done():
				return
			}
		}
	}()

	var cancel bool

	win2.RegisterHook(events.Common.WindowFocus, func(e *application.WindowEvent) {
		app.Logger.Info("[Hook] Window focus!")
		cancel = !cancel
		if cancel {
			e.Cancel()
		}
	})

	win2.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
		app.Logger.Info("[OnWindowEvent] Window focus!")
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
