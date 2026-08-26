//go:build darwin && !ios && !server

package application

import "testing"

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

// A key with no name cannot be bound, so it is named as nothing rather than as
// a bare list of modifiers.
func TestAnUnknownKeyHasNoName(t *testing.T) {
	if got := macAccelerator(999, macModifierCommand, 0, false); got != "cmd" {
		t.Errorf("an unknown key with Command held was named %q, not %q", got, "cmd")
	}
	if got := macAccelerator(999, 0, 0, false); got != "" {
		t.Errorf("an unknown key was named %q, not %q", got, "")
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
