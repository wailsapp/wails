package application

// AutostartManager provides cross-platform control over whether the
// application launches when the user logs in.
//
// Registration takes effect on the next login, not immediately.
//
// Platform behaviour:
//
//   - macOS 13+ (bundled .app):  SMAppService.mainAppService — works for
//     sandboxed and Mac-App-Store apps, no TCC automation prompt.
//   - macOS (older or unbundled): a LaunchAgent plist is written to
//     ~/Library/LaunchAgents/.
//   - Windows: a value is added under
//     HKCU\Software\Microsoft\Windows\CurrentVersion\Run.
//   - Linux: an .desktop file is written to $XDG_CONFIG_HOME/autostart/
//     (defaulting to ~/.config/autostart/).
//   - Android / iOS / server builds: ErrAutostartNotSupported.
type AutostartManager struct {
	app  *App
	impl autostartImpl
}

func newAutostartManager(app *App) *AutostartManager {
	return &AutostartManager{
		app:  app,
		impl: newAutostartImpl(app),
	}
}

// Enable registers the application to launch at user login using default
// options. Calling Enable repeatedly is safe; the registration is overwritten.
func (am *AutostartManager) Enable() error {
	return am.impl.enable(AutostartOptions{})
}

// EnableWithOptions registers the application with the given options.
// See AutostartOptions for the meaning of each field.
func (am *AutostartManager) EnableWithOptions(opts AutostartOptions) error {
	return am.impl.enable(opts)
}

// Disable removes the autostart registration. Returns nil if the application
// was not registered.
func (am *AutostartManager) Disable() error {
	return am.impl.disable()
}

// IsEnabled reports whether the application is currently registered to launch
// at login. It does not verify that the registered executable path still
// points at the running binary; use Status for that.
func (am *AutostartManager) IsEnabled() (bool, error) {
	st, err := am.impl.status()
	if err != nil {
		return false, err
	}
	return st.Enabled, nil
}

// Status returns the full registration state, including the path of the
// on-disk artefact and the platform mechanism used.
func (am *AutostartManager) Status() (AutostartStatus, error) {
	return am.impl.status()
}
