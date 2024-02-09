//go:build darwin && production && !devtools

package application

func (w *macosWebviewWindow) enableDevTools() {}
