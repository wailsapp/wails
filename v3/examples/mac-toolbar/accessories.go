package main

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

// daymarkAccessories owns the native titlebar accessory: a secondary status
// label at the trailing edge of the unified toolbar that mirrors the
// editor's saved state. It is attached before the window exists and follows
// the same editor:dirty event the toolbar's Save button already observes.
// The sidebar's native filter strip lives in split.go beside the sidebar it
// filters.
type daymarkAccessories struct {
	status *application.MacAccessory
	label  *application.MacAccessoryControl
}

func newDaymarkAccessories(app *application.App, window *application.WebviewWindow) *daymarkAccessories {
	result := &daymarkAccessories{
		status: application.NewMacAccessory(application.MacAccessoryLayoutTrailing),
	}
	result.label = result.status.AddLabel("Saved").
		SetSymbol("checkmark.circle").
		SetTooltip("Whether the open note has unsaved changes")
	if err := window.AddTitlebarAccessory(result.status); err != nil {
		window.Error("titlebar accessory: %s", err)
	}

	app.Event.On("editor:dirty", func(event *application.CustomEvent) {
		if dirty, ok := event.Data.(bool); ok {
			result.setDirty(dirty)
		}
	})
	return result
}

func (a *daymarkAccessories) setDirty(dirty bool) {
	if dirty {
		a.label.SetText("Edited").SetSymbol("pencil.circle")
		return
	}
	a.label.SetText("Saved").SetSymbol("checkmark.circle")
}
