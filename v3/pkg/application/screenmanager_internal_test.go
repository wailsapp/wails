package application

import (
	"testing"
)

func TestAlignment_Constants(t *testing.T) {
	if TOP != 0 {
		t.Error("TOP should be 0")
	}
	if RIGHT != 1 {
		t.Error("RIGHT should be 1")
	}
	if BOTTOM != 2 {
		t.Error("BOTTOM should be 2")
	}
	if LEFT != 3 {
		t.Error("LEFT should be 3")
	}
}

func TestOffsetReference_Constants(t *testing.T) {
	if BEGIN != 0 {
		t.Error("BEGIN should be 0")
	}
	if END != 1 {
		t.Error("END should be 1")
	}
}

func TestRect_Origin(t *testing.T) {
	rect := Rect{X: 10, Y: 20, Width: 100, Height: 200}
	origin := rect.Origin()

	if origin.X != 10 {
		t.Errorf("origin.X = %d, want 10", origin.X)
	}
	if origin.Y != 20 {
		t.Errorf("origin.Y = %d, want 20", origin.Y)
	}
}

func TestRect_Corner(t *testing.T) {
	rect := Rect{X: 10, Y: 20, Width: 100, Height: 200}
	corner := rect.Corner()

	if corner.X != 110 { // 10 + 100
		t.Errorf("corner.X = %d, want 110", corner.X)
	}
	if corner.Y != 220 { // 20 + 200
		t.Errorf("corner.Y = %d, want 220", corner.Y)
	}
}

func TestRect_InsideCorner(t *testing.T) {
	rect := Rect{X: 10, Y: 20, Width: 100, Height: 200}
	inside := rect.InsideCorner()

	if inside.X != 109 { // 10 + 100 - 1
		t.Errorf("inside.X = %d, want 109", inside.X)
	}
	if inside.Y != 219 { // 20 + 200 - 1
		t.Errorf("inside.Y = %d, want 219", inside.Y)
	}
}

func TestRect_right(t *testing.T) {
	rect := Rect{X: 10, Y: 20, Width: 100, Height: 200}
	if rect.right() != 110 {
		t.Errorf("right() = %d, want 110", rect.right())
	}
}

func TestRect_bottom(t *testing.T) {
	rect := Rect{X: 10, Y: 20, Width: 100, Height: 200}
	if rect.bottom() != 220 {
		t.Errorf("bottom() = %d, want 220", rect.bottom())
	}
}

func TestRect_Size(t *testing.T) {
	rect := Rect{X: 10, Y: 20, Width: 100, Height: 200}
	size := rect.Size()

	if size.Width != 100 {
		t.Errorf("Width = %d, want 100", size.Width)
	}
	if size.Height != 200 {
		t.Errorf("Height = %d, want 200", size.Height)
	}
}

func TestRect_IsEmpty(t *testing.T) {
	tests := []struct {
		rect     Rect
		expected bool
	}{
		{Rect{X: 0, Y: 0, Width: 0, Height: 0}, true},
		{Rect{X: 0, Y: 0, Width: 100, Height: 0}, true},
		{Rect{X: 0, Y: 0, Width: 0, Height: 100}, true},
		{Rect{X: 0, Y: 0, Width: -1, Height: 100}, true},
		{Rect{X: 0, Y: 0, Width: 100, Height: -1}, true},
		{Rect{X: 0, Y: 0, Width: 100, Height: 200}, false},
		{Rect{X: 10, Y: 20, Width: 1, Height: 1}, false},
	}

	for _, tt := range tests {
		result := tt.rect.IsEmpty()
		if result != tt.expected {
			t.Errorf("Rect%v.IsEmpty() = %v, want %v", tt.rect, result, tt.expected)
		}
	}
}

func TestRect_Contains(t *testing.T) {
	rect := Rect{X: 10, Y: 20, Width: 100, Height: 200}

	tests := []struct {
		point    Point
		expected bool
	}{
		{Point{X: 10, Y: 20}, true},   // top-left corner
		{Point{X: 50, Y: 100}, true},  // inside
		{Point{X: 109, Y: 219}, true}, // inside corner
		{Point{X: 110, Y: 220}, false}, // corner (exclusive)
		{Point{X: 0, Y: 0}, false},    // outside
		{Point{X: 9, Y: 20}, false},   // left of rect
		{Point{X: 10, Y: 19}, false},  // above rect
		{Point{X: 111, Y: 100}, false}, // right of rect
		{Point{X: 50, Y: 221}, false},  // below rect
	}

	for _, tt := range tests {
		result := rect.Contains(tt.point)
		if result != tt.expected {
			t.Errorf("Rect%v.Contains(%v) = %v, want %v", rect, tt.point, result, tt.expected)
		}
	}
}

func TestRect_Intersect(t *testing.T) {
	tests := []struct {
		name     string
		r1       Rect
		r2       Rect
		expected Rect
	}{
		{
			name:     "overlapping",
			r1:       Rect{X: 0, Y: 0, Width: 100, Height: 100},
			r2:       Rect{X: 50, Y: 50, Width: 100, Height: 100},
			expected: Rect{X: 50, Y: 50, Width: 50, Height: 50},
		},
		{
			name:     "no overlap - horizontal",
			r1:       Rect{X: 0, Y: 0, Width: 100, Height: 100},
			r2:       Rect{X: 200, Y: 0, Width: 100, Height: 100},
			expected: Rect{},
		},
		{
			name:     "no overlap - vertical",
			r1:       Rect{X: 0, Y: 0, Width: 100, Height: 100},
			r2:       Rect{X: 0, Y: 200, Width: 100, Height: 100},
			expected: Rect{},
		},
		{
			name:     "contained",
			r1:       Rect{X: 0, Y: 0, Width: 100, Height: 100},
			r2:       Rect{X: 25, Y: 25, Width: 50, Height: 50},
			expected: Rect{X: 25, Y: 25, Width: 50, Height: 50},
		},
		{
			name:     "identical",
			r1:       Rect{X: 10, Y: 20, Width: 100, Height: 100},
			r2:       Rect{X: 10, Y: 20, Width: 100, Height: 100},
			expected: Rect{X: 10, Y: 20, Width: 100, Height: 100},
		},
		{
			name:     "empty rect 1",
			r1:       Rect{},
			r2:       Rect{X: 0, Y: 0, Width: 100, Height: 100},
			expected: Rect{},
		},
		{
			name:     "empty rect 2",
			r1:       Rect{X: 0, Y: 0, Width: 100, Height: 100},
			r2:       Rect{},
			expected: Rect{},
		},
		{
			name:     "touching edges - no intersection",
			r1:       Rect{X: 0, Y: 0, Width: 100, Height: 100},
			r2:       Rect{X: 100, Y: 0, Width: 100, Height: 100},
			expected: Rect{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.r1.Intersect(tt.r2)
			if result != tt.expected {
				t.Errorf("Intersect: got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRect_distanceFromRectSquared(t *testing.T) {
	tests := []struct {
		name     string
		r1       Rect
		r2       Rect
		expected int
	}{
		{
			name:     "overlapping - negative area",
			r1:       Rect{X: 0, Y: 0, Width: 100, Height: 100},
			r2:       Rect{X: 50, Y: 50, Width: 100, Height: 100},
			expected: -(50 * 50), // intersection area
		},
		{
			name:     "horizontal gap",
			r1:       Rect{X: 0, Y: 0, Width: 100, Height: 100},
			r2:       Rect{X: 110, Y: 0, Width: 100, Height: 100},
			expected: 100, // gap of 10, squared
		},
		{
			name:     "vertical gap",
			r1:       Rect{X: 0, Y: 0, Width: 100, Height: 100},
			r2:       Rect{X: 0, Y: 120, Width: 100, Height: 100},
			expected: 400, // gap of 20, squared
		},
		{
			name:     "diagonal gap",
			r1:       Rect{X: 0, Y: 0, Width: 100, Height: 100},
			r2:       Rect{X: 110, Y: 110, Width: 100, Height: 100},
			expected: 200, // dX=10, dY=10, 10^2 + 10^2 = 200
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.r1.distanceFromRectSquared(tt.r2)
			if result != tt.expected {
				t.Errorf("distanceFromRectSquared: got %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestScreen_Origin(t *testing.T) {
	screen := Screen{X: 100, Y: 200}
	origin := screen.Origin()

	if origin.X != 100 {
		t.Errorf("origin.X = %d, want 100", origin.X)
	}
	if origin.Y != 200 {
		t.Errorf("origin.Y = %d, want 200", origin.Y)
	}
}

func TestScreen_scale(t *testing.T) {
	screen := Screen{ScaleFactor: 2.0}

	tests := []struct {
		value    int
		toDip    bool
		expected int
	}{
		{100, true, 50},    // to DIP: 100 / 2 = 50
		{100, false, 200},  // to physical: 100 * 2 = 200
		{101, true, 51},    // to DIP: ceil(101 / 2) = 51
		{101, false, 202},  // to physical: floor(101 * 2) = 202
		{0, true, 0},
		{0, false, 0},
	}

	for _, tt := range tests {
		result := screen.scale(tt.value, tt.toDip)
		if result != tt.expected {
			t.Errorf("scale(%d, %v) = %d, want %d", tt.value, tt.toDip, result, tt.expected)
		}
	}
}

func TestScreen_scale_1_5(t *testing.T) {
	screen := Screen{ScaleFactor: 1.5}

	tests := []struct {
		value    int
		toDip    bool
		expected int
	}{
		{150, true, 100},   // to DIP: ceil(150 / 1.5) = 100
		{100, false, 150},  // to physical: floor(100 * 1.5) = 150
		{100, true, 67},    // to DIP: ceil(100 / 1.5) = 67
		{67, false, 100},   // to physical: floor(67 * 1.5) = 100
	}

	for _, tt := range tests {
		result := screen.scale(tt.value, tt.toDip)
		if result != tt.expected {
			t.Errorf("scale(%d, %v) with factor 1.5 = %d, want %d", tt.value, tt.toDip, result, tt.expected)
		}
	}
}

func TestScreen_right(t *testing.T) {
	screen := Screen{
		Bounds: Rect{X: 100, Y: 0, Width: 200, Height: 100},
	}
	if screen.right() != 300 {
		t.Errorf("right() = %d, want 300", screen.right())
	}
}

func TestScreen_bottom(t *testing.T) {
	screen := Screen{
		Bounds: Rect{X: 0, Y: 100, Width: 100, Height: 200},
	}
	if screen.bottom() != 300 {
		t.Errorf("bottom() = %d, want 300", screen.bottom())
	}
}

func TestScreen_intersects(t *testing.T) {
	screen1 := &Screen{
		X:      0,
		Y:      0,
		Bounds: Rect{X: 0, Y: 0, Width: 100, Height: 100},
	}

	screen2 := &Screen{
		X:      50,
		Y:      50,
		Bounds: Rect{X: 50, Y: 50, Width: 100, Height: 100},
	}

	screen3 := &Screen{
		X:      200,
		Y:      0,
		Bounds: Rect{X: 200, Y: 0, Width: 100, Height: 100},
	}

	if !screen1.intersects(screen2) {
		t.Error("screen1 and screen2 should intersect")
	}
	if screen1.intersects(screen3) {
		t.Error("screen1 and screen3 should not intersect")
	}
}

func TestPoint_Fields(t *testing.T) {
	pt := Point{X: 10, Y: 20}
	if pt.X != 10 || pt.Y != 20 {
		t.Error("Point fields not set correctly")
	}
}

func TestSize_Fields(t *testing.T) {
	size := Size{Width: 100, Height: 200}
	if size.Width != 100 || size.Height != 200 {
		t.Error("Size fields not set correctly")
	}
}

func TestScreen_Fields(t *testing.T) {
	screen := Screen{
		ID:          "display-1",
		Name:        "Primary Display",
		ScaleFactor: 2.0,
		X:           0,
		Y:           0,
		Size:        Size{Width: 1920, Height: 1080},
		Bounds:      Rect{X: 0, Y: 0, Width: 1920, Height: 1080},
		IsPrimary:   true,
		Rotation:    0,
	}

	if screen.ID != "display-1" {
		t.Error("ID not set correctly")
	}
	if screen.Name != "Primary Display" {
		t.Error("Name not set correctly")
	}
	if screen.ScaleFactor != 2.0 {
		t.Error("ScaleFactor not set correctly")
	}
	if !screen.IsPrimary {
		t.Error("IsPrimary not set correctly")
	}
}

func TestScreenPlacement_Fields(t *testing.T) {
	parent := &Screen{ID: "parent"}
	child := &Screen{ID: "child"}

	placement := ScreenPlacement{
		screen:          child,
		parent:          parent,
		alignment:       RIGHT,
		offset:          100,
		offsetReference: BEGIN,
	}

	if placement.screen != child {
		t.Error("screen not set correctly")
	}
	if placement.parent != parent {
		t.Error("parent not set correctly")
	}
	if placement.alignment != RIGHT {
		t.Error("alignment not set correctly")
	}
	if placement.offset != 100 {
		t.Error("offset not set correctly")
	}
	if placement.offsetReference != BEGIN {
		t.Error("offsetReference not set correctly")
	}
}

func TestNewScreenManager(t *testing.T) {
	sm := newScreenManager(nil)
	if sm == nil {
		t.Fatal("newScreenManager returned nil")
	}
	if sm.screens != nil {
		t.Error("screens should be nil initially")
	}
	if sm.primaryScreen != nil {
		t.Error("primaryScreen should be nil initially")
	}
}
