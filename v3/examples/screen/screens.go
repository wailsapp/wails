package main

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

type ScreenService struct {
	screenManager   application.ScreenManager
	isExampleLayout bool
}

// Helper to safely get float64 from interface{}
func getFloat64(v interface{}) float64 {
	if v == nil {
		return 0
	}
	if f, ok := v.(float64); ok {
		return f
	}
	return 0
}

// Helper to safely get int from interface{} (expecting float64 from JSON)
func getInt(v interface{}) int {
	return int(getFloat64(v))
}

func (s *ScreenService) GetSystemScreens() []*application.Screen {
	s.isExampleLayout = false
	screens := application.Get().Screen.GetAll()
	return screens
}

func (s *ScreenService) ProcessExampleScreens(rawScreens []interface{}) []*application.Screen {
	s.isExampleLayout = true

	parseRect := func(m map[string]interface{}) application.Rect {
		if m == nil {
			return application.Rect{}
		}
		return application.Rect{
			X:      getInt(m["X"]),
			Y:      getInt(m["Y"]),
			Width:  getInt(m["Width"]),
			Height: getInt(m["Height"]),
		}
	}

	// Prevent unbounded slice growth by limiting the number of screens
	maxScreens := 32 // Reasonable limit for screen configurations
	if len(rawScreens) > maxScreens {
		rawScreens = rawScreens[:maxScreens]
	}

	screens := make([]*application.Screen, 0, len(rawScreens))
	for _, s := range rawScreens {
		sm, ok := s.(map[string]interface{})
		if !ok {
			continue
		}

		boundsVal, ok := sm["Bounds"]
		if !ok {
			continue
		}
		bounds := parseRect(boundsVal.(map[string]interface{}))

		var id, name string
		var isPrimary bool
		if idVal, ok := sm["ID"].(string); ok {
			id = idVal
		}
		if nameVal, ok := sm["Name"].(string); ok {
			name = nameVal
		}
		if primaryVal, ok := sm["IsPrimary"].(bool); ok {
			isPrimary = primaryVal
		}

		var physicalBounds, workArea, physicalWorkArea application.Rect
		if pb, ok := sm["PhysicalBounds"].(map[string]interface{}); ok {
			physicalBounds = parseRect(pb)
		}
		if wa, ok := sm["WorkArea"].(map[string]interface{}); ok {
			workArea = parseRect(wa)
		}
		if pwa, ok := sm["PhysicalWorkArea"].(map[string]interface{}); ok {
			physicalWorkArea = parseRect(pwa)
		}

		screens = append(screens, &application.Screen{
			ID:               id,
			Name:             name,
			X:                bounds.X,
			Y:                bounds.Y,
			Size:             application.Size{Width: bounds.Width, Height: bounds.Height},
			Bounds:           bounds,
			PhysicalBounds:   physicalBounds,
			WorkArea:         workArea,
			PhysicalWorkArea: physicalWorkArea,
			IsPrimary:        isPrimary,
			ScaleFactor:      float32(getFloat64(sm["ScaleFactor"])),
			Rotation:         0,
		})
	}

	s.screenManager.LayoutScreens(screens)
	return s.screenManager.GetAll()
}

func (s *ScreenService) transformPoint(point application.Point, toDIP bool) application.Point {
	if s.isExampleLayout {
		if toDIP {
			return s.screenManager.PhysicalToDipPoint(point)
		} else {
			return s.screenManager.DipToPhysicalPoint(point)
		}
	} else {
		if toDIP {
			return application.PhysicalToDipPoint(point)
		} else {
			return application.DipToPhysicalPoint(point)
		}
	}
}

func (s *ScreenService) TransformPoint(point map[string]interface{}, toDIP bool) (points [2]application.Point) {
	if point == nil {
		return points
	}

	pt := application.Point{
		X: getInt(point["X"]),
		Y: getInt(point["Y"]),
	}

	ptTransformed := s.transformPoint(pt, toDIP)
	ptDblTransformed := s.transformPoint(ptTransformed, !toDIP)

	// double-transform a limited number of times to catch any double-rounding issues
	// Limit iterations to prevent potential performance issues
	maxIterations := 3 // Reduced from 10 to limit computational overhead
	for i := 0; i < maxIterations; i++ {
		ptTransformed = s.transformPoint(ptDblTransformed, toDIP)
		ptDblTransformed = s.transformPoint(ptTransformed, !toDIP)
	}

	points[0] = ptTransformed
	points[1] = ptDblTransformed
	return points
}

func (s *ScreenService) TransformRect(rect map[string]interface{}, toDIP bool) application.Rect {
	if rect == nil {
		return application.Rect{}
	}

	r := application.Rect{
		X:      getInt(rect["X"]),
		Y:      getInt(rect["Y"]),
		Width:  getInt(rect["Width"]),
		Height: getInt(rect["Height"]),
	}

	if s.isExampleLayout {
		if toDIP {
			return s.screenManager.PhysicalToDipRect(r)
		} else {
			return s.screenManager.DipToPhysicalRect(r)
		}
	} else {
		if toDIP {
			return application.PhysicalToDipRect(r)
		} else {
			return application.DipToPhysicalRect(r)
		}
	}
}
