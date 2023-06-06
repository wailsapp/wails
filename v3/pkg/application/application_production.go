//go:build production

package application

// We use this to patch the application to production mode.
func init() {
	isDebugMode = func() bool {
		return false
	}
}
