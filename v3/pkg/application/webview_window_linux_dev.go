//go:build linux && !production && !android

package application

func (w *linuxWebviewWindow) openDevTools() {
	openDevTools(w.webview)
}

func (w *linuxWebviewWindow) enableDevTools() {
	enableDevTools(w.webview)
}
