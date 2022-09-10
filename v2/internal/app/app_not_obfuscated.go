//go:build !obfuscated

package app

// IsObfuscated returns false if the obfuscated build tag is not set
func IsObfuscated() bool {
	return false
}
