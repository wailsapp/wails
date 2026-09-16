package application

import (
	"fmt"
	"testing"
	"time"
)

func TestOpenRecentRole(t *testing.T) {
	item := NewRole(OpenRecent)
	if item == nil {
		t.Fatal("NewRole(OpenRecent) returned nil")
	}
	if !item.IsSubmenu() {
		t.Fatal("OpenRecent role should be a submenu item")
	}
	if item.Label() != "Open Recent" {
		t.Fatalf("label = %q, want Open Recent", item.Label())
	}
	if item.role != OpenRecent {
		t.Fatalf("role = %v, want OpenRecent", item.role)
	}
	menu := NewMenu()
	menu.AddRole(OpenRecent)
	if menu.FindByRole(OpenRecent) == nil {
		t.Fatal("FindByRole(OpenRecent) did not find the item")
	}
}

func TestMenuItemMixedState(t *testing.T) {
	item := NewMenuItemCheckbox("Wrap", true)
	item.SetMixed()
	if !item.Mixed() {
		t.Fatal("expected Mixed() true after SetMixed")
	}
	if item.Checked() {
		t.Fatal("Checked() must report false while mixed")
	}

	// A click on a mixed checkbox turns it fully on.
	item.handleClick()
	if item.Mixed() {
		t.Fatal("click should leave the mixed state")
	}
	if !item.Checked() {
		t.Fatal("click on a mixed checkbox should turn it on")
	}

	item.SetMixed()
	item.SetChecked(false)
	if item.Mixed() {
		t.Fatal("SetChecked should clear the mixed state")
	}
	if item.Checked() {
		t.Fatal("SetChecked(false) should leave the item unchecked")
	}
}

func TestMenuItemAlternateAndIndentation(t *testing.T) {
	item := NewMenuItem("Close All").SetAccelerator("Cmd+OptionOrAlt+w").SetAlternate(true)
	if !item.Alternate() {
		t.Fatal("expected Alternate() true")
	}
	item.SetAlternate(false)
	if item.Alternate() {
		t.Fatal("expected Alternate() false")
	}

	item.SetIndentationLevel(3)
	if item.IndentationLevel() != 3 {
		t.Fatalf("IndentationLevel = %d, want 3", item.IndentationLevel())
	}
	item.SetIndentationLevel(-4)
	if item.IndentationLevel() != 0 {
		t.Fatalf("negative level should clamp to 0, got %d", item.IndentationLevel())
	}
	item.SetIndentationLevel(99)
	if item.IndentationLevel() != 15 {
		t.Fatalf("large level should clamp to 15, got %d", item.IndentationLevel())
	}
}

func TestMenuItemBadgeModel(t *testing.T) {
	item := NewMenuItem("Inbox")
	if item.HasBadge() {
		t.Fatal("new item should have no badge")
	}
	item.SetBadge(7)
	if !item.HasBadge() || item.BadgeCount() != 7 || item.BadgeText() != "" {
		t.Fatalf("after SetBadge(7): has=%v count=%d text=%q", item.HasBadge(), item.BadgeCount(), item.BadgeText())
	}
	item.SetBadgeText("New")
	if !item.HasBadge() || item.BadgeCount() != 0 || item.BadgeText() != "New" {
		t.Fatalf("after SetBadgeText: has=%v count=%d text=%q", item.HasBadge(), item.BadgeCount(), item.BadgeText())
	}
	item.SetBadge(0)
	if item.HasBadge() || item.BadgeText() != "" || item.BadgeCount() != 0 {
		t.Fatal("SetBadge(0) should clear the badge")
	}
	item.SetBadgeText("x").ClearBadge()
	if item.HasBadge() {
		t.Fatal("ClearBadge should clear the badge")
	}
}

func TestMenuItemSymbolModel(t *testing.T) {
	item := NewMenuItem("Star").SetSymbol("star.fill")
	if item.Symbol() != "star.fill" {
		t.Fatalf("Symbol = %q", item.Symbol())
	}
	clone := item.Clone()
	if clone.Symbol() != "star.fill" {
		t.Fatal("Clone should copy the symbol")
	}
}

func TestMenuSectionHeader(t *testing.T) {
	menu := NewMenu()
	header := menu.AddSectionHeader("Recent")
	if !header.IsSectionHeader() {
		t.Fatal("expected IsSectionHeader() true")
	}
	if header.Enabled() {
		t.Fatal("section headers are not enabled")
	}
	if header.Label() != "Recent" {
		t.Fatalf("label = %q", header.Label())
	}
	if menu.ItemAt(0) != header {
		t.Fatal("header should be appended to the menu")
	}
	if !header.Clone().IsSectionHeader() {
		t.Fatal("Clone should keep the section header kind")
	}
}

func TestMenuAddPalette(t *testing.T) {
	menu := NewMenu()
	got := make(chan int, 1)
	colours := []RGBA{NewRGB(255, 0, 0), NewRGB(0, 255, 0), NewRGB(0, 0, 255)}
	item := menu.AddPalette([]string{"circle.fill"}, colours, 1, func(ctx *Context, index int) {
		if ctx.ClickedMenuItem() == nil {
			got <- -100
			return
		}
		got <- index
	})
	if !item.IsPalette() {
		t.Fatal("expected IsPalette() true")
	}
	if item.PaletteSelected() != 1 {
		t.Fatalf("PaletteSelected = %d, want 1", item.PaletteSelected())
	}
	if NewMenuItem("x").PaletteSelected() != -1 {
		t.Fatal("non-palette items report -1")
	}

	item.handlePaletteSelection(2)
	select {
	case index := <-got:
		if index != 2 {
			t.Fatalf("callback index = %d, want 2", index)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("palette callback did not fire")
	}
	if item.PaletteSelected() != 2 {
		t.Fatalf("PaletteSelected = %d after selection, want 2", item.PaletteSelected())
	}

	// Out of range clears the selection without a callback.
	item.handlePaletteSelection(9)
	if item.PaletteSelected() != -1 {
		t.Fatalf("out-of-range selection should clear, got %d", item.PaletteSelected())
	}
	select {
	case index := <-got:
		t.Fatalf("unexpected callback for out-of-range index: %d", index)
	case <-time.After(50 * time.Millisecond):
	}

	if out := menu.AddPalette(nil, colours, 7, nil); out.PaletteSelected() != -1 {
		t.Fatalf("selected index beyond colours should become -1, got %d", out.PaletteSelected())
	}
	clone := item.Clone()
	if !clone.IsPalette() || len(clone.paletteColours) != 3 {
		t.Fatal("Clone should copy palette colours")
	}
}

func TestMenuManagerRecentDocumentsBookkeeping(t *testing.T) {
	mm := newMenuManager(&App{})
	if len(mm.RecentDocuments()) != 0 {
		t.Fatal("expected an empty recent list")
	}
	mm.AddRecentDocument("/tmp/a.txt")
	mm.AddRecentDocument("/tmp/b.txt")
	mm.AddRecentDocument("/tmp/a.txt")
	mm.AddRecentDocument("")
	got := mm.RecentDocuments()
	if len(got) != 2 || got[0] != "/tmp/a.txt" || got[1] != "/tmp/b.txt" {
		t.Fatalf("recent list = %v, want [/tmp/a.txt /tmp/b.txt]", got)
	}

	// Returned slice is a copy.
	got[0] = "changed"
	if mm.RecentDocuments()[0] != "/tmp/a.txt" {
		t.Fatal("RecentDocuments should return a copy")
	}

	for i := 0; i < maxRecentDocuments+5; i++ {
		mm.AddRecentDocument(fmt.Sprintf("/tmp/file-%d.txt", i))
	}
	if n := len(mm.RecentDocuments()); n != maxRecentDocuments {
		t.Fatalf("list should cap at %d, got %d", maxRecentDocuments, n)
	}
	if mm.RecentDocuments()[0] != fmt.Sprintf("/tmp/file-%d.txt", maxRecentDocuments+4) {
		t.Fatal("most recent entry should be first")
	}

	mm.ClearRecentDocuments()
	if len(mm.RecentDocuments()) != 0 {
		t.Fatal("ClearRecentDocuments should empty the list")
	}
}

func TestMenuManagerDockMenu(t *testing.T) {
	mm := newMenuManager(&App{})
	if mm.resolveDockMenu() != nil {
		t.Fatal("no dock menu by default")
	}
	static := NewMenu()
	static.Add("Static")
	mm.SetDockMenu(static)
	if mm.DockMenu() != static || mm.resolveDockMenu() != static {
		t.Fatal("static dock menu should be returned")
	}

	calls := 0
	var built []*Menu
	mm.OnDockMenu(func() *Menu {
		calls++
		menu := NewMenu()
		menu.Add(fmt.Sprintf("Dynamic %d", calls))
		built = append(built, menu)
		return menu
	})
	first := mm.resolveDockMenu()
	if first == nil || first == static || calls != 1 {
		t.Fatal("dynamic dock menu should win over the static one")
	}
	if first.ItemAt(0) == nil {
		t.Fatal("dynamic menu should keep its items until the next request")
	}
	second := mm.resolveDockMenu()
	if second == first || calls != 2 {
		t.Fatal("each request should build a new menu")
	}
	if first.ItemAt(0) != nil {
		t.Fatal("previous dynamic menu should be destroyed when a new one is built")
	}
	if static.ItemAt(0) == nil {
		t.Fatal("static menu must not be destroyed")
	}

	mm.OnDockMenu(nil)
	if mm.resolveDockMenu() != static {
		t.Fatal("unregistering the builder falls back to the static menu")
	}
	mm.SetDockMenu(nil)
	if mm.resolveDockMenu() != nil {
		t.Fatal("SetDockMenu(nil) removes the dock menu")
	}
}
