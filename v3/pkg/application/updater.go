package application

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/internal/updater"
)

// UpdaterService provides update functionality to the application
type UpdaterService struct {
	updater        *updater.Updater
	config         UpdaterConfig
	currentVersion string
	mu             sync.RWMutex
}

// UpdaterConfig defines the configuration for the updater
type UpdaterConfig struct {
	// UpdateURL is the base URL for update checks
	UpdateURL string

	// CheckInterval is how often to automatically check for updates (0 = disabled)
	CheckInterval time.Duration

	// AllowPrerelease determines whether to include prerelease versions
	AllowPrerelease bool

	// PublicKey is the Ed25519 public key for signature verification (base64 encoded)
	PublicKey string

	// RequireSignature determines whether to require signed updates
	RequireSignature bool

	// Channel is the update channel (e.g., "stable", "beta", "canary")
	Channel string
}

// NewUpdaterService creates a new updater service
func NewUpdaterService(config UpdaterConfig, currentVersion string) *UpdaterService {
	internalConfig := updater.Config{
		UpdateURL:        config.UpdateURL,
		CheckInterval:    config.CheckInterval,
		AllowPrerelease:  config.AllowPrerelease,
		PublicKey:        config.PublicKey,
		RequireSignature: config.RequireSignature,
		Channel:          config.Channel,
	}

	return &UpdaterService{
		updater:        updater.New(internalConfig, currentVersion),
		config:         config,
		currentVersion: currentVersion,
	}
}

// ServiceName returns the name used for binding
func (u *UpdaterService) ServiceName() string {
	return "updater"
}

// ServiceStartup is called when the application starts
func (u *UpdaterService) ServiceStartup(ctx context.Context, options ServiceOptions) error {
	// Start background checking if configured
	if u.config.CheckInterval > 0 {
		u.updater.StartBackgroundChecks(ctx)
	}
	return nil
}

// ServiceShutdown is called when the application shuts down
func (u *UpdaterService) ServiceShutdown() error {
	u.updater.StopBackgroundChecks()
	return nil
}

// GetCurrentVersion returns the current application version
func (u *UpdaterService) GetCurrentVersion() string {
	return u.currentVersion
}

// CheckForUpdate checks if a new version is available
func (u *UpdaterService) CheckForUpdate() (*UpdateInfo, error) {
	ctx := context.Background()
	info, err := u.updater.CheckForUpdate(ctx)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}
	return convertUpdateInfo(info), nil
}

// DownloadUpdate downloads the available update
// Returns progress events via the event system
func (u *UpdaterService) DownloadUpdate() error {
	ctx := context.Background()

	progress := func(p updater.DownloadProgress) {
		// Emit progress event
		globalApplication.EmitEvent("updater:progress", map[string]interface{}{
			"downloaded":     p.Downloaded,
			"total":          p.Total,
			"percentage":     p.Percentage,
			"bytesPerSecond": p.BytesPerSecond,
		})
	}

	return u.updater.DownloadUpdate(ctx, progress)
}

// ApplyUpdate applies the downloaded update
// This will restart the application
func (u *UpdaterService) ApplyUpdate() error {
	ctx := context.Background()
	return u.updater.ApplyUpdate(ctx)
}

// DownloadAndApply downloads and applies an update in one call
func (u *UpdaterService) DownloadAndApply() error {
	ctx := context.Background()

	progress := func(p updater.DownloadProgress) {
		globalApplication.EmitEvent("updater:progress", map[string]interface{}{
			"downloaded":     p.Downloaded,
			"total":          p.Total,
			"percentage":     p.Percentage,
			"bytesPerSecond": p.BytesPerSecond,
		})
	}

	return u.updater.DownloadAndApply(ctx, progress)
}

// GetState returns the current update state
func (u *UpdaterService) GetState() string {
	return string(u.updater.GetState())
}

// GetUpdateInfo returns information about the available update
func (u *UpdaterService) GetUpdateInfo() *UpdateInfo {
	info := u.updater.GetUpdateInfo()
	if info == nil {
		return nil
	}
	return convertUpdateInfo(info)
}

// GetLastError returns the last error message
func (u *UpdaterService) GetLastError() string {
	err := u.updater.GetLastError()
	if err == nil {
		return ""
	}
	return err.Error()
}

// OnUpdateAvailable registers a handler for when an update is available
func (u *UpdaterService) OnUpdateAvailable(handler func(info *UpdateInfo)) {
	u.updater.OnEvent(func(event updater.UpdateEvent) {
		if event.State == updater.StateAvailable {
			handler(convertUpdateInfo(event.Info))
		}
	})
}

// OnDownloadProgress registers a handler for download progress
func (u *UpdaterService) OnDownloadProgress(handler func(downloaded, total int64, percentage float64)) {
	u.updater.OnEvent(func(event updater.UpdateEvent) {
		if event.State == updater.StateDownloading && event.Progress != nil {
			handler(event.Progress.Downloaded, event.Progress.Total, event.Progress.Percentage)
		}
	})
}

// OnError registers a handler for errors
func (u *UpdaterService) OnError(handler func(err string)) {
	u.updater.OnEvent(func(event updater.UpdateEvent) {
		if event.State == updater.StateError {
			handler(event.Error)
		}
	})
}

// Reset resets the updater state
func (u *UpdaterService) Reset() {
	u.updater.Reset()
}

// UpdateInfo represents information about an available update
type UpdateInfo struct {
	Version      string    `json:"version"`
	ReleaseDate  time.Time `json:"releaseDate"`
	ReleaseNotes string    `json:"releaseNotes"`
	Size         int64     `json:"size"`
	PatchSize    int64     `json:"patchSize,omitempty"`
	Mandatory    bool      `json:"mandatory"`
	HasPatch     bool      `json:"hasPatch"`
}

// convertUpdateInfo converts internal update info to the public type
func convertUpdateInfo(info *updater.UpdateInfo) *UpdateInfo {
	if info == nil {
		return nil
	}
	return &UpdateInfo{
		Version:      info.Version,
		ReleaseDate:  info.ReleaseDate,
		ReleaseNotes: info.ReleaseNotes,
		Size:         info.Size,
		PatchSize:    info.PatchSize,
		Mandatory:    info.Mandatory,
		HasPatch:     info.PatchURL != "",
	}
}

// UpdaterOption is a function that configures the updater
type UpdaterOption func(*UpdaterConfig)

// WithUpdateURL sets the update URL
func WithUpdateURL(url string) UpdaterOption {
	return func(c *UpdaterConfig) {
		c.UpdateURL = url
	}
}

// WithCheckInterval sets the automatic check interval
func WithCheckInterval(interval time.Duration) UpdaterOption {
	return func(c *UpdaterConfig) {
		c.CheckInterval = interval
	}
}

// WithAllowPrerelease allows prerelease versions
func WithAllowPrerelease(allow bool) UpdaterOption {
	return func(c *UpdaterConfig) {
		c.AllowPrerelease = allow
	}
}

// WithPublicKey sets the public key for signature verification
func WithPublicKey(key string) UpdaterOption {
	return func(c *UpdaterConfig) {
		c.PublicKey = key
	}
}

// WithRequireSignature requires signed updates
func WithRequireSignature(require bool) UpdaterOption {
	return func(c *UpdaterConfig) {
		c.RequireSignature = require
	}
}

// WithChannel sets the update channel
func WithChannel(channel string) UpdaterOption {
	return func(c *UpdaterConfig) {
		c.Channel = channel
	}
}

// CreateUpdaterService creates an updater service with the given options
func CreateUpdaterService(currentVersion string, opts ...UpdaterOption) (*UpdaterService, error) {
	if currentVersion == "" {
		return nil, fmt.Errorf("current version is required")
	}

	config := UpdaterConfig{}
	for _, opt := range opts {
		opt(&config)
	}

	if config.UpdateURL == "" {
		return nil, fmt.Errorf("update URL is required")
	}

	return NewUpdaterService(config, currentVersion), nil
}
