package main

import (
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/icons"
)

func main() {
	app := application.New(application.Options{
		Name:   "Dialog Test - Question",
		Assets: application.AlphaAssets,
	})

	menu := app.NewMenu()

	testMenu := menu.AddSubmenu("Tests")

	testMenu.Add("Two Buttons").OnClick(func(ctx *application.Context) {
		dialog := app.Dialog.Question().
			SetTitle("Question").
			SetMessage("Do you want to proceed?")
		dialog.AddButton("Yes").OnClick(func() {
			app.Dialog.Info().SetMessage("You clicked Yes").Show()
		})
		dialog.AddButton("No").OnClick(func() {
			app.Dialog.Info().SetMessage("You clicked No").Show()
		})
		dialog.Show()
	})

	testMenu.Add("Three Buttons").OnClick(func(ctx *application.Context) {
		dialog := app.Dialog.Question().
			SetTitle("Save Changes?").
			SetMessage("You have unsaved changes")
		dialog.AddButton("Save").OnClick(func() {
			app.Dialog.Info().SetMessage("Saving...").Show()
		})
		dialog.AddButton("Don't Save").OnClick(func() {
			app.Dialog.Info().SetMessage("Discarding changes").Show()
		})
		dialog.AddButton("Cancel")
		dialog.Show()
	})

	testMenu.Add("With Default Button").OnClick(func(ctx *application.Context) {
		dialog := app.Dialog.Question().
			SetTitle("Confirm").
			SetMessage("Press Enter to select the default button")
		dialog.AddButton("OK")
		no := dialog.AddButton("Cancel")
		dialog.SetDefaultButton(no)
		dialog.Show()
	})

	testMenu.Add("With Cancel Button (Escape)").OnClick(func(ctx *application.Context) {
		dialog := app.Dialog.Question().
			SetTitle("Escape Test").
			SetMessage("Press Escape to cancel")
		ok := dialog.AddButton("OK").OnClick(func() {
			app.Dialog.Info().SetMessage("OK clicked").Show()
		})
		cancel := dialog.AddButton("Cancel").OnClick(func() {
			app.Dialog.Info().SetMessage("Cancelled").Show()
		})
		dialog.SetDefaultButton(ok)
		dialog.SetCancelButton(cancel)
		dialog.Show()
	})

	testMenu.Add("Custom Icon").OnClick(func(ctx *application.Context) {
		dialog := app.Dialog.Question().
			SetTitle("Custom Icon").
			SetMessage("This question has a custom icon").
			SetIcon(icons.WailsLogoWhiteTransparent)
		dialog.AddButton("Nice!").OnClick(func() {
			app.Dialog.Info().SetMessage("Thanks!").Show()
		})
		dialog.AddButton("Meh")
		dialog.Show()
	})

	testMenu.Add("Attached to Window").OnClick(func(ctx *application.Context) {
		dialog := app.Dialog.Question().
			SetTitle("Attached").
			SetMessage("This dialog is attached to the window").
			AttachToWindow(app.Window.Current())
		dialog.AddButton("OK")
		dialog.AddButton("Cancel")
		dialog.Show()
	})

	testMenu.Add("Button Callbacks").OnClick(func(ctx *application.Context) {
		dialog := app.Dialog.Question().
			SetTitle("Callbacks").
			SetMessage("Each button has a callback")
		dialog.AddButton("Option A").OnClick(func() {
			log.Println("Option A selected")
			app.Dialog.Info().SetMessage("You chose Option A").Show()
		})
		dialog.AddButton("Option B").OnClick(func() {
			log.Println("Option B selected")
			app.Dialog.Info().SetMessage("You chose Option B").Show()
		})
		dialog.AddButton("Option C").OnClick(func() {
			log.Println("Option C selected")
			app.Dialog.Info().SetMessage("You chose Option C").Show()
		})
		dialog.Show()
	})

	testMenu.AddSeparator()
	testMenu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Question Dialog Tests",
		Width:  400,
		Height: 200,
		Linux: application.LinuxWindow{
			MenuStyle: application.LinuxMenuStylePrimaryMenu,
		},
	})
	window.SetMenu(menu)

	log.Println("Question Dialog Tests")
	log.Println("Use the Tests menu to run each test case")

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
