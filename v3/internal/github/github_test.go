package github

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"testing/iotest"
)

type tagTransport func(*http.Request) (*http.Response, error)

func (f tagTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestVersionTagSelection(t *testing.T) {
	previousClient := http.DefaultClient
	t.Cleanup(func() { http.DefaultClient = previousClient })
	http.DefaultClient = &http.Client{Transport: tagTransport(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() != "https://api.github.com/repos/wailsapp/wails/tags?per_page=100&page=1" {
			t.Errorf("unexpected request: %s", r.URL)
		}
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`[
			{"name":"v2.99.0"}, {"name":"v3.0.0-beta.2"}, {"name":"v3.0.0"},
			{"name":"v3.1.0-beta.1"}, {"name":"v30.0.0"}, {"name":"v3.0.1"}
		]`)), Header: make(http.Header)}, nil
	})}

	tags, err := GetVersionTags()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"3.1.0-beta.1", "3.0.1", "3.0.0", "3.0.0-beta.2"}
	if len(tags) != len(want) {
		t.Fatalf("got %d tags, want %d: %v", len(tags), len(want), tags)
	}
	for i, tag := range tags {
		if tag.String() != want[i] {
			t.Errorf("tag[%d]=%s, want %s", i, tag, want[i])
		}
	}
	stable, err := GetLatestStableRelease()
	if err != nil || stable == nil || stable.String() != "3.0.1" {
		t.Errorf("latest stable=%v, error=%v", stable, err)
	}
	pre, err := GetLatestPreRelease()
	if err != nil || pre == nil || pre.String() != "3.1.0-beta.1" {
		t.Errorf("latest prerelease=%v, error=%v", pre, err)
	}
	valid, err := IsValidTag("v3.0.1")
	if err != nil || !valid {
		t.Errorf("v3 tag rejected: valid=%t, error=%v", valid, err)
	}
	valid, err = IsValidTag("v2.99.0")
	if err != nil || valid {
		t.Errorf("v2 tag accepted: valid=%t, error=%v", valid, err)
	}
}

func TestVersionTagsPagination(t *testing.T) {
	previous := http.DefaultClient
	t.Cleanup(func() { http.DefaultClient = previous })
	var bodies []*trackedResponseBody
	http.DefaultClient = &http.Client{Transport: tagTransport(func(r *http.Request) (*http.Response, error) {
		if r.URL.Query().Get("per_page") != "100" {
			t.Errorf("page size missing from request: %s", r.URL)
		}
		header := make(http.Header)
		payload := `[{"name":"v3.1.0-beta.2"}]`
		if len(bodies) == 0 {
			header.Set("Link", `<https://api.github.com/repos/wailsapp/wails/tags?page=2>; rel="next"`)
		} else {
			if r.URL.Query().Get("page") != "2" {
				t.Errorf("second page URL=%s", r.URL)
			}
			payload = `[{"name":"v3.0.1"},{"name":"v3.1.0-beta.3"}]`
			header.Set("Link", `<https://api.github.com/repos/wailsapp/wails/tags?page=1>; rel="prev"`)
		}
		body := &trackedResponseBody{Reader: strings.NewReader(payload)}
		bodies = append(bodies, body)
		return &http.Response{StatusCode: http.StatusOK, Body: body, Header: header}, nil
	})}
	tags, err := GetVersionTags()
	if err != nil || len(tags) != 3 || tags[0].String() != "3.1.0-beta.3" || tags[2].String() != "3.0.1" {
		t.Fatalf("paginated tags=%v, error=%v", tags, err)
	}
	if len(bodies) != 2 {
		t.Fatalf("requested %d pages, want 2", len(bodies))
	}
	for _, body := range bodies {
		if body.closeCalls != 1 {
			t.Errorf("page body closed %d times, want 1", body.closeCalls)
		}
	}
}

func TestVersionTagsRejectHTTPFailures(t *testing.T) {
	cases := []struct {
		status int
		reader io.Reader
	}{
		{http.StatusForbidden, strings.NewReader(`{"message":"API rate limit exceeded"}`)},
		{http.StatusTooManyRequests, strings.NewReader(`{"message":"Too many requests"}`)},
		{http.StatusInternalServerError, strings.NewReader("<html>Server error</html>")},
		{http.StatusServiceUnavailable, iotest.ErrReader(errors.New("body read failed"))},
		{http.StatusBadGateway, strings.NewReader(`[{"name":"v3.0.0"}]`)},
	}
	for _, tc := range cases {
		t.Run(http.StatusText(tc.status), func(t *testing.T) {
			previousClient := http.DefaultClient
			t.Cleanup(func() { http.DefaultClient = previousClient })
			body := &trackedResponseBody{Reader: tc.reader}
			http.DefaultClient = &http.Client{Transport: tagTransport(func(*http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: tc.status, Body: body, Header: make(http.Header)}, nil
			})}
			tags, err := GetVersionTags()
			want := fmt.Sprintf("HTTP %d (%s)", tc.status, http.StatusText(tc.status))
			if err == nil || !strings.Contains(err.Error(), want) {
				t.Errorf("error=%v, want HTTP status %q", err, want)
			}
			if len(tags) != 0 {
				t.Errorf("failed HTTP response returned tags: %v", tags)
			}
			if body.closeCalls != 1 {
				t.Errorf("response body closed %d times, want 1", body.closeCalls)
			}
		})
	}
}

type trackedResponseBody struct {
	io.Reader
	closeCalls int
}

func (b *trackedResponseBody) Close() error {
	b.closeCalls++
	return nil
}

func TestGitHubResponsesClose(t *testing.T) {
	apis := []struct {
		name, response string
		run            func()
	}{
		{"tags", `[{"name":"v3.0.0"}]`, func() { _, _ = GetVersionTags() }},
		{"notes", `{"body":"release notes"}`, func() { _, _ = GetReleaseNotes("v3.0.0") }},
	}
	for _, api := range apis {
		cases := []struct {
			name   string
			reader io.Reader
		}{
			{"success", strings.NewReader(api.response)},
			{"invalid_json", strings.NewReader("{")},
			{"read_error", iotest.ErrReader(errors.New("read failed"))},
		}
		for _, tc := range cases {
			t.Run(api.name+"/"+tc.name, func(t *testing.T) {
				previousClient := http.DefaultClient
				t.Cleanup(func() { http.DefaultClient = previousClient })
				body := &trackedResponseBody{Reader: tc.reader}
				http.DefaultClient = &http.Client{Transport: tagTransport(func(*http.Request) (*http.Response, error) {
					return &http.Response{StatusCode: http.StatusOK, Body: body, Header: make(http.Header)}, nil
				})}
				api.run()
				if body.closeCalls != 1 {
					t.Errorf("response body closed %d times, want 1", body.closeCalls)
				}
			})
		}
	}
}

func TestIsValidTagEmptyInput(t *testing.T) {
	for _, input := range []string{"", "v"} {
		t.Run("input="+input, func(t *testing.T) {
			previous := http.DefaultClient
			t.Cleanup(func() { http.DefaultClient = previous })
			requests := 0
			http.DefaultClient = &http.Client{Transport: tagTransport(func(*http.Request) (*http.Response, error) {
				requests++
				return nil, errors.New("unexpected tag lookup")
			})}
			valid, err := IsValidTag(input)
			if valid || err != nil {
				t.Errorf("IsValidTag(%q)=(%t,%v), want false,nil", input, valid, err)
			}
			if requests != 0 {
				t.Errorf("made %d requests for empty version", requests)
			}
		})
	}
}

func TestVersionTagsUseLargerPages(t *testing.T) {
	previous := http.DefaultClient
	t.Cleanup(func() { http.DefaultClient = previous })
	client, calls := tagPageClient(t, 245)
	http.DefaultClient = client
	tags, err := GetVersionTags()
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 245 || *calls != 3 {
		t.Fatalf("tags=%d, requests=%d; want 245, 3", len(tags), *calls)
	}
	for i, tag := range tags {
		want := fmt.Sprintf("3.0.%d", 244-i)
		if tag.String() != want {
			t.Errorf("tag[%d]=%s, want %s", i, tag, want)
		}
	}
}

func TestLatestStableRejectsEmptyVersionSuffix(t *testing.T) {
	previousClient := http.DefaultClient
	t.Cleanup(func() { http.DefaultClient = previousClient })
	for _, tag := range []string{"v3.99.0-", "v3.99.0+", "v3.99.0-+build", "v3.99.0-beta+"} {
		t.Run(tag, func(t *testing.T) {
			http.DefaultClient = &http.Client{Transport: tagTransport(func(*http.Request) (*http.Response, error) {
				body := fmt.Sprintf(`[{"name":%q},{"name":"v3.0.0"}]`, tag)
				return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
			})}
			version, err := GetLatestStableRelease()
			if err != nil || version == nil || version.String() != "3.0.0" {
				t.Fatalf("malformed tag %q returned version=%v, error=%v", tag, version, err)
			}
		})
	}
}
