//go:build ios

package application

/*
#include <stdlib.h>
#include "application_ios.h"
*/
import "C"

import (
	"encoding/json"
	"fmt"
	"sync"
	"unsafe"
)

// Message dialogs are backed by UIAlertController, file dialogs by
// UIDocumentPickerViewController. Save dialogs have no iOS counterpart
// (apps write into their sandbox and share via the share sheet), so they
// return an explicit error instead of failing silently.

// Pending message dialogs keyed by callback ID
var (
	iosDialogsLock    sync.Mutex
	iosPendingDialogs = make(map[uint]*MessageDialog)
	iosNextDialogID   uint = 1
)

type iosDialogButton struct {
	Label     string `json:"label"`
	IsCancel  bool   `json:"isCancel"`
	IsDefault bool   `json:"isDefault"`
}

type iosDialog struct {
	dialog *MessageDialog
}

func newDialogImpl(d *MessageDialog) *iosDialog {
	return &iosDialog{dialog: d}
}

func (d *iosDialog) show() {
	iosDialogsLock.Lock()
	id := iosNextDialogID
	iosNextDialogID++
	iosPendingDialogs[id] = d.dialog
	iosDialogsLock.Unlock()

	buttons := make([]iosDialogButton, 0, len(d.dialog.Buttons))
	for _, b := range d.dialog.Buttons {
		buttons = append(buttons, iosDialogButton{
			Label:     b.Label,
			IsCancel:  b.IsCancel,
			IsDefault: b.IsDefault,
		})
	}
	buttonsJSON, _ := json.Marshal(buttons)

	title := d.dialog.Title
	if title == "" {
		title = defaultTitles[d.dialog.DialogType]
	}

	ctitle := C.CString(title)
	cmessage := C.CString(d.dialog.Message)
	cbuttons := C.CString(string(buttonsJSON))
	defer C.free(unsafe.Pointer(ctitle))
	defer C.free(unsafe.Pointer(cmessage))
	defer C.free(unsafe.Pointer(cbuttons))
	C.ios_show_message_dialog(ctitle, cmessage, cbuttons, C.uint(id))
}

//export iosDialogCallback
func iosDialogCallback(callbackID C.uint, buttonIndex C.int) {
	iosDialogsLock.Lock()
	dialog, ok := iosPendingDialogs[uint(callbackID)]
	delete(iosPendingDialogs, uint(callbackID))
	iosDialogsLock.Unlock()
	if !ok || dialog == nil {
		return
	}
	idx := int(buttonIndex)
	if idx < 0 || idx >= len(dialog.Buttons) {
		return
	}
	button := dialog.Buttons[idx]
	if button.Callback != nil {
		// Run the callback off the main thread, mirroring desktop behaviour
		go func() {
			defer handlePanic()
			button.Callback()
		}()
	}
}

// File dialogs

// Pending file picker channels keyed by dialog ID
var (
	iosFileDialogsLock sync.Mutex
	iosFileResponses   = make(map[uint]chan string)
)

type iosOpenFileDialog struct {
	dialog *OpenFileDialogStruct
}

func newOpenFileDialogImpl(d *OpenFileDialogStruct) openFileDialogImpl {
	return &iosOpenFileDialog{dialog: d}
}

func (d *iosOpenFileDialog) show() (chan string, error) {
	results := make(chan string, 16)

	iosFileDialogsLock.Lock()
	id := d.dialog.id
	iosFileResponses[id] = results
	iosFileDialogsLock.Unlock()

	directories := d.dialog.canChooseDirectories && !d.dialog.canChooseFiles
	C.ios_show_document_picker(C.uint(id), C.bool(directories), C.bool(d.dialog.allowsMultipleSelection))
	return results, nil
}

//export iosOpenFileCallback
func iosOpenFileCallback(callbackID C.uint, cpath *C.char) {
	if cpath == nil {
		return
	}
	path := C.GoString(cpath)
	iosFileDialogsLock.Lock()
	channel, ok := iosFileResponses[uint(callbackID)]
	iosFileDialogsLock.Unlock()
	if ok {
		channel <- path
	}
}

//export iosOpenFileCallbackEnd
func iosOpenFileCallbackEnd(callbackID C.uint) {
	iosFileDialogsLock.Lock()
	channel, ok := iosFileResponses[uint(callbackID)]
	delete(iosFileResponses, uint(callbackID))
	iosFileDialogsLock.Unlock()
	if ok {
		close(channel)
	}
}

// Save dialogs

type iosSaveFileDialog struct{}

func newSaveFileDialogImpl(_ *SaveFileDialogStruct) saveFileDialogImpl {
	return &iosSaveFileDialog{}
}

func (d *iosSaveFileDialog) show() (chan string, error) {
	return nil, fmt.Errorf("save file dialogs are not supported on iOS: write the file inside the app sandbox (e.g. the Documents directory) instead")
}
