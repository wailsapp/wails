//go:build !debug

package app

func IsDebug() bool {
	return false
}
