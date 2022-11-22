//go:build windows
// +build windows

package edge

import (
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/go-webview2/internal/w32"
	"unsafe"
)

func (e *Chromium) SetSize(bounds w32.Rect) {
	if e.controller == nil {
		return
	}

	e.controller.vtbl.PutBounds.Call(
		uintptr(unsafe.Pointer(e.controller)),
		uintptr(unsafe.Pointer(&bounds)),
	)
}
