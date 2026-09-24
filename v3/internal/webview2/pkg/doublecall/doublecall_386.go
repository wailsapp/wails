//go:build windows

package doublecall

import (
	"math"
	"syscall"
)

// Call invokes HRESULT fn(args..., double v) and returns the HRESULT. args
// are the integer words that precede the double in the C signature — the
// interface pointer, plus a by-value aggregate laid out by the caller.
//
// 386: stdcall pushes the double as two words, low half first.
//
//go:uintptrescapes
func Call(fn uintptr, v float64, args ...uintptr) uintptr {
	bits := math.Float64bits(v)
	hr, _, _ := syscall.SyscallN(fn, append(args, uintptr(uint32(bits)), uintptr(uint32(bits>>32)))...)
	return hr
}
