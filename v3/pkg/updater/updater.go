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
	"time"
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
	stagingDir string // os.MkdirTemp parent of resolved, removed on Restart / re-Check
	lastDigest []byte // digest computed streaming during the last successful download
	skipped    string // version recorded by SkipVersion / the default window Skip button

	dlMu    sync.Mutex     // serialises concurrent DownloadAndInstall calls
	sessMu  sync.Mutex     // protects session pointer separately from u.mu
	session *windowSession // current window session, if any

	periodicCtx    context.Context
	periodicCancel context.CancelFunc
	periodicDone   chan struct{}
}

// Host is the minimal surface the Updater needs from the application that
// owns it. The application package implements this via an adapter; tests
// stub it.
type Host interface {
	// Emit a custom event with the supplied data.
	Emit(name string, data ...any) bool

	// OnEvent registers a callback for a custom event. The returned function
	// removes the listener.
	OnEvent(name string, callback func(payload any)) func()

	// OpenWindow creates and shows an update window. Implementations build
	// a real Wails webview window from opts; tests return a recorder.
	OpenWindow(opts WindowOptions) WindowHandle

	// Quit asks the host application to begin its normal shutdown sequence.
	// Called by Restart after the helper has been spawned so the helper's
	// "wait for parent to exit" step actually completes.
	Quit()
}

// WindowHandle is the minimal API the Updater drives once a window is open.
// The application adapter satisfies this around a *WebviewWindow.
//
// SetHTML is deliberately omitted — loading HTML after construction puts it
// on about:blank where the Wails runtime isn't injected (so JS event emits
// no-op on webkit2gtk and in some WebView2 configurations). Templates are
// passed via WindowOptions.InitialHTML at construction time instead.
type WindowHandle interface {
	EmitEvent(name string, data ...any) bool
	Show()
	Close()
}

// WindowSizer is an optional capability a WindowHandle may implement to
// allow the Updater to resize the window in response to state changes.
// The default window template uses this to shrink the Up-to-Date state
// to a compact card; the framework adapter for *application.WebviewWindow
// implements it transparently. BYO windows can opt in by adding the
// SetSize method themselves; if they don't, the Updater silently skips
// the resize and the window stays whatever size the host opened it at.
type WindowSizer interface {
	SetSize(width, height int)
}

// WindowOptions describes the chrome and starting content for a window the
// Updater asks the host to open. Maps to (a subset of)
// application.WebviewWindowOptions on the host side.
type WindowOptions struct {
	Title         string
	Width, Height int
	Frameless     bool
	AlwaysOnTop   bool
	DisableResize bool
	InitialHTML   string
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
	if cfg.CheckInterval > 0 {
		u.periodicCtx, u.periodicCancel = context.WithCancel(context.Background())
		u.periodicDone = make(chan struct{})
		go u.periodicCheckLoop(cfg.CheckInterval)
	}
	return nil
}

// periodicCheckLoop polls the provider chain on the configured interval.
// Each tick runs CheckAndInstall so the default UI (or the user's headless
// subscribers) sees the found update. Ticks that arrive while another flow
// is already in progress are dropped — concurrent state machines are not
// supported and CheckAndInstall's own session-setup lock would defer the
// second call anyway.
func (u *Updater) periodicCheckLoop(d time.Duration) {
	defer close(u.periodicDone)
	t := time.NewTicker(d)
	defer t.Stop()
	for {
		select {
		case <-u.periodicCtx.Done():
			return
		case <-t.C:
			s := u.State()
			if s == StateChecking || s == StateDownloading || s == StateVerifying || s == StateInstalling {
				continue
			}
			_ = u.CheckAndInstall(u.periodicCtx)
		}
	}
}

// StopPeriodicCheck cancels the timer started by Init when
// Config.CheckInterval > 0 and blocks until the polling goroutine has
// returned (so callers can safely inspect provider state afterward). Safe
// to call when no periodic check was configured.
func (u *Updater) StopPeriodicCheck() {
	u.mu.Lock()
	cancel := u.periodicCancel
	done := u.periodicDone
	u.periodicCancel = nil
	u.periodicDone = nil
	u.mu.Unlock()
	if cancel != nil {
		cancel()
	}
	if done != nil {
		<-done
	}
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
		if u.shouldSkip(rel.Version) {
			// User has explicitly skipped this version — surface as up-to-date.
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
// did not produce one. Returns ErrDownloadInProgress if another download is
// already running. The actual binary swap is performed in a follow-up
// commit; v1 of this branch stages the verified file and reports
// StateReady + EventUpdateReady.
func (u *Updater) DownloadAndInstall(ctx context.Context) error {
	if !u.dlMu.TryLock() {
		return ErrDownloadInProgress
	}
	defer u.dlMu.Unlock()

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

	// Drop any stale staging dir from a previous DownloadAndInstall the
	// caller didn't follow up with Restart.
	u.discardStaging()

	tmpPath, tmpDir, err := u.download(ctx, provider, pending)
	if err != nil {
		u.transition(StateError)
		u.host.Emit(EventError, ErrorInfo{Stage: StageDownload, Message: err.Error(), Provider: provider.Name()})
		return err
	}

	u.transition(StateVerifying)
	u.host.Emit(EventVerifying, pending)

	if err := u.verify(tmpPath, pending); err != nil {
		_ = os.RemoveAll(tmpDir)
		u.transition(StateError)
		u.host.Emit(EventError, ErrorInfo{Stage: StageVerify, Message: err.Error(), Provider: provider.Name()})
		return err
	}

	u.transition(StateInstalling)
	u.host.Emit(EventInstalling, pending)

	finalPath, err := finaliseDownload(tmpPath, pending.Artifact.Filename)
	if err != nil {
		_ = os.RemoveAll(tmpDir)
		u.transition(StateError)
		u.host.Emit(EventError, ErrorInfo{Stage: StageInstall, Message: err.Error(), Provider: provider.Name()})
		return err
	}

	// If the artifact is an archive (.zip / .tar.gz), unpack it now so the
	// helper has a real binary or .app bundle to rename into place. Most
	// macOS distributions ship the .app inside a .zip; without this step the
	// helper would replace /Applications/MyApp.app (a directory) with the
	// downloaded .zip (a file). Non-archive artifacts pass through unchanged.
	finalPath, _, err = maybeExtractInto(finalPath)
	if err != nil {
		_ = os.RemoveAll(tmpDir)
		u.transition(StateError)
		u.host.Emit(EventError, ErrorInfo{Stage: StageInstall, Message: err.Error(), Provider: provider.Name()})
		return err
	}

	u.mu.Lock()
	u.resolved = finalPath
	u.stagingDir = tmpDir
	u.state = StateReady
	u.mu.Unlock()

	u.host.Emit(EventUpdateReady, pending)
	return nil
}

// CheckAndInstall is the convenience method: it opens the update window
// (unless Config.Window == WindowNone) and runs Check + DownloadAndInstall.
// Returns nil with no side effects if the application is already up to date.
//
// The window stays open for the duration of the flow AND across the
// "up-to-date" / error terminal states — the user dismisses it via the
// Close button. Opening + immediately closing on the no-update branch
// produced a perceptible flicker on every check; keeping the window up so
// the "You're up to date" panel actually renders matches what users expect
// from system-style updaters.
//
// Apps that want silent background polling should use Config.Window =
// updater.WindowNone (no window ever opens) or invoke Check() directly and
// subscribe to EventNoUpdate / EventUpdateAvailable themselves.
func (u *Updater) CheckAndInstall(ctx context.Context) error {
	// Tear down any session from a previous CheckAndInstall before opening a
	// fresh one — otherwise its listeners and window leak and stale callbacks
	// still fire on user-action events. Hold sessMu across the close → open →
	// assign sequence so concurrent callers can't orphan each other's
	// listeners (caller A opens, caller B closes A then opens its own, and
	// A's assignment lands on top, leaking B's listeners).
	u.sessMu.Lock()
	if u.session != nil {
		u.session.close()
		u.session = nil
	}
	sess := u.openSession(ctx)
	u.session = sess
	u.sessMu.Unlock()

	rel, err := u.Check(ctx)
	if err != nil {
		// Window stays open showing the error; the user can dismiss it via
		// the Cancel button which fires updater:user:cancel.
		return err
	}
	if rel == nil {
		// "Up to date" — leave the window open showing the up-to-date panel
		// (window.html's onNoUpdate handler renders "You're Up to Date" and
		// the current version). Closing here caused a flicker on every
		// check that found nothing.
		return nil
	}
	return u.DownloadAndInstall(ctx)
}

// Restart performs the full restart-into-the-new-version dance: it spawns a
// helper-mode child (the same binary with sentinel env vars set) and then
// asks the host application to begin its shutdown sequence via Host.Quit.
// Once the running process exits, the helper performs the binary swap and
// relaunches the (now-replaced) application.
//
// Returns ErrNotReady if DownloadAndInstall has not produced an installed
// artifact yet. If the helper spawn fails the error is surfaced and Quit is
// not called — the caller's process stays alive on the old binary.
//
// On success Restart returns once the helper has started and Quit has been
// dispatched. The caller's process will exit asynchronously as the host's
// normal shutdown unwinds.
func (u *Updater) Restart(_ context.Context) error {
	u.mu.RLock()
	staged := u.resolved
	u.mu.RUnlock()
	if staged == "" {
		return ErrNotReady
	}

	self, err := selfExecutable()
	if err != nil {
		return fmt.Errorf("updater: resolve self: %w", err)
	}

	target := bundleTarget(self)
	// Include PID so concurrent helpers (e.g. test runs, multiple installed
	// Wails apps updating at the same time) don't truncate each other's logs.
	logPath := filepath.Join(os.TempDir(), fmt.Sprintf("wails-update-%d.log", os.Getpid()))
	env := append(os.Environ(),
		envHelperMode+"=1",
		envHelperTarget+"="+target,
		envHelperNew+"="+staged,
		envHelperPID+"="+itoa(os.Getpid()),
		envHelperLog+"="+logPath,
	)

	cmd := newDetachedCommand(self)
	cmd.Env = env
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("updater: spawn helper: %w", err)
	}
	// Helper is detached and now blocking on waitForPID(os.Getpid()). Hand
	// off to the host's shutdown sequence so the wait completes and the
	// swap proceeds.
	u.host.Quit()
	return nil
}

// DownloadedPath returns the on-disk path of the last successfully-installed
// (staged) update, or "" if none.
func (u *Updater) DownloadedPath() string {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.resolved
}

// --- internals ---

// Window-sizing constants for the built-in template:
//
//  * upToDateWidth/Height — the small "compact" card used for Checking,
//    Up-to-Date, and Error states. The default window opens at this size,
//    so the most common flow (Check → Up-to-Date) involves zero visible
//    resize: the window is already its target size the moment it appears.
//
//  * availableWidth/Height — the larger "full-flow" card with room for
//    Markdown-rendered release notes, the progress bar, and the Restart
//    & Apply primary action. The Updater grows the window into this size
//    via WindowSizer when state transitions to Available / Downloading /
//    Verifying / Installing / Ready, then shrinks back if a fresh check
//    later returns Up-to-Date.
const (
	upToDateWidth     = 348
	upToDateHeight    = 161
	availableWidth    = 520
	availableHeight   = 540
)

// statesNeedingFullSize lists the states whose layout requires the larger
// window (notes panel, progress bar, or per-button row). Any state not in
// this set fits the compact upToDateWidth×upToDateHeight card.
func stateWantsFullSize(s State) bool {
	switch s {
	case StateAvailable, StateDownloading, StateVerifying, StateInstalling, StateReady:
		return true
	}
	return false
}

func (u *Updater) transition(s State) {
	u.mu.Lock()
	u.state = s
	u.mu.Unlock()
	// Resize the default window for states that need more (or less) room.
	// The window opens at the compact size and grows when it has to; if a
	// later transition takes us back to a compact-sized state the window
	// shrinks again. Handles that don't implement WindowSizer (e.g. BYO
	// windows whose owners didn't add SetSize) are silently skipped.
	u.sessMu.Lock()
	sess := u.session
	u.sessMu.Unlock()
	if sess == nil {
		return
	}
	sizer, ok := sess.handle.(WindowSizer)
	if !ok {
		return
	}
	if stateWantsFullSize(s) {
		sizer.SetSize(availableWidth, availableHeight)
	} else {
		sizer.SetSize(upToDateWidth, upToDateHeight)
	}
}

// discardStaging removes any temp directory the previous DownloadAndInstall
// left behind. Called before a new download begins and on Check when an old
// pending release becomes stale. The helper process is responsible for
// cleaning up its own staging dir post-swap; this is for the cases the
// helper never starts.
func (u *Updater) discardStaging() {
	u.mu.Lock()
	dir := u.stagingDir
	u.stagingDir = ""
	u.resolved = ""
	u.mu.Unlock()
	if dir != "" {
		_ = os.RemoveAll(dir)
	}
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
	// filepath.Base normalises away directory components, but it still passes
	// through "." and "..", and on Windows ":" / drive prefixes get reduced to
	// the suffix. Neutralise the ones that would otherwise resolve outside
	// the staging dir or that filepath.Join would not handle sanely.
	if base == "" || base == "." || base == ".." || base == "/" || strings.ContainsAny(base, `\/`) {
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
	ErrAlreadyConfigured  = errors.New("updater: Init already called")
	ErrNotConfigured      = errors.New("updater: Init has not been called")
	ErrNoPendingRelease   = errors.New("updater: no pending release (call Check first)")
	ErrDownloadInProgress = errors.New("updater: download already in progress")
)
