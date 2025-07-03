//go:build !runtimedevtools

package windows

func (f *Frontend) OpenDevTools() {
	// Runtime devtools not enabled - this method does nothing
}