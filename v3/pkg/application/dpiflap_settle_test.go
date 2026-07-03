package application

import "testing"

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

// TestDpiflapSettleNeedsCorrectivePut pins the v200.0.23 settle guard (#5701):
// when native monitor-scale detection goes SILENT (no RasterizationScaleChanged
// event after the last DPI flip), the settle must force one corrective scale
// put iff the raster is stale relative to the target. Platform-independent so
// the regression runs on mac/linux CI, not only on Windows.
func TestDpiflapSettleNeedsCorrectivePut(t *testing.T) {
	tests := []struct {
		name        string
		raster      float64
		target      float64
		dpi         uint32
		isMinimizing bool
		want        bool
	}{
		// v200.0.23 field defect: raster stuck at 1.25 under dpi 216 (target 2.25)
		// for 85s after native detection went silent — the exact case the guard
		// exists to catch.
		{name: "v200.0.23 stuck raster (216/2.25, raster 1.25)", raster: 1.25, target: 2.25, dpi: 216, want: true},
		{name: "stale raster the other way (120/1.25, raster 2.25)", raster: 2.25, target: 1.25, dpi: 120, want: true},

		// Normal settles: native detection corrected, raster already matches.
		{name: "native-corrected on 216", raster: 2.25, target: 2.25, dpi: 216, want: false},
		{name: "native-corrected on 120", raster: 1.25, target: 1.25, dpi: 120, want: false},
		{name: "half-step 1.5x corrected", raster: 1.5, target: 1.5, dpi: 144, want: false},

		// Tolerance boundary: float jitter inside tolerance must NOT trip a put
		// (syncWebviewRasterizationScale would no-op anyway, but the gate stays
		// conservative); just-over-tolerance must.
		{name: "jitter well under tolerance", raster: 2.245, target: 2.25, dpi: 216, want: false},
		{name: "just under tolerance (0.009)", raster: 2.259, target: 2.25, dpi: 216, want: false},
		{name: "just over tolerance (0.011)", raster: 2.261, target: 2.25, dpi: 216, want: true},

		// Safety gates: never touch the controller in these states.
		{name: "controller unavailable (raster 0)", raster: 0, target: 2.25, dpi: 216, want: false},
		{name: "dpi unknown (0)", raster: 2.25, target: 0, dpi: 0, want: false},
		{name: "minimising (#5605 gate)", raster: 1.25, target: 2.25, dpi: 216, isMinimizing: true, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dpiflapSettleNeedsCorrectivePut(tt.raster, tt.target, tt.dpi, tt.isMinimizing)
			if got != tt.want {
				t.Errorf("dpiflapSettleNeedsCorrectivePut(raster=%.3f, target=%.3f, dpi=%d, isMinimizing=%v) = %v, want %v",
					tt.raster, tt.target, tt.dpi, tt.isMinimizing, got, tt.want)
			}
		})
	}
}
