package main

import (
	"strings"
	"time"
)

// Scenario is one load profile. Rate 0 means unthrottled: send as fast as the
// transport accepts, which measures the ceiling and exercises the blocking
// backpressure path.
type Scenario struct {
	Name     string
	Rate     int // frames per second, 0 = unthrottled
	Size     int // total frame bytes including the 16-byte header
	Duration time.Duration
	Note     string
}

var allScenarios = []Scenario{
	// Rate sweep at a small frame. The event work showed call count was free on
	// the eval path; this asks the same question of the poll path, where the
	// cost per response is what matters rather than the cost per frame.
	{Name: "rate-100", Rate: 100, Size: 256},
	{Name: "rate-1000", Rate: 1000, Size: 256},
	{Name: "rate-5000", Rate: 5000, Size: 256},
	{Name: "rate-20000", Rate: 20000, Size: 256},

	// Size sweep at a fixed rate. On the eval path this is where retention
	// appeared, with a knee between 8 and 16 KB on macOS.
	{Name: "size-1KB", Rate: 100, Size: 1024},
	{Name: "size-64KB", Rate: 100, Size: 64 * 1024},
	{Name: "size-1MB", Rate: 100, Size: 1024 * 1024},

	// Constant byte rate (~4 MB/s), varying only frame size. Any step function
	// between these rows is attributable to size alone — the design that found
	// the eval knee.
	{Name: "iso-1KB", Rate: 4096, Size: 1024},
	{Name: "iso-8KB", Rate: 512, Size: 8 * 1024},
	{Name: "iso-64KB", Rate: 64, Size: 64 * 1024},
	{Name: "iso-1MB", Rate: 4, Size: 1024 * 1024},

	// Throughput ceiling and the blocking path. Swept across frame size
	// because the ceiling is set by different things at each end: per-frame
	// allocation at small sizes, memory bandwidth at large ones.
	{Name: "unthrottled-256B", Rate: 0, Size: 256},
	{Name: "unthrottled-4KB", Rate: 0, Size: 4 * 1024},
	{Name: "unthrottled-64KB", Rate: 0, Size: 64 * 1024},
	{Name: "unthrottled-256KB", Rate: 0, Size: 256 * 1024},
	{Name: "unthrottled-1MB", Rate: 0, Size: 1024 * 1024},
	{Name: "unthrottled-4MB", Rate: 0, Size: 4 * 1024 * 1024},
}

// uploadVariant is one JS→Go throughput measurement. Connection count matters
// there in a way it does not Go→JS: the client serialises sends per connection,
// so one connection can have only one frame in flight and its ceiling is
// frameSize/roundTrip. Concurrency is the only way past that.
type uploadVariant struct {
	Name  string
	Size  int
	Conns int
}

var uploadVariants = []uploadVariant{
	{Name: "up-1KB-x1", Size: 1024, Conns: 1},
	{Name: "up-64KB-x1", Size: 64 * 1024, Conns: 1},
	{Name: "up-512KB-x1", Size: 512 * 1024, Conns: 1},
	{Name: "up-4MB-x1", Size: 4 * 1024 * 1024, Conns: 1},
	{Name: "up-64KB-x4", Size: 64 * 1024, Conns: 4},
	{Name: "up-512KB-x4", Size: 512 * 1024, Conns: 4},
	{Name: "up-512KB-x8", Size: 512 * 1024, Conns: 8},
}

func selectedScenarios(only string) []Scenario {
	if strings.TrimSpace(only) == "" {
		return allScenarios
	}
	want := map[string]bool{}
	for _, n := range strings.Split(only, ",") {
		want[strings.TrimSpace(n)] = true
	}
	var out []Scenario
	for _, sc := range allScenarios {
		if want[sc.Name] {
			out = append(out, sc)
		}
	}
	return out
}
