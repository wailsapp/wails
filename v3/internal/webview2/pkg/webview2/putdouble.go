//go:build windows
// +build windows

package webview2

import "errors"

// ErrDoubleArgUnsupported is returned by setters that take a C double BY
// VALUE on architectures where Go cannot marshal such an argument
// (windows/arm64: doubles go in d0-d7, which Go's syscall path cannot
// populate — golang.org/issue/62583). Callers that can tolerate the setting
// staying at its default should treat this error as non-fatal.
var ErrDoubleArgUnsupported = errors.New("COM by-value double arguments are not supported on windows/arm64 (golang.org/issue/62583)")
