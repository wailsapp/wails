//go:build linux && !android

package operatingsystem

import "fmt"

func (v WebkitVersion) IsAtLeast(major int, minor int, micro int) bool {
	if v.Major != uint(major) {
		return v.Major > uint(major)
	}
	if v.Minor != uint(minor) {
		return v.Minor > uint(minor)
	}
	return v.Micro >= uint(micro)
}

func (v WebkitVersion) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Micro)
}
