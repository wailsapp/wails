package main

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type daymarkNote struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Body     string `json:"body"`
	Category string `json:"category"`
	Pinned   bool   `json:"pinned"`
}

// daymarkSplit owns the native NSSplitViewController layout, source-list
// sidebar, and trailing property inspector. The only WKWebView in the window
// is the primary editor pane.
type daymarkSplit struct {
	app         *application.App
	split       *application.MacSplitView
	sidebar     *application.MacSidebar
	sidePane    *application.MacSplitPane
	primaryPane *application.MacSplitWebviewPane
	inspector   *daymarkInspector
	inspectPane *application.MacSplitPane
	section     *application.MacSidebarSection
	items       []*application.MacSidebarItem
	notes       []daymarkNote
	active      int
	stateLock   sync.Mutex

	observersLock      sync.Mutex
	observers          []func(collapsed bool)
	inspectorObservers []func(collapsed bool)
}

func newDaymarkSplit(app *application.App) *daymarkSplit {
	result := &daymarkSplit{
		app:     app,
		split:   application.NewMacSplitView().SetAutosaveName("daymark.main-window"),
		sidebar: application.NewMacSidebar(),
		notes: []daymarkNote{
			{Title: "Saturday, slowly.", Subtitle: "A slow day is still a day well spent.", Body: `A good day has room around it.

Leave the phone at home. Walk until the city sounds different. Buy something warm on the way back.

# The long way home

Take the street with the old trees, even though it adds twenty minutes. Notice the windows left open above the bakery and the bicycles leaning against every second fence.

At the market, choose the peaches by scent instead of colour. Carry them carefully. Let the paper bag warm in your hands.

There is no prize for arriving early today. Sit near the water until the light changes, and allow the afternoon to become evening without asking it to be useful.

When you finally turn home, leave enough quiet in the day to hear your own footsteps.`, Category: "Personal / Field Notes"},
			{Title: "Things worth noticing", Subtitle: "A list for paying closer attention.", Body: "The first warm patch of sunlight on the kitchen floor.\n\nThe sound of a neighbour watering plants.\n\nA good question asked at exactly the right time.\n\nSteam gathering at the edge of a café window.\n\nThe moment a familiar street looks new after rain.", Category: "Ideas / Observations"},
			{Title: "A smaller promise", Subtitle: "Start with what can be carried today.", Body: "Make the bed. Drink the water. Send the message.\n\nThen see what the day is willing to become.\n\nDo one thing slowly enough to notice that you are doing it.\n\nLeave tomorrow somewhere to begin.", Category: "Personal / Drafts"},
		},
	}

	result.section = result.sidebar.AddSection("Notes")
	for index := range result.notes {
		result.addNativeNoteItem(index)
	}
	result.sidebar.SetSelectedItem(result.items[0])

	result.sidePane = result.split.AddSidebar(result.sidebar)
	result.sidePane.
		SetMinimumThickness(210).
		SetMaximumThickness(340).
		SetCollapsible(true)
	result.primaryPane = result.split.AddPrimaryContent().
		SetContentLayout(application.MacContentLayoutEdgeToEdge)
	result.inspector = newDaymarkInspector(app)
	result.inspectPane = result.split.AddInspector(result.inspector.NativeInspector())
	result.inspectPane.
		SetMinimumThickness(240).
		SetMaximumThickness(360).
		SetPreferredThicknessFraction(.24).
		SetCollapsible(true).
		SetCanCollapseFromWindowResize(false)
	result.inspector.SetPane(result.inspectPane)
	result.inspector.OnTitleChange(result.renameActiveNote)
	result.inspector.OnCategoryChange(result.changeActiveCategory)
	result.inspector.OnPinnedChange(result.changeActivePinned)

	result.sidePane.OnCollapsedChange(func(_ *application.Context, collapsed bool) {
		result.observersLock.Lock()
		observers := append([]func(bool){}, result.observers...)
		result.observersLock.Unlock()
		for _, observer := range observers {
			observer(collapsed)
		}
	})
	result.inspectPane.OnCollapsedChange(func(_ *application.Context, collapsed bool) {
		result.observersLock.Lock()
		observers := append([]func(bool){}, result.inspectorObservers...)
		result.observersLock.Unlock()
		for _, observer := range observers {
			observer(collapsed)
		}
	})

	app.Event.On("editor:ready", func(*application.CustomEvent) {
		result.emitActiveNote()
	})
	app.Event.On("editor:note-updated", result.handleNoteUpdated)
	return result
}

func (s *daymarkSplit) addNativeNoteItem(index int) {
	s.stateLock.Lock()
	note := s.notes[index]
	s.stateLock.Unlock()
	item := s.section.AddItem(note.Title).
		SetSymbol("doc.text").
		SetTooltip(note.Subtitle)
	item.OnClick(func(*application.Context) {
		s.stateLock.Lock()
		s.active = index
		note := s.notes[index]
		s.stateLock.Unlock()
		s.app.Event.Emit("sidebar:note-selected", map[string]any{"index": index, "note": note})
	})
	s.stateLock.Lock()
	s.items = append(s.items, item)
	s.stateLock.Unlock()
}

func (s *daymarkSplit) emitActiveNote() {
	s.stateLock.Lock()
	index := s.active
	note := s.notes[index]
	s.stateLock.Unlock()
	s.app.Event.Emit("sidebar:note-selected", map[string]any{"index": index, "note": note})
}

func (s *daymarkSplit) handleNoteUpdated(event *application.CustomEvent) {
	raw, ok := event.Data.(string)
	if !ok {
		return
	}
	var update struct {
		Index    int    `json:"index"`
		Title    string `json:"title"`
		Subtitle string `json:"subtitle"`
		Body     string `json:"body"`
		Category string `json:"category"`
		Pinned   bool   `json:"pinned"`
	}
	if json.Unmarshal([]byte(raw), &update) != nil {
		return
	}
	s.stateLock.Lock()
	if update.Index < 0 || update.Index >= len(s.notes) {
		s.stateLock.Unlock()
		return
	}
	title := update.Title
	if strings.TrimSpace(title) == "" {
		title = "Untitled note"
	}
	s.notes[update.Index].Title = title
	s.notes[update.Index].Subtitle = strings.TrimSpace(update.Subtitle)
	s.notes[update.Index].Body = update.Body
	if category := strings.TrimSpace(update.Category); category != "" {
		s.notes[update.Index].Category = category
	}
	s.notes[update.Index].Pinned = update.Pinned
	item := s.items[update.Index]
	s.stateLock.Unlock()
	item.SetLabel(title).SetTooltip(update.Subtitle)
	if update.Pinned {
		item.SetSymbol("pin.fill")
	} else {
		item.SetSymbol("doc.text")
	}
}

func (s *daymarkSplit) NewNote() {
	s.stateLock.Lock()
	index := len(s.notes)
	s.notes = append(s.notes, daymarkNote{
		Title: "Untitled note", Subtitle: "Start with one honest sentence.", Category: "Personal / Drafts",
	})
	s.stateLock.Unlock()
	s.addNativeNoteItem(index)
	s.stateLock.Lock()
	s.active = index
	item := s.items[index]
	s.stateLock.Unlock()
	s.sidebar.SetSelectedItem(item)
	s.emitActiveNote()
}

func (s *daymarkSplit) Filter(query string) {
	query = strings.ToLower(strings.TrimSpace(query))
	s.stateLock.Lock()
	notes := append([]daymarkNote(nil), s.notes...)
	items := append([]*application.MacSidebarItem(nil), s.items...)
	s.stateLock.Unlock()
	for index, note := range notes {
		haystack := strings.ToLower(note.Title + " " + note.Subtitle + " " + note.Body)
		items[index].SetHidden(query != "" && !strings.Contains(haystack, query))
	}
}

func (s *daymarkSplit) NativeSplitView() *application.MacSplitView { return s.split }
func (s *daymarkSplit) SetContentLayout(layout application.MacContentLayout) {
	s.primaryPane.SetContentLayout(layout)
	if layout == application.MacContentLayoutBelowToolbar {
		s.app.Event.Emit("layout:content", "below-toolbar")
		return
	}
	s.app.Event.Emit("layout:content", "edge-to-edge")
}
func (s *daymarkSplit) ToggleSidebar()                       { s.sidePane.Toggle() }
func (s *daymarkSplit) SetSidebarCollapsed(collapsed bool)   { s.sidePane.SetCollapsed(collapsed) }
func (s *daymarkSplit) ToggleInspector()                     { s.inspector.Toggle() }
func (s *daymarkSplit) SetInspectorCollapsed(collapsed bool) { s.inspector.SetCollapsed(collapsed) }

func (s *daymarkSplit) renameActiveNote(title string) {
	if strings.TrimSpace(title) == "" {
		title = "Untitled note"
	}
	s.stateLock.Lock()
	index := s.active
	s.notes[index].Title = title
	item := s.items[index]
	s.stateLock.Unlock()
	item.SetLabel(title)
	s.app.Event.Emit("inspector:title-changed", title)
}

func (s *daymarkSplit) changeActiveCategory(category string) {
	s.stateLock.Lock()
	s.notes[s.active].Category = category
	s.stateLock.Unlock()
	s.app.Event.Emit("inspector:category-changed", category)
}

func (s *daymarkSplit) changeActivePinned(pinned bool) {
	s.stateLock.Lock()
	index := s.active
	s.notes[index].Pinned = pinned
	item := s.items[index]
	s.stateLock.Unlock()
	if pinned {
		item.SetSymbol("pin.fill")
	} else {
		item.SetSymbol("doc.text")
	}
	s.app.Event.Emit("inspector:pinned-changed", pinned)
}

func (s *daymarkSplit) OnSidebarCollapsedChange(callback func(collapsed bool)) {
	s.observersLock.Lock()
	s.observers = append(s.observers, callback)
	s.observersLock.Unlock()
}

func (s *daymarkSplit) OnInspectorCollapsedChange(callback func(collapsed bool)) {
	s.observersLock.Lock()
	s.inspectorObservers = append(s.inspectorObservers, callback)
	s.observersLock.Unlock()
}
