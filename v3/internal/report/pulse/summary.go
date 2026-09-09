package pulse

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/internal/report"
)

// osGetwd is a tiny indirection so highlightBody's cwd-resolution can be
// stubbed in tests if we ever need to.
var osGetwd = os.Getwd

// renderSummaryLocked draws the end-of-build summary block. This is where the
// build's *story* gets told: the verdict, the counts, where time went, and
// (on failure) one panel per failed step.
//
// Layout, top to bottom:
//
//	[verdict]  build (succeeded|failed)    <duration>
//
//	[counts]   N ran  ·  N cached  ·  N skipped[  ·  N failed]
//
//	[slowest]  Slowest
//	           task              ━━━━━━━━━━━━━━  1.4s  ◀ critical
//	           ...
//
//	[panels]   (on failure) one rounded-border panel per failed step
//
// Each region only renders if it has data — the success path collapses to
// verdict + counts + slowest.
func (r *Reporter) renderSummaryLocked(dur time.Duration, ok bool) {
	// Closing rule — mirrors the header signature. Two horizontal accent
	// strokes bracket the build, framing the scrollback log as a single
	// visual unit.
	cols, _ := r.termWidthLocked()
	ruleW := cols - 2
	if ruleW < 10 {
		ruleW = 10
	}
	// Closing rule fades the opposite way — full saturation on the left,
	// trailing off into dim on the right. Mirrors the opening curtain-rise
	// with a closing curtain-fall. Colour is determined by the verdict so
	// the eye gets red on failure even at peripheral vision. Width matches
	// the right edge of the elapsed-time column for alignment with the rest
	// of the build's scrollback.
	verdictCol := Success
	if !ok {
		verdictCol = Failure
	}
	rule := r.s.gradientRule(GlyphRule, ruleW, verdictCol, Dim)

	fmt.Fprintln(r.w)
	fmt.Fprintf(r.w, "  %s\n", rule)
	fmt.Fprintln(r.w)
	r.writeVerdictLocked(dur, ok)

	if len(r.completed) > 0 {
		fmt.Fprintln(r.w)
		r.writePhaseBarLocked()
		// "Slowest" + "◀ critical" annotation is a debugging tool: it tells
		// you where to focus when triaging a slow build. At Normal verbosity
		// most users don't need it (the phase bar already shows where time
		// went at the category level), so we gate the bottleneck callout
		// behind Debug.
		if r.level >= report.Debug {
			fmt.Fprintln(r.w)
			r.writeSlowestLocked()
		}
	}

	if len(r.artifacts) > 0 {
		fmt.Fprintln(r.w)
		r.writeArtifactsLocked()
	}

	if !ok && len(r.failures) > 0 {
		fmt.Fprintln(r.w)
		r.writeFailuresLocked()
	}

	fmt.Fprintln(r.w)
}

// writeArtifactsLocked draws the "Output" section: one row per registered
// build artifact, with the path on the left and a right-aligned
// human-readable size. Paths are wrapped in OSC 8 hyperlinks pointing at
// the resolved absolute file URI so users can click to open in their editor
// or file manager. Kind, if set, appears as a dim suffix on the path.
//
// Visual style mirrors the slowest / phase rows above: a faint section
// header, then rows indented two cells with the right-aligned size column
// extending to the same right edge as the rules above.
func (r *Reporter) writeArtifactsLocked() {
	cols, _ := r.termWidthLocked()
	cwd, _ := osGetwd()

	// Right column width: enough for "999.9 MiB" + a small breathing room.
	const sizeW = 10

	fmt.Fprintf(r.w, "  %s\n", r.s.faint("Output"))
	for _, a := range r.artifacts {
		uri := ""
		if cwd != "" {
			if filepath.IsAbs(a.Path) {
				uri = "file://" + a.Path
			} else {
				uri = "file://" + filepath.Join(cwd, a.Path)
			}
		}
		path := r.s.link(uri, a.Path)
		if a.Kind != "" {
			path += "  " + r.s.faint(a.Kind)
		}
		size := r.s.fg(Subtle, humanSize(a.Size))
		gap := cols - 2 - visibleWidth(path) - sizeW
		if gap < 2 {
			gap = 2
		}
		fmt.Fprintf(r.w, "  %s%s%s\n",
			path,
			strings.Repeat(" ", gap),
			padLeft(size, sizeW))
	}
}

// humanSize formats a byte count as "1.2 KiB", "12.1 MiB", etc. We use
// binary units (1024-based) because that's what `ls -lh` and most tooling
// emit; it also gives slightly larger numerator values which read as more
// concrete ("12 MiB" vs "12.6 MB").
func humanSize(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for n2 := n / unit; n2 >= unit; n2 /= unit {
		div *= unit
		exp++
	}
	units := []string{"KiB", "MiB", "GiB", "TiB", "PiB"}
	if exp >= len(units) {
		exp = len(units) - 1
	}
	return fmt.Sprintf("%.1f %s", float64(n)/float64(div), units[exp])
}

// writePhaseBarLocked renders a single-line stacked horizontal bar showing
// where the build's cpu time went, grouped by phase. Phases are inferred from
// task-name prefixes (a coarse heuristic, but the prefix conventions in wake
// Taskfiles are stable enough to make this useful at a glance).
//
// The bar reads:
//
//	Where time went  ▰▰▰▰▰▰ compile 5.2s ▰▰▰▰ test 3.1s ▰▰ package 1.2s
//
// Each segment gets its own colour: compile=accent, test=success,
// package=warning, prepare=cached, other=dim. Reading the bar tells you what
// kind of work dominated this build without scanning the slowest table.
func (r *Reporter) writePhaseBarLocked() {
	type phase struct {
		name  string
		col   Colour
		total time.Duration
	}
	phases := []*phase{
		{name: "compile", col: Accent},
		{name: "test", col: Success},
		{name: "package", col: Warning},
		{name: "prepare", col: Cached},
		{name: "other", col: Dim},
	}
	for _, c := range r.completed {
		if c.status == report.StatusCached || c.status == report.StatusSkipped {
			continue
		}
		idx := classifyPhase(c.name)
		phases[idx].total += c.duration
	}

	var total time.Duration
	for _, p := range phases {
		total += p.total
	}
	if total == 0 {
		return
	}

	cols, _ := r.termWidthLocked()
	// The bar and the legend line need to start at the same column. Compute
	// the prefix width explicitly so we can pad the legend line to match.
	const (
		labelText = "Where time went"
		indent    = 2
		labelW    = 15 // visible width of labelText
		barGap    = 2
		prefixW   = indent + labelW + barGap
	)
	barW := cols - prefixW - 2
	if barW < 20 {
		barW = 20
	}
	if barW > 60 {
		barW = 60
	}

	// Compute per-phase cell counts proportional to time share. We round to
	// nearest cell; any rounding leftover gets added to the largest phase so
	// the bar fills exactly barW cells.
	cells := make([]int, len(phases))
	used := 0
	largest, largestN := 0, 0
	for i, p := range phases {
		n := int(float64(barW)*float64(p.total)/float64(total) + 0.5)
		cells[i] = n
		used += n
		if n > largestN {
			largestN = n
			largest = i
		}
	}
	cells[largest] += barW - used

	var bar strings.Builder
	for i, p := range phases {
		if cells[i] <= 0 {
			continue
		}
		bar.WriteString(r.s.fg(p.col, strings.Repeat(BarFilled, cells[i])))
	}

	// Legend entries are joined with two-space separators so the first entry
	// can sit flush with the bar's left edge.
	var entries []string
	for _, p := range phases {
		if p.total == 0 {
			continue
		}
		entries = append(entries, fmt.Sprintf("%s %s",
			r.s.fg(p.col, p.name),
			r.s.faint(fmtDur(p.total))))
	}
	legend := strings.Join(entries, "  ")

	prefix := strings.Repeat(" ", indent) +
		r.s.faint(padRight(labelText, labelW)) +
		strings.Repeat(" ", barGap)
	pad := strings.Repeat(" ", prefixW)
	fmt.Fprintf(r.w, "%s%s\n", prefix, bar.String())
	fmt.Fprintf(r.w, "%s%s\n", pad, legend)
}

// classifyPhase returns the index into the phases slice that task should
// be counted under. Heuristic — matches the prefix conventions used in wake
// Taskfiles: compile/build, test/vet/lint, package/sign/dmg/notarize,
// prepare (mod/generate/bindings/install), other.
func classifyPhase(name string) int {
	n := strings.ToLower(name)
	// Strip leading platform/common namespace so "darwin:test:unit" still
	// classifies as test.
	for _, p := range []string{"darwin:", "linux:", "windows:", "ios:", "android:", "common:"} {
		n = strings.TrimPrefix(n, p)
	}
	switch {
	case strings.HasPrefix(n, "build"), strings.HasPrefix(n, "compile"),
		strings.HasPrefix(n, "link"), strings.HasPrefix(n, "lipo"):
		return 0 // compile
	case strings.HasPrefix(n, "test"), strings.HasPrefix(n, "vet"),
		strings.HasPrefix(n, "lint"), strings.HasPrefix(n, "check"):
		return 1 // test
	case strings.HasPrefix(n, "package"), strings.HasPrefix(n, "codesign"),
		strings.HasPrefix(n, "notarize"), strings.HasPrefix(n, "sign"),
		strings.HasPrefix(n, "dmg"), strings.HasPrefix(n, "deb"),
		strings.HasPrefix(n, "rpm"), strings.HasPrefix(n, "msi"):
		return 2 // package
	case strings.HasPrefix(n, "go:mod"), strings.HasPrefix(n, "generate"),
		strings.HasPrefix(n, "bindings"), strings.HasPrefix(n, "install"),
		strings.HasPrefix(n, "deps"), strings.HasPrefix(n, "frontend"):
		return 3 // prepare
	default:
		return 4 // other
	}
}

func (r *Reporter) writeVerdictLocked(dur time.Duration, ok bool) {
	// Sum the per-step durations to get "cpu time". For parallel builds the
	// ratio of cpu/wall is the speedup; for serial builds it equals 1 and we
	// hide the redundant detail.
	var cpu time.Duration
	for _, c := range r.completed {
		cpu += c.duration
	}
	speedup := ""
	if dur > 0 && cpu > dur+150*time.Millisecond {
		factor := float64(cpu) / float64(dur)
		speedup = "  " + r.s.faint(fmt.Sprintf("%s cpu · %.1f× speedup", fmtDur(cpu), factor))
	}

	var headline, durTail string
	if ok {
		headline = fmt.Sprintf("  %s  %s",
			r.s.fg(Success, GlyphOK),
			r.s.bold(r.s.fg(Success, "build succeeded")))
		durTail = "    " + r.s.faint(fmtDur(dur)) + speedup
	} else {
		headline = fmt.Sprintf("  %s  %s",
			r.s.fg(Failure, GlyphFail),
			r.s.bold(r.s.fg(Failure, "build failed")))
		durTail = "    " + r.s.faint("after "+fmtDur(dur)) + speedup
	}

	// Counts ride at the right edge of the verdict line — saves a whole
	// vertical row that an isolated counts line would have eaten.
	counts := r.renderCountsInline()
	left := headline + durTail
	cols, _ := r.termWidthLocked()
	gap := cols - visibleWidth(left) - visibleWidth(counts)
	if gap < 2 {
		gap = 2
	}
	fmt.Fprintf(r.w, "%s%s%s\n", left, strings.Repeat(" ", gap), counts)
}

// renderCountsInline formats "N ran · N cached · N skipped · N failed" for
// the right-side trailing of the verdict line. Zero-count entries are
// hidden; "ran" always shows (zero ran is itself informative).
func (r *Reporter) renderCountsInline() string {
	var ran, cached, skipped, failed int
	for _, c := range r.completed {
		switch c.status {
		case report.StatusCached:
			cached++
		case report.StatusSkipped:
			skipped++
		case report.StatusFailed:
			failed++
		default:
			ran++
		}
	}
	parts := []string{r.s.fg(Success, fmt.Sprintf("%d ran", ran))}
	if cached > 0 {
		parts = append(parts, r.s.fg(Cached, fmt.Sprintf("%d cached", cached)))
	}
	if skipped > 0 {
		parts = append(parts, r.s.fg(Dim, fmt.Sprintf("%d skipped", skipped)))
	}
	if failed > 0 {
		parts = append(parts, r.s.fg(Failure, fmt.Sprintf("%d failed", failed)))
	}
	sep := "  " + r.s.fg(Dim, GlyphBullet) + "  "
	return strings.Join(parts, sep)
}

// writeSlowestLocked renders the per-step duration bar chart, sorted by
// duration descending, capped at the top 5. The slowest step is annotated
// with a "◀ critical" marker — for a serial build this *is* the critical path;
// for parallel it's the bottleneck, which is the more interesting datum.
func (r *Reporter) writeSlowestLocked() {
	type row struct {
		name string
		dur  time.Duration
	}
	rows := make([]row, 0, len(r.completed))
	for _, c := range r.completed {
		if c.status == report.StatusCached || c.status == report.StatusSkipped {
			continue
		}
		name := c.label
		if name == "" {
			name = displayName(c.name)
		}
		rows = append(rows, row{name: name, dur: c.duration})
	}
	if len(rows) == 0 {
		return
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].dur > rows[j].dur })

	const topN = 5
	if len(rows) > topN {
		rows = rows[:topN]
	}
	nameW := 0
	for _, row := range rows {
		if w := visibleWidth(row.name); w > nameW {
			nameW = w
		}
	}
	maxDur := rows[0].dur

	cols, _ := r.termWidthLocked()
	// barW is whatever's left after name(=nameW), duration(=6 padded),
	// the "◀ critical" annotation (~14), three two-space gaps, and indent (2).
	barW := cols - nameW - 6 - 14 - 6 - 2
	if barW < 8 {
		barW = 8
	}
	if barW > 40 {
		barW = 40
	}

	fmt.Fprintf(r.w, "  %s\n", r.s.faint("Slowest"))
	for i, row := range rows {
		bar := r.s.barchart(float64(row.dur), float64(maxDur), barW, i)
		annotation := ""
		if i == 0 {
			annotation = "  " + r.s.fg(Accent, GlyphCrit+" critical")
		}
		fmt.Fprintf(r.w, "  %s  %s  %s%s\n",
			padRight(row.name, nameW),
			bar,
			padLeft(fmtDur(row.dur), 6),
			annotation)
	}
}

// writeFailuresLocked draws one rounded-border panel per failed step. With
// more than one failure it leads with a banner counting them so the reader
// gets the headline before scanning panels.
func (r *Reporter) writeFailuresLocked() {
	if len(r.failures) > 1 {
		fmt.Fprintf(r.w, "  %s  %s\n",
			r.s.fg(Failure, GlyphFail),
			r.s.bold(r.s.fg(Failure, fmt.Sprintf("%d tasks failed", len(r.failures)))))
		fmt.Fprintln(r.w)
	}
	for i, fs := range r.failures {
		if i > 0 {
			fmt.Fprintln(r.w)
		}
		r.writePanelLocked(fs.name, fs.f)
	}
}

// writePanelLocked draws one failure panel with a rounded border, a header
// row containing the task name, and a body containing command/status/output.
// We compute the border width from the widest body line so the panel never
// stretches edge-to-edge — that "right-sized box" is what makes it feel like
// a panel rather than a banner.
//
// Body lines pass through highlightBody, which bolds FAIL/PASS/ERROR
// keywords and wraps file:line references in OSC 8 hyperlinks. Both
// transformations preserve visible width, so the box math below is
// unaffected.
func (r *Reporter) writePanelLocked(name string, f report.Failure) {
	const labelW = 9 // "command  ", "status   ", "error    "
	cwd, _ := osGetwd()
	var body strings.Builder
	if f.Command != "" {
		fmt.Fprintf(&body, "%s %s\n",
			r.s.bold(r.s.fg(Failure, padRight("command", labelW))),
			f.Command)
	}
	// Skip the "status exited N" body row when there's an exit code — the
	// panel header carries the "exit N" badge already, and repeating it in
	// the body just adds chrome. Keep the "error" row when there's an err
	// because the body needs to carry the error *message*, not just its
	// type.
	if f.ExitCode == 0 && f.Err != nil {
		fmt.Fprintf(&body, "%s %s\n",
			r.s.bold(r.s.fg(Failure, padRight("error", labelW))),
			f.Err.Error())
	}
	out := strings.TrimRight(f.Output, "\n")
	if out != "" {
		body.WriteString("\n")
		body.WriteString(tailLines(out, 20))
	}
	bodyStr := strings.TrimRight(body.String(), "\n")

	bodyLines := strings.Split(bodyStr, "\n")
	for i, ln := range bodyLines {
		bodyLines[i] = r.s.highlightBody(ln, cwd)
	}
	cols, _ := r.termWidthLocked()
	inner := 0
	for _, ln := range bodyLines {
		if w := visibleWidth(ln); w > inner {
			inner = w
		}
	}
	if w := visibleWidth(name) + 4; w > inner {
		inner = w
	}
	maxInner := cols - 4
	if inner > maxInner {
		inner = maxInner
	}

	// Box geometry. A body row is "│ <padded(inner)> │" — visible width
	// inner+4. Top and bottom borders match that exact width.
	//
	// Top: "╭─ name ─...─ <badge> ─╮" — name on left, optional exit/error
	// badge tucked into the right side of the top border so the panel header
	// carries the headline in one line instead of two.
	badge := ""
	badgeW := 0
	if f.ExitCode != 0 {
		s := fmt.Sprintf("exit %d", f.ExitCode)
		badge = r.s.bold(r.s.fg(Failure, s))
		badgeW = visibleWidth(s)
	} else if f.Err != nil {
		s := "error"
		badge = r.s.bold(r.s.fg(Failure, s))
		badgeW = visibleWidth(s)
	}

	// rule cells between the name and either the badge (if present) or the
	// right corner. Layout when there's a badge:
	//   ╭─ name ─...─ badge ─╮
	// Visible: 5 + nameW + ruleLeft + 1 + badgeW + 2 = inner + 4
	// So: ruleLeft = inner - nameW - badgeW - 4
	var top string
	if badgeW > 0 {
		ruleLeft := inner - visibleWidth(name) - badgeW - 4
		if ruleLeft < 1 {
			ruleLeft = 1
		}
		top = "╭─ " + r.s.bold(r.s.fg(Failure, name)) + " " +
			strings.Repeat("─", ruleLeft) + " " + badge + " ─╮"
	} else {
		rule := inner - 1 - visibleWidth(name)
		if rule < 1 {
			rule = 1
		}
		top = "╭─ " + r.s.bold(r.s.fg(Failure, name)) + " " +
			strings.Repeat("─", rule) + "╮"
	}
	bot := "╰" + strings.Repeat("─", inner+2) + "╯"

	fmt.Fprintf(r.w, "  %s\n", r.s.fg(Failure, top))
	for _, ln := range bodyLines {
		// Body lines past the available inner width get truncated with an
		// ellipsis — wrapping a compile-error line at a panel boundary tends
		// to mangle the file:line:col pattern that the reader actually cares
		// about, so we'd rather show "…" than break the column.
		clipped := truncate(ln, inner)
		padded := padRight(clipped, inner)
		fmt.Fprintf(r.w, "  %s %s %s\n",
			r.s.fg(Failure, "│"),
			padded,
			r.s.fg(Failure, "│"))
	}
	fmt.Fprintf(r.w, "  %s\n", r.s.fg(Failure, bot))
}

// tailLines returns at most n lines from the end of s, prefixed by an "…"
// marker if any were dropped.
func tailLines(s string, n int) string {
	lines := strings.Split(s, "\n")
	if len(lines) <= n {
		return s
	}
	return "…\n" + strings.Join(lines[len(lines)-n:], "\n")
}
