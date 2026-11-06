package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	githubAPIURL        = "https://api.github.com"
	githubStarMediaType = "application/vnd.github.star+json"
	githubPageSize      = 100
)

// Star is the portion of a GitHub stargazer record needed by the chart.
type Star struct {
	StarredAt time.Time `json:"starred_at"`
}

// repositorySummary is the repository API subset used for reconciliation.
type repositorySummary struct {
	StargazersCount int `json:"stargazers_count"`
}

// GitHubClient reads the timestamped stargazers for a repository.
type GitHubClient struct {
	HTTPClient *http.Client
	BaseURL    string
}

// FetchStars returns all currently visible stargazers, oldest first or until
// GitHub returns a short page. Access requires a token belonging to a repo
// admin or collaborator under GitHub's current stargazer API policy.
func (c *GitHubClient) FetchStars(ctx context.Context, token, repo string) ([]Star, error) {
	if token == "" {
		return nil, fmt.Errorf("GitHub token is required to read stargazer timestamps")
	}
	parts := strings.Split(repo, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("repository must be in owner/name form, got %q", repo)
	}

	client := c.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	baseURL := strings.TrimRight(c.BaseURL, "/")
	if baseURL == "" {
		baseURL = githubAPIURL
	}

	fetchPage := func(page int) ([]Star, string, error) {
		endpoint := fmt.Sprintf("%s/repos/%s/%s/stargazers?per_page=%d&page=%d", baseURL,
			url.PathEscape(parts[0]), url.PathEscape(parts[1]), githubPageSize, page)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			return nil, "", fmt.Errorf("create GitHub request: %w", err)
		}
		req.Header.Set("Accept", githubStarMediaType)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("User-Agent", "wails-star-history")

		resp, err := client.Do(req)
		if err != nil {
			return nil, "", fmt.Errorf("fetch GitHub stargazers page %d: %w", page, err)
		}
		body, readErr := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if readErr != nil {
			return nil, "", fmt.Errorf("read GitHub stargazers page %d: %w", page, readErr)
		}
		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
			message := strings.TrimSpace(string(body))
			if len(message) > 300 {
				message = message[:300] + "…"
			}
			return nil, "", fmt.Errorf("GitHub stargazers page %d returned %s: %s", page, resp.Status, message)
		}

		var pageStars []Star
		if err := json.Unmarshal(body, &pageStars); err != nil {
			return nil, "", fmt.Errorf("decode GitHub stargazers page %d: %w", page, err)
		}
		return pageStars, resp.Header.Get("Link"), nil
	}

	firstPage, linkHeader, err := fetchPage(1)
	if err != nil {
		return nil, err
	}
	stars := append([]Star(nil), firstPage...)
	if len(firstPage) < githubPageSize {
		return stars, nil
	}

	lastPage := lastPageFromLink(linkHeader)
	if lastPage <= 1 {
		for page := 2; ; page++ {
			pageStars, _, err := fetchPage(page)
			if err != nil {
				return nil, err
			}
			stars = append(stars, pageStars...)
			if len(pageStars) < githubPageSize {
				return stars, nil
			}
		}
	}

	// GitHub's Link header tells us the final page. Fetching those pages in a
	// small bounded pool keeps a large repository's weekly refresh practical
	// without creating an unbounded burst of API requests.
	pageJobs := make(chan int)
	pageResults := make(chan []Star, lastPage-1)
	pageErrors := make(chan error, 1)
	workerCount := lastPage - 1
	if workerCount > 8 {
		workerCount = 8
	}
	var workers sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for page := range pageJobs {
				pageStars, _, err := fetchPage(page)
				if err != nil {
					select {
					case pageErrors <- err:
					default:
					}
					continue
				}
				pageResults <- pageStars
			}
		}()
	}
	go func() {
		for page := 2; page <= lastPage; page++ {
			pageJobs <- page
		}
		close(pageJobs)
		workers.Wait()
		close(pageResults)
	}()
	for pageStars := range pageResults {
		stars = append(stars, pageStars...)
	}
	select {
	case err := <-pageErrors:
		return nil, err
	default:
		return stars, nil
	}
}

// FetchStarCount returns GitHub's current repository total. The stargazer list
// is timestamped history, while this endpoint is the authoritative live count
// used to reconcile a concurrent star/unstar or a restricted list response.
func (c *GitHubClient) FetchStarCount(ctx context.Context, token, repo string) (int, error) {
	if token == "" {
		return 0, fmt.Errorf("GitHub token is required to read the repository star count")
	}
	parts := strings.Split(repo, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return 0, fmt.Errorf("repository must be in owner/name form, got %q", repo)
	}
	client := c.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	baseURL := strings.TrimRight(c.BaseURL, "/")
	if baseURL == "" {
		baseURL = githubAPIURL
	}
	endpoint := fmt.Sprintf("%s/repos/%s/%s", baseURL, url.PathEscape(parts[0]), url.PathEscape(parts[1]))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return 0, fmt.Errorf("create GitHub repository request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "wails-star-history")
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("fetch GitHub repository: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return 0, fmt.Errorf("GitHub repository returned %s", resp.Status)
	}
	var summary repositorySummary
	if err := json.NewDecoder(resp.Body).Decode(&summary); err != nil {
		return 0, fmt.Errorf("decode GitHub repository: %w", err)
	}
	return summary.StargazersCount, nil
}

// lastPageFromLink extracts the final page number from GitHub's Link header.
func lastPageFromLink(linkHeader string) int {
	for _, link := range strings.Split(linkHeader, ",") {
		if !strings.Contains(link, `rel="last"`) {
			continue
		}
		start := strings.Index(link, "<")
		end := strings.Index(link, ">")
		if start < 0 || end <= start {
			continue
		}
		lastURL, err := url.Parse(link[start+1 : end])
		if err == nil {
			page, err := strconv.Atoi(lastURL.Query().Get("page"))
			if err == nil {
				return page
			}
		}
	}
	return 0
}
