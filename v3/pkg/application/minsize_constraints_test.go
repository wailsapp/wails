package application

import (
	"testing"
)

func TestEnforceMinSizeConstraintsBelowMinimum(t *testing.T) {
	testCases := []struct {
		name          string
		width         int
		height        int
		minWidth      int
		minHeight     int
		expectedW     int
		expectedH     int
		shouldEnforce bool
	}{
		{"both below min", 200, 150, 400, 300, 400, 300, true},
		{"width below min", 200, 400, 400, 0, 400, 400, true},
		{"height below min", 500, 100, 0, 300, 500, 300, true},
		{"both above min", 500, 400, 400, 300, 500, 400, false},
		{"exactly at min", 400, 300, 400, 300, 400, 300, false},
		{"no constraints", 200, 100, 0, 0, 200, 100, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			newW, newH, changed := enforceMinSize(tc.width, tc.height, tc.minWidth, tc.minHeight)
			if newW != tc.expectedW || newH != tc.expectedH {
				t.Errorf("got (%d, %d), want (%d, %d)",
					newW, newH,
					tc.expectedW, tc.expectedH)
			}
			if changed != tc.shouldEnforce {
				t.Errorf("changed = %v, want %v", changed, tc.shouldEnforce)
			}
		})
	}
}

func enforceMinSize(currentW, currentH, minWidth, minHeight int) (int, int, bool) {
	changed := false
	if minWidth > 0 && currentW < minWidth {
		currentW = minWidth
		changed = true
	}
	if minHeight > 0 && currentH < minHeight {
		currentH = minHeight
		changed = true
	}
	return currentW, currentH, changed
}
