// Package idlversion implements WebView2 SDK / runtime version parsing and
// comparison. The version strings used by Microsoft mix the SDK number
// (e.g. "1.0.2903.40") and the runtime number (e.g. "121.0.2277.83") in
// the same format; both are handled here.
package idlversion

import (
	"fmt"
	"strconv"
	"strings"
)

// Version is a parsed WebView2-style version (major.minor.patch.build with
// an optional channel suffix such as "1.0.515-prerelease").
type Version struct {
	Major, Minor, Patch, Build int
	Channel                    string
}

// Parse parses a version string. Returns an error if the version has more
// than four numeric segments or any segment fails to parse.
func Parse(v string) (Version, error) {
	var p Version

	if i := strings.Index(v, " "); i > 0 {
		p.Channel = v[i+1:]
		v = v[:i]
	}
	// "1.0.515-prerelease" — strip the trailing channel suffix on the build.
	if i := strings.Index(v, "-"); i > 0 {
		p.Channel = v[i+1:]
		v = v[:i]
	}

	parts := strings.Split(v, ".")
	if len(parts) > 4 {
		return p, fmt.Errorf("too many version parts in %q", v)
	}

	targets := []*int{&p.Major, &p.Minor, &p.Patch, &p.Build}
	for i, s := range parts {
		n, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return p, fmt.Errorf("bad version segment %d in %q: %w", i, v, err)
		}
		*targets[i] = int(n)
	}

	return p, nil
}

// String renders the version back to its canonical form.
func (v Version) String() string {
	s := fmt.Sprintf("%d.%d.%d.%d", v.Major, v.Minor, v.Patch, v.Build)
	if v.Channel != "" {
		s += " " + v.Channel
	}
	return s
}

// Compare returns -1, 0, or +1 to express v <, =, or > w. The channel
// suffix is ignored (the runtime treats "prerelease" identically to
// stable when numerics match — matches webviewloader.CompareBrowserVersions).
func (v Version) Compare(w Version) int {
	for _, pair := range [][2]int{
		{v.Major, w.Major},
		{v.Minor, w.Minor},
		{v.Patch, w.Patch},
		{v.Build, w.Build},
	} {
		if pair[0] < pair[1] {
			return -1
		}
		if pair[0] > pair[1] {
			return 1
		}
	}
	return 0
}

// Compare is a convenience wrapper that parses both strings.
func Compare(a, b string) (int, error) {
	va, err := Parse(a)
	if err != nil {
		return 0, fmt.Errorf("first version: %w", err)
	}
	vb, err := Parse(b)
	if err != nil {
		return 0, fmt.Errorf("second version: %w", err)
	}
	return va.Compare(vb), nil
}
