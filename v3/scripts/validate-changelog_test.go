//go:build ignore

package main

import "testing"

func TestIsSameSourceCorrection(t *testing.T) {
	deleted := []changelogEntry{
		{
			Line:    "- Fixes issue with incorrect handling of empty strings in JSON parsing in [PR](https://github.com/wailsapp/wails/pull/5985) by @taliesin-ai",
			Section: "v3.0.0-beta.9",
		},
	}
	added := "- Update website nanoid lockfiles to patched 3.3.18 in [PR](https://github.com/wailsapp/wails/pull/5985) by @taliesin-ai"
	if !isSameSourceCorrection(added, "v3.0.0-beta.9", deleted) {
		t.Fatal("isSameSourceCorrection() rejected a replacement citing the same PR")
	}
}

func TestIsSameSourceCorrectionRejectsDifferentSource(t *testing.T) {
	deleted := []changelogEntry{
		{Line: "- Old description in [PR](https://github.com/wailsapp/wails/pull/5984) by @contributor", Section: "v3.0.0-beta.9"},
	}
	added := "- New description in [PR](https://github.com/wailsapp/wails/pull/5985) by @contributor"
	if isSameSourceCorrection(added, "v3.0.0-beta.9", deleted) {
		t.Fatal("isSameSourceCorrection() accepted a replacement citing a different PR")
	}
}

func TestIsSameSourceCorrectionRejectsDifferentSection(t *testing.T) {
	deleted := []changelogEntry{
		{Line: "- Old description in [PR](https://github.com/wailsapp/wails/pull/5985) by @contributor", Section: "Unreleased"},
		{Line: "- Older description in [PR](https://github.com/wailsapp/wails/pull/5985) by @contributor", Section: "v3.0.0-beta.8"},
	}
	added := "- New description in [PR](https://github.com/wailsapp/wails/pull/5985) by @contributor"
	if isSameSourceCorrection(added, "v3.0.0-beta.9", deleted) {
		t.Fatal("isSameSourceCorrection() accepted a replacement from another section")
	}
}

func TestIsSameSourceCorrectionRejectsNewOrSourceLessEntry(t *testing.T) {
	tests := []struct {
		name    string
		added   string
		deleted []changelogEntry
	}{
		{
			name:  "new entry",
			added: "- New description in [PR](https://github.com/wailsapp/wails/pull/5985) by @contributor",
		},
		{
			name:    "source-less replacement",
			added:   "- Correct an old release note",
			deleted: []changelogEntry{{Line: "- Incorrect old release note", Section: "v3.0.0-beta.9"}},
		},
		{
			name:    "unchanged line",
			added:   "- Existing entry in [PR](https://github.com/wailsapp/wails/pull/5985) by @contributor",
			deleted: []changelogEntry{{Line: "- Existing entry in [PR](https://github.com/wailsapp/wails/pull/5985) by @contributor", Section: "v3.0.0-beta.9"}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if isSameSourceCorrection(test.added, "v3.0.0-beta.9", test.deleted) {
				t.Fatal("isSameSourceCorrection() accepted a non-correction")
			}
		})
	}
}

func TestDeletedChangelogEntriesPreserveSection(t *testing.T) {
	oldLine := "- Old description in [PR](https://github.com/wailsapp/wails/pull/5985) by @contributor"
	base := "## [Unreleased]\n\n" + oldLine + "\n\n## v3.0.0-beta.9 - 2026-08-16\n\n" + oldLine + "\n"
	current := "## [Unreleased]\n\n" + oldLine + "\n\n## v3.0.0-beta.9 - 2026-08-16\n\n- New description in [PR](https://github.com/wailsapp/wails/pull/5985) by @contributor\n"

	got := deletedChangelogEntries(base, current, []string{oldLine})
	want := []changelogEntry{{Line: oldLine, Section: "v3.0.0-beta.9"}}
	if len(got) != len(want) || got[0] != want[0] {
		t.Fatalf("deletedChangelogEntries() = %#v, want %#v", got, want)
	}
}

func TestPullRequestReferenceFromLineRejectsEmbeddedTrustedURL(t *testing.T) {
	tests := []string{
		"- Entry in [PR](https://attacker.example/https://github.com/wailsapp/wails/pull/5985) by @contributor",
		"- Entry in [PR](https://github.com/wailsapp/wails/pull/5985.attacker.example) by @contributor",
		"- Entry in [PR](https://github.com/wailsapp/wails/pull/5985?redirect=https://attacker.example) by @contributor",
	}
	for _, line := range tests {
		if got := pullRequestReferenceFromLine(line); got != "" {
			t.Fatalf("pullRequestReferenceFromLine(%q) = %q, want rejection", line, got)
		}
	}
}

func TestPublishedBackfill(t *testing.T) {
	bullet := "- Fix keyboard shortcuts in [PR](https://github.com/wailsapp/wails/pull/5902) by @julianstorer"
	notes := map[string]string{"v3.0.0-beta.9": "## Fixed\n" + bullet + "\n"}
	for _, tc := range []struct {
		name, line, section string
		want                bool
	}{
		{"same release exact entry", bullet, "v3.0.0-beta.9", true},
		{"wrong release", bullet, "v3.0.0-beta.8", false},
		{"altered description", "- New feature in [PR](https://github.com/wailsapp/wails/pull/5902) by @julianstorer", "v3.0.0-beta.9", false},
		{"substring", bullet + " and another change", "v3.0.0-beta.9", false},
		{"missing evidence", bullet, "v3.0.0-beta.10", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := isPublishedBackfill(tc.line, tc.section, notes); got != tc.want {
				t.Fatalf("isPublishedBackfill() = %v, want %v", got, tc.want)
			}
		})
	}
	sourceLess := "- Unverified note"
	if isPublishedBackfill(sourceLess, "v3.0.0-beta.9", map[string]string{"v3.0.0-beta.9": sourceLess}) {
		t.Fatal("accepted an entry without a canonical PR reference")
	}
}
