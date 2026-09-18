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
//
// The Sidebar section is collapsible and drives the native source list
// directly: its slider becomes the active row's badge and its colour well
// tints the row symbol.
type daymarkInspector struct {
	app       *application.App
	inspector *application.MacInspector
	pane      *application.MacSplitPane

	title       *application.MacInspectorControl
	category    *application.MacInspectorControl
	pinned      *application.MacInspectorControl
	priority    *application.MacInspectorControl
	tint        *application.MacInspectorControl
	useTint     *application.MacInspectorControl
	reset       *application.MacInspectorControl
	words       *application.MacInspectorControl
	characters  *application.MacInspectorControl
	readingTime *application.MacInspectorControl
	status      *application.MacInspectorControl

	callbacksLock    sync.RWMutex
	onTitleChange    func(string)
	onCategoryChange func(string)
	onPinnedChange   func(bool)
	onPriorityChange func(int)
	onTintChange     func(*application.RGBA)
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

var daymarkDefaultTint = application.NewRGB(0, 122, 255)

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

	sidebar := result.inspector.AddSection("Sidebar").SetCollapsible(true)
	result.priority = sidebar.AddSlider("Priority", 0, 5, 0).
		SetTooltip("Shown as the badge on the note's sidebar row")
	result.useTint = sidebar.AddCheckbox("Tint the sidebar symbol", false)
	result.tint = sidebar.AddColorWell("Tint", daymarkDefaultTint).
		SetTooltip("Colour of the note's sidebar symbol").
		SetEnabled(false)
	result.reset = sidebar.AddButton("Reset Sidebar Row")

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
	result.priority.OnValueChange(func(_ *application.Context, value float64) {
		result.notifyPriority(int(value + 0.5))
	})
	result.useTint.OnToggle(func(_ *application.Context, checked bool) {
		result.tint.SetEnabled(checked)
		result.notifyTint()
	})
	result.tint.OnColorChange(func(*application.Context, application.RGBA) {
		result.notifyTint()
	})
	result.reset.OnClick(func(*application.Context) {
		result.priority.SetFloatValue(0)
		result.useTint.SetChecked(false)
		result.tint.SetColor(daymarkDefaultTint).SetEnabled(false)
		result.notifyPriority(0)
		result.notifyTint()
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

// OnPriorityChange observes the Sidebar section's slider as a whole number.
func (i *daymarkInspector) OnPriorityChange(callback func(int)) {
	i.callbacksLock.Lock()
	i.onPriorityChange = callback
	i.callbacksLock.Unlock()
}

// OnTintChange observes the colour well. It receives nil when tinting is off.
func (i *daymarkInspector) OnTintChange(callback func(*application.RGBA)) {
	i.callbacksLock.Lock()
	i.onTintChange = callback
	i.callbacksLock.Unlock()
}

// SetTitle and SetPinned reflect changes made from the sidebar (inline rename
// and the row context menu) without invoking the inspector callbacks.
func (i *daymarkInspector) SetTitle(title string) { i.title.SetValue(title) }
func (i *daymarkInspector) SetPinned(pinned bool) { i.pinned.SetChecked(pinned) }

// SetSidebarState shows the active note's badge and tint in the Sidebar
// section. Programmatic setters never fire callbacks.
func (i *daymarkInspector) SetSidebarState(priority int, tint *application.RGBA) {
	i.priority.SetFloatValue(float64(priority))
	i.useTint.SetChecked(tint != nil)
	i.tint.SetEnabled(tint != nil)
	if tint != nil {
		i.tint.SetColor(*tint)
	}
}

func (i *daymarkInspector) notifyPriority(priority int) {
	i.callbacksLock.RLock()
	callback := i.onPriorityChange
	i.callbacksLock.RUnlock()
	if callback != nil {
		callback(priority)
	}
}

func (i *daymarkInspector) notifyTint() {
	i.callbacksLock.RLock()
	callback := i.onTintChange
	i.callbacksLock.RUnlock()
	if callback == nil {
		return
	}
	if !i.useTint.Checked() {
		callback(nil)
		return
	}
	colour := i.tint.Color()
	callback(&colour)
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
