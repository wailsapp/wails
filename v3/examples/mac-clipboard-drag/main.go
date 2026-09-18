package main

import (
	"bytes"
	"embed"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed assets
var assets embed.FS

// DragService is bound to the page. The page calls StartDrag from a
// mousedown handler on the draggable card, which is the only moment the OS
// lets an application begin a drag session.
type DragService struct {
	app    *application.App
	window *application.WebviewWindow
}

// StartDrag begins a system drag of the given kind: "text", "file" or
// "promise". It must be called while the mouse button is down.
func (s *DragService) StartDrag(kind string) error {
	var items application.DragItems
	switch kind {
	case "text":
		items.Text = "Dragged out of a Wails window at " + time.Now().Format(time.Kitchen)
	case "file":
		path := filepath.Join(os.TempDir(), "wails-drag-out.txt")
		if err := os.WriteFile(path, []byte("This file was dragged out of a Wails window.\n"), 0o644); err != nil {
			return err
		}
		items.Files = []string{path}
	case "promise":
		items.Promises = []application.DragPromise{{
			Filename: "wails-swatch.png",
			Data: func() ([]byte, error) {
				// Produced only when the destination accepts the drop.
				return swatchPNG(48, color.RGBA{R: 0xE5, G: 0x3E, B: 0x3E, A: 0xFF}), nil
			},
		}}
		items.Image = swatchPNG(48, color.RGBA{R: 0xE5, G: 0x3E, B: 0x3E, A: 0xFF})
		items.ImageOffset = application.Point{X: 24, Y: 24}
	default:
		return fmt.Errorf("unknown drag kind %q", kind)
	}
	err := s.window.StartDrag(items)
	if err != nil {
		s.app.Event.Emit("log", "StartDrag: "+err.Error())
	}
	return err
}

// ClipboardService exposes the rich clipboard to the page.
type ClipboardService struct {
	app *application.App
}

func (s *ClipboardService) CopyText() error {
	if !s.app.Clipboard.SetText("Plain text from Wails at " + time.Now().Format(time.Kitchen)) {
		return fmt.Errorf("SetText failed")
	}
	return nil
}

func (s *ClipboardService) CopyImage() error {
	return s.app.Clipboard.SetImage(swatchPNG(64, color.RGBA{R: 0x2B, G: 0x8A, B: 0xF7, A: 0xFF}))
}

func (s *ClipboardService) CopyHTML() error {
	return s.app.Clipboard.SetHTML("<p>Rich <b>HTML</b> from <i>Wails</i></p>", "Rich HTML from Wails")
}

func (s *ClipboardService) CopyFile() error {
	path := filepath.Join(os.TempDir(), "wails-clipboard-file.txt")
	if err := os.WriteFile(path, []byte("Copied as a file reference by Wails.\n"), 0o644); err != nil {
		return err
	}
	return s.app.Clipboard.SetFiles([]string{path})
}

func (s *ClipboardService) CopyCustom() error {
	return s.app.Clipboard.SetData("com.wails.example.record", []byte(`{"id":42,"name":"custom record"}`))
}

func (s *ClipboardService) Clear() {
	s.app.Clipboard.Clear()
}

// Inspect describes what is on the clipboard right now.
func (s *ClipboardService) Inspect() map[string]any {
	result := map[string]any{
		"changeCount": s.app.Clipboard.ChangeCount(),
		"types":       s.app.Clipboard.Types(),
	}
	if text, ok := s.app.Clipboard.Text(); ok && text != "" {
		result["text"] = text
	}
	if html, err := s.app.Clipboard.HTML(); err == nil {
		result["html"] = html
	}
	if files, err := s.app.Clipboard.Files(); err == nil && len(files) > 0 {
		result["files"] = files
	}
	if data, err := s.app.Clipboard.Data("com.wails.example.record"); err == nil {
		result["custom"] = string(data)
	}
	if img, err := s.app.Clipboard.Image(); err == nil {
		result["image"] = "data:image/png;base64," + base64.StdEncoding.EncodeToString(img)
	}
	return result
}

func main() {
	app := application.New(application.Options{
		Name:        "Clipboard and Drag",
		Description: "Rich clipboard, drag out and non-file drops on macOS",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	dragService := &DragService{app: app}
	app.RegisterService(application.NewService(dragService))
	app.RegisterService(application.NewService(&ClipboardService{app: app}))

	menu := app.NewMenu()
	menu.AddRole(application.AppMenu)
	menu.AddRole(application.EditMenu)
	app.Menu.Set(menu)

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Clipboard and Drag",
		Width:  900,
		Height: 640,
		// DropFiles keeps the normal file drop pipeline; the other types
		// are delivered to OnDrop below instead of the page.
		DropTypes: []application.DropType{
			application.DropFiles,
			application.DropText,
			application.DropURLs,
			application.DropImages,
		},
		Mac: application.MacWindow{
			TitleBar:                application.MacTitleBarHiddenInsetUnified,
			InvisibleTitleBarHeight: 44,
		},
	})
	dragService.window = window

	// Files still arrive through the file drop event.
	window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
		files := event.Context().DroppedFiles()
		app.Event.Emit("dropped", map[string]any{"files": files})
	})

	// Text, URLs and images from other apps arrive here.
	window.OnDrop(func(ctx *application.Context, data application.DropData) {
		payload := map[string]any{
			"text": data.Text,
			"urls": data.URLs,
			"x":    data.X,
			"y":    data.Y,
		}
		if len(data.Images) > 0 {
			payload["image"] = "data:image/png;base64," + base64.StdEncoding.EncodeToString(data.Images[0])
		}
		app.Event.Emit("dropped", payload)
	})

	window.OnDragEnd(func(operation application.DragOperation) {
		app.Event.Emit("log", "drag ended: "+operation.String())
	})

	// Clipboard changes made anywhere, including other apps, are polled
	// every 500 ms while a listener is registered.
	app.Clipboard.OnChange(func() {
		app.Event.Emit("clipboard-changed", app.Clipboard.ChangeCount())
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

// swatchPNG renders a rounded square of the given colour with a white
// initial, used both as clipboard content and as a drag image.
func swatchPNG(size int, fill color.RGBA) []byte {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	radius := size / 6
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if insideRoundedRect(x, y, size, radius) {
				img.SetRGBA(x, y, fill)
			}
		}
	}
	// A simple "W" drawn as three diagonal strokes.
	white := color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	for i := 0; i < size/2; i++ {
		for t := 0; t < 3; t++ {
			img.SetRGBA(size/6+i/2+t, size/4+i, white)
			img.SetRGBA(size/2-i/2+t, size/4+i, white)
			img.SetRGBA(size/2+i/2+t, size/4+i, white)
			img.SetRGBA(size*5/6-i/2+t, size/4+i, white)
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil
	}
	return buf.Bytes()
}

func insideRoundedRect(x, y, size, radius int) bool {
	dx, dy := 0, 0
	if x < radius {
		dx = radius - x
	} else if x >= size-radius {
		dx = x - (size - radius - 1)
	}
	if y < radius {
		dy = radius - y
	} else if y >= size-radius {
		dy = y - (size - radius - 1)
	}
	return dx*dx+dy*dy <= radius*radius
}
