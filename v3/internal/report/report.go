// Package report is the build-reporting contract for wake.
//
// It is a leaf package: it has no dependencies on the executor, on lipgloss, or
// on any producer. Three kinds of code use it:
//
//   - The executor (and anything orchestrating a build) drives a [Reporter]
//     through the build lifecycle.
//   - A renderer (see report/pulse) implements [Reporter] to draw the build.
//   - A producer (e.g. the bindings generator, which runs as a child process)
//     pushes live feedback into the active build UI via [Emit] — over a small
//     wire protocol when it is a subprocess, or directly when in-process.
//
// Keeping the contract here lets the code generator give live feedback to the
// build UI instead of printing to its own logger, without depending on wake.
package report

import (
	"encoding/json"
	"io"
	"strings"
	"time"
)

// Verbosity selects how much of a build is shown.
type Verbosity int

const (
	// Silent shows nothing but failures.
	Silent Verbosity = iota
	// Normal shows the plan, one line per step, and failures. The default.
	Normal
	// Verbose additionally shows each command and streams subprocess output live.
	Verbose
	// Debug additionally shows resolver internals (dep wiring, var refs, DAG order).
	Debug
)

// Status is the terminal outcome of a step.
type Status int

const (
	StatusOK      Status = iota // ran successfully
	StatusCached                // skipped, cache hit
	StatusSkipped               // skipped (up-to-date, platform mismatch, method:none)
	StatusFailed                // a command exited non-zero
)

// Failure carries everything needed to render a clean error for a failed step.
type Failure struct {
	Task     string // fully-qualified task name
	Command  string // the command that failed (already template-expanded)
	ExitCode int    // process exit code, or -1 if unknown
	Output   string // captured stdout+stderr of the failing command
	Err      error  // the underlying error
}

// Artifact is one build output the executor wants surfaced in the summary
// (binary, bundle, archive, etc.). The wake executor builds these from the
// `generates:` declarations of Taskfile tasks and registers them via
// [Reporter.Artifact] when the task succeeds.
//
// Size is the file's size in bytes; renderers format it human-readably. If
// Size is 0 the renderer may stat Path itself; pass a non-zero value to skip
// the stat (useful in tests and dry-runs where the artifact may not yet exist
// on disk by the time the renderer needs it).
//
// Kind is an optional one-word label ("binary", "bundle", "archive", "icon")
// shown next to the path; empty Kind hides the label.
type Artifact struct {
	Path string
	Size int64
	Kind string
}

// StepID identifies a step. It is returned by [Reporter.StepStart] and
// passed back to every step-scoped method (Info/Command/Output/End/Failed),
// so the renderer can address the correct row even when several steps are
// in flight concurrently. A zero ID is reserved for "no step" and is a
// no-op in every method that accepts it.
type StepID uint64

// Reporter receives the lifecycle of a build. Implementations render it.
//
// Multiple steps may be in flight simultaneously under parallel execution:
// each call to StepStart returns a fresh StepID that the caller threads
// through StepInfo / StepCommand / StepOutput / StepEnd / StepFailed for
// that specific step. A nil-safe no-op is available as [Nop].
type Reporter interface {
	// BuildStart begins a build. verb is the command the user ran (e.g. "build",
	// may be empty); target is the resolved task name the DAG says has totalSteps
	// tasks.
	BuildStart(verb, target string, totalSteps int)
	// StepStart announces a task is starting and returns the StepID that
	// subsequent step-scoped methods must reference.
	StepStart(name, label string) StepID
	// StepInfo is a sub-line attributed to the given step, typically pushed
	// by a producer such as the code generator ("generated 12 bindings").
	StepInfo(id StepID, msg string)
	// StepCommand reports the command a step is about to run (shown when Verbose).
	StepCommand(id StepID, cmd string)
	// StepOutput is one line of live subprocess output (shown when Verbose).
	StepOutput(id StepID, line string)
	// StepEnd closes the step with its outcome and duration.
	StepEnd(id StepID, status Status, dur time.Duration)
	// StepFailed closes the step with a failure to render.
	StepFailed(id StepID, f Failure)
	// BuildEnd closes the build.
	BuildEnd(dur time.Duration, ok bool)
	// Artifact registers a build output (binary, bundle, archive, etc.) for
	// display in the end-of-build summary.
	Artifact(a Artifact)
	// Debug renders one diagnostic line (only shown at Debug verbosity).
	Debug(line DebugLine)
	// Level reports the reporter's verbosity so callers can skip expensive work.
	Level() Verbosity
}

// DebugLine is a structured diagnostic for Debug verbosity. Rendering it as
// category badge + subject + arrow + key/value fields keeps the otherwise
// undifferentiated debug stream scannable.
type DebugLine struct {
	Category string       // "dag", "dep", "var", "exec" — drives the badge colour
	Subject  string       // primary identifier (task or var name)
	Arrow    string       // optional "→ <target>" detail
	Fields   []DebugField // optional key=value pairs
}

// DebugField is one key=value pair on a DebugLine.
type DebugField struct{ Key, Val string }

// Nop is a Reporter that does nothing. It is the zero value used by tests and
// by any Executor that has no reporter wired in.
type Nop struct{}

func (Nop) BuildStart(string, string, int)        {}
func (Nop) StepStart(string, string) StepID       { return 0 }
func (Nop) StepInfo(StepID, string)               {}
func (Nop) StepCommand(StepID, string)            {}
func (Nop) StepOutput(StepID, string)             {}
func (Nop) StepEnd(StepID, Status, time.Duration) {}
func (Nop) StepFailed(StepID, Failure)            {}
func (Nop) BuildEnd(time.Duration, bool)          {}
func (Nop) Artifact(Artifact)                     {}
func (Nop) Debug(DebugLine)                       {}
func (Nop) Level() Verbosity                      { return Normal }

// active is the reporter for the current in-process build. Producers that run
// in-process (not as a subprocess) reach the UI through it.
var active Reporter = Nop{}

// SetActive installs r as the reporter for in-process producers. Passing nil
// resets to a no-op. The executor sets this for the duration of a build.
func SetActive(r Reporter) {
	if r == nil {
		active = Nop{}
		return
	}
	active = r
}

// Active returns the in-process reporter (never nil).
func Active() Reporter { return active }

// --- Wire protocol -------------------------------------------------------
//
// A producer running as a child process cannot reach the parent's active
// Reporter, so it serialises events onto its stdout as single lines that the
// parent's output reader recognises and routes. The sentinel is built from
// ASCII Unit Separator (0x1f), which does not occur in ordinary program output,
// so a non-wake consumer simply prints these lines verbatim and a wake consumer
// strips them.

const sentinel = "\x1fwake\x1f"

// EventKind classifies a producer event.
type EventKind string

const (
	KindInfo   EventKind = "info"
	KindStatus EventKind = "status"
	KindWarn   EventKind = "warn"
	KindError  EventKind = "error"
)

// Event is a single piece of live feedback from a producer.
type Event struct {
	Kind EventKind `json:"k"`
	Msg  string    `json:"m,omitempty"`
}

// Encode renders ev as a single wire line (no trailing newline).
func Encode(ev Event) string {
	b, err := json.Marshal(ev)
	if err != nil {
		return ""
	}
	return sentinel + string(b)
}

// Decode parses a wire line produced by Encode. ok is false for ordinary output.
func Decode(line string) (ev Event, ok bool) {
	rest, found := strings.CutPrefix(line, sentinel)
	if !found {
		return Event{}, false
	}
	if err := json.Unmarshal([]byte(rest), &ev); err != nil {
		return Event{}, false
	}
	return ev, true
}

// Emit writes ev to w as a wire line. A subprocess producer (e.g. the bindings
// generator running under wake) calls this with os.Stdout; the parent build
// reader recognises the line and routes it to the live UI.
func Emit(w io.Writer, ev Event) {
	if line := Encode(ev); line != "" {
		_, _ = io.WriteString(w, line+"\n")
	}
}

// Route applies a decoded producer event to the named step of r.
func Route(r Reporter, id StepID, ev Event) {
	switch ev.Kind {
	case KindError, KindWarn, KindInfo, KindStatus:
		r.StepInfo(id, ev.Msg)
	}
}
