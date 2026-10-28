package assetserver

import (
	"bytes"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"testing/fstest"
)

type nonSeekingTestFS struct {
	fs.FS
	bytesRead *int
	maxRead   int
}

func (f nonSeekingTestFS) Open(name string) (fs.File, error) {
	file, err := f.FS.Open(name)
	if err != nil {
		return nil, err
	}
	return nonSeekingTestFile{File: file, bytesRead: f.bytesRead, maxRead: f.maxRead}, nil
}

type nonSeekingTestFile struct {
	fs.File
	bytesRead *int
	maxRead   int
}

func (f nonSeekingTestFile) Read(p []byte) (int, error) {
	if f.maxRead > 0 && len(p) > f.maxRead {
		p = p[:f.maxRead]
	}
	n, err := f.File.Read(p)
	*f.bytesRead += n
	return n, err
}

func TestAssetFileServerShortReadMIME(t *testing.T) {
	png := []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}
	longPNG := append(bytes.Clone(png), make([]byte, 600)...)
	tests := []struct {
		name     string
		method   string
		data     []byte
		wantType string
		wantRead int
	}{
		{"short GET", http.MethodGet, png, "image/png", len(png)},
		{"long GET", http.MethodGet, longPNG, "image/png", len(longPNG)},
		{"short HEAD", http.MethodHead, png, "image/png", len(png)},
		{"long HEAD", http.MethodHead, longPNG, "image/png", 512},
		{"empty GET", http.MethodGet, nil, "text/plain; charset=utf-8", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytesRead := 0
			server := &assetFileServer{fs: nonSeekingTestFS{
				FS:        fstest.MapFS{"asset.unknown": {Data: tt.data}},
				bytesRead: &bytesRead,
				maxRead:   1,
			}}
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, httptest.NewRequest(tt.method, "/asset.unknown", nil))
			wantBody := tt.data
			if tt.method == http.MethodHead {
				wantBody = nil
			}
			if recorder.Code != http.StatusOK || !bytes.Equal(recorder.Body.Bytes(), wantBody) {
				t.Errorf("status=%d, body length=%d; want 200, body length=%d", recorder.Code, recorder.Body.Len(), len(wantBody))
			}
			if recorder.Header().Get(HeaderContentType) != tt.wantType || recorder.Header().Get(HeaderContentLength) != strconv.Itoa(len(tt.data)) {
				t.Errorf("unexpected headers: %v", recorder.Header())
			}
			if bytesRead != tt.wantRead {
				t.Errorf("read %d bytes, want %d", bytesRead, tt.wantRead)
			}
		})
	}
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
