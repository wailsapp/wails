//go:build darwin && !ios && purego

// Package mac provides a set of functions to interact with the macOS platform.
//
// This is the CGO-free (purego) implementation. Instead of compiling the
// Objective-C snippet through cgo, we drive the Objective-C runtime directly
// via github.com/ebitengine/purego and its objc subpackage. This keeps the
// package buildable with CGO_ENABLED=0.
//
// It mirrors the exported surface of mac.go exactly.
package mac

import (
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/objc"
)

// frameworkFoundation is the path to the Foundation framework, which defines
// NSBundle and NSString.
const frameworkFoundation = "/System/Library/Frameworks/Foundation.framework/Foundation"

// loadFoundation dlopen's Foundation so that objc_getClass can resolve the
// NSBundle class. Referencing a class only succeeds once the framework that
// defines it has been mapped into the process, so we do it eagerly (and once).
var loadFoundationOnce sync.Once

func loadFoundation() {
	loadFoundationOnce.Do(func() {
		// RTLD_GLOBAL so the class symbols become visible process-wide.
		_, _ = purego.Dlopen(frameworkFoundation, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	})
}

// goStringFromC copies a NUL-terminated C string at the given address into a Go
// string without cgo.
func goStringFromC(p uintptr) string {
	if p == 0 {
		return ""
	}
	var n int
	for {
		b := *(*byte)(unsafe.Pointer(p + uintptr(n)))
		if b == 0 {
			break
		}
		n++
	}
	if n == 0 {
		return ""
	}
	buf := make([]byte, n)
	copy(buf, unsafe.Slice((*byte)(unsafe.Pointer(p)), n))
	return string(buf)
}

// GetBundleID returns the bundle ID of the application.
//
// Equivalent to `[[NSBundle mainBundle] bundleIdentifier]` converted to a Go
// string via -[NSString UTF8String].
func GetBundleID() string {
	loadFoundation()

	// Callers are arbitrary goroutines with no ambient autorelease pool; wrap
	// one around the UTF8String conversion so nothing autoreleased leaks.
	pool := objc.ID(objc.GetClass("NSAutoreleasePool")).
		Send(objc.RegisterName("alloc")).Send(objc.RegisterName("init"))
	defer pool.Send(objc.RegisterName("drain"))

	nsBundle := objc.ID(objc.GetClass("NSBundle"))
	mainBundle := nsBundle.Send(objc.RegisterName("mainBundle"))
	if mainBundle == 0 {
		return ""
	}

	bundleID := mainBundle.Send(objc.RegisterName("bundleIdentifier"))
	if bundleID == 0 {
		return ""
	}

	cstr := objc.Send[uintptr](bundleID, objc.RegisterName("UTF8String"))
	return goStringFromC(cstr)
}
