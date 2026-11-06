//go:build darwin || linux

package main

import (
	"github.com/ebitengine/purego"
)

// loadLibraries dlopens every library the backend's sources name. A miss is
// expected (sonames come with fallbacks); what matters is that the symbols
// resolve afterwards.
func loadLibraries(r *report) {
	for _, name := range r.libCandidates {
		if h, err := purego.Dlopen(name, purego.RTLD_NOW|purego.RTLD_GLOBAL); err == nil && h != 0 {
			r.Libraries = append(r.Libraries, name)
		}
	}
}

// resolveSymbols looks up every C symbol in the global namespace of the
// libraries just loaded.
func resolveSymbols(r *report) {
	for _, e := range r.Symbols {
		if sym, err := purego.Dlsym(purego.RTLD_DEFAULT, e.Name); err == nil && sym != 0 {
			continue
		}
		if e.Optional {
			r.Absent = append(r.Absent, "symbol "+e.Name)
		} else {
			r.Missing = append(r.Missing, "symbol "+e.Name+" ("+e.Where[0]+")")
		}
	}
}
