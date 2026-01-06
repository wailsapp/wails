package tui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	doctorng "github.com/wailsapp/wails/v3/pkg/doctor-ng"
)

type installState int

const (
	installConfirm installState = iota
	installRunning
	installDone
	installError
)

type InstallModel struct {
	deps         doctorng.DependencyList
	state        installState
	currentIndex int
	output       []string
	err          error
	selectedBtn  int
}

type installStartMsg struct{}
type installProgressMsg struct {
	output string
}
type installDoneMsg struct {
	err error
}

func NewInstallModel(deps doctorng.DependencyList) InstallModel {
	return InstallModel{
		deps:  deps,
		state: installConfirm,
	}
}

func (m InstallModel) Init() tea.Cmd {
	return nil
}

func (m InstallModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "left", "h":
			if m.state == installConfirm && m.selectedBtn > 0 {
				m.selectedBtn--
			}
		case "right", "l":
			if m.state == installConfirm && m.selectedBtn < 1 {
				m.selectedBtn++
			}
		case "tab":
			if m.state == installConfirm {
				m.selectedBtn = (m.selectedBtn + 1) % 2
			}
		case "enter":
			if m.state == installConfirm {
				if m.selectedBtn == 0 {
					return m, tea.Quit
				}
				m.state = installRunning
				return m, m.runInstall()
			}
			if m.state == installDone || m.state == installError {
				return m, tea.Quit
			}
		}

	case installProgressMsg:
		m.output = append(m.output, msg.output)

	case installDoneMsg:
		if msg.err != nil {
			m.state = installError
			m.err = msg.err
		} else {
			m.state = installDone
		}
	}

	return m, nil
}

func (m InstallModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(" Install Dependencies "))
	b.WriteString("\n\n")

	switch m.state {
	case installConfirm:
		b.WriteString(m.viewConfirm())
	case installRunning:
		b.WriteString(m.viewRunning())
	case installDone:
		b.WriteString(m.viewDone())
	case installError:
		b.WriteString(m.viewError())
	}

	return b.String()
}

func (m InstallModel) viewConfirm() string {
	var b strings.Builder

	b.WriteString("The following packages will be installed:\n\n")

	for _, dep := range m.deps {
		b.WriteString(fmt.Sprintf("  %s %s\n",
			errStyle.Render("•"),
			dep.Name))
		if dep.InstallCommand != "" {
			b.WriteString(fmt.Sprintf("    %s\n", mutedStyle.Render(dep.InstallCommand)))
		}
	}

	b.WriteString("\n")
	b.WriteString(warnStyle.Render("This will require sudo privileges"))
	b.WriteString("\n\n")

	cancelBtn := buttonStyle.Render(" Cancel ")
	installBtn := buttonStyle.Render(" Install ")

	if m.selectedBtn == 0 {
		cancelBtn = buttonActiveStyle.Render(" Cancel ")
	} else {
		installBtn = buttonActiveStyle.Render(" Install ")
	}

	b.WriteString(cancelBtn)
	b.WriteString("  ")
	b.WriteString(installBtn)
	b.WriteString("\n\n")
	b.WriteString(mutedStyle.Render("Use Tab or ←/→ to select, Enter to confirm"))

	return b.String()
}

func (m InstallModel) viewRunning() string {
	var b strings.Builder

	b.WriteString("Installing packages...\n\n")

	for _, line := range m.output {
		b.WriteString("  " + line + "\n")
	}

	b.WriteString("\n")
	b.WriteString(mutedStyle.Render("Please wait..."))

	return b.String()
}

func (m InstallModel) viewDone() string {
	var b strings.Builder

	b.WriteString(okStyle.Render("✓ Installation complete!"))
	b.WriteString("\n\n")

	for _, line := range m.output {
		b.WriteString("  " + line + "\n")
	}

	b.WriteString("\n")
	b.WriteString(mutedStyle.Render("Press Enter to exit and re-run doctor to verify"))

	return b.String()
}

func (m InstallModel) viewError() string {
	var b strings.Builder

	b.WriteString(errStyle.Render("✗ Installation failed"))
	b.WriteString("\n\n")

	for _, line := range m.output {
		b.WriteString("  " + line + "\n")
	}

	if m.err != nil {
		b.WriteString("\n")
		b.WriteString(errStyle.Render(fmt.Sprintf("Error: %v", m.err)))
	}

	b.WriteString("\n\n")
	b.WriteString(mutedStyle.Render("Press Enter to exit"))

	return b.String()
}

func (m InstallModel) runInstall() tea.Cmd {
	return func() tea.Msg {
		for _, dep := range m.deps {
			if dep.InstallCommand == "" {
				continue
			}

			cmd := exec.Command("sh", "-c", dep.InstallCommand)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				return installDoneMsg{err: fmt.Errorf("failed to install %s: %w", dep.Name, err)}
			}
		}

		return installDoneMsg{}
	}
}

func RunInstaller(deps doctorng.DependencyList) error {
	if len(deps) == 0 {
		fmt.Println(okStyle.Render("All dependencies are already installed!"))
		return nil
	}

	p := tea.NewProgram(NewInstallModel(deps))
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running installer: %w", err)
	}
	return nil
}
