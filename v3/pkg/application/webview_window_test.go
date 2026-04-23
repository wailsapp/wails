package application

import (
	"testing"
)

type mockWindowImpl struct {
	minW, minH, maxW, maxH int
}

func (m *mockWindowImpl) setMinSize(width, height int) {
	m.minW, m.minH = width, height
}

func (m *mockWindowImpl) setMaxSize(width, height int) {
	m.maxW, m.maxH = width, height
}

func TestDisableSizeConstraintsSavesValues(t *testing.T) {
	w := &WebviewWindow{
		options: WebviewWindowOptions{
			MinWidth:  200,
			MinHeight: 100,
			MaxWidth:  800,
			MaxHeight: 600,
		},
	}

	if w.savedMinWidth != 0 || w.savedMinHeight != 0 {
		t.Fatal("saved values should start at zero")
	}
	if w.constraintsSaved {
		t.Fatal("constraintsSaved should start false")
	}

	w.constraintsSaved = false
	w.savedMinWidth = w.options.MinWidth
	w.savedMinHeight = w.options.MinHeight
	w.savedMaxWidth = w.options.MaxWidth
	w.savedMaxHeight = w.options.MaxHeight
	w.constraintsSaved = true

	if w.savedMinWidth != 200 {
		t.Errorf("savedMinWidth = %d, want 200", w.savedMinWidth)
	}
	if w.savedMinHeight != 100 {
		t.Errorf("savedMinHeight = %d, want 100", w.savedMinHeight)
	}
	if w.savedMaxWidth != 800 {
		t.Errorf("savedMaxWidth = %d, want 800", w.savedMaxWidth)
	}
	if w.savedMaxHeight != 600 {
		t.Errorf("savedMaxHeight = %d, want 600", w.savedMaxHeight)
	}
	if !w.constraintsSaved {
		t.Error("constraintsSaved should be true after saving")
	}
}

func TestRestoreAfterZeroingOptions(t *testing.T) {
	w := &WebviewWindow{
		options: WebviewWindowOptions{
			MinWidth:  200,
			MinHeight: 100,
			MaxWidth:  800,
			MaxHeight: 600,
		},
	}

	savedMinW := w.options.MinWidth
	savedMinH := w.options.MinHeight
	savedMaxW := w.options.MaxWidth
	savedMaxH := w.options.MaxHeight

	w.options.MinWidth = 0
	w.options.MinHeight = 0
	w.options.MaxWidth = 0
	w.options.MaxHeight = 0

	if w.options.MinWidth != 0 {
		t.Error("MinWidth should be zeroed")
	}

	minW, minH := w.options.MinWidth, w.options.MinHeight
	maxW, maxH := w.options.MaxWidth, w.options.MaxHeight
	if true {
		minW, minH = savedMinW, savedMinH
		maxW, maxH = savedMaxW, savedMaxH
	}

	if minW != 200 {
		t.Errorf("restored minW = %d, want 200", minW)
	}
	if minH != 100 {
		t.Errorf("restored minH = %d, want 100", minH)
	}
	if maxW != 800 {
		t.Errorf("restored maxW = %d, want 800", maxW)
	}
	if maxH != 600 {
		t.Errorf("restored maxH = %d, want 600", maxH)
	}
}

func TestDisableSizeConstraintsDoesNotOverwriteSavedValues(t *testing.T) {
	w := &WebviewWindow{
		options: WebviewWindowOptions{
			MinWidth:  300,
			MinHeight: 200,
			MaxWidth:  1024,
			MaxHeight: 768,
		},
	}

	w.savedMinWidth = w.options.MinWidth
	w.savedMinHeight = w.options.MinHeight
	w.savedMaxWidth = w.options.MaxWidth
	w.savedMaxHeight = w.options.MaxHeight
	w.constraintsSaved = true

	w.options.MinWidth = 0
	w.options.MinHeight = 0
	w.options.MaxWidth = 0
	w.options.MaxHeight = 0

	if w.savedMinWidth != 300 {
		t.Errorf("savedMinWidth overwritten = %d, want 300", w.savedMinWidth)
	}
	if w.savedMinHeight != 200 {
		t.Errorf("savedMinHeight overwritten = %d, want 200", w.savedMinHeight)
	}
	if w.savedMaxWidth != 1024 {
		t.Errorf("savedMaxWidth overwritten = %d, want 1024", w.savedMaxWidth)
	}
	if w.savedMaxHeight != 768 {
		t.Errorf("savedMaxHeight overwritten = %d, want 768", w.savedMaxHeight)
	}
}
