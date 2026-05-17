//go:build windows

package application

import (
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v3/internal/go-common-file-dialog/cfd"
	"github.com/wailsapp/wails/v3/pkg/w32"
	"golang.org/x/sys/windows"
)

func (m *windowsApp) showAboutDialog(title string, message string, _ []byte) {
	about := newDialogImpl(&MessageDialog{
		MessageDialogOptions: MessageDialogOptions{
			DialogType: InfoDialogType,
			Title:      title,
			Message:    message,
		},
	})
	about.UseAppIcon = true
	about.show()
}

type windowsDialog struct {
	dialog *MessageDialog

	//dialogImpl unsafe.Pointer
	UseAppIcon bool
}

func (m *windowsDialog) show() {

	title := w32.MustStringToUTF16Ptr(m.dialog.Title)
	message := w32.MustStringToUTF16Ptr(m.dialog.Message)
	flags := calculateMessageDialogFlags(m.dialog.MessageDialogOptions)
	var button int32
	var err error

	var parentWindow uintptr
	if m.dialog.window != nil {
		nativeWindow := m.dialog.window.NativeWindow()
		if nativeWindow != nil {
			parentWindow = uintptr(nativeWindow)
		}
	}

	if m.UseAppIcon || m.dialog.Icon != nil {
		// 3 is the application icon
		button, err = w32.MessageBoxWithIcon(parentWindow, message, title, 3, windows.MB_OK|windows.MB_USERICON)
		if err != nil {
			globalApplication.handleFatalError(err)
		}
	} else {
		button, err = windows.MessageBox(windows.HWND(parentWindow), message, title, flags|windows.MB_SYSTEMMODAL)
		if err != nil {
			globalApplication.handleFatalError(err)
		}
	}
	// This maps MessageBox return values to strings
	responses := []string{"", "Ok", "Cancel", "Abort", "Retry", "Ignore", "Yes", "No", "", "", "Try Again", "Continue"}
	result := "Error"
	if int(button) < len(responses) {
		result = responses[button]
	}
	// Check if there's a callback for the button pressed
	for _, buttonInDialog := range m.dialog.Buttons {
		if buttonInDialog.Label == result {
			if buttonInDialog.Callback != nil {
				buttonInDialog.Callback()
			}
		}
	}
}

func newDialogImpl(d *MessageDialog) *windowsDialog {
	return &windowsDialog{
		dialog: d,
	}
}

type windowOpenFileDialog struct {
	dialog *OpenFileDialogStruct
}

func newOpenFileDialogImpl(d *OpenFileDialogStruct) *windowOpenFileDialog {
	return &windowOpenFileDialog{
		dialog: d,
	}
}

func getDefaultFolder(folder string) (string, error) {
	if folder == "" {
		return "", nil
	}
	return filepath.Abs(folder)
}

func (m *windowOpenFileDialog) show() (chan string, error) {

	defaultFolder, err := getDefaultFolder(m.dialog.directory)
	if err != nil {
		return nil, err
	}

	config := cfd.DialogConfig{
		Title:       m.dialog.title,
		Role:        "PickFolder",
		FileFilters: convertFilters(m.dialog.filters),
		Folder:      defaultFolder,
	}

	var result []string
	if m.dialog.allowsMultipleSelection && !m.dialog.canChooseDirectories {
		temp, err := showCfdDialog(
			func() (cfd.Dialog, error) {
				return cfd.NewOpenMultipleFilesDialog(config)
			}, true, m.dialog.window)
		if err != nil {
			return nil, err
		}
		result = temp.([]string)
	} else {
		if m.dialog.canChooseDirectories {
			temp, err := showCfdDialog(
				func() (cfd.Dialog, error) {
					return cfd.NewSelectFolderDialog(config)
				}, false, m.dialog.window)
			if err != nil {
				return nil, err
			}
			result = []string{temp.(string)}
		} else {
			temp, err := showCfdDialog(
				func() (cfd.Dialog, error) {
					return cfd.NewOpenFileDialog(config)
				}, false, m.dialog.window)
			if err != nil {
				return nil, err
			}
			result = []string{temp.(string)}
		}
	}

	files := make(chan string)
	go func() {
		defer handlePanic()
		for _, file := range result {
			files <- file
		}
		close(files)
	}()
	return files, nil
}

type windowSaveFileDialog struct {
	dialog *SaveFileDialogStruct
}

func newSaveFileDialogImpl(d *SaveFileDialogStruct) *windowSaveFileDialog {
	return &windowSaveFileDialog{
		dialog: d,
	}
}

func (m *windowSaveFileDialog) show() (chan string, error) {
	files := make(chan string)
	defaultFolder, err := getDefaultFolder(m.dialog.directory)
	if err != nil {
		close(files)
		return files, err
	}

	config := cfd.DialogConfig{
		Title:       m.dialog.title,
		Role:        "SaveFile",
		FileFilters: convertFilters(m.dialog.filters),
		FileName:    m.dialog.filename,
		Folder:      defaultFolder,
	}

	// Original PR for v2 by @almas1992: https://github.com/wailsapp/wails/pull/3205
	if len(m.dialog.filters) > 0 {
		config.DefaultExtension = strings.TrimPrefix(strings.Split(m.dialog.filters[0].Pattern, ";")[0], "*")
	}

	result, err := showCfdDialog(
		func() (cfd.Dialog, error) {
			return cfd.NewSaveFileDialog(config)
		}, false, m.dialog.window)
	if err != nil {
		close(files)
		return files, err
	}
	go func() {
		defer handlePanic()
		f, ok := result.(string)
		if ok {
			files <- f
		}
		close(files)
	}()
	return files, err
}

func calculateMessageDialogFlags(options MessageDialogOptions) uint32 {
	var flags uint32

	switch options.DialogType {
	case InfoDialogType:
		flags = windows.MB_OK | windows.MB_ICONINFORMATION
	case ErrorDialogType:
		flags = windows.MB_ICONERROR | windows.MB_OK
	case QuestionDialogType:
		flags = windows.MB_YESNO
		for _, button := range options.Buttons {
			if strings.TrimSpace(strings.ToLower(button.Label)) == "no" && button.IsDefault {
				flags |= windows.MB_DEFBUTTON2
			}
		}
	case WarningDialogType:
		flags = windows.MB_OK | windows.MB_ICONWARNING
	}

	return flags
}

func convertFilters(filters []FileFilter) []cfd.FileFilter {
	var result []cfd.FileFilter
	for _, filter := range filters {
		result = append(result, cfd.FileFilter(filter))
	}
	return result
}

func showCfdDialog(newDlg func() (cfd.Dialog, error), isMultiSelect bool, parentWindow Window) (any, error) {
	dlg, err := newDlg()
	if err != nil {
		return nil, err
	}

	// Set parent window if provided
	if parentWindow != nil {
		nativeWindow := parentWindow.NativeWindow()
		if nativeWindow != nil {
			dlg.SetParentWindowHandle(uintptr(nativeWindow))
		}
	}

	defer func() {
		err := dlg.Release()
		if err != nil {
			globalApplication.error("unable to release dialog: %w", err)
		}
	}()

	if multi, _ := dlg.(cfd.OpenMultipleFilesDialog); multi != nil && isMultiSelect {
		paths, err := multi.ShowAndGetResults()
		if err != nil {
			return nil, err
		}

		for i, path := range paths {
			paths[i] = filepath.Clean(path)
		}
		return paths, nil
	}

	path, err := dlg.ShowAndGetResult()
	if err != nil {
		return nil, err
	}
	return filepath.Clean(path), nil
}
