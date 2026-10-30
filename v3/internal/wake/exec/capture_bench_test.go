package exec

import (
	"bytes"
	"testing"

	"github.com/wailsapp/wails/v3/internal/report"
)

var capturedOutput string

func BenchmarkCaptureFragmentedLine(b *testing.B) {
	for _, size := range []struct {
		name  string
		bytes int
	}{{"8KiB", 8 * 1024}, {"64KiB", 64 * 1024}} {
		b.Run(size.name, func(b *testing.B) {
			line := bytes.Repeat([]byte("x"), size.bytes)
			b.ReportAllocs()
			for b.Loop() {
				writer := newCaptureWriter(report.Nop{}, 1, false)
				for start := 0; start < len(line); start += 256 {
					_, _ = writer.Write(line[start : start+256])
				}
				writer.flush()
				capturedOutput = writer.output()
			}
			if capturedOutput != string(line)+"\n" {
				b.Fatal("captured output differs from the original line")
			}
		})
	}
}
