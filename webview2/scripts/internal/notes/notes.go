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
	// interfaceRE matches an ICoreWebView2 interface name anywhere in a bullet
	// (e.g. "added the `ICoreWebView2_17` interface", "promoted ICoreWebView2Profile").
	interfaceRE = regexp.MustCompile(`(ICoreWebView2[A-Za-z0-9_]*)`)
)

// InterfaceMinimumVersions walks the parsed releases and returns the
// earliest SDK version in which each ICoreWebView2 interface name is
// mentioned. Older releases naturally win because we iterate the list
// backwards (oldest first) and only record the first occurrence.
func InterfaceMinimumVersions(releases []Release) map[string]string {
	out := make(map[string]string)
	// releases are newest-first; walk in reverse so the oldest wins.
	for i := len(releases) - 1; i >= 0; i-- {
		r := releases[i]
		for _, note := range r.Notes {
			for _, m := range interfaceRE.FindAllString(note, -1) {
				if _, seen := out[m]; !seen {
					out[m] = r.SDKVersion
				}
			}
		}
	}
	return out
}
