package commands

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
	"github.com/wailsapp/wails/v3/internal/windowspackage"
)

func commandSigningFixture(t *testing.T) (string, string) {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	cert, err := os.ReadFile("../codesignature/testdata/modern.pfx")
	if err != nil {
		t.Fatal(err)
	}
	certPath := filepath.Join(root, "signing.pfx")
	if err := os.WriteFile(certPath, cert, 0600); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(root, "package")
	if err := os.Mkdir(dir, 0700); err != nil {
		t.Fatal(err)
	}
	manifest := `<Package xmlns="http://schemas.microsoft.com/appx/manifest/foundation/windows10"><Identity Name="test" Publisher="CN=Wails Portable Signing Test Only" Version="1.0.0.0" ProcessorArchitecture="x64"/></Package>`
	if err := os.WriteFile(filepath.Join(dir, "AppxManifest.xml"), []byte(manifest), 0600); err != nil {
		t.Fatal(err)
	}
	input := filepath.Join(root, "input.msix")
	if err := windowspackage.Create(t.Context(), windowspackage.Options{Directory: dir, Output: input}); err != nil {
		t.Fatal(err)
	}
	return input, certPath
}
func TestWindowsBuiltinSignsWithoutSDKAndPreservesSource(t *testing.T) {
	input, cert := commandSigningFixture(t)
	before, err := os.ReadFile(input)
	if err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(filepath.Dir(input), "signed.msix")
	if err := signWindowsBuiltin(t.Context(), &flags.Sign{Input: input, Output: output, Certificate: cert}, "fixture-only"); err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(input)
	if !bytes.Equal(before, after) {
		t.Fatal("input changed")
	}
	archive, err := zip.OpenReader(output)
	if err != nil {
		t.Fatal(err)
	}
	defer archive.Close()
	for _, file := range archive.File {
		if file.Name == "AppxSignature.p7x" {
			return
		}
	}
	t.Fatal("signature missing")
}
func TestWindowsSigningBadPasswordAndTimestampFailClosed(t *testing.T) {
	input, cert := commandSigningFixture(t)
	before, _ := os.ReadFile(input)
	if err := signWindowsWithPassword(t.Context(), &flags.Sign{Input: input, Certificate: cert}, "secret-invalid-password"); err == nil || strings.Contains(err.Error(), "secret-invalid-password") {
		t.Fatalf("unsafe password error: %v", err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/timestamp-reply")
		_, _ = w.Write([]byte("invalid timestamp"))
	}))
	defer server.Close()
	if err := signWindowsWithPassword(t.Context(), &flags.Sign{Input: input, Certificate: cert, Timestamp: server.URL}, "fixture-only"); err == nil {
		t.Fatal("fake timestamp accepted")
	}
	after, _ := os.ReadFile(input)
	if !bytes.Equal(before, after) {
		t.Fatal("failed signing changed input")
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if err := signWindowsWithPassword(ctx, &flags.Sign{Input: input, Certificate: cert}, "fixture-only"); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation: %v", err)
	}
}
func TestManifestWindowsUsesMergedCredentialsAndTransaction(t *testing.T) {
	input, cert := commandSigningFixture(t)
	before, _ := os.ReadFile(input)
	handler := &manifestHandler{root: filepath.Dir(input), environment: []string{"WINDOWS_TEST_PASSWORD=fixture-only"}}
	spec := pipeline.SignSpec{TargetOS: "windows", Format: "msix", Input: filepath.Base(input), Config: manifest.SigningPlatform{Enabled: true, Certificate: filepath.Base(cert), Credential: "WINDOWS_TEST_PASSWORD"}}
	previous := manifestWindowsSign
	t.Cleanup(func() { manifestWindowsSign = previous })
	manifestWindowsSign = func(ctx context.Context, options *flags.Sign) error {
		if options.Password != "fixture-only" || options.Certificate != cert || options.Input != input || options.Output != input+".signed" {
			t.Fatalf("wrong signing settings: certificate=%q input=%q", options.Certificate, options.Input)
		}
		return signWindowsBuiltin(ctx, options, options.Password)
	}
	if _, err := handler.sign(t.Context(), spec); err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(input)
	if !bytes.Equal(before, after) {
		t.Fatal("manifest changed original")
	}
	if _, err := os.Stat(input + ".signed"); err != nil {
		t.Fatal(err)
	}
	handler.environment = []string{"WINDOWS_TEST_PASSWORD="}
	if _, err := handler.sign(t.Context(), spec); err == nil {
		t.Fatal("missing credential accepted")
	}
}

func TestFindSDKSigningTool(t *testing.T) {
	root := t.TempDir()
	for _, version := range []string{"10.0.22621.0", "10.0.26100.0"} {
		path := filepath.Join(root, version, "x64", "signtool.exe")
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("fixture"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	got, err := findSDKSigningTool(root, "x64")
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(root, "10.0.26100.0", "x64", "signtool.exe"); got != want {
		t.Fatalf("got %s want %s", got, want)
	}
	if _, err := findSDKSigningTool(root, "arm64"); err == nil {
		t.Fatal("wrong host architecture accepted")
	}
}
