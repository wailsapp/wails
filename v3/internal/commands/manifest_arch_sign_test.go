package commands

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func TestArchSignaturePipelineWithIndependentGPGVerification(t *testing.T) {
	gpg, err := exec.LookPath("gpg")
	if err != nil {
		t.Skip("gpg is needed for independent detached signature verification")
	}
	gpgv, err := exec.LookPath("gpgv")
	if err != nil {
		t.Skip("gpgv is needed for independent detached signature verification")
	}
	root := executableFixture(t)
	t.Chdir(root)
	t.Setenv("GOFLAGS", "-buildvcs=false")
	t.Setenv("ARCH_TEST_PASSWORD", "test-only-arch-passphrase")
	privateKey, publicKey := filepath.Join(root, "signing.asc"), filepath.Join(root, "public.asc")
	if err := generatePGPKeyForSetup("Wake signing fixture", "fixture@example.invalid", "test-only-arch-passphrase", privateKey, publicKey); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(root, "wails.yaml")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	data = append(data, []byte("\nlinux:\n  signing:\n    certificate: signing.asc\n    credential: ARCH_TEST_PASSWORD\n")...)
	if err := os.WriteFile(manifestPath, data, 0644); err != nil {
		t.Fatal(err)
	}
	options := manifestRunOptions{Verb: "sign", TargetOS: "linux", TargetArch: "arm64", Formats: []string{"archlinux"}, SkipFrontend: true, SkipBindings: true, Environment: []string{"CGO_ENABLED=0"}}
	if _, err := runManifestPipelineResult(options); err != nil {
		t.Fatal(err)
	}
	packages, err := filepath.Glob(filepath.Join(root, "bin", "*.pkg.tar.zst"))
	if err != nil || len(packages) != 1 {
		t.Fatalf("signed package outputs: %v, %v", packages, err)
	}
	archive := packages[0]
	signature := archive + ".sig"
	original, err := os.ReadFile(archive)
	if err != nil {
		t.Fatal(err)
	}
	keyring := "public.gpg"
	output, err := exec.Command(gpg, "--batch", "--dearmor", "--output", keyring, filepath.Base(publicKey)).CombinedOutput()
	if err != nil {
		t.Fatalf("decode fixture public key: %s, %v", output, err)
	}
	verify := func(path string, valid bool) {
		t.Helper()
		relative, err := filepath.Rel(root, path)
		if err != nil {
			t.Fatal(err)
		}
		output, err := exec.Command(gpgv, "--homedir", ".", "--keyring", "./"+keyring, filepath.ToSlash(filepath.Join("bin", filepath.Base(signature))), filepath.ToSlash(relative)).CombinedOutput()
		if (err == nil) != valid {
			t.Fatalf("GPG verification expected valid=%v: %s, %v", valid, output, err)
		}
		if valid {
			t.Logf("Independent GPG verification: %s", output)
		}
	}
	verify(archive, true)
	tampered := filepath.Join(root, "tampered.pkg.tar.zst")
	if err := os.WriteFile(tampered, append(bytes.Clone(original), 'x'), 0644); err != nil {
		t.Fatal(err)
	}
	verify(tampered, false)
	receipt, err := pipeline.VerifyArtifactReceipt(root, ".wails/artifacts/receipt.json")
	if err != nil {
		t.Fatal(err)
	}
	found := map[string]bool{}
	for _, artifact := range receipt.Artifacts {
		found[artifact.Format] = true
	}
	if !found["archlinux"] || !found["sig"] {
		t.Fatalf("receipt omitted package or signature: %#v", receipt)
	}
	oldSignature, err := os.ReadFile(signature)
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("ARCH_TEST_PASSWORD", "wrong-test-password")
	if _, err := runManifestPipelineResult(options); err == nil {
		t.Fatal("wrong passphrase succeeded")
	}
	for name, expected := range map[string][]byte{archive: original, signature: oldSignature} {
		actual, err := os.ReadFile(name)
		if err != nil || !bytes.Equal(actual, expected) {
			t.Fatalf("failed signing replaced %s: %v", name, err)
		}
	}
	t.Setenv("ARCH_TEST_PASSWORD", "test-only-arch-passphrase")
	if err := os.Remove(archive); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(signature); err != nil {
		t.Fatal(err)
	}
	if _, err := runManifestPipelineResult(options); err != nil {
		t.Fatal(err)
	}
	verify(archive, true)
	if _, err := pipeline.VerifyArtifactReceipt(root, ".wails/artifacts/receipt.json"); err != nil {
		t.Fatal(err)
	}
	_, _, plan, err := resolveManifestPlan(options)
	if err != nil {
		t.Fatal(err)
	}
	for _, node := range plan.Nodes {
		if node.Kind == pipeline.SignArtifact && (node.Cache != pipeline.CacheNever || strings.Contains(node.Output, "signing.asc")) {
			t.Fatalf("private-key/signature caching policy: %#v", node)
		}
	}
	options.TargetArch = "amd64"
	if _, err := runManifestPipelineResult(options); err != nil {
		t.Fatal(err)
	}
	archives, err := filepath.Glob(filepath.Join(root, "bin", "*-amd64.pkg.tar.zst"))
	if err != nil || len(archives) != 1 {
		t.Fatalf("amd64 package: %v, %v", archives, err)
	}
	signature = archives[0] + ".sig"
	verify(archives[0], true)
	if err := Sign(&flags.Sign{Input: archives[0], PGPKey: privateKey, PGPPassword: "test-only-arch-passphrase"}); err != nil {
		t.Fatal(err)
	}
	verify(archives[0], true)
}
