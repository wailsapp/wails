package application

import (
	"math"
	"testing"
)

func TestZoomFactor(t *testing.T) {
	tests := []struct {
		name      string
		input     float64
		frameless bool
		want      float64
		ok        bool
	}{
		{"zero rejected", 0, false, 0, false},
		{"negative rejected", -0.5, false, 0, false},
		{"NaN rejected", math.NaN(), false, 0, false},
		{"+Inf rejected", math.Inf(1), false, 0, false},
		{"-Inf rejected", math.Inf(-1), false, 0, false},
		{"0.5 passes through", 0.5, false, 0.5, true},
		{"0.99 passes through", 0.99, false, 0.99, true},
		{"1.0 stays 1.0", 1.0, false, 1.0, true},
		{"2.0 passes through", 2.0, false, 2.0, true},
		{"frameless: 0.5 raised to 1.0", 0.5, true, 1.0, true},
		{"frameless: 1.0 stays 1.0", 1.0, true, 1.0, true},
		{"frameless: 1.5 passes through", 1.5, true, 1.5, true},
		{"frameless: zero rejected", 0, true, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := zoomFactor(tt.input, tt.frameless)
			if ok != tt.ok || got != tt.want {
				t.Errorf("zoomFactor(%v, %v) = (%v, %v), want (%v, %v)", tt.input, tt.frameless, got, ok, tt.want, tt.ok)
			}
		})
	}
}
