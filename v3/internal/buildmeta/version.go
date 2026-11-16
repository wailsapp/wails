// Package buildmeta resolves application version inputs without executing source code.
package buildmeta

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/mod/semver"
)

type Info struct {
	Version string
	Core    string
	Commit  string
}

func (info Info) Prerelease() bool { return semver.Prerelease("v"+info.Version) != "" }

func Parse(value string) (Info, error) {
	version := "v" + strings.TrimPrefix(value, "v")
	core, _, _ := strings.Cut(strings.TrimPrefix(version, "v"), "-")
	core, _, _ = strings.Cut(core, "+")
	if !semver.IsValid(version) || len(strings.Split(core, ".")) != 3 {
		return Info{}, fmt.Errorf("invalid application semantic version %q (expected major.minor.patch)", value)
	}
	return Info{Version: strings.TrimPrefix(version, "v"), Core: core}, nil
}

// Windows returns a four-component numeric resource version. The complete
// semantic version remains available separately for display and comparison.
func (info Info) Windows() ([4]uint16, error) {
	var result [4]uint16
	parts := strings.Split(info.Core, ".")
	if len(parts) != 3 {
		return result, fmt.Errorf("invalid numeric version %q", info.Core)
	}
	for i, part := range parts {
		value, err := strconv.ParseUint(part, 10, 16)
		if err != nil {
			return result, fmt.Errorf("Windows version component %q must be between 0 and 65535", part)
		}
		result[i] = uint16(value)
	}
	return result, nil
}
