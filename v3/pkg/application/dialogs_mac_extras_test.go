package application

import (
	"strings"
	"testing"
)

func TestMessageDialog_SuppressionDefaults(t *testing.T) {
	dialog := newMessageDialog(InfoDialogType)
	if dialog.showsSuppression {
		t.Error("suppression should be off by default")
	}
	if dialog.Suppressed() {
		t.Error("Suppressed should be false before the dialog is shown")
	}
	if dialog.helpCallback != nil || dialog.accessoryView != nil {
		t.Error("help and accessory view should be unset by default")
	}
}

func TestMessageDialog_SetSuppression(t *testing.T) {
	dialog := newMessageDialog(WarningDialogType)
	result := dialog.SetSuppression("Never warn me again")
	if result != dialog {
		t.Error("SetSuppression should return the dialog for chaining")
	}
	if !dialog.showsSuppression {
		t.Error("SetSuppression should enable the checkbox")
	}
	if dialog.suppressionLabel != "Never warn me again" {
		t.Errorf("unexpected label %q", dialog.suppressionLabel)
	}

	// An empty label keeps the checkbox but uses the system title.
	dialog.SetSuppression("")
	if !dialog.showsSuppression || dialog.suppressionLabel != "" {
		t.Error("SetSuppression(\"\") should keep the checkbox with the default title")
	}
}

func TestMessageDialog_OnSuppression(t *testing.T) {
	dialog := newMessageDialog(QuestionDialogType).SetSuppression("Do not ask again")
	var got []bool
	dialog.OnSuppression(func(ticked bool) {
		got = append(got, ticked)
	})

	dialog.setSuppressed(true)
	if !dialog.Suppressed() {
		t.Error("Suppressed should reflect the recorded state")
	}
	dialog.setSuppressed(false)
	if dialog.Suppressed() {
		t.Error("Suppressed should be cleared when the checkbox is not ticked")
	}
	if len(got) != 2 || !got[0] || got[1] {
		t.Errorf("OnSuppression should be called with each state, got %v", got)
	}
}

func TestMessageDialog_SetHelp(t *testing.T) {
	dialog := newMessageDialog(InfoDialogType)
	called := false
	if dialog.SetHelp(func() { called = true }) != dialog {
		t.Error("SetHelp should return the dialog for chaining")
	}
	if dialog.helpCallback == nil {
		t.Fatal("SetHelp should record the callback")
	}
	dialog.helpCallback()
	if !called {
		t.Error("help callback should be invocable")
	}
}

func TestOpenFileDialog_AddContentType(t *testing.T) {
	dialog := newOpenFileDialog()
	defer freeDialogID(dialog.id)

	dialog.AddContentType(" public.png ").AddContentType("").AddContentType("com.adobe.pdf")
	if len(dialog.contentTypes) != 2 {
		t.Fatalf("expected 2 content types, got %v", dialog.contentTypes)
	}
	if dialog.contentTypes[0] != "public.png" || dialog.contentTypes[1] != "com.adobe.pdf" {
		t.Errorf("content types should be trimmed and kept in order, got %v", dialog.contentTypes)
	}

	if dialogContentTypesNative {
		if len(dialog.filters) != 0 {
			t.Errorf("native content types must not add extension filters, got %v", dialog.filters)
		}
	} else {
		if len(dialog.filters) != 2 {
			t.Fatalf("content types should be translated to filters, got %v", dialog.filters)
		}
		if dialog.filters[0].Pattern != "*.png" || dialog.filters[1].Pattern != "*.pdf" {
			t.Errorf("unexpected translated patterns %v", dialog.filters)
		}
	}
}

func TestOpenFileDialog_AddContentTypeCoexistsWithFilters(t *testing.T) {
	dialog := newOpenFileDialog()
	defer freeDialogID(dialog.id)

	dialog.AddFilter("Text", "*.txt").AddContentType("public.image")
	if len(dialog.filters) == 0 || dialog.filters[0].Pattern != "*.txt" {
		t.Errorf("AddFilter entries must survive AddContentType, got %v", dialog.filters)
	}
	if len(dialog.contentTypes) != 1 || dialog.contentTypes[0] != "public.image" {
		t.Errorf("unexpected content types %v", dialog.contentTypes)
	}
}

func TestSaveFileDialog_AddContentType(t *testing.T) {
	dialog := newSaveFileDialog()
	defer freeDialogID(dialog.id)

	dialog.AddContentType("public.json")
	if len(dialog.contentTypes) != 1 || dialog.contentTypes[0] != "public.json" {
		t.Errorf("unexpected content types %v", dialog.contentTypes)
	}
	if !dialogContentTypesNative && (len(dialog.filters) != 1 || dialog.filters[0].Pattern != "*.json") {
		t.Errorf("expected a translated json filter, got %v", dialog.filters)
	}
}

func TestSaveFileDialog_SetFormats(t *testing.T) {
	dialog := newSaveFileDialog()
	defer freeDialogID(dialog.id)

	if dialog.SelectedFormat() != -1 {
		t.Errorf("SelectedFormat should be -1 before SetFormats, got %d", dialog.SelectedFormat())
	}

	formats := []DialogFormat{
		{Label: "PNG image", Extension: "png", UTI: "public.png"},
		{Label: "PDF document", Extension: "pdf", UTI: "com.adobe.pdf"},
		{Label: "Plain text", Extension: "txt"},
	}
	var changes []int
	result := dialog.SetFormats(formats, 1, func(index int) {
		changes = append(changes, index)
	})
	if result != dialog {
		t.Error("SetFormats should return the dialog for chaining")
	}
	if dialog.SelectedFormat() != 1 {
		t.Errorf("SelectedFormat should report the initial selection, got %d", dialog.SelectedFormat())
	}
	if len(dialog.formats) != 3 {
		t.Fatalf("expected 3 formats, got %d", len(dialog.formats))
	}

	// The dialog keeps its own copy of the slice.
	formats[0].Label = "changed"
	if dialog.formats[0].Label != "PNG image" {
		t.Error("SetFormats should copy the formats slice")
	}

	// Out of range initial selections fall back to the first entry.
	dialog.SetFormats(formats, 7, nil)
	if dialog.SelectedFormat() != 0 {
		t.Errorf("out of range selection should clamp to 0, got %d", dialog.SelectedFormat())
	}
	dialog.SetFormats(formats, -3, nil)
	if dialog.SelectedFormat() != 0 {
		t.Errorf("negative selection should clamp to 0, got %d", dialog.SelectedFormat())
	}

	if !dialogContentTypesNative {
		if len(dialog.filters) == 0 {
			t.Error("formats should become extension filters off macOS")
		}
	} else if len(dialog.filters) != 0 {
		t.Errorf("formats must not add extension filters on macOS, got %v", dialog.filters)
	}
}

func TestSaveFileDialog_SetSelectedFormat(t *testing.T) {
	dialog := newSaveFileDialog()
	defer freeDialogID(dialog.id)

	changed := make(chan int, 4)
	dialog.SetFormats([]DialogFormat{{Label: "A", Extension: "a"}, {Label: "B", Extension: "b"}}, 0, func(index int) {
		changed <- index
	})

	dialog.setSelectedFormat(1)
	if dialog.SelectedFormat() != 1 {
		t.Errorf("setSelectedFormat should update the selection, got %d", dialog.SelectedFormat())
	}
	if got := <-changed; got != 1 {
		t.Errorf("onChange should receive the new index, got %d", got)
	}

	// Out of range indexes are ignored and do not fire the callback.
	dialog.setSelectedFormat(5)
	dialog.setSelectedFormat(-1)
	if dialog.SelectedFormat() != 1 {
		t.Errorf("out of range indexes must be ignored, got %d", dialog.SelectedFormat())
	}
	select {
	case got := <-changed:
		t.Errorf("onChange must not fire for out of range index, got %d", got)
	default:
	}
}

func TestSaveFileDialog_NameFieldLabelAndTags(t *testing.T) {
	dialog := newSaveFileDialog()
	defer freeDialogID(dialog.id)

	dialog.SetNameFieldLabel("Export As:").SetTags([]string{" Work ", "", "Draft"})
	if dialog.nameFieldLabel != "Export As:" {
		t.Errorf("unexpected name field label %q", dialog.nameFieldLabel)
	}
	if len(dialog.tags) != 2 || dialog.tags[0] != "Work" || dialog.tags[1] != "Draft" {
		t.Errorf("tags should be trimmed with empties dropped, got %v", dialog.tags)
	}
	dialog.SetTags(nil)
	if len(dialog.tags) != 0 {
		t.Errorf("SetTags(nil) should clear the tags, got %v", dialog.tags)
	}
}

func TestContentTypeFilter(t *testing.T) {
	tests := []struct {
		uti     string
		pattern string
		ok      bool
	}{
		{"public.png", "*.png", true},
		{"PUBLIC.JPEG", "*.jpg;*.jpeg", true},
		{" com.adobe.pdf ", "*.pdf", true},
		{"public.item", "*", true},
		{"public.data", "*", true},
		{"com.example.custom-thing", "", false},
		{"", "", false},
	}
	for _, tt := range tests {
		filter, ok := contentTypeFilter(tt.uti)
		if ok != tt.ok {
			t.Errorf("contentTypeFilter(%q) ok = %v, want %v", tt.uti, ok, tt.ok)
			continue
		}
		if !ok {
			continue
		}
		if filter.Pattern != tt.pattern {
			t.Errorf("contentTypeFilter(%q) pattern = %q, want %q", tt.uti, filter.Pattern, tt.pattern)
		}
		if filter.DisplayName == "" {
			t.Errorf("contentTypeFilter(%q) should carry a display name", tt.uti)
		}
	}
}

func TestContentTypeExtensionsTable(t *testing.T) {
	for uti, extensions := range contentTypeExtensions {
		if len(extensions) == 0 {
			t.Errorf("%s has no extensions", uti)
		}
		if uti != strings.ToLower(uti) {
			t.Errorf("%s should be lower case so lookups are case insensitive", uti)
		}
		for _, ext := range extensions {
			if strings.HasPrefix(ext, ".") || strings.Contains(ext, "*") || ext == "" {
				t.Errorf("%s: extension %q must be bare (no dot or wildcard)", uti, ext)
			}
		}
	}
}

func TestDialogOptionDefaults(t *testing.T) {
	var prompt PromptOptions
	if prompt.Secure || prompt.Window != nil || prompt.OKLabel != "" || prompt.CancelLabel != "" {
		t.Error("PromptOptions zero value should be an insecure, unattached prompt with system button labels")
	}

	var color ColorPickerOptions
	if color.ShowsAlpha || color.OnChange != nil || color.Initial != (RGBA{}) {
		t.Error("ColorPickerOptions zero value should hide alpha and have no live callback")
	}

	var font FontPickerOptions
	if font.Family != "" || font.Size != 0 || font.OnChange != nil {
		t.Error("FontPickerOptions zero value should select the system font")
	}

	var descriptor FontDescriptor
	if descriptor != (FontDescriptor{}) {
		t.Error("FontDescriptor zero value should be empty")
	}

	if ErrDialogNotSupported == nil || ErrDialogInProgress == nil || ErrDialogNotSupported == ErrDialogInProgress {
		t.Error("dialog errors should be distinct sentinel values")
	}
}
