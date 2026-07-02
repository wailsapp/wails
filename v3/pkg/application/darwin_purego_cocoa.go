//go:build darwin && purego && !ios && !server

// Package application - CGO-free macOS backend.
//
// This file provides the foundational Objective-C runtime + Cocoa helper layer
// used by every other `darwin && purego` file. Instead of compiling Objective-C
// through cgo, we drive the Objective-C runtime directly via
// github.com/ebitengine/purego and its objc subpackage. This keeps the macOS
// backend buildable with CGO_ENABLED=0.
//
// The public surface deliberately mirrors the idioms used by the cgo backend so
// that the higher-level files read like their Objective-C counterparts:
//
//	obj := class("NSObject").send("alloc").send("init")
//	obj.send("release")
//	str := nsString("hello")
//	go2 := str.string()
package application

import (
	"reflect"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/objc"
)

// ---------------------------------------------------------------------------
// Foundation / AppKit framework loading
// ---------------------------------------------------------------------------

// The Objective-C runtime plus the frameworks we message into. dlopen keeps the
// symbols resident; the objc runtime itself is loaded by the objc package.
var (
	frameworksOnce sync.Once
)

const (
	frameworkFoundation = "/System/Library/Frameworks/Foundation.framework/Foundation"
	frameworkAppKit     = "/System/Library/Frameworks/AppKit.framework/AppKit"
	frameworkWebKit     = "/System/Library/Frameworks/WebKit.framework/WebKit"
	frameworkCoreGfx    = "/System/Library/Frameworks/CoreGraphics.framework/CoreGraphics"
	frameworkUniType    = "/System/Library/Frameworks/UniformTypeIdentifiers.framework/UniformTypeIdentifiers"
)

// loadFrameworks ensures the AppKit/WebKit class hierarchies are registered with
// the runtime. Referencing a class via objc_getClass only succeeds once the
// framework that defines it has been mapped into the process, so we dlopen them
// eagerly on first use.
func loadFrameworks() {
	frameworksOnce.Do(func() {
		for _, fw := range []string{
			frameworkFoundation,
			frameworkAppKit,
			frameworkWebKit,
			frameworkCoreGfx,
		} {
			// RTLD_GLOBAL so the class symbols become visible process-wide.
			// A failure here is fatal: every subsequent class() lookup would
			// return nil and all sends would become silent no-ops (no window,
			// no diagnostic).
			if _, err := purego.Dlopen(fw, purego.RTLD_NOW|purego.RTLD_GLOBAL); err != nil {
				panic("wails/purego: failed to load " + fw + ": " + err.Error())
			}
		}
		// UniformTypeIdentifiers only exists on macOS 11+; callers guard UTType
		// usage with classExists, so tolerate a failed load.
		_, _ = purego.Dlopen(frameworkUniType, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	})
}

// ---------------------------------------------------------------------------
// Thin id / class wrappers
// ---------------------------------------------------------------------------

// id is a convenience wrapper around objc.ID that lets the higher-level code
// send messages fluently. The zero value is a nil object.
type id objc.ID

func (o id) raw() objc.ID      { return objc.ID(o) }
func (o id) isNil() bool       { return uintptr(o) == 0 }
func (o id) ptr() uintptr      { return uintptr(o) }
func (o id) class() objc.Class { return objc.ID(o).Class() }

// send dispatches a selector (given by name) with the supplied arguments and
// returns the result as an id. Selector lookups are cached, so passing the name
// as a string here is cheap after the first call.
func (o id) send(sel string, args ...any) id {
	return id(objc.ID(o).Send(sel_(sel), args...))
}

// sendSuper dispatches to the superclass implementation.
//
// CAVEAT: purego resolves "super" from the receiver's DYNAMIC class. If an
// instance of one of our registered classes is ever isa-swizzled (e.g. KVO's
// NSKVONotifying_* subclassing), super-dispatch would resolve to the class
// itself and recurse forever. Do not KVO-observe instances of runtime
// registered delegate classes.
func (o id) sendSuper(sel string, args ...any) id {
	return id(objc.ID(o).SendSuper(sel_(sel), args...))
}

// get sends a selector and returns a typed result (BOOL, integers, floats,
// structs such as CGRect, ...).
func get[T any](o id, sel string, args ...any) T {
	return objc.Send[T](objc.ID(o), sel_(sel), args...)
}

// class looks up a registered Objective-C class by name and returns it as an id
// so class methods can be sent to it directly.
func class(name string) id {
	loadFrameworks()
	return id(objc.ID(objc.GetClass(name)))
}

// ---------------------------------------------------------------------------
// Selector cache
// ---------------------------------------------------------------------------

var (
	selMu    sync.RWMutex
	selCache = map[string]objc.SEL{}
)

// sel_ resolves (and caches) a selector by name. RegisterName grabs the global
// Objective-C lock, so caching matters on hot paths.
func sel_(name string) objc.SEL {
	selMu.RLock()
	s, ok := selCache[name]
	selMu.RUnlock()
	if ok {
		return s
	}
	selMu.Lock()
	defer selMu.Unlock()
	if s, ok = selCache[name]; ok {
		return s
	}
	s = objc.RegisterName(name)
	selCache[name] = s
	return s
}

// ---------------------------------------------------------------------------
// Core Graphics / Foundation geometry types
//
// These map to the C structs by memory layout so objc.Send can return them by
// value across the message-send boundary. On 64-bit macOS CGFloat is float64.
// ---------------------------------------------------------------------------

type CGFloat = float64

type CGPoint struct{ X, Y CGFloat }
type CGSize struct{ Width, Height CGFloat }
type CGRect struct {
	Origin CGPoint
	Size   CGSize
}

// NSPoint/NSSize/NSRect are typedefs of the CG equivalents on modern macOS.
type (
	NSPoint = CGPoint
	NSSize  = CGSize
	NSRect  = CGRect
)

func rect(x, y, w, h CGFloat) CGRect {
	return CGRect{Origin: CGPoint{X: x, Y: y}, Size: CGSize{Width: w, Height: h}}
}

// NSRange mirrors the Foundation struct (two NSUInteger fields).
type NSRange struct {
	Location uint
	Length   uint
}

// NSEdgeInsets mirrors the AppKit struct.
type NSEdgeInsets struct {
	Top, Left, Bottom, Right CGFloat
}

// ---------------------------------------------------------------------------
// NSString helpers
// ---------------------------------------------------------------------------

// nsString creates an autoreleased NSString from a Go string.
func nsString(s string) id {
	return class("NSString").send("stringWithUTF8String:", s)
}

// goString converts an NSString id back to a Go string.
func (o id) string() string {
	if o.isNil() {
		return ""
	}
	cstr := get[uintptr](o, "UTF8String")
	if cstr == 0 {
		return ""
	}
	return goStringFromC(cstr)
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

// ---------------------------------------------------------------------------
// NSURL / NSData helpers
// ---------------------------------------------------------------------------

func nsURL(s string) id {
	return class("NSURL").send("URLWithString:", nsString(s))
}

func fileURL(path string) id {
	return class("NSURL").send("fileURLWithPath:", nsString(path))
}

// nsData wraps a Go byte slice in an NSData that copies the bytes (safe to use
// after the slice is collected).
func nsData(b []byte) id {
	if len(b) == 0 {
		return class("NSData").send("data")
	}
	return class("NSData").send("dataWithBytes:length:", unsafe.Pointer(&b[0]), uint(len(b)))
}

// ---------------------------------------------------------------------------
// Autorelease pool
// ---------------------------------------------------------------------------

// withAutoreleasePool runs fn inside a fresh NSAutoreleasePool and drains it
// afterwards. Handy around bursts of autoreleased object creation off the main
// runloop.
func withAutoreleasePool(fn func()) {
	pool := class("NSAutoreleasePool").send("alloc").send("init")
	defer pool.send("drain")
	fn()
}

// ---------------------------------------------------------------------------
// BOOL / numeric helpers
// ---------------------------------------------------------------------------

// Objective-C BOOL is a signed char; purego marshals Go bool correctly, but we
// provide helpers for the NSNumber boxing used by KVC-style setters.
func nsNumberBool(b bool) id {
	return class("NSNumber").send("numberWithBool:", b)
}

func nsNumberInt(i int) id {
	return class("NSNumber").send("numberWithInteger:", i)
}

// ---------------------------------------------------------------------------
// Runtime class registration helpers
// ---------------------------------------------------------------------------

// registerDelegateClass is a thin wrapper over objc.RegisterClass that returns
// the created class as an id (so `new`/`alloc` can be sent to it) and panics on
// failure — a failed class registration is a programming error, not a runtime
// condition, so failing fast surfaces it during development.
func registerDelegateClass(name string, super string, ivars []objc.FieldDef, methods []objc.MethodDef) id {
	loadFrameworks()
	sup := objc.GetClass(super)
	if sup == 0 {
		// objc_allocateClassPair with a Nil superclass "succeeds" by creating
		// a methodless ROOT class, deferring the failure to an uncatchable
		// doesNotRecognizeSelector on the first alloc — fail fast instead.
		panic("wails/purego: superclass not found registering " + name + ": " + super)
	}
	cls, err := objc.RegisterClass(name, sup, nil, ivars, methods)
	if err != nil {
		panic("wails/purego: failed to register class " + name + ": " + err.Error())
	}
	return id(objc.ID(cls))
}

// ptrField is the reflect.Type used for uintptr-sized ivars that hold Go handle
// values (window ids, indices) bridged into Objective-C instances.
var ptrField = reflect.TypeOf(uintptr(0))

// invokeForeignBlock calls an Objective-C block received FROM a framework
// (e.g. a WebKit completion handler). objc.Block.Invoke only works for blocks
// created by Go via objc.NewBlock, so we call the block's invoke function
// pointer directly; it lives at offset 16 of the block literal on 64-bit
// (isa 8 bytes, flags 4, reserved 4). The block itself is the implicit first
// argument.
func invokeForeignBlock(block objc.ID, args ...uintptr) {
	if block == 0 {
		return
	}
	invoke := *(*uintptr)(unsafe.Pointer(uintptr(block) + 16))
	purego.SyscallN(invoke, append([]uintptr{uintptr(block)}, args...)...)
}

// ---------------------------------------------------------------------------
// Version / capability guards
//
// A purego build never compiles against a macOS SDK, so the two Objective-C
// guard mechanisms disappear and are replaced by runtime checks:
//
//   #if MAC_OS_X_VERSION_MAX_ALLOWED >= N   -> gone entirely (no SDK ceiling;
//                                              symbols are resolved by name at
//                                              runtime, i.e. always weak-linked).
//   if (@available(macOS X, *))             -> a runtime check, below.
//
// Prefer FEATURE detection over VERSION detection:
//   - respondsTo(obj, "setFoo:")   — the method/property exists on this build
//   - classExists("NSGlassEffectView") — the class exists on this OS
// Fall back to macOSAtLeast() only when a selector exists across versions but
// its behaviour or a constant's meaning changed.
//
// IMPORTANT: messaging a selector that does not exist raises an NSException
// (doesNotRecognizeSelector:) which CANNOT be recovered in pure Go — so every
// version-gated call MUST be guarded first; there is no @try/@catch safety net.
// ---------------------------------------------------------------------------

// nsOperatingSystemVersion mirrors Foundation's NSOperatingSystemVersion
// (three NSInteger fields).
type nsOperatingSystemVersion struct {
	Major int
	Minor int
	Patch int
}

var (
	osVersionOnce  sync.Once
	osVersionValue nsOperatingSystemVersion
)

func osVersion() nsOperatingSystemVersion {
	osVersionOnce.Do(func() {
		pi := class("NSProcessInfo").send("processInfo")
		osVersionValue = get[nsOperatingSystemVersion](pi, "operatingSystemVersion")
	})
	return osVersionValue
}

// macOSAtLeast reports whether the running OS is at least major.minor. Use only
// when respondsTo/classExists feature detection is not possible.
func macOSAtLeast(major, minor int) bool {
	v := osVersion()
	if v.Major != major {
		return v.Major > major
	}
	return v.Minor >= minor
}
