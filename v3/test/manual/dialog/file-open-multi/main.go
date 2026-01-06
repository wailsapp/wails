package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	app := application.New(application.Options{
		Name:   "Dialog Test - Open Multiple Files",
		Assets: application.AlphaAssets,
	})

	menu := app.NewMenu()

	testMenu := menu.AddSubmenu("Tests")

	testMenu.Add("Select Multiple Files").OnClick(func(ctx *application.Context) {
		results, err := app.Dialog.OpenFile().
			SetTitle("Select Multiple Files").
			CanChooseFiles(true).
			PromptForMultipleSelection()
		showResults(app, "Multi Select", results, err)
	})

	testMenu.Add("With Hidden Files").OnClick(func(ctx *application.Context) {
		results, err := app.Dialog.OpenFile().
			SetTitle("Select Files (Hidden Visible)").
			CanChooseFiles(true).
			ShowHiddenFiles(true).
			PromptForMultipleSelection()
		showResults(app, "Hidden Files", results, err)
	})

	testMenu.Add("Filter: Source Code").OnClick(func(ctx *application.Context) {
		results, err := app.Dialog.OpenFile().
			SetTitle("Select Source Files").
			CanChooseFiles(true).
			AddFilter("Go Files", "*.go").
			AddFilter("All Source", "*.go;*.js;*.ts;*.py;*.rs").
			PromptForMultipleSelection()
		showResults(app, "Source Filter", results, err)
	})

	testMenu.Add("Filter: Documents").OnClick(func(ctx *application.Context) {
		results, err := app.Dialog.OpenFile().
			SetTitle("Select Documents").
			CanChooseFiles(true).
			AddFilter("Documents", "*.pdf;*.doc;*.docx;*.txt;*.md").
			PromptForMultipleSelection()
		showResults(app, "Doc Filter", results, err)
	})

	testMenu.Add("Attached to Window").OnClick(func(ctx *application.Context) {
		results, err := app.Dialog.OpenFile().
			SetTitle("Select Files (Attached)").
			CanChooseFiles(true).
			AttachToWindow(app.Window.Current()).
			PromptForMultipleSelection()
		showResults(app, "Attached", results, err)
	})

	testMenu.AddSeparator()
	testMenu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Open Multiple Files Tests",
		Width:  400,
		Height: 200,
		Linux: application.LinuxWindow{
			MenuStyle: application.LinuxMenuStylePrimaryMenu,
		},
	})
	window.SetMenu(menu)

	log.Println("Open Multiple Files Tests")
	log.Println("Use the Tests menu to run each test case")

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func showResults(app *application.App, test string, results []string, err error) {
	if err != nil {
		log.Printf("[%s] Error: %v", test, err)
		app.Dialog.Error().SetTitle("Error").SetMessage(err.Error()).Show()
		return
	}
	if len(results) == 0 {
		log.Printf("[%s] Cancelled", test)
		app.Dialog.Info().SetTitle(test).SetMessage("No files selected").Show()
		return
	}
	log.Printf("[%s] Selected %d files:", test, len(results))
	for _, r := range results {
		log.Printf("  - %s", r)
	}
	msg := fmt.Sprintf("Selected %d files:\n%s", len(results), strings.Join(results, "\n"))
	app.Dialog.Info().SetTitle(test).SetMessage(msg).Show()
}
