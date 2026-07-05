//go:build windows
// +build windows

package webview2

import "math"

// appendDoubleArg appends the word(s) that pass a by-value C double argument
// to a COM method on this architecture, returning ok=false when the
// architecture cannot express the call.
//
// 386: stdcall passes a double as two little-endian 32-bit words on the
// stack. A single truncated word would corrupt both the value and the
// callee's stack cleanup (stdcall callees pop their own arguments).
func appendDoubleArg(args []uintptr, v float64) ([]uintptr, bool) {
	bits := math.Float64bits(v)
	return append(args, uintptr(uint32(bits)), uintptr(uint32(bits>>32))), true
}
