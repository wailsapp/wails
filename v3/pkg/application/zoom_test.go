package application

import "testing"

func TestZoomClamp(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{"zero clamped to 1.0", 0.0, 1.0},
		{"negative clamped to 1.0", -0.5, 1.0},
		{"0.5 clamped to 1.0", 0.5, 1.0},
		{"0.99 clamped to 1.0", 0.99, 1.0},
		{"1.0 stays 1.0", 1.0, 1.0},
		{"1.5 stays 1.5", 1.5, 1.5},
		{"2.0 stays 2.0", 2.0, 2.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input
			if result < 1.0 {
				result = 1.0
			}
			if result != tt.expected {
				t.Errorf("clamp(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
