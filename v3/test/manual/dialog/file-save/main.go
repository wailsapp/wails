package main

import (
	"log"
	"os"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	app := application.New(application.Options{
		Name:   "Dialog Test - Save File",
		Assets: application.AlphaAssets,
	})

	menu := app.NewMenu()

	

	menu.Add("Basic Save").OnClick(func(ctx *application.Context) {
		result, err := app.Dialog.SaveFile().
			PromptForSingleSelection()
		showResult(app, "Basic Save", result, err)
	})

	menu.Add("With Message").OnClick(func(ctx *application.Context) {
		result, err := app.Dialog.SaveFile().
			SetMessage("Save your file").
			PromptForSingleSelection()
		showResult(app, "With Message", result, err)
	})

	menu.Add("With Default Filename").OnClick(func(ctx *application.Context) {
		result, err := app.Dialog.SaveFile().
			SetMessage("Save Document").
			SetFilename("document.txt").
			PromptForSingleSelection()
		showResult(app, "Default Filename", result, err)
	})

	menu.Add("Start in Home").OnClick(func(ctx *application.Context) {
		home, _ := os.UserHomeDir()
		result, err := app.Dialog.SaveFile().
			SetMessage("Save to Home").
			SetDirectory(home).
			SetFilename("myfile.txt").
			PromptForSingleSelection()
		showResult(app, "Home Directory", result, err)
	})

	menu.Add("Start in /tmp").OnClick(func(ctx *application.Context) {
		result, err := app.Dialog.SaveFile().
			SetMessage("Save to /tmp").
			SetDirectory("/tmp").
			SetFilename("temp_file.txt").
			PromptForSingleSelection()
		showResult(app, "/tmp Directory", result, err)
	})

	menu.Add("Show Hidden Files").OnClick(func(ctx *application.Context) {
		result, err := app.Dialog.SaveFile().
			SetMessage("Save (Hidden Visible)").
			ShowHiddenFiles(true).
			PromptForSingleSelection()
		showResult(app, "Show Hidden", result, err)
	})

	menu.Add("Can Create Directories").OnClick(func(ctx *application.Context) {
		result, err := app.Dialog.SaveFile().
			SetMessage("Save (Can Create Dirs)").
			CanCreateDirectories(true).
			PromptForSingleSelection()
		showResult(app, "Create Dirs", result, err)
	})

	menu.Add("Cannot Create Directories").OnClick(func(ctx *application.Context) {
		result, err := app.Dialog.SaveFile().
			SetMessage("Save (No Create Dirs)").
			CanCreateDirectories(false).
			PromptForSingleSelection()
		showResult(app, "No Create Dirs", result, err)
	})

	menu.Add("Custom Button Text").OnClick(func(ctx *application.Context) {
		result, err := app.Dialog.SaveFile().
			SetMessage("Export File").
			SetButtonText("Export").
			SetFilename("export.json").
			PromptForSingleSelection()
		showResult(app, "Custom Button", result, err)
	})

	menu.Add("Attached to Window").OnClick(func(ctx *application.Context) {
		result, err := app.Dialog.SaveFile().
			SetMessage("Save (Attached)").
			AttachToWindow(app.Window.Current()).
			PromptForSingleSelection()
		showResult(app, "Attached", result, err)
	})

	menu.AddSeparator()
	menu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Save File Dialog Tests",
		Width:  400,
		Height: 200,
		Linux: application.LinuxWindow{
			MenuStyle: application.LinuxMenuStylePrimaryMenu,
		},
	})
	window.SetMenu(menu)

	log.Println("Save File Dialog Tests")
	log.Println("Use the Tests menu to run each test case")

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func showResult(app *application.App, test string, result string, err error) {
	if err != nil {
		log.Printf("[%s] Error: %v", test, err)
		app.Dialog.Error().SetTitle("Error").SetMessage(err.Error()).Show()
		return
	}
	if result == "" {
		log.Printf("[%s] Cancelled", test)
		app.Dialog.Info().SetTitle(test).SetMessage("No file selected").Show()
		return
	}
	log.Printf("[%s] Selected: %s", test, result)
	app.Dialog.Info().SetTitle(test).SetMessage(result).Show()
}
