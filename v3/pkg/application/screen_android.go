//go:build android

package application

// getScreens returns the available screens for Android
func getScreens() ([]*Screen, error) {
	// Android typically has one main display
	// TODO: Support for multi-display via DisplayManager
	return []*Screen{
		{
			ID:        "main",
			Name:      "Main Display",
			IsPrimary: true,
			Size: Size{
				Width:  1080,
				Height: 2400,
			},
			Bounds: Rect{
				X:      0,
				Y:      0,
				Width:  1080,
				Height: 2400,
			},
			WorkArea: Rect{
				X:      0,
				Y:      0,
				Width:  1080,
				Height: 2340, // Minus navigation bar
			},
		},
	}, nil
}
