package endpoint

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/updater"
)

func checkReq() updater.CheckRequest {
	return updater.CheckRequest{
		CurrentVersion: "1.0.0",
		Platform:       "darwin",
		Arch:           "arm64",
	}
}

func manifestJSON(t *testing.T, m map[string]any) []byte {
	t.Helper()
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestNew_RequiresURL(t *testing.T) {
	if _, err := New(Config{}); err == nil {
		t.Fatal("expected error for empty URL")
	}
	if _, err := New(Config{URL: "   "}); err == nil {
		t.Fatal("expected error for blank URL")
	}
}

func TestCheck_DynamicQueryParams(t *testing.T) {
	var gotQuery map[string]string
	var gotAuth, gotAccept string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = map[string]string{}
		for k := range r.URL.Query() {
			gotQuery[k] = r.URL.Query().Get(k)
		}
		gotAuth = r.Header.Get("Authorization")
		gotAccept = r.Header.Get("Accept")
		w.Write(manifestJSON(t, map[string]any{
			"version": "2.0.0",
			"channel": "stable",
			"name":    "Big Release",
			"notes":   "## Changes",
			"artifacts": []map[string]any{
				{"url": "/dl/app-darwin-arm64.zip", "platform": "darwin", "arch": "arm64", "size": 42},
			},
		}))
	}))
	defer srv.Close()

	p, err := New(Config{
		URL:     srv.URL + "/manifest",
		Channel: "stable",
		Headers: map[string]string{"Authorization": "License abc-123"},
	})
	if err != nil {
		t.Fatal(err)
	}
	rel, err := p.Check(context.Background(), checkReq())
	if err != nil {
		t.Fatal(err)
	}
	if rel == nil {
		t.Fatal("expected a release")
	}

	want := map[string]string{"platform": "darwin", "arch": "arm64", "version": "1.0.0", "channel": "stable"}
	for k, v := range want {
		if gotQuery[k] != v {
			t.Errorf("query %s = %q, want %q", k, gotQuery[k], v)
		}
	}
	if gotAuth != "License abc-123" {
		t.Errorf("Authorization = %q", gotAuth)
	}
	if gotAccept != "application/json" {
		t.Errorf("Accept = %q", gotAccept)
	}
	if rel.Version != "2.0.0" || rel.Name != "Big Release" || rel.Notes != "## Changes" {
		t.Errorf("release fields: %+v", rel)
	}
	if rel.Artifact.Filename != "app-darwin-arm64.zip" || rel.Artifact.Filetype != "zip" || rel.Artifact.Size != 42 {
		t.Errorf("artifact: %+v", rel.Artifact)
	}
	wantURL := srv.URL + "/dl/app-darwin-arm64.zip"
	if got := rel.Metadata["endpoint.artifact.url"]; got != wantURL {
		t.Errorf("artifact url = %v, want %s", got, wantURL)
	}
}

func TestCheck_PlaceholdersConsumeParams(t *testing.T) {
	var gotPath, gotRawQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotRawQuery = r.URL.RawQuery
		w.Write(manifestJSON(t, map[string]any{
			"version":   "2.0.0",
			"artifacts": []map[string]any{{"url": "app.zip"}},
		}))
	}))
	defer srv.Close()

	p, err := New(Config{URL: srv.URL + "/updates/{{platform}}/{{arch}}/{{channel}}.json", Channel: "beta"})
	if err != nil {
		t.Fatal(err)
	}
	rel, err := p.Check(context.Background(), checkReq())
	if err != nil {
		t.Fatal(err)
	}
	if rel == nil {
		t.Fatal("expected a release")
	}
	if gotPath != "/updates/darwin/arm64/beta.json" {
		t.Errorf("path = %q", gotPath)
	}
	// version has no placeholder, so it still travels as a query parameter;
	// platform/arch/channel were consumed and must not be duplicated.
	if gotRawQuery != "version=1.0.0" {
		t.Errorf("query = %q", gotRawQuery)
	}
	// Relative artifact URL resolves against the manifest location.
	wantURL := srv.URL + "/updates/darwin/arm64/app.zip"
	if got := rel.Metadata["endpoint.artifact.url"]; got != wantURL {
		t.Errorf("artifact url = %v, want %s", got, wantURL)
	}
}

func TestCheck_UpToDate(t *testing.T) {
	for _, version := range []string{"1.0.0", "0.9.9", "v1.0.0"} {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(manifestJSON(t, map[string]any{
				"version":   version,
				"artifacts": []map[string]any{{"url": "app.zip"}},
			}))
		}))
		p, _ := New(Config{URL: srv.URL})
		rel, err := p.Check(context.Background(), checkReq())
		srv.Close()
		if err != nil {
			t.Fatalf("version %s: %v", version, err)
		}
		if rel != nil {
			t.Errorf("version %s: expected nil release, got %+v", version, rel)
		}
	}
}

func TestCheck_NoContentAndNotFound(t *testing.T) {
	for _, status := range []int{http.StatusNoContent, http.StatusNotFound} {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(status)
		}))
		p, _ := New(Config{URL: srv.URL})
		rel, err := p.Check(context.Background(), checkReq())
		srv.Close()
		if err != nil {
			t.Fatalf("status %d: %v", status, err)
		}
		if rel != nil {
			t.Errorf("status %d: expected nil release", status)
		}
	}
}

func TestCheck_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()
	p, _ := New(Config{URL: srv.URL})
	if _, err := p.Check(context.Background(), checkReq()); err == nil {
		t.Fatal("expected error on HTTP 500")
	}
}

func TestCheck_ChannelMismatch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(manifestJSON(t, map[string]any{
			"version":   "2.0.0",
			"channel":   "beta",
			"artifacts": []map[string]any{{"url": "app.zip"}},
		}))
	}))
	defer srv.Close()
	p, _ := New(Config{URL: srv.URL, Channel: "stable"})
	rel, err := p.Check(context.Background(), checkReq())
	if err != nil {
		t.Fatal(err)
	}
	if rel != nil {
		t.Fatalf("expected nil release for channel mismatch, got %+v", rel)
	}
}

func TestCheck_SchemaVersionTooNew(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(manifestJSON(t, map[string]any{
			"schemaVersion": 2,
			"version":       "2.0.0",
			"artifacts":     []map[string]any{{"url": "app.zip"}},
		}))
	}))
	defer srv.Close()
	p, _ := New(Config{URL: srv.URL})
	if _, err := p.Check(context.Background(), checkReq()); err == nil {
		t.Fatal("expected error for unsupported schemaVersion")
	}
}

func TestCheck_ArtifactSelection(t *testing.T) {
	// Aliases on the manifest side must match Go runtime values, and the
	// darwin/arm64 request must skip non-matching entries.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(manifestJSON(t, map[string]any{
			"version": "2.0.0",
			"artifacts": []map[string]any{
				{"url": "win.exe", "platform": "win", "arch": "x86_64"},
				{"url": "mac.zip", "platform": "macos", "arch": "aarch64"},
				{"url": "linux.tgz", "platform": "linux", "arch": "x86_64"},
			},
		}))
	}))
	defer srv.Close()
	p, _ := New(Config{URL: srv.URL})
	rel, err := p.Check(context.Background(), checkReq())
	if err != nil {
		t.Fatal(err)
	}
	if rel == nil || rel.Artifact.Filename != "mac.zip" {
		t.Fatalf("expected mac.zip artifact, got %+v", rel)
	}
}

func TestCheck_WildcardArtifact(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(manifestJSON(t, map[string]any{
			"version":   "2.0.0",
			"artifacts": []map[string]any{{"url": "universal.zip"}},
		}))
	}))
	defer srv.Close()
	p, _ := New(Config{URL: srv.URL})
	rel, err := p.Check(context.Background(), checkReq())
	if err != nil {
		t.Fatal(err)
	}
	if rel == nil || rel.Artifact.Filename != "universal.zip" {
		t.Fatalf("expected universal.zip artifact, got %+v", rel)
	}
}

func TestCheck_NoMatchingArtifact(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(manifestJSON(t, map[string]any{
			"version":   "2.0.0",
			"artifacts": []map[string]any{{"url": "win.exe", "platform": "windows", "arch": "amd64"}},
		}))
	}))
	defer srv.Close()
	p, _ := New(Config{URL: srv.URL})
	if _, err := p.Check(context.Background(), checkReq()); err == nil {
		t.Fatal("expected error when no artifact matches the platform")
	}
}

func TestCheck_Verification(t *testing.T) {
	payload := []byte("artifact bytes")
	digest := sha512.Sum512(payload)
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	sig, err := priv.Sign(nil, digest[:], &ed25519.Options{Hash: crypto.SHA512})
	if err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(manifestJSON(t, map[string]any{
			"version": "2.0.0",
			"artifacts": []map[string]any{{
				"url":           "app.zip",
				"digestAlgo":    "SHA512",
				"digest":        base64.RawStdEncoding.EncodeToString(digest[:]),
				"signatureAlgo": "ed25519ph",
				"signature":     base64.StdEncoding.EncodeToString(sig),
			}},
		}))
	}))
	defer srv.Close()

	p, _ := New(Config{URL: srv.URL})
	rel, err := p.Check(context.Background(), checkReq())
	if err != nil {
		t.Fatal(err)
	}
	if rel == nil || rel.Verification == nil {
		t.Fatal("expected release with verification")
	}
	v := rel.Verification
	if v.DigestAlgo != "sha512" || !bytes.Equal(v.Digest, digest[:]) {
		t.Errorf("digest mapping wrong: algo=%q", v.DigestAlgo)
	}
	if v.SignatureAlgo != "ed25519ph" || !bytes.Equal(v.Signature, sig) {
		t.Errorf("signature mapping wrong: algo=%q", v.SignatureAlgo)
	}
	// The mapped signature must verify under the framework's ed25519ph rules.
	if err := ed25519.VerifyWithOptions(pub, digest[:], v.Signature, &ed25519.Options{Hash: crypto.SHA512}); err != nil {
		t.Errorf("round-tripped signature does not verify: %v", err)
	}
}

func TestCheck_SignatureWithoutAlgoFails(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(manifestJSON(t, map[string]any{
			"version": "2.0.0",
			"artifacts": []map[string]any{{
				"url":       "app.zip",
				"signature": base64.StdEncoding.EncodeToString([]byte("some signature")),
			}},
		}))
	}))
	defer srv.Close()
	p, _ := New(Config{URL: srv.URL})
	if _, err := p.Check(context.Background(), checkReq()); err == nil {
		t.Fatal("expected error for signature without signatureAlgo, got nil")
	}
}

func TestCheck_UndecodableVerificationFails(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(manifestJSON(t, map[string]any{
			"version": "2.0.0",
			"artifacts": []map[string]any{{
				"url":        "app.zip",
				"digestAlgo": "sha512",
				"digest":     "!!! not base64 !!!",
			}},
		}))
	}))
	defer srv.Close()
	p, _ := New(Config{URL: srv.URL})
	if _, err := p.Check(context.Background(), checkReq()); err == nil {
		t.Fatal("expected error for undecodable digest, got nil")
	}
}

func TestDownload_StreamsWithProgress(t *testing.T) {
	payload := bytes.Repeat([]byte("wails"), 100_000) // 500 KB
	mux := http.NewServeMux()
	mux.HandleFunc("/manifest", func(w http.ResponseWriter, r *http.Request) {
		w.Write(manifestJSON(t, map[string]any{
			"version":   "2.0.0",
			"artifacts": []map[string]any{{"url": "/dl/app.zip", "size": len(payload)}},
		}))
	})
	var gotAuth string
	mux.HandleFunc("/dl/app.zip", func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.Write(payload)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	p, _ := New(Config{URL: srv.URL + "/manifest", Headers: map[string]string{"Authorization": "License k"}})
	rel, err := p.Check(context.Background(), checkReq())
	if err != nil || rel == nil {
		t.Fatalf("check: rel=%v err=%v", rel, err)
	}

	var buf bytes.Buffer
	var lastWritten, lastTotal int64
	calls := 0
	err = p.Download(context.Background(), rel, &buf, func(written, total int64) {
		lastWritten, lastTotal = written, total
		calls++
	})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf.Bytes(), payload) {
		t.Fatal("downloaded bytes differ")
	}
	if calls == 0 || lastWritten != int64(len(payload)) || lastTotal != int64(len(payload)) {
		t.Errorf("progress: calls=%d written=%d total=%d", calls, lastWritten, lastTotal)
	}
	// Same host as the manifest, so the configured header is sent.
	if gotAuth != "License k" {
		t.Errorf("same-host download Authorization = %q", gotAuth)
	}
}

func TestDownload_CrossOriginDropsAuth(t *testing.T) {
	var cdnAuth string
	gotCDN := false
	cdn := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotCDN = true
		cdnAuth = r.Header.Get("Authorization")
		w.Write([]byte("bytes"))
	}))
	defer cdn.Close()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(manifestJSON(t, map[string]any{
			"version":   "2.0.0",
			"artifacts": []map[string]any{{"url": cdn.URL + "/app.zip"}},
		}))
	}))
	defer srv.Close()

	p, _ := New(Config{URL: srv.URL, Headers: map[string]string{"Authorization": "License secret"}})
	rel, err := p.Check(context.Background(), checkReq())
	if err != nil || rel == nil {
		t.Fatalf("check: rel=%v err=%v", rel, err)
	}
	var buf bytes.Buffer
	if err := p.Download(context.Background(), rel, &buf, func(int64, int64) {}); err != nil {
		t.Fatal(err)
	}
	if !gotCDN {
		t.Fatal("CDN was never hit")
	}
	if cdnAuth != "" {
		t.Errorf("Authorization leaked cross-origin: %q", cdnAuth)
	}
}

func TestDownload_RedirectStripsAuth(t *testing.T) {
	var cdnAuth string
	cdn := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cdnAuth = r.Header.Get("Authorization")
		w.Write([]byte("bytes"))
	}))
	defer cdn.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/manifest", func(w http.ResponseWriter, r *http.Request) {
		w.Write(manifestJSON(t, map[string]any{
			"version":   "2.0.0",
			"artifacts": []map[string]any{{"url": "/dl/app.zip"}},
		}))
	})
	mux.HandleFunc("/dl/app.zip", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, cdn.URL+"/app.zip", http.StatusSeeOther)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	p, _ := New(Config{URL: srv.URL + "/manifest", Headers: map[string]string{"Authorization": "License secret"}})
	rel, err := p.Check(context.Background(), checkReq())
	if err != nil || rel == nil {
		t.Fatalf("check: rel=%v err=%v", rel, err)
	}
	var buf bytes.Buffer
	if err := p.Download(context.Background(), rel, &buf, func(int64, int64) {}); err != nil {
		t.Fatal(err)
	}
	if cdnAuth != "" {
		t.Errorf("Authorization leaked on cross-origin redirect: %q", cdnAuth)
	}
}

func TestBuildURL_PreservesExistingQuery(t *testing.T) {
	p, _ := New(Config{URL: "https://example.com/updates?key=abc"})
	got, err := p.buildURL(checkReq())
	if err != nil {
		t.Fatal(err)
	}
	u := mustParse(t, got)
	q := u.Query()
	if q.Get("key") != "abc" || q.Get("platform") != "darwin" || q.Get("arch") != "arm64" || q.Get("version") != "1.0.0" {
		t.Errorf("query = %q", u.RawQuery)
	}
	if q.Has("channel") {
		t.Error("channel sent despite not being configured")
	}
}

func TestResolveURL_RejectsNonHTTP(t *testing.T) {
	if _, err := resolveURL("https://example.com/m.json", "file:///etc/passwd"); err == nil {
		t.Fatal("expected error for file: scheme")
	}
}

func mustParse(t *testing.T, s string) *url.URL {
	t.Helper()
	u, err := url.Parse(s)
	if err != nil {
		t.Fatal(err)
	}
	return u
}

func TestHeadersAllowedFor(t *testing.T) {
	tests := []struct {
		manifest, target string
		want             bool
	}{
		{"https://u.example.com/m.json", "https://u.example.com/app.zip", true},
		{"https://u.example.com/m.json", "http://u.example.com/app.zip", false}, // https→http downgrade leaks auth in cleartext
		{"http://u.example.com/m.json", "https://u.example.com/app.zip", true},  // upgrade is fine
		{"http://u.example.com/m.json", "http://u.example.com/app.zip", true},
		{"https://u.example.com/m.json", "https://cdn.example.com/app.zip", false},
		{"https://u.example.com:8443/m.json", "https://u.example.com/app.zip", false}, // port is part of the host
	}
	for _, tt := range tests {
		if got := headersAllowedFor(tt.manifest, tt.target); got != tt.want {
			t.Errorf("headersAllowedFor(%q, %q) = %v, want %v", tt.manifest, tt.target, got, tt.want)
		}
	}
}

func TestRedirectStripsAuthOnSchemeDowngrade(t *testing.T) {
	// Same host, https → http: the transport would happily copy the
	// Authorization header (Go only strips on host changes), so our
	// CheckRedirect must drop it.
	c := wrapStripAuthOnRedirect(&http.Client{})
	mkReq := func(u string) *http.Request {
		r, err := http.NewRequest(http.MethodGet, u, nil)
		if err != nil {
			t.Fatal(err)
		}
		return r
	}
	next := mkReq("http://updates.example.com/app.zip")
	next.Header.Set("Authorization", "License secret")
	via := []*http.Request{mkReq("https://updates.example.com/app.zip")}
	if err := c.CheckRedirect(next, via); err != nil {
		t.Fatal(err)
	}
	if got := next.Header.Get("Authorization"); got != "" {
		t.Errorf("Authorization survived an https→http downgrade redirect: %q", got)
	}

	// Same host and scheme: the header must survive.
	next = mkReq("https://updates.example.com/other.zip")
	next.Header.Set("Authorization", "License secret")
	if err := c.CheckRedirect(next, via); err != nil {
		t.Fatal(err)
	}
	if next.Header.Get("Authorization") == "" {
		t.Error("Authorization dropped on a same-origin redirect")
	}
}
