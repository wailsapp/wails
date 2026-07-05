//go:build windows
// +build windows

package webview2

// appendDoubleArg appends the word(s) that pass a by-value C double argument
// to a COM method on this architecture, returning ok=false when the
// architecture cannot express the call.
//
// arm64: the Windows ARM64 ABI passes doubles in d0-d7, which Go's syscall
// path cannot populate (golang.org/issue/62583 — asmstdcall loads R0-R7
// only). There is no correct way to make this call from Go, so report
// ok=false and let callers no-op instead of passing garbage bits through an
// integer register (a garbage rasterization scale is how mis-scaled or blank
// content of the #5732 class arises).
func appendDoubleArg(args []uintptr, _ float64) ([]uintptr, bool) {
	return args, false
}
