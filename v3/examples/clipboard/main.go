package main

import (
	_ "embed"
	"log"
	"runtime"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {

	app := application.New(application.Options{
		Name:        "Clipboard Demo",
		Description: "A demo of the clipboard API",
		Assets:      application.AlphaAssets,
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Create a custom menu
	menu := app.NewMenu()
	if runtime.GOOS == "darwin" {
		menu.AddRole(application.AppMenu)
	}

	setClipboardMenu := menu.AddSubmenu("Set Clipboard")
	setClipboardMenu.Add("Set Text 'Hello'").OnClick(func(ctx *application.Context) {
		success := app.Clipboard.SetText("Hello")
		if !success {
			_ = app.Dialog.Info().SetMessage("Failed to set clipboard text").Show()
		}
	})
	setClipboardMenu.Add("Set Text 'World'").OnClick(func(ctx *application.Context) {
		success := app.Clipboard.SetText("World")
		if !success {
			_ = app.Dialog.Info().SetMessage("Failed to set clipboard text").Show()
		}
	})
	setClipboardMenu.Add("Set Text (current time)").OnClick(func(ctx *application.Context) {
		success := app.Clipboard.SetText(time.Now().String())
		if !success {
			_ = app.Dialog.Info().SetMessage("Failed to set clipboard text").Show()
		}
	})
	getClipboardMenu := menu.AddSubmenu("Get Clipboard")
	getClipboardMenu.Add("Get Text").OnClick(func(ctx *application.Context) {
		result, ok := app.Clipboard.Text()
		if !ok {
			_ = app.Dialog.Info().SetMessage("Failed to get clipboard text").Show()
		} else {
			_ = app.Dialog.Info().SetMessage("Got:\n\n" + result).Show()
		}
	})

	clearClipboardMenu := menu.AddSubmenu("Clear Clipboard")
	clearClipboardMenu.Add("Clear Text").OnClick(func(ctx *application.Context) {
		success := app.Clipboard.SetText("")
		if success {
			_ = app.Dialog.Info().SetMessage("Clipboard text cleared").Show()
		} else {
			_ = app.Dialog.Info().SetMessage("Clipboard text not cleared").Show()
		}
	})

	app.Menu.Set(menu)

	app.Window.New()

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
