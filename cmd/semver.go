package cmd

import (
	"regexp"

	"github.com/masterminds/semver"
)

type SemanticVersion struct {
	Version              *semver.Version
	validPrereleaseRegex *regexp.Regexp
	validReleaseRegex    *regexp.Regexp
}

// NewSemanticVersion creates a new SemanticVersion object with the given version string
func NewSemanticVersion(version string) (*SemanticVersion, error) {
	semverVersion, err := semver.NewVersion(version)
	if err != nil {
		return nil, err
	}
	const SemVerRegex string = `v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?`
	const SemPreVerRegex string = `v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*)`

	return &SemanticVersion{
		Version:              semverVersion,
		validPrereleaseRegex: regexp.MustCompile(SemPreVerRegex),
		validReleaseRegex:    regexp.MustCompile(SemVerRegex),
	}, nil
}

// IsRelease returns true if it's a release version
func (s *SemanticVersion) IsRelease() bool {
	return s.validReleaseRegex.MatchString(s.Version.String())
}

// IsPreRelease returns true if it's a prerelease version
func (s *SemanticVersion) IsPreRelease() bool {
	return s.validPrereleaseRegex.MatchString(s.Version.String())
}

func (s *SemanticVersion) String() string {
	return s.Version.String()
}

func (s *SemanticVersion) IsGreaterThan(version *SemanticVersion) (bool, error) {
	// Set up new constraint
	constraint, err := semver.NewConstraint("> " + version.Version.String())
	if err != nil {
		return false, err
	}

	// Check if the desired one is greater than the requested on
	success, msgs := constraint.Validate(s.Version)
	if !success {
		return false, msgs[0]
	}
	return true, nil
}

// SemverCollection is a collection of SemanticVersion objects
type SemverCollection []*SemanticVersion

// Len returns the length of a collection. The number of Version instances
// on the slice.
func (c SemverCollection) Len() int {
	return len(c)
}

// Less is needed for the sort interface to compare two Version objects on the
// slice. If checks if one is less than the other.
func (c SemverCollection) Less(i, j int) bool {
	return c[i].Version.LessThan(c[j].Version)
}

// Swap is needed for the sort interface to replace the Version objects
// at two different positions in the slice.
func (c SemverCollection) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
