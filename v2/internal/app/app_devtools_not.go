//go:build !devtools

package app

func IsDevtoolsEnabled() bool {
	return false
}
