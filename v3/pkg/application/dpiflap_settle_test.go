package application

import "testing"

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
