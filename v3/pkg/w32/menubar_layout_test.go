//go:build windows

package w32

import (
	"testing"
	"unsafe"
)

// The UAH* structs mirror undocumented Windows structs (uxtheme's UAH menu
// messages); their layout must match the C ABI exactly or the menubar code
// reads the wrong fields (#6311).
func TestUAHMenuStructLayout(t *testing.T) {
	if got := unsafe.Offsetof(UAHDRAWMENUITEM{}.UM); got != unsafe.Sizeof(DRAWITEMSTRUCT{}) {
		t.Errorf("UAHDRAWMENUITEM.UM at %d, want right after DRAWITEMSTRUCT (%d)", got, unsafe.Sizeof(DRAWITEMSTRUCT{}))
	}
	// C: int iPosition (4) + UAHMENUITEMMETRICS (32) + UAHMENUPOPUPMETRICS (20).
	if got := unsafe.Offsetof(UAHMENUITEM{}.Umim); got != 4 {
		t.Errorf("UAHMENUITEM.Umim at %d, want 4 (iPosition is a C int)", got)
	}
	if got := unsafe.Sizeof(UAHMENUITEM{}); got != 56 {
		t.Errorf("sizeof(UAHMENUITEM) = %d, want 56", got)
	}
}

// A bad device context makes the drawing helpers report failure instead of
// panicking inside a window procedure (#6311).
func TestDrawingHelpersNeverPanic(t *testing.T) {
	if IsValidDC(0) {
		t.Error("IsValidDC(0) = true")
	}
	if IsValidDC(^HDC(0)) {
		t.Error("IsValidDC(all-ones handle) = true")
	}
	if TrySetTextColor(0, 0xdedede) {
		t.Error("TrySetTextColor on a nil DC reported success")
	}
	dc := memoryDC(t)
	defer DeleteDC(dc)
	if !IsValidDC(dc) {
		t.Error("IsValidDC(memory DC) = false")
	}
	if !TrySetTextColor(dc, 0xdedede) {
		t.Error("TrySetTextColor on a memory DC failed")
	}
}

// A DC whose previous text colour is CLR_INVALID makes SetTextColor return
// CLR_INVALID even though the call worked; that is not a failure (#6311).
func TestTrySetTextColorAcceptsInvalidPreviousColour(t *testing.T) {
	dc := memoryDC(t)
	defer DeleteDC(dc)
	procSetTextColor.Call(uintptr(dc), uintptr(CLR_INVALID))
	if got, _, _ := procGetTextColor.Call(uintptr(dc)); uint32(got) != CLR_INVALID {
		t.Skip("this Windows build does not let a DC hold CLR_INVALID")
	}
	if !TrySetTextColor(dc, 0xdedede) {
		t.Error("TrySetTextColor reported failure for a DC whose previous colour was CLR_INVALID")
	}
	if got, _, _ := procGetTextColor.Call(uintptr(dc)); uint32(got) != 0xdedede {
		t.Errorf("text colour = %#x, want 0xdedede", got)
	}
}

// MenuBarWndProc must contain a panic and a missing message parameter.
func TestMenuBarWndProcContainsFailures(t *testing.T) {
	if handled, _ := MenuBarWndProc(0, WM_UAHDRAWMENUITEM, 0, 0, &MenuBarTheme{}); handled {
		t.Error("a nil lParam was reported as handled")
	}

	orig := menuBarImpl
	defer func() { menuBarImpl = orig }()
	menuBarImpl = func(HWND, uint32, WPARAM, LPARAM, *MenuBarTheme) (bool, LRESULT) {
		panic("boom")
	}
	if handled, result := MenuBarWndProc(0, WM_UAHDRAWMENUITEM, 0, 1, &MenuBarTheme{}); handled || result != 0 {
		t.Errorf("a panic was not contained: handled=%v result=%v", handled, result)
	}
}

// memoryDC creates a memory DC with the raw API call, skipping the test when
// none is available (the CreateCompatibleDC wrapper panics on failure).
func memoryDC(t *testing.T) HDC {
	t.Helper()
	ret, _, _ := procCreateCompatibleDC.Call(0)
	if ret == 0 {
		t.Skip("no memory DC available")
	}
	return HDC(ret)
}

// The Try helpers report failure on an unusable DC instead of panicking.
func TestTrySetBkModeReportsFailure(t *testing.T) {
	if TrySetBkMode(0, TRANSPARENT) {
		t.Error("TrySetBkMode on a nil DC reported success")
	}
	if TrySetBkMode(^HDC(0), TRANSPARENT) {
		t.Error("TrySetBkMode on an all-ones handle reported success")
	}
	dc := memoryDC(t)
	defer DeleteDC(dc)
	if !TrySetBkMode(dc, TRANSPARENT) {
		t.Error("TrySetBkMode on a memory DC failed")
	}
}
