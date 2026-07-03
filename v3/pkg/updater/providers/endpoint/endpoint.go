// Package endpoint implements an updater.Provider against the Wails Update
// Manifest protocol: a single JSON document that describes the latest
// release and its per-platform artifacts.
//
// The protocol is deliberately host-agnostic. A static file server (S3,
// GitHub Pages, any CDN) can publish one manifest per channel listing every
// platform's artifact, while a dynamic update server can read the
// platform / arch / version query parameters this provider sends and return
// a manifest containing exactly one artifact. Both shapes are the same
// document; the provider picks the artifact matching the running platform
// either way.
//
// See the "Update Manifest Protocol" reference page in the Wails docs for
// the full wire format.
package endpoint

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/pkg/updater"
	"github.com/wailsapp/wails/v3/pkg/updater/internal/semver"
)

// schemaVersionMax is the newest manifest schemaVersion this provider
// understands. Manifests with a higher schemaVersion are rejected rather
// than half-parsed.
const schemaVersionMax = 1

// Config configures the endpoint provider.
type Config struct {
	// URL of the update manifest. Required.
	//
	// The URL may embed the placeholders {{platform}}, {{arch}}, {{version}}
	// and {{channel}}, which are substituted per-check (useful for static
	// hosting layouts like /updates/{{platform}}/{{arch}}/stable.json).
	// Whatever the placeholders do not consume is appended as query
	// parameters (platform, arch, version, and channel when set), so a
	// dynamic update server can read them without any placeholder syntax in
	// the configured URL.
	URL string

	// Channel optionally requests a release channel. Sent as the channel
	// query parameter / {{channel}} placeholder, and additionally enforced
	// client-side: a manifest that declares a different non-empty channel is
	// treated as "no update".
	Channel string

	// Headers are added to every manifest request, e.g. an Authorization
	// header for license-gated update servers. Artifact downloads reuse them
	// only when the artifact URL is on the same host as the manifest URL and
	// does not downgrade from https to http; cross-origin downloads never
	// carry the Authorization header.
	Headers map[string]string

	// HTTPClient lets callers inject a custom client. Nil uses a 30s-timeout
	// client. The provider always installs a CheckRedirect that strips the
	// Authorization header on cross-origin hops; caller-supplied clients
	// have their CheckRedirect wrapped, not replaced.
	HTTPClient *http.Client
}

// Provider implements updater.Provider against a Wails update manifest.
type Provider struct {
	cfg    Config
	client *http.Client
}

// New returns a configured Provider. URL is required.
func New(cfg Config) (*Provider, error) {
	if strings.TrimSpace(cfg.URL) == "" {
		return nil, errors.New("endpoint: URL is required")
	}
	c := cfg.HTTPClient
	if c == nil {
		c = &http.Client{Timeout: 30 * time.Second}
	}
	c = wrapStripAuthOnRedirect(c)
	return &Provider{cfg: cfg, client: c}, nil
}

// Name implements updater.Provider.
func (p *Provider) Name() string { return "endpoint" }

// Check implements updater.Provider. It fetches the manifest, compares the
// advertised version against req.CurrentVersion under semver precedence,
// and selects the artifact matching the running platform.
func (p *Provider) Check(ctx context.Context, req updater.CheckRequest) (*updater.Release, error) {
	if req.CurrentVersion == "" {
		return nil, errors.New("endpoint: CurrentVersion is required")
	}

	endpointURL, err := p.buildURL(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, endpointURL, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Accept", "application/json")
	for k, v := range p.cfg.Headers {
		httpReq.Header.Set(k, v)
	}

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("endpoint: fetch manifest: %w", err)
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode == http.StatusNoContent, resp.StatusCode == http.StatusNotFound:
		// Both signal "nothing newer for you": 204 from dynamic servers that
		// compared versions server-side, 404 from static hosts with no
		// manifest published yet.
		return nil, nil
	case resp.StatusCode < 200 || resp.StatusCode >= 300:
		return nil, fmt.Errorf("endpoint: manifest request failed: HTTP %d", resp.StatusCode)
	}

	// 8 MiB cap matches the limit the appcast and keygen providers apply to
	// feed responses; a manifest is typically well under 100 KiB.
	body, err := io.ReadAll(io.LimitReader(resp.Body, 8*1024*1024))
	if err != nil {
		return nil, err
	}
	var m manifest
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, fmt.Errorf("endpoint: decode manifest: %w", err)
	}
	if m.SchemaVersion > schemaVersionMax {
		return nil, fmt.Errorf("endpoint: manifest schemaVersion %d is newer than supported version %d",
			m.SchemaVersion, schemaVersionMax)
	}
	if m.Version == "" {
		return nil, errors.New("endpoint: manifest missing version")
	}
	if p.cfg.Channel != "" && m.Channel != "" && !strings.EqualFold(p.cfg.Channel, m.Channel) {
		return nil, nil
	}
	if !semver.IsNewer(m.Version, req.CurrentVersion) {
		return nil, nil
	}

	art, ok := pickArtifact(m.Artifacts, req.Platform, req.Arch)
	if !ok {
		return nil, fmt.Errorf("endpoint: manifest %s has no artifact for %s/%s",
			m.Version, req.Platform, req.Arch)
	}
	if art.URL == "" {
		return nil, fmt.Errorf("endpoint: artifact for %s/%s has no url", req.Platform, req.Arch)
	}
	resolved, err := resolveURL(endpointURL, art.URL)
	if err != nil {
		return nil, fmt.Errorf("endpoint: artifact url: %w", err)
	}

	filename := art.Filename
	if filename == "" {
		filename = filenameFromURL(resolved)
	}
	filetype := art.Filetype
	if filetype == "" {
		filetype = filetypeFromFilename(filename)
	}

	rel := &updater.Release{
		Version:     semver.TrimPrefix(m.Version),
		Channel:     m.Channel,
		Name:        m.Name,
		Notes:       m.Notes,
		PublishedAt: m.published(),
		Artifact: updater.Artifact{
			Filename: filename,
			Filetype: filetype,
			Size:     art.Size,
			Platform: req.Platform,
			Arch:     req.Arch,
		},
		Metadata: m.Metadata,
	}
	if rel.Metadata == nil {
		rel.Metadata = map[string]any{}
	}
	rel.Metadata["endpoint.artifact.url"] = resolved

	v, err := buildVerification(art)
	if err != nil {
		return nil, fmt.Errorf("endpoint: manifest %s, artifact for %s/%s: %w",
			m.Version, req.Platform, req.Arch, err)
	}
	if v != nil {
		rel.Verification = v
	}
	return rel, nil
}

// Download implements updater.Provider. It streams the artifact URL that
// Check stashed in the release metadata. Configured headers are sent only
// when the artifact lives on the same host as the manifest over an equal or
// upgraded scheme; the Authorization header never crosses origins or rides
// an https-to-http downgrade, either on the first request or on redirects.
func (p *Provider) Download(ctx context.Context, rel *updater.Release, dst io.Writer, onProgress func(int64, int64)) error {
	if rel == nil || rel.Metadata == nil {
		return errors.New("endpoint: release missing metadata")
	}
	urlStr, ok := rel.Metadata["endpoint.artifact.url"].(string)
	if !ok || urlStr == "" {
		return errors.New("endpoint: release metadata missing artifact URL")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/octet-stream")
	if headersAllowedFor(p.cfg.URL, urlStr) {
		for k, v := range p.cfg.Headers {
			req.Header.Set(k, v)
		}
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("endpoint: download: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("endpoint: download: HTTP %d", resp.StatusCode)
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

// --- request building ---

// buildURL substitutes {{platform}} / {{arch}} / {{version}} / {{channel}}
// placeholders into the configured URL, then appends whatever the
// placeholders did not consume as query parameters. channel is only sent
// when configured.
func (p *Provider) buildURL(req updater.CheckRequest) (string, error) {
	raw := p.cfg.URL
	params := []struct {
		name, value string
		send        bool
	}{
		{"platform", req.Platform, true},
		{"arch", req.Arch, true},
		{"version", req.CurrentVersion, true},
		{"channel", p.cfg.Channel, p.cfg.Channel != ""},
	}

	pending := url.Values{}
	for _, prm := range params {
		placeholder := "{{" + prm.name + "}}"
		if strings.Contains(raw, placeholder) {
			raw = strings.ReplaceAll(raw, placeholder, url.PathEscape(prm.value))
			continue
		}
		if prm.send {
			pending.Set(prm.name, prm.value)
		}
	}

	u, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("endpoint: invalid URL: %w", err)
	}
	if len(pending) > 0 {
		q := u.Query()
		for k, vs := range pending {
			q.Set(k, vs[0])
		}
		u.RawQuery = q.Encode()
	}
	return u.String(), nil
}

// --- manifest schema ---

type manifest struct {
	SchemaVersion int                `json:"schemaVersion"`
	Version       string             `json:"version"`
	Channel       string             `json:"channel"`
	Name          string             `json:"name"`
	Notes         string             `json:"notes"`
	PublishedAt   string             `json:"publishedAt"`
	Artifacts     []manifestArtifact `json:"artifacts"`
	Metadata      map[string]any     `json:"metadata"`
}

// published parses the RFC 3339 publishedAt field, tolerating absence and
// malformed values (both yield the zero time rather than failing the check).
func (m *manifest) published() time.Time {
	if m.PublishedAt == "" {
		return time.Time{}
	}
	if t, err := time.Parse(time.RFC3339, m.PublishedAt); err == nil {
		return t
	}
	return time.Time{}
}

type manifestArtifact struct {
	URL           string `json:"url"`
	Filename      string `json:"filename"`
	Filetype      string `json:"filetype"`
	Size          int64  `json:"size"`
	Platform      string `json:"platform"`
	Arch          string `json:"arch"`
	DigestAlgo    string `json:"digestAlgo"`
	Digest        string `json:"digest"`
	SignatureAlgo string `json:"signatureAlgo"`
	Signature     string `json:"signature"`
}

// pickArtifact returns the first artifact matching platform and arch.
// An artifact with an empty platform or arch matches any value, so a single
// universal artifact (e.g. a JS bundle or a fat binary) needs no per-platform
// entries. Document order breaks ties, letting publishers put preferred
// filetypes first.
func pickArtifact(arts []manifestArtifact, platform, arch string) (*manifestArtifact, bool) {
	wantPlat := strings.ToLower(platform)
	wantArch := strings.ToLower(arch)
	for i := range arts {
		a := &arts[i]
		if a.Platform != "" && !platformMatches(a.Platform, wantPlat) {
			continue
		}
		if a.Arch != "" && !archMatches(a.Arch, wantArch) {
			continue
		}
		return a, true
	}
	return nil, false
}

// platformMatches reports whether a manifest artifact's platform corresponds
// to a Go runtime.GOOS value. Manifests should publish GOOS values, but the
// common vernacular aliases are accepted. Case-insensitive.
func platformMatches(published, want string) bool {
	p := strings.ToLower(published)
	if p == want {
		return true
	}
	switch want {
	case "darwin":
		return p == "macos" || p == "mac" || p == "osx"
	case "windows":
		return p == "win" || p == "win32" || p == "win64"
	}
	return false
}

// archMatches reports whether a manifest artifact's arch corresponds to a Go
// runtime.GOARCH value. Manifests should publish GOARCH values, but the
// common vernacular aliases are accepted. Case-insensitive.
func archMatches(published, want string) bool {
	p := strings.ToLower(published)
	if p == want {
		return true
	}
	switch want {
	case "amd64":
		return p == "x86_64" || p == "x64"
	case "arm64":
		return p == "aarch64"
	case "386":
		return p == "i386" || p == "x86" || p == "ia32"
	}
	return false
}

// buildVerification maps an artifact's digest and signature fields into the
// framework's Verification struct. Both are base64; raw (unpadded) and
// standard encodings are accepted. Malformed values are errors, not
// omissions: silently dropping a digest or signature the publisher wrote
// would downgrade verification the client was meant to perform.
func buildVerification(a *manifestArtifact) (*updater.Verification, error) {
	if a.Digest == "" && a.Signature == "" {
		return nil, nil
	}
	v := &updater.Verification{}
	if a.Digest != "" {
		d := decodeB64(a.Digest)
		if d == nil {
			return nil, errors.New("artifact digest is not valid base64")
		}
		v.DigestAlgo = strings.ToLower(a.DigestAlgo)
		v.Digest = d
	}
	if a.Signature != "" {
		if a.SignatureAlgo == "" {
			return nil, errors.New("artifact has a signature but no signatureAlgo")
		}
		s := decodeB64(a.Signature)
		if s == nil {
			return nil, errors.New("artifact signature is not valid base64")
		}
		v.SignatureAlgo = strings.ToLower(a.SignatureAlgo)
		v.Signature = s
	}
	return v, nil
}

// --- helpers ---

func decodeB64(s string) []byte {
	if b, err := base64.RawStdEncoding.DecodeString(s); err == nil {
		return b
	}
	if b, err := base64.StdEncoding.DecodeString(s); err == nil {
		return b
	}
	return nil
}

// resolveURL resolves ref (possibly relative) against base, so a static
// manifest can point at sibling artifact files with plain filenames.
func resolveURL(base, ref string) (string, error) {
	b, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	r, err := url.Parse(ref)
	if err != nil {
		return "", err
	}
	resolved := b.ResolveReference(r)
	if resolved.Scheme != "http" && resolved.Scheme != "https" {
		return "", fmt.Errorf("unsupported scheme %q", resolved.Scheme)
	}
	return resolved.String(), nil
}

// headersAllowedFor reports whether the configured headers may accompany a
// request to target: the host must match the manifest's, and the scheme must
// not downgrade from https to http. Without the scheme rule, credentials
// configured for a secure endpoint could be replayed in cleartext to an
// http URL on the same host.
func headersAllowedFor(manifestURL, target string) bool {
	m, err := url.Parse(manifestURL)
	if err != nil {
		return false
	}
	t, err := url.Parse(target)
	if err != nil {
		return false
	}
	if !strings.EqualFold(m.Host, t.Host) {
		return false
	}
	return !(strings.EqualFold(m.Scheme, "https") && strings.EqualFold(t.Scheme, "http"))
}

func filenameFromURL(u string) string {
	if i := strings.LastIndex(u, "/"); i >= 0 && i+1 < len(u) {
		name := u[i+1:]
		if j := strings.Index(name, "?"); j >= 0 {
			name = name[:j]
		}
		return name
	}
	return u
}

func filetypeFromFilename(name string) string {
	if i := strings.LastIndex(name, "."); i >= 0 {
		return strings.ToLower(name[i+1:])
	}
	return ""
}

// wrapStripAuthOnRedirect installs a CheckRedirect on c that drops the
// Authorization header when a redirect crosses hosts (e.g. update server to
// CDN or object storage) or downgrades from https to http on the same host.
func wrapStripAuthOnRedirect(c *http.Client) *http.Client {
	clone := *c
	prev := clone.CheckRedirect
	clone.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) > 0 {
			last := via[len(via)-1].URL
			crossHost := !strings.EqualFold(last.Host, req.URL.Host)
			downgrade := strings.EqualFold(last.Scheme, "https") && strings.EqualFold(req.URL.Scheme, "http")
			if crossHost || downgrade {
				req.Header.Del("Authorization")
			}
		}
		if prev != nil {
			return prev(req, via)
		}
		if len(via) >= 10 {
			return errors.New("stopped after 10 redirects")
		}
		return nil
	}
	return &clone
}
