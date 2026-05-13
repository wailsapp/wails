//go:build linux && !android && !server

package application

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type linuxAutostart struct {
	app *App
}

func newAutostartImpl(app *App) autostartImpl {
	return &linuxAutostart{app: app}
}

func (a *linuxAutostart) enable(opts AutostartOptions) error {
	if err := validateAutostartIdentifier(opts.Identifier); err != nil {
		return err
	}
	// A newline inside the executable path or any argument would break the
	// .desktop file format (line-based key=value) and could be used to inject
	// arbitrary Desktop Entry keys when Arguments are user-influenced.
	for i, arg := range opts.Arguments {
		if err := validateDesktopExecToken(arg); err != nil {
			return fmt.Errorf("autostart argument %d: %w", i, err)
		}
	}

	exe, err := resolvedExecutable()
	if err != nil {
		return err
	}
	if err := validateDesktopExecToken(exe); err != nil {
		return fmt.Errorf("autostart executable path: %w", err)
	}

	dir, err := a.autostartDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create autostart dir: %w", err)
	}

	id := opts.Identifier
	if id == "" {
		id = autostartSlug(a.app.options.Name)
	}
	path := filepath.Join(dir, id+".desktop")

	// Remove any stale .desktop file pointing at this binary under a
	// different identifier so a previous Enable() with a different
	// Identifier (or a slug derived from a renamed Options.Name) doesn't
	// leave a second entry behind.
	if existing, ferr := a.findDesktopFile(dir); ferr == nil && existing != "" && existing != path {
		_ = os.Remove(existing)
	}

	body := buildDesktopEntry(a.app.options.Name, exe, opts.Arguments)
	if err := writeFileAtomic(path, []byte(body), 0644); err != nil {
		return fmt.Errorf("write desktop file %s: %w", path, err)
	}
	return nil
}

func (a *linuxAutostart) disable() error {
	dir, err := a.autostartDir()
	if err != nil {
		return err
	}
	path, err := a.findDesktopFile(dir)
	if err != nil {
		return err
	}
	if path == "" {
		return nil
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove desktop file: %w", err)
	}
	return nil
}

func (a *linuxAutostart) status() (AutostartStatus, error) {
	dir, err := a.autostartDir()
	if err != nil {
		return AutostartStatus{}, err
	}
	path, err := a.findDesktopFile(dir)
	if err != nil {
		return AutostartStatus{}, err
	}
	if path == "" {
		return AutostartStatus{}, nil
	}
	return AutostartStatus{
		Enabled:  true,
		Path:     path,
		Strategy: AutostartStrategyXDGAutostart,
	}, nil
}

func (a *linuxAutostart) autostartDir() (string, error) {
	cfg := os.Getenv("XDG_CONFIG_HOME")
	if cfg == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("autostart: %w", err)
		}
		cfg = filepath.Join(home, ".config")
	}
	return filepath.Join(cfg, "autostart"), nil
}

// findDesktopFile looks for a .desktop file in dir whose Exec= entry points at
// the current executable. Returns empty path with no error if none found.
// This survives identifier changes between Enable() calls.
func (a *linuxAutostart) findDesktopFile(dir string) (string, error) {
	exe, err := resolvedExecutable()
	if err != nil {
		return "", err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("read autostart dir: %w", err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".desktop") {
			continue
		}
		full := filepath.Join(dir, e.Name())
		data, err := os.ReadFile(full)
		if err != nil {
			continue
		}
		if desktopExecPath(string(data)) == exe {
			return full, nil
		}
	}
	return "", nil
}

func buildDesktopEntry(appName, exe string, args []string) string {
	if appName == "" {
		appName = filepath.Base(exe)
	}
	execLine := quoteExec(exe)
	for _, a := range args {
		execLine += " " + quoteExec(a)
	}
	return fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=%s
Exec=%s
X-GNOME-Autostart-enabled=true
Hidden=false
NoDisplay=false
Terminal=false
`, escapeDesktopValue(appName), execLine)
}

// validateDesktopExecToken rejects control characters that would break the
// .desktop file format or allow Desktop Entry key injection when interpolated
// into an Exec= line. Allowed: spaces and tabs (which quoteExec handles by
// double-quoting); rejected: all other ASCII control characters including
// CR/LF.
func validateDesktopExecToken(s string) error {
	for _, r := range s {
		if r == '\t' || r == ' ' {
			continue
		}
		if r < 0x20 || r == 0x7f {
			return fmt.Errorf("control character %U not allowed in Exec field", r)
		}
	}
	return nil
}

// quoteExec escapes a single Exec field token per the freedesktop.org spec:
// reserved chars are " ` $ \ → escape with backslash; if the token contains
// any reserved or whitespace, double-quote it.
func quoteExec(s string) string {
	needQuote := false
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '"', '`', '$', '\\':
			b.WriteByte('\\')
			b.WriteRune(r)
			needQuote = true
		case ' ', '\t':
			// Newlines are rejected by validateDesktopExecToken before we
			// get here, so any whitespace remaining is safely quotable.
			b.WriteRune(r)
			needQuote = true
		default:
			b.WriteRune(r)
		}
	}
	if needQuote {
		return `"` + b.String() + `"`
	}
	return b.String()
}

// escapeDesktopValue escapes characters that are not allowed in raw Desktop
// Entry values (newlines, leading/trailing whitespace).
func escapeDesktopValue(s string) string {
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	return strings.TrimSpace(s)
}

func desktopExecPath(contents string) string {
	for _, line := range strings.Split(contents, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "Exec=") {
			continue
		}
		val := strings.TrimPrefix(line, "Exec=")
		// First token, possibly quoted.
		val = strings.TrimSpace(val)
		if strings.HasPrefix(val, `"`) {
			end := strings.Index(val[1:], `"`)
			if end < 0 {
				return ""
			}
			return unescapeDesktopToken(val[1 : 1+end])
		}
		if i := strings.IndexAny(val, " \t"); i >= 0 {
			return val[:i]
		}
		return val
	}
	return ""
}

func unescapeDesktopToken(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			b.WriteByte(s[i+1])
			i++
			continue
		}
		b.WriteByte(s[i])
	}
	return b.String()
}
