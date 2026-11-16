package semver

import (
	"fmt"
	"strings"
	"testing"
)

func BenchmarkVersionString(b *testing.B) {
	for _, input := range []string{
		"3.0.1",
		"10.10.10",
		"v3.1.0-nightly.20261029",
		"v3.1.0-nightly.20261029+build.abcdef",
		"18446744073709551615.18446744073709551615.18446744073709551615",
		"3.1.0-nightly+" + strings.Repeat("build", 64),
	} {
		b.Run(input, func(b *testing.B) {
			version, err := NewVersion(input)
			if err != nil {
				b.Fatal(err)
			}
			want := strings.TrimPrefix(input, "v")
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				if got := version.String(); got != want {
					b.Fatalf("String()=%q, want %q", got, want)
				}
			}
		})
	}
}

func FuzzVersionString(f *testing.F) {
	f.Add(uint64(3), uint64(0), uint64(1), "", "")
	f.Add(^uint64(0), ^uint64(0), ^uint64(0), "nightly", strings.Repeat("build", 64))
	f.Add(uint64(0), uint64(0), uint64(0), "日本語", "\x00")
	f.Fuzz(func(t *testing.T, major, minor, patch uint64, pre, meta string) {
		version := Version{Major: major, Minor: minor, Patch: patch, Prerelease: pre, Metadata: meta}
		want := fmt.Sprintf("%d.%d.%d", major, minor, patch)
		if pre != "" {
			want += "-" + pre
		}
		if meta != "" {
			want += "+" + meta
		}
		if got := version.String(); got != want {
			t.Fatalf("String()=%q, want %q", got, want)
		}
	})
}
