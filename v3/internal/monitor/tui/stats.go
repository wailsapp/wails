package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"

	"github.com/atterpac/dado/components"
)

// statsView shows aggregate IPC metrics: top methods by call count, with error
// counts and average latency.
type statsView struct {
	*components.ComponentBase
	m     *Model
	table *components.Table
}

func newStatsView(m *Model) *statsView {
	v := &statsView{m: m}

	v.table = components.NewTable()
	v.table.SetHeaders("METHOD", "CALLS", "ERRORS", "AVG", "ERR%")
	v.table.ConfigureEmpty("∅", "No data yet", "Stats appear as binding calls are made")

	v.ComponentBase = components.NewComponentBase(
		components.NewPanel().SetTitle("Stats").SetContent(v.table).SetFocused(true)).
		SetName("Stats").
		AddHint("j/k", "Move").
		AddHint("s", "Back").
		AddHint("c", "Clear").
		AddHint("q", "Back").
		SetOnStart(v.rebuild)
	v.ComponentBase.SetInputHandler(func(ev *tcell.EventKey) bool {
		return vimTableNav(v.table, ev)
	})
	return v
}

func (v *statsView) rebuild() {
	stats := v.m.topMethods(50)
	v.table.ClearRows()
	for _, s := range stats {
		avg := ""
		if s.calls > 0 && s.totalMS > 0 {
			avg = fmt.Sprintf("%.1fms", s.totalMS/float64(s.calls))
		}
		errPct := ""
		if s.calls > 0 {
			errPct = fmt.Sprintf("%.0f%%", 100*float64(s.errs)/float64(s.calls))
		}
		v.table.AddRow(
			s.method,
			fmt.Sprintf("%d", s.calls),
			fmt.Sprintf("%d", s.errs),
			avg,
			errPct,
		)
	}
}
