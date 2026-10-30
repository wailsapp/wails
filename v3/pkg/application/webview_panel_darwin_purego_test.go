//go:build darwin && purego && !ios && !server

package application

import (
	"testing"

	"github.com/ebitengine/purego/objc"
)

// Expected values are the AppKit SDK constants (NSWindow.h):
// Titled 1<<0, Closable 1<<1, Miniaturizable 1<<2, Resizable 1<<3,
// UtilityWindow 1<<4, NonactivatingPanel 1<<7, FullSizeContentView 1<<15.
const (
	testStandardMask = 1<<0 | 1<<1 | 1<<2 | 1<<3 // 15
	testUtilityBit   = 1 << 4                    // 16
	testNonactBit    = 1 << 7                    // 128
	testFullSizeBit  = 1 << 15                   // 32768
)

func TestPanelStyleMask(t *testing.T) {
	tests := []struct {
		name  string
		base  uint
		prefs MacPanelPreferences
		want  uint
	}{
		{name: "no flags keeps the base mask", base: testStandardMask, want: testStandardMask},
		{name: "borderless base stays borderless without flags", base: 0, want: 0},
		{name: "non-activating adds NonactivatingPanel", base: testStandardMask,
			prefs: MacPanelPreferences{NonActivating: true}, want: testStandardMask | testNonactBit},
		{name: "utility adds UtilityWindow", base: testStandardMask,
			prefs: MacPanelPreferences{UtilityWindow: true}, want: testStandardMask | testUtilityBit},
		{name: "both flags", base: testStandardMask,
			prefs: MacPanelPreferences{NonActivating: true, UtilityWindow: true},
			want:  testStandardMask | testNonactBit | testUtilityBit},
		{name: "panel flags do not clear other base bits", base: testStandardMask | testFullSizeBit,
			prefs: MacPanelPreferences{NonActivating: true}, want: testStandardMask | testFullSizeBit | testNonactBit},
		{name: "FloatingPanel and BecomesKeyOnlyIfNeeded set no mask bits", base: testStandardMask,
			prefs: MacPanelPreferences{FloatingPanel: true, BecomesKeyOnlyIfNeeded: true}, want: testStandardMask},
		{name: "applying the flags twice is idempotent", base: testStandardMask | testNonactBit,
			prefs: MacPanelPreferences{NonActivating: true}, want: testStandardMask | testNonactBit},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := panelStyleMask(tt.base, tt.prefs); got != tt.want {
				t.Fatalf("panelStyleMask(%#x, %+v) = %#x, want %#x", tt.base, tt.prefs, got, tt.want)
			}
		})
	}
}

func TestWebviewPanelClassIsNSPanelSubclass(t *testing.T) {
	cls := registerWebviewPanelClass()
	if cls.isNil() {
		t.Fatal("registerWebviewPanelClass returned nil")
	}
	if !get[bool](cls, "isSubclassOfClass:", class("NSPanel")) {
		t.Fatal("WailsWebviewNSPanel is not an NSPanel subclass")
	}
	if !get[bool](cls, "isSubclassOfClass:", class("NSWindow")) {
		t.Fatal("WailsWebviewNSPanel is not an NSWindow subclass")
	}
	if registerWebviewPanelClass() != cls {
		t.Fatal("registerWebviewPanelClass is not cached")
	}
}

func TestWebviewPanelClassRespondsToOverriddenSelectors(t *testing.T) {
	cls := registerWebviewPanelClass()
	for _, sel := range []string{"sendEvent:", "canBecomeMainWindow", "canBecomeKeyWindow", "keyDown:", "cancelOperation:"} {
		if !get[bool](cls, "instancesRespondToSelector:", sel_(sel)) {
			t.Errorf("WailsWebviewNSPanel does not respond to %s", sel)
		}
	}
}

func TestPanelMethodSetOverridesCanBecomeMainWindow(t *testing.T) {
	// The window set answers YES for canBecomeMainWindow; the panel must
	// register exactly one entry, the NO override. A duplicate would be
	// silently ignored by the runtime, so count entries rather than trust order.
	count := 0
	var overrideIsPanelFn bool
	for _, m := range panelMethodSet() {
		if m.Cmd == sel_("canBecomeMainWindow") {
			count++
			fn, ok := m.Fn.(func(objc.ID, objc.SEL) bool)
			if !ok {
				t.Fatalf("canBecomeMainWindow IMP has type %T", m.Fn)
			}
			overrideIsPanelFn = !fn(0, 0)
		}
	}
	if count != 1 {
		t.Fatalf("panelMethodSet has %d canBecomeMainWindow entries, want 1", count)
	}
	if !overrideIsPanelFn {
		t.Fatal("panel canBecomeMainWindow returns YES, want NO")
	}
}
