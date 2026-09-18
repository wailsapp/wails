package main

import (
	"embed"
	"fmt"
	"log"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/*
var assets embed.FS

// DialogService exposes each native dialog feature to the frontend. Every
// method returns a short description of what the user did so the page can
// display it. Bound methods run on their own goroutine, which is what the
// blocking helpers (Prompt, PickColor, PickFont) require.
type DialogService struct {
	app    *application.App
	window *application.WebviewWindow
}

// Suppression shows a warning with a "Do not warn me again" checkbox and a
// help button. The result reports the button and the checkbox state.
func (s *DialogService) Suppression() string {
	done := make(chan string, 1)
	dialog := s.app.Dialog.Warning().
		SetTitle("Delete 3 items?").
		SetMessage("The items will be moved to the Bin.").
		SetSuppression("Do not warn me again").
		SetHelp(func() {
			log.Println("help button pressed")
			s.app.Event.Emit("dialog:help", "Help requested from the suppression alert")
		}).
		AttachToWindow(s.window)
	dialog.OnSuppression(func(ticked bool) {
		log.Printf("suppression checkbox ticked: %v", ticked)
	})
	dialog.AddButton("Delete").SetAsDefault().OnClick(func() {
		done <- fmt.Sprintf("Delete pressed, suppressed=%v", dialog.Suppressed())
	})
	dialog.AddButton("Cancel").SetAsCancel().OnClick(func() {
		done <- fmt.Sprintf("Cancel pressed, suppressed=%v", dialog.Suppressed())
	})
	dialog.Show()
	return <-done
}

// Prompt asks for a name with a plain text field.
func (s *DialogService) Prompt() string {
	value, ok, err := s.app.Dialog.Prompt(application.PromptOptions{
		Title:        "Name this document",
		Message:      "The name is used for the exported file.",
		Placeholder:  "Untitled",
		DefaultValue: "Quarterly report",
		OKLabel:      "Rename",
		Window:       s.window,
	})
	if err != nil {
		return "error: " + err.Error()
	}
	if !ok {
		return "cancelled"
	}
	return "entered: " + value
}

// SecurePrompt asks for a passphrase with a secure text field, as a modal
// alert not attached to the window.
func (s *DialogService) SecurePrompt() string {
	value, ok, err := s.app.Dialog.Prompt(application.PromptOptions{
		Title:       "Unlock archive",
		Message:     "Enter the passphrase for archive.zip.",
		Placeholder: "Passphrase",
		Secure:      true,
		OKLabel:     "Unlock",
	})
	if err != nil {
		return "error: " + err.Error()
	}
	if !ok {
		return "cancelled"
	}
	return fmt.Sprintf("entered %d characters", len(value))
}

// OpenImages opens files conforming to public.image (any image type,
// including formats not listed by extension) plus PDF by extension.
func (s *DialogService) OpenImages() string {
	paths, err := s.app.Dialog.OpenFile().
		SetTitle("Choose images").
		SetMessage("Any image type, or a PDF").
		AddFilter("PDF", "*.pdf").
		AddContentType("public.image").
		AttachToWindow(s.window).
		PromptForMultipleSelection()
	if err != nil {
		return "error: " + err.Error()
	}
	if len(paths) == 0 {
		return "cancelled"
	}
	names := make([]string, len(paths))
	for i, path := range paths {
		names[i] = filepath.Base(path)
	}
	return "chose: " + strings.Join(names, ", ")
}

// Export shows a save panel with a Format pop-up, a custom name field
// label and Finder tags.
func (s *DialogService) Export() string {
	formats := []application.DialogFormat{
		{Label: "PNG image", Extension: "png", UTI: "public.png"},
		{Label: "PDF document", Extension: "pdf", UTI: "com.adobe.pdf"},
		{Label: "Markdown", Extension: "md", UTI: "net.daringfireball.markdown"},
	}
	dialog := s.app.Dialog.SaveFile().
		SetMessage("Choose a format for the export.").
		SetFilename("Quarterly report.png").
		SetNameFieldLabel("Export As:").
		SetTags([]string{"Reports", "Draft"}).
		SetButtonText("Export").
		AttachToWindow(s.window)
	dialog.SetFormats(formats, 0, func(index int) {
		log.Printf("format changed to %s", formats[index].Label)
	})
	path, err := dialog.PromptForSingleSelection()
	if err != nil {
		return "error: " + err.Error()
	}
	if path == "" {
		return fmt.Sprintf("cancelled (format %s)", formats[dialog.SelectedFormat()].Label)
	}
	return fmt.Sprintf("export %s as %s", filepath.Base(path), formats[dialog.SelectedFormat()].Label)
}

// PickColor opens the system colour panel with live updates.
func (s *DialogService) PickColor() string {
	color, ok, err := s.app.Dialog.PickColor(application.ColorPickerOptions{
		Initial:    application.NewRGB(52, 120, 246),
		ShowsAlpha: true,
		Title:      "Accent colour",
		OnChange: func(color application.RGBA) {
			s.app.Event.Emit("dialog:color", cssColor(color))
		},
	})
	if err != nil {
		return "error: " + err.Error()
	}
	if !ok {
		return "closed without changing the colour"
	}
	return "picked " + cssColor(color)
}

// PickFont opens the system font panel with live updates.
func (s *DialogService) PickFont() string {
	font, ok, err := s.app.Dialog.PickFont(application.FontPickerOptions{
		Family: "Helvetica Neue",
		Size:   18,
		OnChange: func(font application.FontDescriptor) {
			s.app.Event.Emit("dialog:font", font)
		},
	})
	if err != nil {
		return "error: " + err.Error()
	}
	if !ok {
		return "closed without changing the font"
	}
	return fmt.Sprintf("picked %s %s (%s) at %.1fpt", font.Family, font.Face, font.PostScriptName, font.Size)
}

func cssColor(color application.RGBA) string {
	return fmt.Sprintf("rgba(%d, %d, %d, %.2f)", color.Red, color.Green, color.Blue, float64(color.Alpha)/255)
}

func main() {
	service := &DialogService{}

	app := application.New(application.Options{
		Name:        "Native Dialogs",
		Description: "macOS alert, panel and picker extras",
		Logger:      application.DefaultLogger(slog.LevelDebug),
		Services: []application.Service{
			application.NewService(service),
		},
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})
	service.app = app

	service.window = app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Native Dialogs",
		Width:  720,
		Height: 560,
		URL:    "/",
	})

	menu := app.NewMenu()
	menu.AddRole(application.AppMenu)
	menu.AddRole(application.EditMenu)
	menu.AddRole(application.WindowMenu)
	dialogs := menu.AddSubmenu("Dialogs")
	dialogs.Add("Alert with suppression and help").OnClick(func(*application.Context) {
		log.Println(service.Suppression())
	})
	dialogs.Add("Text prompt").OnClick(func(*application.Context) {
		log.Println(service.Prompt())
	})
	dialogs.Add("Secure prompt").OnClick(func(*application.Context) {
		log.Println(service.SecurePrompt())
	})
	dialogs.Add("Open images (content types)").OnClick(func(*application.Context) {
		log.Println(service.OpenImages())
	})
	dialogs.Add("Export with format pop-up").OnClick(func(*application.Context) {
		log.Println(service.Export())
	})
	dialogs.Add("Colour panel").OnClick(func(*application.Context) {
		log.Println(service.PickColor())
	})
	dialogs.Add("Font panel").OnClick(func(*application.Context) {
		log.Println(service.PickFont())
	})
	app.Menu.Set(menu)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
