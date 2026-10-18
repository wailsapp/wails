//go:build windows

package edge

import (
	"testing"
	"unsafe"
)

// Runtime COM failures after startup must be logged and survived, never routed
// through errorCallback, which terminates the process (wailsapp/wails#5580).
// Each test installs a fatal error callback that fails the test, then makes the
// COM call behind one runtime site fail.

const runtimeFailure = uintptr(0x80004005) // E_FAIL

func fatalOnErrorCallback(t *testing.T) *Chromium {
	t.Helper()
	e := NewChromium()
	e.SetErrorCallback(func(err error) { t.Fatalf("fatal callback: %v", err) })
	return e
}

func TestCursorChangedSurvivesCOMFailures(t *testing.T) {
	cases := []struct {
		name           string
		cursorHR, idHR uintptr
	}{
		{"GetCursor fails", runtimeFailure, 0},
		{"GetSystemCursorId fails", 0, runtimeFailure},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sender := &ICoreWebView2CompositionController{Vtbl: &ICoreWebView2CompositionControllerVtbl{
				GetCursor:         NewComProc(func(uintptr, uintptr) uintptr { return tc.cursorHR }),
				GetSystemCursorId: NewComProc(func(uintptr, uintptr) uintptr { return tc.idHR }),
			}}
			e := fatalOnErrorCallback(t)
			e.CursorChangedCallback = func(HCURSOR, uint32) { t.Error("callback ran with an unreadable cursor") }
			e.CursorChanged(sender, nil)
		})
	}
}

func TestSetBackgroundColourSurvivesCOMFailure(t *testing.T) {
	called := false
	controller2 := &ICoreWebView2Controller2{vtbl: &_ICoreWebView2Controller2Vtbl{
		_IUnknownVtbl: _IUnknownVtbl{Release: NewComProc(func(uintptr) uintptr { return 0 })},
		PutDefaultBackgroundColor: NewComProc(func(uintptr, uintptr) uintptr {
			called = true
			return runtimeFailure
		}),
	}}
	controller := &ICoreWebView2Controller{vtbl: &_ICoreWebView2ControllerVtbl{
		_IUnknownVtbl: _IUnknownVtbl{
			QueryInterface: NewComProc(func(_, _, ppvObject uintptr) uintptr {
				*(**ICoreWebView2Controller2)(unsafe.Pointer(ppvObject)) = controller2
				return 0
			}),
		},
	}}
	e := fatalOnErrorCallback(t)
	e.controller = controller
	e.SetBackgroundColour(1, 2, 3, 255)
	if !called {
		t.Fatal("PutDefaultBackgroundColor was not called")
	}
}

func TestAddWebResourceRequestedFilterSurvivesCOMFailure(t *testing.T) {
	webview := &ICoreWebView2{vtbl: &iCoreWebView2Vtbl{
		AddWebResourceRequestedFilter: NewComProc(func(uintptr, uintptr, uintptr) uintptr { return runtimeFailure }),
	}}
	e := fatalOnErrorCallback(t)
	e.webview = webview
	e.AddWebResourceRequestedFilter("*", COREWEBVIEW2_WEB_RESOURCE_CONTEXT_ALL)
}

func TestPermissionRequestedSurvivesCOMFailures(t *testing.T) {
	allow := CoreWebView2PermissionStateAllow
	cases := []struct {
		name    string
		kindHR  uintptr
		putHR   uintptr
		global  *CoreWebView2PermissionState
		want    CoreWebView2PermissionState
		wantPut bool
	}{
		{"unreadable kind is denied", runtimeFailure, 0, nil, CoreWebView2PermissionStateDeny, true},
		{"unreadable kind uses global permission", runtimeFailure, 0, &allow, CoreWebView2PermissionStateAllow, true},
		{"PutState fails", 0, runtimeFailure, nil, CoreWebView2PermissionStateDefault, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			put := false
			var got CoreWebView2PermissionState
			args := &iCoreWebView2PermissionRequestedEventArgs{vtbl: &iCoreWebView2PermissionRequestedEventArgsVtbl{
				GetPermissionKind: NewComProc(func(uintptr, uintptr) uintptr { return tc.kindHR }),
				PutState: NewComProc(func(_ uintptr, state uintptr) uintptr {
					put, got = true, CoreWebView2PermissionState(state)
					return tc.putHR
				}),
			}}
			e := fatalOnErrorCallback(t)
			e.globalPermission = tc.global
			e.PermissionRequested(nil, args)
			if put != tc.wantPut || got != tc.want {
				t.Fatalf("PutState called=%v with %v, want called=%v with %v", put, got, tc.wantPut, tc.want)
			}
		})
	}
}

func TestAcceleratorKeyPressedSurvivesPutHandledFailure(t *testing.T) {
	calls := 0
	args := &ICoreWebView2AcceleratorKeyPressedEventArgs{vtbl: &_ICoreWebView2AcceleratorKeyPressedEventArgsVtbl{
		// Zeroed out-params: a KEY_DOWN with WasKeyDown false.
		GetKeyEventKind:      NewComProc(func(uintptr, uintptr) uintptr { return 0 }),
		GetVirtualKey:        NewComProc(func(uintptr, uintptr) uintptr { return 0 }),
		GetPhysicalKeyStatus: NewComProc(func(uintptr, uintptr) uintptr { return 0 }),
		PutHandled:           NewComProc(func(uintptr, uintptr) uintptr { calls++; return runtimeFailure }),
	}}
	e := fatalOnErrorCallback(t)
	e.AcceleratorKeyCallback = func(uint) bool { return true }
	e.AcceleratorKeyPressed(nil, args)
	if calls == 0 {
		t.Fatal("PutHandled was not called")
	}
}

func TestBoundsSurvivesCOMFailure(t *testing.T) {
	e := fatalOnErrorCallback(t)
	e.controller = &ICoreWebView2Controller{vtbl: &_ICoreWebView2ControllerVtbl{
		GetBounds: NewComProc(func(uintptr, uintptr) uintptr { return runtimeFailure }),
	}}
	if r := e.Bounds(); r != nil {
		t.Fatalf("Bounds() = %v after a failed GetBounds, want nil", r)
	}
}
