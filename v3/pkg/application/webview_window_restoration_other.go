//go:build !darwin || ios || server

package application

import "unsafe"

// Off macOS window state restoration is a documented no-op: the setters
// store their values, OnRestore is never called, and the interaction state
// methods return ErrMacOnly.

func macWindowRestorationSupported() bool { return false }

func macWindowRestorationConfigure(unsafe.Pointer, string)                {}
func macWindowRestorationSetData(unsafe.Pointer, string)                  {}
func macWindowRestorationComplete(uint64, unsafe.Pointer, string)         {}
func macWindowRestorationAwaitNative(Window) unsafe.Pointer               { return nil }
func macWindowRestorationInteractionState(unsafe.Pointer) ([]byte, error) { return nil, ErrMacOnly }
func macWindowRestorationSetInteractionState(unsafe.Pointer, []byte) error {
	return ErrMacOnly
}
