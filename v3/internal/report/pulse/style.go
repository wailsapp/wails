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

// escState walks an ANSI/OSC-bearing string and reports for each rune
// whether it should count as a visible cell or be ignored as part of an
// escape sequence. Shared by truncate and visibleWidth so the two stay
// consistent — they used to drift apart.
//
// Terminator handling for OSC sequences is precise: BEL (0x07) ends an
// OSC, and an ESC followed by `\` (the so-called "string terminator")
// also ends it. The previous implementation collapsed those into "any
// `\` ends an OSC", which truncated hyperlink-wrapped text at the first
// backslash inside the URI.
type escState struct {
	inEsc bool // last rune was ESC (CSI/OSC opener)
	inOsc bool // inside an OSC sequence (ESC ] … ST)
	inSt  bool // inside an OSC string terminator (ESC waiting for \)
}

// step advances the state and returns whether r should be counted as
// visible. r is always written through to the output regardless; the
// caller decides what to do with the visibility flag.
func (e *escState) step(r rune) (visible bool) {
	if e.inOsc {
		switch {
		case r == '\x07': // BEL ends OSC immediately
			e.inOsc = false
		case e.inSt && r == '\\': // ESC \ ends OSC
			e.inOsc = false
			e.inSt = false
		case r == '\x1b': // ESC: maybe the ST starts here
			e.inSt = true
		default:
			e.inSt = false
		}
		return false
	}
	if e.inEsc {
		// ESC ] starts an OSC; otherwise treat as a CSI/SGR-shaped escape
		// that ends on the first byte in [0x40, 0x7e] (other than `[`,
		// which is the CSI parameter introducer).
		if r == ']' {
			e.inOsc = true
			e.inEsc = false
			return false
		}
		if r >= 0x40 && r <= 0x7e && r != '[' {
			e.inEsc = false
		}
		return false
	}
	if r == '\x1b' {
		e.inEsc = true
		return false
	}
	return true
}

// truncate returns s truncated to at most width visible cells, with an
// ellipsis at the end if anything was removed. SGR escapes and OSC 8
// hyperlink wrappers pass through unmodified — they don't count toward
// width because they don't occupy cells on screen.
//
// A trailing reset (\x1b[0m) is appended only if the input actually
// contained any SGR — so a plain-text caller (e.g. the failure-panel
// body in NO_COLOR mode) doesn't have stray "\x1b[0m" bytes appear in
// its output.
func truncate(s string, width int) string {
	if width <= 0 {
		return ""
	}
	if visibleWidth(s) <= width {
		return s
	}
	var (
		b       strings.Builder
		st      escState
		hadSGR  bool
		visible int
	)
	for _, r := range s {
		isVis := st.step(r)
		if !isVis {
			if r == '\x1b' {
				hadSGR = true
			}
			b.WriteRune(r)
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

// visibleWidth returns the number of visible cells s would occupy,
// skipping SGR escapes and OSC 8 hyperlink sequences. East Asian wide
// characters are not honoured — wake task names are ASCII in practice.
func visibleWidth(s string) int {
	var (
		n  int
		st escState
	)
	for _, r := range s {
		if st.step(r) {
			n++
		}
	}
	return n
}
