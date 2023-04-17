//go:build windows
// +build windows

package windows

import (
	"path/filepath"
	"strings"
	"syscall"

	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
	"github.com/wailsapp/wails/v2/internal/go-common-file-dialog/cfd"
	"golang.org/x/sys/windows"
)

func (f *Frontend) getHandleForDialog() w32.HWND {
	if f.mainWindow.IsVisible() {
		return f.mainWindow.Handle()
	}
	return 0
}

func getDefaultFolder(folder string) (string, error) {
	if folder == "" {
		return "", nil
	}
	return filepath.Abs(folder)
}

// OpenDirectoryDialog prompts the user to select a directory
func (f *Frontend) OpenDirectoryDialog(options frontend.OpenDialogOptions) (string, error) {

	defaultFolder, err := getDefaultFolder(options.DefaultDirectory)
	if err != nil {
		return "", err
	}

	config := cfd.DialogConfig{
		Title:  options.Title,
		Role:   "PickFolder",
		Folder: defaultFolder,
	}

	result, err := f.showCfdDialog(
		func() (cfd.Dialog, error) {
			return cfd.NewSelectFolderDialog(config)
		}, false)

	if err != nil && err != cfd.ErrorCancelled {
		return "", err
	}
	return result.(string), nil
}

// OpenFileDialog prompts the user to select a file
func (f *Frontend) OpenFileDialog(options frontend.OpenDialogOptions) (string, error) {
	defaultFolder, err := getDefaultFolder(options.DefaultDirectory)
	if err != nil {
		return "", err
	}

	config := cfd.DialogConfig{
		Folder:      defaultFolder,
		FileFilters: convertFilters(options.Filters),
		FileName:    options.DefaultFilename,
		Title:       options.Title,
	}

	result, err := f.showCfdDialog(
		func() (cfd.Dialog, error) {
			return cfd.NewOpenFileDialog(config)
		}, false)

	if err != nil && err != cfd.ErrorCancelled {
		return "", err
	}
	return result.(string), nil
}

// OpenMultipleFilesDialog prompts the user to select a file
func (f *Frontend) OpenMultipleFilesDialog(options frontend.OpenDialogOptions) ([]string, error) {

	defaultFolder, err := getDefaultFolder(options.DefaultDirectory)
	if err != nil {
		return nil, err
	}

	config := cfd.DialogConfig{
		Title:       options.Title,
		Role:        "OpenMultipleFiles",
		FileFilters: convertFilters(options.Filters),
		FileName:    options.DefaultFilename,
		Folder:      defaultFolder,
	}

	result, err := f.showCfdDialog(
		func() (cfd.Dialog, error) {
			return cfd.NewOpenMultipleFilesDialog(config)
		}, true)

	if err != nil && err != cfd.ErrorCancelled {
		return nil, err
	}
	return result.([]string), nil
}

// SaveFileDialog prompts the user to select a file
func (f *Frontend) SaveFileDialog(options frontend.SaveDialogOptions) (string, error) {

	defaultFolder, err := getDefaultFolder(options.DefaultDirectory)
	if err != nil {
		return "", err
	}

	config := cfd.DialogConfig{
		Title:       options.Title,
		Role:        "SaveFile",
		FileFilters: convertFilters(options.Filters),
		FileName:    options.DefaultFilename,
		Folder:      defaultFolder,
	}

	result, err := f.showCfdDialog(
		func() (cfd.Dialog, error) {
			return cfd.NewSaveFileDialog(config)
		}, false)

	if err != nil && err != cfd.ErrorCancelled {
		return "", err
	}
	return result.(string), nil
}

func (f *Frontend) showCfdDialog(newDlg func() (cfd.Dialog, error), isMultiSelect bool) (any, error) {
	return invokeSync(f.mainWindow, func() (any, error) {
		dlg, err := newDlg()
		if err != nil {
			return nil, err
		}
		defer func() {
			err := dlg.Release()
			if err != nil {
				println("ERROR: Unable to release dialog:", err.Error())
			}
		}()

		dlg.SetParentWindowHandle(f.getHandleForDialog())
		if multi, _ := dlg.(cfd.OpenMultipleFilesDialog); multi != nil && isMultiSelect {
			return multi.ShowAndGetResults()
		}
		return dlg.ShowAndGetResult()
	})
}

func calculateMessageDialogFlags(options frontend.MessageDialogOptions) uint32 {
	var flags uint32

	switch options.Type {
	case frontend.InfoDialog:
		flags = windows.MB_OK | windows.MB_ICONINFORMATION
	case frontend.ErrorDialog:
		flags = windows.MB_ICONERROR | windows.MB_OK
	case frontend.QuestionDialog:
		flags = windows.MB_YESNO
		if strings.TrimSpace(strings.ToLower(options.DefaultButton)) == "no" {
			flags |= windows.MB_DEFBUTTON2
		}
	case frontend.WarningDialog:
		flags = windows.MB_OK | windows.MB_ICONWARNING
	}

	return flags
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

	flags := calculateMessageDialogFlags(options)

	button, _ := windows.MessageBox(windows.HWND(f.getHandleForDialog()), message, title, flags|windows.MB_SYSTEMMODAL)
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
