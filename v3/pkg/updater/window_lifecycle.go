package updater

import "context"

// User-action event names emitted by the default window and any custom
// template that follows the same contract. The Updater listens to these to
// drive the flow without ever calling into the window directly.
const (
	EventWindowReady = "wails:updater:window:ready"
	EventUserInstall = "wails:updater:user:install"
	EventUserSkip    = "wails:updater:user:skip"
	EventUserRemind  = "wails:updater:user:remind"
	EventUserCancel  = "wails:updater:user:cancel"
	EventUserRestart = "wails:updater:user:restart"
)

// windowSession captures the window the Updater opened for the current
// CheckAndInstall flow, plus all the event-subscription cancel funcs that
// need to be invoked when the session ends.
type windowSession struct {
	mode   windowMode
	handle WindowHandle
	cancel []func()
}

func (s *windowSession) close() {
	if s == nil {
		return
	}
	for _, c := range s.cancel {
		if c != nil {
			c()
		}
	}
	s.cancel = nil
	if s.handle != nil {
		s.handle.Close()
		s.handle = nil
	}
}

// openSession opens (or attaches to) a window for the current flow and
// wires up the user-action listeners. The returned session must be Closed
// when the flow ends (success, error, or user dismiss).
func (u *Updater) openSession(ctx context.Context) *windowSession {
	u.mu.RLock()
	cfg := u.cfg
	u.mu.RUnlock()
	if cfg == nil {
		return &windowSession{mode: windowModeNone}
	}

	mode, bw, byo := classifyWindowOption(cfg.Window)
	sess := &windowSession{mode: mode}

	switch mode {
	case windowModeNone:
		// No window. Listeners still wire up so user-action events fired
		// from a custom-built UI continue to work.
	case windowModeBYO:
		sess.handle = byo
		// User-managed window — we don't show/hide it. We do load the same
		// default HTML so the bindings work out of the box, but only when
		// the user has not already populated it themselves. We can't detect
		// that, so we leave the contents alone.
	case windowModeBuiltin:
		opts := composeWindowOptions(bw)
		opts.InitialHTML = composeHTML(bw)
		sess.handle = u.host.OpenWindow(opts)
	}

	sess.cancel = append(sess.cancel,
		u.host.OnEvent(EventUserInstall, func(any) {
			go func() { _ = u.DownloadAndInstall(ctx) }()
		}),
		u.host.OnEvent(EventUserCancel, func(any) {
			u.closeWindow()
		}),
		u.host.OnEvent(EventUserSkip, func(any) {
			u.handleSkip()
		}),
		u.host.OnEvent(EventUserRemind, func(any) {
			u.closeWindow()
		}),
		u.host.OnEvent(EventUserRestart, func(any) {
			go func() { _ = u.Restart(ctx) }()
		}),
		// The default window template fires updater:window:ready on load to
		// rehydrate from current state — e.g. when a periodic-check timer
		// already advanced the flow before the window opened. Re-emit the
		// state-appropriate lifecycle event so the same handlers that drive
		// live updates also drive initial paint.
		u.host.OnEvent(EventWindowReady, func(any) {
			u.replayStateSnapshot()
		}),
	)
	return sess
}

// replayStateSnapshot re-emits the lifecycle event corresponding to the
// current state. The default window subscribes to these events for normal
// updates; replaying the latest one whenever the window asks for a snapshot
// is enough for it to render correctly on (re)open.
//
// EventMeta is emitted first so the page has the host-side context
// (currentVersion, skipped version) before the state-specific event lands —
// the default template's renderSubtitle uses currentVersion to draw the
// "from" version in the Update Available pill and the "v1.2.3 · This is
// the latest version" pill in the Up-to-Date state.
func (u *Updater) replayStateSnapshot() {
	u.mu.RLock()
	state := u.state
	pending := u.pending
	currentVersion := ""
	if u.cfg != nil {
		currentVersion = u.cfg.CurrentVersion
	}
	skipped := u.skipped
	u.mu.RUnlock()

	u.host.Emit(EventMeta, Meta{
		CurrentVersion: currentVersion,
		SkippedVersion: skipped,
	})

	switch state {
	case StateChecking:
		u.host.Emit(EventCheckStarted)
	case StateAvailable:
		u.host.Emit(EventUpdateAvailable, pending)
	case StateDownloading:
		u.host.Emit(EventDownloadStarted, pending)
	case StateVerifying:
		u.host.Emit(EventVerifying, pending)
	case StateInstalling:
		u.host.Emit(EventInstalling, pending)
	case StateReady:
		u.host.Emit(EventUpdateReady, pending)
	case StateUpToDate:
		u.host.Emit(EventNoUpdate)
	}
	// StateIdle / StateUnconfigured / StateError: nothing useful to replay.
}

func (u *Updater) closeWindow() {
	u.sessMu.Lock()
	sess := u.session
	u.session = nil
	u.sessMu.Unlock()
	sess.close()
}

func (u *Updater) handleSkip() {
	u.mu.Lock()
	if u.cfg != nil && u.pending != nil {
		u.skipped = u.pending.Version
	}
	u.mu.Unlock()
	u.closeWindow()
}

// SkipVersion records the supplied version as skipped — subsequent Checks
// will treat that version as "no update." Used by the default window's Skip
// This Version button and by callers driving headless flows.
func (u *Updater) SkipVersion(v string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.skipped = v
}

// SkippedVersion returns the currently-skipped version, if any.
func (u *Updater) SkippedVersion() string {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.skipped
}

// shouldSkip returns whether the supplied release should be ignored because
// its version has been marked skipped.
func (u *Updater) shouldSkip(version string) bool {
	if version == "" {
		return false
	}
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.skipped != "" && u.skipped == version
}
