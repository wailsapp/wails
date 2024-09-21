package main

import (
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type ScreenService struct {
	screenManager   application.ScreenManager
	isExampleLayout bool
}

func (s *ScreenService) GetSystemScreens() []*application.Screen {
	s.isExampleLayout = false
	screens, _ := application.Get().GetScreens()
	return screens
}

func (s *ScreenService) ProcessExampleScreens(rawScreens []interface{}) []*application.Screen {
	s.isExampleLayout = true

	parseRect := func(m map[string]interface{}) application.Rect {
		return application.Rect{
			X:      int(m["X"].(float64)),
			Y:      int(m["Y"].(float64)),
			Width:  int(m["Width"].(float64)),
			Height: int(m["Height"].(float64)),
		}
	}

	screens := []*application.Screen{}
	for _, s := range rawScreens {
		s := s.(map[string]interface{})

		bounds := parseRect(s["Bounds"].(map[string]interface{}))

		screens = append(screens, &application.Screen{
			ID:               s["ID"].(string),
			Name:             s["Name"].(string),
			X:                bounds.X,
			Y:                bounds.Y,
			Size:             application.Size{Width: bounds.Width, Height: bounds.Height},
			Bounds:           bounds,
			PhysicalBounds:   parseRect(s["PhysicalBounds"].(map[string]interface{})),
			WorkArea:         parseRect(s["WorkArea"].(map[string]interface{})),
			PhysicalWorkArea: parseRect(s["PhysicalWorkArea"].(map[string]interface{})),
			IsPrimary:        s["IsPrimary"].(bool),
			ScaleFactor:      float32(s["ScaleFactor"].(float64)),
			Rotation:         0,
		})
	}

	s.screenManager.LayoutScreens(screens)
	return s.screenManager.Screens()
}

func (s *ScreenService) transformPoint(point application.Point, toDIP bool) application.Point {
	if s.isExampleLayout {
		if toDIP {
			return s.screenManager.PhysicalToDipPoint(point)
		} else {
			return s.screenManager.DipToPhysicalPoint(point)
		}
	} else {
		// =======================
		// TODO: remove this block when DPI is implemented in Linux & Mac
		if runtime.GOOS != "windows" {
			println("DPI not implemented yet!")
			return point
		}
		// =======================
		if toDIP {
			return application.PhysicalToDipPoint(point)
		} else {
			return application.DipToPhysicalPoint(point)
		}
	}
}

func (s *ScreenService) TransformPoint(point map[string]interface{}, toDIP bool) (points [2]application.Point) {
	pt := application.Point{
		X: int(point["X"].(float64)),
		Y: int(point["Y"].(float64)),
	}

	ptTransformed := s.transformPoint(pt, toDIP)
	ptDblTransformed := s.transformPoint(ptTransformed, !toDIP)

	// double-transform multiple times to catch any double-rounding issues
	for i := 0; i < 10; i++ {
		ptTransformed = s.transformPoint(ptDblTransformed, toDIP)
		ptDblTransformed = s.transformPoint(ptTransformed, !toDIP)
	}

	points[0] = ptTransformed
	points[1] = ptDblTransformed
	return points
}

func (s *ScreenService) TransformRect(rect map[string]interface{}, toDIP bool) application.Rect {
	r := application.Rect{
		X:      int(rect["X"].(float64)),
		Y:      int(rect["Y"].(float64)),
		Width:  int(rect["Width"].(float64)),
		Height: int(rect["Height"].(float64)),
	}

	if s.isExampleLayout {
		if toDIP {
			return s.screenManager.PhysicalToDipRect(r)
		} else {
			return s.screenManager.DipToPhysicalRect(r)
		}
	} else {
		// =======================
		// TODO: remove this block when DPI is implemented in Linux & Mac
		if runtime.GOOS != "windows" {
			println("DPI not implemented yet!")
			return r
		}
		// =======================
		if toDIP {
			return application.PhysicalToDipRect(r)
		} else {
			return application.DipToPhysicalRect(r)
		}
	}
}
