package application_test

import (
	"fmt"
	"html"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func ExampleMacShareProviderFunc() {
	toolbar := application.NewMacToolbar()

	share := toolbar.AddShare("Share").SetProvider(application.MacShareProviderFunc{
		Available: []application.MacShareRepresentation{
			{ContentType: application.MacShareTypePlainText},
		},
		Load: func(request application.MacShareRequest) ([]byte, error) {
			return []byte("Hello from Wails"), nil
		},
	})
	share.SetSubject("A Wails note").SetSuggestedName("Wails Note")
	share.OnShared(func(_ *application.Context, service string) {
		fmt.Printf("shared with %s\n", service)
	})
	share.OnShareError(func(_ *application.Context, service string, err error) {
		fmt.Printf("%s failed: %v\n", service, err)
	})

	// Attach the completed toolbar to a WebviewWindow:
	// window.SetToolbar(toolbar)
}

func ExampleMacToolbar_AddMenu() {
	toolbar := application.NewMacToolbar()

	// A dropdown toolbar item reuses the ordinary Menu model. Item clicks fire
	// MenuItem.OnClick exactly as they do from the application menu.
	actions := application.NewMenu()
	actions.Add("Duplicate Note").OnClick(func(*application.Context) {
		fmt.Println("duplicate")
	})
	actions.AddSeparator()
	actions.Add("Move to Trash").OnClick(func(*application.Context) {
		fmt.Println("trash")
	})

	toolbar.AddMenu("Actions", actions).
		SetSymbol("ellipsis.circle").
		SetShowsIndicator(true)

	// Items can be added, moved and removed after the toolbar is attached;
	// the native toolbar follows the Go model.
	// window.SetToolbar(toolbar)
	// later := toolbar.AddButton("Later").OnClick(func(*application.Context) {})
	// toolbar.Move(later, 0)
	// toolbar.Remove(later)
}

func ExampleMacToolbar_SetCustomizable() {
	toolbar := application.NewMacToolbar()

	// Stable persistence keys let AppKit restore a user-customised layout on
	// the next launch. Items without a key keep generated identifiers and are
	// not restored.
	back := toolbar.AddButton("Back").
		SetSymbol("chevron.backward").
		SetPersistenceKey("back").
		SetNavigational(true).
		SetVisibilityPriority(application.MacToolbarVisibilityPriorityLow)
	back.OnClick(func(*application.Context) {})

	toolbar.AddFlexibleSpace()

	search := toolbar.AddSearch("Search").
		SetPersistenceKey("search").
		SetSearchPlaceholder("Search notes").
		SetSearchRecentsKey("example.search.recents")
	search.OnSearch(func(_ *application.Context, query string) {
		fmt.Println("search:", query)
	})

	// Offered in the customisation palette but hidden until the user adds it.
	toolbar.AddButton("Statistics").
		SetSymbol("chart.bar").
		SetPersistenceKey("statistics").
		SetInDefaultSet(false).
		OnClick(func(*application.Context) {})

	// Enables the standard "Customize Toolbar..." sheet and autosave. Call it
	// before attaching the toolbar so the persistence key becomes the
	// NSToolbar identifier.
	toolbar.SetCustomizable("example.main-window")

	// window.SetToolbar(toolbar)
	// toolbar.RunCustomizationPalette() opens the sheet programmatically.
}

type exampleNote struct {
	Title string
	Body  string
}

type exampleNoteShareProvider struct {
	lock sync.RWMutex
	note exampleNote
}

func (p *exampleNoteShareProvider) ShareRepresentations() []application.MacShareRepresentation {
	return []application.MacShareRepresentation{
		{ContentType: application.MacShareTypeHTML},
		{ContentType: application.MacShareTypePlainText},
	}
}

func (p *exampleNoteShareProvider) ShareData(request application.MacShareRequest) ([]byte, error) {
	p.lock.RLock()
	note := p.note
	p.lock.RUnlock()

	switch request.ContentType {
	case application.MacShareTypePlainText:
		return []byte(note.Title + "\n\n" + note.Body), nil
	case application.MacShareTypeHTML:
		return []byte("<article><h1>" + html.EscapeString(note.Title) +
			"</h1><p>" + strings.ReplaceAll(html.EscapeString(note.Body), "\n", "<br>") +
			"</p></article>"), nil
	default:
		return nil, fmt.Errorf("unsupported share type %q", request.ContentType)
	}
}

func (p *exampleNoteShareProvider) Update(note exampleNote) {
	p.lock.Lock()
	p.note = note
	p.lock.Unlock()
}

func ExampleMacShareProvider_stateful() {
	provider := &exampleNoteShareProvider{}
	provider.Update(exampleNote{
		Title: "Saturday, slowly.",
		Body:  "A good day has room around it.",
	})

	toolbar := application.NewMacToolbar()
	toolbar.AddShare("Share").
		SetProvider(provider).
		SetSubject("Saturday, slowly.").
		SetSuggestedName("Daymark Note")

	// The provider can be updated without rebuilding the toolbar. The next
	// native share request reads the latest thread-safe snapshot.
	provider.Update(exampleNote{Title: "A new note", Body: "Updated content"})
}
