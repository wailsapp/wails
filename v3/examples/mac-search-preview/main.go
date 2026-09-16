package main

import (
	"embed"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

// spotlightDomain groups the example's items so they can be removed together.
const spotlightDomain = "com.wails.mac-search-preview.notes"

// Note is one of the sample documents. Path points at the file written to
// the sample directory at startup so Quick Look and Spotlight have something
// real to show.
type Note struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Body     string   `json:"body"`
	Keywords []string `json:"keywords"`
	Path     string   `json:"path"`
}

var sampleNotes = []Note{
	{ID: "alpha", Title: "Alpha release checklist", Keywords: []string{"release", "checklist"},
		Body: "Alpha release checklist\n\n- Tag the build\n- Run the smoke tests\n- Notify the beta group\n"},
	{ID: "beta", Title: "Beta feedback summary", Keywords: []string{"feedback", "beta"},
		Body: "Beta feedback summary\n\nMost testers asked for keyboard shortcuts and a dark theme.\n"},
	{ID: "gamma", Title: "Gamma launch plan", Keywords: []string{"launch", "plan"},
		Body: "Gamma launch plan\n\nAnnounce on Monday, publish the changelog, then open the forum thread.\n"},
}

// SearchPreviewService is bound to the frontend and wraps the Spotlight,
// QuickLook and AppleEvents managers.
type SearchPreviewService struct {
	app *application.App

	dir       string
	notes     []Note
	imagePath string

	lock    sync.Mutex
	indexed bool
}

// Notes lists the sample notes.
func (s *SearchPreviewService) Notes() []Note {
	return s.notes
}

// Files lists every sample path: the notes plus the generated image.
func (s *SearchPreviewService) Files() []string {
	paths := make([]string, 0, len(s.notes)+1)
	for _, note := range s.notes {
		paths = append(paths, note.Path)
	}
	return append(paths, s.imagePath)
}

// SpotlightAvailable reports whether the system index accepts items.
func (s *SearchPreviewService) SpotlightAvailable() bool {
	return s.app.Spotlight.IsAvailable()
}

// IndexNotes adds the sample notes to Spotlight.
func (s *SearchPreviewService) IndexNotes() error {
	items := make([]application.SearchableItem, 0, len(s.notes))
	for _, note := range s.notes {
		items = append(items, application.SearchableItem{
			ID:          note.ID,
			Domain:      spotlightDomain,
			Title:       note.Title,
			Description: strings.SplitN(strings.TrimSpace(strings.TrimPrefix(note.Body, note.Title)), "\n", 2)[0],
			Keywords:    note.Keywords,
			ContentType: "public.plain-text",
			URL:         "mac-search-preview://notes/" + note.ID,
		})
	}
	if err := s.app.Spotlight.Index(items); err != nil {
		return err
	}
	s.lock.Lock()
	s.indexed = true
	s.lock.Unlock()
	return nil
}

// RemoveIndex deletes the example's Spotlight domain.
func (s *SearchPreviewService) RemoveIndex() error {
	if err := s.app.Spotlight.DeleteDomain(spotlightDomain); err != nil {
		return err
	}
	s.lock.Lock()
	s.indexed = false
	s.lock.Unlock()
	return nil
}

// Preview opens Quick Look on one note (by ID) or on every sample file.
func (s *SearchPreviewService) Preview(id string) error {
	if id == "" {
		return s.app.QuickLook.Preview(s.Files())
	}
	note, ok := s.note(id)
	if !ok {
		return fmt.Errorf("no note %q", id)
	}
	return s.app.QuickLook.Preview([]string{note.Path})
}

// ClosePreview closes the Quick Look panel.
func (s *SearchPreviewService) ClosePreview() {
	s.app.QuickLook.ClosePreview()
}

// IsPreviewOpen reports whether the Quick Look panel is visible.
func (s *SearchPreviewService) IsPreviewOpen() bool {
	return s.app.QuickLook.IsPreviewOpen()
}

// Thumbnail renders a Quick Look thumbnail for a sample path and returns it
// as a data URL for an <img>.
func (s *SearchPreviewService) Thumbnail(path string, size int, iconMode bool) (string, error) {
	data, err := s.app.QuickLook.Thumbnail(path, application.ThumbnailOptions{Width: size, Height: size, Scale: 2, IconMode: iconMode})
	if err != nil {
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(data), nil
}

// ScriptingDefinition returns the sdef for the registered Apple Event handlers.
func (s *SearchPreviewService) ScriptingDefinition() string {
	return s.app.AppleEvents.ScriptingDefinition()
}

// SendNote sends the WAIL/note event to the application with the given
// bundle identifier (this app when bundled) and returns the reply.
func (s *SearchPreviewService) SendNote(bundleID, text string) (string, error) {
	reply, err := s.app.AppleEvents.Send(bundleID, "WAIL", "note", text)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(reply), nil
}

func (s *SearchPreviewService) note(id string) (Note, bool) {
	id = strings.ToLower(strings.TrimSpace(id))
	for _, note := range s.notes {
		if note.ID == id || strings.ToLower(note.Title) == id {
			return note, true
		}
	}
	return Note{}, false
}

// open tells the page which note was opened and how, and brings the window
// forward.
func (s *SearchPreviewService) open(id, source string) {
	if note, ok := s.note(id); ok {
		s.app.Event.Emit("note:opened", map[string]string{"id": note.ID, "title": note.Title, "source": source})
	} else {
		s.app.Event.Emit("note:opened", map[string]string{"id": id, "title": "", "source": source})
	}
	if window := s.app.Window.Current(); window != nil {
		window.Focus()
	}
}

// writeSamples materialises the notes and a generated PNG in a temporary
// directory so Quick Look and the thumbnail generator have files to read.
func writeSamples() (string, []Note, string, error) {
	dir, err := os.MkdirTemp("", "mac-search-preview-")
	if err != nil {
		return "", nil, "", err
	}
	notes := make([]Note, len(sampleNotes))
	for i, note := range sampleNotes {
		note.Path = filepath.Join(dir, note.ID+".txt")
		if err := os.WriteFile(note.Path, []byte(note.Body), 0o644); err != nil {
			return "", nil, "", err
		}
		notes[i] = note
	}

	imagePath := filepath.Join(dir, "cover.png")
	img := image.NewRGBA(image.Rect(0, 0, 320, 200))
	for y := 0; y < 200; y++ {
		for x := 0; x < 320; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x * 255 / 320), G: uint8(y * 255 / 200), B: 180, A: 255})
		}
	}
	file, err := os.Create(imagePath)
	if err != nil {
		return "", nil, "", err
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		return "", nil, "", err
	}
	return dir, notes, imagePath, nil
}

func main() {
	dir, notes, imagePath, err := writeSamples()
	if err != nil {
		log.Fatal(err)
	}
	service := &SearchPreviewService{dir: dir, notes: notes, imagePath: imagePath}

	app := application.New(application.Options{
		Name:        "mac-search-preview",
		Description: "Spotlight indexing, Quick Look previews and thumbnails, Apple Events",
		Services: []application.Service{
			application.NewService(service),
		},
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		// Leave the system index and the temporary directory clean on exit.
		OnShutdown: func() {
			service.lock.Lock()
			indexed := service.indexed
			service.lock.Unlock()
			if indexed {
				if err := service.app.Spotlight.DeleteDomain(spotlightDomain); err != nil {
					log.Println("spotlight cleanup:", err)
				}
			}
			os.RemoveAll(dir)
		},
	})
	service.app = app

	// Apple Events: "WAIL"/"note" opens a note by ID or title and replies with
	// its text. From a terminal (bundled app):
	//   osascript -e 'tell application id "com.wails.mac-search-preview" to «event WAILnote» "alpha"'
	err = app.AppleEvents.Handle("WAIL", "note", func(ctx *application.Context, event application.AppleEvent) (application.AppleEventReply, error) {
		id, _ := event.DirectObject.(string)
		note, ok := service.note(id)
		if !ok {
			return application.AppleEventReply{ErrorNumber: -1728, ErrorString: fmt.Sprintf("no note %q", id)}, nil
		}
		service.open(note.ID, "apple event")
		return application.AppleEventReply{Result: note.Body}, nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Spotlight: selecting an indexed note (or "Search in App") lands here.
	app.Spotlight.OnOpen(func(ctx *application.Context, id string, query string) {
		if id != "" {
			service.open(id, "spotlight")
			return
		}
		app.Event.Emit("note:search", query)
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Search and Preview",
		Width:  900,
		Height: 720,
		URL:    "/",
	})

	if err := app.Run(); err != nil {
		log.Fatal(fmt.Errorf("run: %w", err))
	}
}
