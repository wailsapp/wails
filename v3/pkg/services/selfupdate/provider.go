// Package selfupdate provides a service for self-updating Wails applications.
package selfupdate

import (
	"context"
	"io"
	"time"
)

// UpdateProvider is the interface that all update backends must implement.
// This abstraction allows the selfupdate service to work with different
// update sources (GitHub, Keygen, HTTP, etc.) through a unified API.
type UpdateProvider interface {
	// Name returns the provider identifier (e.g., "github", "keygen", "http").
	Name() string

	// Configure initializes the provider with the given configuration.
	// This is called once during service startup.
	Configure(ctx context.Context, config *ProviderConfig) error

	// CheckForUpdate queries the update source for available updates.
	// Returns nil if no update is available, or UpdateInfo if an update exists.
	CheckForUpdate(ctx context.Context, opts *CheckOptions) (*UpdateResult, error)

	// DownloadUpdate downloads the update binary.
	// The progress callback is called periodically with download progress.
	// Returns a reader for the downloaded content.
	DownloadUpdate(ctx context.Context, update *UpdateResult, progress ProgressFunc) (io.ReadCloser, error)

	// VerifyUpdate validates the downloaded update using signatures or checksums.
	// The data reader should be the downloaded update content.
	VerifyUpdate(ctx context.Context, update *UpdateResult, data io.Reader) error

	// Close releases any resources held by the provider.
	Close() error
}

// ProviderConfig holds configuration common to all providers.
type ProviderConfig struct {
	// CurrentVersion is the current version of the application (e.g., "1.0.0").
	CurrentVersion string

	// Channel is the release channel to check (e.g., "stable", "beta", "alpha").
	// If empty, defaults to "stable".
	Channel string

	// Platform is the target platform (e.g., "darwin", "linux", "windows").
	// If empty, defaults to runtime.GOOS.
	Platform string

	// Arch is the target architecture (e.g., "amd64", "arm64").
	// If empty, defaults to runtime.GOARCH.
	Arch string

	// Variant is an optional build variant (e.g., "webkit2_41").
	// Used for platform-specific builds with different dependencies.
	Variant string

	// AssetPattern is a template for matching release assets.
	// Supports: {name}, {version}, {goos}, {goarch}, {platform}, {variant}, {ext}
	// Default: "{name}_{goos}_{goarch}{ext}"
	AssetPattern string

	// PublicKey is the Ed25519 public key for signature verification (base64 encoded).
	// If empty, signature verification is skipped (not recommended for production).
	PublicKey string

	// Settings holds provider-specific configuration.
	// Each provider documents its required and optional settings.
	Settings map[string]any
}

// CheckOptions controls how update checks are performed.
type CheckOptions struct {
	// IncludePrerelease includes pre-release versions in the check.
	IncludePrerelease bool

	// Force bypasses any caching and forces a fresh check.
	Force bool
}

// UpdateResult contains information about an available update.
type UpdateResult struct {
	// UpdateAvailable is true if a newer version is available.
	UpdateAvailable bool `json:"updateAvailable"`

	// Version is the version string of the available update.
	Version string `json:"version"`

	// CurrentVersion is the currently installed version.
	CurrentVersion string `json:"currentVersion"`

	// ReleaseDate is when the update was published.
	ReleaseDate time.Time `json:"releaseDate"`

	// ReleaseNotes contains the changelog or release description.
	ReleaseNotes string `json:"releaseNotes"`

	// ReleaseURL is a link to the release page (for manual download).
	ReleaseURL string `json:"releaseUrl"`

	// DownloadURL is the direct download URL for the update asset.
	DownloadURL string `json:"downloadUrl"`

	// Size is the size of the update in bytes.
	Size int64 `json:"size"`

	// AssetName is the filename of the update asset.
	AssetName string `json:"assetName"`

	// Signature is the Ed25519 signature of the asset (base64 encoded).
	Signature string `json:"signature,omitempty"`

	// Checksum is the SHA256 checksum of the asset (hex encoded).
	Checksum string `json:"checksum,omitempty"`

	// Channel is the release channel this update belongs to.
	Channel string `json:"channel,omitempty"`

	// Metadata holds provider-specific data.
	Metadata map[string]any `json:"metadata,omitempty"`
}

// ProgressFunc is called during download to report progress.
type ProgressFunc func(info *ProgressInfo)

// ProgressInfo contains download progress information.
type ProgressInfo struct {
	// State is the current state: "started", "downloading", "finished", "error".
	State string `json:"state"`

	// TotalBytes is the total size of the download.
	TotalBytes int64 `json:"totalBytes"`

	// DownloadedBytes is the number of bytes downloaded so far.
	DownloadedBytes int64 `json:"downloadedBytes"`

	// Percentage is the download progress as a percentage (0-100).
	Percentage float64 `json:"percentage"`

	// BytesPerSecond is the current download speed.
	BytesPerSecond float64 `json:"bytesPerSecond"`

	// Error contains any error that occurred (when State is "error").
	Error string `json:"error,omitempty"`
}

// ProviderWithLicensing is an optional interface for providers that support
// license-gated updates (e.g., Keygen).
type ProviderWithLicensing interface {
	UpdateProvider

	// SetLicense sets the license key for authenticated update checks.
	SetLicense(licenseKey string) error

	// ValidateLicense checks if the current license is valid.
	ValidateLicense(ctx context.Context) error

	// IsLicenseRequired returns true if this provider requires a license.
	IsLicenseRequired() bool
}

// ProviderWithChannels is an optional interface for providers that support
// multiple release channels.
type ProviderWithChannels interface {
	UpdateProvider

	// AvailableChannels returns the list of available release channels.
	AvailableChannels(ctx context.Context) ([]string, error)

	// SetChannel changes the active release channel.
	SetChannel(channel string) error
}

// PrepareFunc is called after download to prepare the update for installation.
// This handles platform-specific tasks like extracting archives or locating
// the actual binary within a bundle (e.g., macOS .app).
//
// Parameters:
//   - ctx: Context for cancellation
//   - downloadPath: Path to the downloaded file
//   - config: Current provider configuration
//
// Returns:
//   - string: Path to the prepared executable/bundle ready for installation
//   - error: Any error that occurred during preparation
type PrepareFunc func(ctx context.Context, downloadPath string, config *ProviderConfig) (string, error)
