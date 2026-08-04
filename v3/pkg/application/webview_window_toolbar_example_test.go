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
