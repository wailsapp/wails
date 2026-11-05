//go:build windows && !server

package application

import (
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"unsafe"

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

	// MessageBox cannot show custom button captions, so question dialogs with
	// custom buttons use a TaskDialog (#5794). It needs comctl32 v6; if that is
	// unavailable (no manifest) or the call fails, fall back to MessageBox.
	if m.dialog.DialogType == QuestionDialogType && len(m.dialog.Buttons) > 0 {
		if m.showTaskDialog(parentWindow) {
			return
		}
	}

	if m.UseAppIcon || m.dialog.Icon != nil {
		// Use the application's embedded icon resource (ID 3). MessageBoxIndirect
		// cannot render arbitrary icon bytes, so a custom Icon also maps to the
		// app icon — the closest supported behaviour (full custom-icon support
		// needs a TaskDialog implementation).
		//
		button, err = w32.MessageBoxWithIcon(parentWindow, message, title, 3, messageDialogUserIconFlags(flags))
		if err != nil {
			// Dev binaries (`go run`, `wails3 dev`) have no embedded icon
			// resource, which makes MessageBoxIndirect fail outright — the
			// dialog never appeared and the app aborted via fatal error
			// (#4233). Fall back to a standard dialog instead.
			button, err = windows.MessageBox(windows.HWND(parentWindow), message, title, flags|windows.MB_SYSTEMMODAL)
		}
		if err != nil {
			globalApplication.handleFatalError(err)
		}
	} else {
		button, err = windows.MessageBox(windows.HWND(parentWindow), message, title, flags|windows.MB_SYSTEMMODAL)
		if err != nil {
			globalApplication.handleFatalError(err)
		}
	}
	// Check if there's a callback for the button pressed
	if b := messageBoxButtonFor(m.dialog.Buttons, button); b != nil && b.Callback != nil {
		b.Callback()
	}
}

// taskDialogCreated preserves MessageBox's MB_SYSTEMMODAL topmost behavior.
// TaskDialogIndirect itself provides modality relative to hwndParent.
var taskDialogCreated = syscall.NewCallback(func(hwnd, notification, wparam, lparam, data uintptr) uintptr {
	const tdnCreated = 0
	if notification == tdnCreated {
		w32.SetWindowPos(hwnd, w32.HWND_TOPMOST, 0, 0, 0, 0, w32.SWP_NOMOVE|w32.SWP_NOSIZE|w32.SWP_NOACTIVATE)
	}
	return 0 // S_OK
})

// taskDialogIcon uses the same embedded app icon as MessageBoxWithIcon.
// Missing development resources omit the icon while preserving custom captions.
func (m *windowsDialog) taskDialogIcon() (instance, resource uint64) {
	if m.UseAppIcon || m.dialog.Icon != nil {
		handle := w32.GetApplicationHandle()
		if w32.LoadIconWithResourceID(handle, 3) != 0 {
			return uint64(handle), 3
		}
	}
	return 0, 0
}

// showTaskDialog shows the dialog with its own button captions. It returns
// false, without having shown anything, when TaskDialogIndirect cannot be used.
func (m *windowsDialog) showTaskDialog(parentWindow uintptr) bool {
	if !w32.TaskDialogAvailable() {
		return false
	}
	plan := planTaskDialog(m.dialog.Buttons)
	if len(plan.buttons) == 0 {
		return false
	}

	ptrSize := int(unsafe.Sizeof(uintptr(0)))
	title := w32.MustStringToUTF16Ptr(m.dialog.Title)
	message := w32.MustStringToUTF16Ptr(m.dialog.Message)
	labels := make([]*uint16, len(plan.buttons))
	labelAddrs := make([]uint64, len(plan.buttons))
	ids := make([]int32, len(plan.buttons))
	for i, b := range plan.buttons {
		labels[i] = w32.MustStringToUTF16Ptr(b.Label)
		labelAddrs[i] = uint64(uintptr(unsafe.Pointer(labels[i])))
		ids[i] = b.ID
	}
	buttons := encodeTaskDialogButtons(ptrSize, ids, labelAddrs)

	var flags uint32
	if plan.cancel != nil {
		flags |= tdfAllowDialogCancellation
	}
	if parentWindow != 0 {
		flags |= tdfPositionRelativeToWindow
	}
	hInstance, mainIcon := m.taskDialogIcon()
	config := encodeTaskDialogConfig(ptrSize, taskDialogConfig{
		hwndParent:    uint64(parentWindow),
		hInstance:     hInstance,
		mainIcon:      mainIcon,
		callback:      uint64(taskDialogCreated),
		flags:         flags,
		title:         uint64(uintptr(unsafe.Pointer(title))),
		content:       uint64(uintptr(unsafe.Pointer(message))),
		buttonCount:   uint32(len(plan.buttons)),
		buttons:       uint64(uintptr(unsafe.Pointer(&buttons[0]))),
		defaultButton: plan.defaultID,
	})

	result, err := w32.TaskDialogIndirect(config)
	// The config only holds raw addresses; keep what they point to alive.
	runtime.KeepAlive(title)
	runtime.KeepAlive(message)
	runtime.KeepAlive(labels)
	runtime.KeepAlive(buttons)
	if err != nil {
		globalApplication.warning("TaskDialogIndirect failed, falling back to MessageBox: %v", err)
		return false
	}
	if b := plan.buttonFor(result); b != nil && b.Callback != nil {
		b.Callback()
	}
	return true
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
	switch openDialogKindFor(m.dialog.canChooseDirectories, m.dialog.allowsMultipleSelection) {
	case openMultipleFiles:
		temp, err := showCfdDialog(
			func() (cfd.Dialog, error) {
				return cfd.NewOpenMultipleFilesDialog(config)
			}, true, m.dialog.window)
		if err != nil {
			return nil, err
		}
		result = temp.([]string)
	case openMultipleFolders:
		temp, err := showCfdDialog(
			func() (cfd.Dialog, error) {
				return cfd.NewSelectMultipleFoldersDialog(config)
			}, true, m.dialog.window)
		if err != nil {
			return nil, err
		}
		result = temp.([]string)
	case openSingleFolder:
		temp, err := showCfdDialog(
			func() (cfd.Dialog, error) {
				return cfd.NewSelectFolderDialog(config)
			}, false, m.dialog.window)
		if err != nil {
			return nil, err
		}
		result = []string{temp.(string)}
	default:
		temp, err := showCfdDialog(
			func() (cfd.Dialog, error) {
				return cfd.NewOpenFileDialog(config)
			}, false, m.dialog.window)
		if err != nil {
			return nil, err
		}
		result = []string{temp.(string)}
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

// openDialogKind is the Common Item Dialog an OpenFile dialog is shown as.
type openDialogKind int

const (
	openSingleFile openDialogKind = iota
	openMultipleFiles
	openSingleFolder
	openMultipleFolders
)

// openDialogKindFor picks the dialog for an OpenFile request. The Windows
// dialog cannot mix files and folders, so directories take precedence when
// both are allowed (unchanged); multiple selection applies to folders as well
// as files, matching macOS and Linux.
func openDialogKindFor(canChooseDirectories, allowsMultipleSelection bool) openDialogKind {
	switch {
	case canChooseDirectories && allowsMultipleSelection:
		return openMultipleFolders
	case canChooseDirectories:
		return openSingleFolder
	case allowsMultipleSelection:
		return openMultipleFiles
	default:
		return openSingleFile
	}
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

// messageDialogUserIconFlags converts standard message dialog flags for use
// with MB_USERICON. The user icon replaces the standard one, so the MB_ICON*
// bits are stripped — but the button configuration is preserved: forcing
// MB_OK here (the old behaviour) silently destroyed Yes/No buttons on
// question dialogs shown with an icon (#4233).
func messageDialogUserIconFlags(flags uint32) uint32 {
	const mbIconMask = w32.MB_ICONHAND | w32.MB_ICONQUESTION | w32.MB_ICONASTERISK | w32.MB_USERICON
	return (flags &^ uint32(mbIconMask)) | windows.MB_USERICON | windows.MB_SYSTEMMODAL
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
