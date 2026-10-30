package updater

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// These tests run the real codesign against bundles built around a real
// Mach-O. A text file standing in for the executable never meets the cases
// that matter: on arm64 the Go linker signs every binary it produces, so a
// bundle assembled without a signing step is not "unsigned" to codesign at
// all — it is a signed executable whose signature promises sealed resources
// the bundle does not have.

// requireCodesign skips unless the real codesign can be run.
func requireCodesign(t *testing.T) string {
	t.Helper()
	if runtime.GOOS != "darwin" {
		t.Skip("codesign verification only runs on darwin")
	}
	tool, err := exec.LookPath("codesign")
	if err != nil {
		t.Skip("codesign unavailable")
	}
	return tool
}

const realBundlePlist = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleExecutable</key>
	<string>app</string>
	<key>CFBundleIdentifier</key>
	<string>com.example.updater-test</string>
	<key>CFBundlePackageType</key>
	<string>APPL</string>
</dict>
</plist>
`

// makeRealBundle builds a bundle at root whose main executable is a copy of
// the running test binary, plus one resource file for a signature to seal.
// Nothing signs it beyond what the linker already did.
func makeRealBundle(t *testing.T, root string) {
	t.Helper()
	self, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	src, err := os.Open(self)
	if err != nil {
		t.Fatal(err)
	}
	defer src.Close()

	exe := filepath.Join(root, "Contents", "MacOS", "app")
	if err := os.MkdirAll(filepath.Dir(exe), 0o755); err != nil {
		t.Fatal(err)
	}
	dst, err := os.OpenFile(exe, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(dst, src); err != nil {
		t.Fatal(err)
	}
	if err := dst.Close(); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(root, "Contents", "Info.plist"), []byte(realBundlePlist))
	writeFile(t, filepath.Join(root, "Contents", "Resources", "data.txt"), []byte("sealed"))
}

// adhocSign seals root with an ad-hoc signature, as a default build does.
func adhocSign(t *testing.T, tool, root string) {
	t.Helper()
	if out, err := exec.Command(tool, "--force", "--deep", "-s", "-", root).CombinedOutput(); err != nil {
		t.Fatalf("codesign -s -: %v: %s", err, out)
	}
}

// tamper modifies a sealed resource after signing.
func tamper(t *testing.T, root string) {
	t.Helper()
	writeFile(t, filepath.Join(root, "Contents", "Resources", "data.txt"), []byte("modified"))
}

// A bundle around a linker-signed binary, with no signing step of its own,
// fails codesign --verify with "code has no resources but signature indicates
// they must be present". That is the shape of every cross-compiled or
// hand-assembled arm64 bundle, and it is unverifiable, not damaged.
func TestVerifySignature_LinkerSignedBundleIsNotAFailure(t *testing.T) {
	tool := requireCodesign(t)
	p := filepath.Join(t.TempDir(), "Linker.app")
	makeRealBundle(t, p)

	out, _ := exec.Command(tool, "-dv", p).CombinedOutput()
	t.Logf("codesign -dv: %s", strings.TrimSpace(string(out)))

	note, err := verifySignature(p, "")
	if err != nil {
		t.Fatalf("a linker-signed bundle must not be treated as corrupt: %v", err)
	}
	t.Logf("note: %q", note)
}

// Tolerating the missing seal must not tolerate a damaged executable: the
// linker's signature still hashes every page of the binary, and codesign
// reports a modified page before it reports the missing resources.
func TestVerifySignature_CorruptLinkerSignedExecutableRejected(t *testing.T) {
	requireCodesign(t)
	p := filepath.Join(t.TempDir(), "Corrupt.app")
	makeRealBundle(t, p)
	exe := filepath.Join(p, "Contents", "MacOS", "app")
	f, err := os.OpenFile(exe, os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteAt([]byte("corrupt"), 64*1024); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	if _, err := verifySignature(p, ""); err == nil {
		t.Fatal("a linker-signed executable modified after linking must be rejected")
	}
}

// End to end: an application shipped as linker-signed bundles must be able to
// update from one to the next. A rejection here is a rollback the next check
// repeats forever.
func TestRunHelperSwap_LinkerSignedBundleUpdates(t *testing.T) {
	requireCodesign(t)
	dir := t.TempDir()
	target := filepath.Join(dir, "App.app")
	newPath := filepath.Join(dir, "staging", "App.app")
	makeRealBundle(t, target)
	makeRealBundle(t, newPath)

	l := &fakeLauncher{}
	code := runHelperSwap(target, newPath, 0, filepath.Join(dir, "log"), instantWaiter, l)
	if code != 0 {
		log, _ := os.ReadFile(filepath.Join(dir, "log"))
		t.Fatalf("code = %d, want 0; log:\n%s", code, log)
	}
}

// An ad-hoc signature seals the bundle's resources, so it is a real integrity
// check and damage after signing must still be caught.
func TestVerifySignature_TamperedAdhocSignedBundleRejected(t *testing.T) {
	tool := requireCodesign(t)
	dir := t.TempDir()
	p := filepath.Join(dir, "Tampered.app")
	makeRealBundle(t, p)
	adhocSign(t, tool, p)
	tamper(t, p)

	intact := filepath.Join(dir, "Intact.app")
	makeRealBundle(t, intact)
	adhocSign(t, tool, intact)

	_, err := verifySignature(p, intact)
	if err == nil {
		t.Fatal("a bundle modified after signing must be rejected")
	}
	if !strings.Contains(err.Error(), "a sealed resource is missing or invalid") {
		t.Errorf("error should carry codesign's verdict, got: %v", err)
	}
}

// The same damage, through the helper: the original verifies, the update does
// not, so the update is rolled back.
func TestRunHelperSwap_TamperedUpdateOverVerifiedOriginal_RollsBack(t *testing.T) {
	tool := requireCodesign(t)
	dir := t.TempDir()
	target := filepath.Join(dir, "App.app")
	newPath := filepath.Join(dir, "staging", "App.app")
	makeRealBundle(t, target)
	adhocSign(t, tool, target)
	makeRealBundle(t, newPath)
	adhocSign(t, tool, newPath)
	tamper(t, newPath)

	code := runHelperSwap(target, newPath, 0, filepath.Join(dir, "log"), instantWaiter, &fakeLauncher{})
	if code != 19 {
		t.Fatalf("code = %d, want 19 (rejected after the swap, backup restored)", code)
	}
	if got := readFile(t, filepath.Join(target, "Contents", "Resources", "data.txt")); string(got) != "sealed" {
		t.Errorf("original application not recovered: %q", got)
	}
}

// A check the installed application cannot pass is not one its update can be
// held to. When the original fails verification too, the verdict on the update
// says nothing about the update, and enforcing it would roll back every update
// to that installation, forever.
func TestRunHelperSwap_OriginalFailsSameCheck_NotEnforced(t *testing.T) {
	tool := requireCodesign(t)
	dir := t.TempDir()
	target := filepath.Join(dir, "App.app")
	newPath := filepath.Join(dir, "staging", "App.app")
	makeRealBundle(t, target)
	adhocSign(t, tool, target)
	tamper(t, target)
	makeRealBundle(t, newPath)
	adhocSign(t, tool, newPath)
	tamper(t, newPath)

	code := runHelperSwap(target, newPath, 0, filepath.Join(dir, "log"), instantWaiter, &fakeLauncher{})
	if code != 0 {
		log, _ := os.ReadFile(filepath.Join(dir, "log"))
		t.Fatalf("code = %d, want 0; log:\n%s", code, log)
	}
}
