//go:build darwin

package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed notes/*.txt
var sampleNotes embed.FS

type nativeFile struct {
	name string
	path string
	item *application.MacSidebarItem
}

type nativeEditorApp struct {
	window      *application.NativeWindow
	editor      *application.MacTextEditor
	sidebar     *application.MacSidebar
	sidebarPane *application.MacSplitPane
	saveItem    *application.MacToolbarItem
	directory   string

	lock      sync.Mutex
	files     []nativeFile
	active    string
	dirty     bool
	switching bool
}

func newNativeEditorApp(app *application.App, benchmark benchmarkConfig) (*nativeEditorApp, error) {
	directory, err := materialiseSampleNotes(benchmark.documentBytes)
	if err != nil {
		return nil, err
	}

	result := &nativeEditorApp{
		directory: directory,
		editor:    application.NewMacTextEditor(),
		sidebar:   application.NewMacSidebar(),
	}
	if err := result.loadSidebar(); err != nil {
		return nil, err
	}

	split := application.NewMacSplitView().SetAutosaveName("wails.native-notes")
	result.sidebarPane = split.AddSidebar(result.sidebar)
	result.sidebarPane.
		SetMinimumThickness(210).
		SetMaximumThickness(340).
		SetCollapsible(true)
	contentPane := split.AddTextEditor(result.editor)
	contentPane.SetMinimumThickness(420)

	result.window = app.NativeWindow.NewWithOptions(application.NativeWindowOptions{
		Name:      "native-notes",
		Title:     "Native Notes",
		Width:     980,
		Height:    680,
		MinWidth:  640,
		MinHeight: 420,
		// Open once on launch so the native hierarchy is immediately visible.
		// Closing the window hides it; the status item opens it again.
		Hidden:          benchmark.hidden,
		HideOnClose:     true,
		InitialPosition: application.WindowCentered,
		Mac: application.MacWindow{
			Backdrop: application.MacBackdropNormal,
			TitleBar: application.MacTitleBar{
				FullSizeContent:      true,
				HideToolbarSeparator: true,
				ToolbarStyle:         application.MacToolbarStyleUnified,
			},
		},
	})
	if err := result.window.SetSplitView(split); err != nil {
		return nil, err
	}

	toolbar := application.NewMacToolbar().SetDisplayMode(application.MacToolbarDisplayModeIconOnly)
	toolbar.AddSidebarToggle()
	toolbar.AddSidebarTrackingSeparator()
	search := toolbar.AddSearch("Search Files").SetTooltip("Filter text files")
	search.OnSearch(func(_ *application.Context, query string) { result.filter(query) })
	toolbar.AddFlexibleSpace()
	result.saveItem = toolbar.AddButton("Save").
		SetSymbol("square.and.arrow.down").
		SetTooltip("Save the current text file").
		SetBordered(true).
		SetEnabled(false)
	result.saveItem.OnClick(func(*application.Context) {
		if err := result.save(); err != nil {
			result.window.Error("save: %s", err)
		}
	})
	if err := result.window.SetToolbar(toolbar); err != nil {
		return nil, err
	}

	result.editor.OnChange(func(*application.Context) {
		result.lock.Lock()
		if result.switching {
			result.lock.Unlock()
			return
		}
		result.dirty = true
		result.lock.Unlock()
		result.saveItem.SetEnabled(true).SetBadgeCount(1)
	})

	if len(result.files) > 0 {
		result.sidebar.SetSelectedItem(result.files[0].item)
		if err := result.open(result.files[0].path); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func materialiseSampleNotes(benchmarkDocumentBytes int) (string, error) {
	directory, err := os.MkdirTemp("", "wails-native-notes-")
	if err != nil {
		return "", err
	}
	entries, err := fs.ReadDir(sampleNotes, "notes")
	if err != nil {
		return "", err
	}
	for _, entry := range entries {
		data, readErr := sampleNotes.ReadFile("notes/" + entry.Name())
		if readErr != nil {
			return "", readErr
		}
		if writeErr := os.WriteFile(filepath.Join(directory, entry.Name()), data, 0o644); writeErr != nil {
			return "", writeErr
		}
	}
	if benchmarkDocumentBytes > 0 {
		line := []byte("A native editor benchmark line with predictable UTF-8 text.\n")
		data := make([]byte, benchmarkDocumentBytes)
		for offset := 0; offset < len(data); offset += len(line) {
			copy(data[offset:], line)
		}
		if err := os.WriteFile(filepath.Join(directory, "00 Benchmark.txt"), data, 0o644); err != nil {
			return "", err
		}
	}
	return directory, nil
}

func (a *nativeEditorApp) loadSidebar() error {
	entries, err := os.ReadDir(a.directory)
	if err != nil {
		return err
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	section := a.sidebar.AddSection("Text Files")
	for _, entry := range entries {
		if entry.IsDir() || strings.ToLower(filepath.Ext(entry.Name())) != ".txt" {
			continue
		}
		path := filepath.Join(a.directory, entry.Name())
		item := section.AddItem(strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))).
			SetSymbol("doc.plaintext").
			SetTooltip(path)
		item.OnClick(func(*application.Context) {
			if err := a.open(path); err != nil {
				a.window.Error("open %s: %s", filepath.Base(path), err)
			}
		})
		a.files = append(a.files, nativeFile{name: entry.Name(), path: path, item: item})
	}
	return nil
}

func (a *nativeEditorApp) open(path string) error {
	a.lock.Lock()
	dirty := a.dirty
	active := a.active
	a.lock.Unlock()
	if dirty && active != "" && active != path {
		if err := a.save(); err != nil {
			return err
		}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	a.lock.Lock()
	a.switching = true
	a.active = path
	a.dirty = false
	a.lock.Unlock()
	a.editor.SetText(string(data))
	a.lock.Lock()
	a.switching = false
	a.lock.Unlock()
	a.saveItem.SetEnabled(false).SetBadgeCount(0)
	a.window.SetTitle(fmt.Sprintf("%s — Native Notes", filepath.Base(path)))
	a.editor.Focus()
	return nil
}

func (a *nativeEditorApp) save() error {
	a.lock.Lock()
	path := a.active
	a.lock.Unlock()
	if path == "" {
		return nil
	}
	if err := os.WriteFile(path, []byte(a.editor.Text()), 0o644); err != nil {
		return err
	}
	a.lock.Lock()
	a.dirty = false
	a.lock.Unlock()
	a.saveItem.SetEnabled(false).SetBadgeCount(0)
	return nil
}

func (a *nativeEditorApp) filter(query string) {
	query = strings.ToLower(strings.TrimSpace(query))
	for _, file := range a.files {
		file.item.SetHidden(query != "" && !strings.Contains(strings.ToLower(file.name), query))
	}
}

func (a *nativeEditorApp) show() {
	a.window.Show()
	a.editor.Focus()
}
