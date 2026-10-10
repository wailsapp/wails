//go:build windows

package doublecall

import "unsafe"

// doubleCall is the argument block of the trampoline in doublecall_arm64.s;
// go_asm.h hands the assembly its field offsets.
type doubleCall struct {
	fn             uintptr
	a1, a2, a3, a4 uintptr
	v              float64
	hr             uintptr
}

// callDoubleABI0 is the ABI0 entry point of the trampoline, exported by
// doublecall_arm64.s the way purego exports its own.
var callDoubleABI0 uintptr

//go:linkname runtime_cgocall runtime.cgocall
func runtime_cgocall(fn uintptr, arg unsafe.Pointer) int32

// Call invokes HRESULT fn(args..., double v) and returns the HRESULT. args
// are the integer words that precede the double in the C signature — the
// interface pointer, plus a by-value aggregate laid out by the caller; at
// most four fit the trampoline.
//
// arm64: the Windows ARM64 ABI passes the double in d0, which
// syscall.SyscallN cannot load (golang.org/issue/62583: asmstdcall fills
// r0-r7 only). The call goes through runtime.cgocall instead — the same
// system-stack switch SyscallN itself uses — into a trampoline that loads
// args into r0-r3 and v into d0. The HRESULT comes back in w0 and the upper
// half of x0 is unspecified, hence the truncation.
//
//go:uintptrescapes
func Call(fn uintptr, v float64, args ...uintptr) uintptr {
	if len(args) > 4 {
		panic("doublecall: at most four integer words may precede the double")
	}
	var words [4]uintptr
	copy(words[:], args)
	call := doubleCall{fn: fn, a1: words[0], a2: words[1], a3: words[2], a4: words[3], v: v}
	runtime_cgocall(callDoubleABI0, unsafe.Pointer(&call))
	return uintptr(uint32(call.hr))
}
