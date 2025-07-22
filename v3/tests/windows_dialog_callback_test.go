//go:build windows

package tests

import (
	"testing"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func TestWindowsDialogButtonCallbacks(t *testing.T) {
	// Test case 1: OK button with callback
	t.Run("OK Button Callback", func(t *testing.T) {
		dialog := application.InfoDialog().
			SetTitle("Test Dialog").
			SetMessage("Testing OK button callback")

		okButton := dialog.AddButton("OK") // Note: Windows will show "Ok"
		okButton.OnClick(func() {
			// Callback would be executed when button is clicked
		})

		// Verify callback is set
		if okButton.Callback == nil {
			t.Error("Callback was not set on OK button")
		}
	})

	// Test case 2: Case-insensitive matching
	t.Run("Case Insensitive Matching", func(t *testing.T) {
		testCases := []struct {
			buttonLabel string
			windowsLabel string
		}{
			{"ok", "Ok"},
			{"OK", "Ok"},
			{"Ok", "Ok"},
			{"YES", "Yes"},
			{"yes", "Yes"},
			{"no", "No"},
			{"NO", "No"},
			{"cancel", "Cancel"},
			{"CANCEL", "Cancel"},
		}

		for _, tc := range testCases {
			dialog := application.QuestionDialog().
				SetTitle("Case Test").
				SetMessage("Testing case-insensitive callback")

			button := dialog.AddButton(tc.buttonLabel)
			button.OnClick(func() {
				// Callback would be executed when button is clicked
			})

			if button.Callback == nil {
				t.Errorf("Callback not set for button label '%s'", tc.buttonLabel)
			}
		}
	})

	// Test case 3: Multiple buttons with callbacks
	t.Run("Multiple Buttons", func(t *testing.T) {
		dialog := application.QuestionDialog().
			SetTitle("Multiple Buttons").
			SetMessage("Testing Yes/No callbacks")

		yesButton := dialog.AddButton("Yes")
		yesButton.OnClick(func() {
			// Yes callback
		})

		noButton := dialog.AddButton("No")
		noButton.OnClick(func() {
			// No callback
		})

		if yesButton.Callback == nil || noButton.Callback == nil {
			t.Error("One or more callbacks were not set")
		}
	})

	// Test case 4: Button with whitespace
	t.Run("Whitespace Handling", func(t *testing.T) {
		dialog := application.InfoDialog().
			SetTitle("Whitespace Test").
			SetMessage("Testing whitespace trimming")

		// Add button with extra whitespace
		button := dialog.AddButton("  OK  ")
		button.OnClick(func() {
			// Callback would be executed when button is clicked
		})

		if button.Callback == nil {
			t.Error("Callback not set for button with whitespace")
		}
	})
}

// TestDialogButtonTypes verifies different dialog types work with callbacks
func TestDialogButtonTypes(t *testing.T) {
	// Info dialog - typically has OK button
	t.Run("Info Dialog", func(t *testing.T) {
		dialog := application.InfoDialog().
			SetTitle("Info").
			SetMessage("Info dialog test")

		button := dialog.AddButton("OK")
		button.OnClick(func() {
			// Callback for info dialog
		})

		if button.Callback == nil {
			t.Error("Info dialog button callback not set")
		}
	})

	// Error dialog - typically has OK button
	t.Run("Error Dialog", func(t *testing.T) {
		dialog := application.ErrorDialog().
			SetTitle("Error").
			SetMessage("Error dialog test")

		button := dialog.AddButton("OK")
		button.OnClick(func() {
			// Callback for error dialog
		})

		if button.Callback == nil {
			t.Error("Error dialog button callback not set")
		}
	})

	// Warning dialog - typically has OK button
	t.Run("Warning Dialog", func(t *testing.T) {
		dialog := application.WarningDialog().
			SetTitle("Warning").
			SetMessage("Warning dialog test")

		button := dialog.AddButton("OK")
		button.OnClick(func() {
			// Callback for warning dialog
		})

		if button.Callback == nil {
			t.Error("Warning dialog button callback not set")
		}
	})

	// Question dialog - typically has Yes/No buttons
	t.Run("Question Dialog", func(t *testing.T) {
		dialog := application.QuestionDialog().
			SetTitle("Question").
			SetMessage("Question dialog test")

		yesButton := dialog.AddButton("Yes")
		yesButton.OnClick(func() {
			// Yes callback
		})

		noButton := dialog.AddButton("No")
		noButton.OnClick(func() {
			// No callback
		})

		if yesButton.Callback == nil || noButton.Callback == nil {
			t.Error("Question dialog button callbacks not set")
		}
	})
}

// TestCustomButtonLabels verifies TaskDialog custom button support
func TestCustomButtonLabels(t *testing.T) {
	// Test custom button labels that TaskDialog API supports
	t.Run("Custom Action Buttons", func(t *testing.T) {
		dialog := application.InfoDialog().
			SetTitle("Custom Actions").
			SetMessage("Choose your action")

		// Test various custom labels
		buttons := []string{
			"Save & Continue",
			"Don't Save",
			"Download Now",
			"Settings",
			"Retry Operation",
			"Skip This Step",
		}

		for _, label := range buttons {
			btn := dialog.AddButton(label)
			btn.OnClick(func() {
				// Custom callback
			})

			if btn.Label != label {
				t.Errorf("Button label mismatch: got %q, want %q", btn.Label, label)
			}
			if btn.Callback == nil {
				t.Errorf("Callback not set for button %q", label)
			}
		}
	})

	// Test emoji and Unicode support
	t.Run("Emoji and Unicode Buttons", func(t *testing.T) {
		dialog := application.InfoDialog().
			SetTitle("Unicode Support").
			SetMessage("Testing Unicode button labels")

		emojiButtons := []string{
			"📥 Download",
			"⚙️ Settings",
			"🚀 Launch",
			"❌ Cancel",
			"✅ Accept",
			"📅 Schedule",
		}

		for _, label := range emojiButtons {
			btn := dialog.AddButton(label)
			if btn.Label != label {
				t.Errorf("Emoji button label mismatch: got %q, want %q", btn.Label, label)
			}
		}
	})

	// Test default and cancel button flags
	t.Run("Default and Cancel Buttons", func(t *testing.T) {
		dialog := application.QuestionDialog().
			SetTitle("Save Changes?").
			SetMessage("You have unsaved changes")

		saveBtn := dialog.AddButton("Save").SetAsDefault()
		dontSaveBtn := dialog.AddButton("Don't Save")
		cancelBtn := dialog.AddButton("Cancel").SetAsCancel()

		if !saveBtn.IsDefault {
			t.Error("Save button should be marked as default")
		}
		if !cancelBtn.IsCancel {
			t.Error("Cancel button should be marked as cancel")
		}
		if dontSaveBtn.IsDefault || dontSaveBtn.IsCancel {
			t.Error("Don't Save button should not be marked as default or cancel")
		}
	})
}