package semver

import "testing"

func TestIsNewer(t *testing.T) {
	cases := []struct {
		tag, current string
		want         bool
	}{
		// Basic ordering.
		{"v2.0.0", "1.0.0", true},
		{"2.0.0", "v1.0.0", true},
		{"v1.0.0", "v1.0.0", false},
		{"v0.9.0", "v1.0.0", false},

		// Numeric segments must compare numerically, not lexically.
		{"v1.10.0", "v1.2.0", true},
		{"v1.2.0", "v1.10.0", false},

		// SemVer prerelease precedence: a prerelease is older than the
		// corresponding release. Previous hand-rolled comparator got this
		// backwards (treated rc.1 as newer than the release).
		{"v1.2.3-rc.1", "v1.2.3", false},
		{"v1.2.3", "v1.2.3-rc.1", true},
		{"v1.2.3-rc.2", "v1.2.3-rc.1", true},
		{"v1.2.3-beta", "v1.2.3-alpha", true},

		// Build metadata must be ignored.
		{"v1.2.3+build.1", "v1.2.3", false},
		{"v1.2.3+build.2", "v1.2.3+build.1", false},

		// Empty handling.
		{"", "1.0.0", false},
		{"1.0.0", "", true},
		{"v", "1.0.0", false}, // empty after trim
	}
	for _, c := range cases {
		t.Run(c.tag+"_vs_"+c.current, func(t *testing.T) {
			if got := IsNewer(c.tag, c.current); got != c.want {
				t.Errorf("IsNewer(%q, %q) = %v, want %v", c.tag, c.current, got, c.want)
			}
		})
	}
}

func TestTrimPrefix(t *testing.T) {
	cases := map[string]string{
		"v1.2.3": "1.2.3",
		"V1.2.3": "1.2.3",
		"1.2.3":  "1.2.3",
		"":       "",
		"v":      "",
	}
	for in, want := range cases {
		if got := TrimPrefix(in); got != want {
			t.Errorf("TrimPrefix(%q) = %q, want %q", in, got, want)
		}
	}
}
