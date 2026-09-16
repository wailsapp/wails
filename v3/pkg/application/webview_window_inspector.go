package application

import (
	"math"
	"slices"
	"sync"
	"sync/atomic"
	"time"
)

// MacInspector describes a native AppKit property inspector hosted by an
// inspector NSSplitViewItem. Sections and controls are rendered by AppKit;
// adding an inspector does not create another WebView.
//
// Keep the control handles returned by the AddXxx methods to update values or
// register callbacks. Control identifiers are generated internally and are
// never part of the application API.
type MacInspector struct {
	lock sync.RWMutex

	sections []*MacInspectorSection
	pane     *MacSplitPane
	dead     bool
}

// MacInspectorSection groups related controls under a native section heading.
// Collapsible sections show a disclosure triangle beside the heading.
type MacInspectorSection struct {
	lock sync.RWMutex

	internalID  uint64
	label       string
	collapsible bool
	collapsed   bool
	removed     bool
	inspector   *MacInspector
	controls    []*MacInspectorControl
}

// MacInspectorControlKind identifies the native AppKit control used for one
// inspector property.
type MacInspectorControlKind uint8

const (
	// MacInspectorLabel is a read-only value displayed beside a property name.
	MacInspectorLabel MacInspectorControlKind = iota
	// MacInspectorTextField is an editable native text field.
	MacInspectorTextField
	// MacInspectorCheckbox is a native checkbox.
	MacInspectorCheckbox
	// MacInspectorPopup is a native pop-up button containing string options.
	MacInspectorPopup
	// MacInspectorSlider is a native horizontal slider over a float range.
	MacInspectorSlider
	// MacInspectorStepper is a native numeric field paired with a stepper.
	MacInspectorStepper
	// MacInspectorSegmented is a native segmented control selecting one option.
	MacInspectorSegmented
	// MacInspectorColorWell is a native colour well backed by the colour panel.
	MacInspectorColorWell
	// MacInspectorDatePicker is a native date and time picker.
	MacInspectorDatePicker
	// MacInspectorButton is a native push button.
	MacInspectorButton
)

// MacInspectorControl is a live handle for one native inspector control.
// Kind-specific setters are safe no-ops when used with the wrong kind.
type MacInspectorControl struct {
	lock sync.RWMutex

	internalID    uint64
	kind          MacInspectorControlKind
	label         string
	value         string
	checked       bool
	options       []string
	selected      int
	number        float64
	minimum       float64
	maximum       float64
	step          float64
	colour        RGBA
	date          time.Time
	tooltip       string
	disabled      bool
	hidden        bool
	removed       bool
	inspector     *MacInspector
	onTextChange  func(*Context, string)
	onToggle      func(*Context, bool)
	onSelection   func(*Context, int, string)
	onValueChange func(*Context, float64)
	onColorChange func(*Context, RGBA)
	onDateChange  func(*Context, time.Time)
	onClick       func(*Context)
}

var macInspectorNodeID uint64

func nextMacInspectorNodeID() uint64 {
	return atomic.AddUint64(&macInspectorNodeID, 1)
}

// NewMacInspector creates an empty native property inspector.
func NewMacInspector() *MacInspector {
	return &MacInspector{}
}

// AddSection appends a native inspector section.
func (i *MacInspector) AddSection(label string) *MacInspectorSection {
	if i == nil {
		return nil
	}
	section := &MacInspectorSection{
		internalID: nextMacInspectorNodeID(),
		label:      label,
		inspector:  i,
	}
	i.lock.Lock()
	if i.dead {
		i.lock.Unlock()
		return nil
	}
	i.sections = append(i.sections, section)
	i.lock.Unlock()
	macInspectorRegisterSectionIfInstalled(section)
	macInspectorApplySnapshot(i)
	return section
}

// Sections returns the inspector's sections in display order.
func (i *MacInspector) Sections() []*MacInspectorSection {
	if i == nil {
		return nil
	}
	i.lock.RLock()
	defer i.lock.RUnlock()
	return append([]*MacInspectorSection(nil), i.sections...)
}

// RemoveSection detaches a section and its controls, unregistering their
// callbacks. The handles are inert afterwards.
func (i *MacInspector) RemoveSection(section *MacInspectorSection) *MacInspector {
	if i == nil || section == nil || section.inspector != i || section.isDead() {
		return i
	}
	i.lock.Lock()
	if i.dead {
		i.lock.Unlock()
		return i
	}
	for index, candidate := range i.sections {
		if candidate == section {
			i.sections = append(i.sections[:index:index], i.sections[index+1:]...)
			break
		}
	}
	i.lock.Unlock()
	section.lock.Lock()
	section.removed = true
	controls := section.controls
	section.controls = nil
	section.lock.Unlock()
	unregisterMacInspectorSection(section.internalID)
	for _, control := range controls {
		control.retire()
	}
	macInspectorApplySnapshot(i)
	return i
}

// MoveSection reorders a section to index, clamped to the inspector's bounds.
func (i *MacInspector) MoveSection(section *MacInspectorSection, index int) *MacInspector {
	if i == nil || section == nil || section.inspector != i || section.isDead() {
		return i
	}
	i.lock.Lock()
	if i.dead {
		i.lock.Unlock()
		return i
	}
	position := slices.Index(i.sections, section)
	if position < 0 {
		i.lock.Unlock()
		return i
	}
	remaining := append(i.sections[:position:position], i.sections[position+1:]...)
	index = clampMacInspectorIndex(index, len(remaining))
	sections := make([]*MacInspectorSection, 0, len(remaining)+1)
	sections = append(sections, remaining[:index]...)
	sections = append(sections, section)
	sections = append(sections, remaining[index:]...)
	i.sections = sections
	i.lock.Unlock()
	macInspectorApplySnapshot(i)
	return i
}

func clampMacInspectorIndex(index, length int) int {
	if index < 0 {
		return 0
	}
	if index > length {
		return length
	}
	return index
}

// SetLabel updates the native section heading.
func (s *MacInspectorSection) SetLabel(label string) *MacInspectorSection {
	if s == nil || s.isDead() {
		return s
	}
	s.lock.Lock()
	if s.label == label {
		s.lock.Unlock()
		return s
	}
	s.label = label
	s.lock.Unlock()
	macInspectorApplySnapshot(s.inspector)
	return s
}

// SetCollapsible shows a disclosure triangle beside the heading so the user
// can collapse the section's controls. Turning it off also expands the
// section.
func (s *MacInspectorSection) SetCollapsible(collapsible bool) *MacInspectorSection {
	if s == nil || s.isDead() {
		return s
	}
	s.lock.Lock()
	if s.collapsible == collapsible {
		s.lock.Unlock()
		return s
	}
	s.collapsible = collapsible
	if !collapsible {
		s.collapsed = false
	}
	s.lock.Unlock()
	macInspectorApplySnapshot(s.inspector)
	return s
}

// SetCollapsed hides or shows the section's controls. It has no effect on a
// section that is not collapsible.
func (s *MacInspectorSection) SetCollapsed(collapsed bool) *MacInspectorSection {
	if s == nil || s.isDead() {
		return s
	}
	s.lock.Lock()
	if !s.collapsible || s.collapsed == collapsed {
		s.lock.Unlock()
		return s
	}
	s.collapsed = collapsed
	s.lock.Unlock()
	macInspectorApplySnapshot(s.inspector)
	return s
}

// IsCollapsed reports whether the section's controls are hidden.
func (s *MacInspectorSection) IsCollapsed() bool {
	if s == nil {
		return false
	}
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.collapsed
}

// Controls returns the section's controls in display order.
func (s *MacInspectorSection) Controls() []*MacInspectorControl {
	if s == nil {
		return nil
	}
	s.lock.RLock()
	defer s.lock.RUnlock()
	return append([]*MacInspectorControl(nil), s.controls...)
}

// Remove detaches a control from this section and unregisters its callbacks.
// The handle is inert afterwards.
func (s *MacInspectorSection) Remove(control *MacInspectorControl) *MacInspectorSection {
	if s == nil || control == nil || s.isDead() {
		return s
	}
	s.lock.Lock()
	position := slices.Index(s.controls, control)
	if position < 0 {
		s.lock.Unlock()
		return s
	}
	s.controls = append(s.controls[:position:position], s.controls[position+1:]...)
	s.lock.Unlock()
	control.retire()
	macInspectorApplySnapshot(s.inspector)
	return s
}

// Move reorders a control within this section to index, clamped to the
// section's bounds.
func (s *MacInspectorSection) Move(control *MacInspectorControl, index int) *MacInspectorSection {
	if s == nil || control == nil || s.isDead() {
		return s
	}
	s.lock.Lock()
	position := slices.Index(s.controls, control)
	if position < 0 {
		s.lock.Unlock()
		return s
	}
	remaining := append(s.controls[:position:position], s.controls[position+1:]...)
	index = clampMacInspectorIndex(index, len(remaining))
	controls := make([]*MacInspectorControl, 0, len(remaining)+1)
	controls = append(controls, remaining[:index]...)
	controls = append(controls, control)
	controls = append(controls, remaining[index:]...)
	s.controls = controls
	s.lock.Unlock()
	macInspectorApplySnapshot(s.inspector)
	return s
}

// AddLabel adds a read-only property value.
func (s *MacInspectorSection) AddLabel(label, value string) *MacInspectorControl {
	return s.addControl(MacInspectorLabel, label, func(control *MacInspectorControl) {
		control.value = value
	})
}

// AddTextField adds an editable native text field. Use OnTextChange to
// observe edits made by the user.
func (s *MacInspectorSection) AddTextField(label, value string) *MacInspectorControl {
	return s.addControl(MacInspectorTextField, label, func(control *MacInspectorControl) {
		control.value = value
	})
}

// AddCheckbox adds a native checkbox.
func (s *MacInspectorSection) AddCheckbox(label string, checked bool) *MacInspectorControl {
	return s.addControl(MacInspectorCheckbox, label, func(control *MacInspectorControl) {
		control.checked = checked
	})
}

// AddPopup adds a native pop-up button. selectedIndex may be -1 for no
// selection; invalid positive indexes are normalized to the first option.
func (s *MacInspectorSection) AddPopup(label string, options []string, selectedIndex int) *MacInspectorControl {
	return s.addControl(MacInspectorPopup, label, func(control *MacInspectorControl) {
		control.options = append([]string(nil), options...)
		control.selected = normalizeMacInspectorSelection(selectedIndex, len(control.options))
	})
}

// AddSlider adds a native horizontal slider. Use OnValueChange to observe
// the user dragging it; the value is clamped to [min, max].
func (s *MacInspectorSection) AddSlider(label string, min, max, value float64) *MacInspectorControl {
	return s.addControl(MacInspectorSlider, label, func(control *MacInspectorControl) {
		control.minimum, control.maximum = normalizeMacInspectorRange(min, max)
		control.step = 1
		control.number = clampMacInspectorNumber(value, control.minimum, control.maximum)
	})
}

// AddStepper adds a native numeric field with a stepper. step must be
// positive; other values fall back to 1.
func (s *MacInspectorSection) AddStepper(label string, min, max, step, value float64) *MacInspectorControl {
	return s.addControl(MacInspectorStepper, label, func(control *MacInspectorControl) {
		control.minimum, control.maximum = normalizeMacInspectorRange(min, max)
		control.step = normalizeMacInspectorStep(step)
		control.number = clampMacInspectorNumber(value, control.minimum, control.maximum)
	})
}

// AddSegmented adds a native segmented control with one segment per option.
// Use OnSelectionChange, SetSelectedIndex, and SetOptions as with a pop-up.
func (s *MacInspectorSection) AddSegmented(label string, options []string, selected int) *MacInspectorControl {
	return s.addControl(MacInspectorSegmented, label, func(control *MacInspectorControl) {
		control.options = append([]string(nil), options...)
		control.selected = normalizeMacInspectorSelection(selected, len(control.options))
	})
}

// AddColorWell adds a native colour well. Use OnColorChange to observe the
// colour panel.
func (s *MacInspectorSection) AddColorWell(label string, colour RGBA) *MacInspectorControl {
	return s.addControl(MacInspectorColorWell, label, func(control *MacInspectorControl) {
		control.colour = colour
	})
}

// AddDatePicker adds a native date and time picker. A zero time shows the
// current time.
func (s *MacInspectorSection) AddDatePicker(label string, t time.Time) *MacInspectorControl {
	return s.addControl(MacInspectorDatePicker, label, func(control *MacInspectorControl) {
		if t.IsZero() {
			t = time.Now()
		}
		control.date = t
	})
}

// AddButton adds a native push button titled label. Use OnClick to observe
// presses.
func (s *MacInspectorSection) AddButton(label string) *MacInspectorControl {
	return s.addControl(MacInspectorButton, label, nil)
}

func (s *MacInspectorSection) addControl(kind MacInspectorControlKind, label string, configure func(*MacInspectorControl)) *MacInspectorControl {
	if s == nil || s.inspector == nil || s.isDead() {
		return nil
	}
	control := &MacInspectorControl{
		internalID: nextMacInspectorNodeID(),
		kind:       kind,
		label:      label,
		selected:   -1,
		inspector:  s.inspector,
	}
	if configure != nil {
		configure(control)
	}
	s.lock.Lock()
	s.controls = append(s.controls, control)
	s.lock.Unlock()
	macInspectorRegisterControlIfInstalled(control)
	macInspectorApplySnapshot(s.inspector)
	return control
}

func normalizeMacInspectorSelection(index, count int) int {
	if count == 0 || index < 0 {
		return -1
	}
	if index >= count {
		return 0
	}
	return index
}

func normalizeMacInspectorRange(min, max float64) (float64, float64) {
	if math.IsNaN(min) || math.IsInf(min, 0) {
		min = 0
	}
	if math.IsNaN(max) || math.IsInf(max, 0) {
		max = min
	}
	if max < min {
		min, max = max, min
	}
	return min, max
}

func normalizeMacInspectorStep(step float64) float64 {
	if math.IsNaN(step) || math.IsInf(step, 0) || step <= 0 {
		return 1
	}
	return step
}

func clampMacInspectorNumber(value, min, max float64) float64 {
	if math.IsNaN(value) {
		return min
	}
	return math.Min(math.Max(value, min), max)
}

// Kind returns the control's immutable native kind.
func (c *MacInspectorControl) Kind() MacInspectorControlKind {
	if c == nil {
		return MacInspectorLabel
	}
	return c.kind
}

func (c *MacInspectorControl) hasKind(kinds ...MacInspectorControlKind) bool {
	return c != nil && slices.Contains(kinds, c.kind)
}

// SetLabel updates the property's user-visible name.
func (c *MacInspectorControl) SetLabel(label string) *MacInspectorControl {
	if c == nil || c.isDead() {
		return c
	}
	c.lock.Lock()
	if c.label == label {
		c.lock.Unlock()
		return c
	}
	c.label = label
	c.lock.Unlock()
	macInspectorApplyControl(c)
	return c
}

// SetValue updates a label or text field. It does not invoke OnTextChange.
func (c *MacInspectorControl) SetValue(value string) *MacInspectorControl {
	if c == nil || c.isDead() || !c.hasKind(MacInspectorLabel, MacInspectorTextField) {
		return c
	}
	c.lock.Lock()
	if c.value == value {
		c.lock.Unlock()
		return c
	}
	c.value = value
	c.lock.Unlock()
	macInspectorApplyControl(c)
	return c
}

// Value returns the latest programmatic or user-entered string value.
func (c *MacInspectorControl) Value() string {
	if c == nil {
		return ""
	}
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.value
}

// SetChecked updates a checkbox. It does not invoke OnToggle.
func (c *MacInspectorControl) SetChecked(checked bool) *MacInspectorControl {
	if c == nil || c.isDead() || c.kind != MacInspectorCheckbox {
		return c
	}
	c.lock.Lock()
	if c.checked == checked {
		c.lock.Unlock()
		return c
	}
	c.checked = checked
	c.lock.Unlock()
	macInspectorApplyControl(c)
	return c
}

// Checked returns the latest checkbox state.
func (c *MacInspectorControl) Checked() bool {
	if c == nil {
		return false
	}
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.checked
}

// SetOptions replaces a pop-up button's or segmented control's options. The
// current selection is retained when possible and otherwise normalized to the
// first option.
func (c *MacInspectorControl) SetOptions(options []string) *MacInspectorControl {
	if c == nil || c.isDead() || !c.hasKind(MacInspectorPopup, MacInspectorSegmented) {
		return c
	}
	c.lock.Lock()
	if slices.Equal(c.options, options) {
		c.lock.Unlock()
		return c
	}
	c.options = append([]string(nil), options...)
	c.selected = normalizeMacInspectorSelection(c.selected, len(c.options))
	c.lock.Unlock()
	macInspectorApplyControl(c)
	return c
}

// SetSelectedIndex updates a pop-up or segmented selection without invoking
// OnSelectionChange. Invalid indexes are ignored.
func (c *MacInspectorControl) SetSelectedIndex(index int) *MacInspectorControl {
	if c == nil || c.isDead() || !c.hasKind(MacInspectorPopup, MacInspectorSegmented) {
		return c
	}
	c.lock.Lock()
	if index < -1 || index >= len(c.options) {
		c.lock.Unlock()
		return c
	}
	if c.selected == index {
		c.lock.Unlock()
		return c
	}
	c.selected = index
	c.lock.Unlock()
	macInspectorApplyControl(c)
	return c
}

// SelectedIndex returns the latest pop-up or segmented selection, or -1 when
// unselected.
func (c *MacInspectorControl) SelectedIndex() int {
	if c == nil {
		return -1
	}
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.selected
}

// SetFloatValue updates a slider or stepper, clamped to its range. It does
// not invoke OnValueChange.
func (c *MacInspectorControl) SetFloatValue(value float64) *MacInspectorControl {
	if c == nil || c.isDead() || !c.hasKind(MacInspectorSlider, MacInspectorStepper) {
		return c
	}
	c.lock.Lock()
	value = clampMacInspectorNumber(value, c.minimum, c.maximum)
	if c.number == value {
		c.lock.Unlock()
		return c
	}
	c.number = value
	c.lock.Unlock()
	macInspectorApplyControl(c)
	return c
}

// FloatValue returns the latest slider or stepper value.
func (c *MacInspectorControl) FloatValue() float64 {
	if c == nil {
		return 0
	}
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.number
}

// SetRange updates a slider's or stepper's bounds, clamping the current value.
func (c *MacInspectorControl) SetRange(min, max float64) *MacInspectorControl {
	if c == nil || c.isDead() || !c.hasKind(MacInspectorSlider, MacInspectorStepper) {
		return c
	}
	min, max = normalizeMacInspectorRange(min, max)
	c.lock.Lock()
	if c.minimum == min && c.maximum == max {
		c.lock.Unlock()
		return c
	}
	c.minimum, c.maximum = min, max
	c.number = clampMacInspectorNumber(c.number, min, max)
	c.lock.Unlock()
	macInspectorApplyControl(c)
	return c
}

// Range returns a slider's or stepper's bounds.
func (c *MacInspectorControl) Range() (min, max float64) {
	if c == nil {
		return 0, 0
	}
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.minimum, c.maximum
}

// SetStep updates a stepper's increment. Non-positive values are ignored.
func (c *MacInspectorControl) SetStep(step float64) *MacInspectorControl {
	if c == nil || c.isDead() || c.kind != MacInspectorStepper || normalizeMacInspectorStep(step) != step {
		return c
	}
	c.lock.Lock()
	if c.step == step {
		c.lock.Unlock()
		return c
	}
	c.step = step
	c.lock.Unlock()
	macInspectorApplyControl(c)
	return c
}

// SetColor updates a colour well. It does not invoke OnColorChange.
func (c *MacInspectorControl) SetColor(colour RGBA) *MacInspectorControl {
	if c == nil || c.isDead() || c.kind != MacInspectorColorWell {
		return c
	}
	c.lock.Lock()
	if c.colour == colour {
		c.lock.Unlock()
		return c
	}
	c.colour = colour
	c.lock.Unlock()
	macInspectorApplyControl(c)
	return c
}

// Color returns the latest colour well value.
func (c *MacInspectorControl) Color() RGBA {
	if c == nil {
		return RGBA{}
	}
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.colour
}

// SetDate updates a date picker. It does not invoke OnDateChange.
func (c *MacInspectorControl) SetDate(t time.Time) *MacInspectorControl {
	if c == nil || c.isDead() || c.kind != MacInspectorDatePicker || t.IsZero() {
		return c
	}
	c.lock.Lock()
	if c.date.Equal(t) {
		c.lock.Unlock()
		return c
	}
	c.date = t
	c.lock.Unlock()
	macInspectorApplyControl(c)
	return c
}

// Date returns the latest date picker value.
func (c *MacInspectorControl) Date() time.Time {
	if c == nil {
		return time.Time{}
	}
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.date
}

// SetTooltip sets the native control tooltip.
func (c *MacInspectorControl) SetTooltip(tooltip string) *MacInspectorControl {
	if c == nil || c.isDead() {
		return c
	}
	c.lock.Lock()
	if c.tooltip == tooltip {
		c.lock.Unlock()
		return c
	}
	c.tooltip = tooltip
	c.lock.Unlock()
	macInspectorApplyControl(c)
	return c
}

// SetEnabled controls whether an interactive inspector control accepts input.
func (c *MacInspectorControl) SetEnabled(enabled bool) *MacInspectorControl {
	if c == nil || c.isDead() {
		return c
	}
	c.lock.Lock()
	if c.disabled == !enabled {
		c.lock.Unlock()
		return c
	}
	c.disabled = !enabled
	c.lock.Unlock()
	macInspectorApplyControl(c)
	return c
}

// SetHidden includes or removes the property row from native layout.
func (c *MacInspectorControl) SetHidden(hidden bool) *MacInspectorControl {
	if c == nil || c.isDead() {
		return c
	}
	c.lock.Lock()
	if c.hidden == hidden {
		c.lock.Unlock()
		return c
	}
	c.hidden = hidden
	c.lock.Unlock()
	macInspectorApplyControl(c)
	return c
}

// OnTextChange sets the callback for user edits to a text field. Passing nil
// clears it. Programmatic SetValue calls never invoke the callback.
func (c *MacInspectorControl) OnTextChange(callback func(*Context, string)) *MacInspectorControl {
	if c == nil || c.isDead() || c.kind != MacInspectorTextField {
		return c
	}
	c.lock.Lock()
	c.onTextChange = callback
	c.lock.Unlock()
	return c
}

// OnToggle sets the callback for user changes to a checkbox.
func (c *MacInspectorControl) OnToggle(callback func(*Context, bool)) *MacInspectorControl {
	if c == nil || c.isDead() || c.kind != MacInspectorCheckbox {
		return c
	}
	c.lock.Lock()
	c.onToggle = callback
	c.lock.Unlock()
	return c
}

// OnSelectionChange sets the callback for user changes to a pop-up button or
// segmented control. The callback receives both the selected index and value.
func (c *MacInspectorControl) OnSelectionChange(callback func(*Context, int, string)) *MacInspectorControl {
	if c == nil || c.isDead() || !c.hasKind(MacInspectorPopup, MacInspectorSegmented) {
		return c
	}
	c.lock.Lock()
	c.onSelection = callback
	c.lock.Unlock()
	return c
}

// OnValueChange sets the callback for user changes to a slider or stepper.
func (c *MacInspectorControl) OnValueChange(callback func(*Context, float64)) *MacInspectorControl {
	if c == nil || c.isDead() || !c.hasKind(MacInspectorSlider, MacInspectorStepper) {
		return c
	}
	c.lock.Lock()
	c.onValueChange = callback
	c.lock.Unlock()
	return c
}

// OnColorChange sets the callback for user changes to a colour well. It runs
// for every colour panel change while the panel is open.
func (c *MacInspectorControl) OnColorChange(callback func(*Context, RGBA)) *MacInspectorControl {
	if c == nil || c.isDead() || c.kind != MacInspectorColorWell {
		return c
	}
	c.lock.Lock()
	c.onColorChange = callback
	c.lock.Unlock()
	return c
}

// OnDateChange sets the callback for user changes to a date picker.
func (c *MacInspectorControl) OnDateChange(callback func(*Context, time.Time)) *MacInspectorControl {
	if c == nil || c.isDead() || c.kind != MacInspectorDatePicker {
		return c
	}
	c.lock.Lock()
	c.onDateChange = callback
	c.lock.Unlock()
	return c
}

// OnClick sets the callback for presses of a button.
func (c *MacInspectorControl) OnClick(callback func(*Context)) *MacInspectorControl {
	if c == nil || c.isDead() || c.kind != MacInspectorButton {
		return c
	}
	c.lock.Lock()
	c.onClick = callback
	c.lock.Unlock()
	return c
}

func (s *MacInspectorSection) isDead() bool {
	if s == nil || s.inspector == nil {
		return true
	}
	s.lock.RLock()
	removed := s.removed
	s.lock.RUnlock()
	if removed {
		return true
	}
	s.inspector.lock.RLock()
	defer s.inspector.lock.RUnlock()
	return s.inspector.dead
}

func (c *MacInspectorControl) isDead() bool {
	if c == nil || c.inspector == nil {
		return true
	}
	c.lock.RLock()
	removed := c.removed
	c.lock.RUnlock()
	if removed {
		return true
	}
	c.inspector.lock.RLock()
	defer c.inspector.lock.RUnlock()
	return c.inspector.dead
}

// retire marks a control removed, drops its callbacks, and unregisters it.
func (c *MacInspectorControl) retire() {
	if c == nil {
		return
	}
	c.lock.Lock()
	c.removed = true
	c.onTextChange = nil
	c.onToggle = nil
	c.onSelection = nil
	c.onValueChange = nil
	c.onColorChange = nil
	c.onDateChange = nil
	c.onClick = nil
	c.lock.Unlock()
	unregisterMacInspectorControl(c.internalID)
}

type macInspectorControlSnapshot struct {
	internalID uint64
	kind       MacInspectorControlKind
	label      string
	value      string
	checked    bool
	options    []string
	selected   int
	number     float64
	minimum    float64
	maximum    float64
	step       float64
	colour     RGBA
	date       time.Time
	tooltip    string
	disabled   bool
	hidden     bool
}

type macInspectorSectionSnapshot struct {
	internalID  uint64
	label       string
	collapsible bool
	collapsed   bool
	controls    []macInspectorControlSnapshot
}

type macInspectorSnapshot struct {
	sections []macInspectorSectionSnapshot
}

func (i *MacInspector) snapshot() macInspectorSnapshot {
	if i == nil {
		return macInspectorSnapshot{}
	}
	i.lock.RLock()
	sections := append([]*MacInspectorSection(nil), i.sections...)
	i.lock.RUnlock()
	result := macInspectorSnapshot{sections: make([]macInspectorSectionSnapshot, 0, len(sections))}
	for _, section := range sections {
		section.lock.RLock()
		entry := macInspectorSectionSnapshot{
			internalID:  section.internalID,
			label:       section.label,
			collapsible: section.collapsible,
			collapsed:   section.collapsed,
			controls:    make([]macInspectorControlSnapshot, 0, len(section.controls)),
		}
		for _, control := range section.controls {
			entry.controls = append(entry.controls, snapshotMacInspectorControl(control))
		}
		section.lock.RUnlock()
		result.sections = append(result.sections, entry)
	}
	return result
}

func snapshotMacInspectorControl(control *MacInspectorControl) macInspectorControlSnapshot {
	control.lock.RLock()
	defer control.lock.RUnlock()
	return macInspectorControlSnapshot{
		internalID: control.internalID,
		kind:       control.kind,
		label:      control.label,
		value:      control.value,
		checked:    control.checked,
		options:    append([]string(nil), control.options...),
		selected:   control.selected,
		number:     control.number,
		minimum:    control.minimum,
		maximum:    control.maximum,
		step:       control.step,
		colour:     control.colour,
		date:       control.date,
		tooltip:    control.tooltip,
		disabled:   control.disabled,
		hidden:     control.hidden,
	}
}

func (i *MacInspector) controlHandles() []*MacInspectorControl {
	if i == nil {
		return nil
	}
	i.lock.RLock()
	sections := append([]*MacInspectorSection(nil), i.sections...)
	i.lock.RUnlock()
	var result []*MacInspectorControl
	for _, section := range sections {
		section.lock.RLock()
		result = append(result, section.controls...)
		section.lock.RUnlock()
	}
	return result
}

func (i *MacInspector) registerControls() {
	for _, section := range i.Sections() {
		registerMacInspectorSection(section)
	}
	for _, control := range i.controlHandles() {
		registerMacInspectorControl(control)
	}
}

func (i *MacInspector) markDead() {
	if i == nil {
		return
	}
	for _, section := range i.Sections() {
		unregisterMacInspectorSection(section.internalID)
	}
	for _, control := range i.controlHandles() {
		unregisterMacInspectorControl(control.internalID)
		control.lock.Lock()
		control.onTextChange = nil
		control.onToggle = nil
		control.onSelection = nil
		control.onValueChange = nil
		control.onColorChange = nil
		control.onDateChange = nil
		control.onClick = nil
		control.lock.Unlock()
	}
	i.lock.Lock()
	i.dead = true
	i.pane = nil
	i.lock.Unlock()
}

var macInspectorControlRegistry = make(map[uint64]*MacInspectorControl)
var macInspectorControlRegistryLock sync.RWMutex

func registerMacInspectorControl(control *MacInspectorControl) {
	if control == nil {
		return
	}
	macInspectorControlRegistryLock.Lock()
	macInspectorControlRegistry[control.internalID] = control
	macInspectorControlRegistryLock.Unlock()
}

func unregisterMacInspectorControl(id uint64) {
	macInspectorControlRegistryLock.Lock()
	delete(macInspectorControlRegistry, id)
	macInspectorControlRegistryLock.Unlock()
}

var macInspectorSectionRegistry = make(map[uint64]*MacInspectorSection)
var macInspectorSectionRegistryLock sync.RWMutex

func registerMacInspectorSection(section *MacInspectorSection) {
	if section == nil {
		return
	}
	macInspectorSectionRegistryLock.Lock()
	macInspectorSectionRegistry[section.internalID] = section
	macInspectorSectionRegistryLock.Unlock()
}

func unregisterMacInspectorSection(id uint64) {
	macInspectorSectionRegistryLock.Lock()
	delete(macInspectorSectionRegistry, id)
	macInspectorSectionRegistryLock.Unlock()
}

// macInspectorControlEvent is the union of native control callbacks. kind
// selects which payload field is meaningful.
type macInspectorControlEvent struct {
	controlID uint64
	kind      MacInspectorControlKind
	value     string
	checked   bool
	selected  int
	number    float64
	colour    RGBA
	date      time.Time
}

var macInspectorControlEvents = make(chan macInspectorControlEvent, 128)

func handleMacInspectorControlEvent(event macInspectorControlEvent) {
	defer handlePanic()
	macInspectorControlRegistryLock.RLock()
	control := macInspectorControlRegistry[event.controlID]
	macInspectorControlRegistryLock.RUnlock()
	if control == nil || control.isDead() || control.kind != event.kind {
		return
	}

	control.lock.Lock()
	switch event.kind {
	case MacInspectorTextField:
		control.value = event.value
		callback := control.onTextChange
		control.lock.Unlock()
		if callback != nil {
			callback(newContext(), event.value)
		}
	case MacInspectorCheckbox:
		control.checked = event.checked
		callback := control.onToggle
		control.lock.Unlock()
		if callback != nil {
			callback(newContext(), event.checked)
		}
	case MacInspectorPopup, MacInspectorSegmented:
		if event.selected < 0 || event.selected >= len(control.options) {
			control.lock.Unlock()
			return
		}
		control.selected = event.selected
		value := control.options[event.selected]
		callback := control.onSelection
		control.lock.Unlock()
		if callback != nil {
			callback(newContext(), event.selected, value)
		}
	case MacInspectorSlider, MacInspectorStepper:
		value := clampMacInspectorNumber(event.number, control.minimum, control.maximum)
		control.number = value
		callback := control.onValueChange
		control.lock.Unlock()
		if callback != nil {
			callback(newContext(), value)
		}
	case MacInspectorColorWell:
		control.colour = event.colour
		callback := control.onColorChange
		control.lock.Unlock()
		if callback != nil {
			callback(newContext(), event.colour)
		}
	case MacInspectorDatePicker:
		control.date = event.date
		callback := control.onDateChange
		control.lock.Unlock()
		if callback != nil {
			callback(newContext(), event.date)
		}
	case MacInspectorButton:
		callback := control.onClick
		control.lock.Unlock()
		if callback != nil {
			callback(newContext())
		}
	default:
		control.lock.Unlock()
	}
}

// macInspectorSectionEvent reports the user toggling a section's disclosure
// triangle. The native view has already collapsed or expanded.
type macInspectorSectionEvent struct {
	sectionID uint64
	collapsed bool
}

var macInspectorSectionEvents = make(chan macInspectorSectionEvent, 64)

func handleMacInspectorSectionEvent(event macInspectorSectionEvent) {
	defer handlePanic()
	macInspectorSectionRegistryLock.RLock()
	section := macInspectorSectionRegistry[event.sectionID]
	macInspectorSectionRegistryLock.RUnlock()
	if section == nil || section.isDead() {
		return
	}
	section.lock.Lock()
	if section.collapsible {
		section.collapsed = event.collapsed
	}
	section.lock.Unlock()
}

func init() {
	registerChromeEventLoop(func(*App) {
		for {
			event := <-macInspectorSectionEvents
			go handleMacInspectorSectionEvent(event)
		}
	})
}
