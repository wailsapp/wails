//go:build linux && production && !devtools && !android && !server

package application

func (w *linuxWebviewWindow) openDevTools() {}

func (w *linuxWebviewWindow) enableDevTools() {}
