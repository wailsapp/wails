//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13
#cgo LDFLAGS: -framework Carbon

#include <Carbon/Carbon.h>

extern void globalShortcutCallback(int id);

static EventHandlerRef gHotKeyHandler = NULL;

// hotKeyHandlerProc is the single Carbon event handler that receives every
// kEventHotKeyPressed event. It extracts the EventHotKeyID we set at
// registration time and forwards the numeric id back into Go.
static OSStatus hotKeyHandlerProc(EventHandlerCallRef nextHandler, EventRef theEvent, void *userData) {
    EventHotKeyID hkID;
    OSStatus status = GetEventParameter(theEvent, kEventParamDirectObject, typeEventHotKeyID,
                                        NULL, sizeof(hkID), NULL, &hkID);
    if (status == noErr) {
        globalShortcutCallback((int)hkID.id);
    }
    return noErr;
}

// installHotKeyHandler installs the shared handler exactly once.
static void installHotKeyHandler(void) {
    if (gHotKeyHandler != NULL) {
        return;
    }
    EventTypeSpec evt;
    evt.eventClass = kEventClassKeyboard;
    evt.eventKind = kEventHotKeyPressed;
    InstallApplicationEventHandler(&hotKeyHandlerProc, 1, &evt, NULL, &gHotKeyHandler);
}

// registerHotKey binds (keyCode, modifiers) to id and returns the OSStatus.
// On success *outRef receives the hot key reference used to unregister later.
static int registerHotKey(unsigned int keyCode, unsigned int modifiers, int id, EventHotKeyRef *outRef) {
    installHotKeyHandler();
    EventHotKeyID hkID;
    hkID.signature = 'WLgs'; // "Wails global shortcut"
    hkID.id = (unsigned int)id;
    OSStatus status = RegisterEventHotKey(keyCode, modifiers, hkID,
                                          GetApplicationEventTarget(), 0, outRef);
    return (int)status;
}

static int unregisterHotKey(EventHotKeyRef ref) {
    if (ref == NULL) {
        return 0;
    }
    return (int)UnregisterEventHotKey(ref);
}
*/
import "C"

import (
	"fmt"
)

// Carbon classic modifier masks (Events.h). RegisterEventHotKey expects these,
// not the Cocoa NSEventModifierFlag values.
const (
	carbonCmdKey     = 0x0100
	carbonShiftKey   = 0x0200
	carbonOptionKey  = 0x0800
	carbonControlKey = 0x1000
)

// macosGlobalShortcuts implements globalShortcutImpl using the Carbon Event
// Manager's RegisterEventHotKey API. Carbon's hot key API is the standard,
// still-supported mechanism for system-wide hot keys on macOS and - unlike a
// CGEventTap or an NSEvent global monitor - does not require Accessibility
// permission.
type macosGlobalShortcuts struct {
	manager *GlobalShortcutManager
	refs    map[int]C.EventHotKeyRef
}

func newGlobalShortcutImpl(manager *GlobalShortcutManager) globalShortcutImpl {
	return &macosGlobalShortcuts{
		manager: manager,
		refs:    make(map[int]C.EventHotKeyRef),
	}
}

func (g *macosGlobalShortcuts) register(id int, accel *accelerator) error {
	keyCode, ok := macKeyCodes[accel.Key]
	if !ok {
		return fmt.Errorf("key %q is not supported as a global shortcut on macOS", accel.Key)
	}

	var mods C.uint
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

	var ref C.EventHotKeyRef
	status := C.registerHotKey(C.uint(keyCode), mods, C.int(id), &ref)
	if status != 0 {
		// -9878 (eventHotKeyExistsErr) means the combination is already taken.
		if status == -9878 {
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
	if status := C.unregisterHotKey(ref); status != 0 {
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

//export globalShortcutCallback
func globalShortcutCallback(id C.int) {
	if globalApplication != nil && globalApplication.GlobalShortcut != nil {
		globalApplication.GlobalShortcut.dispatch(int(id))
	}
}

// macKeyCodes maps Wails accelerator key names (already lower-cased by
// parseAccelerator) to macOS hardware virtual key codes (kVK_* from Carbon's
// HIToolbox/Events.h).
//
// NOTE: macOS hot keys are bound to *hardware* key codes, not characters, so
// this table assumes a standard ANSI/QWERTY physical layout. On non-QWERTY
// layouts the physical key in the QWERTY position is what triggers the
// shortcut. This matches the behaviour of essentially every macOS global hot
// key implementation; a layout-aware mapping (via UCKeyTranslate) could be
// added later if required.
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
