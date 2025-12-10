//go:build ios

package application

// getScreens returns all screens on iOS - Screen type is defined in screenmanager.go

// getScreens returns all screens on iOS
func getScreens() ([]*Screen, error) {
	// iOS typically has one screen
	// This would need proper implementation with UIScreen
	mainRect := Rect{
		X:      0,
		Y:      0,
		Width:  1170,  // iPhone 12 Pro width
		Height: 2532,  // iPhone 12 Pro height
	}
	return []*Screen{
		{
			ID:               "main",
			Name:             "Main Screen",
			ScaleFactor:      3.0,  // iPhone 12 Pro scale
			X:                0,
			Y:                0,
			Size:             Size{Width: 1170, Height: 2532},
			Bounds:           mainRect,
			PhysicalBounds:   mainRect,
			WorkArea:         mainRect,
			PhysicalWorkArea: mainRect,
			IsPrimary:        true,
			Rotation:         0,
		},
	}, nil
}