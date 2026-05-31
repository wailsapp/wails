//go:build darwin && production && !devtools && !server

package application

func (w *macosWebviewWindow) enableDevTools() {}
func (w *macosWebviewWindow) openDevTools()   {}
