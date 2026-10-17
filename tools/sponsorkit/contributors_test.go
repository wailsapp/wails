package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeChangelog(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "changelog.mdx")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestChangelogMentions(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    map[string]int
	}{
		{
			name:    "bare mentions, case-folded and counted",
			content: "- Fixed X by @Alice\n- Fixed Y by @alice and @bob",
			want:    map[string]int{"alice": 2, "bob": 1},
		},
		{
			name:    "profile link trusts the URL, not the (typo'd) link text",
			content: "- Fixed X [@hkere](https://github.com/hkhere)",
			want:    map[string]int{"hkhere": 1},
		},
		{
			name:    "profile link without @ in the text is not a credit",
			content: "- See [the repo](https://github.com/wailsapp) for details",
			want:    map[string]int{},
		},
		{
			name:    "npm scopes and display handles are not logins",
			content: "- Bump @wailsio/runtime\n- Fixed by @ronaldinho_x86",
			want:    map[string]int{},
		},
		{
			name:    "emails and mid-word @ are not mentions",
			content: "- Contact lea.anthony@gmail.com or user@@host",
			want:    map[string]int{},
		},
		{
			name:    "trailing hyphen is trimmed",
			content: "- Fixed by @carol-.",
			want:    map[string]int{"carol": 1},
		},
		{
			name:    "bracketed mention without a URL still counts",
			content: "- Fixed by [@superDingda] in #4236",
			want:    map[string]int{"superdingda": 1},
		},
		{
			name:    "pull-request links are not profile credits",
			content: "- Fix Z by @dave in [PR](https://github.com/wailsapp/wails/pull/1)",
			want:    map[string]int{"dave": 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := changelogMentions([]string{writeChangelog(t, tt.content)})
			if len(got) != len(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
			for login, n := range tt.want {
				if got[login] != n {
					t.Errorf("mentions[%q] = %d, want %d (full: %v)", login, got[login], n, got)
				}
			}
		})
	}
}

func TestCredit(t *testing.T) {
	defer func(m string) { creditMetric = m }(creditMetric)

	tests := []struct {
		name   string
		metric string
		c      Contributor
		want   int
	}{
		{"commits metric takes max of commits and mentions", "commits", Contributor{Commits: 10, Mentions: 3}, 10},
		{"mentions win when larger", "commits", Contributor{Commits: 2, Mentions: 5}, 5},
		{"prs metric uses merged PRs", "prs", Contributor{Commits: 231, PRs: 26}, 26},
		{"prs metric still lets mentions win", "prs", Contributor{PRs: 25, Mentions: 90}, 90},
		{"commit-only contributors floored into the tail", "prs", Contributor{Commits: 50}, 7},
		{"small commit-only counts kept as-is", "prs", Contributor{Commits: 2}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creditMetric = tt.metric
			if got := tt.c.Credit(); got != tt.want {
				t.Errorf("Credit() = %d, want %d", got, tt.want)
			}
		})
	}
}
