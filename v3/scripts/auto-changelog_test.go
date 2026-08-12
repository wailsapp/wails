//go:build ignore

package main

import "testing"

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
