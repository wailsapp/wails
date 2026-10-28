package exec

import (
	"bytes"
	"testing"

	"github.com/wailsapp/wails/v3/internal/report"
)

var capturedOutput string

func BenchmarkCaptureFragmentedLine(b *testing.B) {
	line := bytes.Repeat([]byte("x"), 8*1024)
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
}
