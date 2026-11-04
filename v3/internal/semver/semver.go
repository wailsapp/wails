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
		if meta == "" {
			return nil, fmt.Errorf("invalid semver metadata: empty identifier")
		}
		s = s[:idx]
	}

	// Split prerelease
	pre := ""
	if idx := strings.Index(s, "-"); idx != -1 {
		pre = s[idx+1:]
		if pre == "" {
			return nil, fmt.Errorf("invalid semver prerelease: empty identifier")
		}
		s = s[:idx]
	}

	majorPart, rest, _ := strings.Cut(s, ".")
	minorPart, patchPart, ok := strings.Cut(rest, ".")
	if !ok {
		return nil, fmt.Errorf("invalid semver %q: expected MAJOR.MINOR.PATCH", s)
	}

	major, err := strconv.ParseUint(majorPart, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid semver major %q: %w", majorPart, err)
	}
	minor, err := strconv.ParseUint(minorPart, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid semver minor %q: %w", minorPart, err)
	}
	patch, err := strconv.ParseUint(patchPart, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid semver patch %q: %w", patchPart, err)
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
	var buffer [64]byte
	text := strconv.AppendUint(buffer[:0], v.Major, 10)
	text = append(text, '.')
	text = strconv.AppendUint(text, v.Minor, 10)
	text = append(text, '.')
	text = strconv.AppendUint(text, v.Patch, 10)
	if v.Prerelease == "" && v.Metadata == "" {
		return string(text)
	}
	size := len(text) + len(v.Prerelease) + len(v.Metadata)
	if v.Prerelease != "" {
		size++
	}
	if v.Metadata != "" {
		size++
	}
	var output strings.Builder
	output.Grow(size)
	output.Write(text)
	if v.Prerelease != "" {
		output.WriteByte('-')
		output.WriteString(v.Prerelease)
	}
	if v.Metadata != "" {
		output.WriteByte('+')
		output.WriteString(v.Metadata)
	}
	return output.String()
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
	for {
		ai, aRest, aMore := strings.Cut(a, ".")
		bi, bRest, bMore := strings.Cut(b, ".")
		aNumeric := isNumericIdentifier(ai)
		bNumeric := isNumericIdentifier(bi)
		switch {
		case aNumeric && bNumeric:
			// Compare decimal strings without imposing an integer size limit.
			ai, bi = strings.TrimLeft(ai, "0"), strings.TrimLeft(bi, "0")
			if c := cmpUint(uint64(len(ai)), uint64(len(bi))); c != 0 {
				return c
			}
		case aNumeric:
			return -1 // numeric < alphanumeric
		case bNumeric:
			return 1
		}
		if c := strings.Compare(ai, bi); c != 0 {
			return c
		}
		if !aMore {
			if bMore {
				return -1
			}
			return 0
		}
		if !bMore {
			return 1
		}
		a, b = aRest, bRest
	}
}

func isNumericIdentifier(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
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
