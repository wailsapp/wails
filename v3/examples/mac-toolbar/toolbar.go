package main

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// daymarkToolbar owns the native toolbar and the small amount of UI state
// shared by its controls. Keeping it separate from main makes the example's
// application setup easy to scan while leaving the complete toolbar lifecycle
// in one place.
type daymarkToolbar struct {
	app      *application.App
	split    *daymarkSplit
	toolbar  *application.MacToolbar
	share    *application.MacToolbarShareItem
	provider *daymarkShareProvider
	save     *application.MacToolbarItem
	focus    *application.MacToolbarItem

	stateLock sync.Mutex
	focused   bool
	dirty     bool
	subject   string
}

func newDaymarkToolbar(app *application.App, split *daymarkSplit) *daymarkToolbar {
	result := &daymarkToolbar{
		app:     app,
		split:   split,
		toolbar: application.NewMacToolbar().SetDisplayMode(application.MacToolbarDisplayModeIconOnly),
		provider: newDaymarkPDFShareProvider(sharePayload{
			Title:    "Saturday, slowly.",
			Subtitle: "A good day has room around it.",
			Body:     "Leave the phone at home. Walk until the city sounds different. Buy something warm on the way back.",
		}),
		subject: "Saturday, slowly.",
	}

	result.addItems()
	result.observeEditor()
	result.observeNativePanes()
	return result
}

// NativeToolbar returns the toolbar that should be attached to the window.
func (t *daymarkToolbar) NativeToolbar() *application.MacToolbar {
	return t.toolbar
}

func (t *daymarkToolbar) addItems() {
	// The leading section sits above the native sidebar: AppKit keeps every
	// item before the tracking separator aligned with the sidebar as its
	// divider moves.
	t.toolbar.AddSidebarToggle()

	newNote := t.toolbar.AddButton("New").SetSymbol("square.and.pencil").SetBordered(true)
	newNote.OnClick(func(*application.Context) {
		t.split.NewNote()
	})

	// Everything after the tracking separator operates on the editor content.
	t.toolbar.AddSidebarTrackingSeparator()

	// Search belongs to the primary toolbar section. NSSearchToolbarItem is
	// horizontally expandable; placing it in the sidebar section would allow
	// it to outgrow that pane as the divider moves.
	search := t.toolbar.AddSearch("Search notes").SetTooltip("Search your notes")
	search.OnSearch(func(_ *application.Context, query string) {
		t.split.Filter(query)
	})

	mode := t.toolbar.AddGroup("Mode", application.ToolbarGroupSelectOne)
	mode.SetBordered(true)
	mode.AddButton("Write").SetSymbol("pencil").OnClick(func(*application.Context) {
		mode.SetSelectedIndex(0)
		t.app.Event.Emit("toolbar:mode", "write")
	})
	mode.AddButton("Preview").SetSymbol("doc.richtext").OnClick(func(*application.Context) {
		mode.SetSelectedIndex(1)
		t.app.Event.Emit("toolbar:mode", "preview")
	})

	t.toolbar.AddFlexibleSpace()

	t.share = t.toolbar.AddShare("Share PDF")
	t.share.SetTooltip("Share the current note as a PDF")
	t.share.SetBordered(true)
	t.share.SetProvider(t.provider).SetSubject(t.subject).SetSuggestedName("Daymark Note")
	t.share.OnShared(func(_ *application.Context, service string) {
		t.app.Event.Emit("toolbar:share-complete", service)
	})
	t.share.OnShareError(func(_ *application.Context, service string, err error) {
		t.app.Event.Emit("toolbar:share-error", map[string]string{
			"service": service,
			"message": err.Error(),
		})
	})

	// Related document actions share one native glass group. AppKit owns the
	// capsule, hit testing, highlighting, and Liquid Glass presentation.
	documentActions := t.toolbar.AddGroup("Document", application.ToolbarGroupMomentary)
	documentActions.SetBordered(true)
	t.save = documentActions.AddButton("Save").SetSymbol("checkmark.circle").SetBordered(true).SetBadgeCount(0)
	t.save.OnClick(func(*application.Context) {
		t.saveNote()
	})

	t.focus = documentActions.AddButton("Focus").SetSymbol("arrow.up.left.and.arrow.down.right").SetBordered(true)
	t.focus.OnClick(func(*application.Context) {
		t.toggleFocus()
	})

	// The trailing section tracks the native inspector divider. AppKit owns
	// both standard identifiers on macOS 14+; Wails supplies a native toggle
	// fallback on earlier releases.
	t.toolbar.AddInspectorTrackingSeparator()
	t.toolbar.AddInspectorToggle()
}

func (t *daymarkToolbar) observeEditor() {
	t.app.Event.On("editor:dirty", func(event *application.CustomEvent) {
		isDirty, ok := event.Data.(bool)
		if ok {
			t.setDirty(isDirty)
		}
	})

	t.app.Event.On("editor:share-content", func(event *application.CustomEvent) {
		raw, ok := event.Data.(string)
		if !ok {
			return
		}

		var note sharePayload
		if err := json.Unmarshal([]byte(raw), &note); err != nil {
			return
		}

		t.provider.update(note)
		t.setShareSubject(note.Title)
	})
}

func (t *daymarkToolbar) saveNote() {
	t.stateLock.Lock()
	t.dirty = false
	t.stateLock.Unlock()

	t.save.SetBadgeCount(0).SetLabel("Saved").SetProminent(false)
	t.app.Event.Emit("toolbar:save")
}

func (t *daymarkToolbar) setDirty(dirty bool) {
	t.stateLock.Lock()
	if dirty == t.dirty {
		t.stateLock.Unlock()
		return
	}
	t.dirty = dirty
	t.stateLock.Unlock()

	if dirty {
		t.save.SetBadgeCount(1).SetLabel("Save").SetProminent(true)
		return
	}
	t.save.SetBadgeCount(0).SetLabel("Saved").SetProminent(false)
}

func (t *daymarkToolbar) toggleFocus() {
	t.stateLock.Lock()
	t.focused = !t.focused
	focused := t.focused
	t.stateLock.Unlock()

	t.applyFocusPresentation(focused)
	// Focus mode collapses the native sidebar instead of hiding an HTML
	// column; the collapse animation and divider state are AppKit's.
	t.split.SetSidebarCollapsed(focused)
	if focused {
		t.split.SetInspectorCollapsed(true)
	}
	t.app.Event.Emit("toolbar:focus", focused)
}

// observeNativePanes keeps focus mode synchronized with native collapse state:
// expanding either auxiliary pane while focused leaves focus mode.
func (t *daymarkToolbar) observeNativePanes() {
	t.split.OnSidebarCollapsedChange(func(collapsed bool) {
		t.exitFocusWhenPaneExpands(collapsed)
	})
	t.split.OnInspectorCollapsedChange(func(collapsed bool) {
		t.exitFocusWhenPaneExpands(collapsed)
	})
}

func (t *daymarkToolbar) exitFocusWhenPaneExpands(collapsed bool) {
	if collapsed {
		return
	}
	t.stateLock.Lock()
	wasFocused := t.focused
	t.focused = false
	t.stateLock.Unlock()
	if !wasFocused {
		return
	}
	t.applyFocusPresentation(false)
	t.app.Event.Emit("toolbar:focus", false)
}

func (t *daymarkToolbar) applyFocusPresentation(focused bool) {
	if focused {
		t.focus.SetLabel("Exit Focus").SetSymbol("arrow.down.right.and.arrow.up.left")
	} else {
		t.focus.SetLabel("Focus").SetSymbol("arrow.up.left.and.arrow.down.right")
	}
}

func (t *daymarkToolbar) setShareSubject(title string) {
	title = strings.TrimSpace(title)
	if title == "" {
		title = "Untitled note"
	}

	t.stateLock.Lock()
	changed := t.subject != title
	t.subject = title
	t.stateLock.Unlock()
	if changed {
		t.share.SetSubject(title)
	}
}
