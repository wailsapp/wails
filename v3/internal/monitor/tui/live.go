package tui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"

	"github.com/atterpac/dado/components"
	"github.com/atterpac/dado/core"
	"github.com/atterpac/dado/theme"
)

// liveView is the main monitor screen: a table of binding calls on the left and
// a JSON detail pane on the right. Events are NOT shown here (see eventsView).
type liveView struct {
	*components.ComponentBase
	m *Model

	split  *components.Split
	table  *components.Table
	detail *components.Label
	body   *components.Label
}

func newLiveView(m *Model) *liveView {
	v := &liveView{m: m}

	v.table = components.NewTable()
	v.table.SetHeaders("TIME", "", "METHOD", "WINDOW", "DUR", "ST")
	v.table.ConfigureEmpty("…", "Waiting for binding calls",
		"Call a bound method from the frontend. Press E for events.")

	v.detail = components.NewLabel("").SetDynamicColors(true)
	v.body = components.NewLabel("").SetDynamicColors(true).SetScrollable(true)

	detailPane := core.NewFlex().SetDirection(core.Column).
		AddItem(v.detail, 6, 0, false).
		AddItem(v.body, 0, 1, false)

	tablePanel := components.NewPanel().SetTitle("Calls").SetContent(v.table).SetFocused(true)
	detailPanel := components.NewPanel().SetTitle("Detail").SetContent(detailPane)

	v.split = components.NewSplit()
	v.split.SetDirection(components.SplitHorizontal).SetRatio(0.55).
		SetLeft(tablePanel).SetRight(detailPanel)

	v.table.SetSelectionChangedFunc(func(row, col int) { v.showDetailForRow(row) })

	v.ComponentBase = components.NewComponentBase(v.split).
		SetName("Live").
		AddHint("j/k", "Move").
		AddHint("G", "Follow tail").
		AddHint("Enter", "Detail").
		AddHint("y", "Copy").
		AddHint("/", "Filter").
		AddHint("E", "Events").
		AddHint("s", "Stats").
		AddHint("w", "Windows").
		AddHint("b", "Bindings").
		AddHint("i", "Info").
		AddHint("t", "Timeline").
		AddHint("f", "Follow").
		AddHint("p", "Pause").
		AddHint("e", "Errors").
		AddHint("c", "Clear").
		AddHint("?", "Help").
		AddHint("T", "Theme").
		AddHint("q", "Quit")
	v.ComponentBase.SetInputHandler(v.handleKey)
	return v
}

func (v *liveView) handleKey(ev *tcell.EventKey) bool {
	if ev.Key() == tcell.KeyEnter {
		v.m.showDetailModal(v.selectedRecord())
		return true
	}
	if ev.Key() == tcell.KeyRune && ev.Rune() == 'y' {
		v.m.copyRecord(v.selectedRecord())
		return true
	}
	// Scrolling up means the user wants to inspect history: stop follow so new
	// traces don't yank the cursor to the bottom. `G` / `f` re-enable it.
	if isScrollUpKey(ev) && v.m.follow {
		v.m.follow = false
		v.m.updateStatus()
	}
	if ev.Key() == tcell.KeyRune && ev.Rune() == 'G' {
		v.m.follow = true
		v.selectLast()
		v.m.updateStatus()
		return true
	}
	return vimTableNav(v.table, ev)
}

// isScrollUpKey reports whether the event moves the selection toward older rows.
func isScrollUpKey(ev *tcell.EventKey) bool {
	switch ev.Key() {
	case tcell.KeyUp, tcell.KeyPgUp, tcell.KeyHome:
		return true
	case tcell.KeyRune:
		return ev.Rune() == 'k' || ev.Rune() == 'g'
	}
	return false
}

func (v *liveView) selectedRecord() *record {
	idx := v.table.SelectedRow()
	v.m.mu.Lock()
	defer v.m.mu.Unlock()
	if idx < 0 || idx >= len(v.m.visibleCalls) {
		return nil
	}
	return v.m.visibleCalls[idx]
}

// rebuild repopulates the table from the filtered call snapshot (UI thread).
// When following, the cursor tracks the newest row; otherwise it stays pinned
// to whatever record (by seq) the user had selected, so incoming traces don't
// move it out from under them.
func (v *liveView) rebuild() {
	var keepSeq uint64
	if !v.m.follow {
		if r := v.selectedRecord(); r != nil {
			keepSeq = r.seq
		}
	}

	recs := v.m.snapshotCalls()
	v.table.ClearRows()
	keepIdx := -1
	for i, r := range recs {
		if r.seq == keepSeq {
			keepIdx = i
		}
		v.table.AddRowWithColor(rowColor(r),
			r.t.Format("15:04:05.000"),
			dirArrow(r.dir),
			r.method,
			r.window,
			durStr(r),
			statusGlyph(r),
		)
	}

	if v.m.follow {
		v.selectLast()
	} else if keepIdx >= 0 {
		v.table.SelectRow(keepIdx)
	}
}

func rowColor(r *record) tcell.Color {
	switch r.status {
	case "error":
		return theme.Error()
	case "pending":
		return theme.FgDim()
	case "cancelled":
		return theme.Warning()
	}
	return theme.Fg()
}

func (v *liveView) selectLast() {
	n := v.table.GetDataRowCount()
	if n > 0 {
		v.table.SelectRow(n - 1)
		v.showDetailForRow(n) // table row = header + (n-1) data index
	}
}

func (v *liveView) showDetailForRow(tableRow int) {
	idx := tableRow - 1
	v.m.mu.Lock()
	recs := v.m.visibleCalls
	v.m.mu.Unlock()
	if idx < 0 || idx >= len(recs) {
		v.detail.SetText("")
		v.body.SetText("")
		return
	}
	r := recs[idx]
	v.detail.SetText(callHeader(r))
	v.body.SetText(callBody(r))
}

// statusBadge renders a colored status pill for the detail header.
func statusBadge(r *record) string {
	switch r.status {
	case "ok":
		return "[green]✓ ok[-]"
	case "error":
		return "[red]✗ error[-]"
	case "pending":
		return "[yellow]⏳ pending[-]"
	case "cancelled":
		return "[gray]⊘ cancelled[-]"
	default:
		return ""
	}
}

// callHeader renders the metadata block above the JSON body.
func callHeader(r *record) string {
	var b strings.Builder
	fmt.Fprintf(&b, "[::b]%s[-:-:-]  %s\n", escapeTags(methodOr(r)), statusBadge(r))
	fmt.Fprintf(&b, "[gray]window[-] %s   [gray]took[-] %s\n", orDash(r.window), orDash(durStr(r)))
	fmt.Fprintf(&b, "[gray]call id[-] %s", orDash(r.callID))
	if r.errMsg != "" {
		fmt.Fprintf(&b, "\n[red]✗ %s[-]", escapeTags(r.errMsg))
		if r.errKind != "" {
			fmt.Fprintf(&b, " [gray](%s)[-]", escapeTags(r.errKind))
		}
	}
	b.WriteString("\n[gray]enter expand · y copy[-]")
	return b.String()
}

// callBody renders the args/result with clear Input/Output sections. Shown in a
// JSON CodeView, so dividers are written as comments.
func callBody(r *record) string {
	var b strings.Builder
	b.WriteString(section("Input"))
	if len(r.args) > 0 {
		b.WriteString(prettyJSON(r.args))
	} else {
		b.WriteString("// (no arguments)")
	}

	b.WriteString("\n\n")
	switch r.status {
	case "error":
		b.WriteString(section("Error"))
		if r.errMsg != "" {
			b.WriteString("// " + strings.ReplaceAll(r.errMsg, "\n", "\n// "))
		} else {
			b.WriteString("// (call failed)")
		}
	case "pending":
		b.WriteString(section("Output"))
		b.WriteString("// awaiting response…")
	default:
		b.WriteString(section("Output"))
		if len(r.result) > 0 {
			b.WriteString(prettyJSON(r.result))
		} else {
			b.WriteString("// (no return value)")
		}
	}
	return b.String()
}

// section is a comment-style label for the JSON body.
func section(label string) string {
	return "// " + label + "\n"
}
