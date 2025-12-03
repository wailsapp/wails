package updater

import (
	"time"
)

// UpdateInfo contains information about an available update
type UpdateInfo struct {
	// Version is the new version string (e.g., "1.2.0")
	Version string `json:"version"`

	// ReleaseDate is when this version was released
	ReleaseDate time.Time `json:"release_date"`

	// ReleaseNotes contains markdown-formatted release notes
	ReleaseNotes string `json:"release_notes,omitempty"`

	// DownloadURL is the URL to download the full update
	DownloadURL string `json:"download_url"`

	// Size is the size of the full download in bytes
	Size int64 `json:"size"`

	// Checksum is the SHA256 checksum of the download
	Checksum string `json:"checksum"`

	// Signature is the Ed25519 signature of the checksum
	Signature string `json:"signature,omitempty"`

	// PatchURL is the URL to download a delta patch (optional)
	PatchURL string `json:"patch_url,omitempty"`

	// PatchSize is the size of the patch in bytes
	PatchSize int64 `json:"patch_size,omitempty"`

	// PatchChecksum is the SHA256 checksum of the patch
	PatchChecksum string `json:"patch_checksum,omitempty"`

	// PatchFrom is the version this patch can be applied from
	PatchFrom string `json:"patch_from,omitempty"`

	// Mandatory indicates if this update is required
	Mandatory bool `json:"mandatory,omitempty"`
}

// Manifest is the update manifest file structure
type Manifest struct {
	// Version is the latest version available
	Version string `json:"version"`

	// ReleaseDate is when this version was released
	ReleaseDate time.Time `json:"release_date"`

	// ReleaseNotes contains markdown-formatted release notes
	ReleaseNotes string `json:"release_notes,omitempty"`

	// Platforms contains platform-specific download information
	Platforms map[string]PlatformUpdate `json:"platforms"`

	// MinimumVersion is the minimum version that can update to this version
	MinimumVersion string `json:"minimum_version,omitempty"`

	// Mandatory indicates if this update is required
	Mandatory bool `json:"mandatory,omitempty"`
}

// PlatformUpdate contains update information for a specific platform
type PlatformUpdate struct {
	// URL is the download URL for this platform
	URL string `json:"url"`

	// Size is the download size in bytes
	Size int64 `json:"size"`

	// Checksum is the SHA256 checksum
	Checksum string `json:"checksum"`

	// Signature is the Ed25519 signature of the checksum
	Signature string `json:"signature,omitempty"`

	// Patches contains available delta patches from previous versions
	Patches []PatchInfo `json:"patches,omitempty"`
}

// PatchInfo contains information about a delta patch
type PatchInfo struct {
	// From is the version this patch applies from
	From string `json:"from"`

	// URL is the download URL for the patch
	URL string `json:"url"`

	// Size is the patch size in bytes
	Size int64 `json:"size"`

	// Checksum is the SHA256 checksum of the patch
	Checksum string `json:"checksum"`

	// Signature is the Ed25519 signature
	Signature string `json:"signature,omitempty"`
}

// Config defines the updater configuration
type Config struct {
	// UpdateURL is the base URL for update checks (e.g., "https://updates.example.com/myapp/")
	UpdateURL string `yaml:"url" json:"url"`

	// CheckInterval is how often to automatically check for updates (0 = disabled)
	CheckInterval time.Duration `yaml:"check_interval" json:"check_interval"`

	// AllowPrerelease determines whether to include prerelease versions
	AllowPrerelease bool `yaml:"allow_prerelease" json:"allow_prerelease"`

	// PublicKey is the Ed25519 public key for signature verification
	PublicKey string `yaml:"public_key" json:"public_key,omitempty"`

	// RequireSignature determines whether to require signed updates
	RequireSignature bool `yaml:"require_signature" json:"require_signature"`

	// Channel is the update channel (e.g., "stable", "beta", "canary")
	Channel string `yaml:"channel" json:"channel,omitempty"`
}

// DownloadProgress is passed to progress callbacks during download
type DownloadProgress struct {
	// Downloaded is the number of bytes downloaded so far
	Downloaded int64

	// Total is the total size in bytes
	Total int64

	// Percentage is the download percentage (0-100)
	Percentage float64

	// BytesPerSecond is the current download speed
	BytesPerSecond float64
}

// ProgressCallback is called periodically during download
type ProgressCallback func(progress DownloadProgress)

// UpdateState represents the current state of an update
type UpdateState string

const (
	// StateIdle means no update is in progress
	StateIdle UpdateState = "idle"

	// StateChecking means we're checking for updates
	StateChecking UpdateState = "checking"

	// StateAvailable means an update is available
	StateAvailable UpdateState = "available"

	// StateDownloading means we're downloading the update
	StateDownloading UpdateState = "downloading"

	// StateReady means the update is downloaded and ready to install
	StateReady UpdateState = "ready"

	// StateInstalling means the update is being installed
	StateInstalling UpdateState = "installing"

	// StateError means an error occurred
	StateError UpdateState = "error"
)

// UpdateEvent is emitted when update state changes
type UpdateEvent struct {
	// State is the current update state
	State UpdateState `json:"state"`

	// Info contains update information (when State is StateAvailable or later)
	Info *UpdateInfo `json:"info,omitempty"`

	// Progress contains download progress (when State is StateDownloading)
	Progress *DownloadProgress `json:"progress,omitempty"`

	// Error contains error information (when State is StateError)
	Error string `json:"error,omitempty"`
}
