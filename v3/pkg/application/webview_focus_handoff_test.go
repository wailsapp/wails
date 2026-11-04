package application

import "testing"

func TestFocusHandoffTabOutAndIn(t *testing.T) {
	var f focusHandoff
	if !f.shouldRestoreWebview() {
		t.Fatal("fresh state must restore webview focus")
	}
	if f.tabOut(false, false) {
		t.Fatal("programmatic request must not be handled")
	}
	if !f.shouldRestoreWebview() {
		t.Fatal("unhandled request must not change state")
	}
	if !f.tabOut(true, false) {
		t.Fatal("Tab out must be handled")
	}
	if f.shouldRestoreWebview() {
		t.Fatal("webview must not steal focus back after Tab out")
	}
	if prev, ok := f.tabIn(false, false); !ok || prev {
		t.Fatalf("Tab in: got prev=%v ok=%v, want false,true", prev, ok)
	}
	if !f.shouldRestoreWebview() {
		t.Fatal("state must clear after Tab in")
	}
	if _, ok := f.tabIn(false, false); ok {
		t.Fatal("Tab with focus not on host must not be consumed")
	}
}

func TestFocusHandoffShiftTabAndReset(t *testing.T) {
	var f focusHandoff
	f.tabOut(false, true)
	if prev, ok := f.tabIn(true, false); !ok || !prev {
		t.Fatalf("Shift+Tab in: got prev=%v ok=%v, want true,true", prev, ok)
	}
	f.tabOut(true, false)
	f.reset()
	if !f.shouldRestoreWebview() {
		t.Fatal("reset (deactivation) must allow focus restore")
	}
}

// Ctrl+Tab shortcuts must not consume the host handoff, including Shift+Ctrl+Tab.
func TestFocusHandoffCtrlTabPreservesKeyBindings(t *testing.T) {
	for _, shift := range []bool{false, true} {
		var f focusHandoff
		f.tabOut(true, false)
		if _, handled := f.tabIn(shift, true); handled {
			t.Fatalf("Ctrl+Tab (shift=%v) must reach key bindings", shift)
		}
		if f.shouldRestoreWebview() {
			t.Fatal("Ctrl+Tab must preserve the handoff")
		}
		if previous, handled := f.tabIn(shift, false); !handled || previous != shift {
			t.Fatalf("plain Tab after Ctrl+Tab: previous=%v handled=%v", previous, handled)
		}
	}
}
