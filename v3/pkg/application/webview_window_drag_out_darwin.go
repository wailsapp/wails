//go:build darwin && !ios && !server && !wails_native

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework WebKit
#include "webview_window_drag_out_darwin.h"
#include <stdlib.h>
*/
import "C"

import (
	"errors"
	"os"
	"sync"
	"time"
	"unsafe"
)

// dragOutSession keeps the Go side of a StartDrag call alive until the
// destination has finished with it. Promise data is produced from here when
// the native provider asks for it.
type dragOutSession struct {
	windowID uint
	items    DragItems
}

var (
	dragOutSessionsLock sync.Mutex
	dragOutSessions     = map[uint]*dragOutSession{}
	dragOutSessionSeq   uint
)

// dragOutSessionGrace is how long a finished session is kept so late
// promise writes still find their data.
const dragOutSessionGrace = 2 * time.Minute

type dragOutEndEvent struct {
	windowID  uint
	sessionID uint
	operation DragOperation
}

// dragOutEnded carries NSDraggingSession completions from the main thread to
// the drain loop registered below.
var dragOutEnded = make(chan dragOutEndEvent, 16)

func init() {
	registerChromeEventLoop(func(*App) {
		for {
			event := <-dragOutEnded
			go handleDragOutEnded(event)
		}
	})
}

func newDragOutSession(windowID uint, items DragItems) uint {
	dragOutSessionsLock.Lock()
	defer dragOutSessionsLock.Unlock()
	dragOutSessionSeq++
	id := dragOutSessionSeq
	dragOutSessions[id] = &dragOutSession{windowID: windowID, items: items}
	return id
}

func removeDragOutSession(sessionID uint) {
	dragOutSessionsLock.Lock()
	defer dragOutSessionsLock.Unlock()
	delete(dragOutSessions, sessionID)
}

func lookupDragOutSession(sessionID uint) *dragOutSession {
	dragOutSessionsLock.Lock()
	defer dragOutSessionsLock.Unlock()
	return dragOutSessions[sessionID]
}

func handleDragOutEnded(event dragOutEndEvent) {
	dispatchDragEnd(event.windowID, event.operation)
	if event.operation == DragOperationNone {
		removeDragOutSession(event.sessionID)
		return
	}
	time.AfterFunc(dragOutSessionGrace, func() { removeDragOutSession(event.sessionID) })
}

//export macosDragOutEnded
func macosDragOutEnded(windowID C.uint, sessionID C.uint, operation C.uint) {
	event := dragOutEndEvent{
		windowID:  uint(windowID),
		sessionID: uint(sessionID),
		operation: DragOperation(operation),
	}
	select {
	case dragOutEnded <- event:
	default:
		go func() { dragOutEnded <- event }()
	}
}

// macosDragOutWritePromise produces a promised file. It runs on the promise
// provider's background queue, never on the main thread. It returns NULL on
// success or a malloc'd error message the caller frees.
//
//export macosDragOutWritePromise
func macosDragOutWritePromise(windowID C.uint, sessionID C.uint, index C.int, path *C.char) *C.char {
	session := lookupDragOutSession(uint(sessionID))
	if session == nil || session.windowID != uint(windowID) {
		return C.CString("drag out: the drag session is no longer available")
	}
	if int(index) < 0 || int(index) >= len(session.items.Promises) {
		return C.CString("drag out: unknown promised file")
	}
	promise := session.items.Promises[index]
	data, err := promise.Data()
	if err != nil {
		return C.CString(err.Error())
	}
	if err := os.WriteFile(C.GoString(path), data, 0o644); err != nil {
		return C.CString(err.Error())
	}
	return nil
}

// cStringArray converts Go strings to a NULL-free C array. The returned
// function frees everything. A nil pointer is returned for empty input.
func cStringArray(values []string) (**C.char, func()) {
	if len(values) == 0 {
		return nil, func() {}
	}
	items := make([]*C.char, len(values))
	for i, value := range values {
		items[i] = C.CString(value)
	}
	return (**C.char)(unsafe.Pointer(&items[0])), func() {
		for _, item := range items {
			C.free(unsafe.Pointer(item))
		}
	}
}

func (w *WebviewWindow) startDragOut(items DragItems) error {
	nsWindow := w.NativeWindow()
	if nsWindow == nil {
		return ErrDragOutWindowNotCreated
	}
	sessionID := newDragOutSession(w.id, items)

	cFiles, freeFiles := cStringArray(items.Files)
	defer freeFiles()
	promiseNames := make([]string, len(items.Promises))
	promiseTypes := make([]string, len(items.Promises))
	for i, promise := range items.Promises {
		promiseNames[i] = promise.Filename
		promiseTypes[i] = promise.uti()
	}
	cPromiseNames, freePromiseNames := cStringArray(promiseNames)
	defer freePromiseNames()
	cPromiseTypes, freePromiseTypes := cStringArray(promiseTypes)
	defer freePromiseTypes()
	var cText *C.char
	if items.Text != "" {
		cText = C.CString(items.Text)
		defer C.free(unsafe.Pointer(cText))
	}
	var cImage unsafe.Pointer
	if len(items.Image) > 0 {
		cImage = C.CBytes(items.Image)
		defer C.free(cImage)
	}

	var result C.int
	InvokeSync(func() {
		result = C.dragOutBegin(
			nsWindow,
			C.uint(w.id),
			C.uint(sessionID),
			cFiles, C.int(len(items.Files)),
			cPromiseNames, cPromiseTypes, C.int(len(items.Promises)),
			cText,
			cImage, C.int(len(items.Image)),
			C.int(items.ImageOffset.X), C.int(items.ImageOffset.Y),
			C.uint(items.operations()),
		)
	})

	switch result {
	case C.WailsDragOutStarted:
		return nil
	case C.WailsDragOutNoGesture:
		removeDragOutSession(sessionID)
		return ErrDragOutNoGesture
	case C.WailsDragOutNoWebView:
		removeDragOutSession(sessionID)
		return ErrDragOutWindowNotCreated
	case C.WailsDragOutNoItems:
		removeDragOutSession(sessionID)
		return ErrDragItemsEmpty
	default:
		removeDragOutSession(sessionID)
		return errors.New("drag out: unable to start the drag session")
	}
}

// Non-file drops. The WebviewDrag overlay (webview_window_darwin_drag.m)
// asks which pasteboard types to register and delivers text, URL and image
// drops here.

//export macosDropTypesForWindow
func macosDropTypesForWindow(windowID C.uint) C.int {
	if globalApplication == nil || globalApplication.Window == nil {
		return C.int(dropMaskFiles)
	}
	window, ok := globalApplication.Window.GetByID(uint(windowID))
	if !ok || window == nil {
		return C.int(dropMaskFiles)
	}
	webview, ok := window.(*WebviewWindow)
	if !ok || webview == nil {
		return C.int(dropMaskFiles)
	}
	return C.int(dropTypeMask(webview.options.DropTypes))
}

//export macosOnDrop
func macosOnDrop(windowID C.uint, text *C.char, urls **C.char, urlCount C.int, files **C.char, fileCount C.int, images *unsafe.Pointer, imageLengths *C.int, imageCount C.int, x C.int, y C.int) C.bool {
	data := DropData{X: int(x), Y: int(y)}
	if text != nil {
		data.Text = C.GoString(text)
	}
	if urls != nil && urlCount > 0 {
		for _, url := range unsafe.Slice(urls, int(urlCount)) {
			if url != nil {
				data.URLs = append(data.URLs, C.GoString(url))
			}
		}
	}
	if files != nil && fileCount > 0 {
		for _, file := range unsafe.Slice(files, int(fileCount)) {
			if file != nil {
				data.Files = append(data.Files, C.GoString(file))
			}
		}
	}
	if images != nil && imageLengths != nil && imageCount > 0 {
		pointers := unsafe.Slice(images, int(imageCount))
		lengths := unsafe.Slice(imageLengths, int(imageCount))
		for i, pointer := range pointers {
			if pointer != nil && lengths[i] > 0 {
				data.Images = append(data.Images, C.GoBytes(pointer, lengths[i]))
			}
		}
	}
	if data.IsEmpty() {
		return C.bool(false)
	}
	return C.bool(dispatchDrop(uint(windowID), data))
}
