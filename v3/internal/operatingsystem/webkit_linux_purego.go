//go:build linux && purego && !gtk3 && !android

package operatingsystem

import (
	"fmt"
	"sync"

	"github.com/ebitengine/purego"
)

type WebkitVersion struct {
	Major uint
	Minor uint
	Micro uint
}

var (
	webkitVersionOnce sync.Once

	webkit_get_major_version func() uint32
	webkit_get_minor_version func() uint32
	webkit_get_micro_version func() uint32
)

// loadWebkitVersionFuncs binds the three version getters from the runtime
// WebKitGTK library. On failure the funcs stay nil and GetWebkitVersion
// reports 0.0.0 — this is diagnostic-only code (wails doctor), so a missing
// library must not crash it.
func loadWebkitVersionFuncs() {
	webkitVersionOnce.Do(func() {
		var lib uintptr
		for _, name := range []string{"libwebkitgtk-6.0.so.4", "libwebkitgtk-6.0.so"} {
			handle, err := purego.Dlopen(name, purego.RTLD_NOW|purego.RTLD_GLOBAL)
			if err == nil && handle != 0 {
				lib = handle
				break
			}
		}
		if lib == 0 {
			return
		}
		reg := func(fptr any, name string) {
			sym, err := purego.Dlsym(lib, name)
			if err == nil && sym != 0 {
				purego.RegisterFunc(fptr, sym)
			}
		}
		reg(&webkit_get_major_version, "webkit_get_major_version")
		reg(&webkit_get_minor_version, "webkit_get_minor_version")
		reg(&webkit_get_micro_version, "webkit_get_micro_version")
	})
}

func GetWebkitVersion() WebkitVersion {
	loadWebkitVersionFuncs()
	if webkit_get_major_version == nil ||
		webkit_get_minor_version == nil ||
		webkit_get_micro_version == nil {
		return WebkitVersion{}
	}
	return WebkitVersion{
		Major: uint(webkit_get_major_version()),
		Minor: uint(webkit_get_minor_version()),
		Micro: uint(webkit_get_micro_version()),
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
