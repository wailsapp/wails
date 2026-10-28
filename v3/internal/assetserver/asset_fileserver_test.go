package assetserver

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"testing/fstest"
)

func TestNonSeekingAssetPrefixWriteFailure(t *testing.T) {
	writeError := errors.New("prefix write failed")
	for name, writeErr := range map[string]error{"short write": nil, "write error": writeError} {
		t.Run(name, func(t *testing.T) {
			bytesRead := 0
			server := &assetFileServer{fs: nonSeekingTestFS{
				FS: fstest.MapFS{"asset.txt": {Data: make([]byte, 1024)}}, bytesRead: &bytesRead,
			}}
			writer := &prefixFailureWriter{ResponseWriter: httptest.NewRecorder(), err: writeErr}
			err := server.serveFSFile(writer, httptest.NewRequest(http.MethodGet, "/asset.txt", nil), "asset.txt")
			want := writeErr
			if want == nil {
				want = io.ErrShortWrite
			}
			if !errors.Is(err, want) {
				t.Errorf("error = %v, want %v", err, want)
			}
			if bytesRead != 512 || writer.writes != 1 {
				t.Errorf("continued after prefix failure: read %d bytes, wrote %d times", bytesRead, writer.writes)
			}
		})
	}
}

type prefixFailureWriter struct {
	http.ResponseWriter
	err    error
	writes int
}

func (w *prefixFailureWriter) Write(p []byte) (int, error) {
	w.writes++
	return len(p) / 2, w.err
}

func TestNonSeekingAssetsConcurrent(t *testing.T) {
	for i := 0; i < 16; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()
			data := bytes.Repeat([]byte{byte('a' + i)}, 64*1024)
			bytesRead := 0
			server := &assetFileServer{fs: nonSeekingTestFS{
				FS:        fstest.MapFS{"asset.txt": {Data: data}},
				bytesRead: &bytesRead,
			}}
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/asset.txt", nil))
			if recorder.Code != http.StatusOK || !bytes.Equal(recorder.Body.Bytes(), data) {
				t.Errorf("concurrent response status=%d, body length=%d", recorder.Code, recorder.Body.Len())
			}
			if bytesRead != len(data) {
				t.Errorf("read %d bytes, want %d", bytesRead, len(data))
			}
		})
	}
}

func TestEmptyNonSeekingAssetCommitsSuccess(t *testing.T) {
	for contentType, wantType := range map[string]string{"": "text/plain; charset=utf-8", "text/custom": "text/custom"} {
		t.Run("type="+contentType, func(t *testing.T) {
			bytesRead := 0
			files := &assetFileServer{fs: nonSeekingTestFS{
				FS:        fstest.MapFS{"slow": {}},
				bytesRead: &bytesRead,
			}}
			server, err := NewAssetServer(&Options{Logger: slog.Default(), Handler: files})
			if err != nil {
				t.Fatal(err)
			}
			response := &contextWebViewResponse{}
			if contentType != "" {
				response.Header().Set(HeaderContentType, contentType)
			}
			request := &contextWebViewRequest{ctx: context.Background(), rw: response}
			server.processWebViewRequest(request)
			if response.Code() != http.StatusOK || response.body.Len() != 0 {
				t.Errorf("response status=%d, body length=%d; want 200 with no body", response.Code(), response.body.Len())
			}
			if response.Header().Get(HeaderContentLength) != "0" || response.Header().Get(HeaderContentType) != wantType {
				t.Errorf("unexpected headers: %v", response.Header())
			}
			if response.finishCount.Load() != 1 || request.closeCount.Load() != 1 {
				t.Error("request was not finished and closed exactly once")
			}
		})
	}
}

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
