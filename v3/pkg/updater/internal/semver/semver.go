// Package semver wraps golang.org/x/mod/semver with a small, provider-friendly
// interface. Providers (GitHub Releases, AppCast, keygen.sh) all need the same
// "is the just-fetched tag newer than the running version" comparison; before
// this package each rolled its own hand-tuned numeric-segment comparator with
// subtly wrong prerelease handling.
//
// All exported helpers tolerate the conventional leading "v"/"V" prefix.
package semver

import (
	"strings"

	xsemver "golang.org/x/mod/semver"
)

// TrimPrefix strips a leading "v" or "V" from tag — the conventional release
// tag form used by GitHub, Sparkle AppCast, and most Go projects.
func TrimPrefix(tag string) string {
	if strings.HasPrefix(tag, "v") || strings.HasPrefix(tag, "V") {
		return tag[1:]
	}
	return tag
}

// Compare returns -1, 0, or +1 reflecting whether a < b, a == b, or a > b
// under SemVer 2.0.0 precedence rules. Either input may include the
// conventional leading "v"/"V".
//
// SemVer rules to keep in mind for callers:
//   - Build metadata (+...) is ignored.
//   - A prerelease ("-rc.1", "-alpha") is *lower* precedence than the
//     corresponding release.
//   - Invalid semver strings are all equal to each other and less than every
//     valid semver string, matching x/mod/semver's stable behaviour.
func Compare(a, b string) int {
	return xsemver.Compare(canonical(a), canonical(b))
}

// IsNewer reports whether tag is strictly newer than current. Returns false
// when tag is empty or equal to current; returns true when tag is non-empty
// and current is empty.
func IsNewer(tag, current string) bool {
	if TrimPrefix(tag) == "" {
		return false
	}
	if TrimPrefix(current) == "" {
		return true
	}
	return Compare(tag, current) > 0
}

// canonical converts the input into the form x/mod/semver expects: a single
// leading "v". Inputs already in that form pass through unchanged.
func canonical(s string) string {
	s = TrimPrefix(s)
	if s == "" {
		return ""
	}
	return "v" + s
}
