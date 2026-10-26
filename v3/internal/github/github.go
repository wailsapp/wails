package github

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/glamour/styles"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/charmbracelet/glamour"
)

func GetReleaseNotes(tagVersion string, noColour bool) string {
	resp, err := http.Get("https://api.github.com/repos/wailsapp/wails/releases/tags/" + url.PathEscape(tagVersion))
	if err != nil {
		return "Unable to retrieve release notes. Please check your network connection"
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Unable to retrieve release notes. Please check your network connection"
	}

	data := map[string]interface{}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "Unable to retrieve release notes. Please check your network connection"
	}

	if data["body"] == nil {
		return "No release notes found"
	}

	result := "# Release Notes for " + tagVersion + "\n" + data["body"].(string)
	var renderer *glamour.TermRenderer

	if noColour {
		renderer, err = glamour.NewTermRenderer(glamour.WithStyles(styles.NoTTYStyleConfig))
	} else {
		renderer, err = glamour.NewTermRenderer(glamour.WithAutoStyle())
	}
	if err != nil {
		return result
	}
	result, err = renderer.Render(result)
	if err != nil {
		return err.Error()
	}
	return result
}

// GetVersionTags gets the v3 release tags of the Wails repo, sorted in
// descending semver order. It pages through the tags API (100 per page, up to
// 10 pages) and ignores tags that are not v3 semver tags (v2.*, stray tags).
func GetVersionTags() ([]*SemanticVersion, error) {
	result := []*SemanticVersion{}

	for page := 1; page <= 10; page++ {
		resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/wailsapp/wails/tags?per_page=100&page=%d", page))
		if err != nil {
			return result, err
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return result, err
		}
		if resp.StatusCode != http.StatusOK {
			return result, fmt.Errorf("GitHub tags API returned HTTP %d", resp.StatusCode)
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
		if len(data) < 100 {
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

	return nil, fmt.Errorf("no release tag found")
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
	if tagVersion[0] == 'v' {
		tagVersion = tagVersion[1:]
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
