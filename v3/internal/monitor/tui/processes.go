package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"

	"github.com/atterpac/dado/components"
	"github.com/atterpac/dado/core"
	"github.com/atterpac/dado/theme"
	"github.com/atterpac/refresh/engine"
)

// processesView is the home screen for `wails3 dev -tui`: a table of the dev
// processes (build/app/etc.) on the left and the selected process's live log
// output on the right. State and logs come from the refresh engine SDK via the
// shared ProcStore; engine controls (reload, pause/resume) are wired here.
type processesView struct {
	*components.ComponentBase
	m *Model

	split   *components.Split
	table   *components.Table
	logPane *core.TextView

	procs []engine.ProcessInfo // last rendered snapshot (UI thread)
	sig   string               // signature of the last rendered table
	stop  chan struct{}
}

func newProcessesView(m *Model) *processesView {
	v := &processesView{m: m}

	v.table = components.NewTable()
	v.table.SetHeaders("NAME", "TYPE", "STATE", "PID", "UPTIME")
	v.table.ConfigureEmpty("…", "No processes",
		"Waiting for the dev engine to start your processes.")

	v.logPane = core.NewTextView()
	v.logPane.SetDynamicColors(true).SetScrollable(true)

	tablePanel := components.NewPanel().SetTitle("Processes").SetContent(v.table).SetFocused(true)
	logPanel := components.NewPanel().SetTitle("Logs").SetContent(v.logPane)

	v.split = components.NewSplit()
	v.split.SetDirection(components.SplitHorizontal).SetRatio(0.4).
		SetLeft(tablePanel).SetRight(logPanel)

	v.table.SetSelectionChangedFunc(func(row, col int) { v.showLogForRow(row) })

	v.ComponentBase = components.NewComponentBase(v.split).
		SetName("Processes").
		AddHint("j/k", "Move").
		AddHint("r", "Reload").
		AddHint("space", "Pause/Resume").
		AddHint("l", "Calls").
		AddHint("E", "Events").
		AddHint("b", "Bindings").
		AddHint("w", "Windows").
		AddHint("i", "Info").
		AddHint("t", "Timeline").
		AddHint("?", "Help").
		AddHint("q", "Quit").
		SetOnStart(v.start).
		SetOnStop(v.stopView)
	v.ComponentBase.SetInputHandler(v.handleKey)
	return v
}

func (v *processesView) handleKey(ev *tcell.EventKey) bool {
	if ev.Key() == tcell.KeyRune {
		switch ev.Rune() {
		case 'r':
			v.m.eng.Reload()
			v.m.toasts.Info("Reload triggered")
			return true
		case ' ':
			if v.m.eng.Paused() {
				v.m.eng.Resume()
				v.m.toasts.Info("Engine resumed")
			} else {
				v.m.eng.Pause()
				v.m.toasts.Warning("Engine paused")
			}
			v.m.updateStatus()
			return true
		}
	}
	// vim movement (j/k/g/G); arrows fall through to the table natively.
	return vimTableNav(v.table, ev)
}

// start kicks off a periodic refresh while the view is active.
func (v *processesView) start() {
	v.rebuild()
	v.stop = make(chan struct{})
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-v.stop:
				return
			case <-ticker.C:
				theme.QueueUpdateDraw(func() {
					v.rebuild()
					v.refreshLog()
				})
			}
		}
	}()
}

func (v *processesView) stopView() {
	if v.stop != nil {
		close(v.stop)
		v.stop = nil
	}
}

// rebuild repopulates the process table from the engine snapshot, but ONLY when
// the rendered content actually changed (tracked by a signature). This keeps the
// table — and the user's selection — stable while logs stream in, instead of
// clearing and rebuilding on every tick. Selection is preserved by process name.
func (v *processesView) rebuild() {
	if v.m.eng == nil {
		return
	}
	procs := v.m.eng.Processes()
	sig := procSignature(procs)
	if sig == v.sig && len(v.procs) == len(procs) {
		return // nothing visible changed; leave the cursor where the user put it
	}

	keepName := v.selectedName()
	v.procs = procs
	v.sig = sig
	v.table.ClearRows()
	keepIdx := -1
	for i, p := range procs {
		if p.Name == keepName {
			keepIdx = i
		}
		v.table.AddRowWithColor(stateColor(p.State),
			procName(p),
			string(p.Type),
			string(p.State),
			pidStr(p.PID),
			uptimeStr(p),
		)
	}

	switch {
	case keepIdx >= 0:
		v.table.SelectRow(keepIdx)
	case v.table.GetDataRowCount() > 0 && keepName == "":
		v.table.SelectRow(0)
	}
	v.showLogForRow(v.table.SelectedRow() + 1)
}

// refreshLog repaints the log pane for the current selection without touching the
// table (cheap; safe to call on every output tick).
func (v *processesView) refreshLog() {
	v.showLogForRow(v.table.SelectedRow() + 1)
}

func (v *processesView) selectedName() string {
	idx := v.table.SelectedRow()
	if idx < 0 || idx >= len(v.procs) {
		return ""
	}
	return v.procs[idx].Name
}

// procSignature captures the table-visible state (name/type/state/pid + uptime to
// the second) so rebuild can skip no-op redraws.
func procSignature(procs []engine.ProcessInfo) string {
	var b strings.Builder
	for _, p := range procs {
		fmt.Fprintf(&b, "%s|%s|%s|%d|%s\n", p.Name, p.Type, p.State, p.PID, uptimeStr(p))
	}
	return b.String()
}

func (v *processesView) showLogForRow(tableRow int) {
	idx := tableRow - 1
	if idx < 0 || idx >= len(v.procs) {
		v.logPane.SetText("")
		return
	}
	lines := v.m.store.Lines(v.procs[idx].Name)
	v.logPane.SetText(strings.Join(lines, "\n"))
	// Follow the tail; ScrollTo clamps past-end offsets to the bottom.
	v.logPane.ScrollTo(len(lines), 0)
}

// ---- formatting helpers ----

func stateColor(s engine.ProcessState) tcell.Color {
	switch s {
	case engine.StateRunning:
		return theme.Success()
	case engine.StateFailed:
		return theme.Error()
	case engine.StatePending:
		return theme.Warning()
	default: // exited / killed
		return theme.FgDim()
	}
}

func procName(p engine.ProcessInfo) string {
	if p.Name != "" {
		return p.Name
	}
	return p.Exec
}

func pidStr(pid int) string {
	if pid == 0 {
		return "—"
	}
	return fmt.Sprintf("%d", pid)
}

func uptimeStr(p engine.ProcessInfo) string {
	if p.State != engine.StateRunning || p.StartedAt.IsZero() {
		return "—"
	}
	d := time.Since(p.StartedAt).Truncate(time.Second)
	return d.String()
}
