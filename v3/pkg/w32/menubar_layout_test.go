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
	if IsValidDC(HDC(0xffffffff94011e00)) {
		t.Error("IsValidDC(garbage) = true")
	}
	if TrySetTextColor(0, 0xdedede) {
		t.Error("TrySetTextColor on a nil DC reported success")
	}
	dc := CreateCompatibleDC(0)
	if dc == 0 {
		t.Skip("no memory DC available")
	}
	defer DeleteDC(dc)
	if !IsValidDC(dc) {
		t.Error("IsValidDC(memory DC) = false")
	}
	if !TrySetTextColor(dc, 0xdedede) {
		t.Error("TrySetTextColor on a memory DC failed")
	}
}
