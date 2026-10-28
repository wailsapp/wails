package assetserver

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

type nonSeekingTestFS struct {
	fs.FS
	bytesRead *int
}

func (f nonSeekingTestFS) Open(name string) (fs.File, error) {
	file, err := f.FS.Open(name)
	if err != nil {
		return nil, err
	}
	return nonSeekingTestFile{File: file, bytesRead: f.bytesRead}, nil
}

type nonSeekingTestFile struct {
	fs.File
	bytesRead *int
}

func (f nonSeekingTestFile) Read(p []byte) (int, error) {
	n, err := f.File.Read(p)
	*f.bytesRead += n
	return n, err
}

func TestAssetFileServerNonSeekingHEAD(t *testing.T) {
	data := strings.Repeat("body", 1024)
	tests := []struct {
		name        string
		method      string
		contentType string
		wantType    string
		wantBody    string
		wantRead    int
	}{
		{"HEAD/sniffed", http.MethodHead, "", "text/plain; charset=utf-8", "", 512},
		{"HEAD/explicit", http.MethodHead, "text/custom", "text/custom", "", 0},
		{"GET/sniffed", http.MethodGet, "", "text/plain; charset=utf-8", data, len(data)},
		{"GET/explicit", http.MethodGet, "text/custom", "text/custom", data, len(data)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytesRead := 0
			server := &assetFileServer{fs: nonSeekingTestFS{
				FS:        fstest.MapFS{"asset.txt": {Data: []byte(data)}},
				bytesRead: &bytesRead,
			}}
			recorder := httptest.NewRecorder()
			if tt.contentType != "" {
				recorder.Header().Set(HeaderContentType, tt.contentType)
			}
			server.ServeHTTP(recorder, httptest.NewRequest(tt.method, "/asset.txt", nil))
			if recorder.Code != http.StatusOK || recorder.Body.String() != tt.wantBody {
				t.Errorf("status=%d, body length=%d; want 200, body length=%d", recorder.Code, recorder.Body.Len(), len(tt.wantBody))
			}
			if recorder.Header().Get(HeaderContentLength) != "4096" || recorder.Header().Get(HeaderContentType) != tt.wantType {
				t.Errorf("unexpected headers: %v", recorder.Header())
			}
			if bytesRead != tt.wantRead {
				t.Errorf("read %d bytes, want %d", bytesRead, tt.wantRead)
			}
		})
	}
}
