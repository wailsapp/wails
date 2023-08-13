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
			FS: assets,
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Custom event handling
	app.Events.On("myevent", func(e *application.WailsEvent) {
		log.Printf("[Go] WailsEvent received: %+v\n", e)
	})

	// OS specific application events
	app.On(events.Mac.ApplicationDidFinishLaunching, func(event *application.Event) {
		for {
			log.Println("Sending event")
			app.Events.Emit(&application.WailsEvent{
				Name: "myevent",
				Data: "hello",
			})
			time.Sleep(10 * time.Second)
		}
	})

	// Platform agnostic events
	app.On(events.Common.ApplicationStarted, func(event *application.Event) {
		println("events.Common.ApplicationStarted fired!")
	})

	win1 := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title: "Events Demo",
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
			println("Closing!")
			return
		}
		println("Nope! Not closing!")
		e.Cancel()
	})

	win2 := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title: "Events Demo",
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInsetUnified,
			InvisibleTitleBarHeight: 50,
		},
	})

	var cancel bool

	win2.RegisterHook(events.Common.WindowFocus, func(e *application.WindowEvent) {
		println("---------------\n[Hook] Window focus!")
		cancel = !cancel
		if cancel {
			e.Cancel()
		}
	})

	win2.On(events.Common.WindowFocus, func(e *application.WindowEvent) {
		println("[Event] Window focus!")
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
