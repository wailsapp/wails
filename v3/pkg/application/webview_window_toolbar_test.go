package application

import (
	"errors"
	"strings"
	"sync"
	"testing"
)

func TestValidateToolbarItems(t *testing.T) {
	toolbar := NewMacToolbar()
	button := toolbar.AddButton("Save").SetSymbol("checkmark").OnClick(func(*Context) {})
	search := toolbar.AddSearch("Search").OnSearch(func(*Context, string) {})
	share := toolbar.AddShare("Share").SetProvider(MacShareProviderFunc{
		Available: []MacShareRepresentation{{ContentType: MacShareTypePlainText}},
		Load:      func(MacShareRequest) ([]byte, error) { return []byte("A note"), nil },
	})
	group := toolbar.AddGroup("View", ToolbarGroupSelectOne)
	group.AddButton("Write").OnClick(func(*Context) {})

	if err := validateToolbarItems(toolbar.itemSnapshot()); err != nil {
		t.Fatalf("valid toolbar rejected: %v", err)
	}
	if toolbar.identifier == "" || button.identifier == "" || search.identifier == "" || share.identifier == "" || group.identifier == "" {
		t.Fatal("toolbar and items should receive internal identifiers")
	}
	if button.identifier == search.identifier || button.identifier == group.identifier {
		t.Fatal("toolbar item identifiers should be unique")
	}
	if toolbar.identifier == NewMacToolbar().identifier {
		t.Fatal("toolbar identifiers should be unique per toolbar")
	}
}

func TestToolbarDisplayModeDefaultsAndValidation(t *testing.T) {
	toolbar := NewMacToolbar()
	if toolbar.displayMode != MacToolbarDisplayModeIconAndLabel {
		t.Fatalf("default display mode = %d, want icon and label", toolbar.displayMode)
	}
	toolbar.SetDisplayMode(MacToolbarDisplayModeIconOnly)
	if toolbar.displayMode != MacToolbarDisplayModeIconOnly {
		t.Fatal("SetDisplayMode should update the pending toolbar")
	}
	toolbar.SetDisplayMode(MacToolbarDisplayMode(99))
	if toolbar.displayMode != MacToolbarDisplayModeIconOnly {
		t.Fatal("SetDisplayMode should ignore invalid values")
	}
}

func TestToolbarShareProviderIsNormalisedAndInvokedLazily(t *testing.T) {
	toolbar := NewMacToolbar()
	formats := []MacShareRepresentation{
		{ContentType: MacShareTypeHTML},
		{ContentType: MacShareTypePlainText},
		{ContentType: MacShareTypeHTML},
		{},
	}
	var requested MacShareRequest
	share := toolbar.AddShare("Share").SetProvider(MacShareProviderFunc{
		Available: formats,
		Load: func(request MacShareRequest) ([]byte, error) {
			requested = request
			return []byte("<strong>A note</strong>"), nil
		},
	}).SetSuggestedName("Daymark Note")
	formats[0].ContentType = MacShareTypePDF

	snapshot := snapshotToolbarItemForTest(share.MacToolbarItem)
	if len(snapshot.shareFormats) != 2 {
		t.Fatalf("share formats = %#v, want two unique non-empty formats", snapshot.shareFormats)
	}
	if snapshot.shareFormats[0].ContentType != MacShareTypeHTML {
		t.Fatal("share representation slices must be copied")
	}
	if err := validateToolbarItems(toolbar.itemSnapshot()); err != nil {
		t.Fatalf("share item should not require a click callback: %v", err)
	}

	providerID := registerToolbarShareProvider(share.shareProvider, share.shareFormats, share.shareSuggestedName)
	t.Cleanup(func() { releaseToolbarShareProvider(providerID) })
	data, err := handleToolbarShareData(providerID, MacShareTypeHTML)
	if err != nil {
		t.Fatalf("load HTML representation: %v", err)
	}
	if string(data) != "<strong>A note</strong>" {
		t.Fatalf("share data = %q", data)
	}
	if requested.ContentType != MacShareTypeHTML || requested.SuggestedName != "Daymark Note" {
		t.Fatalf("request = %#v", requested)
	}
}

func TestToolbarShareProviderErrors(t *testing.T) {
	toolbar := NewMacToolbar()
	share := toolbar.AddShare("Share").SetProvider(MacShareProviderFunc{
		Available: []MacShareRepresentation{{ContentType: MacShareTypePDF}},
		Load:      func(MacShareRequest) ([]byte, error) { panic("renderer failed") },
	})
	providerID := registerToolbarShareProvider(share.shareProvider, share.shareFormats, share.shareSuggestedName)
	t.Cleanup(func() { releaseToolbarShareProvider(providerID) })

	if _, err := handleToolbarShareData(providerID, MacShareTypeHTML); err == nil {
		t.Fatal("an unadvertised representation should fail")
	}
	if _, err := handleToolbarShareData(providerID, MacShareTypePDF); err == nil || !strings.Contains(err.Error(), "panicked") {
		t.Fatalf("provider panic was not converted to an error: %v", err)
	}

	share.SetProvider(MacShareProviderFunc{
		Available: []MacShareRepresentation{{ContentType: MacShareTypePDF}},
		Load:      func(MacShareRequest) ([]byte, error) { return nil, errors.New("PDF unavailable") },
	})
	if _, err := handleToolbarShareData(providerID, MacShareTypePDF); err == nil || !strings.Contains(err.Error(), "panicked") {
		t.Fatalf("existing registration did not retain its provider snapshot: %v", err)
	}
	newProviderID := registerToolbarShareProvider(share.shareProvider, share.shareFormats, share.shareSuggestedName)
	t.Cleanup(func() { releaseToolbarShareProvider(newProviderID) })
	if _, err := handleToolbarShareData(newProviderID, MacShareTypePDF); err == nil || err.Error() != "PDF unavailable" {
		t.Fatalf("provider error = %v", err)
	}
}

func TestToolbarShareCallbacks(t *testing.T) {
	toolbar := NewMacToolbar()
	share := toolbar.AddShare("Share")
	var sharedService string
	var failedService string
	var failure string
	share.OnShared(func(_ *Context, service string) { sharedService = service })
	share.OnShareError(func(_ *Context, service string, err error) {
		failedService = service
		failure = err.Error()
	})
	id := nextToolbarNativeID()
	addToToolbarItemMap(id, share.MacToolbarItem)
	t.Cleanup(func() { removeFromToolbarItemMap(id) })

	handleToolbarShareResult(toolbarShareEvent{itemID: id, service: "Mail"})
	if sharedService != "Mail" {
		t.Fatalf("shared service = %q, want Mail", sharedService)
	}
	handleToolbarShareResult(toolbarShareEvent{itemID: id, service: "AirDrop", err: "Unavailable"})
	if failedService != "AirDrop" || failure != "Unavailable" {
		t.Fatalf("failure callback = %q, %q", failedService, failure)
	}
}

func TestToolbarSidebarItemsRequireNoCallbacks(t *testing.T) {
	toolbar := NewMacToolbar()
	toolbar.AddSidebarToggle()
	toolbar.AddButton("New").OnClick(func(*Context) {})
	toolbar.AddSidebarTrackingSeparator()

	if err := validateToolbarItems(toolbar.itemSnapshot()); err != nil {
		t.Fatalf("sidebar items should not require callbacks: %v", err)
	}
	if !toolbar.hasSidebarTrackingSeparator() {
		t.Fatal("hasSidebarTrackingSeparator should report the separator")
	}
	if NewMacToolbar().hasSidebarTrackingSeparator() {
		t.Fatal("an empty toolbar must not report a tracking separator")
	}
}

func TestToolbarRejectsDuplicateSidebarItems(t *testing.T) {
	toggleToolbar := NewMacToolbar()
	toggleToolbar.AddSidebarToggle()
	toggleToolbar.AddSidebarToggle()
	if err := validateToolbarItems(toggleToolbar.itemSnapshot()); err == nil || !strings.Contains(err.Error(), "one sidebar toggle") {
		t.Fatalf("duplicate sidebar toggles should be rejected, got %v", err)
	}

	separatorToolbar := NewMacToolbar()
	separatorToolbar.AddSidebarTrackingSeparator()
	separatorToolbar.AddSidebarTrackingSeparator()
	if err := validateToolbarItems(separatorToolbar.itemSnapshot()); err == nil || !strings.Contains(err.Error(), "one sidebar tracking separator") {
		t.Fatalf("duplicate tracking separators should be rejected, got %v", err)
	}
}

func TestToolbarInspectorItemsRequireNoApplicationCallbacks(t *testing.T) {
	toolbar := NewMacToolbar()
	toolbar.AddInspectorTrackingSeparator()
	toolbar.AddInspectorToggle()
	if err := validateToolbarItems(toolbar.itemSnapshot()); err != nil {
		t.Fatalf("inspector items should own their actions: %v", err)
	}
	if !toolbar.hasInspectorChrome() {
		t.Fatal("toolbar should report native inspector chrome")
	}

	toggleToolbar := NewMacToolbar()
	toggleToolbar.AddInspectorToggle()
	toggleToolbar.AddInspectorToggle()
	if err := validateToolbarItems(toggleToolbar.itemSnapshot()); err == nil ||
		!strings.Contains(err.Error(), "one inspector toggle") {
		t.Fatalf("duplicate inspector toggles should be rejected, got %v", err)
	}

	separatorToolbar := NewMacToolbar()
	separatorToolbar.AddInspectorTrackingSeparator()
	separatorToolbar.AddInspectorTrackingSeparator()
	if err := validateToolbarItems(separatorToolbar.itemSnapshot()); err == nil ||
		!strings.Contains(err.Error(), "one inspector tracking separator") {
		t.Fatalf("duplicate inspector tracking separators should be rejected, got %v", err)
	}
}

func TestValidateToolbarItemsRejectsMissingCallbacks(t *testing.T) {
	toolbar := NewMacToolbar()
	toolbar.AddButton("Save")
	if err := validateToolbarItems(toolbar.itemSnapshot()); err == nil {
		t.Fatal("button without callback should be rejected")
	}

	searchToolbar := NewMacToolbar()
	searchToolbar.AddSearch("Search")
	if err := validateToolbarItems(searchToolbar.itemSnapshot()); err == nil {
		t.Fatal("search field without callback should be rejected")
	}
}

func TestToolbarGroupRejectsInvalidMembers(t *testing.T) {
	toolbar := NewMacToolbar()
	group := toolbar.AddGroup("View", ToolbarGroupSelectOne)
	group.lock.Lock()
	group.items = append(group.items, newMacToolbarItem(toolbar, toolbarSearchField, "Search"))
	group.lock.Unlock()
	if err := validateToolbarItems(toolbar.itemSnapshot()); err == nil {
		t.Fatal("group should reject non-button members")
	}
}

func TestToolbarMutatorsBeforeInstallation(t *testing.T) {
	toolbar := NewMacToolbar()
	item := toolbar.AddButton("Details").OnClick(func(*Context) {})
	color := &RGBA{Red: 12, Green: 34, Blue: 56, Alpha: 255}
	item.SetLabel("Details (3)").
		SetSymbol("info.circle").
		SetTooltip("Show details").
		SetBordered(true).
		SetProminent(true).
		SetTintColor(color).
		SetEnabled(false).
		SetHidden(true).
		SetBadgeCount(3)

	// Mutating the caller's color after SetTintColor must not mutate the item.
	color.Red = 200
	snapshot := snapshotToolbarItemForTest(item)
	if snapshot.label != "Details (3)" || snapshot.symbolName != "info.circle" || snapshot.tooltip != "Show details" {
		t.Fatal("text and symbol mutators should update the item before installation")
	}
	if !snapshot.bordered || !snapshot.prominent || !snapshot.disabled || !snapshot.hidden || snapshot.badgeCount != 3 {
		t.Fatal("state mutators should update the item before installation")
	}
	if snapshot.tintColor == nil || snapshot.tintColor.Red != 12 {
		t.Fatal("tint color should be copied when assigned")
	}
}

func TestToolbarGroupSelectionValidation(t *testing.T) {
	toolbar := NewMacToolbar()
	group := toolbar.AddGroup("Mode", ToolbarGroupSelectOne)
	group.AddButton("Write").OnClick(func(*Context) {})
	group.AddButton("Preview").OnClick(func(*Context) {})

	group.SetSelectedIndex(1)
	if got := snapshotToolbarItemForTest(group.MacToolbarItem).selectedIndex; got != 1 {
		t.Fatalf("selected index = %d, want 1", got)
	}
	group.SetSelectedIndex(9)
	if got := snapshotToolbarItemForTest(group.MacToolbarItem).selectedIndex; got != 1 {
		t.Fatalf("invalid index changed selection to %d", got)
	}
	group.SetSelectionMode(ToolbarGroupMomentary)
	if got := snapshotToolbarItemForTest(group.MacToolbarItem).selectionMode; got != ToolbarGroupMomentary {
		t.Fatalf("selection mode = %d, want momentary", got)
	}
}

func TestToolbarOwnershipReleasedOnDetach(t *testing.T) {
	firstWindow := &WebviewWindow{}
	secondWindow := &WebviewWindow{}
	toolbar := NewMacToolbar()
	toolbar.AddButton("Save").OnClick(func(*Context) {})

	firstWindow.SetToolbar(toolbar)
	toolbar.stateLock.RLock()
	owner := toolbar.state.window
	toolbar.stateLock.RUnlock()
	if owner != firstWindow {
		t.Fatal("first window should own the stashed toolbar")
	}

	firstWindow.SetToolbar(nil)
	toolbar.stateLock.RLock()
	owner = toolbar.state.window
	toolbar.stateLock.RUnlock()
	if owner != nil {
		t.Fatal("detaching should release toolbar ownership")
	}

	secondWindow.SetToolbar(toolbar)
	toolbar.stateLock.RLock()
	owner = toolbar.state.window
	toolbar.stateLock.RUnlock()
	if owner != secondWindow {
		t.Fatal("a detached toolbar should be reusable on another window")
	}
}

func TestToolbarConcurrentConfiguration(t *testing.T) {
	toolbar := NewMacToolbar()
	item := toolbar.AddButton("Save").OnClick(func(*Context) {})

	var wait sync.WaitGroup
	for index := 0; index < 20; index++ {
		wait.Add(1)
		go func(value int) {
			defer wait.Done()
			item.SetBadgeCount(value)
			item.SetEnabled(value%2 == 0)
			item.SetHidden(value%3 == 0)
			_ = snapshotToolbarItemForTest(item)
		}(index)
	}
	wait.Wait()
}

func TestToolbarMenuAndSpaceItems(t *testing.T) {
	toolbar := NewMacToolbar()
	space := toolbar.AddSpace()
	menu := NewMenu()
	menu.Add("Duplicate").OnClick(func(*Context) {})
	actions := toolbar.AddMenu("Actions", menu).SetSymbol("ellipsis.circle")

	if space.kind != toolbarSpace || actions.kind != toolbarMenu {
		t.Fatal("AddSpace and AddMenu should create their dedicated item kinds")
	}
	if actions.menu != menu || !actions.showsIndicator {
		t.Fatal("a menu item should keep its menu and show the indicator by default")
	}
	if err := validateToolbarItems(toolbar.itemSnapshot()); err != nil {
		t.Fatalf("menu and space items need no callbacks: %v", err)
	}
	actions.SetShowsIndicator(false)
	if actions.showsIndicator {
		t.Fatal("SetShowsIndicator should update the pending item")
	}

	missing := NewMacToolbar()
	missing.AddMenu("Actions", nil)
	if err := validateToolbarItems(missing.itemSnapshot()); err == nil || !strings.Contains(err.Error(), "requires a Menu") {
		t.Fatalf("a menu item without a menu should be rejected, got %v", err)
	}
}

func toolbarLabels(items []*MacToolbarItem) []string {
	labels := make([]string, 0, len(items))
	for _, item := range items {
		labels = append(labels, item.label)
	}
	return labels
}

func TestToolbarRemoveAndMoveOrdering(t *testing.T) {
	toolbar := NewMacToolbar()
	first := toolbar.AddButton("First").OnClick(func(*Context) {})
	second := toolbar.AddButton("Second").OnClick(func(*Context) {})
	third := toolbar.AddButton("Third").OnClick(func(*Context) {})
	toolbar.SetCenteredItems(second)

	toolbar.Move(third, 0)
	if got := strings.Join(toolbarLabels(toolbar.itemSnapshot()), ","); got != "Third,First,Second" {
		t.Fatalf("after Move(third, 0) order = %s", got)
	}
	toolbar.Move(third, 99)
	if got := strings.Join(toolbarLabels(toolbar.itemSnapshot()), ","); got != "First,Second,Third" {
		t.Fatalf("Move should clamp a large index to the end, order = %s", got)
	}
	toolbar.Move(first, -5)
	if got := strings.Join(toolbarLabels(toolbar.itemSnapshot()), ","); got != "First,Second,Third" {
		t.Fatalf("Move should clamp a negative index to the start, order = %s", got)
	}

	toolbar.Remove(second)
	if got := strings.Join(toolbarLabels(toolbar.itemSnapshot()), ","); got != "First,Third" {
		t.Fatalf("after Remove(second) order = %s", got)
	}
	toolbar.itemsLock.RLock()
	centered := len(toolbar.centeredItems)
	toolbar.itemsLock.RUnlock()
	if centered != 0 {
		t.Fatal("removing an item should also drop it from the centred set")
	}

	foreign := NewMacToolbar().AddButton("Foreign").OnClick(func(*Context) {})
	toolbar.Remove(foreign)
	toolbar.Move(foreign, 0)
	if got := strings.Join(toolbarLabels(toolbar.itemSnapshot()), ","); got != "First,Third" {
		t.Fatalf("items from another toolbar must be ignored, order = %s", got)
	}
}

func TestToolbarGroupMembersRemoveAndMove(t *testing.T) {
	toolbar := NewMacToolbar()
	group := toolbar.AddGroup("Mode", ToolbarGroupSelectOne)
	write := group.AddButton("Write").OnClick(func(*Context) {})
	preview := group.AddButton("Preview").OnClick(func(*Context) {})
	split := group.AddButton("Split").OnClick(func(*Context) {})

	toolbar.Move(split, 0)
	if got := strings.Join(toolbarLabels(snapshotToolbarItemForTest(group.MacToolbarItem).items), ","); got != "Split,Write,Preview" {
		t.Fatalf("group order after Move = %s", got)
	}
	toolbar.Remove(write)
	if got := strings.Join(toolbarLabels(snapshotToolbarItemForTest(group.MacToolbarItem).items), ","); got != "Split,Preview" {
		t.Fatalf("group order after Remove = %s", got)
	}
	if len(toolbar.itemSnapshot()) != 1 {
		t.Fatal("group member changes must not alter the top-level item list")
	}
	_ = preview
}

func TestToolbarPersistenceKeys(t *testing.T) {
	toolbar := NewMacToolbar()
	save := toolbar.AddButton("Save").OnClick(func(*Context) {})
	share := toolbar.AddShare("Share")
	generated := save.identifier
	if !strings.HasPrefix(generated, "wails.toolbar.item.") {
		t.Fatalf("generated identifier = %q", generated)
	}

	save.SetPersistenceKey("save")
	if save.identifier != "save" {
		t.Fatalf("identifier after SetPersistenceKey = %q", save.identifier)
	}
	if err := validateToolbarItems(toolbar.itemSnapshot()); err != nil {
		t.Fatalf("unique persistence keys should validate: %v", err)
	}

	share.SetPersistenceKey("save")
	if err := validateToolbarItems(toolbar.itemSnapshot()); err == nil || !strings.Contains(err.Error(), "share the persistence key") {
		t.Fatalf("duplicate persistence keys should be rejected, got %v", err)
	}

	save.SetPersistenceKey("")
	if save.identifier != generated {
		t.Fatal("an empty persistence key should restore the generated identifier")
	}

	// Standard AppKit identifiers legitimately repeat.
	spaces := NewMacToolbar()
	spaces.AddSpace()
	spaces.AddSpace()
	spaces.AddFlexibleSpace()
	spaces.AddFlexibleSpace()
	if err := validateToolbarItems(spaces.itemSnapshot()); err != nil {
		t.Fatalf("repeated spaces should validate: %v", err)
	}
}

func TestToolbarCustomizationModel(t *testing.T) {
	toolbar := NewMacToolbar()
	if toolbar.customizable {
		t.Fatal("toolbars must not be customizable by default")
	}
	toolbar.SetCustomizable("example.main")
	if !toolbar.customizable || toolbar.persistenceKey != "example.main" {
		t.Fatal("SetCustomizable should record the persistence key")
	}
	toolbar.SetCustomizable("")
	if toolbar.customizable {
		t.Fatal("an empty persistence key should disable customization")
	}

	first := toolbar.AddButton("First").OnClick(func(*Context) {})
	second := toolbar.AddButton("Second").OnClick(func(*Context) {})
	third := toolbar.AddButton("Third").OnClick(func(*Context) {})
	if !first.inDefaultSet || !second.inDefaultSet || !third.inDefaultSet {
		t.Fatal("items should be in the default set by default")
	}
	second.SetInDefaultSet(false)
	if second.inDefaultSet {
		t.Fatal("SetInDefaultSet(false) should remove the item from the default layout")
	}
	toolbar.SetDefaultItems(first, third)
	if !first.inDefaultSet || second.inDefaultSet || !third.inDefaultSet {
		t.Fatal("SetDefaultItems should mark exactly the listed items")
	}
	toolbar.SetDefaultItems()
	if !first.inDefaultSet || !second.inDefaultSet || !third.inDefaultSet {
		t.Fatal("SetDefaultItems with no items should restore every item")
	}
	// Nothing is attached, so the palette request is a no-op.
	toolbar.RunCustomizationPalette()
}

func TestToolbarItemSemantics(t *testing.T) {
	if MacToolbarVisibilityPriorityStandard != 0 || MacToolbarVisibilityPriorityLow != -1000 ||
		MacToolbarVisibilityPriorityHigh != 1000 || MacToolbarVisibilityPriorityUser != 2000 {
		t.Fatal("visibility priorities must match NSToolbarItemVisibilityPriority")
	}
	toolbar := NewMacToolbar()
	back := toolbar.AddButton("Back").OnClick(func(*Context) {}).
		SetNavigational(true).
		SetVisibilityPriority(MacToolbarVisibilityPriorityLow)
	if !back.navigational || back.visibilityPriority != MacToolbarVisibilityPriorityLow {
		t.Fatal("navigational and visibility priority should update the pending item")
	}

	centred := toolbar.AddButton("Centred").OnClick(func(*Context) {})
	foreign := NewMacToolbar().AddButton("Foreign").OnClick(func(*Context) {})
	toolbar.SetCenteredItems(centred, foreign, nil)
	toolbar.itemsLock.RLock()
	centeredItems := append([]*MacToolbarItem(nil), toolbar.centeredItems...)
	toolbar.itemsLock.RUnlock()
	if len(centeredItems) != 1 || centeredItems[0] != centred {
		t.Fatalf("centred items should keep only this toolbar's items, got %d", len(centeredItems))
	}
}

func TestToolbarSearchOptions(t *testing.T) {
	toolbar := NewMacToolbar()
	menu := NewMenu()
	menu.Add("Titles only").OnClick(func(*Context) {})
	search := toolbar.AddSearch("Search").OnSearch(func(*Context, string) {}).
		SetSearchRecentsKey("example.search").
		SetSearchMenu(menu).
		SetSearchIncremental(true).
		SetSearchPlaceholder("Search notes")

	search.lock.RLock()
	defer search.lock.RUnlock()
	if search.searchRecentsKey != "example.search" || search.searchMenu != menu ||
		!search.searchIncremental || search.searchPlaceholder != "Search notes" {
		t.Fatal("search options should update the pending item")
	}
}

func TestValidateToolbarLayoutRelaxesCallbacksForLiveAdditions(t *testing.T) {
	toolbar := NewMacToolbar()
	toolbar.AddButton("Later")
	toolbar.AddSearch("Later")
	group := toolbar.AddGroup("Mode", ToolbarGroupSelectOne)
	group.AddButton("Write")

	if err := validateToolbarLayout(toolbar.itemSnapshot(), false); err != nil {
		t.Fatalf("live additions may receive callbacks after construction: %v", err)
	}
	if err := validateToolbarLayout(toolbar.itemSnapshot(), true); err == nil {
		t.Fatal("attachment must still require callbacks")
	}
}

func TestToolbarClickWithoutCallbackIsIgnored(t *testing.T) {
	toolbar := NewMacToolbar()
	button := toolbar.AddButton("Later")
	search := toolbar.AddSearch("Later")
	buttonID := nextToolbarNativeID()
	searchID := nextToolbarNativeID()
	addToToolbarItemMap(buttonID, button)
	addToToolbarItemMap(searchID, search)
	t.Cleanup(func() {
		removeFromToolbarItemMap(buttonID)
		removeFromToolbarItemMap(searchID)
	})

	// Without an application the missing callback is reported nowhere; the
	// handlers must still return cleanly.
	handleToolbarItemClicked(buttonID)
	handleToolbarSearch(searchID, "query")

	clicked := false
	button.OnClick(func(*Context) { clicked = true })
	handleToolbarItemClicked(buttonID)
	if !clicked {
		t.Fatal("a callback chained after construction should fire on the next click")
	}
}

type toolbarItemTestSnapshot struct {
	label         string
	symbolName    string
	tooltip       string
	bordered      bool
	prominent     bool
	tintColor     *RGBA
	badgeCount    int
	disabled      bool
	hidden        bool
	selectionMode MacToolbarGroupSelectionMode
	selectedIndex int
	items         []*MacToolbarItem
	shareProvider MacShareProvider
	shareFormats  []MacShareRepresentation
	shareSubject  string
}

func snapshotToolbarItemForTest(item *MacToolbarItem) toolbarItemTestSnapshot {
	item.lock.RLock()
	defer item.lock.RUnlock()
	result := toolbarItemTestSnapshot{
		label:         item.label,
		symbolName:    item.symbolName,
		tooltip:       item.tooltip,
		bordered:      item.bordered,
		prominent:     item.prominent,
		badgeCount:    item.badgeCount,
		disabled:      item.disabled,
		hidden:        item.hidden,
		selectionMode: item.selectionMode,
		selectedIndex: item.selectedIndex,
		items:         append([]*MacToolbarItem(nil), item.items...),
		shareProvider: item.shareProvider,
		shareFormats:  append([]MacShareRepresentation(nil), item.shareFormats...),
		shareSubject:  item.shareSubject,
	}
	if item.tintColor != nil {
		copyOfColor := *item.tintColor
		result.tintColor = &copyOfColor
	}
	return result
}
