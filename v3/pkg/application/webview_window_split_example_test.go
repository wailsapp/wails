package application_test

import (
	"fmt"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func ExampleNewMacSplitView() {
	// The minimal layout: a native sidebar next to the window's own WebView.
	split := application.NewMacSplitView().
		SetAutosaveName("example.main-window")

	sourceList := application.NewMacSidebar()
	notes := sourceList.AddSection("Notes")
	first := notes.AddItem("First note").SetSymbol("doc.text")
	first.OnClick(func(*application.Context) {
		// Update the primary WebView through application state or an event.
	})
	sourceList.SetSelectedItem(first)

	sidePane := split.AddSidebar(sourceList)
	sidePane.
		SetMinimumThickness(210).
		SetMaximumThickness(340).
		SetCollapsible(true)

	split.AddPrimaryContent().
		SetContentLayout(application.MacContentLayoutEdgeToEdge)

	// Attach the layout before the window is shown:
	// window.SetSplitView(split)
}

func ExampleMacToolbar_AddSidebarTrackingSeparator() {
	// AppKit keeps every toolbar item before the tracking separator aligned
	// above the native sidebar as its divider moves. The separator requires
	// the same window to have a sidebar split layout.
	toolbar := application.NewMacToolbar()
	toolbar.AddSidebarToggle()
	toolbar.AddButton("New").OnClick(func(*application.Context) {})
	toolbar.AddSidebarTrackingSeparator()
	toolbar.AddButton("Save").OnClick(func(*application.Context) {})

	// window.SetSplitView(split) // installed before the toolbar
	// window.SetToolbar(toolbar)
}

func ExampleMacSplitPane_OnCollapsedChange() {
	split := application.NewMacSplitView()
	sidePane := split.AddSidebar(application.NewMacSidebar()).SetCollapsible(true)
	split.AddPrimaryContent()

	// The callback observes every collapse source: the toolbar toggle, the
	// View menu, a divider gesture, window resizing, and SetCollapsed.
	sidePane.OnCollapsedChange(func(_ *application.Context, collapsed bool) {
		fmt.Printf("sidebar collapsed: %v\n", collapsed)
	})

	// Programmatic collapse animates and reports through the same callback:
	sidePane.SetCollapsed(true)
	// A menu item can flip the state without tracking it:
	sidePane.Toggle()
}

func ExampleMacSplitView_AddInspector() {
	split := application.NewMacSplitView()
	split.AddSidebar(application.NewMacSidebar())
	split.AddPrimaryContent()

	inspector := application.NewMacInspector()
	document := inspector.AddSection("Document")
	title := document.AddTextField("Title", "A note")
	document.AddPopup("Category", []string{"Personal", "Work"}, 0)
	document.AddCheckbox("Pinned", false)
	words := inspector.AddSection("Statistics").AddLabel("Words", "12")

	title.OnTextChange(func(_ *application.Context, value string) {
		// Update application state. The handle already contains value here.
	})
	inspectPane := split.AddInspector(inspector)
	inspectPane.SetMinimumThickness(240).SetMaximumThickness(360).SetCollapsible(true)

	// Programmatic state changes do not invoke user callbacks.
	words.SetValue("13")
}

func ExampleMacToolbar_AddInspectorToggle() {
	toolbar := application.NewMacToolbar()
	toolbar.AddFlexibleSpace()
	toolbar.AddInspectorTrackingSeparator()
	toolbar.AddInspectorToggle()

	// The window must contain a native inspector split pane before attaching
	// toolbar. On macOS 14+ AppKit owns both standard toolbar identifiers;
	// older releases receive a native toggle fallback.
}
