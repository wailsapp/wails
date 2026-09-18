package application

import (
	"errors"
	"sync"
	"unsafe"
)

// This file is the cross-platform surface of macOS window sheets: presenting
// a second window as a window-modal sheet (beginSheet:completionHandler: and
// beginCriticalSheet:), ending it with a response code, and finding the
// sheet or its parent. The platform functions it calls live in
// webview_window_sheet_darwin.go, with documented stubs in
// webview_window_sheet_other.go, mirroring webview_window_tabs.go.
//
// A sheet window is an ordinary window: create it with Frameless false (a
// titled window; AppKit hides the title bar buttons while it is a sheet) and
// Hidden true so it does not flash on screen before it is attached. AppKit
// orders it out again when the sheet ends, so the same window can be
// presented repeatedly.

var (
	// ErrMacSheetRequired means a nil or not-yet-created window was passed
	// to PresentSheet.
	ErrMacSheetRequired = errors.New("a created window to present as a sheet is required")
	// ErrMacSheetSelf means a window was asked to present itself as a sheet.
	ErrMacSheetSelf = errors.New("a window cannot be presented as a sheet of itself")
	// ErrMacSheetUnsupported means window sheets are unavailable on the
	// current platform.
	ErrMacSheetUnsupported = errors.New("window sheets are unavailable on this platform")
)

// Sheet response codes. EndSheet accepts any int; these mirror the
// NSModalResponse values AppKit itself uses.
const (
	// MacSheetResponseCancel is NSModalResponseCancel.
	MacSheetResponseCancel = 0
	// MacSheetResponseOK is NSModalResponseOK.
	MacSheetResponseOK = 1
	// MacSheetResponseStop is NSModalResponseStop, which AppKit delivers
	// when the sheet window is closed instead of ended.
	MacSheetResponseStop = -1000
	// MacSheetResponseAbort is NSModalResponseAbort.
	MacSheetResponseAbort = -1001
)

// macSheetEndEvent is delivered by the native completion handler when a
// sheet ends. sheet is the sheet's NSWindow.
type macSheetEndEvent struct {
	sheet unsafe.Pointer
	code  int
}

var macSheetEnded = make(chan macSheetEndEvent, 16)

// macSheetCallback is one OnSheetEnd registration.
type macSheetCallback struct {
	callback func(code int)
}

var (
	macSheetCallbacksLock sync.RWMutex
	// macSheetCallbacks is keyed by the sheet's Go window (*WebviewWindow
	// or *NativeWindow) so callbacks can be registered before the native
	// window exists.
	macSheetCallbacks = map[any][]*macSheetCallback{}
)

func init() {
	registerChromeEventLoop(func(*App) {
		for {
			event := <-macSheetEnded
			go handleMacSheetEnded(event)
		}
	})
}

func handleMacSheetEnded(event macSheetEndEvent) {
	defer handlePanic()
	var key any
	if window := windowForNSWindow(event.sheet); window != nil {
		key = window
	} else if native := nativeWindowForNSWindow(event.sheet); native != nil {
		key = native
	}
	if key == nil {
		return
	}
	macSheetCallbacksLock.RLock()
	callbacks := append([]*macSheetCallback(nil), macSheetCallbacks[key]...)
	macSheetCallbacksLock.RUnlock()
	for _, entry := range callbacks {
		entry.callback(event.code)
	}
}

func addMacSheetCallback(key any, callback func(code int)) func() {
	if callback == nil {
		return func() {}
	}
	entry := &macSheetCallback{callback: callback}
	macSheetCallbacksLock.Lock()
	macSheetCallbacks[key] = append(macSheetCallbacks[key], entry)
	macSheetCallbacksLock.Unlock()
	return func() {
		macSheetCallbacksLock.Lock()
		defer macSheetCallbacksLock.Unlock()
		remaining := macSheetCallbacks[key][:0]
		for _, existing := range macSheetCallbacks[key] {
			if existing != entry {
				remaining = append(remaining, existing)
			}
		}
		if len(remaining) == 0 {
			delete(macSheetCallbacks, key)
		} else {
			macSheetCallbacks[key] = remaining
		}
	}
}

// OnSheetEnd registers a callback that runs with the response code when
// this window, presented as a sheet, is ended with EndSheet or closed
// (MacSheetResponseStop). It returns a function that removes the callback.
// It can be called before the native window exists.
func (w *WebviewWindow) OnSheetEnd(callback func(code int)) func() {
	if w == nil {
		return func() {}
	}
	return addMacSheetCallback(w, callback)
}

// OnSheetEnd registers a callback that runs with the response code when
// this native window, presented as a sheet, ends. See
// WebviewWindow.OnSheetEnd.
func (w *NativeWindow) OnSheetEnd(callback func(code int)) func() {
	if w == nil {
		return func() {}
	}
	return addMacSheetCallback(w, callback)
}

// PresentSheet attaches other to this window as a window-modal sheet. Both
// native windows must exist; a window cannot be a sheet of itself. When
// this window already shows a sheet the new one is queued behind it, as
// AppKit does. The sheet ends with other.EndSheet or when other is closed,
// and other's OnSheetEnd callbacks receive the response code. Off macOS it
// returns ErrMacSheetUnsupported.
func (w *WebviewWindow) PresentSheet(other Window) error {
	return w.presentSheet(other, false)
}

// PresentCriticalSheet attaches other as a critical sheet
// (beginCriticalSheet:completionHandler:), which is shown in front of any
// ordinary sheet already attached instead of waiting behind it. Use it for
// sheets the user must see now, such as an error that interrupts a task.
func (w *WebviewWindow) PresentCriticalSheet(other Window) error {
	return w.presentSheet(other, true)
}

func (w *WebviewWindow) presentSheet(other Window, critical bool) error {
	if !macSheetsSupported() {
		return ErrMacSheetUnsupported
	}
	if w == nil {
		return ErrMacWindowNotCreated
	}
	if other == nil {
		return ErrMacSheetRequired
	}
	if otherWebview, ok := other.(*WebviewWindow); ok {
		if otherWebview == nil {
			return ErrMacSheetRequired
		}
		if otherWebview == w {
			return ErrMacSheetSelf
		}
	}
	nsWindow := w.NativeWindow()
	if nsWindow == nil {
		return ErrMacWindowNotCreated
	}
	sheet := other.NativeWindow()
	if sheet == nil {
		return ErrMacSheetRequired
	}
	return macSheetBegin(nsWindow, sheet, critical)
}

// PresentNativeSheet attaches a NativeWindow to this window as a sheet, or
// as a critical sheet when critical is true. See PresentSheet.
func (w *WebviewWindow) PresentNativeSheet(other *NativeWindow, critical bool) error {
	if !macSheetsSupported() {
		return ErrMacSheetUnsupported
	}
	if w == nil {
		return ErrMacWindowNotCreated
	}
	if other == nil {
		return ErrMacSheetRequired
	}
	nsWindow := w.NativeWindow()
	if nsWindow == nil {
		return ErrMacWindowNotCreated
	}
	sheet := other.NativeWindow()
	if sheet == nil {
		return ErrMacSheetRequired
	}
	return macSheetBegin(nsWindow, sheet, critical)
}

// EndSheet ends this window's presentation as a sheet with the given
// response code and orders it out; the parent's completion delivers the
// code to OnSheetEnd. It is a no-op when the window is not a sheet, before
// the native window exists, and off macOS.
func (w *WebviewWindow) EndSheet(code int) Window {
	if w == nil {
		return w
	}
	if nsWindow := w.NativeWindow(); nsWindow != nil {
		macSheetEnd(nsWindow, code)
	}
	return w
}

// EndSheet ends this native window's presentation as a sheet. See
// WebviewWindow.EndSheet.
func (w *NativeWindow) EndSheet(code int) *NativeWindow {
	if w == nil {
		return w
	}
	if nsWindow := w.NativeWindow(); nsWindow != nil {
		macSheetEnd(nsWindow, code)
	}
	return w
}

// IsSheet reports whether the window is currently presented as a sheet.
func (w *WebviewWindow) IsSheet() bool {
	if w == nil {
		return false
	}
	nsWindow := w.NativeWindow()
	return nsWindow != nil && macSheetParent(nsWindow) != nil
}

// SheetParent returns the WebviewWindow this window is currently attached
// to as a sheet, or nil when it is not a sheet, when the parent is not a
// Wails WebView window, and off macOS.
func (w *WebviewWindow) SheetParent() Window {
	if w == nil {
		return nil
	}
	nsWindow := w.NativeWindow()
	if nsWindow == nil {
		return nil
	}
	parent := macSheetParent(nsWindow)
	if parent == nil {
		return nil
	}
	if window := windowForNSWindow(parent); window != nil {
		return window
	}
	return nil
}

// AttachedSheet returns the WebviewWindow currently attached to this window
// as a sheet, or nil when there is none, when the sheet is not a Wails
// WebView window (see AttachedNativeSheet), and off macOS.
func (w *WebviewWindow) AttachedSheet() Window {
	sheet := w.attachedSheetPointer()
	if sheet == nil {
		return nil
	}
	if window := windowForNSWindow(sheet); window != nil {
		return window
	}
	return nil
}

// AttachedNativeSheet returns the NativeWindow currently attached to this
// window as a sheet, or nil.
func (w *WebviewWindow) AttachedNativeSheet() *NativeWindow {
	sheet := w.attachedSheetPointer()
	if sheet == nil {
		return nil
	}
	return nativeWindowForNSWindow(sheet)
}

// HasAttachedSheet reports whether any sheet is attached to this window.
func (w *WebviewWindow) HasAttachedSheet() bool {
	return w.attachedSheetPointer() != nil
}

func (w *WebviewWindow) attachedSheetPointer() unsafe.Pointer {
	if w == nil {
		return nil
	}
	nsWindow := w.NativeWindow()
	if nsWindow == nil {
		return nil
	}
	return macSheetAttached(nsWindow)
}
