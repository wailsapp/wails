package webview

import (
	"context"
	"testing"
)

func BenchmarkRequestTrackerLifecycle(b *testing.B) {
	for _, mode := range []string{"tracked", "keepalive", "unmarked"} {
		b.Run(mode, func(b *testing.B) {
			tracker := NewRequestTracker(context.Background())
			defer tracker.Close()
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				var token string
				switch mode {
				case "tracked":
					token = tracker.Start("benchmark")
				case "keepalive":
					token = tracker.StartKeepalive("benchmark")
				}
				request := tracker.Wrap(trackedTestRequest(token))
				if err := request.Close(); err != nil {
					b.Fatal(err)
				}
				if mode == "keepalive" {
					tracker.Cancel("benchmark")
				}
			}
		})
	}
}
