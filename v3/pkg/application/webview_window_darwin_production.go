//go:build darwin && production && !devtools && !server && !runtimedevtools

package application

func (w *macosWebviewWindow) enableDevTools() {}
func (w *macosWebviewWindow) openDevTools()   {}
