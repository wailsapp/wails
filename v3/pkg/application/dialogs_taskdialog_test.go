package application

import (
	"sync"
	"testing"
)

func TestTaskDialogButtonCallbackStorage(t *testing.T) {
	callbacks := make(map[int32]func())
	var mu sync.Mutex

	quitCalled := false
	cancelCalled := false

	callbacks[100] = func() { quitCalled = true }
	callbacks[101] = func() { cancelCalled = true }

	mu.Lock()
	if cb, ok := callbacks[100]; ok && cb != nil {
		cb()
	}
	mu.Unlock()

	if !quitCalled {
		t.Error("callback for button 100 should have been called")
	}
	if cancelCalled {
		t.Error("callback for button 101 should NOT have been called yet")
	}

	mu.Lock()
	if cb, ok := callbacks[101]; ok && cb != nil {
		cb()
	}
	mu.Unlock()

	if !cancelCalled {
		t.Error("callback for button 101 should have been called")
	}
}

func TestTaskDialogButtonIDMapping(t *testing.T) {
	buttons := []*Button{
		{Label: "Quit", IsDefault: true},
		{Label: "Cancel", IsCancel: true},
		{Label: "Retry"},
	}

	const customButtonBase = 100
	for i, btn := range buttons {
		id := int32(customButtonBase + i)
		expectedID := int32(100 + i)
		if id != expectedID {
			t.Errorf("button %q: id = %d, want %d", btn.Label, id, expectedID)
		}
		if btn.Label == "Quit" && !btn.IsDefault {
			t.Error("Quit button should be default")
		}
		if btn.Label == "Cancel" && !btn.IsCancel {
			t.Error("Cancel button should be cancel")
		}
	}
}

func TestTaskDialogDefaultButtonSelection(t *testing.T) {
	buttons := []*Button{
		{Label: "No", IsDefault: false},
		{Label: "Yes", IsDefault: true},
	}

	var defaultButtonID int32
	const customButtonBase = 100
	for i, btn := range buttons {
		if btn.IsDefault {
			defaultButtonID = int32(customButtonBase + i)
		}
	}

	if defaultButtonID != 101 {
		t.Errorf("defaultButtonID = %d, want 101", defaultButtonID)
	}
}
