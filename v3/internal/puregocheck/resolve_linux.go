package main

import (
	"fmt"
	"unsafe"

	"github.com/ebitengine/purego"
)

func resolve(r *report) error {
	loadLibraries(r)
	r.Loaded = loadedVersions()
	resolveSymbols(r)
	return nil
}

// loadedVersions asks each loaded library for its own version, which is what
// decides the functions the backend may call on this system.
func loadedVersions() []string {
	var out []string
	call := func(prefix string) (string, bool) {
		var v [3]uint32
		for i, part := range []string{"major", "minor", "micro"} {
			sym, err := purego.Dlsym(purego.RTLD_DEFAULT, prefix+"_get_"+part+"_version")
			if err != nil || sym == 0 {
				return "", false
			}
			var fn func() uint32
			purego.RegisterFunc(&fn, sym)
			v[i] = fn()
		}
		return fmt.Sprintf("%d.%d.%d", v[0], v[1], v[2]), true
	}
	for _, lib := range []struct{ name, prefix string }{{"gtk", "gtk"}, {"webkit", "webkit"}, {"jsc", "jsc"}, {"soup", "soup"}} {
		if v, ok := call(lib.prefix); ok {
			out = append(out, lib.name+" "+v)
		}
	}
	// GLib exports its version as data, not functions.
	var v [3]uint32
	for i, part := range []string{"major", "minor", "micro"} {
		sym, err := purego.Dlsym(purego.RTLD_DEFAULT, "glib_"+part+"_version")
		if err != nil || sym == 0 {
			return out
		}
		v[i] = *(*uint32)(unsafe.Pointer(sym))
	}
	return append(out, fmt.Sprintf("glib %d.%d.%d", v[0], v[1], v[2]))
}
