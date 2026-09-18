package application

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// resolvedExecutable returns os.Executable() after resolving symlinks so
// registrations don't break when the binary is installed via a symlink farm
// (Homebrew, Scoop). Falls back to the unresolved path if EvalSymlinks fails.
func resolvedExecutable() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("autostart: get executable path: %w", err)
	}
	if resolved, err := filepath.EvalSymlinks(exe); err == nil {
		return resolved, nil
	}
	return exe, nil
}

// writeFileAtomic writes data to path by way of a tempfile + rename in the
// same directory, so a partial write never leaves a half-formed plist or
// .desktop file in place.
func writeFileAtomic(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, "."+filepath.Base(path)+".*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	cleanup := func() { _ = os.Remove(tmpName) }
	// os.File.Write is documented to return an error on short writes, but we
	// double-check n == len(data) so a future change of writer type can't
	// silently rename a truncated artefact into place.
	n, err := tmp.Write(data)
	if err == nil && n != len(data) {
		err = io.ErrShortWrite
	}
	if err != nil {
		_ = tmp.Close()
		cleanup()
		return err
	}
	if err := tmp.Chmod(perm); err != nil {
		_ = tmp.Close()
		cleanup()
		return err
	}
	if err := tmp.Close(); err != nil {
		cleanup()
		return err
	}
	if err := os.Rename(tmpName, path); err != nil {
		cleanup()
		return err
	}
	return nil
}

// validateAutostartIdentifier rejects identifiers that contain characters
// that would be unsafe as a filename, registry value, or launchd Label.
func validateAutostartIdentifier(id string) error {
	if id == "" {
		return nil
	}
	if len(id) > 200 {
		return fmt.Errorf("autostart identifier too long (max 200): %q", id)
	}
	for _, r := range id {
		switch {
		case r >= 'a' && r <= 'z',
			r >= 'A' && r <= 'Z',
			r >= '0' && r <= '9',
			r == '.', r == '_', r == '-':
		default:
			return fmt.Errorf("autostart identifier contains invalid character %q (allowed: A-Za-z0-9._-)", r)
		}
	}
	return nil
}

// autostartSlug turns a free-form application name into something usable as
// the basename of a registration artefact. Empty input is rejected by the
// caller; this helper never returns an empty string for non-empty input.
func autostartSlug(name string) string {
	var b strings.Builder
	b.Grow(len(name))
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z',
			r >= '0' && r <= '9',
			r == '.', r == '_', r == '-':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r + ('a' - 'A'))
		case r == ' ', r == '\t':
			b.WriteByte('-')
		}
	}
	out := strings.Trim(b.String(), "-._")
	if out == "" {
		return "wails-app"
	}
	return out
}

// ErrAutostartNotSupported is returned when autostart is not available on
// the current platform (mobile, server builds).
var ErrAutostartNotSupported = errors.New("autostart is not supported on this platform")

// AutostartOptions configures how the application is registered to launch at login.
type AutostartOptions struct {
	// Identifier overrides the auto-derived registration ID.
	//
	//   macOS:   launchd Label / SMAppService key (reverse-DNS recommended).
	//   Windows: registry value name under HKCU\…\Run.
	//   Linux:   .desktop filename (without extension).
	//
	// If empty, a sensible default is derived: on macOS the application's
	// bundle identifier (when running from a bundle) or "wails.autostart.<slug>";
	// on Windows and Linux a slugified form of the application's Options.Name
	// (i.e. application.Options.Name from application.New).
	Identifier string

	// Arguments are appended to the executable path when launched at login.
	Arguments []string
}

// AutostartStrategy names the underlying mechanism a registration used.
//
// On macOS this distinguishes between SMAppService (bundled .app on macOS 13+)
// and a LaunchAgent plist (the fallback path). On Windows it is always
// AutostartStrategyRegistryRun and on Linux always AutostartStrategyXDGAutostart.
// Empty when AutostartStatus.Enabled is false.
type AutostartStrategy string

const (
	AutostartStrategyNone          AutostartStrategy = ""
	AutostartStrategySMAppService  AutostartStrategy = "smappservice"
	AutostartStrategyLaunchAgent   AutostartStrategy = "launchagent"
	AutostartStrategyRegistryRun   AutostartStrategy = "registry-run"
	AutostartStrategyXDGAutostart  AutostartStrategy = "xdg-autostart"
)

// AutostartStatus describes the current autostart registration.
type AutostartStatus struct {
	// Enabled reports whether a registration exists.
	Enabled bool
	// Path is the on-disk location of the registration artefact, when
	// applicable (plist path, .desktop path, registry sub-key path). Empty if
	// Enabled is false.
	Path string
	// Strategy names the mechanism that registered the application. Empty if
	// Enabled is false or the platform has only one mechanism.
	Strategy AutostartStrategy
}

type autostartImpl interface {
	enable(opts AutostartOptions) error
	disable() error
	status() (AutostartStatus, error)
}
