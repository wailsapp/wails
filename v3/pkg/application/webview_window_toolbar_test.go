package application

import (
	"sync"
	"testing"
)

func TestValidateToolbarItems(t *testing.T) {
	toolbar := NewMacToolbar()
	button := toolbar.AddButton("Save").SetSymbol("checkmark").OnClick(func(*Context) {})
	search := toolbar.AddSearch("Search").OnSearch(func(*Context, string) {})
	group := toolbar.AddGroup("View", ToolbarGroupSelectOne)
	group.AddButton("Write").OnClick(func(*Context) {})

	if err := validateToolbarItems(toolbar.itemSnapshot()); err != nil {
		t.Fatalf("valid toolbar rejected: %v", err)
	}
	if toolbar.identifier == "" || button.identifier == "" || search.identifier == "" || group.identifier == "" {
		t.Fatal("toolbar and items should receive internal identifiers")
	}
	if button.identifier == search.identifier || button.identifier == group.identifier {
		t.Fatal("toolbar item identifiers should be unique")
	}
	if toolbar.identifier == NewMacToolbar().identifier {
		t.Fatal("toolbar identifiers should be unique per toolbar")
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
	}
	if item.tintColor != nil {
		copyOfColor := *item.tintColor
		result.tintColor = &copyOfColor
	}
	return result
}
