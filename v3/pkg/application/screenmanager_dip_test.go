package application

import (
	"testing"
)

// dipScreen describes a display the way macOS and GTK report one: Bounds and
// WorkArea in logical coordinates, against a single origin at the primary's
// top-left.
type dipScreen struct {
	id       string
	scale    float32
	bounds   Rect
	workArea Rect
	primary  bool
}

func layoutDIP(t *testing.T, screens []dipScreen) *ScreenManager {
	t.Helper()
	list := make([]*Screen, 0, len(screens))
	for _, s := range screens {
		list = append(list, &Screen{
			ID:          s.id,
			Name:        s.id,
			ScaleFactor: s.scale,
			IsPrimary:   s.primary,
			Bounds:      s.bounds,
			WorkArea:    s.workArea,
		})
	}
	manager := &ScreenManager{}
	if err := manager.layoutDIPScreens(list); err != nil {
		t.Fatalf("layoutDIPScreens: %v", err)
	}
	return manager
}

func rect(x, y, w, h int) Rect {
	return Rect{X: x, Y: y, Width: w, Height: h}
}

// TestLayoutDIPScreensKeepsMixedDPICoordinates pins the contract of the
// already-DIP input path: a display's Bounds and WorkArea come back in the
// coordinates the platform reported, whatever its scale factor.
//
// The mixed-DPI arrangements are the point. Rebuilding these positions from
// synthesised device pixels (LayoutScreens) cannot reproduce them, because
// displays that abut in logical space stop abutting once each origin is
// multiplied by its own scale factor: the screen is never reached by the
// placement walk and keeps coordinates no window is ever reported in.
func TestLayoutDIPScreensKeepsMixedDPICoordinates(t *testing.T) {
	cases := []struct {
		name    string
		screens []dipScreen
	}{
		{
			name: "2x primary with 1x display to the right",
			screens: []dipScreen{
				{id: "builtin", scale: 2, bounds: rect(0, 0, 1512, 982), workArea: rect(0, 37, 1512, 945), primary: true},
				{id: "external", scale: 1, bounds: rect(1512, 0, 1920, 1080), workArea: rect(1512, 0, 1920, 1080)},
			},
		},
		{
			name: "1x primary with 2x display to the right",
			screens: []dipScreen{
				{id: "external", scale: 1, bounds: rect(0, 0, 1920, 1080), workArea: rect(0, 25, 1920, 1055), primary: true},
				{id: "builtin", scale: 2, bounds: rect(1920, 0, 1512, 982), workArea: rect(1920, 37, 1512, 945)},
			},
		},
		{
			name: "1x primary with 2x display below",
			screens: []dipScreen{
				{id: "external", scale: 1, bounds: rect(0, 0, 1920, 1080), workArea: rect(0, 25, 1920, 1055), primary: true},
				{id: "builtin", scale: 2, bounds: rect(0, 1080, 1512, 982), workArea: rect(0, 1117, 1512, 945)},
			},
		},
		{
			name: "1x primary with 2x display to the left",
			screens: []dipScreen{
				{id: "external", scale: 1, bounds: rect(0, 0, 1920, 1080), workArea: rect(0, 25, 1920, 1055), primary: true},
				{id: "builtin", scale: 2, bounds: rect(-1512, 0, 1512, 982), workArea: rect(-1512, 37, 1512, 945)},
			},
		},
		{
			// The arrangement reported in #5409.
			name: "1x primary with 2x display above",
			screens: []dipScreen{
				{id: "external", scale: 1, bounds: rect(0, 0, 1920, 1080), workArea: rect(0, 25, 1920, 1055), primary: true},
				{id: "builtin", scale: 2, bounds: rect(0, -982, 1512, 982), workArea: rect(0, -945, 1512, 945)},
			},
		},
		{
			name: "three displays, three scale factors",
			screens: []dipScreen{
				{id: "left", scale: 1, bounds: rect(-2560, 0, 2560, 1440), workArea: rect(-2560, 0, 2560, 1440)},
				{id: "middle", scale: 2, bounds: rect(0, 0, 1512, 982), workArea: rect(0, 37, 1512, 945), primary: true},
				{id: "right", scale: 3, bounds: rect(1512, 0, 1280, 800), workArea: rect(1512, 0, 1280, 800)},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			manager := layoutDIP(t, tc.screens)
			got := manager.GetAll()
			if len(got) != len(tc.screens) {
				t.Fatalf("got %d screens, want %d", len(got), len(tc.screens))
			}
			for i, want := range tc.screens {
				screen := got[i]
				if screen.Bounds != want.bounds {
					t.Errorf("%s: Bounds = %+v, want %+v", want.id, screen.Bounds, want.bounds)
				}
				if screen.WorkArea != want.workArea {
					t.Errorf("%s: WorkArea = %+v, want %+v", want.id, screen.WorkArea, want.workArea)
				}
				if screen.X != want.bounds.X || screen.Y != want.bounds.Y {
					t.Errorf("%s: X,Y = %d,%d, want %d,%d", want.id, screen.X, screen.Y, want.bounds.X, want.bounds.Y)
				}
			}
		})
	}
}

// TestLayoutDIPScreensFindsWindowOnSecondDisplay is the same property stated the
// way an application meets it. A window dragged onto the second display has to
// be found on that display: code asking "which screen is this window on" is what
// decides whether a window has been left somewhere unreachable, and a display
// held in the wrong coordinate space answers "none of them" for every window on
// it.
func TestLayoutDIPScreensFindsWindowOnSecondDisplay(t *testing.T) {
	manager := layoutDIP(t, []dipScreen{
		{id: "external", scale: 1, bounds: rect(0, 0, 1920, 1080), workArea: rect(0, 25, 1920, 1055), primary: true},
		{id: "builtin", scale: 2, bounds: rect(1920, 0, 1512, 982), workArea: rect(1920, 37, 1512, 945)},
	})

	window := rect(2076, 128, 1200, 800)

	var found *Screen
	for _, screen := range manager.GetAll() {
		if !screen.WorkArea.Intersect(window).IsEmpty() {
			found = screen
			break
		}
	}
	if found == nil {
		t.Fatalf("window %+v was found on no display", window)
	}
	if found.ID != "builtin" {
		t.Errorf("window %+v found on %q, want \"builtin\"", window, found.ID)
	}
}

// TestLayoutDIPScreensCompletesBookkeeping covers the fields the DIP path fills
// in for a platform that reports only the logical rectangles, and the refusal to
// overwrite a Physical* rectangle a platform measured for itself.
func TestLayoutDIPScreensCompletesBookkeeping(t *testing.T) {
	measured := rect(7, 9, 11, 13)
	screens := []*Screen{
		{ID: "derived", ScaleFactor: 2, IsPrimary: true, Bounds: rect(0, 0, 1512, 982), WorkArea: rect(0, 37, 1512, 945)},
		{ID: "measured", ScaleFactor: 2, Bounds: rect(1512, 0, 1280, 800), WorkArea: rect(1512, 0, 1280, 800), PhysicalBounds: measured},
	}

	manager := &ScreenManager{}
	if err := manager.layoutDIPScreens(screens); err != nil {
		t.Fatalf("layoutDIPScreens: %v", err)
	}

	derived := manager.GetAll()[0]
	if want := (Size{Width: 1512, Height: 982}); derived.Size != want {
		t.Errorf("Size = %+v, want %+v", derived.Size, want)
	}
	if want := rect(0, 0, 3024, 1964); derived.PhysicalBounds != want {
		t.Errorf("PhysicalBounds = %+v, want %+v", derived.PhysicalBounds, want)
	}
	if want := rect(0, 74, 3024, 1890); derived.PhysicalWorkArea != want {
		t.Errorf("PhysicalWorkArea = %+v, want %+v", derived.PhysicalWorkArea, want)
	}

	if got := manager.GetAll()[1].PhysicalBounds; got != measured {
		t.Errorf("PhysicalBounds = %+v, want the measured %+v", got, measured)
	}

	if primary := manager.GetPrimary(); primary == nil || primary.ID != "derived" {
		t.Errorf("GetPrimary() = %v, want the screen marked primary", primary)
	}
}

// TestLayoutDIPScreensRejectsUnusableInput keeps the two refusals that make a
// bad screen list a reported error rather than a cached one.
func TestLayoutDIPScreensRejectsUnusableInput(t *testing.T) {
	manager := &ScreenManager{}
	if err := manager.layoutDIPScreens(nil); err == nil {
		t.Error("layoutDIPScreens(nil) = nil, want an error")
	}

	noPrimary := []*Screen{{ID: "only", ScaleFactor: 1, Bounds: rect(0, 0, 800, 600)}}
	if err := manager.layoutDIPScreens(noPrimary); err == nil {
		t.Error("layoutDIPScreens with no primary = nil, want an error")
	}
}

// TestLayoutDIPScreensDropsNilEntries keeps a nil entry in a platform's list
// out of the cache, where GetAll and the ScreenNearest* lookups would meet it.
func TestLayoutDIPScreensDropsNilEntries(t *testing.T) {
	manager := &ScreenManager{}
	primary := &Screen{ID: "primary", ScaleFactor: 1, IsPrimary: true, Bounds: rect(0, 0, 800, 600)}
	if err := manager.layoutDIPScreens([]*Screen{nil, primary, nil}); err != nil {
		t.Fatalf("layoutDIPScreens: %v", err)
	}
	all := manager.GetAll()
	if len(all) != 1 || all[0] != primary {
		t.Fatalf("GetAll() = %v, want only the primary", all)
	}
}
