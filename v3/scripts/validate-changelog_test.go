//go:build ignore

package main

import "testing"

func TestHistoricalShorthandReferenceCorrection(t *testing.T) {
	old := "- Fix build hangs with 50-100k+ files (#4939)"
	deleted := []changelogEntry{{Line: old, Section: "v3.0.0-alpha.79"}}
	for _, tc := range []struct {
		name, line, section string
		want                bool
	}{
		{"same reference", "- Fix build hangs with 50000-100000+ files (#4939)", "v3.0.0-alpha.79", true},
		{"other reference", "- Fix build hangs with 50000-100000+ files (#4940)", "v3.0.0-alpha.79", false},
		{"other release", "- Fix build hangs with 50000-100000+ files (#4939)", "v3.0.0-alpha.78", false},
		{"missing reference", "- Fix build hangs with 50000-100000+ files", "v3.0.0-alpha.79", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := isSameSourceCorrection(tc.line, tc.section, deleted); got != tc.want {
				t.Fatalf("isSameSourceCorrection() = %v, want %v", got, tc.want)
			}
		})
	}
	for _, line := range []string{"- New note (#4939) with more text", "- New note #4939", "- New note (other/repo#4939)"} {
		if got := changelogReferenceFromLine(line); got != "" {
			t.Fatalf("changelogReferenceFromLine(%q) = %q, want rejection", line, got)
		}
	}
}

func TestHistoricalTitleFormattingCorrection(t *testing.T) {
	original := "- Added `One Time Handlers` in the docs for `Listening to Events in JavaScript` by @contributor"
	formatted := "- Added *One Time Handlers* in the docs for *Listening to Events in JavaScript* by @contributor"
	for _, tc := range []struct {
		name, old, added, oldSection string
		want                         bool
	}{
		{"title formatting", original, formatted, "v3.0.0-alpha.79", true},
		{"changed title", original, "- Added *New Event Handlers* in the docs for *Listening to Events in JavaScript* by @contributor", "v3.0.0-alpha.79", false},
		{"moved release", original, formatted, "v3.0.0-alpha.78", false},
		{"unchanged line", original, original, "v3.0.0-alpha.79", false},
		{"numeric change", "- Handle 50 files", "- Handle 50000 files", "v3.0.0-alpha.79", false},
		{"code operator removed", "- Support `a*b`", "- Support *ab*", "v3.0.0-alpha.79", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			deleted := []changelogEntry{{Line: tc.old, Section: tc.oldSection}}
			if got := isSameSourceCorrection(tc.added, "v3.0.0-alpha.79", deleted); got != tc.want {
				t.Fatalf("isSameSourceCorrection() = %v, want %v", got, tc.want)
			}
		})
	}
}

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
