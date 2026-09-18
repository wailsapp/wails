//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework WebKit
#include "webview_window_mac_extras_darwin.h"
#include <stdlib.h>
*/
import "C"

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

func macWindowExtrasSupported() bool { return true }

func macWindowExtrasSetRepresentedFile(nsWindow unsafe.Pointer, path string) {
	pathC := C.CString(path)
	defer C.free(unsafe.Pointer(pathC))
	InvokeSync(func() { C.windowExtrasSetRepresentedFile(nsWindow, pathC) })
}

func macWindowExtrasRepresentedFile(nsWindow unsafe.Pointer) string {
	return InvokeSyncWithResult(func() string {
		path := C.windowExtrasRepresentedFile(nsWindow)
		if path == nil {
			return ""
		}
		defer C.free(unsafe.Pointer(path))
		return C.GoString(path)
	})
}

func macWindowExtrasSetDocumentEdited(nsWindow unsafe.Pointer, edited bool) {
	InvokeSync(func() { C.windowExtrasSetDocumentEdited(nsWindow, C.bool(edited)) })
}

func macWindowExtrasIsDocumentEdited(nsWindow unsafe.Pointer) bool {
	return InvokeSyncWithResult(func() bool {
		return bool(C.windowExtrasIsDocumentEdited(nsWindow))
	})
}

// macWindowExtrasSetSubtitle returns false when the running macOS release
// has no window subtitles (before 11.0).
func macWindowExtrasSetSubtitle(nsWindow unsafe.Pointer, subtitle string) bool {
	subtitleC := C.CString(subtitle)
	defer C.free(unsafe.Pointer(subtitleC))
	return InvokeSyncWithResult(func() bool {
		return bool(C.windowExtrasSetSubtitle(nsWindow, subtitleC))
	})
}

// macWindowExtrasCascadeNext must run on the application thread; it is
// called from the window creation paths. The first cascade steps away from
// the most recently created other window.
func macWindowExtrasCascadeNext(nsWindow unsafe.Pointer) {
	C.windowExtrasCascadeNext(nsWindow, macWindowExtrasCascadeAnchor(nsWindow))
}

// macWindowExtrasCascadeAnchor returns the native handle of the most recently
// created window other than exclude, preferring WebView windows, or nil.
func macWindowExtrasCascadeAnchor(exclude unsafe.Pointer) unsafe.Pointer {
	if globalApplication == nil {
		return nil
	}
	var anchor unsafe.Pointer
	var highest uint
	if globalApplication.Window != nil {
		for _, window := range globalApplication.Window.GetAll() {
			if window == nil {
				continue
			}
			pointer := window.NativeWindow()
			if pointer == nil || pointer == exclude || window.ID() < highest {
				continue
			}
			anchor, highest = pointer, window.ID()
		}
	}
	if anchor != nil {
		return anchor
	}
	globalApplication.nativeWindowsLock.RLock()
	nativeWindows := make([]*NativeWindow, 0, len(globalApplication.nativeWindows))
	for _, window := range globalApplication.nativeWindows {
		nativeWindows = append(nativeWindows, window)
	}
	globalApplication.nativeWindowsLock.RUnlock()
	for _, window := range nativeWindows {
		if window == nil {
			continue
		}
		pointer := window.NativeWindow()
		if pointer == nil || pointer == exclude || window.ID() < highest {
			continue
		}
		anchor, highest = pointer, window.ID()
	}
	return anchor
}

func macWindowExtrasCascadeFrom(nsWindow, from unsafe.Pointer) {
	InvokeSync(func() { C.windowExtrasCascadeFrom(nsWindow, from) })
}

func macWindowExtrasSetFrameAutosaveName(nsWindow unsafe.Pointer, name string) {
	nameC := C.CString(name)
	defer C.free(unsafe.Pointer(nameC))
	InvokeSync(func() { C.windowExtrasSetFrameAutosaveName(nsWindow, nameC) })
}

// macWindowExtrasRestoreAutosavedFrame must run on the application thread
// during window creation. It reports whether a saved frame was restored.
func macWindowExtrasRestoreAutosavedFrame(nsWindow unsafe.Pointer, name string) bool {
	if name == "" {
		return false
	}
	nameC := C.CString(name)
	defer C.free(unsafe.Pointer(nameC))
	return bool(C.windowExtrasSetFrameAutosaveName(nsWindow, nameC))
}

func macWindowExtrasSetWindowButtonsOffset(nsWindow unsafe.Pointer, x, y int, enabled bool) {
	InvokeSync(func() {
		C.windowExtrasSetWindowButtonsOffset(nsWindow, C.int(x), C.int(y), C.bool(enabled))
	})
}

func macWindowExtrasRequestAttention(critical bool) int {
	return InvokeSyncWithResult(func() int {
		return int(C.windowExtrasRequestAttention(C.bool(critical)))
	})
}

func macWindowExtrasCancelAttention(id int) {
	InvokeSync(func() { C.windowExtrasCancelAttention(C.long(id)) })
}

func macWindowExtrasPrint(nsWindow unsafe.Pointer, options PrintOptions) error {
	printerC := C.CString(options.PrinterName)
	defer C.free(unsafe.Pointer(printerC))
	paperC := C.CString(options.PaperName)
	defer C.free(unsafe.Pointer(paperC))
	margins := options.Margins
	result := InvokeSyncWithResult(func() C.int {
		return C.windowExtrasPrint(nsWindow,
			C.int(options.Orientation), C.bool(!margins.isZero()),
			C.double(margins.Top), C.double(margins.Left), C.double(margins.Bottom), C.double(margins.Right),
			C.bool(options.Silent), printerC, C.double(options.Scale), paperC)
	})
	switch result {
	case C.WailsWindowExtrasPrintStarted:
		return nil
	case C.WailsWindowExtrasPrintPrinterNotFound:
		return fmt.Errorf("%w: %q", ErrMacPrinterNotFound, options.PrinterName)
	case C.WailsWindowExtrasPrintUnsupported:
		return errors.New("printing requires macOS 11 or later")
	default:
		return ErrMacWindowNotCreated
	}
}

// Export requests. WebKit delivers the result asynchronously on the
// application thread; each request owns a buffered channel so the caller
// can wait off that thread with a timeout.

type macWindowExtrasExportResult struct {
	data []byte
	err  error
}

var (
	macWindowExtrasExportLock     sync.Mutex
	macWindowExtrasExportRequests = map[uint64]chan macWindowExtrasExportResult{}
	macWindowExtrasExportNextID   atomic.Uint64
)

func macWindowExtrasNewExportRequest() (uint64, chan macWindowExtrasExportResult) {
	id := macWindowExtrasExportNextID.Add(1)
	ch := make(chan macWindowExtrasExportResult, 1)
	macWindowExtrasExportLock.Lock()
	macWindowExtrasExportRequests[id] = ch
	macWindowExtrasExportLock.Unlock()
	return id, ch
}

func macWindowExtrasTakeExportRequest(id uint64) chan macWindowExtrasExportResult {
	macWindowExtrasExportLock.Lock()
	defer macWindowExtrasExportLock.Unlock()
	ch := macWindowExtrasExportRequests[id]
	delete(macWindowExtrasExportRequests, id)
	return ch
}

//export processMacWindowExtrasExport
func processMacWindowExtrasExport(requestID C.ulonglong, bytes unsafe.Pointer, length C.int, errorMessage *C.char) {
	var result macWindowExtrasExportResult
	if bytes != nil {
		result.data = C.GoBytes(bytes, length)
		C.free(bytes)
	}
	if errorMessage != nil {
		result.err = errors.New(C.GoString(errorMessage))
	} else if len(result.data) == 0 {
		result.err = errors.New("WebKit returned no data")
	}
	// A request that timed out has already been removed; its late result is
	// dropped here.
	if ch := macWindowExtrasTakeExportRequest(uint64(requestID)); ch != nil {
		ch <- result
	}
}

// macWindowExtrasAwaitExport starts the export on the application thread and
// waits for the result on the calling goroutine.
func macWindowExtrasAwaitExport(timeout time.Duration, start func(id uint64) C.int) ([]byte, error) {
	id, ch := macWindowExtrasNewExportRequest()
	started := make(chan C.int, 1)
	InvokeAsync(func() { started <- start(id) })
	timer := time.NewTimer(exportTimeout(timeout))
	defer timer.Stop()
	select {
	case code := <-started:
		switch code {
		case 0:
		case -3:
			macWindowExtrasTakeExportRequest(id)
			return nil, ErrMacExportUnsupported
		default:
			macWindowExtrasTakeExportRequest(id)
			return nil, ErrMacWindowNotCreated
		}
	case <-timer.C:
		macWindowExtrasTakeExportRequest(id)
		return nil, ErrMacExportTimeout
	}
	select {
	case result := <-ch:
		return result.data, result.err
	case <-timer.C:
		macWindowExtrasTakeExportRequest(id)
		return nil, ErrMacExportTimeout
	}
}

func macWindowExtrasExportPDF(nsWindow unsafe.Pointer, options PDFExportOptions) ([]byte, error) {
	rect := options.Rect
	return macWindowExtrasAwaitExport(options.Timeout, func(id uint64) C.int {
		if rect == nil {
			return C.windowExtrasExportPDF(nsWindow, C.ulonglong(id), C.bool(false), 0, 0, 0, 0)
		}
		return C.windowExtrasExportPDF(nsWindow, C.ulonglong(id), C.bool(true),
			C.double(rect.X), C.double(rect.Y), C.double(rect.Width), C.double(rect.Height))
	})
}

func macWindowExtrasSnapshot(nsWindow unsafe.Pointer, options SnapshotOptions) ([]byte, error) {
	rect := options.Rect
	width := C.int(options.Width)
	return macWindowExtrasAwaitExport(options.Timeout, func(id uint64) C.int {
		if rect == nil {
			return C.windowExtrasSnapshot(nsWindow, C.ulonglong(id), C.bool(false), 0, 0, 0, 0, width)
		}
		return C.windowExtrasSnapshot(nsWindow, C.ulonglong(id), C.bool(true),
			C.double(rect.X), C.double(rect.Y), C.double(rect.Width), C.double(rect.Height), width)
	})
}

// macWindowExtrasFrame returns the window frame in screen points (AppKit
// coordinates, origin bottom left). It must run on the application thread.
func macWindowExtrasFrame(nsWindow unsafe.Pointer) (x, y, width, height float64) {
	var cx, cy, cw, ch C.double
	C.windowExtrasFrame(nsWindow, &cx, &cy, &cw, &ch)
	return float64(cx), float64(cy), float64(cw), float64(ch)
}

// applyMacWindowExtras runs on the application thread during window
// creation, after the toolbar and accessories are attached and before the
// initial position is applied. It restores an autosaved frame, positions the
// window buttons and applies values queued before the window existed. It
// reports whether an autosaved frame was restored, in which case the
// options' initial position must not override it.
func (w *macosWebviewWindow) applyMacWindowExtras() bool {
	if w == nil || w.parent == nil || w.nsWindow == nil {
		return false
	}
	macOptions := w.parent.options.Mac
	if offset := macOptions.TitleBar.WindowButtonsOffset; offset != nil {
		C.windowExtrasSetWindowButtonsOffset(w.nsWindow, C.int(offset.X), C.int(offset.Y), C.bool(true))
	}
	if state := takeMacWindowExtrasPending(w.parent); state != nil {
		if state.representedFile != nil {
			pathC := C.CString(*state.representedFile)
			C.windowExtrasSetRepresentedFile(w.nsWindow, pathC)
			C.free(unsafe.Pointer(pathC))
		}
		if state.documentEdited != nil {
			C.windowExtrasSetDocumentEdited(w.nsWindow, C.bool(*state.documentEdited))
		}
		if state.subtitle != nil {
			subtitleC := C.CString(*state.subtitle)
			if !bool(C.windowExtrasSetSubtitle(w.nsWindow, subtitleC)) && globalApplication != nil {
				globalApplication.debug("SetSubtitle requires macOS 11 or later; ignored", "sender", w.parent.Name())
			}
			C.free(unsafe.Pointer(subtitleC))
		}
	}
	return macWindowExtrasRestoreAutosavedFrame(w.nsWindow, macOptions.FrameAutosaveName)
}
