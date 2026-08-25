//go:build windows

package w32_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/wailsapp/wails/v3/pkg/w32"
)

// TestGetStockObjectReturnsValidHandle is a regression test for GetStockObject
// in gdi32.go.
//
// GetStockObject called procGetDeviceCaps instead of procGetStockObject, even
// though procGetStockObject was declared alongside it. GetDeviceCaps expects
// (HDC, index); handing it a stock-object constant as the HDC returns 0, so
// GetStockObject returned 0 for every input. It had no callers in v3, which is
// why the bug survived.
//
// This test fails on the pre-fix code, where every call returns 0. It runs on
// Windows CI via the //go:build windows constraint.
func TestGetStockObjectReturnsValidHandle(t *testing.T) {
	i := is.New(t)

	for _, tc := range []struct {
		name   string
		object int
	}{
		{"DEFAULT_GUI_FONT", w32.DEFAULT_GUI_FONT},
		{"SYSTEM_FONT", w32.SYSTEM_FONT},
		{"BLACK_BRUSH", w32.BLACK_BRUSH},
		{"WHITE_BRUSH", w32.WHITE_BRUSH},
	} {
		t.Run(tc.name, func(t *testing.T) {
			i := is.New(t)
			// Stock objects are always available; a 0 handle means the wrapper
			// called the wrong export.
			i.True(w32.GetStockObject(tc.object) != 0)
		})
	}

	// The handle must also be usable, not merely non-zero. Select it into a
	// memory DC - not the screen DC - so the test has no side effects on the
	// desktop, and check a previous object comes back.
	hdc := w32.CreateCompatibleDC(0)
	i.True(hdc != 0)
	defer w32.DeleteDC(hdc)

	font := w32.GetStockObject(w32.DEFAULT_GUI_FONT)
	i.True(font != 0)

	previous := w32.SelectObject(hdc, font)
	i.True(previous != 0)
	restored := w32.SelectObject(hdc, previous)
	i.Equal(restored, font)
}
