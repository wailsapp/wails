//go:build windows && runtimedevtools

package application

func (w *windowsWebviewWindow) openDevTools() {
	w.chromium.OpenDevToolsWindow()
}

func (w *windowsWebviewWindow) enableDevTools(settings *edge.ICoreWebViewSettings) {
	err := settings.PutAreDevToolsEnabled(true)
	if err != nil {
		globalApplication.handleFatalError(err)
	}
}
