//go:build darwin && purego && !ios && !server

package application

import (
	"testing"
)

// The event helper and macAccelerator must name a press identically: the
// helper only reduces AppKit's inputs to macAccelerator's, so any difference
// here is a bug in that reduction.
func TestKeyEquivalentAcceleratorAgreesWithMacAccelerator(t *testing.T) {
	cases := []struct {
		name       string
		keyCode    uint16
		modifiers  uint
		characters string
		want       string
	}{
		{"Ctrl+Tab", 48, macModifierControl, "\t", "ctrl+tab"},
		{"Cmd+Shift+K", 40, macModifierCommand | macModifierShift, "K", "shift+cmd+k"},
		{"F12", 111, 0, "", "f12"},
		{"F21 through its character", 0xFFFF, 0, "", "f21"},
		{"unnamed key with Command", 999, macModifierCommand, "", ""},
		{"unnamed key with no characters", 999, 0, "", ""},
		{"Option-A produces a letter", 0, macModifierOption, "å", "option+a"},
		{"Control-A produces U+0001", 0, macModifierControl, "\x01", "ctrl+a"},
		{"Control-Space produces NUL", 49, macModifierControl, "\x00", "ctrl+space"},
		{"Option-Shift-Enter", 36, macModifierOption | macModifierShift, "\r", "shift+option+enter"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := keyEquivalentAccelerator(c.keyCode, c.modifiers, c.characters)
			if got != c.want {
				t.Fatalf("keyEquivalentAccelerator = %q, want %q", got, c.want)
			}

			// Reduce the characters the way the helper must: first UTF-16 unit,
			// and no character at all for an empty string.
			var character rune
			hasCharacter := c.characters != ""
			for _, r := range c.characters {
				character = r
				break
			}
			if direct := macAccelerator(c.keyCode, c.modifiers, character, hasCharacter); direct != got {
				t.Fatalf("macAccelerator = %q, keyEquivalentAccelerator = %q: they disagree", direct, got)
			}
		})
	}
}

// Bindings are stored under the normalised form (modifiers sorted, key upper
// case, darwin modifier names). The lookup must be asked in that form, or a
// binding written in another order can never match.
func TestKeyEquivalentConsumedAsksTheNormalisedForm(t *testing.T) {
	cases := []struct {
		name         string
		accelerator  string
		bindings     map[string]bool
		wantConsumed bool
		wantAsked    string // empty: the lookup must not be asked at all
	}{
		{"bound Ctrl+Tab is consumed", "ctrl+tab", map[string]bool{"Ctrl+TAB": true}, true, "Ctrl+TAB"},
		{"modifiers are normalised", "shift+cmd+k", map[string]bool{"Cmd+Shift+K": true}, true, "Cmd+Shift+K"},
		{"unbound press is not consumed", "cmd+k", map[string]bool{"Ctrl+TAB": true}, false, "Cmd+K"},
		{"no bindings consumes nothing", "ctrl+tab", map[string]bool{}, false, "Ctrl+TAB"},
		{"unnamed press is not looked up", "", map[string]bool{"Ctrl+TAB": true}, false, ""},
		{"unparseable name is not looked up", "ctrl+notakey", map[string]bool{"Ctrl+NOTAKEY": true}, false, ""},
		{"unknown modifier is not looked up", "hyper+k", map[string]bool{"Hyper+K": true}, false, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			asked := ""
			lookup := func(s string) bool {
				asked = s
				return c.bindings[s]
			}
			if got := keyEquivalentConsumed(c.accelerator, lookup); got != c.wantConsumed {
				t.Fatalf("consumed = %v, want %v", got, c.wantConsumed)
			}
			if asked != c.wantAsked {
				t.Fatalf("lookup was asked %q, want %q", asked, c.wantAsked)
			}
		})
	}
}

// AppKit-level: the registered window class must answer performKeyEquivalent:,
// or WKWebView would swallow Ctrl+Tab before any binding ran. Checked against
// the method table, so a registration that drops the entry fails here.
func TestWebviewWindowClassAnswersKeyEquivalents(t *testing.T) {
	loadFrameworks()
	cls := registerWebviewWindowClass()
	if !get[bool](cls, "isSubclassOfClass:", class("NSWindow")) {
		t.Fatal("WailsWebviewNSWindow is not an NSWindow subclass")
	}
	for _, sel := range []string{"performKeyEquivalent:", "keyDown:"} {
		if !get[bool](cls, "instancesRespondToSelector:", sel_(sel)) {
			t.Errorf("WailsWebviewNSWindow does not respond to %s", sel)
		}
	}
}
