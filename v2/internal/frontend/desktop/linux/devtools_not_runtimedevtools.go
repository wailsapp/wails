//go:build !runtimedevtools

package linux

func (f *Frontend) OpenDevTools() {
	// Runtime devtools not enabled - this method does nothing
}