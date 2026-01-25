package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type GreetService struct {
	app *application.App
}

type WindowSummary struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type CurrentWindowReport struct {
	Caller               WindowSummary  `json:"caller"`
	Current              *WindowSummary `json:"current,omitempty"`
	CurrentMatchesCaller bool           `json:"currentMatchesCaller"`
}

func (g *GreetService) Greet(name string) string {
	return "Hello " + name + "!"
}

// OpenChildWindow creates a new window that loads the same frontend entrypoint,
// but with query parameters enabling the "child window" UI.
func (g *GreetService) OpenChildWindow() string {
	if g.app == nil {
		return ""
	}

	name := fmt.Sprintf("child-%d", time.Now().UnixNano())
	childURL := "/?child=1&name=" + url.QueryEscape(name)

	g.app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:             name,
		Title:            "Child Window (" + name + ")",
		Width:            520,
		Height:           420,
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              childURL,
	})

	return name
}

func (g *GreetService) ListWindows() []WindowSummary {
	if g.app == nil {
		return nil
	}

	ws := g.app.Window.GetAll()
	out := make([]WindowSummary, 0, len(ws))
	for _, w := range ws {
		out = append(out, WindowSummary{
			ID:   w.ID(),
			Name: w.Name(),
		})
	}
	return out
}

func (g *GreetService) CloseByName(name string) bool {
	if g.app == nil {
		return false
	}
	w, ok := g.app.Window.GetByName(name)
	if !ok {
		return false
	}
	w.Close()
	return true
}

// CloseAfterCurrentByName calls App.Window.Current() first, then closes the window by name.
// This is the variant intended to reproduce the bug.
func (g *GreetService) CloseAfterCurrentByName(name string) bool {
	if g.app == nil {
		return false
	}

	// The problematic call:
	_ = g.app.Window.Current()

	return g.CloseByName(name)
}

// CloseUsingCurrent closes whichever window Wails considers "current".
func (g *GreetService) CloseUsingCurrent() WindowSummary {
	if g.app == nil {
		return WindowSummary{}
	}

	w := g.app.Window.Current()
	info := WindowSummary{ID: w.ID(), Name: w.Name()}
	w.Close()
	return info
}

// ReportCurrent returns which window made this call (Caller) and what
// App.Window.Current() returns at that moment (Current).
func (g *GreetService) ReportCurrent(ctx context.Context) CurrentWindowReport {
	var caller application.Window
	if ctx != nil {
		if w, ok := ctx.Value(application.WindowKey).(application.Window); ok {
			caller = w
		}
	}

	report := CurrentWindowReport{
		Caller: WindowSummary{
			ID:   0,
			Name: "(unknown)",
		},
	}
	if caller != nil {
		report.Caller = WindowSummary{ID: caller.ID(), Name: caller.Name()}
	}

	if g.app == nil {
		return report
	}

	current := g.app.Window.Current()
	if current != nil {
		s := WindowSummary{ID: current.ID(), Name: current.Name()}
		report.Current = &s
		report.CurrentMatchesCaller = caller != nil && caller.ID() == current.ID()
	}

	log.Printf("[ReportCurrent] caller=%+v current=%+v match=%v", report.Caller, report.Current, report.CurrentMatchesCaller)
	return report
}
