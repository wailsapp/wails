//go:build darwin && !ios && !server

package application

import "strings"

// Naming a key press for the key binding system.
//
// AppKit reports a press three ways: a virtual key code, the modifiers held,
// and the characters the press produces. Only the key code identifies the
// physical key; the characters depend on what else is held down, and with
// Control held macOS collapses letters to their control codes - Control-A
// becomes U+0001. Anything that reads the characters first therefore cannot
// tell Control-A from Home, so the key code is consulted first here and the
// characters are only a fallback for keys with no code of their own.

// AppKit's modifier bits, from NSEvent.h. Declared here so the naming can be
// tested without AppKit.
const (
	macModifierShift   = 1 << 17
	macModifierControl = 1 << 18
	macModifierOption  = 1 << 19
	macModifierCommand = 1 << 20
)

// macKeyNames maps macOS virtual key codes to the names key bindings use. A
// standard US layout; keys whose name would depend on the layout are named
// from the characters instead.
var macKeyNames = map[uint16]string{
	// Function keys
	122: "f1", 120: "f2", 99: "f3", 118: "f4", 96: "f5",
	97: "f6", 98: "f7", 100: "f8", 101: "f9", 109: "f10",
	103: "f11", 111: "f12", 105: "f13", 107: "f14", 113: "f15",
	106: "f16", 64: "f17", 79: "f18", 80: "f19", 90: "f20",

	// Letters
	0: "a", 11: "b", 8: "c", 2: "d", 14: "e", 3: "f", 5: "g",
	4: "h", 34: "i", 38: "j", 40: "k", 37: "l", 46: "m", 45: "n",
	31: "o", 35: "p", 12: "q", 15: "r", 1: "s", 17: "t", 32: "u",
	9: "v", 13: "w", 7: "x", 16: "y", 6: "z",

	// Digits
	29: "0", 18: "1", 19: "2", 20: "3", 21: "4",
	23: "5", 22: "6", 26: "7", 28: "8", 25: "9",

	// Editing and navigation
	36: "enter", 76: "enter", 51: "delete", 117: "forward delete",
	123: "left", 124: "right", 126: "up", 125: "down",
	115: "home", 119: "end", 116: "page up", 121: "page down",
	48: "tab", 53: "escape", 49: "space",

	// Punctuation, on a standard US layout
	33: "[", 30: "]", 43: ",", 27: "-", 39: "'",
	44: "/", 47: ".", 41: ";", 24: "=", 50: "`", 42: "\\",
}

// macCharacterNames names the keys AppKit only reports through the characters
// they produce. Reached only when the key code is not one this build knows,
// so a Control-letter press can no longer be mistaken for one of them.
var macCharacterNames = map[rune]string{
	'\r':   "enter",
	'\b':   "backspace",
	'\x1b': "escape",
	'\x0b': "page down",
	'\x0e': "page up",
	'\x01': "home",
	'\x04': "end",
}

// macKeyName names the key that was pressed, or returns an empty string for a
// key this build has no name for.
func macKeyName(keyCode uint16, character rune, hasCharacter bool) string {
	if name, ok := macKeyNames[keyCode]; ok {
		return name
	}
	if !hasCharacter {
		return ""
	}
	return macCharacterNames[character]
}

// macAccelerator names a key press the way a key binding is written: the
// modifiers held, in a fixed order, then the key, joined with "+". Empty if
// the key has no name, since a binding could not refer to it.
func macAccelerator(keyCode uint16, modifiers uint, character rune, hasCharacter bool) string {
	var parts []string
	if modifiers&macModifierShift != 0 {
		parts = append(parts, "shift")
	}
	if modifiers&macModifierControl != 0 {
		parts = append(parts, "ctrl")
	}
	if modifiers&macModifierOption != 0 {
		parts = append(parts, "option")
	}
	if modifiers&macModifierCommand != 0 {
		parts = append(parts, "cmd")
	}
	if key := macKeyName(keyCode, character, hasCharacter); key != "" {
		parts = append(parts, key)
	}
	return strings.Join(parts, "+")
}
