package application

import "testing"

func TestValidateToolbarItems(t *testing.T) {
	toolbar := NewMacToolbar()
	button := toolbar.AddButton("Save").OnClick(func(*Context) {})
	search := toolbar.AddSearch("Search").OnSearch(func(*Context, string) {})
	group := toolbar.AddGroup("View", ToolbarGroupSelectOne)
	group.AddButton("Write").OnClick(func(*Context) {})

	if err := validateToolbarItems(toolbar.items); err != nil {
		t.Fatalf("valid toolbar rejected: %v", err)
	}
	if button.identifier == "" || search.identifier == "" || group.identifier == "" {
		t.Fatal("toolbar items should receive internal identifiers")
	}
	if button.identifier == search.identifier || button.identifier == group.identifier {
		t.Fatal("toolbar item identifiers should be unique")
	}
}

func TestValidateToolbarItemsRejectsMissingCallbacks(t *testing.T) {
	toolbar := NewMacToolbar()
	toolbar.AddButton("Save")
	if err := validateToolbarItems(toolbar.items); err == nil {
		t.Fatal("button without callback should be rejected")
	}

	search := NewMacToolbar().AddSearch("Search")
	if err := validateToolbarItems([]*MacToolbarItem{search}); err == nil {
		t.Fatal("search field without callback should be rejected")
	}
}

func TestToolbarGroupRejectsInvalidMembers(t *testing.T) {
	toolbar := NewMacToolbar()
	group := toolbar.AddGroup("View", ToolbarGroupSelectOne)
	group.items = append(group.items, newMacToolbarItem(toolbar, ToolbarSearchField, "Search"))
	if err := validateToolbarItems(toolbar.items); err == nil {
		t.Fatal("group should reject non-button members")
	}
}

func TestToolbarMutatorsBeforeInstallation(t *testing.T) {
	toolbar := NewMacToolbar()
	item := toolbar.AddButton("Details").OnClick(func(*Context) {})
	item.SetLabel("Details (3)").SetEnabled(false).SetHidden(true).SetBadgeCount(3)
	if item.label != "Details (3)" || !item.disabled || !item.hidden || item.badgeCount != 3 {
		t.Fatal("mutators should update the item before installation")
	}
}
