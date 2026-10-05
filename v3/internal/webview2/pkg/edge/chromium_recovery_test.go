//go:build windows

package edge

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"

	"github.com/wailsapp/wails/v3/internal/webview2/internal/w32"
	"golang.org/x/sys/windows"
)

func TestEmbedWithErrorInvalidBrowserPath(t *testing.T) {
	e := NewChromium()
	e.DataPath = t.TempDir()
	e.BrowserPath = filepath.Join(e.DataPath, "missing-browser")
	e.SetErrorCallback(func(err error) { t.Fatalf("fatal callback: %v", err) })
	if err := e.EmbedWithError(0); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("EmbedWithError error = %v, want missing browser", err)
	}
	if !e.shuttingDown || e.IsReady() {
		t.Fatal("failed instance must be abandoned and not ready")
	}
}

func TestControllerInitializationFailures(t *testing.T) {
	const failure = uintptr(0x80004005)
	cases := []struct {
		name     string
		complete func(*Chromium)
	}{
		{"environment failure", func(e *Chromium) { e.EnvironmentCompleted(failure, nil) }},
		{"missing environment", func(e *Chromium) { e.EnvironmentCompleted(0, nil) }},
		{"controller failure", func(e *Chromium) { e.CreateCoreWebView2ControllerCompleted(failure, nil) }},
		{"missing controller", func(e *Chromium) { e.CreateCoreWebView2ControllerCompleted(0, nil) }},
		{"composition fallback failure", func(e *Chromium) { e.CreateCoreWebView2CompositionControllerCompleted(failure, nil) }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e := NewChromium()
			e.SetErrorCallback(func(err error) { t.Fatalf("fatal callback: %v", err) })
			tc.complete(e)
			if e.initErr == nil || e.IsReady() {
				t.Fatal("callback must record failure without initializing the controller")
			}
			first := e.initErr
			e.initializationFailed(errors.New("later error"))
			if e.initErr != first {
				t.Fatal("first failure was overwritten")
			}
			if e.pumpUntilInited(0) {
				t.Fatal("failed initialization reported ready")
			}
		})
	}
}

func TestAbandonedChromiumIgnoresLateCallbacks(t *testing.T) {
	e := NewChromium()
	e.ShuttingDown()
	e.NavigationStartingCallback = func(*ICoreWebView2) { t.Error("late navigation starting") }
	e.NavigationCompletedCallback = func(*ICoreWebView2, *ICoreWebView2NavigationCompletedEventArgs) { t.Error("late navigation completed") }
	e.ProcessFailedCallback = func(*ICoreWebView2, *ICoreWebView2ProcessFailedEventArgs) { t.Error("late process failure") }
	e.EnvironmentCompleted(0, nil)
	e.CreateCoreWebView2CompositionControllerCompleted(0, nil)
	e.NavigationStarting(nil, nil)
	e.NavigationCompleted(nil, nil)
	e.ProcessFailed(nil, nil)
	closed := false
	controller := &ICoreWebView2Controller{vtbl: &_ICoreWebView2ControllerVtbl{
		Close: NewComProc(func(uintptr) uintptr { closed = true; return 0 }),
	}}
	e.CreateCoreWebView2ControllerCompleted(0, controller)
	if !closed || e.controller != nil || e.initErr != nil {
		t.Fatal("late controller must be closed without reviving the instance")
	}
}

func TestEmbedPumpPreviouslySeenMessage(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	e := NewChromium()
	const readyMessage = 0x8001 // WM_APP + 1
	className, _ := windows.UTF16PtrFromString(t.Name())
	var instance windows.Handle
	if err := windows.GetModuleHandleEx(0, nil, &instance); err != nil {
		t.Fatal(err)
	}
	wc := w32.WndClassExW{
		CbSize:        uint32(unsafe.Sizeof(w32.WndClassExW{})),
		HInstance:     instance,
		LpszClassName: className,
		LpfnWndProc: windows.NewCallback(func(hwnd uintptr, message uint32, wparam, lparam uintptr) uintptr {
			if message == readyMessage {
				atomic.StoreUintptr(&e.inited, 1)
				return 0
			}
			return w32.DefWindowProc(hwnd, uintptr(message), wparam, lparam)
		}),
	}
	if result, _, err := w32.User32RegisterClassExW.Call(uintptr(unsafe.Pointer(&wc))); result == 0 {
		t.Fatal(err)
	}
	user32 := windows.NewLazySystemDLL("user32.dll")
	defer user32.NewProc("UnregisterClassW").Call(uintptr(unsafe.Pointer(className)), uintptr(instance))
	hwnd, _, err := w32.User32CreateWindowExW.Call(0, uintptr(unsafe.Pointer(className)), 0, 0, 0, 0, 0, 0, 0, 0, uintptr(instance), 0)
	if hwnd == 0 {
		t.Fatal(err)
	}
	defer w32.DestroyWindow(hwnd)
	if result, _, err := user32.NewProc("PostMessageW").Call(hwnd, readyMessage, 0, 0); result == 0 {
		t.Fatal(err)
	}
	var msg w32.Msg
	// PM_NOREMOVE marks the message as seen without consuming it. The old wait
	// ignores this input and times out instead of dispatching the completion.
	if got, _, _ := w32.User32PeekMessageW.Call(uintptr(unsafe.Pointer(&msg)), hwnd, readyMessage, readyMessage, 0); got == 0 {
		t.Fatal("completion message was not queued")
	}
	if !e.pumpUntilInited(time.Second) {
		t.Fatal("previously seen completion was not dispatched")
	}
}

func TestEmbedPumpPreservesQuit(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	e := NewChromium()
	w32.User32PostQuitMessage.Call(23)
	if e.pumpUntilInited(time.Second) {
		t.Fatal("quit reported successful initialization")
	}
	var msg w32.Msg
	got, _, _ := w32.User32PeekMessageW.Call(uintptr(unsafe.Pointer(&msg)), 0, w32.WM_QUIT, w32.WM_QUIT, w32.PM_REMOVE)
	if got == 0 || msg.WParam != 23 {
		t.Fatal("outer loop lost its quit message")
	}
}

func TestInitializeControllerFailureReleasesPartialController(t *testing.T) {
	e := NewChromium()
	e.DataPath = t.TempDir()
	e.BrowserPath = filepath.Join(e.DataPath, "missing-browser")
	e.SetErrorCallback(func(err error) { t.Fatalf("fatal callback: %v", err) })
	refs, closes := 0, 0
	controller := &ICoreWebView2Controller{vtbl: &_ICoreWebView2ControllerVtbl{
		_IUnknownVtbl: _IUnknownVtbl{
			QueryInterface: NewComProc(func(uintptr, uintptr, uintptr) uintptr { return 0x80004002 }),
			AddRef:         NewComProc(func(uintptr) uintptr { refs++; return uintptr(refs) }),
			Release:        NewComProc(func(uintptr) uintptr { refs--; return uintptr(refs) }),
		},
		GetCoreWebView2: NewComProc(func(uintptr, uintptr) uintptr { return 0x80004005 }),
		Close:           NewComProc(func(uintptr) uintptr { closes++; return 0 }),
	}}
	e.initializeController(controller)
	if e.initErr == nil || e.webview != nil || e.IsReady() {
		t.Fatal("failed GetCoreWebView2 must stop initialization")
	}
	// Drive the failed-embed cleanup with a deterministic loader failure.
	if e.EmbedWithError(0) == nil {
		t.Fatal("missing browser unexpectedly initialized")
	}
	if refs != 0 || closes != 1 || e.controller != nil {
		t.Fatalf("partial controller leaked: refs=%d closes=%d", refs, closes)
	}
}
