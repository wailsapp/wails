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
// amd64: the double travels as its bit pattern in the next call word; the
// runtime mirrors the first four words into XMM0-3 (asm_windows_amd64.s),
// which is where the Windows x64 ABI makes the callee read it.
//
//go:uintptrescapes
func Call(fn uintptr, v float64, args ...uintptr) uintptr {
	hr, _, _ := syscall.SyscallN(fn, append(args, uintptr(math.Float64bits(v)))...)
	return hr
}
