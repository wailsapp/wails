package application

import "testing"

func newTestInspector() (*MacInspector, *MacInspectorControl, *MacInspectorControl, *MacInspectorControl) {
	inspector := NewMacInspector()
	document := inspector.AddSection("Document")
	title := document.AddTextField("Title", "A note")
	pinned := document.AddCheckbox("Pinned", false)
	category := document.AddPopup("Category", []string{"Personal", "Work"}, 0)
	inspector.AddSection("Statistics").AddLabel("Words", "12")
	return inspector, title, pinned, category
}

func TestMacInspectorBuildsNativeControlsWithoutUserIDs(t *testing.T) {
	inspector, title, pinned, category := newTestInspector()
	snapshot := inspector.snapshot()
	if len(snapshot.sections) != 2 || len(snapshot.sections[0].controls) != 3 {
		t.Fatalf("unexpected inspector structure: %#v", snapshot)
	}
	if title.internalID == 0 || pinned.internalID == 0 || category.internalID == 0 ||
		title.internalID == pinned.internalID || pinned.internalID == category.internalID {
		t.Fatal("inspector controls should receive unique generated internal IDs")
	}
	if snapshot.sections[0].controls[2].selected != 0 ||
		snapshot.sections[0].controls[2].options[1] != "Work" {
		t.Fatal("popup options and selection should be retained in the snapshot")
	}
}

func TestMacInspectorLiveStateAndDefensiveCopies(t *testing.T) {
	inspector, title, pinned, category := newTestInspector()
	options := []string{"One", "Two", "Three"}
	title.SetValue("Renamed").SetTooltip("Edit title").SetEnabled(false).SetHidden(true)
	pinned.SetChecked(true)
	category.SetOptions(options).SetSelectedIndex(2)
	options[2] = "Mutated"

	snapshot := inspector.snapshot().sections[0].controls
	if snapshot[0].value != "Renamed" || snapshot[0].tooltip != "Edit title" ||
		!snapshot[0].disabled || !snapshot[0].hidden {
		t.Fatalf("unexpected text field snapshot: %#v", snapshot[0])
	}
	if !snapshot[1].checked {
		t.Fatal("checkbox state was not retained")
	}
	if snapshot[2].selected != 2 || snapshot[2].options[2] != "Three" {
		t.Fatal("popup should copy caller-owned option slices")
	}
}

func TestMacInspectorCallbacksUpdateHandles(t *testing.T) {
	inspector, title, pinned, category := newTestInspector()
	inspector.registerControls()
	t.Cleanup(func() {
		for _, control := range inspector.controlHandles() {
			unregisterMacInspectorControl(control.internalID)
		}
	})

	var changedText string
	var changedToggle bool
	var changedIndex int
	var changedSelection string
	title.OnTextChange(func(_ *Context, value string) { changedText = value })
	pinned.OnToggle(func(_ *Context, value bool) { changedToggle = value })
	category.OnSelectionChange(func(_ *Context, index int, value string) {
		changedIndex, changedSelection = index, value
	})

	handleMacInspectorControlEvent(macInspectorControlEvent{
		controlID: title.internalID, kind: MacInspectorTextField, value: "Native edit",
	})
	handleMacInspectorControlEvent(macInspectorControlEvent{
		controlID: pinned.internalID, kind: MacInspectorCheckbox, checked: true,
	})
	handleMacInspectorControlEvent(macInspectorControlEvent{
		controlID: category.internalID, kind: MacInspectorPopup, selected: 1,
	})

	if title.Value() != "Native edit" || changedText != "Native edit" {
		t.Fatal("text callback should update the handle before dispatch")
	}
	if !pinned.Checked() || !changedToggle {
		t.Fatal("toggle callback should update the handle before dispatch")
	}
	if category.SelectedIndex() != 1 || changedIndex != 1 || changedSelection != "Work" {
		t.Fatal("selection callback should include the current index and value")
	}
}

func TestMacInspectorDeadHandlesAreSafe(t *testing.T) {
	inspector, title, pinned, category := newTestInspector()
	inspector.registerControls()
	inspector.markDead()
	title.SetValue("Late").OnTextChange(func(*Context, string) {})
	pinned.SetChecked(true).OnToggle(func(*Context, bool) {})
	category.SetSelectedIndex(1).OnSelectionChange(func(*Context, int, string) {})
	if !inspector.dead || macInspectorControlRegistry[title.internalID] != nil {
		t.Fatal("teardown should invalidate native inspector controls and callbacks")
	}
}
