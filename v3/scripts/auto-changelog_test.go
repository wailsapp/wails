//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExtractCodeRabbitWalkthrough(t *testing.T) {
	body := `<!-- This is an auto-generated comment: summarize by coderabbit.ai -->
<!-- recent_review_start -->
No actionable comments were generated in the recent review.
<!-- recent_review_end -->
<!-- walkthrough_start -->
## Walkthrough

The website lockfiles update nanoid to a patched release.
<!-- walkthrough_end -->`

	got, ok := extractCodeRabbitWalkthrough(body)
	if !ok {
		t.Fatal("extractCodeRabbitWalkthrough() rejected a complete walkthrough")
	}
	want := "## Walkthrough\n\nThe website lockfiles update nanoid to a patched release."
	if got != want {
		t.Fatalf("extractCodeRabbitWalkthrough() = %q, want %q", got, want)
	}
}

func TestExtractCodeRabbitWalkthroughRejectsSkippedReview(t *testing.T) {
	// This is the status-only comment shape from PR #5985. Treating it as a
	// change summary allowed an unrelated JSON-parsing entry into beta.9.
	body := `<!-- This is an auto-generated comment: summarize by coderabbit.ai -->
<!-- This is an auto-generated comment: skip review by coderabbit.ai -->

> [!IMPORTANT]
> ## Review skipped
>
> Review was skipped due to path filters
>
> * website/package-lock.json is excluded
> * website/pnpm-lock.yaml is excluded`

	if got, ok := extractCodeRabbitWalkthrough(body); ok {
		t.Fatalf("extractCodeRabbitWalkthrough() = %q, want rejection", got)
	}
}

func TestChangelogContextAlwaysIncludesPRTitle(t *testing.T) {
	const title = "fix(security): update website nanoid lockfiles"
	if got, want := changelogContext(title, ""), "PR Title: "+title; got != want {
		t.Fatalf("changelogContext() = %q, want %q", got, want)
	}

	got := changelogContext(title, "The lockfiles update nanoid.")
	want := "PR Title: " + title + "\n\nCodeRabbit Walkthrough:\nThe lockfiles update nanoid."
	if got != want {
		t.Fatalf("changelogContext() = %q, want %q", got, want)
	}
}

func TestFetchCodeRabbitWalkthroughPaginatesComments(t *testing.T) {
	type comment struct {
		User struct {
			Login string `json:"login"`
		} `json:"user"`
		Body string `json:"body"`
	}

	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		page := r.URL.Query().Get("page")
		var comments []comment
		switch page {
		case "1":
			comments = make([]comment, 100)
			for i := range comments {
				comments[i].User.Login = "contributor"
				comments[i].Body = fmt.Sprintf("comment %d", i)
			}
		case "2":
			comments = make([]comment, 1)
			comments[0].User.Login = "coderabbitai[bot]"
			comments[0].Body = "<!-- walkthrough_start -->\nSecond-page walkthrough\n<!-- walkthrough_end -->"
		default:
			t.Fatalf("unexpected comments page %q", page)
		}
		if err := json.NewEncoder(w).Encode(comments); err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()

	oldBaseURL := githubAPIBaseURL
	githubAPIBaseURL = server.URL
	defer func() { githubAPIBaseURL = oldBaseURL }()

	got, err := fetchCodeRabbitWalkthrough("wailsapp/wails", "5993", "token")
	if err != nil {
		t.Fatal(err)
	}
	if got != "Second-page walkthrough" {
		t.Fatalf("fetchCodeRabbitWalkthrough() = %q, want second-page walkthrough", got)
	}
	if requests != 2 {
		t.Fatalf("fetchCodeRabbitWalkthrough() made %d requests, want 2", requests)
	}
}

func TestDocumentationURLForFile(t *testing.T) {
	tests := []struct {
		name string
		file string
		want string
	}{
		{
			name: "regular page",
			file: "docs/src/content/docs/features/windows/options.mdx",
			want: "https://v3.wails.io/features/windows/options",
		},
		{
			name: "index page",
			file: "docs/src/content/docs/guides/mobile/index.mdx",
			want: "https://v3.wails.io/guides/mobile",
		},
		{
			name: "localized page",
			file: "docs/src/content/docs/de/quick-start/installation.mdx",
			want: "https://v3.wails.io/de/quick-start/installation",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := documentationURLFromPath(test.file, "")
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Fatalf("documentationURLFromPath(%q) = %q, want %q", test.file, got, test.want)
			}
		})
	}
}

func TestDocumentationURLForSlug(t *testing.T) {
	got, err := documentationURLFromPath("docs/src/content/docs/blog/legacy-name.md", "blog/the-road-to-wails-v3")
	if err != nil {
		t.Fatal(err)
	}
	if got != "https://v3.wails.io/blog/the-road-to-wails-v3" {
		t.Fatalf("got %q", got)
	}
}

func TestDocumentationURLForMissingFile(t *testing.T) {
	got, err := documentationURLForFile("docs/src/content/docs/removed-by-pr.mdx")
	if err != nil {
		t.Fatal(err)
	}
	if got != "" {
		t.Fatalf("documentationURLForFile() = %q, want empty URL for a missing file", got)
	}
}
