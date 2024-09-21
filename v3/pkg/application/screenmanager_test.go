package application_test

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"testing"

	"github.com/matryer/is"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type ScreenDef struct {
	id     int
	w, h   int
	s      float32
	parent ScreenDefParent
	name   string
}

type ScreenDefParent struct {
	id     int
	align  string
	offset int
}

type ScreensLayout struct {
	name    string
	screens []ScreenDef
}

type ParsedLayout struct {
	name    string
	screens []*application.Screen
}

func exampleLayouts() []ParsedLayout {
	layouts := [][]ScreensLayout{
		{
			// Normal examples (demonstrate real life scenarios)
			{
				name: "Single 4k monitor",
				screens: []ScreenDef{
					{id: 1, w: 3840, h: 2160, s: 163.0 / 96, name: `27" 4K UHD 163DPI`},
				},
			},
			{
				name: "Two monitors",
				screens: []ScreenDef{
					{id: 1, w: 3840, h: 2160, s: 163.0 / 96, name: `27" 4K UHD 163DPI`},
					{id: 2, w: 1920, h: 1080, s: 1, parent: ScreenDefParent{id: 1, align: "r", offset: 0}, name: `23" FHD 96DPI`},
				},
			},
			{
				name: "Two monitors (2)",
				screens: []ScreenDef{
					{id: 1, w: 1920, h: 1080, s: 1, name: `23" FHD 96DPI`},
					{id: 2, w: 1920, h: 1080, s: 1.25, parent: ScreenDefParent{id: 1, align: "r", offset: 0}, name: `23" FHD 96DPI (125%)`},
				},
			},
			{
				name: "Three monitors",
				screens: []ScreenDef{
					{id: 1, w: 3840, h: 2160, s: 163.0 / 96, name: `27" 4K UHD 163DPI`},
					{id: 2, w: 1920, h: 1080, s: 1, parent: ScreenDefParent{id: 1, align: "r", offset: 0}, name: `23" FHD 96DPI`},
					{id: 3, w: 1920, h: 1080, s: 1.25, parent: ScreenDefParent{id: 1, align: "l", offset: 0}, name: `23" FHD 96DPI (125%)`},
				},
			},
			{
				name: "Four monitors",
				screens: []ScreenDef{
					{id: 1, w: 3840, h: 2160, s: 163.0 / 96, name: `27" 4K UHD 163DPI`},
					{id: 2, w: 1920, h: 1080, s: 1, parent: ScreenDefParent{id: 1, align: "r", offset: 0}, name: `23" FHD 96DPI`},
					{id: 3, w: 1920, h: 1080, s: 1.25, parent: ScreenDefParent{id: 2, align: "b", offset: 0}, name: `23" FHD 96DPI (125%)`},
					{id: 4, w: 1080, h: 1920, s: 1, parent: ScreenDefParent{id: 1, align: "l", offset: 0}, name: `23" FHD (90deg)`},
				},
			},
		},
		{
			// Test cases examples (demonstrate the algorithm basics)
			{
				name: "Child scaled, Start offset",
				screens: []ScreenDef{
					{id: 1, w: 1200, h: 1200, s: 1, name: "Parent"},
					{id: 2, w: 1200, h: 1200, s: 1.5, parent: ScreenDefParent{id: 1, align: "r", offset: 600}, name: "Child"},
				},
			},
			{
				name: "Child scaled, End offset",
				screens: []ScreenDef{
					{id: 1, w: 1200, h: 1200, s: 1, name: "Parent"},
					{id: 2, w: 1200, h: 1200, s: 1.5, parent: ScreenDefParent{id: 1, align: "r", offset: -600}, name: "Child"},
				},
			},
			{
				name: "Parent scaled, Start offset percent",
				screens: []ScreenDef{
					{id: 1, w: 1200, h: 1200, s: 1.5, name: "Parent"},
					{id: 2, w: 1200, h: 1200, s: 1, parent: ScreenDefParent{id: 1, align: "r", offset: 600}, name: "Child"},
				},
			},
			{
				name: "Parent scaled, End offset percent",
				screens: []ScreenDef{
					{id: 1, w: 1200, h: 1200, s: 1.5, name: "Parent"},
					{id: 2, w: 1200, h: 1200, s: 1, parent: ScreenDefParent{id: 1, align: "r", offset: -600}, name: "Child"},
				},
			},
			{
				name: "Parent scaled, Start align",
				screens: []ScreenDef{
					{id: 1, w: 1200, h: 1200, s: 1.5, name: "Parent"},
					{id: 2, w: 1200, h: 1100, s: 1, parent: ScreenDefParent{id: 1, align: "r", offset: 0}, name: "Child"},
				},
			},
			{
				name: "Parent scaled, End align",
				screens: []ScreenDef{
					{id: 1, w: 1200, h: 1200, s: 1.5, name: "Parent"},
					{id: 2, w: 1200, h: 1200, s: 1, parent: ScreenDefParent{id: 1, align: "r", offset: 0}, name: "Child"},
				},
			},
			{
				name: "Parent scaled, in-between",
				screens: []ScreenDef{
					{id: 1, w: 1200, h: 1200, s: 1.5, name: "Parent"},
					{id: 2, w: 1200, h: 1500, s: 1, parent: ScreenDefParent{id: 1, align: "r", offset: -250}, name: "Child"},
				},
			},
		},
		{
			// Edge cases examples
			{
				name: "Parent order (5 is parent of 4)",
				screens: []ScreenDef{
					{id: 1, w: 1920, h: 1080, s: 1},
					{id: 2, w: 1024, h: 600, s: 1.25, parent: ScreenDefParent{id: 1, align: "r", offset: -200}},
					{id: 3, w: 800, h: 800, s: 1.25, parent: ScreenDefParent{id: 2, align: "b", offset: 0}},
					{id: 4, w: 800, h: 1080, s: 1.5, parent: ScreenDefParent{id: 2, align: "re", offset: 100}},
					{id: 5, w: 600, h: 600, s: 1, parent: ScreenDefParent{id: 3, align: "r", offset: 100}},
				},
			},
			{
				name: "de-intersection reparent",
				screens: []ScreenDef{
					{id: 1, w: 1920, h: 1080, s: 1},
					{id: 2, w: 1680, h: 1050, s: 1.25, parent: ScreenDefParent{id: 1, align: "r", offset: 10}},
					{id: 3, w: 1440, h: 900, s: 1.5, parent: ScreenDefParent{id: 1, align: "le", offset: 150}},
					{id: 4, w: 1024, h: 768, s: 1, parent: ScreenDefParent{id: 3, align: "bc", offset: -200}},
					{id: 5, w: 1024, h: 768, s: 1.25, parent: ScreenDefParent{id: 4, align: "r", offset: 400}},
				},
			},
			{
				name: "de-intersection (unattached child)",
				screens: []ScreenDef{
					{id: 1, w: 1920, h: 1080, s: 1},
					{id: 2, w: 1024, h: 768, s: 1.5, parent: ScreenDefParent{id: 1, align: "le", offset: 10}},
					{id: 3, w: 1024, h: 768, s: 1.25, parent: ScreenDefParent{id: 2, align: "b", offset: 100}},
					{id: 4, w: 1024, h: 768, s: 1, parent: ScreenDefParent{id: 3, align: "r", offset: 500}},
				},
			},
			{
				name: "Multiple de-intersection",
				screens: []ScreenDef{
					{id: 1, w: 1920, h: 1080, s: 1},
					{id: 2, w: 1024, h: 768, s: 1, parent: ScreenDefParent{id: 1, align: "be", offset: 0}},
					{id: 3, w: 1024, h: 768, s: 1, parent: ScreenDefParent{id: 2, align: "b", offset: 300}},
					{id: 4, w: 1024, h: 768, s: 1.5, parent: ScreenDefParent{id: 2, align: "le", offset: 100}},
					{id: 5, w: 1024, h: 768, s: 1, parent: ScreenDefParent{id: 4, align: "be", offset: 100}},
				},
			},
			{
				name: "Multiple de-intersection (left-side)",
				screens: []ScreenDef{
					{id: 1, w: 1920, h: 1080, s: 1},
					{id: 2, w: 1024, h: 768, s: 1, parent: ScreenDefParent{id: 1, align: "le", offset: 0}},
					{id: 3, w: 1024, h: 768, s: 1, parent: ScreenDefParent{id: 2, align: "b", offset: 300}},
					{id: 4, w: 1024, h: 768, s: 1.5, parent: ScreenDefParent{id: 2, align: "le", offset: 100}},
					{id: 5, w: 1024, h: 768, s: 1, parent: ScreenDefParent{id: 4, align: "be", offset: 100}},
				},
			},
			{
				name: "Parent de-intersection child offset",
				screens: []ScreenDef{
					{id: 1, w: 1600, h: 1600, s: 1.5},
					{id: 2, w: 800, h: 800, s: 1, parent: ScreenDefParent{id: 1, align: "r", offset: 0}},
					{id: 3, w: 800, h: 800, s: 1, parent: ScreenDefParent{id: 1, align: "r", offset: 800}},
					{id: 4, w: 800, h: 1600, s: 1, parent: ScreenDefParent{id: 2, align: "r", offset: 0}},
				},
			},
		},
	}

	parsedLayouts := []ParsedLayout{}

	for _, section := range layouts {
		for _, layout := range section {
			parsedLayouts = append(parsedLayouts, parseLayout(layout))
		}
	}

	return parsedLayouts
}

// Parse screens layout from easy-to-define ScreenDef for testing to actual Screens layout
func parseLayout(layout ScreensLayout) ParsedLayout {
	screens := []*application.Screen{}

	for _, screen := range layout.screens {
		var x, y int
		w := screen.w
		h := screen.h

		if screen.parent.id > 0 {
			idx := slices.IndexFunc(screens, func(s *application.Screen) bool { return s.ID == strconv.Itoa(screen.parent.id) })
			parent := screens[idx].Bounds
			offset := screen.parent.offset
			align := screen.parent.align
			align2 := ""

			if len(align) == 2 {
				align2 = string(align[1])
				align = string(align[0])
			}

			x = parent.X
			y = parent.Y
			// t: top, b: bottom, l: left, r: right, e: edge, c: corner
			if align == "t" || align == "b" {
				x += offset
				if align2 == "e" || align2 == "c" {
					x += parent.Width
				}
				if align2 == "e" {
					x -= w
				}
				if align == "t" {
					y -= h
				} else {
					y += parent.Height
				}
			} else {
				y += offset
				if align2 == "e" || align2 == "c" {
					y += parent.Height
				}
				if align2 == "e" {
					y -= h
				}
				if align == "l" {
					x -= w
				} else {
					x += parent.Width
				}
			}
		}
		name := screen.name
		if name == "" {
			name = "Display" + strconv.Itoa(screen.id)
		}
		screens = append(screens, &application.Screen{
			ID:               strconv.Itoa(screen.id),
			Name:             name,
			ScaleFactor:      float32(math.Round(float64(screen.s)*100) / 100),
			X:                x,
			Y:                y,
			Size:             application.Size{Width: w, Height: h},
			Bounds:           application.Rect{X: x, Y: y, Width: w, Height: h},
			PhysicalBounds:   application.Rect{X: x, Y: y, Width: w, Height: h},
			WorkArea:         application.Rect{X: x, Y: y, Width: w, Height: h - int(40*screen.s)},
			PhysicalWorkArea: application.Rect{X: x, Y: y, Width: w, Height: h - int(40*screen.s)},
			IsPrimary:        screen.id == 1,
			Rotation:         0,
		})
	}
	return ParsedLayout{
		name:    layout.name,
		screens: screens,
	}
}

func matchRects(r1, r2 application.Rect) error {
	threshold := 1.0
	if math.Abs(float64(r1.X-r2.X)) > threshold ||
		math.Abs(float64(r1.Y-r2.Y)) > threshold ||
		math.Abs(float64(r1.Width-r2.Width)) > threshold ||
		math.Abs(float64(r1.Height-r2.Height)) > threshold {
		return fmt.Errorf("%v != %v", r1, r2)
	}
	return nil
}

// Test screens layout (DPI transformation)
func TestScreenManager_ScreensLayout(t *testing.T) {
	sm := application.ScreenManager{}

	t.Run("Child scaled", func(t *testing.T) {
		is := is.New(t)

		layout := parseLayout(ScreensLayout{screens: []ScreenDef{
			{id: 1, w: 1200, h: 1200, s: 1},
			{id: 2, w: 1200, h: 1200, s: 1.5, parent: ScreenDefParent{id: 1, align: "r", offset: 600}},
		}})
		err := sm.LayoutScreens(layout.screens)
		is.NoErr(err)

		screens := sm.Screens()
		is.Equal(len(screens), 2)                                                                         // 2 screens
		is.Equal(screens[0].PhysicalBounds, application.Rect{X: 0, Y: 0, Width: 1200, Height: 1200})      // Parent physical bounds
		is.Equal(screens[0].Bounds, screens[0].PhysicalBounds)                                            // Parent no scaling
		is.Equal(screens[1].PhysicalBounds, application.Rect{X: 1200, Y: 600, Width: 1200, Height: 1200}) // Child physical bounds
		is.Equal(screens[1].Bounds, application.Rect{X: 1200, Y: 600, Width: 800, Height: 800})           // Child DIP bounds
	})

	t.Run("Parent scaled", func(t *testing.T) {
		is := is.New(t)

		layout := parseLayout(ScreensLayout{screens: []ScreenDef{
			{id: 1, w: 1200, h: 1200, s: 1.5},
			{id: 2, w: 1200, h: 1200, s: 1, parent: ScreenDefParent{id: 1, align: "r", offset: 600}},
		}})

		err := sm.LayoutScreens(layout.screens)
		is.NoErr(err)

		screens := sm.Screens()
		is.Equal(len(screens), 2)                                                                         // 2 screens
		is.Equal(screens[0].PhysicalBounds, application.Rect{X: 0, Y: 0, Width: 1200, Height: 1200})      // Parent physical bounds
		is.Equal(screens[0].Bounds, application.Rect{X: 0, Y: 0, Width: 800, Height: 800})                // Parent DIP bounds
		is.Equal(screens[1].PhysicalBounds, application.Rect{X: 1200, Y: 600, Width: 1200, Height: 1200}) // Child physical bounds
		is.Equal(screens[1].Bounds, application.Rect{X: 800, Y: 400, Width: 1200, Height: 1200})          // Child DIP bounds
	})
}

// Test basic transformation between physical and DIP coordinates
func TestScreenManager_BasicTranformation(t *testing.T) {
	sm := application.ScreenManager{}
	is := is.New(t)
	layout := parseLayout(ScreensLayout{screens: []ScreenDef{
		{id: 1, w: 1200, h: 1200, s: 1},
		{id: 2, w: 1200, h: 1200, s: 1.5, parent: ScreenDefParent{id: 1, align: "r", offset: 600}},
	}})

	err := sm.LayoutScreens(layout.screens)
	is.NoErr(err)

	pt := application.Point{X: 100, Y: 100}
	is.Equal(sm.DipToPhysicalPoint(pt), pt) // DipToPhysicalPoint screen1
	is.Equal(sm.PhysicalToDipPoint(pt), pt) // PhysicalToDipPoint screen1

	ptDip := application.Point{X: 1300, Y: 700}
	ptPhysical := application.Point{X: 1350, Y: 750}
	is.Equal(sm.DipToPhysicalPoint(ptDip), ptPhysical) // DipToPhysicalPoint screen2
	is.Equal(sm.PhysicalToDipPoint(ptPhysical), ptDip) // PhysicalToDipPoint screen2

	rect := application.Rect{X: 100, Y: 100, Width: 200, Height: 300}
	is.Equal(sm.DipToPhysicalRect(rect), rect) // DipToPhysicalRect screen1
	is.Equal(sm.PhysicalToDipRect(rect), rect) // DipToPhysicalRect screen1

	rectDip := application.Rect{X: 1300, Y: 700, Width: 200, Height: 300}
	rectPhysical := application.Rect{X: 1350, Y: 750, Width: 300, Height: 450}
	is.Equal(sm.DipToPhysicalRect(rectDip), rectPhysical) // DipToPhysicalRect screen2
	is.Equal(sm.PhysicalToDipRect(rectPhysical), rectDip) // DipToPhysicalRect screen2

	rectDip = application.Rect{X: 2200, Y: 250, Width: 200, Height: 300}
	rectPhysical = application.Rect{X: 2700, Y: 75, Width: 300, Height: 450}
	is.Equal(sm.DipToPhysicalRect(rectDip), rectPhysical) // DipToPhysicalRect outside screen2
	is.Equal(sm.PhysicalToDipRect(rectPhysical), rectDip) // DipToPhysicalRect outside screen2
}

func TestScreenManager_PrimaryScreen(t *testing.T) {
	sm := application.ScreenManager{}
	is := is.New(t)

	for _, layout := range exampleLayouts() {
		err := sm.LayoutScreens(layout.screens)
		is.NoErr(err)
		is.Equal(sm.PrimaryScreen(), layout.screens[0]) // Primary screen
	}

	layout := parseLayout(ScreensLayout{screens: []ScreenDef{
		{id: 1, w: 1200, h: 1200, s: 1.5},
		{id: 2, w: 1200, h: 1200, s: 1, parent: ScreenDefParent{id: 1, align: "r", offset: 600}},
	}})

	layout.screens[0], layout.screens[1] = layout.screens[1], layout.screens[0]
	err := sm.LayoutScreens(layout.screens)
	is.NoErr(err)
	is.Equal(sm.PrimaryScreen(), layout.screens[1]) // Primary screen

	layout.screens[1].IsPrimary = false
	err = sm.LayoutScreens(layout.screens)
	is.True(err != nil) // Should error when no primary screen found
}

// Test edge alignment between transformation
// (points and rects on the screen edge should transform to the same precise edge position)
func TestScreenManager_EdgeAlign(t *testing.T) {
	sm := application.ScreenManager{}
	is := is.New(t)

	for _, layout := range exampleLayouts() {
		err := sm.LayoutScreens(layout.screens)
		is.NoErr(err)
		for _, screen := range sm.Screens() {
			ptOriginDip := screen.Bounds.Origin()
			ptOriginPhysical := screen.PhysicalBounds.Origin()
			ptCornerDip := screen.Bounds.InsideCorner()
			ptCornerPhysical := screen.PhysicalBounds.InsideCorner()

			is.Equal(sm.DipToPhysicalPoint(ptOriginDip), ptOriginPhysical) // DipToPhysicalPoint Origin
			is.Equal(sm.PhysicalToDipPoint(ptOriginPhysical), ptOriginDip) // PhysicalToDipPoint Origin
			is.Equal(sm.DipToPhysicalPoint(ptCornerDip), ptCornerPhysical) // DipToPhysicalPoint Corner
			is.Equal(sm.PhysicalToDipPoint(ptCornerPhysical), ptCornerDip) // PhysicalToDipPoint Corner

			rectOriginDip := application.Rect{X: ptOriginDip.X, Y: ptOriginDip.Y, Width: 100, Height: 100}
			rectOriginPhysical := application.Rect{X: ptOriginPhysical.X, Y: ptOriginPhysical.Y, Width: 100, Height: 100}
			rectCornerDip := application.Rect{X: ptCornerDip.X - 99, Y: ptCornerDip.Y - 99, Width: 100, Height: 100}
			rectCornerPhysical := application.Rect{X: ptCornerPhysical.X - 99, Y: ptCornerPhysical.Y - 99, Width: 100, Height: 100}

			is.Equal(sm.DipToPhysicalRect(rectOriginDip).Origin(), rectOriginPhysical.Origin()) // DipToPhysicalRect Origin
			is.Equal(sm.PhysicalToDipRect(rectOriginPhysical).Origin(), rectOriginDip.Origin()) // PhysicalToDipRect Origin
			is.Equal(sm.DipToPhysicalRect(rectCornerDip).Corner(), rectCornerPhysical.Corner()) // DipToPhysicalRect Corner
			is.Equal(sm.PhysicalToDipRect(rectCornerPhysical).Corner(), rectCornerDip.Corner()) // PhysicalToDipRect Corner
		}
	}
}

func TestScreenManager_ProbePoints(t *testing.T) {
	sm := application.ScreenManager{}
	is := is.New(t)
	threshold := 1.0
	steps := 3

	for _, layout := range exampleLayouts() {
		err := sm.LayoutScreens(layout.screens)
		is.NoErr(err)
		for _, screen := range sm.Screens() {
			for i := 0; i <= 1; i++ {
				isDip := (i == 0)

				var b application.Rect
				if isDip {
					b = screen.Bounds
				} else {
					b = screen.PhysicalBounds
				}

				xStep := b.Width / steps
				yStep := b.Height / steps
				if xStep < 1 {
					xStep = 1
				}
				if yStep < 1 {
					yStep = 1
				}
				pt := b.Origin()
				xDone := false
				yDone := false

				for !yDone {
					if pt.Y > b.InsideCorner().Y {
						pt.Y = b.InsideCorner().Y
						yDone = true
					}

					pt.X = b.X
					xDone = false

					for !xDone {
						if pt.X > b.InsideCorner().X {
							pt.X = b.InsideCorner().X
							xDone = true
						}
						var ptDblTransformed application.Point

						if isDip {
							ptDblTransformed = sm.PhysicalToDipPoint(sm.DipToPhysicalPoint(pt))
						} else {
							ptDblTransformed = sm.DipToPhysicalPoint(sm.PhysicalToDipPoint(pt))
						}

						is.True(math.Abs(float64(ptDblTransformed.X-pt.X)) <= threshold)
						is.True(math.Abs(float64(ptDblTransformed.Y-pt.Y)) <= threshold)
						pt.X += xStep
					}
					pt.Y += yStep
				}
			}
		}
	}
}

// Test transformation drift over time
func TestScreenManager_TransformationDrift(t *testing.T) {
	sm := application.ScreenManager{}
	is := is.New(t)

	for _, layout := range exampleLayouts() {
		err := sm.LayoutScreens(layout.screens)
		is.NoErr(err)
		for _, screen := range sm.Screens() {
			rectPhysicalOriginal := application.Rect{
				X:      screen.PhysicalBounds.X + 100,
				Y:      screen.PhysicalBounds.Y + 100,
				Width:  123,
				Height: 123,
			}

			// Slide the position to catch any rounding errors
			for i := 0; i < 10; i++ {
				rectPhysicalOriginal.X++
				rectPhysicalOriginal.Y++
				rectPhysical := rectPhysicalOriginal
				// Transform back and forth several times to make sure no drift is introduced over time
				for j := 0; j < 10; j++ {
					rectDip := sm.PhysicalToDipRect(rectPhysical)
					rectPhysical = sm.DipToPhysicalRect(rectDip)
				}
				is.NoErr(matchRects(rectPhysical, rectPhysicalOriginal))
			}
		}
	}
}

func TestScreenManager_ScreenNearestRect(t *testing.T) {
	sm := application.ScreenManager{}
	is := is.New(t)

	layout := parseLayout(ScreensLayout{screens: []ScreenDef{
		{id: 1, w: 3840, h: 2160, s: 163.0 / 96, name: `27" 4K UHD 163DPI`},
		{id: 2, w: 1920, h: 1080, s: 1, parent: ScreenDefParent{id: 1, align: "r", offset: 0}, name: `23" FHD 96DPI`},
		{id: 3, w: 1920, h: 1080, s: 1.25, parent: ScreenDefParent{id: 1, align: "l", offset: 0}, name: `23" FHD 96DPI (125%)`},
	}})
	err := sm.LayoutScreens(layout.screens)
	is.NoErr(err)

	type Rects map[string][]application.Rect

	t.Run("DIP rects", func(t *testing.T) {
		is := is.New(t)
		rects := Rects{
			"1": []application.Rect{
				{X: -150, Y: 260, Width: 400, Height: 300},
				{X: -250, Y: 750, Width: 400, Height: 300},
				{X: -450, Y: 950, Width: 400, Height: 300},
				{X: 800, Y: 1350, Width: 400, Height: 300},
				{X: 2000, Y: 100, Width: 400, Height: 300},
				{X: 2100, Y: 950, Width: 400, Height: 300},
				{X: 2350, Y: 1200, Width: 400, Height: 300},
			},
			"2": []application.Rect{
				{X: 2100, Y: 50, Width: 400, Height: 300},
				{X: 2150, Y: 950, Width: 400, Height: 300},
				{X: 2450, Y: 1150, Width: 400, Height: 300},
				{X: 4300, Y: 400, Width: 400, Height: 300},
			},
			"3": []application.Rect{
				{X: -2000, Y: 100, Width: 400, Height: 300},
				{X: -220, Y: 200, Width: 400, Height: 300},
				{X: -300, Y: 750, Width: 400, Height: 300},
				{X: -500, Y: 900, Width: 400, Height: 300},
			},
		}

		for screenID, screenRects := range rects {
			for _, rect := range screenRects {
				screen := sm.ScreenNearestDipRect(rect)
				is.Equal(screen.ID, screenID)
			}
		}
	})
	t.Run("Physical rects", func(t *testing.T) {
		is := is.New(t)
		rects := Rects{
			"1": []application.Rect{
				{X: -150, Y: 100, Width: 400, Height: 300},
				{X: -250, Y: 1500, Width: 400, Height: 300},
				{X: 3600, Y: 100, Width: 400, Height: 300},
			},
			"2": []application.Rect{
				{X: 3700, Y: 100, Width: 400, Height: 300},
				{X: 4000, Y: 1150, Width: 400, Height: 300},
			},
			"3": []application.Rect{
				{X: -250, Y: 100, Width: 400, Height: 300},
				{X: -300, Y: 950, Width: 400, Height: 300},
				{X: -1000, Y: 1000, Width: 400, Height: 300},
			},
		}

		for screenID, screenRects := range rects {
			for _, rect := range screenRects {
				screen := sm.ScreenNearestPhysicalRect(rect)
				is.Equal(screen.ID, screenID)
			}
		}
	})

	// DIP rect is near screen1 but when transformed becomes near screen2.
	// To have a consistent transformation back & forth, screen nearest physical rect
	// should be the one given by ScreenNearestDipRect
	t.Run("Edge case 1", func(t *testing.T) {
		is := is.New(t)
		layout := parseLayout(ScreensLayout{screens: []ScreenDef{
			{id: 1, w: 1200, h: 1200, s: 1},
			{id: 2, w: 1200, h: 1300, s: 1.5, parent: ScreenDefParent{id: 1, align: "r", offset: -20}},
		}})
		err := sm.LayoutScreens(layout.screens)
		is.NoErr(err)

		rectDip := application.Rect{X: 1020, Y: 800, Width: 400, Height: 300}
		rectPhysical := sm.DipToPhysicalRect(rectDip)

		screenDip := sm.ScreenNearestDipRect(rectDip)
		screenPhysical := sm.ScreenNearestPhysicalRect(rectPhysical)
		is.Equal(screenDip.ID, "2")      // screenDip
		is.Equal(screenPhysical.ID, "2") // screenPhysical

		rectDblTransformed := sm.PhysicalToDipRect(rectPhysical)
		is.NoErr(matchRects(rectDblTransformed, rectDip)) // double transformation
	})
}

// Unsolved edge cases
func TestScreenManager_UnsolvedEdgeCases(t *testing.T) {
	sm := application.ScreenManager{}
	is := is.New(t)

	// Edge case 1: invalid DIP rect location
	// there could be a setup where some dip rects locations are invalid, meaning that there's no
	// physical rect that could produce that dip rect at this location
	// Not sure how to solve this scenario
	t.Run("Edge case 1: invalid dip rect", func(t *testing.T) {
		t.Skip("Unsolved edge case")
		is := is.New(t)
		layout := parseLayout(ScreensLayout{screens: []ScreenDef{
			{id: 1, w: 1200, h: 1200, s: 1},
			{id: 2, w: 1200, h: 1100, s: 1.5, parent: ScreenDefParent{id: 1, align: "r", offset: 0}},
		}})
		err := sm.LayoutScreens(layout.screens)
		is.NoErr(err)

		rectDip := application.Rect{X: 1050, Y: 700, Width: 400, Height: 300}
		rectPhysical := sm.DipToPhysicalRect(rectDip)

		screenDip := sm.ScreenNearestDipRect(rectDip)
		screenPhysical := sm.ScreenNearestPhysicalRect(rectPhysical)
		is.Equal(screenDip.ID, screenPhysical.ID)

		rectDblTransformed := sm.PhysicalToDipRect(rectPhysical)
		is.NoErr(matchRects(rectDblTransformed, rectDip)) // double transformation
	})

	// Edge case 2: physical rect that changes when double transformed
	// there could be a setup where a dip rect at some locations could be produced by two different physical rects
	// causing one of these physical rects to be changed to the other when double transformed
	// Not sure how to solve this scenario
	t.Run("Edge case 2: changed physical rect", func(t *testing.T) {
		t.Skip("Unsolved edge case")
		is := is.New(t)
		layout := parseLayout(ScreensLayout{screens: []ScreenDef{
			{id: 1, w: 1200, h: 1200, s: 1.5},
			{id: 2, w: 1200, h: 900, s: 1, parent: ScreenDefParent{id: 1, align: "r", offset: 0}},
		}})
		err := sm.LayoutScreens(layout.screens)
		is.NoErr(err)

		rectPhysical := application.Rect{X: 1050, Y: 890, Width: 400, Height: 300}
		rectDblTransformed := sm.DipToPhysicalRect(sm.PhysicalToDipRect(rectPhysical))
		is.NoErr(matchRects(rectDblTransformed, rectPhysical)) // double transformation
	})
}

func BenchmarkScreenManager_LayoutScreens(b *testing.B) {
	sm := application.ScreenManager{}
	layouts := exampleLayouts()
	screens := layouts[3].screens

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sm.LayoutScreens(screens)
	}
}

func BenchmarkScreenManager_TransformPoint(b *testing.B) {
	sm := application.ScreenManager{}
	layouts := exampleLayouts()
	screens := layouts[3].screens
	sm.LayoutScreens(screens)

	pt := application.Point{X: 500, Y: 500}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sm.DipToPhysicalPoint(pt)
	}
}

func BenchmarkScreenManager_TransformRect(b *testing.B) {
	sm := application.ScreenManager{}
	layouts := exampleLayouts()
	screens := layouts[3].screens
	sm.LayoutScreens(screens)

	rect := application.Rect{X: 500, Y: 500, Width: 800, Height: 600}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sm.DipToPhysicalRect(rect)
	}
}
