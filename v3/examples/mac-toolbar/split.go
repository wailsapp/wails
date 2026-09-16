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
	Priority int    `json:"priority"`
	// Tint colours the note's sidebar symbol. It stays native-only.
	Tint *application.RGBA `json:"-"`
}

// daymarkTag is one nested Tags row. The prefix matches note categories such
// as "Personal / Drafts".
type daymarkTag struct {
	label  string
	symbol string
	tint   application.RGBA
	items  []string
}

var daymarkTags = []daymarkTag{
	{label: "Personal", symbol: "person", tint: application.NewRGB(255, 149, 0), items: []string{"Field Notes", "Drafts"}},
	{label: "Ideas", symbol: "lightbulb", tint: application.NewRGB(255, 204, 0), items: []string{"Observations"}},
	{label: "Work", symbol: "briefcase", tint: application.NewRGB(0, 122, 255), items: []string{"Notes"}},
}

// daymarkSplit owns the native NSSplitViewController layout, source-list
// sidebar, and trailing property inspector. The only WKWebView in the window
// is the primary editor pane.
//
// The sidebar shows three native features beyond a flat list: an "All Notes"
// root row with a badge, a reorderable and renameable Notes section with a
// per-row context menu, and a Tags section of nested rows whose badges count
// the notes in each category.
type daymarkSplit struct {
	app         *application.App
	split       *application.MacSplitView
	sidebar     *application.MacSidebar
	sidePane    *application.MacSplitPane
	primaryPane *application.MacSplitWebviewPane
	inspector   *daymarkInspector
	inspectPane *application.MacSplitPane
	section     *application.MacSidebarSection
	allNotes    *application.MacSidebarItem
	tagRows     map[string]*application.MacSidebarItem
	items       []*application.MacSidebarItem
	notes       []daymarkNote
	active      int
	query       string
	category    string
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
		tagRows: map[string]*application.MacSidebarItem{},
		notes: []daymarkNote{
			{Title: "Saturday, slowly.", Subtitle: "A slow day is still a day well spent.", Body: `A good day has room around it.

Leave the phone at home. Walk until the city sounds different. Buy something warm on the way back.

# The long way home

Take the street with the old trees, even though it adds twenty minutes. Notice the windows left open above the bakery and the bicycles leaning against every second fence.

At the market, choose the peaches by scent instead of colour. Carry them carefully. Let the paper bag warm in your hands.

There is no prize for arriving early today. Sit near the water until the light changes, and allow the afternoon to become evening without asking it to be useful.

When you finally turn home, leave enough quiet in the day to hear your own footsteps.`, Category: "Personal / Field Notes", Priority: 2},
			{Title: "Things worth noticing", Subtitle: "A list for paying closer attention.", Body: "The first warm patch of sunlight on the kitchen floor.\n\nThe sound of a neighbour watering plants.\n\nA good question asked at exactly the right time.\n\nSteam gathering at the edge of a café window.\n\nThe moment a familiar street looks new after rain.", Category: "Ideas / Observations"},
			{Title: "A smaller promise", Subtitle: "Start with what can be carried today.", Body: "Make the bed. Drink the water. Send the message.\n\nThen see what the day is willing to become.\n\nDo one thing slowly enough to notice that you are doing it.\n\nLeave tomorrow somewhere to begin.", Category: "Personal / Drafts", Priority: 1},
		},
	}

	// A root row above the sections, badged with the note count.
	result.allNotes = result.sidebar.AddItem("All Notes").
		SetSymbol("tray.full").
		SetTooltip("Show every note").
		OnClick(func(*application.Context) { result.setCategory("") })

	result.section = result.sidebar.AddSection("Notes")
	for index := range result.notes {
		result.addNativeNoteItem(index)
	}
	result.sidebar.SetSelectedItem(result.items[0])

	// Nested rows: each tag expands into its categories. Badges count notes.
	tags := result.sidebar.AddSection("Tags")
	for _, tag := range daymarkTags {
		tint := tag.tint
		parent := tags.AddItem(tag.label).
			SetSymbol(tag.symbol).
			SetTintColor(&tint).
			SetExpanded(true)
		prefix := tag.label
		parent.OnClick(func(*application.Context) { result.setCategory(prefix) })
		result.tagRows[prefix] = parent
		for _, child := range tag.items {
			category := tag.label + " / " + child
			row := parent.AddItem(child).SetSymbol("tag")
			row.OnClick(func(*application.Context) { result.setCategory(category) })
			result.tagRows[category] = row
		}
	}
	result.refreshBadges()

	// Drag reorder within the Notes section keeps the Go model in step.
	result.sidebar.SetReorderable(true)
	result.sidebar.OnMove(result.handleNoteMoved)

	// Right-click: notes get Pin/Unpin and Delete; the empty area offers New Note.
	result.sidebar.OnContextMenu(result.contextMenuForRow)
	emptyArea := application.NewMenu()
	emptyArea.Add("New Note").OnClick(func(*application.Context) { result.NewNote() })
	result.sidebar.SetContextMenu(emptyArea)

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
	result.inspector.OnPriorityChange(result.changeActivePriority)
	result.inspector.OnTintChange(result.changeActiveTint)
	result.pushSidebarStateToInspector()

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

// addNativeNoteItem appends a renameable, badged note row. Row callbacks look
// the note up by handle so drag reordering never leaves a stale index behind.
func (s *daymarkSplit) addNativeNoteItem(index int) {
	s.stateLock.Lock()
	note := s.notes[index]
	s.stateLock.Unlock()
	item := s.section.AddItem(note.Title).
		SetSymbol(noteSymbol(note)).
		SetTooltip(note.Subtitle).
		SetBadge(note.Priority).
		SetTintColor(note.Tint).
		SetEditable(true)
	item.OnClick(func(*application.Context) {
		s.stateLock.Lock()
		index := s.indexOfItemLocked(item)
		if index < 0 {
			s.stateLock.Unlock()
			return
		}
		s.active = index
		note := s.notes[index]
		s.stateLock.Unlock()
		s.pushSidebarStateToInspector()
		s.app.Event.Emit("sidebar:note-selected", map[string]any{"index": index, "note": note})
	})
	item.OnRename(func(_ *application.Context, title string) {
		s.renameNote(item, title)
	})
	s.stateLock.Lock()
	s.items = append(s.items, item)
	s.stateLock.Unlock()
}

func noteSymbol(note daymarkNote) string {
	if note.Pinned {
		return "pin.fill"
	}
	return "doc.text"
}

func (s *daymarkSplit) indexOfItemLocked(item *application.MacSidebarItem) int {
	for index, candidate := range s.items {
		if candidate == item {
			return index
		}
	}
	return -1
}

// contextMenuForRow builds the right-click menu for one sidebar row. Tag rows
// return nil so the sidebar's fallback menu appears instead.
func (s *daymarkSplit) contextMenuForRow(_ *application.Context, item *application.MacSidebarItem) *application.Menu {
	s.stateLock.Lock()
	index := s.indexOfItemLocked(item)
	var note daymarkNote
	if index >= 0 {
		note = s.notes[index]
	}
	s.stateLock.Unlock()
	if index < 0 {
		return nil
	}
	menu := application.NewMenu()
	pinLabel := "Pin Note"
	if note.Pinned {
		pinLabel = "Unpin Note"
	}
	menu.Add(pinLabel).OnClick(func(*application.Context) { s.setPinned(item, !note.Pinned) })
	menu.Add("Clear Priority").SetEnabled(note.Priority > 0).OnClick(func(*application.Context) {
		s.setPriority(item, 0)
	})
	menu.AddSeparator()
	menu.Add("Delete Note").OnClick(func(*application.Context) { s.deleteNote(item) })
	return menu
}

// handleNoteMoved mirrors a native drag reorder into the notes slice.
func (s *daymarkSplit) handleNoteMoved(_ *application.Context, item *application.MacSidebarItem, section *application.MacSidebarSection, newIndex int) {
	if section != s.section {
		return
	}
	s.stateLock.Lock()
	oldIndex := s.indexOfItemLocked(item)
	if oldIndex < 0 || newIndex < 0 || newIndex >= len(s.items) {
		s.stateLock.Unlock()
		return
	}
	activeItem := s.items[s.active]
	note := s.notes[oldIndex]
	s.notes = append(s.notes[:oldIndex], s.notes[oldIndex+1:]...)
	s.items = append(s.items[:oldIndex], s.items[oldIndex+1:]...)
	s.notes = append(s.notes[:newIndex], append([]daymarkNote{note}, s.notes[newIndex:]...)...)
	s.items = append(s.items[:newIndex], append([]*application.MacSidebarItem{item}, s.items[newIndex:]...)...)
	s.active = s.indexOfItemLocked(activeItem)
	s.stateLock.Unlock()
	// The editor addresses notes by index, so tell it where the open note went.
	s.emitActiveNote()
}

func (s *daymarkSplit) renameNote(item *application.MacSidebarItem, title string) {
	if strings.TrimSpace(title) == "" {
		title = "Untitled note"
		item.SetLabel(title)
	}
	s.stateLock.Lock()
	index := s.indexOfItemLocked(item)
	if index < 0 {
		s.stateLock.Unlock()
		return
	}
	s.notes[index].Title = title
	isActive := index == s.active
	s.stateLock.Unlock()
	if isActive {
		s.inspector.SetTitle(title)
		s.app.Event.Emit("inspector:title-changed", title)
	}
}

func (s *daymarkSplit) setPinned(item *application.MacSidebarItem, pinned bool) {
	s.stateLock.Lock()
	index := s.indexOfItemLocked(item)
	if index < 0 {
		s.stateLock.Unlock()
		return
	}
	s.notes[index].Pinned = pinned
	isActive := index == s.active
	s.stateLock.Unlock()
	item.SetSymbol(noteSymbol(daymarkNote{Pinned: pinned}))
	if isActive {
		s.inspector.SetPinned(pinned)
		s.app.Event.Emit("inspector:pinned-changed", pinned)
	}
}

func (s *daymarkSplit) setPriority(item *application.MacSidebarItem, priority int) {
	s.stateLock.Lock()
	index := s.indexOfItemLocked(item)
	if index < 0 {
		s.stateLock.Unlock()
		return
	}
	s.notes[index].Priority = priority
	isActive := index == s.active
	s.stateLock.Unlock()
	item.SetBadge(priority)
	if isActive {
		s.pushSidebarStateToInspector()
	}
}

// deleteNote removes a note and its native row. The last note is never
// deleted so the editor always has a document.
func (s *daymarkSplit) deleteNote(item *application.MacSidebarItem) {
	s.stateLock.Lock()
	index := s.indexOfItemLocked(item)
	if index < 0 || len(s.notes) == 1 {
		s.stateLock.Unlock()
		return
	}
	s.notes = append(s.notes[:index], s.notes[index+1:]...)
	s.items = append(s.items[:index], s.items[index+1:]...)
	if s.active >= len(s.items) {
		s.active = len(s.items) - 1
	} else if s.active > index {
		s.active--
	}
	next := s.items[s.active]
	s.stateLock.Unlock()
	s.section.Remove(item)
	s.sidebar.SetSelectedItem(next)
	s.refreshBadges()
	s.pushSidebarStateToInspector()
	s.emitActiveNote()
}

// refreshBadges recounts the All Notes and Tags badges.
func (s *daymarkSplit) refreshBadges() {
	s.stateLock.Lock()
	notes := append([]daymarkNote(nil), s.notes...)
	s.stateLock.Unlock()
	s.allNotes.SetBadge(len(notes))
	for prefix, row := range s.tagRows {
		count := 0
		for _, note := range notes {
			if strings.HasPrefix(note.Category, prefix) {
				count++
			}
		}
		row.SetBadge(count)
	}
}

func (s *daymarkSplit) setCategory(prefix string) {
	s.stateLock.Lock()
	s.category = prefix
	s.stateLock.Unlock()
	s.applyFilters()
}

// applyFilters hides note rows that match neither the search query nor the
// selected tag. Hidden rows stay in the model and keep their handles.
func (s *daymarkSplit) applyFilters() {
	s.stateLock.Lock()
	query, category := s.query, s.category
	notes := append([]daymarkNote(nil), s.notes...)
	items := append([]*application.MacSidebarItem(nil), s.items...)
	s.stateLock.Unlock()
	for index, note := range notes {
		haystack := strings.ToLower(note.Title + " " + note.Subtitle + " " + note.Body)
		hidden := (query != "" && !strings.Contains(haystack, query)) ||
			(category != "" && !strings.HasPrefix(note.Category, category))
		items[index].SetHidden(hidden)
	}
}

func (s *daymarkSplit) pushSidebarStateToInspector() {
	s.stateLock.Lock()
	note := s.notes[s.active]
	s.stateLock.Unlock()
	s.inspector.SetSidebarState(note.Priority, note.Tint)
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
	note := s.notes[update.Index]
	item := s.items[update.Index]
	s.stateLock.Unlock()
	item.SetLabel(title).SetTooltip(update.Subtitle).SetSymbol(noteSymbol(note))
	s.refreshBadges()
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
	s.refreshBadges()
	s.pushSidebarStateToInspector()
	s.emitActiveNote()
}

func (s *daymarkSplit) Filter(query string) {
	s.stateLock.Lock()
	s.query = strings.ToLower(strings.TrimSpace(query))
	s.stateLock.Unlock()
	s.applyFilters()
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
	s.refreshBadges()
	s.app.Event.Emit("inspector:category-changed", category)
}

func (s *daymarkSplit) changeActivePinned(pinned bool) {
	s.stateLock.Lock()
	index := s.active
	s.notes[index].Pinned = pinned
	item := s.items[index]
	s.stateLock.Unlock()
	item.SetSymbol(noteSymbol(daymarkNote{Pinned: pinned}))
	s.app.Event.Emit("inspector:pinned-changed", pinned)
}

// changeActivePriority shows the inspector slider's value as the row badge.
func (s *daymarkSplit) changeActivePriority(priority int) {
	s.stateLock.Lock()
	index := s.active
	s.notes[index].Priority = priority
	item := s.items[index]
	s.stateLock.Unlock()
	item.SetBadge(priority)
}

// changeActiveTint colours the row symbol from the inspector colour well.
func (s *daymarkSplit) changeActiveTint(tint *application.RGBA) {
	s.stateLock.Lock()
	index := s.active
	s.notes[index].Tint = tint
	item := s.items[index]
	s.stateLock.Unlock()
	item.SetTintColor(tint)
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
