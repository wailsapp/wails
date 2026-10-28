package assetserver

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"testing/fstest"
)

func BenchmarkServeNonSeekingAsset(b *testing.B) {
	for _, size := range []int{0, 1024, 65536} {
		b.Run(strconv.Itoa(size), func(b *testing.B) {
			bytesRead := 0
			server := &assetFileServer{fs: nonSeekingTestFS{
				FS:        fstest.MapFS{"asset.txt": {Data: make([]byte, size)}},
				bytesRead: &bytesRead,
			}}
			request := httptest.NewRequest(http.MethodGet, "/asset.txt", nil)
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				bytesRead = 0
				response := httptest.NewRecorder()
				response.Body = nil
				server.ServeHTTP(response, request)
			}
		})
	}
}
