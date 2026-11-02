//go:build darwin

package main

import "testing"

// The expected values are the macOS constants (HIToolbox/Events.h kVK_*,
// NSEvent.h flags, types and function keys), written out literally so that a
// typo in the table under test fails here.

func TestChordStroke(t *testing.T) {
	tests := []struct {
		name     string
		chord    string
		code     uint16
		flags    uint
		chars    string
		ignoring string
	}{
		{"cmd+m is the menu accelerator key", "cmd+m", 0x2E, 1 << 20, "m", "m"},
		{"cmd+shift+k types upper case", "cmd+shift+k", 0x28, 1<<20 | 1<<17, "K", "k"},
		{"cmd+shift+g keeps the base key", "cmd+shift+g", 0x05, 1<<20 | 1<<17, "G", "g"},
		{"return is the carriage return", "return", 0x24, 0, "\r", "\r"},
		{"enter is the same key as return", "enter", 0x24, 0, "\r", "\r"},
		{"escape", "escape", 0x35, 0, "\x1b", "\x1b"},
		{"tab", "tab", 0x30, 0, "\t", "\t"},
		{"space", "space", 0x31, 0, " ", " "},
		{"down is an arrow with the function flag", "down", 0x7D, 1 << 23, "", ""},
		{"up", "up", 0x7E, 1 << 23, "", ""},
		{"left", "left", 0x7B, 1 << 23, "", ""},
		{"right", "right", 0x7C, 1 << 23, "", ""},
		{"f1 starts the NSFn range at 0xF704", "f1", 0x7A, 1 << 23, "", ""},
		{"f10 is kVK_F10 with NSF10 characters", "f10", 0x6D, 1 << 23, "", ""},
		{"f12", "f12", 0x6F, 1 << 23, "", ""},
		{"ctrl+l is the control character", "ctrl+l", 0x25, 1 << 18, "\x0c", "l"},
		{"ctrl+shift+k", "ctrl+shift+k", 0x28, 1<<18 | 1<<17, "\x0b", "k"},
		{"shift+1 types the shifted symbol", "shift+1", 0x12, 1 << 17, "!", "1"},
		{"shift+/ types a question mark", "shift+/", 0x2C, 1 << 17, "?", "/"},
		{"alt+tab uses option", "alt+tab", 0x30, 1 << 19, "\t", "\t"},
		{"option+a", "option+a", 0x00, 1 << 19, "a", "a"},
		{"cmd+/", "cmd+/", 0x2C, 1 << 20, "/", "/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := parseChord(tt.chord)
			if err != nil {
				t.Fatal(err)
			}
			s, err := chordStroke(c)
			if err != nil {
				t.Fatal(err)
			}
			if s.code != tt.code {
				t.Errorf("keyCode = %#x, want %#x", s.code, tt.code)
			}
			if s.flags != tt.flags {
				t.Errorf("flags = %#x, want %#x", s.flags, tt.flags)
			}
			if s.chars != tt.chars {
				t.Errorf("characters = %q, want %q", s.chars, tt.chars)
			}
			if s.ignoring != tt.ignoring {
				t.Errorf("charactersIgnoringModifiers = %q, want %q", s.ignoring, tt.ignoring)
			}
		})
	}
}

func TestChordStrokeErrors(t *testing.T) {
	for _, chord := range []string{"banana", "f13", "cmd+", "cmd+é", "cmd+§"} {
		t.Run(chord, func(t *testing.T) {
			c, err := parseChord(chord)
			if err != nil {
				return // rejected by the parser: fine
			}
			if _, err := chordStroke(c); err == nil {
				t.Errorf("chordStroke(%q) succeeded, want an error", chord)
			}
		})
	}
}

func TestCmdOrCtrlIsCommandOnDarwin(t *testing.T) {
	c, err := parseChord("cmdorctrl+shift+k")
	if err != nil {
		t.Fatal(err)
	}
	s, err := chordStroke(c)
	if err != nil {
		t.Fatal(err)
	}
	if want := uint(1<<20 | 1<<17); s.flags != want {
		t.Errorf("cmdorctrl on darwin: flags = %#x, want command+shift %#x", s.flags, want)
	}
}

func TestTypeStroke(t *testing.T) {
	tests := []struct {
		name     string
		r        rune
		code     uint16
		flags    uint
		chars    string
		ignoring string
	}{
		{"lower case letter", 'a', 0x00, 0, "a", "a"},
		{"upper case letter is shift plus the base key", 'A', 0x00, 1 << 17, "A", "a"},
		{"digit", '7', 0x1A, 0, "7", "7"},
		{"slash in a path", '/', 0x2C, 0, "/", "/"},
		{"dot in a file name", '.', 0x2F, 0, ".", "."},
		{"dash in a file name", '-', 0x1B, 0, "-", "-"},
		{"underscore is shift plus minus", '_', 0x1B, 1 << 17, "_", "-"},
		{"colon is shift plus semicolon", ':', 0x29, 1 << 17, ":", ";"},
		{"space", ' ', 0x31, 0, " ", " "},
		{"character with no key on the ANSI layout", 'é', 0, 0, "é", "é"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := typeStroke(tt.r)
			if s.code != tt.code || s.flags != tt.flags || s.chars != tt.chars || s.ignoring != tt.ignoring {
				t.Errorf("typeStroke(%q) = {code %#x flags %#x chars %q ignoring %q}, want {code %#x flags %#x chars %q ignoring %q}",
					tt.r, s.code, s.flags, s.chars, s.ignoring, tt.code, tt.flags, tt.chars, tt.ignoring)
			}
		})
	}
}

// Every printable key must have its own virtual key code: a duplicated value
// would make two keys send the same event.
func TestKeyCodesAreUnique(t *testing.T) {
	seen := map[uint16]rune{}
	for r, code := range keyCodes {
		if prev, ok := seen[code]; ok {
			t.Errorf("keys %q and %q share key code %#x", prev, r, code)
		}
		seen[code] = r
	}
	for _, r := range "abcdefghijklmnopqrstuvwxyz" {
		if _, ok := keyCodes[r]; !ok {
			t.Errorf("no key code for %q", r)
		}
	}
}

func TestContentToWindow(t *testing.T) {
	tests := []struct {
		name         string
		x, y, height float64
		wantX, wantY float64
	}{
		{"top left corner is the top of the window", 0, 0, 500, 0, 500},
		{"bottom left corner", 0, 500, 500, 0, 0},
		{"a point inside", 100, 60, 500, 100, 440},
		{"fractional point", 20.5, 139.25, 480, 20.5, 340.75},
		{"tall content", 10, 10, 1000, 10, 990},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gx, gy := contentToWindow(tt.x, tt.y, tt.height)
			if gx != tt.wantX || gy != tt.wantY {
				t.Errorf("contentToWindow(%v, %v, %v) = (%v, %v), want (%v, %v)",
					tt.x, tt.y, tt.height, gx, gy, tt.wantX, tt.wantY)
			}
		})
	}
}
