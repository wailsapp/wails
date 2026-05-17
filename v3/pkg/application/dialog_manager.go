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
