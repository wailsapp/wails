package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Contributor is a single contributor with their credit tally.
type Contributor struct {
	Login     string
	AvatarURL string
	URL       string
	// Commits is the default-branch commit count from the contributors API.
	Commits int
	// PRs is the number of merged pull requests authored (only populated
	// with -metric prs).
	PRs int
	// Mentions counts @login credits found in the changelogs. Squash-merged
	// or hand-applied patches are often credited only there.
	Mentions int
}

// creditMetric selects how contributors are ranked and sized: "commits"
// (default-branch commit count) or "prs" (merged pull requests authored).
// Commit counts reward granular unsquashed histories; PR counts treat one
// merged PR as one unit of work regardless of merge style.
var creditMetric = "commits"

// Credit is the number a contributor is ranked and sized by. The changelog
// mention count mostly describes the same work as the primary metric, so
// take the larger of the two rather than double counting.
func (c Contributor) Credit() int {
	primary := c.Commits
	if creditMetric == "prs" {
		primary = c.PRs
		// Contributors with commits but no recorded PRs (early direct
		// pushes) stay visible in the tail bands rather than vanishing.
		if primary == 0 && c.Mentions == 0 && c.Commits > 0 {
			primary = min(c.Commits, 7)
		}
	}
	if c.Mentions > primary {
		return c.Mentions
	}
	return primary
}

type restContributor struct {
	Login         string `json:"login"`
	AvatarURL     string `json:"avatar_url"`
	HTMLURL       string `json:"html_url"`
	Contributions int    `json:"contributions"`
	Type          string `json:"type"`
}

// FetchContributors returns the repo's contributors (up to the API's cap of
// 500), bots excluded, augmented with changelog @mention credits so that
// people whose work arrived via squashed or hand-applied patches still
// appear. Sorted by credit (descending) then login for deterministic output.
func FetchContributors(client *http.Client, token, repo string, changelogs []string) ([]Contributor, error) {
	byLogin := map[string]*Contributor{}
	var order []string

	for page := 1; ; page++ {
		url := fmt.Sprintf("https://api.github.com/repos/%s/contributors?per_page=100&page=%d", repo, page)
		var batch []restContributor
		if err := getJSON(client, token, url, &batch); err != nil {
			return nil, fmt.Errorf("fetching contributors: %w", err)
		}
		for _, c := range batch {
			if c.Login == "" || isBot(c.Login, c.Type) {
				continue
			}
			key := strings.ToLower(c.Login)
			if _, ok := byLogin[key]; ok {
				continue
			}
			byLogin[key] = &Contributor{
				Login:     c.Login,
				AvatarURL: c.AvatarURL,
				URL:       c.HTMLURL,
				Commits:   c.Contributions,
			}
			order = append(order, key)
		}
		if len(batch) < 100 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	if creditMetric == "prs" {
		prCounts, err := fetchMergedPRCounts(client, token, repo)
		if err != nil {
			return nil, fmt.Errorf("fetching merged PRs: %w", err)
		}
		for key, n := range prCounts {
			if c, ok := byLogin[key]; ok {
				c.PRs = n
				continue
			}
			// Authored merged PRs but absent from the commit history (e.g.
			// unlinked commit email). Same validation as changelog mentions.
			u, err := lookupUser(client, token, key)
			if err != nil || u.Type != "User" || isBot(u.Login, u.Type) {
				continue
			}
			byLogin[key] = &Contributor{
				Login:     u.Login,
				AvatarURL: u.AvatarURL,
				URL:       u.HTMLURL,
				PRs:       n,
			}
			order = append(order, key)
		}
	}

	mentions := changelogMentions(changelogs)
	for key, n := range mentions {
		if c, ok := byLogin[key]; ok {
			c.Mentions = n
			continue
		}
		// Credited in a changelog but absent from the commit history (e.g.
		// their PR was squashed under someone else's authorship). Look the
		// login up so typos and non-user tokens are dropped.
		u, err := lookupUser(client, token, key)
		if err != nil {
			fmt.Printf("  skipping changelog mention @%s: %v\n", key, err)
			continue
		}
		// Organisations and bots cannot be contributors; requiring a real
		// user account also drops tokens that merely happen to start with @.
		if u.Type != "User" || isBot(u.Login, u.Type) {
			fmt.Printf("  skipping changelog mention @%s: %s account\n", key, u.Type)
			continue
		}
		fmt.Printf("  changelog-only contributor: @%s (%d mentions)\n", u.Login, n)
		byLogin[key] = &Contributor{
			Login:     u.Login,
			AvatarURL: u.AvatarURL,
			URL:       u.HTMLURL,
			Mentions:  n,
		}
		order = append(order, key)
	}

	out := make([]Contributor, 0, len(order))
	for _, key := range order {
		out = append(out, *byLogin[key])
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Credit() != out[j].Credit() {
			return out[i].Credit() > out[j].Credit()
		}
		return strings.ToLower(out[i].Login) < strings.ToLower(out[j].Login)
	})
	return out, nil
}

// mentionRe matches @login tokens. GitHub logins are alphanumeric with
// single hyphens, up to 39 characters.
var mentionRe = regexp.MustCompile(`(^|[^\w./@])@([A-Za-z0-9][A-Za-z0-9-]{0,38})`)

// profileLinkRe matches markdown links to GitHub profiles, e.g.
// [@someone](https://github.com/someone).
var profileLinkRe = regexp.MustCompile(`\[([^\]]*)\]\(https://github\.com/([A-Za-z0-9][A-Za-z0-9-]{0,38})/?\)`)

// changelogMentions counts @login credits across the given files, keyed by
// lowercased login. Files that cannot be read are skipped with a warning so
// the tool still works from checkouts without the docs trees.
func changelogMentions(paths []string) map[string]int {
	counts := map[string]int{}
	for _, path := range paths {
		if path == "" {
			continue
		}
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("  warning: changelog %s: %v\n", path, err)
			continue
		}
		// Credits written as markdown profile links first. The link text is
		// sometimes a typo or display handle, so trust the URL, then strip
		// the whole link so the text is not counted again.
		data = profileLinkRe.ReplaceAllFunc(data, func(link []byte) []byte {
			m := profileLinkRe.FindSubmatch(link)
			if bytes.Contains(m[1], []byte("@")) {
				counts[strings.ToLower(string(m[2]))]++
			}
			return []byte(" ")
		})
		for _, m := range mentionRe.FindAllSubmatchIndex(data, -1) {
			login := strings.TrimRight(string(data[m[4]:m[5]]), "-")
			if login == "" {
				continue
			}
			// Package scopes like @wailsio/runtime are not people, and an
			// underscore means the token is a display handle rather than a
			// GitHub login.
			if m[5] < len(data) && (data[m[5]] == '/' || data[m[5]] == '_') {
				continue
			}
			counts[strings.ToLower(login)]++
		}
	}
	return counts
}

// mergedPRQuery pages through every merged PR's author. Unlike the search
// API this has no 1000-result cap.
const mergedPRQuery = `query($owner: String!, $name: String!, $after: String) {
  repository(owner: $owner, name: $name) {
    pullRequests(states: MERGED, first: 100, after: $after) {
      pageInfo { hasNextPage endCursor }
      nodes { author { login } }
    }
  }
}`

type mergedPRResponse struct {
	Data struct {
		Repository struct {
			PullRequests struct {
				PageInfo struct {
					HasNextPage bool   `json:"hasNextPage"`
					EndCursor   string `json:"endCursor"`
				} `json:"pageInfo"`
				Nodes []struct {
					Author *struct {
						Login string `json:"login"`
					} `json:"author"`
				} `json:"nodes"`
			} `json:"pullRequests"`
		} `json:"repository"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// fetchMergedPRCounts returns merged-PR counts keyed by lowercased author
// login, bots excluded.
func fetchMergedPRCounts(client *http.Client, token, repo string) (map[string]int, error) {
	owner, name, ok := strings.Cut(repo, "/")
	if !ok {
		return nil, fmt.Errorf("repo %q is not owner/name", repo)
	}
	counts := map[string]int{}
	after := ""
	for {
		variables := map[string]any{"owner": owner, "name": name}
		if after != "" {
			variables["after"] = after
		}
		body, err := json.Marshal(map[string]any{"query": mergedPRQuery, "variables": variables})
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequest(http.MethodPost, "https://api.github.com/graphql", bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", "wails-sponsorkit")

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		data, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("github graphql: %s: %s", resp.Status, truncate(string(data), 300))
		}
		var parsed mergedPRResponse
		if err := json.Unmarshal(data, &parsed); err != nil {
			return nil, fmt.Errorf("parsing graphql response: %w", err)
		}
		if len(parsed.Errors) > 0 {
			return nil, fmt.Errorf("github graphql: %s", parsed.Errors[0].Message)
		}

		conn := parsed.Data.Repository.PullRequests
		for _, node := range conn.Nodes {
			if node.Author == nil || node.Author.Login == "" || isBot(node.Author.Login, "") {
				continue
			}
			counts[strings.ToLower(node.Author.Login)]++
		}
		if !conn.PageInfo.HasNextPage {
			break
		}
		after = conn.PageInfo.EndCursor
		time.Sleep(100 * time.Millisecond)
	}
	return counts, nil
}

type restUser struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
	Type      string `json:"type"`
}

func lookupUser(client *http.Client, token, login string) (restUser, error) {
	var u restUser
	err := getJSON(client, token, "https://api.github.com/users/"+login, &u)
	return u, err
}

func isBot(login, typ string) bool {
	l := strings.ToLower(login)
	return typ == "Bot" ||
		strings.HasSuffix(l, "[bot]") ||
		strings.HasSuffix(l, "-bot") ||
		l == "dependabot" || l == "allcontributors" || l == "github-actions"
}

func getJSON(client *http.Client, token, url string, into any) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	if token != "" {
		req.Header.Set("Authorization", "bearer "+token)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "wails-sponsorkit")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET %s: %s: %s", url, resp.Status, truncate(string(data), 200))
	}
	return json.Unmarshal(data, into)
}
