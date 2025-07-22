//go:build windows

package tests

import (
	"fmt"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// ExampleWindowsDialogCallbacks demonstrates how to use dialog button callbacks on Windows
// Now with full custom button support using TaskDialog API!
func ExampleWindowsDialogCallbacks() {
	// Example 1: Simple OK button with callback (backwards compatible)
	infoDialog := application.InfoDialog().
		SetTitle("Information").
		SetMessage("Click OK to continue")
	
	okButton := infoDialog.AddButton("OK")
	okButton.OnClick(func() {
		fmt.Println("OK button was clicked!")
	})
	
	// Example 2: Custom buttons with any label (NEW!)
	customDialog := application.InfoDialog().
		SetTitle("Custom Actions").
		SetMessage("Choose your action:")
	
	customDialog.AddButton("📥 Download").OnClick(func() {
		fmt.Println("Download started!")
	})
	
	customDialog.AddButton("⚙️ Settings").OnClick(func() {
		fmt.Println("Opening settings...")
	})
	
	customDialog.AddButton("Cancel").OnClick(func() {
		fmt.Println("Operation cancelled")
	}).SetAsCancel()
	
	// Example 3: Complex dialog with multiple custom buttons
	saveDialog := application.QuestionDialog().
		SetTitle("Save Changes?").
		SetMessage("You have unsaved changes. What would you like to do?")
	
	saveDialog.AddButton("Save").OnClick(func() {
		fmt.Println("Saving changes...")
	}).SetAsDefault()
	
	saveDialog.AddButton("Don't Save").OnClick(func() {
		fmt.Println("Discarding changes...")
	})
	
	saveDialog.AddButton("Cancel").OnClick(func() {
		fmt.Println("Cancelled - returning to editor")
	}).SetAsCancel()
	
	// Example 4: Dialog with emoji and special characters
	updateDialog := application.InfoDialog().
		SetTitle("Update Available").
		SetMessage("A new version is available!")
	
	updateDialog.AddButton("🚀 Update Now").OnClick(func() {
		fmt.Println("Starting update...")
	}).SetAsDefault()
	
	updateDialog.AddButton("📅 Remind Me Later").OnClick(func() {
		fmt.Println("Will remind later...")
	})
	
	updateDialog.AddButton("❌ Skip This Version").OnClick(func() {
		fmt.Println("Skipping this version")
	})
	
	// Example 5: Error dialog with custom retry logic
	errorDialog := application.ErrorDialog().
		SetTitle("Connection Failed").
		SetMessage("Unable to connect to the server")
	
	errorDialog.AddButton("Retry").OnClick(func() {
		fmt.Println("Retrying connection...")
	})
	
	errorDialog.AddButton("Work Offline").OnClick(func() {
		fmt.Println("Switching to offline mode...")
	})
	
	errorDialog.AddButton("Exit").OnClick(func() {
		fmt.Println("Exiting application...")
	})
	
	// Note: With TaskDialog API (Windows Vista+), we now support:
	// - Unlimited custom button labels
	// - Emoji and Unicode in button text
	// - Default and cancel button designation
	// - Full callback support for all buttons
	// 
	// Falls back to MessageBox API on older Windows versions
	// with limited button options (OK, Yes/No, etc.)
}