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
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Events.On("myevent", func(e *application.CustomEvent) {
		log.Printf("[Go] CustomEvent received: %+v\n", e)
	})

	app.On(events.Mac.ApplicationDidFinishLaunching, func() {
		for {
			log.Println("Sending event")
			app.Events.Emit(&application.CustomEvent{
				Name: "myevent",
				Data: "hello",
			})
			time.Sleep(10 * time.Second)
		}
	})

	app.NewWebviewWindowWithOptions(&application.WebviewWindowOptions{
		Title: "Events Demo",
		Assets: application.AssetOptions{
			FS: assets,
		},
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInsetUnified,
			InvisibleTitleBarHeight: 50,
		},
	})
	app.NewWebviewWindowWithOptions(&application.WebviewWindowOptions{
		Title: "Events Demo",
		Assets: application.AssetOptions{
			FS: assets,
		},
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
