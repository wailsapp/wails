package tui

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	doctorng "github.com/wailsapp/wails/v3/pkg/doctor-ng"
)

type state int

const (
	stateLoading state = iota
	stateReport
)

type Model struct {
	state        state
	report       *doctorng.Report
	spinner      spinner.Model
	err          error
	width        int
	height       int
	showHelp     bool
	copiedNotice bool
}

type reportReadyMsg struct {
	report *doctorng.Report
	err    error
}

type installCompleteMsg struct {
	err error
}

func NewModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	return Model{
		state:   stateLoading,
		spinner: s,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		runDoctor,
	)
}

func runDoctor() tea.Msg {
	d := doctorng.New()
	report, err := d.Run()
	return reportReadyMsg{report: report, err: err}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "?":
			m.showHelp = !m.showHelp

		case "i":
			if m.state == stateReport && m.report != nil {
				missing := m.report.Dependencies.RequiredMissing()
				if len(missing) > 0 {
					return m, tea.ExecProcess(
						createInstallCmd(missing),
						func(err error) tea.Msg { return installCompleteMsg{err: err} },
					)
				}
			}

		case "r":
			if m.state == stateReport {
				m.state = stateLoading
				return m, tea.Batch(m.spinner.Tick, runDoctor)
			}
		case "c":
			if m.state == stateReport && m.report != nil {
				text := m.generateClipboardText()
				if err := clipboard.WriteAll(text); err == nil {
					m.copiedNotice = true
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case reportReadyMsg:
		m.report = msg.report
		m.err = msg.err
		m.state = stateReport

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case installCompleteMsg:
		m.state = stateLoading
		return m, tea.Batch(m.spinner.Tick, runDoctor)
	}

	return m, nil
}

func (m Model) View() string {
	switch m.state {
	case stateLoading:
		return m.viewLoading()
	case stateReport:
		return m.viewReport()
	default:
		return ""
	}
}

func (m Model) viewLoading() string {
	return fmt.Sprintf("\n  %s Scanning system...\n", m.spinner.View())
}

func (m Model) viewReport() string {
	if m.err != nil {
		return errStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}
	if m.report == nil {
		return ""
	}

	var b strings.Builder

	b.WriteString(titleStyle.Render(" Wails Doctor "))
	b.WriteString("\n\n")

	b.WriteString(m.renderSystemInfo())
	b.WriteString(m.renderBuildInfo())
	b.WriteString(m.renderDependencies())
	b.WriteString(m.renderDiagnostics())
	b.WriteString(m.renderSummary())

	if m.showHelp {
		b.WriteString(m.renderHelp())
	} else if m.copiedNotice {
		b.WriteString(okStyle.Render("Copied to clipboard!") + " " + helpStyle.Render("Press ? for help, q to quit"))
	} else {
		b.WriteString(helpStyle.Render("Press c to copy, ? for help, q to quit"))
	}

	return b.String()
}

func (m Model) renderSystemInfo() string {
	var b strings.Builder
	b.WriteString(sectionStyle.Render("System"))
	b.WriteString("\n")

	sys := m.report.System
	rows := [][]string{
		{"OS", fmt.Sprintf("%s %s", sys.OS.Name, sys.OS.Version)},
		{"Platform", fmt.Sprintf("%s/%s", sys.OS.Platform, sys.OS.Arch)},
	}

	if len(sys.Hardware.CPUs) > 0 {
		rows = append(rows, []string{"CPU", sys.Hardware.CPUs[0].Model})
	}
	if len(sys.Hardware.GPUs) > 0 {
		gpuInfo := sys.Hardware.GPUs[0].Name
		if sys.Hardware.GPUs[0].Vendor != "" {
			gpuInfo += " (" + sys.Hardware.GPUs[0].Vendor + ")"
		}
		rows = append(rows, []string{"GPU", gpuInfo})
	}
	rows = append(rows, []string{"Memory", sys.Hardware.Memory})

	keys := make([]string, 0, len(sys.PlatformExtras))
	for k := range sys.PlatformExtras {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		rows = append(rows, []string{k, sys.PlatformExtras[k]})
	}

	b.WriteString(renderTable(rows))
	return b.String()
}

func (m Model) renderBuildInfo() string {
	var b strings.Builder
	b.WriteString(sectionStyle.Render("Build Environment"))
	b.WriteString("\n")

	build := m.report.Build
	rows := [][]string{
		{"Wails", build.WailsVersion},
		{"Go", build.GoVersion},
		{"CGO", fmt.Sprintf("%v", build.CGOEnabled)},
	}

	b.WriteString(renderTable(rows))
	return b.String()
}

func (m Model) renderDependencies() string {
	var b strings.Builder
	b.WriteString(sectionStyle.Render("Dependencies"))
	b.WriteString("\n")

	for _, dep := range m.report.Dependencies {
		icon := statusIconTri(dep.Status.String())
		name := dep.Name
		version := dep.Version
		if version == "" {
			version = mutedStyle.Render("not installed")
		}

		row := fmt.Sprintf("  %s %-25s %s", icon, name, version)

		if !dep.Required {
			row += mutedStyle.Render(" (optional)")
		}

		b.WriteString(row)
		b.WriteString("\n")
	}

	return b.String()
}

func (m Model) renderDiagnostics() string {
	if len(m.report.Diagnostics) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString(sectionStyle.Render("Issues Found"))
	b.WriteString("\n")

	for _, diag := range m.report.Diagnostics {
		var icon string
		var style lipgloss.Style
		switch diag.Severity {
		case doctorng.SeverityError:
			icon = "✗"
			style = errStyle
		case doctorng.SeverityWarning:
			icon = "!"
			style = warnStyle
		default:
			icon = "i"
			style = mutedStyle
		}

		b.WriteString(fmt.Sprintf("  %s %s: %s\n",
			style.Render(icon),
			style.Render(diag.Name),
			diag.Message))

		if diag.Fix != nil && diag.Fix.Command != "" {
			b.WriteString(fmt.Sprintf("    Fix: %s\n", mutedStyle.Render(diag.Fix.Command)))
		}
	}

	return b.String()
}

func (m Model) renderSummary() string {
	var b strings.Builder
	b.WriteString("\n")

	if m.report.Ready {
		b.WriteString(okStyle.Render("✓ " + m.report.Summary))
	} else {
		b.WriteString(errStyle.Render("✗ " + m.report.Summary))

		missing := m.report.Dependencies.RequiredMissing()
		if len(missing) > 0 {
			b.WriteString("\n\n")
			b.WriteString(mutedStyle.Render("Press 'i' to install missing dependencies"))
		}
	}

	b.WriteString("\n")
	return b.String()
}

func (m Model) renderHelp() string {
	var b strings.Builder
	b.WriteString("\n")

	help := "Keyboard Shortcuts:\n"
	help += "  c           Copy report to clipboard\n"
	help += "  r           Refresh / re-scan system\n"

	if m.report != nil && len(m.report.Dependencies.RequiredMissing()) > 0 {
		help += "  i           Install missing dependencies\n"
	}

	help += "  ?           Toggle help\n"
	help += "  q           Quit"

	b.WriteString(boxStyle.Render(help))
	return b.String()
}

func renderTable(rows [][]string) string {
	var b strings.Builder
	maxKeyLen := 0
	for _, row := range rows {
		if len(row[0]) > maxKeyLen {
			maxKeyLen = len(row[0])
		}
	}

	for _, row := range rows {
		key := tableCellStyle.Render(fmt.Sprintf("%-*s", maxKeyLen, row[0]))
		val := row[1]
		b.WriteString(fmt.Sprintf("  %s  %s\n", mutedStyle.Render(key), val))
	}

	return b.String()
}

func createInstallCmd(deps doctorng.DependencyList) *exec.Cmd {
	var commands []string
	for _, dep := range deps {
		if dep.InstallCommand != "" {
			commands = append(commands, dep.InstallCommand)
		}
	}

	if len(commands) == 0 {
		return exec.Command("echo", "Nothing to install")
	}

	combined := strings.Join(commands, " && ")
	return exec.Command("sh", "-c", combined)
}

func (m Model) generateClipboardText() string {
	if m.report == nil {
		return ""
	}

	var b strings.Builder
	b.WriteString("Wails Doctor Report\n")
	b.WriteString("===================\n\n")

	sys := m.report.System
	b.WriteString("System:\n")
	b.WriteString(fmt.Sprintf("  OS: %s %s\n", sys.OS.Name, sys.OS.Version))
	b.WriteString(fmt.Sprintf("  Platform: %s/%s\n", sys.OS.Platform, sys.OS.Arch))
	if len(sys.Hardware.CPUs) > 0 {
		b.WriteString(fmt.Sprintf("  CPU: %s\n", sys.Hardware.CPUs[0].Model))
	}
	if len(sys.Hardware.GPUs) > 0 {
		b.WriteString(fmt.Sprintf("  GPU: %s\n", sys.Hardware.GPUs[0].Name))
	}
	b.WriteString(fmt.Sprintf("  Memory: %s\n", sys.Hardware.Memory))

	keys := make([]string, 0, len(sys.PlatformExtras))
	for k := range sys.PlatformExtras {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		b.WriteString(fmt.Sprintf("  %s: %s\n", k, sys.PlatformExtras[k]))
	}

	b.WriteString("\nBuild Environment:\n")
	b.WriteString(fmt.Sprintf("  Wails: %s\n", m.report.Build.WailsVersion))
	b.WriteString(fmt.Sprintf("  Go: %s\n", m.report.Build.GoVersion))

	b.WriteString("\nDependencies:\n")
	for _, dep := range m.report.Dependencies {
		status := "✓"
		if dep.Status != doctorng.StatusOK {
			status = "✗"
		}
		version := dep.Version
		if version == "" {
			version = "not installed"
		}
		optional := ""
		if !dep.Required {
			optional = " (optional)"
		}
		b.WriteString(fmt.Sprintf("  %s %s: %s%s\n", status, dep.Name, version, optional))
	}

	if len(m.report.Diagnostics) > 0 {
		b.WriteString("\nIssues:\n")
		for _, diag := range m.report.Diagnostics {
			b.WriteString(fmt.Sprintf("  - %s: %s\n", diag.Name, diag.Message))
		}
	}

	b.WriteString(fmt.Sprintf("\nStatus: %s\n", m.report.Summary))

	return b.String()
}
