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
	title := toolbar.AddTitle("Document")
	toolbar.AddSeparator()
	share := toolbar.AddShare("Share").SetProvider(MacShareProviderFunc{
		Available: []MacShareRepresentation{{ContentType: MacShareTypePlainText}},
		Load:      func(MacShareRequest) ([]byte, error) { return []byte("A note"), nil },
	})
	group := toolbar.AddGroup("View", ToolbarGroupSelectOne)
	group.AddButton("Write").OnClick(func(*Context) {})

	if err := validateToolbarItems(toolbar.itemSnapshot()); err != nil {
		t.Fatalf("valid toolbar rejected: %v", err)
	}
	if toolbar.identifier == "" || button.identifier == "" || search.identifier == "" || title.identifier == "" || share.identifier == "" || group.identifier == "" {
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
		SetNavigational(true).
		SetProminent(true).
		SetTintColor(color).
		SetEnabled(false).
		SetHidden(true).
		SetBadgeCount(3).
		SetDraggable(true)

	// Mutating the caller's color after SetTintColor must not mutate the item.
	color.Red = 200
	snapshot := snapshotToolbarItemForTest(item)
	if snapshot.label != "Details (3)" || snapshot.symbolName != "info.circle" || snapshot.tooltip != "Show details" {
		t.Fatal("text and symbol mutators should update the item before installation")
	}
	if !snapshot.bordered || !snapshot.navigational || !snapshot.prominent || !snapshot.disabled || !snapshot.hidden ||
		snapshot.badgeCount != 3 || !snapshot.draggable {
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

func TestToolbarMomentaryGroupDoesNotRetainSelection(t *testing.T) {
	toolbar := NewMacToolbar()
	group := toolbar.AddGroup("Navigation", ToolbarGroupMomentary)
	group.AddButton("Back").OnClick(func(*Context) {})
	group.AddButton("Forward").OnClick(func(*Context) {})

	if got := snapshotToolbarItemForTest(group.MacToolbarItem).selectedIndex; got != -1 {
		t.Fatalf("momentary group selected index = %d, want -1", got)
	}

	group.SetSelectedIndex(0)
	group.SetSelectionMode(ToolbarGroupMomentary)
	if got := snapshotToolbarItemForTest(group.MacToolbarItem).selectedIndex; got != -1 {
		t.Fatalf("momentary selection mode retained index %d, want -1", got)
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

type toolbarItemTestSnapshot struct {
	label         string
	symbolName    string
	tooltip       string
	bordered      bool
	navigational  bool
	prominent     bool
	tintColor     *RGBA
	badgeCount    int
	draggable     bool
	disabled      bool
	hidden        bool
	selectionMode MacToolbarGroupSelectionMode
	selectedIndex int
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
		navigational:  item.navigational,
		prominent:     item.prominent,
		badgeCount:    item.badgeCount,
		draggable:     item.draggable,
		disabled:      item.disabled,
		hidden:        item.hidden,
		selectionMode: item.selectionMode,
		selectedIndex: item.selectedIndex,
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
