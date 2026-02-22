// Cross-platform.

// Common File Dialogs
package cfd

type Dialog interface {
	// Show the dialog to the user.
	// Blocks until the user has closed the dialog.
	Show() error
	// Sets the dialog's parent window. Use 0 to set the dialog to have no parent window.
	SetParentWindowHandle(hwnd uintptr)
	// Show the dialog to the user.
	// Blocks until the user has closed the dialog and returns their selection.
	// Returns an error if the user cancelled the dialog.
	// Do not use for the Open Multiple Files dialog. Use ShowAndGetResults instead.
	ShowAndGetResult() (string, error)
	// Sets the title of the dialog window.
	SetTitle(title string) error
	// Sets the "role" of the dialog. This is used to derive the dialog's GUID, which the
	// OS will use to differentiate it from dialogs that are intended for other purposes.
	// This means that, for example, a dialog with role "Import" will have a different
	// previous location that it will open to than a dialog with role "Open". Can be any string.
	SetRole(role string) error
	// Sets the folder used as a default if there is not a recently used folder value available
	SetDefaultFolder(defaultFolder string) error
	// Sets the folder that the dialog always opens to.
	// If this is set, it will override the "default folder" behaviour and the dialog will always open to this folder.
	SetFolder(folder string) error
	// Gets the selected file or folder path, as an absolute path eg. "C:\Folder\file.txt"
	// Do not use for the Open Multiple Files dialog. Use GetResults instead.
	GetResult() (string, error)
	// Sets the file name, I.E. the contents of the file name text box.
	// For Select Folder Dialog, sets folder name.
	SetFileName(fileName string) error
	// Release the resources allocated to this Dialog.
	// Should be called when the dialog is finished with.
	Release() error
}

type FileDialog interface {
	Dialog
	// Set the list of file filters that the user can select.
	SetFileFilters(fileFilter []FileFilter) error
	// Set the selected item from the list of file filters (set using SetFileFilters) by its index. Defaults to 0 (the first item in the list) if not called.
	SetSelectedFileFilterIndex(index uint) error
	// Sets the default extension applied when a user does not provide one as part of the file name.
	// If the user selects a different file filter, the default extension will be automatically updated to match the new file filter.
	// For Open / Open Multiple File Dialog, this only has an effect when the user specifies a file name with no extension and a file with the default extension exists.
	// For Save File Dialog, this extension will be used whenever a user does not specify an extension.
	SetDefaultExtension(defaultExtension string) error
}

type OpenFileDialog interface {
	FileDialog
}

type OpenMultipleFilesDialog interface {
	FileDialog
	// Show the dialog to the user.
	// Blocks until the user has closed the dialog and returns the selected files.
	ShowAndGetResults() ([]string, error)
	// Gets the selected file paths, as absolute paths eg. "C:\Folder\file.txt"
	GetResults() ([]string, error)
}

type SelectFolderDialog interface {
	Dialog
}

type SaveFileDialog interface { // TODO Properties
	FileDialog
}
