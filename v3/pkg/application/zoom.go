package application

import "math"

// zoomFactor validates a zoom factor requested through SetZoom and returns the
// factor the platform should apply. Non-positive, NaN and Inf factors are
// rejected (ok == false): WebView2 returns an error for them and WKWebView
// raises NSInvalidArgumentException. Frameless windows are kept at 1.0 or
// above because their drag and resize regions are not adjusted for zoom
// (https://github.com/wailsapp/wails/issues/4590).
func zoomFactor(zoom float64, frameless bool) (factor float64, ok bool) {
	if zoom <= 0 || math.IsNaN(zoom) || math.IsInf(zoom, 0) {
		return 0, false
	}
	if frameless && zoom < 1 {
		return 1, true
	}
	return zoom, true
}
