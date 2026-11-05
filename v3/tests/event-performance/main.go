// Command eventperf is the Phase -1 measurement harness for Wails v3 Go→JS
// event dispatch.
//
// It exists to answer one question with numbers instead of argument: does
// high-frequency event dispatch leak memory in the WebKit content process,
// and if so is the growth driven by call COUNT or by payload BYTES?
//
// It drives events straight at WebviewWindow.DispatchWailsEvent — the exact
// entry point that performs the evaluateJavaScript call — so the measurement
// is of the transport, not of the runtime's listener plumbing or the
// application-level fanout.
//
// Build/run (macOS):
//
//	go run ./tests/event-performance -duration 30s
package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed assets/index.html
var indexHTML []byte

var (
	flagDuration = flag.Duration("duration", 30*time.Second, "measurement window per scenario")
	flagSettle   = flag.Duration("settle", 3*time.Second, "settle period before each scenario (discarded)")
	flagOut      = flag.String("out", "", "output directory (default: ./eventperf-results)")
	flagOnly     = flag.String("only", "", "comma-separated scenario names to run (default: all)")
	flagSample   = flag.Duration("sample", 250*time.Millisecond, "footprint sample interval")
)

// perfPayload is the event body. Pad controls payload size; Seq and TMS let the
// JS side detect drops/reorders and compute delivery latency.
type perfPayload struct {
	Seq int64   `json:"seq"`
	TMS float64 `json:"tms"` // Go monotonic ms since process start
	Pad string  `json:"pad,omitempty"`
}

// ctlPayload resets the JS-side counters at a scenario boundary.
type ctlPayload struct {
	Ctl   string `json:"ctl"`
	Epoch int    `json:"epoch"`
}

// jsReport is what the page POSTs to /harness/report every 250ms.
// Latency values are raw deltas (jsPerformanceNow - goTMS); the clock offset
// between the two domains is removed later by subtracting DeltaMin.
type jsReport struct {
	Epoch      int     `json:"epoch"`
	Received   int64   `json:"received"`
	Drops      int64   `json:"drops"`
	Reorders   int64   `json:"reorders"`
	Frames     int64   `json:"frames"`
	LongFrames int64   `json:"longFrames"`
	DeltaMin   float64 `json:"deltaMin"`
	P50        float64 `json:"p50"`
	P95        float64 `json:"p95"`
	P99        float64 `json:"p99"`
	Max        float64 `json:"max"`
}

type harness struct {
	start time.Time

	basePids map[int]bool // WebContent pids that existed before we launched
	ourPids  []int        // WebContent pids attributable to us
	hostPid  int

	mu        sync.Mutex
	latestJS  jsReport
	emitUS    []float64 // emit wall times for the current sample interval
	epoch     int
	readyOnce sync.Once

	results []*scenarioResult
}

func main() {
	flag.Parse()

	if !samplerSupported {
		fmt.Println("NOTE: memory sampling is unsupported on this platform;")
		fmt.Println("timing/ordering metrics will still be collected.")
	}

	h := &harness{
		start:   time.Now(),
		hostPid: os.Getpid(),
	}

	// Snapshot WebContent pids BEFORE anything WebKit-related exists, so the
	// differential later attributes only our own content process.
	h.basePids = map[int]bool{}
	if samplerSupported {
		pids, err := webContentPids()
		if err != nil {
			// A failed pre-launch enumeration would classify every pre-existing
			// content process as ours, so refuse to run rather than report
			// someone else's memory.
			log.Fatalf("pre-launch process enumeration failed: %v", err)
		}
		for _, p := range pids {
			h.basePids[p] = true
		}
		log.Printf("pre-launch content processes on this machine: %d", len(h.basePids))
	}

	app := application.New(application.Options{
		Name:        "eventperf",
		Description: "Wails v3 Go→JS event dispatch measurement harness",
		Assets: application.AssetOptions{
			Handler:        h.handler(),
			DisableLogging: true,
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	win := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:   "eventperf",
		Title:  "eventperf",
		Width:  520,
		Height: 360,
		URL:    "/",
	})

	win.OnWindowEvent(events.Common.WindowRuntimeReady, func(*application.WindowEvent) {
		h.readyOnce.Do(func() {
			go h.runAll(app, win)
		})
	})

	if err := app.Run(); err != nil {
		log.Fatalf("app.Run: %v", err)
	}
}

// handler serves the harness page and receives JS-side reports.
func (h *harness) handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/harness/report", func(w http.ResponseWriter, r *http.Request) {
		var rep jsReport
		if err := json.NewDecoder(r.Body).Decode(&rep); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		h.mu.Lock()
		// Ignore reports from a previous scenario's epoch.
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

// nowMS is the Go-side monotonic clock in milliseconds, stamped into every event.
func (h *harness) nowMS() float64 {
	return float64(time.Since(h.start).Nanoseconds()) / 1e6
}

// runAll discovers the content process, runs every scenario, writes results.
func (h *harness) runAll(app *application.App, win *application.WebviewWindow) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("harness panic: %v", r)
		}
		app.Quit()
	}()

	// Give WebKit a moment to finish bringing up the content process.
	time.Sleep(2 * time.Second)

	if samplerSupported {
		if err := h.discoverWebContent(); err != nil {
			// Per the plan: fail loudly rather than emit a plausible-looking CSV
			// built from some other application's WebContent processes.
			log.Printf("FATAL: %v", err)
			fmt.Println()
			fmt.Println("Could not attribute a WebContent process to this app.")
			fmt.Println("WebKit may have adopted a prewarmed process that predates our snapshot.")
			fmt.Println("Aborting rather than reporting another app's memory.")
			return
		}
		log.Printf("our WebContent pid(s): %v  (host pid %d)", h.ourPids, h.hostPid)
	}

	outDir := *flagOut
	if outDir == "" {
		outDir = "eventperf-results"
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
		res := h.runScenario(win, sc)
		h.results = append(h.results, res)

		if err := writeScenarioCSV(outDir, res); err != nil {
			log.Printf("  csv: %v", err)
		}
		log.Printf("  → %s", res.oneLine())

		if res.Crashed {
			log.Printf("  web process died — stopping run (later scenarios would measure a dead webview)")
			break
		}
	}

	if err := writeSummary(outDir, h.results); err != nil {
		log.Printf("summary.json: %v", err)
	}
	if err := writeReport(outDir, h.results); err != nil {
		log.Printf("REPORT.md: %v", err)
	}
	fmt.Printf("\nResults written to %s\n", abs)
	fmt.Println(renderConsoleTable(h.results))
}

// runScenario executes one scenario: reset JS state, settle, then measure.
func (h *harness) runScenario(win *application.WebviewWindow, sc Scenario) *scenarioResult {
	dur := *flagDuration
	if sc.Duration > 0 {
		dur = sc.Duration
	}

	// New epoch: clears JS counters so this scenario measures only itself.
	h.mu.Lock()
	h.epoch++
	epoch := h.epoch
	h.latestJS = jsReport{Epoch: epoch}
	h.emitUS = nil
	h.mu.Unlock()

	win.DispatchWailsEvent(&application.CustomEvent{
		Name: "perfctl",
		Data: ctlPayload{Ctl: "reset", Epoch: epoch},
	})

	// Settle: run the load but discard measurements, so we measure steady state.
	stopSettle := h.startEmitter(win, sc)
	time.Sleep(*flagSettle)
	stopSettle()
	sc.resetCounter() // settle traffic must not count toward the measured totals

	// Reset counters again after settle so the measurement window is clean.
	h.mu.Lock()
	h.epoch++
	epoch = h.epoch
	h.latestJS = jsReport{Epoch: epoch}
	h.emitUS = nil
	h.mu.Unlock()
	win.DispatchWailsEvent(&application.CustomEvent{
		Name: "perfctl",
		Data: ctlPayload{Ctl: "reset", Epoch: epoch},
	})
	time.Sleep(300 * time.Millisecond) // let the reset land before sampling

	res := &scenarioResult{Scenario: sc, Duration: dur}
	res.StartFootprintWeb, _ = h.webFootprint()
	res.StartFootprintHost, _ = footprint(h.hostPid)

	stop := h.startEmitter(win, sc)
	deadline := time.Now().Add(dur)
	ticker := time.NewTicker(*flagSample)
	defer ticker.Stop()
	t0 := time.Now()

	for time.Now().Before(deadline) {
		<-ticker.C

		webFP, alive := h.webFootprint()
		hostFP, _ := footprint(h.hostPid)

		h.mu.Lock()
		js := h.latestJS
		emits := h.emitUS
		h.emitUS = nil
		h.mu.Unlock()

		p50, p99 := percentile(emits, 50), percentile(emits, 99)

		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)

		res.Samples = append(res.Samples, sample{
			TSms:      time.Since(t0).Milliseconds(),
			HostFP:    hostFP,
			WebFP:     webFP,
			Sent:      sc.sentSoFar(),
			EvalCalls: sc.sentSoFar(), // 1 eval per event on master; recorded so batching factor is visible later
			Received:  js.Received,
			EmitP50us: p50,
			EmitP99us: p99,

			GoHeapAlloc:    ms.HeapAlloc,
			GoHeapSys:      ms.HeapSys,
			GoHeapIdle:     ms.HeapIdle,
			GoHeapReleased: ms.HeapReleased,
			GoNumGC:        ms.NumGC,
		})

		if !alive {
			res.Crashed = true
			break
		}
	}
	stop()

	// Let the last in-flight events land before reading final JS state.
	time.Sleep(500 * time.Millisecond)

	h.mu.Lock()
	res.JS = h.latestJS
	h.mu.Unlock()

	res.Sent = sc.sentSoFar()
	res.EndFootprintWeb, _ = h.webFootprint()
	res.EndFootprintHost, _ = footprint(h.hostPid)
	res.finalise()
	sc.resetCounter()
	return res
}

// startEmitter drives events at the scenario's rate until the returned stop is called.
// It paces in 2ms groups because a per-event ticker cannot hold 5000 ev/s.
func (h *harness) startEmitter(win *application.WebviewWindow, sc Scenario) func() {
	if sc.Rate <= 0 {
		return func() {} // idle scenario: no events at all
	}

	done := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)

	pad := strings.Repeat("x", sc.PayloadBytes)
	altPad := pad
	if sc.AltPayloadBytes > 0 {
		altPad = strings.Repeat("x", sc.AltPayloadBytes)
	}

	// Second emitter running on the main thread, concurrent with the goroutine
	// emitter below. Both take sequence numbers from the same counter, so if
	// the inline main-thread path lets a later event overtake an earlier queued
	// one, the JS side records a reorder.
	if sc.MixedSource {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ticker := time.NewTicker(2 * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					// The sequence number must be taken inside the callback,
					// at the moment DispatchWailsEvent is actually called.
					// Taking it out here and deferring the dispatch would make
					// seq order differ from call order by construction, and
					// every reorder the JS side counted would be this loop's
					// own race rather than anything the framework did.
					// Acquired here, on the ticker goroutine, and released by
					// the UI callback below. Locking inside the callback
					// instead would put the UI thread behind a lock that a
					// goroutine can be holding while it waits for queue
					// capacity — and since the UI thread is the drainer, that
					// wait could never be relieved. Handing the lock over means
					// only this goroutine ever waits for it.
					//
					// sync.Mutex permits unlocking from a different goroutine.
					sc.emitMu.Lock()
					application.InvokeAsync(func() {
						defer sc.emitMu.Unlock()
						seq := sc.nextSeq()
						ev := &application.CustomEvent{
							Name: "perf",
							Data: perfPayload{Seq: seq, TMS: h.nowMS(), Pad: pad},
						}
						win.DispatchWailsEvent(ev)
					})
				}
			}
		}()
	}

	go func() {
		defer wg.Done()

		groupInterval := 2 * time.Millisecond
		perGroup := float64(sc.Rate) * groupInterval.Seconds()
		if sc.Burst > 0 {
			// Deliver the same mean rate but in bursts of sc.Burst.
			groupInterval = time.Duration(float64(sc.Burst) / float64(sc.Rate) * float64(time.Second))
			perGroup = float64(sc.Burst)
		}

		ticker := time.NewTicker(groupInterval)
		defer ticker.Stop()
		var credit float64

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				credit += perGroup
				n := int(credit)
				credit -= float64(n)
				for i := 0; i < n; i++ {
					if sc.MixedSource {
						sc.emitMu.Lock()
					}
					seq := sc.nextSeq()
					body := pad
					if seq%2 == 1 {
						body = altPad
					}
					ev := &application.CustomEvent{
						Name: "perf",
						Data: perfPayload{Seq: seq, TMS: h.nowMS(), Pad: body},
					}
					t := time.Now()
					func() {
						// A panic in dispatch must not leave emitMu held, or
						// every later emit in this scenario would wedge.
						if sc.MixedSource {
							defer sc.emitMu.Unlock()
						}
						win.DispatchWailsEvent(ev)
					}()
					us := float64(time.Since(t).Nanoseconds()) / 1e3

					h.mu.Lock()
					h.emitUS = append(h.emitUS, us)
					h.mu.Unlock()

					select {
					case <-done:
						return
					default:
					}
				}
			}
		}
	}()

	return func() { close(done); wg.Wait() }
}

// webFootprint sums our WebContent processes. alive=false means the process
// set changed or a read failed — i.e. the content process was replaced.
func (h *harness) webFootprint() (uint64, bool) {
	if !samplerSupported || len(h.ourPids) == 0 {
		return 0, true
	}
	var total uint64
	for _, pid := range h.ourPids {
		fp, ok := footprint(pid)
		if !ok {
			return total, false
		}
		total += fp
	}
	return total, true
}

// discoverWebContent diffs the current WebContent pid set against the
// pre-launch snapshot. An empty diff is a hard error, never a silent fallback.
func (h *harness) discoverWebContent() error {
	pids, err := webContentPids()
	if err != nil {
		return fmt.Errorf("process enumeration failed: %w", err)
	}
	var ours []int
	for _, p := range pids {
		if !h.basePids[p] {
			ours = append(ours, p)
		}
	}
	if len(ours) == 0 {
		return fmt.Errorf("no new com.apple.WebKit.WebContent process appeared after window creation "+
			"(%d existed before launch); cannot attribute one to this app", len(h.basePids))
	}
	h.ourPids = ours
	return nil
}
