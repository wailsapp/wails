package pulse

import (
	"fmt"
	"strings"
)

// Block-element glyphs we use for visuals. Single-cell, monospace-safe,
// supported by every font shipped with a terminal emulator. We deliberately
// avoid emoji (variable width, inconsistent fonts) and the heavier box-drawing
// patterns; the wake TUI's identity is "calm Unicode, not maximalist".
const (
	GlyphOK      = "✓"
	GlyphFail    = "✗"
	GlyphCached  = "⊙"
	GlyphSkipped = "·"
	GlyphArrow   = "›"
	GlyphBullet  = "•"
	GlyphActive  = "◆"
	GlyphQueued  = "◇"
	GlyphCrit    = "◀"
	GlyphRule    = "━" // header signature rule

	// Progress bar: thick horizontal lines, with a transition glyph between
	// filled and empty so the bar reads as one continuous stroke rather than
	// two abutting blocks. The leading "╾" gives the head a soft point.
	BarFilled = "━"
	BarHead   = "╾"
	BarEmpty  = "─"

	// Sparkline scale: 8 evenly-stepped block heights.
	SparkRunes = "▁▂▃▄▅▆▇█"

	// Sub-cell horizontal block fill (eighths): each glyph is one cell wide
	// but increases left-anchored fill in 1/8 increments. Used by the
	// fine-grained progress bar so the head advances smoothly even when the
	// overall ratio changes by tiny amounts.
	HEighths = " ▏▎▍▌▋▊▉█"
)

// Spinner — Wake's signature "heartbeat": a vertical block that swells from
// a low bar to full height and back down. This is what earns the renderer
// its name. 14 frames at 80 ms = ~1.1 s per beat, roughly a resting human
// heart rate. All in-flight steps share a single frame counter so the build
// pulses in unison, like a chorus rather than a crowd.
//
// We deliberately don't use the more common braille spinner ("⠋⠙⠹⠸"…) —
// those signify "loading" generically; this signifies *wake*.
var spinnerFrames = []string{
	"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█", "▇", "▆", "▅", "▄", "▃", "▂",
}

// progressBar renders a single-line progress bar of width visible cells.
// ratio is clamped to [0, 1]. The bar reads "━━━━━╾─────  47%  ETA 5s".
//
// Sub-cell precision: the filled region uses eighths-block fill so the head
// advances in 1/8-cell steps as ratio changes between integer cell crossings.
// This makes the progress feel "alive" even on slow builds where a normal
// block-resolution bar would sit unchanged for seconds at a time.
//
// The leading edge gets BarHead drawn in accent, while everything to the
// left is solid accent-coloured BarFilled — the contrast pulls the eye to
// the head, which is where the progress information lives.
func (s *styler) progressBar(ratio float64, width int) string {
	if width < 4 {
		width = 4
	}
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	// Compute fill in eighths of a cell. 8 sub-cells per cell.
	totalEighths := int(float64(width*8) * ratio)
	fullCells := totalEighths / 8
	remainder := totalEighths % 8
	if fullCells > width {
		fullCells = width
		remainder = 0
	}
	var b strings.Builder
	cells := fullCells
	switch {
	case remainder > 0 && cells < width:
		// Full cells, then a partial eighths-block head, then dim track.
		if cells > 0 {
			b.WriteString(s.fg(Accent, strings.Repeat(BarFilled, cells)))
		}
		runes := []rune(HEighths)
		b.WriteString(s.fg(Accent, string(runes[remainder])))
		cells++
	case cells > 0:
		// Cell-aligned head: use BarHead as the soft point at the leading edge.
		b.WriteString(s.fg(Accent, strings.Repeat(BarFilled, cells-1)+BarHead))
	}
	if cells < width {
		b.WriteString(s.fg(Dim, strings.Repeat(BarEmpty, width-cells)))
	}
	return b.String()
}

// progressLine returns the full progress strip: bar + percent + ETA, padded.
// elapsed and eta are formatted by the caller; pass "" for unknown ETA.
func (s *styler) progressLine(ratio float64, width int, elapsed, eta string) string {
	// Reserve the right-side metrics column. "  100%  ETA 99m99s" worst case = 18.
	const metricsW = 22
	barW := width - metricsW
	if barW < 8 {
		barW = 8
	}
	pct := fmt.Sprintf("%3d%%", int(ratio*100+0.5))
	right := s.bold(pct)
	if eta != "" {
		right += "  " + s.faint("ETA "+eta)
	} else if elapsed != "" {
		right += "  " + s.faint(elapsed)
	}
	return s.progressBar(ratio, barW) + "  " + right
}

// sparkline renders values as Unicode block heights. The range is computed
// against the values themselves: small absolute differences still produce
// useful peaks-and-troughs, which is what we want for "is the build moving?"
// rather than absolute throughput numbers.
func (s *styler) sparkline(values []float64) string {
	if len(values) == 0 {
		return ""
	}
	maxV := values[0]
	for _, v := range values {
		if v > maxV {
			maxV = v
		}
	}
	if maxV <= 0 {
		return strings.Repeat(string([]rune(SparkRunes)[0]), len(values))
	}
	runes := []rune(SparkRunes)
	var b strings.Builder
	for _, v := range values {
		if v < 0 {
			v = 0
		}
		idx := int(v / maxV * float64(len(runes)-1))
		if idx >= len(runes) {
			idx = len(runes) - 1
		}
		b.WriteRune(runes[idx])
	}
	return s.fg(Accent, b.String())
}

// barchart renders one horizontal bar (used in the end-summary "slowest
// steps" table). width is the total cell budget; the bar is sized to
// "value / max * width" and rendered with the same accent fill as the
// progress bar so the visual idiom is consistent.
//
// rank is the row's position in the sorted list (0 = slowest). The slowest
// row uses bold accent — it is *the* bottleneck and deserves the visual
// weight; subsequent rows fall back to plain accent, then dim on the empty
// track. This creates a vertical hierarchy without needing per-row colour.
func (s *styler) barchart(value, max float64, width, rank int) string {
	if max <= 0 || width <= 0 {
		return ""
	}
	ratio := value / max
	if ratio > 1 {
		ratio = 1
	}
	n := int(float64(width)*ratio + 0.5)
	if n < 1 && value > 0 {
		n = 1
	}
	if n > width {
		n = width
	}
	fill := strings.Repeat(BarFilled, n)
	if rank == 0 {
		fill = s.bold(s.fg(Accent, fill))
	} else {
		fill = s.fg(Accent, fill)
	}
	return fill + s.fg(Dim, strings.Repeat(BarEmpty, width-n))
}
