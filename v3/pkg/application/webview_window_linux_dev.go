//go:build linux && (!production || devtools) && !android && !openharmony && !server

package application

func (w *linuxWebviewWindow) openDevTools() {
	openDevTools(w.webview)
}

func (w *linuxWebviewWindow) enableDevTools() {
	enableDevTools(w.webview)
}
