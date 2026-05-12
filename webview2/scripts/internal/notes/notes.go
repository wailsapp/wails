// Package notes scrapes the WebView2 SDK release-notes markdown to extract
// version metadata and interface-to-version mappings.
//
// Source of truth:
//
//	https://raw.githubusercontent.com/MicrosoftDocs/edge-developer/master/microsoft-edge/webview2/release-notes/index.md
//
// The markdown is a per-release section; each release records its SDK
// version (the NuGet number, like 1.0.2903.40), the minimum runtime version
// ("requires WebView2 Runtime version 121.0.2277.83 or higher"), and a
// block of "Promoted to Stable" / new-API bullets that name the ICoreWebView2
// interface that gained capability in that release.
package notes

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	// SourceURL is the canonical location of the SDK release notes.
	SourceURL = "https://raw.githubusercontent.com/MicrosoftDocs/edge-developer/master/microsoft-edge/webview2/release-notes/index.md"

	// FetchTimeout bounds the release-notes fetch.
	FetchTimeout = 30 * time.Second
)

// Release is one row in the release-notes index — one SDK version.
type Release struct {
	// SDKVersion is the NuGet number, e.g. "1.0.2903.40".
	SDKVersion string

	// RuntimeVersion is the minimum runtime stated in the notes.
	RuntimeVersion string

	// URL is a deep link to the release section in the published notes.
	URL string

	// Notes are the raw bullet lines (in order) under the release header.
	Notes []string
}

var versionRE = regexp.MustCompile(`\d+\.\d+\.\d+(?:\.\d+|-prerelease)`)

// Fetch downloads the release-notes markdown.
func Fetch() ([]byte, error) {
	client := &http.Client{Timeout: FetchTimeout}
	resp, err := client.Get(SourceURL)
	if err != nil {
		return nil, fmt.Errorf("get release notes: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get release notes: HTTP %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

// Parse extracts the per-release sections from the notes markdown.
// Releases are returned newest-first (matching the source ordering).
func Parse(md []byte) ([]Release, error) {
	r := bufio.NewReader(bytes.NewReader(md))

	var releases []Release
	var cur *Release
	var inNotes bool

	for {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read line: %w", err)
		}
		l := string(line)

		// A new "## Heading" starts a non-release section: stop appending notes.
		if strings.HasPrefix(l, "## ") {
			inNotes = false
			continue
		}

		// "[NuGet package for WebView2 1.0.2903.40](...)" — start of a release.
		if strings.HasPrefix(l, "[NuGet package for WebView2 ") {
			v := versionRE.FindString(l)
			if v == "" {
				continue
			}
			if cur != nil {
				releases = append(releases, *cur)
			}
			cur = &Release{
				SDKVersion: v,
				URL:        "https://learn.microsoft.com/en-us/microsoft-edge/webview2/release-notes?tabs=win32cpp#" + strings.ReplaceAll(v, ".", ""),
			}
			inNotes = false
			continue
		}

		// "This release requires WebView2 Runtime version 121.0.2277.83 or higher."
		if strings.HasSuffix(strings.TrimSpace(l), "or higher.") {
			if cur != nil {
				cur.RuntimeVersion = versionRE.FindString(l)
				inNotes = true
			}
			continue
		}

		if cur != nil && inNotes {
			cur.Notes = append(cur.Notes, l)
		}
	}
	if cur != nil {
		releases = append(releases, *cur)
	}
	return releases, nil
}

var (
	// interfaceLinkRE matches the canonical "new interface" link in the
	// release-notes API listing, e.g.
	//   * [ICoreWebView2_19 interface](/microsoft-edge/webview2/reference/...
	// This is much stricter than scanning every textual mention: only
	// top-level interface announcements (not method links, not renames
	// in prose) carry this exact `_N interface]` suffix.
	interfaceLinkRE = regexp.MustCompile(`\[(ICoreWebView2[A-Za-z0-9_]*) interface\]`)

	// prereleaseRE detects prerelease SDK versions. Prerelease versions
	// are not shipped to end users, so they should not drive capability
	// gating — we want the stable version that first contained the API.
	prereleaseRE = regexp.MustCompile(`(?i)(prerelease|preview)`)
)

// IsPrerelease reports whether an SDK version string represents a
// prerelease (e.g. "1.0.1305-prerelease").
func IsPrerelease(sdkVersion string) bool { return prereleaseRE.MatchString(sdkVersion) }

// InterfaceMinimumVersions walks the parsed releases and returns the
// earliest stable SDK version that explicitly added each
// ICoreWebView2 interface. Prerelease versions are skipped — they
// frequently announce an interface and then re-announce it in the
// stable release that follows, and "stable" is what consumers gate on.
//
// Releases are walked oldest-first so the first stable mention wins.
func InterfaceMinimumVersions(releases []Release) map[string]string {
	out := make(map[string]string)
	// releases are newest-first; walk in reverse so the oldest wins.
	for i := len(releases) - 1; i >= 0; i-- {
		r := releases[i]
		if IsPrerelease(r.SDKVersion) {
			continue
		}
		for _, note := range r.Notes {
			for _, m := range interfaceLinkRE.FindAllStringSubmatch(note, -1) {
				name := m[1]
				if _, seen := out[name]; !seen {
					out[name] = r.SDKVersion
				}
			}
		}
	}
	return out
}
