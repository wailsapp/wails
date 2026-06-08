package tui

import (
	"io"
	"strings"
	"sync"

	"github.com/atterpac/refresh/engine"
)

// maxProcLines is the per-process log ring-buffer cap.
const maxProcLines = 1000

// ProcStore is the goroutine-safe bridge between the refresh engine's SDK taps
// (Output / OnProcessEvent, called from engine-owned goroutines) and the TUI's
// processes view. The engine writes into it; the view reads from it on the UI
// thread. The OnLog / OnEvent callbacks let the TUI marshal a redraw onto the
// draw thread (they are invoked from engine goroutines, so the TUI wraps them
// in theme.QueueUpdateDraw).
type ProcStore struct {
	mu        sync.Mutex
	lines     map[string][]string // process name -> capped log lines
	partial   map[string]string   // process name -> trailing partial line
	lastEvent map[string]engine.ProcessEvent

	// OnLog fires after new output is appended. OnEvent fires on a lifecycle
	// transition. Both may be nil; set by the TUI.
	OnLog   func()
	OnEvent func(engine.ProcessEvent)
}

// NewProcStore returns an empty store ready to receive engine taps.
func NewProcStore() *ProcStore {
	return &ProcStore{
		lines:     map[string][]string{},
		partial:   map[string]string{},
		lastEvent: map[string]engine.ProcessEvent{},
	}
}

// WriterFor implements engine.OutputFunc. It returns a non-nil writer for every
// process so that no output leaks to the terminal the TUI owns. Both streams of
// a process share one buffer (stderr is not treated as an error stream).
func (s *ProcStore) WriterFor(info engine.ProcessInfo, _ string) io.Writer {
	return &procWriter{store: s, name: info.Name}
}

// HandleEvent implements engine.EventFunc.
func (s *ProcStore) HandleEvent(ev engine.ProcessEvent) {
	s.mu.Lock()
	s.lastEvent[ev.Info.Name] = ev
	s.mu.Unlock()
	if s.OnEvent != nil {
		s.OnEvent(ev)
	}
}

// Lines returns a copy of the buffered log lines for a process.
func (s *ProcStore) Lines(name string) []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	src := s.lines[name]
	out := make([]string, len(src))
	copy(out, src)
	if p := s.partial[name]; p != "" {
		out = append(out, p)
	}
	return out
}

// append splits b into lines (carrying a trailing partial across writes) and
// stores them, capped at maxProcLines.
func (s *ProcStore) append(name string, b []byte) {
	s.mu.Lock()
	text := s.partial[name] + string(b)
	parts := strings.Split(text, "\n")
	// The last element is the (possibly empty) trailing partial.
	s.partial[name] = parts[len(parts)-1]
	complete := parts[:len(parts)-1]
	for _, ln := range complete {
		// stderr is NOT an error stream — dev tools (vite, esbuild, go) log
		// normal output there. Render it like stdout and escape any tag markup
		// so the process's own "[...]" text isn't parsed as color codes.
		ln = escapeTags(strings.TrimRight(ln, "\r"))
		s.lines[name] = append(s.lines[name], ln)
	}
	if over := len(s.lines[name]) - maxProcLines; over > 0 {
		s.lines[name] = s.lines[name][over:]
	}
	s.mu.Unlock()
	if len(complete) > 0 && s.OnLog != nil {
		s.OnLog()
	}
}

// procWriter is the per-process io.Writer handed to the engine.
type procWriter struct {
	store *ProcStore
	name  string
}

func (w *procWriter) Write(b []byte) (int, error) {
	w.store.append(w.name, b)
	return len(b), nil
}
