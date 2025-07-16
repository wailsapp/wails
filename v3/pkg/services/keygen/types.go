package keygen

import (
	"time"

	"github.com/keygen-sh/keygen-go/v3"
)

// Config contains configuration options for the Keygen service
type Config struct {
	// Account is the Keygen account ID
	Account string `json:"account"`

	// Product is the Keygen product ID
	Product string `json:"product"`

	// LicenseKey is the user's license key
	LicenseKey string `json:"licenseKey"`

	// PublicKey is the Ed25519 public key for signature verification
	PublicKey string `json:"publicKey"`

	// CacheDir is the directory for storing offline license data
	CacheDir string `json:"cacheDir,omitempty"`

	// AutoCheck enables automatic update checking
	AutoCheck bool `json:"autoCheck"`

	// CheckInterval is the duration between automatic update checks
	CheckInterval time.Duration `json:"checkInterval,omitempty"`

	// UpdateChannel specifies the release channel (stable, beta, etc.)
	UpdateChannel string `json:"updateChannel,omitempty"`

	// Environment specifies the Keygen environment (production, staging, etc.)
	Environment string `json:"environment,omitempty"`
}

// LicenseState tracks the current state of the license
type LicenseState struct {
	// Valid indicates if the license is currently valid
	Valid bool `json:"valid"`

	// License contains the Keygen license object
	License *keygen.License `json:"license,omitempty"`

	// LastChecked is the timestamp of the last validation
	LastChecked time.Time `json:"lastChecked"`

	// OfflineMode indicates if the license is being validated offline
	OfflineMode bool `json:"offlineMode"`

	// Entitlements contains feature flags and permissions
	Entitlements map[string]interface{} `json:"entitlements,omitempty"`

	// ExpiresAt is when the license expires (nil if no expiration)
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`

	// Machine contains the current machine activation
	Machine *keygen.Machine `json:"machine,omitempty"`
}

// LicenseValidationResult represents the result of a license validation
type LicenseValidationResult struct {
	// Valid indicates if the license is valid
	Valid bool `json:"valid"`

	// Message provides a human-readable validation message
	Message string `json:"message"`

	// Code is the error/status code if validation failed
	Code string `json:"code,omitempty"`

	// State contains the current license state
	State *LicenseState `json:"state,omitempty"`

	// RequiresActivation indicates if machine activation is needed
	RequiresActivation bool `json:"requiresActivation"`
}

// MachineActivationResult represents the result of a machine activation
type MachineActivationResult struct {
	// Success indicates if the activation was successful
	Success bool `json:"success"`

	// Message provides a human-readable message
	Message string `json:"message"`

	// Machine contains the activated machine details
	Machine *keygen.Machine `json:"machine,omitempty"`

	// Fingerprint is the machine's unique identifier
	Fingerprint string `json:"fingerprint"`
}

// LicenseInfo provides detailed license information
type LicenseInfo struct {
	// Key is the license key
	Key string `json:"key"`

	// Name is the licensee name
	Name string `json:"name,omitempty"`

	// Email is the licensee email
	Email string `json:"email,omitempty"`

	// Company is the licensee company
	Company string `json:"company,omitempty"`

	// Status is the current license status
	Status string `json:"status"`

	// ExpiresAt is when the license expires
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`

	// CreatedAt is when the license was created
	CreatedAt time.Time `json:"createdAt"`

	// LastValidatedAt is when the license was last validated
	LastValidatedAt *time.Time `json:"lastValidatedAt,omitempty"`

	// Entitlements contains feature flags
	Entitlements map[string]interface{} `json:"entitlements,omitempty"`

	// Metadata contains custom license metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// MaxMachines is the maximum number of allowed machines
	MaxMachines *int `json:"maxMachines,omitempty"`

	// MachineCount is the current number of activated machines
	MachineCount int `json:"machineCount"`
}

// UpdateInfo contains information about available updates
type UpdateInfo struct {
	// Available indicates if an update is available
	Available bool `json:"available"`

	// CurrentVersion is the currently installed version
	CurrentVersion string `json:"currentVersion"`

	// LatestVersion is the latest available version
	LatestVersion string `json:"latestVersion"`

	// ReleaseID is the Keygen release ID
	ReleaseID string `json:"releaseId,omitempty"`

	// ReleaseName is the human-readable release name
	ReleaseName string `json:"releaseName,omitempty"`

	// ReleaseNotes contains the release notes/changelog
	ReleaseNotes string `json:"releaseNotes,omitempty"`

	// PublishedAt is when the release was published
	PublishedAt *time.Time `json:"publishedAt,omitempty"`

	// Critical indicates if this is a critical update
	Critical bool `json:"critical"`

	// Size is the download size in bytes
	Size int64 `json:"size,omitempty"`

	// Channel is the release channel (stable, beta, etc.)
	Channel string `json:"channel,omitempty"`

	// Artifacts contains platform-specific download artifacts
	Artifacts []ReleaseArtifact `json:"artifacts,omitempty"`
}

// ReleaseArtifact represents a downloadable release artifact
type ReleaseArtifact struct {
	// ID is the artifact ID
	ID string `json:"id"`

	// Platform is the target platform (darwin, windows, linux)
	Platform string `json:"platform"`

	// Arch is the target architecture (amd64, arm64, etc.)
	Arch string `json:"arch"`

	// Filename is the artifact filename
	Filename string `json:"filename"`

	// Size is the file size in bytes
	Size int64 `json:"size"`

	// Checksum is the file checksum
	Checksum string `json:"checksum"`

	// SignatureURL is the URL for the signature file
	SignatureURL string `json:"signatureUrl,omitempty"`
}

// DownloadProgress tracks the progress of an update download
type DownloadProgress struct {
	// ID is the download ID
	ID string `json:"id"`

	// State is the current download state (pending, downloading, completed, failed)
	State string `json:"state"`

	// BytesDownloaded is the number of bytes downloaded
	BytesDownloaded int64 `json:"bytesDownloaded"`

	// TotalBytes is the total size in bytes
	TotalBytes int64 `json:"totalBytes"`

	// Percentage is the download percentage (0-100)
	Percentage float64 `json:"percentage"`

	// Speed is the current download speed in bytes per second
	Speed int64 `json:"speed"`

	// ETA is the estimated time remaining in seconds
	ETA int64 `json:"eta,omitempty"`

	// Error contains any error message
	Error string `json:"error,omitempty"`

	// FilePath is the local path where the file is being downloaded
	FilePath string `json:"filePath,omitempty"`
}

// ReleaseInfo provides detailed information about a specific release
type ReleaseInfo struct {
	// ID is the release ID
	ID string `json:"id"`

	// Version is the release version
	Version string `json:"version"`

	// Name is the release name
	Name string `json:"name,omitempty"`

	// Description contains release notes
	Description string `json:"description,omitempty"`

	// PublishedAt is when the release was published
	PublishedAt time.Time `json:"publishedAt"`

	// Channel is the release channel
	Channel string `json:"channel"`

	// Metadata contains custom release metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// Artifacts contains the release artifacts
	Artifacts []ReleaseArtifact `json:"artifacts"`
}

// UpdateInstalledEvent is emitted when an update is successfully installed
type UpdateInstalledEvent struct {
	Version     string    `json:"version"`
	InstalledAt time.Time `json:"installedAt"`
}

// Download states
const (
	DownloadStatePending     = "pending"
	DownloadStateDownloading = "downloading"
	DownloadStateCompleted   = "completed"
	DownloadStateFailed      = "failed"
	DownloadStateCancelled   = "cancelled"
)

// License validation codes
const (
	ValidationCodeValid               = "VALID"
	ValidationCodeInvalid             = "INVALID"
	ValidationCodeExpired             = "EXPIRED"
	ValidationCodeSuspended           = "SUSPENDED"
	ValidationCodeOverdue             = "OVERDUE"
	ValidationCodeNoMachines          = "NO_MACHINES"
	ValidationCodeMachineLimitReached = "MACHINE_LIMIT_REACHED"
	ValidationCodeFingerprintMismatch = "FINGERPRINT_MISMATCH"
)
