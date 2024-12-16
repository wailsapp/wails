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
		dialog := application.InfoDialog()
		dialog.SetTitle("Custom Title")
		dialog.SetMessage("This is a custom message")
		dialog.Show()
	})

	infoMenu.Add("Info (Title only)").OnClick(func(ctx *application.Context) {
		dialog := application.InfoDialog()
		dialog.SetTitle("Custom Title")
		dialog.Show()
	})
	infoMenu.Add("Info (Message only)").OnClick(func(ctx *application.Context) {
		dialog := application.InfoDialog()
		dialog.SetMessage("This is a custom message")
		dialog.Show()
	})
	infoMenu.Add("Info (Custom Icon)").OnClick(func(ctx *application.Context) {
		dialog := application.InfoDialog()
		dialog.SetTitle("Custom Icon Example")
		dialog.SetMessage("Using a custom icon")
		dialog.SetIcon(icons.ApplicationDarkMode256)
		dialog.Show()
	})
	infoMenu.Add("About").OnClick(func(ctx *application.Context) {
		app.ShowAboutDialog()
	})

	questionMenu := menu.AddSubmenu("Question")
	questionMenu.Add("Question (No default)").OnClick(func(ctx *application.Context) {
		dialog := application.QuestionDialog()
		dialog.SetMessage("No default button")
		dialog.AddButton("Yes")
		dialog.AddButton("No")
		dialog.Show()
	})
	questionMenu.Add("Question (Attached to Window)").OnClick(func(ctx *application.Context) {
		dialog := application.QuestionDialog()
		dialog.AttachToWindow(app.CurrentWindow())
		dialog.SetMessage("No default button")
		dialog.AddButton("Yes")
		dialog.AddButton("No")
		dialog.Show()
	})
	questionMenu.Add("Question (With Default)").OnClick(func(ctx *application.Context) {
		dialog := application.QuestionDialog()
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
		dialog := application.QuestionDialog().
			SetTitle("Update").
			SetMessage("The cancel button is selected when pressing escape")
		download := dialog.AddButton("📥 Download")
		download.OnClick(func() {
			application.InfoDialog().SetMessage("Downloading...").Show()
		})
		no := dialog.AddButton("Cancel")
		dialog.SetDefaultButton(download)
		dialog.SetCancelButton(no)
		dialog.Show()
	})
	questionMenu.Add("Question (Custom Icon)").OnClick(func(ctx *application.Context) {
		dialog := application.QuestionDialog()
		dialog.SetTitle("Custom Icon Example")
		dialog.SetMessage("Using a custom icon")
		dialog.SetIcon(icons.WailsLogoWhiteTransparent)
		likeIt := dialog.AddButton("I like it!").OnClick(func() {
			application.InfoDialog().SetMessage("Thanks!").Show()
		})
		dialog.AddButton("Not so keen...").OnClick(func() {
			application.InfoDialog().SetMessage("Too bad!").Show()
		})
		dialog.SetDefaultButton(likeIt)
		dialog.Show()
	})

	warningMenu := menu.AddSubmenu("Warning")
	warningMenu.Add("Warning").OnClick(func(ctx *application.Context) {
		application.WarningDialog().
			SetTitle("Custom Title").
			SetMessage("This is a custom message").
			Show()
	})
	warningMenu.Add("Warning (Title only)").OnClick(func(ctx *application.Context) {
		dialog := application.WarningDialog()
		dialog.SetTitle("Custom Title")
		dialog.Show()
	})
	warningMenu.Add("Warning (Message only)").OnClick(func(ctx *application.Context) {
		dialog := application.WarningDialog()
		dialog.SetMessage("This is a custom message")
		dialog.Show()
	})
	warningMenu.Add("Warning (Custom Icon)").OnClick(func(ctx *application.Context) {
		dialog := application.WarningDialog()
		dialog.SetTitle("Custom Icon Example")
		dialog.SetMessage("Using a custom icon")
		dialog.SetIcon(icons.ApplicationLightMode256)
		dialog.Show()
	})

	errorMenu := menu.AddSubmenu("Error")
	errorMenu.Add("Error").OnClick(func(ctx *application.Context) {
		dialog := application.ErrorDialog()
		dialog.SetTitle("Ooops")
		dialog.SetMessage("I accidentally the whole of Twitter")
		dialog.Show()
	})
	errorMenu.Add("Error (Title Only)").OnClick(func(ctx *application.Context) {
		dialog := application.ErrorDialog()
		dialog.SetTitle("Custom Title")
		dialog.Show()
	})
	errorMenu.Add("Error (Custom Message)").OnClick(func(ctx *application.Context) {
		application.ErrorDialog().
			SetMessage("This is a custom message").
			Show()
	})
	errorMenu.Add("Error (Custom Icon)").OnClick(func(ctx *application.Context) {
		dialog := application.ErrorDialog()
		dialog.SetTitle("Custom Icon Example")
		dialog.SetMessage("Using a custom icon")
		dialog.SetIcon(icons.WailsLogoWhite)
		dialog.Show()
	})

	openMenu := menu.AddSubmenu("Open")
	openMenu.Add("Open File").OnClick(func(ctx *application.Context) {
		result, _ := application.OpenFileDialog().
			CanChooseFiles(true).
			PromptForSingleSelection()
		if result != "" {
			application.InfoDialog().SetMessage(result).Show()
		} else {
			application.InfoDialog().SetMessage("No file selected").Show()
		}
	})
	openMenu.Add("Open File (Show Hidden Files)").OnClick(func(ctx *application.Context) {
		result, _ := application.OpenFileDialog().
			CanChooseFiles(true).
			CanCreateDirectories(true).
			ShowHiddenFiles(true).
			PromptForSingleSelection()
		if result != "" {
			application.InfoDialog().SetMessage(result).Show()
		} else {
			application.InfoDialog().SetMessage("No file selected").Show()
		}
	})
	openMenu.Add("Open File (Attach to window)").OnClick(func(ctx *application.Context) {
		result, _ := application.OpenFileDialog().
			CanChooseFiles(true).
			CanCreateDirectories(true).
			ShowHiddenFiles(true).
			AttachToWindow(app.CurrentWindow()).
			PromptForSingleSelection()
		if result != "" {
			application.InfoDialog().SetMessage(result).Show()
		} else {
			application.InfoDialog().SetMessage("No file selected").Show()
		}
	})
	openMenu.Add("Open Multiple Files (Show Hidden Files)").OnClick(func(ctx *application.Context) {
		result, _ := application.OpenFileDialog().
			CanChooseFiles(true).
			CanCreateDirectories(true).
			ShowHiddenFiles(true).
			PromptForMultipleSelection()
		if len(result) > 0 {
			application.InfoDialog().SetMessage(strings.Join(result, ",")).Show()
		} else {
			application.InfoDialog().SetMessage("No file selected").Show()
		}
	})
	openMenu.Add("Open Directory").OnClick(func(ctx *application.Context) {
		result, _ := application.OpenFileDialog().
			CanChooseDirectories(true).
			CanChooseFiles(false).
			PromptForSingleSelection()
		if result != "" {
			application.InfoDialog().SetMessage(result).Show()
		} else {
			application.InfoDialog().SetMessage("No directory selected").Show()
		}
	})
	openMenu.Add("Open Directory (Create Directories)").OnClick(func(ctx *application.Context) {
		result, _ := application.OpenFileDialog().
			CanChooseDirectories(true).
			CanCreateDirectories(true).
			CanChooseFiles(false).
			PromptForSingleSelection()
		if result != "" {
			application.InfoDialog().SetMessage(result).Show()
		} else {
			application.InfoDialog().SetMessage("No directory selected").Show()
		}
	})
	openMenu.Add("Open Directory (Resolves Aliases)").OnClick(func(ctx *application.Context) {
		result, _ := application.OpenFileDialog().
			CanChooseDirectories(true).
			CanCreateDirectories(true).
			CanChooseFiles(false).
			ResolvesAliases(true).
			PromptForSingleSelection()
		if result != "" {
			application.InfoDialog().SetMessage(result).Show()
		} else {
			application.InfoDialog().SetMessage("No directory selected").Show()
		}
	})
	openMenu.Add("Open File/Directory (Set Title)").OnClick(func(ctx *application.Context) {
		dialog := application.OpenFileDialog().
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
			application.InfoDialog().SetMessage(result).Show()
		} else {
			application.InfoDialog().SetMessage("No file/directory selected").Show()
		}
	})
	openMenu.Add("Open (Full Example)").OnClick(func(ctx *application.Context) {
		cwd, _ := os.Getwd()
		dialog := application.OpenFileDialog().
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
			application.InfoDialog().SetMessage(result).Show()
		} else {
			application.InfoDialog().SetMessage("No file selected").Show()
		}
	})

	saveMenu := menu.AddSubmenu("Save")
	saveMenu.Add("Select File (Defaults)").OnClick(func(ctx *application.Context) {
		result, _ := application.SaveFileDialog().
			PromptForSingleSelection()
		if result != "" {
			application.InfoDialog().SetMessage(result).Show()
		}
	})
	saveMenu.Add("Select File (Attach To WebviewWindow)").OnClick(func(ctx *application.Context) {
		result, _ := application.SaveFileDialog().
			AttachToWindow(app.CurrentWindow()).
			PromptForSingleSelection()
		if result != "" {
			application.InfoDialog().SetMessage(result).Show()
		}
	})
	saveMenu.Add("Select File (Show Hidden Files)").OnClick(func(ctx *application.Context) {
		result, _ := application.SaveFileDialog().
			ShowHiddenFiles(true).
			PromptForSingleSelection()
		if result != "" {
			application.InfoDialog().SetMessage(result).Show()
		}
	})
	saveMenu.Add("Select File (Cannot Create Directories)").OnClick(func(ctx *application.Context) {
		result, _ := application.SaveFileDialog().
			CanCreateDirectories(false).
			PromptForSingleSelection()
		if result != "" {
			application.InfoDialog().SetMessage(result).Show()
		}
	})
	saveMenu.Add("Select File (Full Example)").OnClick(func(ctx *application.Context) {
		result, _ := application.SaveFileDialog().
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
			application.InfoDialog().SetMessage(result).Show()
		}
	})

	app.SetMenu(menu)

	app.NewWebviewWindow()

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
