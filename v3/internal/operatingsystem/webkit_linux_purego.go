//go:build linux && purego && !android

package operatingsystem

import (
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
		for _, name := range []string{puregoWebKitLibrary} {
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
