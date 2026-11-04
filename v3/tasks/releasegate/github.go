package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	githubDefaultAPI = "https://api.github.com"
	githubAPIVersion = "2022-11-28"
	// maxPages stops a misbehaving server from looping forever through
	// `Link: rel="next"`. 1000 pages of 100 items is far beyond anything the
	// gate reads.
	maxPages = 1000
)

// client is a minimal GitHub REST client. Every non-2xx response is an
// error, so callers can fail closed on anything they could not read.
type client struct {
	base  string
	token string
	http  *http.Client
}

func newClient(base, token string) *client {
	if base == "" {
		base = githubDefaultAPI
	}
	return &client{
		base:  strings.TrimRight(base, "/"),
		token: token,
		http:  &http.Client{Timeout: 60 * time.Second},
	}
}

type apiError struct {
	method string
	url    string
	status int
	body   string
}

func (e *apiError) Error() string {
	return fmt.Sprintf("%s %s returned %d: %s", e.method, e.url, e.status, e.body)
}

// do sends one request. The body, if any, is JSON.
func (c *client) do(method, rawURL string, body any) ([]byte, http.Header, error) {
	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, nil, err
		}
		reader = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, rawURL, reader)
	if err != nil {
		return nil, nil, err
	}
	// Without a token the request is sent unauthenticated; the private
	// endpoints then answer 401/403/404 and the gate fails closed.
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", githubAPIVersion)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, nil, &apiError{method: method, url: rawURL, status: resp.StatusCode, body: strings.TrimSpace(string(data))}
	}
	return data, resp.Header, nil
}

// get reads a single (non-paginated) JSON object into v.
func (c *client) get(path string, v any) error {
	data, _, err := c.do(http.MethodGet, c.base+"/"+path, nil)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// send issues a POST/PATCH with a JSON body and discards the response.
func (c *client) send(method, path string, body any) error {
	_, _, err := c.do(method, c.base+"/"+path, body)
	return err
}

// list reads every page of a JSON array endpoint, following the `Link`
// header's rel="next" URL (used as given: it is absolute), and returns the
// concatenation, like `gh api --paginate | jq -s 'add // []'`. Any page that
// fails or is not an array makes the whole list an error.
func list[T any](c *client, path string) ([]T, error) {
	baseURL, err := url.Parse(c.base)
	if err != nil {
		return nil, err
	}
	next := c.base + "/" + path
	var all []T
	for page := 0; next != ""; page++ {
		if page == maxPages {
			return nil, fmt.Errorf("GET %s: more than %d pages", path, maxPages)
		}
		data, header, err := c.do(http.MethodGet, next, nil)
		if err != nil {
			return nil, err
		}
		var items []T
		if err := json.Unmarshal(data, &items); err != nil {
			return nil, fmt.Errorf("GET %s: %w", next, err)
		}
		all = append(all, items...)
		next = nextLink(header.Get("Link"))
		if next != "" {
			// Never send the token to a host other than the API's.
			u, err := url.Parse(next)
			if err != nil || u.Host != baseURL.Host {
				return nil, fmt.Errorf("GET %s: refusing pagination link %q", path, next)
			}
		}
	}
	return all, nil
}

// nextLink extracts the rel="next" URL from a Link header, or "".
//
//	Link: <https://api.github.com/...&page=2>; rel="next", <...>; rel="last"
func nextLink(header string) string {
	for _, part := range strings.Split(header, ",") {
		segs := strings.Split(part, ";")
		if len(segs) < 2 {
			continue
		}
		target := strings.TrimSpace(segs[0])
		if !strings.HasPrefix(target, "<") || !strings.HasSuffix(target, ">") {
			continue
		}
		for _, param := range segs[1:] {
			param = strings.TrimSpace(param)
			if param == `rel="next"` || param == "rel=next" {
				return target[1 : len(target)-1]
			}
		}
	}
	return ""
}
