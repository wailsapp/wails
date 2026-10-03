package idl

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// rtFunc adapts a function to an http.RoundTripper so Download's hardcoded
// NuGet URL can be intercepted and served a canned response.
type rtFunc func(*http.Request) (*http.Response, error)

func (f rtFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// clientServing returns an *http.Client whose transport answers every request
// with the given status and body, regardless of URL.
func clientServing(status int, body []byte) *http.Client {
	return &http.Client{Transport: rtFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: status,
			Body:       io.NopCloser(bytes.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})}
}

func TestNewFetcher(t *testing.T) {
	s := NewStore(t.TempDir())
	f := NewFetcher(s)
	if f.Store != s {
		t.Error("NewFetcher should retain the provided store")
	}
	if f.HTTPClient == nil {
		t.Fatal("NewFetcher should construct an HTTP client")
	}
	if f.HTTPClient.Timeout != FetchTimeout {
		t.Errorf("HTTPClient.Timeout = %v, want %v", f.HTTPClient.Timeout, FetchTimeout)
	}
}

func TestDownloadNetworkSuccess(t *testing.T) {
	pkg := makePackage(t, []byte(fakeIDL), "")
	dir := t.TempDir()
	s := NewStore(dir)
	f := &Fetcher{HTTPClient: clientServing(http.StatusOK, pkg), Store: s}

	got, err := f.Download("1.2.3.4")
	if err != nil {
		t.Fatalf("Download: %v", err)
	}
	if string(got) != fakeIDL {
		t.Errorf("Download returned %q, want %q", got, fakeIDL)
	}
	// The fetched IDL must be written back to the cache.
	if !s.Has("1.2.3.4") {
		t.Error("Download should populate the cache")
	}
	cached, err := os.ReadFile(s.CachePath("1.2.3.4"))
	if err != nil || string(cached) != fakeIDL {
		t.Errorf("cache file = (%q, %v), want %q", cached, err, fakeIDL)
	}
}

func TestDownloadNoStore(t *testing.T) {
	pkg := makePackage(t, []byte(fakeIDL), "")
	f := &Fetcher{HTTPClient: clientServing(http.StatusOK, pkg), Store: nil}

	got, err := f.Download("1.2.3.4")
	if err != nil {
		t.Fatalf("Download with nil store: %v", err)
	}
	if string(got) != fakeIDL {
		t.Errorf("Download returned %q, want %q", got, fakeIDL)
	}
}

func TestDownloadHTTPError(t *testing.T) {
	wantErr := errors.New("dial tcp: boom")
	f := &Fetcher{HTTPClient: &http.Client{Transport: rtFunc(func(r *http.Request) (*http.Response, error) {
		return nil, wantErr
	})}}
	if _, err := f.Download("1.2.3.4"); err == nil || !strings.Contains(err.Error(), "fetch") {
		t.Errorf("Download transport error = %v, want a fetch error", err)
	}
}

func TestDownloadNon200(t *testing.T) {
	f := &Fetcher{HTTPClient: clientServing(http.StatusNotFound, []byte("nope"))}
	if _, err := f.Download("1.2.3.4"); err == nil || !strings.Contains(err.Error(), "HTTP 404") {
		t.Errorf("Download non-200 error = %v, want HTTP 404", err)
	}
}

// errReader fails partway through, exercising Download's io.ReadAll error arm.
type errReader struct{}

func (errReader) Read([]byte) (int, error) { return 0, errors.New("read failed") }

func TestDownloadBodyReadError(t *testing.T) {
	f := &Fetcher{HTTPClient: &http.Client{Transport: rtFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(errReader{}),
			Header:     make(http.Header),
		}, nil
	})}}
	if _, err := f.Download("1.2.3.4"); err == nil || !strings.Contains(err.Error(), "read body") {
		t.Errorf("Download body-read error = %v, want read body error", err)
	}
}

func TestDownloadExtractError(t *testing.T) {
	// Not a zip → extractIDL fails.
	f := &Fetcher{HTTPClient: clientServing(http.StatusOK, []byte("definitely not a zip"))}
	if _, err := f.Download("1.2.3.4"); err == nil || !strings.Contains(err.Error(), "extract IDL") {
		t.Errorf("Download extract error = %v, want extract IDL error", err)
	}
}

func TestDownloadMkdirError(t *testing.T) {
	pkg := makePackage(t, []byte(fakeIDL), "")
	// Place the cache dir beneath a regular file so MkdirAll cannot create it.
	blocker := filepath.Join(t.TempDir(), "afile")
	if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	s := NewStore(filepath.Join(blocker, "cache"))
	f := &Fetcher{HTTPClient: clientServing(http.StatusOK, pkg), Store: s}
	if _, err := f.Download("1.2.3.4"); err == nil || !strings.Contains(err.Error(), "create cache dir") {
		t.Errorf("Download mkdir error = %v, want create cache dir error", err)
	}
}

func TestDownloadWriteError(t *testing.T) {
	pkg := makePackage(t, []byte(fakeIDL), "")
	dir := t.TempDir()
	s := NewStore(dir)
	// Pre-create a directory exactly where the cache file should be written so
	// WriteFile fails. Has() ignores it (it's a directory), so the network path runs.
	if err := os.Mkdir(s.CachePath("1.2.3.4"), 0o755); err != nil {
		t.Fatal(err)
	}
	f := &Fetcher{HTTPClient: clientServing(http.StatusOK, pkg), Store: s}
	if _, err := f.Download("1.2.3.4"); err == nil || !strings.Contains(err.Error(), "write cache") {
		t.Errorf("Download write error = %v, want write cache error", err)
	}
}

func TestListMissingDir(t *testing.T) {
	s := NewStore(filepath.Join(t.TempDir(), "does-not-exist"))
	got, err := s.List()
	if err != nil {
		t.Fatalf("List on missing dir: %v", err)
	}
	if got != nil {
		t.Errorf("List on missing dir = %v, want nil", got)
	}
}

func TestListReadDirError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("unix permission bits don't block directory reads on windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("root bypasses directory permission checks")
	}
	dir := filepath.Join(t.TempDir(), "locked")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	// Remove read permission so ReadDir fails with a non-IsNotExist error.
	if err := os.Chmod(dir, 0o000); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(dir, 0o755) //nolint:errcheck // restore for cleanup

	s := NewStore(dir)
	if _, err := s.List(); err == nil {
		t.Error("List on an unreadable directory should return an error")
	}
}

// deflateIDLZip builds a single-entry zip (WebView2.idl, deflate-compressed)
// using content compressible enough to produce a real deflate stream that can
// be corrupted. It returns the bytes plus the byte offset of the compression
// method field in the local file header and the start of the compressed data.
func deflateIDLZip(t *testing.T) (raw []byte, methodOff, dataOff int) {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, err := zw.CreateHeader(&zip.FileHeader{Name: "WebView2.idl", Method: zip.Deflate})
	if err != nil {
		t.Fatal(err)
	}
	// Repetitive content guarantees a non-trivial deflate stream.
	if _, err := w.Write(bytes.Repeat([]byte("interface IFake { };\n"), 64)); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	raw = buf.Bytes()

	// Locate the local file header: signature PK\x03\x04. Method is at +8,
	// name length at +26, extra length at +28, data begins after both.
	i := bytes.Index(raw, []byte("PK\x03\x04"))
	if i < 0 {
		t.Fatal("local file header signature not found")
	}
	nameLen := int(raw[i+26]) | int(raw[i+27])<<8
	extraLen := int(raw[i+28]) | int(raw[i+29])<<8
	return raw, i + 8, i + 30 + nameLen + extraLen
}

// TestExtractIDLOpenError exercises the file.Open() error arm by patching the
// stored compression method to an unregistered value the reader rejects.
func TestExtractIDLOpenError(t *testing.T) {
	raw, methodOff, _ := deflateIDLZip(t)
	// Patch the local-header method field. The reader opens entries from the
	// central directory but reads the method from there too; patch both copies.
	raw[methodOff] = 99
	raw[methodOff+1] = 0
	if c := bytes.Index(raw, []byte("PK\x01\x02")); c >= 0 {
		raw[c+10] = 99
		raw[c+11] = 0
	}
	if _, err := extractIDL(raw); err == nil ||
		!strings.Contains(err.Error(), "open") {
		t.Errorf("extractIDL with bad compression method = %v, want an open error", err)
	}
}

// TestExtractIDLReadError exercises the io.ReadAll error arm by corrupting the
// deflate stream so decompression fails after the entry opens.
func TestExtractIDLReadError(t *testing.T) {
	raw, _, dataOff := deflateIDLZip(t)
	// Corrupt several bytes inside the compressed payload.
	for k := 0; k < 8 && dataOff+k < len(raw); k++ {
		raw[dataOff+k] ^= 0xFF
	}
	if _, err := extractIDL(raw); err == nil ||
		!strings.Contains(err.Error(), "read") {
		t.Errorf("extractIDL with corrupt deflate = %v, want a read error", err)
	}
}

func TestExtractIDLBadZip(t *testing.T) {
	if _, err := extractIDL([]byte("PK\x03\x04 not really a zip")); err == nil ||
		!strings.Contains(err.Error(), "parse zip") {
		t.Errorf("extractIDL on bad zip = %v, want parse zip error", err)
	}
}

func TestListSkipsEmptyVersion(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(dir)
	// "WebView2..idl" trims to an empty version and must be skipped.
	if err := os.WriteFile(filepath.Join(dir, "WebView2..idl"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(s.CachePath("1.0.0.0"), []byte(fakeIDL), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 1 || got[0] != "1.0.0.0" {
		t.Errorf("List() = %v, want [1.0.0.0] (empty version skipped)", got)
	}
}
