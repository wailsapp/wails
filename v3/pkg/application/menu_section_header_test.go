package application

import "testing"

func TestMenuSectionHeader(t *testing.T) {
	m := NewMenu()
	defer m.Destroy()
	first := m.AddRadio("First", true)
	m.AddSectionHeader("mac-node1")
	second := m.AddRadio("Second", false)
	m.processRadioGroups()
	if !menuSectionHeadersSupported() {
		if len(m.items) != 2 || m.ItemAt(1) != second || len(first.radioGroupMembers) != 2 {
			t.Fatal("unsupported section headers must not insert items or split radio groups")
		}
		return
	}
	heading := m.ItemAt(1)
	if len(m.items) != 3 || heading.itemType != sectionHeader || heading.label != "mac-node1" || !heading.disabled || heading.callback != nil {
		t.Fatalf("unexpected heading: %+v", heading)
	}
	if len(first.radioGroupMembers) != 1 || len(second.radioGroupMembers) != 1 {
		t.Fatal("a visible section header must separate radio groups")
	}
	clone := m.Clone()
	if got := clone.ItemAt(1); got == heading || got.itemType != sectionHeader || got.label != heading.label || !got.disabled {
		t.Fatalf("cloned heading: %+v", got)
	}
}
