//go:build openharmony && !server

package application

import "encoding/json"

// Screen information is provided by the WailsBridge (getScreenInfoJson).
// Sizes are reported in density-independent pixels (dp), which match the
// ArkWeb CSS pixels; PhysicalBounds is in hardware pixels. WorkArea is the
// display minus the system bar insets.

type openharmonyScreenInfo struct {
	WidthPx     int     `json:"widthPx"`
	HeightPx    int     `json:"heightPx"`
	Density     float64 `json:"density"`
	InsetTop    int     `json:"insetTop"`
	InsetBottom int     `json:"insetBottom"`
	InsetLeft   int     `json:"insetLeft"`
	InsetRight  int     `json:"insetRight"`
}

// getScreens returns the device display. Multi-display setups are not
// supported: the ArkWeb always occupies the default display.
func getScreens() ([]*Screen, error) {
	jsonStr, ok := ohBridgeString("getScreenInfoJson", "{}")
	if !ok || jsonStr == "" {
		return fallbackScreens(), nil
	}

	var info openharmonyScreenInfo
	if err := json.Unmarshal([]byte(jsonStr), &info); err != nil || info.WidthPx <= 0 || info.HeightPx <= 0 {
		return fallbackScreens(), nil
	}

	scale := info.Density
	if scale <= 0 {
		scale = 1
	}

	dp := func(px int) int { return int(float64(px) / scale) }

	width := dp(info.WidthPx)
	height := dp(info.HeightPx)

	physicalBounds := Rect{X: 0, Y: 0, Width: info.WidthPx, Height: info.HeightPx}
	physicalWorkArea := Rect{
		X:      info.InsetLeft,
		Y:      info.InsetTop,
		Width:  info.WidthPx - info.InsetLeft - info.InsetRight,
		Height: info.HeightPx - info.InsetTop - info.InsetBottom,
	}

	workArea := Rect{
		X:      dp(info.InsetLeft),
		Y:      dp(info.InsetTop),
		Width:  width - dp(info.InsetLeft) - dp(info.InsetRight),
		Height: height - dp(info.InsetTop) - dp(info.InsetBottom),
	}

	return []*Screen{
		{
			ID:               "main",
			Name:             "Main Display",
			IsPrimary:        true,
			ScaleFactor:      float32(scale),
			Size:             Size{Width: width, Height: height},
			Bounds:           Rect{X: 0, Y: 0, Width: width, Height: height},
			PhysicalBounds:   physicalBounds,
			WorkArea:         workArea,
			PhysicalWorkArea: physicalWorkArea,
		},
	}, nil
}

func fallbackScreens() []*Screen {
	// Sensible defaults (dp) when the bridge isn't available yet
	return []*Screen{
		{
			ID:               "main",
			Name:             "Main Display",
			IsPrimary:        true,
			ScaleFactor:      2.75,
			Size:             Size{Width: 393, Height: 873},
			Bounds:           Rect{X: 0, Y: 0, Width: 393, Height: 873},
			PhysicalBounds:   Rect{X: 0, Y: 0, Width: 1080, Height: 2400},
			WorkArea:         Rect{X: 0, Y: 24, Width: 393, Height: 825},
			PhysicalWorkArea: Rect{X: 0, Y: 66, Width: 1080, Height: 2268},
		},
	}
}
