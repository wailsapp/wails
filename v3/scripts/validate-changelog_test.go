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
