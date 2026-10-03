// Package idl fetches, caches, and inspects WebView2 SDK IDL files.
//
// The SDK is distributed as a NuGet package
// (https://www.nuget.org/api/v2/package/Microsoft.Web.WebView2/<version>),
// which is a zip with the WebView2.idl file inside. The fetcher downloads
// the package, extracts the IDL, and writes it to a local cache directory
// so subsequent runs are offline.
package idl

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// NuGetPackageURL is the template for the SDK NuGet package download.
	NuGetPackageURL = "https://www.nuget.org/api/v2/package/Microsoft.Web.WebView2/%s"

	// IDLFilenameInPackage is the path inside the .nupkg where the IDL lives.
	IDLFilenameInPackage = "WebView2.idl"

	// FetchTimeout bounds a single network fetch.
	FetchTimeout = 60 * time.Second
)

// Store is a local cache of IDL files keyed by SDK version.
type Store struct {
	Dir string
}

// NewStore returns a Store backed by dir. The directory is created on demand.
func NewStore(dir string) *Store { return &Store{Dir: dir} }

// CachePath returns the absolute path where a given version's IDL would
// be stored. It does not check existence.
func (s *Store) CachePath(version string) string {
	return filepath.Join(s.Dir, "WebView2."+version+".idl")
}

// Has reports whether the cache contains an IDL for the given version.
func (s *Store) Has(version string) bool {
	info, err := os.Stat(s.CachePath(version))
	return err == nil && !info.IsDir() && info.Size() > 0
}

// List returns the versions present in the cache, parsed from filenames
// of the form `WebView2.<version>.idl`. Versions are returned unsorted;
// the caller may sort them with idlversion.Compare if needed.
func (s *Store) List() ([]string, error) {
	entries, err := os.ReadDir(s.Dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var versions []string
	for _, e := range entries {
		name := e.Name()
		if !strings.HasPrefix(name, "WebView2.") || !strings.HasSuffix(name, ".idl") {
			continue
		}
		v := strings.TrimSuffix(strings.TrimPrefix(name, "WebView2."), ".idl")
		if v != "" {
			versions = append(versions, v)
		}
	}
	return versions, nil
}

// Read returns the IDL bytes for version, reading from the cache only.
// Returns an error wrapping os.ErrNotExist if the version isn't cached.
func (s *Store) Read(version string) ([]byte, error) {
	return os.ReadFile(s.CachePath(version))
}

// Fetcher downloads SDK IDL files from NuGet.
type Fetcher struct {
	HTTPClient *http.Client
	Store      *Store
}

// NewFetcher constructs a Fetcher with a default HTTP client.
func NewFetcher(store *Store) *Fetcher {
	return &Fetcher{
		HTTPClient: &http.Client{Timeout: FetchTimeout},
		Store:      store,
	}
}

// Download fetches the SDK IDL for the given version. If the IDL is
// already cached the cached bytes are returned. The fetched bytes are
// always written back to the cache.
func (f *Fetcher) Download(version string) ([]byte, error) {
	if f.Store != nil && f.Store.Has(version) {
		return f.Store.Read(version)
	}

	url := fmt.Sprintf(NuGetPackageURL, version)
	resp, err := f.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch %s: HTTP %d", url, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	idl, err := extractIDL(body)
	if err != nil {
		return nil, fmt.Errorf("extract IDL from package %s: %w", version, err)
	}

	if f.Store != nil {
		if err := os.MkdirAll(f.Store.Dir, 0o755); err != nil {
			return nil, fmt.Errorf("create cache dir: %w", err)
		}
		if err := os.WriteFile(f.Store.CachePath(version), idl, 0o644); err != nil {
			return nil, fmt.Errorf("write cache: %w", err)
		}
	}
	return idl, nil
}

// extractIDL pulls WebView2.idl out of a .nupkg byte slice.
func extractIDL(pkg []byte) ([]byte, error) {
	zr, err := zip.NewReader(bytes.NewReader(pkg), int64(len(pkg)))
	if err != nil {
		return nil, fmt.Errorf("parse zip: %w", err)
	}
	for _, file := range zr.File {
		// The IDL lives at the root of the package as WebView2.idl. Some
		// older or arch-specific .nupkgs nest it; accept any path whose
		// basename matches.
		if filepath.Base(file.Name) != IDLFilenameInPackage {
			continue
		}
		r, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("open %s: %w", file.Name, err)
		}
		idl, err := io.ReadAll(r)
		r.Close()
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", file.Name, err)
		}
		return idl, nil
	}
	return nil, fmt.Errorf("WebView2.idl not found in package")
}
