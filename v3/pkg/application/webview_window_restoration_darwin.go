//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework WebKit
#include "webview_window_restoration_darwin.h"
#include <stdlib.h>
*/
import "C"

import (
	"time"
	"unsafe"
)

func macWindowRestorationSupported() bool { return true }

func macWindowRestorationConfigure(nsWindow unsafe.Pointer, identifier string) {
	identifierC := C.CString(identifier)
	defer C.free(unsafe.Pointer(identifierC))
	InvokeSync(func() { C.windowRestorationConfigure(nsWindow, identifierC) })
}

func macWindowRestorationSetData(nsWindow unsafe.Pointer, encoded string) {
	encodedC := C.CString(encoded)
	defer C.free(unsafe.Pointer(encodedC))
	InvokeSync(func() { C.windowRestorationSetData(nsWindow, encodedC) })
}

func macWindowRestorationComplete(requestID uint64, nsWindow unsafe.Pointer, message string) {
	if globalApplication == nil || globalApplication.impl == nil {
		// AppKit cannot have asked for a restore without a running app.
		return
	}
	messageC := C.CString(message)
	defer C.free(unsafe.Pointer(messageC))
	InvokeSync(func() { C.windowRestorationComplete(C.ulonglong(requestID), nsWindow, messageC) })
}

// macWindowRestorationAwaitNative waits for the native window of a window
// returned by the OnRestore handler. Creation is dispatched to the
// application thread ahead of this call, so the first InvokeSync normally
// finds it; the loop covers a window whose creation was deferred.
func macWindowRestorationAwaitNative(window Window) unsafe.Pointer {
	deadline := time.Now().Add(10 * time.Second)
	for {
		nsWindow := InvokeSyncWithResult(func() unsafe.Pointer { return window.NativeWindow() })
		if nsWindow != nil {
			return nsWindow
		}
		if webview, ok := window.(*WebviewWindow); ok && (webview == nil || webview.isDestroyed()) {
			return nil
		}
		if time.Now().After(deadline) {
			return nil
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func macWindowRestorationInteractionState(nsWindow unsafe.Pointer) ([]byte, error) {
	var data []byte
	status := InvokeSyncWithResult(func() C.int {
		var bytes unsafe.Pointer
		var length C.int
		result := C.windowRestorationInteractionState(nsWindow, &bytes, &length)
		if bytes != nil {
			data = C.GoBytes(bytes, length)
			C.free(bytes)
		}
		return result
	})
	switch status {
	case C.WailsRestorationStateOK:
		return data, nil
	case C.WailsRestorationStateUnsupported:
		return nil, ErrMacInteractionStateUnsupported
	case C.WailsRestorationStateInvalid:
		return nil, ErrMacInteractionStateInvalid
	default:
		return nil, ErrMacWindowNotCreated
	}
}

func macWindowRestorationSetInteractionState(nsWindow unsafe.Pointer, state []byte) error {
	if len(state) == 0 {
		return ErrMacInteractionStateInvalid
	}
	status := InvokeSyncWithResult(func() C.int {
		return C.windowRestorationSetInteractionState(nsWindow, unsafe.Pointer(&state[0]), C.int(len(state)))
	})
	switch status {
	case C.WailsRestorationStateOK:
		return nil
	case C.WailsRestorationStateUnsupported:
		return ErrMacInteractionStateUnsupported
	case C.WailsRestorationStateInvalid:
		return ErrMacInteractionStateInvalid
	default:
		return ErrMacWindowNotCreated
	}
}

// processMacWindowRestore is called by WailsWindowRestoration on the
// application thread at launch, once per saved window. The request is
// queued so the OnRestore handler can create windows without blocking
// AppKit; the completion handler is called later.
//
//export processMacWindowRestore
func processMacWindowRestore(requestID C.ulonglong, identifier *C.char, encoded *C.char) {
	noteMacWindowRestoreRequest(macWindowRestoreRequest{
		requestID: uint64(requestID),
		id:        C.GoString(identifier),
		encoded:   C.GoString(encoded),
	})
}

// processMacWindowRestorationVisible is called on the application thread
// the first time a window comes on screen. It applies a restoration
// identifier configured before the native window existed (through
// MacWindow.RestorationID or an early SetRestorationID).
//
//export processMacWindowRestorationVisible
func processMacWindowRestorationVisible(nsWindow unsafe.Pointer) {
	window, ok := windowForNSWindow(nsWindow).(*WebviewWindow)
	if !ok || window == nil {
		return
	}
	if id := window.RestorationID(); id != "" {
		window.applyRestorationID(nsWindow, id)
	}
}

// HandleSupportsSecureRestorableState answers the application delegate's
// applicationSupportsSecureRestorableState: from
// MacOptions.SupportsSecureRestorableState.
//
//export HandleSupportsSecureRestorableState
func HandleSupportsSecureRestorableState() C.bool {
	if globalApplication == nil {
		return C.bool(false)
	}
	return C.bool(globalApplication.options.Mac.SupportsSecureRestorableState)
}
