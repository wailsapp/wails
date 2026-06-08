package monitor

import (
	"context"
	"runtime"
	"sync"
	"time"
)

// Sample is one periodic resource-usage reading for the monitored process. It
// shares a time axis with Trace records, so a consumer can correlate a span of
// binding calls/events with the RAM/CPU the process was using at that moment.
type Sample struct {
	Time time.Time `json:"time"`

	// Process-wide resident memory, in bytes (RSS from the OS). Zero if the
	// platform probe is unavailable.
	RSS uint64 `json:"rss"`
	// CPU usage since the previous sample, as a percentage of one core. May
	// exceed 100 on multi-core workloads.
	CPUPct float64 `json:"cpuPct"`

	// Go runtime stats — always available, no syscall.
	HeapAlloc  uint64 `json:"heapAlloc"`
	HeapSys    uint64 `json:"heapSys"`
	Goroutines int    `json:"goroutines"`
	NumGC      uint32 `json:"numGC"`

	// Thread / file-descriptor counts from the OS (zero if unavailable).
	Threads int `json:"threads"`
	FDs     int `json:"fds"`
}

// procProbe is the OS-specific reader. It returns RSS (bytes), cumulative CPU
// time used by the process, thread count and open-fd count. cpuTime is used by
// the sampler to derive a percentage between ticks. Unavailable fields are zero.
type procProbe struct {
	RSS     uint64
	CPUTime time.Duration
	Threads int
	FDs     int
}

// sampler owns the periodic sampling goroutine.
type sampler struct {
	interval time.Duration
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

// startSampler launches a goroutine that broadcasts a Sample every interval
// until ctx is cancelled. It is started by Start when a sink comes up.
func startSampler(s *Sink, interval time.Duration) *sampler {
	if interval <= 0 {
		interval = time.Second
	}
	ctx, cancel := context.WithCancel(context.Background())
	sp := &sampler{interval: interval, cancel: cancel}
	sp.wg.Add(1)
	go sp.run(ctx, s)
	return sp
}

func (sp *sampler) stop() {
	if sp == nil {
		return
	}
	sp.cancel()
	sp.wg.Wait()
}

func (sp *sampler) run(ctx context.Context, s *Sink) {
	defer sp.wg.Done()
	t := time.NewTicker(sp.interval)
	defer t.Stop()

	var (
		prevCPU  time.Duration
		prevTime time.Time
		havePrev bool
	)

	emit := func(now time.Time) {
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)

		smp := Sample{
			Time:       now,
			HeapAlloc:  ms.HeapAlloc,
			HeapSys:    ms.HeapSys,
			Goroutines: runtime.NumGoroutine(),
			NumGC:      ms.NumGC,
		}

		if p, ok := readProc(); ok {
			smp.RSS = p.RSS
			smp.Threads = p.Threads
			smp.FDs = p.FDs
			if havePrev {
				dt := now.Sub(prevTime).Seconds()
				if dt > 0 {
					smp.CPUPct = (p.CPUTime - prevCPU).Seconds() / dt * 100
				}
			}
			prevCPU = p.CPUTime
			prevTime = now
			havePrev = true
		}

		s.broadcast(Envelope{Type: MsgSample, Sample: &smp})
	}

	// Prime the CPU delta immediately so the first emitted sample (one tick
	// later) already has a percentage.
	if p, ok := readProc(); ok {
		prevCPU = p.CPUTime
		havePrev = true
	}
	prevTime = sinceStart()

	for {
		select {
		case <-ctx.Done():
			return
		case now := <-t.C:
			emit(now)
		}
	}
}

// sinceStart returns a monotonic-ish now for the very first delta. Scripts in
// this package cannot call time.Now in tests that need determinism, but the
// sampler is runtime-only, so a direct call is fine here.
func sinceStart() time.Time { return time.Now() }
