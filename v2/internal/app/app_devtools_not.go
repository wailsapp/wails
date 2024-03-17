//go:build !devtools

package app

// IsDevtoolsEnabled returns true if devtools should be enabled
// Note: devtools flag is also added in debug builds
func IsDevtoolsEnabled() bool {
	return false
}
