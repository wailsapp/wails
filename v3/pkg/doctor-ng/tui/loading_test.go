package tui

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
	doctorng "github.com/wailsapp/wails/v3/pkg/doctor-ng"
)

func TestLoadingFrames(t *testing.T) {
	model := NewModel()
	for _, glyph := range []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷", "⣾"} {
		if got, want := ansi.Strip(model.View()), fmt.Sprintf("\n  %s  Scanning system...\n", glyph); got != want {
			t.Fatalf("loading view = %q, want %q", got, want)
		}
		next, command := model.Update(loadingTickMsg(model.scanGeneration))
		if command == nil {
			t.Fatal("active scan stopped its clock")
		}
		model = next.(Model)
	}
}

func TestLoadingInitClock(t *testing.T) {
	model := NewModel()
	started := time.Now()
	commands := model.Init()().(tea.BatchMsg)
	if len(commands) != 2 {
		t.Fatalf("initial commands = %d, want clock and scan", len(commands))
	}
	message := commands[0]() // Do not execute the real doctor scan.
	if tick, ok := message.(loadingTickMsg); !ok || uint64(tick) != model.scanGeneration {
		t.Fatalf("clock message = %#v", message)
	}
	if elapsed := time.Since(started); elapsed < 90*time.Millisecond {
		t.Fatalf("clock fired too early after %v", elapsed)
	}
}

func TestLoadingClockStopsAfterReport(t *testing.T) {
	for _, scanError := range []error{nil, errors.New("scan failed")} {
		t.Run(fmt.Sprint(scanError), func(t *testing.T) {
			model := NewModel()
			report := &doctorng.Report{Summary: "finished"}
			ready, command := model.Update(reportReadyMsg{report: report, err: scanError})
			done := ready.(Model)
			if command != nil || done.state != stateReport || done.report != report || done.err != scanError {
				t.Fatal("scan completion did not retain its result and stop")
			}
			next, command := done.Update(loadingTickMsg(model.scanGeneration))
			if command != nil || next.(Model).loadingFrame != done.loadingFrame {
				t.Fatal("completed report continued its loading animation")
			}
			if got := ansi.Strip(next.View()); !strings.Contains(got, "Wails Doctor") && !strings.Contains(got, "scan failed") {
				t.Fatalf("completed view = %q", got)
			}
		})
	}
}

func TestLoadingRestartRejectsOldClock(t *testing.T) {
	for _, restart := range []struct {
		name    string
		message tea.Msg
	}{
		{"refresh", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}},
		{"install completed", installCompleteMsg{}},
	} {
		t.Run(restart.name, func(t *testing.T) {
			initial := NewModel()
			finished, _ := initial.Update(reportReadyMsg{})
			restarted, command := finished.Update(restart.message)
			model := restarted.(Model)
			if model.state != stateLoading || model.scanGeneration == initial.scanGeneration || command == nil {
				t.Fatal("restart did not start a new scan and clock generation")
			}
			commands := command().(tea.BatchMsg)
			if len(commands) != 2 {
				t.Fatalf("restart commands = %d, want clock and scan", len(commands))
			}
			currentTick := commands[0]()
			next, command := model.Update(loadingTickMsg(initial.scanGeneration))
			if command != nil || next.(Model).loadingFrame != model.loadingFrame {
				t.Fatal("stale clock advanced or restarted animation")
			}
			next, command = model.Update(currentTick)
			if command == nil || next.(Model).loadingFrame == model.loadingFrame {
				t.Fatal("new scan clock did not advance")
			}
		})
	}
}

func TestLoadingAcceptsNavigation(t *testing.T) {
	model := NewModel()
	sized, command := model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	if command != nil || sized.(Model).width != 80 || sized.(Model).height != 24 {
		t.Fatal("resize was not retained")
	}
	help, _ := sized.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	if !help.(Model).showHelp {
		t.Fatal("help key was not retained")
	}
	_, quit := help.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if quit == nil {
		t.Fatal("interrupt did not quit")
	}
	if _, ok := quit().(tea.QuitMsg); !ok {
		t.Fatal("interrupt command did not produce quit")
	}
}

func TestLoadingClockIsolation(t *testing.T) {
	first, second := NewModel(), NewModel()
	next, command := second.Update(loadingTickMsg(first.scanGeneration))
	if command != nil || next.(Model).loadingFrame != second.loadingFrame {
		t.Fatal("another doctor's clock advanced this model")
	}
}
