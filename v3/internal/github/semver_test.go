package github

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSemanticVersion_IsGreaterThan(t *testing.T) {
	is2 := require.New(t)

	alpha1, err := NewSemanticVersion("v3.0.0-alpha.1")
	is2.NoError(err)

	beta1, err := NewSemanticVersion("v3.0.0-beta.1")
	is2.NoError(err)

	v2, err := NewSemanticVersion("v3.0.0")
	is2.NoError(err)

	is2.True(alpha1.IsPreRelease())
	is2.True(beta1.IsPreRelease())
	is2.True(!v2.IsPreRelease())
	is2.True(v2.IsRelease())

	result, err := beta1.IsGreaterThan(alpha1)
	is2.NoError(err)
	is2.True(result)

	result, err = v2.IsGreaterThan(beta1)
	is2.NoError(err)
	is2.True(result)

	beta44, err := NewSemanticVersion("v2.0.0-beta.44.2")
	is2.NoError(err)

	rc1, err := NewSemanticVersion("v2.0.0-rc.1")
	is2.NoError(err)

	result, err = rc1.IsGreaterThan(beta44)
	is2.NoError(err)
	is2.True(result)

}

func TestReleaseShapes(t *testing.T) {
	cases := []struct {
		v                       string
		weekly, hotfix, nightly bool
	}{
		{"v3.0.0", true, false, false},
		{"v3.1.0", true, false, false},
		{"v3.1.1", false, true, false},
		{"v3.2.0-nightly.20261102", false, false, true},
		{"v3.2.0-nightly.20261102.1", false, false, true},
		{"v3.0.0-beta.34", false, false, false},
		{"v2.12.0", false, false, false},
	}
	for _, c := range cases {
		s, err := NewSemanticVersion(c.v)
		if err != nil {
			t.Fatal(err)
		}
		if s.IsWeekly() != c.weekly || s.IsHotfix() != c.hotfix || s.IsNightly() != c.nightly {
			t.Errorf("%s: weekly=%v hotfix=%v nightly=%v", c.v, s.IsWeekly(), s.IsHotfix(), s.IsNightly())
		}
	}
	// Go @latest simulation: the newest release must beat every nightly of itself and earlier.
	w, _ := NewSemanticVersion("v3.1.0")
	n, _ := NewSemanticVersion("v3.1.0-nightly.20261031")
	if ok, _ := w.IsGreaterThan(n); !ok {
		t.Error("weekly must sort above its nightlies")
	}
}

func TestSemanticVersionReleaseClassification(t *testing.T) {
	for _, tc := range []struct {
		input               string
		release, prerelease bool
	}{
		{"v3.0.0", true, false},
		{"v3.0.0+build.123", true, false},
		{"v3.0.0+nightly", true, false},
		{"v3.0.0-beta.1+build.123", false, true},
		{"v2.0.0+build.123", false, false},
	} {
		t.Run(tc.input, func(t *testing.T) {
			version, err := NewSemanticVersion(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if version.IsRelease() != tc.release || version.IsPreRelease() != tc.prerelease {
				t.Errorf("release=%t, prerelease=%t; want %t, %t", version.IsRelease(), version.IsPreRelease(), tc.release, tc.prerelease)
			}
		})
	}
}
