//go:build runtimedevtools

package windows

func (f *Frontend) OpenDevTools() {
	f.chromium.OpenDevToolsWindow()
}