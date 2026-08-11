//go:build android

package application

import (
	"encoding/json"
	"fmt"
	"sync"
)

// Message dialogs are backed by AlertDialog, file dialogs by the Storage
// Access Framework document picker (selected documents are copied into the
// app's cache directory so callers receive real filesystem paths). Save
// dialogs have no Android counterpart that yields a filesystem path (apps
// write into their sandbox and share via intents), so they return an
// explicit error instead of failing silently.

// Pending message dialogs keyed by callback ID
var (
	androidDialogsLock    sync.Mutex
	androidPendingDialogs = make(map[uint]*MessageDialog)
	androidNextDialogID   uint = 1
)

type androidDialogButton struct {
	Label     string `json:"label"`
	IsCancel  bool   `json:"isCancel"`
	IsDefault bool   `json:"isDefault"`
}

type androidDialogOptions struct {
	Title   string                `json:"title"`
	Message string                `json:"message"`
	Buttons []androidDialogButton `json:"buttons"`
}

type androidDialog struct {
	dialog *MessageDialog
}

func newDialogImpl(d *MessageDialog) *androidDialog {
	return &androidDialog{dialog: d}
}

func (d *androidDialog) show() {
	androidDialogsLock.Lock()
	id := androidNextDialogID
	androidNextDialogID++
	androidPendingDialogs[id] = d.dialog
	androidDialogsLock.Unlock()

	buttons := make([]androidDialogButton, 0, len(d.dialog.Buttons))
	for _, b := range d.dialog.Buttons {
		buttons = append(buttons, androidDialogButton{
			Label:     b.Label,
			IsCancel:  b.IsCancel,
			IsDefault: b.IsDefault,
		})
	}

	title := d.dialog.Title
	if title == "" {
		title = defaultTitles[d.dialog.DialogType]
	}

	optionsJSON, _ := json.Marshal(androidDialogOptions{
		Title:   title,
		Message: d.dialog.Message,
		Buttons: buttons,
	})

	androidBridgeVoidIntString("showMessageDialog", int(id), string(optionsJSON))
}

// androidDialogCallback is invoked from JNI when a dialog button is pressed
// (buttonIndex is the index into the dialog's button slice, or -1 for a
// dismissal with no matching button).
func androidDialogCallback(callbackID uint, buttonIndex int) {
	androidDialogsLock.Lock()
	dialog, ok := androidPendingDialogs[callbackID]
	delete(androidPendingDialogs, callbackID)
	androidDialogsLock.Unlock()
	if !ok || dialog == nil {
		return
	}
	if buttonIndex < 0 || buttonIndex >= len(dialog.Buttons) {
		return
	}
	button := dialog.Buttons[buttonIndex]
	if button.Callback != nil {
		// Run the callback off the JNI thread, mirroring desktop behaviour
		go func() {
			defer handlePanic()
			button.Callback()
		}()
	}
}

// File dialogs

// Pending file picker channels keyed by dialog ID
var (
	androidFileDialogsLock sync.Mutex
	androidFileResponses   = make(map[uint]chan string)
)

type androidOpenFileDialog struct {
	dialog *OpenFileDialogStruct
}

func newOpenFileDialogImpl(d *OpenFileDialogStruct) openFileDialogImpl {
	return &androidOpenFileDialog{dialog: d}
}

type androidFilePickerOptions struct {
	Multiple bool `json:"multiple"`
}

func (d *androidOpenFileDialog) show() (chan string, error) {
	if d.dialog.canChooseDirectories && !d.dialog.canChooseFiles {
		return nil, fmt.Errorf("directory selection is not supported on Android: the Storage Access Framework returns document-tree URIs, not filesystem paths")
	}

	results := make(chan string, 16)

	androidFileDialogsLock.Lock()
	id := d.dialog.id
	androidFileResponses[id] = results
	androidFileDialogsLock.Unlock()

	optionsJSON, _ := json.Marshal(androidFilePickerOptions{
		Multiple: d.dialog.allowsMultipleSelection,
	})

	androidBridgeVoidIntString("showFilePicker", int(id), string(optionsJSON))
	return results, nil
}

// androidFilePickerResult is invoked from JNI once per selected file.
func androidFilePickerResult(callbackID uint, path string) {
	if path == "" {
		return
	}
	androidFileDialogsLock.Lock()
	channel, ok := androidFileResponses[callbackID]
	androidFileDialogsLock.Unlock()
	if ok {
		channel <- path
	}
}

// androidFilePickerDone is invoked from JNI when the picker finishes
// (after all results, or immediately on cancellation).
func androidFilePickerDone(callbackID uint) {
	androidFileDialogsLock.Lock()
	channel, ok := androidFileResponses[callbackID]
	delete(androidFileResponses, callbackID)
	androidFileDialogsLock.Unlock()
	if ok {
		close(channel)
	}
}

// Save dialogs

type androidSaveFileDialog struct{}

func newSaveFileDialogImpl(_ *SaveFileDialogStruct) saveFileDialogImpl {
	return &androidSaveFileDialog{}
}

func (d *androidSaveFileDialog) show() (chan string, error) {
	return nil, fmt.Errorf("save file dialogs are not supported on Android: write the file inside the app sandbox (e.g. the app's files directory) instead")
}
