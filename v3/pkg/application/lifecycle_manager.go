package application

import "sync"

// lifecycleImpl is the platform side of LifecycleManager.
type lifecycleImpl interface {
	// suddenTerminationDefault reports whether the process starts with sudden
	// termination enabled (NSSupportsSuddenTermination in Info.plist).
	suddenTerminationDefault() bool
	// disableTermination and enableTermination bracket one hold. They are
	// counted by the platform, so each disable needs exactly one enable.
	disableTermination(reason string)
	enableTermination(reason string)
	setSuddenTerminationEnabled(enabled bool)
}

// LifecycleManager controls process termination behaviour such as sudden and automatic termination.
//
// macOS has two related mechanisms:
//
//   - Sudden termination lets the system kill the process with SIGKILL at
//     logout or shutdown instead of asking it to quit, which makes those
//     operations instant. It is opt-in through the NSSupportsSuddenTermination
//     Info.plist key (set it to true) and can be toggled at runtime with
//     SetSuddenTerminationEnabled. An app that opts in must never have
//     unsaved state while sudden termination is enabled.
//   - Automatic termination lets the system quit an idle app whose windows
//     are all closed to reclaim resources; it is opt-in through
//     NSSupportsAutomaticTermination.
//
// HoldTermination suspends both for the duration of a critical section such
// as saving a file or finishing an upload.
//
// Platform support is documented on each method. Methods without support on
// the current platform are safe no-ops or return a documented error.
type LifecycleManager struct {
	app  *App
	impl lifecycleImpl

	lock          sync.Mutex
	holds         int
	suddenKnown   bool
	suddenEnabled bool
}

func newLifecycleManager(app *App) *LifecycleManager {
	return &LifecycleManager{app: app, impl: newLifecycleImpl(app)}
}

// HoldTermination prevents the system from terminating the process without
// asking until the returned release func is called. Holds are reference
// counted: termination is allowed again only when every hold has been
// released. Calling release more than once is harmless. reason is a short
// description for diagnostics ("Saving document").
//
// macOS: NSProcessInfo disableSuddenTermination and
// disableAutomaticTermination:, balanced on release. This is the runtime
// equivalent of the NSSupportsSuddenTermination Info.plist key being false
// for the duration of the hold.
//
// Other platforms: no-op release.
func (lm *LifecycleManager) HoldTermination(reason string) (release func()) {
	if lm.impl == nil {
		return func() {}
	}
	if reason == "" {
		reason = "Wails application activity"
	}
	lm.lock.Lock()
	lm.holds++
	lm.lock.Unlock()
	lm.impl.disableTermination(reason)

	var once sync.Once
	return func() {
		once.Do(func() {
			lm.impl.enableTermination(reason)
			lm.lock.Lock()
			lm.holds--
			lm.lock.Unlock()
		})
	}
}

// TerminationHoldCount returns the number of HoldTermination holds not yet
// released.
func (lm *LifecycleManager) TerminationHoldCount() int {
	lm.lock.Lock()
	defer lm.lock.Unlock()
	return lm.holds
}

// SetSuddenTerminationEnabled turns sudden termination on or off for the
// process, independently of any HoldTermination holds. The call is
// idempotent: enabling twice performs one native enable, so the platform's
// enable/disable counter stays balanced. The starting state comes from the
// NSSupportsSuddenTermination Info.plist key (false when absent).
//
// macOS: NSProcessInfo enableSuddenTermination / disableSuddenTermination.
//
// Other platforms: no-op.
func (lm *LifecycleManager) SetSuddenTerminationEnabled(enabled bool) {
	if lm.impl == nil {
		return
	}
	lm.lock.Lock()
	defer lm.lock.Unlock()
	if !lm.suddenKnown {
		lm.suddenKnown = true
		lm.suddenEnabled = lm.impl.suddenTerminationDefault()
	}
	if lm.suddenEnabled == enabled {
		return
	}
	lm.suddenEnabled = enabled
	lm.impl.setSuddenTerminationEnabled(enabled)
}

// SuddenTerminationEnabled reports the state last established through
// SetSuddenTerminationEnabled, or the Info.plist default before any call.
// It does not account for active HoldTermination holds.
func (lm *LifecycleManager) SuddenTerminationEnabled() bool {
	if lm.impl == nil {
		return false
	}
	lm.lock.Lock()
	defer lm.lock.Unlock()
	if !lm.suddenKnown {
		return lm.impl.suddenTerminationDefault()
	}
	return lm.suddenEnabled
}
