package tui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"

	"github.com/atterpac/dado/components"
	"github.com/atterpac/dado/core"
	"github.com/atterpac/dado/theme"
	mon "github.com/wailsapp/wails/v3/internal/monitor"
)

// graphPoints caps how many samples feed the braille charts.
const graphPoints = 240

// timelineView: three stacked graphs on the left (RAM, CPU, calls/s), detail
// label + activity table on the right. h/l scrubs all graphs together;
// Enter shifts focus to the call list; q/Esc returns to the graphs.
type timelineView struct {
	*components.ComponentBase
	m *Model

	// outer is a horizontal split: graphs left, panel right.
	outer *core.Flex

	// three auto-scaled graphs, stacked vertically on the left.
	ramGraph   *components.LineGraph
	cpuGraph   *components.LineGraph
	callGraph  *components.LineGraph
	graphsFlex *core.Flex

	// right panel: detail label on top, activity table on bottom.
	rightFlex *core.Flex
	detail    *components.Label
	callTable *components.Table

	// which side has focus: false = graphs, true = call list.
	listFocused bool

	// cursor index from end (0=latest, -1=follow).
	cursor int

	// records (calls + events) in the currently selected bucket.
	bucketRecs []*record
}

func newLineGraph(title string) *components.LineGraph {
	return components.NewLineGraph().
		SetTitle(title).
		SetStyle(components.LineGraphSolid).
		SetShowGrid(true).
		SetShowLegend(true).
		SetAutoScale(true)
}

func newTimelineView(m *Model) *timelineView {
	v := &timelineView{m: m, cursor: -1}

	v.ramGraph = newLineGraph("RAM (MiB)")
	v.cpuGraph = newLineGraph("CPU (%)")
	v.callGraph = newLineGraph("calls/s")

	v.graphsFlex = core.NewFlex().SetDirection(core.Column).
		AddItem(v.ramGraph, 0, 1, false).
		AddItem(v.cpuGraph, 0, 1, false).
		AddItem(v.callGraph, 0, 1, false)

	v.detail = components.NewLabel("").SetDynamicColors(true).SetScrollable(true)

	v.callTable = components.NewTable()
	v.callTable.SetHeaders("TIME", "KIND", "METHOD", "DUR/COUNT", "ST")
	v.callTable.ConfigureEmpty("—", "No activity in bucket", "Scrub to a busier window.")

	v.rightFlex = core.NewFlex().SetDirection(core.Column).
		AddItem(v.detail, 0, 2, false).
		AddItem(v.callTable, 0, 3, false)

	graphsPanel := components.NewPanel().SetTitle("Resources").SetContent(v.graphsFlex).SetFocused(true)
	activityPanel := components.NewPanel().SetTitle("Activity").SetContent(v.rightFlex)

	v.outer = core.NewFlex().SetDirection(core.Row).
		AddItem(graphsPanel, 0, 3, true).
		AddItem(activityPanel, 0, 2, false)

	v.ComponentBase = components.NewComponentBase(v.outer).
		SetName("Timeline").
		AddHint("h/l", "Scrub").
		AddHint("G", "Latest").
		AddHint("Enter", "Focus list / Detail").
		AddHint("q", "Back")
	v.ComponentBase.SetInputHandler(v.handleKey)
	return v
}

func (v *timelineView) handleKey(ev *tcell.EventKey) bool {
	n := len(v.m.snapshotSamples())

	// When the call list is focused, Enter opens detail; q/Esc returns to graphs.
	if v.listFocused {
		switch {
		case ev.Key() == tcell.KeyEscape || (ev.Key() == tcell.KeyRune && ev.Rune() == 'q'):
			v.listFocused = false
			v.m.app.SetFocus(v.graphsFlex)
			return true
		case ev.Key() == tcell.KeyEnter:
			idx := v.callTable.SelectedRow() - 1
			if idx >= 0 && idx < len(v.bucketRecs) {
				v.m.showDetailModal(v.bucketRecs[idx])
			}
			return true
		}
		return false
	}

	switch ev.Key() {
	case tcell.KeyLeft:
		v.moveCursor(1, n)
		return true
	case tcell.KeyRight:
		v.moveCursor(-1, n)
		return true
	case tcell.KeyEnter:
		v.listFocused = true
		v.m.app.SetFocus(v.callTable)
		return true
	case tcell.KeyRune:
		switch ev.Rune() {
		case 'h':
			v.moveCursor(1, n)
			return true
		case 'l':
			v.moveCursor(-1, n)
			return true
		case 'G':
			v.cursor = -1
			v.rebuild()
			return true
		}
	}
	return false
}

func (v *timelineView) moveCursor(delta, n int) {
	if n == 0 {
		return
	}
	cur := v.cursor
	if cur < 0 {
		cur = 0
	}
	cur += delta
	if cur < 0 {
		cur = 0
	}
	if cur > n-1 {
		cur = n - 1
	}
	v.cursor = cur
	v.rebuild()
}

func (v *timelineView) selectedIndex(n int) int {
	if v.cursor < 0 {
		return n - 1
	}
	idx := n - 1 - v.cursor
	if idx < 0 {
		idx = 0
	}
	return idx
}

func (v *timelineView) rebuild() {
	samples := v.m.snapshotSamples()
	if len(samples) == 0 {
		v.ramGraph.SetSeries()
		v.cpuGraph.SetSeries()
		v.callGraph.SetSeries()
		v.detail.SetText("\n  [gray]Waiting for resource samples…[-]\n  Samples arrive ~1/s once the app is running.")
		v.callTable.ClearRows()
		return
	}
	v.rebuildGraphs(samples)
	v.rebuildBottom(samples)
}

// rebuildGraphs feeds all three LineGraphs and syncs the cursor across them.
func (v *timelineView) rebuildGraphs(samples []mon.Sample) {
	view := tailN(samples, graphPoints)
	ram := make([]float64, len(view))
	cpu := make([]float64, len(view))
	for i, s := range view {
		ram[i] = float64(s.RSS) / (1024 * 1024)
		cpu[i] = s.CPUPct
	}
	rate := callRate(view, v.m.callsInRange)
	rateF := make([]float64, len(rate))
	for i, r := range rate {
		rateF[i] = float64(r)
	}

	// Cursor position within the view window.
	pos := len(view) - 1
	if v.cursor >= 0 {
		full := len(samples)
		idx := v.selectedIndex(full)
		start := full - len(view)
		pos = idx - start
		if pos < 0 {
			pos = 0
		}
		if pos > len(view)-1 {
			pos = len(view) - 1
		}
	}
	frac := 1.0
	if len(view) > 1 {
		frac = float64(pos) / float64(len(view)-1)
	}

	cur := view[pos]
	curRate := 0
	if pos < len(rate) {
		curRate = rate[pos]
	}

	v.ramGraph.SetSeries(components.DataSeries{
		Label:  fmt.Sprintf("%.1f MiB", float64(cur.RSS)/(1024*1024)),
		Values: ram,
		Color:  theme.Accent(),
	})
	v.cpuGraph.SetSeries(components.DataSeries{
		Label:  fmt.Sprintf("%.1f%%", cur.CPUPct),
		Values: cpu,
		Color:  theme.Warning(),
	})
	v.callGraph.SetSeries(components.DataSeries{
		Label:  fmt.Sprintf("%d/s", curRate),
		Values: rateF,
		Color:  theme.Success(),
	})

	// Build cursor card shown on the RAM graph (top). Shows all three metrics +
	// call/event counts for the selected bucket so the card is self-contained.
	var fromBucket time.Time
	if pos == 0 {
		fromBucket = cur.Time.Add(-time.Second)
	} else {
		fromBucket = view[pos-1].Time
	}
	_, evts, errs := v.m.callsInRange(fromBucket, cur.Time)
	tag := ""
	if v.cursor < 0 {
		tag = " (latest)"
	}
	card := []string{
		cur.Time.Format("15:04:05") + tag,
		fmt.Sprintf("RAM  %.1f MiB", float64(cur.RSS)/(1024*1024)),
		fmt.Sprintf("CPU  %.1f%%", cur.CPUPct),
		fmt.Sprintf("calls %d/s  ev %d  err %d", curRate, evts, errs),
	}
	if errs > 0 {
		card[3] = fmt.Sprintf("calls %d/s  ev %d  [red]err %d[-]", curRate, evts, errs)
	}

	v.ramGraph.SetCursorFrac(frac)
	v.ramGraph.SetCursorLabel(card)
	v.cpuGraph.SetCursorFrac(frac)
	v.cpuGraph.SetCursorLabel(nil)
	v.callGraph.SetCursorFrac(frac)
	v.callGraph.SetCursorLabel(nil)
}

// rebuildBottom populates the detail label and call table for the selected bucket.
func (v *timelineView) rebuildBottom(samples []mon.Sample) {
	n := len(samples)
	idx := v.selectedIndex(n)
	s := samples[idx]

	var from time.Time
	if idx == 0 {
		from = s.Time.Add(-time.Second)
	} else {
		from = samples[idx-1].Time
	}
	calls, events, errs := v.m.callsInRange(from, s.Time)
	methods := v.m.topMethodsInRange(from, s.Time, 5)

	accent := theme.TagAccent()
	var b strings.Builder
	fmt.Fprintf(&b, "[%s::b]%s[-:-:-]", accent, s.Time.Format("15:04:05"))
	if v.cursor < 0 {
		b.WriteString("  [gray](latest)[-]")
	}
	b.WriteString("\n\n")
	b.WriteString(kv("RAM", humanBytes(s.RSS)))
	b.WriteString(kv("CPU", fmt.Sprintf("%.1f%%", s.CPUPct)))
	b.WriteString(kv("Heap", humanBytes(s.HeapAlloc)))
	if s.Goroutines > 0 {
		b.WriteString(kv("Goroutines", fmt.Sprintf("%d", s.Goroutines)))
	}
	b.WriteString("\n")
	b.WriteString(kv("Calls", fmt.Sprintf("%d", calls)))
	if errs > 0 {
		b.WriteString(kv("Errors", fmt.Sprintf("[red]%d[-]", errs)))
	}
	b.WriteString(kv("Events", fmt.Sprintf("%d", events)))
	if len(methods) > 0 {
		b.WriteString("\n")
		fmt.Fprintf(&b, "[%s::b]Top methods[-:-:-]\n", accent)
		for _, mc := range methods {
			fmt.Fprintf(&b, "  [gray]×%d[-]  %s\n", mc.count, escapeTags(mc.method))
		}
	}
	b.WriteString("\n[gray]Enter → list  (Enter again → detail)[-]")
	v.detail.SetText(b.String())

	// Populate table with calls and events in [from, s.Time], sorted by time.
	v.m.mu.Lock()
	v.bucketRecs = nil
	for _, r := range v.m.callRecs {
		if !r.t.Before(from) && !r.t.After(s.Time) {
			v.bucketRecs = append(v.bucketRecs, r)
		}
	}
	for _, r := range v.m.eventRecs {
		if !r.t.Before(from) && !r.t.After(s.Time) {
			v.bucketRecs = append(v.bucketRecs, r)
		}
	}
	v.m.mu.Unlock()
	sort.Slice(v.bucketRecs, func(i, j int) bool {
		return v.bucketRecs[i].t.Before(v.bucketRecs[j].t)
	})

	v.callTable.ClearRows()
	for _, r := range v.bucketRecs {
		var kind, extra, st string
		switch r.kind {
		case "event":
			kind = "evt"
			if r.count > 1 {
				extra = fmt.Sprintf("×%d", r.count)
			}
		default:
			kind = "call"
			extra = durStr(r)
			st = statusGlyph(r)
		}
		v.callTable.AddRowWithColor(rowColor(r),
			r.t.Format("15:04:05.000"),
			kind,
			shortMethod(r.method),
			extra,
			st,
		)
	}
}

// --- shared helpers ---

func callRate(samples []mon.Sample, calls func(from, to time.Time) (int, int, int)) []int {
	rate := make([]int, len(samples))
	for i := range samples {
		var from time.Time
		if i == 0 {
			from = samples[i].Time.Add(-time.Second)
		} else {
			from = samples[i-1].Time
		}
		c, _, _ := calls(from, samples[i].Time)
		rate[i] = c
	}
	return rate
}

func tailN(s []mon.Sample, n int) []mon.Sample {
	if len(s) <= n {
		return s
	}
	return s[len(s)-n:]
}

func humanBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

func kv(k, val string) string {
	return fmt.Sprintf("  [gray]%-12s[-] %s\n", k, val)
}

// --- Model methods ---

const maxSamples = 1024

func (m *Model) applySample(s mon.Sample) {
	m.mu.Lock()
	trimmed := len(m.samples) >= maxSamples
	m.samples = append(m.samples, s)
	if trimmed {
		m.samples = m.samples[len(m.samples)-maxSamples:]
	}
	n := len(m.samples)
	m.mu.Unlock()
	if m.timeline != nil && m.app.Pages().Current() == m.timeline {
		// When the user has pinned the cursor, advance it by one to keep the same
		// sample selected as new data arrives. Don't advance past the end.
		if m.timeline.cursor >= 0 && !trimmed {
			m.timeline.cursor++
			if m.timeline.cursor > n-1 {
				m.timeline.cursor = n - 1
			}
		}
		m.timeline.rebuild()
	}
	m.updateStatus()
}

func (m *Model) snapshotSamples() []mon.Sample {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]mon.Sample, len(m.samples))
	copy(out, m.samples)
	return out
}

func (m *Model) callsInRange(from, to time.Time) (calls, events, errs int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, r := range m.callRecs {
		if !r.t.Before(from) && !r.t.After(to) {
			calls++
			if r.status == "error" {
				errs++
			}
		}
	}
	for _, r := range m.eventRecs {
		if !r.t.Before(from) && !r.t.After(to) {
			events += maxInt(r.count, 1)
		}
	}
	return
}

func (m *Model) showTimeline() {
	if m.app.Pages().Current() == m.timeline && m.app.Pages().CanPop() {
		m.app.Pages().Pop()
		return
	}
	if m.timeline == nil {
		m.timeline = newTimelineView(m)
	}
	m.timeline.rebuild()
	m.app.Pages().Push(m.timeline)
}

func (m *Model) topMethodsInRange(from, to time.Time, limit int) []methodCount {
	m.mu.Lock()
	tally := map[string]int{}
	for _, r := range m.callRecs {
		if !r.t.Before(from) && !r.t.After(to) {
			tally[r.method]++
		}
	}
	m.mu.Unlock()
	out := make([]methodCount, 0, len(tally))
	for k, c := range tally {
		out = append(out, methodCount{method: k, count: c})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].count > out[j].count })
	if len(out) > limit {
		out = out[:limit]
	}
	return out
}
