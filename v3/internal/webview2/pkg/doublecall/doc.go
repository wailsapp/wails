//go:build windows

// Package doublecall calls WebView2 COM methods that take a C double BY VALUE
// (put_ZoomFactor, put_RasterizationScale, put_Expires, the print settings).
// syscall.SyscallN passes every argument as an integer word, which is only
// right for such a double on amd64 (the runtime mirrors the call words into
// XMM0-3) and 386 (stdcall pushes both halves); on arm64 the ABI wants it in
// d0, which no syscall path can load (golang.org/issue/62583), so Call goes
// through a runtime.cgocall trampoline there.
package doublecall
