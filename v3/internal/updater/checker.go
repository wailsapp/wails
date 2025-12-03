package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"

	"github.com/Masterminds/semver"
)

// Checker checks for available updates
type Checker struct {
	config     Config
	httpClient *http.Client
	lastCheck  time.Time
	cachedInfo *UpdateInfo
}

// NewChecker creates a new update checker
func NewChecker(config Config) *Checker {
	return &Checker{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CheckForUpdate checks if a new version is available
func (c *Checker) CheckForUpdate(ctx context.Context, currentVersion string) (*UpdateInfo, error) {
	manifest, err := c.fetchManifest(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch manifest: %w", err)
	}

	// Parse versions for comparison
	current, err := semver.NewVersion(currentVersion)
	if err != nil {
		return nil, fmt.Errorf("invalid current version: %w", err)
	}

	latest, err := semver.NewVersion(manifest.Version)
	if err != nil {
		return nil, fmt.Errorf("invalid manifest version: %w", err)
	}

	// Check if update is available
	if !latest.GreaterThan(current) {
		return nil, nil // No update available
	}

	// Check prerelease policy
	if latest.Prerelease() != "" && !c.config.AllowPrerelease {
		return nil, nil // Prerelease not allowed
	}

	// Check minimum version requirement
	if manifest.MinimumVersion != "" {
		minVersion, err := semver.NewVersion(manifest.MinimumVersion)
		if err == nil && current.LessThan(minVersion) {
			// Current version is too old for incremental update
			// User needs to do a full reinstall
			return nil, fmt.Errorf("current version %s is too old, minimum required is %s. Please reinstall the application",
				currentVersion, manifest.MinimumVersion)
		}
	}

	// Get platform-specific update info
	platform := getPlatformKey()
	platformUpdate, ok := manifest.Platforms[platform]
	if !ok {
		return nil, fmt.Errorf("no update available for platform: %s", platform)
	}

	info := &UpdateInfo{
		Version:      manifest.Version,
		ReleaseDate:  manifest.ReleaseDate,
		ReleaseNotes: manifest.ReleaseNotes,
		DownloadURL:  platformUpdate.URL,
		Size:         platformUpdate.Size,
		Checksum:     platformUpdate.Checksum,
		Signature:    platformUpdate.Signature,
		Mandatory:    manifest.Mandatory,
	}

	// Check for applicable patch
	for _, patch := range platformUpdate.Patches {
		if patch.From == currentVersion {
			info.PatchURL = patch.URL
			info.PatchSize = patch.Size
			info.PatchChecksum = patch.Checksum
			info.PatchFrom = patch.From
			break
		}
	}

	c.lastCheck = time.Now()
	c.cachedInfo = info

	return info, nil
}

// GetCachedInfo returns the cached update info from the last check
func (c *Checker) GetCachedInfo() *UpdateInfo {
	return c.cachedInfo
}

// GetLastCheckTime returns when the last check was performed
func (c *Checker) GetLastCheckTime() time.Time {
	return c.lastCheck
}

// fetchManifest retrieves the update manifest from the server
func (c *Checker) fetchManifest(ctx context.Context) (*Manifest, error) {
	manifestURL, err := c.buildManifestURL()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, manifestURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Wails-Updater/1.0")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var manifest Manifest
	if err := json.Unmarshal(body, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	return &manifest, nil
}

// buildManifestURL builds the URL for the update manifest
func (c *Checker) buildManifestURL() (string, error) {
	baseURL := strings.TrimSuffix(c.config.UpdateURL, "/")

	// Add channel to path if specified
	if c.config.Channel != "" {
		baseURL = fmt.Sprintf("%s/%s", baseURL, c.config.Channel)
	}

	manifestURL, err := url.JoinPath(baseURL, "update.json")
	if err != nil {
		return "", err
	}

	return manifestURL, nil
}

// getPlatformKey returns the platform identifier for the current system
func getPlatformKey() string {
	os := runtime.GOOS
	arch := runtime.GOARCH

	// Normalize platform names
	switch os {
	case "darwin":
		os = "macos"
	}

	return fmt.Sprintf("%s-%s", os, arch)
}

// CompareVersions compares two semantic versions
// Returns: -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func CompareVersions(v1, v2 string) (int, error) {
	ver1, err := semver.NewVersion(v1)
	if err != nil {
		return 0, fmt.Errorf("invalid version %s: %w", v1, err)
	}

	ver2, err := semver.NewVersion(v2)
	if err != nil {
		return 0, fmt.Errorf("invalid version %s: %w", v2, err)
	}

	return ver1.Compare(ver2), nil
}
