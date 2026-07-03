package application

import (
	"math"
	"time"
)

// dpiFlapReleaseSettle is the quiet threshold once the user has released
// the drag (!inSizeMove). A parked window that is not straddling produces
// no further transitions, so the frozen rasterization scale can snap back
// almost immediately — the v200.0.7 field test showed the ~1.6 s tail is
// still human-visible when a drag ends inside a suppression window. A
// parked STRADDLING window keeps oscillating: a premature settle there
// costs one transition before the parked fast path re-trips and the
// resolver ends the straddle.
//
// Declared here — not with the other dpiFlap timing constants in
// webview_window_windows.go — because decideScaleReconcile uses it as the
// post-flip quiet gate and must compile on every OS for CI.
const dpiFlapReleaseSettle = 400 * time.Millisecond

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

// scaleReconcileAction is what a DPI verify pass decided to do about the
// raster-vs-target diff it observed. Every verify logs its action — including
// OK — so field sessions show the diffs we detected AND the ones we ruled in
// sync, instead of leaving wrong-scale exposure to be inferred from
// screenshots (the v200.0.24 62s blind window, #5701).
type scaleReconcileAction int

const (
	scaleReconcileOK    scaleReconcileAction = iota // in sync — nothing to do
	scaleReconcilePut                               // mismatch and safe — one corrective scale-put
	scaleReconcileDefer                             // mismatch but a gate blocks the put right now
	scaleReconcileSkip                              // cannot evaluate — reason says why
)

// scaleReconcileInput is the full decision state for one verify pass,
// controller-free so the gate order is unit-testable on every platform.
// raster <= 0 means the live scale was unavailable (or deliberately unread
// while minimising — COM into a possibly-suspended controller is the #5605
// restore-crash class, so callers must not read it then). sinceLastFlip /
// sinceLastPut are -1 when the event has never happened.
type scaleReconcileInput struct {
	raster, target    float64
	dpi               uint32
	isMinimizing      bool
	stormActive       bool
	inSizeMove        bool
	sinceLastFlip     time.Duration
	sinceLastPut      time.Duration
	rebuildInProgress bool
	controllerReady   bool
	allowPut          bool
}

// scaleReconcilePutMinInterval rate-limits corrective puts: reconciliation
// only ever needs one put to converge (the put itself echoes a
// RasterizationScaleChanged that then matches the target), so anything faster
// than 1/s means something else is rewriting the scale and hammering puts
// would fight it — per-flip-cadence puts are the field-proven browser-kill
// pattern this ladder exists to avoid (#5701, Steps 18/19).
const scaleReconcilePutMinInterval = time.Second

// decideScaleReconcile is the single decision point for "the rasterization
// scale disagrees with the window's DPI — may we correct it?". The returned
// reason string is the gate that decided (logged verbatim in the DPI verify
// breadcrumb). Gate order matters and is pinned by tests:
//
//  1. evaluability (controller / minimising / rebuild / unreadable dpi or
//     raster) — SKIP: we cannot even say whether the scale is wrong;
//  2. |raster-target| <= tolerance — OK: the invariant holds;
//  3. quiet gates (storm / in-drag / <settle-threshold since the last flip) —
//     DEFER: WM_DPICHANGED churn is still in flight, the settle chain owns
//     the correction, and putting mid-storm is the crash-adjacent pattern;
//  4. allowPut=false — DEFER: a log-only observation pass (e.g. the post-put
//     re-verify) never writes;
//  5. put rate limit — DEFER: see scaleReconcilePutMinInterval;
//  6. PUT.
func decideScaleReconcile(in scaleReconcileInput) (scaleReconcileAction, string) {
	if !in.controllerReady {
		return scaleReconcileSkip, "controller"
	}
	if in.isMinimizing {
		return scaleReconcileSkip, "minimising"
	}
	if in.rebuildInProgress {
		return scaleReconcileSkip, "rebuild"
	}
	if in.dpi == 0 || in.target <= 0 {
		return scaleReconcileSkip, "dpi-unreadable"
	}
	if in.raster <= 0 {
		return scaleReconcileSkip, "raster-unavailable"
	}
	if math.Abs(in.raster-in.target) <= dpiFlapSettleScaleTolerance {
		return scaleReconcileOK, "in-sync"
	}
	if in.stormActive {
		return scaleReconcileDefer, "storm"
	}
	if in.inSizeMove {
		return scaleReconcileDefer, "in-drag"
	}
	if in.sinceLastFlip >= 0 && in.sinceLastFlip < dpiFlapReleaseSettle {
		return scaleReconcileDefer, "recent-flip"
	}
	if !in.allowPut {
		return scaleReconcileDefer, "log-only"
	}
	if in.sinceLastPut >= 0 && in.sinceLastPut < scaleReconcilePutMinInterval {
		return scaleReconcileDefer, "rate-limited"
	}
	return scaleReconcilePut, "mismatch"
}

// scaleReconcileRetryGates lists the Defer reasons that arm a +1s retry
// probe: transient conditions with no other owner watching them. "storm" is
// excluded (the settle chain re-verifies when the storm ends) and "log-only"
// is terminal by definition.
func scaleReconcileShouldRetry(reason string) bool {
	switch reason {
	case "in-drag", "recent-flip", "rate-limited":
		return true
	}
	return false
}
