//go:build linux && !android && gtk3 && !server

package application

// Pure GTK3 keyboard and dialog decisions. They hold no C state, so the cgo
// and purego GTK3 builds share them and the tests run without a display.

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

// GdkModifierType bits (gdk/gdktypes.h, GTK3).
const (
	gtk3ShiftMask   uint32 = 1 << 0
	gtk3ControlMask uint32 = 1 << 2
	gtk3Mod1Mask    uint32 = 1 << 3 // Alt
	gtk3SuperMask   uint32 = 1 << 26
)

// gtk3ResponseDeleteEvent is GTK_RESPONSE_DELETE_EVENT. gtk_dialog_run returns
// it when the dialog is closed without a button: Escape, or the window's close
// control.
const gtk3ResponseDeleteEvent = -4

// gtk3AcceleratorKeys maps the named keys of keys.go to GDK keysyms. Keypad
// keysyms (0xff9x) are not used: they are a different key from the main
// keyboard's Home, End, Delete and PageUp.
var gtk3AcceleratorKeys = map[string]uint32{
	"backspace":      0xff08, // GDK_KEY_BackSpace
	"delete":         0xffff, // GDK_KEY_Delete, as GTK4, Windows (VK_DELETE) and keys_linux.go
	"forward delete": 0xffff, // GDK_KEY_Delete
	"tab":            0xff09,
	"return":         0xff0d,
	"enter":          0xff0d,
	"escape":         0xff1b,
	"left":           0xff51,
	"up":             0xff52,
	"right":          0xff53,
	"down":           0xff54,
	"page up":        0xff55,
	"page down":      0xff56,
	"end":            0xff57,
	"home":           0xff50,
	"space":          0x0020,
	"numlock":        0xff7f,
}

// gtk3AcceleratorKeyval returns the GDK keysym for a Wails accelerator key, or
// 0 if the key has no GDK keysym. Named keys match keys.go, case-insensitively. A single character is its Unicode code point; GDK maps code
// points above 0xff to 0x01000000 plus the code point.
func gtk3AcceleratorKeyval(key string) uint32 {
	key = strings.ToLower(key)
	if keyval, ok := gtk3AcceleratorKeys[key]; ok {
		return keyval
	}
	if n, ok := gtk3FunctionKey(key); ok {
		return 0xffbd + uint32(n) // GDK_KEY_F1 is 0xffbe
	}
	r, size := utf8.DecodeRuneInString(key)
	if r == utf8.RuneError || size == 0 || size != len(key) {
		return 0
	}
	if r >= 'A' && r <= 'Z' {
		r += 'a' - 'A' // GTK canonicalises letters to lower case, with Shift as a modifier
	}
	if r < 0x100 {
		return uint32(r)
	}
	return 0x01000000 | uint32(r)
}

// gtk3FunctionKey parses "f1".."f35" (keys.go accepts up to f35).
func gtk3FunctionKey(key string) (int, bool) {
	if len(key) < 2 || key[0] != 'f' {
		return 0, false
	}
	n, err := strconv.Atoi(key[1:])
	if err != nil || n < 1 || n > 35 {
		return 0, false
	}
	return n, true
}

// gtk3AcceleratorModifiers maps Wails modifiers to GDK modifier bits.
// CmdOrCtrl is Control on Linux, matching the rest of the Linux backend.
func gtk3AcceleratorModifiers(modifiers []modifier) uint32 {
	var result uint32
	for _, mod := range modifiers {
		switch mod {
		case ShiftKey:
			result |= gtk3ShiftMask
		case ControlKey, CmdOrCtrlKey:
			result |= gtk3ControlMask
		case OptionOrAltKey:
			result |= gtk3Mod1Mask
		case SuperKey:
			result |= gtk3SuperMask
		}
	}
	return result
}

// dialogResponseIndex maps a gtk_dialog_run response to a button index, or -1
// for no button. A button response is its index. Closing the dialog (Escape or
// the close control) gives GTK_RESPONSE_DELETE_EVENT, which means the button
// marked SetAsCancel, if there is one. Otherwise the close does nothing.
func dialogResponseIndex(response int, buttons []*Button) int {
	if response >= 0 {
		return response
	}
	if response != gtk3ResponseDeleteEvent {
		return -1
	}
	for i, button := range buttons {
		if button != nil && button.IsCancel {
			return i
		}
	}
	return -1
}

// keyRepeatFilter tells X auto-repeat from new key presses. Auto-repeat sends
// another key press for a key that is still held, with no release in between,
// so a press of a key the window already sees as down is a repeat. The timing
// of the presses does not matter: a fast double tap has a release between the
// two presses and is a new press, and a slow repeat rate is still a repeat.
// The window calls release from the key-release handler, and reset when it
// loses focus, since a key released while another window has focus gives no
// release event here.
type keyRepeatFilter struct {
	held map[uint]bool // hardware keycodes down in this window
}

// press records a key press and reports whether it is a repeat of a key that
// is already down.
func (f *keyRepeatFilter) press(keycode uint) bool {
	if f.held[keycode] {
		return true
	}
	if f.held == nil {
		f.held = make(map[uint]bool)
	}
	f.held[keycode] = true
	return false
}

// release records a key release.
func (f *keyRepeatFilter) release(keycode uint) {
	delete(f.held, keycode)
}

// reset forgets every held key.
func (f *keyRepeatFilter) reset() {
	f.held = nil
}
