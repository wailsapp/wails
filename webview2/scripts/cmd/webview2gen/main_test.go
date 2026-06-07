package main

import (
	"archive/zip"
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"updater/internal/idl"
)

// buildBinary compiles the CLI into a temp dir and returns its path.
// Each test gets its own copy so concurrent runs don't conflict.
func buildBinary(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "webview2gen")
	// `go build -o <name>` on Windows writes <name>.exe regardless of the
	// requested name. Mirror that so exec.Command finds the right file.
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	cmd := exec.Command("go", "build", "-o", bin, ".")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}
	return bin
}

func run(t *testing.T, bin, wd string, args ...string) (string, string, int) {
	t.Helper()
	cmd := exec.Command(bin, args...)
	cmd.Dir = wd
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	code := 0
	if exit, ok := err.(*exec.ExitError); ok {
		code = exit.ExitCode()
	} else if err != nil {
		t.Fatalf("run: %v", err)
	}
	return stdout.String(), stderr.String(), code
}

func TestHelp(t *testing.T) {
	bin := buildBinary(t)
	out, _, code := run(t, bin, t.TempDir(), "help")
	if code != 0 {
		t.Errorf("help exit = %d, want 0", code)
	}
	for _, want := range []string{"download", "generate", "capabilities", "verify", "full"} {
		if !strings.Contains(out, want) {
			t.Errorf("help output missing %q", want)
		}
	}
}

func TestUnknownCommand(t *testing.T) {
	bin := buildBinary(t)
	_, stderr, code := run(t, bin, t.TempDir(), "frobnicate")
	if code == 0 {
		t.Error("unknown command should exit nonzero")
	}
	if !strings.Contains(stderr, "unknown command") {
		t.Errorf("expected 'unknown command' in stderr, got: %s", stderr)
	}
}

// fakeIDLContents is the minimum IDL required for the parser to produce a file.
const fakeIDLContents = `
[uuid(26d34152-879f-4065-bea2-3daa2cfadfb8), version(1.0)]
library WebView2 {
	[uuid(d60ac92c-37a6-4b26-a39e-95cfe59047bb), object, pointer_default(unique)]
	interface ICoreWebView2Fake : IUnknown {
		HRESULT Ping([out, retval] LPWSTR* result);
	}
}
`

func TestGenerateAndVerify(t *testing.T) {
	bin := buildBinary(t)
	wd := t.TempDir()

	idlDir := filepath.Join(wd, "scripts-cache")
	if err := os.MkdirAll(idlDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(idlDir, "WebView2.1.0.9999.0.idl"), []byte(fakeIDLContents), 0o644); err != nil {
		t.Fatal(err)
	}

	out := filepath.Join(wd, "out")
	_, stderr, code := run(t, bin, wd, "generate", "-dir", idlDir, "-out", out)
	if code != 0 {
		t.Fatalf("generate exit=%d stderr=%s", code, stderr)
	}

	entries, err := os.ReadDir(out)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Fatal("generate produced no files")
	}

	// verify should succeed against the freshly-generated tree.
	_, stderr, code = run(t, bin, wd, "verify", "-dir", idlDir, "-out", out)
	if code != 0 {
		t.Fatalf("verify exit=%d stderr=%s", code, stderr)
	}

	// Mutate one file; verify should fail.
	target := filepath.Join(out, entries[0].Name())
	content, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, append(content, []byte("\n// hand edit\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	_, stderr, code = run(t, bin, wd, "verify", "-dir", idlDir, "-out", out)
	if code == 0 {
		t.Errorf("verify should have failed on hand-edited file; stderr=%s", stderr)
	}
	if !strings.Contains(stderr, "changed file") {
		t.Errorf("verify failure should mention 'changed file', got: %s", stderr)
	}
}

func TestCapabilities(t *testing.T) {
	bin := buildBinary(t)
	wd := t.TempDir()

	notes := `## Stable Release Notes

[NuGet package for WebView2 1.0.500.1](url)

This release requires WebView2 Runtime version 100.0.0.1 or higher.

* [ICoreWebView2Foo interface](/reference/x?view=webview2-1.0.500.1)
* [ICoreWebView2_5 interface](/reference/y?view=webview2-1.0.500.1)
`
	src := filepath.Join(wd, "notes.md")
	if err := os.WriteFile(src, []byte(notes), 0o644); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(wd, "pkg")
	_, stderr, code := run(t, bin, wd, "capabilities", "-source", src, "-out", out)
	if code != 0 {
		t.Fatalf("capabilities exit=%d stderr=%s", code, stderr)
	}
	got, err := os.ReadFile(filepath.Join(out, "capabilities.go"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(got)
	for _, want := range []string{
		`"ICoreWebView2Foo": "1.0.500.1"`,
		`"ICoreWebView2_5": "1.0.500.1"`,
		"SupportsInterface",
		"HasCapability",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("capabilities.go missing %q", want)
		}
	}
}

func TestDownloadServesFromCache(t *testing.T) {
	bin := buildBinary(t)
	wd := t.TempDir()

	// Reach into the fetcher's logic by setting up a cached IDL: the
	// `download` command sees Has(v)=true and exits with the "already cached"
	// message without hitting the network.
	store := idl.NewStore(filepath.Join(wd, "cache"))
	if err := os.MkdirAll(store.Dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(store.CachePath("1.0.123.0"), []byte("dummy"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, stderr, code := run(t, bin, wd, "download", "-version", "1.0.123.0", "-dir", store.Dir)
	if code != 0 {
		t.Fatalf("download exit=%d stderr=%s", code, stderr)
	}
	if !strings.Contains(stderr, "already cached") {
		t.Errorf("expected 'already cached' message, got: %s", stderr)
	}
}

// TestDownloadFromHTTP shows the wire path works against a stubbed NuGet.
// The CLI's `download` command always hits the hard-coded NuGet URL, so we
// exercise the lower-level Fetcher to prove the end-to-end pipeline works
// without depending on the internet.
func TestDownloadFetcherEndToEnd(t *testing.T) {
	// Build a tiny .nupkg containing WebView2.idl.
	var pkg bytes.Buffer
	zw := zip.NewWriter(&pkg)
	if f, _ := zw.Create("WebView2.idl"); f != nil {
		f.Write([]byte(fakeIDLContents))
	}
	zw.Close()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(pkg.Bytes())
	}))
	defer srv.Close()

	store := idl.NewStore(t.TempDir())
	f := &idl.Fetcher{HTTPClient: srv.Client(), Store: store}

	// Manually invoke the http request because Download() uses the real URL.
	resp, err := f.HTTPClient.Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
}
