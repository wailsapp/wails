package github

import (
	"github.com/matryer/is"
	"testing"
)

func TestSemanticVersion_IsGreaterThan(t *testing.T) {
	is2 := is.New(t)

	alpha1, err := NewSemanticVersion("v3.0.0-alpha.1")
	is2.NoErr(err)

	beta1, err := NewSemanticVersion("v3.0.0-beta.1")
	is2.NoErr(err)

	v2, err := NewSemanticVersion("v3.0.0")
	is2.NoErr(err)

	is2.True(alpha1.IsPreRelease())
	is2.True(beta1.IsPreRelease())
	is2.True(!v2.IsPreRelease())
	is2.True(v2.IsRelease())

	result, err := beta1.IsGreaterThan(alpha1)
	is2.NoErr(err)
	is2.True(result)

	result, err = v2.IsGreaterThan(beta1)
	is2.NoErr(err)
	is2.True(result)

	beta44, err := NewSemanticVersion("v2.0.0-beta.44.2")
	is2.NoErr(err)

	rc1, err := NewSemanticVersion("v2.0.0-rc.1")
	is2.NoErr(err)

	result, err = rc1.IsGreaterThan(beta44)
	is2.NoErr(err)
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
