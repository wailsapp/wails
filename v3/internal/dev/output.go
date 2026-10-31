package dev

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/internal/report"
)

// Diagnostic retains the cause and launch/build identity for CLI diagnostics.
type Diagnostic struct {
	Component, Phase string
	Generation       uint64
	Command          []string
	ExitCode         int
	Output           string
	Err              error
}

func (d *Diagnostic) Error() string { return fmt.Sprintf("%s %s: %v", d.Component, d.Phase, d.Err) }
func (d *Diagnostic) Unwrap() error { return d.Err }
func exitCode(err error) int {
	var e *exec.ExitError
	if errors.As(err, &e) {
		return e.ExitCode()
	}
	return -1
}

// Event is a Wails dev lifecycle notification. Child output remains text so
// multiline compiler diagnostics and application stack traces stay intact.
type Event struct {
	Generation                uint64
	Component, Phase, Message string
	Diagnostic                *Diagnostic
}

// Output serialises builds, lifecycle events and child output on one sink.
// It deliberately uses stable lines: a long-lived application must never race
// a build spinner for ownership of terminal cursor positions.
type Output struct {
	mu       sync.Mutex
	writer   io.Writer
	level    report.Verbosity
	observer func(Event)
}

func NewOutput(w io.Writer, level report.Verbosity) *Output { return &Output{writer: w, level: level} }

// Observe must be configured before the session starts. The callback must not
// re-enter Output; delivery is ordered with terminal output.
func (o *Output) Observe(f func(Event)) { o.observer = f }
func (o *Output) Status(g uint64, c, p, m string) {
	o.emit(Event{Generation: g, Component: c, Phase: p, Message: m})
}
func (o *Output) Failure(g uint64, c, p string, err error, message string) {
	d := &Diagnostic{Component: c, Phase: p, Generation: g, ExitCode: exitCode(err), Err: err}
	var cause *Diagnostic
	if errors.As(err, &cause) {
		*d = *cause
		if g != 0 {
			d.Generation = g
		}
	}
	o.emit(Event{Generation: g, Component: c, Phase: p, Message: message, Diagnostic: d})
}
func (o *Output) emit(e Event) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.observer != nil {
		o.observer(e)
	}
	if e.Diagnostic != nil {
		if !report.IsReported(e.Diagnostic.Err) {
			fmt.Fprintf(o.writer, "[%s] %v\n", e.Component, e.Diagnostic)
			if o.level == report.Silent && e.Diagnostic.Output != "" {
				fmt.Fprintln(o.writer, e.Diagnostic.Output)
			}
		}
		if e.Message != "" {
			fmt.Fprintln(o.writer, e.Message)
		}
	} else if o.level != report.Silent {
		fmt.Fprintf(o.writer, "[%s] %s\n", e.Component, e.Message)
	}
}
func (o *Output) Writer(component, stream string) io.WriteCloser {
	return &lineWriter{out: o, label: component + "/" + stream}
}

// lineWriter bounds incomplete lines to 64 KiB, including output without any
// newline. Close flushes the final fragment after the process is reaped.
type lineWriter struct {
	mu      sync.Mutex
	out     *Output
	label   string
	pending []byte
}

func (w *lineWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	n := len(p)
	for len(p) > 0 {
		room := 65536 - len(w.pending)
		take := len(p)
		if take > room {
			take = room
		}
		w.pending = append(w.pending, p[:take]...)
		p = p[take:]
		for {
			i := strings.IndexByte(string(w.pending), '\n')
			if i < 0 {
				break
			}
			w.flush(w.pending[:i])
			w.pending = w.pending[i+1:]
		}
		if len(w.pending) == 65536 {
			w.flush(w.pending)
			w.pending = nil
		}
	}
	return n, nil
}
func (w *lineWriter) flush(p []byte) {
	w.out.mu.Lock()
	defer w.out.mu.Unlock()
	if w.out.level != report.Silent {
		fmt.Fprintf(w.out.writer, "[%s] %s\n", w.label, p)
	}
}
func (w *lineWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if len(w.pending) > 0 {
		w.flush(w.pending)
		w.pending = nil
	}
	return nil
}

type outputKey struct{}
type generationKey struct{}

func WithGeneration(ctx context.Context, g uint64) context.Context {
	return context.WithValue(ctx, generationKey{}, g)
}
func Generation(ctx context.Context) uint64 {
	if ctx == nil {
		return 0
	}
	g, _ := ctx.Value(generationKey{}).(uint64)
	return g
}

func WithOutput(ctx context.Context, o *Output) context.Context {
	return context.WithValue(ctx, outputKey{}, o)
}
func OutputFrom(ctx context.Context) *Output {
	if ctx == nil {
		return nil
	}
	o, _ := ctx.Value(outputKey{}).(*Output)
	return o
}
func (o *Output) BuildReporter(generation ...uint64) report.Reporter {
	var g uint64
	if len(generation) > 0 {
		g = generation[0]
	}
	return &buildReporter{output: o, generation: g, names: map[report.StepID]string{}}
}

type buildReporter struct {
	generation uint64
	mu         sync.Mutex
	output     *Output
	next       report.StepID
	names      map[report.StepID]string
}

func (r *buildReporter) Level() report.Verbosity { return r.output.level }
func (r *buildReporter) BuildStart(verb, target string, n int) {
	r.output.Status(r.generation, "build", "started", fmt.Sprintf("%s %s (%d steps)", verb, target, n))
}
func (r *buildReporter) StepStart(name, label string) report.StepID {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.next++
	r.names[r.next] = name
	return r.next
}
func (r *buildReporter) name(id report.StepID) string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.names[id]
}
func (r *buildReporter) StepInfo(id report.StepID, m string) {
	r.output.Status(r.generation, r.name(id), "info", m)
}
func (r *buildReporter) StepCommand(id report.StepID, c string) {
	if r.Level() >= report.Verbose {
		r.StepInfo(id, "$ "+c)
	}
}
func (r *buildReporter) StepOutput(id report.StepID, l string) {
	if r.Level() >= report.Verbose {
		r.StepInfo(id, l)
	}
}
func (r *buildReporter) StepEnd(id report.StepID, s report.Status, d time.Duration) {
	status := "built"
	switch s {
	case report.StatusCached:
		status = "cached"
	case report.StatusSkipped:
		status = "skipped"
	case report.StatusFailed:
		status = "failed"
	}
	r.StepInfo(id, fmt.Sprintf("%s (%s)", status, d.Round(time.Millisecond)))
}
func (r *buildReporter) StepFailed(id report.StepID, f report.Failure) {
	d := &Diagnostic{Generation: r.generation, Component: r.name(id), Phase: "build", Command: []string{f.Command}, ExitCode: f.ExitCode, Output: f.Output, Err: f.Err}
	r.output.mu.Lock()
	defer r.output.mu.Unlock()
	if r.output.observer != nil {
		r.output.observer(Event{Generation: r.generation, Component: d.Component, Phase: d.Phase, Diagnostic: d})
	}
	fmt.Fprintf(r.output.writer, "[%s] %v\n$ %s\n%s\n", d.Component, f.Err, f.Command, f.Output)
}
func (r *buildReporter) BuildEnd(d time.Duration, ok bool) {
	message := "Build complete"
	if !ok {
		message = "Build failed"
	}
	r.output.Status(r.generation, "build", "finished", fmt.Sprintf("%s in %s", message, d.Round(time.Millisecond)))
}
func (r *buildReporter) BuildCanceled(d time.Duration) {
	r.output.Status(r.generation, "build", "canceled", "Build superseded or session stopped")
}
func (r *buildReporter) Artifact(a report.Artifact) {
	if r.Level() >= report.Verbose {
		r.output.Status(r.generation, "build", "artifact", a.Path)
	}
}
func (r *buildReporter) Debug(d report.DebugLine) {
	if r.Level() >= report.Debug {
		r.output.Status(r.generation, "build", "debug", fmt.Sprint(d))
	}
}

// Tail retains bounded diagnostic context even when ordinary output is quiet.
type Tail struct {
	mu   sync.Mutex
	data []byte
}

func (t *Tail) Write(p []byte) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	n := len(p)
	if n >= 65536 {
		t.data = append(t.data[:0], p[n-65536:]...)
	} else {
		drop := len(t.data) + n - 65536
		if drop > 0 {
			copy(t.data, t.data[drop:])
			t.data = t.data[:len(t.data)-drop]
		}
		t.data = append(t.data, p...)
	}
	return n, nil
}
func (t *Tail) String() string { t.mu.Lock(); defer t.mu.Unlock(); return string(t.data) }

func processFailure(p Process, component, phase, message string) *Diagnostic {
	cause := p.Err()
	diagnostic := &Diagnostic{Component: component, Phase: phase, ExitCode: exitCode(cause)}
	var previous *Diagnostic
	if errors.As(cause, &previous) {
		diagnostic.Command = previous.Command
		diagnostic.Output = previous.Output
		cause = previous.Err
	}
	if cause == nil {
		diagnostic.Err = errors.New(message)
	} else {
		diagnostic.Err = fmt.Errorf("%s: %w", message, cause)
	}
	return diagnostic
}
