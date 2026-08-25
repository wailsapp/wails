package keygen_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/updater"
	"github.com/wailsapp/wails/v3/pkg/updater/providers/keygen"
)

const sampleUpgradeOK = `{
  "data": {
    "id": "30c64dcd-a74d-4f0d-8479-8745172a4817",
    "type": "releases",
    "attributes": {
      "name": "App v2.0.0",
      "description": "Major release",
      "channel": "stable",
      "status": "PUBLISHED",
      "tag": "latest",
      "version": "2.0.0",
      "metadata": {"sha256": "abc123"},
      "created": "2022-05-31T14:26:09.319Z"
    }
  },
  "included": [
    {
      "id": "0dad8516-f071-4573-bcea-d774e81c4a37",
      "type": "artifacts",
      "attributes": {
        "filename": "App-darwin-arm64.dmg",
        "filetype": "dmg",
        "filesize": 12345,
        "platform": "darwin",
        "arch": "arm64",
        "signature": "c2lnLWJ5dGVz",
        "checksum": "Y2hlY2tzdW0tYnl0ZXM",
        "status": "UPLOADED"
      }
    },
    {
      "id": "1aaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
      "type": "artifacts",
      "attributes": {
        "filename": "App-linux-amd64.tar.gz",
        "filetype": "gz",
        "filesize": 22222,
        "platform": "linux",
        "arch": "amd64",
        "status": "UPLOADED"
      }
    }
  ]
}`

func TestCheck_PicksMatchingArtifact_PopulatesVerification(t *testing.T) {
	var sawAuth, sawAccept string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if want := "/v1/accounts/acct/releases/1.0.0/upgrade"; r.URL.Path != want {
			t.Errorf("path: want %q, got %q", want, r.URL.Path)
		}
		sawAuth = r.Header.Get("Authorization")
		sawAccept = r.Header.Get("Accept")
		io.WriteString(w, sampleUpgradeOK)
	}))
	defer srv.Close()

	p, err := keygen.New(keygen.Config{
		Account: "acct",
		Token:   "admi-secret",
		BaseURL: srv.URL,
	})
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
		t.Fatal("expected release")
	}
	if rel.Version != "2.0.0" {
		t.Errorf("version: %q", rel.Version)
	}
	if rel.Artifact.Filename != "App-darwin-arm64.dmg" {
		t.Errorf("artifact: %q", rel.Artifact.Filename)
	}

	// Verification block must reflect keygen's SHA-512 + Ed25519ph mapping.
	if rel.Verification == nil {
		t.Fatal("expected Verification populated")
	}
	wantDigest, _ := base64.RawStdEncoding.DecodeString("Y2hlY2tzdW0tYnl0ZXM")
	wantSig, _ := base64.RawStdEncoding.DecodeString("c2lnLWJ5dGVz")
	if rel.Verification.DigestAlgo != "sha512" {
		t.Errorf("digest algo: %s", rel.Verification.DigestAlgo)
	}
	if !bytes.Equal(rel.Verification.Digest, wantDigest) {
		t.Errorf("digest: %x", rel.Verification.Digest)
	}
	if rel.Verification.SignatureAlgo != "ed25519ph" {
		t.Errorf("sig algo: %s", rel.Verification.SignatureAlgo)
	}
	if !bytes.Equal(rel.Verification.Signature, wantSig) {
		t.Errorf("signature: %x", rel.Verification.Signature)
	}

	if sawAuth != "Bearer admi-secret" {
		t.Errorf("auth: %q", sawAuth)
	}
	if !strings.Contains(sawAccept, "application/vnd.api+json") {
		t.Errorf("accept: %q", sawAccept)
	}
}

func TestCheck_UpToDate_404_TreatsAsNil(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, `{"errors":[{"title":"Not found","code":"NOT_FOUND"}]}`)
	}))
	defer srv.Close()
	p, _ := keygen.New(keygen.Config{Account: "acct", BaseURL: srv.URL})
	rel, err := p.Check(context.Background(), updater.CheckRequest{CurrentVersion: "1.0.0", Platform: "darwin", Arch: "arm64"})
	if err != nil {
		t.Fatal(err)
	}
	if rel != nil {
		t.Errorf("expected nil release, got %+v", rel)
	}
}

func TestCheck_LicenseAuth(t *testing.T) {
	var sawAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sawAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()
	p, _ := keygen.New(keygen.Config{Account: "acct", LicenseKey: "LKEY", BaseURL: srv.URL})
	_, _ = p.Check(context.Background(), updater.CheckRequest{CurrentVersion: "1.0.0"})
	if sawAuth != "License LKEY" {
		t.Errorf("auth: %q", sawAuth)
	}
}

func TestCheck_TokenWinsOverLicense(t *testing.T) {
	var sawAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sawAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()
	p, _ := keygen.New(keygen.Config{Account: "acct", Token: "prod-x", LicenseKey: "lk", BaseURL: srv.URL})
	_, _ = p.Check(context.Background(), updater.CheckRequest{CurrentVersion: "1.0.0"})
	if sawAuth != "Bearer prod-x" {
		t.Errorf("auth: %q", sawAuth)
	}
}

func TestCheck_APIError_Typed(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, `{"errors":[{"title":"Unauthorized","detail":"bad token","code":"TOKEN_INVALID"}]}`)
	}))
	defer srv.Close()
	p, _ := keygen.New(keygen.Config{Account: "acct", Token: "bad", BaseURL: srv.URL})
	_, err := p.Check(context.Background(), updater.CheckRequest{CurrentVersion: "1.0.0"})
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *keygen.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("want *keygen.APIError, got %T: %v", err, err)
	}
	if apiErr.StatusCode != 401 || apiErr.Code != "TOKEN_INVALID" {
		t.Errorf("unexpected: %+v", apiErr)
	}
}

func TestCheck_NoMatchingArtifactReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, sampleUpgradeOK)
	}))
	defer srv.Close()
	p, _ := keygen.New(keygen.Config{Account: "acct", BaseURL: srv.URL})
	_, err := p.Check(context.Background(), updater.CheckRequest{
		CurrentVersion: "1.0.0",
		Platform:       "windows",
		Arch:           "amd64",
	})
	if err == nil || !strings.Contains(err.Error(), "no artifact for windows/amd64") {
		t.Fatalf("expected no-artifact error, got %v", err)
	}
}

func TestCheck_FallbackToArtifactsList_WhenNoIncluded(t *testing.T) {
	const releaseID = "30c64dcd-a74d-4f0d-8479-8745172a4817"
	const upgradeNoArts = `{
      "data": {
        "id": "30c64dcd-a74d-4f0d-8479-8745172a4817",
        "type": "releases",
        "attributes": {
          "name": "v2", "description": "", "channel": "stable",
          "status": "PUBLISHED", "tag": null, "version": "2.0.0",
          "metadata": {}, "created": "2022-05-31T14:26:09.319Z"
        }
      }
    }`
	const artsList = `{
      "data": [{
        "id": "0dad8516-f071-4573-bcea-d774e81c4a37",
        "type": "artifacts",
        "attributes": {
          "filename": "App-darwin-arm64.dmg", "filetype": "dmg",
          "filesize": 12345, "platform": "darwin", "arch": "arm64",
          "status": "UPLOADED"
        }
      }]
    }`

	var calls []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.URL.Path)
		switch r.URL.Path {
		case "/v1/accounts/acct/releases/1.0.0/upgrade":
			io.WriteString(w, upgradeNoArts)
		case "/v1/accounts/acct/releases/" + releaseID + "/artifacts":
			io.WriteString(w, artsList)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	p, _ := keygen.New(keygen.Config{Account: "acct", BaseURL: srv.URL})
	rel, err := p.Check(context.Background(), updater.CheckRequest{CurrentVersion: "1.0.0", Platform: "darwin", Arch: "arm64"})
	if err != nil {
		t.Fatal(err)
	}
	if rel == nil || rel.Artifact.Filename != "App-darwin-arm64.dmg" {
		t.Fatalf("unexpected release: %+v", rel)
	}
	if len(calls) != 2 {
		t.Errorf("expected 2 API calls, got %v", calls)
	}
}

func TestDownload_StripsAuthOnRedirect(t *testing.T) {
	var s3Auth string
	var s3Hits int32
	s3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&s3Hits, 1)
		s3Auth = r.Header.Get("Authorization")
		w.Header().Set("Content-Length", "11")
		w.Write([]byte("hello-world"))
	}))
	defer s3.Close()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/accounts/acct/artifacts/App-darwin-arm64.dmg" {
			t.Errorf("path: %s", r.URL.Path)
		}
		w.Header().Set("Location", s3.URL+"/signed")
		w.WriteHeader(http.StatusSeeOther)
	}))
	defer srv.Close()

	p, _ := keygen.New(keygen.Config{Account: "acct", Token: "admi-secret", BaseURL: srv.URL})

	var buf bytes.Buffer
	var ticks int32
	err := p.Download(context.Background(), &updater.Release{
		Artifact: updater.Artifact{Filename: "App-darwin-arm64.dmg", Size: 11},
	}, &buf, func(_, _ int64) { atomic.AddInt32(&ticks, 1) })
	if err != nil {
		t.Fatal(err)
	}
	if buf.String() != "hello-world" {
		t.Errorf("body: %q", buf.String())
	}
	if atomic.LoadInt32(&s3Hits) != 1 {
		t.Errorf("s3 hits: %d", s3Hits)
	}
	if s3Auth != "" {
		t.Errorf("Authorization leaked: %q", s3Auth)
	}
	if atomic.LoadInt32(&ticks) == 0 {
		t.Error("expected progress")
	}
}

// Check must stash the chosen artifact's keygen.sh ID under
// rel.Metadata["keygen.artifact.id"] so a follow-up Download targets the
// artifact by its unique ID rather than its filename. Filenames are not
// unique across platforms (installer.exe for both amd64 and arm64).
func TestCheck_StashesArtifactID(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, sampleUpgradeOK)
	}))
	defer srv.Close()
	p, _ := keygen.New(keygen.Config{Account: "acct", BaseURL: srv.URL})
	rel, err := p.Check(context.Background(), updater.CheckRequest{
		CurrentVersion: "1.0.0",
		Platform:       "darwin",
		Arch:           "arm64",
	})
	if err != nil {
		t.Fatal(err)
	}
	if rel == nil || rel.Metadata == nil {
		t.Fatalf("missing rel/metadata: %+v", rel)
	}
	const wantID = "0dad8516-f071-4573-bcea-d774e81c4a37"
	if got, _ := rel.Metadata["keygen.artifact.id"].(string); got != wantID {
		t.Errorf("keygen.artifact.id: got %q, want %q", got, wantID)
	}
}

// Download prefers the artifact ID over the filename when Metadata carries
// one, so two artifacts sharing a filename across platforms can still be
// fetched deterministically.
func TestDownload_PrefersArtifactID(t *testing.T) {
	const artifactID = "0dad8516-f071-4573-bcea-d774e81c4a37"
	var path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		w.Header().Set("Content-Length", "2")
		w.Write([]byte("ok"))
	}))
	defer srv.Close()

	p, _ := keygen.New(keygen.Config{Account: "acct", BaseURL: srv.URL})
	rel := &updater.Release{
		Artifact: updater.Artifact{Filename: "installer.exe", Size: 2},
		Metadata: map[string]any{"keygen.artifact.id": artifactID},
	}
	var buf bytes.Buffer
	if err := p.Download(context.Background(), rel, &buf, func(_, _ int64) {}); err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(path, "/artifacts/"+artifactID) {
		t.Errorf("expected download by ID, got path %q", path)
	}
}

// pickArtifact must tolerate the operator-defined platform/arch aliases
// keygen.sh release operators commonly use ("macos", "x86_64", "aarch64").
// Previously a strict equality check dropped these releases silently.
func TestCheck_NormalisesPlatformArchAliases(t *testing.T) {
	const aliasFeed = `{
      "data": {
        "id": "30c64dcd-a74d-4f0d-8479-8745172a4817",
        "type": "releases",
        "attributes": {
          "name": "v2", "description": "", "channel": "stable",
          "status": "PUBLISHED", "tag": "latest", "version": "2.0.0",
          "metadata": {}, "created": "2022-05-31T14:26:09.319Z"
        }
      },
      "included": [{
        "id": "ABC",
        "type": "artifacts",
        "attributes": {
          "filename": "App-macos-aarch64.dmg", "filetype": "dmg",
          "filesize": 12345,
          "platform": "macos", "arch": "aarch64",
          "status": "UPLOADED"
        }
      }]
    }`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, aliasFeed)
	}))
	defer srv.Close()
	p, _ := keygen.New(keygen.Config{Account: "acct", BaseURL: srv.URL})
	rel, err := p.Check(context.Background(), updater.CheckRequest{
		CurrentVersion: "1.0.0",
		Platform:       "darwin", // alias for macos
		Arch:           "arm64",  // alias for aarch64
	})
	if err != nil {
		t.Fatal(err)
	}
	if rel == nil || rel.Artifact.Filename != "App-macos-aarch64.dmg" {
		t.Fatalf("expected alias match, got %+v", rel)
	}
}

func TestNew_RequiresAccount(t *testing.T) {
	if _, err := keygen.New(keygen.Config{}); err == nil {
		t.Fatal("expected error")
	}
}

func TestCheck_RejectsEmptyVersion(t *testing.T) {
	p, _ := keygen.New(keygen.Config{Account: "acct"})
	if _, err := p.Check(context.Background(), updater.CheckRequest{}); err == nil {
		t.Fatal("expected error")
	}
}

func TestProviderInterface(t *testing.T) {
	var _ updater.Provider = (*keygen.Provider)(nil)
}
