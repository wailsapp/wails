//go:build !production

package assetserver

import (
	"bytes"
	"testing"
)

var defaultIndexOutput []byte

func BenchmarkDefaultIndexHTML(b *testing.B) {
	for _, tt := range []struct{ language, file string }{
		{"en", "defaults/index.en.html"},
		{"zh", "defaults/index.zh.html"},
		{"unknown", ""},
	} {
		b.Run(tt.language, func(b *testing.B) {
			want := []byte("index.html not found")
			if tt.file != "" {
				var err error
				want, err = defaultHTML.ReadFile(tt.file)
				if err != nil {
					b.Fatal(err)
				}
			}
			b.ReportAllocs()
			for b.Loop() {
				defaultIndexOutput = defaultIndexHTML(tt.language)
			}
			if !bytes.Equal(defaultIndexOutput, want) {
				b.Fatal("fallback content differs from its embedded fixture")
			}
		})
	}
}
