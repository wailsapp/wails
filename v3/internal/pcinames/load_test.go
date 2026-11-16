package pcinames

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestDiscoveryOrder(t *testing.T) {
	root := t.TempDir()
	c := config{root: root, path: filepath.Join(root, "override"), cache: filepath.Join(root, "cached")}
	want := []string{c.path, c.cache, filepath.Join(root, "usr/share/hwdata/pci.ids"), filepath.Join(root, "usr/share/misc/pci.ids"), filepath.Join(root, "usr/share/hwdata/pci.ids.gz"), filepath.Join(root, "usr/share/misc/pci.ids.gz")}
	if got := c.paths(); !reflect.DeepEqual(got, want) {
		t.Fatalf("paths %+v; want %+v", got, want)
	}
	for _, path := range want {
		data := []byte("1234  " + filepath.Base(path) + "\n\t5678  Card\n")
		if strings.HasSuffix(path, ".gz") {
			data = compressed(t, data)
		}
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	for _, path := range want {
		db, err := load(c, nil, "")
		if err != nil {
			t.Fatal(err)
		}
		vendor, _ := db.Lookup("1234", "5678")
		if vendor != filepath.Base(path) {
			t.Fatalf("selected %q before %s", vendor, path)
		}
		if err := os.Remove(path); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := load(c, nil, ""); err == nil {
		t.Fatal("missing databases accepted")
	}
}

func TestEnvironment(t *testing.T) {
	for k, v := range map[string]string{"PCIDB_PATH": "/custom", "PCIDB_CHROOT": "/host", "PCIDB_CACHE_PATH": "/cached", "PCIDB_CACHE_ONLY": "1", "PCIDB_ENABLE_NETWORK_FETCH": "1"} {
		t.Setenv(k, v)
	}
	c := environment()
	if c.root != "/host" || c.path != "/custom" || c.cache != "/cached" || !c.cacheOnly || !c.network {
		t.Fatalf("unexpected config %+v", c)
	}
	if !reflect.DeepEqual(c.paths(), []string{"/custom", "/cached"}) {
		t.Fatal(c.paths())
	}
	t.Setenv("PCIDB_ENABLE_NETWORK_FETCH", "0")
	if environment().network {
		t.Fatal("network enabled for zero")
	}
	t.Setenv("PCIDB_CACHE_PATH", "")
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	if environment().cache != filepath.Join(home, ".cache/pci.ids") {
		t.Fatal("default cache changed")
	}
}

func TestFetchAndCache(t *testing.T) {
	requests := 0
	data := compressed(t, []byte(fixture))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.UserAgent() != "Wails doctor PCI names" {
			t.Error("missing user agent")
		}
		_, _ = w.Write(data)
	}))
	defer server.Close()
	c := config{cache: filepath.Join(t.TempDir(), "nested/pci.ids"), cacheOnly: true, network: true}
	db, err := load(c, server.Client(), server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if vendor, _ := db.Lookup("1234", "5678"); vendor != "Example  vendor" {
		t.Fatal(vendor)
	}
	cached, err := os.ReadFile(c.cache)
	if err != nil || !bytes.Equal(cached, data) {
		t.Fatalf("cache %d bytes, %v", len(cached), err)
	}
	if _, err := load(c, nil, ""); err != nil {
		t.Fatal(err)
	}
	if requests != 1 {
		t.Fatalf("unexpected requests %d", requests)
	}
	// Cache failures must not discard a valid response for the current report.
	blocker := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(blocker, []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := fetch(server.Client(), server.URL, filepath.Join(blocker, "pci.ids")); err != nil {
		t.Fatal(err)
	}
}

func TestFetchFailures(t *testing.T) {
	for _, tt := range []struct {
		name   string
		status int
		data   []byte
	}{
		{"status", 503, nil}, {"corrupt gzip", 200, []byte{0x1f, 0x8b}}, {"oversized", 200, bytes.Repeat([]byte{'x'}, maxDownloadBytes+1)},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(tt.status); _, _ = w.Write(tt.data) }))
			defer server.Close()
			cache := filepath.Join(t.TempDir(), "pci.ids")
			if _, err := fetch(server.Client(), server.URL, cache); err == nil {
				t.Fatal("bad response accepted")
			}
			if _, err := os.Stat(cache); !os.IsNotExist(err) {
				t.Fatalf("failed response cached: %v", err)
			}
		})
	}
	if _, err := fetch(http.DefaultClient, ":bad-url", ""); err == nil {
		t.Fatal("invalid endpoint accepted")
	}
}

func TestStorePreservesCacheOnFailure(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "existing-dir")
	if err := os.Mkdir(target, 0700); err != nil {
		t.Fatal(err)
	}
	if err := store(target, []byte(fixture)); err == nil {
		t.Fatal("replaced directory")
	}
	entries, err := os.ReadDir(root)
	if err != nil || len(entries) != 1 {
		t.Fatalf("temporary files leaked: %v %v", entries, err)
	}
}

type transportFunc func(*http.Request) (*http.Response, error)

func (f transportFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestLoadHonoursConfiguredClient(t *testing.T) {
	prior := http.DefaultClient
	t.Cleanup(func() { http.DefaultClient = prior })
	for key, value := range map[string]string{"PCIDB_CACHE_ONLY": "true", "PCIDB_PATH": "", "PCIDB_ENABLE_NETWORK_FETCH": "true"} {
		t.Setenv(key, value)
	}
	for _, timeout := range []time.Duration{100 * time.Millisecond, time.Minute, 0} {
		t.Setenv("PCIDB_CACHE_PATH", filepath.Join(t.TempDir(), "pci.ids"))
		calls := 0
		client := &http.Client{Timeout: timeout, Transport: transportFunc(func(r *http.Request) (*http.Response, error) {
			calls++
			deadline, ok := r.Context().Deadline()
			upper := 15 * time.Second
			if timeout > 0 && timeout < upper {
				upper = timeout
			}
			remaining := time.Until(deadline)
			if !ok || remaining <= 0 || remaining > upper {
				t.Errorf("unexpected request timeout %v, present %v", remaining, ok)
			}
			return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(fixture)), Request: r}, nil
		})}
		http.DefaultClient = client
		if _, err := Load(); err != nil {
			t.Fatal(err)
		}
		if calls != 1 || client.Timeout != timeout {
			t.Fatalf("calls=%d; shared client timeout changed to %v", calls, client.Timeout)
		}
	}
}
