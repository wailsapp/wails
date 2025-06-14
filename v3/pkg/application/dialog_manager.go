package application

// DialogManager manages dialog-related operations
type DialogManager struct {
	app *App
}

// NewDialogManager creates a new DialogManager instance
func NewDialogManager(app *App) *DialogManager {
	return &DialogManager{
		app: app,
	}
}

// OpenFile opens a file dialog
func (dm *DialogManager) OpenFile() *OpenFileDialogStruct {
	return OpenFileDialog()
}

// OpenFileWithOptions opens a file dialog with options
func (dm *DialogManager) OpenFileWithOptions(options *OpenFileDialogOptions) *OpenFileDialogStruct {
	result := OpenFileDialog()
	result.SetOptions(options)
	return result
}

// SaveFile opens a save file dialog
func (dm *DialogManager) SaveFile() *SaveFileDialogStruct {
	return SaveFileDialog()
}

// SaveFileWithOptions opens a save file dialog with options
func (dm *DialogManager) SaveFileWithOptions(options *SaveFileDialogOptions) *SaveFileDialogStruct {
	result := SaveFileDialog()
	result.SetOptions(options)
	return result
}

// Info shows an info dialog
func (dm *DialogManager) Info() *MessageDialog {
	return InfoDialog()
}

// Question shows a question dialog
func (dm *DialogManager) Question() *MessageDialog {
	return QuestionDialog()
}

// Warning shows a warning dialog
func (dm *DialogManager) Warning() *MessageDialog {
	return WarningDialog()
}

// Error shows an error dialog
func (dm *DialogManager) Error() *MessageDialog {
	return ErrorDialog()
}
