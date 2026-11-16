//go:build linux

package libpath

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindFileWithLdconfigMatchesLibraryName(t *testing.T) {
	bin := t.TempDir()
	script := "#!/bin/sh\nprintf '%s\\n' '" +
		"libwails-test-helper.so.1 (libc6,x86-64) => /libs/libwails-test-helper.so.1\n" +
		"libwails-test.so.20 (libc6,x86-64) => /libs/libwails-test.so.20\n" +
		"libwails-test.so.2 (libc6,x86-64) => /libs/libwails-test.so.2\n" +
		"libother.so.1 (libc6,x86-64) => /libs/libwails-missing.so/libother.so.1'\n"
	if err := os.WriteFile(filepath.Join(bin, "ldconfig"), []byte(script), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin)
	for _, tt := range []struct {
		name string
		want string
	}{
		{"libwails-test.so", "/libs/libwails-test.so.20"},
		{"libwails-test.so.2", "/libs/libwails-test.so.2"},
		{"libwails-missing.so", ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findFileWithLdconfig(tt.name)
			if got != tt.want || (err != nil) != (tt.want == "") {
				t.Errorf("findFileWithLdconfig = %q, %v; want %q", got, err, tt.want)
			}
		})
	}
}
