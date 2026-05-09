//go:build linux && production && !devtools && !android && !server && !runtimedevtools

package application

func (w *linuxWebviewWindow) openDevTools() {}

func (w *linuxWebviewWindow) enableDevTools() {}
