package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

type sample struct {
	TSms      int64
	HostFP    uint64
	WebFP     uint64
	Sent      int64
	EvalCalls int64
	Received  int64
	EmitP50us float64
	EmitP99us float64

	// Go runtime memory, to separate Go-heap growth from native/cgo growth
	// when the host process footprint moves.
	GoHeapAlloc    uint64
	GoHeapSys      uint64
	GoHeapIdle     uint64
	GoHeapReleased uint64
	GoNumGC        uint32
}

type scenarioResult struct {
	Name     string        `json:"name"`
	Note     string        `json:"note"`
	Rate     int           `json:"rate_per_sec"`
	Payload  int           `json:"payload_bytes"`
	Burst    int           `json:"burst"`
	Duration time.Duration `json:"-"`

	Scenario Scenario `json:"-"`
	Samples  []sample `json:"-"`
	JS       jsReport `json:"js"`

	Sent    int64 `json:"events_sent"`
	Crashed bool  `json:"web_process_died"`

	StartFootprintWeb  uint64 `json:"start_footprint_web"`
	EndFootprintWeb    uint64 `json:"end_footprint_web"`
	StartFootprintHost uint64 `json:"start_footprint_host"`
	EndFootprintHost   uint64 `json:"end_footprint_host"`

	WebSlopeBps   float64 `json:"web_slope_bytes_per_sec"`
	HostSlopeBps  float64 `json:"host_slope_bytes_per_sec"`
	BytesPerEvent float64 `json:"web_bytes_per_event"`

	// Footprint under GC is a sawtooth, so a least-squares slope over a short
	// window mostly measures which phase of the cycle the window happened to
	// start and end in. The floor is the robust signal: a process that is
	// leaking cannot return to its earlier minimum.
	NetDeltaBytes   int64   `json:"net_delta_bytes"`
	SwingBytes      int64   `json:"swing_bytes"`
	FloorRiseBytes  int64   `json:"floor_rise_bytes"`
	FloorRiseBps    float64 `json:"floor_rise_bytes_per_sec"`
	FloorRisePerEvt float64 `json:"floor_rise_bytes_per_event"`

	// The host process matters as much as WebContent: pending eval payloads
	// queue natively on the host side with no backpressure, so a byte-heavy
	// workload grows the HOST, not the web process. Grading only WebContent
	// misses it completely.
	HostPeakBytes      int64   `json:"host_peak_bytes"`
	HostRiseBytes      int64   `json:"host_rise_bytes"`
	ByteRatePerSec     float64 `json:"emitted_byte_rate_per_sec"`
	GoHeapSysPeakBytes int64   `json:"go_heap_sys_peak_bytes"`
	AchievedRate       float64 `json:"achieved_rate_per_sec"`
	EmitP50us          float64 `json:"emit_p50_us"`
	EmitP99us          float64 `json:"emit_p99_us"`
	LatP50ms           float64 `json:"delivery_p50_ms"`
	LatP99ms           float64 `json:"delivery_p99_ms"`
	Valid              bool    `json:"valid"`
	Verdict            string  `json:"verdict"`
}

func (r *scenarioResult) finalise() {
	sc := r.Scenario
	r.Name, r.Note, r.Rate, r.Payload, r.Burst = sc.Name, sc.Note, sc.Rate, sc.PayloadBytes, sc.Burst

	secs := r.Duration.Seconds()
	if secs > 0 {
		r.AchievedRate = float64(r.Sent) / secs
	}

	xs := make([]float64, 0, len(r.Samples))
	web := make([]float64, 0, len(r.Samples))
	host := make([]float64, 0, len(r.Samples))
	var emit50, emit99 []float64
	for _, s := range r.Samples {
		xs = append(xs, float64(s.TSms)/1000)
		web = append(web, float64(s.WebFP))
		host = append(host, float64(s.HostFP))
		if s.EmitP50us > 0 {
			emit50 = append(emit50, s.EmitP50us)
		}
		if s.EmitP99us > 0 {
			emit99 = append(emit99, s.EmitP99us)
		}
	}
	r.WebSlopeBps = slope(xs, web)
	r.HostSlopeBps = slope(xs, host)
	r.EmitP50us = percentile(emit50, 50)
	r.EmitP99us = percentile(emit99, 99)

	if r.Sent > 0 {
		r.BytesPerEvent = r.WebSlopeBps * secs / float64(r.Sent)
	}

	// Floor analysis: compare the minimum footprint of the first half of the
	// window against the minimum of the second half. Sawtooth returns to the
	// same floor; a leak raises it.
	if n := len(web); n >= 4 {
		half := n / 2
		firstMin, secondMin := web[0], web[half]
		for _, v := range web[:half] {
			if v < firstMin {
				firstMin = v
			}
		}
		for _, v := range web[half:] {
			if v < secondMin {
				secondMin = v
			}
		}
		lo, hi := web[0], web[0]
		for _, v := range web {
			if v < lo {
				lo = v
			}
			if v > hi {
				hi = v
			}
		}
		r.SwingBytes = int64(hi - lo)
		r.NetDeltaBytes = int64(web[n-1] - web[0])
		r.FloorRiseBytes = int64(secondMin - firstMin)
		if secs > 0 {
			r.FloorRiseBps = float64(r.FloorRiseBytes) / (secs / 2)
		}
		if r.Sent > 0 {
			r.FloorRisePerEvt = float64(r.FloorRiseBytes) / float64(r.Sent)
		}
	}

	// Host-side native growth, and the Go heap alongside it so the two can be
	// told apart: flat Go heap + rising host footprint means native buffering.
	if len(host) > 0 {
		hlo, hhi := host[0], host[0]
		for _, v := range host {
			if v < hlo {
				hlo = v
			}
			if v > hhi {
				hhi = v
			}
		}
		r.HostPeakBytes = int64(hhi)
		r.HostRiseBytes = int64(hhi - hlo)
	}
	for _, s := range r.Samples {
		if int64(s.GoHeapSys) > r.GoHeapSysPeakBytes {
			r.GoHeapSysPeakBytes = int64(s.GoHeapSys)
		}
	}
	if secs > 0 {
		// interleave alternates two payload sizes, so a single-size figure
		// would understate its byte rate by ~32x and mislead the attribution
		// column, which is the main signal in REPORT.md.
		avg := float64(r.Payload)
		if alt := sc.AltPayloadBytes; alt > 0 {
			avg = (float64(r.Payload) + float64(alt)) / 2
		}
		r.ByteRatePerSec = float64(r.Sent) * avg / secs
	}

	// Clock domains differ but both are monotonic, so a constant offset cancels:
	// subtract the smallest observed transit from the percentiles.
	if r.JS.Received > 0 {
		r.LatP50ms = r.JS.P50 - r.JS.DeltaMin
		r.LatP99ms = r.JS.P99 - r.JS.DeltaMin
	}

	// A scenario that sent events but received none measured nothing. Reporting
	// that as "flat" would be a false negative.
	r.Valid = r.Rate == 0 || r.JS.Received > 0
}

func (r *scenarioResult) oneLine() string {
	switch {
	case r.Crashed:
		return fmt.Sprintf("CRASH after %d events (web process died)", r.Sent)
	case !r.Valid:
		return fmt.Sprintf("INVALID — sent %d, received 0 (runtime not mounted?)", r.Sent)
	}
	return fmt.Sprintf("sent=%d recv=%d achieved=%.0f/s floor=%s swing=%s emit_p50=%.0fµs emit_p99=%.0fµs drops=%d",
		r.Sent, r.JS.Received, r.AchievedRate, signedBytes(r.FloorRiseBytes),
		humanBytes(r.SwingBytes), r.EmitP50us, r.EmitP99us, r.JS.Drops)
}

// applyVerdict grades a scenario against the idle baseline slope.
func (r *scenarioResult) applyVerdict(baselineBps float64) {
	switch {
	case r.Crashed:
		r.Verdict = "CRASH — web process terminated (reproduces #215729's terminal symptom)"
		return
	case !r.Valid:
		r.Verdict = "INVALID — zero events received; do not read as flat"
		return
	}
	// Graded on floor rise, not slope. Requires BOTH a rate above the idle
	// baseline AND an absolute magnitude, so a short window full of sawtooth
	// cannot produce a leak verdict on its own.
	const (
		leakBps      = 50 * 1024 // 50 KB/s sustained floor rise
		leakAbsBytes = 4 << 20   // and at least 4 MB of it
		elevBps      = 25 * 1024
		elevAbsBytes = 1 << 20
	)
	rise, riseBps := r.FloorRiseBytes, r.FloorRiseBps
	ctx := fmt.Sprintf("floor %s over %ds (swing %s, net %s)",
		signedBytes(rise), int(r.Duration.Seconds()), humanBytes(r.SwingBytes), signedBytes(r.NetDeltaBytes))

	switch {
	case riseBps > baselineBps+leakBps && rise > leakAbsBytes:
		r.Verdict = "WEBKIT LEAK — " + ctx
	case riseBps > baselineBps+elevBps && rise > elevAbsBytes:
		r.Verdict = "WEBKIT ELEVATED — " + ctx
	default:
		r.Verdict = "webkit flat — " + ctx
	}

	// Host-side (UI process) retention inside evaluateJavaScript, proportional
	// to JS source bytes. Attributed by experiment: with the eval call removed
	// but the CString malloc/free and NSString construction kept, the same 1 MB
	// workload holds the host flat at ~50 MB. Go heap stays flat throughout, so
	// this is native memory, not Go allocation.
	if r.HostRiseBytes > 64<<20 && r.GoHeapSysPeakBytes < r.HostRiseBytes/4 {
		r.Verdict += fmt.Sprintf("  ||  HOST RETENTION — host +%s (peak %s, Go heap only %s) at %s/s of JS source",
			humanBytes(r.HostRiseBytes), humanBytes(r.HostPeakBytes),
			humanBytes(r.GoHeapSysPeakBytes), humanBytes(int64(r.ByteRatePerSec)))
	}
}

func writeScenarioCSV(dir string, r *scenarioResult) error {
	var b strings.Builder
	b.WriteString("timestamp_ms,footprint_web_bytes,footprint_host_bytes,events_sent,events_received,eval_calls,emit_p50_us,emit_p99_us," +
		"go_heap_alloc,go_heap_sys,go_heap_idle,go_heap_released,go_num_gc\n")
	for _, s := range r.Samples {
		fmt.Fprintf(&b, "%d,%d,%d,%d,%d,%d,%.1f,%.1f,%d,%d,%d,%d,%d\n",
			s.TSms, s.WebFP, s.HostFP, s.Sent, s.Received, s.EvalCalls, s.EmitP50us, s.EmitP99us,
			s.GoHeapAlloc, s.GoHeapSys, s.GoHeapIdle, s.GoHeapReleased, s.GoNumGC)
	}
	return os.WriteFile(filepath.Join(dir, r.Scenario.Name+".csv"), []byte(b.String()), 0o644)
}

type summary struct {
	GeneratedAt string            `json:"generated_at"`
	Environment map[string]string `json:"environment"`
	Scenarios   []*scenarioResult `json:"scenarios"`
}

func writeSummary(dir string, results []*scenarioResult) error {
	grade(results)
	s := summary{
		GeneratedAt: time.Now().Format(time.RFC3339),
		Environment: environment(),
		Scenarios:   results,
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "summary.json"), b, 0o644)
}

func grade(results []*scenarioResult) {
	var baseline float64
	for _, r := range results {
		if r.Scenario.Name == "idle" {
			baseline = r.FloorRiseBps
		}
	}
	for _, r := range results {
		r.applyVerdict(baseline)
	}
}

func writeReport(dir string, results []*scenarioResult) error {
	grade(results)
	env := environment()

	var b strings.Builder
	b.WriteString("# Wails v3 Go→JS event dispatch — Phase -1 results\n\n")
	b.WriteString("Measured at the `DispatchWailsEvent` → `evaluateJavaScript` boundary.\n")
	b.WriteString("`footprint` is `ri_phys_footprint` of the WebKit content process — what\n")
	b.WriteString("Activity Monitor shows and what jetsam acts on.\n\n")

	b.WriteString("## Environment\n\n")
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(&b, "- **%s:** %s\n", k, env[k])
	}

	b.WriteString("\n## Results\n\n")
	b.WriteString("| scenario | rate | payload | byte rate | sent | recv | web floor | web swing | **host peak** | **host rise** | go heap | emit p50 | emit p99 | drops |\n")
	b.WriteString("|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|\n")
	for _, r := range results {
		fmt.Fprintf(&b, "| `%s` | %d/s | %s | %s/s | %d | %d | %s | %s | **%s** | **%s** | %s | %.0f µs | %.0f µs | %d |\n",
			r.Scenario.Name, r.Rate, humanBytes(int64(r.Payload)), humanBytes(int64(r.ByteRatePerSec)),
			r.Sent, r.JS.Received, signedBytes(r.FloorRiseBytes), humanBytes(r.SwingBytes),
			humanBytes(r.HostPeakBytes), signedBytes(r.HostRiseBytes), humanBytes(r.GoHeapSysPeakBytes),
			r.EmitP50us, r.EmitP99us, r.JS.Drops)
	}

	b.WriteString("\n## Verdicts\n\n")
	for _, r := range results {
		fmt.Fprintf(&b, "- **`%s`** — %s\n", r.Scenario.Name, r.Verdict)
	}

	b.WriteString("\n## Reading this\n\n")
	b.WriteString("- **`floor rise`** is the leak signal: the minimum footprint of the second half\n")
	b.WriteString("  of the window minus the minimum of the first half. Memory under GC is a\n")
	b.WriteString("  sawtooth, so `net Δ` and any regression slope mostly report which phase of\n")
	b.WriteString("  the cycle the window happened to start and end in. A process that is\n")
	b.WriteString("  genuinely leaking cannot return to its earlier floor.\n")
	b.WriteString("- `swing` is the peak-to-trough range: read `floor rise` against it. A floor\n")
	b.WriteString("  rise much smaller than the swing is noise.\n")
	b.WriteString("- Floor rise concentrated in the **rate sweep** implicates a per-call leak;\n")
	b.WriteString("  in the **size sweep**, payload retention (WebKit #215729).\n")
	b.WriteString("- An `INVALID` row received zero events and must not be read as \"flat\".\n")
	b.WriteString("\nAll scenarios shared one window; footprint is therefore cumulative across\n")
	b.WriteString("rows and only the per-scenario *slope* is comparable, not absolute values.\n")

	return os.WriteFile(filepath.Join(dir, "REPORT.md"), []byte(b.String()), 0o644)
}

func renderConsoleTable(results []*scenarioResult) string {
	grade(results)
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%-14s %9s %9s %10s %11s %10s %9s\n",
		"SCENARIO", "SENT", "RECV", "ACHIEVED", "FLOOR RISE", "SWING", "EMIT p99"))
	for _, r := range results {
		b.WriteString(fmt.Sprintf("%-14s %9d %9d %9.0f/s %11s %10s %8.0fµs  %s\n",
			r.Scenario.Name, r.Sent, r.JS.Received, r.AchievedRate,
			signedBytes(r.FloorRiseBytes), humanBytes(r.SwingBytes), r.EmitP99us, r.Verdict))
	}
	return b.String()
}

func environment() map[string]string {
	env := map[string]string{
		"go":       runtime.Version(),
		"platform": runtime.GOOS + "/" + runtime.GOARCH,
	}
	if out, err := exec.Command("sw_vers", "-productVersion").Output(); err == nil {
		env["macos"] = strings.TrimSpace(string(out))
	}
	if out, err := exec.Command("sysctl", "-n", "machdep.cpu.brand_string").Output(); err == nil {
		env["cpu"] = strings.TrimSpace(string(out))
	}
	if out, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output(); err == nil {
		env["wails_commit"] = strings.TrimSpace(string(out))
	}
	if out, err := exec.Command("defaults", "read",
		"/System/Library/Frameworks/WebKit.framework/Resources/Info.plist", "CFBundleVersion").Output(); err == nil {
		env["webkit"] = strings.TrimSpace(string(out))
	}
	return env
}

// slope is a least-squares fit of y over x, in y-units per x-unit.
func slope(xs, ys []float64) float64 {
	n := float64(len(xs))
	if n < 2 {
		return 0
	}
	var sx, sy, sxy, sxx float64
	for i := range xs {
		sx += xs[i]
		sy += ys[i]
		sxy += xs[i] * ys[i]
		sxx += xs[i] * xs[i]
	}
	den := n*sxx - sx*sx
	if den == 0 {
		return 0
	}
	return (n*sxy - sx*sy) / den
}

func percentile(v []float64, p float64) float64 {
	if len(v) == 0 {
		return 0
	}
	s := append([]float64(nil), v...)
	sort.Float64s(s)
	idx := int(p / 100 * float64(len(s)-1))
	if idx < 0 {
		idx = 0
	}
	if idx >= len(s) {
		idx = len(s) - 1
	}
	return s[idx]
}

// signedBytes always shows a sign, so a falling floor is unmistakable.
func signedBytes(b int64) string {
	if b >= 0 {
		return "+" + humanBytes(b)
	}
	return humanBytes(b)
}

func humanBytes(b int64) string {
	neg := b < 0
	if neg {
		b = -b
	}
	var s string
	switch {
	case b >= 1<<30:
		s = fmt.Sprintf("%.2f GB", float64(b)/(1<<30))
	case b >= 1<<20:
		s = fmt.Sprintf("%.2f MB", float64(b)/(1<<20))
	case b >= 1<<10:
		s = fmt.Sprintf("%.1f KB", float64(b)/(1<<10))
	default:
		s = fmt.Sprintf("%d B", b)
	}
	if neg {
		return "-" + s
	}
	return s
}
