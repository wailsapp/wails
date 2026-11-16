//go:build linux

package libpath

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindLibraryFileFallback(t *testing.T) {
	bin := t.TempDir()
	if err := os.WriteFile(filepath.Join(bin, "ldconfig"), []byte("#!/bin/sh\nexit 1\n"), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin)
	for _, tt := range []struct {
		name      string
		request   string
		candidate string
		isDir     bool
		wantFound bool
	}{
		{"exact", "libwails-fallback-test.so", "libwails-fallback-test.so", false, true},
		{"version suffix", "libwails-fallback-test.so", "libwails-fallback-test.so.2", false, true},
		{"subversion", "libwails-fallback-test.so.2", "libwails-fallback-test.so.2.1", false, true},
		{"different ABI", "libwails-fallback-test.so.2", "libwails-fallback-test.so.20", false, false},
		{"backup", "libwails-fallback-test.so", "libwails-fallback-test.so.backup", false, false},
		{"exact directory", "libwails-fallback-test.so", "libwails-fallback-test.so", true, false},
		{"version directory", "libwails-fallback-test.so", "libwails-fallback-test.so.2", true, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			t.Setenv("LD_LIBRARY_PATH", dir)
			candidate := filepath.Join(dir, tt.candidate)
			var err error
			if tt.isDir {
				err = os.Mkdir(candidate, 0755)
			} else {
				err = os.WriteFile(candidate, nil, 0644)
			}
			if err != nil {
				t.Fatal(err)
			}
			want := ""
			if tt.wantFound {
				want = candidate
			}
			got, err := FindLibraryFile(tt.request)
			if got != want || (err == nil) != tt.wantFound {
				t.Errorf("FindLibraryFile = %q, %v; want %q, found %v", got, err, want, tt.wantFound)
			}
		})
	}
}
