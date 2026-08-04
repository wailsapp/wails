package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"runtime"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

type sharePayload struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Body     string `json:"body"`
}

type daymarkShareProvider struct {
	lock sync.RWMutex
	note sharePayload
}

func (p *daymarkShareProvider) ShareRepresentations() []application.MacShareRepresentation {
	return []application.MacShareRepresentation{
		{ContentType: application.MacShareTypePDF},
		{ContentType: application.MacShareTypeHTML},
		{ContentType: application.MacShareTypePlainText},
	}
}

func (p *daymarkShareProvider) ShareData(request application.MacShareRequest) ([]byte, error) {
	p.lock.RLock()
	note := p.note
	p.lock.RUnlock()

	title := strings.TrimSpace(note.Title)
	if title == "" {
		title = "Untitled note"
	}
	subtitle := strings.TrimSpace(note.Subtitle)
	body := strings.TrimSpace(note.Body)
	log.Printf("native share requested %s", request.ContentType)

	switch request.ContentType {
	case application.MacShareTypePlainText:
		parts := []string{title}
		if subtitle != "" {
			parts = append(parts, subtitle)
		}
		if body != "" {
			parts = append(parts, body)
		}
		return []byte(strings.Join(parts, "\n\n")), nil
	case application.MacShareTypeHTML:
		return []byte(fmt.Sprintf(`<!doctype html>
<html lang="en"><head><meta charset="utf-8"><title>%s</title></head>
<body><article><h1>%s</h1><p><em>%s</em></p><p>%s</p></article></body></html>`,
			html.EscapeString(title), html.EscapeString(title), html.EscapeString(subtitle),
			strings.ReplaceAll(html.EscapeString(body), "\n", "<br>\n"))), nil
	case application.MacShareTypePDF:
		return renderDaymarkPDF(note)
	default:
		return nil, fmt.Errorf("Daymark cannot provide %q", request.ContentType)
	}
}

func (p *daymarkShareProvider) update(note sharePayload) {
	p.lock.Lock()
	p.note = note
	p.lock.Unlock()
}

func main() {
	app := application.New(application.Options{
		Name:        "Daymark",
		Description: "A native macOS toolbar editor",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Daymark",
		Width:  1180,
		Height: 760,
		URL:    "/",
		Mac: application.MacWindow{
			TitleBar: application.MacTitleBar{
				ToolbarStyle: application.MacToolbarStyleExpanded,
			},
		},
	})

	toolbar := application.NewMacToolbar()
	newNote := toolbar.AddButton("New").SetSymbol("square.and.pencil").SetBordered(true)
	mode := toolbar.AddGroup("Mode", application.ToolbarGroupSelectOne)
	mode.AddButton("Write").SetSymbol("pencil").OnClick(func(*application.Context) {
		mode.SetSelectedIndex(0)
		app.Event.Emit("toolbar:mode", "write")
	})
	mode.AddButton("Preview").SetSymbol("doc.richtext").OnClick(func(*application.Context) {
		mode.SetSelectedIndex(1)
		app.Event.Emit("toolbar:mode", "preview")
	})
	toolbar.AddFlexibleSpace()
	search := toolbar.AddSearch("Search notes").SetTooltip("Search your notes")
	share := toolbar.AddShare("Share")
	share.SetTooltip("Share the current note with another app")
	shareProvider := &daymarkShareProvider{note: sharePayload{
		Title:    "Saturday, slowly.",
		Subtitle: "A good day has room around it.",
		Body:     "Leave the phone at home. Walk until the city sounds different. Buy something warm on the way back.",
	}}
	share.SetProvider(shareProvider).SetSubject("Saturday, slowly.").SetSuggestedName("Daymark Note")
	details := toolbar.AddButton("Details").SetSymbol("info.circle").SetBordered(true)
	save := toolbar.AddButton("Save").SetSymbol("checkmark.circle").SetProminent(true).SetBadgeCount(0)
	focus := toolbar.AddButton("Focus").SetSymbol("arrow.up.left.and.arrow.down.right").SetBordered(true)

	var stateLock sync.Mutex
	focused := false
	dirty := false
	setSaved := func() {
		stateLock.Lock()
		dirty = false
		stateLock.Unlock()
		save.SetBadgeCount(0).SetLabel("Saved")
		app.Event.Emit("toolbar:save")
	}
	toggleFocus := func() {
		stateLock.Lock()
		focused = !focused
		isFocused := focused
		stateLock.Unlock()
		details.SetHidden(isFocused)
		if isFocused {
			focus.SetLabel("Exit Focus").SetSymbol("arrow.down.right.and.arrow.up.left")
		} else {
			focus.SetLabel("Focus").SetSymbol("arrow.up.left.and.arrow.down.right")
		}
		app.Event.Emit("toolbar:focus", isFocused)
	}

	newNote.OnClick(func(*application.Context) { app.Event.Emit("toolbar:new") })
	search.OnSearch(func(_ *application.Context, query string) { app.Event.Emit("toolbar:search", query) })
	share.OnShared(func(_ *application.Context, service string) {
		app.Event.Emit("toolbar:share-complete", service)
	})
	share.OnShareError(func(_ *application.Context, service string, err error) {
		app.Event.Emit("toolbar:share-error", map[string]string{"service": service, "message": err.Error()})
	})
	details.OnClick(func(*application.Context) { app.Event.Emit("toolbar:details") })
	save.OnClick(func(*application.Context) { setSaved() })
	focus.OnClick(func(*application.Context) { toggleFocus() })

	menu := app.NewMenu()
	if runtime.GOOS == "darwin" {
		menu.AddRole(application.AppMenu)
	}
	fileMenu := menu.AddSubmenu("File")
	fileMenu.Add("New Note").SetAccelerator("CmdOrCtrl+n").OnClick(func(*application.Context) {
		app.Event.Emit("toolbar:new")
	})
	fileMenu.Add("Save Note").SetAccelerator("CmdOrCtrl+s").OnClick(func(*application.Context) {
		setSaved()
	})
	fileMenu.AddSeparator()
	fileMenu.AddRole(application.CloseWindow)
	menu.AddRole(application.EditMenu)
	viewMenu := menu.AddSubmenu("View")
	viewMenu.Add("Toggle Focus").SetAccelerator("CmdOrCtrl+Shift+f").OnClick(func(*application.Context) {
		toggleFocus()
	})
	viewMenu.Add("Toggle Inspector").SetAccelerator("CmdOrCtrl+Option+i").OnClick(func(*application.Context) {
		app.Event.Emit("toolbar:details")
	})
	menu.AddRole(application.WindowMenu)
	menu.AddRole(application.HelpMenu)
	app.Menu.Set(menu)

	app.Event.On("editor:dirty", func(e *application.CustomEvent) {
		isDirty, ok := e.Data.(bool)
		if !ok {
			return
		}
		stateLock.Lock()
		if isDirty == dirty {
			stateLock.Unlock()
			return
		}
		dirty = isDirty
		stateLock.Unlock()
		if isDirty {
			save.SetBadgeCount(1).SetLabel("Save")
		} else {
			save.SetBadgeCount(0).SetLabel("Saved")
		}
	})

	shareSubject := "Saturday, slowly."
	var shareSubjectLock sync.Mutex
	app.Event.On("editor:share-content", func(e *application.CustomEvent) {
		raw, ok := e.Data.(string)
		if !ok {
			return
		}
		var payload sharePayload
		if json.Unmarshal([]byte(raw), &payload) != nil {
			return
		}
		title := strings.TrimSpace(payload.Title)
		if title == "" {
			title = "Untitled note"
		}
		shareProvider.update(payload)

		shareSubjectLock.Lock()
		subjectChanged := shareSubject != title
		shareSubject = title
		shareSubjectLock.Unlock()
		if subjectChanged {
			share.SetSubject(title)
		}
	})

	window.SetToolbar(toolbar)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
