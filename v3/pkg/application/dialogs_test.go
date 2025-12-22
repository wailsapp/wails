package application

import (
	"testing"
)

func TestGetDialogID(t *testing.T) {
	// Get first dialog ID
	id1 := getDialogID()

	// Get second dialog ID - should be different
	id2 := getDialogID()
	if id1 == id2 {
		t.Error("getDialogID should return unique IDs")
	}

	// Free first ID
	freeDialogID(id1)

	// Get another ID - could be the freed one
	id3 := getDialogID()
	if id3 == id2 {
		t.Error("getDialogID should not return the same ID as an active dialog")
	}

	// Cleanup
	freeDialogID(id2)
	freeDialogID(id3)
}

func TestFreeDialogID(t *testing.T) {
	id := getDialogID()
	freeDialogID(id)

	// Should be able to get the same ID again after freeing
	newID := getDialogID()
	freeDialogID(newID)
	// Just verify it doesn't panic
}

func TestButton_OnClick(t *testing.T) {
	button := &Button{Label: "Test"}
	called := false

	result := button.OnClick(func() {
		called = true
	})

	// Should return the same button for chaining
	if result != button {
		t.Error("OnClick should return the same button")
	}

	// Callback should be set
	if button.Callback == nil {
		t.Error("Callback should be set")
	}

	// Call the callback
	button.Callback()
	if !called {
		t.Error("Callback should have been called")
	}
}

func TestButton_SetAsDefault(t *testing.T) {
	button := &Button{Label: "Test"}

	result := button.SetAsDefault()

	// Should return the same button for chaining
	if result != button {
		t.Error("SetAsDefault should return the same button")
	}

	if !button.IsDefault {
		t.Error("IsDefault should be true")
	}
}

func TestButton_SetAsCancel(t *testing.T) {
	button := &Button{Label: "Test"}

	result := button.SetAsCancel()

	// Should return the same button for chaining
	if result != button {
		t.Error("SetAsCancel should return the same button")
	}

	if !button.IsCancel {
		t.Error("IsCancel should be true")
	}
}

func TestButton_Chaining(t *testing.T) {
	button := &Button{Label: "OK"}

	button.SetAsDefault().SetAsCancel().OnClick(func() {})

	if !button.IsDefault {
		t.Error("IsDefault should be true after chaining")
	}
	if !button.IsCancel {
		t.Error("IsCancel should be true after chaining")
	}
	if button.Callback == nil {
		t.Error("Callback should be set after chaining")
	}
}

func TestDialogType_Constants(t *testing.T) {
	// Verify dialog type constants are distinct
	types := []DialogType{InfoDialogType, QuestionDialogType, WarningDialogType, ErrorDialogType}
	seen := make(map[DialogType]bool)

	for _, dt := range types {
		if seen[dt] {
			t.Errorf("DialogType %d is duplicated", dt)
		}
		seen[dt] = true
	}
}

func TestMessageDialogOptions_Fields(t *testing.T) {
	opts := MessageDialogOptions{
		DialogType: InfoDialogType,
		Title:      "Test Title",
		Message:    "Test Message",
		Buttons:    []*Button{{Label: "OK"}},
		Icon:       []byte{1, 2, 3},
	}

	if opts.DialogType != InfoDialogType {
		t.Error("DialogType not set correctly")
	}
	if opts.Title != "Test Title" {
		t.Error("Title not set correctly")
	}
	if opts.Message != "Test Message" {
		t.Error("Message not set correctly")
	}
	if len(opts.Buttons) != 1 {
		t.Error("Buttons not set correctly")
	}
	if len(opts.Icon) != 3 {
		t.Error("Icon not set correctly")
	}
}

func TestFileFilter_Fields(t *testing.T) {
	filter := FileFilter{
		DisplayName: "Image Files (*.jpg, *.png)",
		Pattern:     "*.jpg;*.png",
	}

	if filter.DisplayName != "Image Files (*.jpg, *.png)" {
		t.Error("DisplayName not set correctly")
	}
	if filter.Pattern != "*.jpg;*.png" {
		t.Error("Pattern not set correctly")
	}
}

func TestOpenFileDialogOptions_Fields(t *testing.T) {
	opts := OpenFileDialogOptions{
		CanChooseDirectories:    true,
		CanChooseFiles:          true,
		CanCreateDirectories:    true,
		ShowHiddenFiles:         true,
		ResolvesAliases:         true,
		AllowsMultipleSelection: true,
		Title:                   "Open",
		Message:                 "Select a file",
		ButtonText:              "Choose",
		Directory:               "/home",
		Filters: []FileFilter{
			{DisplayName: "All Files", Pattern: "*"},
		},
	}

	if !opts.CanChooseDirectories {
		t.Error("CanChooseDirectories not set correctly")
	}
	if !opts.CanChooseFiles {
		t.Error("CanChooseFiles not set correctly")
	}
	if opts.Title != "Open" {
		t.Error("Title not set correctly")
	}
	if len(opts.Filters) != 1 {
		t.Error("Filters not set correctly")
	}
}

func TestSaveFileDialogOptions_Fields(t *testing.T) {
	opts := SaveFileDialogOptions{
		CanCreateDirectories: true,
		ShowHiddenFiles:      true,
		Title:                "Save",
		Message:              "Save as",
		Directory:            "/home",
		Filename:             "file.txt",
		ButtonText:           "Save",
		Filters: []FileFilter{
			{DisplayName: "Text Files", Pattern: "*.txt"},
		},
	}

	if !opts.CanCreateDirectories {
		t.Error("CanCreateDirectories not set correctly")
	}
	if opts.Title != "Save" {
		t.Error("Title not set correctly")
	}
	if opts.Filename != "file.txt" {
		t.Error("Filename not set correctly")
	}
}
