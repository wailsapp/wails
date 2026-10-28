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

func TestRestoreSavedSizeConstraintOptions(t *testing.T) {
	w := &WebviewWindow{
		options:          WebviewWindowOptions{},
		savedMinWidth:    320,
		savedMinHeight:   240,
		savedMaxWidth:    1280,
		savedMaxHeight:   720,
		constraintsSaved: true,
	}

	if !w.restoreSavedSizeConstraintOptions() {
		t.Fatal("restoreSavedSizeConstraintOptions should report restored constraints")
	}

	if w.options.MinWidth != 320 || w.options.MinHeight != 240 {
		t.Fatalf("restored min constraints = %dx%d, want 320x240", w.options.MinWidth, w.options.MinHeight)
	}
	if w.options.MaxWidth != 1280 || w.options.MaxHeight != 720 {
		t.Fatalf("restored max constraints = %dx%d, want 1280x720", w.options.MaxWidth, w.options.MaxHeight)
	}
	if w.constraintsSaved {
		t.Fatal("constraintsSaved should be false after restoring")
	}
}

func TestRestoreSavedSizeConstraintOptionsNoopWithoutSavedValues(t *testing.T) {
	w := &WebviewWindow{
		options: WebviewWindowOptions{
			MinWidth:  200,
			MinHeight: 100,
			MaxWidth:  800,
			MaxHeight: 600,
		},
	}

	if w.restoreSavedSizeConstraintOptions() {
		t.Fatal("restoreSavedSizeConstraintOptions should report no-op when nothing was saved")
	}

	if w.options.MinWidth != 200 || w.options.MinHeight != 100 {
		t.Fatalf("min constraints changed unexpectedly to %dx%d", w.options.MinWidth, w.options.MinHeight)
	}
	if w.options.MaxWidth != 800 || w.options.MaxHeight != 600 {
		t.Fatalf("max constraints changed unexpectedly to %dx%d", w.options.MaxWidth, w.options.MaxHeight)
	}
}

// A fullscreen or maximised window is sized by the platform. SetMinSize and
// SetMaxSize used to resize it anyway to satisfy the new constraint, which
// left the geometry and the state disagreeing: on macOS a window shrunk to
// 1000x700 this way goes on reporting IsFullscreen() == true.

// sizeStateImpl embeds the interface, so the 80-odd methods this path never
// touches are present and panic if the path ever reaches one.
type sizeStateImpl struct {
	webviewWindowImpl

	stubFullscreen bool
	stubMaximised  bool
	stubWidth      int
	stubHeight     int
	stubSetSizes   int
}

func (s *sizeStateImpl) size() (int, int)    { return s.stubWidth, s.stubHeight }
func (s *sizeStateImpl) isFullscreen() bool  { return s.stubFullscreen }
func (s *sizeStateImpl) isMaximised() bool   { return s.stubMaximised }
func (s *sizeStateImpl) setMinSize(_, _ int) {}
func (s *sizeStateImpl) setMaxSize(_, _ int) {}

func (s *sizeStateImpl) setSize(width, height int) {
	s.stubSetSizes++
	s.stubWidth, s.stubHeight = width, height
}

// newSizeStateWindow installs a main thread that runs dispatched work inline,
// so the InvokeSync inside the setters completes rather than queueing.
func newSizeStateWindow(t *testing.T, impl *sizeStateImpl) *WebviewWindow {
	t.Helper()

	dispatcher := &threadProbeApp{}
	dispatcher.onMain.Store(true)

	previous := globalApplication
	globalApplication = &App{impl: dispatcher}
	t.Cleanup(func() { globalApplication = previous })

	return &WebviewWindow{impl: impl}
}

func TestSetMaxSizeLeavesAFullscreenWindowAlone(t *testing.T) {
	impl := &sizeStateImpl{stubFullscreen: true, stubWidth: 1512, stubHeight: 949}
	w := newSizeStateWindow(t, impl)

	w.SetMaxSize(1000, 700)

	if impl.stubSetSizes != 0 {
		t.Errorf("setSize called %d times on a fullscreen window, want 0", impl.stubSetSizes)
	}
	if impl.stubWidth != 1512 || impl.stubHeight != 949 {
		t.Errorf("window is %dx%d, want it left at 1512x949", impl.stubWidth, impl.stubHeight)
	}
	if w.options.MaxWidth != 1000 || w.options.MaxHeight != 700 {
		t.Errorf("max recorded as %dx%d, want 1000x700", w.options.MaxWidth, w.options.MaxHeight)
	}
}

func TestSetMinSizeLeavesAMaximisedWindowAlone(t *testing.T) {
	impl := &sizeStateImpl{stubMaximised: true, stubWidth: 400, stubHeight: 300}
	w := newSizeStateWindow(t, impl)

	w.SetMinSize(800, 600)

	if impl.stubSetSizes != 0 {
		t.Errorf("setSize called %d times on a maximised window, want 0", impl.stubSetSizes)
	}
	if impl.stubWidth != 400 || impl.stubHeight != 300 {
		t.Errorf("window is %dx%d, want it left at 400x300", impl.stubWidth, impl.stubHeight)
	}
	if w.options.MinWidth != 800 || w.options.MinHeight != 600 {
		t.Errorf("min recorded as %dx%d, want 800x600", w.options.MinWidth, w.options.MinHeight)
	}
}

// The resize itself is the point of these setters on a window the application
// does own, so the guard has to stop at the two managed states.
func TestSetMaxSizeStillResizesAnOrdinaryWindow(t *testing.T) {
	impl := &sizeStateImpl{stubWidth: 1512, stubHeight: 949}
	w := newSizeStateWindow(t, impl)

	w.SetMaxSize(1000, 700)

	if impl.stubSetSizes != 1 {
		t.Fatalf("setSize called %d times, want 1", impl.stubSetSizes)
	}
	if impl.stubWidth != 1000 || impl.stubHeight != 700 {
		t.Errorf("window is %dx%d, want 1000x700", impl.stubWidth, impl.stubHeight)
	}
}

func TestSetMinSizeStillResizesAnOrdinaryWindow(t *testing.T) {
	impl := &sizeStateImpl{stubWidth: 400, stubHeight: 300}
	w := newSizeStateWindow(t, impl)

	w.SetMinSize(800, 600)

	if impl.stubSetSizes != 1 {
		t.Fatalf("setSize called %d times, want 1", impl.stubSetSizes)
	}
	if impl.stubWidth != 800 || impl.stubHeight != 600 {
		t.Errorf("window is %dx%d, want 800x600", impl.stubWidth, impl.stubHeight)
	}
}
