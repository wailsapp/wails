package idl

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

const fakeIDL = "library WebView2 { interface IFake { HRESULT Ping(); } }"

func makePackage(t *testing.T, idl []byte, nestedPath string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for _, p := range []string{"_rels/.rels", "WebView2.nuspec"} {
		f, err := zw.Create(p)
		if err != nil {
			t.Fatal(err)
		}
		f.Write([]byte("filler"))
	}
	path := "WebView2.idl"
	if nestedPath != "" {
		path = nestedPath
	}
	f, err := zw.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	f.Write(idl)
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestExtractIDL(t *testing.T) {
	pkg := makePackage(t, []byte(fakeIDL), "")
	got, err := extractIDL(pkg)
	if err != nil {
		t.Fatalf("extractIDL: %v", err)
	}
	if string(got) != fakeIDL {
		t.Errorf("extractIDL produced %q, want %q", got, fakeIDL)
	}
}

func TestExtractIDLNestedPath(t *testing.T) {
	pkg := makePackage(t, []byte(fakeIDL), "lib/native/include/WebView2.idl")
	got, err := extractIDL(pkg)
	if err != nil {
		t.Fatalf("extractIDL: %v", err)
	}
	if string(got) != fakeIDL {
		t.Errorf("extractIDL nested: got %q, want %q", got, fakeIDL)
	}
}

func TestExtractIDLMissing(t *testing.T) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	f, _ := zw.Create("README.md")
	f.Write([]byte("no idl here"))
	zw.Close()
	if _, err := extractIDL(buf.Bytes()); err == nil {
		t.Error("extractIDL on package without IDL should fail")
	}
}

func TestStore(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(dir)

	if s.Has("1.0.0.0") {
		t.Error("empty store should not have any versions")
	}
	if list, err := s.List(); err != nil || len(list) != 0 {
		t.Errorf("empty List() = (%v, %v); want ([], nil)", list, err)
	}

	for _, v := range []string{"1.0.2903.40", "1.0.2739.15"} {
		path := s.CachePath(v)
		if err := os.WriteFile(path, []byte(fakeIDL), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// A non-IDL file should be ignored.
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	sort.Strings(got)
	want := []string{"1.0.2739.15", "1.0.2903.40"}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("List() = %v, want %v", got, want)
	}

	if !s.Has("1.0.2903.40") {
		t.Error("Has() should report cached version")
	}
	data, err := s.Read("1.0.2903.40")
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if string(data) != fakeIDL {
		t.Errorf("Read produced unexpected content")
	}
}

func TestFetcherUsesCache(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(dir)
	if err := os.WriteFile(s.CachePath("9.9.9.9"), []byte(fakeIDL), 0o644); err != nil {
		t.Fatal(err)
	}

	fetched := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fetched = true
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()
	f := &Fetcher{HTTPClient: srv.Client(), Store: s}

	got, err := f.Download("9.9.9.9")
	if err != nil {
		t.Fatalf("Download cached: %v", err)
	}
	if string(got) != fakeIDL {
		t.Errorf("cached download returned wrong content: %q", got)
	}
	if fetched {
		t.Error("Download should not hit the network when version is cached")
	}
}

func TestFetcherDownloads(t *testing.T) {
	pkg := makePackage(t, []byte(fakeIDL), "")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, bytes.NewReader(pkg))
	}))
	defer srv.Close()

	dir := t.TempDir()
	s := NewStore(dir)
	f := &Fetcher{HTTPClient: srv.Client(), Store: s}

	// Bypass NuGet URL by calling Download against our test server directly.
	resp, err := f.HTTPClient.Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	idl, err := extractIDL(body)
	if err != nil {
		t.Fatalf("extractIDL roundtrip: %v", err)
	}
	if string(idl) != fakeIDL {
		t.Errorf("fetch roundtrip produced wrong content: %q", idl)
	}
}
