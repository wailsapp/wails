//go:build darwin && !ios && production && !devtools && !server && !wails_native

package application

func (w *macosWebviewWindow) enableDevTools() {}
func (w *macosWebviewWindow) openDevTools()   {}
