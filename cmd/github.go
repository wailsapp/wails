package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
)

// GitHubHelper is a utility class for interacting with GitHub
type GitHubHelper struct {
}

// NewGitHubHelper returns a new GitHub Helper
func NewGitHubHelper() *GitHubHelper {
	return &GitHubHelper{}
}

// GetVersionTags gets the list of tags on the Wails repo
// It returns a list of sorted tags in descending order
func (g *GitHubHelper) GetVersionTags() ([]*SemanticVersion, error) {

	result := []*SemanticVersion{}
	var err error

	resp, err := http.Get("https://api.github.com/repos/wailsapp/wails/releases")
	if err != nil {
		return result, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	data := []map[string]interface{}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return result, err
	}

	// Convert tag data to Version structs
	for _, tag := range data {
		version := tag["name"].(string)
		semver, err := NewSemanticVersion(version)
		if err != nil {
			return result, err
		}
		result = append(result, semver)
	}

	// Reverse Sort
	sort.Sort(sort.Reverse(SemverCollection(result)))

	return result, err
}

// GetLatestStableRelease gets the latest stable release on GitHub
func (g *GitHubHelper) GetLatestStableRelease() (result *SemanticVersion, err error) {

	tags, err := g.GetVersionTags()
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

// GetLatestPreRelease gets the latest prerelease on GitHub
func (g *GitHubHelper) GetLatestPreRelease() (result *SemanticVersion, err error) {

	tags, err := g.GetVersionTags()
	if err != nil {
		return nil, err
	}

	for _, tag := range tags {
		if tag.IsPreRelease() {
			return tag, nil
		}
	}

	return nil, fmt.Errorf("no prerelease tag found")
}

// IsValidTag returns true if the given string is a valid tag
func (g *GitHubHelper) IsValidTag(tagVersion string) (bool, error) {
	if tagVersion[0] == 'v' {
		tagVersion = tagVersion[1:]
	}
	tags, err := g.GetVersionTags()
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
