//go:build windows

package combridge

import (
	"testing"

	"golang.org/x/sys/windows"
)

// parseGUID and formatGUID must agree with ole32, which they replace.
func TestParseGUIDMatchesGUIDFromString(t *testing.T) {
	for _, s := range []string{
		iUnknownGUID,
		"{00000122-0000-0000-C000-000000000046}", // IDropTarget
		"{5c4889f0-5ef6-4c5a-952c-d8f1b92d0574}",
		"{e2fda4be-5456-406c-a261-3d452138362c}",
		"{FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF}",
		"{00000000-0000-0000-0000-000000000000}",
	} {
		want, err := windows.GUIDFromString(s)
		if err != nil {
			t.Fatalf("GUIDFromString(%q): %v", s, err)
		}
		got, err := parseGUID(s)
		if err != nil {
			t.Fatalf("parseGUID(%q): %v", s, err)
		}
		if got != want {
			t.Errorf("parseGUID(%q) = %v, want %v", s, got, want)
		}
		if f, w := formatGUID(got), want.String(); f != w {
			t.Errorf("formatGUID(%q) = %s, want %s", s, f, w)
		}
	}
}

// TestParseGUIDRejectsMalformed checks the shape is enforced.
func TestParseGUIDRejectsMalformed(t *testing.T) {
	for _, s := range []string{
		"",
		"00000000-0000-0000-C000-000000000046",
		"{00000000-0000-0000-C000-00000000004}",
		"{00000000-0000-0000-C000-0000000000467}",
		"{00000000_0000-0000-C000-000000000046}",
		"{0000000G-0000-0000-C000-000000000046}",
	} {
		if _, err := parseGUID(s); err == nil {
			t.Errorf("parseGUID(%q) accepted a malformed GUID", s)
		}
	}
}
