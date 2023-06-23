//go:build windows

package application

import (
	"github.com/wailsapp/wails/v3/internal/go-common-file-dialog/cfd"
	"github.com/wailsapp/wails/v3/pkg/w32"
	"golang.org/x/sys/windows"
	"path/filepath"
	"strings"
)

func (m *windowsApp) showAboutDialog(title string, message string, icon []byte) {
	about := newDialogImpl(&MessageDialog{
		MessageDialogOptions: MessageDialogOptions{
			DialogType: InfoDialog,
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

	var parentWindow uintptr
	var err error
	if m.dialog.window != nil {
		parentWindow, err = m.dialog.window.NativeWindowHandle()
		if err != nil {
			w32.Fatal(err.Error())
		}
	}

	if m.UseAppIcon {
		// 3 is the application icon
		button, _ = w32.MessageBoxWithIcon(parentWindow, message, title, 3, windows.MB_OK|windows.MB_USERICON)
	} else {
		button, _ = windows.MessageBox(windows.HWND(parentWindow), message, title, flags|windows.MB_SYSTEMMODAL)
	}
	// This maps MessageBox return values to strings
	responses := []string{"", "Ok", "Cancel", "Abort", "Retry", "Ignore", "Yes", "No", "", "", "Try Again", "Continue"}
	result := "Error"
	if int(button) < len(responses) {
		result = responses[button]
	}
	// Check if there's a callback for the button pressed
	for _, button := range m.dialog.Buttons {
		if button.Label == result {
			if button.Callback != nil {
				button.Callback()
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
	dialog *OpenFileDialog
}

func newOpenFileDialogImpl(d *OpenFileDialog) *windowOpenFileDialog {
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

func (m *windowOpenFileDialog) show() ([]string, error) {

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

	if m.dialog.window != nil {
		config.ParentWindowHandle, err = m.dialog.window.NativeWindowHandle()
		if err != nil {
			w32.Fatal(err.Error())
		}
	}

	var result []string
	if m.dialog.allowsMultipleSelection {
		temp, err := showCfdDialog(
			func() (cfd.Dialog, error) {
				return cfd.NewOpenMultipleFilesDialog(config)
			}, true)
		if err != nil {
			return nil, err
		}
		result = temp.([]string)
	} else {
		temp, err := showCfdDialog(
			func() (cfd.Dialog, error) {
				return cfd.NewOpenFileDialog(config)
			}, false)
		if err != nil {
			return nil, err
		}
		result = []string{temp.(string)}
	}

	return result, nil
}

type windowSaveFileDialog struct {
	dialog *SaveFileDialog
}

func newSaveFileDialogImpl(d *SaveFileDialog) *windowSaveFileDialog {
	return &windowSaveFileDialog{
		dialog: d,
	}
}

func (m *windowSaveFileDialog) show() (string, error) {
	defaultFolder, err := getDefaultFolder(m.dialog.directory)
	if err != nil {
		return "", err
	}

	config := cfd.DialogConfig{
		Title:       m.dialog.title,
		Role:        "SaveFile",
		FileFilters: convertFilters(m.dialog.filters),
		FileName:    m.dialog.filename,
		Folder:      defaultFolder,
	}

	result, err := showCfdDialog(
		func() (cfd.Dialog, error) {
			return cfd.NewSaveFileDialog(config)
		}, false)
	return result.(string), nil
}

func calculateMessageDialogFlags(options MessageDialogOptions) uint32 {
	var flags uint32

	switch options.DialogType {
	case InfoDialog:
		flags = windows.MB_OK | windows.MB_ICONINFORMATION
	case ErrorDialog:
		flags = windows.MB_ICONERROR | windows.MB_OK
	case QuestionDialog:
		flags = windows.MB_YESNO
		for _, button := range options.Buttons {
			if strings.TrimSpace(strings.ToLower(button.Label)) == "no" && button.IsDefault {
				flags |= windows.MB_DEFBUTTON2
			}
		}
	case WarningDialog:
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

func showCfdDialog(newDlg func() (cfd.Dialog, error), isMultiSelect bool) (any, error) {
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

	if multi, _ := dlg.(cfd.OpenMultipleFilesDialog); multi != nil && isMultiSelect {
		return multi.ShowAndGetResults()
	}
	return dlg.ShowAndGetResult()
}
