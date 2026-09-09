//go:build windows

package w32

import (
	"strings"
	"testing"
)

// TestGetAllScreensToleratesCursorFailure is the regression test for the
// startup crash where a failed GetCursorPos (workstation locked / secure
// desktop at launch) aborted screen enumeration. processAndCacheScreens()
// treats that error as fatal (os.Exit) during app.Run, so the whole app died
// before it ever showed. GetAllScreens must now degrade gracefully.
//
// The assertion is display-independent: on a headless runner EnumDisplayMonitors
// may legitimately fail, but the returned error must NEVER be the cursor error,
// and no screen may be marked current when the cursor position is unknown.
func TestGetAllScreensToleratesCursorFailure(t *testing.T) {
	orig := cursorPosForScreens
	t.Cleanup(func() { cursorPosForScreens = orig })

	// Simulate GetCursorPos returning FALSE (ERROR_ACCESS_DENIED).
	cursorPosForScreens = func() (POINT, bool) { return POINT{}, false }

	screens, err := GetAllScreens()

	if err != nil && strings.Contains(err.Error(), "GetCursorPos") {
		t.Fatalf("GetAllScreens must not fail on a cursor-position error; got: %v", err)
	}
	for i, s := range screens {
		if s.IsCurrent {
			t.Errorf("screen[%d] marked IsCurrent despite the cursor query failing", i)
		}
	}
}

// TestCursorPosForScreensDefaultCallable confirms the default seam is wired to
// the real syscall and is safe to call (ok is true or false depending on the
// session's desktop state; both are valid).
func TestCursorPosForScreensDefaultCallable(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("default cursorPosForScreens panicked: %v", r)
		}
	}()
	_, _ = cursorPosForScreens()
}
