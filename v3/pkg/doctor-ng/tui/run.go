package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func Run() error {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running TUI: %w", err)
	}
	return nil
}

func RunSimple() error {
	p := tea.NewProgram(NewModel())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running TUI: %w", err)
	}
	return nil
}

func RunNonInteractive() error {
	m := NewModel()

	msg := runDoctor()
	reportMsg, ok := msg.(reportReadyMsg)
	if !ok {
		return fmt.Errorf("unexpected message type")
	}

	if reportMsg.err != nil {
		return reportMsg.err
	}

	m.report = reportMsg.report
	m.state = stateReport

	fmt.Fprint(os.Stdout, m.View())
	return nil
}
