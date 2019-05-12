package cmd

import (
	"github.com/masterminds/semver"
)

// SemanticVersion is a struct containing a semantic version
type SemanticVersion struct {
	Version *semver.Version
}

// NewSemanticVersion creates a new SemanticVersion object with the given version string
func NewSemanticVersion(version string) (*SemanticVersion, error) {
	semverVersion, err := semver.NewVersion(version)
	if err != nil {
		return nil, err
	}
	return &SemanticVersion{
		Version: semverVersion,
	}, nil
}

// IsRelease returns true if it's a release version
func (s *SemanticVersion) IsRelease() bool {
	return len(s.Version.Prerelease()) == 0 && len(s.Version.Metadata()) == 0
}

// IsPreRelease returns true if it's a prerelease version
func (s *SemanticVersion) IsPreRelease() bool {
	return len(s.Version.Prerelease()) > 0
}

func (s *SemanticVersion) String() string {
	return s.Version.String()
}

// IsGreaterThan returns true if this version is greater than the given version
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
