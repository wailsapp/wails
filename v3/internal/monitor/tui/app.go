// Package tui implements the dado-based terminal UI for `wails3 monitor`.
//
// Layout/navigation follow the conventions used by the qry and gxt dado apps:
//   - a StatusBar top bar with colored sections,
//   - a Menu bottom bar whose hints are driven by the active nav.Component,
//   - an ActionRegistry as the single source of truth for global keys + help,
//   - a global SetInputCapture for quit / escape / modal dismissal,
//   - a ToastManager for transient feedback.
//
// DevX choices:
//   - binding calls and events are SEPARATE streams. The default Live view shows
//     calls only (the firehose of window/mouse events lives in its own view).
//   - consecutive identical events collapse into one row with a ×N count.
//   - a mute set (seeded with the noisy default window events) hides chosen
//     event names from the events view.
package tui

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"

	"github.com/atterpac/dado/clipboard"
	"github.com/atterpac/dado/components"
	"github.com/atterpac/dado/core"
	"github.com/atterpac/dado/help"
	"github.com/atterpac/dado/input"
	"github.com/atterpac/dado/layout"
	"github.com/atterpac/dado/nav"
	"github.com/atterpac/dado/theme"
	"github.com/atterpac/refresh/engine"

	monitor "github.com/wailsapp/wails/v3/internal/monitor"
)

// defaultMutedEvents are high-frequency Wails window events hidden by default in
// the events view (still counted). Toggle per-name with `m`.
var defaultMutedEvents = []string{
	"common:WindowDidResize",
	"common:WindowDidMove",
	"common:WindowDPIChanged",
	"common:WindowFocus",
	"common:WindowLostFocus",
	"common:WindowZoom",
	"linux:WindowDidResize",
	"linux:WindowDidMove",
	"linux:WindowFocusIn",
	"linux:WindowFocusOut",
}

// record is one displayed line. Calls carry a result/error merged in later via
// callID; events may carry a collapse count.
type record struct {
	seq     uint64
	t       time.Time
	kind    string // call | event | cancel
	dir     string
	callID  string
	method  string
	window  string
	args    json.RawMessage
	result  json.RawMessage
	errMsg  string
	errKind string
	durMS   float64
	status  string // pending | ok | error | cancelled | ""
	count   int    // collapse count for repeated events (1 = single)
}

// methodStat aggregates per-method counters for the stats view.
type methodStat struct {
	method  string
	calls   int
	errs    int
	totalMS float64
}

// Model holds all UI state. All mutation happens on the UI thread.
type Model struct {
	app       *layout.App
	statusBar *layout.StatusBar
	menu      *layout.Menu
	toasts    *components.ToastManager
	actions   *input.ActionRegistry

	processes *processesView
	live      *liveView
	events    *eventsView
	stats     *statsView
	windows   *windowsView
	bindings  *bindingsView
	info      *infoView
	timeline  *timelineView

	// dev-mode engine integration (set by RunDev).
	eng   *engine.Engine
	store *ProcStore

	client   *monitor.Client
	target   monitor.DiscoveryEntry
	snapshot *monitor.Snapshot

	// samples is the rolling resource-usage history (oldest first), capped at
	// maxSamples. Shares a time axis with callRecs/eventRecs for correlation.
	samples []monitor.Sample

	mu        sync.Mutex
	callRecs  []*record          // binding calls (Live view)
	eventRecs []*record          // events (Events view), collapsed
	byCall    map[string]*record // callID -> call record
	methodAgg map[string]*methodStat
	muted     map[string]bool

	filter     string
	errorsOnly bool
	paused     bool
	follow     bool
	connected  bool

	// Time-range filter set by the timeline scrubber; when rangeActive the
	// stream views show only records within [rangeFrom,rangeTo].
	rangeActive bool
	rangeFrom   time.Time
	rangeTo     time.Time

	visibleCalls  []*record // filtered calls shown in the live table
	visibleEvents []*record // filtered events shown in the events table

	// counters
	calls, errs, events_, pending, mutedHits int
}

// RunDev runs the TUI for `wails3 dev -tui`. The processes view is the home
// screen, fed by the refresh engine SDK. As soon as the running app exposes its
// IPC monitor socket, the TUI attaches and the bindings / events / calls views
// come alive; it reattaches across reloads (each restart is a new PID/socket).
func RunDev(ctx context.Context, eng *engine.Engine, store *ProcStore) error {
	m := &Model{
		eng:       eng,
		store:     store,
		byCall:    map[string]*record{},
		methodAgg: map[string]*methodStat{},
		muted:     map[string]bool{},
		follow:    true,
	}
	for _, e := range defaultMutedEvents {
		m.muted[e] = true
	}
	m.buildApp()
	m.setup()

	// Engine output / lifecycle changes refresh the processes view + status.
	store.OnLog = func() {
		theme.QueueUpdateDraw(func() {
			if m.processes != nil {
				m.processes.refreshLog()
			}
		})
	}
	store.OnEvent = func(ev engine.ProcessEvent) {
		theme.QueueUpdateDraw(func() {
			if m.processes != nil {
				m.processes.rebuild()
			}
			if ev.Info.State == engine.StateFailed && ev.Err != nil {
				m.toasts.Error(ev.Info.Name + ": " + ev.Err.Error())
			}
		})
	}

	streamCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go m.discoverLoop(streamCtx)

	return m.app.Run()
}

// discoverLoop watches for the dev app's IPC monitor socket and attaches when it
// appears, reattaching after each reload.
//
// It cannot match on the engine's Primary PID: refresh launches the primary via
// `/bin/sh -c <cmd>`, so the engine records the SHELL's PID while the app writes
// its discovery file with its own (child) PID. Instead it ignores any discovery
// entries that already existed when the TUI started (unrelated apps) and attaches
// to the newest entry that appears afterward — in dev that is the app we just
// launched, and on reload the restarted app's socket is new again (the old one's
// process dies and is pruned by monitor.List).
func (m *Model) discoverLoop(ctx context.Context) {
	// Baseline: sockets present before we launched anything.
	baseline := map[string]bool{}
	if entries, err := monitor.List(); err == nil {
		for _, e := range entries {
			baseline[e.Sock] = true
		}
	}

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.mu.Lock()
			connected := m.connected
			m.mu.Unlock()
			if connected {
				continue
			}
			entries, err := monitor.List()
			if err != nil {
				continue
			}
			var best *monitor.DiscoveryEntry
			for i := range entries {
				if baseline[entries[i].Sock] {
					continue // pre-existing / unrelated app
				}
				if best == nil || entries[i].StartedAt.After(best.StartedAt) {
					best = &entries[i]
				}
			}
			if best != nil {
				m.attach(ctx, *best)
			}
		}
	}
}

// attach connects to a discovered app's monitor socket and starts the trace /
// snapshot / sample stream goroutines. On disconnect it clears the connected
// flag so discoverLoop reattaches after the next reload.
func (m *Model) attach(ctx context.Context, target monitor.DiscoveryEntry) {
	client, err := monitor.Connect(ctx, target.Sock)
	if err != nil {
		return
	}
	m.mu.Lock()
	m.client = client
	m.target = target
	m.connected = true
	m.mu.Unlock()
	theme.QueueUpdateDraw(func() {
		m.updateStatus()
		m.toasts.Success("Attached to " + target.Name)
	})

	go func() {
		for {
			select {
			case t, ok := <-client.Traces():
				if !ok {
					theme.QueueUpdateDraw(func() {
						m.mu.Lock()
						m.connected = false
						m.mu.Unlock()
						m.updateStatus()
						m.toasts.Warning("Disconnected from app")
					})
					return
				}
				theme.QueueUpdateDraw(func() { m.apply(t) })
			case <-ctx.Done():
				return
			}
		}
	}()
	go func() {
		for {
			select {
			case snap, ok := <-client.Snapshots():
				if !ok {
					return
				}
				theme.QueueUpdateDraw(func() { m.applySnapshot(snap) })
			case <-ctx.Done():
				return
			}
		}
	}()
	go func() {
		for {
			select {
			case s, ok := <-client.Samples():
				if !ok {
					return
				}
				theme.QueueUpdateDraw(func() { m.applySample(s) })
			case <-ctx.Done():
				return
			}
		}
	}()
	go func() {
		if e := <-client.Errors(); e != nil {
			theme.QueueUpdateDraw(func() {
				m.mu.Lock()
				m.connected = false
				m.mu.Unlock()
				m.updateStatus()
				m.toasts.Error("stream error: " + e.Error())
			})
		}
	}()

	// Pull an initial snapshot for the dashboard views.
	_ = client.Describe()
}

// buildApp constructs the shell (status bar, menu, toasts) — qry style.
func (m *Model) buildApp() {
	m.statusBar = layout.NewStatusBar()
	m.statusBar.SetTitle("wails monitor")
	m.statusBar.SetTitleAlign(components.AlignLeft)
	m.menu = layout.NewMenu()

	m.app = layout.NewApp(layout.AppConfig{
		TopBar:          m.statusBar,
		BottomBar:       m.menu,
		TopBarHeight:    3,
		BottomBarHeight: 1,
		OnComponentChange: func(c nav.Component) {
			if c != nil {
				m.menu.SetHints(c.Hints())
			}
		},
	})
	m.app.EnableThemes(layout.ThemeOptions{Default: "tokyonight-night"})

	m.toasts = components.NewToastManager()
	m.toasts.SetPosition(components.ToastBottomRight)
	m.toasts.SetMaxVisible(3)
	m.toasts.SetDefaultDuration(3 * time.Second)
	m.app.GetApp().SetAfterDrawFunc(func(screen tcell.Screen) {
		w, h := screen.Size()
		m.toasts.Draw(screen, w, h)
	})

	m.updateStatus()
}

// setup wires global keys and pushes the initial view.
func (m *Model) setup() {
	m.actions = input.NewActionRegistry().
		AddSimple("help", '?', "Help", m.showHelp).
		AddSimple("calls", 'l', "Live calls view", m.showLive).
		AddSimple("events", 'E', "Events view", m.showEvents).
		AddSimple("stats", 's', "Stats view", m.toggleStats).
		AddSimple("windows", 'w', "Windows", m.showWindows).
		AddSimple("bindings", 'b', "Bindings", m.showBindings).
		AddSimple("info", 'i', "App info", m.showInfo).
		AddSimple("timeline", 't', "Timeline", m.showTimeline).
		AddSimple("follow", 'f', "Toggle follow", m.toggleFollow).
		AddSimple("pause", 'p', "Toggle pause", m.togglePause).
		AddSimple("errors", 'e', "Errors only", m.toggleErrors).
		AddSimple("clear", 'c', "Clear", m.clear).
		AddSimple("filter", '/', "Filter", m.startFilter)

	m.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if m.statusBar.IsCommandMode() {
			return event
		}
		// Open theme selector on T (Shift+T). Ctrl+T also works via EnableThemes.
		if event.Key() == tcell.KeyRune && event.Rune() == 'T' {
			m.app.QueueUpdateDraw(m.app.OpenThemeSelector)
			return nil
		}
		isModal := m.app.Pages().CurrentIsModal()

		switch {
		case event.Rune() == 'q' && !isModal:
			if m.app.Pages().CanPop() {
				m.app.Pages().Pop()
			} else {
				m.app.Stop()
			}
			return nil
		case event.Key() == tcell.KeyEscape:
			if isModal {
				m.app.Pages().DismissModal()
				return nil
			}
			if m.app.Pages().CanPop() {
				m.app.Pages().Pop()
				return nil
			}
		}

		// The timeline's panes are scrollable Labels that swallow nav keys
		// before the component handler runs, so route its scrub keys here (the
		// app capture fires before any focused primitive).
		if !isModal && m.timeline != nil && m.app.Pages().Current() == m.timeline {
			if m.timeline.handleKey(event) {
				return nil
			}
		}

		if !isModal && m.actions.Handle(event) {
			return nil
		}
		return event
	})

	// Live is created eagerly (the trace-ingest path mutates it unconditionally)
	// but the processes view is the home screen for dev mode.
	m.live = newLiveView(m)
	m.processes = newProcessesView(m)
	m.app.Pages().Push(m.processes)
}

// showLive opens (or toggles) the binding-call Live view. It is no longer the
// home screen, so it is reachable via the `l` action.
func (m *Model) showLive() {
	if m.app.Pages().Current() == m.live && m.app.Pages().CanPop() {
		m.app.Pages().Pop()
		return
	}
	m.live.rebuild()
	m.app.Pages().Push(m.live)
}

// ---- trace ingestion (UI thread) ----

func (m *Model) apply(t monitor.Trace) {
	m.mu.Lock()
	switch t.Kind {
	case "call":
		r := &record{
			seq: t.Seq, t: t.Time, kind: "call", dir: t.Dir, callID: t.CallID,
			method: t.Method, window: t.Window, args: t.Args, status: "pending", count: 1,
		}
		m.callRecs = append(m.callRecs, r)
		m.byCall[t.CallID] = r
		m.calls++
		m.pending++
		m.aggFor(r.method).calls++
	case "result":
		if r := m.byCall[t.CallID]; r != nil {
			r.status = "ok"
			r.result = t.Result
			r.durMS = t.DurationMS
			m.pending--
			m.aggFor(r.method).totalMS += t.DurationMS
		}
	case "error":
		if r := m.byCall[t.CallID]; r != nil {
			r.status = "error"
			if t.Error != nil {
				r.errMsg = t.Error.Message
				r.errKind = t.Error.Kind
			}
			r.durMS = t.DurationMS
			m.pending--
			m.errs++
			a := m.aggFor(r.method)
			a.errs++
			a.totalMS += t.DurationMS
		}
	case "cancel":
		if r := m.byCall[t.CallID]; r != nil {
			if r.status == "pending" {
				m.pending--
			}
			r.status = "cancelled"
		}
	case "event":
		m.events_++
		if m.muted[t.Method] {
			m.mutedHits++
		}
		// Collapse consecutive identical events (same name + direction).
		if n := len(m.eventRecs); n > 0 {
			last := m.eventRecs[n-1]
			if last.method == t.Method && last.dir == t.Dir {
				last.count++
				last.t = t.Time
				last.args = t.Args
				m.mu.Unlock()
				if !m.paused && m.events != nil {
					m.events.rebuild()
				}
				m.updateStatus()
				return
			}
		}
		m.eventRecs = append(m.eventRecs, &record{
			seq: t.Seq, t: t.Time, kind: "event", dir: t.Dir, callID: t.CallID,
			method: t.Method, window: t.Window, args: t.Args, count: 1,
		})
	}
	m.mu.Unlock()

	if !m.paused {
		m.live.rebuild()
		if m.events != nil {
			m.events.rebuild()
		}
		if m.stats != nil && m.app.Pages().Current() == m.stats {
			m.stats.rebuild()
		}
	}
	m.updateStatus()
}

func (m *Model) aggFor(method string) *methodStat {
	if method == "" {
		method = "(event)"
	}
	a := m.methodAgg[method]
	if a == nil {
		a = &methodStat{method: method}
		m.methodAgg[method] = a
	}
	return a
}

// snapshotCalls returns the filtered call list (UI thread).
func (m *Model) snapshotCalls() []*record {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]*record, 0, len(m.callRecs))
	for _, r := range m.callRecs {
		if m.errorsOnly && r.status != "error" {
			continue
		}
		if m.rangeActive && (r.t.Before(m.rangeFrom) || r.t.After(m.rangeTo)) {
			continue
		}
		if m.filter != "" && !recordMatches(r, m.filter) {
			continue
		}
		out = append(out, r)
	}
	m.visibleCalls = out
	return out
}

// snapshotEvents returns the filtered, non-muted event list (UI thread).
func (m *Model) snapshotEvents() []*record {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]*record, 0, len(m.eventRecs))
	for _, r := range m.eventRecs {
		if m.muted[r.method] {
			continue
		}
		if m.rangeActive && (r.t.Before(m.rangeFrom) || r.t.After(m.rangeTo)) {
			continue
		}
		if m.filter != "" && !recordMatches(r, m.filter) {
			continue
		}
		out = append(out, r)
	}
	m.visibleEvents = out
	return out
}

func (m *Model) topMethods(limit int) []methodStat {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]methodStat, 0, len(m.methodAgg))
	for _, a := range m.methodAgg {
		out = append(out, *a)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].calls != out[j].calls {
			return out[i].calls > out[j].calls
		}
		return out[i].method < out[j].method
	})
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out
}

// windowActivity is per-window IPC activity derived from the trace stream.
type windowActivity struct {
	calls, errs, events int
	topMethods          []methodCount // most-called methods by this window
	recent              []*record     // most recent calls (newest last)
}

type methodCount struct {
	method string
	count  int
}

// activityForWindow aggregates the calls/events a given window has generated.
// Matching is by the trace's window field; bindings themselves are global.
func (m *Model) activityForWindow(name string) windowActivity {
	m.mu.Lock()
	defer m.mu.Unlock()

	var act windowActivity
	counts := map[string]int{}
	for _, r := range m.callRecs {
		if r.window != name {
			continue
		}
		act.calls++
		if r.status == "error" {
			act.errs++
		}
		counts[r.method]++
		act.recent = append(act.recent, r)
	}
	for _, r := range m.eventRecs {
		if r.window == name {
			act.events += maxInt(r.count, 1)
		}
	}

	for method, c := range counts {
		act.topMethods = append(act.topMethods, methodCount{method, c})
	}
	sort.Slice(act.topMethods, func(i, j int) bool {
		if act.topMethods[i].count != act.topMethods[j].count {
			return act.topMethods[i].count > act.topMethods[j].count
		}
		return act.topMethods[i].method < act.topMethods[j].method
	})
	if len(act.topMethods) > 5 {
		act.topMethods = act.topMethods[:5]
	}
	if len(act.recent) > 6 {
		act.recent = act.recent[len(act.recent)-6:]
	}
	return act
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func recordMatches(r *record, q string) bool {
	q = strings.ToLower(q)
	return strings.Contains(strings.ToLower(r.method), q) ||
		strings.Contains(strings.ToLower(r.window), q) ||
		strings.Contains(strings.ToLower(r.kind), q)
}

// updateStatus rebuilds the status bar sections (no QueueUpdateDraw inside).
func (m *Model) updateStatus() {
	m.mu.Lock()
	calls, errs, events, pending, muted := m.calls, m.errs, m.events_, m.pending, m.mutedHits
	paused, errorsOnly, filter, connected := m.paused, m.errorsOnly, m.filter, m.connected
	m.mu.Unlock()

	m.statusBar.ClearSections()
	if connected {
		m.statusBar.SetConnectionStatus(true, m.target.Name)
		m.statusBar.AddSection(layout.StatusSection{Text: fmt.Sprintf("pid %d", m.target.PID), Color: theme.FgDim()})
	} else {
		m.statusBar.SetConnectionStatus(false, "")
		m.statusBar.AddSection(layout.StatusSection{Text: "monitor: waiting…", Color: theme.FgDim()})
	}

	// Engine (refresh) pause state, when running under `dev -tui`.
	if m.eng != nil {
		if m.eng.Paused() {
			m.statusBar.AddSection(layout.StatusSection{Text: "⏸ engine paused", Color: theme.Warning()})
		} else {
			m.statusBar.AddSection(layout.StatusSection{Text: "▶ engine", Color: theme.Success()})
		}
	}
	if paused {
		m.statusBar.AddSection(layout.StatusSection{Text: "❚❚ paused", Color: theme.Warning()})
	} else {
		m.statusBar.AddSection(layout.StatusSection{Text: "● live", Color: theme.Success()})
	}

	m.statusBar.SetRightSections(nil)
	m.statusBar.AddRightSection(layout.StatusSection{Text: fmt.Sprintf("%d calls", calls), Color: theme.Fg()})
	m.statusBar.AddRightSection(layout.StatusSection{Text: fmt.Sprintf("%d err", errs), Color: theme.Error()})
	m.statusBar.AddRightSection(layout.StatusSection{Text: fmt.Sprintf("%d evt", events), Color: theme.Accent()})
	m.statusBar.AddRightSection(layout.StatusSection{Text: fmt.Sprintf("%d pending", pending), Color: theme.FgDim()})
	if muted > 0 {
		m.statusBar.AddRightSection(layout.StatusSection{Text: fmt.Sprintf("%d muted", muted), Color: theme.FgDim()})
	}
	if errorsOnly {
		m.statusBar.AddRightSection(layout.StatusSection{Text: "errors-only", Color: theme.Warning()})
	}
	if filter != "" {
		m.statusBar.AddRightSection(layout.StatusSection{Text: "/" + filter, Color: theme.Accent()})
	}
}

// ---- global actions ----

func (m *Model) toggleFollow() {
	m.follow = !m.follow
	if m.follow {
		m.live.selectLast()
	}
	m.toasts.Info(boolMsg("Follow", m.follow))
}

func (m *Model) togglePause() {
	m.paused = !m.paused
	if !m.paused {
		m.live.rebuild()
		if m.events != nil {
			m.events.rebuild()
		}
	}
	m.updateStatus()
	m.toasts.Info(boolMsg("Paused", m.paused))
}

func (m *Model) toggleErrors() {
	m.errorsOnly = !m.errorsOnly
	m.live.rebuild()
	m.updateStatus()
	m.toasts.Info(boolMsg("Errors-only", m.errorsOnly))
}

func (m *Model) clear() {
	m.mu.Lock()
	m.callRecs = nil
	m.eventRecs = nil
	m.byCall = map[string]*record{}
	m.methodAgg = map[string]*methodStat{}
	m.calls, m.errs, m.events_, m.pending, m.mutedHits = 0, 0, 0, 0, 0
	m.mu.Unlock()
	m.live.rebuild()
	if m.events != nil {
		m.events.rebuild()
	}
	if m.stats != nil {
		m.stats.rebuild()
	}
	m.updateStatus()
	m.toasts.Info("Cleared")
}

func (m *Model) startFilter() {
	m.statusBar.SetCommandPrompt("/ ")
	m.statusBar.SetCommandPlaceholder("filter by method / window / kind…")
	m.statusBar.EnterCommandMode()
	m.statusBar.SetCommandText(m.filter)
	m.app.SetFocus(m.statusBar)

	apply := func() {
		m.live.rebuild()
		if m.events != nil {
			m.events.rebuild()
		}
		m.updateStatus()
	}
	m.statusBar.SetOnCommandChange(func(text string) {
		m.mu.Lock()
		m.filter = text
		m.mu.Unlock()
		apply()
	})
	m.statusBar.SetOnCommandSubmit(func(string) {
		m.statusBar.SetOnCommandChange(nil)
		m.statusBar.ExitCommandMode()
		m.refocusCurrent()
	})
	m.statusBar.SetOnCommandCancel(func() {
		m.statusBar.SetOnCommandChange(nil)
		m.mu.Lock()
		m.filter = ""
		m.mu.Unlock()
		m.statusBar.ExitCommandMode()
		apply()
		m.refocusCurrent()
	})
}

func (m *Model) showEvents() {
	if m.app.Pages().Current() == m.events && m.app.Pages().CanPop() {
		m.app.Pages().Pop()
		return
	}
	if m.events == nil {
		m.events = newEventsView(m)
	}
	m.events.rebuild()
	m.app.Pages().Push(m.events)
}

// applySnapshot stores the latest snapshot and refreshes any open dashboard
// views (UI thread).
func (m *Model) applySnapshot(snap *monitor.Snapshot) {
	m.snapshot = snap
	if m.windows != nil {
		m.windows.rebuild()
	}
	if m.bindings != nil {
		m.bindings.rebuild()
	}
	if m.info != nil {
		m.info.rebuild()
	}
}

// requestSnapshot asks the app for a fresh snapshot (non-blocking).
func (m *Model) requestSnapshot() {
	if m.client != nil {
		_ = m.client.Describe()
	}
}

func (m *Model) showWindows() {
	if m.app.Pages().Current() == m.windows && m.app.Pages().CanPop() {
		m.app.Pages().Pop()
		return
	}
	if m.windows == nil {
		m.windows = newWindowsView(m)
	}
	m.requestSnapshot() // refresh on open
	m.windows.rebuild()
	m.app.Pages().Push(m.windows)
}

func (m *Model) showBindings() {
	if m.app.Pages().Current() == m.bindings && m.app.Pages().CanPop() {
		m.app.Pages().Pop()
		return
	}
	if m.bindings == nil {
		m.bindings = newBindingsView(m)
	}
	m.requestSnapshot()
	m.bindings.rebuild()
	m.app.Pages().Push(m.bindings)
}

func (m *Model) showInfo() {
	if m.app.Pages().Current() == m.info && m.app.Pages().CanPop() {
		m.app.Pages().Pop()
		return
	}
	if m.info == nil {
		m.info = newInfoView(m)
	}
	m.requestSnapshot()
	m.info.rebuild()
	m.app.Pages().Push(m.info)
}

func (m *Model) toggleStats() {
	if m.app.Pages().Current() == m.stats && m.app.Pages().CanPop() {
		m.app.Pages().Pop()
		return
	}
	if m.stats == nil {
		m.stats = newStatsView(m)
	}
	m.stats.rebuild()
	m.app.Pages().Push(m.stats)
}

// toggleMute flips the mute state of an event name (called from events view).
func (m *Model) toggleMute(name string) {
	if name == "" {
		return
	}
	m.mu.Lock()
	m.muted[name] = !m.muted[name]
	muted := m.muted[name]
	m.mu.Unlock()
	if m.events != nil {
		m.events.rebuild()
	}
	m.updateStatus()
	if muted {
		m.toasts.Info("Muted " + name)
	} else {
		m.toasts.Info("Unmuted " + name)
	}
}

// jumpToWindowStream pops to the Live view and filters the trace stream to the
// given window name, so the user can see exactly what that window is doing.
// jumpToRange pins the stream views to a time window and switches to live.
func (m *Model) jumpToRange(from, to time.Time) {
	m.mu.Lock()
	m.rangeActive = true
	m.rangeFrom = from
	m.rangeTo = to
	m.follow = false
	m.mu.Unlock()
	for m.app.Pages().Current() != m.live && m.app.Pages().CanPop() {
		m.app.Pages().Pop()
	}
	m.live.rebuild()
	m.updateStatus()
	m.toasts.Info("Filtered to selected time range")
}

// clearRange removes any active time-range filter.
func (m *Model) clearRange() {
	m.mu.Lock()
	m.rangeActive = false
	m.mu.Unlock()
	m.live.rebuild()
	m.events.rebuild()
	m.updateStatus()
}

func (m *Model) jumpToWindowStream(name string) {
	m.mu.Lock()
	m.filter = name
	m.errorsOnly = false
	m.mu.Unlock()
	// Return to the Live view (pop any pushed dashboard views).
	for m.app.Pages().Current() != m.live && m.app.Pages().CanPop() {
		m.app.Pages().Pop()
	}
	m.live.rebuild()
	m.updateStatus()
	m.toasts.Info("Filtered to window: " + name)
}

func (m *Model) refocusCurrent() {
	m.app.SetFocus(m.app.Pages())
}

// copyRecord copies a record's args+result JSON to the clipboard.
func (m *Model) copyRecord(r *record) {
	if r == nil {
		return
	}
	var b strings.Builder
	if len(r.args) > 0 {
		b.WriteString(string(r.args))
	}
	if len(r.result) > 0 {
		if b.Len() > 0 {
			b.WriteString("\n")
		}
		b.WriteString(string(r.result))
	}
	if b.Len() == 0 {
		m.toasts.Warning("Nothing to copy")
		return
	}
	if err := clipboard.Copy(b.String()); err != nil {
		m.toasts.Error("Copy failed: " + err.Error())
		return
	}
	m.toasts.Success("Copied to clipboard")
}

// showDetailModal opens a scrollable fullscreen detail of a record.
func (m *Model) showDetailModal(r *record) {
	if r == nil {
		return
	}
	tv := core.NewTextView()
	tv.SetDynamicColors(true).SetScrollable(true)
	tv.SetText(detailText(r))

	modal := components.NewModal(components.ModalConfig{
		Title:    methodOr(r),
		Width:    80,
		Height:   28,
		Backdrop: true,
	}).SetContent(tv).
		SetHints([]components.KeyHint{
			{Key: "j/k", Description: "Scroll"},
			{Key: "Esc", Description: "Close"},
		})
	m.app.Pages().Push(modal)
}

func (m *Model) showHelp() {
	var b strings.Builder
	b.WriteString("[::b]Wails IPC Monitor[::-]\n")
	for _, section := range m.helpModel().GetSections() {
		b.WriteString("\n[::b]" + section.Name + "[::-]\n")
		for _, act := range section.Actions {
			fmt.Fprintf(&b, "  %-8s %s\n", act.Key, act.Description)
		}
	}
	tv := core.NewTextView()
	tv.SetDynamicColors(true).SetScrollable(true)
	tv.SetText(strings.TrimRight(b.String(), "\n"))

	modal := components.NewModal(components.ModalConfig{
		Title: "Help", Width: 56, Height: 22,
	}).SetContent(tv)
	m.app.Pages().Push(modal)
}

func (m *Model) helpModel() *help.Help {
	return help.New().
		SetAppName("wails monitor").
		AddSection("Navigation", []help.ActionInfo{
			{Key: "j/k", Description: "Move selection (scrolling up pauses follow)"},
			{Key: "G", Description: "Jump to tail + resume follow"},
			{Key: "Enter", Description: "Expand detail"},
			{Key: "y", Description: "Copy args/result"},
			{Key: "q", Description: "Quit / pop view"},
			{Key: "Esc", Description: "Pop view / dismiss"},
		}).
		AddSection("Events view", []help.ActionInfo{
			{Key: "m", Description: "Mute/unmute event name"},
		}).
		AddRegistry("Global Hotkeys", m.actions)
}

// ---- shared helpers ----

// tview color tag helpers that delegate to dado's theme so colors follow the
// active theme (tokyonight-night or any user override) instead of hardcoded names.
func tagDim() string    { return "[" + theme.TagFgDim() + "]" }
func tagErr() string    { return "[" + theme.TagError() + "]" }
func tagWarn() string   { return "[" + theme.TagWarning() + "]" }
func tagOK() string     { return "[" + theme.TagSuccess() + "]" }
func tagAccent() string { return "[" + theme.TagAccent() + "]" }
func tagReset() string  { return "[-]" }

// escapeTags neutralizes dado color-tag markup ("[...]") in user-supplied text
// so it is rendered literally by core.TextView / components.Label. A "[" is
// rewritten to "[ " so it can never open a tag (the dado parser only treats a
// well-formed "[tag]" as markup).
func escapeTags(s string) string {
	if !strings.Contains(s, "[") {
		return s
	}
	return strings.ReplaceAll(s, "[", "[​")
}

func boolMsg(label string, on bool) string {
	if on {
		return label + ": on"
	}
	return label + ": off"
}

func prettyJSON(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var buf bytes.Buffer
	if err := json.Indent(&buf, raw, "", "  "); err != nil {
		return string(raw)
	}
	return buf.String()
}

// detailText builds the full dynamic-color detail body for a record.
func detailText(r *record) string {
	st := r.status
	if st == "" {
		st = "—"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "[::b]%s[-:-:-]\n", escapeTags(methodOr(r)))
	fmt.Fprintf(&b, "[gray]callId[-] %s   [gray]window[-] %s\n", orDash(r.callID), orDash(r.window))
	fmt.Fprintf(&b, "[gray]dir[-] %s   [gray]dur[-] %s   [gray]status[-] %s", r.dir, orDash(durStr(r)), st)
	if r.count > 1 {
		fmt.Fprintf(&b, "   [gray]count[-] ×%d", r.count)
	}
	if r.errMsg != "" {
		fmt.Fprintf(&b, "\n[red]error[-] %s", escapeTags(r.errMsg))
		if r.errKind != "" {
			fmt.Fprintf(&b, " [gray](%s)[-]", r.errKind)
		}
	}
	if len(r.args) > 0 {
		b.WriteString("\n\n[gray]// args[-]\n" + escapeTags(prettyJSON(r.args)))
	}
	if len(r.result) > 0 {
		b.WriteString("\n\n[gray]// result[-]\n" + escapeTags(prettyJSON(r.result)))
	}
	return b.String()
}

func statusGlyph(r *record) string {
	switch r.status {
	case "ok":
		return "✓"
	case "error":
		return "✗"
	case "pending":
		return "⏳"
	case "cancelled":
		return "⊘"
	default:
		return ""
	}
}

func dirArrow(dir string) string {
	if dir == "out" {
		return "←"
	}
	return "→"
}

func durStr(r *record) string {
	if r.durMS == 0 {
		return ""
	}
	return fmt.Sprintf("%.1fms", r.durMS)
}

func methodOr(r *record) string {
	if r.method != "" {
		return r.method
	}
	return "(" + r.kind + ")"
}

func orDash(s string) string {
	if s == "" {
		return "—"
	}
	return s
}
