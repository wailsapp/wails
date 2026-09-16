//go:build darwin

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// quickNote is a small floating capture panel. It exists to show the
// NativeWindowOptions.Mac options that the main document window has no use
// for: a translucent backdrop, an NSPanel at an explicit floating window
// level, a hidden titlebar with a custom corner radius, and a collection
// behaviour that keeps the panel available on every Space.
type quickNote struct {
	window *application.NativeWindow
	editor *application.MacTextEditor
	count  *application.MacInspectorControl
	main   *nativeEditorApp
	path   string
}

func newQuickNote(app *application.App, main *nativeEditorApp) (*quickNote, error) {
	note := &quickNote{
		editor: application.NewMacTextEditor(),
		main:   main,
		path:   filepath.Join(main.directory, "Quick Note.txt"),
	}
	if data, err := os.ReadFile(note.path); err == nil {
		note.editor.SetText(string(data))
	}

	inspector := application.NewMacInspector()
	section := inspector.AddSection("Quick Note")
	note.count = section.AddLabel("Characters", note.characterCount())
	section.AddButton("Append to Note").OnClick(func(*application.Context) { note.append() })
	section.AddButton("Save").OnClick(func(*application.Context) {
		if err := note.save(); err != nil {
			note.window.Error("save: %s", err)
		}
	})
	section.AddButton("Hide").OnClick(func(*application.Context) { note.window.Hide() })

	split := application.NewMacSplitView().SetAutosaveName("wails.native-notes.quick")
	split.AddTextEditor(note.editor).SetMinimumThickness(240)
	split.AddInspector(inspector).SetMinimumThickness(180).SetMaximumThickness(220)

	note.window = app.NativeWindow.NewWithOptions(application.NativeWindowOptions{
		Name:            "quick-note",
		Title:           "Quick Note",
		Width:           520,
		Height:          300,
		MinWidth:        420,
		MinHeight:       200,
		Hidden:          true,
		HideOnClose:     true,
		InitialPosition: application.WindowCentered,
		Mac: application.MacWindow{
			// A floating NSPanel that stays above the document window and
			// follows the user across Spaces and fullscreen applications.
			WindowClass:        application.MacWindowClassPanel,
			PanelPreferences:   application.MacPanelPreferences{FloatingPanel: true},
			WindowLevel:        application.MacWindowLevelFloating,
			CollectionBehavior: application.MacWindowCollectionBehaviorCanJoinAllSpaces | application.MacWindowCollectionBehaviorFullScreenAuxiliary,
			// Frosted panel with a custom rounded shape. TitleBar.Hide is the
			// native frameless switch; CornerRadius then applies.
			Backdrop:     application.MacBackdropTranslucent,
			TitleBar:     application.MacTitleBar{Hide: true},
			CornerRadius: 14,
		},
	})
	if err := note.window.SetSplitView(split); err != nil {
		return nil, err
	}
	note.editor.OnChange(func(*application.Context) {
		note.count.SetValue(note.characterCount())
	})
	return note, nil
}

func (n *quickNote) characterCount() string {
	return fmt.Sprintf("%d", utf8.RuneCountInString(n.editor.Text()))
}

// toggle shows the panel with the editor focused, or hides it when visible.
func (n *quickNote) toggle() {
	if n.window.IsVisible() {
		n.window.Hide()
		return
	}
	n.window.Show()
	n.editor.Focus()
}

func (n *quickNote) save() error {
	return os.WriteFile(n.path, []byte(n.editor.Text()), 0o644)
}

// append moves the captured text into the document open in the main window.
func (n *quickNote) append() {
	text := strings.TrimSpace(n.editor.Text())
	if text == "" {
		return
	}
	n.main.appendText(text)
	n.editor.SetText("")
	n.count.SetValue(n.characterCount())
	if err := n.save(); err != nil {
		n.window.Error("save: %s", err)
	}
	n.main.show()
}
