package pulse

import (
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/report/pulse/ansi"
)

func TestVisibleWidth(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"", 0},
		{"hello", 5},
		{"hi\x1b[31m world\x1b[0m", 8},                                   // "hi world" = 8 visible cells
		{"\x1b]8;;file:///x\x1b\\click\x1b]8;;\x1b\\", 5},                // OSC 8 hyperlink wraps "click"
		{"\x1b]8;;file:///C:\\path\\to.go\x1b\\click\x1b]8;;\x1b\\", 5}, // URI with `\` characters: regression
		{"━━━━━╾─", 7},                                                   // block chars
		{strings.Repeat("a", 100), 100},
	}
	for _, c := range cases {
		if got := visibleWidth(c.in); got != c.want {
			t.Errorf("visibleWidth(%q) = %d, want %d", c.in, got, c.want)
		}
	}
}

// Truncating across an OSC 8 wrapper must not eat the wrapper or the
// terminating sequence. Regression for the OSC state machine.
func TestTruncatePreservesOSC8Wrapper(t *testing.T) {
	link := "\x1b]8;;file:///x.go\x1b\\click here\x1b]8;;\x1b\\"
	got := truncate(link, 6)
	if !strings.Contains(got, "\x1b]8;;file:///x.go\x1b\\") {
		t.Errorf("truncate dropped the OSC 8 opener: %q", got)
	}
}

func TestTruncate(t *testing.T) {
	// Plain text past width: ellipsis, no trailing reset.
	if got := truncate("hello world", 7); got != "hello …" {
		t.Errorf("truncate plain past width = %q, want %q", got, "hello …")
	}
	// Plain text inside width: pass through unchanged.
	if got := truncate("hi", 7); got != "hi" {
		t.Errorf("truncate plain inside width = %q, want %q", got, "hi")
	}
	// Styled input past width: ellipsis with trailing reset.
	in := "\x1b[31mhello world\x1b[0m"
	got := truncate(in, 7)
	if !strings.HasSuffix(got, "…"+ansi.Reset) {
		t.Errorf("truncate styled past width = %q, want suffix %q", got, "…"+ansi.Reset)
	}
	// Width=0: empty out.
	if got := truncate("anything", 0); got != "" {
		t.Errorf("truncate width 0 = %q, want empty", got)
	}
}

func TestPad(t *testing.T) {
	if got := padRight("hi", 5); got != "hi   " {
		t.Errorf("padRight short = %q", got)
	}
	if got := padRight("hello", 3); got != "hello" {
		t.Errorf("padRight overflow = %q, want unchanged", got)
	}
	if got := padLeft("42", 5); got != "   42" {
		t.Errorf("padLeft short = %q", got)
	}
	// SGR-styled input still measures by visible cells, not bytes.
	in := "\x1b[31mhi\x1b[0m"
	if got := padRight(in, 5); got != in+"   " {
		t.Errorf("padRight styled = %q, want suffix of 3 spaces", got)
	}
}

func TestEraseLines(t *testing.T) {
	if got := ansi.EraseLines(0); got != "" {
		t.Errorf("EraseLines(0) = %q, want empty", got)
	}
	// n=1 → clear current line only (no up-and-clear).
	got := ansi.EraseLines(1)
	if strings.Contains(got, "1A") {
		t.Errorf("EraseLines(1) should not move up, got %q", got)
	}
	// n=3 → one current-line clear plus 2 up-and-clears.
	got = ansi.EraseLines(3)
	upMoves := strings.Count(got, ansi.CSI+"1A")
	if upMoves != 2 {
		t.Errorf("EraseLines(3) up-moves = %d, want 2", upMoves)
	}
	clears := strings.Count(got, ansi.ClearLine)
	if clears != 3 {
		t.Errorf("EraseLines(3) clears = %d, want 3", clears)
	}
}
