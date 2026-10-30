package semver

import "testing"

func BenchmarkParseVersion(b *testing.B) {
	for _, input := range []string{
		"3.0.1",
		"v3.1.0-nightly.20261029",
		"v3.1.0-nightly.20261029+build.abcdef",
		"18446744073709551615.18446744073709551615.18446744073709551615",
	} {
		b.Run(input, func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				version, err := NewVersion(input)
				if err != nil || version == nil {
					b.Fatalf("parse %q: %v", input, err)
				}
			}
		})
	}
}
