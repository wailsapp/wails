package updater

import (
	"context"
	"sync"
)

// User-action event names emitted by the default window and any custom
// template that follows the same contract. The Updater listens to these to
// drive the flow without ever calling into the window directly.
const (
	EventWindowReady = "updater:window:ready"
	EventUserInstall = "updater:user:install"
	EventUserSkip    = "updater:user:skip"
	EventUserRemind  = "updater:user:remind"
	EventUserCancel  = "updater:user:cancel"
	EventUserRestart = "updater:user:restart"
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
	)
	return sess
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

// sessionMu guards the active session pointer separately from the main
// state mutex to avoid lock-ordering hazards when a user-action callback
// (fired by Host.OnEvent on a separate goroutine) needs to close the
// session while the main flow holds u.mu.
type sessionMu struct {
	sync.Mutex
}
