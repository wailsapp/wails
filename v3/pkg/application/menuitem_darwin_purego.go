//go:build darwin && purego && !ios && !server

// Package application - CGO-free macOS menu-item backend.
//
// This is the purego counterpart of menuitem_darwin.go. Instead of subclassing
// NSMenuItem in Objective-C (the `MenuItem` class with a `menuItemID` property
// and a `handleClick` method), we use a stock NSMenuItem plus:
//
//   - the item's `tag`, which we set to the Wails menu-item id, and
//   - a single shared target object (WailsMenuItemTarget) whose
//     `menuItemClicked:` method reads the sender's tag and forwards it to
//     processMenuItemClick.
//
// Role-based items keep target == nil so their action travels up the responder
// chain, exactly like the cgo build.
package application

import (
	"sync"
	"unsafe"

	"github.com/ebitengine/purego/objc"
)

// ---------------------------------------------------------------------------
// unsafe.Pointer <-> objc id bridging
//
// The struct fields (nsMenu, nsMenuItem) are typed unsafe.Pointer to stay
// source-compatible with the rest of the package (application_darwin,
// systemtray_darwin, ...). We store the raw Objective-C id inside the pointer
// and convert back on use. This is the documented pattern for this backend.
// ---------------------------------------------------------------------------

func idFromPtr(p unsafe.Pointer) id { return id(uintptr(p)) }

func ptrFromID(o id) unsafe.Pointer { return unsafe.Pointer(o.ptr()) }

// menuIsMainThread reports whether the caller is running on the main (AppKit)
// thread, without depending on the app-level helper.
func menuIsMainThread() bool {
	return get[bool](class("NSThread"), "isMainThread")
}

// runOnMainMenu runs fn on the main thread. If already on the main thread it
// runs inline (InvokeSync would deadlock), otherwise it dispatches and waits.
// This mirrors the isMainThread / dispatch_sync pattern the cgo setters use to
// avoid the stale-menu-state race (wailsapp/wails#5002).
func runOnMainMenu(fn func()) {
	if menuIsMainThread() {
		fn()
		return
	}
	InvokeSync(fn)
}

// ---------------------------------------------------------------------------
// Shared click target
// ---------------------------------------------------------------------------

var (
	menuTargetOnce     sync.Once
	menuTargetInstance id
)

// sharedMenuTarget lazily registers the WailsMenuItemTarget class and returns a
// singleton instance used as the target for every custom-callback menu item.
func sharedMenuTarget() id {
	menuTargetOnce.Do(func() {
		cls := registerDelegateClass("WailsMenuItemTarget", "NSObject", nil, []objc.MethodDef{
			{Cmd: sel_("menuItemClicked:"), Fn: menuItemClickedIMP},
		})
		menuTargetInstance = cls.send("new")
	})
	return menuTargetInstance
}

// menuItemClickedIMP is the IMP backing -[WailsMenuItemTarget menuItemClicked:].
// AppKit invokes it with the NSMenuItem as the sender; we read the tag we set
// (the Wails menu-item id) and forward it to processMenuItemClick.
func menuItemClickedIMP(self objc.ID, cmd objc.SEL, sender objc.ID) {
	tag := objc.Send[int](sender, sel_("tag"))
	processMenuItemClick(uint(tag))
}

// processMenuItemClick delivers a menu-item click to the application's menu
// dispatch loop. In the cgo backend this is a //export'd C callback; here it is
// a plain Go function invoked from menuItemClickedIMP.
func processMenuItemClick(menuID uint) {
	menuItemClicked <- menuID
}

// ---------------------------------------------------------------------------
// macosMenuItem
// ---------------------------------------------------------------------------

type macosMenuItem struct {
	menuItem *MenuItem

	nsMenuItem unsafe.Pointer

	// customCallback records whether this item drives a user callback (target
	// is the shared click target) versus a role-based action (target nil, sent
	// up the responder chain). It replaces the cgo build's action-selector
	// comparison in setDisabled.
	customCallback bool
}

func (m macosMenuItem) setTooltip(tooltip string) {
	item := idFromPtr(m.nsMenuItem)
	runOnMainMenu(func() {
		item.send("setToolTip:", nsString(tooltip))
	})
}

func (m macosMenuItem) setLabel(s string) {
	item := idFromPtr(m.nsMenuItem)
	runOnMainMenu(func() {
		item.send("setTitle:", nsString(s))
	})
}

func (m macosMenuItem) setDisabled(disabled bool) {
	item := idFromPtr(m.nsMenuItem)
	runOnMainMenu(func() {
		item.send("setEnabled:", !disabled)
		// Handle target based on whether the item uses a custom callback or a
		// role-based action.
		if disabled {
			item.send("setTarget:", id(0))
			return
		}
		if m.customCallback {
			// Custom callback: target the shared click handler.
			item.send("setTarget:", sharedMenuTarget())
		} else {
			// Role-based: leave target nil so the action is sent up the
			// responder chain.
			item.send("setTarget:", id(0))
		}
	})
}

func (m macosMenuItem) setChecked(checked bool) {
	item := idFromPtr(m.nsMenuItem)
	// NSControlStateValueOn == 1, NSControlStateValueOff == 0.
	state := 0
	if checked {
		state = 1
	}
	runOnMainMenu(func() {
		item.send("setState:", state)
	})
}

func (m macosMenuItem) setHidden(hidden bool) {
	item := idFromPtr(m.nsMenuItem)
	runOnMainMenu(func() {
		item.send("setHidden:", hidden)
	})
}

func (m macosMenuItem) setBitmap(bitmap []byte) {
	if len(bitmap) == 0 {
		return
	}
	item := idFromPtr(m.nsMenuItem)
	image := class("NSImage").send("alloc").send("initWithData:", nsData(bitmap))
	item.send("setImage:", image)
	// The menu item retains its image; drop the creation reference.
	image.send("release")
}

func (m macosMenuItem) setAccelerator(accelerator *accelerator) {
	item := idFromPtr(m.nsMenuItem)

	// Default to no accelerator.
	keyEquivalent := ""
	modifier := 0
	if accelerator != nil {
		keyEquivalent = translateKey(accelerator.Key)
		modifier = toMacModifier(accelerator.Modifiers)
	}

	item.send("setKeyEquivalent:", nsString(keyEquivalent))
	item.send("setKeyEquivalentModifierMask:", uint(modifier))
}

// destroy releases the owning +1 from `new`. Idempotent: both a menu rebuild
// (processMenu) and MenuItem.Destroy() may call it.
func (m *macosMenuItem) destroy() {
	if m.nsMenuItem == nil {
		return
	}
	idFromPtr(m.nsMenuItem).send("release")
	m.nsMenuItem = nil
}

func newMenuItemImpl(item *MenuItem) *macosMenuItem {
	result := &macosMenuItem{
		menuItem: item,
	}

	// Stock NSMenuItem; `new` == alloc+init and returns a +1 retained object
	// that stays alive for the lifetime of the (retained) menu.
	nsItem := class("NSMenuItem").send("new")

	// Label.
	nsItem.send("setTitle:", nsString(item.label))

	// Action / target wiring.
	selector := getSelectorForRole(item.role)
	if selector != "" {
		// Role-based action: send it up the responder chain (target nil).
		nsItem.send("setAction:", sel_(selector))
		nsItem.send("setTarget:", id(0))
		result.customCallback = false
	} else {
		// Custom callback: route through the shared target unless disabled.
		nsItem.send("setAction:", sel_("menuItemClicked:"))
		if item.disabled {
			nsItem.send("setTarget:", id(0))
		} else {
			nsItem.send("setTarget:", sharedMenuTarget())
		}
		result.customCallback = true
	}

	// Enabled state.
	nsItem.send("setEnabled:", !item.disabled)

	// Tooltip (always applied, matching the cgo build).
	nsItem.send("setToolTip:", nsString(item.tooltip))

	// Tag carries the Wails menu-item id so the click handler can resolve it.
	nsItem.send("setTag:", int(item.id))

	result.nsMenuItem = ptrFromID(nsItem)

	switch item.itemType {
	case checkbox, radio:
		result.setChecked(item.checked)
	}

	if item.accelerator != nil {
		result.setAccelerator(item.accelerator)
	}

	return result
}

// translateKey maps a Wails accelerator key name to the single-character string
// AppKit expects as an NSMenuItem keyEquivalent. Special keys map to their
// function-key unicode code points; anything else is returned verbatim.
//
// This mirrors translateKey() in menuitem_darwin.go.
func translateKey(key string) string {
	if key == "" {
		return ""
	}
	if r, ok := keyEquivalentRunes[key]; ok {
		return string(r)
	}
	return key
}

// keyEquivalentRunes maps special key names to their AppKit unicode code points
// (see NSText.h function-key constants and the cgo translateKey()).
var keyEquivalentRunes = map[string]rune{
	"backspace": 0x0008,
	"tab":       0x0009,
	"return":    0x000d,
	"enter":     0x000d,
	"escape":    0x001b,
	"left":      0xf702,
	"right":     0xf703,
	"up":        0xf700,
	"down":      0xf701,
	"space":     0x0020,
	"delete":    0x007f,
	"home":      0x2196,
	"end":       0x2198,
	"page up":   0x21de,
	"page down": 0x21df,
	"f1":        0xf704,
	"f2":        0xf705,
	"f3":        0xf706,
	"f4":        0xf707,
	"f5":        0xf708,
	"f6":        0xf709,
	"f7":        0xf70a,
	"f8":        0xf70b,
	"f9":        0xf70c,
	"f10":       0xf70d,
	"f11":       0xf70e,
	"f12":       0xf70f,
	"f13":       0xf710,
	"f14":       0xf711,
	"f15":       0xf712,
	"f16":       0xf713,
	"f17":       0xf714,
	"f18":       0xf715,
	"f19":       0xf716,
	"f20":       0xf717,
	"f21":       0xf718,
	"f22":       0xf719,
	"f23":       0xf71a,
	"f24":       0xf71b,
	"f25":       0xf71c,
	"f26":       0xf71d,
	"f27":       0xf71e,
	"f28":       0xf71f,
	"f29":       0xf720,
	"f30":       0xf721,
	"f31":       0xf722,
	"f32":       0xf723,
	"f33":       0xf724,
	"f34":       0xf725,
	"f35":       0xf726,
	"numLock":   0xf739,
}
