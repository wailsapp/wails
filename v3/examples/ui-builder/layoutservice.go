package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// LayoutService is the Go side of the UI builder. The frontend keeps the
// design in memory and calls these methods to persist it, reopen it and export
// it as a Wails frontend. All file access goes through native dialogs so the
// user always picks where their files live.
type LayoutService struct{}

// LayoutFile is what OpenLayout hands back to the frontend: the path the file
// was read from (so subsequent saves can skip the dialog) and its JSON payload.
type LayoutFile struct {
	Path string `json:"path"`
	Data string `json:"data"`
}

// ExportFile is one file of an exported frontend, with a path relative to the
// export directory (e.g. "src/main.js").
type ExportFile struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

// SaveLayout writes the layout JSON to disk. If path is empty a native "Save
// As" dialog is shown, seeded with a filename derived from name. The path the
// file was written to is returned so the frontend can remember it; an empty
// string means the user cancelled.
func (s *LayoutService) SaveLayout(path string, name string, data string) (string, error) {
	if !json.Valid([]byte(data)) {
		return "", errors.New("layout is not valid JSON")
	}
	if path == "" {
		chosen, err := application.Get().Dialog.
			SaveFileWithOptions(&application.SaveFileDialogOptions{Title: "Save layout"}).
			SetFilename(safeFilename(name, "layout")+".uib.json").
			AddFilter("UI Builder layout", "*.uib.json;*.json").
			CanCreateDirectories(true).
			PromptForSingleSelection()
		if err != nil {
			return "", err
		}
		if chosen == "" {
			return "", nil
		}
		path = chosen
	}
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		return "", fmt.Errorf("could not write %s: %w", filepath.Base(path), err)
	}
	return path, nil
}

// OpenLayout shows a native open dialog and returns the selected layout. A nil
// result means the user cancelled.
func (s *LayoutService) OpenLayout() (*LayoutFile, error) {
	path, err := application.Get().Dialog.OpenFile().
		SetTitle("Open layout").
		AddFilter("UI Builder layout", "*.uib.json;*.json").
		CanChooseFiles(true).
		CanChooseDirectories(false).
		PromptForSingleSelection()
	if err != nil {
		return nil, err
	}
	if path == "" {
		return nil, nil
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read %s: %w", filepath.Base(path), err)
	}
	if !json.Valid(raw) {
		return nil, fmt.Errorf("%s is not a valid layout file", filepath.Base(path))
	}
	return &LayoutFile{Path: path, Data: string(raw)}, nil
}

// ExportFrontend asks for a destination folder and writes the generated
// frontend into a new "<name>-frontend" directory inside it. Returns the
// directory written, or an empty string if the user cancelled.
func (s *LayoutService) ExportFrontend(name string, files []ExportFile) (string, error) {
	if len(files) == 0 {
		return "", errors.New("nothing to export")
	}
	parent, err := application.Get().Dialog.OpenFile().
		SetTitle("Choose where to export the frontend").
		SetButtonText("Export here").
		CanChooseDirectories(true).
		CanChooseFiles(false).
		CanCreateDirectories(true).
		PromptForSingleSelection()
	if err != nil {
		return "", err
	}
	if parent == "" {
		return "", nil
	}
	dir := filepath.Join(parent, safeFilename(name, "app")+"-frontend")
	for _, f := range files {
		rel := filepath.Clean(filepath.FromSlash(f.Path))
		if rel == "." || strings.HasPrefix(rel, "..") || filepath.IsAbs(rel) {
			return "", fmt.Errorf("refusing to write outside the export directory: %q", f.Path)
		}
		target := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return "", err
		}
		if err := os.WriteFile(target, []byte(f.Content), 0o644); err != nil {
			return "", fmt.Errorf("could not write %s: %w", rel, err)
		}
	}
	return dir, nil
}

// safeFilename turns a free-form design name into something every filesystem
// accepts, falling back to fallback when nothing usable is left.
func safeFilename(name, fallback string) string {
	var b strings.Builder
	prevDash := false
	for _, r := range strings.ToLower(strings.TrimSpace(name)) {
		switch {
		case (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'):
			b.WriteRune(r)
			prevDash = false
		case !prevDash && b.Len() > 0:
			b.WriteRune('-')
			prevDash = true
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		return fallback
	}
	return out
}
