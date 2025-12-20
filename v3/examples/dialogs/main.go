package main

import (
	_ "embed"
	"log"
	"log/slog"
	"os"
	"runtime"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/icons"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {

	app := application.New(application.Options{
		Name:        "Dialogs Demo",
		Description: "A demo of the dialogs API",
		Assets:      application.AlphaAssets,
		Logger:      application.DefaultLogger(slog.LevelDebug),
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Create a custom menu
	menu := app.NewMenu()
	menu.AddRole(application.AppMenu)
	menu.AddRole(application.EditMenu)
	menu.AddRole(application.WindowMenu)
	menu.AddRole(application.ServicesMenu)
	menu.AddRole(application.HelpMenu)

	// Let's make a "Demo" menu
	infoMenu := menu.AddSubmenu("Info")
	infoMenu.Add("Info").OnClick(func(ctx *application.Context) {
		dialog := app.Dialog.Info()
		dialog.SetTitle("Custom Title")
		dialog.SetMessage("This is a custom message")
		dialog.Show()
	})

	infoMenu.Add("Info (Title only)").OnClick(func(ctx *application.Context) {
		dialog := app.Dialog.Info()
		dialog.SetTitle("Custom Title")
		dialog.Show()
	})
	infoMenu.Add("Info (Message only)").OnClick(func(ctx *application.Context) {
		dialog := app.Dialog.Info()
		dialog.SetMessage("This is a custom message")
		dialog.Show()
	})
	infoMenu.Add("Info (Custom Icon)").OnClick(func(ctx *application.Context) {
		dialog := app.Dialog.Info()
		dialog.SetTitle("Custom Icon Example")
		dialog.SetMessage("Using a custom icon")
		dialog.SetIcon(icons.ApplicationDarkMode256)
		dialog.Show()
	})
	infoMenu.Add("About").OnClick(func(ctx *application.Context) {
		app.Menu.ShowAbout()
	})

	questionMenu := menu.AddSubmenu("Question")
	questionMenu.Add("Question (No default)").OnClick(func(ctx *application.Context) {
		dialog := app.Dialog.Question()
		dialog.SetMessage("No default button")
		dialog.AddButton("Yes")
		dialog.AddButton("No")
		dialog.Show()
	})
	questionMenu.Add("Question (Attached to Window)").OnClick(func(ctx *application.Context) {
		dialog := app.Dialog.Question()
		dialog.AttachToWindow(app.Window.Current())
		dialog.SetMessage("No default button")
		dialog.AddButton("Yes")
		dialog.AddButton("No")
		dialog.Show()
	})
	questionMenu.Add("Question (With Default)").OnClick(func(ctx *application.Context) {
		dialog := app.Dialog.Question()
		dialog.SetTitle("Quit")
		dialog.SetMessage("You have unsaved work. Are you sure you want to quit?")
		dialog.AddButton("Yes")
		no := dialog.AddButton("No")
		dialog.SetDefaultButton(no)
		// Show() returns the clicked button's label
		result, _ := dialog.Show()
		if result == "Yes" {
			app.Quit()
		}
	})
	questionMenu.Add("Question (With Cancel)").OnClick(func(ctx *application.Context) {
		dialog := app.Dialog.Question().
			SetTitle("Update").
			SetMessage("The cancel button is selected when pressing escape")
		download := dialog.AddButton("📥 Download")
		no := dialog.AddButton("Cancel")
		dialog.SetDefaultButton(download)
		dialog.SetCancelButton(no)
		// Show() returns the clicked button's label - use switch for multiple options
		switch result, _ := dialog.Show(); result {
		case "📥 Download":
			app.Dialog.Info().SetMessage("Downloading...").Show()
		case "Cancel":
			// User cancelled
		}
	})
	questionMenu.Add("Question (Fluent Cancel/Default)").OnClick(func(ctx *application.Context) {
		// Fluent API using WithDefaultButton and WithCancelButton
		result, _ := app.Dialog.Question().
			SetTitle("Save Changes?").
			SetMessage("You have unsaved changes. What would you like to do?").
			WithDefaultButton("Save").
			WithButton("Don't Save").
			WithCancelButton("Cancel").
			Show()

		switch result {
		case "Save":
			app.Dialog.Info().SetMessage("Changes saved!").Show()
		case "Don't Save":
			app.Dialog.Info().SetMessage("Changes discarded.").Show()
		case "Cancel":
			// User cancelled - do nothing
		}
	})
	questionMenu.Add("Question (Custom Icon)").OnClick(func(ctx *application.Context) {
		dialog := app.Dialog.Question()
		dialog.SetTitle("Custom Icon Example")
		dialog.SetMessage("Using a custom icon")
		dialog.SetIcon(icons.WailsLogoWhiteTransparent)
		likeIt := dialog.AddButton("I like it!")
		dialog.AddButton("Not so keen...")
		dialog.SetDefaultButton(likeIt)
		// Show() returns the clicked button's label
		switch result, _ := dialog.Show(); result {
		case "I like it!":
			app.Dialog.Info().SetMessage("Thanks!").Show()
		case "Not so keen...":
			app.Dialog.Info().SetMessage("Too bad!").Show()
		}
	})

	warningMenu := menu.AddSubmenu("Warning")
	warningMenu.Add("Warning").OnClick(func(ctx *application.Context) {
		app.Dialog.Warning().
			SetTitle("Custom Title").
			SetMessage("This is a custom message").
			Show()
	})
	warningMenu.Add("Warning (Title only)").OnClick(func(ctx *application.Context) {
		dialog := app.Dialog.Warning()
		dialog.SetTitle("Custom Title")
		dialog.Show()
	})
	warningMenu.Add("Warning (Message only)").OnClick(func(ctx *application.Context) {
		dialog := app.Dialog.Warning()
		dialog.SetMessage("This is a custom message")
		dialog.Show()
	})
	warningMenu.Add("Warning (Custom Icon)").OnClick(func(ctx *application.Context) {
		dialog := app.Dialog.Warning()
		dialog.SetTitle("Custom Icon Example")
		dialog.SetMessage("Using a custom icon")
		dialog.SetIcon(icons.ApplicationLightMode256)
		dialog.Show()
	})

	errorMenu := menu.AddSubmenu("Error")
	errorMenu.Add("Error").OnClick(func(ctx *application.Context) {
		dialog := app.Dialog.Error()
		dialog.SetTitle("Ooops")
		dialog.SetMessage("I accidentally the whole of Twitter")
		dialog.Show()
	})
	errorMenu.Add("Error (Title Only)").OnClick(func(ctx *application.Context) {
		dialog := app.Dialog.Error()
		dialog.SetTitle("Custom Title")
		dialog.Show()
	})
	errorMenu.Add("Error (Custom Message)").OnClick(func(ctx *application.Context) {
		app.Dialog.Error().
			SetMessage("This is a custom message").
			Show()
	})
	errorMenu.Add("Error (Custom Icon)").OnClick(func(ctx *application.Context) {
		dialog := app.Dialog.Error()
		dialog.SetTitle("Custom Icon Example")
		dialog.SetMessage("Using a custom icon")
		dialog.SetIcon(icons.WailsLogoWhite)
		dialog.Show()
	})

	openMenu := menu.AddSubmenu("Open")
	openMenu.Add("Open File").OnClick(func(ctx *application.Context) {
		result, _ := app.Dialog.OpenFile().
			CanChooseFiles(true).
			PromptForSingleSelection()
		if result != "" {
			app.Dialog.Info().SetMessage(result).Show()
		} else {
			app.Dialog.Info().SetMessage("No file selected").Show()
		}
	})
	openMenu.Add("Open File (Show Hidden Files)").OnClick(func(ctx *application.Context) {
		result, _ := app.Dialog.OpenFile().
			CanChooseFiles(true).
			CanCreateDirectories(true).
			ShowHiddenFiles(true).
			PromptForSingleSelection()
		if result != "" {
			app.Dialog.Info().SetMessage(result).Show()
		} else {
			app.Dialog.Info().SetMessage("No file selected").Show()
		}
	})
	openMenu.Add("Open File (Attach to window)").OnClick(func(ctx *application.Context) {
		result, _ := app.Dialog.OpenFile().
			CanChooseFiles(true).
			CanCreateDirectories(true).
			ShowHiddenFiles(true).
			AttachToWindow(app.Window.Current()).
			PromptForSingleSelection()
		if result != "" {
			app.Dialog.Info().SetMessage(result).Show()
		} else {
			app.Dialog.Info().SetMessage("No file selected").Show()
		}
	})
	openMenu.Add("Open Multiple Files (Show Hidden Files)").OnClick(func(ctx *application.Context) {
		result, _ := app.Dialog.OpenFile().
			CanChooseFiles(true).
			CanCreateDirectories(true).
			ShowHiddenFiles(true).
			PromptForMultipleSelection()
		if len(result) > 0 {
			app.Dialog.Info().SetMessage(strings.Join(result, ",")).Show()
		} else {
			app.Dialog.Info().SetMessage("No file selected").Show()
		}
	})
	openMenu.Add("Open Directory").OnClick(func(ctx *application.Context) {
		result, _ := app.Dialog.OpenFile().
			CanChooseDirectories(true).
			CanChooseFiles(false).
			PromptForSingleSelection()
		if result != "" {
			app.Dialog.Info().SetMessage(result).Show()
		} else {
			app.Dialog.Info().SetMessage("No directory selected").Show()
		}
	})
	openMenu.Add("Open Directory (Create Directories)").OnClick(func(ctx *application.Context) {
		result, _ := app.Dialog.OpenFile().
			CanChooseDirectories(true).
			CanCreateDirectories(true).
			CanChooseFiles(false).
			PromptForSingleSelection()
		if result != "" {
			app.Dialog.Info().SetMessage(result).Show()
		} else {
			app.Dialog.Info().SetMessage("No directory selected").Show()
		}
	})
	openMenu.Add("Open Directory (Resolves Aliases)").OnClick(func(ctx *application.Context) {
		result, _ := app.Dialog.OpenFile().
			CanChooseDirectories(true).
			CanCreateDirectories(true).
			CanChooseFiles(false).
			ResolvesAliases(true).
			PromptForSingleSelection()
		if result != "" {
			app.Dialog.Info().SetMessage(result).Show()
		} else {
			app.Dialog.Info().SetMessage("No directory selected").Show()
		}
	})
	openMenu.Add("Open File/Directory (Set Title)").OnClick(func(ctx *application.Context) {
		dialog := app.Dialog.OpenFile().
			CanChooseDirectories(true).
			CanCreateDirectories(true).
			ResolvesAliases(true)
		if runtime.GOOS == "darwin" {
			dialog.SetMessage("Select a file/directory")
		} else {
			dialog.SetTitle("Select a file/directory")
		}

		result, _ := dialog.PromptForSingleSelection()
		if result != "" {
			app.Dialog.Info().SetMessage(result).Show()
		} else {
			app.Dialog.Info().SetMessage("No file/directory selected").Show()
		}
	})
	openMenu.Add("Open (Full Example)").OnClick(func(ctx *application.Context) {
		cwd, _ := os.Getwd()
		dialog := app.Dialog.OpenFile().
			SetTitle("Select a file").
			SetMessage("Select a file to open").
			SetButtonText("Let's do this!").
			SetDirectory(cwd).
			CanCreateDirectories(true).
			ResolvesAliases(true).
			AllowsOtherFileTypes(true).
			TreatsFilePackagesAsDirectories(true).
			ShowHiddenFiles(true).
			CanSelectHiddenExtension(true).
			AddFilter("Text Files", "*.txt; *.md").
			AddFilter("Video Files", "*.mov; *.mp4; *.avi")

		if runtime.GOOS == "darwin" {
			dialog.SetMessage("Select a file")
		} else {
			dialog.SetTitle("Select a file")
		}

		result, _ := dialog.PromptForSingleSelection()
		if result != "" {
			app.Dialog.Info().SetMessage(result).Show()
		} else {
			app.Dialog.Info().SetMessage("No file selected").Show()
		}
	})

	saveMenu := menu.AddSubmenu("Save")
	saveMenu.Add("Select File (Defaults)").OnClick(func(ctx *application.Context) {
		result, _ := app.Dialog.SaveFile().
			PromptForSingleSelection()
		if result != "" {
			app.Dialog.Info().SetMessage(result).Show()
		}
	})
	saveMenu.Add("Select File (Attach To WebviewWindow)").OnClick(func(ctx *application.Context) {
		result, _ := app.Dialog.SaveFile().
			AttachToWindow(app.Window.Current()).
			PromptForSingleSelection()
		if result != "" {
			app.Dialog.Info().SetMessage(result).Show()
		}
	})
	saveMenu.Add("Select File (Show Hidden Files)").OnClick(func(ctx *application.Context) {
		result, _ := app.Dialog.SaveFile().
			ShowHiddenFiles(true).
			PromptForSingleSelection()
		if result != "" {
			app.Dialog.Info().SetMessage(result).Show()
		}
	})
	saveMenu.Add("Select File (Cannot Create Directories)").OnClick(func(ctx *application.Context) {
		result, _ := app.Dialog.SaveFile().
			CanCreateDirectories(false).
			PromptForSingleSelection()
		if result != "" {
			app.Dialog.Info().SetMessage(result).Show()
		}
	})
	saveMenu.Add("Select File (Full Example)").OnClick(func(ctx *application.Context) {
		result, _ := app.Dialog.SaveFile().
			CanCreateDirectories(false).
			ShowHiddenFiles(true).
			SetMessage("Select a file").
			SetDirectory("/Applications").
			SetButtonText("Let's do this!").
			SetFilename("README.md").
			HideExtension(true).
			AllowsOtherFileTypes(true).
			TreatsFilePackagesAsDirectories(true).
			ShowHiddenFiles(true).
			PromptForSingleSelection()
		if result != "" {
			app.Dialog.Info().SetMessage(result).Show()
		}
	})

	app.Menu.Set(menu)

	app.Window.New()

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
