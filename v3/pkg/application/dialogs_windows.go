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
	// If we have custom buttons, use TaskDialog
	if len(m.dialog.Buttons) > 0 && hasCustomButtons(m.dialog.Buttons) {
		m.showTaskDialog()
		return
	}

	// Fallback to MessageBox for standard dialogs
	m.showMessageBox()
}

func (m *windowsDialog) showTaskDialog() {
	title := w32.MustStringToUTF16Ptr(m.dialog.Title)
	message := w32.MustStringToUTF16Ptr(m.dialog.Message)

	var parentWindow uintptr
	var err error
	if m.dialog.window != nil {
		parentWindow, err = m.dialog.window.NativeWindowHandle()
		if err != nil {
			globalApplication.handleFatalError(err)
		}
	}

	// Create custom buttons
	buttons := make([]w32.TASKDIALOG_BUTTON, len(m.dialog.Buttons))
	for i, btn := range m.dialog.Buttons {
		buttons[i] = w32.TASKDIALOG_BUTTON{
			NButtonID:     int32(i + 1000), // Use unique IDs starting from 1000
			PszButtonText: w32.MustStringToUTF16Ptr(btn.Label),
		}
	}

	// Determine default button
	var defaultButton int32 = 1000 // Default to first button
	for i, btn := range m.dialog.Buttons {
		if btn.IsDefault {
			defaultButton = int32(i + 1000)
			break
		}
	}

	// Get appropriate icon
	icon := getTaskDialogIcon(m.dialog.DialogType, m.dialog.Icon)

	// Show TaskDialog
	buttonPressed, err := w32.CustomTaskDialog(
		w32.HWND(parentWindow),
		title,
		message,
		nil, // content (we put everything in instruction for simplicity)
		buttons,
		icon,
		defaultButton,
	)

	if err != nil {
		globalApplication.handleFatalError(err)
		return
	}

	// Find and execute callback for pressed button
	buttonIndex := int(buttonPressed - 1000)
	if buttonIndex >= 0 && buttonIndex < len(m.dialog.Buttons) {
		if m.dialog.Buttons[buttonIndex].Callback != nil {
			m.dialog.Buttons[buttonIndex].Callback()
		}
	}
}

func (m *windowsDialog) showMessageBox() {
	title := w32.MustStringToUTF16Ptr(m.dialog.Title)
	message := w32.MustStringToUTF16Ptr(m.dialog.Message)
	flags := calculateMessageDialogFlags(m.dialog.MessageDialogOptions)
	var button int32

	var parentWindow uintptr
	var err error
	if m.dialog.window != nil {
		parentWindow, err = m.dialog.window.NativeWindowHandle()
		if err != nil {
			globalApplication.handleFatalError(err)
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

// hasCustomButtons checks if the dialog has custom (non-standard) buttons
func hasCustomButtons(buttons []*Button) bool {
	standardButtons := map[string]bool{
		"Ok":       true,
		"Cancel":   true,
		"Yes":      true,
		"No":       true,
		"Abort":    true,
		"Retry":    true,
		"Ignore":   true,
		"Continue": true,
	}
	
	for _, btn := range buttons {
		if !standardButtons[btn.Label] {
			return true
		}
	}
	return false
}

// getTaskDialogIcon returns the appropriate icon for TaskDialog
func getTaskDialogIcon(dialogType DialogType, customIcon []byte) uintptr {
	if customIcon != nil {
		// For custom icons, we would need to load from memory
		// This is complex and requires additional implementation
		// For now, fall back to system icons
	}
	
	switch dialogType {
	case InfoDialogType:
		return w32.TD_INFORMATION_ICON
	case WarningDialogType:
		return w32.TD_WARNING_ICON
	case ErrorDialogType:
		return w32.TD_ERROR_ICON
	case QuestionDialogType:
		return w32.TD_INFORMATION_ICON // Question uses info icon typically
	default:
		return w32.TD_INFORMATION_ICON
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

	if m.dialog.window != nil {
		config.ParentWindowHandle, err = m.dialog.window.NativeWindowHandle()
		if err != nil {
			globalApplication.handleFatalError(err)
		}
	}

	var result []string
	if m.dialog.allowsMultipleSelection && !m.dialog.canChooseDirectories {
		temp, err := showCfdDialog(
			func() (cfd.Dialog, error) {
				return cfd.NewOpenMultipleFilesDialog(config)
			}, true)
		if err != nil {
			return nil, err
		}
		result = temp.([]string)
	} else {
		if m.dialog.canChooseDirectories {
			temp, err := showCfdDialog(
				func() (cfd.Dialog, error) {
					return cfd.NewSelectFolderDialog(config)
				}, false)
			if err != nil {
				return nil, err
			}
			result = []string{temp.(string)}
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
		}, false)
	go func() {
		defer handlePanic()
		files <- result.(string)
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

func showCfdDialog(newDlg func() (cfd.Dialog, error), isMultiSelect bool) (any, error) {
	dlg, err := newDlg()
	if err != nil {
		return nil, err
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
