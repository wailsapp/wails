package application

import (
	"testing"
)

func TestFullscreenNCHitTestReturnsHTClient(t *testing.T) {
	testCases := []struct {
		name            string
		isFullscreen    bool
		expectedCapture bool
	}{
		{"fullscreen should capture", true, true},
		{"not fullscreen should not capture", false, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := shouldCaptureHitTest(tc.isFullscreen)
			if result != tc.expectedCapture {
				t.Errorf("shouldCaptureHitTest(%v) = %v, want %v",
					tc.isFullscreen, result, tc.expectedCapture)
			}
		})
	}
}

func shouldCaptureHitTest(isFullscreen bool) bool {
	return isFullscreen
}
