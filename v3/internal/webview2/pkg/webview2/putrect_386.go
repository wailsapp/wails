//go:build windows
// +build windows

package webview2

import "unsafe"

// appendRectArg appends the word(s) that pass a by-value RECT argument to a
// COM method on this architecture.
//
// 386: stdcall passes the 16-byte struct as four stack words.
func appendRectArg(args []uintptr, bounds *RECT) ([]uintptr, bool) {
	words := (*[4]uintptr)(unsafe.Pointer(bounds))
	return append(args, words[0], words[1], words[2], words[3]), true
}
