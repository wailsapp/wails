package updater

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Limits applied to archive extraction. Apps ship code, not blobs — a single
// macOS .app bundle is typically <100 MiB and a few thousand entries. These
// caps give comfortable headroom while keeping malformed / hostile archives
// from exhausting disk or file descriptors.
const (
	maxArchiveEntries   = 50_000
	maxArchiveTotalSize = 2 << 30 // 2 GiB total uncompressed
)

// archiveKind is the set of archive formats the framework can unpack inline
// before handing the artifact to the helper for swap.
type archiveKind int

const (
	archiveNone archiveKind = iota
	archiveZip
	archiveTarGz
)

// detectArchive classifies path by filename extension. Magic-byte sniffing
// would be more robust but every shipping CDN preserves extensions, and a
// false-positive on detection silently breaks the swap — so we prefer the
// conservative reading: only extract when the extension says so.
func detectArchive(path string) archiveKind {
	lower := strings.ToLower(path)
	switch {
	case strings.HasSuffix(lower, ".zip"):
		return archiveZip
	case strings.HasSuffix(lower, ".tar.gz"), strings.HasSuffix(lower, ".tgz"):
		return archiveTarGz
	}
	return archiveNone
}

// maybeExtractInto unpacks an archive at stagedPath in place: the archive is
// expanded into a scratch subdirectory of its parent, the original archive is
// removed, and the single top-level entry of the archive is moved up to where
// the archive used to live. The returned newPath either equals stagedPath
// (non-archive — nothing was done) or names the unpacked entry.
//
// Constraints enforced:
//   - Exactly one top-level entry. Archives with multiple top-level files or
//     directories are rejected: the helper has nothing meaningful to swap
//     into a single target path. Distribute one .app bundle (or one binary)
//     per artifact.
//   - No path traversal. Entries whose normalised path escapes the extraction
//     root are rejected (zip-slip).
//   - No escaping symlinks. Symlinks whose target resolves outside the
//     extraction root are rejected; in-archive symlinks pointing to sibling
//     entries are preserved.
//   - Total uncompressed size and entry count are capped (see constants).
func maybeExtractInto(stagedPath string) (newPath string, didExtract bool, err error) {
	kind := detectArchive(stagedPath)
	if kind == archiveNone {
		return stagedPath, false, nil
	}

	parent := filepath.Dir(stagedPath)
	scratch, err := os.MkdirTemp(parent, ".payload-*")
	if err != nil {
		return "", false, fmt.Errorf("updater: extract scratch dir: %w", err)
	}
	cleanup := func() { _ = os.RemoveAll(scratch) }

	switch kind {
	case archiveZip:
		err = extractZip(stagedPath, scratch)
	case archiveTarGz:
		err = extractTarGz(stagedPath, scratch)
	}
	if err != nil {
		cleanup()
		return "", false, err
	}

	entries, err := os.ReadDir(scratch)
	if err != nil {
		cleanup()
		return "", false, fmt.Errorf("updater: extract: %w", err)
	}
	if len(entries) != 1 {
		cleanup()
		return "", false, fmt.Errorf("updater: archive must contain exactly one top-level entry, got %d", len(entries))
	}
	entry := entries[0]

	// Remove the original archive so the staging dir holds only the extracted
	// payload — the helper's post-swap cleanup walks one level up.
	if err := os.Remove(stagedPath); err != nil {
		cleanup()
		return "", false, fmt.Errorf("updater: remove staged archive: %w", err)
	}

	src := filepath.Join(scratch, entry.Name())
	dst := filepath.Join(parent, entry.Name())
	if err := os.Rename(src, dst); err != nil {
		cleanup()
		return "", false, fmt.Errorf("updater: promote extracted entry: %w", err)
	}
	cleanup()
	return dst, true, nil
}

// extractZip unpacks src into dst. Path traversal, symlink escape, entry
// count, and total-size caps are all enforced.
//
// Symlinks are deferred to a second pass: planting a symlink first and then
// writing through it (e.g. archive contains "link → /etc" followed by
// "link/passwd") is the standard zip-slip-via-symlink escalation. Two
// passes mean every file write happens before any symlinks exist in the
// destination tree, so writes can't be redirected through them.
func extractZip(src, dst string) error {
	zr, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("updater: zip open: %w", err)
	}
	defer zr.Close()

	if len(zr.File) > maxArchiveEntries {
		return fmt.Errorf("updater: zip has %d entries (cap %d)", len(zr.File), maxArchiveEntries)
	}

	rootClean := filepath.Clean(dst)

	// Pass 1 — directories and regular files. Symlinks recorded for pass 2.
	type pendingLink struct{ src *zip.File; target string }
	var symlinks []pendingLink
	var written int64
	for _, f := range zr.File {
		target, err := safeJoin(rootClean, f.Name)
		if err != nil {
			return err
		}
		mode := f.Mode()
		switch {
		case mode&os.ModeSymlink != 0:
			// Reject escapes up front so we don't waste pass 1 on a doomed archive.
			if err := validateSymlinkTarget(f, target, rootClean); err != nil {
				return err
			}
			symlinks = append(symlinks, pendingLink{src: f, target: target})
		case f.FileInfo().IsDir():
			if err := os.MkdirAll(target, dirModeFrom(mode)); err != nil {
				return fmt.Errorf("updater: zip mkdir: %w", err)
			}
		default:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return fmt.Errorf("updater: zip mkdir: %w", err)
			}
			n, err := writeArchiveFile(f, target, mode, &written)
			if err != nil {
				return err
			}
			written = n
		}
	}

	// Pass 2 — symlinks only. By the time any of these are created, the file
	// tree on disk matches what was in the archive (sans symlinks), so a link
	// that resolves to a directory under root can never have been used to
	// redirect a write from pass 1.
	for _, sl := range symlinks {
		if err := writeArchiveSymlink(sl.src, sl.target, rootClean); err != nil {
			return err
		}
	}
	return nil
}

// extractTarGz unpacks src (a gzipped tar) into dst, enforcing the same
// invariants as extractZip. See that function for the rationale behind the
// two-pass extraction (regular files + dirs first, symlinks last).
func extractTarGz(src, dst string) error {
	f, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("updater: tar.gz open: %w", err)
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("updater: gzip: %w", err)
	}
	defer gz.Close()
	tr := tar.NewReader(gz)

	rootClean := filepath.Clean(dst)
	type pendingLink struct{ linkname, target string }
	var symlinks []pendingLink
	var entries int
	var written int64
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("updater: tar: %w", err)
		}
		entries++
		if entries > maxArchiveEntries {
			return fmt.Errorf("updater: tar has more than %d entries", maxArchiveEntries)
		}
		target, err := safeJoin(rootClean, hdr.Name)
		if err != nil {
			return err
		}
		mode := os.FileMode(hdr.Mode)
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, dirModeFrom(mode)); err != nil {
				return fmt.Errorf("updater: tar mkdir: %w", err)
			}
		case tar.TypeSymlink:
			// Reject escapes up front; defer creation until pass 2 so no
			// file write can be redirected through a planted symlink.
			if err := validateSymlinkPath(hdr.Linkname, target, rootClean); err != nil {
				return err
			}
			symlinks = append(symlinks, pendingLink{linkname: hdr.Linkname, target: target})
		case tar.TypeReg, tar.TypeRegA:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return fmt.Errorf("updater: tar mkdir: %w", err)
			}
			out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fileModeFrom(mode))
			if err != nil {
				return fmt.Errorf("updater: tar create: %w", err)
			}
			n, copyErr := io.CopyN(out, tr, maxArchiveTotalSize-written+1)
			closeErr := out.Close()
			if copyErr != nil && !errors.Is(copyErr, io.EOF) {
				return fmt.Errorf("updater: tar copy: %w", copyErr)
			}
			if closeErr != nil {
				return fmt.Errorf("updater: tar close: %w", closeErr)
			}
			written += n
			if written > maxArchiveTotalSize {
				return fmt.Errorf("updater: tar uncompressed size exceeds %d bytes", maxArchiveTotalSize)
			}
		default:
			// Block-special, char-special, FIFO, etc. — skip silently. App
			// bundles never contain these.
		}
	}

	// Pass 2 — symlinks last (see extractZip rationale).
	for _, sl := range symlinks {
		if err := writeSymlinkPath(sl.linkname, sl.target, rootClean); err != nil {
			return err
		}
	}
	return nil
}

// writeArchiveFile streams one zip entry to disk, accumulating into the
// running total-uncompressed-size budget.
func writeArchiveFile(f *zip.File, target string, mode os.FileMode, written *int64) (int64, error) {
	rc, err := f.Open()
	if err != nil {
		return *written, fmt.Errorf("updater: zip read: %w", err)
	}
	defer rc.Close()
	out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fileModeFrom(mode))
	if err != nil {
		return *written, fmt.Errorf("updater: zip create: %w", err)
	}
	// Limit each copy so a single bomb entry can't blow past the budget.
	remaining := maxArchiveTotalSize - *written + 1
	n, err := io.CopyN(out, rc, remaining)
	closeErr := out.Close()
	if err != nil && !errors.Is(err, io.EOF) {
		return *written, fmt.Errorf("updater: zip copy: %w", err)
	}
	if closeErr != nil {
		return *written, fmt.Errorf("updater: zip close: %w", closeErr)
	}
	total := *written + n
	if total > maxArchiveTotalSize {
		return total, fmt.Errorf("updater: zip uncompressed size exceeds %d bytes", maxArchiveTotalSize)
	}
	return total, nil
}

// validateSymlinkTarget reads a zip entry's symlink body and validates that
// the link target stays inside the extraction root. Used in pass 1 so we can
// reject malicious archives before any file write happens.
func validateSymlinkTarget(f *zip.File, target, root string) error {
	rc, err := f.Open()
	if err != nil {
		return fmt.Errorf("updater: zip symlink read: %w", err)
	}
	defer rc.Close()
	link, err := io.ReadAll(io.LimitReader(rc, 4096))
	if err != nil {
		return fmt.Errorf("updater: zip symlink body: %w", err)
	}
	return validateSymlinkPath(string(link), target, root)
}

// validateSymlinkPath checks that creating target → link would not escape
// root. Returns nil if safe.
func validateSymlinkPath(link, target, root string) error {
	if filepath.IsAbs(link) {
		return fmt.Errorf("updater: archive symlink has absolute target: %s", link)
	}
	resolved := filepath.Join(filepath.Dir(target), link)
	rel, err := filepath.Rel(root, resolved)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return fmt.Errorf("updater: archive symlink escapes root: %s -> %s", target, link)
	}
	return nil
}

// writeArchiveSymlink writes a previously-validated zip symlink entry. Called
// only from pass 2 of extractZip, after all regular files and directories
// have been written, so a planted symlink can't be used to redirect any
// earlier write.
func writeArchiveSymlink(f *zip.File, target, root string) error {
	rc, err := f.Open()
	if err != nil {
		return fmt.Errorf("updater: zip symlink read: %w", err)
	}
	defer rc.Close()
	link, err := io.ReadAll(io.LimitReader(rc, 4096))
	if err != nil {
		return fmt.Errorf("updater: zip symlink body: %w", err)
	}
	return writeSymlinkPath(string(link), target, root)
}

// writeSymlinkPath creates target as a symlink to link. Re-validates the
// path even though pass 1 already checked it — defence in depth.
func writeSymlinkPath(link, target, root string) error {
	if err := validateSymlinkPath(link, target, root); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("updater: archive symlink mkdir: %w", err)
	}
	// Remove any existing entry — happens on re-extraction in tests.
	_ = os.Remove(target)
	if err := os.Symlink(link, target); err != nil {
		return fmt.Errorf("updater: archive symlink create: %w", err)
	}
	return nil
}

// safeJoin resolves name relative to root and returns the cleaned absolute
// path, rejecting any entry whose normalised path escapes root (zip-slip).
func safeJoin(root, name string) (string, error) {
	// Normalise separators so a Windows-produced zip with "Foo\\bar" doesn't
	// slip past the prefix check on POSIX hosts.
	clean := filepath.Clean(strings.ReplaceAll(name, `\`, "/"))
	if clean == "" || clean == "." {
		return root, nil
	}
	if filepath.IsAbs(clean) || strings.HasPrefix(clean, "..") {
		return "", fmt.Errorf("updater: archive entry escapes root: %s", name)
	}
	target := filepath.Join(root, clean)
	rel, err := filepath.Rel(root, target)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("updater: archive entry escapes root: %s", name)
	}
	return target, nil
}

// fileModeFrom returns a safe file mode for an extracted regular file.
// Archives produced on Windows often report mode 0, which would create a
// file with no permission bits; substitute 0644 in that case. The executable
// bit is preserved when the archive provides it (e.g. the binary inside a
// macOS .app bundle).
func fileModeFrom(m os.FileMode) os.FileMode {
	perm := m.Perm()
	if perm == 0 {
		return 0o644
	}
	return perm
}

// dirModeFrom is the directory analogue of fileModeFrom.
func dirModeFrom(m os.FileMode) os.FileMode {
	perm := m.Perm()
	if perm == 0 {
		return 0o755
	}
	return perm
}
