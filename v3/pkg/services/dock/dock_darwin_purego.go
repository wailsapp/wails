//go:build darwin && !ios && purego

package dock

// CGO-free macOS dock service implementation.
//
// This mirrors the behaviour of the cgo dock_darwin.go by driving NSApplication
// / NSDockTile directly through the Objective-C runtime (via
// github.com/ebitengine/purego) instead of compiling Objective-C. It keeps the
// dock service buildable with CGO_ENABLED=0.
//
//	HideAppIcon -> [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory]
//	ShowAppIcon -> [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular]
//	setBadge    -> require Regular policy, then
//	               [[NSApp dockTile] setBadgeLabel:label]; [[NSApp dockTile] display]
//
// All AppKit access is marshalled onto the main thread with
// dispatch_sync(dispatch_get_main_queue(), ...), exactly like the cgo version.

import (
	"context"
	"fmt"
	"sync"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/objc"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// NSApplicationActivationPolicy values (AppKit).
const (
	nsApplicationActivationPolicyRegular   = 0
	nsApplicationActivationPolicyAccessory = 1
)

// ---------------------------------------------------------------------------
// Minimal Objective-C runtime helpers (package-local; not shared across
// packages by design).
// ---------------------------------------------------------------------------

const (
	frameworkFoundation = "/System/Library/Frameworks/Foundation.framework/Foundation"
	frameworkAppKit     = "/System/Library/Frameworks/AppKit.framework/AppKit"
)

var frameworksOnce sync.Once

func loadFrameworks() {
	frameworksOnce.Do(func() {
		for _, fw := range []string{frameworkFoundation, frameworkAppKit} {
			_, _ = purego.Dlopen(fw, purego.RTLD_NOW|purego.RTLD_GLOBAL)
		}
	})
}

// id is a thin wrapper around objc.ID for fluent message sends.
type id objc.ID

func (o id) isNil() bool { return uintptr(o) == 0 }

func (o id) send(sel string, args ...any) id {
	return id(objc.ID(o).Send(sel_(sel), args...))
}

func get[T any](o id, sel string, args ...any) T {
	return objc.Send[T](objc.ID(o), sel_(sel), args...)
}

func class(name string) id {
	loadFrameworks()
	return id(objc.ID(objc.GetClass(name)))
}

var (
	selMu    sync.RWMutex
	selCache = map[string]objc.SEL{}
)

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

func nsString(s string) id {
	return class("NSString").send("stringWithUTF8String:", s)
}

// nsApp returns the shared NSApplication instance ([NSApplication sharedApplication]).
func nsApp() id {
	return class("NSApplication").send("sharedApplication")
}

// ---------------------------------------------------------------------------
// libdispatch: dispatch_sync(dispatch_get_main_queue(), block)
// ---------------------------------------------------------------------------

var (
	dispatchMainQueue uintptr
	dispatchSync      func(queue uintptr, block objc.Block)
	dispatchOnce      sync.Once
)

func initDispatch() {
	dispatchOnce.Do(func() {
		lib, err := purego.Dlopen("/usr/lib/libSystem.B.dylib", purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if err != nil {
			panic("wails/purego: failed to load libSystem: " + err.Error())
		}
		mainQ, err := purego.Dlsym(lib, "_dispatch_main_q")
		if err != nil {
			panic("wails/purego: failed to resolve _dispatch_main_q: " + err.Error())
		}
		dispatchMainQueue = mainQ
		purego.RegisterLibFunc(&dispatchSync, lib, "dispatch_sync")
	})
}

// onMain runs fn synchronously on the process main thread, matching the cgo
// dispatch_sync(dispatch_get_main_queue(), ^{ ... }) usage. If already on the
// main thread it runs fn directly to avoid a deadlock.
func onMain(fn func()) {
	initDispatch()
	if get[bool](class("NSThread"), "isMainThread") {
		fn()
		return
	}
	block := objc.NewBlock(func(objc.Block) {
		fn()
	})
	dispatchSync(dispatchMainQueue, block)
}

// ---------------------------------------------------------------------------
// Dock service
// ---------------------------------------------------------------------------

type darwinDock struct {
	mu    sync.RWMutex
	Badge *string
}

// New creates a new Dock Service.
func New() *DockService {
	return &DockService{
		impl: &darwinDock{
			Badge: nil,
		},
	}
}

// NewWithOptions creates a new dock service with badge options.
// Currently, options are not available on macOS and are ignored.
func NewWithOptions(options BadgeOptions) *DockService {
	return New()
}

func (d *darwinDock) Startup(ctx context.Context, options application.ServiceOptions) error {
	return nil
}

func (d *darwinDock) Shutdown() error {
	return nil
}

// HideAppIcon hides the app icon in the macOS Dock.
func (d *darwinDock) HideAppIcon() {
	onMain(func() {
		nsApp().send("setActivationPolicy:", nsApplicationActivationPolicyAccessory)
	})
}

// ShowAppIcon shows the app icon in the macOS Dock.
// Note: After showing the dock icon, you may need to call SetBadge again
// to reapply any previously set badge, as changing activation policies clears the badge.
func (d *darwinDock) ShowAppIcon() {
	onMain(func() {
		nsApp().send("setActivationPolicy:", nsApplicationActivationPolicyRegular)
	})
}

// setBadge handles the native call and updates the internal badge state with locking.
func (d *darwinDock) setBadge(label *string) error {
	var success bool
	onMain(func() {
		// Ensure the app is in Regular activation policy (dock icon visible).
		app := nsApp()
		if get[int](app, "activationPolicy") != nsApplicationActivationPolicyRegular {
			success = false
			return
		}

		var nsLabel id
		if label != nil {
			nsLabel = nsString(*label)
		}
		dockTile := app.send("dockTile")
		dockTile.send("setBadgeLabel:", nsLabel)
		dockTile.send("display")
		success = true
	})

	if !success {
		return fmt.Errorf("failed to set badge")
	}

	d.mu.Lock()
	d.Badge = label
	d.mu.Unlock()

	return nil
}

// SetBadge sets the badge label on the application icon.
// Available default badge labels:
// Single space " " empty badge
// Empty string "" dot "●" indeterminate badge
func (d *darwinDock) SetBadge(label string) error {
	// Always pick a label (use "●" if empty).
	if label == "" {
		label = "●" // Default badge character
	}
	return d.setBadge(&label)
}

// SetCustomBadge is not supported on macOS, SetBadge is called instead.
func (d *darwinDock) SetCustomBadge(label string, options BadgeOptions) error {
	return d.SetBadge(label)
}

// RemoveBadge removes the badge label from the application icon.
func (d *darwinDock) RemoveBadge() error {
	return d.setBadge(nil)
}

// GetBadge returns the badge label on the application icon.
func (d *darwinDock) GetBadge() *string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.Badge
}
