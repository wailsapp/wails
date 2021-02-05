package cmd

import (
	"fmt"

	"github.com/Masterminds/semver"
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
	// Limit to v1
	if s.Version.Major() != 1 {
		return false
	}
	return len(s.Version.Prerelease()) == 0 && len(s.Version.Metadata()) == 0
}

// IsPreRelease returns true if it's a prerelease version
func (s *SemanticVersion) IsPreRelease() bool {
	// Limit to v1
	if s.Version.Major() != 1 {
		return false
	}
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

// IsGreaterThanOrEqual returns true if this version is greater than or equal the given version
func (s *SemanticVersion) IsGreaterThanOrEqual(version *SemanticVersion) (bool, error) {
	// Set up new constraint
	constraint, err := semver.NewConstraint(">= " + version.Version.String())
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

// MainVersion returns the main version of any version+prerelease+metadata
// EG: MainVersion("1.2.3-pre") => "1.2.3"
func (s *SemanticVersion) MainVersion() *SemanticVersion {
	mainVersion := fmt.Sprintf("%d.%d.%d", s.Version.Major(), s.Version.Minor(), s.Version.Patch())
	result, _ := NewSemanticVersion(mainVersion)
	return result
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
