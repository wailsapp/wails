package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	app := application.New(application.Options{
		Name:        "Dialog Test",
		Description: "Test application for macOS dialogs",
		Logger:      application.DefaultLogger(slog.LevelDebug),
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Create main window
	mainWindow := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:     "Dialog Tests",
		Width:     800,
		Height:    600,
		MinWidth:  800,
		MinHeight: 600,
	})
	mainWindow.SetAlwaysOnTop(true)

	// Create main menu
	menu := app.NewMenu()
	app.SetMenu(menu)
	menu.AddRole(application.AppMenu)
	menu.AddRole(application.EditMenu)
	menu.AddRole(application.WindowMenu)

	// Add test menu
	testMenu := menu.AddSubmenu("Tests")

	// Test 1: Basic file open with no filters (no window)
	testMenu.Add("1. Basic Open (No Window)").OnClick(func(ctx *application.Context) {
		result, err := application.OpenFileDialog().
			CanChooseFiles(true).
			PromptForSingleSelection()
		showResult("Basic Open", result, err, nil)
	})

	// Test 1b: Basic file open with window
	testMenu.Add("1b. Basic Open (With Window)").OnClick(func(ctx *application.Context) {
		result, err := application.OpenFileDialog().
			CanChooseFiles(true).
			AttachToWindow(mainWindow).
			PromptForSingleSelection()
		showResult("Basic Open", result, err, mainWindow)
	})

	// Test 2: Open with single extension filter
	testMenu.Add("2. Single Filter").OnClick(func(ctx *application.Context) {
		result, err := application.OpenFileDialog().
			CanChooseFiles(true).
			AddFilter("Text Files", "*.txt").
			AttachToWindow(mainWindow).
			PromptForSingleSelection()
		showResult("Single Filter", result, err, mainWindow)
	})

	// Test 3: Open with multiple extension filter
	testMenu.Add("3. Multiple Filter").OnClick(func(ctx *application.Context) {
		result, err := application.OpenFileDialog().
			CanChooseFiles(true).
			AddFilter("Documents", "*.txt;*.md;*.doc;*.docx").
			AttachToWindow(mainWindow).
			PromptForSingleSelection()
		showResult("Multiple Filter", result, err, mainWindow)
	})

	// Test 4: Multiple file selection
	testMenu.Add("4. Multiple Selection").OnClick(func(ctx *application.Context) {
		results, err := application.OpenFileDialog().
			CanChooseFiles(true).
			AddFilter("Images", "*.png;*.jpg;*.jpeg").
			AttachToWindow(mainWindow).
			PromptForMultipleSelection()
		if err != nil {
			showError("Multiple Selection", err, mainWindow)
			return
		}
		showResults("Multiple Selection", results, mainWindow)
	})

	// Test 5: Directory selection
	testMenu.Add("5. Directory Selection").OnClick(func(ctx *application.Context) {
		result, err := application.OpenFileDialog().
			CanChooseDirectories(true).
			CanChooseFiles(false).
			AttachToWindow(mainWindow).
			PromptForSingleSelection()
		showResult("Directory Selection", result, err, mainWindow)
	})

	// Test 6: Save dialog with extension
	testMenu.Add("6. Save Dialog").OnClick(func(ctx *application.Context) {
		result, err := application.SaveFileDialog().
			SetFilename("test.txt").
			AddFilter("Text Files", "*.txt").
			AttachToWindow(mainWindow).
			PromptForSingleSelection()
		showResult("Save Dialog", result, err, mainWindow)
	})

	// Test 7: Complex filters
	testMenu.Add("7. Complex Filters").OnClick(func(ctx *application.Context) {
		result, err := application.OpenFileDialog().
			CanChooseFiles(true).
			AddFilter("All Documents", "*.txt;*.md;*.doc;*.docx;*.pdf").
			AddFilter("Text Files", "*.txt").
			AddFilter("Markdown", "*.md").
			AddFilter("Word Documents", "*.doc;*.docx").
			AddFilter("PDF Files", "*.pdf").
			AttachToWindow(mainWindow).
			PromptForSingleSelection()
		showResult("Complex Filters", result, err, mainWindow)
	})

	// Test 8: Hidden files
	testMenu.Add("8. Show Hidden").OnClick(func(ctx *application.Context) {
		result, err := application.OpenFileDialog().
			CanChooseFiles(true).
			ShowHiddenFiles(true).
			AttachToWindow(mainWindow).
			PromptForSingleSelection()
		showResult("Show Hidden", result, err, mainWindow)
	})

	// Test 9: Default directory
	testMenu.Add("9. Default Directory").OnClick(func(ctx *application.Context) {
		home, _ := os.UserHomeDir()
		result, err := application.OpenFileDialog().
			CanChooseFiles(true).
			SetDirectory(home).
			AttachToWindow(mainWindow).
			PromptForSingleSelection()
		showResult("Default Directory", result, err, mainWindow)
	})

	// Test 10: Full featured dialog
	testMenu.Add("10. Full Featured").OnClick(func(ctx *application.Context) {
		home, _ := os.UserHomeDir()
		dialog := application.OpenFileDialog().
			SetTitle("Full Featured Dialog").
			SetDirectory(home).
			CanChooseFiles(true).
			CanCreateDirectories(true).
			ShowHiddenFiles(true).
			ResolvesAliases(true).
			AllowsOtherFileTypes(true).
			AttachToWindow(mainWindow)

		if runtime.GOOS == "darwin" {
			dialog.SetMessage("Please select files")
		}

		dialog.AddFilter("All Supported", "*.txt;*.md;*.pdf;*.png;*.jpg")
		dialog.AddFilter("Documents", "*.txt;*.md;*.pdf")
		dialog.AddFilter("Images", "*.png;*.jpg;*.jpeg")

		results, err := dialog.PromptForMultipleSelection()
		if err != nil {
			showError("Full Featured", err, mainWindow)
			return
		}
		showResults("Full Featured", results, mainWindow)
	})

	// Show the window
	mainWindow.Show()

	// Run the app
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func showResult(test string, result string, err error, window *application.WebviewWindow) {
	if err != nil {
		showError(test, err, window)
		return
	}
	if result == "" {
		dialog := application.InfoDialog().
			SetTitle(test).
			SetMessage("No file selected")
		if window != nil {
			dialog.AttachToWindow(window)
		}
		dialog.Show()
		return
	}
	dialog := application.InfoDialog().
		SetTitle(test).
		SetMessage(fmt.Sprintf("Selected: %s\nType: %s", result, getFileType(result)))
	if window != nil {
		dialog.AttachToWindow(window)
	}
	dialog.Show()
}

func showResults(test string, results []string, window *application.WebviewWindow) {
	if len(results) == 0 {
		dialog := application.InfoDialog().
			SetTitle(test).
			SetMessage("No files selected")
		if window != nil {
			dialog.AttachToWindow(window)
		}
		dialog.Show()
		return
	}
	var message strings.Builder
	message.WriteString(fmt.Sprintf("Selected %d files:\n\n", len(results)))
	for _, result := range results {
		message.WriteString(fmt.Sprintf("%s (%s)\n", result, getFileType(result)))
	}
	dialog := application.InfoDialog().
		SetTitle(test).
		SetMessage(message.String())
	if window != nil {
		dialog.AttachToWindow(window)
	}
	dialog.Show()
}

func showError(test string, err error, window *application.WebviewWindow) {
	dialog := application.ErrorDialog().
		SetTitle(test).
		SetMessage(fmt.Sprintf("Error: %v", err))
	if window != nil {
		dialog.AttachToWindow(window)
	}
	dialog.Show()
}

func getFileType(path string) string {
	if path == "" {
		return "unknown"
	}
	ext := strings.ToLower(filepath.Ext(path))
	if ext == "" {
		if fi, err := os.Stat(path); err == nil && fi.IsDir() {
			return "directory"
		}
		return "no extension"
	}
	return ext
}
