package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"

	"github.com/Masterminds/semver"
)

// GitHubHelper is a utility class for interacting with GitHub
type GitHubHelper struct {
	validPrereleaseRegex *regexp.Regexp
	validReleaseRegex    *regexp.Regexp
}

// NewGitHubHelper returns a new GitHub Helper
func NewGitHubHelper() *GitHubHelper {
	const SemVerRegex string = `v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?`
	const SemPreVerRegex string = `v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*)`

	return &GitHubHelper{
		validPrereleaseRegex: regexp.MustCompile(SemPreVerRegex),
		validReleaseRegex:    regexp.MustCompile(SemVerRegex),
	}
}

// GetVersionTags gets the list of tags on the Wails repo
// It retuns a list of sorted tags in descending order
func (g *GitHubHelper) GetVersionTags() ([]*semver.Version, error) {

	result := []*semver.Version{}
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
		semver, err := semver.NewVersion(version)
		if err != nil {
			return result, err
		}
		result = append(result, semver)
	}

	// Reverse Sort
	sort.Sort(sort.Reverse(semver.Collection(result)))

	return result, err
}

func (g *GitHubHelper) isRelease(tag *semver.Version) bool {
	return g.validReleaseRegex.MatchString(tag.String())
}

func (g *GitHubHelper) isPreRelease(tag *semver.Version) bool {
	return g.validPrereleaseRegex.MatchString(tag.String())
}

// GetLatestStableRelease gets the latest stable release on GitHub
func (g *GitHubHelper) GetLatestStableRelease() (result string, err error) {

	tags, err := g.GetVersionTags()
	if err != nil {
		return "", err
	}

	for _, tag := range tags {
		if g.isRelease(tag) {
			return "v" + tag.String(), nil
		}
	}

	return "", fmt.Errorf("no release tag found")
}
