//go:build darwin && purego && !ios && !server

package application

import (
	"math"
	"testing"
)

// Screen geometry of a 14-inch MacBook Pro display (1512x982 points, menu bar
// 25 points) and a housing gap from x=600 to x=912 (756 is its centre).
var (
	testScreenFrame   = NSRect{Origin: CGPoint{X: 0, Y: 0}, Size: CGSize{Width: 1512, Height: 982}}
	testVisibleFrame  = NSRect{Origin: CGPoint{X: 0, Y: 0}, Size: CGSize{Width: 1512, Height: 957}}
	testLeftArea      = NSRect{Origin: CGPoint{X: 0, Y: 950}, Size: CGSize{Width: 600, Height: 32}}
	testRightArea     = NSRect{Origin: CGPoint{X: 912, Y: 950}, Size: CGSize{Width: 600, Height: 32}}
	testHousingCentre = 756.0
	testNotchFrame    = NSRect{Origin: CGPoint{X: 0, Y: 0}, Size: CGSize{Width: 700, Height: 104}}
)

func approxEqual(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestNotchBackgroundOutline(t *testing.T) {
	// Expected points are computed by hand for a 200x100 view with wing 20 and
	// body radius 28: top edge, right wing, right flank, bottom edge, left
	// flank, left wing. Control points use the 0.58/0.42/0.45 factors of
	// webview_notch_window_darwin.m.
	segs := notchBackgroundOutline(200, 100)
	if len(segs) != 9 {
		t.Fatalf("len(outline) = %d, want 9", len(segs))
	}
	type want struct {
		curve  bool
		c1, c2 CGPoint
		end    CGPoint
	}
	wants := []want{
		0: {end: CGPoint{X: 0, Y: 100}},                                                                            // moveTo
		1: {end: CGPoint{X: 200, Y: 100}},                                                                          // top edge
		2: {curve: true, c1: CGPoint{X: 188.4, Y: 100}, c2: CGPoint{X: 180, Y: 91.6}, end: CGPoint{X: 180, Y: 80}}, // right wing
		3: {end: CGPoint{X: 180, Y: 28}},                                                                           // right flank
		4: {curve: true, c1: CGPoint{X: 180, Y: 12.6}, c2: CGPoint{X: 167.4, Y: 0}, end: CGPoint{X: 152, Y: 0}},    // bottom right
		5: {end: CGPoint{X: 48, Y: 0}},                                                                             // bottom edge
		6: {curve: true, c1: CGPoint{X: 32.6, Y: 0}, c2: CGPoint{X: 20, Y: 12.6}, end: CGPoint{X: 20, Y: 28}},      // bottom left
		7: {end: CGPoint{X: 20, Y: 80}},                                                                            // left flank
		8: {curve: true, c1: CGPoint{X: 20, Y: 91.6}, c2: CGPoint{X: 11.6, Y: 100}, end: CGPoint{X: 0, Y: 100}},    // left wing
	}
	for i, w := range wants {
		s := segs[i]
		if s.Curve != w.curve {
			t.Errorf("segment %d Curve = %v, want %v", i, s.Curve, w.curve)
		}
		if !approxEqual(s.End.X, w.end.X) || !approxEqual(s.End.Y, w.end.Y) {
			t.Errorf("segment %d End = %+v, want %+v", i, s.End, w.end)
		}
		if w.curve {
			if !approxEqual(s.C1.X, w.c1.X) || !approxEqual(s.C1.Y, w.c1.Y) ||
				!approxEqual(s.C2.X, w.c2.X) || !approxEqual(s.C2.Y, w.c2.Y) {
				t.Errorf("segment %d controls = %+v %+v, want %+v %+v", i, s.C1, s.C2, w.c1, w.c2)
			}
		}
	}
}

func TestNotchBackgroundOutlineStartsAtTopLeft(t *testing.T) {
	// The outline must start at the top-left corner so the first moveToPoint:
	// is at (0, height) for any size; a zero-value first segment is a bug.
	for _, size := range []struct{ w, h float64 }{{1, 1}, {660, 92}, {800, 120}} {
		segs := notchBackgroundOutline(size.w, size.h)
		if segs[0].Curve || segs[0].End != (CGPoint{X: 0, Y: size.h}) {
			t.Errorf("size %vx%v: first segment = %+v, want move to (0,%v)", size.w, size.h, segs[0], size.h)
		}
	}
}

func TestNotchCameraHousingCentre(t *testing.T) {
	tests := []struct {
		name        string
		left, right NSRect
		safeTop     float64
		available   bool
		wantCentre  float64
		wantFound   bool
	}{
		{name: "housing gap centre", left: testLeftArea, right: testRightArea, safeTop: 32, available: true,
			wantCentre: testHousingCentre, wantFound: true},
		{name: "no safe-area inset means no housing", left: testLeftArea, right: testRightArea, safeTop: 0, available: true},
		{name: "negative safe-area inset means no housing", left: testLeftArea, right: testRightArea, safeTop: -1, available: true},
		{name: "areas unavailable before macOS 12", left: testLeftArea, right: testRightArea, safeTop: 32, available: false},
		{name: "empty left area", left: NSRect{Origin: CGPoint{X: 0, Y: 950}}, right: testRightArea, safeTop: 32, available: true},
		{name: "empty right area", left: testLeftArea, right: NSRect{Origin: CGPoint{X: 912, Y: 950}}, safeTop: 32, available: true},
		{name: "areas overlap", left: NSRect{Origin: CGPoint{X: 0, Y: 950}, Size: CGSize{Width: 1000, Height: 32}},
			right: testRightArea, safeTop: 32, available: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			centre, found := notchCameraHousingCentre(tt.left, tt.right, tt.safeTop, tt.available)
			if found != tt.wantFound {
				t.Fatalf("found = %v, want %v", found, tt.wantFound)
			}
			if found && !approxEqual(centre, tt.wantCentre) {
				t.Fatalf("centre = %v, want %v", centre, tt.wantCentre)
			}
		})
	}
}

func TestNotchTargetFrame(t *testing.T) {
	secondScreen := NSRect{Origin: CGPoint{X: -1512, Y: 0}, Size: CGSize{Width: 1512, Height: 982}}
	secondVisible := NSRect{Origin: CGPoint{X: -1512, Y: 0}, Size: CGSize{Width: 1512, Height: 957}}
	tests := []struct {
		name         string
		frame        NSRect
		screen       NSRect
		visible      NSRect
		centre       float64
		hasHousing   bool
		wantX, wantY float64
	}{
		{
			name:  "camera housing: flush with the screen top, centred on the housing",
			frame: testNotchFrame, screen: testScreenFrame, visible: testVisibleFrame,
			centre: testHousingCentre, hasHousing: true,
			wantX: 406, wantY: 878, // 756 - 350, 982 - 104
		},
		{
			name:  "no housing: flush with the visible top (below the menu bar), centred on the screen",
			frame: testNotchFrame, screen: testScreenFrame, visible: testVisibleFrame,
			hasHousing: false,
			wantX:      406, wantY: 853, // 957 - 104
		},
		{
			name:  "odd width rounds half away from zero like the native round()",
			frame: NSRect{Size: CGSize{Width: 701, Height: 104}}, screen: testScreenFrame, visible: testVisibleFrame,
			centre: testHousingCentre, hasHousing: true,
			wantX: 406, wantY: 878, // 756 - 350.5 = 405.5 -> 406
		},
		{
			name:  "housing centre fraction rounds to the nearest point",
			frame: testNotchFrame, screen: testScreenFrame, visible: testVisibleFrame,
			centre: 756.4, hasHousing: true,
			wantX: 406, wantY: 878, // 756.4 - 350 = 406.4 -> 406
		},
		{
			name:  "secondary display with negative origin",
			frame: testNotchFrame, screen: secondScreen, visible: secondVisible,
			hasHousing: false,
			wantX:      -1106, wantY: 853, // -756 - 350
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := notchTargetFrame(tt.frame, tt.screen, tt.visible, tt.centre, tt.hasHousing)
			if !approxEqual(got.Origin.X, tt.wantX) || !approxEqual(got.Origin.Y, tt.wantY) {
				t.Fatalf("origin = (%v,%v), want (%v,%v)", got.Origin.X, got.Origin.Y, tt.wantX, tt.wantY)
			}
			if got.Size != tt.frame.Size {
				t.Fatalf("size changed to %+v, want %+v", got.Size, tt.frame.Size)
			}
		})
	}
}

func TestNotchHiddenFrame(t *testing.T) {
	tests := []struct {
		name   string
		target NSRect
		wantY  float64
	}{
		{name: "docked at the screen top", target: NSRect{Origin: CGPoint{X: 406, Y: 878}, Size: CGSize{Width: 700, Height: 104}}, wantY: 982},
		{name: "docked below the menu bar", target: NSRect{Origin: CGPoint{X: 406, Y: 853}, Size: CGSize{Width: 700, Height: 104}}, wantY: 957},
		{name: "zero height keeps the origin", target: NSRect{Origin: CGPoint{X: 0, Y: 10}}, wantY: 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hidden := notchHiddenFrame(tt.target)
			if hidden.Origin.X != tt.target.Origin.X || hidden.Size != tt.target.Size {
				t.Fatalf("hidden frame moved horizontally or resized: %+v from %+v", hidden, tt.target)
			}
			if !approxEqual(hidden.Origin.Y, tt.wantY) {
				t.Fatalf("hidden Y = %v, want %v", hidden.Origin.Y, tt.wantY)
			}
		})
	}
}

func TestNotchFrameForNilScreenKeepsFrame(t *testing.T) {
	if got := notchFrameForScreen(testNotchFrame, id(0)); got != testNotchFrame {
		t.Fatalf("notchFrameForScreen(nil screen) = %+v, want unchanged %+v", got, testNotchFrame)
	}
}

func TestNotchTargetScreenFallsBackToPrimary(t *testing.T) {
	// An unknown ID (and the empty ID) must resolve to the primary screen, the
	// first of [NSScreen screens], as cgo notchTargetScreen does.
	screens := class("NSScreen").send("screens")
	primary := screens.send("firstObject")
	if got := notchTargetScreen("no-such-display"); got.ptr() != primary.ptr() {
		t.Fatalf("notchTargetScreen(unknown) = %#x, want primary %#x", got.ptr(), primary.ptr())
	}
	if got := notchTargetScreen(""); got.ptr() != primary.ptr() {
		t.Fatalf("notchTargetScreen(\"\") = %#x, want primary %#x", got.ptr(), primary.ptr())
	}
}

func TestNotchBackgroundBezierPathElements(t *testing.T) {
	// Checks the production path builder against the outline element by
	// element: a move, then a line or curve per outline segment, then the
	// close. NSBezierPath also appends an implicit moveTo after closePath, so
	// the element count is not checked exactly.
	const (
		moveTo    = 0 // NSBezierPathElementMoveTo
		lineTo    = 1 // NSBezierPathElementLineTo
		curveTo   = 2 // NSBezierPathElementCurveTo
		closePath = 3 // NSBezierPathElementClosePath
	)
	segments := notchBackgroundOutline(200, 100)
	path := notchBackgroundBezierPath(200, 100)
	if got := get[int](path, "elementCount"); got < len(segments)+1 {
		t.Fatalf("elementCount = %d, want at least %d", got, len(segments)+1)
	}
	for i, s := range segments {
		want := lineTo
		switch {
		case i == 0:
			want = moveTo
		case s.Curve:
			want = curveTo
		}
		if got := get[int](path, "elementAtIndex:", i); got != want {
			t.Errorf("element %d = %d, want %d", i, got, want)
		}
	}
	if got := get[int](path, "elementAtIndex:", len(segments)); got != closePath {
		t.Errorf("element %d = %d, want closePath", len(segments), got)
	}
}

func TestNotchBackgroundClassMethods(t *testing.T) {
	cls := registerNotchBackgroundClass()
	if !get[bool](cls, "isSubclassOfClass:", class("NSView")) {
		t.Fatal("WailsNotchBackgroundView is not an NSView subclass")
	}
	for _, sel := range []string{"drawRect:", "isOpaque"} {
		if !get[bool](cls, "instancesRespondToSelector:", sel_(sel)) {
			t.Errorf("WailsNotchBackgroundView does not respond to %s", sel)
		}
	}
}

func TestNotchWindowClassIsPanelSubclass(t *testing.T) {
	cls := registerNotchWindowClass()
	if !get[bool](cls, "isSubclassOfClass:", class("NSPanel")) {
		t.Fatal("WailsWebviewNotchPanel is not an NSPanel subclass")
	}
	if !get[bool](cls, "isSubclassOfClass:", registerWebviewPanelClass()) {
		t.Fatal("WailsWebviewNotchPanel does not inherit from the panel class")
	}
	if !get[bool](cls, "instancesRespondToSelector:", sel_("performWindowDragWithEvent:")) {
		t.Fatal("WailsWebviewNotchPanel does not respond to performWindowDragWithEvent:")
	}
	if !get[bool](cls, "instancesRespondToSelector:", sel_("canBecomeKeyWindow")) {
		t.Fatal("WailsWebviewNotchPanel lost the inherited canBecomeKeyWindow")
	}
}
