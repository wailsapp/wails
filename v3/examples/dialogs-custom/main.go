package main

import (
	"log"
	"log/slog"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {

	app := application.New(application.Options{
		Name:        "Custom Dialogs Demo",
		Description: "A demo of custom dialog functionality with custom buttons and icons",
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

	// Custom Dialog Examples
	customMenu := menu.AddSubmenu("Custom Dialogs")

	// Example 1: Custom button text
	customMenu.Add("Custom Button Text").OnClick(func(ctx *application.Context) {
		dialog := application.QuestionDialog()
		dialog.SetTitle("Save Changes?")
		dialog.SetMessage("You have unsaved changes. What would you like to do?")
		
		saveBtn := dialog.AddButton("💾 Save Now")
		saveBtn.OnClick(func() {
			application.InfoDialog().SetMessage("Changes saved!").Show()
		})
		
		discardBtn := dialog.AddButton("🗑️ Discard Changes")
		discardBtn.OnClick(func() {
			application.InfoDialog().SetMessage("Changes discarded!").Show()
		})
		
		cancelBtn := dialog.AddButton("❌ Cancel")
		cancelBtn.OnClick(func() {
			application.InfoDialog().SetMessage("Operation cancelled!").Show()
		})
		
		dialog.SetDefaultButton(saveBtn)
		dialog.SetCancelButton(cancelBtn)
		dialog.Show()
	})

	// Example 2: Action buttons with emojis
	customMenu.Add("Action Buttons").OnClick(func(ctx *application.Context) {
		dialog := application.InfoDialog()
		dialog.SetTitle("Choose Your Action")
		dialog.SetMessage("What would you like to do next?")
		
		downloadBtn := dialog.AddButton("📥 Download")
		downloadBtn.OnClick(func() {
			application.InfoDialog().SetMessage("Download started...").Show()
		})
		
		uploadBtn := dialog.AddButton("📤 Upload")
		uploadBtn.OnClick(func() {
			application.InfoDialog().SetMessage("Upload started...").Show()
		})
		
		settingsBtn := dialog.AddButton("⚙️ Settings")
		settingsBtn.OnClick(func() {
			application.InfoDialog().SetMessage("Opening settings...").Show()
		})
		
		exitBtn := dialog.AddButton("🚪 Exit")
		exitBtn.OnClick(func() {
			app.Quit()
		})
		
		dialog.SetDefaultButton(downloadBtn)
		dialog.Show()
	})

	// Example 3: Confirmation with custom text
	customMenu.Add("Custom Confirmation").OnClick(func(ctx *application.Context) {
		dialog := application.WarningDialog()
		dialog.SetTitle("Delete Files")
		dialog.SetMessage("This will permanently delete 42 files from your system. This action cannot be undone.")
		
		deleteBtn := dialog.AddButton("🔥 Delete Permanently")
		deleteBtn.OnClick(func() {
			application.InfoDialog().SetMessage("Files deleted successfully!").Show()
		})
		
		backupBtn := dialog.AddButton("💾 Backup First")
		backupBtn.OnClick(func() {
			application.InfoDialog().SetMessage("Creating backup...").Show()
		})
		
		cancelBtn := dialog.AddButton("🛡️ Keep Files")
		cancelBtn.OnClick(func() {
			application.InfoDialog().SetMessage("Files kept safe!").Show()
		})
		
		dialog.SetDefaultButton(backupBtn)
		dialog.SetCancelButton(cancelBtn)
		dialog.Show()
	})

	// Example 4: Multiple choice
	customMenu.Add("Multiple Choice").OnClick(func(ctx *application.Context) {
		dialog := application.QuestionDialog()
		dialog.SetTitle("Select Theme")
		dialog.SetMessage("Which theme would you prefer for the application?")
		
		lightBtn := dialog.AddButton("☀️ Light Theme")
		lightBtn.OnClick(func() {
			application.InfoDialog().SetMessage("Light theme selected!").Show()
		})
		
		darkBtn := dialog.AddButton("🌙 Dark Theme")
		darkBtn.OnClick(func() {
			application.InfoDialog().SetMessage("Dark theme selected!").Show()
		})
		
		autoBtn := dialog.AddButton("🔄 Auto (System)")
		autoBtn.OnClick(func() {
			application.InfoDialog().SetMessage("Auto theme selected!").Show()
		})
		
		highContrastBtn := dialog.AddButton("🔳 High Contrast")
		highContrastBtn.OnClick(func() {
			application.InfoDialog().SetMessage("High contrast theme selected!").Show()
		})
		
		dialog.SetDefaultButton(autoBtn)
		dialog.Show()
	})

	// Example 5: Game-style dialog
	customMenu.Add("Game Style").OnClick(func(ctx *application.Context) {
		dialog := application.InfoDialog()
		dialog.SetTitle("Level Complete!")
		dialog.SetMessage("Congratulations! You've completed Level 5 with a score of 9,847 points.")
		
		nextBtn := dialog.AddButton("▶️ Next Level")
		nextBtn.OnClick(func() {
			application.InfoDialog().SetMessage("Loading Level 6...").Show()
		})
		
		replayBtn := dialog.AddButton("🔄 Replay Level")
		replayBtn.OnClick(func() {
			application.InfoDialog().SetMessage("Restarting Level 5...").Show()
		})
		
		menuBtn := dialog.AddButton("🏠 Main Menu")
		menuBtn.OnClick(func() {
			application.InfoDialog().SetMessage("Returning to main menu...").Show()
		})
		
		shareBtn := dialog.AddButton("📱 Share Score")
		shareBtn.OnClick(func() {
			application.InfoDialog().SetMessage("Opening share options...").Show()
		})
		
		dialog.SetDefaultButton(nextBtn)
		dialog.Show()
	})

	// Traditional dialogs for comparison
	traditionalMenu := menu.AddSubmenu("Traditional Dialogs")

	traditionalMenu.Add("Standard Question").OnClick(func(ctx *application.Context) {
		dialog := application.QuestionDialog()
		dialog.SetMessage("Do you want to continue?")
		dialog.AddButton("Yes").OnClick(func() {
			application.InfoDialog().SetMessage("Continuing...").Show()
		})
		dialog.AddButton("No").OnClick(func() {
			application.InfoDialog().SetMessage("Cancelled.").Show()
		})
		dialog.Show()
	})

	traditionalMenu.Add("Standard Info").OnClick(func(ctx *application.Context) {
		application.InfoDialog().
			SetTitle("Information").
			SetMessage("This is a standard information dialog.").
			Show()
	})

	app.Menu.Set(menu)

	app.Window.New()

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}