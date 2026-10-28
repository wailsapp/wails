//go:build windows && !server && !cef

package application

import "testing"

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

// The public EnableSizeConstraints may already have restored the shared copy
// before the platform unfullscreen runs.
func TestWindowsRestoringUnsuspendedConstraintsLeavesOptionsAlone(t *testing.T) {
	w := newConstraintTestWindow()

	w.enableSizeConstraints()

	options := w.parent.options
	if options.MinWidth != 200 || options.MinHeight != 100 || options.MaxWidth != 800 || options.MaxHeight != 600 {
		t.Errorf("options changed to min %dx%d max %dx%d, want 200x100 and 800x600",
			options.MinWidth, options.MinHeight, options.MaxWidth, options.MaxHeight)
	}
}

// Stub native fullscreen transitions while retaining the Windows constraint code.
type constraintTestWindow struct {
	*windowsWebviewWindow
}

func (w *constraintTestWindow) fullscreen() {
	w.disableSizeConstraints()
	w.isCurrentlyFullscreen = true
}

func (w *constraintTestWindow) unfullscreen() {
	w.isCurrentlyFullscreen = false
	w.enableSizeConstraints()
}

func (w *constraintTestWindow) size() (int, int) { return 640, 480 }
func (w *constraintTestWindow) setSize(int, int) {}

func TestWindowsFullscreenConstraintRoundTrip(t *testing.T) {
	for _, startup := range []bool{false, true} {
		entry := "public"
		if startup {
			entry = "startup"
		}
		t.Run(entry, func(t *testing.T) {
			for _, tt := range []struct {
				name   string
				update func(*WebviewWindow)
				want   [4]int
			}{
				{"unchanged", func(*WebviewWindow) {}, [4]int{200, 100, 800, 600}},
				{"minimum", func(w *WebviewWindow) { w.SetMinSize(300, 250) }, [4]int{300, 250, 800, 600}},
				{"maximum", func(w *WebviewWindow) { w.SetMaxSize(1200, 900) }, [4]int{200, 100, 1200, 900}},
				{"both", func(w *WebviewWindow) {
					w.SetMinSize(300, 250)
					w.SetMaxSize(1200, 900)
				}, [4]int{300, 250, 1200, 900}},
				{"clear minimum", func(w *WebviewWindow) { w.SetMinSize(0, 0) }, [4]int{0, 0, 800, 600}},
				{"clear maximum", func(w *WebviewWindow) { w.SetMaxSize(0, 0) }, [4]int{200, 100, 0, 0}},
				{"clear both", func(w *WebviewWindow) {
					w.SetMinSize(0, 0)
					w.SetMaxSize(0, 0)
				}, [4]int{}},
			} {
				t.Run(tt.name, func(t *testing.T) {
					previous := globalApplication
					dispatcher := &threadProbeApp{}
					dispatcher.onMain.Store(true)
					globalApplication = &App{impl: dispatcher}
					t.Cleanup(func() { globalApplication = previous })

					platform := &constraintTestWindow{newConstraintTestWindow()}
					w := platform.parent
					w.impl = platform
					if startup {
						platform.fullscreen()
					} else {
						w.Fullscreen()
					}
					w.Fullscreen() // Re-entering must preserve the saved copy.
					tt.update(w)
					w.UnFullscreen()

					for cycle := range 2 {
						options := w.options
						got := [4]int{options.MinWidth, options.MinHeight, options.MaxWidth, options.MaxHeight}
						if got != tt.want {
							t.Fatalf("cycle %d: constraints = %v, want %v", cycle, got, tt.want)
						}
						if w.constraintsSaved {
							t.Fatal("constraints still saved after leaving fullscreen")
						}
						if cycle == 0 {
							w.Fullscreen()
							w.UnFullscreen()
						}
					}
				})
			}
		})
	}
}
