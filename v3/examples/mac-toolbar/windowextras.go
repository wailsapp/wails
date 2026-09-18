package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// daymarkWindowExtras wires the document-style window behaviours to the
// editor window: the close button's edited dot follows the editor's dirty
// state, File > Export PDF renders the WebView through WebKit and shows the
// exported file's proxy icon in the titlebar, File > Print prints the note
// with explicit page settings, and Window > New Cascaded Window opens an
// editor window offset from the previous one, as document windows are.
type daymarkWindowExtras struct {
	app      *application.App
	window   *application.WebviewWindow
	cascaded atomic.Int32
}

func newDaymarkWindowExtras(app *application.App, window *application.WebviewWindow) *daymarkWindowExtras {
	result := &daymarkWindowExtras{app: app, window: window}
	window.SetSubtitle("Notes")
	app.Event.On("editor:dirty", func(event *application.CustomEvent) {
		if dirty, ok := event.Data.(bool); ok {
			window.SetDocumentEdited(dirty)
		}
	})
	app.Event.On("toolbar:save", func(*application.CustomEvent) {
		window.SetDocumentEdited(false)
	})
	return result
}

// addFileItems adds Export PDF and Print to the File menu.
func (e *daymarkWindowExtras) addFileItems(fileMenu *application.Menu) {
	fileMenu.Add("Export PDF...").SetAccelerator("CmdOrCtrl+Shift+e").OnClick(func(*application.Context) {
		// ExportPDF waits for WebKit, so it must run off the application
		// thread that delivers menu clicks.
		go e.exportPDF()
	})
	fileMenu.Add("Print...").SetAccelerator("CmdOrCtrl+p").OnClick(func(*application.Context) {
		err := e.window.PrintWithOptions(application.PrintOptions{
			Orientation: application.PrintOrientationPortrait,
			Margins:     application.PrintMargins{Top: 36, Left: 36, Bottom: 36, Right: 36},
		})
		if err != nil {
			e.window.Error("print: %s", err)
		}
	})
}

// addWindowItems adds New Cascaded Window to the Window menu.
func (e *daymarkWindowExtras) addWindowItems(windowMenu *application.Menu) {
	windowMenu.Add("New Cascaded Window").SetAccelerator("CmdOrCtrl+Shift+c").OnClick(func(*application.Context) {
		e.newCascadedWindow()
	})
}

func (e *daymarkWindowExtras) exportPDF() {
	pdf, err := e.window.ExportPDF(application.PDFExportOptions{})
	if err != nil {
		e.window.Error("export PDF: %s", err)
		return
	}
	path, err := e.app.Dialog.SaveFile().
		SetMessage("Export the note as a PDF document").
		SetFilename("Daymark Note.pdf").
		AttachToWindow(e.window).
		PromptForSingleSelection()
	if err != nil || path == "" {
		return
	}
	if filepath.Ext(path) == "" {
		path += ".pdf"
	}
	if err := os.WriteFile(path, pdf, 0o644); err != nil {
		e.window.Error("export PDF: %s", err)
		return
	}
	// The titlebar now shows the exported file's proxy icon; it can be
	// dragged to Finder or command-clicked to reveal the folder.
	e.window.SetRepresentedFile(path)
	e.window.SetSubtitle(filepath.Base(path))
}

// newCascadedWindow opens another editor window. WindowCascade places it
// down and to the right of the previous cascaded window (the main window the
// first time), so a stack of windows fans out instead of piling up.
func (e *daymarkWindowExtras) newCascadedWindow() *application.WebviewWindow {
	n := e.cascaded.Add(1)
	return e.app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:           fmt.Sprintf("Daymark %d", n+1),
		Width:           900,
		Height:          600,
		URL:             "/editor.html",
		InitialPosition: application.WindowCascade,
		Mac: application.MacWindow{
			Backdrop: application.MacBackdropNormal,
			TitleBar: application.MacTitleBar{
				HideToolbarSeparator: true,
			},
		},
	})
}
