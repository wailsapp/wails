package changelog

import (
	"strconv"
	"strings"
	"testing"
)

var validatedChangelog ValidationResult

func BenchmarkValidateContent(b *testing.B) {
	validator := NewValidator()
	for _, fixture := range []struct {
		entries     int
		capitalized bool
	}{{32, false}, {256, false}, {1024, false}, {1024, true}} {
		name := strconv.Itoa(fixture.entries)
		bullet := "- preserve existing application behavior\n"
		if fixture.capitalized {
			name += "Capitalized"
			bullet = "- Preserve existing application behavior\n"
		}
		b.Run(name, func(b *testing.B) {
			content := "## Fixed\n" + strings.Repeat(bullet, fixture.entries)
			if result := validator.ValidateContent(content); !result.Valid || len(result.Errors) != 0 {
				b.Fatalf("invalid benchmark fixture: %+v", result)
			}
			b.ReportAllocs()
			b.SetBytes(int64(len(content)))
			for b.Loop() {
				validatedChangelog = validator.ValidateContent(content)
			}
		})
	}
}
