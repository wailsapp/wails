//go:build windows && !server && !cef

package application

import "testing"

// A window that starts fullscreen reaches the platform fullscreen() directly
// from the StartState switch in run(), so WebviewWindow.DisableSizeConstraints
// - which keeps its own saved copy - never runs. On this platform
// setMinSize/setMaxSize write nothing but parent.options, the clamp being read
// back from there in WM_GETMINMAXINFO, so suspending it by zeroing those
// options used to destroy the only copy and the window never got its limits
// back.

func newConstraintTestWindow() *windowsWebviewWindow {
	return &windowsWebviewWindow{
		parent: &WebviewWindow{
			options: WebviewWindowOptions{
				MinWidth:  200,
				MinHeight: 100,
				MaxWidth:  800,
				MaxHeight: 600,
			},
		},
	}
}

func TestWindowsSizeConstraintsSurviveSuspension(t *testing.T) {
	w := newConstraintTestWindow()

	w.disableSizeConstraints()

	options := w.parent.options
	if options.MinWidth != 0 || options.MinHeight != 0 || options.MaxWidth != 0 || options.MaxHeight != 0 {
		t.Fatalf("constraints not suspended: min %dx%d max %dx%d",
			options.MinWidth, options.MinHeight, options.MaxWidth, options.MaxHeight)
	}

	w.enableSizeConstraints()

	options = w.parent.options
	if options.MinWidth != 200 || options.MinHeight != 100 {
		t.Errorf("min restored as %dx%d, want 200x100", options.MinWidth, options.MinHeight)
	}
	if options.MaxWidth != 800 || options.MaxHeight != 600 {
		t.Errorf("max restored as %dx%d, want 800x600", options.MaxWidth, options.MaxHeight)
	}
}

// Suspending twice must not overwrite the saved copy with the zeroes the first
// suspension wrote. fullscreen() calls disableSizeConstraints unconditionally,
// and the public Fullscreen() has already called the window-level one, so a
// second call on an already-suspended window is the normal path, not an edge
// case.
func TestWindowsSuspendingConstraintsTwiceKeepsTheFirstCopy(t *testing.T) {
	w := newConstraintTestWindow()

	w.disableSizeConstraints()
	w.disableSizeConstraints()
	w.enableSizeConstraints()

	options := w.parent.options
	if options.MinWidth != 200 || options.MaxWidth != 800 {
		t.Errorf("after two suspensions min=%d max=%d, want 200 and 800", options.MinWidth, options.MaxWidth)
	}
}

// With nothing suspended, restoring still has to read the options - that is
// the path taken when the window-level EnableSizeConstraints has already put
// them back before unfullscreen() runs.
func TestWindowsRestoringUnsuspendedConstraintsLeavesOptionsAlone(t *testing.T) {
	w := newConstraintTestWindow()

	w.enableSizeConstraints()

	options := w.parent.options
	if options.MinWidth != 200 || options.MinHeight != 100 || options.MaxWidth != 800 || options.MaxHeight != 600 {
		t.Errorf("options changed to min %dx%d max %dx%d, want 200x100 and 800x600",
			options.MinWidth, options.MinHeight, options.MaxWidth, options.MaxHeight)
	}
}
