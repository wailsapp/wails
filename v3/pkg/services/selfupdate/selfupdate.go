// Package selfupdate provides a service for self-updating Wails applications.
//
// The service supports multiple sources for update discovery:
//   - GitHub Releases (default)
//   - GitLab Releases
//   - Gitea Releases
//
// Example usage:
//
//	app := application.New(application.Options{
//	    Services: []application.Service{
//	        application.NewService(selfupdate.NewWithConfig(&selfupdate.Config{
//	            CurrentVersion: "1.0.0",
//	            Source:         selfupdate.SourceGitHub,
//	            Repository:     "owner/repo",
//	        })),
//	    },
//	})
package selfupdate

import (
	"context"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/creativeprojects/go-selfupdate"
	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// Source represents the update source type.
type Source string

const (
	// SourceGitHub fetches updates from GitHub Releases.
	SourceGitHub Source = "github"
	// SourceGitLab fetches updates from GitLab Releases.
	SourceGitLab Source = "gitlab"
	// SourceGitea fetches updates from Gitea Releases.
	SourceGitea Source = "gitea"
)

// SignatureType represents the type of signature verification to use.
type SignatureType string

const (
	// SignatureNone disables signature verification (not recommended for production).
	SignatureNone SignatureType = ""
	// SignatureChecksum verifies updates using SHA256 checksums.
	SignatureChecksum SignatureType = "checksum"
	// SignatureECDSA verifies updates using ECDSA signatures.
	SignatureECDSA SignatureType = "ecdsa"
	// SignaturePGP verifies updates using PGP signatures.
	SignaturePGP SignatureType = "pgp"
)

// Config holds the configuration for the selfupdate service.
type Config struct {
	// CurrentVersion is the current version of the application (e.g., "1.0.0" or "v1.0.0").
	// This is required.
	CurrentVersion string

	// Source is the update source type (github, gitlab, gitea).
	// Defaults to SourceGitHub.
	Source Source

	// Repository is the repository identifier in "owner/repo" format.
	// This is required.
	Repository string

	// BaseURL is the base URL for self-hosted GitLab or Gitea instances.
	// For GitHub Enterprise, this is the API base URL (e.g., "https://your-org/api/v3/").
	// Not required for public GitHub or public GitLab.
	BaseURL string

	// Token is an optional access token for private repositories or to increase rate limits.
	Token string

	// AllowPrerelease allows updating to pre-release versions.
	AllowPrerelease bool

	// AllowDowngrade allows downgrading to older versions.
	AllowDowngrade bool

	// Signature configures update signature verification.
	Signature *SignatureConfig

	// UI configures the update progress UI appearance.
	UI *UIConfig
}

// SignatureConfig holds signature verification configuration.
type SignatureConfig struct {
	// Type is the signature verification type (checksum, ecdsa, pgp).
	Type SignatureType

	// PublicKey is the public key for signature verification.
	// For ECDSA: PEM-encoded public key.
	// For PGP: ASCII-armored public key.
	// For Checksum: not used.
	PublicKey string

	// ChecksumFilename is the name of the checksums file for checksum validation.
	// Defaults to "checksums.txt".
	ChecksumFilename string
}

// UIConfig holds UI theming configuration for the update progress window.
type UIConfig struct {
	// Title is the window title for the update dialog.
	// Defaults to "Update Available".
	Title string

	// IconPath is the path to a custom icon for the update dialog.
	IconPath string

	// BackgroundColor is the background color in hex format (e.g., "#ffffff").
	BackgroundColor string

	// TextColor is the text color in hex format (e.g., "#000000").
	TextColor string

	// AccentColor is the accent/progress bar color in hex format (e.g., "#007bff").
	AccentColor string

	// ProgressBarColor is the progress bar fill color in hex format.
	// Defaults to AccentColor if not set.
	ProgressBarColor string

	// ButtonColor is the button background color in hex format.
	ButtonColor string

	// ButtonTextColor is the button text color in hex format.
	ButtonTextColor string
}

// UpdateInfo contains information about an available update.
type UpdateInfo struct {
	// UpdateAvailable is true if an update is available.
	UpdateAvailable bool `json:"updateAvailable"`

	// CurrentVersion is the current version of the application.
	CurrentVersion string `json:"currentVersion"`

	// LatestVersion is the latest available version.
	LatestVersion string `json:"latestVersion"`

	// ReleaseNotes contains the release notes/changelog for the update.
	ReleaseNotes string `json:"releaseNotes"`

	// ReleaseURL is a link to the release page.
	ReleaseURL string `json:"releaseUrl"`

	// PublishedAt is the publication date of the release.
	PublishedAt string `json:"publishedAt"`

	// AssetName is the name of the asset that will be downloaded.
	AssetName string `json:"assetName"`

	// AssetURL is the direct download URL for the asset.
	AssetURL string `json:"assetUrl"`

	// AssetSize is the size of the asset in bytes.
	AssetSize int `json:"assetSize"`
}

// DownloadProgress contains information about download progress.
type DownloadProgress struct {
	// State is the current download state: "started", "progress", or "finished".
	State string `json:"state"`

	// TotalBytes is the total size of the download in bytes.
	TotalBytes int64 `json:"totalBytes"`

	// DownloadedBytes is the number of bytes downloaded so far.
	DownloadedBytes int64 `json:"downloadedBytes"`

	// Percentage is the download progress as a percentage (0-100).
	Percentage float64 `json:"percentage"`
}

// Service provides self-update functionality for Wails applications.
type Service struct {
	lock          sync.RWMutex
	config        *Config
	configErr     error // Stores any configuration error for reporting at startup.
	source        selfupdate.Source
	app           *application.App
	pendingUpdate *selfupdate.Release
	updater       *selfupdate.Updater
}

// New creates a new selfupdate service with default configuration.
// Note: You must call Configure with a valid Config before using the service.
func New() *Service {
	return NewWithConfig(nil)
}

// NewWithConfig creates a new selfupdate service with the given configuration.
// If config is nil, you must call Configure before using the service.
// If the config is invalid, the error will be returned when ServiceStartup is called.
func NewWithConfig(config *Config) *Service {
	s := &Service{}
	if config != nil {
		if err := s.Configure(config); err != nil {
			s.configErr = err
		}
	}
	return s
}

// ServiceName returns the name of the service.
func (s *Service) ServiceName() string {
	return "github.com/wailsapp/wails/v3/services/selfupdate"
}

// ServiceStartup initializes the service.
func (s *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.app = application.Get()

	// Return any configuration error from NewWithConfig.
	if s.configErr != nil {
		return errors.Wrap(s.configErr, "selfupdate configuration error")
	}

	if s.config == nil {
		return errors.New("selfupdate service not configured; call Configure with a valid Config")
	}

	return s.initSource()
}

// ServiceShutdown cleans up the service.
func (s *Service) ServiceShutdown() error {
	return nil
}

// Configure sets the service configuration.
// This can be called before registering the service or during runtime.
//
//wails:ignore
func (s *Service) Configure(config *Config) error {
	if config == nil {
		return errors.New("config cannot be nil")
	}
	if config.Repository == "" {
		return errors.New("repository is required")
	}
	if config.CurrentVersion == "" {
		return errors.New("currentVersion is required")
	}

	// Clone to prevent changes from the outside.
	clone := new(Config)
	*clone = *config

	// Deep clone nested pointer fields.
	if config.Signature != nil {
		sig := *config.Signature
		clone.Signature = &sig
	}
	if config.UI != nil {
		ui := *config.UI
		clone.UI = &ui
	}

	// Set defaults.
	if clone.Source == "" {
		clone.Source = SourceGitHub
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	s.config = clone
	return s.initSource()
}

// initSource initializes the update source based on configuration.
// Must be called with lock held.
func (s *Service) initSource() error {
	if s.config == nil {
		return nil
	}

	var source selfupdate.Source
	var err error

	switch s.config.Source {
	case SourceGitHub:
		ghConfig := selfupdate.GitHubConfig{}
		if s.config.Token != "" {
			ghConfig.APIToken = s.config.Token
		}
		if s.config.BaseURL != "" {
			ghConfig.EnterpriseBaseURL = s.config.BaseURL
		}
		source, err = selfupdate.NewGitHubSource(ghConfig)
	case SourceGitLab:
		glConfig := selfupdate.GitLabConfig{}
		if s.config.BaseURL != "" {
			glConfig.BaseURL = s.config.BaseURL
		}
		if s.config.Token != "" {
			glConfig.APIToken = s.config.Token
		}
		source, err = selfupdate.NewGitLabSource(glConfig)
	case SourceGitea:
		if s.config.BaseURL == "" {
			return errors.New("baseURL is required for Gitea source")
		}
		gtConfig := selfupdate.GiteaConfig{
			BaseURL: s.config.BaseURL,
		}
		if s.config.Token != "" {
			gtConfig.APIToken = s.config.Token
		}
		source, err = selfupdate.NewGiteaSource(gtConfig)
	default:
		return errors.Errorf("unsupported source: %s", s.config.Source)
	}

	if err != nil {
		return errors.Wrap(err, "failed to create update source")
	}

	s.source = source

	// Create the updater config.
	updaterConfig := selfupdate.Config{
		Source:     source,
		Prerelease: s.config.AllowPrerelease,
	}

	// Configure signature validation if specified.
	if s.config.Signature != nil {
		validator, err := s.createValidator(s.config.Signature)
		if err != nil {
			return errors.Wrap(err, "failed to create signature validator")
		}
		if validator != nil {
			updaterConfig.Validator = validator
		}
	}

	// Create the updater.
	s.updater, err = selfupdate.NewUpdater(updaterConfig)
	if err != nil {
		return errors.Wrap(err, "failed to create updater")
	}

	return nil
}

// createValidator creates a signature validator based on the configuration.
func (s *Service) createValidator(config *SignatureConfig) (validator selfupdate.Validator, err error) {
	if config == nil {
		return nil, nil
	}

	// Recover from panics in the validator creation (the library panics on invalid keys).
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("failed to create validator: %v", r)
			validator = nil
		}
	}()

	switch config.Type {
	case SignatureNone:
		return nil, nil

	case SignatureChecksum:
		filename := config.ChecksumFilename
		if filename == "" {
			filename = "checksums.txt"
		}
		return &selfupdate.ChecksumValidator{
			UniqueFilename: filename,
		}, nil

	case SignatureECDSA:
		if config.PublicKey == "" {
			return nil, errors.New("public key is required for ECDSA signature verification")
		}
		v := &selfupdate.ECDSAValidator{}
		return v.WithPublicKey([]byte(config.PublicKey)), nil

	case SignaturePGP:
		if config.PublicKey == "" {
			return nil, errors.New("public key is required for PGP signature verification")
		}
		v := &selfupdate.PGPValidator{}
		return v.WithArmoredKeyRing([]byte(config.PublicKey)), nil

	default:
		return nil, errors.Errorf("unsupported signature type: %s", config.Type)
	}
}

// parseRepository splits "owner/repo" into owner and repo parts.
func parseRepository(repo string) (owner, name string) {
	parts := strings.SplitN(repo, "/", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return repo, ""
}

// Check checks if an update is available.
// Returns UpdateInfo with details about the available update.
// This is the recommended first step in the update flow.
func (s *Service) Check(ctx context.Context) (*UpdateInfo, error) {
	s.lock.RLock()
	config := s.config
	updater := s.updater
	s.lock.RUnlock()

	if config == nil || updater == nil {
		return nil, errors.New("selfupdate service not configured")
	}

	owner, repo := parseRepository(config.Repository)
	latest, found, err := updater.DetectLatest(ctx, selfupdate.NewRepositorySlug(owner, repo))
	if err != nil {
		return nil, errors.Wrap(err, "failed to detect latest version")
	}

	info := &UpdateInfo{
		CurrentVersion: config.CurrentVersion,
	}

	if !found {
		info.UpdateAvailable = false
		info.LatestVersion = config.CurrentVersion
		return info, nil
	}

	info.LatestVersion = latest.Version()
	info.ReleaseNotes = latest.ReleaseNotes
	info.ReleaseURL = latest.URL
	info.AssetName = latest.AssetName
	info.AssetURL = latest.AssetURL
	info.AssetSize = latest.AssetByteSize
	if !latest.PublishedAt.IsZero() {
		info.PublishedAt = latest.PublishedAt.Format("2006-01-02T15:04:05Z")
	}

	// Compare versions using the library's version comparison.
	if config.AllowDowngrade {
		info.UpdateAvailable = !latest.Equal(config.CurrentVersion)
	} else {
		info.UpdateAvailable = latest.GreaterThan(config.CurrentVersion)
	}

	// Store the pending update for later download/install.
	if info.UpdateAvailable {
		s.lock.Lock()
		s.pendingUpdate = latest
		s.lock.Unlock()
	}

	return info, nil
}

// DownloadAndInstall downloads and installs the pending update in one step.
// You must call Check first to detect an available update.
// Returns true if the update was successful.
// The application should be restarted after a successful update.
func (s *Service) DownloadAndInstall(ctx context.Context) (bool, error) {
	s.lock.RLock()
	config := s.config
	updater := s.updater
	pending := s.pendingUpdate
	s.lock.RUnlock()

	if config == nil || updater == nil {
		return false, errors.New("selfupdate service not configured")
	}

	if pending == nil {
		// No pending update, try to check and update directly.
		owner, repo := parseRepository(config.Repository)
		latest, err := updater.UpdateSelf(ctx, config.CurrentVersion, selfupdate.NewRepositorySlug(owner, repo))
		if err != nil {
			return false, errors.Wrap(err, "failed to update")
		}
		if latest == nil {
			return false, nil
		}
		return latest.GreaterThan(config.CurrentVersion), nil
	}

	// Get the executable path.
	exe, err := selfupdate.ExecutablePath()
	if err != nil {
		return false, errors.Wrap(err, "failed to get executable path")
	}

	// Emit download started event.
	s.emitProgress(DownloadProgress{
		State:      "started",
		TotalBytes: int64(pending.AssetByteSize),
	})

	// Download and apply the update.
	if err := updater.UpdateTo(ctx, pending, exe); err != nil {
		return false, errors.Wrap(err, "failed to apply update")
	}

	// Emit download finished event.
	s.emitProgress(DownloadProgress{
		State:           "finished",
		TotalBytes:      int64(pending.AssetByteSize),
		DownloadedBytes: int64(pending.AssetByteSize),
		Percentage:      100,
	})

	// Clear the pending update.
	s.lock.Lock()
	s.pendingUpdate = nil
	s.lock.Unlock()

	return true, nil
}

// emitProgress emits a progress event to the frontend.
func (s *Service) emitProgress(progress DownloadProgress) {
	if s.app != nil {
		s.app.Event.Emit("selfupdate:progress", progress)
	}
}

// GetCurrentVersion returns the current application version.
func (s *Service) GetCurrentVersion() string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if s.config == nil {
		return ""
	}
	return s.config.CurrentVersion
}

// GetPlatformInfo returns information about the current platform.
// This can be useful for debugging update issues or for determining
// which asset to expect.
func (s *Service) GetPlatformInfo() map[string]string {
	return map[string]string{
		"os":   runtime.GOOS,
		"arch": runtime.GOARCH,
	}
}

// CanUpdate returns true if the current process has permission to update itself.
// This checks if the executable is writable.
func (s *Service) CanUpdate() bool {
	exe, err := selfupdate.ExecutablePath()
	if err != nil {
		return false
	}

	// Check if we can write to the executable.
	file, err := os.OpenFile(exe, os.O_WRONLY, 0)
	if err != nil {
		return false
	}
	file.Close()
	return true
}

// Restart restarts the application after an update.
// This is a convenience method that can be called after DownloadAndInstall.
// Note: This method may not return if the restart is successful.
func (s *Service) Restart() error {
	if s.app == nil {
		return errors.New("application not available")
	}

	// Get the executable path.
	exe, err := os.Executable()
	if err != nil {
		return errors.Wrap(err, "failed to get executable path")
	}

	// Start a new instance of the application.
	cmd := &os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}

	_, err = os.StartProcess(exe, os.Args, cmd)
	if err != nil {
		return errors.Wrap(err, "failed to start new process")
	}

	// Quit the current application.
	s.app.Quit()

	return nil
}

// CheckForUpdate is an alias for Check for backwards compatibility.
func (s *Service) CheckForUpdate(ctx context.Context) (*UpdateInfo, error) {
	return s.Check(ctx)
}

// PerformUpdate is an alias for DownloadAndInstall for backwards compatibility.
func (s *Service) PerformUpdate(ctx context.Context) (bool, error) {
	return s.DownloadAndInstall(ctx)
}

// GetLatestRelease returns information about the latest release.
// This is similar to Check but doesn't store the pending update.
func (s *Service) GetLatestRelease(ctx context.Context) (*UpdateInfo, error) {
	s.lock.RLock()
	config := s.config
	updater := s.updater
	s.lock.RUnlock()

	if config == nil || updater == nil {
		return nil, errors.New("selfupdate service not configured")
	}

	owner, repo := parseRepository(config.Repository)
	latest, found, err := updater.DetectLatest(ctx, selfupdate.NewRepositorySlug(owner, repo))
	if err != nil {
		return nil, errors.Wrap(err, "failed to detect latest version")
	}

	info := &UpdateInfo{
		CurrentVersion: config.CurrentVersion,
	}

	if !found {
		info.UpdateAvailable = false
		return info, nil
	}

	info.LatestVersion = latest.Version()
	info.ReleaseNotes = latest.ReleaseNotes
	info.ReleaseURL = latest.URL
	info.AssetName = latest.AssetName
	info.AssetURL = latest.AssetURL
	info.AssetSize = latest.AssetByteSize
	if !latest.PublishedAt.IsZero() {
		info.PublishedAt = latest.PublishedAt.Format("2006-01-02T15:04:05Z")
	}

	if config.AllowDowngrade {
		info.UpdateAvailable = !latest.Equal(config.CurrentVersion)
	} else {
		info.UpdateAvailable = latest.GreaterThan(config.CurrentVersion)
	}

	return info, nil
}

// GetUIConfig returns the UI configuration for theming the update dialog.
// Returns nil if no UI configuration is set.
func (s *Service) GetUIConfig() *UIConfig {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if s.config == nil || s.config.UI == nil {
		return nil
	}

	// Return a copy to prevent external modification.
	ui := *s.config.UI
	return &ui
}
