//go:build ignore

package main

import "testing"

func TestIsSameSourceCorrection(t *testing.T) {
	deleted := []string{
		"- Fixes issue with incorrect handling of empty strings in JSON parsing in [PR](https://github.com/wailsapp/wails/pull/5985) by @taliesin-ai",
	}
	added := "- Update website nanoid lockfiles to patched 3.3.18 in [PR](https://github.com/wailsapp/wails/pull/5985) by @taliesin-ai"
	if !isSameSourceCorrection(added, deleted) {
		t.Fatal("isSameSourceCorrection() rejected a replacement citing the same PR")
	}
}

func TestIsSameSourceCorrectionRejectsDifferentSource(t *testing.T) {
	deleted := []string{
		"- Old description in [PR](https://github.com/wailsapp/wails/pull/5984) by @contributor",
	}
	added := "- New description in [PR](https://github.com/wailsapp/wails/pull/5985) by @contributor"
	if isSameSourceCorrection(added, deleted) {
		t.Fatal("isSameSourceCorrection() accepted a replacement citing a different PR")
	}
}

func TestIsSameSourceCorrectionRejectsNewOrSourceLessEntry(t *testing.T) {
	tests := []struct {
		name    string
		added   string
		deleted []string
	}{
		{
			name:  "new entry",
			added: "- New description in [PR](https://github.com/wailsapp/wails/pull/5985) by @contributor",
		},
		{
			name:    "source-less replacement",
			added:   "- Correct an old release note",
			deleted: []string{"- Incorrect old release note"},
		},
		{
			name:    "unchanged line",
			added:   "- Existing entry in [PR](https://github.com/wailsapp/wails/pull/5985) by @contributor",
			deleted: []string{"- Existing entry in [PR](https://github.com/wailsapp/wails/pull/5985) by @contributor"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if isSameSourceCorrection(test.added, test.deleted) {
				t.Fatal("isSameSourceCorrection() accepted a non-correction")
			}
		})
	}
}
