package application

import (
	"testing"
	"time"
)

// TestRasterizationTargetForDPI pins the target formula (#5701): the correct
// rasterization scale is dpi/96 × the Windows text scale factor. The
// text-scale term is what the be8e16f7e settle guard missed — comparing
// raster against a bare dpi/96 mis-flags every text-scaling user's correct
// scale as stale and "corrects" their text smaller.
func TestRasterizationTargetForDPI(t *testing.T) {
	tests := []struct {
		name      string
		dpi       uint32
		textScale float64
		want      float64
	}{
		{name: "216 dpi, no text scaling", dpi: 216, textScale: 1.0, want: 2.25},
		{name: "120 dpi, no text scaling", dpi: 120, textScale: 1.0, want: 1.25},
		{name: "216 dpi with 125% text (the be8e16f7e defect case)", dpi: 216, textScale: 1.25, want: 2.8125},
		{name: "96 dpi with 225% text (max text size)", dpi: 96, textScale: 2.25, want: 2.25},
		{name: "dpi unreadable yields 0 so callers can gate", dpi: 0, textScale: 1.25, want: 0},
		{name: "text scale unreadable (0) treated as 1.0", dpi: 216, textScale: 0, want: 2.25},
		{name: "text scale negative treated as 1.0", dpi: 120, textScale: -1, want: 1.25},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rasterizationTargetForDPI(tt.dpi, tt.textScale)
			if got != tt.want {
				t.Errorf("rasterizationTargetForDPI(%d, %.2f) = %v, want %v", tt.dpi, tt.textScale, got, tt.want)
			}
		})
	}
}

// quietInput is a scaleReconcileInput at rest: controller live, no storm, no
// drag, last flip long past, no prior put, puts permitted. Tests override the
// field under test so each case reads as its delta from "quiet and healthy".
func quietInput() scaleReconcileInput {
	return scaleReconcileInput{
		raster:          2.25,
		target:          2.25,
		dpi:             216,
		sinceLastFlip:   5 * time.Second,
		sinceLastPut:    -1,
		controllerReady: true,
		allowPut:        true,
	}
}

// TestDecideScaleReconcile pins the single decision point of the DPI verify
// ladder (#5701): every gate, in order, plus the field-trace regressions.
// Platform-independent so the matrix runs on mac/linux CI, not only Windows.
func TestDecideScaleReconcile(t *testing.T) {
	mismatch := func(mutate func(*scaleReconcileInput)) scaleReconcileInput {
		in := quietInput()
		in.raster = 1.25 // vs target 2.25 — the v200.0.23/24 stuck value
		if mutate != nil {
			mutate(&in)
		}
		return in
	}
	tests := []struct {
		name       string
		in         scaleReconcileInput
		wantAction scaleReconcileAction
		wantReason string
	}{
		// The field defects this ladder exists for: raster stuck at 1.25 under
		// dpi 216/target 2.25 (v200.0.23: 85s; v200.0.24: 62s), quiet system.
		{"v200.0.23/24 stuck raster corrects", mismatch(nil), scaleReconcilePut, "mismatch"},
		{"stale raster the other way", func() scaleReconcileInput {
			in := quietInput()
			in.raster, in.target, in.dpi = 2.25, 1.25, 120
			return in
		}(), scaleReconcilePut, "mismatch"},
		{"spontaneous mismatch, no flip ever seen", mismatch(func(in *scaleReconcileInput) {
			in.sinceLastFlip = -1
		}), scaleReconcilePut, "mismatch"},

		// The invariant holding — the routine OK lines field logs must show.
		{"in sync on 216", quietInput(), scaleReconcileOK, "in-sync"},
		{"in sync on 120", func() scaleReconcileInput {
			in := quietInput()
			in.raster, in.target, in.dpi = 1.25, 1.25, 120
			return in
		}(), scaleReconcileOK, "in-sync"},
		{"in sync even when puts are forbidden", func() scaleReconcileInput {
			in := quietInput()
			in.allowPut = false
			return in
		}(), scaleReconcileOK, "in-sync"},
		// Text-scale regression pin (the be8e16f7e defect): a 125%-text user's
		// raster of 2.8125 on dpi 216 is CORRECT when the target includes text
		// scale — and a mismatch against a naive dpi/96 target, which is why
		// rasterizationTargetForDPI exists.
		{"text-scaled raster vs text-aware target", func() scaleReconcileInput {
			in := quietInput()
			in.raster, in.target = 2.8125, rasterizationTargetForDPI(216, 1.25)
			return in
		}(), scaleReconcileOK, "in-sync"},
		{"text-scaled raster vs naive dpi/96 target mis-flags", func() scaleReconcileInput {
			in := quietInput()
			in.raster, in.target = 2.8125, 2.25
			return in
		}(), scaleReconcilePut, "mismatch"},

		// Tolerance boundary: float jitter must not put; a real step must.
		{"jitter well under tolerance", func() scaleReconcileInput {
			in := quietInput()
			in.raster = 2.245
			return in
		}(), scaleReconcileOK, "in-sync"},
		{"just under tolerance (0.009)", func() scaleReconcileInput {
			in := quietInput()
			in.raster = 2.259
			return in
		}(), scaleReconcileOK, "in-sync"},
		{"just over tolerance (0.011)", func() scaleReconcileInput {
			in := quietInput()
			in.raster = 2.261
			return in
		}(), scaleReconcilePut, "mismatch"},

		// Evaluability gates, in precedence order.
		{"controller not ready", mismatch(func(in *scaleReconcileInput) {
			in.controllerReady = false
		}), scaleReconcileSkip, "controller"},
		{"minimising (#5605 gate, raster unread)", mismatch(func(in *scaleReconcileInput) {
			in.isMinimizing = true
			in.raster = 0 // callers must not COM-read while minimising
		}), scaleReconcileSkip, "minimising"},
		{"rebuild in progress", mismatch(func(in *scaleReconcileInput) {
			in.rebuildInProgress = true
		}), scaleReconcileSkip, "rebuild"},
		{"dpi unreadable", mismatch(func(in *scaleReconcileInput) {
			in.dpi, in.target = 0, 0
		}), scaleReconcileSkip, "dpi-unreadable"},
		{"raster unavailable", mismatch(func(in *scaleReconcileInput) {
			in.raster = 0
		}), scaleReconcileSkip, "raster-unavailable"},

		// Quiet gates: mismatches wait for churn to end.
		{"mismatch during storm defers to the settle chain", mismatch(func(in *scaleReconcileInput) {
			in.stormActive = true
		}), scaleReconcileDefer, "storm"},
		{"mismatch mid-drag defers", mismatch(func(in *scaleReconcileInput) {
			in.inSizeMove = true
		}), scaleReconcileDefer, "in-drag"},
		{"mismatch 300ms after a flip defers", mismatch(func(in *scaleReconcileInput) {
			in.sinceLastFlip = 300 * time.Millisecond
		}), scaleReconcileDefer, "recent-flip"},
		{"mismatch 401ms after a flip puts", mismatch(func(in *scaleReconcileInput) {
			in.sinceLastFlip = 401 * time.Millisecond
		}), scaleReconcilePut, "mismatch"},

		// Write gates.
		{"log-only pass never puts", mismatch(func(in *scaleReconcileInput) {
			in.allowPut = false
		}), scaleReconcileDefer, "log-only"},
		{"put 500ms after the last put is rate-limited", mismatch(func(in *scaleReconcileInput) {
			in.sinceLastPut = 500 * time.Millisecond
		}), scaleReconcileDefer, "rate-limited"},
		{"put 1.1s after the last put proceeds", mismatch(func(in *scaleReconcileInput) {
			in.sinceLastPut = 1100 * time.Millisecond
		}), scaleReconcilePut, "mismatch"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action, reason := decideScaleReconcile(tt.in)
			if action != tt.wantAction || reason != tt.wantReason {
				t.Errorf("decideScaleReconcile(%+v) = (%v, %q), want (%v, %q)",
					tt.in, action, reason, tt.wantAction, tt.wantReason)
			}
		})
	}
}

// TestScaleReconcileShouldRetry pins which deferred gates arm a +1s retry
// probe: transient conditions with no other watcher. "storm" must NOT retry
// (the settle chain re-verifies) and "log-only" is terminal by definition.
func TestScaleReconcileShouldRetry(t *testing.T) {
	want := map[string]bool{
		"in-drag":      true,
		"recent-flip":  true,
		"rate-limited": true,
		"storm":        false,
		"log-only":     false,
		"in-sync":      false,
		"mismatch":     false,
		"controller":   false,
	}
	for reason, expected := range want {
		if got := scaleReconcileShouldRetry(reason); got != expected {
			t.Errorf("scaleReconcileShouldRetry(%q) = %v, want %v", reason, got, expected)
		}
	}
}
