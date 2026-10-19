//go:build !windows

package doctor

import (
	"os"
	"path/filepath"
	"testing"
)

// TestCommonDependenciesPreserveFailures guards earlier platform/sandbox failures.
func TestCommonDependenciesPreserveFailures(t *testing.T) {
	for _, version := range []string{"11.0.0", "6.0.0"} {
		for _, initiallyReady := range []bool{false, true} {
			dir := t.TempDir()
			if err := os.WriteFile(filepath.Join(dir, "npm"), []byte("#!/bin/sh\necho "+version+"\n"), 0755); err != nil {
				t.Fatal(err)
			}
			t.Setenv("PATH", dir)
			t.Setenv("ANDROID_HOME", "")
			t.Setenv("ANDROID_SDK_ROOT", "")
			t.Setenv("HOME", dir)
			ready := initiallyReady
			result := map[string]string{}
			checkCommonDependencies(result, &ready)
			want := initiallyReady && version == "11.0.0"
			if ready != want {
				t.Errorf("npm=%s initiallyReady=%v: ready=%v, want %v", version, initiallyReady, ready, want)
			}
		}
	}
}
