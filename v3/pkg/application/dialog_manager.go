package application

// DialogManager manages dialog-related operations
type DialogManager struct {
	app *App
}

// newDialogManager creates a new DialogManager instance
func newDialogManager(app *App) *DialogManager {
	return &DialogManager{
		app: app,
	}
}

// OpenFile creates a file dialog for selecting files
func (dm *DialogManager) OpenFile() *OpenFileDialogStruct {
	return newOpenFileDialog()
}

// OpenFileWithOptions creates a file dialog with options
func (dm *DialogManager) OpenFileWithOptions(options *OpenFileDialogOptions) *OpenFileDialogStruct {
	result := newOpenFileDialog()
	result.SetOptions(options)
	return result
}

// SaveFile creates a save file dialog
func (dm *DialogManager) SaveFile() *SaveFileDialogStruct {
	return newSaveFileDialog()
}

// SaveFileWithOptions creates a save file dialog with options
func (dm *DialogManager) SaveFileWithOptions(options *SaveFileDialogOptions) *SaveFileDialogStruct {
	result := newSaveFileDialog()
	result.SetOptions(options)
	return result
}

// Info creates an information dialog
func (dm *DialogManager) Info() *MessageDialog {
	return newMessageDialog(InfoDialogType)
}

// Question creates a question dialog
func (dm *DialogManager) Question() *MessageDialog {
	return newMessageDialog(QuestionDialogType)
}

// Warning creates a warning dialog
func (dm *DialogManager) Warning() *MessageDialog {
	return newMessageDialog(WarningDialogType)
}

// Error creates an error dialog
func (dm *DialogManager) Error() *MessageDialog {
	return newMessageDialog(ErrorDialogType)
}

// Prompt shows a native alert with a text field and blocks until it is
// dismissed. It returns the entered value and true when the OK button was
// pressed, or "" and false when cancelled. The alert is a sheet when
// options.Window is set. Call it from a goroutine, not from the main
// thread. On platforms without a native implementation it returns
// ErrDialogNotSupported.
func (dm *DialogManager) Prompt(options PromptOptions) (string, bool, error) {
	return dialogPrompt(options)
}

// PickColor opens the system colour panel and blocks until it is closed.
// It returns the final colour and true when the user changed the colour,
// or the initial colour and false when the panel was closed untouched.
// Call it from a goroutine, not from the main thread. On platforms without
// a native colour panel it returns ErrDialogNotSupported.
func (dm *DialogManager) PickColor(options ColorPickerOptions) (RGBA, bool, error) {
	return dialogPickColor(options)
}

// PickFont opens the system font panel and blocks until it is closed. It
// returns the final font and true when the user changed the selection, or
// the initial font and false when the panel was closed untouched. Call it
// from a goroutine, not from the main thread. On platforms without a native
// font panel it returns ErrDialogNotSupported.
func (dm *DialogManager) PickFont(options FontPickerOptions) (FontDescriptor, bool, error) {
	return dialogPickFont(options)
}
