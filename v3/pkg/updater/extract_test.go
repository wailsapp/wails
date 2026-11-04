package updater

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// Happy path: a .zip carrying a single top-level .app bundle becomes a
// directory on disk at the staging location. This is the load-bearing path
// for macOS distribution.
func TestMaybeExtract_Zip_AppBundle(t *testing.T) {
	dir := t.TempDir()
	archive := filepath.Join(dir, "MyApp-darwin.zip")
	writeZip(t, archive, []zipEntry{
		{name: "MyApp.app/", isDir: true},
		{name: "MyApp.app/Contents/", isDir: true},
		{name: "MyApp.app/Contents/Info.plist", body: []byte("<plist/>")},
		{name: "MyApp.app/Contents/MacOS/", isDir: true},
		{name: "MyApp.app/Contents/MacOS/MyApp", body: []byte("binary"), mode: 0o755},
	})

	out, did, err := maybeExtractInto(archive)
	if err != nil {
		t.Fatalf("maybeExtractInto: %v", err)
	}
	if !did {
		t.Fatalf("expected extraction to occur")
	}
	if got, want := filepath.Base(out), "MyApp.app"; got != want {
		t.Fatalf("output name: got %q, want %q", got, want)
	}
	if _, err := os.Stat(archive); !os.IsNotExist(err) {
		t.Errorf("archive should have been removed: %v", err)
	}
	if body := readF(t, filepath.Join(out, "Contents", "Info.plist")); string(body) != "<plist/>" {
		t.Errorf("Info.plist contents: %q", body)
	}
	if runtime.GOOS != "windows" {
		info, err := os.Stat(filepath.Join(out, "Contents", "MacOS", "MyApp"))
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode().Perm()&0o111 == 0 {
			t.Errorf("binary should be executable, got mode %o", info.Mode().Perm())
		}
	}
}

// Non-archive artifacts must pass through unchanged so the existing
// flat-binary distribution model keeps working.
func TestMaybeExtract_NonArchive_Passthrough(t *testing.T) {
	dir := t.TempDir()
	staged := filepath.Join(dir, "binary")
	if err := os.WriteFile(staged, []byte("ELF"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, did, err := maybeExtractInto(staged)
	if err != nil {
		t.Fatalf("maybeExtractInto: %v", err)
	}
	if did {
		t.Errorf("expected no extraction for non-archive")
	}
	if out != staged {
		t.Errorf("output path: got %q, want %q", out, staged)
	}
}

// Tar.gz support exercises the same machinery as zip via a different code
// path. Useful for Linux distribution.
func TestMaybeExtract_TarGz_FlatBinary(t *testing.T) {
	dir := t.TempDir()
	archive := filepath.Join(dir, "myapp.tar.gz")
	writeTarGz(t, archive, []tarEntry{
		{name: "myapp", body: []byte("ELF"), mode: 0o755},
	})
	out, did, err := maybeExtractInto(archive)
	if err != nil {
		t.Fatalf("maybeExtractInto: %v", err)
	}
	if !did {
		t.Fatal("expected extraction")
	}
	if filepath.Base(out) != "myapp" {
		t.Errorf("output name: %q", out)
	}
	if string(readF(t, out)) != "ELF" {
		t.Errorf("binary contents")
	}
}

// A zip whose top level contains more than one entry is ambiguous: the
// helper has nothing meaningful to swap into the single target path.
func TestMaybeExtract_MultipleTopLevel_Rejected(t *testing.T) {
	dir := t.TempDir()
	archive := filepath.Join(dir, "two.zip")
	writeZip(t, archive, []zipEntry{
		{name: "one", body: []byte("a")},
		{name: "two", body: []byte("b")},
	})
	_, _, err := maybeExtractInto(archive)
	if err == nil || !strings.Contains(err.Error(), "one top-level entry") {
		t.Fatalf("expected multi-entry rejection, got %v", err)
	}
}

// Classic zip-slip: an entry whose name resolves outside the extraction root
// must be rejected before any bytes hit disk.
func TestMaybeExtract_Zip_SlipRejected(t *testing.T) {
	dir := t.TempDir()
	archive := filepath.Join(dir, "evil.zip")
	writeZip(t, archive, []zipEntry{
		{name: "../escape.txt", body: []byte("pwned")},
	})
	_, _, err := maybeExtractInto(archive)
	if err == nil || !strings.Contains(err.Error(), "escapes root") {
		t.Fatalf("expected zip-slip rejection, got %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "escape.txt")); !os.IsNotExist(err) {
		t.Errorf("escape file should not exist: %v", err)
	}
}

// Tar variant of the same.
func TestMaybeExtract_Tar_SlipRejected(t *testing.T) {
	dir := t.TempDir()
	archive := filepath.Join(dir, "evil.tar.gz")
	writeTarGz(t, archive, []tarEntry{
		{name: "../escape", body: []byte("pwned")},
	})
	_, _, err := maybeExtractInto(archive)
	if err == nil || !strings.Contains(err.Error(), "escapes root") {
		t.Fatalf("expected tar-slip rejection, got %v", err)
	}
}

// Symlinks whose target resolves outside the extraction root must be
// rejected — they would let a crafted archive point at /etc/passwd or any
// other host file after extraction completes.
func TestMaybeExtract_Tar_SymlinkEscape_Rejected(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlinks on Windows require elevation; the path is the same")
	}
	dir := t.TempDir()
	archive := filepath.Join(dir, "evil.tar.gz")
	writeTarGz(t, archive, []tarEntry{
		{name: "bundle/", isDir: true, mode: 0o755},
		{name: "bundle/link", symlinkTarget: "../../../../etc/passwd"},
	})
	_, _, err := maybeExtractInto(archive)
	if err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("expected symlink-escape rejection, got %v", err)
	}
}

// Sanity: a symlink that stays inside the bundle (the common case for macOS
// `.app/Contents/MacOS` → `Frameworks` style links) is allowed.
func TestMaybeExtract_Tar_InternalSymlink_Allowed(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlinks on Windows require elevation")
	}
	dir := t.TempDir()
	archive := filepath.Join(dir, "bundle.tar.gz")
	writeTarGz(t, archive, []tarEntry{
		{name: "bundle/", isDir: true, mode: 0o755},
		{name: "bundle/real", body: []byte("payload"), mode: 0o644},
		{name: "bundle/alias", symlinkTarget: "real"},
	})
	out, _, err := maybeExtractInto(archive)
	if err != nil {
		t.Fatalf("internal symlink rejected: %v", err)
	}
	target, err := os.Readlink(filepath.Join(out, "alias"))
	if err != nil {
		t.Fatalf("readlink: %v", err)
	}
	if target != "real" {
		t.Errorf("symlink target: %q", target)
	}
}

// Zip-slip variants beyond the plain "../" prefix: nested traversal that
// only escapes after cleaning, Windows-style separators, and absolute paths.
// Each must be rejected for both zip and tar.gz and leave nothing outside the
// extraction root.
func TestMaybeExtract_SlipVariants_Rejected(t *testing.T) {
	names := []string{
		"bundle/../../escape.txt",
		`..\escape.txt`,
		`bundle\..\..\escape.txt`,
		"/escape.txt",
		"./../escape.txt",
	}
	for _, name := range names {
		t.Run("zip "+name, func(t *testing.T) {
			dir := t.TempDir()
			archive := filepath.Join(dir, "evil.zip")
			writeZip(t, archive, []zipEntry{{name: name, body: []byte("pwned")}})
			_, _, err := maybeExtractInto(archive)
			if err == nil || !strings.Contains(err.Error(), "escapes root") {
				t.Fatalf("expected zip-slip rejection for %q, got %v", name, err)
			}
			if _, err := os.Stat(filepath.Join(dir, "escape.txt")); !os.IsNotExist(err) {
				t.Errorf("escape file should not exist: %v", err)
			}
		})
		t.Run("tar "+name, func(t *testing.T) {
			dir := t.TempDir()
			archive := filepath.Join(dir, "evil.tar.gz")
			writeTarGz(t, archive, []tarEntry{{name: name, body: []byte("pwned")}})
			_, _, err := maybeExtractInto(archive)
			if err == nil || !strings.Contains(err.Error(), "escapes root") {
				t.Fatalf("expected tar-slip rejection for %q, got %v", name, err)
			}
			if _, err := os.Stat(filepath.Join(dir, "escape.txt")); !os.IsNotExist(err) {
				t.Errorf("escape file should not exist: %v", err)
			}
		})
	}
}

// Zip variant of the symlink-escape check: a zip entry with the symlink mode
// bit whose body points outside the root must be rejected.
func TestMaybeExtract_Zip_SymlinkEscape_Rejected(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlinks on Windows require elevation; the path is the same")
	}
	for _, link := range []string{"../../../../etc/passwd", "/etc/passwd"} {
		dir := t.TempDir()
		archive := filepath.Join(dir, "evil.zip")
		writeZip(t, archive, []zipEntry{
			{name: "bundle/", isDir: true},
			{name: "bundle/link", body: []byte(link), mode: os.ModeSymlink | 0o777},
		})
		_, _, err := maybeExtractInto(archive)
		if err == nil || !strings.Contains(err.Error(), "symlink") {
			t.Fatalf("expected zip symlink-escape rejection for %q, got %v", link, err)
		}
	}
}

// Symlink-then-write-through: an archive plants a symlink pointing at a
// directory outside the root and then writes a file "through" it. The
// symlink is rejected before any write, so nothing lands outside the root.
func TestMaybeExtract_Tar_SymlinkWriteThrough_Rejected(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlinks on Windows require elevation; the path is the same")
	}
	outside := t.TempDir()
	dir := t.TempDir()
	archive := filepath.Join(dir, "evil.tar.gz")
	writeTarGz(t, archive, []tarEntry{
		{name: "bundle/", isDir: true, mode: 0o755},
		{name: "bundle/link", symlinkTarget: outside},
		{name: "bundle/link/pwned.txt", body: []byte("pwned")},
	})
	if _, _, err := maybeExtractInto(archive); err == nil {
		t.Fatal("expected rejection of symlink write-through archive")
	}
	if _, err := os.Stat(filepath.Join(outside, "pwned.txt")); !os.IsNotExist(err) {
		t.Errorf("write-through file should not exist outside root: %v", err)
	}
}

func TestSafeJoin(t *testing.T) {
	root := filepath.Clean(t.TempDir())
	ok := map[string]string{
		"a.txt":        filepath.Join(root, "a.txt"),
		"dir/b.txt":    filepath.Join(root, "dir", "b.txt"),
		`dir\c.txt`:    filepath.Join(root, "dir", "c.txt"),
		"dir/../d.txt": filepath.Join(root, "d.txt"),
		".":            root,
		"":             root,
		"..foo":        filepath.Join(root, "..foo"),
	}
	for name, want := range ok {
		got, err := safeJoin(root, name)
		if err != nil || got != want {
			t.Errorf("safeJoin(%q) = %q, %v; want %q", name, got, err, want)
		}
	}
	for _, name := range []string{"..", "../x", "a/../../x", "/abs", `..\x`} {
		if got, err := safeJoin(root, name); err == nil {
			t.Errorf("safeJoin(%q) = %q; want escape error", name, got)
		}
	}
}

// --- test fixtures ---

type zipEntry struct {
	name  string
	body  []byte
	isDir bool
	mode  os.FileMode
}

func writeZip(t *testing.T, path string, entries []zipEntry) {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for _, e := range entries {
		hdr := &zip.FileHeader{Name: e.name, Method: zip.Deflate}
		mode := e.mode
		if mode == 0 {
			if e.isDir {
				mode = 0o755 | os.ModeDir
			} else {
				mode = 0o644
			}
		} else if e.isDir {
			mode |= os.ModeDir
		}
		hdr.SetMode(mode)
		w, err := zw.CreateHeader(hdr)
		if err != nil {
			t.Fatal(err)
		}
		if !e.isDir {
			if _, err := w.Write(e.body); err != nil {
				t.Fatal(err)
			}
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
}

type tarEntry struct {
	name          string
	body          []byte
	isDir         bool
	mode          os.FileMode
	symlinkTarget string
}

func writeTarGz(t *testing.T, path string, entries []tarEntry) {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for _, e := range entries {
		hdr := &tar.Header{Name: e.name, Mode: int64(e.mode)}
		switch {
		case e.symlinkTarget != "":
			hdr.Typeflag = tar.TypeSymlink
			hdr.Linkname = e.symlinkTarget
		case e.isDir:
			hdr.Typeflag = tar.TypeDir
			if hdr.Mode == 0 {
				hdr.Mode = 0o755
			}
		default:
			hdr.Typeflag = tar.TypeReg
			hdr.Size = int64(len(e.body))
			if hdr.Mode == 0 {
				hdr.Mode = 0o644
			}
		}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}
		if hdr.Typeflag == tar.TypeReg {
			if _, err := tw.Write(e.body); err != nil {
				t.Fatal(err)
			}
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readF(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
