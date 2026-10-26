package updater

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// A staged .app bundle emptied while the helper waits for the parent to exit
// must not be installed. os.Rename succeeds on a directory whether or not it
// still holds anything, so without a check the helper renames the husk over a
// working application, reports success, and deletes the only copy of what it
// replaced. That is how a working install becomes a zero-byte one.
//
// helper_recovery_test.go covers newPath being REMOVED, which makes the rename
// fail and rolls back on its own. Emptying is the case a bundle actually hits,
// because the thing deleting it walks into the directory rather than unlinking
// it whole.
func TestRunHelperSwap_StagedBundleEmptiedDuringWait(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "App.app")
	newPath := filepath.Join(dir, "staging", "App.app")
	makeAppBundle(t, target, "old-bin")
	makeAppBundle(t, newPath, "new-bin")

	// Stand in for a concurrent discardStaging walking the staging directory
	// while the helper is blocked waiting for the parent.
	wait := func(int, time.Duration) error {
		return os.RemoveAll(filepath.Join(newPath, "Contents"))
	}

	launched := 0
	l := &funcLauncher{fn: func(path string) error {
		launched++
		if path != target {
			t.Errorf("relaunched %q, want the untouched original %q", path, target)
		}
		return nil
	}}

	code := runHelperSwap(target, newPath, 1234, filepath.Join(dir, "log"), wait, l)
	if code != 18 {
		t.Fatalf("code = %d, want 18 (payload not installable)", code)
	}
	if got := readFile(t, filepath.Join(target, "Contents", "MacOS", "exe")); string(got) != "old-bin" {
		t.Errorf("original application was damaged: %q", got)
	}
	if launched != 1 {
		t.Errorf("launches = %d, want the original relaunched exactly once", launched)
	}
}

func TestRunHelperSwap_InvalidStagedBundleWithoutAppSuffix(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "App.app")
	newPath := filepath.Join(dir, "App.app.new")
	makeAppBundle(t, target, "old-bin")
	writeFile(t, filepath.Join(newPath, "Contents", "junk"), []byte("x"))
	l := &funcLauncher{fn: func(string) error { return nil }}
	code := runHelperSwap(target, newPath, 0, filepath.Join(dir, "log"), instantWaiter, l)
	if code != 18 {
		t.Fatalf("code = %d, want 18 (rejected before swap)", code)
	}
	if _, err := os.Stat(target + ".bak"); !os.IsNotExist(err) {
		t.Fatalf("backup created before validation: %v", err)
	}
	if got := readFile(t, filepath.Join(target, "Contents", "MacOS", "exe")); string(got) != "old-bin" {
		t.Fatalf("original changed: %q", got)
	}
}

func TestValidatePayloadRejectsMalformedPlist(t *testing.T) {
	for _, plist := range []string{
		"<plist><dict><key>CFBundleExecutable</key>",
		"<plist><dict><key>CFBundleExecutable</key><string>Main</string>",
		"not a plist",
		"<plist><dict></dict></plist><plist/>",
		"<plist><dict></dict></plist>trailing junk",
	} {
		t.Run(plist, func(t *testing.T) {
			bundle := filepath.Join(t.TempDir(), "App.app")
			makeNamedAppBundle(t, bundle, "Main")
			writeFile(t, filepath.Join(bundle, "Contents", "Info.plist"), []byte(plist))
			if err := validatePayload(bundle); err == nil {
				t.Fatal("malformed Info.plist accepted")
			}
			if err := validateInstalled(bundle, "", nil); err == nil {
				t.Fatal("malformed installed Info.plist accepted")
			}
		})
	}
}

func TestReplaceTargetRejectsInvalidBundleWithoutAppSuffix(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "App.app")
	makeAppBundle(t, target, "old-bin")
	newPath := filepath.Join(dir, "staged")
	writeFile(t, filepath.Join(newPath, "junk"), []byte("x"))
	if err := replaceTarget(target, newPath); !errors.Is(err, errNotInstallable) {
		t.Fatalf("got %v, want validation refusal", err)
	}
	if got := readFile(t, filepath.Join(target, "Contents", "MacOS", "exe")); string(got) != "old-bin" {
		t.Fatalf("original changed: %q", got)
	}
}

func TestValidatePayload(t *testing.T) {
	dir := t.TempDir()

	t.Run("empty file rejected", func(t *testing.T) {
		p := filepath.Join(dir, "empty.bin")
		writeFile(t, p, nil)
		if err := validatePayload(p); err == nil {
			t.Fatal("an empty binary must not be installable")
		}
	})

	t.Run("non-empty file accepted", func(t *testing.T) {
		p := filepath.Join(dir, "ok.bin")
		writeFile(t, p, []byte("ELF"))
		if err := validatePayload(p); err != nil {
			t.Fatalf("unexpected: %v", err)
		}
	})

	t.Run("empty directory rejected", func(t *testing.T) {
		p := filepath.Join(dir, "hollow")
		if err := os.MkdirAll(p, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := validatePayload(p); err == nil {
			t.Fatal("an empty directory must not be installable")
		}
	})

	t.Run("directory of empty files rejected", func(t *testing.T) {
		p := filepath.Join(dir, "husk")
		writeFile(t, filepath.Join(p, "a"), nil)
		writeFile(t, filepath.Join(p, "b"), nil)
		if err := validatePayload(p); err == nil {
			t.Fatal("a directory holding nothing but empty files must not be installable")
		}
	})

	t.Run("bundle with zero-length binary rejected", func(t *testing.T) {
		p := filepath.Join(dir, "Zero.app")
		makeAppBundle(t, p, "real")
		exe := filepath.Join(p, "Contents", "MacOS", "exe")
		if err := os.Truncate(exe, 0); err != nil {
			t.Fatal(err)
		}
		err := validatePayload(p)
		if err == nil {
			t.Fatal("a bundle whose only executable is zero-length must not be installable")
		}
		if !strings.Contains(err.Error(), "MacOS") {
			t.Errorf("error should name the missing executable, got: %v", err)
		}
	})

	t.Run("bundle without Info.plist rejected", func(t *testing.T) {
		p := filepath.Join(dir, "NoPlist.app")
		makeAppBundle(t, p, "real")
		if err := os.Remove(filepath.Join(p, "Contents", "Info.plist")); err != nil {
			t.Fatal(err)
		}
		if err := validatePayload(p); err == nil {
			t.Fatal("a bundle without an Info.plist must not be installable")
		}
	})

	t.Run("well-formed bundle accepted", func(t *testing.T) {
		p := filepath.Join(dir, "Good.app")
		makeAppBundle(t, p, "real")
		if err := validatePayload(p); err != nil {
			t.Fatalf("unexpected: %v", err)
		}
	})

	// The executable the bundle launches is the one Info.plist names. An intact
	// helper beside a truncated main binary is still a bundle that cannot start.
	t.Run("truncated CFBundleExecutable beside an intact helper rejected", func(t *testing.T) {
		p := filepath.Join(dir, "Named.app")
		makeNamedAppBundle(t, p, "main")
		helper := filepath.Join(p, "Contents", "MacOS", "helper")
		writeFile(t, helper, []byte("helper-bin"))
		if err := os.Chmod(helper, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := validatePayload(p); err != nil {
			t.Fatalf("intact bundle rejected: %v", err)
		}
		if err := os.Truncate(filepath.Join(p, "Contents", "MacOS", "main"), 0); err != nil {
			t.Fatal(err)
		}
		err := validatePayload(p)
		if err == nil {
			t.Fatal("a bundle whose CFBundleExecutable is zero-length must not be installable")
		}
		if !strings.Contains(err.Error(), "main") {
			t.Errorf("error should name the bundle executable, got: %v", err)
		}
	})

	t.Run("missing CFBundleExecutable rejected", func(t *testing.T) {
		p := filepath.Join(dir, "Missing.app")
		makeNamedAppBundle(t, p, "main")
		other := filepath.Join(p, "Contents", "MacOS", "other")
		writeFile(t, other, []byte("bin"))
		if err := os.Chmod(other, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Remove(filepath.Join(p, "Contents", "MacOS", "main")); err != nil {
			t.Fatal(err)
		}
		if err := validatePayload(p); err == nil {
			t.Fatal("a bundle without the executable its Info.plist names must not be installable")
		}
	})

	// A binary plist cannot be read without plutil, which only exists on macOS.
	t.Run("binary plist names the executable", func(t *testing.T) {
		plutil, err := exec.LookPath("plutil")
		if runtime.GOOS != "darwin" || err != nil {
			t.Skip("plutil unavailable")
		}
		p := filepath.Join(dir, "Binary.app")
		makeNamedAppBundle(t, p, "main")
		helper := filepath.Join(p, "Contents", "MacOS", "helper")
		writeFile(t, helper, []byte("helper-bin"))
		if err := os.Chmod(helper, 0o755); err != nil {
			t.Fatal(err)
		}
		plist := filepath.Join(p, "Contents", "Info.plist")
		if out, err := exec.Command(plutil, "-convert", "binary1", plist).CombinedOutput(); err != nil {
			t.Fatalf("plutil -convert: %v: %s", err, out)
		}
		if err := os.Truncate(filepath.Join(p, "Contents", "MacOS", "main"), 0); err != nil {
			t.Fatal(err)
		}
		if err := validatePayload(p); err == nil {
			t.Fatal("a zero-length CFBundleExecutable named by a binary plist must not be installable")
		}
	})

	// A symlink is renamed into place as a symlink, still pointing into the
	// staging directory the helper deletes once the swap succeeds. Validation
	// never follows links, so the payload itself must not be one.
	t.Run("symlinked payload rejected", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("symlinks on Windows require elevation")
		}
		real := filepath.Join(dir, "real.bin")
		writeFile(t, real, []byte("ELF"))
		link := filepath.Join(dir, "link.bin")
		if err := os.Symlink(real, link); err != nil {
			t.Fatal(err)
		}
		err := validatePayload(link)
		if err == nil {
			t.Fatal("a symlinked payload must not be installable")
		}
		if !strings.Contains(err.Error(), "symbolic link") {
			t.Errorf("error should say the payload is a link, got: %v", err)
		}
	})
}

// replaceTarget refuses a payload that is not installable before it disturbs
// the target, and says so distinctly: retrying cannot repair a payload, so the
// swap loop stops rather than spending its retries on one.
func TestReplaceTarget_InvalidPayloadIsTerminal(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "App.app")
	makeAppBundle(t, target, "old-bin")
	hollow := filepath.Join(dir, "staging", "App.app")
	if err := os.MkdirAll(hollow, 0o755); err != nil {
		t.Fatal(err)
	}

	err := replaceTarget(target, hollow)
	if !errors.Is(err, errNotInstallable) {
		t.Fatalf("err = %v, want errNotInstallable", err)
	}
	if got := readFile(t, filepath.Join(target, "Contents", "MacOS", "exe")); string(got) != "old-bin" {
		t.Errorf("target disturbed: %q", got)
	}
}

// makeNamedAppBundle creates a bundle whose Info.plist names exe as its
// CFBundleExecutable, in the XML form a build writes.
func makeNamedAppBundle(t *testing.T, root, exe string) {
	t.Helper()
	path := filepath.Join(root, "Contents", "MacOS", exe)
	writeFile(t, path, []byte("main-bin"))
	if err := os.Chmod(path, 0o755); err != nil {
		t.Fatal(err)
	}
	plist := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleName</key>
	<string>Named</string>
	<key>CFBundleExecutable</key>
	<string>` + exe + `</string>
</dict>
</plist>
`
	writeFile(t, filepath.Join(root, "Contents", "Info.plist"), []byte(plist))
}

// An unsigned bundle is unverifiable, not corrupt: a build made without signing
// credentials must still be able to update itself.
func TestVerifySignature_UnsignedBundleIsNotAFailure(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "Unsigned.app")
	makeAppBundle(t, p, "real")

	note, err := verifySignature(p, "")
	if err != nil {
		t.Fatalf("an unsigned bundle must not be treated as corrupt: %v", err)
	}
	t.Logf("note: %q", note)
}

func TestIsUnsigned(t *testing.T) {
	if !isUnsigned([]byte("/x/Foo.app: code object is not signed at all")) {
		t.Error("codesign's unsigned verdict must be recognised")
	}
	if isUnsigned([]byte("/x/Foo.app: a sealed resource is missing or invalid")) {
		t.Error("a genuine rejection must not be read as unsigned")
	}
}

// Only binary plists may fall back to plutil on macOS; XML plists, including
// malformed ones, must be judged by the strict XML reader on every platform.
func TestIsBinaryPlist(t *testing.T) {
	for _, tc := range []struct {
		data string
		want bool
	}{
		{"bplist00\x00\x01", true},
		{"<plist><dict></dict></plist>trailing junk", false},
		{"<?xml version=\"1.0\"?><plist><dict></dict></plist>", false},
		{"not a plist", false},
		{"", false},
	} {
		if got := isBinaryPlist([]byte(tc.data)); got != tc.want {
			t.Errorf("isBinaryPlist(%q) = %v, want %v", tc.data, got, tc.want)
		}
	}
}
