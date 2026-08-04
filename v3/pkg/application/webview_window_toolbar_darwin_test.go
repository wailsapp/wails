//go:build darwin && !ios && !server

package application

import (
	"encoding/json"
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
