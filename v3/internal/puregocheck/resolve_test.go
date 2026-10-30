//go:build darwin || linux

package main

import (
	"runtime"
	"strings"
	"testing"
)

// platformLibC is a library every process on this OS has loaded or can load.
func platformLibC() string {
	if runtime.GOOS == "darwin" {
		return "/usr/lib/libSystem.B.dylib"
	}
	return "libc.so.6"
}

func TestLoadLibrariesRecordsOnlyLoadable(t *testing.T) {
	r := &report{libCandidates: []string{"definitely-missing-lib.so.9", platformLibC()}}
	loadLibraries(r)
	if len(r.Libraries) != 1 || r.Libraries[0] != platformLibC() {
		t.Errorf("Libraries = %v, want only %s", r.Libraries, platformLibC())
	}
}

func TestResolveSymbolsRequiredAndOptional(t *testing.T) {
	r := &report{Symbols: []*entry{
		{Name: "malloc", Where: []string{"a.go:1"}},
		{Name: "no_such_symbol_puregocheck", Where: []string{"a.go:2"}},
		{Name: "no_such_optional_puregocheck", Optional: true, Where: []string{"a.go:3"}},
	}}
	loadLibraries(&report{libCandidates: []string{platformLibC()}})
	resolveSymbols(r)
	if len(r.Missing) != 1 || !strings.Contains(r.Missing[0], "no_such_symbol_puregocheck") || !strings.Contains(r.Missing[0], "a.go:2") {
		t.Errorf("Missing = %v", r.Missing)
	}
	if len(r.Absent) != 1 || r.Absent[0] != "symbol no_such_optional_puregocheck" {
		t.Errorf("Absent = %v", r.Absent)
	}
}
