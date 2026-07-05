//go:build windows
// +build windows

package edge

import (
	"log"

	"github.com/wailsapp/wails/v3/internal/webview2/internal/w32"
)

func (e *Chromium) SetSize(bounds w32.Rect) {
	if e.controller == nil {
		return
	}

	err := e.controller.PutBounds(bounds)
	if err != nil {
		// PutBounds can fail transiently while the browser process is
		// reconfiguring — e.g. RESOURCE_NOT_IN_CORRECT_STATE during a DPI
		// transition or after restoring from a minimised state
		// (wailsapp/wails#5544). A dropped resize is recoverable (the next
		// WM_SIZE will re-assert bounds); killing the process is not.
		log.Printf("[WebView2] SetSize failed: %v", err)
	}
}
