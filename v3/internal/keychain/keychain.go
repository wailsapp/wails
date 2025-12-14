// Package keychain provides secure credential storage using the system keychain.
// On macOS it uses Keychain, on Windows it uses Credential Manager,
// and on Linux it uses Secret Service (via D-Bus).
package keychain

import (
	"fmt"
	"os"

	"github.com/zalando/go-keyring"
)

const (
	// ServiceName is the service identifier used for all Wails credentials
	ServiceName = "wails"

	// Credential keys - Windows
	KeyWindowsCertPassword = "windows-cert-password"

	// Credential keys - Linux
	KeyPGPPassword = "pgp-password"

	// Credential keys - macOS (cross-platform signing via Quill)
	KeyMacOSP12Password = "macos-p12-password"
	KeyNotaryKeyID      = "notary-key-id"
	KeyNotaryIssuer     = "notary-issuer"
)

// Set stores a credential in the system keychain.
// The credential is identified by a key and can be retrieved later with Get.
func Set(key, value string) error {
	err := keyring.Set(ServiceName, key, value)
	if err != nil {
		return fmt.Errorf("failed to store credential in keychain: %w", err)
	}
	return nil
}

// Get retrieves a credential from the system keychain.
// Returns the value and nil error if found, or empty string and error if not found.
// Also checks environment variables as a fallback (useful for CI).
func Get(key string) (string, error) {
	// First check environment variable (for CI/automation)
	envKey := "WAILS_" + toEnvName(key)
	if val := os.Getenv(envKey); val != "" {
		return val, nil
	}

	// Try keychain
	value, err := keyring.Get(ServiceName, key)
	if err != nil {
		if err == keyring.ErrNotFound {
			return "", fmt.Errorf("credential %q not found in keychain (set with: wails3 setup signing, or set env var %s)", key, envKey)
		}
		return "", fmt.Errorf("failed to retrieve credential from keychain: %w", err)
	}
	return value, nil
}

// Delete removes a credential from the system keychain.
func Delete(key string) error {
	err := keyring.Delete(ServiceName, key)
	if err != nil && err != keyring.ErrNotFound {
		return fmt.Errorf("failed to delete credential from keychain: %w", err)
	}
	return nil
}

// Exists checks if a credential exists in the keychain or environment.
func Exists(key string) bool {
	// Check environment variable first
	envKey := "WAILS_" + toEnvName(key)
	if os.Getenv(envKey) != "" {
		return true
	}

	// Check keychain
	_, err := keyring.Get(ServiceName, key)
	return err == nil
}

// toEnvName converts a key to an environment variable name.
// e.g., "windows-cert-password" -> "WINDOWS_CERT_PASSWORD"
func toEnvName(key string) string {
	result := make([]byte, len(key))
	for i, c := range key {
		if c == '-' {
			result[i] = '_'
		} else if c >= 'a' && c <= 'z' {
			result[i] = byte(c - 'a' + 'A')
		} else {
			result[i] = byte(c)
		}
	}
	return string(result)
}
