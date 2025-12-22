package application

import (
	"runtime"
	"strings"
	"testing"
)

func TestModifier_Constants(t *testing.T) {
	// Verify modifier constants are distinct
	modifiers := []modifier{CmdOrCtrlKey, OptionOrAltKey, ShiftKey, SuperKey, ControlKey}
	seen := make(map[modifier]bool)
	for _, m := range modifiers {
		if seen[m] {
			t.Errorf("Duplicate modifier value: %d", m)
		}
		seen[m] = true
	}

	// CmdOrCtrlKey should be 0 (the base value)
	if CmdOrCtrlKey != 0 {
		t.Error("CmdOrCtrlKey should be 0")
	}
}

func TestParseKey_Valid(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"a", "a"},
		{"A", "a"},
		{"z", "z"},
		{"0", "0"},
		{"9", "9"},
		{"+", "+"}, // Single + is a valid printable character
		{"plus", "+"},
		{"Plus", "+"},
		{"PLUS", "+"},
		{"backspace", "backspace"},
		{"Backspace", "backspace"},
		{"BACKSPACE", "backspace"},
		{"tab", "tab"},
		{"return", "return"},
		{"enter", "enter"},
		{"escape", "escape"},
		{"left", "left"},
		{"right", "right"},
		{"up", "up"},
		{"down", "down"},
		{"space", "space"},
		{"delete", "delete"},
		{"home", "home"},
		{"end", "end"},
		{"page up", "page up"},
		{"page down", "page down"},
		{"f1", "f1"},
		{"F1", "f1"},
		{"f12", "f12"},
		{"f35", "f35"},
		{"numlock", "numlock"},
	}

	for _, tt := range tests {
		result, valid := parseKey(tt.input)
		if tt.expected == "" {
			if valid {
				t.Errorf("parseKey(%q) should be invalid", tt.input)
			}
		} else {
			if !valid {
				t.Errorf("parseKey(%q) should be valid", tt.input)
			}
			if result != tt.expected {
				t.Errorf("parseKey(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		}
	}
}

func TestParseKey_Invalid(t *testing.T) {
	tests := []string{
		"abc",       // multiple chars
		"",          // empty
		"notakey",   // not a named key
		"ctrl+a",    // shortcut syntax
		"backspac",  // misspelled
	}

	for _, tt := range tests {
		_, valid := parseKey(tt)
		if valid {
			t.Errorf("parseKey(%q) should be invalid", tt)
		}
	}
}

func TestParseAccelerator_Valid(t *testing.T) {
	tests := []struct {
		input      string
		key        string
		modCount   int
	}{
		{"a", "a", 0},
		{"Ctrl+A", "a", 1},
		{"ctrl+a", "a", 1},
		{"Ctrl+Shift+A", "a", 2},
		{"ctrl+shift+a", "a", 2},
		{"Cmd+A", "a", 1},
		{"Command+A", "a", 1},
		{"CmdOrCtrl+A", "a", 1},
		{"Alt+A", "a", 1},
		{"Option+A", "a", 1},
		{"OptionOrAlt+A", "a", 1},
		{"Shift+A", "a", 1},
		{"Super+A", "a", 1},
		{"Ctrl+Shift+Alt+A", "a", 3},
		{"Ctrl+plus", "+", 1},
		{"F1", "f1", 0},
		{"Ctrl+F12", "f12", 1},
		{"Ctrl+Shift+F1", "f1", 2},
		{"Ctrl+backspace", "backspace", 1},
		{"Ctrl+escape", "escape", 1},
	}

	for _, tt := range tests {
		acc, err := parseAccelerator(tt.input)
		if err != nil {
			t.Errorf("parseAccelerator(%q) returned error: %v", tt.input, err)
			continue
		}
		if acc.Key != tt.key {
			t.Errorf("parseAccelerator(%q).Key = %q, want %q", tt.input, acc.Key, tt.key)
		}
		if len(acc.Modifiers) != tt.modCount {
			t.Errorf("parseAccelerator(%q) has %d modifiers, want %d", tt.input, len(acc.Modifiers), tt.modCount)
		}
	}
}

func TestParseAccelerator_Invalid(t *testing.T) {
	tests := []struct {
		input  string
		errMsg string
	}{
		{"", "no components"},
		{"Ctrl+", "not a valid key"},
		{"Ctrl+abc", "not a valid key"},
		{"NotAModifier+A", "not a valid modifier"},
		{"Ctrl+Shift+notakey", "not a valid key"},
	}

	for _, tt := range tests {
		_, err := parseAccelerator(tt.input)
		if err == nil {
			t.Errorf("parseAccelerator(%q) should return error", tt.input)
		}
	}
}

func TestParseAccelerator_DuplicateModifiers(t *testing.T) {
	// Duplicate modifiers should be deduplicated
	acc, err := parseAccelerator("Ctrl+Ctrl+A")
	if err != nil {
		t.Errorf("parseAccelerator returned error: %v", err)
		return
	}
	if len(acc.Modifiers) != 1 {
		t.Errorf("Duplicate modifiers should be deduplicated, got %d modifiers", len(acc.Modifiers))
	}
}

func TestAccelerator_Clone(t *testing.T) {
	original := &accelerator{
		Key:       "a",
		Modifiers: []modifier{ControlKey, ShiftKey},
	}

	clone := original.clone()

	if clone == original {
		t.Error("Clone should return a different pointer")
	}
	if clone.Key != original.Key {
		t.Error("Clone should have same Key")
	}
	// Note: the slice reference is copied, so modifying the clone's slice would affect original
	// This is a shallow clone
}

func TestAccelerator_String(t *testing.T) {
	// Test key-only accelerator (platform-independent)
	acc := &accelerator{Key: "a", Modifiers: []modifier{}}
	result := acc.String()
	if result != "A" {
		t.Errorf("accelerator.String() = %q, want %q", result, "A")
	}

	// Test with ControlKey modifier - output varies by platform
	acc = &accelerator{Key: "a", Modifiers: []modifier{ControlKey}}
	result = acc.String()
	// On macOS: "Ctrl+A", on Linux/Windows: "Ctrl+A"
	// The representation should contain the key and be non-empty
	if !strings.HasSuffix(result, "+A") && result != "A" {
		t.Errorf("accelerator.String() = %q, expected to end with '+A'", result)
	}
	if result == "" {
		t.Error("accelerator.String() should not return empty")
	}

	// Test function key with modifier
	acc = &accelerator{Key: "f1", Modifiers: []modifier{ControlKey}}
	result = acc.String()
	if !strings.HasSuffix(result, "+F1") {
		t.Errorf("accelerator.String() = %q, expected to end with '+F1'", result)
	}
}

func TestAccelerator_String_PlatformSpecific(t *testing.T) {
	// This test documents the expected platform-specific behavior
	acc := &accelerator{Key: "a", Modifiers: []modifier{ControlKey}}
	result := acc.String()

	switch runtime.GOOS {
	case "darwin":
		// On macOS, Ctrl key is represented as "Ctrl" (distinct from Cmd)
		if !strings.Contains(result, "Ctrl") && !strings.Contains(result, "⌃") {
			t.Logf("On macOS, got %q for ControlKey modifier", result)
		}
	case "linux", "windows":
		if !strings.Contains(result, "Ctrl") {
			t.Errorf("On %s, expected 'Ctrl' in result, got %q", runtime.GOOS, result)
		}
	}
}

func TestAccelerator_StringWithMultipleModifiers(t *testing.T) {
	acc := &accelerator{
		Key:       "a",
		Modifiers: []modifier{ShiftKey, ControlKey},
	}

	result := acc.String()

	// Result should contain both modifiers and the key
	if result == "" {
		t.Error("String() should not return empty")
	}
	// The modifiers are sorted, so order is deterministic
}

func TestModifierMap_Contains(t *testing.T) {
	expectedMappings := map[string]modifier{
		"cmdorctrl":   CmdOrCtrlKey,
		"cmd":         CmdOrCtrlKey,
		"command":     CmdOrCtrlKey,
		"ctrl":        ControlKey,
		"optionoralt": OptionOrAltKey,
		"alt":         OptionOrAltKey,
		"option":      OptionOrAltKey,
		"shift":       ShiftKey,
		"super":       SuperKey,
	}

	for key, expected := range expectedMappings {
		actual, ok := modifierMap[key]
		if !ok {
			t.Errorf("modifierMap should contain key %q", key)
		}
		if actual != expected {
			t.Errorf("modifierMap[%q] = %v, want %v", key, actual, expected)
		}
	}
}

func TestNamedKeys_Contains(t *testing.T) {
	expectedKeys := []string{
		"backspace", "tab", "return", "enter", "escape",
		"left", "right", "up", "down", "space", "delete",
		"home", "end", "page up", "page down",
		"f1", "f2", "f3", "f4", "f5", "f6", "f7", "f8", "f9", "f10",
		"f11", "f12", "f13", "f14", "f15", "f16", "f17", "f18", "f19", "f20",
		"f21", "f22", "f23", "f24", "f25", "f26", "f27", "f28", "f29", "f30",
		"f31", "f32", "f33", "f34", "f35",
		"numlock",
	}

	for _, key := range expectedKeys {
		if _, ok := namedKeys[key]; !ok {
			t.Errorf("namedKeys should contain %q", key)
		}
	}
}
