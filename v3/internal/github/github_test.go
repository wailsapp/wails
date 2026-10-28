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
		if r.URL.String() != "https://api.github.com/repos/wailsapp/wails/tags" {
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
		{"notes", `{"body":"release notes"}`, func() { GetReleaseNotes("v3.0.0", true) }},
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
