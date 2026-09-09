//go:build windows
// +build windows

package webview2

import "unsafe"

// appendRectArg appends the word(s) that pass a by-value RECT argument to a
// COM method on this architecture.
//
// arm64: a 16-byte struct is passed by value in two registers.
func appendRectArg(args []uintptr, bounds *RECT) ([]uintptr, bool) {
	words := (*[2]uintptr)(unsafe.Pointer(bounds))
	return append(args, words[0], words[1]), true
}
