package semver

import (
	"slices"
	"strconv"
	"strings"
	"testing"
)

func BenchmarkComparePrerelease(b *testing.B) {
	for _, tc := range []struct {
		name, a, b string
		want       int
	}{
		{"nightly", "nightly.123", "nightly.124", -1},
		{"long_equal_prefix", strings.Repeat("preview.", 16) + "123", strings.Repeat("preview.", 16) + "124", -1},
		{"early_difference", "alpha." + strings.Repeat("1.", 32), "beta." + strings.Repeat("1.", 32), -1},
		{"large_numeric", "nightly." + strings.Repeat("9", 128), "nightly.1" + strings.Repeat("0", 128), -1},
	} {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if got := comparePrerelease(tc.a, tc.b); got != tc.want {
					b.Fatalf("comparison=%d, want %d", got, tc.want)
				}
			}
		})
	}
}

func BenchmarkSortPrereleases(b *testing.B) {
	versions := make([]Version, 256)
	for i := range versions {
		versions[i] = Version{Major: 3, Prerelease: "nightly." + strconv.Itoa((i*73)%256)}
	}
	work := make([]Version, len(versions))
	compare := func(a, b Version) int { return a.compare(&b) }
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(work, versions)
		slices.SortFunc(work, compare)
	}
	b.StopTimer()
	for i, v := range work {
		if v.Prerelease != "nightly."+strconv.Itoa(i) {
			b.Fatalf("version %d=%s", i, v.Prerelease)
		}
	}
}
