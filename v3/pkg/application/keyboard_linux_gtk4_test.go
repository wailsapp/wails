//go:build linux && !android && !gtk3 && !server

package application

import "testing"

// The expected keysyms are written out from X11 keysymdef.h, not copied from
// the map, so a wrong map entry fails here.
func TestGTK4NamedKeysAreMainKeyboardKeysyms(t *testing.T) {
	want := map[string]uint32{
		"backspace":      0xff08,
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
		"delete":         0xffff,
		"forward delete": 0xffff,
		"f1":             0xffbe,
		"f12":            0xffc9,
		"f35":            0xffe0,
		"numlock":        0xff7f,
	}
	for key, keysym := range want {
		if got, ok := gtk4NamedKeys[key]; !ok || got != keysym {
			t.Errorf("gtk4NamedKeys[%q] = %#x (present %v), want %#x", key, got, ok, keysym)
		}
	}
}

// Keypad keysyms are 0xff80..0xffbd. None of the named keys may use one: the
// keypad key is a different key from the main keyboard's.
func TestGTK4NamedKeysAvoidKeypadKeysyms(t *testing.T) {
	for key, keysym := range gtk4NamedKeys {
		if keysym >= 0xff80 && keysym <= 0xffbd {
			t.Errorf("gtk4NamedKeys[%q] = %#x is a keypad keysym", key, keysym)
		}
	}
}
