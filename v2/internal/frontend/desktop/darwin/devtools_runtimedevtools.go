//go:build darwin && runtimedevtools

package darwin

func (f *Frontend) OpenDevTools() {
	showInspector(f.mainWindow.context)
}