// Cross-platform.

package cfd

type FileFilter struct {
	// The display name of the filter (That is shown to the user)
	DisplayName string
	// The filter pattern. Eg. "*.txt;*.png" to select all txt and png files, "*.*" to select any files, etc.
	Pattern string
}

type DialogConfig struct {
	// The title of the dialog
	Title string
	// The role of the dialog. This is used to derive the dialog's GUID, which the
	// OS will use to differentiate it from dialogs that are intended for other purposes.
	// This means that, for example, a dialog with role "Import" will have a different
	// previous location that it will open to than a dialog with role "Open". Can be any string.
	Role string
	// The default folder - the folder that is used the first time the user opens it
	// (after the first time their last used location is used).
	DefaultFolder string
	// The initial folder - the folder that the dialog always opens to if not empty.
	// If this is not empty, it will override the "default folder" behaviour and
	// the dialog will always open to this folder.
	Folder string
	// The file filters that restrict which types of files the dialog is able to choose.
	// Ignored by Select Folder Dialog.
	FileFilters []FileFilter
	// Sets the initially selected file filter. This is an index of FileFilters.
	// Ignored by Select Folder Dialog.
	SelectedFileFilterIndex uint
	// The initial name of the file (I.E. the text in the file name text box) when the user opens the dialog.
	// For the Select Folder Dialog, this sets the initial folder name.
	FileName string
	// The default extension applied when a user does not provide one as part of the file name.
	// If the user selects a different file filter, the default extension will be automatically updated to match the new file filter.
	// For Open / Open Multiple File Dialog, this only has an effect when the user specifies a file name with no extension and a file with the default extension exists.
	// For Save File Dialog, this extension will be used whenever a user does not specify an extension.
	// Ignored by Select Folder Dialog.
	DefaultExtension string
	// ParentWindowHandle is the handle (HWND) to the parent window of the dialog.
	// If left as 0 / nil, the dialog will have no parent window.
	ParentWindowHandle uintptr
}

var defaultFilters = []FileFilter{
	{
		DisplayName: "All Files (*.*)",
		Pattern:     "*.*",
	},
}

func (config *DialogConfig) apply(dialog Dialog) (err error) {
	if config.Title != "" {
		err = dialog.SetTitle(config.Title)
		if err != nil {
			return
		}
	}

	if config.Role != "" {
		err = dialog.SetRole(config.Role)
		if err != nil {
			return
		}
	}

	if config.Folder != "" {
		err = dialog.SetFolder(config.Folder)
		if err != nil {
			return
		}
	}

	if config.DefaultFolder != "" {
		err = dialog.SetDefaultFolder(config.DefaultFolder)
		if err != nil {
			return
		}
	}

	if config.FileName != "" {
		err = dialog.SetFileName(config.FileName)
		if err != nil {
			return
		}
	}

	dialog.SetParentWindowHandle(config.ParentWindowHandle)

	if dialog, ok := dialog.(FileDialog); ok {
		var fileFilters []FileFilter
		if config.FileFilters != nil && len(config.FileFilters) > 0 {
			fileFilters = config.FileFilters
		} else {
			fileFilters = defaultFilters
		}
		err = dialog.SetFileFilters(fileFilters)
		if err != nil {
			return
		}

		if config.SelectedFileFilterIndex != 0 {
			err = dialog.SetSelectedFileFilterIndex(config.SelectedFileFilterIndex)
			if err != nil {
				return
			}
		}

		if config.DefaultExtension != "" {
			err = dialog.SetDefaultExtension(config.DefaultExtension)
			if err != nil {
				return
			}
		}
	}

	return
}
