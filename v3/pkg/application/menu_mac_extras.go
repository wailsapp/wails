package application

import (
	"path/filepath"
	"runtime"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/events"
)

// This file holds the cross-platform model for the macOS menu and dock
// extras: SF Symbol images, NSMenuItemBadge, section headers, palette menus,
// mixed state, alternates, indentation, the dock menu and the Open Recent
// role. The native side lives in menu_mac_extras_darwin.*; on other
// platforms the state is kept but has no visual effect.

// menuItemSymbolSetter is implemented by platform menu items that can show
// a system symbol image (SF Symbols on macOS 11+).
type menuItemSymbolSetter interface {
	setSymbol(name string)
}

// menuItemBadgeSetter is implemented by platform menu items that support a
// trailing badge (NSMenuItemBadge on macOS 14+).
type menuItemBadgeSetter interface {
	setBadge(text string, count int, present bool)
}

// menuItemStateExtras is implemented by platform menu items that support
// the mixed check state, alternate items and indentation.
type menuItemStateExtras interface {
	setMixed(mixed bool)
	setAlternate(alternate bool)
	setIndentationLevel(level int)
}

// recentDocumentsImpl is implemented by platform application impls that
// keep a native recent documents list (NSDocumentController on macOS).
type recentDocumentsImpl interface {
	addRecentDocument(path string)
	clearRecentDocuments()
	recentDocuments() []string
}

// maxRecentDocuments mirrors NSDocumentController's default
// maximumRecentDocumentCount.
const maxRecentDocuments = 10

// SetSymbol sets a system symbol image (an SF Symbol name such as
// "star.fill") on the menu item. Supported on macOS 11 and later; an empty
// name clears the image. On other platforms the name is stored but ignored.
// A symbol set after SetBitmap replaces the bitmap image.
func (m *MenuItem) SetSymbol(name string) *MenuItem {
	m.symbol = name
	if setter, ok := m.impl.(menuItemSymbolSetter); ok && m.impl != nil {
		setter.setSymbol(name)
	}
	return m
}

// Symbol returns the system symbol name set with SetSymbol.
func (m *MenuItem) Symbol() string {
	return m.symbol
}

// SetBadge shows a numeric badge (NSMenuItemBadge) next to the menu item
// title on macOS 14 and later. A count of zero or less clears the badge.
func (m *MenuItem) SetBadge(count int) *MenuItem {
	if count <= 0 {
		return m.ClearBadge()
	}
	m.badgeCount = count
	m.badgeText = ""
	m.hasBadge = true
	m.applyBadge()
	return m
}

// SetBadgeText shows a text badge next to the menu item title on macOS 14
// and later. An empty string clears the badge.
func (m *MenuItem) SetBadgeText(text string) *MenuItem {
	if text == "" {
		return m.ClearBadge()
	}
	m.badgeText = text
	m.badgeCount = 0
	m.hasBadge = true
	m.applyBadge()
	return m
}

// ClearBadge removes the badge set with SetBadge or SetBadgeText.
func (m *MenuItem) ClearBadge() *MenuItem {
	m.badgeText = ""
	m.badgeCount = 0
	m.hasBadge = false
	m.applyBadge()
	return m
}

// HasBadge reports whether a badge is set on the menu item.
func (m *MenuItem) HasBadge() bool {
	return m.hasBadge
}

// BadgeCount returns the numeric badge value, or 0 when the badge is text
// or not set.
func (m *MenuItem) BadgeCount() int {
	return m.badgeCount
}

// BadgeText returns the text badge value, or "" when the badge is numeric
// or not set.
func (m *MenuItem) BadgeText() string {
	return m.badgeText
}

func (m *MenuItem) applyBadge() {
	if setter, ok := m.impl.(menuItemBadgeSetter); ok && m.impl != nil {
		setter.setBadge(m.badgeText, m.badgeCount, m.hasBadge)
	}
}

// SetMixed puts a checkbox item into the mixed (partially checked) state,
// rendered as a dash on macOS. While mixed, Checked() reports false and
// Mixed() reports true. SetChecked, or a click on the item, leaves the mixed
// state: a click on a mixed checkbox turns it fully on, matching AppKit.
func (m *MenuItem) SetMixed() *MenuItem {
	m.mixed = true
	m.checked = false
	if extras, ok := m.impl.(menuItemStateExtras); ok && m.impl != nil {
		extras.setMixed(true)
	}
	return m
}

// Mixed reports whether the item is in the mixed (partially checked) state.
func (m *MenuItem) Mixed() bool {
	return m.mixed
}

// SetAlternate marks the item as an alternate of the item directly above
// it. AppKit shows the alternate in place of the previous item while the
// modifier keys that differ between their key equivalents are held, so the
// two items must share a key and differ in modifiers (for example "Cmd+w"
// and "Cmd+OptionOrAlt+w"). Without a differing modifier both items are
// shown. macOS only.
func (m *MenuItem) SetAlternate(alternate bool) *MenuItem {
	m.alternate = alternate
	if extras, ok := m.impl.(menuItemStateExtras); ok && m.impl != nil {
		extras.setAlternate(alternate)
	}
	return m
}

// Alternate reports whether the item was marked with SetAlternate.
func (m *MenuItem) Alternate() bool {
	return m.alternate
}

// SetIndentationLevel indents the item title by the given number of levels
// (0 to 15). macOS only.
func (m *MenuItem) SetIndentationLevel(level int) *MenuItem {
	if level < 0 {
		level = 0
	}
	if level > 15 {
		level = 15
	}
	m.indentationLevel = level
	if extras, ok := m.impl.(menuItemStateExtras); ok && m.impl != nil {
		extras.setIndentationLevel(level)
	}
	return m
}

// IndentationLevel returns the indentation set with SetIndentationLevel.
func (m *MenuItem) IndentationLevel() int {
	return m.indentationLevel
}

// IsSectionHeader reports whether the item was created by AddSectionHeader.
func (m *MenuItem) IsSectionHeader() bool {
	return m.itemType == sectionHeader
}

// IsPalette reports whether the item was created by AddPalette.
func (m *MenuItem) IsPalette() bool {
	return m.itemType == palette
}

// PaletteSelected returns the index of the selected palette swatch, or -1
// when nothing is selected. Only meaningful for items created by AddPalette.
func (m *MenuItem) PaletteSelected() int {
	if m.itemType != palette {
		return -1
	}
	return m.paletteSelected
}

// NewMenuItemSectionHeader creates a section header item. On macOS 14 and
// later it is rendered as an NSMenuItem section header; elsewhere it is a
// disabled item showing the title.
func NewMenuItemSectionHeader(title string) *MenuItem {
	result := NewMenuItem(title)
	result.itemType = sectionHeader
	result.disabled = true
	return result
}

// AddSectionHeader appends a section header (macOS 14 and later) to the
// menu and returns it. On older systems and other platforms the header is a
// disabled item showing the title.
func (m *Menu) AddSectionHeader(title string) *MenuItem {
	result := NewMenuItemSectionHeader(title)
	m.items = append(m.items, result)
	return result
}

// AddPalette appends a palette of colour swatches (macOS 14 and later)
// backed by NSMenu paletteMenuWithColors:. colours gives one swatch per
// entry. symbols optionally gives the SF Symbol drawn in each swatch: one
// name applies to every swatch, one name per colour applies individually,
// and an empty slice draws filled circles. selected is the initially
// selected index (-1 for none). onSelect receives the clicked index.
//
// The palette is inline in the parent menu when the returned item has no
// label; give it a label with SetLabel to present it as a titled submenu.
// On systems older than macOS 14 and on other platforms the item is hidden.
func (m *Menu) AddPalette(symbols []string, colours []RGBA, selected int, onSelect func(*Context, int)) *MenuItem {
	result := NewMenuItem("")
	result.itemType = palette
	result.paletteSymbols = append([]string(nil), symbols...)
	result.paletteColours = append([]RGBA(nil), colours...)
	if selected < 0 || selected >= len(colours) {
		selected = -1
	}
	result.paletteSelected = selected
	result.paletteCallback = onSelect
	if runtime.GOOS != "darwin" {
		result.hidden = true
	}
	m.items = append(m.items, result)
	return result
}

// handlePaletteSelection records a swatch selection and runs the callback.
func (m *MenuItem) handlePaletteSelection(index int) {
	if m.itemType != palette {
		return
	}
	if index < 0 || index >= len(m.paletteColours) {
		index = -1
	}
	m.paletteSelected = index
	if m.paletteCallback != nil && index >= 0 {
		ctx := newContext().
			withClickedMenuItem(m).
			withContextMenuData(m.contextMenuData)
		go func() {
			defer handlePanic()
			m.paletteCallback(ctx, index)
		}()
	}
}

// SetDockMenu sets the menu shown when the user right-clicks the app's Dock
// icon (macOS). Item callbacks fire as for any other menu. Pass nil to
// remove it. A menu set with OnDockMenu takes precedence.
func (mm *MenuManager) SetDockMenu(menu *Menu) {
	mm.dockMenuLock.Lock()
	mm.dockMenu = menu
	mm.dockMenuLock.Unlock()
}

// DockMenu returns the menu set with SetDockMenu.
func (mm *MenuManager) DockMenu() *Menu {
	mm.dockMenuLock.Lock()
	defer mm.dockMenuLock.Unlock()
	return mm.dockMenu
}

// OnDockMenu registers a function that builds the Dock menu each time it is
// about to be shown (macOS). Returning nil shows no menu. The previously
// built menu is destroyed when the next one is requested, unless the
// function returns the same *Menu again. Pass nil to unregister.
func (mm *MenuManager) OnDockMenu(fn func() *Menu) {
	mm.dockMenuLock.Lock()
	mm.dockMenuFunc = fn
	mm.dockMenuLock.Unlock()
}

// resolveDockMenu returns the menu to show for the Dock, preferring the
// dynamic builder over the static menu.
func (mm *MenuManager) resolveDockMenu() *Menu {
	mm.dockMenuLock.Lock()
	fn := mm.dockMenuFunc
	static := mm.dockMenu
	mm.dockMenuLock.Unlock()
	if fn == nil {
		return static
	}
	menu := fn()
	mm.dockMenuLock.Lock()
	previous := mm.lastDynamicDockMenu
	mm.lastDynamicDockMenu = menu
	mm.dockMenuLock.Unlock()
	if previous != nil && previous != menu && previous != static {
		previous.Destroy()
	}
	return menu
}

// AddRecentDocument adds a file to the application's recent documents list
// (NSDocumentController on macOS, shown by the OpenRecent role and the Dock
// menu). The path is cleaned and moved to the front if already present.
func (mm *MenuManager) AddRecentDocument(path string) {
	if path == "" {
		return
	}
	path = filepath.Clean(path)
	mm.recentLock.Lock()
	list := []string{path}
	for _, existing := range mm.recentDocuments {
		if existing != path {
			list = append(list, existing)
		}
	}
	if len(list) > maxRecentDocuments {
		list = list[:maxRecentDocuments]
	}
	mm.recentDocuments = list
	mm.recentLock.Unlock()
	if impl, ok := mm.app.impl.(recentDocumentsImpl); ok && mm.app.impl != nil {
		impl.addRecentDocument(path)
	}
}

// ClearRecentDocuments empties the recent documents list.
func (mm *MenuManager) ClearRecentDocuments() {
	mm.recentLock.Lock()
	mm.recentDocuments = nil
	mm.recentLock.Unlock()
	if impl, ok := mm.app.impl.(recentDocumentsImpl); ok && mm.app.impl != nil {
		impl.clearRecentDocuments()
	}
}

// RecentDocuments returns the recent documents list, most recent first. On
// macOS this is read from NSDocumentController, so it survives restarts and
// includes files the system noted on the app's behalf.
func (mm *MenuManager) RecentDocuments() []string {
	if impl, ok := mm.app.impl.(recentDocumentsImpl); ok && mm.app.impl != nil {
		if native := impl.recentDocuments(); native != nil {
			return native
		}
	}
	mm.recentLock.Lock()
	defer mm.recentLock.Unlock()
	return append([]string(nil), mm.recentDocuments...)
}

// OpenRecentDocument routes a recent document through the same path as a
// file opened from Finder: it is noted as recent again and delivered as an
// ApplicationOpenedWithFile event.
func (mm *MenuManager) OpenRecentDocument(path string) {
	if path == "" {
		return
	}
	mm.AddRecentDocument(path)
	eventContext := newApplicationEventContext()
	eventContext.setOpenedWithFile(path)
	applicationEvents <- &ApplicationEvent{
		Id:  uint(events.Common.ApplicationOpenedWithFile),
		ctx: eventContext,
	}
}

// NewOpenRecentMenuItem creates the "Open Recent" submenu for the OpenRecent
// role. On macOS the submenu is populated natively from NSDocumentController
// every time it opens, with a "Clear Menu" item. Elsewhere it is filled from
// the Go recent documents list when created.
func NewOpenRecentMenuItem() *MenuItem {
	item := NewSubMenuItem("Open Recent")
	item.role = OpenRecent
	if runtime.GOOS == "darwin" {
		return item
	}
	var recents []string
	if app := globalApplication; app != nil && app.Menu != nil {
		recents = app.Menu.RecentDocuments()
	}
	for _, path := range recents {
		path := path
		item.submenu.Add(filepath.Base(path)).SetTooltip(path).OnClick(func(*Context) {
			if app := globalApplication; app != nil {
				app.Menu.OpenRecentDocument(path)
			}
		})
	}
	if len(recents) > 0 {
		item.submenu.AddSeparator()
	}
	item.submenu.Add("Clear Menu").SetEnabled(len(recents) > 0).OnClick(func(*Context) {
		if app := globalApplication; app != nil {
			app.Menu.ClearRecentDocuments()
		}
	})
	return item
}

// menuManagerExtras holds the dock menu and recent documents state embedded
// in MenuManager.
type menuManagerExtras struct {
	dockMenuLock        sync.Mutex
	dockMenu            *Menu
	dockMenuFunc        func() *Menu
	lastDynamicDockMenu *Menu

	recentLock      sync.Mutex
	recentDocuments []string
}
