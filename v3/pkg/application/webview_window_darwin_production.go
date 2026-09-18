//go:build darwin && !ios && production && !devtools && !server

package application

func (w *macosWebviewWindow) enableDevTools() {}
func (w *macosWebviewWindow) openDevTools()   {}
