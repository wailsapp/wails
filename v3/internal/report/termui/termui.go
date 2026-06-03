// Package termui renders a wake build to a terminal.
//
// It implements report.Reporter. On a real terminal at Normal verbosity it
// draws a single live line per step with a spinner that is rewritten in place
// to its final status; subprocess output is captured by the executor and shown
// only on failure. At Verbose it streams commands and output live. With no TTY
// (CI, pipes) or NO_COLOR it degrades to plain one-line-per-step output.
package termui

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"golang.org/x/term"

	"github.com/wailsapp/wails/v3/internal/report"
)

const (
	glyphOK     = "✓" // ✓
	glyphFail   = "✗" // ✗
	glyphSkip   = "·" // ·
	glyphArrow  = "›" // ›
	glyphBullet = "•" // •
)

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

type styles struct {
	dim     lipgloss.Style
	accent  lipgloss.Style
	ok      lipgloss.Style
	fail    lipgloss.Style
	cached  lipgloss.Style
	counter lipgloss.Style
	label   lipgloss.Style
	panel   lipgloss.Style
	panelHd lipgloss.Style

	debugVal lipgloss.Style
	debugCat map[string]lipgloss.Style
}

func newStyles(r *lipgloss.Renderer) styles {
	return styles{
		dim:     r.NewStyle().Faint(true),
		accent:  r.NewStyle().Foreground(lipgloss.Color("12")).Bold(true),
		ok:      r.NewStyle().Foreground(lipgloss.Color("10")),
		fail:    r.NewStyle().Foreground(lipgloss.Color("9")).Bold(true),
		cached:  r.NewStyle().Foreground(lipgloss.Color("8")),
		counter: r.NewStyle().Foreground(lipgloss.Color("8")),
		label:   r.NewStyle(),
		panel:   r.NewStyle().Foreground(lipgloss.Color("9")).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("9")).Padding(0, 1),
		panelHd: r.NewStyle().Foreground(lipgloss.Color("9")).Bold(true),

		debugVal: r.NewStyle().Foreground(lipgloss.Color("10")),
		debugCat: map[string]lipgloss.Style{
			"dag":  r.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("13")).Bold(true),
			"dep":  r.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("12")).Bold(true),
			"var":  r.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("14")).Bold(true),
			"exec": r.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("8")).Bold(true),
		},
	}
}

// Reporter renders a build to a terminal.
type Reporter struct {
	w       io.Writer
	level   report.Verbosity
	animate bool // single-line spinner rewrite (TTY + Normal)
	st      styles

	mu          sync.Mutex
	total       int
	idx         int
	stepName    string
	stepLabel   string
	detail      string
	headerShown bool // plain mode: a header line was printed for the current step
	stepStart   time.Time
	frame       int
	ticker      *time.Ticker
	stop        chan struct{}
	dirty       bool // a live line is currently on screen awaiting finalize

	okCount     int
	cachedCount int
	skipCount   int
}

// New builds a Reporter writing to w at the given verbosity. It auto-detects
// whether w is a terminal and whether colour is permitted (NO_COLOR).
func New(w io.Writer, level report.Verbosity) *Reporter {
	isTTY := false
	if f, ok := w.(*os.File); ok {
		isTTY = term.IsTerminal(int(f.Fd()))
	}
	color := isTTY && os.Getenv("NO_COLOR") == ""

	rnd := lipgloss.NewRenderer(w)
	if !color {
		rnd.SetColorProfile(termenv.Ascii)
	}

	return &Reporter{
		w:       w,
		level:   level,
		animate: isTTY && level == report.Normal,
		st:      newStyles(rnd),
	}
}

func (r *Reporter) Level() report.Verbosity { return r.level }

func (r *Reporter) BuildStart(verb, target string, total int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.total = total
	if r.level == report.Silent {
		return
	}
	plural := "steps"
	if total == 1 {
		plural = "step"
	}
	count := r.st.dim.Render(fmt.Sprintf("(%d %s)", total, plural))
	if verb == "" {
		fmt.Fprintf(r.w, "\n  %s %s\n\n", r.st.label.Render(target), count)
		return
	}
	fmt.Fprintf(r.w, "\n  %s %s %s %s\n\n",
		r.st.accent.Render(verb),
		r.st.dim.Render(glyphArrow),
		r.st.label.Render(target),
		count,
	)
}

func (r *Reporter) StepStart(name, label string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.idx++
	r.stepName = name
	r.stepLabel = label
	r.detail = ""
	r.stepStart = time.Now()
	if r.level == report.Silent {
		return
	}
	if r.animate {
		r.frame = 0
		r.dirty = true
		r.paintLocked()
		r.startTickerLocked()
		return
	}
	// Non-animated modes print the header lazily: only once a command, output,
	// or info line actually arrives to sit beneath it. A step with nothing to
	// show (e.g. a cache hit) collapses to its single final status line.
	r.headerShown = false
}

func (r *Reporter) printHeaderLocked() {
	fmt.Fprintf(r.w, "  %s %s %s\n", r.st.dim.Render(glyphArrow), r.counterStr(), r.title())
	r.headerShown = true
}

func (r *Reporter) StepInfo(msg string) {
	msg = strings.TrimSpace(msg)
	if msg == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.level == report.Silent {
		return
	}
	if r.animate {
		r.detail = msg
		r.paintLocked()
		return
	}
	if !r.headerShown {
		r.printHeaderLocked()
	}
	fmt.Fprintf(r.w, "      %s %s\n", r.st.dim.Render(glyphBullet), r.st.dim.Render(msg))
}

func (r *Reporter) StepCommand(cmd string) {
	if r.level < report.Verbose {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.headerShown {
		r.printHeaderLocked()
	}
	fmt.Fprintf(r.w, "      %s %s\n", r.st.dim.Render("$"), r.st.dim.Render(cmd))
}

func (r *Reporter) StepOutput(line string) {
	if r.level < report.Verbose {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.headerShown {
		r.printHeaderLocked()
	}
	fmt.Fprintf(r.w, "      %s\n", strings.TrimRight(line, "\n"))
}

func (r *Reporter) StepEnd(status report.Status, dur time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.stopTickerLocked()
	switch status {
	case report.StatusCached:
		r.cachedCount++
	case report.StatusSkipped:
		r.skipCount++
	default:
		r.okCount++
	}
	if r.level == report.Silent {
		return
	}

	var glyph, suffix string
	switch status {
	case report.StatusCached:
		glyph = r.st.cached.Render(glyphSkip)
		suffix = r.st.cached.Render(" cached")
	case report.StatusSkipped:
		glyph = r.st.cached.Render(glyphSkip)
		suffix = r.st.cached.Render(" skipped")
	default:
		glyph = r.st.ok.Render(glyphOK)
		suffix = "  " + r.st.dim.Render(fmtDur(dur))
	}

	line := fmt.Sprintf("  %s %s %s%s", glyph, r.counterStr(), r.title(), suffix)
	r.finalizeLineLocked(line)
}

func (r *Reporter) StepFailed(f report.Failure) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.stopTickerLocked()

	// A failure always surfaces, even at Silent verbosity.
	line := fmt.Sprintf("  %s %s %s", r.st.fail.Render(glyphFail), r.counterStr(), r.title())
	r.finalizeLineLocked(line)
	r.renderFailureLocked(f)
}

func (r *Reporter) BuildEnd(dur time.Duration, ok bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.stopTickerLocked()
	if r.level == report.Silent {
		return
	}

	var head string
	if ok {
		ran := r.okCount
		extra := ""
		if r.cachedCount > 0 {
			extra += fmt.Sprintf(", %d cached", r.cachedCount)
		}
		if r.skipCount > 0 {
			extra += fmt.Sprintf(", %d skipped", r.skipCount)
		}
		head = fmt.Sprintf("  %s %s %s",
			r.st.ok.Render(glyphOK),
			r.st.ok.Render("build succeeded"),
			r.st.dim.Render(fmt.Sprintf("in %s (%d ran%s)", fmtDur(dur), ran, extra)),
		)
	} else {
		head = fmt.Sprintf("  %s %s %s",
			r.st.fail.Render(glyphFail),
			r.st.fail.Render("build failed"),
			r.st.dim.Render("after "+fmtDur(dur)),
		)
	}
	fmt.Fprintf(r.w, "\n%s\n\n", head)
}

// Artifact appends one output line after the build summary. termui's
// rendering is deliberately plain — pulse provides the rich version.
func (r *Reporter) Artifact(a report.Artifact) {
	if a.Path == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.level == report.Silent {
		return
	}
	size := ""
	if a.Size > 0 {
		size = "  " + r.st.dim.Render(humanSize(a.Size))
	}
	kind := ""
	if a.Kind != "" {
		kind = "  " + r.st.dim.Render(a.Kind)
	}
	fmt.Fprintf(r.w, "  %s %s%s%s\n", r.st.dim.Render("·"), a.Path, kind, size)
}

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
	style, ok := r.st.debugCat[d.Category]
	if !ok {
		style = r.st.debugCat["exec"]
	}
	badge := style.Render(fmt.Sprintf(" %-4s ", cat))

	var b strings.Builder
	fmt.Fprintf(&b, "  %s %s", badge, d.Subject)
	if d.Arrow != "" {
		fmt.Fprintf(&b, " %s%s", r.st.dim.Render("→ "), r.st.dim.Render(d.Arrow))
	}
	for _, f := range d.Fields {
		fmt.Fprintf(&b, "  %s%s", r.st.dim.Render(f.Key+"="), r.st.debugVal.Render(f.Val))
	}
	fmt.Fprintln(r.w, b.String())
}

// --- internals ----------------------------------------------------------

func (r *Reporter) counterStr() string {
	return r.st.counter.Render(fmt.Sprintf("[%d/%d]", r.idx, r.total))
}

// title is the human string for the in-flight step: the author's label when set
// (labels conventionally embed the name plus context, e.g.
// "build:frontend (PRODUCTION=false)"), otherwise the canonical task name.
func (r *Reporter) title() string {
	if r.stepLabel != "" {
		return r.st.label.Render(r.stepLabel)
	}
	return r.st.label.Render(displayName(r.stepName))
}

// displayName drops the leading platform/common namespace from a task name for
// display: during a darwin build, "darwin:common:go:mod:tidy" reads better as
// "go:mod:tidy". The canonical name is still used everywhere else.
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

// paintLocked draws the in-flight spinner line in place. Caller holds r.mu.
func (r *Reporter) paintLocked() {
	if !r.animate {
		return
	}
	spin := r.st.accent.Render(spinnerFrames[r.frame%len(spinnerFrames)])
	line := fmt.Sprintf("  %s %s %s", spin, r.counterStr(), r.title())
	if r.detail != "" {
		line += "  " + r.st.dim.Render(glyphBullet+" "+r.detail)
	}
	fmt.Fprintf(r.w, "\r\x1b[K%s", line)
}

// finalizeLineLocked replaces any live line with the final one and ends it with
// a newline. Caller holds r.mu.
func (r *Reporter) finalizeLineLocked(line string) {
	if r.animate && r.dirty {
		fmt.Fprintf(r.w, "\r\x1b[K%s\n", line)
	} else {
		fmt.Fprintf(r.w, "%s\n", line)
	}
	r.dirty = false
}

func (r *Reporter) startTickerLocked() {
	r.stop = make(chan struct{})
	r.ticker = time.NewTicker(90 * time.Millisecond)
	stop := r.stop
	tick := r.ticker
	go func() {
		for {
			select {
			case <-stop:
				return
			case <-tick.C:
				r.mu.Lock()
				if r.dirty {
					r.frame++
					r.paintLocked()
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

func (r *Reporter) renderFailureLocked(f report.Failure) {
	var b strings.Builder
	if f.Command != "" {
		fmt.Fprintf(&b, "%s %s\n", r.st.panelHd.Render("command"), f.Command)
	}
	if f.ExitCode != 0 {
		fmt.Fprintf(&b, "%s exited %d\n", r.st.panelHd.Render("status "), f.ExitCode)
	} else if f.Err != nil {
		fmt.Fprintf(&b, "%s %s\n", r.st.panelHd.Render("error  "), f.Err.Error())
	}
	out := strings.TrimRight(f.Output, "\n")
	if out != "" {
		b.WriteString("\n")
		b.WriteString(tail(out, 20))
	}
	body := strings.TrimRight(b.String(), "\n")
	fmt.Fprintf(r.w, "\n%s\n", r.st.panel.Render(body))
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

func tail(s string, n int) string {
	lines := strings.Split(s, "\n")
	if len(lines) <= n {
		return s
	}
	return "..." + "\n" + strings.Join(lines[len(lines)-n:], "\n")
}
