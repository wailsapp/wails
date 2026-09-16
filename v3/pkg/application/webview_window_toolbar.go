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
	stateLock   sync.RWMutex
	state       *macToolbarState
	displayMode MacToolbarDisplayMode

	// customizable enables the standard "Customize Toolbar..." sheet and
	// autosave. persistenceKey becomes the NSToolbar identifier so the saved
	// layout is stable across launches. Both are protected by stateLock.
	customizable   bool
	persistenceKey string

	// centeredItems lists the items AppKit should keep centred (macOS 13+).
	// Protected by itemsLock.
	centeredItems []*MacToolbarItem
}

var toolbarIdentifier uint64

// NewMacToolbar creates an empty native macOS toolbar.
func NewMacToolbar() *MacToolbar {
	sequence := atomic.AddUint64(&toolbarIdentifier, 1)
	return &MacToolbar{
		identifier:  fmt.Sprintf("wails.toolbar.%d", sequence),
		displayMode: MacToolbarDisplayModeIconAndLabel,
	}
}

// SetCustomizable enables the standard AppKit "Customize Toolbar..." sheet
// (allowsUserCustomization) and automatic layout persistence
// (autosavesConfiguration). persistenceKey becomes the NSToolbar identifier
// under which AppKit stores the user's layout, so it must be unique per
// toolbar within the application and stable across launches. Items keep their
// generated identifiers unless SetPersistenceKey gives them a stable one, and
// only stable item identifiers survive a relaunch. Passing an empty key
// disables customisation again.
//
// Call SetCustomizable before the toolbar is attached: an NSToolbar's
// identifier is fixed at creation, so a key set afterwards is applied the next
// time the toolbar is attached while the customisation sheet is enabled
// immediately.
func (t *MacToolbar) SetCustomizable(persistenceKey string) *MacToolbar {
	if t == nil {
		return t
	}
	t.stateLock.Lock()
	t.customizable = persistenceKey != ""
	t.persistenceKey = persistenceKey
	customizable := t.customizable
	installed := t.state != nil && t.state.native != nil
	t.stateLock.Unlock()
	if !installed {
		return t
	}
	if customizable {
		t.logDebug("SetCustomizable: the toolbar is already attached; persistence key %q applies on the next attachment", persistenceKey)
	}
	InvokeSync(func() {
		t.stateLock.RLock()
		defer t.stateLock.RUnlock()
		if t.state != nil && t.state.native != nil {
			macToolbarSetCustomizable(t.state.native, t.customizable)
		}
	})
	t.syncNative(nil)
	return t
}

// RunCustomizationPalette opens the standard AppKit toolbar customisation
// sheet. It does nothing unless the toolbar is attached and SetCustomizable
// has been called.
func (t *MacToolbar) RunCustomizationPalette() *MacToolbar {
	if t == nil {
		return t
	}
	t.stateLock.RLock()
	installed := t.state != nil && t.state.native != nil
	customizable := t.customizable
	t.stateLock.RUnlock()
	if !installed || !customizable {
		t.logDebug("RunCustomizationPalette: the toolbar must be attached and customizable")
		return t
	}
	InvokeSync(func() {
		t.stateLock.RLock()
		defer t.stateLock.RUnlock()
		if t.state != nil && t.state.native != nil {
			macToolbarRunCustomizationPalette(t.state.native)
		}
	})
	return t
}

// SetCenteredItems asks AppKit to keep the given items centred in the toolbar
// (NSToolbar.centeredItemIdentifiers, macOS 13 and newer; earlier releases
// keep the standard layout). Pass no items to clear the centred set. Items
// must belong to this toolbar; others are ignored.
func (t *MacToolbar) SetCenteredItems(items ...*MacToolbarItem) *MacToolbar {
	if t == nil {
		return t
	}
	centered := make([]*MacToolbarItem, 0, len(items))
	for _, item := range items {
		if item != nil && item.toolbar == t && item.parent == nil {
			centered = append(centered, item)
		}
	}
	t.itemsLock.Lock()
	t.centeredItems = centered
	t.itemsLock.Unlock()
	t.syncNative(nil)
	return t
}

// SetDefaultItems marks the given top-level items as the toolbar's default
// layout: every other item remains available in the customisation palette
// (see SetCustomizable) but is not shown until the user adds it. It is
// equivalent to calling SetInDefaultSet on every item. Passing no items
// restores every item to the default layout.
func (t *MacToolbar) SetDefaultItems(items ...*MacToolbarItem) *MacToolbar {
	if t == nil {
		return t
	}
	defaults := make(map[*MacToolbarItem]bool, len(items))
	for _, item := range items {
		defaults[item] = true
	}
	for _, item := range t.itemSnapshot() {
		item.lock.Lock()
		item.inDefaultSet = len(items) == 0 || defaults[item]
		item.lock.Unlock()
	}
	t.syncNative(nil)
	return t
}

// Remove detaches an item (or a group member) from the toolbar. An attached
// toolbar drops the native item immediately. The item cannot be re-added; its
// setters become no-ops.
func (t *MacToolbar) Remove(item *MacToolbarItem) *MacToolbar {
	if t == nil || item == nil || item.toolbar != t {
		return t
	}
	if item.parent != nil {
		parent := item.parent
		parent.lock.Lock()
		parent.items = removeMacToolbarItem(parent.items, item)
		parent.lock.Unlock()
	} else {
		t.itemsLock.Lock()
		t.items = removeMacToolbarItem(t.items, item)
		t.centeredItems = removeMacToolbarItem(t.centeredItems, item)
		t.itemsLock.Unlock()
	}
	t.syncNative(nil)
	return t
}

// Move places an item (or a group member) at index within its container.
// Indices outside the container are clamped. An attached toolbar reorders the
// native item immediately.
func (t *MacToolbar) Move(item *MacToolbarItem, index int) *MacToolbar {
	if t == nil || item == nil || item.toolbar != t {
		return t
	}
	if item.parent != nil {
		parent := item.parent
		parent.lock.Lock()
		parent.items = moveMacToolbarItem(parent.items, item, index)
		parent.lock.Unlock()
	} else {
		t.itemsLock.Lock()
		t.items = moveMacToolbarItem(t.items, item, index)
		t.itemsLock.Unlock()
	}
	t.syncNative(item)
	return t
}

func removeMacToolbarItem(items []*MacToolbarItem, item *MacToolbarItem) []*MacToolbarItem {
	result := make([]*MacToolbarItem, 0, len(items))
	for _, candidate := range items {
		if candidate != item {
			result = append(result, candidate)
		}
	}
	return result
}

func moveMacToolbarItem(items []*MacToolbarItem, item *MacToolbarItem, index int) []*MacToolbarItem {
	found := false
	for _, candidate := range items {
		if candidate == item {
			found = true
			break
		}
	}
	if !found {
		return items
	}
	result := removeMacToolbarItem(items, item)
	if index < 0 {
		index = 0
	}
	if index > len(result) {
		index = len(result)
	}
	result = append(result, nil)
	copy(result[index+1:], result[index:])
	result[index] = item
	return result
}

// syncNative re-synchronises an attached native toolbar with the Go model
// after a structural change (an item was added, removed, moved, renamed with
// a persistence key, or its default-set membership changed). Callback
// requirements are not enforced here because callbacks are normally chained
// after the constructor; a click on an item without a callback is reported
// when it happens. Structural problems are logged rather than failing
// silently.
func (t *MacToolbar) syncNative(moved *MacToolbarItem) {
	if t == nil {
		return
	}
	t.stateLock.RLock()
	installed := t.state != nil && t.state.native != nil
	t.stateLock.RUnlock()
	if !installed {
		return
	}
	if err := validateToolbarLayout(t.itemSnapshot(), false); err != nil {
		t.logError("toolbar: %s", err)
	}
	InvokeSync(func() {
		t.stateLock.Lock()
		defer t.stateLock.Unlock()
		if t.state == nil || t.state.native == nil {
			return
		}
		macToolbarSyncNative(t, moved)
	})
}

// logError reports a toolbar problem through the owning window when the
// toolbar is attached, and through the application logger otherwise.
func (t *MacToolbar) logError(format string, args ...any) {
	if t == nil || globalApplication == nil {
		return
	}
	t.stateLock.RLock()
	var window macToolbarWindow
	if t.state != nil {
		window = t.state.window
	}
	t.stateLock.RUnlock()
	if window != nil {
		window.Error(format, args...)
		return
	}
	globalApplication.error(format, args...)
}

func (t *MacToolbar) logDebug(format string, args ...any) {
	if t == nil || globalApplication == nil {
		return
	}
	globalApplication.debug(fmt.Sprintf(format, args...))
}

// MacToolbarDisplayMode controls whether an NSToolbar shows item icons,
// labels, or both. Individual controls such as search fields keep their
// native AppKit presentation in every mode.
type MacToolbarDisplayMode int

const (
	// MacToolbarDisplayModeDefault lets AppKit choose the presentation.
	MacToolbarDisplayModeDefault MacToolbarDisplayMode = iota
	// MacToolbarDisplayModeIconAndLabel shows both icons and labels. This is
	// the Wails default for compatibility with existing toolbars.
	MacToolbarDisplayModeIconAndLabel
	// MacToolbarDisplayModeIconOnly shows toolbar icons without labels.
	MacToolbarDisplayModeIconOnly
	// MacToolbarDisplayModeLabelOnly shows labels without icons.
	MacToolbarDisplayModeLabelOnly
)

// SetDisplayMode sets the toolbar's native AppKit display mode. It may be
// called before or after the toolbar is attached; live changes are applied on
// the application thread. Invalid values are ignored.
func (t *MacToolbar) SetDisplayMode(mode MacToolbarDisplayMode) *MacToolbar {
	if t == nil || mode < MacToolbarDisplayModeDefault || mode > MacToolbarDisplayModeLabelOnly {
		return t
	}
	t.stateLock.Lock()
	t.displayMode = mode
	installed := t.state != nil && t.state.native != nil
	t.stateLock.Unlock()
	if !installed {
		return t
	}
	InvokeSync(func() {
		t.stateLock.RLock()
		defer t.stateLock.RUnlock()
		if t.state != nil && t.state.native != nil {
			macToolbarSetDisplayMode(t.state.native, t.displayMode)
		}
	})
	return t
}

type macToolbarItemKind int

const (
	toolbarButton macToolbarItemKind = iota
	toolbarGroup
	toolbarSearchField
	toolbarShare
	toolbarFlexibleSpace
	toolbarSidebarToggle
	toolbarSidebarTrackingSeparator
	toolbarInspectorToggle
	toolbarInspectorTrackingSeparator
	toolbarMenu
	toolbarSpace
)

// MacToolbarVisibilityPriority controls which toolbar items AppKit moves into
// the overflow menu first when the window is too narrow. Values map directly
// to NSToolbarItemVisibilityPriority.
type MacToolbarVisibilityPriority int

const (
	// MacToolbarVisibilityPriorityStandard is the default priority.
	MacToolbarVisibilityPriorityStandard MacToolbarVisibilityPriority = 0
	// MacToolbarVisibilityPriorityLow items are the first to overflow.
	MacToolbarVisibilityPriorityLow MacToolbarVisibilityPriority = -1000
	// MacToolbarVisibilityPriorityHigh items overflow after standard items.
	MacToolbarVisibilityPriorityHigh MacToolbarVisibilityPriority = 1000
	// MacToolbarVisibilityPriorityUser is reserved for items the user has
	// explicitly prioritised; AppKit keeps them visible for as long as possible.
	MacToolbarVisibilityPriorityUser MacToolbarVisibilityPriority = 2000
)

// macToolbarDefaultSearchRecents is the NSSearchField maximumRecents applied
// when a recents key is configured.
const macToolbarDefaultSearchRecents = 10

// MacShareContentType is a Uniform Type Identifier (UTI) advertised to
// NSItemProvider. It tells receiving services how to interpret the bytes
// returned by MacShareProvider.ShareData. Custom UTIs are allowed in addition
// to the common values declared by this package.
type MacShareContentType string

const (
	// MacShareTypePlainText is UTF-8 encoded plain text.
	MacShareTypePlainText MacShareContentType = "public.utf8-plain-text"
	// MacShareTypeHTML is a UTF-8 encoded HTML document or fragment.
	MacShareTypeHTML MacShareContentType = "public.html"
	// MacShareTypePDF is a complete PDF document.
	MacShareTypePDF MacShareContentType = "com.adobe.pdf"
	// MacShareTypePNG is PNG-encoded image data.
	MacShareTypePNG MacShareContentType = "public.png"
	// MacShareTypeJPEG is JPEG-encoded image data.
	MacShareTypeJPEG MacShareContentType = "public.jpeg"
)

// MacShareRepresentation describes one format a provider can generate. The
// ContentType must match the bytes returned when ShareData receives that type.
type MacShareRepresentation struct {
	// ContentType is a Uniform Type Identifier such as MacShareTypePDF.
	ContentType MacShareContentType
}

// MacShareRequest identifies the representation AppKit requested at runtime.
// Switch on ContentType to decide which encoder or renderer to invoke.
type MacShareRequest struct {
	// ContentType is one of the types returned by ShareRepresentations.
	ContentType MacShareContentType
	// SuggestedName is the logical item name configured with SetSuggestedName.
	// It is shared by all representations and may be empty.
	SuggestedName string
}

// MacShareProvider lazily supplies bytes to the macOS share sheet. AppKit
// chooses among ShareRepresentations and Wails calls ShareData only when a
// receiving service requests one. A service may request more than one format;
// implementations may therefore be called concurrently and should be safe to
// repeat. ShareData should render from application state that is already safe
// to read rather than synchronously waiting for the WebView or AppKit thread.
//
// ShareRepresentations is called synchronously by MacToolbarShareItem.SetProvider.
// Its result is snapshotted until SetProvider is called again. ShareData is
// called lazily from a background callback and may return a rendering error;
// Wails forwards that failure to the native sharing service.
type MacShareProvider interface {
	// ShareRepresentations returns every representation this provider can
	// generate. Empty content types and duplicate types are ignored.
	ShareRepresentations() []MacShareRepresentation
	// ShareData returns complete bytes for the requested representation.
	ShareData(MacShareRequest) ([]byte, error)
}

// MacShareProviderFunc adapts a function into a MacShareProvider. It is useful
// for simple or stateless share items; stateful applications can implement
// MacShareProvider directly.
type MacShareProviderFunc struct {
	// Available lists every content type Load supports. The slice is copied
	// when the provider is installed.
	Available []MacShareRepresentation
	// Load generates the requested bytes. A nil Load returns an error to the
	// native sharing service.
	Load func(MacShareRequest) ([]byte, error)
}

func (p MacShareProviderFunc) ShareRepresentations() []MacShareRepresentation {
	return append([]MacShareRepresentation(nil), p.Available...)
}

func (p MacShareProviderFunc) ShareData(request MacShareRequest) ([]byte, error) {
	if p.Load == nil {
		return nil, fmt.Errorf("macOS share provider has no Load function")
	}
	return p.Load(request)
}

func normaliseMacShareRepresentations(representations []MacShareRepresentation) []MacShareRepresentation {
	result := make([]MacShareRepresentation, 0, len(representations))
	seen := make(map[MacShareContentType]struct{}, len(representations))
	for _, representation := range representations {
		if representation.ContentType == "" {
			continue
		}
		if _, exists := seen[representation.ContentType]; exists {
			continue
		}
		seen[representation.ContentType] = struct{}{}
		result = append(result, representation)
	}
	return result
}

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

	// identifier is the NSToolbarItemIdentifier: generatedIdentifier unless
	// SetPersistenceKey supplied a stable key.
	identifier          string
	generatedIdentifier string
	kind                macToolbarItemKind
	toolbar             *MacToolbar
	parent              *MacToolbarItem

	label      string
	symbolName string
	tooltip    string
	bordered   bool
	prominent  bool
	tintColor  *RGBA
	badgeCount int
	disabled   bool
	hidden     bool

	visibilityPriority MacToolbarVisibilityPriority
	navigational       bool
	inDefaultSet       bool

	menu           *Menu
	showsIndicator bool

	searchRecentsKey  string
	searchMenu        *Menu
	searchIncremental bool
	searchPlaceholder string

	items              []*MacToolbarItem
	selectionMode      MacToolbarGroupSelectionMode
	selectedIndex      int
	shareProvider      MacShareProvider
	shareFormats       []MacShareRepresentation
	shareSubject       string
	shareSuggestedName string
	onClick            func(*Context)
	onSearch           func(*Context, string)
	onShared           func(*Context, string)
	onShareError       func(*Context, string, error)
}

// MacToolbarGroup is a type-safe handle for a native segmented toolbar group.
// Its embedded MacToolbarItem supports the common item setters.
type MacToolbarGroup struct {
	*MacToolbarItem
}

// MacToolbarShareItem is a native macOS sharing-service toolbar item. Its
// embedded MacToolbarItem supports the common item setters.
type MacToolbarShareItem struct {
	*MacToolbarItem
}

var toolbarItemIdentifier uint64

func newMacToolbarItem(toolbar *MacToolbar, kind macToolbarItemKind, label string) *MacToolbarItem {
	sequence := atomic.AddUint64(&toolbarItemIdentifier, 1)
	identifier := fmt.Sprintf("wails.toolbar.item.%d", sequence)
	return &MacToolbarItem{
		identifier:          identifier,
		generatedIdentifier: identifier,
		kind:                kind,
		toolbar:             toolbar,
		label:               label,
		inDefaultSet:        true,
		showsIndicator:      true,
	}
}

// add appends a top-level item. When the toolbar is already attached the
// native toolbar is updated immediately; callbacks chained after the
// constructor take effect on the next click.
func (t *MacToolbar) add(item *MacToolbarItem) *MacToolbarItem {
	t.itemsLock.Lock()
	t.items = append(t.items, item)
	t.itemsLock.Unlock()
	t.syncNative(nil)
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

// AddShare adds the system sharing-service toolbar item. AppKit presents the
// native share sheet when it is clicked. The item is disabled until SetProvider
// advertises at least one representation. The returned handle owns callbacks
// and metadata; callers do not need to supply an identifier.
func (t *MacToolbar) AddShare(label string) *MacToolbarShareItem {
	item := newMacToolbarItem(t, toolbarShare, label)
	t.add(item)
	return &MacToolbarShareItem{MacToolbarItem: item}
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

// AddSpace adds the standard fixed-width space item
// (NSToolbarSpaceItemIdentifier).
func (t *MacToolbar) AddSpace() *MacToolbarItem {
	return t.add(newMacToolbarItem(t, toolbarSpace, ""))
}

// AddMenu adds a native dropdown toolbar item (NSMenuToolbarItem, macOS 10.15
// and newer; earlier releases omit the item) that opens menu when clicked.
// Menu item clicks fire the usual MenuItem.OnClick callbacks, and Menu.Update
// refreshes the native menu in place. Use SetSymbol for an icon and
// SetShowsIndicator to hide the dropdown chevron.
func (t *MacToolbar) AddMenu(label string, menu *Menu) *MacToolbarItem {
	item := newMacToolbarItem(t, toolbarMenu, label)
	item.menu = menu
	return t.add(item)
}

// AddSidebarToggle adds the standard AppKit toolbar item that sends
// toggleSidebar: through the responder chain. It requires no callback or ID
// because AppKit owns its action; collapse changes are observable on the
// sidebar pane through MacSplitPane.OnCollapsedChange. A toolbar may contain
// at most one sidebar toggle. AppKit owns this standard item, so there is no
// mutable Wails item handle to return.
func (t *MacToolbar) AddSidebarToggle() {
	t.add(newMacToolbarItem(t, toolbarSidebarToggle, ""))
}

// AddSidebarTrackingSeparator divides the toolbar into a leading sidebar
// section and a content section. AppKit keeps the leading section aligned
// above the first native sidebar divider as it moves, on macOS 11 and newer;
// older supported releases omit the separator while the sidebar and toggle
// keep working. Attaching a toolbar that contains a tracking separator to a
// window without a sidebar split layout is an error. A toolbar may contain at
// most one tracking separator. AppKit owns this standard item, so there is no
// mutable Wails item handle to return.
func (t *MacToolbar) AddSidebarTrackingSeparator() {
	t.add(newMacToolbarItem(t, toolbarSidebarTrackingSeparator, ""))
}

// AddInspectorToggle adds AppKit's standard inspector toggle on macOS 14 and
// newer. On older supported releases Wails supplies a native toolbar button
// with the same behavior. It requires a MacSplitView containing an inspector
// pane and needs no application callback or identifier.
func (t *MacToolbar) AddInspectorToggle() {
	item := newMacToolbarItem(t, toolbarInspectorToggle, "Inspector")
	item.symbolName = "sidebar.trailing"
	item.tooltip = "Show or hide the inspector"
	item.bordered = true
	item.onClick = func(*Context) {
		if pane := t.attachedInspectorPane(); pane != nil {
			pane.Toggle()
		}
	}
	t.add(item)
}

// AddInspectorTrackingSeparator adds AppKit's inspector tracking separator on
// macOS 14 and newer. It follows the leading edge of the native inspector as
// its divider moves. Older supported releases omit the separator while the
// inspector and its toggle continue to work.
func (t *MacToolbar) AddInspectorTrackingSeparator() {
	t.add(newMacToolbarItem(t, toolbarInspectorTrackingSeparator, ""))
}

func (t *MacToolbar) attachedInspectorPane() *MacSplitPane {
	if t == nil {
		return nil
	}
	t.stateLock.RLock()
	var window macToolbarWindow
	if t.state != nil {
		window = t.state.window
	}
	t.stateLock.RUnlock()
	if window == nil {
		return nil
	}
	return window.macInspectorPane()
}

// hasSidebarTrackingSeparator reports whether the toolbar contains a sidebar
// tracking separator, which is only meaningful on a window with a native
// sidebar split layout.
func (t *MacToolbar) hasSidebarTrackingSeparator() bool {
	for _, item := range t.itemSnapshot() {
		if item == nil {
			continue
		}
		item.lock.RLock()
		kind := item.kind
		item.lock.RUnlock()
		if kind == toolbarSidebarTrackingSeparator {
			return true
		}
	}
	return false
}

func (t *MacToolbar) hasInspectorChrome() bool {
	for _, item := range t.itemSnapshot() {
		if item == nil {
			continue
		}
		item.lock.RLock()
		kind := item.kind
		item.lock.RUnlock()
		if kind == toolbarInspectorToggle || kind == toolbarInspectorTrackingSeparator {
			return true
		}
	}
	return false
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
	// A native segmented group is rebuilt when its members change.
	g.toolbar.syncNative(nil)
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

// SetShowsIndicator controls whether a menu item shows the dropdown chevron
// (NSMenuToolbarItem.showsIndicator). It is ignored for other item kinds.
func (i *MacToolbarItem) SetShowsIndicator(shows bool) *MacToolbarItem {
	if i == nil {
		return i
	}
	i.lock.Lock()
	i.showsIndicator = shows
	i.lock.Unlock()
	i.update(func(native unsafe.Pointer) { macToolbarMenuItemSetShowsIndicator(native, i.identifier, shows) })
	return i
}

// SetVisibilityPriority controls which items AppKit moves into the overflow
// menu first when the window narrows (NSToolbarItem.visibilityPriority).
func (i *MacToolbarItem) SetVisibilityPriority(priority MacToolbarVisibilityPriority) *MacToolbarItem {
	if i == nil {
		return i
	}
	i.lock.Lock()
	i.visibilityPriority = priority
	i.lock.Unlock()
	i.update(func(native unsafe.Pointer) { macToolbarItemSetVisibilityPriority(native, i.identifier, priority) })
	return i
}

// SetNavigational marks an item as navigational (NSToolbarItem.navigational,
// macOS 11 and newer): AppKit keeps it at the leading edge of the toolbar,
// before the sidebar section, as back and forward buttons are in Safari.
func (i *MacToolbarItem) SetNavigational(navigational bool) *MacToolbarItem {
	if i == nil {
		return i
	}
	i.lock.Lock()
	i.navigational = navigational
	i.lock.Unlock()
	i.update(func(native unsafe.Pointer) { macToolbarItemSetNavigational(native, i.identifier, navigational) })
	return i
}

// SetPersistenceKey gives the item a stable NSToolbarItemIdentifier so a
// user-customised layout (see MacToolbar.SetCustomizable) can restore it
// across launches. Keys must be unique within a toolbar; an empty key
// restores the generated identifier. Changing the key of an attached item
// rebuilds its native item.
func (i *MacToolbarItem) SetPersistenceKey(key string) *MacToolbarItem {
	if i == nil {
		return i
	}
	i.lock.Lock()
	if key == "" {
		key = i.generatedIdentifier
	}
	changed := i.identifier != key
	i.identifier = key
	i.lock.Unlock()
	if changed && i.toolbar != nil {
		i.toolbar.syncNative(nil)
	}
	return i
}

// SetInDefaultSet controls whether a top-level item is part of the default
// toolbar layout. Items outside the default set are still offered in the
// customisation palette (see MacToolbar.SetCustomizable) and are otherwise
// hidden. The default is true.
func (i *MacToolbarItem) SetInDefaultSet(inDefaultSet bool) *MacToolbarItem {
	if i == nil {
		return i
	}
	i.lock.Lock()
	changed := i.inDefaultSet != inDefaultSet
	i.inDefaultSet = inDefaultSet
	i.lock.Unlock()
	if changed && i.toolbar != nil {
		i.toolbar.syncNative(nil)
	}
	return i
}

// SetSearchRecentsKey enables the native recent-searches menu of a search
// item and persists the recent queries under key
// (NSSearchField.recentsAutosaveName). An empty key disables recents. It is
// ignored for other item kinds.
func (i *MacToolbarItem) SetSearchRecentsKey(key string) *MacToolbarItem {
	if i == nil {
		return i
	}
	i.lock.Lock()
	i.searchRecentsKey = key
	menu := i.searchMenu
	i.lock.Unlock()
	i.update(func(native unsafe.Pointer) {
		macToolbarSearchItemSetRecents(native, i.identifier, key, macToolbarDefaultSearchRecents)
		macToolbarSearchItemSetMenu(native, i.identifier, menu, key != "")
	})
	return i
}

// SetSearchMenu installs a custom menu behind the search field's magnifier
// (NSSearchField.searchMenuTemplate). MenuItem.OnClick callbacks fire as
// usual. When a recents key is also set, the standard recent-searches section
// is appended after the custom items. Passing nil removes the custom menu.
// Calling Menu.Update after attachment rebuilds the custom items; call
// SetSearchMenu again to restore the recents section.
func (i *MacToolbarItem) SetSearchMenu(menu *Menu) *MacToolbarItem {
	if i == nil {
		return i
	}
	i.lock.Lock()
	i.searchMenu = menu
	recentsKey := i.searchRecentsKey
	i.lock.Unlock()
	i.update(func(native unsafe.Pointer) {
		macToolbarSearchItemSetMenu(native, i.identifier, menu, recentsKey != "")
	})
	return i
}

// SetSearchIncremental controls when OnSearch fires: incrementally as the
// user types (NSSearchField.sendsSearchStringImmediately) or only when the
// user presses Return or clears the field (the default,
// sendsWholeSearchString).
func (i *MacToolbarItem) SetSearchIncremental(incremental bool) *MacToolbarItem {
	if i == nil {
		return i
	}
	i.lock.Lock()
	i.searchIncremental = incremental
	i.lock.Unlock()
	i.update(func(native unsafe.Pointer) { macToolbarSearchItemSetIncremental(native, i.identifier, incremental) })
	return i
}

// SetSearchPlaceholder sets the placeholder text shown in an empty search
// field. It is ignored for other item kinds.
func (i *MacToolbarItem) SetSearchPlaceholder(placeholder string) *MacToolbarItem {
	if i == nil {
		return i
	}
	i.lock.Lock()
	i.searchPlaceholder = placeholder
	i.lock.Unlock()
	i.update(func(native unsafe.Pointer) { macToolbarSearchItemSetPlaceholder(native, i.identifier, placeholder) })
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

// ShowPopover shows a MacPopover anchored to this item's toolbar control.
// The item must be installed in a live toolbar; on macOS 14 and newer the
// popover anchors to the item itself, on earlier releases it anchors to the
// item's custom view (search fields) and other items return an error. See
// webview_window_popover.go.
func (i *MacToolbarItem) ShowPopover(popover *MacPopover) error {
	return macToolbarItemShowPopover(i, popover)
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
	nativeEnabled := enabled && (i.kind != toolbarShare || len(i.shareFormats) > 0)
	i.lock.Unlock()
	i.update(func(native unsafe.Pointer) { macToolbarItemSetEnabled(native, i.identifier, nativeEnabled) })
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

// SetProvider installs the lazy data provider used by the native share item.
// It calls ShareRepresentations immediately and copies the returned formats.
// Passing nil, or a provider with no valid representations, disables the item.
// Replacing a provider is safe while sharing is in progress: an in-flight
// native request retains the provider snapshot with which it started.
func (s *MacToolbarShareItem) SetProvider(provider MacShareProvider) *MacToolbarShareItem {
	if s == nil || s.MacToolbarItem == nil {
		return s
	}
	var formats []MacShareRepresentation
	if provider != nil {
		formats = normaliseMacShareRepresentations(provider.ShareRepresentations())
	}
	s.lock.Lock()
	s.shareProvider = provider
	s.shareFormats = formats
	subject := s.shareSubject
	suggestedName := s.shareSuggestedName
	enabled := !s.disabled && len(formats) > 0
	s.lock.Unlock()
	s.update(func(native unsafe.Pointer) {
		macToolbarShareItemSetProvider(native, s.identifier, provider, subject, suggestedName, formats)
		macToolbarItemSetEnabled(native, s.identifier, enabled)
	})
	return s
}

// SetSubject sets the subject used by sharing services that support one, such
// as Mail. It does not alter the bytes returned by the provider.
func (s *MacToolbarShareItem) SetSubject(subject string) *MacToolbarShareItem {
	if s == nil || s.MacToolbarItem == nil {
		return s
	}
	s.lock.Lock()
	s.shareSubject = subject
	provider := s.shareProvider
	formats := append([]MacShareRepresentation(nil), s.shareFormats...)
	suggestedName := s.shareSuggestedName
	s.lock.Unlock()
	s.update(func(native unsafe.Pointer) {
		macToolbarShareItemSetProvider(native, s.identifier, provider, subject, suggestedName, formats)
	})
	return s
}

// SetSuggestedName sets the human-readable name receiving services use for
// the shared logical item. The name applies to every advertised representation;
// AppKit derives the appropriate file type from the requested content type.
// This maps to NSItemProvider.suggestedName on macOS 10.14 and newer.
func (s *MacToolbarShareItem) SetSuggestedName(name string) *MacToolbarShareItem {
	if s == nil || s.MacToolbarItem == nil {
		return s
	}
	s.lock.Lock()
	s.shareSuggestedName = name
	provider := s.shareProvider
	subject := s.shareSubject
	formats := append([]MacShareRepresentation(nil), s.shareFormats...)
	s.lock.Unlock()
	s.update(func(native unsafe.Pointer) {
		macToolbarShareItemSetProvider(native, s.identifier, provider, subject, name, formats)
	})
	return s
}

// OnShared sets the callback invoked after a sharing service successfully
// shares the content. service is the localized AppKit service title. The
// callback is asynchronous and may safely use live toolbar-item setters.
func (s *MacToolbarShareItem) OnShared(callback func(*Context, string)) *MacToolbarShareItem {
	if s != nil && s.MacToolbarItem != nil {
		s.lock.Lock()
		s.onShared = callback
		s.lock.Unlock()
	}
	return s
}

// OnShareError sets the callback invoked when a selected sharing service or
// provider fails. service is the localized AppKit service title and may be
// empty when AppKit cannot associate the failure with a named service. The
// callback is asynchronous.
func (s *MacToolbarShareItem) OnShareError(callback func(*Context, string, error)) *MacToolbarShareItem {
	if s != nil && s.MacToolbarItem != nil {
		s.lock.Lock()
		s.onShareError = callback
		s.lock.Unlock()
	}
	return s
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
	window  macToolbarWindow
	native  unsafe.Pointer
	itemIDs []uint
	// nativeItems records which top-level items have been built into the
	// native toolbar, keyed by identifier, so a structural change can rebuild
	// only what changed.
	nativeItems map[string]macToolbarNativeItem
}

// macToolbarNativeItem describes one built native item. signature captures
// the properties that require a rebuild when changed (a group's member
// identifiers); itemIDs are the callback registrations owned by the item.
type macToolbarNativeItem struct {
	signature         string
	itemIDs           []uint
	memberIdentifiers []string
}

type macToolbarWindow interface {
	ID() uint
	Error(message string, args ...any)
	macInspectorPane() *MacSplitPane
}

func claimMacToolbar(toolbar *MacToolbar, window macToolbarWindow) (bool, error) {
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

func releaseMacToolbarOwnership(toolbar *MacToolbar, window macToolbarWindow) {
	if toolbar == nil {
		return
	}
	toolbar.stateLock.Lock()
	if toolbar.state != nil && toolbar.state.window == window && toolbar.state.native == nil {
		toolbar.state.window = nil
	}
	toolbar.stateLock.Unlock()
}

// validateToolbarItems checks a toolbar before it is attached: every item
// must be well formed and every button or search field must have its
// callback.
func validateToolbarItems(items []*MacToolbarItem) error {
	return validateToolbarLayout(items, true)
}

// validateToolbarLayout checks item structure and identifier uniqueness.
// requireCallbacks additionally enforces OnClick and OnSearch; a live
// structural change skips that because callbacks are chained after the
// constructor, and a missing callback is reported when the item is used.
func validateToolbarLayout(items []*MacToolbarItem, requireCallbacks bool) error {
	sidebarToggles := 0
	trackingSeparators := 0
	inspectorToggles := 0
	inspectorTrackingSeparators := 0
	identifiers := make(map[string]string)
	var validate func([]*MacToolbarItem, string) error
	validate = func(items []*MacToolbarItem, parent string) error {
		for _, item := range items {
			if item == nil {
				return fmt.Errorf("toolbar item in %q is nil", parent)
			}
			item.lock.RLock()
			kind := item.kind
			label := item.label
			identifier := item.identifier
			onClick := item.onClick
			onSearch := item.onSearch
			menu := item.menu
			members := append([]*MacToolbarItem(nil), item.items...)
			item.lock.RUnlock()

			if kind != toolbarFlexibleSpace && kind != toolbarSpace &&
				kind != toolbarSidebarToggle && kind != toolbarSidebarTrackingSeparator &&
				kind != toolbarInspectorTrackingSeparator {
				if previous, exists := identifiers[identifier]; exists {
					return fmt.Errorf("toolbar items %q and %q share the persistence key %q", previous, label, identifier)
				}
				identifiers[identifier] = label
			}

			switch kind {
			case toolbarButton:
				if requireCallbacks && onClick == nil {
					return fmt.Errorf("toolbar button %q requires OnClick", label)
				}
			case toolbarSearchField:
				if requireCallbacks && onSearch == nil {
					return fmt.Errorf("toolbar search field %q requires OnSearch", label)
				}
			case toolbarMenu:
				if menu == nil {
					return fmt.Errorf("toolbar menu %q requires a Menu", label)
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
					valid := member.kind == toolbarButton && (!requireCallbacks || member.onClick != nil)
					member.lock.RUnlock()
					if !valid {
						return fmt.Errorf("toolbar group %q members must be buttons with OnClick", label)
					}
				}
				if err := validate(members, label); err != nil {
					return err
				}
			case toolbarShare, toolbarFlexibleSpace, toolbarSpace:
			case toolbarSidebarToggle:
				// AppKit owns the toggle's action, so no OnClick is required.
				sidebarToggles++
				if sidebarToggles > 1 {
					return fmt.Errorf("toolbar may only contain one sidebar toggle")
				}
			case toolbarSidebarTrackingSeparator:
				trackingSeparators++
				if trackingSeparators > 1 {
					return fmt.Errorf("toolbar may only contain one sidebar tracking separator")
				}
			case toolbarInspectorToggle:
				// AppKit owns the action on macOS 14+. The internal callback is
				// used by Wails' native fallback on earlier releases.
				inspectorToggles++
				if inspectorToggles > 1 {
					return fmt.Errorf("toolbar may only contain one inspector toggle")
				}
			case toolbarInspectorTrackingSeparator:
				inspectorTrackingSeparators++
				if inspectorTrackingSeparators > 1 {
					return fmt.Errorf("toolbar may only contain one inspector tracking separator")
				}
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

type toolbarShareEvent struct {
	itemID  uint
	service string
	err     string
}

var toolbarShareCompleted = make(chan toolbarShareEvent, 32)

type macShareProviderRegistration struct {
	provider      MacShareProvider
	formats       []MacShareRepresentation
	suggestedName string
}

var toolbarShareProviderID uint32
var toolbarShareProviderMap = make(map[uint]macShareProviderRegistration)
var toolbarShareProviderMapLock sync.RWMutex

func registerToolbarShareProvider(provider MacShareProvider, formats []MacShareRepresentation, suggestedName string) uint {
	if provider == nil || len(formats) == 0 {
		return 0
	}
	id := uint(atomic.AddUint32(&toolbarShareProviderID, 1))
	if id == 0 {
		id = uint(atomic.AddUint32(&toolbarShareProviderID, 1))
	}
	registration := macShareProviderRegistration{
		provider:      provider,
		formats:       append([]MacShareRepresentation(nil), formats...),
		suggestedName: suggestedName,
	}
	toolbarShareProviderMapLock.Lock()
	toolbarShareProviderMap[id] = registration
	toolbarShareProviderMapLock.Unlock()
	return id
}

func releaseToolbarShareProvider(id uint) {
	if id == 0 {
		return
	}
	toolbarShareProviderMapLock.Lock()
	delete(toolbarShareProviderMap, id)
	toolbarShareProviderMapLock.Unlock()
}

func handleToolbarItemClicked(itemID uint) {
	defer handlePanic()
	item := getToolbarItemByID(itemID)
	if item == nil {
		return
	}
	item.lock.RLock()
	callback := item.onClick
	label := item.label
	item.lock.RUnlock()
	if callback == nil {
		// Items added after attachment are installed before their callback
		// is chained, so the missing callback is reported here instead.
		item.toolbar.logError("toolbar button %q was clicked but has no OnClick", label)
		return
	}
	callback(newContext())
}

func handleToolbarSearch(itemID uint, query string) {
	defer handlePanic()
	item := getToolbarItemByID(itemID)
	if item == nil {
		return
	}
	item.lock.RLock()
	callback := item.onSearch
	label := item.label
	item.lock.RUnlock()
	if callback == nil {
		item.toolbar.logError("toolbar search field %q was used but has no OnSearch", label)
		return
	}
	callback(newContext(), query)
}

func handleToolbarShareResult(event toolbarShareEvent) {
	defer handlePanic()
	item := getToolbarItemByID(event.itemID)
	if item == nil {
		return
	}
	item.lock.RLock()
	onShared := item.onShared
	onError := item.onShareError
	item.lock.RUnlock()
	if event.err == "" {
		if onShared != nil {
			onShared(newContext(), event.service)
		}
		return
	}
	if onError != nil {
		onError(newContext(), event.service, fmt.Errorf("%s", event.err))
	}
}

func handleToolbarShareData(providerID uint, contentType MacShareContentType) (data []byte, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			data = nil
			err = fmt.Errorf("macOS share provider panicked: %v", recovered)
		}
	}()
	toolbarShareProviderMapLock.RLock()
	registration, exists := toolbarShareProviderMap[providerID]
	toolbarShareProviderMapLock.RUnlock()
	if !exists || registration.provider == nil {
		return nil, fmt.Errorf("macOS share provider is no longer available")
	}
	request := MacShareRequest{ContentType: contentType, SuggestedName: registration.suggestedName}
	for _, representation := range registration.formats {
		if representation.ContentType == contentType {
			return registration.provider.ShareData(request)
		}
	}
	return nil, fmt.Errorf("macOS share provider does not offer %q", contentType)
}
