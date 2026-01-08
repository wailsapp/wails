package main

import (
	"log"
	"os"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	app := application.New(application.Options{
		Name:   "Dialog Test - Open File",
		Assets: application.AlphaAssets,
	})

	menu := app.NewMenu()

	

	menu.Add("Basic Open").OnClick(func(ctx *application.Context) {
		result, err := app.Dialog.OpenFile().
			CanChooseFiles(true).
			PromptForSingleSelection()
		showResult(app, "Basic Open", result, err)
	})

	menu.Add("With Title").OnClick(func(ctx *application.Context) {
		result, err := app.Dialog.OpenFile().
			SetTitle("Select a File").
			CanChooseFiles(true).
			PromptForSingleSelection()
		showResult(app, "With Title", result, err)
	})

	menu.Add("Show Hidden Files").OnClick(func(ctx *application.Context) {
		result, err := app.Dialog.OpenFile().
			SetTitle("Select File (Hidden Visible)").
			CanChooseFiles(true).
			ShowHiddenFiles(true).
			PromptForSingleSelection()
		showResult(app, "Show Hidden", result, err)
	})

	menu.Add("Start in Home Directory").OnClick(func(ctx *application.Context) {
		home, _ := os.UserHomeDir()
		result, err := app.Dialog.OpenFile().
			SetTitle("Select from Home").
			SetDirectory(home).
			CanChooseFiles(true).
			PromptForSingleSelection()
		showResult(app, "Home Directory", result, err)
	})

	menu.Add("Start in /tmp").OnClick(func(ctx *application.Context) {
		result, err := app.Dialog.OpenFile().
			SetTitle("Select from /tmp").
			SetDirectory("/tmp").
			CanChooseFiles(true).
			PromptForSingleSelection()
		showResult(app, "/tmp Directory", result, err)
	})

	menu.Add("Filter: Text Files").OnClick(func(ctx *application.Context) {
		result, err := app.Dialog.OpenFile().
			SetTitle("Select Text File").
			CanChooseFiles(true).
			AddFilter("Text Files", "*.txt;*.md;*.log").
			PromptForSingleSelection()
		showResult(app, "Text Filter", result, err)
	})

	menu.Add("Filter: Images").OnClick(func(ctx *application.Context) {
		result, err := app.Dialog.OpenFile().
			SetTitle("Select Image").
			CanChooseFiles(true).
			AddFilter("Images", "*.png;*.jpg;*.jpeg;*.gif;*.webp").
			PromptForSingleSelection()
		showResult(app, "Image Filter", result, err)
	})

	menu.Add("Multiple Filters").OnClick(func(ctx *application.Context) {
		result, err := app.Dialog.OpenFile().
			SetTitle("Select File").
			CanChooseFiles(true).
			AddFilter("Documents", "*.txt;*.md;*.pdf").
			AddFilter("Images", "*.png;*.jpg;*.gif").
			AddFilter("All Files", "*").
			PromptForSingleSelection()
		showResult(app, "Multi Filter", result, err)
	})

	menu.Add("Custom Button Text").OnClick(func(ctx *application.Context) {
		result, err := app.Dialog.OpenFile().
			SetTitle("Select File").
			SetButtonText("Choose This One").
			CanChooseFiles(true).
			PromptForSingleSelection()
		showResult(app, "Custom Button", result, err)
	})

	menu.Add("Attached to Window").OnClick(func(ctx *application.Context) {
		result, err := app.Dialog.OpenFile().
			SetTitle("Attached Dialog").
			CanChooseFiles(true).
			AttachToWindow(app.Window.Current()).
			PromptForSingleSelection()
		showResult(app, "Attached", result, err)
	})

	menu.AddSeparator()
	menu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Open File Dialog Tests",
		Width:  400,
		Height: 200,
		Linux: application.LinuxWindow{
			MenuStyle: application.LinuxMenuStylePrimaryMenu,
		},
	})
	window.SetMenu(menu)

	log.Println("Open File Dialog Tests")
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
