package selfupdate

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// maxAPIResponseSize is the maximum size of a GitHub API response (10MB).
const maxAPIResponseSize = 10 * 1024 * 1024

// validGitHubName matches valid GitHub owner/repo names.
var validGitHubName = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

func init() {
	RegisterProvider("github", func() UpdateProvider {
		return NewGitHubProvider()
	})
}

// GitHubProvider implements UpdateProvider for GitHub Releases.
type GitHubProvider struct {
	config     *ProviderConfig
	httpClient *http.Client
	owner      string
	repo       string
	token      string
	baseURL    string
}

// NewGitHubProvider creates a new GitHub provider.
func NewGitHubProvider() *GitHubProvider {
	return &GitHubProvider{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://api.github.com",
	}
}

// Name returns the provider identifier.
func (p *GitHubProvider) Name() string {
	return "github"
}

// Configure initializes the provider with the given configuration.
func (p *GitHubProvider) Configure(_ context.Context, config *ProviderConfig) error {
	if config == nil {
		return fmt.Errorf("config is required")
	}

	// Reset fields to avoid state leakage from previous Configure calls
	p.owner = ""
	p.repo = ""
	p.token = ""

	p.config = config

	// Extract GitHub-specific settings
	if config.Settings != nil {
		if owner, ok := config.Settings["owner"].(string); ok {
			p.owner = owner
		}
		if repo, ok := config.Settings["repo"].(string); ok {
			p.repo = repo
		}
		if token, ok := config.Settings["token"].(string); ok {
			p.token = token
		}
		if baseURL, ok := config.Settings["baseURL"].(string); ok && baseURL != "" {
			p.baseURL = strings.TrimSuffix(baseURL, "/")
		}
	}

	// Validate required fields
	if p.owner == "" {
		return fmt.Errorf("github: owner is required in settings")
	}
	if p.repo == "" {
		return fmt.Errorf("github: repo is required in settings")
	}

	// Validate owner/repo format to prevent URL injection (M5)
	if !validGitHubName.MatchString(p.owner) {
		return fmt.Errorf("github: invalid owner format: %q", p.owner)
	}
	if !validGitHubName.MatchString(p.repo) {
		return fmt.Errorf("github: invalid repo format: %q", p.repo)
	}

	// Validate baseURL is a valid HTTPS URL
	if p.baseURL != "" {
		parsed, err := url.Parse(p.baseURL)
		if err != nil {
			return fmt.Errorf("github: invalid baseURL: %w", err)
		}
		if parsed.Scheme != "https" {
			return fmt.Errorf("github: baseURL must use HTTPS")
		}
	}

	return nil
}

// CheckForUpdate queries GitHub for available updates.
func (p *GitHubProvider) CheckForUpdate(ctx context.Context, opts *CheckOptions) (*UpdateResult, error) {
	if p.config == nil {
		return nil, fmt.Errorf("provider not configured")
	}

	// Fetch latest release
	release, err := p.getLatestRelease(ctx, opts)
	if err != nil {
		return nil, err
	}

	if release == nil {
		return &UpdateResult{
			UpdateAvailable: false,
			CurrentVersion:  p.config.CurrentVersion,
		}, nil
	}

	// Find matching asset
	asset := p.findMatchingAsset(release)
	if asset == nil {
		return nil, fmt.Errorf("no matching asset found for %s/%s",
			p.getPlatform(), p.getArch())
	}

	// Compare versions
	currentVersion := strings.TrimPrefix(p.config.CurrentVersion, "v")
	latestVersion := strings.TrimPrefix(release.TagName, "v")

	updateAvailable := compareVersions(currentVersion, latestVersion) < 0

	result := &UpdateResult{
		UpdateAvailable: updateAvailable,
		Version:         latestVersion,
		CurrentVersion:  currentVersion,
		ReleaseNotes:    release.Body,
		ReleaseURL:      release.HTMLURL,
		DownloadURL:     asset.BrowserDownloadURL,
		Size:            asset.Size,
		AssetName:       asset.Name,
		Channel:         p.config.Channel,
		Metadata: map[string]any{
			"release_id": release.ID,
			"asset_id":   asset.ID,
			"tag_name":   release.TagName,
		},
	}

	if !release.PublishedAt.IsZero() {
		result.ReleaseDate = release.PublishedAt
	}

	return result, nil
}

// DownloadUpdate downloads the update binary.
func (p *GitHubProvider) DownloadUpdate(ctx context.Context, update *UpdateResult, progress ProgressFunc) (io.ReadCloser, error) {
	if update == nil || update.DownloadURL == "" {
		return nil, fmt.Errorf("no download URL available")
	}

	// Validate download URL (C2)
	if err := p.validateDownloadURL(update.DownloadURL); err != nil {
		return nil, fmt.Errorf("invalid download URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, update.DownloadURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Accept binary data
	req.Header.Set("Accept", "application/octet-stream")
	if p.token != "" {
		req.Header.Set("Authorization", "token "+p.token)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Notify progress start
	if progress != nil {
		progress(&ProgressInfo{
			State:      "started",
			TotalBytes: update.Size,
		})
	}

	// Wrap reader with progress tracking (L2: throttled)
	if progress != nil {
		return &progressReader{
			reader:       resp.Body,
			totalBytes:   update.Size,
			progress:     progress,
			startTime:    time.Now(),
			lastProgress: time.Now(),
		}, nil
	}

	return resp.Body, nil
}

// validateDownloadURL checks that a download URL is safe to fetch from (C2).
func (p *GitHubProvider) validateDownloadURL(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("malformed URL: %w", err)
	}

	if parsed.Scheme != "https" {
		return fmt.Errorf("URL must use HTTPS, got %q", parsed.Scheme)
	}

	// Build list of allowed hosts
	allowedHosts := []string{
		"github.com",
		"objects.githubusercontent.com",
		"github-releases.githubusercontent.com",
	}

	// If a custom baseURL is configured, allow its host too
	if p.baseURL != "" && p.baseURL != "https://api.github.com" {
		if base, err := url.Parse(p.baseURL); err == nil {
			allowedHosts = append(allowedHosts, base.Hostname())
		}
	}

	host := parsed.Hostname()
	for _, allowed := range allowedHosts {
		if host == allowed || strings.HasSuffix(host, "."+allowed) {
			return nil
		}
	}

	return fmt.Errorf("host %q is not an allowed GitHub domain", host)
}

// VerifyUpdate validates the downloaded update.
func (p *GitHubProvider) VerifyUpdate(_ context.Context, update *UpdateResult, data io.Reader) error {
	// Read data with size limit (C3)
	content, err := io.ReadAll(io.LimitReader(data, MaxVerifySize+1))
	if err != nil {
		return fmt.Errorf("failed to read update data: %w", err)
	}
	if int64(len(content)) > MaxVerifySize {
		return fmt.Errorf("update data exceeds maximum size of %d bytes", MaxVerifySize)
	}

	// Verify checksum if available
	if update.Checksum != "" {
		if err := VerifyChecksum(content, update.Checksum); err != nil {
			return err
		}
	}

	// Verify signature - mandatory when public key is configured (C1)
	if p.config.PublicKey != "" {
		if update.Signature == "" {
			return fmt.Errorf("signature is required: public key is configured but release has no signature")
		}
		verifier, err := NewVerifier(p.config.PublicKey)
		if err != nil {
			return fmt.Errorf("failed to create verifier: %w", err)
		}
		if err := verifier.VerifySignature(content, update.Signature); err != nil {
			return err
		}
	}

	return nil
}

// Close releases any resources.
func (p *GitHubProvider) Close() error {
	return nil
}

// GitHub API types

type githubRelease struct {
	ID          int64         `json:"id"`
	TagName     string        `json:"tag_name"`
	Name        string        `json:"name"`
	Body        string        `json:"body"`
	Draft       bool          `json:"draft"`
	Prerelease  bool          `json:"prerelease"`
	PublishedAt time.Time     `json:"published_at"`
	HTMLURL     string        `json:"html_url"`
	Assets      []githubAsset `json:"assets"`
}

type githubAsset struct {
	ID                 int64  `json:"id"`
	Name               string `json:"name"`
	Size               int64  `json:"size"`
	BrowserDownloadURL string `json:"browser_download_url"`
	ContentType        string `json:"content_type"`
}

// getLatestRelease fetches the latest release from GitHub.
func (p *GitHubProvider) getLatestRelease(ctx context.Context, opts *CheckOptions) (*githubRelease, error) {
	// If we need prereleases, we have to list all releases
	if opts != nil && opts.IncludePrerelease {
		return p.getLatestReleaseFromList(ctx, true)
	}

	// Otherwise, use the /latest endpoint
	apiURL := fmt.Sprintf("%s/repos/%s/%s/releases/latest", p.baseURL, p.owner, p.repo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	p.setRequestHeaders(req)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // No releases
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	// Size-limited reader for API response (M6)
	limitedBody := io.LimitReader(resp.Body, maxAPIResponseSize)

	var release githubRelease
	if err := json.NewDecoder(limitedBody).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to parse release: %w", err)
	}

	return &release, nil
}

// getLatestReleaseFromList fetches releases and finds the latest one.
func (p *GitHubProvider) getLatestReleaseFromList(ctx context.Context, includePrereleases bool) (*githubRelease, error) {
	apiURL := fmt.Sprintf("%s/repos/%s/%s/releases?per_page=20", p.baseURL, p.owner, p.repo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	p.setRequestHeaders(req)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	// Size-limited reader for API response (M6)
	limitedBody := io.LimitReader(resp.Body, maxAPIResponseSize)

	var releases []githubRelease
	if err := json.NewDecoder(limitedBody).Decode(&releases); err != nil {
		return nil, fmt.Errorf("failed to parse releases: %w", err)
	}

	for _, release := range releases {
		if release.Draft {
			continue
		}
		if !includePrereleases && release.Prerelease {
			continue
		}
		return &release, nil
	}

	return nil, nil
}

// setRequestHeaders sets common headers for GitHub API requests.
func (p *GitHubProvider) setRequestHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "wails-selfupdate")
	if p.token != "" {
		req.Header.Set("Authorization", "token "+p.token)
	}
}

// findMatchingAsset finds the release asset that matches the current platform.
func (p *GitHubProvider) findMatchingAsset(release *githubRelease) *githubAsset {
	if release == nil || len(release.Assets) == 0 {
		return nil
	}

	// Build list of asset names
	assetNames := make([]string, len(release.Assets))
	for i, asset := range release.Assets {
		assetNames[i] = asset.Name
	}

	// Use pattern matching to find the asset
	pattern := p.config.AssetPattern
	if pattern == "" {
		pattern = DefaultAssetPattern
	}

	vars := PatternVariables{
		Name:    p.repo,
		Version: strings.TrimPrefix(release.TagName, "v"),
		GOOS:    p.getPlatform(),
		GOARCH:  p.getArch(),
		Variant: p.config.Variant,
	}

	matchedName, _ := FindMatchingAsset(assetNames, pattern, vars)
	if matchedName == "" {
		return nil
	}

	// Find the asset object
	for i := range release.Assets {
		if release.Assets[i].Name == matchedName {
			return &release.Assets[i]
		}
	}

	return nil
}

// getPlatform returns the target platform.
func (p *GitHubProvider) getPlatform() string {
	if p.config != nil && p.config.Platform != "" {
		return p.config.Platform
	}
	return runtime.GOOS
}

// getArch returns the target architecture.
func (p *GitHubProvider) getArch() string {
	if p.config != nil && p.config.Arch != "" {
		return p.config.Arch
	}
	return runtime.GOARCH
}

// progressReader wraps an io.ReadCloser to track download progress.
type progressReader struct {
	reader          io.ReadCloser
	totalBytes      int64
	downloadedBytes int64
	progress        ProgressFunc
	startTime       time.Time
	lastProgress    time.Time
}

// progressThrottle is the minimum interval between progress callbacks (L2).
const progressThrottle = 150 * time.Millisecond

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	if n > 0 {
		pr.downloadedBytes += int64(n)

		// Throttle progress updates (L2)
		now := time.Now()
		if now.Sub(pr.lastProgress) >= progressThrottle {
			pr.lastProgress = now

			elapsed := now.Sub(pr.startTime).Seconds()
			var speed float64
			if elapsed > 0 {
				speed = float64(pr.downloadedBytes) / elapsed
			}

			var percentage float64
			if pr.totalBytes > 0 {
				percentage = float64(pr.downloadedBytes) / float64(pr.totalBytes) * 100
			}

			pr.progress(&ProgressInfo{
				State:           "downloading",
				TotalBytes:      pr.totalBytes,
				DownloadedBytes: pr.downloadedBytes,
				Percentage:      percentage,
				BytesPerSecond:  speed,
			})
		}
	}

	if err == io.EOF {
		pr.progress(&ProgressInfo{
			State:           "finished",
			TotalBytes:      pr.totalBytes,
			DownloadedBytes: pr.downloadedBytes,
			Percentage:      100,
		})
	}

	return n, err
}

func (pr *progressReader) Close() error {
	return pr.reader.Close()
}

// compareVersions compares two semantic versions.
// Returns -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2.
// Handles pre-release suffixes: 1.0.0-alpha < 1.0.0 (M1).
func compareVersions(v1, v2 string) int {
	// Split off pre-release suffix
	base1, pre1 := splitPrerelease(v1)
	base2, pre2 := splitPrerelease(v2)

	// Compare base version parts
	parts1 := strings.Split(base1, ".")
	parts2 := strings.Split(base2, ".")

	maxLen := max(len(parts1), len(parts2))

	for i := range maxLen {
		var n1, n2 int
		if i < len(parts1) {
			n1, _ = strconv.Atoi(parts1[i])
		}
		if i < len(parts2) {
			n2, _ = strconv.Atoi(parts2[i])
		}

		if n1 < n2 {
			return -1
		}
		if n1 > n2 {
			return 1
		}
	}

	// Base versions are equal - compare pre-release
	// Per semver: version with pre-release has lower precedence than release
	if pre1 == "" && pre2 != "" {
		return 1 // v1 is release, v2 is pre-release
	}
	if pre1 != "" && pre2 == "" {
		return -1 // v1 is pre-release, v2 is release
	}
	// Both have pre-release or both don't: compare lexicographically
	if pre1 < pre2 {
		return -1
	}
	if pre1 > pre2 {
		return 1
	}

	return 0
}

// splitPrerelease splits a version into base and pre-release components.
// "1.0.0-beta.1" -> ("1.0.0", "beta.1")
// "1.0.0" -> ("1.0.0", "")
func splitPrerelease(version string) (base, prerelease string) {
	if idx := strings.IndexByte(version, '-'); idx >= 0 {
		return version[:idx], version[idx+1:]
	}
	return version, ""
}
