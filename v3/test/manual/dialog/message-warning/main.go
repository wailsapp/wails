package main

import (
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/icons"
)

func main() {
	app := application.New(application.Options{
		Name:   "Dialog Test - Warning",
		Assets: application.AlphaAssets,
	})

	menu := app.NewMenu()

	testMenu := menu.AddSubmenu("Tests")

	testMenu.Add("Basic Warning").OnClick(func(ctx *application.Context) {
		app.Dialog.Warning().
			SetTitle("Warning").
			SetMessage("This is a warning message").
			Show()
	})

	testMenu.Add("Title Only").OnClick(func(ctx *application.Context) {
		app.Dialog.Warning().
			SetTitle("Warning - Title Only").
			Show()
	})

	testMenu.Add("Message Only").OnClick(func(ctx *application.Context) {
		app.Dialog.Warning().
			SetMessage("Warning message without title").
			Show()
	})

	testMenu.Add("Custom Icon").OnClick(func(ctx *application.Context) {
		app.Dialog.Warning().
			SetTitle("Custom Warning Icon").
			SetMessage("This warning has a custom icon").
			SetIcon(icons.ApplicationLightMode256).
			Show()
	})

	testMenu.Add("Long Warning").OnClick(func(ctx *application.Context) {
		app.Dialog.Warning().
			SetTitle("Important Warning").
			SetMessage("This is an important warning that contains a lot of text. " +
				"You should read this carefully before proceeding. " +
				"Ignoring this warning may result in unexpected behavior.").
			Show()
	})

	testMenu.Add("Attached to Window").OnClick(func(ctx *application.Context) {
		app.Dialog.Warning().
			SetTitle("Attached Warning").
			SetMessage("This warning is attached to the main window").
			AttachToWindow(app.Window.Current()).
			Show()
	})

	testMenu.AddSeparator()
	testMenu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Warning Dialog Tests",
		Width:  400,
		Height: 200,
		Linux: application.LinuxWindow{
			MenuStyle: application.LinuxMenuStylePrimaryMenu,
		},
	})
	window.SetMenu(menu)

	log.Println("Warning Dialog Tests")
	log.Println("Use the Tests menu to run each test case")

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
