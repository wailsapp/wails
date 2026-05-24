//go:build windows && (!production || devtools)

package application

import "github.com/wailsapp/wails/webview2/pkg/edge"

func (w *windowsWebviewWindow) openDevTools() {
	w.chromium.OpenDevToolsWindow()
}

func (w *windowsWebviewWindow) enableDevTools(settings *edge.ICoreWebViewSettings) {
	err := settings.PutAreDevToolsEnabled(true)
	if err != nil {
		globalApplication.handleFatalError(err)
	}
}
