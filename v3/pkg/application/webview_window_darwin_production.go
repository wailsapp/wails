//go:build darwin && !ios && production && !devtools && !server && !purego

package application

func (w *macosWebviewWindow) enableDevTools() {}
func (w *macosWebviewWindow) openDevTools()   {}
