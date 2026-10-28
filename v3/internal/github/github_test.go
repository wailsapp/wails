package github

import (
	"io"
	"net/http"
	"strings"
	"testing"
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
