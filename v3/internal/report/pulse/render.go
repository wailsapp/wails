package pulse

import (
	"fmt"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/internal/report"
	"github.com/wailsapp/wails/v3/internal/report/pulse/ansi"
)

// repaintLocked rebuilds the pinned region from current state and atomically
// replaces what's on screen. Caller holds r.mu.
func (r *Reporter) repaintLocked() {
	if !r.animate {
		return
	}
	r.region.paint(r.buildPinnedLocked())
}

// buildPinnedLocked computes the pinned region's lines from current state.
//
// Layout, skeleton mode (Normal verbosity):
//
//	N slot rows (one per planned step) + blank + progress row
//	+ optional sparkline once parallel throughput samples are meaningful.
//
// Layout, scrollback mode (Verbose/Debug or unknown total):
//
//	serial   — 1 spinner row + 1 progress row
//	parallel — "Active" header, M spinner rows, blank, progress row,
//	           optional sparkline row.
//	idle     — empty (nil), which clears the region.
//
// Caller holds r.mu.
func (r *Reporter) buildPinnedLocked() []string {
	if !r.animate {
		return nil
	}
	if r.skeleton {
		return r.buildSkeletonLocked()
	}
	if len(r.active) == 0 {
		return nil
	}
	cols, _ := r.termWidthLocked()
	spin := r.spinnerLocked()
	var lines []string

	if len(r.active) == 1 {
		a := r.active[0]
		lines = append(lines, r.renderActiveLine(spin, a, cols))
	} else {
		lines = append(lines, "  "+r.s.faint("Active"))
		for _, a := range r.active {
			if a == nil {
				continue
			}
			lines = append(lines, r.renderActiveLine(spin, a, cols))
		}
		lines = append(lines, "")
	}

	lines = append(lines, r.renderProgressLine(cols))

	if len(r.active) > 1 && len(r.throughput) >= 4 {
		lines = append(lines, r.renderThroughputLine(cols))
	}

	return lines
}

// buildSkeletonLocked renders the full skeleton: one row per planned step,
// each in its current state (pending / active / completed). The progress
// strip sits beneath, separated by a blank line.
//
// The skeleton can grow taller than the terminal — for now we render every
// slot and trust the terminal to scroll. A future refinement could window
// around the currently-active rows.
func (r *Reporter) buildSkeletonLocked() []string {
	cols, _ := r.termWidthLocked()
	lines := r.buildSkeletonRowsLocked()
	lines = append(lines, "")
	lines = append(lines, r.renderProgressLine(cols))
	if len(r.throughput) >= 4 {
		lines = append(lines, r.renderThroughputLine(cols))
	}
	return lines
}

// buildSkeletonRowsLocked returns just the per-step rows (no trailing
// progress strip). Used at BuildEnd in skeleton mode where the step rows
// already carry the final status of every step and a 100% progress bar
// would be visual chrome.
func (r *Reporter) buildSkeletonRowsLocked() []string {
	cols, _ := r.termWidthLocked()
	spin := r.spinnerLocked()
	lines := make([]string, 0, len(r.slots))
	for _, slot := range r.slots {
		if slot == nil {
			continue
		}
		lines = append(lines, r.renderSlotLine(spin, slot, cols))
	}
	return lines
}

// renderSlotLine renders one skeleton row in the visual style appropriate
// for its current status. The four shapes are:
//
//	pending   "  · [N/M]" in dim
//	active    "▎ ⠦ [N/M] name • detail                              1.4s"
//	completed " ✓ [N/M] name                                       480ms"
//
// All right-aligned to the duration column so the rows stack into a single
// vertical column.
func (r *Reporter) renderSlotLine(spin string, slot *stepRow, cols int) string {
	counter := fmt.Sprintf("[%d/%d]", slot.idx, r.total)
	switch slot.status {
	case stepPending:
		return fmt.Sprintf("  %s %s",
			r.s.fg(Dim, GlyphSkipped),
			r.s.fg(Dim, counter))
	case stepActive:
		dur := time.Since(slot.startedAt)
		elapsed := fmtDur(dur)
		col := Subtle
		switch {
		case dur >= 30*time.Second:
			col = Failure
		case dur >= 5*time.Second:
			col = Warning
		}
		stripe := r.s.fg(Accent, "▎")
		name := slot.label
		if name == "" {
			name = displayName(slot.name)
		}
		left := fmt.Sprintf("%s %s %s %s",
			stripe,
			spin,
			r.s.fg(Dim, counter),
			name)
		if slot.detail != "" {
			left += "  " + r.s.fg(Dim, GlyphBullet+" "+slot.detail)
		}
		right := r.s.fg(col, elapsed)
		left = truncate(left, cols-visibleWidth(right)-2)
		gap := cols - visibleWidth(left) - visibleWidth(right)
		if gap < 1 {
			gap = 1
		}
		return left + strings.Repeat(" ", gap) + right
	default:
		// Completed states reuse renderCompletedLine via a synthesized
		// activeStep so the styling is identical to scrollback mode's
		// completed rows.
		as := &activeStep{id: slot.id, name: slot.name, label: slot.label}
		status := slotStatusToReportStatus(slot.status)
		return r.renderCompletedLine(as, status, slot.duration)
	}
}

func slotStatusToReportStatus(s stepStatus) report.Status {
	switch s {
	case stepCached:
		return report.StatusCached
	case stepSkipped:
		return report.StatusSkipped
	case stepFailed:
		return report.StatusFailed
	default:
		return report.StatusOK
	}
}

// renderActiveLine renders one in-flight step row. Width-aware: it allocates
// the visible columns to "▎ ⠦ [N/M] name detail elapsed" with right-aligned
// elapsed, truncating the detail before the name and the name before the
// counter if space runs out.
//
// The leading "▎" stripe is an accent-coloured quarter-block — a thin
// vertical bar on the left edge of each active row. It's a single-character
// "this is live" marker that differentiates active rows from completed ones
// at a glance, the way IDEs mark modified files with a stripe.
func (r *Reporter) renderActiveLine(spin string, a *activeStep, cols int) string {
	counter := r.counterFor(a.id)
	name := r.titleFor(a)
	dur := time.Since(a.startedAt)
	elapsed := fmtDur(dur)

	// Escalate the elapsed-time colour past humane thresholds so long-running
	// steps flag themselves before the user has to read the number. The
	// hand-tuned breakpoints are: under 5 s = boring, 5–30 s = notable,
	// over 30 s = a step the user probably wants to look at.
	col := Subtle
	switch {
	case dur >= 30*time.Second:
		col = Failure
	case dur >= 5*time.Second:
		col = Warning
	}

	// One column of indent + stripe + space + spinner + space + counter +
	// space + name + optional detail + elapsed (right-aligned at the end).
	const leftPad = 2
	right := r.s.fg(col, elapsed)
	rightW := visibleWidth(elapsed)

	// The stripe lives in column 0 — the gutter that's otherwise blank for
	// completed rows. That way the spinner stays at column 2, aligned
	// vertically with completed rows' status glyph, and the stripe doesn't
	// shift any content rightward.
	stripe := r.s.fg(Accent, "▎")
	left := fmt.Sprintf("%s %s %s %s",
		stripe,
		spin,
		r.s.fg(Dim, counter),
		name)
	if a.detail != "" {
		left += "  " + r.s.fg(Dim, GlyphBullet+" "+a.detail)
	}

	maxLeft := cols - rightW - 2
	if maxLeft < 10 {
		maxLeft = cols - leftPad
	}
	left = truncate(left, maxLeft)
	gap := cols - visibleWidth(left) - rightW
	if gap < 1 {
		gap = 1
	}
	return left + strings.Repeat(" ", gap) + right
}

func (r *Reporter) renderCompletedLine(a *activeStep, status report.Status, dur time.Duration) string {
	var glyph, tail string
	switch status {
	case report.StatusCached:
		glyph = r.s.fg(Cached, GlyphCached)
		tail = r.s.fg(Cached, "cached")
	case report.StatusSkipped:
		glyph = r.s.fg(Dim, GlyphSkipped)
		tail = r.s.fg(Dim, "skipped")
	case report.StatusFailed:
		glyph = r.s.fg(Failure, GlyphFail)
		tail = r.s.fg(Failure, fmtDur(dur))
	default:
		glyph = r.s.fg(Success, GlyphOK)
		// Tier the duration colour so scrollback reads as a heatmap. Snappy
		// completions (<200 ms) get a green tint that pairs visually with
		// the ✓ glyph; slow ones (>10 s) get a warning yellow that flags
		// the row for inspection on a post-hoc scroll-up. The middle tier
		// stays faint — most completions live there and shouldn't compete.
		switch {
		case dur < 200*time.Millisecond:
			tail = r.s.fg(Success, fmtDur(dur))
		case dur >= 10*time.Second:
			tail = r.s.fg(Warning, fmtDur(dur))
		default:
			tail = r.s.faint(fmtDur(dur))
		}
	}
	// Right-align the tail so completed lines stack into a clean column down
	// the right edge — the eye should be able to scan the duration column
	// without re-fixating on each row.
	cols, _ := r.termWidthLocked()
	counter := r.s.fg(Dim, r.counterFor(a.id))
	title := r.titleFor(a)
	left := fmt.Sprintf("  %s %s %s", glyph, counter, title)
	tailW := visibleWidth(tail)
	gap := cols - visibleWidth(left) - tailW
	if gap < 2 {
		gap = 2
	}
	return left + strings.Repeat(" ", gap) + tail
}

func (r *Reporter) renderProgressLine(cols int) string {
	var ratio float64
	if r.total > 0 {
		ratio = float64(len(r.completed)) / float64(r.total)
	}
	elapsed := fmtDur(time.Since(r.buildStart))
	eta := r.estimateETALocked()
	// Indented by 2 to align under the spinner column.
	return "  " + r.s.progressLine(ratio, cols-2, elapsed, eta)
}

func (r *Reporter) renderThroughputLine(cols int) string {
	spark := r.s.sparkline(r.throughput)
	label := r.s.faint("throughput")
	return "  " + label + "  " + spark
}

// estimateETALocked produces a rough ETA in seconds using mean-per-step from
// what we've seen so far. Not science — but adequate for "is this almost done"
// without a calibration model. Returns "" while we don't yet have a useful
// sample.
func (r *Reporter) estimateETALocked() string {
	if r.total <= 0 || len(r.completed) < 2 {
		return ""
	}
	done := len(r.completed)
	if done >= r.total {
		return ""
	}
	mean := time.Since(r.buildStart) / time.Duration(done)
	remaining := mean * time.Duration(r.total-done)
	return fmtDur(remaining)
}

func (r *Reporter) spinnerLocked() string {
	return r.s.fg(Accent, spinnerFrames[r.frame%len(spinnerFrames)])
}

func (r *Reporter) counterFor(id StepID) string {
	// We render the absolute index relative to total. Per-step indices are
	// stable because we assign idx at StepStart and the active step keeps
	// that position. For a quick path we use id as a proxy; this is fine
	// because nextID increases monotonically and matches r.idx semantics.
	return fmt.Sprintf("[%d/%d]", int(id), r.total)
}

func (r *Reporter) titleFor(a *activeStep) string {
	name := a.label
	if name == "" {
		name = displayName(a.name)
	}
	return name
}

func (r *Reporter) termWidthLocked() (int, int) {
	w, h := ansi.TermSize(r.f)
	if w > 140 {
		// Cap line length so the pinned region never feels sparse on
		// ultra-wide terminals.
		w = 140
	}
	return w, h
}

// displayName drops the leading platform / common namespace from a task name,
// matching the existing termui behaviour: under a darwin build,
// "darwin:common:go:mod:tidy" reads better as "go:mod:tidy".
func displayName(name string) string {
	for _, p := range []string{"darwin:", "linux:", "windows:", "ios:", "android:"} {
		if rest, ok := strings.CutPrefix(name, p); ok {
			name = rest
			break
		}
	}
	if rest, ok := strings.CutPrefix(name, "common:"); ok {
		name = rest
	}
	return name
}

func fmtDur(d time.Duration) string {
	switch {
	case d < time.Millisecond:
		return "0ms"
	case d < time.Second:
		return fmt.Sprintf("%dms", d.Milliseconds())
	case d < time.Minute:
		return fmt.Sprintf("%.1fs", d.Seconds())
	default:
		return fmt.Sprintf("%dm%02ds", int(d.Minutes()), int(d.Seconds())%60)
	}
}
