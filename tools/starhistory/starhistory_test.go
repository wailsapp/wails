package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestBuildWeeklyHistoryAccumulatesStarsByWeek verifies weekly cumulative grouping.
func TestBuildWeeklyHistoryAccumulatesStarsByWeek(t *testing.T) {
	stars := []Star{
		{StarredAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)}, // Monday
		{StarredAt: time.Date(2024, 1, 7, 12, 0, 0, 0, time.UTC)}, // Same week
		{StarredAt: time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC)}, // Next week
	}

	history := BuildWeeklyHistory(stars)

	if len(history) != 2 {
		t.Fatalf("expected 2 weekly points, got %d", len(history))
	}
	if got, want := history[0].Date.Format("2006-01-02"), "2024-01-01"; got != want {
		t.Fatalf("first point date = %s, want %s", got, want)
	}
	if history[0].Count != 2 || history[1].Count != 3 {
		t.Fatalf("counts = %d, %d; want 2, 3", history[0].Count, history[1].Count)
	}
}

// TestReconcileCurrentCountAddsAuthoritativeCurrentPoint verifies live-count reconciliation.
func TestReconcileCurrentCountAddsAuthoritativeCurrentPoint(t *testing.T) {
	now := time.Date(2024, 1, 9, 0, 0, 0, 0, time.UTC)
	history := []WeeklyPoint{{Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), Count: 2}}
	history = ReconcileCurrentCount(history, 3, now)
	if len(history) != 2 || history[1].Count != 3 || !history[1].Date.Equal(now) {
		t.Fatalf("reconciled history = %#v, want current point at %s with count 3", history, now)
	}
}

// TestFetchStarsPaginatesAndAuthenticates verifies API pagination and request headers.
func TestFetchStarsPaginatesAndAuthenticates(t *testing.T) {
	var requests int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if got, want := r.Header.Get("Authorization"), "Bearer secret"; got != want {
			t.Errorf("authorization = %q, want %q", got, want)
		}
		if got, want := r.Header.Get("Accept"), githubStarMediaType; got != want {
			t.Errorf("accept = %q, want %q", got, want)
		}

		var payload []map[string]string
		if r.URL.Query().Get("page") == "2" {
			payload = []map[string]string{{"starred_at": "2024-01-08T12:00:00Z"}}
		} else {
			payload = make([]map[string]string, 100)
			for i := range payload {
				payload[i] = map[string]string{"starred_at": "2024-01-01T12:00:00Z"}
			}
		}
		_ = json.NewEncoder(w).Encode(payload)
	}))
	defer server.Close()

	client := &GitHubClient{HTTPClient: server.Client(), BaseURL: server.URL}
	stars, err := client.FetchStars(context.Background(), "secret", "wailsapp/wails")
	if err != nil {
		t.Fatalf("FetchStars() error = %v", err)
	}
	if len(stars) != 101 {
		t.Fatalf("got %d stars, want 101", len(stars))
	}
	if requests != 2 {
		t.Fatalf("made %d requests, want 2", requests)
	}
}

// TestLastPageFromLink verifies parsing of GitHub's pagination metadata.
func TestLastPageFromLink(t *testing.T) {
	link := `<https://api.github.com/repositories/161951219/stargazers?per_page=100&page=2>; rel="next", <https://api.github.com/repositories/161951219/stargazers?per_page=100&page=359>; rel="last"`
	if got, want := lastPageFromLink(link), 359; got != want {
		t.Fatalf("last page = %d, want %d", got, want)
	}
}

// TestRenderSVGIncludesFadedBackgroundAndRedChart verifies the chart's embedded assets and styling.
func TestRenderSVGIncludesFadedBackgroundAndRedChart(t *testing.T) {
	history := []WeeklyPoint{
		{Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), Count: 2},
		{Date: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC), Count: 8},
	}

	svg := RenderSVG(RenderOptions{
		Repo:           "wailsapp/wails",
		History:        history,
		BackgroundData: []byte("background-bytes"),
		LogoData:       []byte("<svg>logo</svg>"),
		RefreshedAt:    time.Date(2024, 1, 9, 0, 0, 0, 0, time.UTC),
	})

	for _, want := range []string{
		`data:image/webp;base64,YmFja2dyb3VuZC1ieXRlcw==`,
		`data:image/svg+xml;base64,PHN2Zz5sb2dvPC9zdmc+`,
		`fill="url(#chartFill)"`,
		`stroke="#ff5364"`,
		"Wails",
		"8 stars",
	} {
		if !strings.Contains(svg, want) {
			t.Errorf("SVG does not contain %q", want)
		}
	}
}

// TestFormatAxisCountUsesCompactThousands verifies abbreviated y-axis labels.
func TestFormatAxisCountUsesCompactThousands(t *testing.T) {
	for value, want := range map[int]string{
		0:     "0",
		999:   "999",
		1000:  "1k",
		12500: "12.5k",
		50000: "50k",
	} {
		if got := formatAxisCount(value); got != want {
			t.Errorf("formatAxisCount(%d) = %q, want %q", value, got, want)
		}
	}
}
