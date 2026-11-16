package main

import (
	"runtime"
	"testing"
)

func TestParseChord(t *testing.T) {
	linuxOnly := runtime.GOOS != "darwin"
	tests := []struct {
		in      string
		want    chord
		wantErr bool
	}{
		{in: "ctrl+shift+k", want: chord{ctrl: true, shift: true, key: "k"}},
		{in: "CTRL+SHIFT+K", want: chord{ctrl: true, shift: true, key: "k"}},
		{in: "control+x", want: chord{ctrl: true, key: "x"}},
		{in: "alt+x", want: chord{alt: true, key: "x"}},
		{in: "option+x", want: chord{alt: true, key: "x"}},
		{in: "super+l", want: chord{cmd: true, key: "l"}},
		{in: "cmd+m", want: chord{cmd: true, key: "m"}},
		{in: "return", want: chord{key: "return"}},
		{in: "f10", want: chord{key: "f10"}},
		{in: "ctrl+plus", want: chord{ctrl: true, key: "plus"}},
		{in: "ctrl+ shift + k", want: chord{ctrl: true, shift: true, key: "k"}},
		{in: "ctrl+", wantErr: true},
		{in: "", wantErr: true},
		{in: "hyper+k", wantErr: true},
	}
	for _, tc := range tests {
		got, err := parseChord(tc.in)
		if (err != nil) != tc.wantErr {
			t.Fatalf("parseChord(%q) error = %v, want error %v", tc.in, err, tc.wantErr)
		}
		if err == nil && got != tc.want {
			t.Errorf("parseChord(%q) = %+v, want %+v", tc.in, got, tc.want)
		}
	}

	// cmdorctrl is cmd on macOS and ctrl elsewhere, like Wails accelerators.
	got, err := parseChord("CmdOrCtrl+Shift+K")
	if err != nil {
		t.Fatal(err)
	}
	if linuxOnly {
		if got != (chord{ctrl: true, shift: true, key: "k"}) {
			t.Errorf("CmdOrCtrl+Shift+K on %s = %+v, want ctrl+shift", runtime.GOOS, got)
		}
	} else if got != (chord{cmd: true, shift: true, key: "k"}) {
		t.Errorf("CmdOrCtrl+Shift+K on darwin = %+v, want cmd+shift", got)
	}
}

func TestKeysymForKey(t *testing.T) {
	tests := map[string]string{
		"k": "k", "7": "7", "return": "Return", "enter": "Return",
		"escape": "Escape", "tab": "Tab", "space": "space", "up": "Up",
		"down": "Down", "left": "Left", "right": "Right", "backspace": "BackSpace",
		"f1": "F1", "f10": "F10", "f12": "F12", "/": "slash", "=": "equal",
	}
	for in, want := range tests {
		got, err := keysymForKey(in)
		if err != nil || got != want {
			t.Errorf("keysymForKey(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
	for _, bad := range []string{"", "f13", "f0", "ab", "nosuchkey", "€"} {
		if got, err := keysymForKey(bad); err == nil {
			t.Errorf("keysymForKey(%q) = %q, want an error", bad, got)
		}
	}
}

func TestKeysymForRune(t *testing.T) {
	tests := map[rune]string{
		'a': "a", 'Z': "Z", '0': "0", ' ': "space", '\n': "Return", '\t': "Tab",
		'/': "slash", '.': "period", '-': "minus", '_': "underscore", ':': "colon",
		'?': "question", '"': "quotedbl",
	}
	for in, want := range tests {
		got, ok := keysymForRune(in)
		if !ok || got != want {
			t.Errorf("keysymForRune(%q) = %q, %v; want %q", in, got, ok, want)
		}
	}
	for _, bad := range []rune{'é', '€', '\r', 0} {
		if got, ok := keysymForRune(bad); ok {
			t.Errorf("keysymForRune(%q) = %q, want no keysym", bad, got)
		}
	}
}

// The punctuation table must map each character to a distinct keysym name;
// a copy-paste duplicate would type the wrong character.
func TestPunctuationNamesAreDistinct(t *testing.T) {
	seen := map[string]rune{}
	for r, name := range punctKeysyms {
		if prev, dup := seen[name]; dup {
			t.Errorf("keysym %q is mapped from both %q and %q", name, prev, r)
		}
		seen[name] = r
	}
}

// Keysym values used in the shift decision: 'a' 0x61, 'A' 0x41, '/' 0x2f,
// '?' 0x3f, Return 0xff0d.
func TestNeedsShift(t *testing.T) {
	const (
		a, upperA   uint64 = 0x61, 0x41
		slash, ques uint64 = 0x2f, 0x3f
		ret         uint64 = 0xff0d
		unknown     uint64 = 0x1234
	)
	tests := []struct {
		name                   string
		level0, level1, want   uint64
		wantShift, wantKeyable bool
	}{
		{"letter unshifted", a, upperA, a, false, true},
		{"capital needs shift", a, upperA, upperA, true, true},
		{"slash unshifted", slash, ques, slash, false, true},
		{"question needs shift", slash, ques, ques, true, true},
		{"same on both levels prefers unshifted", ret, ret, ret, false, true},
		{"keysym on neither level", a, upperA, slash, false, false},
		{"unknown keysym", a, upperA, unknown, false, false},
		{"zero keysym", a, upperA, 0, false, false},
	}
	for _, tc := range tests {
		shift, ok := needsShift(tc.level0, tc.level1, tc.want)
		if shift != tc.wantShift || ok != tc.wantKeyable {
			t.Errorf("%s: needsShift = (%v, %v), want (%v, %v)", tc.name, shift, ok, tc.wantShift, tc.wantKeyable)
		}
	}
}

func TestModifierKeysym(t *testing.T) {
	want := map[string]string{"ctrl": "Control_L", "shift": "Shift_L", "alt": "Alt_L", "cmd": "Super_L", "hyper": ""}
	for in, w := range want {
		if got := modifierKeysym(in); got != w {
			t.Errorf("modifierKeysym(%q) = %q, want %q", in, got, w)
		}
	}
}
