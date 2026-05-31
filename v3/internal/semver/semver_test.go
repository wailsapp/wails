package semver

import (
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// NewVersion – parse paths
// ---------------------------------------------------------------------------

func TestNewVersion_Valid(t *testing.T) {
	cases := []struct {
		input             string
		major, minor, patch uint64
		pre, meta         string
	}{
		{"1.2.3", 1, 2, 3, "", ""},
		{"v1.2.3", 1, 2, 3, "", ""},
		{"0.0.0", 0, 0, 0, "", ""},
		{"10.20.30", 10, 20, 30, "", ""},
		{"1.2.3-alpha", 1, 2, 3, "alpha", ""},
		{"1.2.3-alpha.1", 1, 2, 3, "alpha.1", ""},
		{"1.2.3+build.1", 1, 2, 3, "", "build.1"},
		{"1.2.3-beta.2+exp.sha", 1, 2, 3, "beta.2", "exp.sha"},
		{"v2.0.0-rc.1+meta", 2, 0, 0, "rc.1", "meta"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			v, err := NewVersion(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if v.Major != tc.major || v.Minor != tc.minor || v.Patch != tc.patch {
				t.Errorf("core: got %d.%d.%d want %d.%d.%d",
					v.Major, v.Minor, v.Patch, tc.major, tc.minor, tc.patch)
			}
			if v.Prerelease != tc.pre {
				t.Errorf("pre: got %q want %q", v.Prerelease, tc.pre)
			}
			if v.Metadata != tc.meta {
				t.Errorf("meta: got %q want %q", v.Metadata, tc.meta)
			}
		})
	}
}

func TestNewVersion_Invalid(t *testing.T) {
	cases := []struct {
		input   string
		wantMsg string
	}{
		// wrong part count (only 2 parts)
		{"1.2", "MAJOR.MINOR.PATCH"},
		// only 1 part
		{"1", "MAJOR.MINOR.PATCH"},
		// empty string
		{"", "MAJOR.MINOR.PATCH"},
		// bad major
		{"x.1.2", "major"},
		// bad minor
		{"1.x.2", "minor"},
		// bad patch
		{"1.2.x", "patch"},
		// SplitN(s, ".", 3) makes parts[2]="3.4" which is non-numeric → patch error
		{"1.2.3.4", "patch"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			_, err := NewVersion(tc.input)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantMsg) {
				t.Errorf("error %q does not contain %q", err.Error(), tc.wantMsg)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// String – round-trip and individual branches
// ---------------------------------------------------------------------------

func TestString(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		// no pre, no meta
		{"1.2.3", "1.2.3"},
		// with pre, no meta
		{"1.2.3-alpha.1", "1.2.3-alpha.1"},
		// no pre, with meta (meta is stripped from comparison but preserved in String)
		{"1.2.3+build.42", "1.2.3+build.42"},
		// both pre and meta
		{"1.2.3-beta+exp", "1.2.3-beta+exp"},
		// v-prefix stripped
		{"v4.5.6", "4.5.6"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			v, err := NewVersion(tc.input)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			if got := v.String(); got != tc.want {
				t.Errorf("String() = %q, want %q", got, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// compare – all branch paths
// ---------------------------------------------------------------------------

func mustParse(t *testing.T, s string) *Version {
	t.Helper()
	v, err := NewVersion(s)
	if err != nil {
		t.Fatalf("mustParse(%q): %v", s, err)
	}
	return v
}

func TestCompare_MajorDiffers(t *testing.T) {
	a := mustParse(t, "2.0.0")
	b := mustParse(t, "1.0.0")
	if !a.GreaterThan(b) {
		t.Error("2.0.0 should be > 1.0.0")
	}
	if !b.LessThan(a) {
		t.Error("1.0.0 should be < 2.0.0")
	}
}

func TestCompare_MinorDiffers(t *testing.T) {
	a := mustParse(t, "1.3.0")
	b := mustParse(t, "1.2.0")
	if !a.GreaterThan(b) {
		t.Error("1.3.0 should be > 1.2.0")
	}
	if !b.LessThan(a) {
		t.Error("1.2.0 should be < 1.3.0")
	}
}

func TestCompare_PatchDiffers(t *testing.T) {
	a := mustParse(t, "1.2.4")
	b := mustParse(t, "1.2.3")
	if !a.GreaterThan(b) {
		t.Error("1.2.4 should be > 1.2.3")
	}
	if !b.LessThan(a) {
		t.Error("1.2.3 should be < 1.2.4")
	}
}

// Both prerelease empty → equal
func TestCompare_BothPreEmpty(t *testing.T) {
	a := mustParse(t, "1.0.0")
	b := mustParse(t, "1.0.0")
	if !a.Equal(b) {
		t.Error("1.0.0 should equal 1.0.0")
	}
	if !a.GreaterThanOrEqual(b) {
		t.Error("1.0.0 >= 1.0.0 should be true")
	}
}

// Only v.Prerelease empty → v is greater (normal release > pre-release)
func TestCompare_OnlyVPreEmpty(t *testing.T) {
	release := mustParse(t, "1.0.0")
	pre := mustParse(t, "1.0.0-alpha")
	if !release.GreaterThan(pre) {
		t.Error("1.0.0 should be > 1.0.0-alpha")
	}
}

// Only other.Prerelease empty → v is less
func TestCompare_OnlyOtherPreEmpty(t *testing.T) {
	pre := mustParse(t, "1.0.0-alpha")
	release := mustParse(t, "1.0.0")
	if !pre.LessThan(release) {
		t.Error("1.0.0-alpha should be < 1.0.0")
	}
}

// Both prerelease non-empty → comparePrerelease used
func TestCompare_BothPreNonEmpty(t *testing.T) {
	a := mustParse(t, "1.0.0-beta")
	b := mustParse(t, "1.0.0-alpha")
	if !a.GreaterThan(b) {
		t.Error("1.0.0-beta should be > 1.0.0-alpha")
	}
}

// ---------------------------------------------------------------------------
// comparePrerelease – all §11.4 branches
// ---------------------------------------------------------------------------

func TestComparePrerelease(t *testing.T) {
	type tc struct {
		a, b string
		want int // -1, 0, 1
	}
	cases := []tc{
		// numeric vs numeric, same → continue to next identifier
		{"1.2", "1.3", -1},
		{"1.3", "1.2", 1},
		{"1.1", "1.1", 0},

		// numeric < alphanumeric
		{"1", "alpha", -1},

		// alphanumeric > numeric (bErr == nil branch)
		{"alpha", "1", 1},

		// alphanumeric vs alphanumeric, lexical
		{"alpha", "beta", -1},
		{"beta", "alpha", 1},
		{"alpha", "alpha", 0},

		// alpha==alpha, continue to next identifier
		{"alpha.1", "alpha.2", -1},
		{"alpha.2", "alpha.1", 1},
		{"alpha.1", "alpha.1", 0},

		// longer wins when prefix equal
		{"alpha", "alpha.1", -1},  // shorter < longer
		{"alpha.1", "alpha", 1},   // longer > shorter

		// numeric longer wins
		{"1", "1.1", -1},
		{"1.1", "1", 1},

		// SemVer §11.4 example sequence
		{"alpha", "alpha.1", -1},
		{"alpha.1", "alpha.beta", -1},
		{"alpha.beta", "beta", -1},
		{"beta", "beta.2", -1},
		{"beta.2", "beta.11", -1},
		{"beta.11", "rc.1", -1},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.a+"_vs_"+tc.b, func(t *testing.T) {
			got := comparePrerelease(tc.a, tc.b)
			// normalise to -1/0/1
			if got < 0 {
				got = -1
			} else if got > 0 {
				got = 1
			}
			if got != tc.want {
				t.Errorf("comparePrerelease(%q, %q) = %d, want %d", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Full version ordering (integration)
// ---------------------------------------------------------------------------

func TestVersionOrdering(t *testing.T) {
	// Ordered from lowest to highest
	versions := []string{
		"1.0.0-alpha",
		"1.0.0-alpha.1",
		"1.0.0-alpha.beta",
		"1.0.0-beta",
		"1.0.0-beta.2",
		"1.0.0-beta.11",
		"1.0.0-rc.1",
		"1.0.0",
	}
	parsed := make([]*Version, len(versions))
	for i, s := range versions {
		parsed[i] = mustParse(t, s)
	}
	for i := 0; i < len(parsed)-1; i++ {
		a, b := parsed[i], parsed[i+1]
		if !a.LessThan(b) {
			t.Errorf("%s should be < %s", versions[i], versions[i+1])
		}
		if !b.GreaterThan(a) {
			t.Errorf("%s should be > %s", versions[i+1], versions[i])
		}
		if a.Equal(b) {
			t.Errorf("%s should not equal %s", versions[i], versions[i+1])
		}
		if !b.GreaterThanOrEqual(a) {
			t.Errorf("%s >= %s should be true", versions[i+1], versions[i])
		}
	}
}

// Metadata must be ignored in comparisons.
func TestMetadataIgnored(t *testing.T) {
	a := mustParse(t, "1.0.0+build.1")
	b := mustParse(t, "1.0.0+build.2")
	if !a.Equal(b) {
		t.Error("versions differing only by metadata should be equal")
	}
}

// cmpUint equal branch (indirectly via compare on identical versions).
func TestCmpUintEqual(t *testing.T) {
	a := mustParse(t, "3.3.3")
	b := mustParse(t, "3.3.3")
	if got := a.compare(b); got != 0 {
		t.Errorf("compare identical versions = %d, want 0", got)
	}
}
