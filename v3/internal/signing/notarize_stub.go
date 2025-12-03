//go:build !darwin

package signing

import (
	"fmt"
	"time"
)

// NotarizeOptions defines options for notarization
type NotarizeOptions struct {
	AppPath             string
	AppleID             string
	TeamID              string
	KeychainProfile     string
	AppSpecificPassword string
	Wait                bool
	Verbose             bool
}

// NotarizationStatus represents the status of a notarization submission
type NotarizationStatus struct {
	ID        string
	Status    string
	Message   string
	CreatedAt time.Time
}

// Notarize submits an application for notarization
func Notarize(options NotarizeOptions) (*NotarizationStatus, error) {
	return nil, fmt.Errorf("notarization is only available on macOS")
}

// WaitForNotarization waits for a notarization submission to complete
func WaitForNotarization(submissionID string, options NotarizeOptions) (*NotarizationStatus, error) {
	return nil, fmt.Errorf("notarization is only available on macOS")
}

// GetNotarizationLog retrieves the log for a notarization submission
func GetNotarizationLog(submissionID string, options NotarizeOptions) (string, error) {
	return "", fmt.Errorf("notarization is only available on macOS")
}

// Staple staples the notarization ticket to the application
func Staple(appPath string) error {
	return fmt.Errorf("notarization is only available on macOS")
}

// NotarizeAndStaple performs the complete notarization workflow
func NotarizeAndStaple(options NotarizeOptions) error {
	return fmt.Errorf("notarization is only available on macOS")
}

// StoreCredentials stores notarization credentials in the keychain
func StoreCredentials(profileName, appleID, teamID, password string) error {
	return fmt.Errorf("credential storage is only available on macOS")
}
