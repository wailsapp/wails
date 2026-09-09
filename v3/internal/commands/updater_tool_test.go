package commands

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/updater"
	"github.com/wailsapp/wails/v3/pkg/updater/providers/endpoint"
)

func genTestKey(t *testing.T) (string, ed25519.PublicKey) {
	t.Helper()
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "updater.key")
	if err := UpdaterGenKey(&UpdaterGenKeyOptions{Out: keyPath}); err != nil {
		t.Fatalf("genkey: %v", err)
	}
	pubAny, err := loadUpdaterPublicKey(keyPath + ".pub")
	if err != nil {
		t.Fatalf("load public key: %v", err)
	}
	pub, ok := pubAny.(ed25519.PublicKey)
	if !ok {
		t.Fatalf("genkey produced a %T public key, want ed25519", pubAny)
	}
	return keyPath, pub
}

func writeArtifact(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestUpdaterGenKey(t *testing.T) {
	keyPath, pub := genTestKey(t)

	if runtime.GOOS != "windows" {
		info, err := os.Stat(keyPath)
		if err != nil {
			t.Fatal(err)
		}
		if perm := info.Mode().Perm(); perm != 0o600 {
			t.Errorf("private key permissions = %o, want 600", perm)
		}
	}

	// The keypair must actually correspond.
	priv, err := loadUpdaterPrivateKey(keyPath)
	if err != nil {
		t.Fatal(err)
	}
	if !pub.Equal(priv.Public()) {
		t.Error("public key file does not match the private key")
	}

	// A second run must refuse to clobber the key.
	if err := UpdaterGenKey(&UpdaterGenKeyOptions{Out: keyPath}); err == nil {
		t.Error("expected an error when the key already exists")
	}
	if err := UpdaterGenKey(&UpdaterGenKeyOptions{Out: keyPath, Force: true}); err != nil {
		t.Errorf("force overwrite failed: %v", err)
	}
}

func TestUpdaterPrivateKeyBase64Forms(t *testing.T) {
	// CI secret stores often hold the key as base64 rather than PEM.
	_, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	for name, content := range map[string]string{
		"full.b64": base64.StdEncoding.EncodeToString(priv),
		"seed.b64": base64.StdEncoding.EncodeToString(priv.Seed()) + "\n",
	} {
		p := writeArtifact(t, dir, name, content)
		loaded, err := loadUpdaterPrivateKey(p)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if !loaded.Equal(priv) {
			t.Errorf("%s: loaded key differs", name)
		}
	}
}

func TestSignUpdaterArtifact(t *testing.T) {
	keyPath, pub := genTestKey(t)
	priv, err := loadUpdaterPrivateKey(keyPath)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	content := "artifact payload"
	p := writeArtifact(t, dir, "MyApp-darwin-arm64.zip", content)

	sig, err := signUpdaterArtifact(p, priv)
	if err != nil {
		t.Fatal(err)
	}
	if sig.Size != int64(len(content)) {
		t.Errorf("size = %d, want %d", sig.Size, len(content))
	}
	if sig.DigestAlgo != "sha512" || sig.SignatureAlgo != "ed25519ph" {
		t.Errorf("algos = %s/%s", sig.DigestAlgo, sig.SignatureAlgo)
	}
	wantDigest := sha512.Sum512([]byte(content))
	gotDigest, _ := base64.StdEncoding.DecodeString(sig.Digest)
	if string(gotDigest) != string(wantDigest[:]) {
		t.Error("digest mismatch")
	}
	rawSig, _ := base64.StdEncoding.DecodeString(sig.Signature)
	if err := ed25519.VerifyWithOptions(pub, wantDigest[:], rawSig, &ed25519.Options{Hash: crypto.SHA512}); err != nil {
		t.Errorf("signature does not verify: %v", err)
	}
}

func TestInferUpdaterPlatformArch(t *testing.T) {
	tests := []struct {
		name     string
		platform string
		arch     string
	}{
		{"MyApp-2.1.0-darwin-arm64.zip", "darwin", "arm64"},
		{"MyApp_2.1.0_macOS_aarch64.dmg", "darwin", "arm64"},
		{"MyApp-darwin-universal.zip", "darwin", ""},
		{"MyApp-windows-amd64.zip", "windows", "amd64"},
		{"MyApp-win64-x86_64-setup.exe", "windows", "amd64"},
		{"MyApp-Setup.exe", "windows", ""},
		{"MyApp-linux-x86_64.AppImage", "linux", "amd64"},
		{"myapp_2.1.0_amd64.deb", "linux", "amd64"},
		{"MyApp-linux64.tar.gz", "linux", ""}, // x64 must NOT match inside linux64
		{"MyApp-osx-x64.zip", "darwin", "amd64"},
		{"MyApp-win-ia32.msi", "windows", "386"},
		{"MyApp.tar.gz", "", ""},
	}
	for _, tt := range tests {
		if got := inferUpdaterPlatform(tt.name); got != tt.platform {
			t.Errorf("inferUpdaterPlatform(%q) = %q, want %q", tt.name, got, tt.platform)
		}
		if got, _ := inferUpdaterArch(tt.name); got != tt.arch {
			t.Errorf("inferUpdaterArch(%q) = %q, want %q", tt.name, got, tt.arch)
		}
	}
}

func TestUpdaterManifestAndVerify(t *testing.T) {
	keyPath, _ := genTestKey(t)
	dir := t.TempDir()
	writeArtifact(t, dir, "MyApp-2.1.0-darwin-arm64.zip", "darwin payload")
	writeArtifact(t, dir, "MyApp-2.1.0-windows-amd64.zip", "windows payload")
	writeArtifact(t, dir, "MyApp-2.1.0-linux-x86_64.AppImage", "linux payload")
	writeArtifact(t, dir, "SHA256SUMS.txt", "should be skipped")
	notesPath := writeArtifact(t, dir, "notes.md", "## What's new")
	manifestPath := filepath.Join(dir, "manifest.json")

	err := UpdaterManifest(&UpdaterManifestOptions{
		Version:   "2.1.0",
		Channel:   "stable",
		Name:      "Summer Release",
		NotesFile: notesPath,
		Key:       keyPath,
		URLPrefix: "https://cdn.example.com/myapp/2.1.0",
		Output:    manifestPath,
	}, []string{dir})
	if err != nil {
		t.Fatalf("manifest: %v", err)
	}

	raw, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	var m updaterManifest
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatal(err)
	}
	if m.SchemaVersion != 1 || m.Version != "2.1.0" || m.Channel != "stable" {
		t.Errorf("header fields wrong: %+v", m)
	}
	if m.Notes != "## What's new" {
		t.Errorf("notes = %q", m.Notes)
	}
	if m.PublishedAt == "" {
		t.Error("publishedAt not set")
	}
	if len(m.Artifacts) != 3 {
		t.Fatalf("got %d artifacts, want 3 (sidecar files must be skipped): %+v", len(m.Artifacts), m.Artifacts)
	}
	byPlatform := map[string]updaterManifestArtifact{}
	for _, a := range m.Artifacts {
		byPlatform[a.Platform] = a
		if a.Signature == "" || a.SignatureAlgo != "ed25519ph" || a.DigestAlgo != "sha512" {
			t.Errorf("artifact %s not signed as expected: %+v", a.URL, a)
		}
	}
	if a := byPlatform["darwin"]; a.Arch != "arm64" || a.URL != "https://cdn.example.com/myapp/2.1.0/MyApp-2.1.0-darwin-arm64.zip" {
		t.Errorf("darwin artifact wrong: %+v", a)
	}
	if a := byPlatform["linux"]; a.Arch != "amd64" {
		t.Errorf("linux artifact wrong: %+v", a)
	}

	// verify must pass with the right key...
	err = UpdaterVerify(&UpdaterVerifyOptions{Manifest: manifestPath, PublicKey: keyPath + ".pub"})
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	// ...fail without a key when signatures are present...
	if err := UpdaterVerify(&UpdaterVerifyOptions{Manifest: manifestPath}); err == nil {
		t.Error("verify without key should fail on signed manifest")
	}
	// ...fail with the wrong key...
	otherKey, _ := genTestKey(t)
	if err := UpdaterVerify(&UpdaterVerifyOptions{Manifest: manifestPath, PublicKey: otherKey + ".pub"}); err == nil {
		t.Error("verify with wrong key should fail")
	}
	// ...and fail when an artifact is tampered with.
	writeArtifact(t, dir, "MyApp-2.1.0-windows-amd64.zip", "tampered payload")
	if err := UpdaterVerify(&UpdaterVerifyOptions{Manifest: manifestPath, PublicKey: keyPath + ".pub"}); err == nil {
		t.Error("verify should fail after tampering")
	}
}

func TestUpdaterManifestRequiresValidVersion(t *testing.T) {
	if err := UpdaterManifest(&UpdaterManifestOptions{}, nil); err == nil {
		t.Error("missing version should fail")
	}
	if err := UpdaterManifest(&UpdaterManifestOptions{Version: "not-a-version"}, []string{"x"}); err == nil {
		t.Error("invalid version should fail")
	}
}

// TestUpdaterManifestEndToEnd round-trips a CLI-generated manifest through
// the endpoint provider: serve the output directory over HTTP, run a Check
// as a darwin/arm64 client on an older version, download the artifact and
// verify the digest and signature the provider extracted — everything a
// shipped application would do.
func TestUpdaterManifestEndToEnd(t *testing.T) {
	keyPath, pub := genTestKey(t)
	dir := t.TempDir()
	content := "the real update payload"
	writeArtifact(t, dir, "MyApp-2.1.0-darwin-arm64.zip", content)
	manifestPath := filepath.Join(dir, "manifest.json")

	err := UpdaterManifest(&UpdaterManifestOptions{
		Version: "2.1.0",
		Key:     keyPath,
		Output:  manifestPath,
	}, []string{dir})
	if err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(http.FileServer(http.Dir(dir)))
	defer srv.Close()

	ep, err := endpoint.New(endpoint.Config{URL: srv.URL + "/manifest.json"})
	if err != nil {
		t.Fatal(err)
	}
	rel, err := ep.Check(context.Background(), updater.CheckRequest{
		CurrentVersion: "2.0.0",
		Platform:       "darwin",
		Arch:           "arm64",
	})
	if err != nil {
		t.Fatal(err)
	}
	if rel == nil {
		t.Fatal("expected a release")
	}
	if rel.Version != "2.1.0" {
		t.Errorf("version = %q", rel.Version)
	}
	if rel.Verification == nil || rel.Verification.SignatureAlgo != "ed25519ph" {
		t.Fatalf("verification not carried through: %+v", rel.Verification)
	}

	var buf bytes.Buffer
	if err := ep.Download(context.Background(), rel, &buf, nil); err != nil {
		t.Fatal(err)
	}
	if buf.String() != content {
		t.Errorf("downloaded %q, want %q", buf.String(), content)
	}
	digest := sha512.Sum512(buf.Bytes())
	if !bytes.Equal(digest[:], rel.Verification.Digest) {
		t.Error("digest in release does not match downloaded bytes")
	}
	err = ed25519.VerifyWithOptions(pub, digest[:], rel.Verification.Signature, &ed25519.Options{Hash: crypto.SHA512})
	if err != nil {
		t.Errorf("signature from provider does not verify: %v", err)
	}
}

func TestCollectUpdaterArtifactsSkipsExplicitSidecars(t *testing.T) {
	dir := t.TempDir()
	app := writeArtifact(t, dir, "MyApp-2.1.0-darwin-arm64.zip", "payload")
	notes := writeArtifact(t, dir, "notes.md", "## notes")
	pub := writeArtifact(t, dir, "updater.pub", "key material")
	manifest := writeArtifact(t, dir, "manifest.json", "{}")

	// A shell glob (build/*) arrives as explicit file arguments: sidecars,
	// key material and the output manifest must still be filtered out.
	files, err := collectUpdaterArtifacts([]string{app, notes, pub, manifest}, manifest)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 || files[0] != app {
		t.Fatalf("got %v, want just %s", files, app)
	}
}

func TestDecodeUpdaterB64AcceptsUnpadded(t *testing.T) {
	want := []byte("digest-bytes!")
	for _, enc := range []string{
		base64.StdEncoding.EncodeToString(want),
		base64.RawStdEncoding.EncodeToString(want),
	} {
		got, err := decodeUpdaterB64(enc)
		if err != nil {
			t.Fatalf("decode %q: %v", enc, err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("decode %q = %q, want %q", enc, got, want)
		}
	}
	if _, err := decodeUpdaterB64("!!!"); err == nil {
		t.Fatal("expected error for invalid base64")
	}
}

func TestUpdaterVerifyRejectsUnverifiableArtifact(t *testing.T) {
	// An entry with neither digest nor signature must fail verification:
	// reporting an unchecked file as verified would defeat the CI gate.
	dir := t.TempDir()
	writeArtifact(t, dir, "app.zip", "payload")
	writeArtifact(t, dir, "manifest.json", `{
		"schemaVersion": 1,
		"version": "1.0.0",
		"publishedAt": "2026-01-01T00:00:00Z",
		"artifacts": [{"url": "app.zip"}]
	}`)
	err := UpdaterVerify(&UpdaterVerifyOptions{Manifest: filepath.Join(dir, "manifest.json")})
	if err == nil {
		t.Fatal("expected verify to fail for an artifact with neither digest nor signature")
	}
}

func TestUpdaterVerifyECDSAP256(t *testing.T) {
	// Third-party publishers may sign with ecdsa-p256; verify must support
	// every algorithm the runtime verifier does.
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	content := "ecdsa signed payload"
	writeArtifact(t, dir, "app.zip", content)
	digest := sha256.Sum256([]byte(content))
	sigDER, err := ecdsa.SignASN1(rand.Reader, key, digest[:])
	if err != nil {
		t.Fatal(err)
	}
	pubDER, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	pubPath := filepath.Join(dir, "ec.pub")
	if err := os.WriteFile(pubPath, pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER}), 0o644); err != nil {
		t.Fatal(err)
	}
	manifest := map[string]any{
		"schemaVersion": 1,
		"version":       "1.0.0",
		"publishedAt":   "2026-01-01T00:00:00Z",
		"artifacts": []map[string]any{{
			"url":           "app.zip",
			"digestAlgo":    "sha256",
			"digest":        base64.StdEncoding.EncodeToString(digest[:]),
			"signatureAlgo": "ecdsa-p256",
			"signature":     base64.StdEncoding.EncodeToString(sigDER),
		}},
	}
	raw, _ := json.Marshal(manifest)
	writeArtifact(t, dir, "manifest.json", string(raw))
	manifestPath := filepath.Join(dir, "manifest.json")

	if err := UpdaterVerify(&UpdaterVerifyOptions{Manifest: manifestPath, PublicKey: pubPath}); err != nil {
		t.Fatalf("ecdsa-p256 verify failed: %v", err)
	}
	// Tampering must still fail.
	writeArtifact(t, dir, "app.zip", "tampered")
	if err := UpdaterVerify(&UpdaterVerifyOptions{Manifest: manifestPath, PublicKey: pubPath}); err == nil {
		t.Fatal("ecdsa-p256 verify passed on tampered artifact")
	}
}
