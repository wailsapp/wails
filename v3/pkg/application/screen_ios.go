//go:build ios

package application

/*
#include <stdlib.h>
#include "application_ios.h"
*/
import "C"

import (
	"encoding/json"
	"unsafe"
)

type iosScreenInfo struct {
	PointWidth  int     `json:"pointWidth"`
	PointHeight int     `json:"pointHeight"`
	PixelWidth  int     `json:"pixelWidth"`
	PixelHeight int     `json:"pixelHeight"`
	Scale       float32 `json:"scale"`
	SafeTop     int     `json:"safeTop"`
	SafeBottom  int     `json:"safeBottom"`
	SafeLeft    int     `json:"safeLeft"`
	SafeRight   int     `json:"safeRight"`
}

// getScreens returns the device's screen via UIScreen. iOS has exactly one
// screen from the app's perspective; sizes are reported in points (logical
// pixels), with PhysicalBounds carrying the native pixel size. The work area
// is the screen inset by the current safe area (notch, home indicator).
func getScreens() ([]*Screen, error) {
	cinfo := C.ios_screen_info_json()
	if cinfo == nil {
		return fallbackScreens(), nil
	}
	defer C.free(unsafe.Pointer(cinfo))

	var info iosScreenInfo
	if err := json.Unmarshal([]byte(C.GoString(cinfo)), &info); err != nil || info.PointWidth == 0 {
		return fallbackScreens(), nil
	}

	bounds := Rect{X: 0, Y: 0, Width: info.PointWidth, Height: info.PointHeight}
	physical := Rect{X: 0, Y: 0, Width: info.PixelWidth, Height: info.PixelHeight}
	workArea := Rect{
		X:      info.SafeLeft,
		Y:      info.SafeTop,
		Width:  info.PointWidth - info.SafeLeft - info.SafeRight,
		Height: info.PointHeight - info.SafeTop - info.SafeBottom,
	}
	scale := info.Scale
	if scale == 0 {
		scale = 1
	}
	physicalWorkArea := Rect{
		X:      int(float32(workArea.X) * scale),
		Y:      int(float32(workArea.Y) * scale),
		Width:  int(float32(workArea.Width) * scale),
		Height: int(float32(workArea.Height) * scale),
	}

	return []*Screen{
		{
			ID:               "main",
			Name:             "Main Screen",
			ScaleFactor:      scale,
			X:                0,
			Y:                0,
			Size:             Size{Width: bounds.Width, Height: bounds.Height},
			Bounds:           bounds,
			PhysicalBounds:   physical,
			WorkArea:         workArea,
			PhysicalWorkArea: physicalWorkArea,
			IsPrimary:        true,
			Rotation:         0,
		},
	}, nil
}

// fallbackScreens is used if the native screen query fails (e.g. before
// UIKit is fully up). Callers get a sane portrait phone shape rather than an
// error so layout code can proceed.
func fallbackScreens() []*Screen {
	mainRect := Rect{X: 0, Y: 0, Width: 390, Height: 844}
	return []*Screen{
		{
			ID:               "main",
			Name:             "Main Screen",
			ScaleFactor:      3.0,
			Size:             Size{Width: mainRect.Width, Height: mainRect.Height},
			Bounds:           mainRect,
			PhysicalBounds:   Rect{Width: mainRect.Width * 3, Height: mainRect.Height * 3},
			WorkArea:         mainRect,
			PhysicalWorkArea: Rect{Width: mainRect.Width * 3, Height: mainRect.Height * 3},
			IsPrimary:        true,
		},
	}
}
