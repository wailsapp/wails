//go:build windows

package application

import (
	"time"

	"golang.org/x/sys/windows/registry"
)

const (
	// scaleVerifyPutAckProbe re-verifies (log-only) after a corrective put to
	// confirm the put actually took — the v200.0.24 62s stick begged the
	// question whether a put would even stick, and a put-ack read-back alone
	// can race the browser's async commit.
	scaleVerifyPutAckProbe = time.Second
	// scaleVerifyRetryProbe re-checks a deferred mismatch (in-drag /
	// recent-flip / rate-limited): transient gates with no other owner
	// watching. One live timer per chain; the chain dies on OK/Skip.
	scaleVerifyRetryProbe = time.Second
	// scaleVerifyProbeShort / scaleVerifyProbeLong watch the window after a
	// storm settles. Wrong-scale events provably land up to ~6s after the
	// last flip (v200.0.24 field trace), i.e. AFTER the settle verify — these
	// probes bound that exposure to ≤10s and, critically, make it a logged
	// MISMATCH line instead of a screenshot-only defect.
	scaleVerifyProbeShort = 2 * time.Second
	scaleVerifyProbeLong  = 10 * time.Second
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

// scaleOwnerName names the configured rasterization-scale writer for the
// telemetry lines ("webview2" = native monitor-scale detection).
func scaleOwnerName() string {
	return "webview2"
}

// verifyWebviewScale is the single observer/reconciler of the scale
// invariant: at all quiet times, RasterizationScale == GetDpiForWindow/96 ×
// textScale. It reads the live DPI and live raster, has decideScaleReconcile
// judge the diff, and ALWAYS logs the verdict — the OK lines are the point:
// field sessions must show the diffs we detected and the ones we ruled in
// sync, so a wrong-scale exposure can never again hide in a logging blind
// window (v200.0.24: 62s of wrong UI, visible only via screenshot, #5701).
// On a permitted mismatch it performs the ONE corrective scale-put + bounds
// re-assert, stamps the put for rate limiting, and arms a log-only re-verify
// to confirm the put took. Deferred transient gates arm a +1s retry.
// Main thread only (COM + storm bookkeeping); reason is a short tag naming
// the trigger point, logged verbatim.
func (w *windowsWebviewWindow) verifyWebviewScale(reason string, allowPut bool) {
	if w.hwnd == 0 || w.parent.isDestroyed() || globalApplication.performingShutdown {
		return
	}
	dpi, _ := w.DPI()
	// Live GetDpiForWindow is the target authority, never lastKnownDPI: a
	// disagreement means the OS never delivered the final WM_DPICHANGED (the
	// wrong-TARGET failure mode) — surface it as STALE and re-track.
	lastKnown := w.lastKnownDPI
	stale := ""
	if dpi != 0 && lastKnown != dpi {
		stale = " STALE"
		if !w.isMinimizing {
			w.lastKnownDPI = dpi
		}
	}
	controllerReady := w.chromium != nil && w.chromium.IsReady()
	raster := 0.0
	if controllerReady && !w.isMinimizing {
		// Property read only; skipped while minimising — COM into a
		// possibly-suspended controller is the #5605 restore-crash class.
		raster = w.currentWebviewRasterizationScale()
	}
	textScale := windowsTextScaleFactor()
	target := rasterizationTargetForDPI(uint32(dpi), textScale)
	sinceFlip := time.Duration(-1)
	if !w.lastDPITransitionAt.IsZero() {
		sinceFlip = time.Since(w.lastDPITransitionAt)
	}
	sincePut := time.Duration(-1)
	if !w.lastScalePutAt.IsZero() {
		sincePut = time.Since(w.lastScalePutAt)
	}
	action, gate := decideScaleReconcile(scaleReconcileInput{
		raster:            raster,
		target:            target,
		dpi:               uint32(dpi),
		isMinimizing:      w.isMinimizing,
		stormActive:       !w.dpiFlapSuppressUntil.IsZero(),
		inSizeMove:        w.inSizeMove,
		sinceLastFlip:     sinceFlip,
		sinceLastPut:      sincePut,
		rebuildInProgress: w.webviewRebuildInProgress,
		controllerReady:   controllerReady,
		allowPut:          allowPut,
	})
	var verdict string
	switch action {
	case scaleReconcileOK:
		verdict = "OK"
	case scaleReconcilePut:
		verdict = "MISMATCH → corrective put"
	case scaleReconcileDefer:
		verdict = "MISMATCH deferred (" + gate + ")"
	case scaleReconcileSkip:
		verdict = "SKIP (" + gate + ")"
	}
	globalApplication.warning(
		"DPI verify [%s]: window %d dpi %d (lastKnown %d%s) raster %.3f target %.3f textScale %.2f diff %+.3f owner=%s — %s (#5701)",
		reason, w.parent.id, dpi, lastKnown, stale, raster, target, textScale, raster-target, scaleOwnerName(), verdict)
	switch action {
	case scaleReconcilePut:
		ok, readBack := w.putWebviewRasterizationScale(target)
		w.lastScalePutAt = time.Now()
		relayout := ok && controllerReady
		if relayout {
			// A scale change without a bounds re-assert re-scales the frame
			// WITHOUT re-laying it out (#5677) — the Resize is what makes the
			// corrected scale visible.
			w.chromium.Resize()
		}
		globalApplication.warning(
			"DPI verify put-ack [%s]: window %d put %.3f read-back %.3f ok=%v relayout=%v — re-verify in 1s (#5701)",
			reason, w.parent.id, target, readBack, ok, relayout)
		w.scheduleScaleVerify("put-ack:"+reason, scaleVerifyPutAckProbe, false)
	case scaleReconcileDefer:
		if scaleReconcileShouldRetry(gate) {
			w.scheduleScaleVerify(gate+"-retry", scaleVerifyRetryProbe, allowPut)
		}
	}
}

// scheduleScaleVerify arms a verify pass after delay, hopping to the main
// thread like scheduleDPIFlapSettle. Duplicate pending probes are harmless:
// each pass re-reads live state, no-ops to an OK line when the invariant
// holds, and arms at most one successor.
func (w *windowsWebviewWindow) scheduleScaleVerify(reason string, delay time.Duration, allowPut bool) {
	time.AfterFunc(delay, func() {
		InvokeAsync(func() {
			w.verifyWebviewScale(reason, allowPut)
		})
	})
}

// putWebviewRasterizationScale writes target onto the controller and reads it
// back for the put-ack breadcrumb. Unlike syncWebviewRasterizationScale it
// does not pre-compare — verifyWebviewScale has already established the
// mismatch — so the read-back is the put's outcome, not its precondition.
// ok=false means the controller was unavailable or the put itself failed;
// readBack 0 means the confirmation read failed (put may still have landed —
// the log-only re-verify settles it).
func (w *windowsWebviewWindow) putWebviewRasterizationScale(target float64) (bool, float64) {
	if w.chromium == nil {
		return false, 0
	}
	controller := w.chromium.GetController()
	if controller == nil {
		return false, 0
	}
	controller3 := controller.GetICoreWebView2Controller3()
	if controller3 == nil {
		return false, 0
	}
	if err := controller3.PutRasterizationScale(target); err != nil {
		globalApplication.error("failed to update WebView2 rasterization scale: %s", err)
		return false, 0
	}
	readBack, err := controller3.GetRasterizationScale()
	if err != nil {
		return true, 0
	}
	return true, readBack
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
		"WebView2 scale ownership [%s]: window %d owner=%s AddRasterizationScaleChanged=%s ShouldDetectMonitorScaleChanges=%s textScale %.2f (#5701)",
		reason, w.parent.id, scaleOwnerName(), regStatus, detectionStatus, windowsTextScaleFactor())
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
