package term

import (
	"fmt"
	"testing"
)

var hyperlinkResult string

func BenchmarkHyperlink(b *testing.B) {
	url := "https://wails.io/docs/" + b.Name()
	label := "Wails documentation"
	want := fmt.Sprintf("\x1b]8;;%s\x1b\\\x1b[4m%s\x1b[0m\x1b]8;;\x1b\\", url, label)
	if got := Hyperlink(url, label); got != want {
		b.Fatalf("hyperlink bytes = %q, want %q", got, want)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		hyperlinkResult = Hyperlink(url, label)
	}
}
