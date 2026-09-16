//go:build !darwin || ios || server

package application

// newPermissionsImpl returns nil where no native permission API is wired
// up; PermissionsManager then reports PermissionStatusUnsupported and
// ErrPermissionsUnsupported.
func newPermissionsImpl(app *App) permissionsImpl {
	return nil
}
