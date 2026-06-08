package tui

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gdamore/tcell/v2"

	"github.com/atterpac/dado/components"

	monitor "github.com/wailsapp/wails/v3/internal/monitor"
)

// windowsView is a master/detail dashboard of live windows: the list on the
// left, and on the right the selected window's geometry, screen/DPI info, and
// the IPC activity it has generated (calls/events, by the trace window field).
//
// While focused it auto-polls the app for a fresh snapshot so resize/move
// updates appear live.
type windowsView struct {
	*components.ComponentBase
	m *Model

	split   *components.Split
	table   *components.Table
	detail  *components.Label
	visible []monitor.WindowInfo
	polling atomic.Bool
}

func newWindowsView(m *Model) *windowsView {
	v := &windowsView{m: m}

	v.table = components.NewTable()
	v.table.SetHeaders("ID", "NAME", "SIZE", "STATE")
	v.table.ConfigureEmpty("□", "No windows", "Snapshot not yet received")

	v.detail = components.NewLabel("").SetDynamicColors(true).SetScrollable(true)

	tablePanel := components.NewPanel().SetTitle("Windows").SetContent(v.table).SetFocused(true)
	detailPanel := components.NewPanel().SetTitle("Activity").SetContent(v.detail)

	v.split = components.NewSplit()
	v.split.SetDirection(components.SplitHorizontal).SetRatio(0.4).
		SetLeft(tablePanel).SetRight(detailPanel)

	v.table.SetSelectionChangedFunc(func(row, col int) { v.showDetail(row) })

	v.ComponentBase = components.NewComponentBase(v.split).
		SetName("Windows").
		AddHint("j/k", "Move").
		AddHint("r", "Refresh").
		AddHint("q", "Back").
		SetOnStart(v.start).
		SetOnStop(v.stop)
	v.ComponentBase.SetInputHandler(func(ev *tcell.EventKey) bool {
		if ev.Key() == tcell.KeyRune && ev.Rune() == 'r' {
			v.m.requestSnapshot()
			return true
		}
		if ev.Key() == tcell.KeyEnter {
			idx := v.table.SelectedRow()
			if idx >= 0 && idx < len(v.visible) {
				v.m.jumpToWindowStream(v.visible[idx].Name)
			}
			return true
		}
		return vimTableNav(v.table, ev)
	})
	return v
}

// start kicks off auto-polling while the view is active, so resize/move shows
// up live. The goroutine exits when stop() flips the flag (on view pop).
func (v *windowsView) start() {
	v.rebuild()
	if v.polling.Swap(true) {
		return // already polling
	}
	go func() {
		for v.polling.Load() {
			time.Sleep(time.Second)
			if !v.polling.Load() {
				return
			}
			v.m.requestSnapshot()
		}
	}()
}

func (v *windowsView) stop() {
	v.polling.Store(false)
}

func (v *windowsView) rebuild() {
	v.visible = nil
	if v.m.snapshot != nil {
		v.visible = v.m.snapshot.Windows
	}

	sel := v.table.SelectedRow()
	v.table.ClearRows()
	for _, w := range v.visible {
		state := windowState(w)
		v.table.AddRow(
			fmt.Sprintf("%d", w.ID),
			w.Name,
			fmt.Sprintf("%d×%d", w.Width, w.Height),
			state,
		)
	}
	if sel >= 0 && sel < len(v.visible) {
		v.table.SelectRow(sel)
		v.showDetail(sel + 1)
	} else if len(v.visible) > 0 {
		v.table.SelectRow(0)
		v.showDetail(1)
	} else {
		v.detail.SetText("")
	}
}

func (v *windowsView) showDetail(tableRow int) {
	idx := tableRow - 1
	if idx < 0 || idx >= len(v.visible) {
		v.detail.SetText("")
		return
	}
	w := v.visible[idx]

	var b strings.Builder
	fmt.Fprintf(&b, "[::b]%s[-:-:-]  [gray]id %d[-]\n\n", escapeTags(orDash(w.Name)), w.ID)

	row(&b, "size", fmt.Sprintf("%d×%d", w.Width, w.Height))
	row(&b, "border", fmt.Sprintf("L%d R%d T%d B%d", w.Border.Left, w.Border.Right, w.Border.Top, w.Border.Bottom))
	row(&b, "position", fmt.Sprintf("%d,%d [gray]abs[-]  ·  %d,%d [gray]rel[-]", w.X, w.Y, w.RelX, w.RelY))
	row(&b, "state", windowState(w))
	row(&b, "zoom", fmt.Sprintf("%.0f%%", w.Zoom*100))

	if w.Screen != nil {
		s := w.Screen
		primary := ""
		if s.IsPrimary {
			primary = " [gray](primary)[-]"
		}
		b.WriteString("\n[::b]Screen[-:-:-]\n")
		row(&b, "name", s.Name+primary)
		row(&b, "scale", fmt.Sprintf("%.2gx  (%.0f DPI)", float64(s.ScaleFactor), float64(s.ScaleFactor)*96))
		row(&b, "bounds", fmt.Sprintf("%d,%d %d×%d", s.BoundsX, s.BoundsY, s.BoundsW, s.BoundsH))
		row(&b, "work area", fmt.Sprintf("%d×%d", s.WorkW, s.WorkH))
		if s.Rotation != 0 {
			row(&b, "rotation", fmt.Sprintf("%.0f°", float64(s.Rotation)))
		}
	}

	// Per-window IPC activity (derived from the trace stream).
	act := v.m.activityForWindow(w.Name)
	b.WriteString("\n[::b]Activity[-:-:-] [gray](this window)[-]\n")
	fmt.Fprintf(&b, "  [gray]%d calls · %d err · %d events[-]\n", act.calls, act.errs, act.events)
	if len(act.topMethods) > 0 {
		b.WriteString("  [gray]top:[-] ")
		parts := make([]string, 0, len(act.topMethods))
		for _, mc := range act.topMethods {
			parts = append(parts, fmt.Sprintf("%s ×%d", shortMethod(mc.method), mc.count))
		}
		b.WriteString(escapeTags(strings.Join(parts, ", ")) + "\n")
	}
	if len(act.recent) > 0 {
		b.WriteString("\n  [gray]recent:[-]\n")
		for i := len(act.recent) - 1; i >= 0; i-- {
			r := act.recent[i]
			fmt.Fprintf(&b, "    %s %s %s\n",
				r.t.Format("15:04:05"), statusGlyph(r), escapeTags(shortMethod(r.method)))
		}
	}

	v.detail.SetText(strings.TrimRight(b.String(), "\n"))
}

// shortMethod trims a fully-qualified name to its final segment for display.
func shortMethod(fqn string) string {
	if i := strings.LastIndex(fqn, "."); i >= 0 && i < len(fqn)-1 {
		return fqn[i+1:]
	}
	return fqn
}

func windowState(w monitor.WindowInfo) string {
	parts := []string{}
	if w.Focused {
		parts = append(parts, "focused")
	}
	if w.Fullscreen {
		parts = append(parts, "fullscreen")
	}
	if w.Maximised {
		parts = append(parts, "max")
	}
	if w.Minimised {
		parts = append(parts, "min")
	}
	if w.Resizable {
		parts = append(parts, "resizable")
	}
	if w.IgnoreMouse {
		parts = append(parts, "click-through")
	}
	if len(parts) == 0 {
		return "—"
	}
	return strings.Join(parts, " ")
}
