//go:build darwin && purego && !ios && !server

package application

import "testing"

// RemoveAccelerator must clear what AppKit matches on, not only the Go field:
// the NSMenuItem keyEquivalent and its modifier mask. Runs on macOS (AppKit).
func TestMenuItemRemoveAcceleratorClearsKeyEquivalent(t *testing.T) {
	item := NewMenuItem("Save")
	item.impl = newMenuItemImpl(item)
	native := idFromPtr(item.impl.(*macosMenuItem).nsMenuItem)
	t.Cleanup(func() { item.Destroy() })

	item.SetAccelerator("CmdOrCtrl+Shift+s")
	if key := native.send("keyEquivalent").string(); key == "" {
		t.Fatal("SetAccelerator left no keyEquivalent on the NSMenuItem")
	}
	if mask := get[uint](native, "keyEquivalentModifierMask"); mask == 0 {
		t.Fatal("SetAccelerator left no modifier mask on the NSMenuItem")
	}

	item.RemoveAccelerator()
	if key := native.send("keyEquivalent").string(); key != "" {
		t.Fatalf("after RemoveAccelerator keyEquivalent = %q, want empty", key)
	}
	if mask := get[uint](native, "keyEquivalentModifierMask"); mask != 0 {
		t.Fatalf("after RemoveAccelerator modifier mask = %#x, want 0", mask)
	}
	if item.GetAccelerator() != "" {
		t.Fatalf("GetAccelerator after RemoveAccelerator = %q, want empty", item.GetAccelerator())
	}
}
