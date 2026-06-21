package main

import (
	"log"
	"os"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	app := application.New(application.Options{
		Name:   "Dialog Test - Directory",
		Assets: application.AlphaAssets,
	})

	menu := app.NewMenu()

	menu.Add("Basic Directory").OnClick(func(ctx *application.Context) {
		result, err := app.Dialog.OpenFile().
			SetTitle("Select Directory").
			CanChooseDirectories(true).
			CanChooseFiles(false).
			PromptForSingleSelection()
		showResult(app, "Basic Directory", result, err)
	})

	menu.Add("Start in Home").OnClick(func(ctx *application.Context) {
		home, _ := os.UserHomeDir()
		result, err := app.Dialog.OpenFile().
			SetTitle("Select from Home").
			SetDirectory(home).
			CanChooseDirectories(true).
			CanChooseFiles(false).
			PromptForSingleSelection()
		showResult(app, "Home Directory", result, err)
	})

	menu.Add("Start in /").OnClick(func(ctx *application.Context) {
		result, err := app.Dialog.OpenFile().
			SetTitle("Select from Root").
			SetDirectory("/").
			CanChooseDirectories(true).
			CanChooseFiles(false).
			PromptForSingleSelection()
		showResult(app, "Root Directory", result, err)
	})

	menu.Add("Can Create Directories").OnClick(func(ctx *application.Context) {
		result, err := app.Dialog.OpenFile().
			SetTitle("Select Directory (Can Create)").
			CanChooseDirectories(true).
			CanChooseFiles(false).
			CanCreateDirectories(true).
			PromptForSingleSelection()
		showResult(app, "Create Dirs", result, err)
	})

	menu.Add("Show Hidden").OnClick(func(ctx *application.Context) {
		result, err := app.Dialog.OpenFile().
			SetTitle("Select Directory (Hidden Visible)").
			CanChooseDirectories(true).
			CanChooseFiles(false).
			ShowHiddenFiles(true).
			PromptForSingleSelection()
		showResult(app, "Show Hidden", result, err)
	})

	menu.Add("Resolve Aliases/Symlinks").OnClick(func(ctx *application.Context) {
		result, err := app.Dialog.OpenFile().
			SetTitle("Select Directory (Resolve Symlinks)").
			CanChooseDirectories(true).
			CanChooseFiles(false).
			ResolvesAliases(true).
			PromptForSingleSelection()
		showResult(app, "Resolve Aliases", result, err)
	})

	menu.Add("Custom Button Text").OnClick(func(ctx *application.Context) {
		result, err := app.Dialog.OpenFile().
			SetTitle("Choose Project Folder").
			SetButtonText("Use This Folder").
			CanChooseDirectories(true).
			CanChooseFiles(false).
			PromptForSingleSelection()
		showResult(app, "Custom Button", result, err)
	})

	menu.Add("Multiple Directories").OnClick(func(ctx *application.Context) {
		results, err := app.Dialog.OpenFile().
			SetTitle("Select Multiple Directories").
			CanChooseDirectories(true).
			CanChooseFiles(false).
			PromptForMultipleSelection()
		if err != nil {
			log.Printf("[Multi Dir] Error: %v", err)
			app.Dialog.Error().SetTitle("Error").SetMessage(err.Error()).Show()
			return
		}
		if len(results) == 0 {
			log.Printf("[Multi Dir] Cancelled")
			app.Dialog.Info().SetTitle("Multi Dir").SetMessage("No directories selected").Show()
			return
		}
		log.Printf("[Multi Dir] Selected %d directories", len(results))
		msg := ""
		for _, r := range results {
			msg += r + "\n"
		}
		app.Dialog.Info().SetTitle("Multi Dir").SetMessage(msg).Show()
	})

	menu.Add("Attached to Window").OnClick(func(ctx *application.Context) {
		result, err := app.Dialog.OpenFile().
			SetTitle("Select Directory (Attached)").
			CanChooseDirectories(true).
			CanChooseFiles(false).
			AttachToWindow(app.Window.Current()).
			PromptForSingleSelection()
		showResult(app, "Attached", result, err)
	})

	menu.AddSeparator()
	menu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Directory Dialog Tests",
		Width:  400,
		Height: 200,
		Linux: application.LinuxWindow{
			MenuStyle: application.LinuxMenuStylePrimaryMenu,
		},
	})
	window.SetMenu(menu)

	log.Println("Directory Dialog Tests")
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
		app.Dialog.Info().SetTitle(test).SetMessage("No directory selected").Show()
		return
	}
	log.Printf("[%s] Selected: %s", test, result)
	app.Dialog.Info().SetTitle(test).SetMessage(result).Show()
}
