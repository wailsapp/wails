//go:build windows
// +build windows

package edge

import "math"

// appendDoubleArg appends the word(s) that pass a by-value C double argument
// to a COM method on this architecture, returning ok=false when the
// architecture cannot express the call.
//
// amd64: one word holding the bit pattern. The Go runtime mirrors the first
// four call words into XMM0-3 (asm_windows_amd64.s), which is where the
// Windows x64 ABI makes the callee read a double argument.
func appendDoubleArg(args []uintptr, v float64) ([]uintptr, bool) {
	return append(args, uintptr(math.Float64bits(v))), true
}
