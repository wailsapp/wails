package changelog

import (
	"fmt"
	"testing"
)

var formattedChangelog string

func BenchmarkFormatEntries(b *testing.B) {
	items := make([]string, 20)
	for i := range items {
		items[i] = fmt.Sprintf("Entry %d: preserve 100%% of the original text and Unicode ✨", i)
	}
	entry := &ChangelogEntry{
		Added: items, Changed: items, Fixed: items,
		Deprecated: items, Removed: items, Security: items,
	}
	for _, format := range []struct {
		name string
		call func() string
	}{
		{"changelog", entry.FormatForChangelog},
		{"release", entry.FormatForRelease},
	} {
		b.Run(format.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				formattedChangelog = format.call()
			}
		})
	}
}
