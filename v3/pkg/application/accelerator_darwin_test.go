//go:build darwin && !ios && !server

package application

import (
	"strconv"
	"testing"
)

// With Control held, macOS collapses letters to their control codes: Control-A
// produces U+0001, which is also what an older reading of this took to mean
// Home. Naming a press from its characters before its key code therefore made
// eight combinations impossible to bind, and worse, ran the binding they were
// mistaken for.
func TestControlLetterIsNamedAsTheLetter(t *testing.T) {
	cases := []struct {
		keyCode   uint16
		character rune
		want      string
	}{
		{0, '\x01', "ctrl+a"},
		{2, '\x04', "ctrl+d"},
		{4, '\b', "ctrl+h"},
		{40, '\x0b', "ctrl+k"},
		{37, '\x0c', "ctrl+l"},
		{46, '\r', "ctrl+m"},
		{45, '\x0e', "ctrl+n"},
		{33, '\x1b', "ctrl+["},
		// Control-B and Control-C never collided, and must not regress.
		{11, '\x02', "ctrl+b"},
		{8, '\x03', "ctrl+c"},
	}
	for _, c := range cases {
		got := macAccelerator(c.keyCode, macModifierControl, c.character, true)
		if got != c.want {
			t.Errorf("key code %d with Control held was named %q, not %q", c.keyCode, got, c.want)
		}
	}
}

// The keys the control codes were mistaken for have key codes of their own, so
// naming them does not depend on the characters they produce.
func TestNavigationKeysAreNamedFromTheirKeyCode(t *testing.T) {
	cases := []struct {
		name    string
		keyCode uint16
		want    string
	}{
		{"return", 36, "enter"},
		{"keypad enter", 76, "enter"},
		{"home", 115, "home"},
		{"end", 119, "end"},
		{"page up", 116, "page up"},
		{"page down", 121, "page down"},
		{"delete", 51, "delete"},
		{"escape", 53, "escape"},
		{"tab", 48, "tab"},
		{"space", 49, "space"},
		{"left", 123, "left"},
		{"f12", 111, "f12"},
	}
	for _, c := range cases {
		// No character at all, which is what several of these produce.
		if got := macAccelerator(c.keyCode, 0, 0, false); got != c.want {
			t.Errorf("%s was named %q, not %q", c.name, got, c.want)
		}
	}
}

// Control-Space produces a NUL. It has to stay bindable.
func TestControlSpaceIsNamed(t *testing.T) {
	if got := macAccelerator(49, macModifierControl, 0, true); got != "ctrl+space" {
		t.Errorf("Control-Space was named %q, not %q", got, "ctrl+space")
	}
}

func TestModifiersAreNamedInOrder(t *testing.T) {
	all := uint(macModifierCommand | macModifierOption | macModifierControl | macModifierShift)
	if got := macAccelerator(0, all, 'a', true); got != "shift+ctrl+option+cmd+a" {
		t.Errorf("the press was named %q, not %q", got, "shift+ctrl+option+cmd+a")
	}
	if got := macAccelerator(0, 0, 'a', true); got != "a" {
		t.Errorf("an unmodified press was named %q, not %q", got, "a")
	}
}

// A key with no name cannot be bound to, so the press is named as nothing at
// all. Naming it after the modifiers alone produces something like "cmd",
// which no binding can match and which parseAccelerator rejects - so every
// press of an unnamed key with a modifier held was reported as an error.
func TestAnUnknownKeyHasNoName(t *testing.T) {
	for _, modifiers := range []uint{0, macModifierCommand, macModifierShift | macModifierControl} {
		if got := macAccelerator(999, modifiers, 0, false); got != "" {
			t.Errorf("an unknown key with modifiers %#x was named %q, not %q", modifiers, got, "")
		}
	}
}

// The names this produces are handed to parseAccelerator, so anything it
// produces for a real press has to survive that. A bare list of modifiers does
// not.
func TestNothingProducedIsRejectedByTheParser(t *testing.T) {
	presses := []struct {
		name         string
		keyCode      uint16
		modifiers    uint
		character    rune
		hasCharacter bool
	}{
		{"an unknown key held with Command", 999, macModifierCommand, 0, false},
		{"a key with no name held with Shift", 56, macModifierShift, 0, false},
		{"Control-A", 0, macModifierControl, '\x01', true},
		{"F12", 111, 0, 0, false},
	}
	for _, p := range presses {
		got := macAccelerator(p.keyCode, p.modifiers, p.character, p.hasCharacter)
		if got == "" {
			continue // Nothing is dispatched for an unnamed press.
		}
		if _, err := parseAccelerator(got); err != nil {
			t.Errorf("%s was named %q, which the parser rejects: %s", p.name, got, err)
		}
	}
}

// namedKeys accepts these, so a binding can be written for them; they have to
// be reachable from a real press or the binding can never fire. AppKit reports
// them only through the characters they produce.
func TestExtendedFunctionKeysAreNamed(t *testing.T) {
	// NSF21FunctionKey is 0xF718, and they run consecutively to F35.
	for i := 0; i <= 35-21; i++ {
		character := rune(0xF718 + i)
		want := "f" + strconv.Itoa(21+i)
		if got := macAccelerator(0xFFFF, 0, character, true); got != want {
			t.Errorf("the key producing %U was named %q, not %q", character, got, want)
		}
	}
	// NSClearLineFunctionKey. Wails calls this key numlock.
	if got := macAccelerator(0xFFFF, 0, 0xF739, true); got != "numlock" {
		t.Errorf("the clear key was named %q, not %q", got, "numlock")
	}
}

// Every name this produces has to be one parseAccelerator accepts, or the
// binding it was meant for can never match.
func TestEveryKeyNameCanBeParsed(t *testing.T) {
	for keyCode, name := range macKeyNames {
		if _, err := parseAccelerator(name); err != nil {
			t.Errorf("key code %d is named %q, which cannot be parsed: %s", keyCode, name, err)
		}
	}
	for character, name := range macCharacterNames {
		if _, err := parseAccelerator(name); err != nil {
			t.Errorf("character %q is named %q, which cannot be parsed: %s", character, name, err)
		}
	}
}
