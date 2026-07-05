//go:build windows
// +build windows

package webview2

import "unsafe"

// appendRectArg appends the word(s) that pass a by-value RECT argument to a
// COM method on this architecture.
//
// amd64: structs larger than 8 bytes are passed by reference.
func appendRectArg(args []uintptr, bounds *RECT) ([]uintptr, bool) {
	return append(args, uintptr(unsafe.Pointer(bounds))), true
}
