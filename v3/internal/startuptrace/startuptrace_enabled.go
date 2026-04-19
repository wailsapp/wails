//go:build wails_trace_startup

package startuptrace

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type event struct {
	name     string
	windowID uint
	ts       time.Time
}

var (
	mu     sync.Mutex
	start  = time.Now()
	events = make([]event, 0, 128)
)

// Mark records a named event at the current time.
func Mark(name string) {
	now := time.Now()
	mu.Lock()
	events = append(events, event{name: name, ts: now})
	mu.Unlock()
}

// MarkWindow records a named event scoped to a specific window.
func MarkWindow(windowID uint, name string) {
	now := time.Now()
	mu.Lock()
	events = append(events, event{name: name, windowID: windowID, ts: now})
	mu.Unlock()
}

// Flush writes a Chrome trace JSON document to the file named by the
// WAILS_TRACE_STARTUP_OUTPUT environment variable. Safe to call multiple
// times; only the first call with a non-empty env var writes output.
func Flush() {
	path := os.Getenv("WAILS_TRACE_STARTUP_OUTPUT")
	if path == "" {
		return
	}
	mu.Lock()
	snapshot := make([]event, len(events))
	copy(snapshot, events)
	mu.Unlock()

	type traceEntry struct {
		Name  string `json:"name"`
		Ph    string `json:"ph"`
		Pid   int    `json:"pid"`
		Tid   uint   `json:"tid"`
		Ts    int64  `json:"ts"`
		Scope string `json:"s,omitempty"`
	}

	out := make([]traceEntry, 0, len(snapshot))
	for _, e := range snapshot {
		out = append(out, traceEntry{
			Name:  e.name,
			Ph:    "i",
			Pid:   1,
			Tid:   e.windowID,
			Ts:    e.ts.Sub(start).Microseconds(),
			Scope: "g",
		})
	}

	data, err := json.Marshal(out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "startuptrace: marshal failed: %v\n", err)
		return
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "startuptrace: write %s failed: %v\n", path, err)
	}
}

// Enabled reports whether tracing is compiled in.
func Enabled() bool { return true }
