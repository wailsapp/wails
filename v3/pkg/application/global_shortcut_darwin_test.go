//go:build darwin && !ios && !server

package application

import "testing"

// TestMacKeyCodeForLetters checks that every letter resolves to a valid virtual
// key code on the host's active layout (either translated via UCKeyTranslate or
// via the positional fallback) and that the result is in range.
func TestMacKeyCodeForLetters(t *testing.T) {
	for c := byte('a'); c <= 'z'; c++ {
		key := string(c)
		code, err := macKeyCodeFor(key)
		if err != nil {
			t.Fatalf("macKeyCodeFor(%q) returned error: %v", key, err)
		}
		if code < 0 || code > 127 {
			t.Fatalf("macKeyCodeFor(%q) = %d, want a virtual key code in [0,127]", key, code)
		}
	}
}

// TestMacKeyCodeForPositional checks that non-letter keys use the fixed
// positional table rather than layout translation. In particular digits are
// positional: translating them would mis-bind to the numeric keypad on layouts
// whose number row is shifted (for example AZERTY).
func TestMacKeyCodeForPositional(t *testing.T) {
	cases := map[string]int{
		"f5":    96,
		"left":  123,
		"space": 49,
		"1":     18,
	}
	for key, want := range cases {
		code, err := macKeyCodeFor(key)
		if err != nil {
			t.Fatalf("macKeyCodeFor(%q) returned error: %v", key, err)
		}
		if code != want {
			t.Fatalf("macKeyCodeFor(%q) = %d, want %d", key, code, want)
		}
	}
}

func TestMacKeyCodeForUnknown(t *testing.T) {
	if _, err := macKeyCodeFor("nosuchkey"); err == nil {
		t.Fatal("macKeyCodeFor(\"nosuchkey\") = nil error, want an error")
	}
}
