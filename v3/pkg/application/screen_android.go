//go:build android

package application

import (
	"encoding/json"
)

// ScreenInfo represents the JSON structure returned from Java
type ScreenInfo struct {
	WidthPixels   int     `json:"widthPixels"`
	HeightPixels  int     `json:"heightPixels"`
	Density       float32 `json:"density"`
	DensityDpi    int     `json:"densityDpi"`
	ScaledDensity float32 `json:"scaledDensity"`
	Xdpi          float32 `json:"xdpi"`
	Ydpi          float32 `json:"ydpi"`
}

// getScreens returns the available screens for Android
func getScreens() ([]*Screen, error) {
	// Get screen info from Android via JNI
	screenInfoJSON := AndroidGetScreenInfo()

	var screenInfo ScreenInfo
	err := json.Unmarshal([]byte(screenInfoJSON), &screenInfo)
	if err != nil {
		// Fallback to hardcoded values if JSON parsing fails
		androidLogf("error", "Failed to parse screen info JSON: %v", err)
		return getFallbackScreens(), nil
	}

	// Calculate scale factor (DPI / 160, which is Android's baseline DPI)
	scaleFactor := float32(screenInfo.DensityDpi) / 160.0

	// Estimate work area (subtract status bar and navigation bar)
	// Typical status bar is ~24dp and nav bar is ~48dp at 160dpi
	statusBarPixels := int(24 * screenInfo.Density)
	navBarPixels := int(48 * screenInfo.Density)
	workAreaHeight := screenInfo.HeightPixels - statusBarPixels - navBarPixels

	// Android typically has one main display
	// TODO: Support for multi-display via DisplayManager
	return []*Screen{
		{
			ID:          "main",
			Name:        "Main Display",
			IsPrimary:   true,
			ScaleFactor: scaleFactor,
			X:           0,
			Y:           0,
			Size: Size{
				Width:  screenInfo.WidthPixels,
				Height: screenInfo.HeightPixels,
			},
			Bounds: Rect{
				X:      0,
				Y:      0,
				Width:  screenInfo.WidthPixels,
				Height: screenInfo.HeightPixels,
			},
			WorkArea: Rect{
				X:      0,
				Y:      statusBarPixels,
				Width:  screenInfo.WidthPixels,
				Height: workAreaHeight,
			},
		},
	}, nil
}

// getFallbackScreens returns hardcoded screen values as a fallback
func getFallbackScreens() []*Screen {
	return []*Screen{
		{
			ID:          "main",
			Name:        "Main Display",
			IsPrimary:   true,
			ScaleFactor: 2.0,
			X:           0,
			Y:           0,
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
	}
}
