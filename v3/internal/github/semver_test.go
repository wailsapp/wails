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
