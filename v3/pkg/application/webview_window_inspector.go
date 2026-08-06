package application

import (
	"slices"
	"sync"
	"sync/atomic"
)

// MacInspector describes a native AppKit property inspector hosted by an
// inspector NSSplitViewItem. Sections and controls are rendered by AppKit;
// adding an inspector does not create another WebView.
//
// Keep the control handles returned by AddLabel, AddTextField, AddCheckbox,
// and AddPopup to update values or register callbacks. Control identifiers
// are generated internally and are never part of the application API.
type MacInspector struct {
	lock sync.RWMutex

	sections []*MacInspectorSection
	pane     *MacSplitPane
	dead     bool
}

// MacInspectorSection groups related controls under a native section heading.
type MacInspectorSection struct {
	lock sync.RWMutex

	internalID uint64
	label      string
	inspector  *MacInspector
	controls   []*MacInspectorControl
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
)

// MacInspectorControl is a live handle for one native inspector control.
// Kind-specific setters are safe no-ops when used with the wrong kind.
type MacInspectorControl struct {
	lock sync.RWMutex

	internalID   uint64
	kind         MacInspectorControlKind
	label        string
	value        string
	checked      bool
	options      []string
	selected     int
	tooltip      string
	disabled     bool
	hidden       bool
	inspector    *MacInspector
	onTextChange func(*Context, string)
	onToggle     func(*Context, bool)
	onSelection  func(*Context, int, string)
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
	macInspectorApplySnapshot(i)
	return section
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

// AddLabel adds a read-only property value.
func (s *MacInspectorSection) AddLabel(label, value string) *MacInspectorControl {
	return s.addControl(MacInspectorLabel, label, value, nil)
}

// AddTextField adds an editable native text field. Use OnTextChange to
// observe edits made by the user.
func (s *MacInspectorSection) AddTextField(label, value string) *MacInspectorControl {
	return s.addControl(MacInspectorTextField, label, value, nil)
}

// AddCheckbox adds a native checkbox.
func (s *MacInspectorSection) AddCheckbox(label string, checked bool) *MacInspectorControl {
	control := s.addControl(MacInspectorCheckbox, label, "", nil)
	if control != nil {
		control.lock.Lock()
		control.checked = checked
		control.lock.Unlock()
		macInspectorApplyControl(control)
	}
	return control
}

// AddPopup adds a native pop-up button. selectedIndex may be -1 for no
// selection; invalid positive indexes are normalized to the first option.
func (s *MacInspectorSection) AddPopup(label string, options []string, selectedIndex int) *MacInspectorControl {
	control := s.addControl(MacInspectorPopup, label, "", options)
	if control == nil {
		return nil
	}
	control.lock.Lock()
	control.selected = normalizeMacInspectorSelection(selectedIndex, len(control.options))
	control.lock.Unlock()
	macInspectorApplyControl(control)
	return control
}

func (s *MacInspectorSection) addControl(kind MacInspectorControlKind, label, value string, options []string) *MacInspectorControl {
	if s == nil || s.inspector == nil || s.isDead() {
		return nil
	}
	control := &MacInspectorControl{
		internalID: nextMacInspectorNodeID(),
		kind:       kind,
		label:      label,
		value:      value,
		options:    append([]string(nil), options...),
		selected:   -1,
		inspector:  s.inspector,
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

// Kind returns the control's immutable native kind.
func (c *MacInspectorControl) Kind() MacInspectorControlKind {
	if c == nil {
		return MacInspectorLabel
	}
	return c.kind
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
	if c == nil || c.isDead() || (c.kind != MacInspectorLabel && c.kind != MacInspectorTextField) {
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

// SetOptions replaces a pop-up button's options. The current selection is
// retained when possible and otherwise normalized to the first option.
func (c *MacInspectorControl) SetOptions(options []string) *MacInspectorControl {
	if c == nil || c.isDead() || c.kind != MacInspectorPopup {
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

// SetSelectedIndex updates a pop-up selection without invoking
// OnSelectionChange. Invalid indexes are ignored.
func (c *MacInspectorControl) SetSelectedIndex(index int) *MacInspectorControl {
	if c == nil || c.isDead() || c.kind != MacInspectorPopup {
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

// SelectedIndex returns the latest pop-up selection, or -1 when unselected.
func (c *MacInspectorControl) SelectedIndex() int {
	if c == nil {
		return -1
	}
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.selected
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

// OnSelectionChange sets the callback for user changes to a pop-up button.
// The callback receives both the selected index and value.
func (c *MacInspectorControl) OnSelectionChange(callback func(*Context, int, string)) *MacInspectorControl {
	if c == nil || c.isDead() || c.kind != MacInspectorPopup {
		return c
	}
	c.lock.Lock()
	c.onSelection = callback
	c.lock.Unlock()
	return c
}

func (s *MacInspectorSection) isDead() bool {
	if s == nil || s.inspector == nil {
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
	c.inspector.lock.RLock()
	defer c.inspector.lock.RUnlock()
	return c.inspector.dead
}

type macInspectorControlSnapshot struct {
	internalID uint64
	kind       MacInspectorControlKind
	label      string
	value      string
	checked    bool
	options    []string
	selected   int
	tooltip    string
	disabled   bool
	hidden     bool
}

type macInspectorSectionSnapshot struct {
	internalID uint64
	label      string
	controls   []macInspectorControlSnapshot
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
			internalID: section.internalID,
			label:      section.label,
			controls:   make([]macInspectorControlSnapshot, 0, len(section.controls)),
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
	for _, control := range i.controlHandles() {
		registerMacInspectorControl(control)
	}
}

func (i *MacInspector) markDead() {
	if i == nil {
		return
	}
	for _, control := range i.controlHandles() {
		unregisterMacInspectorControl(control.internalID)
		control.lock.Lock()
		control.onTextChange = nil
		control.onToggle = nil
		control.onSelection = nil
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

type macInspectorControlEvent struct {
	controlID uint64
	kind      MacInspectorControlKind
	value     string
	checked   bool
	selected  int
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
	case MacInspectorPopup:
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
	default:
		control.lock.Unlock()
	}
}
