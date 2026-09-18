package updater_test

import (
	"archive/zip"
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
	"encoding/asn1"
	"errors"
	"io"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/wailsapp/wails/v3/pkg/updater"
)

func removeAllOS(p string) error { return os.RemoveAll(p) }

// fakeHost captures every emit so tests can assert event sequence and
// payloads, dispatches OnEvent callbacks to fired emits, and records
// OpenWindow calls.
type fakeHost struct {
	mu        sync.Mutex
	events    []fakeEvent
	listeners map[string][]func(any)

	openCalls []updater.WindowOptions
	window    *fakeWindow
	quits     int
}

type fakeEvent struct {
	Name string
	Data any
}

func (f *fakeHost) Emit(name string, data ...any) bool {
	f.mu.Lock()
	var d any
	if len(data) == 1 {
		d = data[0]
	} else if len(data) > 1 {
		d = data
	}
	f.events = append(f.events, fakeEvent{Name: name, Data: d})
	cbs := append([]func(any){}, f.listeners[name]...)
	f.mu.Unlock()
	for _, cb := range cbs {
		cb(d)
	}
	return false
}

func (f *fakeHost) OnEvent(name string, cb func(any)) func() {
	f.mu.Lock()
	if f.listeners == nil {
		f.listeners = map[string][]func(any){}
	}
	f.listeners[name] = append(f.listeners[name], cb)
	idx := len(f.listeners[name]) - 1
	f.mu.Unlock()
	return func() {
		f.mu.Lock()
		defer f.mu.Unlock()
		if cbs := f.listeners[name]; idx < len(cbs) {
			f.listeners[name] = append(cbs[:idx], cbs[idx+1:]...)
		}
	}
}

func (f *fakeHost) OpenWindow(opts updater.WindowOptions) updater.WindowHandle {
	f.mu.Lock()
	f.openCalls = append(f.openCalls, opts)
	w := &fakeWindow{}
	f.window = w
	f.mu.Unlock()
	return w
}

func (f *fakeHost) Quit() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.quits++
}

// fakeWindow records every interaction the Updater performs on a window.
type fakeWindow struct {
	mu     sync.Mutex
	closed bool
	shown  int
	events []fakeEvent
}

func (w *fakeWindow) EmitEvent(name string, data ...any) bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	var d any
	if len(data) == 1 {
		d = data[0]
	}
	w.events = append(w.events, fakeEvent{Name: name, Data: d})
	return false
}
func (w *fakeWindow) Show() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.shown++
}
func (w *fakeWindow) Close() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.closed = true
}

func (f *fakeHost) names() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]string, len(f.events))
	for i, e := range f.events {
		out[i] = e.Name
	}
	return out
}

func (f *fakeHost) payloadFor(name string) any {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, e := range f.events {
		if e.Name == name {
			return e.Data
		}
	}
	return nil
}

// fakeProvider is a controllable in-memory Provider.
type fakeProvider struct {
	name      string
	rel       *updater.Release
	checkErr  error
	body      []byte
	dlErr     error
	calls     int
	downloads int
}

func (f *fakeProvider) Name() string { return f.name }
func (f *fakeProvider) Check(ctx context.Context, _ updater.CheckRequest) (*updater.Release, error) {
	f.calls++
	return f.rel, f.checkErr
}
func (f *fakeProvider) Download(ctx context.Context, _ *updater.Release, dst io.Writer, onProgress func(int64, int64)) error {
	f.downloads++
	if f.dlErr != nil {
		return f.dlErr
	}
	if _, err := io.Copy(dst, bytes.NewReader(f.body)); err != nil {
		return err
	}
	if onProgress != nil {
		onProgress(int64(len(f.body)), int64(len(f.body)))
	}
	return nil
}

// --- Init / validation ---

func TestInit_RequiresCurrentVersion(t *testing.T) {
	u := updater.New(&fakeHost{})
	if err := u.Init(updater.Config{Providers: []updater.Provider{&fakeProvider{name: "f"}}}); err == nil {
		t.Fatal("expected error for missing CurrentVersion")
	}
}

func TestInit_RequiresProviders(t *testing.T) {
	u := updater.New(&fakeHost{})
	if err := u.Init(updater.Config{CurrentVersion: "1.0.0"}); err == nil {
		t.Fatal("expected error for empty Providers")
	}
}

func TestInit_RejectsNilProvider(t *testing.T) {
	u := updater.New(&fakeHost{})
	err := u.Init(updater.Config{
		CurrentVersion: "1.0.0",
		Providers:      []updater.Provider{&fakeProvider{name: "f"}, nil},
	})
	if err == nil || !strings.Contains(err.Error(), "nil entry") {
		t.Fatalf("expected nil-entry error, got %v", err)
	}
}

// Config.CheckInterval starts a background poll loop that invokes
// CheckAndInstall on each tick. Verify that an Init with a short interval
// produces ticks against the provider and that StopPeriodicCheck cleanly
// halts the loop.
func TestInit_CheckInterval_TicksProviderAndStops(t *testing.T) {
	host := &fakeHost{}
	rel := &updater.Release{
		Version:  "2.0.0",
		Artifact: updater.Artifact{Filename: "app.bin"},
		Verification: &updater.Verification{DigestAlgo: "sha256", Digest: sha256.New().Sum(nil)},
	}
	// Pre-compute matching digest so Check finds + DownloadAndInstall succeeds.
	body := []byte("payload")
	d := sha256.Sum256(body)
	rel.Verification.Digest = d[:]
	rel.Artifact.Size = int64(len(body))
	p := &fakeProvider{name: "p", rel: rel, body: body}

	u := updater.New(host)
	if err := u.Init(updater.Config{
		CurrentVersion: "1.0.0",
		Providers:      []updater.Provider{p},
		CheckInterval:  50 * time.Millisecond,
		Window:         updater.WindowNone,
	}); err != nil {
		t.Fatal(err)
	}
	defer u.StopPeriodicCheck()

	// Allow at least two ticks. CheckAndInstall is a no-op after the first
	// success (state stays in StateReady which periodic loop ignores) so we
	// only assert at least one Check landed.
	time.Sleep(200 * time.Millisecond)
	u.StopPeriodicCheck()

	if p.calls == 0 {
		t.Fatalf("periodic check never invoked the provider")
	}
}

func TestInit_RejectsDoubleConfigure(t *testing.T) {
	u := updater.New(&fakeHost{})
	cfg := updater.Config{CurrentVersion: "1.0.0", Providers: []updater.Provider{&fakeProvider{name: "f"}}}
	if err := u.Init(cfg); err != nil {
		t.Fatal(err)
	}
	if err := u.Init(cfg); !errors.Is(err, updater.ErrAlreadyConfigured) {
		t.Fatalf("expected ErrAlreadyConfigured, got %v", err)
	}
}

func TestCheck_NotConfigured(t *testing.T) {
	u := updater.New(&fakeHost{})
	if _, err := u.Check(context.Background()); !errors.Is(err, updater.ErrNotConfigured) {
		t.Fatalf("expected ErrNotConfigured, got %v", err)
	}
}

// --- Check + fallback ---

func TestCheck_SingleProvider_UpdateAvailable(t *testing.T) {
	host := &fakeHost{}
	rel := &updater.Release{Version: "2.0.0", Artifact: updater.Artifact{Filename: "app.dmg", Size: 5}}
	p := &fakeProvider{name: "primary", rel: rel}

	u := newConfigured(t, host, p)
	got, err := u.Check(context.Background())
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if got == nil || got.Version != "2.0.0" {
		t.Fatalf("got %+v", got)
	}
	if got.Provider != "primary" {
		t.Errorf("provider tag: want primary, got %q", got.Provider)
	}
	if u.State() != updater.StateAvailable {
		t.Errorf("state: %s", u.State())
	}
	want := []string{updater.EventCheckStarted, updater.EventUpdateAvailable}
	assertEventNames(t, host.names(), want)
}

func TestCheck_UpToDate_StopsFallbackChain(t *testing.T) {
	host := &fakeHost{}
	primary := &fakeProvider{name: "primary"} // returns (nil,nil) → up-to-date
	secondary := &fakeProvider{name: "secondary", rel: &updater.Release{Version: "9.9.9"}}

	u := newConfigured(t, host, primary, secondary)
	got, err := u.Check(context.Background())
	if err != nil || got != nil {
		t.Fatalf("want (nil,nil), got %+v %v", got, err)
	}
	if secondary.calls != 0 {
		t.Errorf("secondary should NOT have been called after primary said up-to-date")
	}
	if u.State() != updater.StateUpToDate {
		t.Errorf("state: %s", u.State())
	}
	assertEventNames(t, host.names(), []string{updater.EventCheckStarted, updater.EventNoUpdate})
}

func TestCheck_FallbackOnError(t *testing.T) {
	host := &fakeHost{}
	primary := &fakeProvider{name: "primary", checkErr: errors.New("network down")}
	rel := &updater.Release{Version: "2.0.0", Artifact: updater.Artifact{Filename: "app.dmg"}}
	secondary := &fakeProvider{name: "secondary", rel: rel}

	u := newConfigured(t, host, primary, secondary)
	got, err := u.Check(context.Background())
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if got.Provider != "secondary" {
		t.Fatalf("expected fallback to secondary, got %q", got.Provider)
	}
}

func TestCheck_AllProvidersError_ReturnsWrappedError(t *testing.T) {
	host := &fakeHost{}
	p1 := &fakeProvider{name: "primary", checkErr: errors.New("dns")}
	p2 := &fakeProvider{name: "secondary", checkErr: errors.New("401")}
	u := newConfigured(t, host, p1, p2)

	_, err := u.Check(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "primary") || !strings.Contains(err.Error(), "secondary") {
		t.Errorf("error should mention both providers: %v", err)
	}
	if u.State() != updater.StateError {
		t.Errorf("state: %s", u.State())
	}
	assertEventNames(t, host.names(), []string{updater.EventCheckStarted, updater.EventError})
}

// --- Download + verify ---

func TestDownloadAndInstall_NoPendingRelease(t *testing.T) {
	u := newConfigured(t, &fakeHost{}, &fakeProvider{name: "p"})
	if err := u.DownloadAndInstall(context.Background()); !errors.Is(err, updater.ErrNoPendingRelease) {
		t.Fatalf("expected ErrNoPendingRelease, got %v", err)
	}
}

func TestDownloadAndInstall_HappyPath_NoVerification(t *testing.T) {
	host := &fakeHost{}
	body := []byte("hello-end-to-end")
	rel := &updater.Release{
		Version:  "2.0.0",
		Artifact: updater.Artifact{Filename: "app.bin", Size: int64(len(body))},
	}
	p := &fakeProvider{name: "p", rel: rel, body: body}
	u := newConfigured(t, host, p)

	if _, err := u.Check(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := u.DownloadAndInstall(context.Background()); err != nil {
		t.Fatalf("DownloadAndInstall: %v", err)
	}
	if u.State() != updater.StateReady {
		t.Errorf("state: %s", u.State())
	}
	wantOrder := []string{
		updater.EventCheckStarted,
		updater.EventUpdateAvailable,
		updater.EventDownloadStarted,
		updater.EventDownloadProgress, // at least one
		updater.EventDownloadComplete,
		updater.EventVerifying,
		updater.EventInstalling,
		updater.EventUpdateReady,
	}
	assertEventSubsequence(t, host.names(), wantOrder)
}

// On any failure between download and ready (verify mismatch, unknown
// digest algo, etc.) the temp staging directory created under
// os.TempDir/wails-update-* must be removed; otherwise repeated update
// attempts accumulate orphan directories.
func TestDownloadAndInstall_FailedFlow_RemovesStagingDir(t *testing.T) {
	host := &fakeHost{}
	body := []byte("real-bytes")
	wrong := sha256.Sum256([]byte("different-bytes"))
	rel := &updater.Release{
		Version:  "2.0.0",
		Artifact: updater.Artifact{Filename: "app.bin", Size: int64(len(body))},
		Verification: &updater.Verification{
			DigestAlgo: "sha256",
			Digest:     wrong[:],
		},
	}
	p := &fakeProvider{name: "p", rel: rel, body: body}
	u := newConfigured(t, host, p)

	before := countStagingDirs(t)

	if _, err := u.Check(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := u.DownloadAndInstall(context.Background()); err == nil {
		t.Fatal("expected error from digest mismatch")
	}

	after := countStagingDirs(t)
	if after > before {
		t.Errorf("staging dir leaked: %d → %d wails-update-* directories under %s", before, after, os.TempDir())
	}
}

// countStagingDirs returns the number of `wails-update-*` directories under
// os.TempDir. Used to detect leaks across the download/install flow without
// being sensitive to absolute paths.
func countStagingDirs(t *testing.T) int {
	t.Helper()
	entries, err := os.ReadDir(os.TempDir())
	if err != nil {
		t.Fatalf("read tempdir: %v", err)
	}
	n := 0
	for _, e := range entries {
		if e.IsDir() && strings.HasPrefix(e.Name(), "wails-update-") {
			n++
		}
	}
	return n
}

func TestDownloadAndInstall_DigestMismatch_FailsClosed(t *testing.T) {
	host := &fakeHost{}
	body := []byte("real-bytes")
	wrongDigest := sha256.Sum256([]byte("different-bytes"))
	rel := &updater.Release{
		Version:  "2.0.0",
		Artifact: updater.Artifact{Filename: "app.bin", Size: int64(len(body))},
		Verification: &updater.Verification{
			DigestAlgo: "sha256",
			Digest:     wrongDigest[:],
		},
	}
	p := &fakeProvider{name: "p", rel: rel, body: body}
	u := newConfigured(t, host, p)

	if _, err := u.Check(context.Background()); err != nil {
		t.Fatal(err)
	}
	err := u.DownloadAndInstall(context.Background())
	if err == nil {
		t.Fatal("expected verification error")
	}
	if !strings.Contains(err.Error(), "digest mismatch") {
		t.Errorf("wrong error: %v", err)
	}
	if u.State() != updater.StateError {
		t.Errorf("state: %s", u.State())
	}
	if u.DownloadedPath() != "" {
		t.Errorf("downloaded path should be empty after failed verify: %q", u.DownloadedPath())
	}
	// Error event must mention the verify stage and provider.
	payload := host.payloadFor(updater.EventError)
	ei, ok := payload.(updater.ErrorInfo)
	if !ok {
		t.Fatalf("payload type: %T", payload)
	}
	if ei.Stage != updater.StageVerify {
		t.Errorf("error stage: %s", ei.Stage)
	}
	if ei.Provider != "p" {
		t.Errorf("error provider: %s", ei.Provider)
	}
}

// Restart must (a) refuse when nothing is staged, (b) on successful spawn
// dispatch Host.Quit so the helper's "wait for parent to exit" step
// completes. Previously Restart spawned the helper and returned, leaving
// the caller responsible for calling app.Quit themselves — which everyone
// missed, so "Restart" looked like a no-op.
func TestRestart_QuitsAfterSpawn(t *testing.T) {
	// Substitute a benign command in place of the real self re-exec so
	// Start() succeeds without actually launching another test binary in
	// helper mode.
	t.Cleanup(updater.SetSelfExecutableForTest(os.Executable))
	t.Cleanup(updater.SetNewDetachedCommandForTest(func(path string) *exec.Cmd {
		return exec.Command(path, "-test.run=^$")
	}))

	host := &fakeHost{}
	body := []byte("payload")
	rel := &updater.Release{
		Version:  "2.0.0",
		Artifact: updater.Artifact{Filename: "app.bin", Size: int64(len(body))},
	}
	p := &fakeProvider{name: "p", rel: rel, body: body}
	u := newConfigured(t, host, p)

	// Not ready yet: Restart should refuse.
	if err := u.Restart(context.Background()); !errors.Is(err, updater.ErrNotReady) {
		t.Fatalf("Restart pre-install: want ErrNotReady, got %v", err)
	}
	if host.quits != 0 {
		t.Fatalf("Quit dispatched too early: %d", host.quits)
	}

	// Stage a release through the normal flow.
	if _, err := u.Check(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := u.DownloadAndInstall(context.Background()); err != nil {
		t.Fatalf("DownloadAndInstall: %v", err)
	}

	if err := u.Restart(context.Background()); err != nil {
		t.Fatalf("Restart: %v", err)
	}
	host.mu.Lock()
	defer host.mu.Unlock()
	if host.quits != 1 {
		t.Errorf("Quit dispatch count: got %d, want 1", host.quits)
	}
}

func TestDownloadAndInstall_DigestMatch_Succeeds(t *testing.T) {
	host := &fakeHost{}
	body := []byte("real-bytes")
	digest := sha256.Sum256(body)
	rel := &updater.Release{
		Version:  "2.0.0",
		Artifact: updater.Artifact{Filename: "app.bin", Size: int64(len(body))},
		Verification: &updater.Verification{
			DigestAlgo: "sha256",
			Digest:     digest[:],
		},
	}
	p := &fakeProvider{name: "p", rel: rel, body: body}
	u := newConfigured(t, host, p)

	if _, err := u.Check(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := u.DownloadAndInstall(context.Background()); err != nil {
		t.Fatalf("DownloadAndInstall: %v", err)
	}
	if u.State() != updater.StateReady {
		t.Errorf("state: %s", u.State())
	}
}

// End-to-end check that .zip artifacts get unpacked between verify and ready.
// macOS distributes .app bundles inside .zip archives; without extraction the
// helper would replace a directory target with a single zip file. This test
// downloads a zip carrying a single top-level .app directory and asserts that
// u.DownloadedPath() points at the extracted directory.
func TestDownloadAndInstall_ZipBundle_Extracted(t *testing.T) {
	host := &fakeHost{}

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	mkEntry := func(name string, body []byte, mode os.FileMode, dir bool) {
		hdr := &zip.FileHeader{Name: name, Method: zip.Deflate}
		if dir {
			hdr.SetMode(mode | os.ModeDir)
		} else {
			hdr.SetMode(mode)
		}
		w, err := zw.CreateHeader(hdr)
		if err != nil {
			t.Fatal(err)
		}
		if !dir {
			if _, err := w.Write(body); err != nil {
				t.Fatal(err)
			}
		}
	}
	mkEntry("MyApp.app/", nil, 0o755, true)
	mkEntry("MyApp.app/Contents/", nil, 0o755, true)
	mkEntry("MyApp.app/Contents/Info.plist", []byte("<plist/>"), 0o644, false)
	mkEntry("MyApp.app/Contents/MacOS/", nil, 0o755, true)
	mkEntry("MyApp.app/Contents/MacOS/MyApp", []byte("ELF"), 0o755, false)
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	body := buf.Bytes()
	digest := sha256.Sum256(body)

	rel := &updater.Release{
		Version:  "2.0.0",
		Artifact: updater.Artifact{Filename: "MyApp-darwin.zip", Size: int64(len(body))},
		Verification: &updater.Verification{
			DigestAlgo: "sha256",
			Digest:     digest[:],
		},
	}
	p := &fakeProvider{name: "p", rel: rel, body: body}
	u := newConfigured(t, host, p)

	if _, err := u.Check(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := u.DownloadAndInstall(context.Background()); err != nil {
		t.Fatalf("DownloadAndInstall: %v", err)
	}
	staged := u.DownloadedPath()
	info, err := os.Stat(staged)
	if err != nil {
		t.Fatalf("stat staged: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("staged path should be a directory after extraction, got %v", info.Mode())
	}
	if filepath.Base(staged) != "MyApp.app" {
		t.Errorf("staged base: got %q, want MyApp.app", filepath.Base(staged))
	}
	if got := readFileT(t, filepath.Join(staged, "Contents", "Info.plist")); string(got) != "<plist/>" {
		t.Errorf("Info.plist after extract: %q", got)
	}
}

func readFileT(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestDownloadAndInstall_Ed25519Signature_Verifies(t *testing.T) {
	host := &fakeHost{}
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	body := []byte("signed-payload")
	// Ed25519 verifier signs the digest (here, SHA-256 of the body since
	// DigestAlgo is sha256). The signer signs the same digest.
	digest := sha256.Sum256(body)
	sig := ed25519.Sign(priv, digest[:])

	rel := &updater.Release{
		Version:  "2.0.0",
		Artifact: updater.Artifact{Filename: "app.bin", Size: int64(len(body))},
		Verification: &updater.Verification{
			DigestAlgo:    "sha256",
			SignatureAlgo: "ed25519",
			Signature:     sig,
		},
	}
	p := &fakeProvider{name: "p", rel: rel, body: body}

	cfg := updater.Config{
		CurrentVersion: "1.0.0",
		Providers:      []updater.Provider{p},
		PublicKey:      []byte(pub),
	}
	u := updater.New(host)
	if err := u.Init(cfg); err != nil {
		t.Fatal(err)
	}
	if _, err := u.Check(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := u.DownloadAndInstall(context.Background()); err != nil {
		t.Fatalf("DownloadAndInstall: %v", err)
	}
	if u.State() != updater.StateReady {
		t.Errorf("state: %s", u.State())
	}
}

func TestDownloadAndInstall_Ed25519phSignature_Verifies(t *testing.T) {
	host := &fakeHost{}
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	body := []byte("signed-payload-ph")
	digest := sha512.Sum512(body)
	sig, err := priv.Sign(rand.Reader, digest[:], &ed25519.Options{Hash: crypto.SHA512})
	if err != nil {
		t.Fatal(err)
	}

	pkix, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}

	rel := &updater.Release{
		Version:  "2.0.0",
		Artifact: updater.Artifact{Filename: "app.bin", Size: int64(len(body))},
		Verification: &updater.Verification{
			DigestAlgo:    "sha512",
			SignatureAlgo: "ed25519ph",
			Signature:     sig,
		},
	}
	p := &fakeProvider{name: "p", rel: rel, body: body}
	u := newConfiguredWithKey(t, host, pkix, p)

	if _, err := u.Check(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := u.DownloadAndInstall(context.Background()); err != nil {
		t.Fatalf("DownloadAndInstall: %v", err)
	}
}

func TestDownloadAndInstall_ECDSAP256Signature_Verifies(t *testing.T) {
	host := &fakeHost{}
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	body := []byte("ecdsa-payload")
	digest := sha256.Sum256(body)
	r, s, err := ecdsa.Sign(rand.Reader, priv, digest[:])
	if err != nil {
		t.Fatal(err)
	}
	sig, err := asn1.Marshal(struct{ R, S *big.Int }{r, s})
	if err != nil {
		t.Fatal(err)
	}
	pkix, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		t.Fatal(err)
	}

	rel := &updater.Release{
		Version:  "2.0.0",
		Artifact: updater.Artifact{Filename: "app.bin", Size: int64(len(body))},
		Verification: &updater.Verification{
			DigestAlgo:    "sha256",
			SignatureAlgo: "ecdsa-p256",
			Signature:     sig,
		},
	}
	p := &fakeProvider{name: "p", rel: rel, body: body}
	u := newConfiguredWithKey(t, host, pkix, p)

	if _, err := u.Check(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := u.DownloadAndInstall(context.Background()); err != nil {
		t.Fatalf("DownloadAndInstall: %v", err)
	}
}

// ASN.1 DER signatures with trailing bytes after the valid SEQUENCE must be
// rejected — accepting them risks signature malleability, and conforming
// signers never emit them.
func TestDownloadAndInstall_ECDSA_TrailingBytes_Rejected(t *testing.T) {
	host := &fakeHost{}
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	body := []byte("ecdsa-payload")
	digest := sha256.Sum256(body)
	r, s, err := ecdsa.Sign(rand.Reader, priv, digest[:])
	if err != nil {
		t.Fatal(err)
	}
	sig, err := asn1.Marshal(struct{ R, S *big.Int }{r, s})
	if err != nil {
		t.Fatal(err)
	}
	// Append junk after the valid SEQUENCE.
	sig = append(sig, 0xDE, 0xAD, 0xBE, 0xEF)

	pkix, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	rel := &updater.Release{
		Version:  "2.0.0",
		Artifact: updater.Artifact{Filename: "app.bin", Size: int64(len(body))},
		Verification: &updater.Verification{
			DigestAlgo:    "sha256",
			SignatureAlgo: "ecdsa-p256",
			Signature:     sig,
		},
	}
	p := &fakeProvider{name: "p", rel: rel, body: body}
	u := newConfiguredWithKey(t, host, pkix, p)

	if _, err := u.Check(context.Background()); err != nil {
		t.Fatal(err)
	}
	err = u.DownloadAndInstall(context.Background())
	if err == nil || !strings.Contains(err.Error(), "trailing") {
		t.Fatalf("expected trailing-data error, got %v", err)
	}
}

func TestDownloadAndInstall_SignatureWithoutPublicKey_Fails(t *testing.T) {
	host := &fakeHost{}
	body := []byte("payload")
	rel := &updater.Release{
		Version:  "2.0.0",
		Artifact: updater.Artifact{Filename: "app.bin", Size: int64(len(body))},
		Verification: &updater.Verification{
			DigestAlgo:    "sha256",
			SignatureAlgo: "ed25519",
			Signature:     []byte("doesn't matter — verify aborts before this"),
		},
	}
	p := &fakeProvider{name: "p", rel: rel, body: body}
	u := newConfigured(t, host, p)

	if _, err := u.Check(context.Background()); err != nil {
		t.Fatal(err)
	}
	err := u.DownloadAndInstall(context.Background())
	if err == nil || !strings.Contains(err.Error(), "requires a public key") {
		t.Fatalf("expected public-key-missing error, got %v", err)
	}
}

// A release whose signature verifies against an attacker-controlled key must
// be rejected when Config.PublicKey is set to a different key. The trust root
// is pinned at Init time; the release source has no say in which key
// authenticates it. Regression test: an earlier API allowed Verification to
// carry its own PublicKey that took precedence over Config.PublicKey, which
// would let a compromised release source self-attest with any key it liked.
func TestDownloadAndInstall_Signature_RejectsAttackerKey(t *testing.T) {
	host := &fakeHost{}

	// Attacker generates a key and signs the payload with it.
	_, attackerPriv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	body := []byte("payload")
	digest := sha256.Sum256(body)
	attackerSig := ed25519.Sign(attackerPriv, digest[:])

	// The user has pinned a different key as the trust root.
	pinnedPub, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	pinnedPKIX, err := x509.MarshalPKIXPublicKey(pinnedPub)
	if err != nil {
		t.Fatal(err)
	}

	rel := &updater.Release{
		Version:  "2.0.0",
		Artifact: updater.Artifact{Filename: "app.bin", Size: int64(len(body))},
		Verification: &updater.Verification{
			DigestAlgo:    "sha256",
			SignatureAlgo: "ed25519",
			Signature:     attackerSig,
		},
	}
	p := &fakeProvider{name: "p", rel: rel, body: body}
	u := newConfiguredWithKey(t, host, pinnedPKIX, p)

	if _, err := u.Check(context.Background()); err != nil {
		t.Fatal(err)
	}
	err = u.DownloadAndInstall(context.Background())
	if err == nil || !strings.Contains(err.Error(), "did not verify") {
		t.Fatalf("expected signature-verification failure under pinned key, got %v", err)
	}
}

// --- CheckAndInstall convenience ---

func TestCheckAndInstall_UpToDate_NoOp(t *testing.T) {
	host := &fakeHost{}
	p := &fakeProvider{name: "p"} // returns (nil,nil)
	u := newConfigured(t, host, p)

	if err := u.CheckAndInstall(context.Background()); err != nil {
		t.Fatalf("CheckAndInstall: %v", err)
	}
	if p.downloads != 0 {
		t.Errorf("download should not have been attempted")
	}
}

func TestCheckAndInstall_HappyPath(t *testing.T) {
	host := &fakeHost{}
	body := []byte("payload")
	digest := sha256.Sum256(body)
	rel := &updater.Release{
		Version:      "2.0.0",
		Artifact:     updater.Artifact{Filename: "app.bin", Size: int64(len(body))},
		Verification: &updater.Verification{DigestAlgo: "sha256", Digest: digest[:]},
	}
	p := &fakeProvider{name: "p", rel: rel, body: body}
	u := newConfigured(t, host, p)

	if err := u.CheckAndInstall(context.Background()); err != nil {
		t.Fatalf("CheckAndInstall: %v", err)
	}
	if u.State() != updater.StateReady {
		t.Errorf("state: %s", u.State())
	}
}

// --- helpers ---

func newConfigured(t *testing.T, host updater.Host, providers ...updater.Provider) *updater.Updater {
	t.Helper()
	return newConfiguredWithKey(t, host, nil, providers...)
}

// newConfiguredWithKey is the same as newConfigured but also sets
// Config.PublicKey — used by tests that verify signed releases.
func newConfiguredWithKey(t *testing.T, host updater.Host, publicKey []byte, providers ...updater.Provider) *updater.Updater {
	t.Helper()
	u := updater.New(host)
	if err := u.Init(updater.Config{
		CurrentVersion: "1.0.0",
		Providers:      providers,
		PublicKey:      publicKey,
	}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if p := u.DownloadedPath(); p != "" {
			// Walk up to the wails-update-* temp dir and remove it whole.
			dir := p
			for dir != "" && dir != "/" && dir != "." {
				if strings.Contains(dir, "wails-update-") {
					_ = osRemoveAll(dir)
					return
				}
				dir = parentOf(dir)
			}
		}
	})
	return u
}

func assertEventNames(t *testing.T, got, want []string) {
	t.Helper()
	// Allow EventDownloadProgress to appear N times in a row.
	if !slicesEqualIgnoringProgressRepeats(got, want) {
		t.Errorf("event sequence mismatch:\n  got  %v\n  want %v", got, want)
	}
}

func assertEventSubsequence(t *testing.T, got, mustAppearInOrder []string) {
	t.Helper()
	idx := 0
	for _, g := range got {
		if idx < len(mustAppearInOrder) && g == mustAppearInOrder[idx] {
			idx++
		}
	}
	if idx != len(mustAppearInOrder) {
		t.Errorf("expected subsequence not present:\n  got      %v\n  required %v\n  matched up to %d/%d",
			got, mustAppearInOrder, idx, len(mustAppearInOrder))
	}
}

func slicesEqualIgnoringProgressRepeats(a, b []string) bool {
	// Compact runs of EventDownloadProgress in a to a single occurrence.
	compact := func(s []string) []string {
		out := make([]string, 0, len(s))
		for _, x := range s {
			if len(out) > 0 && out[len(out)-1] == x && x == updater.EventDownloadProgress {
				continue
			}
			out = append(out, x)
		}
		return out
	}
	ac := compact(a)
	bc := compact(b)
	if len(ac) != len(bc) {
		return false
	}
	for i := range ac {
		if ac[i] != bc[i] {
			return false
		}
	}
	return true
}

func osRemoveAll(p string) error { return removeAllOS(p) }
func parentOf(p string) string  { return filepath.Dir(p) }
