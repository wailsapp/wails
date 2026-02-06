//go:build linux && cgo && gtk4 && !android

package operatingsystem

/*
#cgo linux pkg-config: gtk4 webkitgtk-6.0
#include <webkit/webkit.h>
*/
import "C"
import "fmt"

type WebkitVersion struct {
	Major uint
	Minor uint
	Micro uint
}

func GetWebkitVersion() WebkitVersion {
	var major, minor, micro C.uint
	major = C.webkit_get_major_version()
	minor = C.webkit_get_minor_version()
	micro = C.webkit_get_micro_version()
	return WebkitVersion{
		Major: uint(major),
		Minor: uint(minor),
		Micro: uint(micro),
	}
}

func (v WebkitVersion) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Micro)
}

func (v WebkitVersion) IsAtLeast(major int, minor int, micro int) bool {
	if v.Major != uint(major) {
		return v.Major > uint(major)
	}
	if v.Minor != uint(minor) {
		return v.Minor > uint(minor)
	}
	return v.Micro >= uint(micro)
}
