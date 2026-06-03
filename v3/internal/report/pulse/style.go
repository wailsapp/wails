package pulse

import (
	"strings"

	"github.com/wailsapp/wails/v3/internal/report/pulse/ansi"
)

// styler binds a Profile so call sites can paint without naming the profile
// every time. Methods return strings ready to write to the terminal.
//
// All methods are nil-safe and return their input unchanged if the styler's
// profile is None — this lets unit tests build a no-op styler.
type styler struct {
	p Profile
}

func newStyler(p Profile) *styler { return &styler{p: p} }

func (s *styler) fg(c Colour, text string) string {
	if s == nil || s.p == ProfileNone {
		return text
	}
	return c.FG(s.p) + text + ansi.Reset
}

func (s *styler) bg(c Colour, text string) string {
	if s == nil || s.p == ProfileNone {
		return text
	}
	return c.BG(s.p) + text + ansi.Reset
}

// bold renders text in bold without changing colour.
func (s *styler) bold(text string) string {
	if s == nil || s.p == ProfileNone {
		return text
	}
	return ansi.Bold() + text + ansi.Reset
}

// faint renders text dimmed without changing colour.
func (s *styler) faint(text string) string {
	if s == nil || s.p == ProfileNone {
		return text
	}
	return ansi.Faint() + text + ansi.Reset
}

// accentBold is the brand emphasis: bold accent-coloured. Used only for the
// verb in the build header and similar high-attention text.
func (s *styler) accentBold(text string) string {
	if s == nil || s.p == ProfileNone {
		return text
	}
	return ansi.Bold() + Accent.FG(s.p) + text + ansi.Reset
}

// link wraps text in an OSC 8 hyperlink pointing at uri. Terminals that don't
// understand OSC 8 print text only (the escape sequence is invisible). Empty
// uri returns text unwrapped.
func (s *styler) link(uri, text string) string {
	if s == nil || s.p == ProfileNone || uri == "" {
		return text
	}
	return "\x1b]8;;" + uri + "\x1b\\" + text + "\x1b]8;;\x1b\\"
}

// gradientRule renders a horizontal rule of width cells, with each cell's
// foreground colour linearly interpolated from `from` to `to`. On TrueColor
// terminals this paints a smooth fade across the bar; on 256/16 terminals
// the same logic falls back to a single solid colour (we don't try to step
// through a palette range). The visual is a small luxury that signals
// "wake takes terminal aesthetics seriously" without costing layout.
func (s *styler) gradientRule(glyph string, width int, from, to Colour) string {
	if s == nil || s.p != ProfileTrueCol {
		// Solid fallback — use the destination colour (typically Accent).
		return s.fg(to, strings.Repeat(glyph, width))
	}
	r1, g1, b1 := hexRGB(from.True)
	r2, g2, b2 := hexRGB(to.True)
	var b strings.Builder
	for i := 0; i < width; i++ {
		t := float64(i) / float64(width-1)
		r := uint8(float64(r1) + (float64(r2)-float64(r1))*t)
		g := uint8(float64(g1) + (float64(g2)-float64(g1))*t)
		bl := uint8(float64(b1) + (float64(b2)-float64(b1))*t)
		b.WriteString(ansi.CSI)
		b.WriteString(formatRGB(r, g, bl))
		b.WriteString("m")
		b.WriteString(glyph)
	}
	b.WriteString(ansi.Reset)
	return b.String()
}

// formatRGB returns the "38;2;R;G;B" body for a TrueColor SGR. Pulled out
// of gradientRule so the hot loop avoids a Sprintf per cell.
func formatRGB(r, g, b uint8) string {
	var buf [16]byte
	out := buf[:0]
	out = append(out, "38;2;"...)
	out = appendU8(out, r)
	out = append(out, ';')
	out = appendU8(out, g)
	out = append(out, ';')
	out = appendU8(out, b)
	return string(out)
}

func appendU8(b []byte, v uint8) []byte {
	if v >= 100 {
		b = append(b, '0'+v/100)
		v %= 100
		b = append(b, '0'+v/10)
		b = append(b, '0'+v%10)
		return b
	}
	if v >= 10 {
		b = append(b, '0'+v/10)
		b = append(b, '0'+v%10)
		return b
	}
	b = append(b, '0'+v)
	return b
}

// padRight returns s padded with spaces to at least width visible cells,
// counting runes not bytes. SGR escapes are stripped for width measurement.
func padRight(s string, width int) string {
	v := visibleWidth(s)
	if v >= width {
		return s
	}
	return s + strings.Repeat(" ", width-v)
}

// padLeft returns s padded with spaces on the left to at least width visible
// cells. Used for right-aligned numbers in the summary table.
func padLeft(s string, width int) string {
	v := visibleWidth(s)
	if v >= width {
		return s
	}
	return strings.Repeat(" ", width-v) + s
}

// truncate returns s truncated to at most width visible cells, with an
// ellipsis at the end if anything was removed. SGR escapes pass through; we
// don't try to split inside one (the wake reporter never builds nested-style
// strings, so a simple state machine is enough).
//
// A trailing reset (\x1b[0m) is appended only if the input actually contained
// any SGR — so a plain-text caller (e.g. the failure-panel body in NO_COLOR
// mode) doesn't have stray "\x1b[0m" bytes appear in its output.
func truncate(s string, width int) string {
	if width <= 0 {
		return ""
	}
	if visibleWidth(s) <= width {
		return s
	}
	var (
		b       strings.Builder
		inEsc   bool
		hadSGR  bool
		visible int
	)
	for _, r := range s {
		if r == '\x1b' {
			inEsc = true
			hadSGR = true
			b.WriteRune(r)
			continue
		}
		if inEsc {
			b.WriteRune(r)
			if (r >= 0x40 && r <= 0x7e) && r != '[' {
				inEsc = false
			}
			continue
		}
		if visible+1 > width-1 { // reserve 1 cell for the ellipsis
			b.WriteRune('…')
			if hadSGR {
				b.WriteString(ansi.Reset)
			}
			return b.String()
		}
		b.WriteRune(r)
		visible++
	}
	return b.String()
}

// visibleWidth returns the number of visible cells s would occupy, skipping
// SGR escapes and OSC 8 hyperlink sequences. It does *not* honour East Asian
// wide characters — wake task names are ASCII in practice.
func visibleWidth(s string) int {
	var (
		n     int
		inEsc bool
		inOsc bool
	)
	for _, r := range s {
		if inOsc {
			if r == '\x07' || r == '\\' { // BEL or ESC \
				inOsc = false
			}
			continue
		}
		if r == '\x1b' {
			inEsc = true
			continue
		}
		if inEsc {
			// OSC sequences start ESC ].
			if r == ']' {
				inOsc = true
				inEsc = false
				continue
			}
			// CSI/other sequences end at the final byte (0x40-0x7e).
			if r >= 0x40 && r <= 0x7e && r != '[' {
				inEsc = false
			}
			continue
		}
		n++
	}
	return n
}
