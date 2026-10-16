//go:build darwin && purego && !ios && !server

package application

// CGO-free re-implementation of global_shortcut_darwin.go.
//
// Like the cgo version this uses the Carbon Event Manager's RegisterEventHotKey
// API — the standard, still-supported mechanism for system-wide hot keys on
// macOS which, unlike a CGEventTap or NSEvent global monitor, does not require
// Accessibility permission. The Carbon functions are bound directly through
// purego (dlopen + RegisterLibFunc) instead of cgo, and the hot-key event
// handler is a purego callback rather than a compiled C function.

import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

// Carbon classic modifier masks (Events.h). RegisterEventHotKey expects these,
// not the Cocoa NSEventModifierFlag values.
const (
	carbonCmdKey     = 0x0100
	carbonShiftKey   = 0x0200
	carbonOptionKey  = 0x0800
	carbonControlKey = 0x1000
)

// Carbon event constants (CarbonEvents.h / Events.h), all FourCharCodes or
// small integers.
const (
	kEventClassKeyboard     = 0x6B657962 // 'keyb'
	kEventHotKeyPressed     = 5
	kEventParamDirectObject = 0x2D2D2D2D // '----'
	typeEventHotKeyID       = 0x686B6964 // 'hkid'
	hotKeySignature         = 0x574C6773 // 'WLgs' — "Wails global shortcut"

	// eventHotKeyExistsErr: the combination is already taken.
	eventHotKeyExistsErr = -9878
)

// eventHotKeyID mirrors the Carbon struct { OSType signature; UInt32 id; }.
type eventHotKeyID struct {
	signature uint32
	id        uint32
}

// eventTypeSpec mirrors the Carbon struct { OSType eventClass; UInt32 eventKind; }.
type eventTypeSpec struct {
	eventClass uint32
	eventKind  uint32
}

const frameworkCarbon = "/System/Library/Frameworks/Carbon.framework/Carbon"

var (
	carbonOnce sync.Once

	getApplicationEventTarget func() uintptr
	installEventHandler       func(target, handler uintptr, numTypes uint, list *eventTypeSpec, userData uintptr, outRef *uintptr) int32
	registerEventHotKeyFn     func(keyCode, modifiers uint32, hotKeyID eventHotKeyID, target uintptr, options uint32, outRef *uintptr) int32
	unregisterEventHotKeyFn   func(ref uintptr) int32
	getEventParameter         func(event uintptr, name, typ uint32, outActualType *uint32, bufferSize uint, outActualSize *uint, outData unsafe.Pointer) int32

	hotKeyHandlerOnce sync.Once
)

// loadCarbon binds the Carbon Event Manager symbols we need.
func loadCarbon() {
	carbonOnce.Do(func() {
		handle, err := purego.Dlopen(frameworkCarbon, purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if err != nil || handle == 0 {
			return
		}
		purego.RegisterLibFunc(&getApplicationEventTarget, handle, "GetApplicationEventTarget")
		purego.RegisterLibFunc(&installEventHandler, handle, "InstallEventHandler")
		purego.RegisterLibFunc(&registerEventHotKeyFn, handle, "RegisterEventHotKey")
		purego.RegisterLibFunc(&unregisterEventHotKeyFn, handle, "UnregisterEventHotKey")
		purego.RegisterLibFunc(&getEventParameter, handle, "GetEventParameter")
	})
}

// installHotKeyHandler installs the shared Carbon event handler exactly once.
func installHotKeyHandler() {
	hotKeyHandlerOnce.Do(func() {
		spec := eventTypeSpec{eventClass: kEventClassKeyboard, eventKind: kEventHotKeyPressed}
		var ref uintptr
		cb := purego.NewCallback(hotKeyHandlerProc)
		installEventHandler(getApplicationEventTarget(), cb, 1, &spec, 0, &ref)
	})
}

// hotKeyHandlerProc is the single Carbon event handler that receives every
// kEventHotKeyPressed event. It extracts the EventHotKeyID we set at
// registration time and forwards the numeric id back into Go.
func hotKeyHandlerProc(nextHandler, theEvent, userData uintptr) int32 {
	var hkID eventHotKeyID
	status := getEventParameter(theEvent,
		kEventParamDirectObject,
		typeEventHotKeyID,
		nil,
		uint(unsafe.Sizeof(hkID)),
		nil,
		unsafe.Pointer(&hkID),
	)
	if status == 0 {
		globalShortcutCallback(int(hkID.id))
	}
	return 0 // noErr
}

// macosGlobalShortcuts implements globalShortcutImpl using the Carbon Event
// Manager's RegisterEventHotKey API.
type macosGlobalShortcuts struct {
	manager *GlobalShortcutManager
	refs    map[int]uintptr
}

func newGlobalShortcutImpl(manager *GlobalShortcutManager) globalShortcutImpl {
	return &macosGlobalShortcuts{
		manager: manager,
		refs:    make(map[int]uintptr),
	}
}

func (g *macosGlobalShortcuts) register(id int, accel *accelerator) error {
	loadCarbon()
	if registerEventHotKeyFn == nil {
		return fmt.Errorf("Carbon Event Manager unavailable")
	}

	keyCode, ok := macKeyCodes[accel.Key]
	if !ok {
		return fmt.Errorf("key %q is not supported as a global shortcut on macOS", accel.Key)
	}

	var mods uint32
	for _, m := range accel.Modifiers {
		switch m {
		case CmdOrCtrlKey, SuperKey:
			mods |= carbonCmdKey
		case ControlKey:
			mods |= carbonControlKey
		case OptionOrAltKey:
			mods |= carbonOptionKey
		case ShiftKey:
			mods |= carbonShiftKey
		}
	}

	installHotKeyHandler()

	hkID := eventHotKeyID{signature: hotKeySignature, id: uint32(id)}
	var ref uintptr
	status := registerEventHotKeyFn(uint32(keyCode), mods, hkID, getApplicationEventTarget(), 0, &ref)
	if status != 0 {
		// -9878 (eventHotKeyExistsErr) means the combination is already taken.
		if status == eventHotKeyExistsErr {
			return fmt.Errorf("the shortcut is already registered (possibly by another application) (OSStatus %d)", int(status))
		}
		return fmt.Errorf("RegisterEventHotKey failed (OSStatus %d)", int(status))
	}
	g.refs[id] = ref
	return nil
}

func (g *macosGlobalShortcuts) unregister(id int) error {
	ref, ok := g.refs[id]
	if !ok {
		return nil
	}
	delete(g.refs, id)
	if status := unregisterEventHotKeyFn(ref); status != 0 {
		return fmt.Errorf("UnregisterEventHotKey failed (OSStatus %d)", int(status))
	}
	return nil
}

func (g *macosGlobalShortcuts) unregisterAll() error {
	var firstErr error
	for id := range g.refs {
		if err := g.unregister(id); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// globalShortcutCallback is invoked (from the Carbon event handler) with the
// numeric id of the hot key that fired.
func globalShortcutCallback(id int) {
	if globalApplication != nil && globalApplication.GlobalShortcut != nil {
		globalApplication.GlobalShortcut.dispatch(id)
	}
}

// macKeyCodes maps Wails accelerator key names (already lower-cased by
// parseAccelerator) to macOS hardware virtual key codes (kVK_* from Carbon's
// HIToolbox/Events.h).
//
// NOTE: macOS hot keys are bound to *hardware* key codes, not characters, so
// this table assumes a standard ANSI/QWERTY physical layout.
var macKeyCodes = map[string]int{
	// Letters
	"a": 0, "s": 1, "d": 2, "f": 3, "h": 4, "g": 5, "z": 6, "x": 7,
	"c": 8, "v": 9, "b": 11, "q": 12, "w": 13, "e": 14, "r": 15, "y": 16,
	"t": 17, "o": 31, "u": 32, "i": 34, "p": 35, "l": 37, "j": 38, "k": 40,
	"n": 45, "m": 46,
	// Number row
	"1": 18, "2": 19, "3": 20, "4": 21, "6": 22, "5": 23, "9": 25, "7": 26,
	"8": 28, "0": 29,
	// Punctuation
	"=": 24, "-": 27, "]": 30, "[": 33, "'": 39, ";": 41, "\\": 42,
	",": 43, "/": 44, ".": 47, "`": 50, "+": 24,
	// Named keys
	"return":    36,
	"enter":     36,
	"tab":       48,
	"space":     49,
	"backspace": 51,  // kVK_Delete (labelled "delete" on Mac keyboards)
	"delete":    117, // kVK_ForwardDelete
	"escape":    53,
	"home":      115,
	"page up":   116,
	"end":       119,
	"page down": 121,
	"left":      123,
	"right":     124,
	"down":      125,
	"up":        126,
	// Function keys
	"f1": 122, "f2": 120, "f3": 99, "f4": 118, "f5": 96, "f6": 97,
	"f7": 98, "f8": 100, "f9": 101, "f10": 109, "f11": 103, "f12": 111,
	"f13": 105, "f14": 107, "f15": 113, "f16": 106, "f17": 64, "f18": 79,
	"f19": 80, "f20": 90,
}
