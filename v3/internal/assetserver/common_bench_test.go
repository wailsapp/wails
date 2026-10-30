package assetserver

import (
	"net/http/httptest"
	"strconv"
	"testing"
)

func BenchmarkServeFile(b *testing.B) {
	for _, size := range []int{0, 1024, 65536} {
		b.Run(strconv.Itoa(size), func(b *testing.B) {
			blob := make([]byte, size)
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				response := httptest.NewRecorder()
				if err := ServeFile(response, "asset.js", blob); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
