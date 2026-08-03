package application

import (
	"sync"
	"sync/atomic"
)

// MacToolbar describes a native macOS window toolbar (NSToolbar). Attach it
// with WebviewWindow.SetToolbar. macOS only; a no-op on other platforms.
type MacToolbar struct {
	Items []MacToolbarItem
}

// MacToolbarItemKind selects what a MacToolbarItem renders as.
type MacToolbarItemKind int

const (
	// ToolbarButton is a plain clickable item (NSToolbarItem). Requires OnClick.
	ToolbarButton MacToolbarItemKind = iota
	// ToolbarGroup renders Items as a single segmented control (NSToolbarItemGroup).
	ToolbarGroup
	// ToolbarSearchField is a native search field (NSSearchToolbarItem). Requires OnSearch.
	ToolbarSearchField
	// ToolbarFlexibleSpace is a flexible spacer (NSToolbarFlexibleSpaceItem).
	ToolbarFlexibleSpace
	// ToolbarSidebarToggle is the standard sidebar-toggle action, plus (if
	// TracksSplitPane is set) an NSTrackingSeparatorToolbarItem bound to
	// that pane's divider.
	ToolbarSidebarToggle
)

// MacToolbarItem is one item in a MacToolbar.
type MacToolbarItem struct {
	// ID is the NSToolbarItem identifier. Must be unique within the toolbar.
	ID string

	Kind MacToolbarItemKind

	Label      string
	SymbolName string // SF Symbol name
	Bordered   bool
	Prominent  bool
	TintColor  *RGBA
	BadgeCount int // 0 = no badge

	// Items holds the member items when Kind == ToolbarGroup.
	Items []MacToolbarItem

	// TracksSplitPane, for Kind == ToolbarSidebarToggle: the split-pane
	// Name whose divider this item's tracking separator follows.
	TracksSplitPane string

	// OnClick is required for Kind == ToolbarButton (and for each member
	// item of a ToolbarGroup with the default per-item click behaviour).
	OnClick func(*Context)
	// OnSearch is required for Kind == ToolbarSearchField.
	OnSearch func(*Context, string)
}

// Toolbar item click/search dispatch: same id-tagged-native-object +
// narrow //export bridge + channel + central-dispatch shape as menu items
// (see menuitem.go's menuItemID/menuItemMap and application.go's
// menuItemClicked dispatcher). Declared here, cross-platform, so
// application.go's startup dispatcher registration compiles on every
// platform; only the darwin implementation ever sends on these channels or
// populates this map.

var toolbarItemID uintptr
var toolbarItemMap = make(map[uint]*MacToolbarItem)
var toolbarItemMapLock sync.Mutex

func nextToolbarItemID() uint {
	return uint(atomic.AddUintptr(&toolbarItemID, 1))
}

func addToToolbarItemMap(id uint, item *MacToolbarItem) {
	toolbarItemMapLock.Lock()
	toolbarItemMap[id] = item
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
	item := getToolbarItemByID(itemID)
	if item == nil {
		globalApplication.warning("Toolbar item #%d not found", itemID)
		return
	}
	if item.OnClick != nil {
		item.OnClick(newContext())
	}
}

func handleToolbarSearch(itemID uint, query string) {
	defer handlePanic()
	item := getToolbarItemByID(itemID)
	if item == nil {
		globalApplication.warning("Toolbar item #%d not found", itemID)
		return
	}
	if item.OnSearch != nil {
		item.OnSearch(newContext(), query)
	}
}
