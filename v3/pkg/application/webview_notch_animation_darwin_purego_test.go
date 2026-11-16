//go:build darwin && purego && !ios && !server

package application

import (
	"testing"
)

// notchTestTarget is where testNotchFrame lands on the 14-inch test screen:
// centred on the housing (x = round(756 - 700/2) = 406) and flush with the top
// of the screen frame (y = 982 - 104 = 878).
var notchTestTarget = NSRect{Origin: CGPoint{X: 406, Y: 878}, Size: CGSize{Width: 700, Height: 104}}

func TestNotchTestTargetMatchesPlacement(t *testing.T) {
	got := notchTargetFrame(testNotchFrame, testScreenFrame, testVisibleFrame, testHousingCentre, true)
	if got != notchTestTarget {
		t.Fatalf("placed frame = %+v, want %+v", got, notchTestTarget)
	}
}

func TestNotchShowPlan(t *testing.T) {
	plan := notchShowPlan(notchTestTarget, 0.42)

	// The slide starts from the frame directly above the screen top: it is
	// the target moved up by its own height (NSMaxY of the target).
	wantStart := NSRect{Origin: CGPoint{X: 406, Y: 982}, Size: CGSize{Width: 700, Height: 104}}
	if plan.StartFrame != wantStart {
		t.Errorf("StartFrame = %+v, want %+v", plan.StartFrame, wantStart)
	}
	if plan.StartFrame.Origin.Y != testScreenFrame.Origin.Y+testScreenFrame.Size.Height {
		t.Errorf("StartFrame does not sit on the top edge of the screen")
	}
	if !approxEqual(plan.StartAlpha, 0.96) {
		t.Errorf("StartAlpha = %v, want 0.96", plan.StartAlpha)
	}
	if plan.EndFrame != notchTestTarget {
		t.Errorf("EndFrame = %+v, want the target %+v", plan.EndFrame, notchTestTarget)
	}
	if !approxEqual(plan.EndAlpha, 1.0) {
		t.Errorf("EndAlpha = %v, want 1", plan.EndAlpha)
	}
	if !approxEqual(plan.Duration, 0.42) {
		t.Errorf("Duration = %v, want 0.42", plan.Duration)
	}
	wantTiming := [4]float64{0.22, 1.0, 0.36, 1.0}
	if plan.Timing != wantTiming {
		t.Errorf("Timing = %v, want %v", plan.Timing, wantTiming)
	}
}

func TestNotchHidePlan(t *testing.T) {
	plan := notchHidePlan(notchTestTarget, 0.42)

	// cgo starts the hide from the current frame: the plan sets no start.
	if plan.StartFrame != (NSRect{}) || plan.StartAlpha != 0 {
		t.Errorf("hide plan sets a start state (%+v, %v); want none", plan.StartFrame, plan.StartAlpha)
	}
	wantEnd := NSRect{Origin: CGPoint{X: 406, Y: 982}, Size: CGSize{Width: 700, Height: 104}}
	if plan.EndFrame != wantEnd {
		t.Errorf("EndFrame = %+v, want %+v", plan.EndFrame, wantEnd)
	}
	if !approxEqual(plan.EndAlpha, 0.9) {
		t.Errorf("EndAlpha = %v, want 0.9", plan.EndAlpha)
	}
	if !approxEqual(plan.Duration, 0.28) {
		t.Errorf("Duration = %v, want two thirds of 0.42 = 0.28", plan.Duration)
	}
	wantTiming := [4]float64{0.4, 0.0, 1.0, 1.0}
	if plan.Timing != wantTiming {
		t.Errorf("Timing = %v, want %v", plan.Timing, wantTiming)
	}
}

func TestNotchHideDuration(t *testing.T) {
	tests := []struct {
		show, want float64
	}{
		{show: 0.42, want: 0.28},
		{show: 0.3, want: 0.2},
		{show: 1, want: 2.0 / 3.0},
		{show: 0, want: 0},
	}
	for _, tt := range tests {
		if got := notchHideDuration(tt.show); !approxEqual(got, tt.want) {
			t.Errorf("notchHideDuration(%v) = %v, want %v", tt.show, got, tt.want)
		}
	}
}

func TestNotchShowModeFor(t *testing.T) {
	tests := []struct {
		name             string
		requestedVisible bool
		animated         bool
		duration         float64
		want             notchShowMode
	}{
		{name: "already requested orders front regardless of animation", requestedVisible: true, animated: true, duration: 0.42, want: notchShowAlreadyVisible},
		{name: "already requested without animation", requestedVisible: true, want: notchShowAlreadyVisible},
		{name: "not animated is instant", animated: false, duration: 0.42, want: notchShowInstant},
		{name: "zero duration is instant", animated: true, duration: 0, want: notchShowInstant},
		{name: "negative duration is instant", animated: true, duration: -1, want: notchShowInstant},
		{name: "animated with a positive duration slides", animated: true, duration: 0.42, want: notchShowSlide},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := notchShowModeFor(tt.requestedVisible, tt.animated, tt.duration); got != tt.want {
				t.Errorf("notchShowModeFor(%v, %v, %v) = %v, want %v", tt.requestedVisible, tt.animated, tt.duration, got, tt.want)
			}
		})
	}
}

func TestNotchHideModeFor(t *testing.T) {
	tests := []struct {
		name             string
		requestedVisible bool
		visible          bool
		animated         bool
		duration         float64
		want             notchHideMode
	}{
		{name: "nothing requested and not visible is a no-op", animated: true, duration: 0.42, want: notchHideNoop},
		{name: "a requested show that is not yet on screen still hides", requestedVisible: true, animated: true, duration: 0.42, want: notchHideSlide},
		{name: "visible but not requested still hides", visible: true, animated: true, duration: 0.42, want: notchHideSlide},
		{name: "visible without animation is instant", visible: true, requestedVisible: true, duration: 0.42, want: notchHideInstant},
		{name: "visible with zero duration is instant", visible: true, requestedVisible: true, animated: true, want: notchHideInstant},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := notchHideModeFor(tt.requestedVisible, tt.visible, tt.animated, tt.duration); got != tt.want {
				t.Errorf("notchHideModeFor(%v, %v, %v, %v) = %v, want %v",
					tt.requestedVisible, tt.visible, tt.animated, tt.duration, got, tt.want)
			}
		})
	}
}

func TestNotchTrackingOptionsValues(t *testing.T) {
	// NSTrackingMouseEnteredAndExited | NSTrackingActiveAlways | NSTrackingInVisibleRect.
	const want = 0x01 | 0x80 | 0x200
	if notchTrackingOptions != want {
		t.Fatalf("notchTrackingOptions = %#x, want %#x", notchTrackingOptions, want)
	}
}

// TestNotchWindowClassRespondsToOverrides checks that the registered notch
// class answers the selectors cgo's WebviewNotchWindow implements.
func TestNotchWindowClassRespondsToOverrides(t *testing.T) {
	cls := registerNotchWindowClass()
	for _, sel := range []string{
		"showAnimated:duration:",
		"hideAnimated:duration:",
		"mouseEntered:",
		"mouseExited:",
		"close",
		"performWindowDragWithEvent:",
	} {
		if !get[bool](cls, "instancesRespondToSelector:", sel_(sel)) {
			t.Errorf("WailsWebviewNotchPanel does not respond to %s", sel)
		}
	}
}
