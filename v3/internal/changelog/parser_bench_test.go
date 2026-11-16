package changelog

import (
	"strconv"
	"strings"
	"testing"
)

var parsedChangelog *ChangelogEntry

func BenchmarkParseContent(b *testing.B) {
	parser := NewParser()
	for _, entries := range []int{32, 256, 1024} {
		b.Run(strconv.Itoa(entries), func(b *testing.B) {
			content := "## Fixed\n" + strings.Repeat("- Preserve existing application behavior\n", entries)
			entry, err := parser.ParseContent(strings.NewReader(content))
			if err != nil || len(entry.Fixed) != entries || entry.Fixed[0] != "Preserve existing application behavior" {
				b.Fatalf("invalid benchmark fixture: %+v, %v", entry, err)
			}
			b.ReportAllocs()
			b.SetBytes(int64(len(content)))
			for b.Loop() {
				parsedChangelog, err = parser.ParseContent(strings.NewReader(content))
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
