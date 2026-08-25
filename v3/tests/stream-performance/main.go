// Command streamperf is the load harness for GoStream.
//
// It answers the questions the design could not answer by argument: how much
// throughput a held-poll transport sustains Go→JS, whether anything is lost or
// reordered under load, what it costs in host and content-process memory, and
// how long a Send blocks.
//
// It drives StreamConn.Send directly — the transport entry point — so what is
// measured is the transport, not application-level fanout.
//
// Run (macOS):
//
//	go run ./tests/stream-performance -duration 20s
//
// One scenario only:
//
//	go run ./tests/stream-performance -only rate-5000
package main

import (
	"context"
	_ "embed"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed assets/index.html
var indexHTML []byte

var (
	flagDuration = flag.Duration("duration", 20*time.Second, "measurement window per scenario")
	flagSettle   = flag.Duration("settle", 2*time.Second, "settle period before each scenario (discarded)")
	flagOut      = flag.String("out", "", "output directory (default: ./streamperf-results)")
	flagOnly     = flag.String("only", "", "comma-separated scenario names to run (default: all)")
	flagSample   = flag.Duration("sample", 250*time.Millisecond, "footprint sample interval")
	flagReloads  = flag.Int("reloads", 0, "instead of the sweep, reload the page N times and report connection lifecycle")
	flagUpload   = flag.Bool("upload", false, "instead of the sweep, measure JS→Go throughput across frame sizes and connection counts")
)

// Frame layout on the wire, chosen so neither side pays a JSON parse per frame
// and the measurement is of the transport rather than of encoding:
//
//	seq   int64   big endian
//	tms   float64 big endian — Go monotonic ms since process start
//	pad   n bytes
const frameHeaderBytes = 16

type jsReport struct {
	Epoch      int     `json:"epoch"`
	Received   int64   `json:"received"`
	Bytes      int64   `json:"bytes"`
	Drops      int64   `json:"drops"`
	Reorders   int64   `json:"reorders"`
	Polls      int64   `json:"polls"`
	MaxBatch   int64   `json:"maxBatch"`
	Echoes     int64   `json:"echoes"`
	EchoErrors int64   `json:"echoErrors"`
	DeltaMin   float64 `json:"deltaMin"`
	P50        float64 `json:"p50"`
	P95        float64 `json:"p95"`
	P99        float64 `json:"p99"`
	Max        float64 `json:"max"`
}

type harness struct {
	start time.Time

	basePids map[int]bool
	ourPids  []int
	hostPid  int

	mu       sync.Mutex
	latestJS jsReport
	sendUS   []float64
	epoch    int

	win    *application.WebviewWindow
	connCh chan *application.StreamConn
	ctlCh  chan *application.StreamConn

	// JS→Go throughput counters, filled by the sink stream.
	sinkBytes  atomic.Int64
	sinkFrames atomic.Int64

	echoes atomic.Int64

	// Connection lifecycle, for the reload check: a reload must close the
	// previous page's connection, not leave it live alongside the new one.
	connects    atomic.Int64
	disconnects atomic.Int64
	live        atomic.Int64
	maxLive     atomic.Int64

	results []*scenarioResult
}

func main() {
	flag.Parse()

	if !samplerSupported {
		fmt.Println("NOTE: memory sampling is unsupported on this platform;")
		fmt.Println("throughput and correctness metrics will still be collected.")
	}

	h := &harness{
		start:   time.Now(),
		hostPid: os.Getpid(),
		connCh:  make(chan *application.StreamConn, 4),
		ctlCh:   make(chan *application.StreamConn, 4),
	}

	h.basePids = map[int]bool{}
	if samplerSupported {
		pids, err := webContentPids()
		if err != nil {
			log.Fatalf("pre-launch process enumeration failed: %v", err)
		}
		for _, p := range pids {
			h.basePids[p] = true
		}
		log.Printf("pre-launch content processes on this machine: %d", len(h.basePids))
	}

	app := application.New(application.Options{
		Name:        "streamperf",
		Description: "Wails v3 GoStream load harness",
		Assets: application.AssetOptions{
			Handler:        h.handler(),
			DisableLogging: true,
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// The stream under test. The handler holds the connection open for the
	// whole run; the emitter goroutine below does the sending.
	app.HandleStream("perf", func(c *application.StreamConn) {
		h.connects.Add(1)
		n := h.live.Add(1)
		for {
			current := h.maxLive.Load()
			if n <= current || h.maxLive.CompareAndSwap(current, n) {
				break
			}
		}
		log.Printf("stream connected: %s (live=%d)", c.Name(), h.live.Load())

		select {
		case h.connCh <- c:
		default: // reload mode does not consume connections
		}

		<-c.Context().Done()
		h.live.Add(-1)
		h.disconnects.Add(1)
		log.Printf("stream disconnected (live=%d)", h.live.Load())
	})

	// A second stream, used only to exercise the JS→Go direction: the frontend
	// sends, Go echoes the same bytes back.
	app.HandleStream("echo", func(c *application.StreamConn) {
		defer c.Close()
		for {
			frame, err := c.Receive()
			if err != nil {
				return
			}
			h.echoes.Add(1)
			if err := c.Send(frame); err != nil {
				return
			}
		}
	})

	// Control channel: Go tells the page what upload to run. Using a stream for
	// this also exercises several simultaneous streams on one window, which is
	// the multiplexing the design rests on.
	app.HandleStream("ctl", func(c *application.StreamConn) {
		select {
		case h.ctlCh <- c:
		default:
		}
		<-c.Context().Done()
	})

	// Upload sink: counts what the frontend sends and discards it, so the
	// measurement is of the transport rather than of anything Go does with the
	// bytes afterwards.
	app.HandleStream("sink", func(c *application.StreamConn) {
		for {
			frame, err := c.Receive()
			if err != nil {
				return
			}
			h.sinkFrames.Add(1)
			h.sinkBytes.Add(int64(len(frame)))
		}
	})

	h.win = app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:   "streamperf",
		Title:  "streamperf",
		Width:  560,
		Height: 420,
		URL:    "/",
	})

	win := h.win
	var once sync.Once
	win.OnWindowEvent(events.Common.WindowRuntimeReady, func(*application.WindowEvent) {
		once.Do(func() { go h.runAll(app) })
	})

	if err := app.Run(); err != nil {
		log.Fatalf("app.Run: %v", err)
	}
}

func (h *harness) handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/harness/report", func(w http.ResponseWriter, r *http.Request) {
		var rep jsReport
		if err := json.NewDecoder(r.Body).Decode(&rep); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		h.mu.Lock()
		if rep.Epoch == h.epoch {
			h.latestJS = rep
		}
		h.mu.Unlock()
		w.WriteHeader(http.StatusNoContent)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(indexHTML)
	})

	return mux
}

func (h *harness) nowMS() float64 {
	return float64(time.Since(h.start).Nanoseconds()) / 1e6
}

func (h *harness) runAll(app *application.App) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("harness panic: %v", r)
		}
		app.Quit()
	}()

	if *flagReloads > 0 {
		h.runReloadCheck()
		return
	}

	if *flagUpload {
		h.runUploadCheck()
		return
	}

	// Wait for the page to connect its stream before doing anything else.
	var conn *application.StreamConn
	select {
	case conn = <-h.connCh:
	case <-time.After(15 * time.Second):
		log.Printf("FATAL: frontend never connected the 'perf' stream")
		return
	}

	time.Sleep(1500 * time.Millisecond)

	if samplerSupported {
		if err := h.discoverWebContent(); err != nil {
			log.Printf("FATAL: %v", err)
			fmt.Println("Could not attribute a WebContent process to this app.")
			return
		}
		log.Printf("our WebContent pid(s): %v  (host pid %d)", h.ourPids, h.hostPid)
	}

	outDir := *flagOut
	if outDir == "" {
		outDir = "streamperf-results"
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		log.Printf("mkdir %s: %v", outDir, err)
		return
	}
	abs, _ := filepath.Abs(outDir)

	scenarios := selectedScenarios(*flagOnly)
	log.Printf("running %d scenario(s), %s each, results → %s", len(scenarios), *flagDuration, abs)

	for i, sc := range scenarios {
		log.Printf("[%d/%d] %s", i+1, len(scenarios), sc.Name)
		res := h.runScenario(conn, sc)
		h.results = append(h.results, res)
		if err := writeScenarioCSV(outDir, res); err != nil {
			log.Printf("  csv: %v", err)
		}
		log.Printf("  → %s", res.oneLine())
	}

	if err := writeSummary(outDir, h.results); err != nil {
		log.Printf("summary.json: %v", err)
	}
	fmt.Printf("\nResults written to %s\n", abs)
	fmt.Println(renderConsoleTable(h.results))
}

// runUploadCheck measures JS→Go throughput across frame sizes and connection
// counts.
//
// The shape to expect is different from Go→JS. The client serialises sends per
// connection, because concurrent fetch POSTs do not preserve order, so one
// connection can only have one frame in flight: its ceiling is
// frameSize/roundTrip regardless of how fast either side can move bytes. More
// connections is the only way past that, which is why the matrix has a
// connection axis at all.
func (h *harness) runUploadCheck() {
	var ctl *application.StreamConn
	select {
	case ctl = <-h.ctlCh:
	case <-time.After(15 * time.Second):
		log.Printf("FATAL: frontend never connected the 'ctl' stream")
		return
	}

	dur := *flagDuration
	log.Printf("upload check: %d variants, %s each", len(uploadVariants), dur)
	fmt.Printf("\n%-14s %10s %8s %12s %12s %10s\n",
		"variant", "frameKB", "conns", "frames", "frames/s", "MB/s")
	fmt.Println(strings.Repeat("-", 72))

	for _, v := range uploadVariants {
		cmd, _ := json.Marshal(map[string]any{"cmd": "upload", "size": v.Size, "conns": v.Conns})
		if err := ctl.Send(cmd); err != nil {
			log.Printf("  %s: ctl send failed: %v", v.Name, err)
			continue
		}

		// Settle: let the uploaders reach steady state, then zero the counters
		// so the measured window excludes ramp-up.
		time.Sleep(*flagSettle)
		h.sinkBytes.Store(0)
		h.sinkFrames.Store(0)

		start := time.Now()
		time.Sleep(dur)
		elapsed := time.Since(start).Seconds()
		bytes := h.sinkBytes.Load()
		frames := h.sinkFrames.Load()

		stop, _ := json.Marshal(map[string]any{"cmd": "stop"})
		_ = ctl.Send(stop)
		time.Sleep(1500 * time.Millisecond)

		fmt.Printf("%-14s %10.0f %8d %12d %12.1f %10.1f\n",
			v.Name, float64(v.Size)/1024, v.Conns, frames,
			float64(frames)/elapsed,
			float64(bytes)/elapsed/(1024*1024))
	}
	fmt.Println()
}

// runReloadCheck reloads the page repeatedly and reports what Go observed.
//
// A reload is the closest thing this transport has to a socket closing, and
// there is no cancellation from the platform layer to signal it. The page comes
// back with a new session id, and the previous session has to be superseded, or
// the app holds two live connections to the same stream until the TTL sweep -
// up to streamSessionTTL + streamSessionSweep. This measures which happens, and
// it also checks the assumption the fix rests on: that a reloaded page presents
// the same window id.
func (h *harness) runReloadCheck() {
	log.Printf("reload check: %d reloads", *flagReloads)

	for i := 0; i < *flagReloads; i++ {
		time.Sleep(2 * time.Second)
		log.Printf("  reload %d/%d (live=%d)", i+1, *flagReloads, h.live.Load())
		h.win.ExecJS("location.reload()")
	}

	// Let the last page settle.
	time.Sleep(3 * time.Second)

	log.Printf("reload check complete: connects=%d disconnects=%d live=%d maxLive=%d",
		h.connects.Load(), h.disconnects.Load(), h.live.Load(), h.maxLive.Load())
	if h.maxLive.Load() > 1 {
		log.Printf("  NOTE: %d connections were live at once - a reload did not close its predecessor promptly",
			h.maxLive.Load())
	}
}

func (h *harness) runScenario(conn *application.StreamConn, sc Scenario) *scenarioResult {
	dur := *flagDuration
	if sc.Duration > 0 {
		dur = sc.Duration
	}

	h.mu.Lock()
	h.epoch++
	epoch := h.epoch
	h.latestJS = jsReport{Epoch: epoch}
	h.sendUS = nil
	h.mu.Unlock()

	// Epoch marker: a frame with seq -1 tells the page to clear its counters,
	// so a scenario measures only its own traffic.
	_ = conn.Send(markerFrame(epoch))

	// Settle under load, then discard.
	stop := h.startEmitter(conn, sc)
	time.Sleep(*flagSettle)
	stop()

	h.mu.Lock()
	h.epoch++
	epoch = h.epoch
	h.latestJS = jsReport{Epoch: epoch}
	h.sendUS = nil
	h.mu.Unlock()
	_ = conn.Send(markerFrame(epoch))
	time.Sleep(250 * time.Millisecond)

	res := &scenarioResult{
		Name:        sc.Name,
		PayloadSize: sc.Size,
		TargetRate:  sc.Rate,
		Duration:    dur,
		Epoch:       epoch,
	}

	echoBase := h.echoes.Load()
	sent := h.startEmitterCounted(conn, sc)
	deadline := time.Now().Add(dur)
	ticker := time.NewTicker(*flagSample)
	defer ticker.Stop()

	for time.Now().Before(deadline) {
		<-ticker.C
		res.Samples = append(res.Samples, h.sample())
	}

	stopped := sent()
	// Let the last frames drain to the page before reading its counters.
	time.Sleep(1500 * time.Millisecond)

	h.mu.Lock()
	res.JS = h.latestJS
	sendUS := append([]float64(nil), h.sendUS...)
	h.mu.Unlock()

	res.Sent = stopped
	res.SendP50, res.SendP99 = pctl(sendUS, 50), pctl(sendUS, 99)
	res.Echoes = h.echoes.Load() - echoBase
	res.finalise()
	return res
}

// markerFrame carries seq = -1 and the epoch in the payload slot.
func markerFrame(epoch int) []byte {
	buf := make([]byte, frameHeaderBytes+8)
	var minusOne int64 = -1
	binary.BigEndian.PutUint64(buf[0:8], uint64(minusOne))
	binary.BigEndian.PutUint64(buf[8:16], math.Float64bits(0))
	binary.BigEndian.PutUint64(buf[16:24], uint64(int64(epoch)))
	return buf
}

func (h *harness) startEmitter(conn *application.StreamConn, sc Scenario) func() {
	stop := h.startEmitterCounted(conn, sc)
	return func() { stop() }
}

// startEmitterCounted runs the load and returns a stop function reporting how
// many frames were accepted by Send.
func (h *harness) startEmitterCounted(conn *application.StreamConn, sc Scenario) func() int64 {
	ctx, cancel := context.WithCancel(context.Background())
	var seq atomic.Int64
	var count atomic.Int64
	done := make(chan struct{})

	pad := make([]byte, max(0, sc.Size-frameHeaderBytes))
	for i := range pad {
		pad[i] = byte('a' + i%26)
	}

	go func() {
		defer close(done)

		var next time.Time
		var interval time.Duration
		if sc.Rate > 0 {
			interval = time.Duration(float64(time.Second) / float64(sc.Rate))
			next = time.Now()
		}

		local := make([]float64, 0, 4096)
		flush := time.NewTicker(250 * time.Millisecond)
		defer flush.Stop()

		for {
			select {
			case <-ctx.Done():
				h.mu.Lock()
				h.sendUS = append(h.sendUS, local...)
				h.mu.Unlock()
				return
			case <-flush.C:
				h.mu.Lock()
				h.sendUS = append(h.sendUS, local...)
				h.mu.Unlock()
				local = local[:0]
			default:
			}

			if sc.Rate > 0 {
				now := time.Now()
				if now.Before(next) {
					// Sleeping the whole gap oversleeps at high rates; cap the
					// nap so the pacing error stays small.
					nap := next.Sub(now)
					if nap > 2*time.Millisecond {
						nap = 2 * time.Millisecond
					}
					time.Sleep(nap)
					continue
				}
				next = next.Add(interval)
			}

			frame := make([]byte, frameHeaderBytes+len(pad))
			binary.BigEndian.PutUint64(frame[0:8], uint64(seq.Add(1)))
			binary.BigEndian.PutUint64(frame[8:16], math.Float64bits(h.nowMS()))
			copy(frame[frameHeaderBytes:], pad)

			// Blocking Send exercises the production backpressure path: it parks
			// until the frontend drains and wakes when space frees. The stop
			// function below refuses to wait forever for a parked producer.
			t0 := time.Now()
			err := conn.Send(frame)
			local = append(local, float64(time.Since(t0).Microseconds()))
			if err != nil {
				return
			}
			count.Add(1)
		}
	}()

	return func() int64 {
		cancel()
		// Bounded: a producer parked inside a blocking Send does not observe
		// ctx, so joining unconditionally hung the whole run. The count is
		// atomic, so reading it after abandoning the goroutine is safe.
		select {
		case <-done:
		case <-time.After(2 * time.Second):
		}
		return count.Load()
	}
}

func (h *harness) sample() sampleRow {
	row := sampleRow{TSec: time.Since(h.start).Seconds()}
	if samplerSupported {
		if v, ok := footprint(h.hostPid); ok {
			row.HostMB = float64(v) / (1024 * 1024)
		}
		var total uint64
		for _, p := range h.ourPids {
			if v, ok := footprint(p); ok {
				total += v
			}
		}
		row.WebMB = float64(total) / (1024 * 1024)
	}

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	row.HeapAllocMB = float64(ms.HeapAlloc) / (1024 * 1024)
	row.HeapSysMB = float64(ms.HeapSys) / (1024 * 1024)
	row.NumGC = ms.NumGC

	h.mu.Lock()
	row.Received = h.latestJS.Received
	row.Drops = h.latestJS.Drops
	row.Reorders = h.latestJS.Reorders
	h.mu.Unlock()
	return row
}

func (h *harness) discoverWebContent() error {
	pids, err := webContentPids()
	if err != nil {
		return err
	}
	for _, p := range pids {
		if !h.basePids[p] {
			h.ourPids = append(h.ourPids, p)
		}
	}
	if len(h.ourPids) == 0 {
		return fmt.Errorf("no new WebContent process appeared after launch")
	}
	return nil
}

func pctl(xs []float64, p float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	s := append([]float64(nil), xs...)
	sort.Float64s(s)
	idx := int(math.Ceil(p/100*float64(len(s)))) - 1
	if idx < 0 {
		idx = 0
	}
	if idx >= len(s) {
		idx = len(s) - 1
	}
	return s[idx]
}
