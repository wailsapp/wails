//go:build windows && production && !devtools

package application

import "github.com/wailsapp/go-webview2/pkg/edge"

func (w *windowsWebviewWindow) openDevTools() {}

func (w *windowsWebviewWindow) enableDevTools(settings *edge.ICoreWebViewSettings) {
	err := settings.PutAreDevToolsEnabled(false)
	if err != nil {
		globalApplication.handleFatalError(err)
	}
}
