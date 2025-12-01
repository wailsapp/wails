//go:build android

package application

import "encoding/json"

// dialogsImpl implements dialogs for Android
type dialogsImpl struct {
	// Android-specific fields if needed
}

func newDialogsImpl() *dialogsImpl {
	return &dialogsImpl{}
}

// getDialogTypeString converts DialogType to string for JNI
func getDialogTypeString(dialogType DialogType) string {
	switch dialogType {
	case InfoDialogType:
		return "info"
	case WarningDialogType:
		return "warning"
	case ErrorDialogType:
		return "error"
	case QuestionDialogType:
		return "question"
	default:
		return "info"
	}
}

// getButtonLabels extracts button labels from MessageDialogOptions
func getButtonLabels(param MessageDialogOptions) []string {
	labels := make([]string, 0, len(param.Buttons))
	for _, btn := range param.Buttons {
		labels = append(labels, btn.Label)
	}
	return labels
}

// findButtonCallback finds the callback for the clicked button
func findButtonCallback(param MessageDialogOptions, clickedLabel string) func() {
	for _, btn := range param.Buttons {
		if btn.Label == clickedLabel {
			return btn.Callback
		}
	}
	return nil
}

func (d *dialogsImpl) info(id uint, param MessageDialogOptions) {
	defer freeDialogID(id)

	// Get button labels or use default "OK"
	buttonLabels := getButtonLabels(param)
	if len(buttonLabels) == 0 {
		buttonLabels = []string{"OK"}
	}

	// Convert to JSON array
	buttonsJSON, _ := json.Marshal(buttonLabels)

	// Show the dialog
	dialogType := getDialogTypeString(param.DialogType)
	result := AndroidShowMessageDialog(dialogType, param.Title, param.Message, string(buttonsJSON))

	// Execute callback if button was clicked
	if result != "" {
		if callback := findButtonCallback(param, result); callback != nil {
			callback()
		}
	}
}

func (d *dialogsImpl) warning(id uint, param MessageDialogOptions) {
	defer freeDialogID(id)

	// Get button labels or use default "OK"
	buttonLabels := getButtonLabels(param)
	if len(buttonLabels) == 0 {
		buttonLabels = []string{"OK"}
	}

	// Convert to JSON array
	buttonsJSON, _ := json.Marshal(buttonLabels)

	// Show the dialog
	dialogType := getDialogTypeString(param.DialogType)
	result := AndroidShowMessageDialog(dialogType, param.Title, param.Message, string(buttonsJSON))

	// Execute callback if button was clicked
	if result != "" {
		if callback := findButtonCallback(param, result); callback != nil {
			callback()
		}
	}
}

func (d *dialogsImpl) error(id uint, param MessageDialogOptions) {
	defer freeDialogID(id)

	// Get button labels or use default "OK"
	buttonLabels := getButtonLabels(param)
	if len(buttonLabels) == 0 {
		buttonLabels = []string{"OK"}
	}

	// Convert to JSON array
	buttonsJSON, _ := json.Marshal(buttonLabels)

	// Show the dialog
	dialogType := getDialogTypeString(param.DialogType)
	result := AndroidShowMessageDialog(dialogType, param.Title, param.Message, string(buttonsJSON))

	// Execute callback if button was clicked
	if result != "" {
		if callback := findButtonCallback(param, result); callback != nil {
			callback()
		}
	}
}

func (d *dialogsImpl) question(id uint, param MessageDialogOptions) chan bool {
	defer freeDialogID(id)

	ch := make(chan bool, 1)

	// Get button labels or use default "Yes", "No"
	buttonLabels := getButtonLabels(param)
	if len(buttonLabels) == 0 {
		buttonLabels = []string{"Yes", "No"}
	}

	// Convert to JSON array
	buttonsJSON, _ := json.Marshal(buttonLabels)

	// Show the dialog
	dialogType := getDialogTypeString(param.DialogType)
	result := AndroidShowMessageDialog(dialogType, param.Title, param.Message, string(buttonsJSON))

	// Determine boolean result based on clicked button
	// The first button is considered "Yes" (true), others are "No" (false)
	isYes := false
	if result != "" && len(buttonLabels) > 0 && result == buttonLabels[0] {
		isYes = true
	}

	// Execute callback if button was clicked
	if result != "" {
		if callback := findButtonCallback(param, result); callback != nil {
			callback()
		}
	}

	ch <- isYes
	return ch
}

func (d *dialogsImpl) openFile(id uint, param OpenFileDialogOptions) chan string {
	// TODO: Implement using Android file picker intent
	ch := make(chan string, 1)
	ch <- ""
	return ch
}

func (d *dialogsImpl) openMultipleFiles(id uint, param OpenFileDialogOptions) chan []string {
	// TODO: Implement using Android file picker intent
	ch := make(chan []string, 1)
	ch <- []string{}
	return ch
}

func (d *dialogsImpl) openDirectory(id uint, param OpenFileDialogOptions) chan string {
	// TODO: Implement using Android file picker intent
	ch := make(chan string, 1)
	ch <- ""
	return ch
}

func (d *dialogsImpl) saveFile(id uint, param SaveFileDialogOptions) chan string {
	// TODO: Implement using Android file picker intent
	ch := make(chan string, 1)
	ch <- ""
	return ch
}

type androidDialog struct {
	dialog *MessageDialog
}

func (d *androidDialog) show() {
	// TODO: Implement using AlertDialog
}

func newDialogImpl(d *MessageDialog) *androidDialog {
	return &androidDialog{
		dialog: d,
	}
}

func (d *dialogsImpl) show() (chan string, error) {
	ch := make(chan string, 1)
	ch <- ""
	return ch, nil
}

func newOpenFileDialogImpl(_ *OpenFileDialogStruct) openFileDialogImpl {
	return &dialogsImpl{}
}

func newSaveFileDialogImpl(_ *SaveFileDialogStruct) saveFileDialogImpl {
	return &dialogsImpl{}
}
