package pulse

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/internal/report"
	"github.com/wailsapp/wails/v3/internal/report/pulse/ansi"
)

// Reporter renders a wake build as the "Pulse" TUI.
//
// In Normal verbosity it uses **skeleton mode**: at BuildStart it pre-paints
// one row per planned step (all in "pending" state), and each StepStart /
// StepEnd updates that row in place. The terminal never scrolls during the
// build — the skeleton stays anchored, pending rows fill in as active rows
// fill in as completed rows, and the end-of-build summary unfolds beneath.
//
// In Verbose / Debug verbosities it falls back to **scrollback mode** (the
// original Buck2-Superconsole shape): completed step lines promote into the
// scrollback above a pinned region that holds only the in-flight step(s) and
// the progress strip. Interleaved StepCommand / StepOutput / Debug lines need
// the scrolling history to make sense.
//
// Pulse implements [report.Reporter]. The serial path matches the existing
// termui contract exactly — one in-flight step at a time. The parallel path
// is exposed through the additional methods on this concrete type so callers
// that drive multiple workers can hand each its own step ID. Mixing the two
// in the same build is supported.
type Reporter struct {
	w        io.Writer
	f        *os.File // nil if w isn't an *os.File
	level    report.Verbosity
	animate  bool
	skeleton bool // skeleton mode (Normal verbosity, TTY, known total)
	s        *styler
	region   *liveRegion

	mu sync.Mutex

	// build state
	started    bool
	total      int
	idx        int
	verb       string
	target     string
	buildStart time.Time

	// in-flight steps. nil entries are tombstones from finished steps left in
	// place to preserve render order until the next paint. Used in scrollback
	// mode only — in skeleton mode the slots[] list owns step state.
	active []*activeStep
	nextID StepID

	// skeleton mode: pre-allocated slots, one per planned step. nextSlot is
	// the index of the next slot to be claimed by a StepStart. Active and
	// completed step state both live in the slot rather than in active/
	// completed slices.
	slots    []*stepRow
	nextSlot int

	// finished steps and failures, retained for the end-of-build summary.
	completed []completedStep
	failures  []failedStep

	// artifacts: build outputs (binaries, bundles, archives) registered by
	// the executor as tasks with `generates:` declarations complete. Rendered
	// in their own block at the end of the summary so a user finishing a
	// `wails3 build` sees exactly what came out and how big it is.
	artifacts []report.Artifact

	// throughput samples: completions in each ~250 ms window, for the
	// sparkline. Capped to ~40 samples so older history rolls off.
	throughput []float64
	winStart   time.Time
	winCount   int

	// spinner / ticker
	frame  int
	ticker *time.Ticker
	stop   chan struct{}
}

// StepID is re-exported from the report package.
type StepID = report.StepID

// Artifact is re-exported from the report package so call sites that
// reference *pulse.Reporter need only one import.
type Artifact = report.Artifact

type activeStep struct {
	id        StepID
	name      string
	label     string
	detail    string
	startedAt time.Time
}

type completedStep struct {
	name     string
	label    string
	status   report.Status
	duration time.Duration
}

// stepStatus is the lifecycle state of a step in skeleton mode.
type stepStatus int

const (
	stepPending stepStatus = iota
	stepActive
	stepOK
	stepCached
	stepSkipped
	stepFailed
)

// stepRow is one row in the pre-painted skeleton. Status transitions are
// pending → active → (ok | cached | skipped | failed). Detail is shown only
// while status == active.
type stepRow struct {
	idx       int    // 1-based row number; matches the displayed [N/M] counter
	id        StepID // 0 while pending; assigned at StepStart
	name      string
	label     string
	detail    string
	status    stepStatus
	startedAt time.Time
	duration  time.Duration
}

type failedStep struct {
	name string
	at   time.Time
	f    report.Failure
}

// New builds a Reporter writing to w at the given verbosity. It auto-detects
// whether w is a terminal and whether colour is permitted (TERM, COLORTERM,
// NO_COLOR), and degrades to plain one-line-per-step output if not.
func New(w io.Writer, level report.Verbosity) *Reporter {
	f, _ := w.(*os.File)
	profile := DetectProfile(f)
	isTTY := ansi.IsTerminal(f)
	return &Reporter{
		w:       w,
		f:       f,
		level:   level,
		animate: isTTY && level != report.Silent,
		s:       newStyler(profile),
		region:  newLiveRegion(w),
	}
}

// Level reports the reporter's verbosity, satisfying report.Reporter.
func (r *Reporter) Level() report.Verbosity { return r.level }

// BuildStart begins the build. Prints the header to scrollback, decides
// whether to use skeleton mode, pre-paints the skeleton if so, and starts
// the redraw ticker.
func (r *Reporter) BuildStart(verb, target string, total int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.started = true
	r.total = total
	r.verb = verb
	r.target = target
	r.buildStart = time.Now()
	r.winStart = r.buildStart

	// Skeleton mode requires: a TTY, Normal verbosity (Verbose/Debug need
	// scrolling for interleaved output), and a known step count.
	r.skeleton = r.animate && r.level == report.Normal && total > 0
	if r.skeleton {
		r.slots = make([]*stepRow, total)
		for i := 0; i < total; i++ {
			r.slots[i] = &stepRow{idx: i + 1, status: stepPending}
		}
	}

	if r.level == report.Silent {
		return
	}
	r.printHeaderLocked()
	if r.animate {
		r.startTickerLocked()
		if r.skeleton {
			r.repaintLocked() // paint the initial skeleton immediately
		}
	}
}

func (r *Reporter) printHeaderLocked() {
	var (
		count = ""
		plur  = "steps"
	)
	if r.total == 1 {
		plur = "step"
	}
	if r.total > 0 {
		count = r.s.faint(fmt.Sprintf("(%d %s)", r.total, plur))
	}
	// The header is composed of two lines: the build identity (verb › target
	// count) and a thin accent rule beneath it. The rule is wake's signature
	// stroke — it fades in from dim on the left to full accent on the right,
	// so the build's opening visual moment feels like a curtain rising.
	//
	// Rule width is sized to extend all the way to the right edge that the
	// step-row elapsed times sit against, so the rule and the elapsed column
	// line up exactly.
	cols, _ := r.termWidthLocked()
	ruleW := cols - 2
	if ruleW < 10 {
		ruleW = 10
	}
	rule := r.s.gradientRule(GlyphRule, ruleW, Dim, Accent)

	switch {
	case r.verb != "" && r.target != "":
		fmt.Fprintf(r.w, "\n  %s %s %s %s\n  %s\n\n",
			r.s.accentBold(r.verb),
			r.s.fg(Dim, GlyphArrow),
			r.target,
			count,
			rule)
	case r.target != "":
		fmt.Fprintf(r.w, "\n  %s %s\n  %s\n\n", r.target, count, rule)
	default:
		fmt.Fprintf(r.w, "\n  %s %s\n  %s\n\n", r.s.accentBold("build"), count, rule)
	}
}

// StepStart announces a step. Returns the StepID the caller threads through
// step-scoped methods.
func (r *Reporter) StepStart(name, label string) StepID {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.idx++
	id := r.allocIDLocked()
	if r.skeleton {
		r.claimSlotLocked(id, name, label)
	} else {
		r.active = append(r.active, &activeStep{
			id:        id,
			name:      name,
			label:     label,
			startedAt: time.Now(),
		})
	}
	r.repaintLocked()
	return id
}

// StepInfo attaches a detail line to the named step.
func (r *Reporter) StepInfo(id StepID, msg string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	msg = strings.TrimSpace(msg)
	if r.skeleton {
		if slot := r.findSlotLocked(id); slot != nil {
			slot.detail = msg
			r.repaintLocked()
		}
		return
	}
	if as := r.findLocked(id); as != nil {
		as.detail = msg
		r.repaintLocked()
	}
}

// StepCommand prints the command above the pinned region (Verbose only).
func (r *Reporter) StepCommand(id StepID, cmd string) {
	if r.level < report.Verbose {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	line := fmt.Sprintf("      %s %s",
		r.s.fg(Dim, "$"),
		r.s.fg(Subtle, cmd))
	r.region.promote([]string{line}, r.region.lines)
}

// StepOutput prints one captured output line above the pinned region (Verbose).
func (r *Reporter) StepOutput(id StepID, line string) {
	if r.level < report.Verbose {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	out := "      " + strings.TrimRight(line, "\n")
	r.region.promote([]string{out}, r.region.lines)
}

// StepEnd closes the named step.
func (r *Reporter) StepEnd(id StepID, status report.Status, dur time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.finishLocked(id, status, dur, nil)
}

// StepFailed closes the named step with a failure. Passing `0` for the
// duration lets finishLocked derive it from the step's own startedAt; the
// previous `time.Since(r.buildStart)` form measured the wrong interval —
// a step that started 200 ms in and ran for 50 ms would otherwise report
// a duration of 250 ms in the summary.
func (r *Reporter) StepFailed(id StepID, f report.Failure) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.finishLocked(id, report.StatusFailed, 0, &f)
}

// ParallelStepStart, ParallelStepInfo, ParallelStepEnd, ParallelStepFailed
// were earlier concrete-type aliases of the now-id-returning interface
// methods; with StepStart returning a StepID directly, callers use the
// standard interface for both serial and parallel flows.
func (r *Reporter) ParallelStepStart(name, label string) StepID {
	return r.StepStart(name, label)
}
func (r *Reporter) ParallelStepInfo(id StepID, msg string) { r.StepInfo(id, msg) }
func (r *Reporter) ParallelStepEnd(id StepID, status report.Status, dur time.Duration) {
	r.StepEnd(id, status, dur)
}
func (r *Reporter) ParallelStepFailed(id StepID, f report.Failure) {
	r.StepFailed(id, f)
}

// Artifact registers one build output for display in the end-of-build
// summary. Satisfies the report.Reporter interface. If a.Size is 0, the path
// is stat()ed to fill in the size; a non-zero size short-circuits the stat
// (useful for tests and demos). Safe to call concurrently with step
// lifecycle methods.
func (r *Reporter) Artifact(a report.Artifact) {
	if a.Path == "" {
		return
	}
	if a.Size == 0 {
		if info, err := os.Stat(a.Path); err == nil {
			a.Size = info.Size()
		}
	}
	r.mu.Lock()
	r.artifacts = append(r.artifacts, a)
	r.mu.Unlock()
}

// allocIDLocked returns a fresh non-zero StepID.
func (r *Reporter) allocIDLocked() StepID {
	r.nextID++
	return r.nextID
}

func (r *Reporter) findLocked(id StepID) *activeStep {
	if id == 0 {
		return nil
	}
	for _, a := range r.active {
		if a != nil && a.id == id {
			return a
		}
	}
	return nil
}

// claimSlotLocked finds the next pending slot, fills in name/label/id, and
// flips its status to active. If the build overflows its declared total
// (more StepStart calls than total promised at BuildStart) we append a new
// slot rather than panic, which keeps the demo robust even when scripted
// totals are slightly off.
func (r *Reporter) claimSlotLocked(id StepID, name, label string) *stepRow {
	if r.nextSlot >= len(r.slots) {
		idx := len(r.slots) + 1
		slot := &stepRow{idx: idx}
		r.slots = append(r.slots, slot)
		r.nextSlot = len(r.slots)
		r.populateSlotLocked(slot, id, name, label)
		return slot
	}
	slot := r.slots[r.nextSlot]
	r.nextSlot++
	r.populateSlotLocked(slot, id, name, label)
	return slot
}

func (r *Reporter) populateSlotLocked(slot *stepRow, id StepID, name, label string) {
	slot.id = id
	slot.name = name
	slot.label = label
	slot.status = stepActive
	slot.startedAt = time.Now()
}

func (r *Reporter) findSlotLocked(id StepID) *stepRow {
	if id == 0 {
		return nil
	}
	for _, s := range r.slots {
		if s != nil && s.id == id {
			return s
		}
	}
	return nil
}

// finishLocked closes out the named step. In skeleton mode this updates the
// slot in place (no scrollback motion); in scrollback mode it promotes a
// completed line above the pinned region.
func (r *Reporter) finishLocked(id StepID, status report.Status, dur time.Duration, failure *report.Failure) {
	if id == 0 {
		return
	}
	if r.skeleton {
		r.finishSlotLocked(id, status, dur, failure)
		return
	}
	var done *activeStep
	for i, a := range r.active {
		if a != nil && a.id == id {
			done = a
			r.active = append(r.active[:i], r.active[i+1:]...)
			break
		}
	}
	if done == nil {
		return
	}
	if dur == 0 {
		dur = time.Since(done.startedAt)
	}
	r.completed = append(r.completed, completedStep{
		name:     done.name,
		label:    done.label,
		status:   status,
		duration: dur,
	})
	if failure != nil {
		r.failures = append(r.failures, failedStep{
			name: done.name,
			at:   time.Now(),
			f:    *failure,
		})
	}
	r.recordThroughputLocked()
	if r.level == report.Silent && status != report.StatusFailed {
		return
	}
	completed := r.renderCompletedLine(done, status, dur)
	r.region.promote([]string{completed}, r.buildPinnedLocked())
}

// finishSlotLocked is the skeleton-mode finish path: update the slot's
// status/duration in place and add to the completed list so the end-of-build
// summary still has the data it needs.
func (r *Reporter) finishSlotLocked(id StepID, status report.Status, dur time.Duration, failure *report.Failure) {
	slot := r.findSlotLocked(id)
	if slot == nil {
		return
	}
	if dur == 0 {
		dur = time.Since(slot.startedAt)
	}
	slot.duration = dur
	slot.detail = ""
	switch status {
	case report.StatusCached:
		slot.status = stepCached
	case report.StatusSkipped:
		slot.status = stepSkipped
	case report.StatusFailed:
		slot.status = stepFailed
	default:
		slot.status = stepOK
	}
	r.completed = append(r.completed, completedStep{
		name:     slot.name,
		label:    slot.label,
		status:   status,
		duration: dur,
	})
	if failure != nil {
		r.failures = append(r.failures, failedStep{
			name: slot.name,
			at:   time.Now(),
			f:    *failure,
		})
	}
	r.recordThroughputLocked()
	r.repaintLocked()
}

// recordThroughputLocked bumps the count of the current 250 ms window. The
// sparkline shows the last ~10 s of build pace.
func (r *Reporter) recordThroughputLocked() {
	now := time.Now()
	if now.Sub(r.winStart) > 250*time.Millisecond {
		r.throughput = append(r.throughput, float64(r.winCount))
		if len(r.throughput) > 40 {
			r.throughput = r.throughput[len(r.throughput)-40:]
		}
		r.winStart = now
		r.winCount = 0
	}
	r.winCount++
}

// BuildEnd closes the build and prints the summary.
//
// In scrollback mode the pinned region is erased so the summary can take its
// place. In skeleton mode the settled skeleton stays exactly where it was
// drawn (so the user keeps seeing the per-step ledger), and the summary
// unfolds below it.
func (r *Reporter) BuildEnd(dur time.Duration, ok bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.stopTickerLocked()
	if r.skeleton {
		// Final paint with all slots in their settled status. We pin only
		// the step rows themselves — the progress strip is dropped because
		// at this point every slot already shows its outcome, so the
		// "100% — Xms" line is redundant chrome. Detach without erasing so
		// the settled skeleton flows into scrollback.
		r.region.paint(r.buildSkeletonRowsLocked())
		r.region.detach()
		fmt.Fprintln(r.w)
	} else {
		r.region.close()
	}
	if r.level == report.Silent && ok {
		return
	}
	r.renderSummaryLocked(dur, ok)
}

// Debug renders one diagnostic line. We rely on the existing termui-style
// "BADGE subject → arrow  k=v" layout; it reads well above the pinned region.
func (r *Reporter) Debug(d report.DebugLine) {
	if r.level < report.Debug {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	cat := strings.ToUpper(d.Category)
	if len(cat) > 4 {
		cat = cat[:4]
	}
	col := Dim
	switch d.Category {
	case "dag":
		col = Accent
	case "dep":
		col = Cached
	case "var":
		col = Warning
	}
	badge := r.s.bg(col, fmt.Sprintf(" %-4s ", cat))
	var b strings.Builder
	fmt.Fprintf(&b, "  %s %s", badge, d.Subject)
	if d.Arrow != "" {
		fmt.Fprintf(&b, " %s%s",
			r.s.fg(Dim, "→ "),
			r.s.fg(Subtle, d.Arrow))
	}
	for _, f := range d.Fields {
		fmt.Fprintf(&b, "  %s%s",
			r.s.fg(Dim, f.Key+"="),
			r.s.fg(Success, f.Val))
	}
	r.region.promote([]string{b.String()}, r.region.lines)
}

// startTickerLocked spins up the redraw goroutine.
func (r *Reporter) startTickerLocked() {
	r.stop = make(chan struct{})
	r.ticker = time.NewTicker(80 * time.Millisecond)
	stop := r.stop
	tick := r.ticker
	go func() {
		for {
			select {
			case <-stop:
				return
			case <-tick.C:
				r.mu.Lock()
				if len(r.active) > 0 {
					r.frame++
					r.repaintLocked()
				}
				r.mu.Unlock()
			}
		}
	}()
}

func (r *Reporter) stopTickerLocked() {
	if r.ticker != nil {
		r.ticker.Stop()
		close(r.stop)
		r.ticker = nil
		r.stop = nil
	}
}
