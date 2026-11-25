//go:build linux && (!production || devtools) && !android && !server

package application

func (w *linuxWebviewWindow) openDevTools() {
	openDevTools(w.webview)
}

func (w *linuxWebviewWindow) enableDevTools() {
	enableDevTools(w.webview)
}
