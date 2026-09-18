package application

import (
	"errors"
	"testing"
)

func TestMacSheetResponseCodesMatchNSModalResponse(t *testing.T) {
	if MacSheetResponseCancel != 0 || MacSheetResponseOK != 1 || MacSheetResponseStop != -1000 || MacSheetResponseAbort != -1001 {
		t.Fatal("sheet response codes do not match NSModalResponse")
	}
}

func TestPresentSheetErrorPaths(t *testing.T) {
	parent := &WebviewWindow{options: WebviewWindowOptions{Name: "parent"}}
	sheet := &WebviewWindow{options: WebviewWindowOptions{Name: "sheet"}}
	var nilWindow *WebviewWindow
	var nilTyped *WebviewWindow

	if !macSheetsSupported() {
		if err := parent.PresentSheet(sheet); !errors.Is(err, ErrMacSheetUnsupported) {
			t.Fatalf("PresentSheet off macOS = %v, want ErrMacSheetUnsupported", err)
		}
		if err := parent.PresentNativeSheet(&NativeWindow{}, false); !errors.Is(err, ErrMacSheetUnsupported) {
			t.Fatalf("PresentNativeSheet off macOS = %v, want ErrMacSheetUnsupported", err)
		}
		return
	}
	if err := nilWindow.PresentSheet(sheet); !errors.Is(err, ErrMacWindowNotCreated) {
		t.Fatalf("nil receiver = %v, want ErrMacWindowNotCreated", err)
	}
	if err := parent.PresentSheet(nil); !errors.Is(err, ErrMacSheetRequired) {
		t.Fatalf("nil sheet = %v, want ErrMacSheetRequired", err)
	}
	if err := parent.PresentSheet(nilTyped); !errors.Is(err, ErrMacSheetRequired) {
		t.Fatalf("typed nil sheet = %v, want ErrMacSheetRequired", err)
	}
	if err := parent.PresentSheet(parent); !errors.Is(err, ErrMacSheetSelf) {
		t.Fatalf("self as sheet = %v, want ErrMacSheetSelf", err)
	}
	if err := parent.PresentCriticalSheet(parent); !errors.Is(err, ErrMacSheetSelf) {
		t.Fatalf("self as critical sheet = %v, want ErrMacSheetSelf", err)
	}
	// Neither window has a native window yet: the receiver is checked first.
	if err := parent.PresentSheet(sheet); !errors.Is(err, ErrMacWindowNotCreated) {
		t.Fatalf("uncreated parent = %v, want ErrMacWindowNotCreated", err)
	}
	if err := parent.PresentNativeSheet(nil, false); !errors.Is(err, ErrMacSheetRequired) {
		t.Fatalf("nil native sheet = %v, want ErrMacSheetRequired", err)
	}
	if err := parent.PresentNativeSheet(&NativeWindow{}, true); !errors.Is(err, ErrMacWindowNotCreated) {
		t.Fatalf("uncreated parent with native sheet = %v, want ErrMacWindowNotCreated", err)
	}
}

func TestSheetQueriesAreSafeBeforeCreation(t *testing.T) {
	window := &WebviewWindow{options: WebviewWindowOptions{Name: "uncreated"}}
	var nilWindow *WebviewWindow
	var nilNative *NativeWindow
	for _, w := range []*WebviewWindow{window, nilWindow} {
		if w.EndSheet(MacSheetResponseOK) != w {
			t.Fatal("EndSheet must return its receiver")
		}
		if w.IsSheet() || w.SheetParent() != nil || w.AttachedSheet() != nil || w.AttachedNativeSheet() != nil || w.HasAttachedSheet() {
			t.Fatal("an uncreated window reported sheet state")
		}
	}
	if nilNative.EndSheet(0) != nil {
		t.Fatal("EndSheet on a nil native window must return nil")
	}
	nilWindow.OnSheetEnd(func(int) {})()
	nilNative.OnSheetEnd(func(int) {})()
}

func TestSheetEndCallbacksRegisterAndRemove(t *testing.T) {
	window := &WebviewWindow{options: WebviewWindowOptions{Name: "sheet"}}
	native := &NativeWindow{}
	var received []int
	remove := window.OnSheetEnd(func(code int) { received = append(received, code) })
	keep := window.OnSheetEnd(func(code int) { received = append(received, code*10) })
	defer keep()
	nativeRemove := native.OnSheetEnd(func(int) {})
	defer nativeRemove()

	macSheetCallbacksLock.RLock()
	count := len(macSheetCallbacks[window])
	nativeCount := len(macSheetCallbacks[native])
	macSheetCallbacksLock.RUnlock()
	if count != 2 || nativeCount != 1 {
		t.Fatalf("registered callbacks = %d webview, %d native", count, nativeCount)
	}
	for _, entry := range func() []*macSheetCallback {
		macSheetCallbacksLock.RLock()
		defer macSheetCallbacksLock.RUnlock()
		return append([]*macSheetCallback(nil), macSheetCallbacks[window]...)
	}() {
		entry.callback(MacSheetResponseOK)
	}
	if len(received) != 2 || received[0] != 1 || received[1] != 10 {
		t.Fatalf("callbacks delivered %v", received)
	}
	remove()
	remove() // removing twice is harmless
	macSheetCallbacksLock.RLock()
	count = len(macSheetCallbacks[window])
	macSheetCallbacksLock.RUnlock()
	if count != 1 {
		t.Fatalf("callbacks after removal = %d, want 1", count)
	}
	keep()
	macSheetCallbacksLock.RLock()
	_, present := macSheetCallbacks[window]
	macSheetCallbacksLock.RUnlock()
	if present {
		t.Fatal("empty callback list was not deleted")
	}
	if window.OnSheetEnd(nil) == nil {
		t.Fatal("a nil callback must still return a remover")
	}
}

func TestHandleMacSheetEndedIgnoresUnknownWindows(t *testing.T) {
	// No application, no windows: the event must be dropped without panicking.
	handleMacSheetEnded(macSheetEndEvent{sheet: nil, code: MacSheetResponseStop})
}
