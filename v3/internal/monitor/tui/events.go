package tui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"

	"github.com/atterpac/dado/components"
	"github.com/atterpac/dado/core"
	"github.com/atterpac/dado/theme"
)

// eventsView is the event firehose: every custom event crossing the IPC, with
// consecutive repeats collapsed to a ×N count. High-frequency window events are
// muted by default; toggle a name's mute with `m`.
type eventsView struct {
	*components.ComponentBase
	m *Model

	split  *components.Split
	table  *components.Table
	detail *components.Label
	body   *components.Label
}

func newEventsView(m *Model) *eventsView {
	v := &eventsView{m: m}

	v.table = components.NewTable()
	v.table.SetHeaders("TIME", "", "EVENT", "WINDOW", "×N")
	v.table.ConfigureEmpty("…", "No events yet",
		"Events appear here. Noisy window events are muted by default.")

	v.detail = components.NewLabel("").SetDynamicColors(true)
	v.body = components.NewLabel("").SetDynamicColors(true).SetScrollable(true)

	detailPane := core.NewFlex().SetDirection(core.Column).
		AddItem(v.detail, 5, 0, false).
		AddItem(v.body, 0, 1, false)

	v.split = components.NewSplit()
	v.split.SetDirection(components.SplitHorizontal).SetRatio(0.55).
		SetLeft(v.table).SetRight(detailPane)

	v.table.SetSelectionChangedFunc(func(row, col int) { v.showDetailForRow(row) })

	v.ComponentBase = components.NewComponentBase(v.split).
		SetName("Events").
		AddHint("j/k", "Move").
		AddHint("Enter", "Detail").
		AddHint("y", "Copy").
		AddHint("m", "Mute name").
		AddHint("/", "Filter").
		AddHint("c", "Clear").
		AddHint("q", "Back")
	v.ComponentBase.SetInputHandler(v.handleKey)
	return v
}

func (v *eventsView) handleKey(ev *tcell.EventKey) bool {
	if ev.Key() == tcell.KeyEnter {
		v.m.showDetailModal(v.selectedRecord())
		return true
	}
	if ev.Key() == tcell.KeyRune {
		switch ev.Rune() {
		case 'y':
			v.m.copyRecord(v.selectedRecord())
			return true
		case 'm':
			if r := v.selectedRecord(); r != nil {
				v.m.toggleMute(r.method)
			}
			return true
		case 'G':
			v.m.follow = true
			v.selectLast()
			v.m.updateStatus()
			return true
		}
	}
	if isScrollUpKey(ev) && v.m.follow {
		v.m.follow = false
		v.m.updateStatus()
	}
	return false
}

func (v *eventsView) selectedRecord() *record {
	idx := v.table.SelectedRow()
	v.m.mu.Lock()
	defer v.m.mu.Unlock()
	if idx < 0 || idx >= len(v.m.visibleEvents) {
		return nil
	}
	return v.m.visibleEvents[idx]
}

func (v *eventsView) rebuild() {
	var keepSeq uint64
	if !v.m.follow {
		if r := v.selectedRecord(); r != nil {
			keepSeq = r.seq
		}
	}

	recs := v.m.snapshotEvents()
	v.table.ClearRows()
	keepIdx := -1
	for i, r := range recs {
		if r.seq == keepSeq {
			keepIdx = i
		}
		count := ""
		if r.count > 1 {
			count = fmt.Sprintf("×%d", r.count)
		}
		v.table.AddRowWithColor(theme.Accent(),
			r.t.Format("15:04:05.000"),
			dirArrow(r.dir),
			r.method,
			r.window,
			count,
		)
	}

	if v.m.follow {
		v.selectLast()
	} else if keepIdx >= 0 {
		v.table.SelectRow(keepIdx)
	}
}

func (v *eventsView) selectLast() {
	n := v.table.GetDataRowCount()
	if n > 0 {
		v.table.SelectRow(n - 1)
		v.showDetailForRow(n)
	}
}

func (v *eventsView) showDetailForRow(tableRow int) {
	idx := tableRow - 1
	v.m.mu.Lock()
	recs := v.m.visibleEvents
	muted := false
	if idx >= 0 && idx < len(recs) {
		muted = v.m.muted[recs[idx].method]
	}
	v.m.mu.Unlock()
	if idx < 0 || idx >= len(recs) {
		v.detail.SetText("")
		v.body.SetText("")
		return
	}
	r := recs[idx]

	muteState := "[gray]m mute[-]"
	if muted {
		muteState = "[yellow]m unmute[-]"
	}
	dirLabel := "Go → JS"
	if r.dir == "in" {
		dirLabel = "JS → Go"
	}
	var h strings.Builder
	fmt.Fprintf(&h, "[::b]%s[-:-:-]", escapeTags(r.method))
	if r.count > 1 {
		fmt.Fprintf(&h, "  [gray]×%d[-]", r.count)
	}
	fmt.Fprintf(&h, "\n[gray]direction[-] %s   [gray]window[-] %s\n", dirLabel, orDash(r.window))
	fmt.Fprintf(&h, "[gray]enter expand · y copy · %s[-]", muteState)
	v.detail.SetText(h.String())

	if len(r.args) > 0 {
		v.body.SetText(section("Payload") + escapeTags(prettyJSON(r.args)))
	} else {
		v.body.SetText(section("Payload") + "[gray](no data)[-]")
	}
}
