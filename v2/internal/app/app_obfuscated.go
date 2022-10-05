//go:build obfuscated

package app

// IsObfuscated returns true if the obfuscated build tag is set
func IsObfuscated() bool {
	return true
}
