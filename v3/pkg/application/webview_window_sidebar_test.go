package application

import "testing"

func TestMacSidebarBuildsNativeSourceListWithoutUserIDs(t *testing.T) {
	sidebar := NewMacSidebar()
	recents := sidebar.AddItem("Recents").SetSymbol("clock")
	section := sidebar.AddSection("Favorites")
	documents := section.AddItem("Documents").SetSymbol("doc").SetTooltip("Open documents")

	snapshot := sidebar.snapshot()
	if len(snapshot.entries) != 2 || snapshot.entries[0].item == nil || snapshot.entries[1].section == nil {
		t.Fatal("sidebar should preserve root items and sections in declaration order")
	}
	if recents.internalID == 0 || documents.internalID == 0 || recents.internalID == documents.internalID {
		t.Fatal("sidebar nodes should receive unique generated internal IDs")
	}
	if snapshot.entries[1].section.items[0].symbolName != "doc" {
		t.Fatal("native item presentation should be retained in the snapshot")
	}
}

func TestMacSidebarSelectionAndCallbacks(t *testing.T) {
	sidebar := NewMacSidebar()
	item := sidebar.AddItem("Recents")
	count := 0
	item.OnClick(func(*Context) { count++ })
	sidebar.registerItems()
	sidebar.SetSelectedItem(item)
	if sidebar.snapshot().selectedItemID != item.internalID {
		t.Fatal("programmatic selection should update sidebar state")
	}
	handleMacSidebarItemSelected(item.internalID)
	if count != 1 || sidebar.selected != item {
		t.Fatalf("native selection callback count = %d", count)
	}
	item.OnClick(nil)
	handleMacSidebarItemSelected(item.internalID)
	if count != 1 {
		t.Fatal("OnClick(nil) should clear the callback")
	}
}

func TestMacSidebarLiveItemState(t *testing.T) {
	sidebar := NewMacSidebar()
	item := sidebar.AddItem("Draft")
	item.SetLabel("Published").SetSymbol("checkmark").SetTooltip("Ready").SetEnabled(false).SetHidden(true)
	snapshot := sidebar.snapshot().entries[0].item
	if snapshot.label != "Published" || snapshot.symbolName != "checkmark" || snapshot.tooltip != "Ready" ||
		!snapshot.disabled || !snapshot.hidden {
		t.Fatalf("unexpected sidebar item snapshot: %#v", snapshot)
	}
}

// newRegisteredTestSidebar wires a sidebar into a split pane and registers
// the pane so pane-scoped native callbacks can be exercised without AppKit.
func newRegisteredTestSidebar(t *testing.T) (*MacSidebar, uint64) {
	t.Helper()
	sidebar := NewMacSidebar()
	split := NewMacSplitView()
	pane := split.AddSidebar(sidebar)
	split.AddPrimaryContent()
	internal := split.paneSnapshot()[0]
	registerMacSplitPane(internal)
	t.Cleanup(func() {
		unregisterMacSplitPane(internal.internalID)
		for _, item := range sidebar.itemHandles() {
			unregisterMacSidebarItem(item.internalID)
		}
	})
	return sidebar, pane.internalID
}

func sidebarLabels(items []*MacSidebarItem) string {
	result := ""
	for index, item := range items {
		if index > 0 {
			result += ","
		}
		result += item.label
	}
	return result
}

func TestMacSidebarTreeSnapshot(t *testing.T) {
	sidebar := NewMacSidebar()
	section := sidebar.AddSection("Mailboxes")
	inbox := section.AddItem("Inbox").
		SetBadge(12).
		SetAccessorySymbol("pin.fill").
		SetTintColor(&RGBA{Red: 10, Green: 20, Blue: 30, Alpha: 255}).
		SetEditable(true)
	work := inbox.AddItem("Work")
	invoices := work.AddItem("Invoices")
	inbox.SetExpanded(true)

	if inbox.Parent() != nil || work.Parent() != inbox || invoices.Parent() != work {
		t.Fatal("nested rows should report their parent row")
	}
	if inbox.Section() != section || invoices.Section() != section {
		t.Fatal("Section should walk up nested rows to the owning section")
	}
	if !inbox.IsExpanded() || work.IsExpanded() {
		t.Fatal("nested rows start collapsed until SetExpanded")
	}
	snapshot := sidebar.snapshot().entries[0].section.items[0]
	if snapshot.badge != 12 || snapshot.accessorySymbol != "pin.fill" || !snapshot.editable || !snapshot.expanded ||
		snapshot.tintColor == nil || snapshot.tintColor.Blue != 30 {
		t.Fatalf("row presentation missing from snapshot: %#v", snapshot)
	}
	if len(snapshot.children) != 1 || snapshot.children[0].label != "Work" ||
		len(snapshot.children[0].children) != 1 || snapshot.children[0].children[0].internalID != invoices.internalID {
		t.Fatalf("snapshot should carry the whole tree: %#v", snapshot.children)
	}
	if len(sidebar.itemHandles()) != 3 {
		t.Fatal("itemHandles should include nested rows")
	}
	inbox.SetBadge(-4)
	if inbox.Badge() != 0 {
		t.Fatal("negative badge counts should clear the badge")
	}
	tint := RGBA{Red: 1}
	inbox.SetTintColor(&tint)
	tint.Red = 99
	if sidebar.snapshot().entries[0].section.items[0].tintColor.Red != 1 {
		t.Fatal("SetTintColor should copy the caller's colour")
	}
}

func TestMacSidebarRemovalUnregistersCallbacks(t *testing.T) {
	sidebar := NewMacSidebar()
	section := sidebar.AddSection("Mailboxes")
	inbox := section.AddItem("Inbox")
	work := inbox.AddItem("Work")
	recents := sidebar.AddItem("Recents")
	clicks := 0
	work.OnClick(func(*Context) { clicks++ })
	sidebar.registerItems()
	t.Cleanup(func() {
		for _, id := range []uint64{inbox.internalID, work.internalID, recents.internalID} {
			unregisterMacSidebarItem(id)
		}
	})
	sidebar.SetSelectedItem(work)

	inbox.Remove()
	if macSidebarItemRegistry[inbox.internalID] != nil || macSidebarItemRegistry[work.internalID] != nil {
		t.Fatal("removing a row should unregister its whole subtree")
	}
	handleMacSidebarItemSelected(work.internalID)
	if clicks != 0 {
		t.Fatal("callbacks on removed rows must not fire")
	}
	if sidebar.SelectedItem() != nil || len(sidebar.SelectedItems()) != 0 {
		t.Fatal("removing the selected row should clear the selection")
	}
	if inbox.AddItem("Late") != nil || len(section.Items()) != 0 {
		t.Fatal("removed rows should be inert")
	}
	inbox.SetLabel("Late")
	if inbox.label != "Inbox" {
		t.Fatal("setters on removed rows should be no-ops")
	}

	sidebar.RemoveSection(section)
	if len(sidebar.Sections()) != 0 || section.AddItem("Late") != nil {
		t.Fatal("removed sections should leave the sidebar and become inert")
	}
	sidebar.Remove(recents)
	if len(sidebar.snapshot().entries) != 0 || macSidebarItemRegistry[recents.internalID] != nil {
		t.Fatal("removing a root row should drop it from the snapshot and the registry")
	}
}

func TestMacSidebarMoveOrdering(t *testing.T) {
	sidebar, paneID := newRegisteredTestSidebar(t)
	favorites := sidebar.AddSection("Favorites")
	a := favorites.AddItem("A")
	b := favorites.AddItem("B")
	c := favorites.AddItem("C")
	archive := sidebar.AddSection("Archive")
	sidebar.registerItems()

	var moved *MacSidebarItem
	var movedSection *MacSidebarSection
	movedIndex := -1
	sidebar.OnMove(func(_ *Context, item *MacSidebarItem, section *MacSidebarSection, index int) {
		moved, movedSection, movedIndex = item, section, index
	})

	// Native drop indexes count rows before the dragged row is removed.
	handleMacSidebarItemMoved(macSidebarMoveEvent{paneID: paneID, itemID: a.internalID, parentID: favorites.internalID, index: 2})
	if got := sidebarLabels(favorites.Items()); got != "B,A,C" {
		t.Fatalf("favorites = %s, want B,A,C", got)
	}
	if moved != a || movedSection != favorites || movedIndex != 1 {
		t.Fatalf("OnMove got (%v, %v, %d)", moved == a, movedSection == favorites, movedIndex)
	}

	handleMacSidebarItemMoved(macSidebarMoveEvent{paneID: paneID, itemID: c.internalID, parentID: archive.internalID, index: 0})
	if sidebarLabels(favorites.Items()) != "B,A" || sidebarLabels(archive.Items()) != "C" || c.Section() != archive {
		t.Fatal("rows should move between sections")
	}
	if movedSection != archive || movedIndex != 0 {
		t.Fatal("OnMove should report the destination section and index")
	}

	handleMacSidebarItemMoved(macSidebarMoveEvent{paneID: paneID, itemID: b.internalID, parentID: 0, index: 0})
	entries := sidebar.snapshot().entries
	if entries[0].item == nil || entries[0].item.label != "B" || b.Section() != nil || movedSection != nil {
		t.Fatal("rows should move to the sidebar root")
	}

	child := a.AddItem("Child")
	sidebar.registerItems()
	handleMacSidebarItemMoved(macSidebarMoveEvent{paneID: paneID, itemID: a.internalID, parentID: child.internalID, index: 0})
	if a.Parent() != nil || len(a.Items()) != 1 || child.Parent() != a {
		t.Fatal("a row must not be reparented beneath its own descendant")
	}

	favorites.AddItem("D")
	favorites.Move(a, 5)
	if got := sidebarLabels(favorites.Items()); got != "D,A" {
		t.Fatalf("Move should clamp to the end: %s", got)
	}
}

func TestMacSidebarMultipleSelectionCallbacks(t *testing.T) {
	sidebar, paneID := newRegisteredTestSidebar(t)
	section := sidebar.AddSection("Notes")
	first := section.AddItem("First")
	second := section.AddItem("Second")
	sidebar.SetAllowsMultipleSelection(true)
	sidebar.registerItems()

	clicks := 0
	second.OnClick(func(*Context) { clicks++ })
	var selection []*MacSidebarItem
	calls := 0
	sidebar.OnSelectionChange(func(_ *Context, items []*MacSidebarItem) {
		selection = items
		calls++
	})

	handleMacSidebarSelectionChanged(macSidebarSelectionEvent{paneID: paneID, itemIDs: []uint64{second.internalID, first.internalID}})
	if clicks != 1 || calls != 1 || len(selection) != 2 || selection[0] != second {
		t.Fatalf("clicks=%d calls=%d selection=%d", clicks, calls, len(selection))
	}
	if sidebar.SelectedItem() != second || len(sidebar.SelectedItems()) != 2 {
		t.Fatal("the sidebar should track the clicked row and the full selection")
	}

	handleMacSidebarSelectionChanged(macSidebarSelectionEvent{paneID: paneID})
	if calls != 2 || len(selection) != 0 || sidebar.SelectedItem() != nil || clicks != 1 {
		t.Fatal("clearing the selection should notify without firing OnClick")
	}
	if !sidebar.snapshot().allowsMultipleSelection {
		t.Fatal("multiple selection should reach the snapshot")
	}
}

func TestMacSidebarContextMenuResolution(t *testing.T) {
	sidebar, paneID := newRegisteredTestSidebar(t)
	inbox := sidebar.AddItem("Inbox")
	drafts := sidebar.AddItem("Drafts")
	sidebar.registerItems()

	rowMenu, dynamicMenu, fallback := NewMenu(), NewMenu(), NewMenu()
	inbox.SetContextMenu(rowMenu)
	sidebar.SetContextMenu(fallback)
	var asked *MacSidebarItem
	askedCalls := 0
	sidebar.OnContextMenu(func(_ *Context, item *MacSidebarItem) *Menu {
		asked = item
		askedCalls++
		if item == drafts {
			return dynamicMenu
		}
		return nil
	})

	if resolveMacSidebarContextMenu(paneID, inbox.internalID) != rowMenu || askedCalls != 0 {
		t.Fatal("a row's own menu should win without consulting the callback")
	}
	if resolveMacSidebarContextMenu(paneID, drafts.internalID) != dynamicMenu || asked != drafts {
		t.Fatal("the callback should receive the clicked row")
	}
	if resolveMacSidebarContextMenu(paneID, 0) != fallback || asked != nil {
		t.Fatal("empty-area clicks should pass nil and fall back to the sidebar menu")
	}
	if resolveMacSidebarContextMenu(paneID+1000, 0) != nil {
		t.Fatal("unknown panes should show no menu")
	}
}

func TestMacSidebarRenameAndExpansionCallbacks(t *testing.T) {
	sidebar := NewMacSidebar()
	item := sidebar.AddItem("Draft")
	child := item.AddItem("Child")
	sidebar.registerItems()
	t.Cleanup(func() {
		unregisterMacSidebarItem(item.internalID)
		unregisterMacSidebarItem(child.internalID)
	})

	var renamed string
	item.OnRename(func(_ *Context, label string) { renamed = label })
	handleMacSidebarItemRenamed(macSidebarRenameEvent{itemID: item.internalID, label: "Ignored"})
	if renamed != "" || item.label != "Draft" {
		t.Fatal("rows that are not editable ignore native renames")
	}
	item.SetEditable(true)
	handleMacSidebarItemRenamed(macSidebarRenameEvent{itemID: item.internalID, label: "Final"})
	if renamed != "Final" || item.label != "Final" {
		t.Fatal("renames should update the handle before OnRename runs")
	}

	expandedCalls := 0
	item.OnExpandedChange(func(_ *Context, expanded bool) {
		if expanded {
			expandedCalls++
		}
	})
	handleMacSidebarItemExpanded(macSidebarExpandedEvent{itemID: item.internalID, expanded: true})
	handleMacSidebarItemExpanded(macSidebarExpandedEvent{itemID: item.internalID, expanded: true})
	if expandedCalls != 1 || !item.IsExpanded() {
		t.Fatalf("expansion callback count = %d", expandedCalls)
	}
}
