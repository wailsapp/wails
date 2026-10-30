package term

import (
	"fmt"
	"strings"
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

var tableResult string

func BenchmarkRenderTemplateTable(b *testing.B) {
	rows := [][]string{
		{"Name", "example"},
		{"Directory", "/home/user/example"},
		{"Framework", "React"},
		{"Language", "TypeScript"},
		{"Bindings", "Classes"},
	}
	want := "  Name       example\n  Directory  /home/user/example\n  Framework  React\n  Language   TypeScript\n  Bindings   Classes\n"
	if got := RenderTable(rows); got != want {
		b.Fatalf("table = %q, want %q", got, want)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		tableResult = RenderTable(rows)
	}
}

func BenchmarkPadRight(b *testing.B) {
	want := "Name" + strings.Repeat(" ", 76)
	if got := padRight("Name", 80); got != want {
		b.Fatalf("padding = %q, want %q", got, want)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		tableResult = padRight("Name", 80)
	}
}

var measuredDisplayWidth int

func BenchmarkDisplayWidth(b *testing.B) {
	for _, tc := range []struct {
		name, text string
	}{
		{"ASCII", "Framework"},
		{"Unicode", "界e\u0301👩‍💻"},
		{"ANSI", "\x1b[31mFramework\x1b[0m"},
		{"Multiline", "Framework\n界界\n"},
	} {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				measuredDisplayWidth = displayWidth(tc.text)
			}
		})
	}
}
