package term

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	PrimaryColor  = lipgloss.Color("#CE0000")
	SuccessColor  = lipgloss.Color("#10B981")
	WarningColor  = lipgloss.Color("#F59E0B")
	ErrorColor    = lipgloss.Color("#EF4444")
	MutedColor    = lipgloss.Color("#6B7280")
	HeaderBgColor = lipgloss.Color("#1A1A1A")

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(PrimaryColor).
			Padding(0, 1)

	SectionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(PrimaryColor).
			MarginTop(1).
			MarginBottom(1)

	OkStyle = lipgloss.NewStyle().
		Foreground(SuccessColor)

	WarnStyle = lipgloss.NewStyle().
			Foreground(WarningColor)

	ErrStyle = lipgloss.NewStyle().
			Foreground(ErrorColor)

	MutedStyle = lipgloss.NewStyle().
			Foreground(MutedColor)

	TableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(HeaderBgColor).
				Padding(0, 1)

	TableCellStyle = lipgloss.NewStyle().
			Padding(0, 1)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#374151")).
			Padding(1, 2)
)

func StatusIcon(ok bool) string {
	if ok {
		return OkStyle.Render("✓")
	}
	return ErrStyle.Render("✗")
}

func StatusIconTri(status string) string {
	switch status {
	case "ok":
		return OkStyle.Render("✓")
	case "warning":
		return WarnStyle.Render("!")
	case "missing", "error":
		return ErrStyle.Render("✗")
	default:
		return MutedStyle.Render("?")
	}
}

func RenderTable(rows [][]string) string {
	var b strings.Builder
	maxKeyLen := 0
	for _, row := range rows {
		if len(row[0]) > maxKeyLen {
			maxKeyLen = len(row[0])
		}
	}

	for _, row := range rows {
		key := MutedStyle.Render(padRight(row[0], maxKeyLen))
		val := row[1]
		b.WriteString("  " + key + "  " + val + "\n")
	}

	return b.String()
}

func padRight(s string, length int) string {
	for len(s) < length {
		s += " "
	}
	return s
}
