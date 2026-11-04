//go:build !windows && !darwin

package reportcopy

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func BenchmarkWriteReport(b *testing.B) {
	directory := b.TempDir()
	if err := os.WriteFile(filepath.Join(directory, "wl-copy"), []byte("#!/bin/sh\n/bin/cat > /dev/null\n"), 0755); err != nil {
		b.Fatal(err)
	}
	b.Setenv("PATH", directory)
	b.Setenv("WAYLAND_DISPLAY", "benchmark")
	text := strings.Repeat("Wails doctor report\n", 100)
	b.ReportAllocs()
	for b.Loop() {
		if err := Write(text); err != nil {
			b.Fatal(err)
		}
	}
}
