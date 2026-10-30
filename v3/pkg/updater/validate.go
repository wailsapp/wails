package updater

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Payload validation.
//
// The download is verified by digest while it streams (see download.go), but
// that digest covers the ARCHIVE, which extract.go deletes as soon as it has
// unpacked it. Nothing downstream of that point — the extracted payload
// sitting in a temp directory, and the bundle once it has been renamed into
// place — is covered by it.
//
// That gap is load-bearing rather than theoretical. A staged .app bundle is a
// DIRECTORY, and os.Rename of a directory succeeds whether or not it still has
// anything in it, so a payload emptied between staging and the swap would
// install cleanly and report success. These checks close it at both ends:
// refuse to swap in a payload that cannot be an application, and refuse to
// throw away the rollback until the thing now on disk looks like one.

// codesignTimeout bounds the macOS signature check. It reads every sealed file
// in the bundle, so it is not instant on a large one; a bound means a wedged
// codesign degrades to "unverified" rather than hanging the swap forever.
const codesignTimeout = 90 * time.Second

// errNotInstallable marks every error validatePayload returns. Retrying a swap
// cannot repair a payload, so the swap loop stops on it.
var errNotInstallable = errors.New("payload failed validation")

// validatePayload reports whether path is something worth installing. It is
// deliberately cheap and structural: it answers "could this possibly be an
// application" rather than "is this the right application", which is the
// digest's job and already done by the time we get here.
//
// A regular file must be non-empty. A directory must contain at least one
// non-empty regular file — that alone rejects the emptied-bundle case. A
// directory that looks like a macOS .app is held to the bundle layout too.
//
// Symbolic links are never followed, at the top or anywhere inside, and the
// payload itself may not be one: it would be renamed into place as a link
// still pointing into the staging directory, which is deleted once the swap
// succeeds.
func validatePayload(path string) error {
	return validatePayloadForTarget(path, path)
}

func validatePayloadForTarget(path, target string) error {
	if err := checkPayload(path, filepath.Ext(target) == ".app"); err != nil {
		return fmt.Errorf("%w: %w", errNotInstallable, err)
	}
	return nil
}

func checkPayload(path string, appBundle bool) error {
	info, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("staged payload unreadable: %w", err)
	}

	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("staged payload %s is a symbolic link", filepath.Base(path))
	}

	if !info.IsDir() {
		if appBundle {
			return fmt.Errorf("staged application bundle is not a directory")
		}
		if info.Size() == 0 {
			return fmt.Errorf("staged payload %s is empty", filepath.Base(path))
		}
		return nil
	}

	if appBundle {
		if err := validateAppBundle(path); err != nil {
			return err
		}
	}

	found, err := hasContent(path)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("staged payload %s is an empty directory", filepath.Base(path))
	}
	return nil
}

// hasContent reports whether root holds at least one regular file of non-zero
// length. It stops at the first one rather than walking the whole tree, so the
// cost is a few directory reads and not a stat of every file in the bundle.
func hasContent(root string) (bool, error) {
	found := false
	err := filepath.Walk(root, func(_ string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.Mode().IsRegular() && fi.Size() > 0 {
			found = true
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return false, fmt.Errorf("staged payload unreadable: %w", err)
	}
	return found, nil
}

// validateAppBundle checks the parts of a macOS .app without which the bundle
// cannot launch: an Info.plist with something in it, and a non-empty main
// executable in Contents/MacOS. A zero-length main binary is exactly the shape
// a partially-deleted bundle takes, and is the failure this is here to catch.
//
// The main executable is the one Info.plist names as CFBundleExecutable, so an
// intact helper beside a truncated main binary does not pass. When that key
// is absent in a readable plist, any non-empty executable in Contents/MacOS
// is accepted.
//
// The layout is checked on every platform. A .app is a .app wherever it was
// unpacked, and a build that cross-stages one should not skip the check for
// want of a matching GOOS. "Executable" is the exec bit only where the file
// system records one: Windows has none (os.Stat reports 0666 or 0444), so
// there any non-empty regular file counts.
func validateAppBundle(bundle string) error {
	plist := filepath.Join(bundle, "Contents", "Info.plist")
	switch fi, err := os.Stat(plist); {
	case err != nil:
		return fmt.Errorf("bundle has no readable Contents/Info.plist: %w", err)
	case fi.Size() == 0:
		return fmt.Errorf("bundle Contents/Info.plist is empty")
	}

	macOS := filepath.Join(bundle, "Contents", "MacOS")
	name, readable := bundleExecutable(plist)
	if !readable {
		return fmt.Errorf("bundle Contents/Info.plist is malformed or unreadable")
	}
	if name != "" {
		fi, err := os.Lstat(filepath.Join(macOS, name))
		if err != nil {
			return fmt.Errorf("bundle executable Contents/MacOS/%s unreadable: %w", name, err)
		}
		if !isNonEmptyExecutable(fi) {
			return fmt.Errorf("bundle executable Contents/MacOS/%s is not a non-empty executable file", name)
		}
		return nil
	}

	entries, err := os.ReadDir(macOS)
	if err != nil {
		return fmt.Errorf("bundle has no readable Contents/MacOS: %w", err)
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		fi, err := e.Info()
		if err != nil {
			return fmt.Errorf("bundle executable unreadable: %w", err)
		}
		if isNonEmptyExecutable(fi) {
			return nil
		}
	}
	return fmt.Errorf("bundle Contents/MacOS holds no non-empty executable")
}

func isNonEmptyExecutable(fi os.FileInfo) bool {
	executable := runtime.GOOS == "windows" || fi.Mode().Perm()&0o111 != 0
	return fi.Mode().IsRegular() && fi.Size() > 0 && executable
}

// maxPlistSize bounds how much of an Info.plist is read to find one key. Real
// ones are a few kilobytes.
const maxPlistSize = 1 << 20

// bundleExecutable preserves the distinction between an unreadable plist and
// a readable plist without CFBundleExecutable, which may use the fallback.
func bundleExecutable(path string) (string, bool) {
	f, err := os.Open(path)
	if err != nil {
		return "", false
	}
	data, err := io.ReadAll(io.LimitReader(f, maxPlistSize+1))
	_ = f.Close()
	if err != nil || len(data) > maxPlistSize {
		return "", false
	}
	name, ok := xmlPlistString(data, "CFBundleExecutable")
	// plutil is only a fallback for binary plists. Apple's parser accepts XML
	// that the strict reader above rejects (an unclosed dict, content after
	// </plist>), so an XML plist must not get a second chance through it.
	if !ok && runtime.GOOS == "darwin" && isBinaryPlist(data) {
		name, ok = plutilString(path, "CFBundleExecutable")
	}
	if name == "." || name == ".." || strings.ContainsAny(name, `/\`) {
		return "", false
	}
	return name, ok
}

// xmlPlistString reads a top-level dictionary string and validates the entire
// document, including any content after the requested value.
func xmlPlistString(data []byte, key string) (value string, ok bool) {
	dec := xml.NewDecoder(bytes.NewReader(data))
	depth := 0
	sawPlist, sawDict, matched := false, false, false
	for {
		tok, err := dec.Token()
		if err != nil {
			return value, err == io.EOF && depth == 0 && sawDict
		}
		switch t := tok.(type) {
		case xml.StartElement:
			if depth == 0 {
				if t.Name.Local != "plist" || sawPlist {
					return "", false
				}
				sawPlist = true
			}
			if depth == 1 {
				if t.Name.Local != "dict" || sawDict {
					return "", false
				}
				sawDict = true
			}
			if depth != 2 {
				depth++
				continue
			}
			var text string
			if err := dec.DecodeElement(&text, &t); err != nil {
				return "", false
			}
			switch {
			case t.Name.Local == "key":
				matched = strings.TrimSpace(text) == key
			case matched:
				if t.Name.Local != "string" {
					return "", false
				}
				value = strings.TrimSpace(text)
				matched = false
			}
		case xml.CharData:
			if depth < 2 && strings.TrimSpace(string(t)) != "" {
				return "", false
			}
		case xml.EndElement:
			depth--
		}
	}
}

// isBinaryPlist reports whether data starts with the binary property list
// header ("bplist", followed by a format version such as "00").
func isBinaryPlist(data []byte) bool {
	return bytes.HasPrefix(data, []byte("bplist"))
}

// plutilString validates a binary plist before looking up its executable key.
func plutilString(path, key string) (string, bool) {
	tool, err := exec.LookPath("plutil")
	if err != nil {
		return "", false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := exec.CommandContext(ctx, tool, "-lint", path).Run(); err != nil {
		return "", false
	}
	out, err := exec.CommandContext(ctx, tool, "-extract", key, "raw", "-o", "-", path).Output()
	if ctx.Err() != nil {
		return "", false
	}
	if err != nil {
		return "", true
	}
	return strings.TrimSpace(string(out)), true
}

// validateInstalled runs the post-swap checks over whatever now sits at path:
// the structural check the payload already passed (cheap, and it catches a
// swap that moved the wrong thing), then the platform signature check, judged
// against original — the application that was replaced, still intact as the
// backup.
//
// A signature that could not be run is logged and tolerated; one that ran and
// returned a verdict is decisive, and the caller still holds the backup.
func validateInstalled(path, original string, lg *helperLog) error {
	if err := validatePayload(path); err != nil {
		return err
	}
	note, err := verifySignature(path, original)
	if err != nil {
		return err
	}
	if note != "" {
		lg.logf("%s", note)
	}
	return nil
}

// verifySignature runs the platform's own integrity check over an installed
// payload, and is the per-file checksum the archive digest stopped providing
// the moment extract.go deleted the archive: codesign validates every sealed
// file against the hashes in the bundle's CodeDirectory, so a truncated or
// substituted binary fails it.
//
// The failure modes are kept apart deliberately. A codesign that RAN and
// rejected the bundle is a real verdict and returns an error. Anything that
// merely means "no opinion available" returns nil with a note for the log.
// Refusing on those terms would strand people whose situation is unusual
// rather than broken, and silently: a refused update is downloaded and refused
// again on every later check.
//
//   - codesign missing, or timed out;
//   - the payload is not a macOS bundle at all;
//   - the bundle carries no signature, or only the one the linker gives its
//     executable (see linkerSignedOnly). A build made locally, or by anyone
//     without signing credentials, is in one of those states by construction.
//
// A verdict is enforced only if original, when given, passes the same check.
// The update is held to the standard the installation it replaces already
// meets. Whatever leaves the original failing — a signing state not recognised
// above, a codesign that objects for reasons of its own — says nothing about
// whether the update was damaged, and enforcing it would roll back every
// update to that installation. Where the original does pass, as any bundle
// signed at build time does, ad-hoc included, a failing update is rejected.
// The original is checked only once the update has failed, so a passing update
// costs one codesign run.
func verifySignature(path, original string) (note string, err error) {
	if runtime.GOOS != "darwin" || filepath.Ext(path) != ".app" {
		return "", nil
	}
	tool, lookErr := exec.LookPath("codesign")
	if lookErr != nil {
		return "codesign unavailable, bundle signature not checked", nil
	}

	note, verdict := codesignVerify(tool, path)
	if verdict == nil {
		return note, nil
	}
	if original != "" {
		if baseNote, baseVerdict := codesignVerify(tool, original); baseVerdict != nil || baseNote != "" {
			return fmt.Sprintf("bundle failed signature verification (%v), but so does the application it replaced; not enforced", verdict), nil
		}
	}
	return "", fmt.Errorf("bundle failed signature verification: %w", verdict)
}

// codesignVerify runs codesign --verify --deep --strict over bundle. It returns
// ("", nil) for a pass, a note and a nil verdict when there is no verdict to be
// had, and a non-nil verdict when codesign rejected the bundle.
func codesignVerify(tool, bundle string) (note string, verdict error) {
	ctx, cancel := context.WithTimeout(context.Background(), codesignTimeout)
	defer cancel()

	out, runErr := exec.CommandContext(ctx, tool, "--verify", "--deep", "--strict", bundle).CombinedOutput()
	if ctx.Err() != nil {
		return "codesign timed out, bundle signature not checked", nil
	}
	if runErr == nil {
		return "", nil
	}
	if isUnsigned(out) {
		return "bundle carries no signature, integrity not checked", nil
	}
	if linkerSignedOnly(tool, bundle, out) {
		return "bundle executable carries only its linker signature, bundle resources not checked", nil
	}
	return "", errors.New(firstLine(out))
}

// isUnsigned distinguishes codesign's "there is nothing here to check" from a
// genuine rejection. The wording has been stable for many releases, and the
// consequence of misreading it in either direction is mild: an unsigned bundle
// wrongly called signed is judged against the original instead, and a corrupt
// bundle wrongly called unsigned still has to pass the structural check.
func isUnsigned(out []byte) bool {
	return strings.Contains(string(out), "not signed at all")
}

// noResources is codesign's rejection of a bundle whose executable is signed
// but whose resources were never sealed.
const noResources = "code has no resources but signature indicates they must be present"

// linkerSignedOnly reports whether a verification failure means only that the
// bundle was never signed as a bundle. On arm64 the Go linker, like Apple's,
// gives every executable it writes an ad-hoc signature (flags
// "adhoc,linker-signed"), because the kernel will not run one without; that
// signature hashes the executable alone. Placed in a .app with no signing step
// of its own — a cross-compiled or hand-assembled bundle — it reads as a
// promise of sealed resources the bundle does not have, and codesign rejects
// it with noResources rather than calling it unsigned.
//
// Only that rejection is excused, and only when codesign -dv confirms the
// signature is the linker's. codesign checks the executable's own pages
// first, so a linker-signed binary modified after linking fails with "invalid
// signature" instead and is not excused. Signing the bundle with codesign -s,
// ad-hoc or otherwise, clears the linker-signed flag and seals the resources,
// so such a bundle's verdict always stands.
func linkerSignedOnly(tool, bundle string, verifyOut []byte) bool {
	if !strings.Contains(string(verifyOut), noResources) {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), codesignTimeout)
	defer cancel()
	out, _ := exec.CommandContext(ctx, tool, "-dv", bundle).CombinedOutput()
	for _, line := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(line, "CodeDirectory ") && strings.Contains(line, "linker-signed") {
			return true
		}
	}
	return false
}

// firstLine trims codesign's output down to the line that says what is wrong,
// so a rejection reads as one line in the helper log rather than a page.
func firstLine(b []byte) string {
	s := string(b)
	for i, r := range s {
		if r == '\n' {
			return s[:i]
		}
	}
	if s == "" {
		return "no output"
	}
	return s
}
