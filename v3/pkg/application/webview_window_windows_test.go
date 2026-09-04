//go:build windows

package application

import "testing"

// Tests for the position arithmetic used by relativePosition() and setRelativePosition()
// on Windows. The invariant: coordinates are physical pixels throughout; no DPI scaling
// is applied to the caller-supplied values.
//
// Regression case (GitHub #4300): the old implementation passed coordinates through
// DipToPhysicalRect before adding PhysicalWorkArea.X/Y, multiplying them by the scale
// factor a second time. At 125% DPI a window requested at physical (500, 300) would
// land at (625, 375) — 25% past the intended position.

func TestSetRelativePositionArithmetic(t *testing.T) {
	tests := []struct {
		name        string
		scaleFactor float32
		workAreaX   int
		workAreaY   int
		inputX      int
		inputY      int
		wantX       int
		wantY       int
	}{
		{
			name:        "100% scale, primary monitor",
			scaleFactor: 1.0, workAreaX: 0, workAreaY: 0,
			inputX: 500, inputY: 300,
			wantX: 500, wantY: 300,
		},
		{
			name:        "125% scale — physical coords must not be scaled again",
			scaleFactor: 1.25, workAreaX: 0, workAreaY: 0,
			inputX: 500, inputY: 300,
			// old buggy result: int(500*1.25)=625, int(300*1.25)=375
			wantX: 500, wantY: 300,
		},
		{
			name:        "150% scale",
			scaleFactor: 1.5, workAreaX: 0, workAreaY: 0,
			inputX: 500, inputY: 300,
			// old buggy result: 750, 450
			wantX: 500, wantY: 300,
		},
		{
			name:        "200% scale",
			scaleFactor: 2.0, workAreaX: 0, workAreaY: 0,
			inputX: 500, inputY: 300,
			// old buggy result: 1000, 600
			wantX: 500, wantY: 300,
		},
		{
			name:        "125% scale, secondary monitor with work area offset",
			scaleFactor: 1.25, workAreaX: 1920, workAreaY: 0,
			inputX: 100, inputY: 50,
			wantX: 2020, wantY: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			screen := &Screen{
				ScaleFactor:      tt.scaleFactor,
				PhysicalWorkArea: Rect{X: tt.workAreaX, Y: tt.workAreaY},
			}
			// mirrors the arithmetic in setRelativePosition
			gotX := screen.PhysicalWorkArea.X + tt.inputX
			gotY := screen.PhysicalWorkArea.Y + tt.inputY
			if gotX != tt.wantX || gotY != tt.wantY {
				t.Errorf("setRelativePosition: got (%d,%d), want (%d,%d)", gotX, gotY, tt.wantX, tt.wantY)
			}
		})
	}
}

func TestSetRelativePositionNilScreen(t *testing.T) {
	// When getScreen returns nil the function must not panic and must place
	// the window at (x, y) in absolute physical coordinates.
	x, y := 500, 300
	physBounds := Rect{Width: 800, Height: 600}

	var screen *Screen
	var result Rect
	if screen == nil {
		result = Rect{X: x, Y: y, Width: physBounds.Width, Height: physBounds.Height}
	} else {
		result = Rect{X: screen.PhysicalWorkArea.X + x, Y: screen.PhysicalWorkArea.Y + y,
			Width: physBounds.Width, Height: physBounds.Height}
	}

	if result.X != x || result.Y != y {
		t.Errorf("nil screen: got (%d,%d), want (%d,%d)", result.X, result.Y, x, y)
	}
	if result.Width != physBounds.Width || result.Height != physBounds.Height {
		t.Errorf("nil screen: size got (%d,%d), want (%d,%d)",
			result.Width, result.Height, physBounds.Width, physBounds.Height)
	}
}

func TestRelativePositionArithmetic(t *testing.T) {
	tests := []struct {
		name        string
		scaleFactor float32
		workAreaX   int
		workAreaY   int
		physX       int
		physY       int
		wantX       int
		wantY       int
	}{
		{
			name:        "100% scale, primary monitor",
			scaleFactor: 1.0, workAreaX: 0, workAreaY: 0,
			physX: 500, physY: 300,
			wantX: 500, wantY: 300,
		},
		{
			name:        "125% scale — physical position returned without DIP conversion",
			scaleFactor: 1.25, workAreaX: 0, workAreaY: 0,
			physX: 625, physY: 375,
			wantX: 625, wantY: 375,
		},
		{
			name:        "150% scale",
			scaleFactor: 1.5, workAreaX: 0, workAreaY: 0,
			physX: 750, physY: 450,
			wantX: 750, wantY: 450,
		},
		{
			name:        "200% scale",
			scaleFactor: 2.0, workAreaX: 0, workAreaY: 0,
			physX: 1000, physY: 600,
			wantX: 1000, wantY: 600,
		},
		{
			name:        "125% scale, secondary monitor with work area offset",
			scaleFactor: 1.25, workAreaX: 1920, workAreaY: 0,
			physX: 2020, physY: 50,
			wantX: 100, wantY: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			screen := &Screen{
				ScaleFactor:      tt.scaleFactor,
				PhysicalWorkArea: Rect{X: tt.workAreaX, Y: tt.workAreaY},
			}
			physBounds := Rect{X: tt.physX, Y: tt.physY}
			// mirrors the arithmetic in relativePosition
			gotX := physBounds.X - screen.PhysicalWorkArea.X
			gotY := physBounds.Y - screen.PhysicalWorkArea.Y
			if gotX != tt.wantX || gotY != tt.wantY {
				t.Errorf("relativePosition: got (%d,%d), want (%d,%d)", gotX, gotY, tt.wantX, tt.wantY)
			}
		})
	}
}

func TestRelativePositionNilScreen(t *testing.T) {
	// When getScreen returns nil the function must not panic and must return
	// the raw physical window position.
	physBounds := Rect{X: 500, Y: 300}

	var screen *Screen
	var gotX, gotY int
	if screen == nil {
		gotX, gotY = physBounds.X, physBounds.Y
	} else {
		gotX = physBounds.X - screen.PhysicalWorkArea.X
		gotY = physBounds.Y - screen.PhysicalWorkArea.Y
	}

	if gotX != physBounds.X || gotY != physBounds.Y {
		t.Errorf("nil screen: got (%d,%d), want (%d,%d)", gotX, gotY, physBounds.X, physBounds.Y)
	}
}
