//go:build windows
// +build windows

package windows

import (
	"github.com/leaanthony/go-common-file-dialog/cfd"
	"github.com/wailsapp/wails/v2/internal/frontend"
	"golang.org/x/sys/windows"
	"syscall"
)

// OpenDirectoryDialog prompts the user to select a directory
func (f *Frontend) OpenDirectoryDialog(options frontend.OpenDialogOptions) (string, error) {
	config := cfd.DialogConfig{
		Title:  options.Title,
		Role:   "PickFolder",
		Folder: options.DefaultDirectory,
	}
	thisDialog, err := cfd.NewSelectFolderDialog(config)
	if err != nil {
		return "", err
	}
	thisDialog.SetParentWindowHandle(f.mainWindow.Handle())
	defer func(thisDialog cfd.SelectFolderDialog) {
		err := thisDialog.Release()
		if err != nil {
			println("ERROR: Unable to release dialog:", err.Error())
		}
	}(thisDialog)
	result, err := thisDialog.ShowAndGetResult()
	if err != nil && err != cfd.ErrorCancelled {
		return "", err
	}
	return result, nil
}

// OpenFileDialog prompts the user to select a file
func (f *Frontend) OpenFileDialog(options frontend.OpenDialogOptions) (string, error) {
	config := cfd.DialogConfig{
		Folder:      options.DefaultDirectory,
		FileFilters: convertFilters(options.Filters),
		FileName:    options.DefaultFilename,
		Title:       options.Title,
	}
	thisdialog, err := cfd.NewOpenFileDialog(config)
	if err != nil {
		return "", err
	}
	thisdialog.SetParentWindowHandle(f.mainWindow.Handle())
	defer func(thisdialog cfd.OpenFileDialog) {
		err := thisdialog.Release()
		if err != nil {
			println("ERROR: Unable to release dialog:", err.Error())
		}
	}(thisdialog)
	result, err := thisdialog.ShowAndGetResult()
	if err != nil && err != cfd.ErrorCancelled {
		return "", err
	}
	return result, nil
}

// OpenMultipleFilesDialog prompts the user to select a file
func (f *Frontend) OpenMultipleFilesDialog(dialogOptions frontend.OpenDialogOptions) ([]string, error) {
	config := cfd.DialogConfig{
		Title:       dialogOptions.Title,
		Role:        "OpenMultipleFiles",
		FileFilters: convertFilters(dialogOptions.Filters),
		FileName:    dialogOptions.DefaultFilename,
		Folder:      dialogOptions.DefaultDirectory,
	}
	thisdialog, err := cfd.NewOpenMultipleFilesDialog(config)
	if err != nil {
		return nil, err
	}
	thisdialog.SetParentWindowHandle(f.mainWindow.Handle())
	defer func(thisdialog cfd.OpenMultipleFilesDialog) {
		err := thisdialog.Release()
		if err != nil {
			println("ERROR: Unable to release dialog:", err.Error())
		}
	}(thisdialog)
	result, err := thisdialog.ShowAndGetResults()
	if err != nil && err != cfd.ErrorCancelled {
		return nil, err
	}
	return result, nil
}

// SaveFileDialog prompts the user to select a file
func (f *Frontend) SaveFileDialog(dialogOptions frontend.SaveDialogOptions) (string, error) {
	saveDialog, err := cfd.NewSaveFileDialog(cfd.DialogConfig{
		Title:       dialogOptions.Title,
		Role:        "SaveFile",
		FileFilters: convertFilters(dialogOptions.Filters),
		FileName:    dialogOptions.DefaultFilename,
		Folder:      dialogOptions.DefaultDirectory,
	})
	if err != nil {
		return "", err
	}
	saveDialog.SetParentWindowHandle(f.mainWindow.Handle())
	err = saveDialog.Show()
	if err != nil {
		return "", err
	}
	result, err := saveDialog.GetResult()
	if err != nil && err != cfd.ErrorCancelled {
		return "", err
	}
	return result, nil
}

// MessageDialog show a message dialog to the user
func (f *Frontend) MessageDialog(options frontend.MessageDialogOptions) (string, error) {

	title, err := syscall.UTF16PtrFromString(options.Title)
	if err != nil {
		return "", err
	}
	message, err := syscall.UTF16PtrFromString(options.Message)
	if err != nil {
		return "", err
	}
	var flags uint32
	switch options.Type {
	case frontend.InfoDialog:
		flags = windows.MB_OK | windows.MB_ICONINFORMATION
	case frontend.ErrorDialog:
		flags = windows.MB_ICONERROR | windows.MB_OK
	case frontend.QuestionDialog:
		flags = windows.MB_YESNO
	case frontend.WarningDialog:
		flags = windows.MB_OK | windows.MB_ICONWARNING
	}

	button, _ := windows.MessageBox(windows.HWND(f.mainWindow.Handle()), message, title, flags|windows.MB_SYSTEMMODAL)
	// This maps MessageBox return values to strings
	responses := []string{"", "Ok", "Cancel", "Abort", "Retry", "Ignore", "Yes", "No", "", "", "Try Again", "Continue"}
	result := "Error"
	if int(button) < len(responses) {
		result = responses[button]
	}
	return result, nil
}

func convertFilters(filters []frontend.FileFilter) []cfd.FileFilter {
	var result []cfd.FileFilter
	for _, filter := range filters {
		result = append(result, cfd.FileFilter(filter))
	}
	return result
}
