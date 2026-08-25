// Package appcast implements an updater.Provider against a Sparkle-style
// AppCast RSS feed.
//
// The provider speaks the canonical Sparkle XML vocabulary so existing
// Sparkle/WinSparkle feeds work without modification: items are picked by
// shortVersionString (falling back to sparkle:version), the download URL
// comes from <enclosure>, and the Ed25519 signature in
// sparkle:edSignature populates Release.Verification when present.
//
// Sparkle 1's DSA signatures (sparkle:dsaSignature) are not supported in
// v1 — projects on that signing scheme should rotate to EdDSA (Sparkle 2)
// or supply Config.PublicKey via a different verification path.
package appcast

import (
	"context"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/pkg/updater"
	"github.com/wailsapp/wails/v3/pkg/updater/internal/semver"
)

// Config configures the AppCast provider.
type Config struct {
	// URL of the AppCast XML feed. Required.
	URL string

	// Channel optionally filters items by sparkle:channel (Sparkle 2.x
	// supports channels on the feed). Empty matches every item.
	Channel string

	// HTTPClient lets callers inject a custom client. Nil uses a 30s-timeout
	// client with a redirect-stripping wrapper.
	HTTPClient *http.Client
}

// Provider implements updater.Provider against a Sparkle AppCast feed.
type Provider struct {
	cfg    Config
	client *http.Client
}

// New returns a configured Provider. URL is required.
func New(cfg Config) (*Provider, error) {
	if strings.TrimSpace(cfg.URL) == "" {
		return nil, errors.New("appcast: URL is required")
	}
	c := cfg.HTTPClient
	if c == nil {
		c = &http.Client{Timeout: 30 * time.Second}
	}
	c = wrapStripAuthOnRedirect(c)
	return &Provider{cfg: cfg, client: c}, nil
}

// Name implements updater.Provider.
func (p *Provider) Name() string { return "appcast" }

// Check implements updater.Provider. It fetches the AppCast feed, parses
// the channel/items, and selects the newest item whose version is greater
// than the supplied CurrentVersion (and whose os matches the running
// platform, when sparkle:os is present on the item).
func (p *Provider) Check(ctx context.Context, req updater.CheckRequest) (*updater.Release, error) {
	if req.CurrentVersion == "" {
		return nil, errors.New("appcast: CurrentVersion is required")
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, p.cfg.URL, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Accept", "application/rss+xml, application/xml, text/xml")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("appcast: fetch: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("appcast: HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 8*1024*1024))
	if err != nil {
		return nil, err
	}

	feed, err := parseFeed(body)
	if err != nil {
		return nil, fmt.Errorf("appcast: parse: %w", err)
	}

	best := pickBestItem(feed.Items, req, p.cfg.Channel)
	if best == nil {
		return nil, nil
	}
	if !semver.IsNewer(best.shortVersion(), req.CurrentVersion) {
		return nil, nil
	}
	if best.Enclosure.URL == "" {
		return nil, fmt.Errorf("appcast: matched item %q has no <enclosure>", best.Title)
	}

	rel := &updater.Release{
		Version:     best.shortVersion(),
		Channel:     best.SparkleChannel,
		Name:        best.Title,
		Notes:       best.Description,
		PublishedAt: best.published(),
		Artifact: updater.Artifact{
			Filename: filenameFromURL(best.Enclosure.URL),
			Filetype: filetypeFromURL(best.Enclosure.URL),
			Size:     best.Enclosure.Length,
			Platform: req.Platform,
			Arch:     req.Arch,
		},
		Metadata: map[string]any{
			"appcast.itemTitle":            best.Title,
			"appcast.enclosure.url":        best.Enclosure.URL,
			"appcast.enclosure.type":       best.Enclosure.Type,
			"appcast.sparkleVersion":       best.SparkleVersion,
			"appcast.sparkleShortVersion":  best.SparkleShortVersion,
			"appcast.minimumSystemVersion": best.SparkleMinSystemVersion,
			"appcast.releaseNotesLink":     best.ReleaseNotesLink,
		},
	}

	if best.Enclosure.SparkleEdSignature != "" {
		if sig := decodeB64(best.Enclosure.SparkleEdSignature); sig != nil {
			rel.Verification = &updater.Verification{
				SignatureAlgo: "ed25519",
				Signature:     sig,
			}
		}
	}
	return rel, nil
}

// Download implements updater.Provider. It streams the enclosure URL to
// dst, reporting progress.
func (p *Provider) Download(ctx context.Context, rel *updater.Release, dst io.Writer, onProgress func(int64, int64)) error {
	if rel == nil || rel.Metadata == nil {
		return errors.New("appcast: release missing metadata")
	}
	urlStr, ok := rel.Metadata["appcast.enclosure.url"].(string)
	if !ok || urlStr == "" {
		return errors.New("appcast: release metadata missing enclosure URL")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/octet-stream")
	// Enclosure URLs point at a publisher CDN, often on a different host
	// from the feed. Even though Check never sets Authorization on its
	// requests, defensively scrub it here so a caller-supplied default on
	// http.Request can't leak credentials cross-origin on the initial hop.
	// The redirect wrapper handles subsequent hops.
	req.Header.Del("Authorization")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("appcast: download: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("appcast: download: HTTP %d", resp.StatusCode)
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

// --- XML schema ---

type rssFeed struct {
	XMLName xml.Name `xml:"rss"`
	Items   []item   `xml:"channel>item"`
}

type item struct {
	Title                   string    `xml:"title"`
	Description             string    `xml:"description"`
	PubDate                 string    `xml:"pubDate"`
	ReleaseNotesLink        string    `xml:"http://www.andymatuschak.org/xml-namespaces/sparkle releaseNotesLink"`
	SparkleVersion          string    `xml:"http://www.andymatuschak.org/xml-namespaces/sparkle version"`
	SparkleShortVersion     string    `xml:"http://www.andymatuschak.org/xml-namespaces/sparkle shortVersionString"`
	SparkleMinSystemVersion string    `xml:"http://www.andymatuschak.org/xml-namespaces/sparkle minimumSystemVersion"`
	SparkleOS               string    `xml:"http://www.andymatuschak.org/xml-namespaces/sparkle os"`
	SparkleChannel          string    `xml:"http://www.andymatuschak.org/xml-namespaces/sparkle channel"`
	Enclosure               enclosure `xml:"enclosure"`
}

func (i *item) shortVersion() string {
	if i.SparkleShortVersion != "" {
		return strings.TrimPrefix(strings.TrimPrefix(i.SparkleShortVersion, "v"), "V")
	}
	return strings.TrimPrefix(strings.TrimPrefix(i.SparkleVersion, "v"), "V")
}

func (i *item) published() time.Time {
	if i.PubDate == "" {
		return time.Time{}
	}
	for _, layout := range []string{time.RFC1123Z, time.RFC1123, time.RFC822Z, time.RFC822} {
		if t, err := time.Parse(layout, i.PubDate); err == nil {
			return t
		}
	}
	return time.Time{}
}

type enclosure struct {
	URL                 string `xml:"url,attr"`
	Length              int64  `xml:"length,attr"`
	Type                string `xml:"type,attr"`
	SparkleVersion      string `xml:"http://www.andymatuschak.org/xml-namespaces/sparkle version,attr"`
	SparkleEdSignature  string `xml:"http://www.andymatuschak.org/xml-namespaces/sparkle edSignature,attr"`
	SparkleDSASignature string `xml:"http://www.andymatuschak.org/xml-namespaces/sparkle dsaSignature,attr"`
	SparkleOS           string `xml:"http://www.andymatuschak.org/xml-namespaces/sparkle os,attr"`
}

func parseFeed(body []byte) (*rssFeed, error) {
	var f rssFeed
	if err := xml.Unmarshal(body, &f); err != nil {
		return nil, err
	}
	return &f, nil
}

// pickBestItem returns the newest matching item by parsed semver order.
// "Matching" = (1) channel filter matches if set (an unlabelled item never
// matches a specific channel — Sparkle's contract is that items opt *in* to
// a channel, not the reverse), (2) sparkle:os matches req.Platform when
// sparkle:os is set on the item or enclosure.
func pickBestItem(items []item, req updater.CheckRequest, channel string) *item {
	plat := strings.ToLower(req.Platform)
	var best *item
	for idx := range items {
		it := &items[idx]
		if channel != "" && it.SparkleChannel != channel {
			continue
		}
		// sparkle:os may sit either on the <item> or the <enclosure>. Both
		// are accepted; either must match if present.
		os := strings.ToLower(it.SparkleOS)
		if os == "" {
			os = strings.ToLower(it.Enclosure.SparkleOS)
		}
		if os != "" && plat != "" && !platformMatches(os, plat) {
			continue
		}
		// Tied versions preserve document order: equal shortVersion()s keep
		// the first-seen item, matching Sparkle. A corrected later entry
		// with the same version won't override an earlier one.
		if best == nil || semver.IsNewer(it.shortVersion(), best.shortVersion()) {
			best = it
		}
	}
	return best
}

// platformMatches reconciles sparkle's loose OS naming with Go's GOOS.
// Sparkle uses "macos" / "windows", we get "darwin" / "windows" / "linux".
func platformMatches(sparkleOS, plat string) bool {
	switch sparkleOS {
	case "macos", "mac", "osx":
		return plat == "darwin"
	}
	return sparkleOS == plat
}

func filenameFromURL(u string) string {
	if i := strings.LastIndex(u, "/"); i >= 0 && i+1 < len(u) {
		// Strip query.
		name := u[i+1:]
		if j := strings.Index(name, "?"); j >= 0 {
			name = name[:j]
		}
		return name
	}
	return u
}

func filetypeFromURL(u string) string {
	name := filenameFromURL(u)
	if i := strings.LastIndex(name, "."); i >= 0 {
		return strings.ToLower(name[i+1:])
	}
	return ""
}

func decodeB64(s string) []byte {
	if b, err := base64.RawStdEncoding.DecodeString(s); err == nil {
		return b
	}
	if b, err := base64.StdEncoding.DecodeString(s); err == nil {
		return b
	}
	return nil
}

func wrapStripAuthOnRedirect(c *http.Client) *http.Client {
	clone := *c
	prev := clone.CheckRedirect
	clone.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) > 0 && !strings.EqualFold(via[len(via)-1].URL.Host, req.URL.Host) {
			req.Header.Del("Authorization")
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
