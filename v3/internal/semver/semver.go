// Package semver provides simple semantic version parsing and comparison.
// Implements a subset of the Semantic Versioning 2.0.0 spec sufficient
// for the wails version-check use cases.
package semver

import (
	"fmt"
	"strconv"
	"strings"
)

// Version represents a parsed semantic version.
type Version struct {
	Major      uint64
	Minor      uint64
	Patch      uint64
	Prerelease string
	Metadata   string
}

// NewVersion parses a semver string such as "1.2.3", "1.2.3-beta.1", or "v1.2.3".
func NewVersion(s string) (*Version, error) {
	s = strings.TrimPrefix(s, "v")

	// Split metadata
	meta := ""
	if idx := strings.Index(s, "+"); idx != -1 {
		meta = s[idx+1:]
		s = s[:idx]
	}

	// Split prerelease
	pre := ""
	if idx := strings.Index(s, "-"); idx != -1 {
		pre = s[idx+1:]
		s = s[:idx]
	}

	parts := strings.SplitN(s, ".", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid semver %q: expected MAJOR.MINOR.PATCH", s)
	}

	major, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid semver major %q: %w", parts[0], err)
	}
	minor, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid semver minor %q: %w", parts[1], err)
	}
	patch, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid semver patch %q: %w", parts[2], err)
	}

	return &Version{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		Prerelease: pre,
		Metadata:   meta,
	}, nil
}

// String returns the canonical string representation.
func (v *Version) String() string {
	s := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Prerelease != "" {
		s += "-" + v.Prerelease
	}
	if v.Metadata != "" {
		s += "+" + v.Metadata
	}
	return s
}

// compare returns -1, 0, or 1 depending on whether v is less than, equal to,
// or greater than other. Metadata is ignored per the semver spec.
func (v *Version) compare(other *Version) int {
	if c := cmpUint(v.Major, other.Major); c != 0 {
		return c
	}
	if c := cmpUint(v.Minor, other.Minor); c != 0 {
		return c
	}
	if c := cmpUint(v.Patch, other.Patch); c != 0 {
		return c
	}
	// Pre-release has lower precedence than a normal version.
	switch {
	case v.Prerelease == "" && other.Prerelease == "":
		return 0
	case v.Prerelease == "":
		return 1
	case other.Prerelease == "":
		return -1
	default:
		return comparePrerelease(v.Prerelease, other.Prerelease)
	}
}

// comparePrerelease compares two prerelease strings per SemVer 2.0 §11.4:
// dot-separated identifiers, numeric identifiers compared numerically,
// alphanumeric identifiers compared lexically, numeric < alphanumeric,
// longer wins when all shorter identifiers are equal.
func comparePrerelease(a, b string) int {
	as := strings.Split(a, ".")
	bs := strings.Split(b, ".")
	for i := 0; i < len(as) && i < len(bs); i++ {
		ai, aErr := strconv.ParseUint(as[i], 10, 64)
		bi, bErr := strconv.ParseUint(bs[i], 10, 64)
		switch {
		case aErr == nil && bErr == nil:
			if c := cmpUint(ai, bi); c != 0 {
				return c
			}
		case aErr == nil:
			return -1 // numeric < alphanumeric
		case bErr == nil:
			return 1
		default:
			if c := strings.Compare(as[i], bs[i]); c != 0 {
				return c
			}
		}
	}
	return cmpUint(uint64(len(as)), uint64(len(bs)))
}

func cmpUint(a, b uint64) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

// LessThan reports whether v < other.
func (v *Version) LessThan(other *Version) bool {
	return v.compare(other) < 0
}

// GreaterThan reports whether v > other.
func (v *Version) GreaterThan(other *Version) bool {
	return v.compare(other) > 0
}

// Equal reports whether v == other (metadata ignored).
func (v *Version) Equal(other *Version) bool {
	return v.compare(other) == 0
}

// GreaterThanOrEqual reports whether v >= other.
func (v *Version) GreaterThanOrEqual(other *Version) bool {
	return v.compare(other) >= 0
}

