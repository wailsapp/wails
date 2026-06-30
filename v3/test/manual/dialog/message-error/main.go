package main

import (
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/icons"
)

func main() {
	app := application.New(application.Options{
		Name:   "Dialog Test - Error",
		Assets: application.AlphaAssets,
	})

	menu := app.NewMenu()

	

	menu.Add("Basic Error").OnClick(func(ctx *application.Context) {
		app.Dialog.Error().
			SetTitle("Error").
			SetMessage("An error has occurred").
			Show()
	})

	menu.Add("Title Only").OnClick(func(ctx *application.Context) {
		app.Dialog.Error().
			SetTitle("Error - Something went wrong").
			Show()
	})

	menu.Add("Message Only").OnClick(func(ctx *application.Context) {
		app.Dialog.Error().
			SetMessage("Error message without a title").
			Show()
	})

	menu.Add("Custom Icon").OnClick(func(ctx *application.Context) {
		app.Dialog.Error().
			SetTitle("Custom Error Icon").
			SetMessage("This error dialog has a custom icon").
			SetIcon(icons.WailsLogoBlack).
			Show()
	})

	menu.Add("Technical Error").OnClick(func(ctx *application.Context) {
		app.Dialog.Error().
			SetTitle("Connection Failed").
			SetMessage("Failed to connect to server at localhost:8080. " +
				"Error: connection refused. " +
				"Please check that the server is running and try again.").
			Show()
	})

	menu.Add("Attached to Window").OnClick(func(ctx *application.Context) {
		app.Dialog.Error().
			SetTitle("Attached Error").
			SetMessage("This error dialog is attached to the main window").
			AttachToWindow(app.Window.Current()).
			Show()
	})

	menu.AddSeparator()
	menu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Error Dialog Tests",
		Width:  400,
		Height: 200,
		Linux: application.LinuxWindow{
			MenuStyle: application.LinuxMenuStylePrimaryMenu,
		},
	})
	window.SetMenu(menu)

	log.Println("Error Dialog Tests")
	log.Println("Use the Tests menu to run each test case")

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
