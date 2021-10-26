package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

// GetVersionTags gets the list of tags on the Wails repo
// It returns a list of sorted tags in descending order
func GetVersionTags() ([]*SemanticVersion, error) {

	result := []*SemanticVersion{}
	var err error

	resp, err := http.Get("https://api.github.com/repos/wailsapp/wails/tags")
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
		if !strings.HasPrefix(version, "v2") {
			continue
		}
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

// GetLatestPreRelease gets the latest prerelease on GitHub
func GetLatestPreRelease() (result *SemanticVersion, err error) {

	tags, err := GetVersionTags()
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
