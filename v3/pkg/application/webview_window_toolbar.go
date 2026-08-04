package application

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)

// MacToolbar describes a native macOS window toolbar (NSToolbar). Build one
// with NewMacToolbar and attach it with WebviewWindow.SetToolbar.
//
// Toolbar item identifiers are generated internally. Keep the returned item
// pointers when an item needs to be updated after installation.
type MacToolbar struct {
	items []*MacToolbarItem

	stateLock sync.RWMutex
	state     *macToolbarState
}

// NewMacToolbar creates an empty native macOS toolbar.
func NewMacToolbar() *MacToolbar { return &MacToolbar{} }

// MacToolbarItemKind selects what a MacToolbarItem renders as.
type MacToolbarItemKind int

const (
	ToolbarButton MacToolbarItemKind = iota
	ToolbarGroup
	ToolbarSearchField
	ToolbarFlexibleSpace
	ToolbarSidebarToggle
	ToolbarInspectorToggle
)

type MacToolbarGroupSelectionMode int

const (
	ToolbarGroupSelectOne MacToolbarGroupSelectionMode = iota
	ToolbarGroupMomentary
	ToolbarGroupSelectAny
)

// MacToolbarItem is both the item description and its live update handle.
// Construct items through MacToolbar's Add methods instead of assigning an
// internal identifier.
type MacToolbarItem struct {
	identifier string
	kind       MacToolbarItemKind
	toolbar    *MacToolbar
	parent     *MacToolbarItem

	label      string
	symbolName string
	tooltip    string
	bordered   bool
	prominent  bool
	tintColor  *RGBA
	badgeCount int
	disabled   bool
	hidden     bool

	items         []*MacToolbarItem
	selectionMode MacToolbarGroupSelectionMode
	selectedIndex int
	onClick       func(*Context)
	onSearch      func(*Context, string)
}

var toolbarItemIdentifier uint64

func newMacToolbarItem(toolbar *MacToolbar, kind MacToolbarItemKind, label string) *MacToolbarItem {
	sequence := atomic.AddUint64(&toolbarItemIdentifier, 1)
	return &MacToolbarItem{
		identifier: fmt.Sprintf("wails.toolbar.item.%d", sequence),
		kind:       kind,
		toolbar:    toolbar,
		label:      label,
	}
}

func (t *MacToolbar) add(item *MacToolbarItem) *MacToolbarItem {
	t.items = append(t.items, item)
	return item
}

// AddButton adds a native toolbar button. The optional symbol is an SF Symbol
// name. The returned item is used for callbacks and live state updates.
func (t *MacToolbar) AddButton(label string, symbol ...string) *MacToolbarItem {
	item := newMacToolbarItem(t, ToolbarButton, label)
	if len(symbol) > 0 {
		item.symbolName = symbol[0]
	}
	return t.add(item)
}

// AddSearch adds a native NSSearchToolbarItem.
func (t *MacToolbar) AddSearch(label string) *MacToolbarItem {
	return t.add(newMacToolbarItem(t, ToolbarSearchField, label))
}

// AddGroup adds a native segmented toolbar group. Add members with the
// returned item's AddButton method.
func (t *MacToolbar) AddGroup(label string, mode MacToolbarGroupSelectionMode) *MacToolbarItem {
	item := newMacToolbarItem(t, ToolbarGroup, label)
	item.selectionMode = mode
	return t.add(item)
}

func (t *MacToolbar) AddFlexibleSpace() *MacToolbarItem {
	return t.add(newMacToolbarItem(t, ToolbarFlexibleSpace, ""))
}

func (t *MacToolbar) AddSidebarToggle() *MacToolbarItem {
	return t.add(newMacToolbarItem(t, ToolbarSidebarToggle, "Sidebar"))
}

func (t *MacToolbar) AddInspectorToggle() *MacToolbarItem {
	return t.add(newMacToolbarItem(t, ToolbarInspectorToggle, "Inspector"))
}

// AddButton adds a member to a toolbar group.
func (g *MacToolbarItem) AddButton(label string, symbol ...string) *MacToolbarItem {
	if g == nil || g.kind != ToolbarGroup {
		return nil
	}
	item := newMacToolbarItem(g.toolbar, ToolbarButton, label)
	item.parent = g
	if len(symbol) > 0 {
		item.symbolName = symbol[0]
	}
	g.items = append(g.items, item)
	return item
}

func (i *MacToolbarItem) OnClick(callback func(*Context)) *MacToolbarItem {
	if i != nil {
		i.onClick = callback
	}
	return i
}

func (i *MacToolbarItem) OnSearch(callback func(*Context, string)) *MacToolbarItem {
	if i != nil {
		i.onSearch = callback
	}
	return i
}

func (i *MacToolbarItem) SetLabel(label string) *MacToolbarItem {
	if i == nil {
		return i
	}
	i.label = label
	i.update(func(native unsafe.Pointer) { macToolbarItemSetLabel(native, i.identifier, label) })
	return i
}

func (i *MacToolbarItem) SetSymbol(symbol string) *MacToolbarItem {
	if i != nil {
		i.symbolName = symbol
	}
	return i
}

func (i *MacToolbarItem) SetTooltip(tooltip string) *MacToolbarItem {
	if i != nil {
		i.tooltip = tooltip
	}
	return i
}

func (i *MacToolbarItem) SetBordered(bordered bool) *MacToolbarItem {
	if i != nil {
		i.bordered = bordered
	}
	return i
}

func (i *MacToolbarItem) SetProminent(prominent bool) *MacToolbarItem {
	if i != nil {
		i.prominent = prominent
	}
	return i
}

func (i *MacToolbarItem) SetTintColor(color *RGBA) *MacToolbarItem {
	if i != nil {
		i.tintColor = color
	}
	return i
}

func (i *MacToolbarItem) SetEnabled(enabled bool) *MacToolbarItem {
	if i == nil {
		return i
	}
	i.disabled = !enabled
	i.update(func(native unsafe.Pointer) { macToolbarItemSetEnabled(native, i.identifier, enabled) })
	return i
}

func (i *MacToolbarItem) SetHidden(hidden bool) *MacToolbarItem {
	if i == nil {
		return i
	}
	i.hidden = hidden
	i.update(func(native unsafe.Pointer) { macToolbarItemSetHidden(native, i.identifier, hidden) })
	return i
}

func (i *MacToolbarItem) SetBadgeCount(count int) *MacToolbarItem {
	if i == nil {
		return i
	}
	if count < 0 {
		count = 0
	}
	i.badgeCount = count
	i.update(func(native unsafe.Pointer) { macToolbarItemSetBadgeCount(native, i.identifier, count) })
	return i
}

func (i *MacToolbarItem) SetSelectedIndex(index int) *MacToolbarItem {
	if i == nil || i.kind != ToolbarGroup {
		return i
	}
	i.selectedIndex = index
	i.update(func(native unsafe.Pointer) { macToolbarGroupSetSelectedIndex(native, i.identifier, index) })
	return i
}

func (i *MacToolbarItem) SetSelectionMode(mode MacToolbarGroupSelectionMode) *MacToolbarItem {
	if i != nil && i.kind == ToolbarGroup {
		i.selectionMode = mode
	}
	return i
}

func (i *MacToolbarItem) update(update func(unsafe.Pointer)) {
	if i == nil || i.toolbar == nil {
		return
	}
	i.toolbar.stateLock.RLock()
	state := i.toolbar.state
	i.toolbar.stateLock.RUnlock()
	if state == nil {
		return
	}
	state.lock.RLock()
	native := state.native
	state.lock.RUnlock()
	if native != nil {
		InvokeSync(func() { update(native) })
	}
}

type macToolbarState struct {
	window  *WebviewWindow
	lock    sync.RWMutex
	native  unsafe.Pointer
	itemIDs []uint
}

func validateToolbarItems(items []*MacToolbarItem) error {
	var validate func([]*MacToolbarItem, string) error
	validate = func(items []*MacToolbarItem, parent string) error {
		for _, item := range items {
			if item == nil {
				return fmt.Errorf("toolbar item in %q is nil", parent)
			}
			switch item.kind {
			case ToolbarButton:
				if item.onClick == nil {
					return fmt.Errorf("toolbar button %q requires OnClick", item.label)
				}
			case ToolbarSearchField:
				if item.onSearch == nil {
					return fmt.Errorf("toolbar search field %q requires OnSearch", item.label)
				}
			case ToolbarGroup:
				if len(item.items) == 0 {
					return fmt.Errorf("toolbar group %q requires at least one member", item.label)
				}
				for _, member := range item.items {
					if member == nil || member.kind != ToolbarButton || member.onClick == nil {
						return fmt.Errorf("toolbar group %q members must be buttons with OnClick", item.label)
					}
				}
				if err := validate(item.items, item.label); err != nil {
					return err
				}
			case ToolbarFlexibleSpace, ToolbarSidebarToggle, ToolbarInspectorToggle:
			default:
				return fmt.Errorf("toolbar item %q has unknown kind %d", item.label, item.kind)
			}
		}
		return nil
	}
	return validate(items, "root")
}

var toolbarItemMap = make(map[uint]*MacToolbarItem)
var toolbarItemMapLock sync.Mutex
var toolbarNativeID uintptr

func nextToolbarNativeID() uint { return uint(atomic.AddUintptr(&toolbarNativeID, 1)) }
func addToToolbarItemMap(id uint, item *MacToolbarItem) {
	toolbarItemMapLock.Lock()
	toolbarItemMap[id] = item
	toolbarItemMapLock.Unlock()
}
func removeFromToolbarItemMap(id uint) {
	toolbarItemMapLock.Lock()
	delete(toolbarItemMap, id)
	toolbarItemMapLock.Unlock()
}
func getToolbarItemByID(id uint) *MacToolbarItem {
	toolbarItemMapLock.Lock()
	defer toolbarItemMapLock.Unlock()
	return toolbarItemMap[id]
}

var toolbarItemClicked = make(chan uint, 5)

type toolbarSearchEvent struct {
	itemID uint
	query  string
}

var toolbarSearchTriggered = make(chan toolbarSearchEvent, 5)

func handleToolbarItemClicked(itemID uint) {
	defer handlePanic()
	if item := getToolbarItemByID(itemID); item != nil && item.onClick != nil {
		item.onClick(newContext())
	}
}

func handleToolbarSearch(itemID uint, query string) {
	defer handlePanic()
	if item := getToolbarItemByID(itemID); item != nil && item.onSearch != nil {
		item.onSearch(newContext(), query)
	}
}
