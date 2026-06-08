package monitor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

// DiscoveryEntry describes a running, monitorable Wails app.
type DiscoveryEntry struct {
	Name      string    `json:"name"`
	PID       int       `json:"pid"`
	Sock      string    `json:"sock"`
	StartedAt time.Time `json:"startedAt"`
}

// discoveryDir returns the directory used for sockets and discovery files,
// creating it (0700) if necessary.
func discoveryDir() (string, error) {
	base := os.Getenv("XDG_RUNTIME_DIR")
	if base == "" {
		base = os.TempDir()
	}
	dir := filepath.Join(base, "wails3")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	return dir, nil
}

// DefaultSocketPath returns the default unix socket path for the given app.
func DefaultSocketPath(appName string, pid int) string {
	dir, err := discoveryDir()
	if err != nil {
		// Best-effort fallback; Start will surface any real error.
		dir = filepath.Join(os.TempDir(), "wails3")
	}
	return filepath.Join(dir, fmt.Sprintf("%s-%d.sock", sanitize(appName), pid))
}

// WriteDiscovery writes a discovery file so external tools can find this app.
// The returned cleanup func removes it.
func WriteDiscovery(appName string, pid int, sockPath string) (cleanup func(), err error) {
	dir, err := discoveryDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, fmt.Sprintf("%s-%d.json", sanitize(appName), pid))

	entry := DiscoveryEntry{
		Name:      appName,
		PID:       pid,
		Sock:      sockPath,
		StartedAt: time.Now(),
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return nil, err
	}
	return func() { _ = os.Remove(path) }, nil
}

// ListDiscovery returns all discoverable apps, skipping (and cleaning up)
// entries whose process is no longer alive.
func ListDiscovery() ([]DiscoveryEntry, error) {
	dir, err := discoveryDir()
	if err != nil {
		return nil, err
	}
	matches, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		return nil, err
	}

	var entries []DiscoveryEntry
	for _, path := range matches {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var entry DiscoveryEntry
		if err := json.Unmarshal(data, &entry); err != nil {
			continue
		}
		if !pidAlive(entry.PID) {
			_ = os.Remove(path)
			continue
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func pidAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// Signal 0 performs error checking without actually sending a signal.
	return proc.Signal(syscall.Signal(0)) == nil
}

func sanitize(name string) string {
	out := make([]rune, 0, len(name))
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '-', r == '_':
			out = append(out, r)
		default:
			out = append(out, '_')
		}
	}
	if len(out) == 0 {
		return "wails-app"
	}
	return string(out)
}
