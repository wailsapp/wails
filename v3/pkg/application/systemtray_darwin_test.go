//go:build darwin && !ios && !server

package application

import "testing"

// NSEventType and pressedMouseButtons bit values from AppKit, duplicated
// here so the test does not need to import AppKit.
const (
	testNSEventTypeLeftMouseDown  = 1
	testNSEventTypeRightMouseDown = 3
	testNSEventTypeLeftMouseUp    = 2
	testNSEventTypeMouseMoved     = 5
	testNSEventTypeMouseEntered   = 8
	testNSEventTypeCursorUpdate   = 17

	testPressedNone  = 0
	testPressedLeft  = 1 << 0
	testPressedRight = 1 << 1
)

// TestCoerceStatusItemEventType_MacOS27Regression pins the fix for #5752.
// See systemtray_darwin.h for the rationale.
func TestCoerceStatusItemEventType_MacOS27Regression(t *testing.T) {
	tests := []struct {
		name          string
		rawEventType  int
		pressed       uint64
		wantEventType int
	}{
		// macOS <=26 path: currentEvent still is the mouse-down; helper is
		// a no-op.
		{
			name:          "left mouse-down passes through unchanged (<=26)",
			rawEventType:  testNSEventTypeLeftMouseDown,
			pressed:       testPressedLeft,
			wantEventType: testNSEventTypeLeftMouseDown,
		},
		{
			name:          "right mouse-down passes through unchanged (<=26)",
			rawEventType:  testNSEventTypeRightMouseDown,
			pressed:       testPressedRight,
			wantEventType: testNSEventTypeRightMouseDown,
		},

		// macOS 27 regression path: raw type is *not* a mouse-down, so the
		// helper reads pressedMouseButtons to recover the intended button.
		// The MouseMoved+pressed=1 case is the exact combination observed
		// in field logs on macOS 27 for a plain left click on the tray icon.
		{
			name:          "macOS 27: MouseMoved with left pressed -> LeftMouseDown",
			rawEventType:  testNSEventTypeMouseMoved,
			pressed:       testPressedLeft,
			wantEventType: testNSEventTypeLeftMouseDown,
		},
		{
			name:          "macOS 27: MouseMoved with right pressed -> RightMouseDown",
			rawEventType:  testNSEventTypeMouseMoved,
			pressed:       testPressedRight,
			wantEventType: testNSEventTypeRightMouseDown,
		},
		{
			name:          "macOS 27: CursorUpdate with left pressed -> LeftMouseDown",
			rawEventType:  testNSEventTypeCursorUpdate,
			pressed:       testPressedLeft,
			wantEventType: testNSEventTypeLeftMouseDown,
		},
		{
			name:          "macOS 27: MouseEntered with no button pressed -> LeftMouseDown (default)",
			rawEventType:  testNSEventTypeMouseEntered,
			pressed:       testPressedNone,
			wantEventType: testNSEventTypeLeftMouseDown,
		},
		{
			name:          "macOS 27: LeftMouseUp with left pressed -> LeftMouseDown",
			rawEventType:  testNSEventTypeLeftMouseUp,
			pressed:       testPressedLeft,
			wantEventType: testNSEventTypeLeftMouseDown,
		},

		// Right takes precedence when both bits are set: this matches
		// existing right-click semantics (right-click hides the attached
		// window and opens the menu).
		{
			name:          "both buttons pressed -> RightMouseDown wins",
			rawEventType:  testNSEventTypeMouseMoved,
			pressed:       testPressedLeft | testPressedRight,
			wantEventType: testNSEventTypeRightMouseDown,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := coerceStatusItemEventType(tc.rawEventType, tc.pressed)
			if got != tc.wantEventType {
				t.Fatalf("coerceStatusItemEventType(raw=%d, pressed=%#x) = %d, want %d",
					tc.rawEventType, tc.pressed, got, tc.wantEventType)
			}
		})
	}
}
