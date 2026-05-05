//go:build windows
// +build windows

package edge

import (
	"github.com/wailsapp/wails/webview2/internal/w32"
)

func (e *Chromium) SetSize(bounds w32.Rect) {
	if e.controller == nil {
		return
	}

	err := e.controller.PutBounds(bounds)
	if err != nil {
		e.errorCallback(err)
	}
}
