package main

import (
	_ "embed"
	"log"

	"github.com/wailsapp/wails/exp/pkg/application"
)

func main() {

	app := application.New()
	app.SetName("Dialogs Demo")
	app.SetDescription("A demo of the Wails dialogs")

	// Create a custom menu
	menu := app.NewMenu()
	menu.AddRole(application.AppMenu)

	// Let's make a "Demo" menu
	infoMenu := menu.AddSubmenu("Info")
	infoMenu.Add("Info").OnClick(func(ctx *application.Context) {
		dialog := app.NewInfoDialog()
		dialog.SetTitle("Custom Title")
		dialog.SetMessage("This is a custom message")
		dialog.Show()
	})
	infoMenu.Add("Info (Title only)").OnClick(func(ctx *application.Context) {
		dialog := app.NewInfoDialog()
		dialog.SetTitle("Custom Title")
		dialog.Show()
	})
	infoMenu.Add("Info (Message only)").OnClick(func(ctx *application.Context) {
		dialog := app.NewInfoDialog()
		dialog.SetMessage("This is a custom message")
		dialog.Show()
	})
	infoMenu.Add("Info (Custom Icon)").OnClick(func(ctx *application.Context) {
		dialog := app.NewInfoDialog()
		dialog.SetTitle("Custom Icon Example")
		dialog.SetMessage("Using a custom icon")
		dialog.SetIcon(application.DefaultApplicationIcon)
		dialog.Show()
	})

	questionMenu := menu.AddSubmenu("Question")
	questionMenu.Add("Question (No default)").OnClick(func(ctx *application.Context) {
		dialog := app.NewQuestionDialog()
		dialog.SetMessage("No default button")
		dialog.AddButton("Yes")
		dialog.AddButton("No")
		dialog.Show()
	})
	questionMenu.Add("Question (With Default)").OnClick(func(ctx *application.Context) {
		dialog := app.NewQuestionDialog()
		dialog.SetTitle("Quit")
		dialog.SetMessage("You have unsaved work. Are you sure you want to quit?")
		dialog.AddButton("Yes").OnClick(func() {
			app.Quit()
		})
		no := dialog.AddButton("No")
		dialog.SetDefaultButton(no)
		dialog.Show()
	})
	questionMenu.Add("Question (With Cancel)").OnClick(func(ctx *application.Context) {
		dialog := app.NewQuestionDialog().
			SetTitle("Update").
			SetMessage("The cancel button is selected when pressing escape")
		download := dialog.AddButton("📥 Download")
		download.OnClick(func() {
			app.NewInfoDialog().SetMessage("Downloading...").Show()
		})
		no := dialog.AddButton("Cancel")
		dialog.SetDefaultButton(download)
		dialog.SetCancelButton(no)
		dialog.Show()
	})
	questionMenu.Add("Question (Custom Icon)").OnClick(func(ctx *application.Context) {
		dialog := app.NewQuestionDialog()
		dialog.SetTitle("Custom Icon Example")
		dialog.SetMessage("Using a custom icon")
		dialog.SetIcon(application.WailsLogoWhiteTransparent)
		dialog.SetDefaultButton(dialog.AddButton("I like it!"))
		dialog.AddButton("Not so keen...")
		dialog.Show()
	})

	warningMenu := menu.AddSubmenu("Warning")
	warningMenu.Add("Warning").OnClick(func(ctx *application.Context) {
		dialog := app.NewWarningDialog()
		dialog.SetTitle("Custom Title")
		dialog.SetMessage("This is a custom message")
		dialog.Show()
	})
	warningMenu.Add("Warning (Title only)").OnClick(func(ctx *application.Context) {
		dialog := app.NewWarningDialog()
		dialog.SetTitle("Custom Title")
		dialog.Show()
	})
	warningMenu.Add("Warning (Message only)").OnClick(func(ctx *application.Context) {
		dialog := app.NewWarningDialog()
		dialog.SetMessage("This is a custom message")
		dialog.Show()
	})
	warningMenu.Add("Warning (Custom Icon)").OnClick(func(ctx *application.Context) {
		dialog := app.NewWarningDialog()
		dialog.SetTitle("Custom Icon Example")
		dialog.SetMessage("Using a custom icon")
		dialog.SetIcon(application.DefaultApplicationIcon)
		dialog.Show()
	})

	errorMenu := menu.AddSubmenu("Error")
	errorMenu.Add("Error").OnClick(func(ctx *application.Context) {
		dialog := app.NewErrorDialog()
		dialog.SetTitle("Ooops")
		dialog.SetMessage("I accidentally the whole of Twitter")
		dialog.Show()
	})
	errorMenu.Add("Error (Title Only)").OnClick(func(ctx *application.Context) {
		dialog := app.NewErrorDialog()
		dialog.SetTitle("Custom Title")
		dialog.Show()
	})
	errorMenu.Add("Error (Custom Message)").OnClick(func(ctx *application.Context) {
		dialog := app.NewErrorDialog()
		dialog.SetMessage("This is a custom message")
		dialog.Show()
	})
	errorMenu.Add("Error (Custom Icon)").OnClick(func(ctx *application.Context) {
		dialog := app.NewErrorDialog()
		dialog.SetTitle("Custom Icon Example")
		dialog.SetMessage("Using a custom icon")
		dialog.SetIcon(application.WailsLogoWhite)
		dialog.Show()
	})

	openMenu := menu.AddSubmenu("Open")
	openMenu.Add("Open File").OnClick(func(ctx *application.Context) {
		result, _ := app.NewOpenFileDialog().
			CanChooseFiles(true).
			Show()
		if result != "" {
			app.NewInfoDialog().SetMessage(result).Show()
		} else {
			app.NewInfoDialog().SetMessage("No file selected").Show()
		}
	})
	openMenu.Add("Open File (Show Hidden Files)").OnClick(func(ctx *application.Context) {
		result, _ := app.NewOpenFileDialog().
			CanChooseFiles(true).
			CanCreateDirectories(true).
			ShowHiddenFiles(true).
			Show()
		if result != "" {
			app.NewInfoDialog().SetMessage(result).Show()
		} else {
			app.NewInfoDialog().SetMessage("No file selected").Show()
		}
	})
	//openMenu.Add("Open Multiple Files (Show Hidden Files)").OnClick(func(ctx *application.Context) {
	//	result, _ := app.NewOpenMultipleFilesDialog().
	//		CanChooseFiles(true).
	//		CanCreateDirectories(true).
	//		ShowHiddenFiles(true).
	//		Show()
	//	if len(result) > 0 {
	//		app.NewInfoDialog().SetMessage(strings.Join(result, ",")).Show()
	//	} else {
	//		app.NewInfoDialog().SetMessage("No file selected").Show()
	//	}
	//})
	openMenu.Add("Open Directory").OnClick(func(ctx *application.Context) {
		result, _ := app.NewOpenFileDialog().
			CanChooseDirectories(true).
			Show()
		if result != "" {
			app.NewInfoDialog().SetMessage(result).Show()
		} else {
			app.NewInfoDialog().SetMessage("No directory selected").Show()
		}
	})
	openMenu.Add("Open Directory (Create Directories)").OnClick(func(ctx *application.Context) {
		result, _ := app.NewOpenFileDialog().
			CanChooseDirectories(true).
			CanCreateDirectories(true).
			Show()
		if result != "" {
			app.NewInfoDialog().SetMessage(result).Show()
		} else {
			app.NewInfoDialog().SetMessage("No directory selected").Show()
		}
	})

	app.SetMenu(menu)

	app.NewWindow()

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
