//go:build !darwin || ios || server

package application

import (
	"fmt"
	"unsafe"
)

type unsupportedNativeWindow struct{}

var nativeWindowClosed = make(chan uint)

func newNativeWindowImpl(*NativeWindow) nativeWindowImpl { return &unsupportedNativeWindow{} }
func (*unsupportedNativeWindow) run() error {
	return fmt.Errorf("NativeWindow is currently supported only on macOS")
}
func (*unsupportedNativeWindow) show()                        {}
func (*unsupportedNativeWindow) hide()                        {}
func (*unsupportedNativeWindow) focus()                       {}
func (*unsupportedNativeWindow) close()                       {}
func (*unsupportedNativeWindow) isVisible() bool              { return false }
func (*unsupportedNativeWindow) setTitle(string)              {}
func (*unsupportedNativeWindow) nativeWindow() unsafe.Pointer { return nil }
func (*unsupportedNativeWindow) setToolbar(*MacToolbar) error { return nil }
func (*unsupportedNativeWindow) installSplitView() error      { return nil }
