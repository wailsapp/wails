package application

import (
	"encoding/binary"
	"testing"
)

// Guards #5794: custom button labels must map back to their own callbacks.
func TestPlanTaskDialogMapsButtons(t *testing.T) {
	cancel := &Button{Label: "Cancel", IsCancel: true}
	confirm := &Button{Label: "Confirm", IsDefault: true}
	p := planTaskDialog([]*Button{cancel, confirm})

	if len(p.buttons) != 2 || p.buttons[0].Label != "Cancel" || p.buttons[1].Label != "Confirm" {
		t.Fatalf("unexpected buttons: %+v", p.buttons)
	}
	if p.buttons[0].ID == p.buttons[1].ID || p.buttons[0].ID < 12 {
		t.Errorf("bad IDs: %+v", p.buttons)
	}
	if p.defaultID != p.buttons[1].ID {
		t.Errorf("default = %d, want %d", p.defaultID, p.buttons[1].ID)
	}
	if p.buttonFor(p.buttons[0].ID) != cancel || p.buttonFor(p.buttons[1].ID) != confirm {
		t.Error("button IDs do not map back to buttons")
	}
	if p.buttonFor(idCancel) != cancel {
		t.Error("dismissing the dialog should activate the cancel button")
	}
	if p.buttonFor(12345) != nil {
		t.Error("unknown ID should map to nothing")
	}
}

func TestPlanTaskDialogNoCancelButton(t *testing.T) {
	p := planTaskDialog([]*Button{{Label: "A"}, {Label: "B"}})
	if p.buttonFor(idCancel) != nil || p.defaultID != 0 {
		t.Error("expected no cancel mapping and no default")
	}
}

func TestMessageBoxButtonFor(t *testing.T) {
	cancel := &Button{Label: "Nope", IsCancel: true}
	ok := &Button{Label: "Go"}
	// Custom captions cannot be shown by MessageBox: Yes -> first, No -> second.
	if got := messageBoxButtonFor([]*Button{ok, cancel}, 6); got != ok {
		t.Error("IDYES should map to the first button")
	}
	if got := messageBoxButtonFor([]*Button{ok, cancel}, 7); got != cancel {
		t.Error("IDNO should map to the second button")
	}
	if got := messageBoxButtonFor([]*Button{ok, cancel}, 2); got != cancel {
		t.Error("IDCANCEL should map to the cancel button")
	}
	// A leading custom cancel button must not turn No into confirmation.
	if got := messageBoxButtonFor([]*Button{cancel, ok}, 7); got != cancel {
		t.Error("IDNO should prefer the cancel role regardless of order")
	}
	if got := messageBoxButtonFor([]*Button{cancel, ok}, 6); got != ok {
		t.Error("IDYES should skip the cancel button")
	}
	// Standard captions still match by label, case-insensitively.
	yes := &Button{Label: "yes"}
	no := &Button{Label: "No"}
	if got := messageBoxButtonFor([]*Button{no, yes}, 6); got != yes {
		t.Error("IDYES should match the Yes label")
	}
	okBtn := &Button{Label: "OK"}
	if got := messageBoxButtonFor([]*Button{okBtn}, 1); got != okBtn {
		t.Error("IDOK should match an OK label")
	}
}

func TestEncodeTaskDialogConfigLayout(t *testing.T) {
	c := taskDialogConfig{hwndParent: 0x11, hInstance: 0x55, mainIcon: 3, callback: 0x66, flags: 0x1008, title: 0x22, content: 0x33, buttonCount: 2, buttons: 0x44, defaultButton: 1001}
	b := encodeTaskDialogConfig(8, c)
	if len(b) != 160 || binary.LittleEndian.Uint32(b) != 160 {
		t.Fatalf("amd64 size = %d", len(b))
	}
	checks := map[int]uint64{4: 0x11, 12: 0x55, 36: 3, 140: 0x66, 20: 0x1008, 28: 0x22, 52: 0x33, 60: 2, 64: 0x44, 72: 1001}
	for off, want := range checks {
		var got uint64
		switch off {
		case 20, 60, 72:
			got = uint64(binary.LittleEndian.Uint32(b[off:]))
		default:
			got = binary.LittleEndian.Uint64(b[off:])
		}
		if got != want {
			t.Errorf("offset %d = %#x, want %#x", off, got, want)
		}
	}
	b32 := encodeTaskDialogConfig(4, c)
	if len(b32) != 96 {
		t.Fatalf("386 size = %d, want 96", len(b32))
	}
	for off, want := range map[int]uint32{8: 0x55, 24: 3, 84: 0x66} {
		if got := binary.LittleEndian.Uint32(b32[off:]); got != want {
			t.Errorf("386 offset %d = %#x, want %#x", off, got, want)
		}
	}
	btn := encodeTaskDialogButtons(8, []int32{1000, 1001}, []uint64{0xA, 0xB})
	if len(btn) != 24 || binary.LittleEndian.Uint32(btn[12:]) != 1001 || binary.LittleEndian.Uint64(btn[16:]) != 0xB {
		t.Errorf("bad button layout: %v", btn)
	}
}
