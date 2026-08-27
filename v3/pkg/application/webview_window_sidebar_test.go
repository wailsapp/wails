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
	image := []byte{1, 2, 3}
	item.SetLabel("Published").
		SetSubtitle("Ready for review").
		SetSymbol("checkmark").
		SetImage(image).
		SetTooltip("Ready").
		SetEnabled(false).
		SetHidden(true)
	image[0] = 9
	snapshot := sidebar.snapshot().entries[0].item
	if snapshot.label != "Published" || snapshot.subtitle != "Ready for review" ||
		snapshot.symbolName != "checkmark" || len(snapshot.imageData) != 3 || snapshot.imageData[0] != 1 ||
		snapshot.tooltip != "Ready" ||
		!snapshot.disabled || !snapshot.hidden {
		t.Fatalf("unexpected sidebar item snapshot: %#v", snapshot)
	}
	item.SetImage(nil)
	if len(sidebar.snapshot().entries[0].item.imageData) != 0 {
		t.Fatal("SetImage(nil) should remove the image")
	}
}
