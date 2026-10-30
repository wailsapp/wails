package assetserver

import (
	"strconv"
	"testing"
	"testing/fstest"
)

func BenchmarkFindPathToFile(b *testing.B) {
	for _, count := range []int{1, 128, 512} {
		b.Run(strconv.Itoa(count), func(b *testing.B) {
			assets := make(fstest.MapFS, count)
			for i := 0; i < count; i++ {
				assets["assets/"+strconv.Itoa(i)+"/index.html"] = &fstest.MapFile{}
			}
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				root, err := findPathToFile(assets, indexHTML)
				if err != nil || root != "assets/0/" {
					b.Fatalf("root = %q, error = %v", root, err)
				}
			}
		})
	}
}
