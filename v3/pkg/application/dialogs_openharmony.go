//go:build openharmony && !server

package application

import (
	"encoding/json"
	"fmt"
	"sync"
)

// Message dialogs are backed by the ArkTS prompt (AlertDialog), file dialogs by
// the OHOS picker (selected files are copied into the app's cache so callers
// receive real paths). Save dialogs have no OHOS counterpart that yields a
// filesystem path, so they return an explicit error instead of failing silently.

// Pending message dialogs keyed by callback ID
var (
	openharmonyDialogsLock    sync.Mutex
	openharmonyPendingDialogs      = make(map[uint]*MessageDialog)
	openharmonyNextDialogID   uint = 1
)

type openharmonyDialogButton struct {
	Label     string `json:"label"`
	IsCancel  bool   `json:"isCancel"`
	IsDefault bool   `json:"isDefault"`
}

type openharmonyDialogOptions struct {
	Title   string                    `json:"title"`
	Message string                    `json:"message"`
	Buttons []openharmonyDialogButton `json:"buttons"`
}

type openharmonyDialog struct {
	dialog *MessageDialog
}

func newDialogImpl(d *MessageDialog) *openharmonyDialog {
	return &openharmonyDialog{dialog: d}
}

func (d *openharmonyDialog) show() {
	openharmonyDialogsLock.Lock()
	id := openharmonyNextDialogID
	openharmonyNextDialogID++
	openharmonyPendingDialogs[id] = d.dialog
	openharmonyDialogsLock.Unlock()

	buttons := make([]openharmonyDialogButton, 0, len(d.dialog.Buttons))
	for _, b := range d.dialog.Buttons {
		buttons = append(buttons, openharmonyDialogButton{
			Label:     b.Label,
			IsCancel:  b.IsCancel,
			IsDefault: b.IsDefault,
		})
	}

	title := d.dialog.Title
	if title == "" {
		title = defaultTitles[d.dialog.DialogType]
	}

	optionsJSON, _ := json.Marshal(openharmonyDialogOptions{
		Title:   title,
		Message: d.dialog.Message,
		Buttons: buttons,
	})

	ohBridgeVoidIntString("showMessageDialog", int(id), string(optionsJSON))
}

// openharmonyDialogCallback is invoked from ArkTS when a dialog button is
// pressed (buttonIndex is the index into the dialog's button slice, or -1 for a
// dismissal with no matching button).
func openharmonyDialogCallback(callbackID uint, buttonIndex int) {
	openharmonyDialogsLock.Lock()
	dialog, ok := openharmonyPendingDialogs[callbackID]
	delete(openharmonyPendingDialogs, callbackID)
	openharmonyDialogsLock.Unlock()
	if !ok || dialog == nil {
		return
	}
	if buttonIndex < 0 || buttonIndex >= len(dialog.Buttons) {
		return
	}
	button := dialog.Buttons[buttonIndex]
	if button.Callback != nil {
		go func() {
			defer handlePanic()
			button.Callback()
		}()
	}
}

// File dialogs

var (
	openharmonyFileDialogsLock sync.Mutex
	openharmonyFileResponses   = make(map[uint]chan string)
)

type openharmonyOpenFileDialog struct {
	dialog *OpenFileDialogStruct
}

func newOpenFileDialogImpl(d *OpenFileDialogStruct) openFileDialogImpl {
	return &openharmonyOpenFileDialog{dialog: d}
}

type openharmonyFilePickerOptions struct {
	Multiple bool `json:"multiple"`
}

func (d *openharmonyOpenFileDialog) show() (chan string, error) {
	if d.dialog.canChooseDirectories && !d.dialog.canChooseFiles {
		return nil, fmt.Errorf("directory selection is not supported on OHOS: the picker returns URIs, not filesystem paths")
	}

	results := make(chan string, 16)

	openharmonyFileDialogsLock.Lock()
	id := d.dialog.id
	openharmonyFileResponses[id] = results
	openharmonyFileDialogsLock.Unlock()

	optionsJSON, _ := json.Marshal(openharmonyFilePickerOptions{
		Multiple: d.dialog.allowsMultipleSelection,
	})

	ohBridgeVoidIntString("showFilePicker", int(id), string(optionsJSON))
	return results, nil
}

func openharmonyFilePickerResult(callbackID uint, path string) {
	if path == "" {
		return
	}
	openharmonyFileDialogsLock.Lock()
	channel, ok := openharmonyFileResponses[callbackID]
	openharmonyFileDialogsLock.Unlock()
	if ok {
		channel <- path
	}
}

func openharmonyFilePickerDone(callbackID uint) {
	openharmonyFileDialogsLock.Lock()
	channel, ok := openharmonyFileResponses[callbackID]
	delete(openharmonyFileResponses, callbackID)
	openharmonyFileDialogsLock.Unlock()
	if ok {
		close(channel)
	}
}

// Save dialogs

type openharmonySaveFileDialog struct {
	dialog *SaveFileDialogStruct
}

func newSaveFileDialogImpl(d *SaveFileDialogStruct) saveFileDialogImpl {
	return &openharmonySaveFileDialog{dialog: d}
}

func (d *openharmonySaveFileDialog) show() (chan string, error) {
	results := make(chan string, 1)

	openharmonyFileDialogsLock.Lock()
	id := d.dialog.id
	openharmonyFileResponses[id] = results
	openharmonyFileDialogsLock.Unlock()

	filters := make([]string, 0, len(d.dialog.filters))
	for _, filter := range d.dialog.filters {
		filters = append(filters, filter.Pattern)
	}
	optionsJSON, _ := json.Marshal(struct {
		Title    string   `json:"title"`
		Filename string   `json:"filename"`
		Filters  []string `json:"filters"`
	}{
		Title:    d.dialog.title,
		Filename: d.dialog.filename,
		Filters:  filters,
	})

	ohBridgeVoidIntString("saveFilePicker", int(id), string(optionsJSON))
	return results, nil
}
