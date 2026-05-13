//go:build windows && !server

package application

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/sys/windows/registry"
)

const defaultAutostartRegistrySubKey = `Software\Microsoft\Windows\CurrentVersion\Run`

type windowsAutostart struct {
	app *App

	// registrySubKey is overridable for tests; production code reads/writes
	// HKCU\Software\Microsoft\Windows\CurrentVersion\Run.
	registrySubKey string
}

func newAutostartImpl(app *App) autostartImpl {
	return &windowsAutostart{
		app:            app,
		registrySubKey: defaultAutostartRegistrySubKey,
	}
}

func (a *windowsAutostart) enable(opts AutostartOptions) error {
	if err := validateAutostartIdentifier(opts.Identifier); err != nil {
		return err
	}
	exe, err := resolvedExecutable()
	if err != nil {
		return err
	}
	id := opts.Identifier
	if id == "" {
		id = autostartSlug(a.app.options.Name)
	}

	cmd := quoteWindowsArg(exe)
	for _, arg := range opts.Arguments {
		cmd += " " + quoteWindowsArg(arg)
	}

	key, _, err := registry.CreateKey(registry.CURRENT_USER, a.registrySubKey, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("autostart: open registry key: %w", err)
	}
	defer key.Close()

	// Remove any stale entry pointing at this binary under a different value
	// name, so a previous Enable() with a different Identifier (or a slug
	// derived from a renamed Options.Name) doesn't leave behind a second
	// autostart entry.
	if existing, _, ferr := a.find(); ferr == nil && existing != "" && existing != id {
		_ = key.DeleteValue(existing)
	}

	if err := key.SetStringValue(id, cmd); err != nil {
		return fmt.Errorf("autostart: write registry value: %w", err)
	}
	return nil
}

func (a *windowsAutostart) disable() error {
	id, _, err := a.find()
	if err != nil {
		return err
	}
	if id == "" {
		return nil
	}
	key, err := registry.OpenKey(registry.CURRENT_USER, a.registrySubKey, registry.SET_VALUE)
	if err != nil {
		if errors.Is(err, registry.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("autostart: open registry key: %w", err)
	}
	defer key.Close()
	if err := key.DeleteValue(id); err != nil && !errors.Is(err, registry.ErrNotExist) {
		return fmt.Errorf("autostart: delete registry value: %w", err)
	}
	return nil
}

func (a *windowsAutostart) status() (AutostartStatus, error) {
	id, _, err := a.find()
	if err != nil {
		return AutostartStatus{}, err
	}
	if id == "" {
		return AutostartStatus{}, nil
	}
	return AutostartStatus{
		Enabled:  true,
		Path:     `HKCU\` + a.registrySubKey + `\` + id,
		Strategy: AutostartStrategyRegistryRun,
	}, nil
}

// find returns the value name and command of the registry entry whose first
// token equals our current executable. Empty name means not registered.
func (a *windowsAutostart) find() (string, string, error) {
	exe, err := resolvedExecutable()
	if err != nil {
		return "", "", err
	}
	key, err := registry.OpenKey(registry.CURRENT_USER, a.registrySubKey, registry.QUERY_VALUE)
	if err != nil {
		if errors.Is(err, registry.ErrNotExist) {
			return "", "", nil
		}
		return "", "", fmt.Errorf("autostart: open registry key: %w", err)
	}
	defer key.Close()

	names, err := key.ReadValueNames(-1)
	if err != nil {
		return "", "", fmt.Errorf("autostart: list registry values: %w", err)
	}
	exeLower := strings.ToLower(exe)
	for _, name := range names {
		val, _, err := key.GetStringValue(name)
		if err != nil {
			continue
		}
		if strings.EqualFold(parseWindowsCommandExe(val), exeLower) {
			return name, val, nil
		}
	}
	return "", "", nil
}

// parseWindowsCommandExe returns the first token of a Windows command line,
// honouring surrounding double quotes for paths with spaces. Returned in
// lowercase for case-insensitive comparison.
func parseWindowsCommandExe(cmd string) string {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return ""
	}
	if cmd[0] == '"' {
		end := strings.IndexByte(cmd[1:], '"')
		if end < 0 {
			return strings.ToLower(cmd[1:])
		}
		return strings.ToLower(cmd[1 : 1+end])
	}
	if i := strings.IndexAny(cmd, " \t"); i >= 0 {
		return strings.ToLower(cmd[:i])
	}
	return strings.ToLower(cmd)
}

// quoteWindowsArg wraps an argument in double quotes when it contains
// whitespace or quotes, and escapes embedded quotes. Backslashes preceding a
// quote are doubled per CommandLineToArgvW rules.
func quoteWindowsArg(s string) string {
	if s != "" && !strings.ContainsAny(s, ` "	`) {
		return s
	}
	var b strings.Builder
	b.WriteByte('"')
	backslashes := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case '\\':
			backslashes++
		case '"':
			// Per CommandLineToArgvW: a literal quote preceded by N
			// backslashes must be encoded as (2N+1) backslashes + quote.
			// Earlier versions emitted only (N+1), which made the parser
			// lose the quote (the 2 backslashes toggled quoted state).
			for j := 0; j < 2*backslashes; j++ {
				b.WriteByte('\\')
			}
			b.WriteByte('\\')
			b.WriteByte('"')
			backslashes = 0
		default:
			for j := 0; j < backslashes; j++ {
				b.WriteByte('\\')
			}
			backslashes = 0
			b.WriteByte(c)
		}
	}
	for j := 0; j < backslashes; j++ {
		b.WriteByte('\\')
		b.WriteByte('\\')
	}
	b.WriteByte('"')
	return b.String()
}
