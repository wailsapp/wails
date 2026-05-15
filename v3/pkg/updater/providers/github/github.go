// Package github implements an updater.Provider backed by GitHub Releases.
//
// The provider hits the standard releases API and selects an asset matching
// the running platform by filename heuristics. It supports public repos
// out of the box, private repos via a personal-access token, GitHub
// Enterprise via a base-URL override, and an optional sibling-asset
// checksums file for verification when the release publisher provides one.
package github

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/pkg/updater"
	"github.com/wailsapp/wails/v3/pkg/updater/internal/semver"
)

const (
	defaultBaseURL = "https://api.github.com"
	mediaType      = "application/vnd.github+json"
)

// Config configures the GitHub provider.
type Config struct {
	// Repository is "owner/repo". Required.
	Repository string

	// Token is an optional GitHub PAT. Required for private repos; for public
	// repos it just raises the rate limit from 60 to 5000 req/hour.
	Token string

	// Prerelease, when true, walks /repos/{r}/releases (which includes
	// prereleases) instead of /repos/{r}/releases/latest (which doesn't).
	Prerelease bool

	// BaseURL overrides https://api.github.com. Useful for GitHub Enterprise:
	// pass "https://<host>/api/v3". Trailing slashes are trimmed.
	BaseURL string

	// AssetMatcher decides which asset is the right one for the running
	// platform. Nil falls back to DefaultAssetMatcher.
	AssetMatcher AssetMatcher

	// ChecksumAsset, when non-empty, names a sibling asset (e.g.
	// "checksums.txt" or "SHA256SUMS") the provider fetches and parses to
	// populate Release.Verification. The default is no checksum lookup; the
	// release ships unsigned and the framework warns at startup unless a
	// Config.PublicKey is also provided for a separate signature scheme.
	ChecksumAsset string

	// HTTPClient lets callers inject a custom client. Nil uses a 30s-timeout
	// client.
	HTTPClient *http.Client
}

// AssetMatcher returns the index of the asset (in releaseAssets) that
// should be downloaded for the supplied CheckRequest, or -1 to signal "no
// suitable asset on this release."
type AssetMatcher func(req updater.CheckRequest, assets []ReleaseAsset) int

// ReleaseAsset is the public-facing shape passed to a custom AssetMatcher.
// It mirrors the fields of the GitHub assets API response that matter for
// matching.
type ReleaseAsset struct {
	Name        string
	ContentType string
	Size        int64
	URL         string
}

// Provider implements updater.Provider against the GitHub Releases API.
type Provider struct {
	cfg    Config
	client *http.Client
	base   string
}

// New returns a configured Provider. The only required field is
// Config.Repository ("owner/repo").
func New(cfg Config) (*Provider, error) {
	if strings.TrimSpace(cfg.Repository) == "" {
		return nil, errors.New("github: Repository (\"owner/repo\") is required")
	}
	if !strings.Contains(cfg.Repository, "/") {
		return nil, fmt.Errorf("github: Repository must be in \"owner/repo\" form (got %q)", cfg.Repository)
	}
	base := strings.TrimRight(cfg.BaseURL, "/")
	if base == "" {
		base = defaultBaseURL
	}
	client := cfg.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	if cfg.AssetMatcher == nil {
		cfg.AssetMatcher = DefaultAssetMatcher
	}
	return &Provider{cfg: cfg, client: client, base: base}, nil
}

// Name implements updater.Provider.
func (p *Provider) Name() string { return "github" }

// Check implements updater.Provider. It resolves the latest release (or the
// latest including prereleases when Config.Prerelease is set), picks an
// asset matching the running platform, and decorates the Release with a
// Verification block when the publisher ships a checksum sidecar.
func (p *Provider) Check(ctx context.Context, req updater.CheckRequest) (*updater.Release, error) {
	endpoint := p.base + "/repos/" + p.cfg.Repository + "/releases/latest"
	if p.cfg.Prerelease {
		// Fetch a small page so we can skip any drafts at the top of the
		// list. The /releases endpoint includes drafts when the request is
		// authenticated, and the API does not guarantee the first item is
		// the newest *published* release once drafts are present.
		endpoint = p.base + "/repos/" + p.cfg.Repository + "/releases?per_page=10"
	}
	rel, err := p.fetchRelease(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	if rel == nil {
		return nil, nil
	}
	if !semver.IsNewer(rel.TagName, req.CurrentVersion) {
		return nil, nil
	}

	idx := p.cfg.AssetMatcher(req, asReleaseAssets(rel.Assets))
	if idx < 0 || idx >= len(rel.Assets) {
		return nil, fmt.Errorf("github: release %s has no asset for %s/%s",
			rel.TagName, req.Platform, req.Arch)
	}
	picked := rel.Assets[idx]

	out := &updater.Release{
		Version:     semver.TrimPrefix(rel.TagName),
		Channel:     channelFor(rel),
		Name:        rel.Name,
		Notes:       rel.Body,
		PublishedAt: rel.PublishedAt,
		Artifact: updater.Artifact{
			Filename: picked.Name,
			Filetype: ftypeOf(picked.Name),
			Size:     picked.Size,
			Platform: req.Platform,
			Arch:     req.Arch,
		},
		Metadata: map[string]any{
			"github.asset.id":          picked.ID,
			"github.asset.contentType": picked.ContentType,
			"github.release.tag":       rel.TagName,
			"github.release.htmlURL":   rel.HTMLURL,
		},
	}

	// Stash the download URL on Metadata so Download can find it without
	// re-querying the API.
	out.Metadata["github.asset.url"] = picked.BrowserDownloadURL

	// Optional verification via a sibling checksum asset.
	if p.cfg.ChecksumAsset != "" {
		digest, err := p.fetchChecksumFor(ctx, rel.Assets, p.cfg.ChecksumAsset, picked.Name)
		if err != nil {
			return nil, fmt.Errorf("github: load checksum sidecar: %w", err)
		}
		if digest != nil {
			out.Verification = &updater.Verification{
				DigestAlgo: "sha256",
				Digest:     digest,
			}
		}
	}

	return out, nil
}

// Download implements updater.Provider. It streams the picked asset to dst,
// reporting progress every Write.
func (p *Provider) Download(ctx context.Context, rel *updater.Release, dst io.Writer, onProgress func(written, total int64)) error {
	if rel == nil || rel.Metadata == nil {
		return errors.New("github: release missing metadata (was it produced by this provider?)")
	}
	urlStr, ok := rel.Metadata["github.asset.url"].(string)
	if !ok || urlStr == "" {
		return errors.New("github: release metadata missing asset URL")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return err
	}
	// browser_download_url returns 302 → a presigned URL. The redirect
	// chain may go to a different host; we strip Authorization on the hop.
	req.Header.Set("Accept", "application/octet-stream")
	p.setAuth(req)

	resp, err := p.followAndStrip(req)
	if err != nil {
		return fmt.Errorf("github: download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("github: download: HTTP %d", resp.StatusCode)
	}

	total := rel.Artifact.Size
	if total == 0 && resp.ContentLength > 0 {
		total = resp.ContentLength
	}
	written := int64(0)
	buf := make([]byte, 64*1024)
	for {
		n, rerr := resp.Body.Read(buf)
		if n > 0 {
			if _, werr := dst.Write(buf[:n]); werr != nil {
				return werr
			}
			written += int64(n)
			if onProgress != nil {
				onProgress(written, total)
			}
		}
		if rerr == io.EOF {
			return nil
		}
		if rerr != nil {
			return rerr
		}
	}
}

// followAndStrip executes a request, following redirects with a CheckRedirect
// wrapper that drops the Authorization header on a cross-host hop (otherwise
// the GitHub PAT would be sent to AWS, which fails the download).
//
// A caller-supplied CheckRedirect on p.client is preserved: this wrapper
// performs its strip, then delegates to the prior policy.
func (p *Provider) followAndStrip(req *http.Request) (*http.Response, error) {
	client := *p.client
	prev := client.CheckRedirect
	client.CheckRedirect = func(r *http.Request, via []*http.Request) error {
		if len(via) > 0 && !strings.EqualFold(via[len(via)-1].URL.Host, r.URL.Host) {
			r.Header.Del("Authorization")
		}
		if prev != nil {
			return prev(r, via)
		}
		if len(via) >= 10 {
			return errors.New("stopped after 10 redirects")
		}
		return nil
	}
	return client.Do(req)
}

// --- helpers ---

func (p *Provider) setAuth(req *http.Request) {
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if p.cfg.Token != "" {
		req.Header.Set("Authorization", "Bearer "+p.cfg.Token)
	}
}

func (p *Provider) fetchRelease(ctx context.Context, endpoint string) (*apiRelease, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", mediaType)
	p.setAuth(req)
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github: api request: %w", err)
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		// No release published yet — treat as up to date.
		return nil, nil
	default:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("github: api %d: %s", resp.StatusCode, body)
	}

	if p.cfg.Prerelease {
		var list []apiRelease
		if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
			return nil, fmt.Errorf("github: decode releases list: %w", err)
		}
		for i := range list {
			if list[i].Draft {
				continue
			}
			return &list[i], nil
		}
		return nil, nil
	}
	var rel apiRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, fmt.Errorf("github: decode release: %w", err)
	}
	return &rel, nil
}

func (p *Provider) fetchChecksumFor(ctx context.Context, assets []apiAsset, sidecarName, targetName string) ([]byte, error) {
	idx := -1
	for i, a := range assets {
		if a.Name == sidecarName {
			idx = i
			break
		}
	}
	if idx < 0 {
		return nil, nil
	}
	urlStr := assets[idx].BrowserDownloadURL
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/octet-stream")
	p.setAuth(req)
	resp, err := p.followAndStrip(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("checksum sidecar HTTP %d", resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	return parseChecksumLine(string(body), targetName)
}

// parseChecksumLine extracts the digest for `target` from a `sha256sum`-style
// listing. Each line is "<hex-digest>  <filename>". Filename comparison is
// done on the base name only and tolerates the "*" mode-marker before the
// filename emitted by `sha256sum --binary`.
func parseChecksumLine(body, target string) ([]byte, error) {
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		name := fields[len(fields)-1]
		name = strings.TrimPrefix(name, "*")
		name = strings.TrimPrefix(name, "./")
		if name != target {
			continue
		}
		digest, err := hex.DecodeString(fields[0])
		if err != nil {
			return nil, fmt.Errorf("malformed digest for %s: %w", target, err)
		}
		return digest, nil
	}
	return nil, nil
}

func channelFor(r *apiRelease) string {
	if r.Prerelease {
		return "prerelease"
	}
	return "stable"
}

func ftypeOf(name string) string {
	if i := strings.LastIndex(name, "."); i >= 0 {
		return strings.ToLower(name[i+1:])
	}
	return ""
}

// DefaultAssetMatcher picks the first asset whose lowercase filename
// contains both the platform string AND the architecture string. Empty
// platform/arch matches everything. Optional .sig and checksum sidecars
// are skipped automatically.
func DefaultAssetMatcher(req updater.CheckRequest, assets []ReleaseAsset) int {
	plat := strings.ToLower(req.Platform)
	arch := strings.ToLower(req.Arch)

	// First pass: skip obvious sidecars.
	for i, a := range assets {
		name := strings.ToLower(a.Name)
		if strings.HasSuffix(name, ".sig") || strings.HasSuffix(name, ".asc") {
			continue
		}
		if isChecksumName(name) {
			continue
		}
		if plat != "" && !strings.Contains(name, plat) {
			continue
		}
		if arch != "" && !containsArch(name, arch) {
			continue
		}
		return i
	}
	return -1
}

func containsArch(name, arch string) bool {
	if strings.Contains(name, arch) {
		return true
	}
	// amd64 is also commonly published as "x86_64" / "x64".
	if arch == "amd64" && (strings.Contains(name, "x86_64") || strings.Contains(name, "x64")) {
		return true
	}
	// arm64 is also commonly published as "aarch64".
	if arch == "arm64" && strings.Contains(name, "aarch64") {
		return true
	}
	// 386 is also published as "i386" / "x86" / "ia32" — but bare "x86" is
	// a substring of "x86_64", so any 64-bit alias on the name vetoes a 386
	// match outright.
	if arch == "386" {
		if strings.Contains(name, "x86_64") || strings.Contains(name, "x64") || strings.Contains(name, "amd64") {
			return false
		}
		if strings.Contains(name, "i386") || strings.Contains(name, "ia32") || strings.Contains(name, "x86") {
			return true
		}
	}
	return false
}

// isChecksumName reports whether lowerName looks like a checksum sidecar
// rather than a primary release artifact. Anchored to extensions / whole
// tokens so legitimate artifacts whose names contain a hash algorithm
// (e.g. "myapp-sha256.zip", "win-x64.tar.gz") aren't mistaken for sidecars.
func isChecksumName(lowerName string) bool {
	for _, ext := range []string{".sha256", ".sha512", ".sums", ".checksum", ".checksums"} {
		if strings.HasSuffix(lowerName, ext) {
			return true
		}
	}
	for _, token := range []string{"checksums", "sha256sums", "sha512sums"} {
		if lowerName == token || strings.HasPrefix(lowerName, token+".") {
			return true
		}
	}
	return false
}

func asReleaseAssets(in []apiAsset) []ReleaseAsset {
	out := make([]ReleaseAsset, len(in))
	for i, a := range in {
		out[i] = ReleaseAsset{
			Name:        a.Name,
			ContentType: a.ContentType,
			Size:        a.Size,
			URL:         a.BrowserDownloadURL,
		}
	}
	return out
}

// --- API shapes ---

type apiRelease struct {
	TagName     string     `json:"tag_name"`
	Name        string     `json:"name"`
	Body        string     `json:"body"`
	Prerelease  bool       `json:"prerelease"`
	Draft       bool       `json:"draft"`
	PublishedAt time.Time  `json:"published_at"`
	HTMLURL     string     `json:"html_url"`
	Assets      []apiAsset `json:"assets"`
}

type apiAsset struct {
	ID                 int64  `json:"id"`
	Name               string `json:"name"`
	ContentType        string `json:"content_type"`
	Size               int64  `json:"size"`
	BrowserDownloadURL string `json:"browser_download_url"`
}
