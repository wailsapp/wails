package main

import (
	"strings"
	"testing"
)

// loadedVersions reports only libraries whose version functions it found, as
// "name major.minor.micro". Without GTK installed it reports nothing, which is
// also correct.
func TestLoadedVersionsFormat(t *testing.T) {
	for _, v := range loadedVersions() {
		name, ver, ok := strings.Cut(v, " ")
		if !ok || name == "" || strings.Count(ver, ".") != 2 {
			t.Errorf("loadedVersions entry %q is not 'name x.y.z'", v)
		}
	}
}
