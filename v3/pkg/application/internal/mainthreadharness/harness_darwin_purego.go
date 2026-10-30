//go:build darwin && purego && !ios && !server

package mainthreadharness

// CGO-free counterpart of harness_darwin.go. It exposes the same API and runs
// the same loop: the C helpers become purego bindings to AppKit and
// CoreFoundation, and the cgo malloc'd atomic run flag becomes an atomic.Bool.

import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/ebitengine/purego"
)

const (
	frameworkAppKit         = "/System/Library/Frameworks/AppKit.framework/AppKit"
	frameworkCoreFoundation = "/System/Library/Frameworks/CoreFoundation.framework/CoreFoundation"
)

var (
	harnessOnce sync.Once

	// NSApplicationLoad() BOOL
	nsApplicationLoad func() bool
	// CFRunLoopRunInMode(CFStringRef mode, CFTimeInterval seconds, Boolean returnAfterSourceHandled) SInt32
	cfRunLoopRunInModeFn func(mode uintptr, seconds float64, returnAfterSourceHandled bool) int32

	// Address of the NSApp global (NSApplication *).
	nsAppAddr uintptr
	// Value of the kCFRunLoopDefaultMode global (CFStringRef).
	defaultRunLoopMode uintptr
)

// load binds the AppKit and CoreFoundation symbols the harness uses. AppKit is
// not linked into a purego test binary, so it is dlopen'd here.
func load() {
	harnessOnce.Do(func() {
		appKit, err := purego.Dlopen(frameworkAppKit, purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if err != nil {
			panic("mainthreadharness: failed to load AppKit: " + err.Error())
		}
		cf, err := purego.Dlopen(frameworkCoreFoundation, purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if err != nil {
			panic("mainthreadharness: failed to load CoreFoundation: " + err.Error())
		}

		purego.RegisterLibFunc(&nsApplicationLoad, appKit, "NSApplicationLoad")
		purego.RegisterLibFunc(&cfRunLoopRunInModeFn, cf, "CFRunLoopRunInMode")

		nsAppSym, err := purego.Dlsym(appKit, "NSApp")
		if err != nil {
			panic("mainthreadharness: failed to resolve NSApp: " + err.Error())
		}
		nsAppAddr = nsAppSym

		modeSym, err := purego.Dlsym(cf, "kCFRunLoopDefaultMode")
		if err != nil {
			panic("mainthreadharness: failed to resolve kCFRunLoopDefaultMode: " + err.Error())
		}
		defaultRunLoopMode = *(*uintptr)(unsafe.Pointer(modeSym))
	})
}

// InitAppKit brings AppKit up and reports whether it started. It matters for
// the dispatch tests: until AppKit is loaded the main dispatch queue is drained
// by CFRunLoop, which re-enters and delivers queued blocks inside a nested
// loop. Once AppKit owns the draining it does not re-enter, which is the
// behaviour a real application has.
func InitAppKit() bool {
	load()
	loaded := nsApplicationLoad()
	return loaded && *(*uintptr)(unsafe.Pointer(nsAppAddr)) != 0
}

// Loop turns a run loop until it is stopped.
type Loop struct {
	running atomic.Bool
}

// NewLoop returns a loop which is already in the running state, so that Stop
// may be called before Run.
func NewLoop() *Loop {
	l := &Loop{}
	l.running.Store(true)
	return l
}

// Run turns the run loop of the calling thread until Stop is called. It returns
// straight away if Stop was called first.
func (l *Loop) Run() {
	load()
	for l.running.Load() {
		cfRunLoopRunInModeFn(defaultRunLoopMode, 0.05, false)
	}
}

// Stop ends the loop. It is safe to call from any thread, and before Run.
func (l *Loop) Stop() { l.running.Store(false) }

// Free releases the loop. The purego loop holds no native resources, so this
// does nothing; it is kept so the API matches the cgo harness.
func (l *Loop) Free() {}

// RunNested runs a nested run loop for d without returning, which is what
// [NSAlert runModal] and menu tracking do while they are on screen.
func RunNested(d time.Duration) {
	load()
	cfRunLoopRunInModeFn(defaultRunLoopMode, d.Seconds(), false)
}
