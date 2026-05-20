// Package keygen implements an updater.Provider against the keygen.sh REST
// API.
//
// Auth is via any keygen.sh token value (admin "admi-…", product "prod-…",
// environment "envi-…", user "user-…") set on Config.Token, or via a
// license key set on Config.LicenseKey. The two are mutually exclusive;
// Token wins when both are set.
//
// The provider's Check returns a Release whose Verification block is
// populated from keygen.sh's per-artifact SHA-512 checksum and Ed25519ph
// signature, so the framework's verifier authenticates the download
// automatically when the user has set Config.PublicKey on the Updater.
package keygen

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
)

const (
	defaultBaseURL = "https://api.keygen.sh"
	mediaType      = "application/vnd.api+json"
)

// Config configures the keygen.sh provider.
type Config struct {
	// Account is the keygen.sh account slug or UUID. Required.
	Account string

	// Product optionally scopes upgrade lookups to a single product. Strongly
	// recommended when an account has more than one published product.
	Product string

	// Package optionally narrows further to a specific package within Product.
	Package string

	// Channel selects the release channel ("stable", "rc", "beta", "alpha",
	// "dev"). Empty defaults to "stable".
	Channel string

	// Filetype optionally narrows artifact selection (e.g. "dmg", "exe").
	Filetype string

	// Token is a keygen.sh token value (e.g. "admi-...", "prod-..."). Optional.
	// Wins over LicenseKey when both are set.
	Token string

	// LicenseKey authenticates as a license. Used only if Token is empty.
	LicenseKey string

	// BaseURL overrides the API base. Empty uses https://api.keygen.sh.
	BaseURL string

	// HTTPClient lets callers inject a custom client. Nil uses a 30s-timeout
	// client. The provider always installs a CheckRedirect to strip the
	// Authorization header on cross-origin hops; caller-supplied clients
	// have their CheckRedirect wrapped, not replaced.
	HTTPClient *http.Client
}

// Provider implements updater.Provider against keygen.sh.
type Provider struct {
	cfg    Config
	client *http.Client
	base   string
}

// New returns a configured Provider. Account is required.
func New(cfg Config) (*Provider, error) {
	if strings.TrimSpace(cfg.Account) == "" {
		return nil, errors.New("keygen: Account is required")
	}
	base := strings.TrimRight(cfg.BaseURL, "/")
	if base == "" {
		base = defaultBaseURL
	}
	c := cfg.HTTPClient
	if c == nil {
		c = &http.Client{Timeout: 30 * time.Second}
	}
	c = wrapStripAuthOnRedirect(c)
	return &Provider{cfg: cfg, client: c, base: base}, nil
}

// Name implements updater.Provider.
func (p *Provider) Name() string { return "keygen.sh" }

// Check implements updater.Provider. It calls the /upgrade endpoint and,
// because that endpoint silently ignores ?include=artifacts, follows up
// with a fetch of /releases/{id}/artifacts to find one matching the
// running platform.
func (p *Provider) Check(ctx context.Context, req updater.CheckRequest) (*updater.Release, error) {
	if req.CurrentVersion == "" {
		return nil, errors.New("keygen: CurrentVersion is required")
	}

	endpoint := fmt.Sprintf("%s/v1/accounts/%s/releases/%s/upgrade",
		p.base, url.PathEscape(p.cfg.Account), url.PathEscape(req.CurrentVersion))
	q := url.Values{}
	if p.cfg.Product != "" {
		q.Set("product", p.cfg.Product)
	}
	if p.cfg.Package != "" {
		q.Set("package", p.cfg.Package)
	}
	if p.cfg.Channel != "" {
		q.Set("channel", p.cfg.Channel)
	}
	if len(q) > 0 {
		endpoint += "?" + q.Encode()
	}

	resp, err := p.do(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// fall through
	case http.StatusNoContent, http.StatusNotFound:
		// keygen.sh signals "no upgrade available" with 404 (or 204 on older
		// API versions). Either way: up to date.
		return nil, nil
	default:
		return nil, parseAPIError(resp)
	}

	// 8 MiB is generous for an upgrade envelope (typically a handful of KiB)
	// and matches the cap the appcast provider applies to feed responses.
	// Guards against an unexpectedly large or malicious response OOMing the
	// host.
	body, err := io.ReadAll(io.LimitReader(resp.Body, 8*1024*1024))
	if err != nil {
		return nil, err
	}
	var env upgradeEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("keygen: decode upgrade: %w", err)
	}
	if env.Data == nil {
		return nil, nil
	}

	rel := &updater.Release{
		Version:     env.Data.Attributes.Version,
		Channel:     env.Data.Attributes.Channel,
		Name:        env.Data.Attributes.Name,
		Notes:       env.Data.Attributes.Description,
		PublishedAt: env.Data.Attributes.Created,
		Metadata:    env.Data.Attributes.Metadata,
	}

	// Some keygen.sh deployments sideload artifacts; prefer that. Otherwise
	// fall through to /releases/{id}/artifacts.
	arts := env.Included
	if len(arts) == 0 {
		arts, err = p.listReleaseArtifacts(ctx, env.Data.ID)
		if err != nil {
			return nil, err
		}
	}

	plat, arch := req.Platform, req.Arch
	picked, ok := pickArtifact(arts, plat, arch, p.cfg.Filetype)
	if !ok {
		return nil, fmt.Errorf("keygen: release %s has no artifact for %s/%s",
			rel.Version, plat, arch)
	}
	rel.Artifact = updater.Artifact{
		Filename: picked.Attributes.Filename,
		Filetype: picked.Attributes.Filetype,
		Size:     picked.Attributes.Filesize,
		Platform: picked.Attributes.Platform,
		Arch:     picked.Attributes.Arch,
	}

	// Stash the artifact's keygen.sh ID so Download can fetch by ID rather
	// than filename — filenames are not unique across platforms (a release
	// can ship two installer.exe artifacts, one each for amd64/arm64), but
	// IDs are.
	if picked.ID != "" {
		if rel.Metadata == nil {
			rel.Metadata = map[string]any{}
		}
		rel.Metadata["keygen.artifact.id"] = picked.ID
	}

	// Map keygen.sh's checksum + signature into framework Verification.
	// keygen.sh ships them base64-encoded without padding.
	if v := buildVerification(picked.Attributes.Checksum, picked.Attributes.Signature); v != nil {
		rel.Verification = v
	}
	return rel, nil
}

// Download implements updater.Provider. The artifact endpoint 303s to the
// keygen.sh distribution backend (Cloudflare R2 or S3); the redirect-strip
// wrapper installed at construction time prevents the Authorization header
// from leaking on the cross-origin hop.
//
// When Check stashed the artifact ID under "keygen.artifact.id" the download
// targets it directly; otherwise the filename is used as a fallback. Filename
// lookups are ambiguous when a release ships two artifacts that share a name
// across platforms (e.g. installer.exe for both amd64 and arm64).
func (p *Provider) Download(ctx context.Context, rel *updater.Release, dst io.Writer, onProgress func(int64, int64)) error {
	if rel == nil {
		return errors.New("keygen: nil release")
	}
	identifier := ""
	if rel.Metadata != nil {
		if id, ok := rel.Metadata["keygen.artifact.id"].(string); ok && id != "" {
			identifier = id
		}
	}
	if identifier == "" {
		identifier = rel.Artifact.Filename
	}
	if identifier == "" {
		return errors.New("keygen: release missing artifact id and filename")
	}
	endpoint := fmt.Sprintf("%s/v1/accounts/%s/artifacts/%s",
		p.base, url.PathEscape(p.cfg.Account), url.PathEscape(identifier))

	resp, err := p.do(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return parseAPIError(resp)
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

// --- helpers ---

func (p *Provider) do(ctx context.Context, method, endpoint string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", mediaType)
	if body != nil {
		req.Header.Set("Content-Type", mediaType)
	}
	switch {
	case p.cfg.Token != "":
		req.Header.Set("Authorization", "Bearer "+p.cfg.Token)
	case p.cfg.LicenseKey != "":
		req.Header.Set("Authorization", "License "+p.cfg.LicenseKey)
	}
	return p.client.Do(req)
}

func (p *Provider) listReleaseArtifacts(ctx context.Context, releaseID string) ([]artifactResource, error) {
	endpoint := fmt.Sprintf("%s/v1/accounts/%s/releases/%s/artifacts",
		p.base, url.PathEscape(p.cfg.Account), url.PathEscape(releaseID))
	resp, err := p.do(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, parseAPIError(resp)
	}
	var env struct {
		Data []artifactResource `json:"data"`
	}
	// Cap the artifact list at 8 MiB to match the upgrade-envelope limit.
	// Releases ship a handful of platform artifacts (typically <100 KiB of
	// JSON); an unbounded read here would let a misbehaving / malicious
	// account OOM the host.
	body, err := io.ReadAll(io.LimitReader(resp.Body, 8*1024*1024))
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("keygen: decode artifacts: %w", err)
	}
	return env.Data, nil
}

func pickArtifact(arts []artifactResource, platform, arch, filetype string) (*artifactResource, bool) {
	wantPlat := strings.ToLower(platform)
	wantArch := strings.ToLower(arch)
	wantType := strings.ToLower(filetype)
	for i := range arts {
		a := &arts[i]
		if a.Type != "artifacts" {
			continue
		}
		if a.Attributes.Status != "" && a.Attributes.Status != "UPLOADED" {
			continue
		}
		if a.Attributes.Platform != "" && !platformMatches(a.Attributes.Platform, wantPlat) {
			continue
		}
		if a.Attributes.Arch != "" && !archMatches(a.Attributes.Arch, wantArch) {
			continue
		}
		if wantType != "" && a.Attributes.Filetype != "" && strings.ToLower(a.Attributes.Filetype) != wantType {
			continue
		}
		return a, true
	}
	return nil, false
}

// platformMatches reports whether a keygen.sh artifact's published platform
// (operator-defined, often "macos"/"win"/"linux") corresponds to a Go
// runtime.GOOS value. Case-insensitive.
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

// archMatches reports whether a keygen.sh artifact's published arch
// (operator-defined, often "x86_64"/"aarch64") corresponds to a Go
// runtime.GOARCH value. Case-insensitive.
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

// buildVerification maps keygen.sh's per-artifact checksum + signature
// fields into the framework's Verification struct. Both fields are base64
// (raw, no padding) per keygen's docs.
//
// Algorithms: SHA-512 for the checksum, Ed25519ph for the signature.
func buildVerification(checksumB64, sigB64 string) *updater.Verification {
	if checksumB64 == "" && sigB64 == "" {
		return nil
	}
	v := &updater.Verification{}
	if checksumB64 != "" {
		if d := decodeKeygenB64(checksumB64); d != nil {
			v.DigestAlgo = "sha512"
			v.Digest = d
		}
	}
	if sigB64 != "" {
		if s := decodeKeygenB64(sigB64); s != nil {
			v.SignatureAlgo = "ed25519ph"
			v.Signature = s
		}
	}
	if v.Digest == nil && v.Signature == nil {
		return nil
	}
	return v
}

// decodeKeygenB64 decodes a keygen.sh base64 field. The keygen docs specify
// raw (no-padding) base64, but be tolerant of either form.
func decodeKeygenB64(s string) []byte {
	if b, err := base64.RawStdEncoding.DecodeString(s); err == nil {
		return b
	}
	if b, err := base64.StdEncoding.DecodeString(s); err == nil {
		return b
	}
	return nil
}

// --- JSON:API envelopes ---

type upgradeEnvelope struct {
	Data     *releaseResource   `json:"data"`
	Included []artifactResource `json:"included"`
}

type releaseResource struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Name        string         `json:"name"`
		Description string         `json:"description"`
		Channel     string         `json:"channel"`
		Status      string         `json:"status"`
		Tag         string         `json:"tag"`
		Version     string         `json:"version"`
		Metadata    map[string]any `json:"metadata"`
		Created     time.Time      `json:"created"`
	} `json:"attributes"`
}

type artifactResource struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Filename  string `json:"filename"`
		Filetype  string `json:"filetype"`
		Filesize  int64  `json:"filesize"`
		Platform  string `json:"platform"`
		Arch      string `json:"arch"`
		Signature string `json:"signature"`
		Checksum  string `json:"checksum"`
		Status    string `json:"status"`
	} `json:"attributes"`
}

// --- error envelope ---

type apiErrorEnvelope struct {
	Errors []apiError `json:"errors"`
}

type apiError struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Code   string `json:"code"`
}

// APIError surfaces a non-2xx response from keygen.sh.
type APIError struct {
	StatusCode int
	Title      string
	Detail     string
	Code       string
}

func (e *APIError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("keygen: %d %s: %s (%s)", e.StatusCode, e.Title, e.Detail, e.Code)
	}
	return fmt.Sprintf("keygen: %d %s", e.StatusCode, e.Title)
}

func parseAPIError(resp *http.Response) error {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	apiErr := &APIError{StatusCode: resp.StatusCode, Title: resp.Status}
	var env apiErrorEnvelope
	if err := json.Unmarshal(body, &env); err == nil && len(env.Errors) > 0 {
		first := env.Errors[0]
		if first.Title != "" {
			apiErr.Title = first.Title
		}
		apiErr.Detail = first.Detail
		apiErr.Code = first.Code
	}
	return apiErr
}

// wrapStripAuthOnRedirect installs a CheckRedirect on c that drops the
// Authorization header when crossing origins (keygen.sh → R2/S3).
func wrapStripAuthOnRedirect(c *http.Client) *http.Client {
	clone := *c
	prev := clone.CheckRedirect
	clone.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) > 0 {
			prevReq := via[len(via)-1]
			if !strings.EqualFold(prevReq.URL.Host, req.URL.Host) {
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
