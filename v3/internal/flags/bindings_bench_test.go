package flags

import (
	"strings"
	"testing"
)

func BenchmarkBuildFlags(b *testing.B) {
	for _, tc := range []struct {
		name, input string
		count       int
	}{
		{"single", "-race", 1},
		{"quoted", `-tags 'desktop debug' -ldflags "-X main.version=3.1.0 -s -w"`, 4},
		{"many", strings.Repeat("-tags=desktop ", 32), 32},
	} {
		b.Run(tc.name, func(b *testing.B) {
			options := GenerateBindingsOptions{BuildFlagsString: tc.input}
			b.ReportAllocs()
			for range b.N {
				flags, err := options.BuildFlags()
				if err != nil || len(flags) != tc.count {
					b.Fatalf("flags=%q, error=%v; want %d flags", flags, err, tc.count)
				}
			}
		})
	}
}
