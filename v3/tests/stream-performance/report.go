package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type sampleRow struct {
	TSec   float64 `json:"t"`
	HostMB float64 `json:"hostMB"`
	WebMB  float64 `json:"webMB"`

	// Go heap alongside host footprint, so any host growth can be attributed
	// to our own allocation rather than to the transport. This is what settled
	// the equivalent question for events: the Go heap was flat at 24 MB while
	// the host climbed to 11.5 GB, which located the retention inside WebKit.
	HeapAllocMB float64 `json:"heapAllocMB"`
	HeapSysMB   float64 `json:"heapSysMB"`
	NumGC       uint32  `json:"numGC"`

	Received int64 `json:"received"`
	Drops    int64 `json:"drops"`
	Reorders int64 `json:"reorders"`
}

type scenarioResult struct {
	Name        string        `json:"name"`
	PayloadSize int           `json:"payloadBytes"`
	TargetRate  int           `json:"targetRate"`
	Duration    time.Duration `json:"duration"`
	Epoch       int           `json:"epoch"`

	Sent    int64 `json:"sent"`
	Echoes  int64 `json:"echoes"`
	SendP50 float64
	SendP99 float64

	JS      jsReport    `json:"js"`
	Samples []sampleRow `json:"samples"`

	// Derived
	AchievedRate float64 `json:"achievedRate"`
	MBPerSec     float64 `json:"mbPerSec"`
	HostFirstMB  float64 `json:"hostFirstMB"`
	HostPeakMB   float64 `json:"hostPeakMB"`
	HostLastMB   float64 `json:"hostLastMB"`
	HostRiseMB   float64 `json:"hostRiseMB"`
	WebFirstMB   float64 `json:"webFirstMB"`
	WebPeakMB    float64 `json:"webPeakMB"`
	WebLastMB    float64 `json:"webLastMB"`
	WebRiseMB    float64 `json:"webRiseMB"`
	FramesPerReq float64 `json:"framesPerRequest"`
}

func (r *scenarioResult) finalise() {
	secs := r.Duration.Seconds()
	if secs > 0 {
		r.AchievedRate = float64(r.Sent) / secs
		r.MBPerSec = float64(r.Sent) * float64(r.PayloadSize) / secs / (1024 * 1024)
	}
	if r.JS.Polls > 0 {
		r.FramesPerReq = float64(r.JS.Received) / float64(r.JS.Polls)
	}

	if len(r.Samples) == 0 {
		return
	}
	r.HostFirstMB = r.Samples[0].HostMB
	r.WebFirstMB = r.Samples[0].WebMB
	r.HostLastMB = r.Samples[len(r.Samples)-1].HostMB
	r.WebLastMB = r.Samples[len(r.Samples)-1].WebMB
	for _, s := range r.Samples {
		if s.HostMB > r.HostPeakMB {
			r.HostPeakMB = s.HostMB
		}
		if s.WebMB > r.WebPeakMB {
			r.WebPeakMB = s.WebMB
		}
	}
	// Rise is measured on the floor rather than the peak, deliberately: a
	// least-squares slope over a GC sawtooth mostly reports which phase the
	// window started in, and produced false leak labels in the event work.
	r.HostRiseMB = r.HostLastMB - r.HostFirstMB
	r.WebRiseMB = r.WebLastMB - r.WebFirstMB
}

func (r *scenarioResult) oneLine() string {
	lat := r.JS.P99 - r.JS.DeltaMin
	if lat < 0 {
		lat = 0
	}
	return fmt.Sprintf(
		"sent %d recv %d drops %d reorders %d | %.0f f/s %.1f MB/s | batch avg %.1f max %d | send p50 %.0fus p99 %.0fus | lat p99 %.1fms | host %.0f→%.0fMB web %.0f→%.0fMB",
		r.Sent, r.JS.Received, r.JS.Drops, r.JS.Reorders,
		r.AchievedRate, r.MBPerSec,
		r.FramesPerReq, r.JS.MaxBatch,
		r.SendP50, r.SendP99,
		lat,
		r.HostFirstMB, r.HostLastMB, r.WebFirstMB, r.WebLastMB)
}

func writeScenarioCSV(dir string, r *scenarioResult) error {
	f, err := os.Create(filepath.Join(dir, r.Name+".csv"))
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{"t_sec", "host_mb", "web_mb", "heap_alloc_mb", "heap_sys_mb", "num_gc", "received", "drops", "reorders"}); err != nil {
		return err
	}
	for _, s := range r.Samples {
		if err := w.Write([]string{
			strconv.FormatFloat(s.TSec, 'f', 2, 64),
			strconv.FormatFloat(s.HostMB, 'f', 2, 64),
			strconv.FormatFloat(s.WebMB, 'f', 2, 64),
			strconv.FormatFloat(s.HeapAllocMB, 'f', 2, 64),
			strconv.FormatFloat(s.HeapSysMB, 'f', 2, 64),
			strconv.FormatUint(uint64(s.NumGC), 10),
			strconv.FormatInt(s.Received, 10),
			strconv.FormatInt(s.Drops, 10),
			strconv.FormatInt(s.Reorders, 10),
		}); err != nil {
			return err
		}
	}
	return nil
}

func writeSummary(dir string, results []*scenarioResult) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "summary.json"), data, 0o644)
}

func renderConsoleTable(results []*scenarioResult) string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("%-18s %10s %10s %7s %9s %8s %8s %9s %9s %9s\n",
		"scenario", "sent", "recv", "drops", "reorders", "f/s", "MB/s", "batch", "sendp99", "host→MB"))
	b.WriteString(strings.Repeat("-", 110) + "\n")
	for _, r := range results {
		b.WriteString(fmt.Sprintf("%-18s %10d %10d %7d %9d %8.0f %8.1f %9.1f %9.0f %9.0f\n",
			r.Name, r.Sent, r.JS.Received, r.JS.Drops, r.JS.Reorders,
			r.AchievedRate, r.MBPerSec, r.FramesPerReq, r.SendP99, r.HostLastMB))
	}
	return b.String()
}
