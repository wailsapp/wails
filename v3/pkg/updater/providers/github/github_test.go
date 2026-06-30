package github_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/updater"
	"github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)

const sampleReleaseJSON = `{
  "tag_name": "v2.0.0",
  "name": "v2",
  "body": "Bug fixes and new features.",
  "prerelease": false,
  "draft": false,
  "published_at": "2026-03-01T10:00:00Z",
  "html_url": "https://github.com/o/r/releases/tag/v2.0.0",
  "assets": [
    {"id": 1, "name": "app-darwin-arm64.dmg", "content_type": "application/octet-stream", "size": 12345, "browser_download_url": "%s/dl/darwin-arm64"},
    {"id": 2, "name": "app-linux-amd64.tar.gz", "content_type": "application/gzip", "size": 22222, "browser_download_url": "%s/dl/linux-amd64"},
    {"id": 3, "name": "app-windows-amd64.zip", "content_type": "application/zip", "size": 33333, "browser_download_url": "%s/dl/windows-amd64"}
  ]
}`

func TestCheck_PicksMatchingAsset(t *testing.T) {
	var sawAuth, sawAccept string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if want := "/repos/o/r/releases/latest"; r.URL.Path != want {
			t.Errorf("path: want %q, got %q", want, r.URL.Path)
		}
		sawAuth = r.Header.Get("Authorization")
		sawAccept = r.Header.Get("Accept")
		fmt.Fprintf(w, sampleReleaseJSON, "https://example.invalid", "https://example.invalid", "https://example.invalid")
	}))
	defer srv.Close()

	p, err := github.New(github.Config{Repository: "o/r", Token: "ghp_xxx", BaseURL: srv.URL})
	if err != nil {
		t.Fatal(err)
	}
	rel, err := p.Check(context.Background(), updater.CheckRequest{
		CurrentVersion: "1.0.0",
		Platform:       "darwin",
		Arch:           "arm64",
	})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if rel == nil {
		t.Fatal("expected release, got nil")
	}
	if rel.Version != "2.0.0" {
		t.Errorf("version: want 2.0.0, got %q", rel.Version)
	}
	if rel.Artifact.Filename != "app-darwin-arm64.dmg" {
		t.Errorf("asset: %q", rel.Artifact.Filename)
	}
	if rel.Artifact.Size != 12345 {
		t.Errorf("size: %d", rel.Artifact.Size)
	}
	if sawAuth != "Bearer ghp_xxx" {
		t.Errorf("auth: %q", sawAuth)
	}
	if !strings.Contains(sawAccept, "application/vnd.github+json") {
		t.Errorf("accept: %q", sawAccept)
	}
}

func TestCheck_UpToDate_ReturnsNil(t *testing.T) {
	var host string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, sampleReleaseJSON, host, host, host)
	}))
	defer srv.Close()
	host = srv.URL

	p, _ := github.New(github.Config{Repository: "o/r", BaseURL: srv.URL})
	rel, err := p.Check(context.Background(), updater.CheckRequest{CurrentVersion: "2.0.0", Platform: "darwin", Arch: "arm64"})
	if err != nil {
		t.Fatal(err)
	}
	if rel != nil {
		t.Fatalf("expected nil release (up-to-date), got %+v", rel)
	}
}

func TestCheck_404_TreatedAsUpToDate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()
	p, _ := github.New(github.Config{Repository: "o/r", BaseURL: srv.URL})
	rel, err := p.Check(context.Background(), updater.CheckRequest{CurrentVersion: "1.0.0"})
	if err != nil {
		t.Fatal(err)
	}
	if rel != nil {
		t.Errorf("404 should map to nil release, got %+v", rel)
	}
}

func TestCheck_NoMatchingAsset_Errors(t *testing.T) {
	var host string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, sampleReleaseJSON, host, host, host)
	}))
	defer srv.Close()
	host = srv.URL
	p, _ := github.New(github.Config{Repository: "o/r", BaseURL: srv.URL})
	_, err := p.Check(context.Background(), updater.CheckRequest{
		CurrentVersion: "1.0.0",
		Platform:       "freebsd",
		Arch:           "amd64",
	})
	if err == nil || !strings.Contains(err.Error(), "no asset for freebsd/amd64") {
		t.Fatalf("expected no-asset error, got %v", err)
	}
}

func TestCheck_Prerelease_UsesReleasesList(t *testing.T) {
	var paths []string
	var host string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		// Return a single-item list (the prerelease endpoint).
		fmt.Fprintf(w, `[`+sampleReleaseJSON+`]`, host, host, host)
	}))
	defer srv.Close()
	host = srv.URL
	p, _ := github.New(github.Config{Repository: "o/r", BaseURL: srv.URL, Prerelease: true})
	rel, err := p.Check(context.Background(), updater.CheckRequest{CurrentVersion: "1.0.0", Platform: "darwin", Arch: "arm64"})
	if err != nil {
		t.Fatal(err)
	}
	if rel == nil {
		t.Fatal("expected release")
	}
	if len(paths) != 1 || paths[0] != "/repos/o/r/releases" {
		t.Errorf("path: %v (want /repos/o/r/releases for prerelease list)", paths)
	}
}

func TestCheck_Prerelease_SkipsDrafts(t *testing.T) {
	// The list endpoint includes drafts (visible to authenticated users with
	// push access). Drafts must never be surfaced as available updates.
	var paths []string
	var host string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path+"?"+r.URL.RawQuery)
		fmt.Fprintf(w, `[
		  {"tag_name": "v3.0.0-draft", "draft": true, "prerelease": true, "published_at": "2026-04-01T10:00:00Z", "assets": []},
		  `+sampleReleaseJSON+`
		]`, host, host, host)
	}))
	defer srv.Close()
	host = srv.URL

	p, _ := github.New(github.Config{Repository: "o/r", BaseURL: srv.URL, Prerelease: true, Token: "ghp_xxx"})
	rel, err := p.Check(context.Background(), updater.CheckRequest{
		CurrentVersion: "1.0.0",
		Platform:       "darwin",
		Arch:           "arm64",
	})
	if err != nil {
		t.Fatal(err)
	}
	if rel == nil {
		t.Fatal("expected the v2.0.0 release (drafts must be skipped)")
	}
	if rel.Version != "2.0.0" {
		t.Errorf("draft surfaced as available update: got version %q", rel.Version)
	}
	// Should have requested per_page>1 so we can skip over leading drafts.
	if len(paths) != 1 || !strings.Contains(paths[0], "per_page=") || strings.Contains(paths[0], "per_page=1&") || strings.HasSuffix(paths[0], "per_page=1") {
		t.Errorf("expected per_page>1, got %v", paths)
	}
}

func TestCheck_ChecksumSidecar_PopulatesVerification(t *testing.T) {
	body := []byte("hello-world")
	digest := sha256.Sum256(body)
	digestHex := hex.EncodeToString(digest[:])

	// host is assigned after httptest.NewServer; the handler closures only
	// read it at request time, so the late binding is safe.
	var host string
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/o/r/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{
		  "tag_name": "v2.0.0",
		  "name": "v2", "body": "",
		  "prerelease": false, "draft": false,
		  "published_at": "2026-03-01T10:00:00Z",
		  "html_url": "",
		  "assets": [
		    {"id": 1, "name": "app-linux-amd64.tar.gz", "content_type": "application/gzip", "size": %d, "browser_download_url": "%s/dl/asset"},
		    {"id": 2, "name": "SHA256SUMS", "content_type": "text/plain", "size": 64, "browser_download_url": "%s/dl/sums"}
		  ]
		}`, len(body), host, host)
	})
	mux.HandleFunc("/dl/sums", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s  app-linux-amd64.tar.gz\n", digestHex)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	host = srv.URL

	p, _ := github.New(github.Config{
		Repository:    "o/r",
		BaseURL:       srv.URL,
		ChecksumAsset: "SHA256SUMS",
	})
	rel, err := p.Check(context.Background(), updater.CheckRequest{
		CurrentVersion: "1.0.0",
		Platform:       "linux",
		Arch:           "amd64",
	})
	if err != nil {
		t.Fatal(err)
	}
	if rel == nil || rel.Verification == nil {
		t.Fatalf("expected verification populated, got %+v", rel)
	}
	if rel.Verification.DigestAlgo != "sha256" {
		t.Errorf("algo: %s", rel.Verification.DigestAlgo)
	}
	if !bytes.Equal(rel.Verification.Digest, digest[:]) {
		t.Errorf("digest mismatch: got %x", rel.Verification.Digest)
	}
}

func TestDownload_FollowsRedirect_StripsAuth(t *testing.T) {
	body := []byte("downloaded-bytes")
	var s3Auth string
	var s3Hits int32
	s3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&s3Hits, 1)
		s3Auth = r.Header.Get("Authorization")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
		w.Write(body)
	}))
	defer s3.Close()

	gh := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, s3.URL+"/signed", http.StatusFound)
	}))
	defer gh.Close()

	p, _ := github.New(github.Config{Repository: "o/r", BaseURL: gh.URL, Token: "ghp_secret"})
	rel := &updater.Release{
		Artifact: updater.Artifact{Filename: "app.bin", Size: int64(len(body))},
		Metadata: map[string]any{"github.asset.url": gh.URL + "/dl"},
	}
	var got bytes.Buffer
	var progressCalls int32
	err := p.Download(context.Background(), rel, &got, func(_, _ int64) { atomic.AddInt32(&progressCalls, 1) })
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got.Bytes(), body) {
		t.Errorf("bytes: %q", got.Bytes())
	}
	if atomic.LoadInt32(&s3Hits) != 1 {
		t.Errorf("s3 hits: %d", s3Hits)
	}
	if s3Auth != "" {
		t.Errorf("Authorization should not leak to S3: %q", s3Auth)
	}
	if atomic.LoadInt32(&progressCalls) == 0 {
		t.Error("expected progress callbacks")
	}
}

func TestNew_RequiresRepository(t *testing.T) {
	if _, err := github.New(github.Config{}); err == nil {
		t.Fatal("expected error")
	}
	if _, err := github.New(github.Config{Repository: "no-slash"}); err == nil {
		t.Fatal("expected error for bad repo format")
	}
}

func TestDefaultAssetMatcher(t *testing.T) {
	assets := []github.ReleaseAsset{
		{Name: "app-darwin-arm64.dmg"},
		{Name: "app-linux-amd64.tar.gz"},
		{Name: "app-linux-x86_64.tar.gz"}, // alt amd64 name
		{Name: "app-darwin-aarch64.zip"},  // alt arm64 name
		{Name: "checksums.txt"},
		{Name: "app-darwin-arm64.dmg.sig"},
	}
	cases := []struct {
		plat, arch string
		wantIndex  int
	}{
		{"darwin", "arm64", 0},
		{"linux", "amd64", 1},
		{"linux", "386", -1},     // x86 alias must not match the linux-x86_64 asset
		{"freebsd", "amd64", -1}, // no freebsd asset
		{"", "amd64", 1},         // empty plat picks first matching arch (skipping sidecars)
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("%s/%s", c.plat, c.arch), func(t *testing.T) {
			got := github.DefaultAssetMatcher(updater.CheckRequest{Platform: c.plat, Arch: c.arch}, assets)
			if got != c.wantIndex {
				t.Errorf("got %d, want %d", got, c.wantIndex)
			}
		})
	}
}

// 386 falls back to matching "x86", but "x86" is a substring of "x86_64".
// The matcher must not return an amd64 asset for a 386 request.
func TestDefaultAssetMatcher_386_DoesNotMatchX86_64(t *testing.T) {
	assets := []github.ReleaseAsset{
		{Name: "app-linux-x86_64.tar.gz"},
		{Name: "app-linux-amd64.deb"},
		{Name: "app-linux-i386.tar.gz"},
	}
	if got := github.DefaultAssetMatcher(updater.CheckRequest{Platform: "linux", Arch: "386"}, assets); got != 2 {
		t.Errorf("386 request should pick i386 asset (index 2), got %d", got)
	}
	// And the strict amd64 path is unaffected.
	if got := github.DefaultAssetMatcher(updater.CheckRequest{Platform: "linux", Arch: "amd64"}, assets); got != 0 {
		t.Errorf("amd64 request should still pick x86_64 asset (index 0), got %d", got)
	}
}

// Some projects encode the hash algorithm in the primary artifact's filename
// (e.g. \"myapp-1.0-windows-sha256.zip\"). The previous substring-based
// isChecksumName treated those as sidecars and DefaultAssetMatcher silently
// skipped them, leaving the matcher with no candidate.
func TestDefaultAssetMatcher_NameContainsAlgoButIsNotSidecar(t *testing.T) {
	assets := []github.ReleaseAsset{
		{Name: "myapp-1.0-windows-sha256-amd64.zip"},
		{Name: "SHA256SUMS"},
		{Name: "myapp-1.0-windows-amd64.zip.sha256"},
	}
	got := github.DefaultAssetMatcher(updater.CheckRequest{Platform: "windows", Arch: "amd64"}, assets)
	if got != 0 {
		t.Errorf("want primary artifact (index 0), got %d", got)
	}
}

func TestCheck_APIError_Surfaced(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		io.WriteString(w, `{"message":"rate limit"}`)
	}))
	defer srv.Close()
	p, _ := github.New(github.Config{Repository: "o/r", BaseURL: srv.URL})
	_, err := p.Check(context.Background(), updater.CheckRequest{CurrentVersion: "1.0.0"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "403") {
		t.Errorf("expected status in error, got %v", err)
	}
}

func TestProviderInterface(t *testing.T) {
	var _ updater.Provider = (*github.Provider)(nil)
}

