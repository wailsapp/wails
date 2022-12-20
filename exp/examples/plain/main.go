package main

import (
	_ "embed"
	"log"

	"github.com/wailsapp/wails/exp/pkg/options"

	"github.com/wailsapp/wails/exp/pkg/application"
)

func main() {
	app := application.New()

	// Create window
	app.NewWindowWithOptions(&options.Window{
		Title:          "Plain Bundle",
		EnableDevTools: true,
		HTML:           `<html><head><title>Plain Bundle</title></head><body><div class="main"><h1>Plain Bundle</h1><p>This is a plain bundle. It has no frontend code.</p></div></body></html>`,
		CSS:            `body { background-color: rgba(255, 255, 255, 0); } .main { color: white; margin: 20%; }`,
		Mac: &options.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                options.MacBackdropTranslucent,
			TitleBar:                options.TitleBarHiddenInset,
		},
	})
	// Create window
	app.NewWindowWithOptions(&options.Window{
		Title:          "Plain Bundle",
		EnableDevTools: true,
		HTML:           `<html><head><title>Plain Bundle</title></head><body><div class="main"><h1>Plain Bundle</h1><p>This is a plain bundle. It has no frontend code.</p></div></body></html>`,
		CSS:            `body { background-color: rgba(255, 255, 255, 0); } .main { color: white; margin: 20%; }`,
		Mac: &options.MacWindow{
			Backdrop: options.MacBackdropTranslucent,
			TitleBar: options.TitleBarHiddenInset,
		},
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
