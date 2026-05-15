package updater

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// Updater is the singleton exposed as app.Updater. It is constructed during
// application initialisation but does nothing useful until Init is called.
type Updater struct {
	mu sync.RWMutex

	// host is the bridge to the parent application — used to emit lifecycle
	// events without introducing a circular import on the application package.
	host Host

	cfg *Config

	state      State
	current    string // CurrentVersion, snapshot for State()
	pending    *Release
	resolved   string // resolved download path (after install)
	lastDigest []byte // digest computed streaming during the last successful download
}

// Host is the minimal surface the Updater needs from the application that
// owns it. The application package implements this on *App; tests stub it.
type Host interface {
	// Emit a custom event with the supplied data. Mirrors the signature of
	// (*EventManager).Emit so the application's adapter is trivially thin.
	Emit(name string, data ...any) bool
}

// New is for internal use by the application package and tests. End users
// obtain an Updater via app.Updater (or app.Updater.Init).
func New(host Host) *Updater {
	return &Updater{host: host, state: StateUnconfigured}
}

// Init configures the Updater. Returns ErrAlreadyConfigured if Init has
// already been called, or a validation error if cfg is malformed.
func (u *Updater) Init(cfg Config) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.cfg != nil {
		return ErrAlreadyConfigured
	}
	if err := cfg.validate(); err != nil {
		return err
	}
	if cfg.Platform == "" {
		cfg.Platform = runtime.GOOS
	}
	if cfg.Arch == "" {
		cfg.Arch = runtime.GOARCH
	}
	u.cfg = &cfg
	u.current = cfg.CurrentVersion
	u.state = StateIdle
	return nil
}

// State returns the Updater's current high-level lifecycle phase.
func (u *Updater) State() State {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.state
}

// CurrentVersion returns the version supplied at Init time, or "" if not configured.
func (u *Updater) CurrentVersion() string {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.current
}

// Check walks the provider chain looking for an upgrade. Returns:
//   - (rel, nil): newer release found; payload of EventUpdateAvailable.
//   - (nil, nil): caller is up to date; payload of EventNoUpdate.
//   - (nil, err): all providers errored; err wraps every failure.
//
// Fallback semantics: a provider returning (nil, nil) short-circuits to
// up-to-date — fallback exists for "primary unreachable", not "providers
// disagree about what's latest." A provider returning an error advances to
// the next one. If every provider errors, Check returns an error built from
// the chain.
func (u *Updater) Check(ctx context.Context) (*Release, error) {
	u.mu.RLock()
	cfg := u.cfg
	u.mu.RUnlock()
	if cfg == nil {
		return nil, ErrNotConfigured
	}

	u.transition(StateChecking)
	u.host.Emit(EventCheckStarted)

	req := CheckRequest{
		CurrentVersion: cfg.CurrentVersion,
		Platform:       cfg.Platform,
		Arch:           cfg.Arch,
	}

	var failures []error
	for _, p := range cfg.Providers {
		rel, err := p.Check(ctx, req)
		if err != nil {
			failures = append(failures, fmt.Errorf("%s: %w", p.Name(), err))
			continue
		}
		if rel == nil {
			u.transition(StateUpToDate)
			u.host.Emit(EventNoUpdate)
			return nil, nil
		}
		rel.Provider = p.Name()

		u.mu.Lock()
		u.pending = rel
		u.state = StateAvailable
		u.mu.Unlock()

		u.host.Emit(EventUpdateAvailable, rel)
		return rel, nil
	}

	err := joinErrors("all providers failed", failures)
	u.transition(StateError)
	u.host.Emit(EventError, ErrorInfo{Stage: StageCheck, Message: err.Error()})
	return nil, err
}

// DownloadAndInstall downloads the pending release (set by a previous Check),
// verifies it, and stages it for swap. Returns ErrNoPendingRelease if Check
// did not produce one. The actual binary swap is performed in a follow-up
// commit; v1 of this branch stages the verified file and reports
// StateReady + EventUpdateReady.
func (u *Updater) DownloadAndInstall(ctx context.Context) error {
	u.mu.RLock()
	cfg := u.cfg
	pending := u.pending
	u.mu.RUnlock()
	if cfg == nil {
		return ErrNotConfigured
	}
	if pending == nil {
		return ErrNoPendingRelease
	}

	provider, err := findProvider(cfg.Providers, pending.Provider)
	if err != nil {
		return err
	}

	tmpPath, err := u.download(ctx, provider, pending)
	if err != nil {
		u.transition(StateError)
		u.host.Emit(EventError, ErrorInfo{Stage: StageDownload, Message: err.Error(), Provider: provider.Name()})
		return err
	}

	u.transition(StateVerifying)
	u.host.Emit(EventVerifying, pending)

	if err := u.verify(tmpPath, pending); err != nil {
		_ = os.Remove(tmpPath)
		u.transition(StateError)
		u.host.Emit(EventError, ErrorInfo{Stage: StageVerify, Message: err.Error(), Provider: provider.Name()})
		return err
	}

	u.transition(StateInstalling)
	u.host.Emit(EventInstalling, pending)

	finalPath, err := finaliseDownload(tmpPath, pending.Artifact.Filename)
	if err != nil {
		_ = os.Remove(tmpPath)
		u.transition(StateError)
		u.host.Emit(EventError, ErrorInfo{Stage: StageInstall, Message: err.Error(), Provider: provider.Name()})
		return err
	}

	u.mu.Lock()
	u.resolved = finalPath
	u.state = StateReady
	u.mu.Unlock()

	u.host.Emit(EventUpdateReady, pending)
	return nil
}

// CheckAndInstall is the convenience method: Check + DownloadAndInstall in
// one call. Returns nil with no side effects if the application is already
// up to date.
func (u *Updater) CheckAndInstall(ctx context.Context) error {
	rel, err := u.Check(ctx)
	if err != nil {
		return err
	}
	if rel == nil {
		return nil
	}
	return u.DownloadAndInstall(ctx)
}

// DownloadedPath returns the on-disk path of the last successfully-installed
// (staged) update, or "" if none.
func (u *Updater) DownloadedPath() string {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.resolved
}

// --- internals ---

func (u *Updater) transition(s State) {
	u.mu.Lock()
	u.state = s
	u.mu.Unlock()
}

func findProvider(providers []Provider, name string) (Provider, error) {
	for _, p := range providers {
		if p.Name() == name {
			return p, nil
		}
	}
	return nil, fmt.Errorf("updater: provider %q is no longer registered", name)
}

func joinErrors(prefix string, errs []error) error {
	if len(errs) == 0 {
		return errors.New(prefix)
	}
	parts := make([]string, 0, len(errs))
	for _, e := range errs {
		parts = append(parts, e.Error())
	}
	return fmt.Errorf("updater: %s: %s", prefix, strings.Join(parts, "; "))
}

func finaliseDownload(tmpPath, filename string) (string, error) {
	base := filepath.Base(filename)
	if base == "" || base == "." || base == "/" {
		base = "wails-update.bin"
	}
	final := filepath.Join(filepath.Dir(tmpPath), base)
	if err := os.Rename(tmpPath, final); err != nil {
		return "", fmt.Errorf("updater: finalise: %w", err)
	}
	return final, nil
}

// errors

var (
	ErrAlreadyConfigured = errors.New("updater: Init already called")
	ErrNotConfigured     = errors.New("updater: Init has not been called")
	ErrNoPendingRelease  = errors.New("updater: no pending release (call Check first)")
)
