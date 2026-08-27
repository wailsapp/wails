package application

import (
	"sync"
	"sync/atomic"
)

// MacSidebar describes a native AppKit source list hosted by a sidebar
// NSSplitViewItem. Its sections and rows are NSOutlineView content; no WebView
// is created for the sidebar.
type MacSidebar struct {
	lock sync.RWMutex

	entries  []macSidebarEntry
	selected *MacSidebarItem
	pane     *MacSplitPane
	dead     bool
}

type macSidebarEntry struct {
	section *MacSidebarSection
	item    *MacSidebarItem
}

// MacSidebarSection is a native source-list group. Sections are expanded by
// default and are not selectable.
type MacSidebarSection struct {
	lock sync.RWMutex

	internalID uint64
	label      string
	sidebar    *MacSidebar
	items      []*MacSidebarItem
}

// MacSidebarItem is a selectable native source-list row. Identifiers are
// generated internally; retain the returned handle to update or select a row.
type MacSidebarItem struct {
	lock sync.RWMutex

	internalID uint64
	label      string
	subtitle   string
	symbolName string
	imageData  []byte
	tooltip    string
	disabled   bool
	hidden     bool
	sidebar    *MacSidebar
	onClick    func(*Context)
}

var macSidebarNodeID uint64

func nextMacSidebarNodeID() uint64 {
	return atomic.AddUint64(&macSidebarNodeID, 1)
}

// NewMacSidebar creates an empty native source list.
func NewMacSidebar() *MacSidebar {
	return &MacSidebar{}
}

// AddSection adds a non-selectable source-list group.
func (s *MacSidebar) AddSection(label string) *MacSidebarSection {
	if s == nil {
		return nil
	}
	section := &MacSidebarSection{
		internalID: nextMacSidebarNodeID(),
		label:      label,
		sidebar:    s,
	}
	s.lock.Lock()
	if s.dead {
		s.lock.Unlock()
		return nil
	}
	s.entries = append(s.entries, macSidebarEntry{section: section})
	s.lock.Unlock()
	s.reload()
	return section
}

// AddItem adds a row at the source list's root, before or between sections.
func (s *MacSidebar) AddItem(label string) *MacSidebarItem {
	if s == nil {
		return nil
	}
	item := newMacSidebarItem(s, label)
	s.lock.Lock()
	if s.dead {
		s.lock.Unlock()
		return nil
	}
	s.entries = append(s.entries, macSidebarEntry{item: item})
	s.lock.Unlock()
	s.reload()
	return item
}

// AddItem adds a row to this section.
func (s *MacSidebarSection) AddItem(label string) *MacSidebarItem {
	if s == nil || s.sidebar == nil {
		return nil
	}
	item := newMacSidebarItem(s.sidebar, label)
	s.sidebar.lock.RLock()
	dead := s.sidebar.dead
	s.sidebar.lock.RUnlock()
	if dead {
		return nil
	}
	s.lock.Lock()
	s.items = append(s.items, item)
	s.lock.Unlock()
	s.sidebar.reload()
	return item
}

func newMacSidebarItem(sidebar *MacSidebar, label string) *MacSidebarItem {
	return &MacSidebarItem{
		internalID: nextMacSidebarNodeID(),
		label:      label,
		sidebar:    sidebar,
	}
}

// SetLabel updates the row's primary text.
func (i *MacSidebarItem) SetLabel(label string) *MacSidebarItem {
	if i == nil || i.isDead() {
		return i
	}
	i.lock.Lock()
	i.label = label
	i.lock.Unlock()
	i.sidebar.reload()
	return i
}

// SetSubtitle updates the secondary text in a native sidebar row. A row with
// a subtitle uses the larger two-line AppKit presentation.
func (i *MacSidebarItem) SetSubtitle(subtitle string) *MacSidebarItem {
	if i == nil || i.isDead() {
		return i
	}
	i.lock.Lock()
	i.subtitle = subtitle
	i.lock.Unlock()
	i.sidebar.reload()
	return i
}

// SetSymbol sets an SF Symbol for the row on macOS 11 and newer. Passing an
// empty string removes the symbol. SetImage takes precedence when an image is
// also present.
func (i *MacSidebarItem) SetSymbol(symbol string) *MacSidebarItem {
	if i == nil || i.isDead() {
		return i
	}
	i.lock.Lock()
	i.symbolName = symbol
	i.lock.Unlock()
	i.sidebar.reload()
	return i
}

// SetImage sets an encoded image (for example PNG or JPEG) for the row.
// Passing nil or an empty slice removes the image. An image or subtitle uses
// the larger AppKit presentation, with the image displayed as an avatar.
func (i *MacSidebarItem) SetImage(image []byte) *MacSidebarItem {
	if i == nil || i.isDead() {
		return i
	}
	i.lock.Lock()
	i.imageData = append(i.imageData[:0], image...)
	i.lock.Unlock()
	i.sidebar.reload()
	return i
}

// SetTooltip sets the native row tooltip.
func (i *MacSidebarItem) SetTooltip(tooltip string) *MacSidebarItem {
	if i == nil || i.isDead() {
		return i
	}
	i.lock.Lock()
	i.tooltip = tooltip
	i.lock.Unlock()
	i.sidebar.reload()
	return i
}

// SetEnabled controls whether the row can be selected.
func (i *MacSidebarItem) SetEnabled(enabled bool) *MacSidebarItem {
	if i == nil || i.isDead() {
		return i
	}
	i.lock.Lock()
	i.disabled = !enabled
	i.lock.Unlock()
	i.sidebar.reload()
	return i
}

// SetHidden includes or removes the row from the native source list.
func (i *MacSidebarItem) SetHidden(hidden bool) *MacSidebarItem {
	if i == nil || i.isDead() {
		return i
	}
	i.lock.Lock()
	i.hidden = hidden
	i.lock.Unlock()
	i.sidebar.reload()
	return i
}

// OnClick sets the callback invoked when AppKit selects this row. Passing nil
// clears the callback.
func (i *MacSidebarItem) OnClick(callback func(*Context)) *MacSidebarItem {
	if i == nil || i.isDead() {
		return i
	}
	i.lock.Lock()
	i.onClick = callback
	i.lock.Unlock()
	return i
}

// SetSelectedItem selects a row without invoking its OnClick callback.
// Passing nil clears the selection. Items belonging to another sidebar are
// ignored.
func (s *MacSidebar) SetSelectedItem(item *MacSidebarItem) *MacSidebar {
	if s == nil || (item != nil && item.sidebar != s) {
		return s
	}
	s.lock.Lock()
	if s.dead {
		s.lock.Unlock()
		return s
	}
	s.selected = item
	s.lock.Unlock()
	macSidebarApplySelection(s, item)
	return s
}

func (i *MacSidebarItem) isDead() bool {
	if i == nil || i.sidebar == nil {
		return true
	}
	i.sidebar.lock.RLock()
	defer i.sidebar.lock.RUnlock()
	return i.sidebar.dead
}

type macSidebarItemSnapshot struct {
	internalID uint64
	label      string
	subtitle   string
	symbolName string
	imageData  []byte
	tooltip    string
	disabled   bool
	hidden     bool
}

type macSidebarSectionSnapshot struct {
	internalID uint64
	label      string
	items      []macSidebarItemSnapshot
}

type macSidebarEntrySnapshot struct {
	section *macSidebarSectionSnapshot
	item    *macSidebarItemSnapshot
}

type macSidebarSnapshot struct {
	entries        []macSidebarEntrySnapshot
	selectedItemID uint64
}

func (s *MacSidebar) snapshot() macSidebarSnapshot {
	if s == nil {
		return macSidebarSnapshot{}
	}
	s.lock.RLock()
	entries := append([]macSidebarEntry(nil), s.entries...)
	selected := s.selected
	s.lock.RUnlock()
	result := macSidebarSnapshot{entries: make([]macSidebarEntrySnapshot, 0, len(entries))}
	if selected != nil {
		result.selectedItemID = selected.internalID
	}
	for _, entry := range entries {
		if entry.item != nil {
			item := snapshotMacSidebarItem(entry.item)
			result.entries = append(result.entries, macSidebarEntrySnapshot{item: &item})
			continue
		}
		if entry.section == nil {
			continue
		}
		entry.section.lock.RLock()
		section := macSidebarSectionSnapshot{
			internalID: entry.section.internalID,
			label:      entry.section.label,
			items:      make([]macSidebarItemSnapshot, 0, len(entry.section.items)),
		}
		for _, item := range entry.section.items {
			section.items = append(section.items, snapshotMacSidebarItem(item))
		}
		entry.section.lock.RUnlock()
		result.entries = append(result.entries, macSidebarEntrySnapshot{section: &section})
	}
	return result
}

func snapshotMacSidebarItem(item *MacSidebarItem) macSidebarItemSnapshot {
	item.lock.RLock()
	defer item.lock.RUnlock()
	return macSidebarItemSnapshot{
		internalID: item.internalID,
		label:      item.label,
		subtitle:   item.subtitle,
		symbolName: item.symbolName,
		imageData:  append([]byte(nil), item.imageData...),
		tooltip:    item.tooltip,
		disabled:   item.disabled,
		hidden:     item.hidden,
	}
}

func (s *MacSidebar) reload() {
	if s != nil {
		macSidebarApplySnapshot(s)
	}
}

func (s *MacSidebar) markDead() {
	if s == nil {
		return
	}
	for _, item := range s.itemHandles() {
		unregisterMacSidebarItem(item.internalID)
		item.lock.Lock()
		item.onClick = nil
		item.lock.Unlock()
	}
	s.lock.Lock()
	s.dead = true
	s.pane = nil
	s.selected = nil
	s.lock.Unlock()
}

func (s *MacSidebar) itemHandles() []*MacSidebarItem {
	if s == nil {
		return nil
	}
	s.lock.RLock()
	entries := append([]macSidebarEntry(nil), s.entries...)
	s.lock.RUnlock()
	var result []*MacSidebarItem
	for _, entry := range entries {
		if entry.item != nil {
			result = append(result, entry.item)
		}
		if entry.section != nil {
			entry.section.lock.RLock()
			result = append(result, entry.section.items...)
			entry.section.lock.RUnlock()
		}
	}
	return result
}

func (s *MacSidebar) registerItems() {
	for _, item := range s.itemHandles() {
		registerMacSidebarItem(item)
	}
}

var macSidebarItemRegistry = make(map[uint64]*MacSidebarItem)
var macSidebarItemRegistryLock sync.RWMutex

func registerMacSidebarItem(item *MacSidebarItem) {
	if item == nil {
		return
	}
	macSidebarItemRegistryLock.Lock()
	macSidebarItemRegistry[item.internalID] = item
	macSidebarItemRegistryLock.Unlock()
}

func unregisterMacSidebarItem(id uint64) {
	macSidebarItemRegistryLock.Lock()
	delete(macSidebarItemRegistry, id)
	macSidebarItemRegistryLock.Unlock()
}

func handleMacSidebarItemSelected(id uint64) {
	defer handlePanic()
	macSidebarItemRegistryLock.RLock()
	item := macSidebarItemRegistry[id]
	macSidebarItemRegistryLock.RUnlock()
	if item == nil || item.isDead() {
		return
	}
	item.sidebar.lock.Lock()
	item.sidebar.selected = item
	item.sidebar.lock.Unlock()
	item.lock.RLock()
	callback := item.onClick
	item.lock.RUnlock()
	if callback != nil {
		callback(newContext())
	}
}

var macSidebarItemSelected = make(chan uint64, 64)
