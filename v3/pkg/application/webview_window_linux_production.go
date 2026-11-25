//go:build linux && production && !devtools && !android && !openharmony && !server

package application

func (w *linuxWebviewWindow) openDevTools() {}

func (w *linuxWebviewWindow) enableDevTools() {}
