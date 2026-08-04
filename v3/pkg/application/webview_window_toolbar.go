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
// Toolbar and item identifiers are generated internally. Keep the returned
// item pointers when an item needs to be updated after installation.
//
// A toolbar may only be attached to one window at a time. It can be detached
// with SetToolbar(nil) and subsequently attached to another window.
type MacToolbar struct {
	identifier string

	itemsLock sync.RWMutex
	items     []*MacToolbarItem

	// stateLock protects state and every field in macToolbarState. Native
	// operations take this lock on the application thread so a toolbar cannot
	// be detached while an item update is using its AppKit handle.
	stateLock sync.RWMutex
	state     *macToolbarState
}

var toolbarIdentifier uint64

// NewMacToolbar creates an empty native macOS toolbar.
func NewMacToolbar() *MacToolbar {
	sequence := atomic.AddUint64(&toolbarIdentifier, 1)
	return &MacToolbar{identifier: fmt.Sprintf("wails.toolbar.%d", sequence)}
}

type macToolbarItemKind int

const (
	toolbarButton macToolbarItemKind = iota
	toolbarGroup
	toolbarSearchField
	toolbarFlexibleSpace
)

// MacToolbarGroupSelectionMode controls how a native toolbar group behaves.
type MacToolbarGroupSelectionMode int

const (
	// ToolbarGroupSelectOne keeps one group member selected.
	ToolbarGroupSelectOne MacToolbarGroupSelectionMode = iota
	// ToolbarGroupMomentary makes each group member behave like a momentary
	// button without retaining selection.
	ToolbarGroupMomentary
	// ToolbarGroupSelectAny lets AppKit independently toggle group members.
	ToolbarGroupSelectAny
)

// MacToolbarItem is both an item description and its live update handle.
// Construct items through MacToolbar's Add methods instead of assigning an
// identifier. All Set methods update an installed native item immediately;
// where AppKit does not provide a feature on the current macOS version, the
// requested value remains stored and will be applied on a supported version.
type MacToolbarItem struct {
	lock sync.RWMutex

	identifier string
	kind       macToolbarItemKind
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

// MacToolbarGroup is a type-safe handle for a native segmented toolbar group.
// Its embedded MacToolbarItem supports the common item setters.
type MacToolbarGroup struct {
	*MacToolbarItem
}

var toolbarItemIdentifier uint64

func newMacToolbarItem(toolbar *MacToolbar, kind macToolbarItemKind, label string) *MacToolbarItem {
	sequence := atomic.AddUint64(&toolbarItemIdentifier, 1)
	return &MacToolbarItem{
		identifier: fmt.Sprintf("wails.toolbar.item.%d", sequence),
		kind:       kind,
		toolbar:    toolbar,
		label:      label,
	}
}

func (t *MacToolbar) add(item *MacToolbarItem) *MacToolbarItem {
	t.itemsLock.Lock()
	t.items = append(t.items, item)
	t.itemsLock.Unlock()
	return item
}

func (t *MacToolbar) itemSnapshot() []*MacToolbarItem {
	t.itemsLock.RLock()
	defer t.itemsLock.RUnlock()
	return append([]*MacToolbarItem(nil), t.items...)
}

// AddButton adds a native toolbar button. Use SetSymbol with an SF Symbol name
// when the button should have an icon. The returned item receives callbacks
// and provides live state updates.
func (t *MacToolbar) AddButton(label string) *MacToolbarItem {
	return t.add(newMacToolbarItem(t, toolbarButton, label))
}

// AddSearch adds a native NSSearchToolbarItem on macOS 11 and newer. Wails
// uses an NSToolbarItem containing an NSSearchField on older supported macOS
// releases, preserving the same callback behavior.
func (t *MacToolbar) AddSearch(label string) *MacToolbarItem {
	return t.add(newMacToolbarItem(t, toolbarSearchField, label))
}

// AddGroup adds a native segmented toolbar group. Add its members through the
// returned MacToolbarGroup.
func (t *MacToolbar) AddGroup(label string, mode MacToolbarGroupSelectionMode) *MacToolbarGroup {
	item := newMacToolbarItem(t, toolbarGroup, label)
	item.selectionMode = mode
	t.add(item)
	return &MacToolbarGroup{MacToolbarItem: item}
}

// AddFlexibleSpace adds a native flexible-space item.
func (t *MacToolbar) AddFlexibleSpace() *MacToolbarItem {
	return t.add(newMacToolbarItem(t, toolbarFlexibleSpace, ""))
}

// AddButton adds a member to a native toolbar group.
func (g *MacToolbarGroup) AddButton(label string) *MacToolbarItem {
	if g == nil || g.MacToolbarItem == nil {
		return nil
	}
	item := newMacToolbarItem(g.toolbar, toolbarButton, label)
	item.parent = g.MacToolbarItem
	g.lock.Lock()
	g.items = append(g.items, item)
	g.lock.Unlock()
	return item
}

// OnClick sets the callback for a button or toolbar-group member.
func (i *MacToolbarItem) OnClick(callback func(*Context)) *MacToolbarItem {
	if i != nil {
		i.lock.Lock()
		i.onClick = callback
		i.lock.Unlock()
	}
	return i
}

// OnSearch sets the callback invoked when the user submits a search.
func (i *MacToolbarItem) OnSearch(callback func(*Context, string)) *MacToolbarItem {
	if i != nil {
		i.lock.Lock()
		i.onSearch = callback
		i.lock.Unlock()
	}
	return i
}

func (i *MacToolbarItem) SetLabel(label string) *MacToolbarItem {
	if i == nil {
		return i
	}
	i.lock.Lock()
	i.label = label
	i.lock.Unlock()
	i.update(func(native unsafe.Pointer) { macToolbarItemSetLabel(native, i.identifier, label) })
	return i
}

// SetSymbol sets the item's SF Symbol image on macOS 11 and newer. Passing an
// empty string removes the image.
func (i *MacToolbarItem) SetSymbol(symbol string) *MacToolbarItem {
	if i == nil {
		return i
	}
	i.lock.Lock()
	i.symbolName = symbol
	i.lock.Unlock()
	i.update(func(native unsafe.Pointer) { macToolbarItemSetSymbol(native, i.identifier, symbol) })
	return i
}

func (i *MacToolbarItem) SetTooltip(tooltip string) *MacToolbarItem {
	if i == nil {
		return i
	}
	i.lock.Lock()
	i.tooltip = tooltip
	i.lock.Unlock()
	i.update(func(native unsafe.Pointer) { macToolbarItemSetTooltip(native, i.identifier, tooltip) })
	return i
}

// SetBordered controls the item's bordered presentation on macOS 10.15 and
// newer.
func (i *MacToolbarItem) SetBordered(bordered bool) *MacToolbarItem {
	if i == nil {
		return i
	}
	i.lock.Lock()
	i.bordered = bordered
	i.lock.Unlock()
	i.update(func(native unsafe.Pointer) { macToolbarItemSetBordered(native, i.identifier, bordered) })
	return i
}

// SetProminent controls the prominent toolbar-item style on macOS 26 and
// newer.
func (i *MacToolbarItem) SetProminent(prominent bool) *MacToolbarItem {
	if i == nil {
		return i
	}
	i.lock.Lock()
	i.prominent = prominent
	i.lock.Unlock()
	i.update(func(native unsafe.Pointer) { macToolbarItemSetProminent(native, i.identifier, prominent) })
	return i
}

// SetTintColor controls the toolbar-item background tint on macOS 26 and
// newer. Passing nil removes the tint.
func (i *MacToolbarItem) SetTintColor(color *RGBA) *MacToolbarItem {
	if i == nil {
		return i
	}
	var stored *RGBA
	if color != nil {
		copyOfColor := *color
		stored = &copyOfColor
	}
	i.lock.Lock()
	i.tintColor = stored
	i.lock.Unlock()
	i.update(func(native unsafe.Pointer) { macToolbarItemSetTintColor(native, i.identifier, stored) })
	return i
}

func (i *MacToolbarItem) SetEnabled(enabled bool) *MacToolbarItem {
	if i == nil {
		return i
	}
	i.lock.Lock()
	i.disabled = !enabled
	i.lock.Unlock()
	i.update(func(native unsafe.Pointer) { macToolbarItemSetEnabled(native, i.identifier, enabled) })
	return i
}

func (i *MacToolbarItem) SetHidden(hidden bool) *MacToolbarItem {
	if i == nil {
		return i
	}
	i.lock.Lock()
	i.hidden = hidden
	i.lock.Unlock()
	i.update(func(native unsafe.Pointer) { macToolbarItemSetHidden(native, i.identifier, hidden) })
	return i
}

// SetBadgeCount sets the toolbar-item badge on macOS 26 and newer. Values less
// than zero are treated as zero.
func (i *MacToolbarItem) SetBadgeCount(count int) *MacToolbarItem {
	if i == nil {
		return i
	}
	if count < 0 {
		count = 0
	}
	i.lock.Lock()
	i.badgeCount = count
	i.lock.Unlock()
	i.update(func(native unsafe.Pointer) { macToolbarItemSetBadgeCount(native, i.identifier, count) })
	return i
}

// SetSelectedIndex updates the selected member of a select-one group. Invalid
// indices are ignored. Use -1 to clear the selection.
func (g *MacToolbarGroup) SetSelectedIndex(index int) *MacToolbarGroup {
	if g == nil || g.MacToolbarItem == nil {
		return g
	}
	g.lock.Lock()
	if index < -1 || index >= len(g.items) {
		g.lock.Unlock()
		return g
	}
	g.selectedIndex = index
	g.lock.Unlock()
	g.update(func(native unsafe.Pointer) { macToolbarGroupSetSelectedIndex(native, g.identifier, index) })
	return g
}

// SetSelectionMode changes a group's native selection behavior.
func (g *MacToolbarGroup) SetSelectionMode(mode MacToolbarGroupSelectionMode) *MacToolbarGroup {
	if g == nil || g.MacToolbarItem == nil {
		return g
	}
	g.lock.Lock()
	g.selectionMode = mode
	g.lock.Unlock()
	g.update(func(native unsafe.Pointer) { macToolbarGroupSetSelectionMode(native, g.identifier, mode) })
	return g
}

func (i *MacToolbarItem) update(update func(unsafe.Pointer)) {
	if i == nil || i.toolbar == nil {
		return
	}

	// Avoid scheduling work before installation while still rechecking under
	// the same lock on the application thread before dereferencing native.
	i.toolbar.stateLock.RLock()
	installed := i.toolbar.state != nil && i.toolbar.state.native != nil
	i.toolbar.stateLock.RUnlock()
	if !installed {
		return
	}

	InvokeSync(func() {
		i.toolbar.stateLock.RLock()
		defer i.toolbar.stateLock.RUnlock()
		if i.toolbar.state != nil && i.toolbar.state.native != nil {
			update(i.toolbar.state.native)
		}
	})
}

// macToolbarState is protected exclusively by MacToolbar.stateLock.
type macToolbarState struct {
	window  *WebviewWindow
	native  unsafe.Pointer
	itemIDs []uint
}

func claimMacToolbar(toolbar *MacToolbar, window *WebviewWindow) (bool, error) {
	toolbar.stateLock.Lock()
	defer toolbar.stateLock.Unlock()
	if toolbar.state == nil {
		toolbar.state = &macToolbarState{window: window}
		return true, nil
	}
	if toolbar.state.window != nil && toolbar.state.window != window {
		return false, fmt.Errorf("toolbar is already attached to another window")
	}
	claimed := toolbar.state.window == nil
	toolbar.state.window = window
	return claimed, nil
}

func releaseMacToolbarOwnership(toolbar *MacToolbar, window *WebviewWindow) {
	if toolbar == nil {
		return
	}
	toolbar.stateLock.Lock()
	if toolbar.state != nil && toolbar.state.window == window && toolbar.state.native == nil {
		toolbar.state.window = nil
	}
	toolbar.stateLock.Unlock()
}

func validateToolbarItems(items []*MacToolbarItem) error {
	var validate func([]*MacToolbarItem, string) error
	validate = func(items []*MacToolbarItem, parent string) error {
		for _, item := range items {
			if item == nil {
				return fmt.Errorf("toolbar item in %q is nil", parent)
			}
			item.lock.RLock()
			kind := item.kind
			label := item.label
			onClick := item.onClick
			onSearch := item.onSearch
			members := append([]*MacToolbarItem(nil), item.items...)
			item.lock.RUnlock()

			switch kind {
			case toolbarButton:
				if onClick == nil {
					return fmt.Errorf("toolbar button %q requires OnClick", label)
				}
			case toolbarSearchField:
				if onSearch == nil {
					return fmt.Errorf("toolbar search field %q requires OnSearch", label)
				}
			case toolbarGroup:
				if len(members) == 0 {
					return fmt.Errorf("toolbar group %q requires at least one member", label)
				}
				for _, member := range members {
					if member == nil {
						return fmt.Errorf("toolbar group %q contains a nil member", label)
					}
					member.lock.RLock()
					valid := member.kind == toolbarButton && member.onClick != nil
					member.lock.RUnlock()
					if !valid {
						return fmt.Errorf("toolbar group %q members must be buttons with OnClick", label)
					}
				}
				if err := validate(members, label); err != nil {
					return err
				}
			case toolbarFlexibleSpace:
			default:
				return fmt.Errorf("toolbar item %q has unknown kind %d", label, kind)
			}
		}
		return nil
	}
	return validate(items, "root")
}

var toolbarItemMap = make(map[uint]*MacToolbarItem)
var toolbarItemMapLock sync.RWMutex
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
	toolbarItemMapLock.RLock()
	defer toolbarItemMapLock.RUnlock()
	return toolbarItemMap[id]
}

var toolbarItemClicked = make(chan uint, 32)

type toolbarSearchEvent struct {
	itemID uint
	query  string
}

var toolbarSearchTriggered = make(chan toolbarSearchEvent, 32)

func handleToolbarItemClicked(itemID uint) {
	defer handlePanic()
	item := getToolbarItemByID(itemID)
	if item == nil {
		return
	}
	item.lock.RLock()
	callback := item.onClick
	item.lock.RUnlock()
	if callback != nil {
		callback(newContext())
	}
}

func handleToolbarSearch(itemID uint, query string) {
	defer handlePanic()
	item := getToolbarItemByID(itemID)
	if item == nil {
		return
	}
	item.lock.RLock()
	callback := item.onSearch
	item.lock.RUnlock()
	if callback != nil {
		callback(newContext(), query)
	}
}
