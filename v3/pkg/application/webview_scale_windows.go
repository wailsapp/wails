//go:build windows

package application

import (
	"golang.org/x/sys/windows/registry"
)

// This file owns the "what SHOULD the WebView2 rasterization scale be" side
// of the mixed-DPI work (#5701). The storm/settle state machine lives in
// webview_window_windows.go; the pure decision helpers live in
// dpiflap_settle.go so mac/linux CI can pin them.

// windowsTextScaleFactor reads the Windows accessibility text scale
// (Settings → Accessibility → Text size) as a factor: 1.0 for 100%, 2.25 for
// 225%. WebView2 defines RasterizationScale as monitor DPI scale × text
// scale, and native monitor-scale detection folds text scaling in — so every
// target the app computes must include it as well, or text-scaling users get
// mis-flagged as mismatched and "corrected" smaller (the be8e16f7e
// settle-guard defect). The value lives at
// HKCU\SOFTWARE\Microsoft\Accessibility!TextScaleFactor (100–225) and is
// simply absent at the 100% default; any read failure means 1.0. A registry
// read is microseconds and needs no COM, so callers read it fresh each time —
// no cache to invalidate on WM_SETTINGCHANGE.
func windowsTextScaleFactor() float64 {
	key, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Accessibility`, registry.QUERY_VALUE)
	if err != nil {
		return 1.0
	}
	defer key.Close()
	value, _, err := key.GetIntegerValue("TextScaleFactor")
	if err != nil || value == 0 {
		return 1.0
	}
	return float64(value) / 100.0
}

// targetRasterizationScale is the window's authoritative "what the
// controller's rasterization scale should be right now" for the given DPI:
// dpi/96 × text scale (see rasterizationTargetForDPI for the formula and why
// text scale is non-optional). 0 means the DPI was unreadable.
func (w *windowsWebviewWindow) targetRasterizationScale(dpi uint32) float64 {
	return rasterizationTargetForDPI(dpi, windowsTextScaleFactor())
}

// configureWebviewScaleOwnership selects who writes the controller's
// rasterization scale for this window's controller lifetime and emits the
// one-per-controller ownership breadcrumb. Called from setupChromium right
// after a successful Embed (both the initial run() and every recovery
// rebuild — webviewRebuildInProgress distinguishes them in the breadcrumb).
//
// The breadcrumb ships at warning level because it is the field decoder ring
// for every other DPI line in the session: it records whether the
// RasterizationScaleChanged event stream is even registered (a silent
// registration failure otherwise reads as "no scale changes happened") and
// the text scale factor in effect at startup (#5701).
func (w *windowsWebviewWindow) configureWebviewScaleOwnership() {
	reason := "embed"
	if w.webviewRebuildInProgress {
		reason = "rebuild"
	}
	detectionStatus := w.enableNativeMonitorScaleDetection()
	regStatus := "ok"
	if ok, err := w.chromium.RasterizationScaleEventRegistration(); !ok {
		regStatus = "FAILED: " + errString(err)
	}
	globalApplication.warning(
		"WebView2 scale ownership [%s]: window %d owner=webview2 AddRasterizationScaleChanged=%s ShouldDetectMonitorScaleChanges=%s textScale %.2f (#5701)",
		reason, w.parent.id, regStatus, detectionStatus, windowsTextScaleFactor())
}

// errString renders an error for a breadcrumb without panicking on nil (a
// failed registration always carries an error, but the breadcrumb must never
// be the thing that crashes).
func errString(err error) string {
	if err == nil {
		return "unknown"
	}
	return err.Error()
}
