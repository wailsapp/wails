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
	toolbar  *application.MacToolbar
	share    *application.MacToolbarShareItem
	provider *daymarkShareProvider
	details  *application.MacToolbarItem
	save     *application.MacToolbarItem
	focus    *application.MacToolbarItem

	stateLock sync.Mutex
	focused   bool
	dirty     bool
	subject   string
}

func newDaymarkToolbar(app *application.App) *daymarkToolbar {
	result := &daymarkToolbar{
		app:     app,
		toolbar: application.NewMacToolbar(),
		provider: &daymarkShareProvider{note: sharePayload{
			Title:    "Saturday, slowly.",
			Subtitle: "A good day has room around it.",
			Body:     "Leave the phone at home. Walk until the city sounds different. Buy something warm on the way back.",
		}},
		subject: "Saturday, slowly.",
	}

	result.addItems()
	result.observeEditor()
	return result
}

// NativeToolbar returns the toolbar that should be attached to the window.
func (t *daymarkToolbar) NativeToolbar() *application.MacToolbar {
	return t.toolbar
}

func (t *daymarkToolbar) addItems() {
	newNote := t.toolbar.AddButton("New").SetSymbol("square.and.pencil").SetBordered(true)
	newNote.OnClick(func(*application.Context) {
		t.app.Event.Emit("toolbar:new")
	})

	mode := t.toolbar.AddGroup("Mode", application.ToolbarGroupSelectOne)
	mode.AddButton("Write").SetSymbol("pencil").OnClick(func(*application.Context) {
		mode.SetSelectedIndex(0)
		t.app.Event.Emit("toolbar:mode", "write")
	})
	mode.AddButton("Preview").SetSymbol("doc.richtext").OnClick(func(*application.Context) {
		mode.SetSelectedIndex(1)
		t.app.Event.Emit("toolbar:mode", "preview")
	})

	t.toolbar.AddFlexibleSpace()
	search := t.toolbar.AddSearch("Search notes").SetTooltip("Search your notes")
	search.OnSearch(func(_ *application.Context, query string) {
		t.app.Event.Emit("toolbar:search", query)
	})

	t.share = t.toolbar.AddShare("Share")
	t.share.SetTooltip("Share the current note with another app")
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

	t.details = t.toolbar.AddButton("Details").SetSymbol("info.circle").SetBordered(true)
	t.details.OnClick(func(*application.Context) {
		t.app.Event.Emit("toolbar:details")
	})

	t.save = t.toolbar.AddButton("Save").SetSymbol("checkmark.circle").SetProminent(true).SetBadgeCount(0)
	t.save.OnClick(func(*application.Context) {
		t.saveNote()
	})

	t.focus = t.toolbar.AddButton("Focus").SetSymbol("arrow.up.left.and.arrow.down.right").SetBordered(true)
	t.focus.OnClick(func(*application.Context) {
		t.toggleFocus()
	})
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

	t.save.SetBadgeCount(0).SetLabel("Saved")
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
		t.save.SetBadgeCount(1).SetLabel("Save")
		return
	}
	t.save.SetBadgeCount(0).SetLabel("Saved")
}

func (t *daymarkToolbar) toggleFocus() {
	t.stateLock.Lock()
	t.focused = !t.focused
	focused := t.focused
	t.stateLock.Unlock()

	t.details.SetHidden(focused)
	if focused {
		t.focus.SetLabel("Exit Focus").SetSymbol("arrow.down.right.and.arrow.up.left")
	} else {
		t.focus.SetLabel("Focus").SetSymbol("arrow.up.left.and.arrow.down.right")
	}
	t.app.Event.Emit("toolbar:focus", focused)
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
