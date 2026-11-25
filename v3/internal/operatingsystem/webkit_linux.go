//go:build linux && cgo && !gtk3 && !android && !purego

package operatingsystem

/*
#cgo linux pkg-config: gtk4 webkitgtk-6.0
#include <webkit/webkit.h>
*/
import "C"

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
