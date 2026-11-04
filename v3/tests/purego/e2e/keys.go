package main

import (
	"fmt"
	"regexp"
)

// The X keysym names the Linux input driver uses. They are plain data, so
// this file has no build tag and is unit tested on every OS.

// x11NamedKeys maps the key names chords and tests use to X keysym names.
var x11NamedKeys = map[string]string{
	"return":    "Return",
	"enter":     "Return",
	"escape":    "Escape",
	"esc":       "Escape",
	"tab":       "Tab",
	"space":     "space",
	"backspace": "BackSpace",
	"delete":    "Delete",
	"up":        "Up",
	"down":      "Down",
	"left":      "Left",
	"right":     "Right",
	"home":      "Home",
	"end":       "End",
	"pageup":    "Prior",
	"pagedown":  "Next",
	"insert":    "Insert",
}

// punctKeysyms maps ASCII punctuation to X keysym names. The keysym is the
// unshifted key; the driver decides from the keymap whether Shift is needed.
var punctKeysyms = map[rune]string{
	'!': "exclam", '"': "quotedbl", '#': "numbersign", '$': "dollar",
	'%': "percent", '&': "ampersand", '\'': "apostrophe", '(': "parenleft",
	')': "parenright", '*': "asterisk", '+': "plus", ',': "comma",
	'-': "minus", '.': "period", '/': "slash", ':': "colon",
	';': "semicolon", '<': "less", '=': "equal", '>': "greater",
	'?': "question", '@': "at", '[': "bracketleft", '\\': "backslash",
	']': "bracketright", '^': "asciicircum", '_': "underscore", '`': "grave",
	'{': "braceleft", '|': "bar", '}': "braceright", '~': "asciitilde",
}

var fKey = regexp.MustCompile(`^f([1-9]|1[0-2])$`)

// keysymForKey returns the X keysym name for a chord key name such as
// "k", "7", "return", "f10", "/" or "space".
func keysymForKey(key string) (string, error) {
	if n, ok := x11NamedKeys[key]; ok {
		return n, nil
	}
	if fKey.MatchString(key) {
		return "F" + key[1:], nil
	}
	if r := []rune(key); len(r) == 1 {
		if name, ok := keysymForRune(r[0]); ok {
			return name, nil
		}
	}
	return "", fmt.Errorf("no X keysym for key %q", key)
}

// keysymForRune returns the X keysym name for a character typed with
// typeText. Shift is decided later from the keymap.
func keysymForRune(r rune) (string, bool) {
	switch {
	case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
		return string(r), true
	case r == ' ':
		return "space", true
	case r == '\n':
		return "Return", true
	case r == '\t':
		return "Tab", true
	}
	name, ok := punctKeysyms[r]
	return name, ok
}

// modifierKeysym returns the X keysym of a chord modifier, or "" if unknown.
func modifierKeysym(mod string) string {
	switch mod {
	case "ctrl":
		return "Control_L"
	case "shift":
		return "Shift_L"
	case "alt":
		return "Alt_L"
	case "cmd":
		return "Super_L"
	}
	return ""
}

// needsShift decides how to type a keysym given the keysyms on the two
// levels of its keycode (level0: unshifted, level1: shifted). It returns
// shift=true when the keysym is only reachable with Shift, and ok=false when
// the keycode cannot produce it at all. The unshifted level wins when both
// levels carry the keysym.
func needsShift(level0, level1, want uint64) (shift, ok bool) {
	switch {
	case want == 0:
		return false, false
	case level0 == want:
		return false, true
	case level1 == want:
		return true, true
	}
	return false, false
}
