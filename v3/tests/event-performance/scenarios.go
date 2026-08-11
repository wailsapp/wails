package main

import (
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Scenario sweeps the two axes independently — that separation is the whole
// point of the harness. Growth on the rate sweep implicates a per-call leak;
// growth on the size sweep implicates payload retention.
type Scenario struct {
	Name         string
	Rate         int // events/sec; 0 = idle
	PayloadBytes int // bytes of filler per event
	Burst        int // if >0, deliver the mean rate in bursts of this size
	Duration     time.Duration
	Note         string

	// MixedSource runs a second emitter on the main thread concurrently with
	// the goroutine emitter. dispatchOnMainThread has an inline fast path, so a
	// main-thread emit executes its eval immediately while a goroutine emit is
	// still sitting in the dispatch queue — the two can invert. Both emitters
	// draw from the same monotonic counter, so any inversion shows as a reorder.
	MixedSource bool

	// AltPayloadBytes, when >0, alternates every other event to this size.
	// This is what exercises mixed inline/by-reference delivery: small events
	// dispatch synchronously, large ones go through an async fetch, and the
	// sequence numbers prove nothing overtakes anything.
	AltPayloadBytes int

	counter *atomic.Int64 // shared across copies of the struct

	// emitMu serialises taking a sequence number with the dispatch call that
	// carries it. Without it, two emitters can take seq N and N+1 and then
	// reach the framework in the other order, and the reorders the JS side
	// counts are this harness racing itself rather than anything the framework
	// did. Only used by MixedSource scenarios.
	emitMu *sync.Mutex
}

func (s Scenario) nextSeq() int64   { return s.counter.Add(1) - 1 }
func (s Scenario) sentSoFar() int64 { return s.counter.Load() }
func (s Scenario) resetCounter()    { s.counter.Store(0) }

func newScenario(s Scenario) Scenario {
	s.counter = &atomic.Int64{}
	s.emitMu = &sync.Mutex{}
	return s
}

// allScenarios is ordered deliberately: idle first (it is the baseline every
// other result is read against), pathological LAST because it is expected to
// kill the web process and everything after it would measure a corpse.
func allScenarios() []Scenario {
	var out []Scenario
	add := func(s Scenario) { out = append(out, newScenario(s)) }

	add(Scenario{
		Name: "idle", Rate: 0,
		Note: "baseline: footprint drift from merely having a webview open",
	})

	// Rate sweep — fixed 64B payload. Growth here implicates a per-call leak.
	for _, r := range []int{10, 100, 500, 1000, 2500, 5000} {
		add(Scenario{
			Name: "rate-" + itoa(r), Rate: r, PayloadBytes: 64,
			Note: "rate sweep @64B",
		})
	}

	// Size sweep — fixed 100 ev/s. Growth here implicates payload retention.
	for _, sz := range []struct {
		label string
		bytes int
	}{{"1KB", 1 << 10}, {"16KB", 16 << 10}, {"256KB", 256 << 10}, {"1MB", 1 << 20}} {
		add(Scenario{
			Name: "size-" + sz.label, Rate: 100, PayloadBytes: sz.bytes,
			Note: "size sweep @100 ev/s",
		})
	}

	// Iso-byte-rate sweep: hold ~4 MB/s constant and vary only the payload size.
	// Retention appears to be a step function of payload size rather than a
	// smooth function of bytes, so this is what locates the knee — at a fixed
	// byte rate, any difference between these rows is attributable to size alone.
	for _, iso := range []struct {
		label string
		bytes int
		rate  int
	}{
		{"1KB", 1 << 10, 4096},
		{"2KB", 2 << 10, 2048},
		{"4KB", 4 << 10, 1024},
		{"8KB", 8 << 10, 512},
		{"16KB", 16 << 10, 256},
		{"32KB", 32 << 10, 128},
		{"64KB", 64 << 10, 64},
	} {
		add(Scenario{
			Name: "iso-" + iso.label, Rate: iso.rate, PayloadBytes: iso.bytes,
			Note: "iso-byte-rate ~4 MB/s",
		})
	}

	// Second iso family at ~32 MB/s, covering 32 KB → 1 MB. The switchover size
	// is platform-specific: macOS/WebKit-Cocoa flips between 8 KB and 16 KB,
	// while WebKitGTK is still flat at 64 KB, so locating its knee needs both
	// larger payloads and enough throughput to make retention obvious.
	for _, iso := range []struct {
		label string
		bytes int
		rate  int
	}{
		{"32KB", 32 << 10, 1024},
		{"64KB", 64 << 10, 512},
		{"128KB", 128 << 10, 256},
		{"256KB", 256 << 10, 128},
		{"512KB", 512 << 10, 64},
		{"1MB", 1 << 20, 32},
	} {
		add(Scenario{
			Name: "iso2-" + iso.label, Rate: iso.rate, PayloadBytes: iso.bytes,
			Note: "iso-byte-rate ~32 MB/s",
		})
	}

	// Ordering under mixed delivery modes. Alternating a well-under-threshold
	// payload with a well-over-threshold one means consecutive sequence numbers
	// travel by different mechanisms; any reorder shows up immediately.
	add(Scenario{
		Name: "interleave", Rate: 200, PayloadBytes: 1 << 10, AltPayloadBytes: 64 << 10,
		Note: "alternating 1KB/64KB — proves ordering across inline vs by-reference",
	})

	// Two concurrent emitters, one per thread-of-origin. This is the ordering
	// case the single-emitter scenarios cannot see.
	add(Scenario{
		Name: "mixedsource", Rate: 500, PayloadBytes: 256, MixedSource: true,
		Note: "goroutine + main-thread emitters sharing one sequence",
	})

	// The shape real apps actually produce.
	add(Scenario{
		Name: "burst", Rate: 1000, PayloadBytes: 256, Burst: 20,
		Note: "mean 1000 ev/s in bursts of 20 @256B",
	})

	// Direct reproduction of WebKit #215729's shape. Expected to terminate the
	// web process; that is a recorded result, not an error.
	add(Scenario{
		Name: "pathological", Rate: 10, PayloadBytes: 8 << 20,
		Note: "10 ev/s @8MB — expect web process termination",
	})

	return out
}

func selectedScenarios(only string) []Scenario {
	all := allScenarios()
	if strings.TrimSpace(only) == "" {
		return all
	}
	want := map[string]bool{}
	for _, n := range strings.Split(only, ",") {
		want[strings.TrimSpace(n)] = true
	}
	var out []Scenario
	for _, s := range all {
		if want[s.Name] {
			out = append(out, s)
		}
	}
	return out
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var b [20]byte
	p := len(b)
	for i > 0 {
		p--
		b[p] = byte('0' + i%10)
		i /= 10
	}
	return string(b[p:])
}
