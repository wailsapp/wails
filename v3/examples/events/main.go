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
		Name:        "Events Demo",
		Description: "A demo of the Events API",
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Custom event handling
	app.Events.On("myevent", func(e *application.WailsEvent) {
		app.Logger.Info("[Go] WailsEvent received", "name", e.Name, "data", e.Data, "sender", e.Sender, "cancelled", e.Cancelled)
	})

	// OS specific application events
	app.On(events.Common.ApplicationStarted, func(event *application.Event) {
		for {
			app.Events.Emit(&application.WailsEvent{
				Name: "myevent",
				Data: "hello",
			})
			time.Sleep(10 * time.Second)
		}
	})

	app.On(events.Common.ThemeChanged, func(event *application.Event) {
		app.Logger.Info("System theme changed!")
		if event.Context().IsDarkMode() {
			app.Logger.Info("System is now using dark mode!")
		} else {
			app.Logger.Info("System is now using light mode!")
		}
	})

	// Platform agnostic events
	app.On(events.Common.ApplicationStarted, func(event *application.Event) {
		app.Logger.Info("events.Common.ApplicationStarted fired!")
	})

	win1 := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title: "Events Demo",
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

	win2 := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title: "Events Demo",
		Name:  "Window 2",
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInsetUnified,
			InvisibleTitleBarHeight: 50,
		},
	})

	var cancel bool

	win2.RegisterHook(events.Common.WindowFocus, func(e *application.WindowEvent) {
		app.Logger.Info("[Hook] Window focus!")
		cancel = !cancel
		if cancel {
			e.Cancel()
		}
	})

	win2.On(events.Common.WindowFocus, func(e *application.WindowEvent) {
		app.Logger.Info("[Event] Window focus!")
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
