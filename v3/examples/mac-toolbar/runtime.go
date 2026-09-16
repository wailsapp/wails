package main

import (
	"fmt"
	"sync/atomic"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// daymarkRuntimeWindows creates windows after App.Run, from the Window menu.
// Both paths use options-time construction: the split view and toolbar are
// part of the window options, so a window created from a callback in a
// running application comes up with its native chrome in one step. The same
// layouts could also be installed later with SetSplitView, which installs
// into an already created window on the application thread.
type daymarkRuntimeWindows struct {
	app   *application.App
	count atomic.Int32
}

func newDaymarkRuntimeWindows(app *application.App) *daymarkRuntimeWindows {
	return &daymarkRuntimeWindows{app: app}
}

// NewSplitWindow opens a second WebView window whose native sidebar sits
// beside the Daymark editor page.
func (r *daymarkRuntimeWindows) NewSplitWindow() *application.WebviewWindow {
	n := r.count.Add(1)
	var window *application.WebviewWindow

	sidebar := application.NewMacSidebar()
	section := sidebar.AddSection("Scratch")
	for _, title := range []string{"Draft", "Outline", "Reading list"} {
		item := section.AddItem(title).SetSymbol("doc.text")
		item.OnClick(func(*application.Context) {
			if window != nil {
				window.SetTitle(fmt.Sprintf("Scratch %d: %s", n, title))
			}
		})
	}

	split := application.NewMacSplitView()
	split.AddSidebar(sidebar).
		SetMinimumThickness(200).
		SetMaximumThickness(320).
		SetCollapsible(true)
	split.AddPrimaryContent().SetContentLayout(application.MacContentLayoutEdgeToEdge)

	toolbar := application.NewMacToolbar().SetDisplayMode(application.MacToolbarDisplayModeIconOnly)
	toolbar.AddSidebarToggle()
	toolbar.AddSidebarTrackingSeparator()
	toolbar.AddFlexibleSpace()
	toolbar.AddButton("Close").
		SetSymbol("xmark.circle").
		SetTooltip("Close this scratch window").
		OnClick(func(*application.Context) {
			if window != nil {
				window.Close()
			}
		})

	window = r.app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  fmt.Sprintf("Scratch %d", n),
		Width:  980,
		Height: 640,
		URL:    "/editor.html",
		Mac: application.MacWindow{
			Backdrop:      application.MacBackdropNormal,
			ContentLayout: application.MacContentLayoutEdgeToEdge,
			TitleBar: application.MacTitleBar{
				FullSizeContent:      true,
				HideToolbarSeparator: true,
				ToolbarStyle:         application.MacToolbarStyleUnified,
			},
			// Configured with the window so it is installed during creation,
			// before the toolbar whose tracking separator follows the sidebar.
			SplitView: split,
			Toolbar:   toolbar,
		},
	})
	return window
}

// NewNativeEditorWindow opens a WebView-free window: a native sidebar of
// snippets beside an NSTextView editor.
func (r *daymarkRuntimeWindows) NewNativeEditorWindow() *application.NativeWindow {
	n := r.count.Add(1)
	editor := application.NewMacTextEditor().
		SetText("A native NSTextView, created from the Window menu after the application started.\n")

	snippets := []struct{ title, body string }{
		{"Morning", "Coffee first. Then the walk, then the page.\n"},
		{"Afternoon", "Reply to the two letters. Water the fig.\n"},
		{"Evening", "Read until the light goes. Sleep early.\n"},
	}
	sidebar := application.NewMacSidebar()
	section := sidebar.AddSection("Snippets")
	for _, snippet := range snippets {
		body := snippet.body
		section.AddItem(snippet.title).SetSymbol("text.quote").OnClick(func(*application.Context) {
			editor.SetText(body)
		})
	}

	split := application.NewMacSplitView()
	split.AddSidebar(sidebar).
		SetMinimumThickness(180).
		SetMaximumThickness(300).
		SetCollapsible(true)
	split.AddTextEditor(editor).SetMinimumThickness(360)

	toolbar := application.NewMacToolbar().SetDisplayMode(application.MacToolbarDisplayModeIconOnly)
	toolbar.AddSidebarToggle()
	toolbar.AddSidebarTrackingSeparator()
	toolbar.AddFlexibleSpace()
	toolbar.AddButton("Clear").
		SetSymbol("trash").
		SetTooltip("Clear the editor").
		OnClick(func(*application.Context) { editor.SetText("") })

	// A NativeWindow is created as soon as it has content. With the layout
	// in the options that is now, and the window is shown before this call
	// returns; NativeWindowManager.New followed by SetSplitView would create
	// it at the SetSplitView call instead.
	return r.app.NativeWindow.NewWithOptions(application.NativeWindowOptions{
		Name:            fmt.Sprintf("native-editor-%d", n),
		Title:           fmt.Sprintf("Native Editor %d", n),
		Width:           900,
		Height:          600,
		MinWidth:        560,
		MinHeight:       360,
		InitialPosition: application.WindowCentered,
		Mac: application.MacWindow{
			Backdrop: application.MacBackdropNormal,
			TitleBar: application.MacTitleBar{
				FullSizeContent:      true,
				HideToolbarSeparator: true,
				ToolbarStyle:         application.MacToolbarStyleUnified,
			},
		},
		SplitView: split,
		Toolbar:   toolbar,
	})
}
