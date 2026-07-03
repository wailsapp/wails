package application

import "math"

// dpiFlapSettleScaleTolerance is the |raster - target| above which the settle
// treats native monitor-scale detection as having gone SILENT and
// force-corrects the rasterization scale itself. Native detection owns the
// scale during the storm, but the v200.0.23 field session logged ONE settle at
// raster 1.25 under dpi 216/target 2.25 that stayed stuck for 85s until a
// maximize/restore forced the change — no RasterizationScaleChanged event ever
// fired, so the per-event re-layout never ran. Any genuine DPI mismatch is a
// >=0.25 scale step (1.0 for a 120<->216 flip), so 0.01 catches every real
// wrong-scale state while never tripping on float jitter; the corrective put
// itself no-ops when syncWebviewRasterizationScale finds the scale already in
// sync (its own 0.001 tolerance), so the normal native-corrected case stays at
// zero puts.
//
// Declared here — not next to the other dpiFlap timing constants in
// webview_window_windows.go — because it backs a platform-independent decision
// helper, so it must compile on every OS (the regression test runs on mac/linux
// CI, not just Windows where the live controller exists).
const dpiFlapSettleScaleTolerance = 0.01

// rasterizationTargetForDPI is the single definition of the correct WebView2
// rasterization scale for a window DPI: dpi/96 × the Windows text scale
// factor. Native monitor-scale detection folds text scaling in (the
// RasterizationScale docs define the property as "the combination of the
// monitor DPI scale and text scaling set by the user"), so any target that
// compares against — or is written to — the controller MUST include it too.
// Comparing raster against a bare dpi/96 mis-flags every text-scaling user's
// perfectly-correct scale as a mismatch and "corrects" their text smaller
// (the be8e16f7e settle-guard defect, #5701). textScale <= 0 (unreadable)
// is treated as 1.0; dpi 0 (unreadable) yields 0 so callers can gate on it.
// Pure and platform-independent so mac/linux CI pins the formula.
func rasterizationTargetForDPI(dpi uint32, textScale float64) float64 {
	if dpi == 0 {
		return 0
	}
	if textScale <= 0 {
		textScale = 1.0
	}
	return float64(dpi) / 96.0 * textScale
}

// dpiflapSettleNeedsCorrectivePut decides whether dpiFlapSettleCheck must
// force one corrective rasterization-scale put before its bounds re-assert. It
// is the pure (controller-free) core of the v200.0.23 settle guard (#5701):
// extracted so the regression — native detection going silent and leaving the
// raster stale at settle — is unit-testable on every platform. dpi is the
// window's current DPI (0 = unknown); raster is the controller's live
// rasterization scale (<=0 = unavailable); target is dpi/96; isMinimizing
// gates the controller touch off the #5605 minimised-restore crash class.
func dpiflapSettleNeedsCorrectivePut(raster, target float64, dpi uint32, isMinimizing bool) bool {
	if dpi == 0 || raster <= 0 || target <= 0 || isMinimizing {
		return false
	}
	return math.Abs(raster-target) > dpiFlapSettleScaleTolerance
}
