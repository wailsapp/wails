package keygen

import (
	"fmt"
)

// KeygenError represents a Keygen service error
type KeygenError struct {
	// Code is the error code
	Code string `json:"code"`

	// Message is the human-readable error message
	Message string `json:"message"`

	// Details contains additional error information
	Details map[string]interface{} `json:"details,omitempty"`

	// Err is the underlying error
	Err error `json:"-"`
}

// Error implements the error interface
func (e *KeygenError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *KeygenError) Unwrap() error {
	return e.Err
}

// Error codes
const (
	// License errors
	ErrLicenseInvalid          = "LICENSE_INVALID"
	ErrLicenseExpired          = "LICENSE_EXPIRED"
	ErrLicenseSuspended        = "LICENSE_SUSPENDED"
	ErrLicenseOverdue          = "LICENSE_OVERDUE"
	ErrLicenseNotFound         = "LICENSE_NOT_FOUND"
	ErrLicenseKeyRequired      = "LICENSE_KEY_REQUIRED"
	ErrLicenseValidationFailed = "LICENSE_VALIDATION_FAILED"

	// Machine errors
	ErrMachineLimitReached        = "MACHINE_LIMIT_REACHED"
	ErrMachineNotActivated        = "MACHINE_NOT_ACTIVATED"
	ErrMachineAlreadyActivated    = "MACHINE_ALREADY_ACTIVATED"
	ErrMachineFingerprintMismatch = "MACHINE_FINGERPRINT_MISMATCH"
	ErrMachineActivationFailed    = "MACHINE_ACTIVATION_FAILED"
	ErrMachineDeactivationFailed  = "MACHINE_DEACTIVATION_FAILED"

	// Update errors
	ErrUpdateFailed           = "UPDATE_FAILED"
	ErrUpdateNotAvailable     = "UPDATE_NOT_AVAILABLE"
	ErrUpdateDownloadFailed   = "UPDATE_DOWNLOAD_FAILED"
	ErrUpdateInstallFailed    = "UPDATE_INSTALL_FAILED"
	ErrUpdateSignatureInvalid = "UPDATE_SIGNATURE_INVALID"
	ErrUpdateChannelInvalid   = "UPDATE_CHANNEL_INVALID"
	ErrUpdateInProgress       = "UPDATE_IN_PROGRESS"

	// Network errors
	ErrNetworkError    = "NETWORK_ERROR"
	ErrNetworkTimeout  = "NETWORK_TIMEOUT"
	ErrNetworkOffline  = "NETWORK_OFFLINE"
	ErrAPIError        = "API_ERROR"
	ErrAPIRateLimited  = "API_RATE_LIMITED"
	ErrAPIUnauthorized = "API_UNAUTHORIZED"

	// Configuration errors
	ErrConfigInvalid     = "CONFIG_INVALID"
	ErrConfigMissing     = "CONFIG_MISSING"
	ErrAccountRequired   = "ACCOUNT_REQUIRED"
	ErrProductRequired   = "PRODUCT_REQUIRED"
	ErrPublicKeyRequired = "PUBLIC_KEY_REQUIRED"
	ErrPublicKeyInvalid  = "PUBLIC_KEY_INVALID"

	// Cache errors
	ErrCacheCorrupted   = "CACHE_CORRUPTED"
	ErrCacheReadFailed  = "CACHE_READ_FAILED"
	ErrCacheWriteFailed = "CACHE_WRITE_FAILED"

	// Entitlement errors
	ErrEntitlementNotFound = "ENTITLEMENT_NOT_FOUND"
	ErrEntitlementDenied   = "ENTITLEMENT_DENIED"

	// General errors
	ErrServiceNotInitialized = "SERVICE_NOT_INITIALIZED"
	ErrOperationCancelled    = "OPERATION_CANCELLED"
	ErrUnknownError          = "UNKNOWN_ERROR"
)

// Helper functions for creating specific errors

// NewLicenseInvalidError creates a new license invalid error
func NewLicenseInvalidError(message string) *KeygenError {
	return &KeygenError{
		Code:    ErrLicenseInvalid,
		Message: message,
	}
}

// NewLicenseExpiredError creates a new license expired error
func NewLicenseExpiredError(expiresAt string) *KeygenError {
	return &KeygenError{
		Code:    ErrLicenseExpired,
		Message: fmt.Sprintf("License expired on %s", expiresAt),
		Details: map[string]interface{}{
			"expiresAt": expiresAt,
		},
	}
}

// NewMachineLimitReachedError creates a new machine limit reached error
func NewMachineLimitReachedError(current, max int) *KeygenError {
	return &KeygenError{
		Code:    ErrMachineLimitReached,
		Message: fmt.Sprintf("Machine limit reached: %d of %d machines activated", current, max),
		Details: map[string]interface{}{
			"currentMachines": current,
			"maxMachines":     max,
		},
	}
}

// NewNetworkError creates a new network error
func NewNetworkError(err error) *KeygenError {
	return &KeygenError{
		Code:    ErrNetworkError,
		Message: "Network error occurred",
		Err:     err,
	}
}

// NewConfigError creates a new configuration error
func NewConfigError(field, message string) *KeygenError {
	return &KeygenError{
		Code:    ErrConfigInvalid,
		Message: fmt.Sprintf("Invalid configuration: %s", message),
		Details: map[string]interface{}{
			"field": field,
		},
	}
}

// NewAPIError creates a new API error with status code
func NewAPIError(statusCode int, message string) *KeygenError {
	code := ErrAPIError
	switch statusCode {
	case 401:
		code = ErrAPIUnauthorized
	case 429:
		code = ErrAPIRateLimited
	}

	return &KeygenError{
		Code:    code,
		Message: message,
		Details: map[string]interface{}{
			"statusCode": statusCode,
		},
	}
}

// NewUpdateError creates a new update error
func NewUpdateError(stage, message string, err error) *KeygenError {
	return &KeygenError{
		Code:    ErrUpdateFailed,
		Message: fmt.Sprintf("Update failed during %s: %s", stage, message),
		Details: map[string]interface{}{
			"stage": stage,
		},
		Err: err,
	}
}

// NewEntitlementError creates a new entitlement error
func NewEntitlementError(feature string) *KeygenError {
	return &KeygenError{
		Code:    ErrEntitlementDenied,
		Message: fmt.Sprintf("Access denied to feature: %s", feature),
		Details: map[string]interface{}{
			"feature": feature,
		},
	}
}

// IsRetryableError checks if an error is retryable
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	keygenErr, ok := err.(*KeygenError)
	if !ok {
		return false
	}

	switch keygenErr.Code {
	case ErrNetworkError, ErrNetworkTimeout, ErrAPIRateLimited:
		return true
	default:
		return false
	}
}

// IsLicenseError checks if an error is license-related
func IsLicenseError(err error) bool {
	if err == nil {
		return false
	}

	keygenErr, ok := err.(*KeygenError)
	if !ok {
		return false
	}

	switch keygenErr.Code {
	case ErrLicenseInvalid, ErrLicenseExpired, ErrLicenseSuspended,
		ErrLicenseOverdue, ErrLicenseNotFound, ErrLicenseKeyRequired,
		ErrLicenseValidationFailed:
		return true
	default:
		return false
	}
}

// IsMachineError checks if an error is machine-related
func IsMachineError(err error) bool {
	if err == nil {
		return false
	}

	keygenErr, ok := err.(*KeygenError)
	if !ok {
		return false
	}

	switch keygenErr.Code {
	case ErrMachineLimitReached, ErrMachineNotActivated, ErrMachineAlreadyActivated,
		ErrMachineFingerprintMismatch, ErrMachineActivationFailed, ErrMachineDeactivationFailed:
		return true
	default:
		return false
	}
}

// IsUpdateError checks if an error is update-related
func IsUpdateError(err error) bool {
	if err == nil {
		return false
	}

	keygenErr, ok := err.(*KeygenError)
	if !ok {
		return false
	}

	switch keygenErr.Code {
	case ErrUpdateFailed, ErrUpdateNotAvailable, ErrUpdateDownloadFailed,
		ErrUpdateInstallFailed, ErrUpdateSignatureInvalid, ErrUpdateChannelInvalid,
		ErrUpdateInProgress:
		return true
	default:
		return false
	}
}
