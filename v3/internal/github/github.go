package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// ErrNoStableRelease means the available tags contain no stable release.
var ErrNoStableRelease = errors.New("no release tag found")

const requestTimeout = 15 * time.Second

func getWithTimeout(endpoint string) (*http.Response, error) {
	// Preserve configured transport/redirect behavior without mutating the shared client.
	client := *http.DefaultClient
	if client.Timeout <= 0 || client.Timeout > requestTimeout {
		client.Timeout = requestTimeout
	}
	return client.Get(endpoint)
}

// GetReleaseNotes fetches release Markdown, including its version heading.
func GetReleaseNotes(tagVersion string) (string, error) {
	resp, err := getWithTimeout("https://api.github.com/repos/wailsapp/wails/releases/tags/" + url.PathEscape(tagVersion))
	if err != nil {
		return "", errors.New("Unable to retrieve release notes. Please check your network connection")
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return "", errors.New("No release notes found")
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Unable to retrieve release notes: HTTP %d (%s)", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("Unable to retrieve release notes. Please check your network connection")
	}

	var data struct {
		Body *string `json:"body"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", errors.New("Unable to retrieve release notes. Please check your network connection")
	}

	if data.Body == nil {
		return "", errors.New("No release notes found")
	}

	return "# Release Notes for " + tagVersion + "\n" + *data.Body, nil
}

// GetVersionTags gets the v3 release tags of the Wails repo, sorted in
// descending semver order. It pages through the tags API (100 per page, up to
// 10 pages) and ignores tags that are not v3 semver tags (v2.*, stray tags).
func GetVersionTags() ([]*SemanticVersion, error) {
	result := []*SemanticVersion{}

	for page := 1; page <= 10; page++ {
		resp, err := getWithTimeout(fmt.Sprintf("https://api.github.com/repos/wailsapp/wails/tags?per_page=100&page=%d", page))
		if err != nil {
			return result, err
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return result, fmt.Errorf("GitHub tags request failed: HTTP %d (%s)", resp.StatusCode, http.StatusText(resp.StatusCode))
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return result, err
		}

		data := []map[string]interface{}{}
		if err := json.Unmarshal(body, &data); err != nil {
			return result, err
		}
		if len(data) == 0 {
			break
		}

		for _, tag := range data {
			name, _ := tag["name"].(string)
			if !strings.HasPrefix(name, "v3.") {
				continue
			}
			semver, err := NewSemanticVersion(name)
			if err != nil {
				continue // stray, non-semver tag
			}
			result = append(result, semver)
		}
		if len(data) < 100 && !strings.Contains(resp.Header.Get("Link"), `rel="next"`) {
			break
		}
	}

	sort.Sort(sort.Reverse(SemverCollection(result)))
	return result, nil
}

// GetLatestStableRelease gets the newest non-prerelease tag: a weekly
// (vX.Y.0) or a hotfix (vX.Y.Z). This is exactly what Go resolves "@latest" to.
func GetLatestStableRelease() (result *SemanticVersion, err error) {
	tags, err := GetVersionTags()
	if err != nil {
		return nil, err
	}

	for _, tag := range tags {
		if tag.IsRelease() {
			return tag, nil
		}
	}

	return nil, ErrNoStableRelease
}

// GetLatestNightly gets the newest nightly pre-release (vX.Y.0-nightly.DATE).
func GetLatestNightly() (result *SemanticVersion, err error) {
	tags, err := GetVersionTags()
	if err != nil {
		return nil, err
	}

	for _, tag := range tags {
		if tag.IsNightly() {
			return tag, nil
		}
	}

	return nil, fmt.Errorf("no nightly tag found")
}

// GetLatestPreRelease gets the newest build of any kind (nightly, weekly,
// hotfix or legacy beta/rc).
func GetLatestPreRelease() (result *SemanticVersion, err error) {
	tags, err := GetVersionTags()
	if err != nil {
		return nil, err
	}
	if len(tags) == 0 {
		return nil, fmt.Errorf("no prerelease tag found")
	}
	return tags[0], nil
}

// IsValidTag returns true if the given string is a valid tag
func IsValidTag(tagVersion string) (bool, error) {
	tagVersion = strings.TrimPrefix(tagVersion, "v")
	if tagVersion == "" {
		return false, nil
	}
	tags, err := GetVersionTags()
	if err != nil {
		return false, err
	}

	for _, tag := range tags {
		if tag.String() == tagVersion {
			return true, nil
		}
	}
	return false, nil
}
