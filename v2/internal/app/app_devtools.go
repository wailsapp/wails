//go:build devtools

package app

// Note: devtools flag is also added in debug builds
func IsDevtoolsEnabled() bool {
	return true
}
