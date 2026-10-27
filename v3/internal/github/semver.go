package github

import (
	"fmt"
	"strings"

	"github.com/wailsapp/wails/v3/internal/semver"
)

const majorVersion = 3

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
	if s.Version.Major != majorVersion {
		return false
	}
	return len(s.Version.Prerelease) == 0 && len(s.Version.Metadata) == 0
}

// IsPreRelease returns true if it's a prerelease version
func (s *SemanticVersion) IsPreRelease() bool {
	if s.Version.Major != majorVersion {
		return false
	}
	return len(s.Version.Prerelease) > 0
}

// IsWeekly returns true for a weekly release (vX.Y.0, no suffix). Minor
// numbers advance with time, not features.
func (s *SemanticVersion) IsWeekly() bool {
	return s.IsRelease() && s.Version.Patch == 0
}

// IsHotfix returns true for a hotfix release (vX.Y.Z, Z > 0, no suffix).
func (s *SemanticVersion) IsHotfix() bool {
	return s.IsRelease() && s.Version.Patch > 0
}

// IsNightly returns true for a nightly pre-release (vX.Y.0-nightly.DATE[.N]).
func (s *SemanticVersion) IsNightly() bool {
	return s.IsPreRelease() && strings.HasPrefix(s.Version.Prerelease, "nightly")
}

func (s *SemanticVersion) String() string {
	return s.Version.String()
}

// IsGreaterThan returns true if this version is greater than the given version
func (s *SemanticVersion) IsGreaterThan(version *SemanticVersion) (bool, error) {
	return s.Version.GreaterThan(version.Version), nil
}

// IsGreaterThanOrEqual returns true if this version is greater than or equal the given version
func (s *SemanticVersion) IsGreaterThanOrEqual(version *SemanticVersion) (bool, error) {
	return s.Version.GreaterThanOrEqual(version.Version), nil
}

// MainVersion returns the main version of any version+prerelease+metadata
// EG: MainVersion("1.2.3-pre") => "1.2.3"
func (s *SemanticVersion) MainVersion() *SemanticVersion {
	mainVersion := fmt.Sprintf("%d.%d.%d", s.Version.Major, s.Version.Minor, s.Version.Patch)
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
