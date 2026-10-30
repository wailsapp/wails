//go:build !production && !ios && !android

package application

import (
	"log/slog"
	"testing"
)

func BenchmarkDefaultLogger(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = DefaultLogger(slog.LevelInfo)
	}
}
