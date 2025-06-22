//go:build darwin && !runtimedevtools

package darwin

func (f *Frontend) OpenDevTools() {
	// Runtime devtools not enabled - this method does nothing
}