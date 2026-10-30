package github

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func tagPageClient(tb testing.TB, total int) (*http.Client, *int) {
	tb.Helper()
	names := make([]string, total)
	for i := range names {
		names[i] = fmt.Sprintf(`{"name":"v3.0.%d"}`, i)
	}
	pages := make(map[int][]string)
	for _, size := range []int{30, 100} {
		for start := 0; start < total; start += size {
			end := min(start+size, total)
			pages[size] = append(pages[size], "["+strings.Join(names[start:end], ",")+"]")
		}
	}
	calls := 0
	client := &http.Client{Transport: tagTransport(func(r *http.Request) (*http.Response, error) {
		calls++
		query := r.URL.Query()
		size := 30
		if value := query.Get("per_page"); value != "" {
			size, _ = strconv.Atoi(value)
		}
		page := 1
		if value := query.Get("page"); value != "" {
			page, _ = strconv.Atoi(value)
		}
		if page < 1 || page > len(pages[size]) {
			tb.Fatalf("unexpected tag request: %s", r.URL)
		}
		header := make(http.Header)
		if page < len(pages[size]) {
			header.Set("Link", fmt.Sprintf(`<https://api.github.com/repos/wailsapp/wails/tags?per_page=%d&page=%d>; rel="next"`, size, page+1))
		}
		return &http.Response{StatusCode: http.StatusOK, Header: header, Body: io.NopCloser(strings.NewReader(pages[size][page-1]))}, nil
	})}
	return client, &calls
}

func BenchmarkVersionTagsPagination(b *testing.B) {
	previous := http.DefaultClient
	b.Cleanup(func() { http.DefaultClient = previous })
	client, calls := tagPageClient(b, 245)
	http.DefaultClient = client
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		tags, err := GetVersionTags()
		if err != nil || len(tags) != 245 {
			b.Fatalf("tags=%d, error=%v; want 245, nil", len(tags), err)
		}
	}
	b.ReportMetric(float64(*calls)/float64(b.N), "requests/op")
}
