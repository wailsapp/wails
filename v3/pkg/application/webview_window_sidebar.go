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
	// selectedItems mirrors the native selection in row order. With single
	// selection it holds at most the selected item.
	selectedItems           []*MacSidebarItem
	allowsMultipleSelection bool
	reorderable             bool
	contextMenu             *Menu
	onContextMenu           func(*Context, *MacSidebarItem) *Menu
	onSelectionChange       func(*Context, []*MacSidebarItem)
	onMove                  func(*Context, *MacSidebarItem, *MacSidebarSection, int)
	pane                    *MacSplitPane
	dead                    bool
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
	removed    bool
}

// MacSidebarItem is a selectable native source-list row. Identifiers are
// generated internally; retain the returned handle to update or select a row.
// Rows nest to any depth through AddItem; nested rows start collapsed.
type MacSidebarItem struct {
	lock sync.RWMutex

	internalID      uint64
	label           string
	symbolName      string
	tooltip         string
	disabled        bool
	hidden          bool
	expanded        bool
	editable        bool
	badge           int
	accessorySymbol string
	tintColor       *RGBA
	removed         bool
	sidebar         *MacSidebar
	// A row lives in exactly one container: a section, a parent row, or the
	// sidebar root when both are nil.
	section          *MacSidebarSection
	parent           *MacSidebarItem
	children         []*MacSidebarItem
	contextMenu      *Menu
	onClick          func(*Context)
	onExpandedChange func(*Context, bool)
	onRename         func(*Context, string)
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
	if s == nil || s.sidebar == nil || s.isDead() {
		return nil
	}
	item := newMacSidebarItem(s.sidebar, label)
	item.section = s
	s.lock.Lock()
	s.items = append(s.items, item)
	s.lock.Unlock()
	s.sidebar.reload()
	return item
}

// Items returns the section's direct rows in display order.
func (s *MacSidebarSection) Items() []*MacSidebarItem {
	if s == nil {
		return nil
	}
	s.lock.RLock()
	defer s.lock.RUnlock()
	return append([]*MacSidebarItem(nil), s.items...)
}

// Remove detaches a row from this section, together with any nested rows,
// and unregisters their callbacks. The handle is inert afterwards.
func (s *MacSidebarSection) Remove(item *MacSidebarItem) *MacSidebarSection {
	if s == nil || item == nil {
		return s
	}
	item.lock.RLock()
	owner := item.section
	item.lock.RUnlock()
	if owner == s {
		item.Remove()
	}
	return s
}

// Move reorders a direct row of this section to index, clamped to the
// section's bounds.
func (s *MacSidebarSection) Move(item *MacSidebarItem, index int) *MacSidebarSection {
	if s == nil || item == nil || s.isDead() || item.isDead() {
		return s
	}
	item.lock.RLock()
	owner := item.section
	item.lock.RUnlock()
	if owner != s {
		return s
	}
	s.lock.Lock()
	s.items = moveMacSidebarItem(s.items, item, index)
	s.lock.Unlock()
	s.sidebar.reload()
	return s
}

func (s *MacSidebarSection) isDead() bool {
	if s == nil || s.sidebar == nil {
		return true
	}
	s.lock.RLock()
	removed := s.removed
	s.lock.RUnlock()
	if removed {
		return true
	}
	s.sidebar.lock.RLock()
	defer s.sidebar.lock.RUnlock()
	return s.sidebar.dead
}

// AddItem adds a nested row beneath this row. Nested rows are shown when the
// parent is expanded.
func (i *MacSidebarItem) AddItem(label string) *MacSidebarItem {
	if i == nil || i.isDead() {
		return nil
	}
	child := newMacSidebarItem(i.sidebar, label)
	child.parent = i
	i.lock.Lock()
	i.children = append(i.children, child)
	i.lock.Unlock()
	i.sidebar.reload()
	return child
}

// Items returns the row's direct nested rows in display order.
func (i *MacSidebarItem) Items() []*MacSidebarItem {
	if i == nil {
		return nil
	}
	i.lock.RLock()
	defer i.lock.RUnlock()
	return append([]*MacSidebarItem(nil), i.children...)
}

// Parent returns the row this row nests beneath, or nil for section and root
// rows.
func (i *MacSidebarItem) Parent() *MacSidebarItem {
	if i == nil {
		return nil
	}
	i.lock.RLock()
	defer i.lock.RUnlock()
	return i.parent
}

// Section returns the section containing this row, walking up through nested
// parents. It is nil for rows at the sidebar root.
func (i *MacSidebarItem) Section() *MacSidebarSection {
	for current := i; current != nil; {
		current.lock.RLock()
		section, parent := current.section, current.parent
		current.lock.RUnlock()
		if section != nil {
			return section
		}
		current = parent
	}
	return nil
}

func newMacSidebarItem(sidebar *MacSidebar, label string) *MacSidebarItem {
	return &MacSidebarItem{
		internalID: nextMacSidebarNodeID(),
		label:      label,
		sidebar:    sidebar,
	}
}

func removeMacSidebarItemFrom(items []*MacSidebarItem, item *MacSidebarItem) ([]*MacSidebarItem, int) {
	for index, candidate := range items {
		if candidate == item {
			return append(items[:index:index], items[index+1:]...), index
		}
	}
	return items, -1
}

func insertMacSidebarItem(items []*MacSidebarItem, item *MacSidebarItem, index int) ([]*MacSidebarItem, int) {
	if index < 0 {
		index = 0
	}
	if index > len(items) {
		index = len(items)
	}
	result := make([]*MacSidebarItem, 0, len(items)+1)
	result = append(result, items[:index]...)
	result = append(result, item)
	result = append(result, items[index:]...)
	return result, index
}

func moveMacSidebarItem(items []*MacSidebarItem, item *MacSidebarItem, index int) []*MacSidebarItem {
	remaining, old := removeMacSidebarItemFrom(items, item)
	if old < 0 {
		return items
	}
	if index > len(remaining) {
		index = len(remaining)
	}
	result, _ := insertMacSidebarItem(remaining, item, index)
	return result
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

// SetSymbol sets an SF Symbol for the row on macOS 11 and newer. Passing an
// empty string removes the image.
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

// SetExpanded expands or collapses this row's nested rows. It does not
// invoke OnExpandedChange.
func (i *MacSidebarItem) SetExpanded(expanded bool) *MacSidebarItem {
	if i == nil || i.isDead() {
		return i
	}
	i.lock.Lock()
	i.expanded = expanded
	i.lock.Unlock()
	i.sidebar.reload()
	return i
}

// IsExpanded reports whether the row's nested rows are shown.
func (i *MacSidebarItem) IsExpanded() bool {
	if i == nil {
		return false
	}
	i.lock.RLock()
	defer i.lock.RUnlock()
	return i.expanded
}

// OnExpandedChange sets the callback invoked when the user expands or
// collapses this row. Passing nil clears it.
func (i *MacSidebarItem) OnExpandedChange(callback func(*Context, bool)) *MacSidebarItem {
	if i == nil || i.isDead() {
		return i
	}
	i.lock.Lock()
	i.onExpandedChange = callback
	i.lock.Unlock()
	return i
}

// SetBadge shows a trailing count beside the row, like unread counts in
// Mail. Zero or a negative count removes the badge.
func (i *MacSidebarItem) SetBadge(count int) *MacSidebarItem {
	if i == nil || i.isDead() {
		return i
	}
	if count < 0 {
		count = 0
	}
	i.lock.Lock()
	i.badge = count
	i.lock.Unlock()
	i.sidebar.reload()
	return i
}

// Badge returns the row's current badge count, or zero when none is shown.
func (i *MacSidebarItem) Badge() int {
	if i == nil {
		return 0
	}
	i.lock.RLock()
	defer i.lock.RUnlock()
	return i.badge
}

// SetAccessorySymbol shows a trailing SF Symbol after the badge, for example
// a pin or a cloud status glyph. Passing an empty string removes it.
func (i *MacSidebarItem) SetAccessorySymbol(symbol string) *MacSidebarItem {
	if i == nil || i.isDead() {
		return i
	}
	i.lock.Lock()
	i.accessorySymbol = symbol
	i.lock.Unlock()
	i.sidebar.reload()
	return i
}

// SetTintColor tints the leading symbol on macOS 10.14 and newer. Passing
// nil restores the standard sidebar symbol colour.
func (i *MacSidebarItem) SetTintColor(colour *RGBA) *MacSidebarItem {
	if i == nil || i.isDead() {
		return i
	}
	var copied *RGBA
	if colour != nil {
		value := *colour
		copied = &value
	}
	i.lock.Lock()
	i.tintColor = copied
	i.lock.Unlock()
	i.sidebar.reload()
	return i
}

// SetEditable allows the user to rename the row in place with Return or a
// click on the selected row's label. Commits update the Go label and invoke
// OnRename.
func (i *MacSidebarItem) SetEditable(editable bool) *MacSidebarItem {
	if i == nil || i.isDead() {
		return i
	}
	i.lock.Lock()
	i.editable = editable
	i.lock.Unlock()
	i.sidebar.reload()
	return i
}

// OnRename sets the callback invoked after an inline rename commits. The
// handle's label is already updated when it runs. Passing nil clears it.
func (i *MacSidebarItem) OnRename(callback func(*Context, string)) *MacSidebarItem {
	if i == nil || i.isDead() {
		return i
	}
	i.lock.Lock()
	i.onRename = callback
	i.lock.Unlock()
	return i
}

// SetContextMenu sets the native menu shown when the row is right-clicked.
// Item OnClick handlers on the menu fire through the normal menu plumbing.
// Passing nil removes it; the sidebar's OnContextMenu callback and fallback
// menu are consulted afterwards.
func (i *MacSidebarItem) SetContextMenu(menu *Menu) *MacSidebarItem {
	if i == nil || i.isDead() {
		return i
	}
	i.lock.Lock()
	i.contextMenu = menu
	i.lock.Unlock()
	return i
}

// Remove detaches the row and its nested rows from the sidebar and
// unregisters their callbacks. The handle is inert afterwards.
func (i *MacSidebarItem) Remove() {
	if i == nil || i.isDead() {
		return
	}
	i.sidebar.removeItem(i)
}

// SetSelectedItem selects a row without invoking its OnClick callback.
// Passing nil clears the selection. Items belonging to another sidebar are
// ignored.
func (s *MacSidebar) SetSelectedItem(item *MacSidebarItem) *MacSidebar {
	if s == nil || (item != nil && (item.sidebar != s || item.isDead())) {
		return s
	}
	s.lock.Lock()
	if s.dead {
		s.lock.Unlock()
		return s
	}
	s.selected = item
	s.selectedItems = nil
	if item != nil {
		s.selectedItems = []*MacSidebarItem{item}
	}
	s.lock.Unlock()
	macSidebarApplySelection(s, item)
	return s
}

// SelectedItem returns the primary selected row, or nil when nothing is
// selected.
func (s *MacSidebar) SelectedItem() *MacSidebarItem {
	if s == nil {
		return nil
	}
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.selected
}

// SelectedItems returns every selected row in display order.
func (s *MacSidebar) SelectedItems() []*MacSidebarItem {
	if s == nil {
		return nil
	}
	s.lock.RLock()
	defer s.lock.RUnlock()
	return append([]*MacSidebarItem(nil), s.selectedItems...)
}

// SetAllowsMultipleSelection lets the user select several rows with Command
// and Shift clicks. OnClick still fires for the clicked row; use
// OnSelectionChange to observe the whole selection.
func (s *MacSidebar) SetAllowsMultipleSelection(allowed bool) *MacSidebar {
	if s == nil {
		return s
	}
	s.lock.Lock()
	if s.dead {
		s.lock.Unlock()
		return s
	}
	s.allowsMultipleSelection = allowed
	s.lock.Unlock()
	s.reload()
	return s
}

// OnSelectionChange sets the callback invoked whenever the user changes the
// selection, including clearing it. The slice lists selected rows in display
// order. Passing nil clears the callback.
func (s *MacSidebar) OnSelectionChange(callback func(*Context, []*MacSidebarItem)) *MacSidebar {
	if s == nil {
		return s
	}
	s.lock.Lock()
	s.onSelectionChange = callback
	s.lock.Unlock()
	return s
}

// SetContextMenu sets the native menu shown when the user right-clicks a row
// without its own menu or the empty area below the rows.
func (s *MacSidebar) SetContextMenu(menu *Menu) *MacSidebar {
	if s == nil {
		return s
	}
	s.lock.Lock()
	s.contextMenu = menu
	s.lock.Unlock()
	return s
}

// OnContextMenu sets a callback that builds the context menu for a
// right-click. item is nil for the empty area. Returning nil falls back to
// the sidebar menu set with SetContextMenu. A row's own SetContextMenu takes
// precedence over this callback. The callback runs synchronously on the
// application thread while AppKit waits to show the menu, so keep it quick.
func (s *MacSidebar) OnContextMenu(callback func(*Context, *MacSidebarItem) *Menu) *MacSidebar {
	if s == nil {
		return s
	}
	s.lock.Lock()
	s.onContextMenu = callback
	s.lock.Unlock()
	return s
}

// SetReorderable lets the user drag rows to reorder them. Rows can move
// within their section, between sections, and to or from the sidebar root;
// nested rows reorder beneath their parent. Nothing is accepted from outside
// the sidebar. Drops update the Go model and invoke OnMove.
func (s *MacSidebar) SetReorderable(reorderable bool) *MacSidebar {
	if s == nil {
		return s
	}
	s.lock.Lock()
	if s.dead {
		s.lock.Unlock()
		return s
	}
	s.reorderable = reorderable
	s.lock.Unlock()
	s.reload()
	return s
}

// OnMove sets the callback invoked after a drag reorder commits. section is
// the row's new section, or nil at the sidebar root; index is the row's new
// position inside its container. Passing nil clears the callback.
func (s *MacSidebar) OnMove(callback func(*Context, *MacSidebarItem, *MacSidebarSection, int)) *MacSidebar {
	if s == nil {
		return s
	}
	s.lock.Lock()
	s.onMove = callback
	s.lock.Unlock()
	return s
}

// Remove detaches any row belonging to this sidebar; see MacSidebarItem.Remove.
func (s *MacSidebar) Remove(item *MacSidebarItem) *MacSidebar {
	if s == nil || item == nil || item.sidebar != s {
		return s
	}
	item.Remove()
	return s
}

// RemoveSection detaches a section and all of its rows, unregistering their
// callbacks. The handles are inert afterwards.
func (s *MacSidebar) RemoveSection(section *MacSidebarSection) *MacSidebar {
	if s == nil || section == nil || section.sidebar != s || section.isDead() {
		return s
	}
	s.lock.Lock()
	if s.dead {
		s.lock.Unlock()
		return s
	}
	for index, entry := range s.entries {
		if entry.section == section {
			s.entries = append(s.entries[:index:index], s.entries[index+1:]...)
			break
		}
	}
	s.lock.Unlock()
	section.lock.Lock()
	section.removed = true
	items := section.items
	section.items = nil
	section.lock.Unlock()
	for _, item := range items {
		item.retire()
	}
	s.pruneSelection()
	s.reload()
	return s
}

// Sections returns the sidebar's sections in display order.
func (s *MacSidebar) Sections() []*MacSidebarSection {
	if s == nil {
		return nil
	}
	s.lock.RLock()
	defer s.lock.RUnlock()
	var result []*MacSidebarSection
	for _, entry := range s.entries {
		if entry.section != nil {
			result = append(result, entry.section)
		}
	}
	return result
}

// Items returns the sidebar's root rows in display order.
func (s *MacSidebar) Items() []*MacSidebarItem {
	if s == nil {
		return nil
	}
	s.lock.RLock()
	defer s.lock.RUnlock()
	var result []*MacSidebarItem
	for _, entry := range s.entries {
		if entry.item != nil {
			result = append(result, entry.item)
		}
	}
	return result
}

// removeItem detaches item from its container, retires the subtree, and
// reloads the native list.
func (s *MacSidebar) removeItem(item *MacSidebarItem) {
	s.detachItem(item)
	item.retire()
	s.pruneSelection()
	s.reload()
}

// detachItem removes item from whichever container holds it and returns the
// container and the index it occupied.
func (s *MacSidebar) detachItem(item *MacSidebarItem) (parent *MacSidebarItem, section *MacSidebarSection, index int) {
	item.lock.RLock()
	parent, section = item.parent, item.section
	item.lock.RUnlock()
	index = -1
	switch {
	case parent != nil:
		parent.lock.Lock()
		parent.children, index = removeMacSidebarItemFrom(parent.children, item)
		parent.lock.Unlock()
	case section != nil:
		section.lock.Lock()
		section.items, index = removeMacSidebarItemFrom(section.items, item)
		section.lock.Unlock()
	default:
		s.lock.Lock()
		for position, entry := range s.entries {
			if entry.item == item {
				s.entries = append(s.entries[:position:position], s.entries[position+1:]...)
				index = position
				break
			}
		}
		s.lock.Unlock()
	}
	return parent, section, index
}

// retire marks the subtree removed, drops its callbacks, and unregisters it.
func (i *MacSidebarItem) retire() {
	if i == nil {
		return
	}
	i.lock.Lock()
	i.removed = true
	i.onClick = nil
	i.onExpandedChange = nil
	i.onRename = nil
	i.contextMenu = nil
	children := i.children
	i.lock.Unlock()
	unregisterMacSidebarItem(i.internalID)
	for _, child := range children {
		child.retire()
	}
}

// pruneSelection drops removed rows from the selection.
func (s *MacSidebar) pruneSelection() {
	s.lock.Lock()
	defer s.lock.Unlock()
	kept := s.selectedItems[:0:0]
	for _, item := range s.selectedItems {
		item.lock.RLock()
		removed := item.removed
		item.lock.RUnlock()
		if !removed {
			kept = append(kept, item)
		}
	}
	s.selectedItems = kept
	if s.selected != nil {
		s.selected.lock.RLock()
		removed := s.selected.removed
		s.selected.lock.RUnlock()
		if removed {
			s.selected = nil
			if len(kept) > 0 {
				s.selected = kept[0]
			}
		}
	}
}

func (i *MacSidebarItem) isDead() bool {
	if i == nil || i.sidebar == nil {
		return true
	}
	i.lock.RLock()
	removed := i.removed
	i.lock.RUnlock()
	if removed {
		return true
	}
	i.sidebar.lock.RLock()
	defer i.sidebar.lock.RUnlock()
	return i.sidebar.dead
}

// isDescendantOf reports whether i is ancestor or nests beneath it.
func (i *MacSidebarItem) isDescendantOf(ancestor *MacSidebarItem) bool {
	for current := i; current != nil; {
		if current == ancestor {
			return true
		}
		current.lock.RLock()
		parent := current.parent
		current.lock.RUnlock()
		current = parent
	}
	return false
}

type macSidebarItemSnapshot struct {
	internalID      uint64
	label           string
	symbolName      string
	tooltip         string
	disabled        bool
	hidden          bool
	expanded        bool
	editable        bool
	badge           int
	accessorySymbol string
	tintColor       *RGBA
	children        []macSidebarItemSnapshot
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
	entries                 []macSidebarEntrySnapshot
	selectedItemID          uint64
	allowsMultipleSelection bool
	reorderable             bool
}

func (s *MacSidebar) snapshot() macSidebarSnapshot {
	if s == nil {
		return macSidebarSnapshot{}
	}
	s.lock.RLock()
	entries := append([]macSidebarEntry(nil), s.entries...)
	selected := s.selected
	result := macSidebarSnapshot{
		entries:                 make([]macSidebarEntrySnapshot, 0, len(entries)),
		allowsMultipleSelection: s.allowsMultipleSelection,
		reorderable:             s.reorderable,
	}
	s.lock.RUnlock()
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
		items := append([]*MacSidebarItem(nil), entry.section.items...)
		entry.section.lock.RUnlock()
		for _, item := range items {
			section.items = append(section.items, snapshotMacSidebarItem(item))
		}
		result.entries = append(result.entries, macSidebarEntrySnapshot{section: &section})
	}
	return result
}

func snapshotMacSidebarItem(item *MacSidebarItem) macSidebarItemSnapshot {
	item.lock.RLock()
	result := macSidebarItemSnapshot{
		internalID:      item.internalID,
		label:           item.label,
		symbolName:      item.symbolName,
		tooltip:         item.tooltip,
		disabled:        item.disabled,
		hidden:          item.hidden,
		expanded:        item.expanded,
		editable:        item.editable,
		badge:           item.badge,
		accessorySymbol: item.accessorySymbol,
	}
	if item.tintColor != nil {
		tint := *item.tintColor
		result.tintColor = &tint
	}
	children := append([]*MacSidebarItem(nil), item.children...)
	item.lock.RUnlock()
	if len(children) > 0 {
		result.children = make([]macSidebarItemSnapshot, 0, len(children))
		for _, child := range children {
			result.children = append(result.children, snapshotMacSidebarItem(child))
		}
	}
	return result
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
		item.onExpandedChange = nil
		item.onRename = nil
		item.contextMenu = nil
		item.lock.Unlock()
	}
	s.lock.Lock()
	s.dead = true
	s.pane = nil
	s.selected = nil
	s.selectedItems = nil
	s.contextMenu = nil
	s.onContextMenu = nil
	s.onSelectionChange = nil
	s.onMove = nil
	s.lock.Unlock()
}

// itemHandles returns every row in the sidebar, including nested rows, in
// display order.
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
			result = appendMacSidebarSubtree(result, entry.item)
		}
		if entry.section != nil {
			entry.section.lock.RLock()
			items := append([]*MacSidebarItem(nil), entry.section.items...)
			entry.section.lock.RUnlock()
			for _, item := range items {
				result = appendMacSidebarSubtree(result, item)
			}
		}
	}
	return result
}

func appendMacSidebarSubtree(result []*MacSidebarItem, item *MacSidebarItem) []*MacSidebarItem {
	result = append(result, item)
	item.lock.RLock()
	children := append([]*MacSidebarItem(nil), item.children...)
	item.lock.RUnlock()
	for _, child := range children {
		result = appendMacSidebarSubtree(result, child)
	}
	return result
}

func (s *MacSidebar) registerItems() {
	for _, item := range s.itemHandles() {
		registerMacSidebarItem(item)
	}
}

// findSection returns the live section with the given internal ID.
func (s *MacSidebar) findSection(id uint64) *MacSidebarSection {
	if s == nil || id == 0 {
		return nil
	}
	s.lock.RLock()
	defer s.lock.RUnlock()
	for _, entry := range s.entries {
		if entry.section != nil && entry.section.internalID == id {
			return entry.section
		}
	}
	return nil
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

func macSidebarItemByID(id uint64) *MacSidebarItem {
	macSidebarItemRegistryLock.RLock()
	defer macSidebarItemRegistryLock.RUnlock()
	return macSidebarItemRegistry[id]
}

// macSidebarForPane resolves the sidebar hosted by a native split pane.
func macSidebarForPane(paneID uint64) *MacSidebar {
	pane := macSplitPaneByID(paneID)
	if pane == nil || pane.sidebar == nil || pane.isDead() {
		return nil
	}
	return pane.sidebar
}

func handleMacSidebarItemSelected(id uint64) {
	defer handlePanic()
	item := macSidebarItemByID(id)
	if item == nil || item.isDead() {
		return
	}
	item.sidebar.lock.Lock()
	item.sidebar.selected = item
	item.sidebar.selectedItems = []*MacSidebarItem{item}
	item.sidebar.lock.Unlock()
	item.lock.RLock()
	callback := item.onClick
	item.lock.RUnlock()
	if callback != nil {
		callback(newContext())
	}
}

var macSidebarItemSelected = make(chan uint64, 64)

// macSidebarSelectionEvent carries the native selection. itemIDs lists the
// clicked row first, followed by the rest of the selection in row order.
type macSidebarSelectionEvent struct {
	paneID  uint64
	itemIDs []uint64
}

var macSidebarSelectionEvents = make(chan macSidebarSelectionEvent, 64)

// handleMacSidebarSelectionChanged updates the sidebar's selection, fires
// OnClick for the clicked row, and then OnSelectionChange for the whole set.
func handleMacSidebarSelectionChanged(event macSidebarSelectionEvent) {
	defer handlePanic()
	sidebar := macSidebarForPane(event.paneID)
	if sidebar == nil {
		return
	}
	var clicked *MacSidebarItem
	items := make([]*MacSidebarItem, 0, len(event.itemIDs))
	for index, id := range event.itemIDs {
		item := macSidebarItemByID(id)
		if item == nil || item.sidebar != sidebar || item.isDead() {
			continue
		}
		if index == 0 {
			clicked = item
		}
		items = append(items, item)
	}
	sidebar.lock.Lock()
	if sidebar.dead {
		sidebar.lock.Unlock()
		return
	}
	sidebar.selected = clicked
	sidebar.selectedItems = append([]*MacSidebarItem(nil), items...)
	onSelectionChange := sidebar.onSelectionChange
	sidebar.lock.Unlock()
	if clicked != nil {
		clicked.lock.RLock()
		callback := clicked.onClick
		clicked.lock.RUnlock()
		if callback != nil {
			callback(newContext())
		}
	}
	if onSelectionChange != nil {
		onSelectionChange(newContext(), append([]*MacSidebarItem(nil), items...))
	}
}

type macSidebarExpandedEvent struct {
	itemID   uint64
	expanded bool
}

var macSidebarExpandedEvents = make(chan macSidebarExpandedEvent, 64)

func handleMacSidebarItemExpanded(event macSidebarExpandedEvent) {
	defer handlePanic()
	item := macSidebarItemByID(event.itemID)
	if item == nil || item.isDead() {
		return
	}
	item.lock.Lock()
	changed := item.expanded != event.expanded
	item.expanded = event.expanded
	callback := item.onExpandedChange
	item.lock.Unlock()
	if changed && callback != nil {
		callback(newContext(), event.expanded)
	}
}

type macSidebarRenameEvent struct {
	itemID uint64
	label  string
}

var macSidebarRenameEvents = make(chan macSidebarRenameEvent, 64)

func handleMacSidebarItemRenamed(event macSidebarRenameEvent) {
	defer handlePanic()
	item := macSidebarItemByID(event.itemID)
	if item == nil || item.isDead() {
		return
	}
	item.lock.Lock()
	if !item.editable || item.label == event.label {
		item.lock.Unlock()
		return
	}
	item.label = event.label
	callback := item.onRename
	item.lock.Unlock()
	if callback != nil {
		callback(newContext(), event.label)
	}
}

// macSidebarMoveEvent describes a native drop. parentID is 0 for the sidebar
// root, otherwise a section or row identifier. index counts every child of
// the destination, hidden rows included, before the dragged row is removed.
type macSidebarMoveEvent struct {
	paneID   uint64
	itemID   uint64
	parentID uint64
	index    int
}

var macSidebarMoveEvents = make(chan macSidebarMoveEvent, 64)

func handleMacSidebarItemMoved(event macSidebarMoveEvent) {
	defer handlePanic()
	sidebar := macSidebarForPane(event.paneID)
	item := macSidebarItemByID(event.itemID)
	if sidebar == nil || item == nil || item.sidebar != sidebar || item.isDead() {
		return
	}
	var newParent *MacSidebarItem
	var newSection *MacSidebarSection
	if event.parentID != 0 {
		newParent = macSidebarItemByID(event.parentID)
		if newParent != nil && (newParent.sidebar != sidebar || newParent.isDead()) {
			newParent = nil
		}
		if newParent == nil {
			newSection = sidebar.findSection(event.parentID)
			if newSection == nil {
				return
			}
		}
	}
	section, index, ok := sidebar.moveItem(item, newParent, newSection, event.index)
	if !ok {
		return
	}
	sidebar.reload()
	sidebar.lock.RLock()
	callback := sidebar.onMove
	sidebar.lock.RUnlock()
	if callback != nil {
		callback(newContext(), item, section, index)
	}
}

// moveItem relocates item into the requested container at index and returns
// the section the row now belongs to and its final position. Reparenting a
// row beneath itself or one of its descendants is rejected.
func (s *MacSidebar) moveItem(item, newParent *MacSidebarItem, newSection *MacSidebarSection, index int) (*MacSidebarSection, int, bool) {
	if newParent != nil && newParent.isDescendantOf(item) {
		return nil, 0, false
	}
	if newSection != nil && newSection.isDead() {
		return nil, 0, false
	}
	oldParent, oldSection, oldIndex := s.detachItem(item)
	if oldIndex < 0 {
		return nil, 0, false
	}
	sameContainer := oldParent == newParent && oldSection == newSection
	if sameContainer && oldIndex < index {
		index--
	}
	item.lock.Lock()
	item.parent = newParent
	item.section = newSection
	item.lock.Unlock()
	switch {
	case newParent != nil:
		newParent.lock.Lock()
		newParent.children, index = insertMacSidebarItem(newParent.children, item, index)
		newParent.lock.Unlock()
		return newParent.Section(), index, true
	case newSection != nil:
		newSection.lock.Lock()
		newSection.items, index = insertMacSidebarItem(newSection.items, item, index)
		newSection.lock.Unlock()
		return newSection, index, true
	default:
		s.lock.Lock()
		if index < 0 {
			index = 0
		}
		if index > len(s.entries) {
			index = len(s.entries)
		}
		entries := make([]macSidebarEntry, 0, len(s.entries)+1)
		entries = append(entries, s.entries[:index]...)
		entries = append(entries, macSidebarEntry{item: item})
		entries = append(entries, s.entries[index:]...)
		s.entries = entries
		s.lock.Unlock()
		return nil, index, true
	}
}

// resolveMacSidebarContextMenu picks the menu for a right-click: the row's
// own menu, then the sidebar's OnContextMenu callback, then the sidebar's
// fallback menu. It runs synchronously on the application thread.
func resolveMacSidebarContextMenu(paneID, itemID uint64) (menu *Menu) {
	defer handlePanic()
	sidebar := macSidebarForPane(paneID)
	if sidebar == nil {
		return nil
	}
	var item *MacSidebarItem
	if itemID != 0 {
		item = macSidebarItemByID(itemID)
		if item != nil && (item.sidebar != sidebar || item.isDead()) {
			item = nil
		}
	}
	if item != nil {
		item.lock.RLock()
		menu = item.contextMenu
		item.lock.RUnlock()
		if menu != nil {
			return menu
		}
	}
	sidebar.lock.RLock()
	dynamic := sidebar.onContextMenu
	fallback := sidebar.contextMenu
	sidebar.lock.RUnlock()
	if dynamic != nil {
		if menu = dynamic(newContext(), item); menu != nil {
			return menu
		}
	}
	return fallback
}

func init() {
	registerChromeEventLoop(func(*App) {
		for {
			event := <-macSidebarSelectionEvents
			go handleMacSidebarSelectionChanged(event)
		}
	})
	registerChromeEventLoop(func(*App) {
		for {
			event := <-macSidebarExpandedEvents
			go handleMacSidebarItemExpanded(event)
		}
	})
	registerChromeEventLoop(func(*App) {
		for {
			event := <-macSidebarRenameEvents
			go handleMacSidebarItemRenamed(event)
		}
	})
	registerChromeEventLoop(func(*App) {
		for {
			event := <-macSidebarMoveEvents
			go handleMacSidebarItemMoved(event)
		}
	})
}
