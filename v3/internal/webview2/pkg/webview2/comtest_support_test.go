//go:build windows

// comtest_support_test.go is the hand-written harness behind the generated
// *_gen_test.go files. It builds fake COM objects whose vtbl slots are
// recording trampolines, so generated wrappers can be called without a real
// WebView2 runtime and their raw marshalled words asserted.
//
// The expected encodings implemented here (f64w, u64w, rect handling, ...)
// are written from the Windows calling conventions directly — independently
// of the generator's emission logic — so a marshalling bug in the generator
// fails the generated tests instead of being mirrored into them.
//
// Windows limits a process to ~2000 syscall.NewCallback registrations, so a
// fixed set of shared trampolines (one per arity) records into a package
// global; tests in this package must not use t.Parallel().

package webview2

import (
	"fmt"
	"math"
	"runtime"
	"testing"
	"unsafe"

	"golang.org/x/sys/windows"
)

// wordExpect is one expected argument word. any==true skips the comparison
// (used for pointer-valued words whose addresses are not predictable).
type wordExpect struct {
	v   uintptr
	any bool
}

func xw(v uintptr) wordExpect { return wordExpect{v: v} }
func anyw() wordExpect        { return wordExpect{any: true} }

func anyws(n int) []wordExpect {
	out := make([]wordExpect, n)
	for i := range out {
		out[i] = anyw()
	}
	return out
}

func ptrw(p unsafe.Pointer) wordExpect { return wordExpect{v: uintptr(p)} }

// f64w returns the expected encoding of an [in] double: the IEEE-754 bit
// pattern in one register word on 64-bit, two 4-byte stack words (low first)
// on 386.
func f64w(v float64) []wordExpect {
	bits := math.Float64bits(v)
	if runtime.GOARCH == "386" {
		return []wordExpect{xw(uintptr(uint32(bits))), xw(uintptr(uint32(bits >> 32)))}
	}
	return []wordExpect{xw(uintptr(bits))}
}

func f32w(v float32) wordExpect { return xw(uintptr(math.Float32bits(v))) }

// u64w returns the expected encoding of an 8-byte integer or aggregate.
func u64w(v uint64) []wordExpect {
	if runtime.GOARCH == "386" {
		return []wordExpect{xw(uintptr(uint32(v))), xw(uintptr(uint32(v >> 32)))}
	}
	return []wordExpect{xw(uintptr(v))}
}

// rectWordCount is the number of argument words a by-value RECT occupies:
// one pointer-to-copy on amd64, a register pair on arm64, four stack words
// on 386.
func rectWordCount() int {
	switch runtime.GOARCH {
	case "amd64":
		return 1
	case "386":
		return 4
	default: // arm64
		return 2
	}
}

func makeToken(bits int64) EventRegistrationToken {
	return EventRegistrationToken{value: bits}
}

func tokenBits(t EventRegistrationToken) uint64 { return uint64(t.value) }

// asPointer reconstructs a pointer from a recorded argument word. The memory
// behind it is owned by the in-flight fake call (wrapper locals, CoTaskMem,
// or test fixtures), so the conversion is valid; the indirection keeps go
// vet's unsafeptr heuristic from flagging these deliberate reconstructions.
func asPointer(word uintptr) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Pointer(&word))
}

var fakeTarget [64]byte

// fakeTargetPtr is a stable address tests use as the "object" a fake COM
// method hands back through an out-parameter.
func fakeTargetPtr() unsafe.Pointer { return unsafe.Pointer(&fakeTarget[0]) }

// recorder is the shared state behind the arity trampolines. Tests are
// sequential, so a single global is safe.
var rec struct {
	words    []uintptr
	hr       uintptr
	writers  map[int]func(p unsafe.Pointer)
	captures map[int]func(args []uintptr)
	strings  map[int]string
	guids    map[int]string
	rect     RECT
	hasRect  bool
}

func recReset() {
	rec.words = nil
	rec.hr = 0
	rec.writers = map[int]func(p unsafe.Pointer){}
	rec.captures = map[int]func(args []uintptr){}
	rec.strings = map[int]string{}
	rec.guids = map[int]string{}
	rec.rect = RECT{}
	rec.hasRect = false
}

func recSetHR(hr uintptr) { rec.hr = hr }

// recordCall is the single body behind every arity trampoline.
func recordCall(args []uintptr) uintptr {
	rec.words = append([]uintptr(nil), args...)
	for idx, capture := range rec.captures {
		if idx < len(args) {
			capture(args)
		}
	}
	for idx, write := range rec.writers {
		if idx < len(args) && args[idx] != 0 {
			write(asPointer(args[idx]))
		}
	}
	return rec.hr
}

// One shared trampoline per arity. Created once at init so the process-wide
// NewCallback budget is untouched by the number of tests.
var recProcs = map[int]ComProc{
	1: NewComProc(func(a0 uintptr) uintptr { return recordCall([]uintptr{a0}) }),
	2: NewComProc(func(a0, a1 uintptr) uintptr { return recordCall([]uintptr{a0, a1}) }),
	3: NewComProc(func(a0, a1, a2 uintptr) uintptr { return recordCall([]uintptr{a0, a1, a2}) }),
	4: NewComProc(func(a0, a1, a2, a3 uintptr) uintptr { return recordCall([]uintptr{a0, a1, a2, a3}) }),
	5: NewComProc(func(a0, a1, a2, a3, a4 uintptr) uintptr { return recordCall([]uintptr{a0, a1, a2, a3, a4}) }),
	6: NewComProc(func(a0, a1, a2, a3, a4, a5 uintptr) uintptr {
		return recordCall([]uintptr{a0, a1, a2, a3, a4, a5})
	}),
	7: NewComProc(func(a0, a1, a2, a3, a4, a5, a6 uintptr) uintptr {
		return recordCall([]uintptr{a0, a1, a2, a3, a4, a5, a6})
	}),
	8: NewComProc(func(a0, a1, a2, a3, a4, a5, a6, a7 uintptr) uintptr {
		return recordCall([]uintptr{a0, a1, a2, a3, a4, a5, a6, a7})
	}),
	9: NewComProc(func(a0, a1, a2, a3, a4, a5, a6, a7, a8 uintptr) uintptr {
		return recordCall([]uintptr{a0, a1, a2, a3, a4, a5, a6, a7, a8})
	}),
	10: NewComProc(func(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9 uintptr) uintptr {
		return recordCall([]uintptr{a0, a1, a2, a3, a4, a5, a6, a7, a8, a9})
	}),
	11: NewComProc(func(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10 uintptr) uintptr {
		return recordCall([]uintptr{a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10})
	}),
	12: NewComProc(func(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11 uintptr) uintptr {
		return recordCall([]uintptr{a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11})
	}),
	13: NewComProc(func(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12 uintptr) uintptr {
		return recordCall([]uintptr{a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12})
	}),
	14: NewComProc(func(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13 uintptr) uintptr {
		return recordCall([]uintptr{a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13})
	}),
	15: NewComProc(func(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14 uintptr) uintptr {
		return recordCall([]uintptr{a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14})
	}),
	16: NewComProc(func(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15 uintptr) uintptr {
		return recordCall([]uintptr{a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15})
	}),
}

// recProc returns the shared recording trampoline for the given arity
// (including the `this` pointer).
func recProc(arity int) ComProc {
	p, ok := recProcs[arity]
	if !ok {
		panic(fmt.Sprintf("comtest: no trampoline for arity %d — extend recProcs", arity))
	}
	return p
}

// ── Out-parameter writers ────────────────────────────────────────────────────

func recWrite8(argIdx int, v uint8) {
	rec.writers[argIdx] = func(p unsafe.Pointer) { *(*uint8)(p) = v }
}

func recWrite16(argIdx int, v uint16) {
	rec.writers[argIdx] = func(p unsafe.Pointer) { *(*uint16)(p) = v }
}

func recWrite32(argIdx int, v uint32) {
	rec.writers[argIdx] = func(p unsafe.Pointer) { *(*uint32)(p) = v }
}

func recWrite64(argIdx int, v uint64) {
	rec.writers[argIdx] = func(p unsafe.Pointer) { *(*uint64)(p) = v }
}

func recWritePtr(argIdx int, target unsafe.Pointer) {
	rec.writers[argIdx] = func(p unsafe.Pointer) { *(*unsafe.Pointer)(p) = target }
}

// recWriteUTF16 simulates COM returning a CoTaskMem-allocated string.
func recWriteUTF16(argIdx int, s string) {
	rec.writers[argIdx] = func(p unsafe.Pointer) {
		*(*unsafe.Pointer)(p) = coTaskMemUTF16(s)
	}
}

// recWritePtrArray simulates COM returning a CoTaskMem-allocated array of n
// interface pointers (all pointing at the fake target).
func recWritePtrArray(argIdx int, n int) {
	rec.writers[argIdx] = func(p unsafe.Pointer) {
		arr := coTaskMemAlloc(uintptr(n) * unsafe.Sizeof(uintptr(0)))
		slice := unsafe.Slice((*unsafe.Pointer)(arr), n)
		for i := range slice {
			slice[i] = fakeTargetPtr()
		}
		*(*unsafe.Pointer)(p) = arr
	}
}

// recWriteUTF16Array simulates COM returning a CoTaskMem array of CoTaskMem
// UTF-16 strings.
func recWriteUTF16Array(argIdx int, ss []string) {
	rec.writers[argIdx] = func(p unsafe.Pointer) {
		arr := coTaskMemAlloc(uintptr(len(ss)) * unsafe.Sizeof(uintptr(0)))
		slice := unsafe.Slice((*unsafe.Pointer)(arr), len(ss))
		for i, s := range ss {
			slice[i] = coTaskMemUTF16(s)
		}
		*(*unsafe.Pointer)(p) = arr
	}
}

var procCoTaskMemAlloc = windows.NewLazySystemDLL("ole32.dll").NewProc("CoTaskMemAlloc")

func coTaskMemAlloc(n uintptr) unsafe.Pointer {
	p, _, _ := procCoTaskMemAlloc.Call(n)
	if p == 0 {
		panic("CoTaskMemAlloc failed")
	}
	return asPointer(p)
}

func coTaskMemUTF16(s string) unsafe.Pointer {
	u, err := windows.UTF16FromString(s)
	if err != nil {
		panic(err)
	}
	p := coTaskMemAlloc(uintptr(len(u)) * 2)
	copy(unsafe.Slice((*uint16)(p), len(u)), u)
	return p
}

// ── Captures (snapshots taken during the fake call, while the pointed-at
//    memory is still live) ───────────────────────────────────────────────────

func recCaptureUTF16(argIdx int) {
	rec.captures[argIdx] = func(args []uintptr) {
		if args[argIdx] == 0 {
			rec.strings[argIdx] = "<nil>"
			return
		}
		rec.strings[argIdx] = windows.UTF16PtrToString((*uint16)(asPointer(args[argIdx])))
	}
}

func recCapturedString(argIdx int) string { return rec.strings[argIdx] }

func recCaptureGUID(argIdx int) {
	rec.captures[argIdx] = func(args []uintptr) {
		g := (*windows.GUID)(asPointer(args[argIdx]))
		rec.guids[argIdx] = fmt.Sprintf("%08x-%04x-%04x-%02x%02x-%02x%02x%02x%02x%02x%02x",
			g.Data1, g.Data2, g.Data3,
			g.Data4[0], g.Data4[1], g.Data4[2], g.Data4[3],
			g.Data4[4], g.Data4[5], g.Data4[6], g.Data4[7])
	}
}

func recCapturedGUID(argIdx int) string { return rec.guids[argIdx] }

// recCaptureRECT snapshots a by-value RECT argument starting at argIdx:
// dereferencing the pointer-to-copy on amd64, reassembling register/stack
// words on arm64/386.
func recCaptureRECT(argIdx int) {
	rec.captures[argIdx] = func(args []uintptr) {
		var r RECT
		switch runtime.GOARCH {
		case "amd64":
			r = *(*RECT)(asPointer(args[argIdx]))
		case "386":
			words := [4]uint32{uint32(args[argIdx]), uint32(args[argIdx+1]), uint32(args[argIdx+2]), uint32(args[argIdx+3])}
			r = *(*RECT)(unsafe.Pointer(&words))
		default: // arm64: two 8-byte register words
			words := [2]uintptr{args[argIdx], args[argIdx+1]}
			r = *(*RECT)(unsafe.Pointer(&words))
		}
		rec.rect = r
		rec.hasRect = true
	}
}

func recCapturedRECT() RECT { return rec.rect }

// ── Assertions ───────────────────────────────────────────────────────────────

// expectCall asserts the recorded call: this pointer first, then each
// expected word (anyw() entries are skipped).
func expectCall(t *testing.T, this unsafe.Pointer, want []wordExpect) {
	t.Helper()
	if len(rec.words) == 0 {
		t.Fatal("fake COM method was never called")
	}
	if rec.words[0] != uintptr(this) {
		t.Errorf("this pointer = %#x, want %#x", rec.words[0], uintptr(this))
	}
	if len(rec.words)-1 != len(want) {
		t.Fatalf("call carried %d argument words, want %d (%#v)", len(rec.words)-1, len(want), rec.words[1:])
	}
	for i, w := range want {
		if w.any {
			continue
		}
		if rec.words[1+i] != w.v {
			t.Errorf("argument word %d = %#x, want %#x", i, rec.words[1+i], w.v)
		}
	}
}

func utf16Ptr(t *testing.T, s string) *uint16 {
	t.Helper()
	p, err := windows.UTF16PtrFromString(s)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

// TestComProcCallErrorMapping pins the HRESULT error path once for the whole
// package: a failing fake must surface as syscall.Errno, a succeeding one as
// a nil error (never the always-non-nil GetLastError from SyscallN).
func TestComProcCallErrorMapping(t *testing.T) {
	recReset()
	obj := &ICoreWebView2{Vtbl: &ICoreWebView2Vtbl{}}
	obj.Vtbl.Stop = recProc(1)
	if err := obj.Stop(); err != nil {
		t.Errorf("S_OK call returned error %v, want nil", err)
	}
	recReset()
	recSetHR(0x80004005) // E_FAIL
	obj.Vtbl.Stop = recProc(1)
	err := obj.Stop()
	if err == nil {
		t.Fatal("E_FAIL call returned nil error")
	}
}
