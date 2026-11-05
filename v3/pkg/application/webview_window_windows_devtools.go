//go:build windows && !server && (!production || devtools)

package application

import "github.com/wailsapp/wails/v3/internal/webview2/pkg/edge"

func (w *windowsWebviewWindow) openDevTools() {
	w.chromium.OpenDevToolsWindow()
}

func (w *windowsWebviewWindow) enableDevTools(settings *edge.ICoreWebViewSettings) error {
	return settings.PutAreDevToolsEnabled(true)
}
