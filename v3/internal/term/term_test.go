package term

import (
	"bytes"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/pterm/pterm"
)

func TestLiteralAndFormattedMessages(t *testing.T) {
	originalInfo, originalWarning := pterm.Info, pterm.Warning
	originalError, originalSuccess := pterm.Error, pterm.Success
	t.Cleanup(func() {
		pterm.Info, pterm.Warning = originalInfo, originalWarning
		pterm.Error, pterm.Success = originalError, originalSuccess
	})
	var output bytes.Buffer
	pterm.Info.Writer, pterm.Warning.Writer = &output, &output
	pterm.Error.Writer, pterm.Success.Writer = &output, &output
	message := "100% done: https://example.com/a%20b literal %s %%"
	for name, print := range map[string]func(any){
		"Info": Info, "Warning": Warning, "Error": Error, "Success": Success,
	} {
		t.Run(name, func(t *testing.T) {
			output.Reset()
			print(message)
			if got := output.String(); !strings.Contains(got, message) {
				t.Errorf("message was changed: %q", got)
			}
		})
	}
	for name, print := range map[string]func(any, ...any){"Infof": Infof, "Warningf": Warningf} {
		t.Run(name, func(t *testing.T) {
			output.Reset()
			print("progress %d%%", 42)
			if got := output.String(); !strings.Contains(got, "progress 42%") {
				t.Errorf("formatted message = %q", got)
			}
			output.Reset()
			print("literal %%s")
			if got := output.String(); !strings.Contains(got, "literal %s") {
				t.Errorf("zero-argument format = %q", got)
			}
		})
	}
}

func TestRenderTableAlignsDisplayWidths(t *testing.T) {
	rows := [][]string{
		{"Name", "VALUE"},
		{"é", "VALUE"},
		{"界界界", "VALUE"},
		{"e\u0301", "VALUE"},
		{"\x1b[31mred\x1b[0m", "VALUE"},
		nil,
		{},
	}
	lines := strings.Split(strings.TrimSuffix(RenderTable(rows), "\n"), "\n")
	if len(lines) != 5 {
		t.Fatalf("rendered %d lines, want 5", len(lines))
	}
	for i, line := range lines {
		prefix, value, found := strings.Cut(line, "VALUE")
		if !found || value != "" {
			t.Errorf("row %d value changed: %q", i, line)
		}
		if width := lipgloss.Width(prefix); width != 10 {
			t.Errorf("row %d value starts after %d columns, want 10: %q", i, width, line)
		}
	}
}

func TestDisplayWidth(t *testing.T) {
	for _, tc := range []struct {
		text string
		want int
	}{
		{"", 0},
		{"plain", 5},
		{"é", 1},
		{"e\u0301", 1},
		{"界界", 4},
		{"👩‍💻", 2},
		{"\x1b[31mred\x1b[0m", 3},
		{"wide界\nx", 6},
		{"\n", 0},
		{"long\nshorter\n", 7},
	} {
		if got := displayWidth(tc.text); got != tc.want {
			t.Errorf("displayWidth(%q) = %d, want %d", tc.text, got, tc.want)
		}
	}
}
