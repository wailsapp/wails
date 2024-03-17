//go:build linux && production && !devtools

package application

func (w *linuxWebviewWindow) openDevTools() {}

func (w *linuxWebviewWindow) enableDevTools() {}
