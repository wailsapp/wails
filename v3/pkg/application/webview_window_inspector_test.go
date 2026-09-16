package application

import (
	"testing"
	"time"
)

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

func inspectorLabels(controls []*MacInspectorControl) string {
	result := ""
	for index, control := range controls {
		if index > 0 {
			result += ","
		}
		result += control.label
	}
	return result
}

func TestMacInspectorExtendedControlKinds(t *testing.T) {
	inspector := NewMacInspector()
	section := inspector.AddSection("Appearance")
	opacity := section.AddSlider("Opacity", 1, 0, 5)
	size := section.AddStepper("Size", 8, 72, 0, 12)
	align := section.AddSegmented("Align", []string{"Left", "Center", "Right"}, 7)
	tint := section.AddColorWell("Tint", NewRGB(10, 20, 30))
	when := time.Date(2026, 9, 16, 12, 0, 0, 0, time.UTC)
	due := section.AddDatePicker("Due", when)
	reset := section.AddButton("Reset")

	kinds := []MacInspectorControlKind{opacity.Kind(), size.Kind(), align.Kind(), tint.Kind(), due.Kind(), reset.Kind()}
	want := []MacInspectorControlKind{MacInspectorSlider, MacInspectorStepper, MacInspectorSegmented,
		MacInspectorColorWell, MacInspectorDatePicker, MacInspectorButton}
	for index := range want {
		if kinds[index] != want[index] {
			t.Fatalf("control %d kind = %d, want %d", index, kinds[index], want[index])
		}
	}
	if min, max := opacity.Range(); min != 0 || max != 1 || opacity.FloatValue() != 1 {
		t.Fatal("slider ranges should normalize and clamp the initial value")
	}
	if size.step != 1 {
		t.Fatal("non-positive steps should fall back to 1")
	}
	size.SetStep(4).SetStep(-1)
	if size.step != 4 {
		t.Fatal("SetStep should ignore non-positive steps")
	}
	if align.SelectedIndex() != 0 {
		t.Fatal("out-of-range segmented selections should normalize to the first segment")
	}
	if tint.Color() != NewRGB(10, 20, 30) || !due.Date().Equal(when) {
		t.Fatal("colour and date should be retained")
	}
	opacity.SetFloatValue(9)
	if opacity.FloatValue() != 1 {
		t.Fatal("SetFloatValue should clamp to the range")
	}
	opacity.SetRange(0, 100).SetFloatValue(50)
	if opacity.FloatValue() != 50 {
		t.Fatal("SetRange should widen the accepted values")
	}
	align.SetOptions([]string{"Top", "Bottom"}).SetSelectedIndex(1)
	if align.SelectedIndex() != 1 {
		t.Fatal("segmented controls should accept SetOptions and SetSelectedIndex")
	}

	label := inspector.AddSection("Other").AddLabel("Words", "1")
	label.SetFloatValue(3).SetColor(NewRGB(1, 1, 1)).SetDate(when).SetSelectedIndex(0).SetChecked(true).SetStep(2)
	if label.FloatValue() != 0 || label.Color() != (RGBA{}) || !label.Date().IsZero() ||
		label.SelectedIndex() != -1 || label.Checked() {
		t.Fatal("kind-specific setters must be no-ops on other kinds")
	}
	label.OnValueChange(func(*Context, float64) {}).OnClick(func(*Context) {}).
		OnColorChange(func(*Context, RGBA) {}).OnDateChange(func(*Context, time.Time) {})
	if label.onValueChange != nil || label.onClick != nil || label.onColorChange != nil || label.onDateChange != nil {
		t.Fatal("kind-specific callbacks must be ignored on other kinds")
	}

	snapshot := inspector.snapshot().sections[0].controls
	if snapshot[0].minimum != 0 || snapshot[0].maximum != 100 || snapshot[0].number != 50 ||
		snapshot[1].step != 4 || snapshot[2].options[1] != "Bottom" || snapshot[3].colour.Blue != 30 ||
		!snapshot[4].date.Equal(when) || snapshot[5].label != "Reset" {
		t.Fatalf("snapshot should carry numeric, colour, and date state: %#v", snapshot)
	}
}

func TestMacInspectorExtendedCallbacks(t *testing.T) {
	inspector := NewMacInspector()
	section := inspector.AddSection("Appearance")
	opacity := section.AddSlider("Opacity", 0, 1, 0.5)
	size := section.AddStepper("Size", 8, 72, 2, 12)
	align := section.AddSegmented("Align", []string{"Left", "Right"}, 0)
	tint := section.AddColorWell("Tint", NewRGB(0, 0, 0))
	due := section.AddDatePicker("Due", time.Unix(0, 0))
	reset := section.AddButton("Reset")
	inspector.registerControls()
	t.Cleanup(inspector.markDead)

	var numbers []float64
	gotIndex := -1
	var gotColor RGBA
	var gotDate time.Time
	clicks := 0
	opacity.OnValueChange(func(_ *Context, value float64) { numbers = append(numbers, value) })
	size.OnValueChange(func(_ *Context, value float64) { numbers = append(numbers, value) })
	align.OnSelectionChange(func(_ *Context, index int, _ string) { gotIndex = index })
	tint.OnColorChange(func(_ *Context, colour RGBA) { gotColor = colour })
	due.OnDateChange(func(_ *Context, date time.Time) { gotDate = date })
	reset.OnClick(func(*Context) { clicks++ })

	when := time.Unix(1_800_000_000, 0)
	handleMacInspectorControlEvent(macInspectorControlEvent{controlID: opacity.internalID, kind: MacInspectorSlider, number: 0.75})
	handleMacInspectorControlEvent(macInspectorControlEvent{controlID: size.internalID, kind: MacInspectorStepper, number: 500})
	handleMacInspectorControlEvent(macInspectorControlEvent{controlID: align.internalID, kind: MacInspectorSegmented, selected: 1})
	handleMacInspectorControlEvent(macInspectorControlEvent{controlID: align.internalID, kind: MacInspectorPopup, selected: 0})
	handleMacInspectorControlEvent(macInspectorControlEvent{controlID: tint.internalID, kind: MacInspectorColorWell, colour: NewRGB(1, 2, 3)})
	handleMacInspectorControlEvent(macInspectorControlEvent{controlID: due.internalID, kind: MacInspectorDatePicker, date: when})
	handleMacInspectorControlEvent(macInspectorControlEvent{controlID: reset.internalID, kind: MacInspectorButton})

	if len(numbers) != 2 || numbers[0] != 0.75 || numbers[1] != 72 || opacity.FloatValue() != 0.75 || size.FloatValue() != 72 {
		t.Fatalf("numeric callbacks = %v", numbers)
	}
	if gotIndex != 1 || align.SelectedIndex() != 1 {
		t.Fatal("segmented events should update the handle and ignore mismatched kinds")
	}
	if gotColor != NewRGB(1, 2, 3) || tint.Color() != gotColor {
		t.Fatal("colour events should update the handle before dispatch")
	}
	if !gotDate.Equal(when) || !due.Date().Equal(when) {
		t.Fatal("date events should update the handle before dispatch")
	}
	if clicks != 1 {
		t.Fatalf("button clicks = %d", clicks)
	}
}

func TestMacInspectorSectionsCollapseRemoveMove(t *testing.T) {
	inspector := NewMacInspector()
	document := inspector.AddSection("Document")
	title := document.AddTextField("Title", "A note")
	pinned := document.AddCheckbox("Pinned", false)
	words := document.AddLabel("Words", "1")
	statistics := inspector.AddSection("Statistics")
	inspector.registerControls()
	t.Cleanup(inspector.markDead)

	document.SetCollapsed(true)
	if document.IsCollapsed() {
		t.Fatal("sections that are not collapsible ignore SetCollapsed")
	}
	document.SetCollapsible(true).SetCollapsed(true)
	snapshot := inspector.snapshot().sections[0]
	if !document.IsCollapsed() || !snapshot.collapsible || !snapshot.collapsed {
		t.Fatal("collapsible state should reach the snapshot")
	}
	handleMacInspectorSectionEvent(macInspectorSectionEvent{sectionID: document.internalID, collapsed: false})
	if document.IsCollapsed() {
		t.Fatal("native disclosure clicks should update the section")
	}
	document.SetCollapsed(true).SetCollapsible(false)
	if document.IsCollapsed() {
		t.Fatal("turning collapsible off should expand the section")
	}

	document.Move(words, 0)
	if got := inspectorLabels(document.Controls()); got != "Words,Title,Pinned" {
		t.Fatalf("Move to front: %s", got)
	}
	document.Move(title, 99)
	if got := inspectorLabels(document.Controls()); got != "Words,Pinned,Title" {
		t.Fatalf("Move should clamp to the end: %s", got)
	}

	title.OnTextChange(func(*Context, string) { t.Fatal("callbacks on removed controls must not fire") })
	document.Remove(title)
	if macInspectorControlRegistry[title.internalID] != nil || len(document.Controls()) != 2 {
		t.Fatal("Remove should unregister the control and drop it from the section")
	}
	handleMacInspectorControlEvent(macInspectorControlEvent{controlID: title.internalID, kind: MacInspectorTextField, value: "x"})
	title.SetValue("Late")
	if title.Value() != "A note" {
		t.Fatal("removed controls should be inert")
	}

	inspector.MoveSection(statistics, 0)
	if sections := inspector.Sections(); sections[0] != statistics || sections[1] != document {
		t.Fatal("MoveSection should reorder sections")
	}
	inspector.RemoveSection(document)
	if len(inspector.Sections()) != 1 || macInspectorControlRegistry[pinned.internalID] != nil ||
		macInspectorSectionRegistry[document.internalID] != nil || document.AddLabel("Late", "x") != nil {
		t.Fatal("RemoveSection should unregister its controls and make the section inert")
	}
}
