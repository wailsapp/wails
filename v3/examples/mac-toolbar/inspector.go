package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
)

var daymarkCategories = []string{
	"Personal / Field Notes",
	"Ideas / Observations",
	"Personal / Drafts",
	"Work / Notes",
}

// daymarkInspector owns the native trailing inspector and translates between
// AppKit control callbacks and the editor's application model. It contains no
// HTML and creates no WebView.
type daymarkInspector struct {
	app       *application.App
	inspector *application.MacInspector
	pane      *application.MacSplitPane

	title       *application.MacInspectorControl
	category    *application.MacInspectorControl
	pinned      *application.MacInspectorControl
	words       *application.MacInspectorControl
	characters  *application.MacInspectorControl
	readingTime *application.MacInspectorControl
	status      *application.MacInspectorControl

	callbacksLock    sync.RWMutex
	onTitleChange    func(string)
	onCategoryChange func(string)
	onPinnedChange   func(bool)
}

type daymarkInspectorState struct {
	Title      string `json:"title"`
	Category   string `json:"category"`
	Pinned     bool   `json:"pinned"`
	Words      int    `json:"words"`
	Characters int    `json:"characters"`
	Reading    int    `json:"readingMinutes"`
	Dirty      bool   `json:"dirty"`
}

func newDaymarkInspector(app *application.App) *daymarkInspector {
	result := &daymarkInspector{
		app:       app,
		inspector: application.NewMacInspector(),
	}

	document := result.inspector.AddSection("Document")
	result.title = document.AddTextField("Title", "Saturday, slowly.").
		SetTooltip("Rename the current note")
	result.category = document.AddPopup("Category", daymarkCategories, 0).
		SetTooltip("Choose where this note belongs")
	result.pinned = document.AddCheckbox("Keep this note pinned", false)

	statistics := result.inspector.AddSection("Statistics")
	result.words = statistics.AddLabel("Words", "0")
	result.characters = statistics.AddLabel("Characters", "0")
	result.readingTime = statistics.AddLabel("Reading time", "1 min")

	state := result.inspector.AddSection("State")
	result.status = state.AddLabel("Changes", "Saved")

	result.title.OnTextChange(func(_ *application.Context, value string) {
		result.callbacksLock.RLock()
		callback := result.onTitleChange
		result.callbacksLock.RUnlock()
		if callback != nil {
			callback(value)
		}
	})
	result.category.OnSelectionChange(func(_ *application.Context, _ int, value string) {
		result.callbacksLock.RLock()
		callback := result.onCategoryChange
		result.callbacksLock.RUnlock()
		if callback != nil {
			callback(value)
		}
	})
	result.pinned.OnToggle(func(_ *application.Context, checked bool) {
		result.callbacksLock.RLock()
		callback := result.onPinnedChange
		result.callbacksLock.RUnlock()
		if callback != nil {
			callback(checked)
		}
	})

	app.Event.On("editor:inspector-state", result.handleEditorState)
	return result
}

func (i *daymarkInspector) NativeInspector() *application.MacInspector { return i.inspector }

func (i *daymarkInspector) SetPane(pane *application.MacSplitPane) {
	i.pane = pane
}

func (i *daymarkInspector) Toggle() {
	if i.pane != nil {
		i.pane.Toggle()
	}
}

func (i *daymarkInspector) SetCollapsed(collapsed bool) {
	if i.pane != nil {
		i.pane.SetCollapsed(collapsed)
	}
}

func (i *daymarkInspector) OnTitleChange(callback func(string)) {
	i.callbacksLock.Lock()
	i.onTitleChange = callback
	i.callbacksLock.Unlock()
}

func (i *daymarkInspector) OnCategoryChange(callback func(string)) {
	i.callbacksLock.Lock()
	i.onCategoryChange = callback
	i.callbacksLock.Unlock()
}

func (i *daymarkInspector) OnPinnedChange(callback func(bool)) {
	i.callbacksLock.Lock()
	i.onPinnedChange = callback
	i.callbacksLock.Unlock()
}

func (i *daymarkInspector) handleEditorState(event *application.CustomEvent) {
	raw, ok := event.Data.(string)
	if !ok {
		return
	}
	var state daymarkInspectorState
	if json.Unmarshal([]byte(raw), &state) != nil {
		return
	}
	i.applyState(state)
}

func (i *daymarkInspector) applyState(state daymarkInspectorState) {
	i.title.SetValue(state.Title)
	i.category.SetSelectedIndex(categoryIndex(state.Category))
	i.pinned.SetChecked(state.Pinned)
	i.words.SetValue(fmt.Sprintf("%d", state.Words))
	i.characters.SetValue(fmt.Sprintf("%d", state.Characters))
	i.readingTime.SetValue(fmt.Sprintf("%d min", max(1, state.Reading)))
	if state.Dirty {
		i.status.SetValue("Unsaved")
	} else {
		i.status.SetValue("Saved")
	}
}

func categoryIndex(category string) int {
	category = strings.TrimSpace(category)
	for index, candidate := range daymarkCategories {
		if candidate == category {
			return index
		}
	}
	return 0
}
