package main

import (
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/icons"
)

func main() {
	app := application.New(application.Options{
		Name:   "Dialog Test - Info",
		Assets: application.AlphaAssets,
	})

	menu := app.NewMenu()

	testMenu := menu.AddSubmenu("Tests")

	testMenu.Add("Basic Info").OnClick(func(ctx *application.Context) {
		app.Dialog.Info().
			SetTitle("Information").
			SetMessage("This is a basic info dialog").
			Show()
	})

	testMenu.Add("Title Only").OnClick(func(ctx *application.Context) {
		app.Dialog.Info().
			SetTitle("Title Only - No Message").
			Show()
	})

	testMenu.Add("Message Only").OnClick(func(ctx *application.Context) {
		app.Dialog.Info().
			SetMessage("Message only - no title set").
			Show()
	})

	testMenu.Add("Custom Icon").OnClick(func(ctx *application.Context) {
		app.Dialog.Info().
			SetTitle("Custom Icon").
			SetMessage("This dialog has a custom icon").
			SetIcon(icons.WailsLogoBlackTransparent).
			Show()
	})

	testMenu.Add("Long Message").OnClick(func(ctx *application.Context) {
		app.Dialog.Info().
			SetTitle("Long Message Test").
			SetMessage("This is a very long message that should wrap properly in the dialog. " +
				"It contains multiple sentences to test how the dialog handles longer content. " +
				"The dialog should display this text in a readable manner without truncation.").
			Show()
	})

	testMenu.Add("Attached to Window").OnClick(func(ctx *application.Context) {
		app.Dialog.Info().
			SetTitle("Attached Dialog").
			SetMessage("This dialog is attached to the main window").
			AttachToWindow(app.Window.Current()).
			Show()
	})

	testMenu.AddSeparator()
	testMenu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Info Dialog Tests",
		Width:  400,
		Height: 200,
		Linux: application.LinuxWindow{
			MenuStyle: application.LinuxMenuStylePrimaryMenu,
		},
	})
	window.SetMenu(menu)

	log.Println("Info Dialog Tests")
	log.Println("Use the Tests menu to run each test case")

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
