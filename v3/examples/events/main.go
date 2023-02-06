package main

import (
	"embed"
	_ "embed"
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/events"

	"github.com/wailsapp/wails/v3/pkg/options"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

func main() {

	app := application.New(options.Application{
		Name:        "Events Demo",
		Description: "A demo of the Events API",
		Mac: options.Mac{
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

	app.NewWebviewWindowWithOptions(&options.WebviewWindow{
		Title: "Events Demo",
		Assets: options.Assets{
			FS: assets,
		},
		Mac: options.MacWindow{
			Backdrop:                options.MacBackdropTranslucent,
			TitleBar:                options.TitleBarHiddenInsetUnified,
			InvisibleTitleBarHeight: 50,
		},
	})
	app.NewWebviewWindowWithOptions(&options.WebviewWindow{
		Title: "Events Demo",
		Assets: options.Assets{
			FS: assets,
		},
		Mac: options.MacWindow{
			Backdrop:                options.MacBackdropTranslucent,
			TitleBar:                options.TitleBarHiddenInsetUnified,
			InvisibleTitleBarHeight: 50,
		},
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
