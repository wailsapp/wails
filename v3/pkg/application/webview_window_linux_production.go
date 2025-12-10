//go:build linux && production && !devtools && !android

package application

func (w *linuxWebviewWindow) openDevTools() {}

func (w *linuxWebviewWindow) enableDevTools() {}
