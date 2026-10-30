package cmds

import "testing"

var joinedGoTags string

func BenchmarkJoinTags(b *testing.B) {
	for _, tc := range []struct {
		name string
		tags []string
		want string
	}{
		{"empty", nil, ""},
		{"single", []string{"production"}, "production"},
		{"four", []string{"production", "desktop", "gtk4", "custom_feature"}, "production desktop gtk4 custom_feature"},
		{"eight", []string{"production", "desktop", "gtk4", "custom_feature", "sqlite", "debug", "profile", "integration"}, "production desktop gtk4 custom_feature sqlite debug profile integration"},
	} {
		b.Run(tc.name, func(b *testing.B) {
			if got := joinTags(tc.tags); got != tc.want {
				b.Fatalf("joined tags = %q, want %q", got, tc.want)
			}
			b.ReportAllocs()
			for b.Loop() {
				joinedGoTags = joinTags(tc.tags)
			}
		})
	}
}
