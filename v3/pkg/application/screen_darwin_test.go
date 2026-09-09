//go:build darwin && !ios && !server

package application

import (
	"strconv"
	"testing"
	"unicode/utf8"
)

// TestScreenStringsSurviveAutoreleasePool guards against storing autoreleased
// UTF8String buffers in the C Screen struct (#5556). getAllScreens runs inside
// an explicit autorelease pool that drains before the structs reach Go, so the
// screen id and name must be malloc'd copies — a dangling pointer into a
// drained pool yields garbage here (deterministically so under
// MallocScribble=1, which poisons freed memory with 0x55 bytes).
func TestScreenStringsSurviveAutoreleasePool(t *testing.T) {
	screens := allScreens()
	if len(screens) == 0 {
		t.Skip("no screens attached (headless)")
	}

	for i, screen := range screens {
		// The id is the CGDirectDisplayID formatted with %d, so it must
		// parse back as an integer.
		if _, err := strconv.ParseInt(screen.ID, 10, 64); err != nil {
			t.Errorf("screen %d: ID %q is not numeric — points at freed memory? %v", i, screen.ID, err)
		}
		if !utf8.ValidString(screen.Name) {
			t.Errorf("screen %d: Name %q is not valid UTF-8 — points at freed memory?", i, screen.Name)
		}
		t.Logf("screen %d: ID=%q Name=%q primary=%v", i, screen.ID, screen.Name, screen.IsPrimary)
	}
}
