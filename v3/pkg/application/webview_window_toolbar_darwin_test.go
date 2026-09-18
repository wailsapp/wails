//go:build darwin && !ios && !server

package application

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestMacToolbarShareProviderJSONPreservesRepresentations(t *testing.T) {
	var payload struct {
		ProviderID      uint   `json:"providerID"`
		Subject         string `json:"subject"`
		SuggestedName   string `json:"suggestedName"`
		Representations []struct {
			ContentType MacShareContentType `json:"contentType"`
		} `json:"representations"`
	}
	encoded := macToolbarShareProviderJSON(42, "A note", "Daymark Note", []MacShareRepresentation{
		{ContentType: MacShareTypeHTML},
		{ContentType: MacShareTypePlainText},
	})
	if err := json.Unmarshal([]byte(encoded), &payload); err != nil {
		t.Fatalf("decode share payload: %v", err)
	}
	if payload.ProviderID != 42 || payload.Subject != "A note" || payload.SuggestedName != "Daymark Note" || len(payload.Representations) != 2 {
		t.Fatalf("unexpected share payload: %#v", payload)
	}
	if payload.Representations[0].ContentType != MacShareTypeHTML {
		t.Fatalf("unexpected HTML representation: %#v", payload.Representations[0])
	}
}

func TestClearMacToolbarStateRemovesCallbacksAndOwnership(t *testing.T) {
	window := &WebviewWindow{}
	toolbar := NewMacToolbar()
	item := toolbar.AddButton("Save").OnClick(func(*Context) {})
	firstID := nextToolbarNativeID()
	secondID := nextToolbarNativeID()
	addToToolbarItemMap(firstID, item)
	addToToolbarItemMap(secondID, item)

	toolbar.stateLock.Lock()
	toolbar.state = &macToolbarState{
		window:  window,
		itemIDs: []uint{firstID, secondID},
	}
	toolbar.stateLock.Unlock()

	clearMacToolbarState(toolbar, window, true)

	if getToolbarItemByID(firstID) != nil || getToolbarItemByID(secondID) != nil {
		t.Fatal("clearing toolbar state should remove all callback mappings")
	}
	toolbar.stateLock.RLock()
	defer toolbar.stateLock.RUnlock()
	if toolbar.state.native != nil || len(toolbar.state.itemIDs) != 0 || toolbar.state.window != nil {
		t.Fatal("clearing toolbar state should invalidate the native handle and release ownership")
	}
}

func TestMacToolbarLayoutPayloadDistinguishesDefaultAndAllowed(t *testing.T) {
	toolbar := NewMacToolbar()
	toolbar.AddSidebarToggle()
	toolbar.AddSidebarTrackingSeparator()
	back := toolbar.AddButton("Back").OnClick(func(*Context) {}).SetPersistenceKey("back")
	toolbar.AddSpace()
	toolbar.AddFlexibleSpace()
	palette := toolbar.AddButton("Palette only").OnClick(func(*Context) {}).SetInDefaultSet(false)
	toolbar.AddInspectorTrackingSeparator()
	toolbar.AddInspectorToggle()
	toolbar.SetCenteredItems(back)

	payload := macToolbarLayoutPayloadFor(toolbar, true, "back")
	if !payload.Customizable || payload.Moved != "back" {
		t.Fatalf("payload flags = %#v", payload)
	}
	kinds := func(entries []macToolbarLayoutEntry) []string {
		result := make([]string, 0, len(entries))
		for _, entry := range entries {
			result = append(result, entry.Kind)
		}
		return result
	}
	wantAllowed := []string{"sidebarToggle", "sidebarSeparator", "item", "space", "flexibleSpace", "item", "inspectorSeparator", "inspectorToggle"}
	if got := kinds(payload.Allowed); strings.Join(got, ",") != strings.Join(wantAllowed, ",") {
		t.Fatalf("allowed kinds = %v", got)
	}
	if len(payload.Default) != len(payload.Allowed)-1 {
		t.Fatalf("an item outside the default set must still be allowed: default=%d allowed=%d", len(payload.Default), len(payload.Allowed))
	}
	for _, entry := range payload.Default {
		if entry.ID == palette.identifier {
			t.Fatal("palette-only item leaked into the default layout")
		}
	}
	if payload.Allowed[2].ID != "back" {
		t.Fatalf("persistence key should be the layout identifier, got %q", payload.Allowed[2].ID)
	}
	if len(payload.Centered) != 1 || payload.Centered[0].ID != "back" {
		t.Fatalf("centred entries = %#v", payload.Centered)
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("encode layout: %v", err)
	}
	if !strings.Contains(string(encoded), `"kind":"space"`) {
		t.Fatalf("layout JSON should carry entry kinds: %s", encoded)
	}
}

func TestMacToolbarItemSignatureTracksGroupMembers(t *testing.T) {
	toolbar := NewMacToolbar()
	button := toolbar.AddButton("Save").OnClick(func(*Context) {})
	group := toolbar.AddGroup("Mode", ToolbarGroupSelectOne)
	write := group.AddButton("Write").OnClick(func(*Context) {})

	if signature, members := macToolbarItemSignature(snapshotMacToolbarItem(button)); signature != "" || members != nil {
		t.Fatal("non-group items have no rebuild signature")
	}
	before, members := macToolbarItemSignature(snapshotMacToolbarItem(group.MacToolbarItem))
	if len(members) != 1 || members[0] != write.identifier {
		t.Fatalf("group members = %v", members)
	}
	group.AddButton("Preview").OnClick(func(*Context) {})
	after, _ := macToolbarItemSignature(snapshotMacToolbarItem(group.MacToolbarItem))
	if before == after {
		t.Fatal("adding a member must change the group signature so the native group is rebuilt")
	}
	// Relabelling a member does not require a rebuild.
	write.SetLabel("Compose")
	relabelled, _ := macToolbarItemSignature(snapshotMacToolbarItem(group.MacToolbarItem))
	if relabelled != after {
		t.Fatal("label changes must not force a group rebuild")
	}
}

func TestMacToolbarLayoutKinds(t *testing.T) {
	cases := map[macToolbarItemKind]string{
		toolbarButton:                     "item",
		toolbarGroup:                      "item",
		toolbarSearchField:                "item",
		toolbarShare:                      "item",
		toolbarMenu:                       "item",
		toolbarSpace:                      "space",
		toolbarFlexibleSpace:              "flexibleSpace",
		toolbarSidebarToggle:              "sidebarToggle",
		toolbarSidebarTrackingSeparator:   "sidebarSeparator",
		toolbarInspectorToggle:            "inspectorToggle",
		toolbarInspectorTrackingSeparator: "inspectorSeparator",
	}
	for kind, want := range cases {
		if got := macToolbarLayoutKind(kind); got != want {
			t.Fatalf("layout kind for %d = %q, want %q", kind, got, want)
		}
	}
	if macToolbarKindIsWailsOwned(toolbarSpace) || !macToolbarKindIsWailsOwned(toolbarMenu) {
		t.Fatal("ownership classification is wrong")
	}
}

func TestClearMacToolbarStateIgnoresDifferentOwner(t *testing.T) {
	owner := &WebviewWindow{}
	other := &WebviewWindow{}
	toolbar := NewMacToolbar()
	toolbar.state = &macToolbarState{window: owner}

	clearMacToolbarState(toolbar, other, true)

	toolbar.stateLock.RLock()
	defer toolbar.stateLock.RUnlock()
	if toolbar.state.window != owner {
		t.Fatal("a different window must not clear toolbar ownership")
	}
}
