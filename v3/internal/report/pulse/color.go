// Package pulse renders a wake build to a terminal as the "Pulse" TUI: a
// pre-painted skeleton of one row per planned step that updates each row
// in place as steps progress, plus a pinned progress strip and throughput
// sparkline beneath. Written without third-party styling libraries.
package pulse

import (
	"fmt"
	"os"
	"strings"

	"github.com/wailsapp/wails/v3/internal/report/pulse/ansi"
)

// Profile is what the terminal can render. We fall through TrueColor → 256 →
// 16-colour ANSI → no colour at all, picking the richest profile the terminal
// advertises. Terminals routinely lie about this (tmux historically claims
// less than it can render; some CI systems claim more) but TERM/COLORTERM/
// NO_COLOR are the conventions, so we follow them.
type Profile int

const (
	ProfileNone     Profile = iota // monochrome
	ProfileANSI                    // 16 colours
	ProfileANSI256                 // 256 colours
	ProfileTrueCol                 // 24-bit
)

// DetectProfile returns the colour profile the receiving terminal advertises.
// f is the output stream; if it isn't a TTY the profile is forced to None.
// NO_COLOR (any value) forces None regardless.
func DetectProfile(f *os.File) Profile {
	if !ansi.IsTerminal(f) {
		return ProfileNone
	}
	if os.Getenv("NO_COLOR") != "" {
		return ProfileNone
	}
	colorTerm := strings.ToLower(os.Getenv("COLORTERM"))
	if colorTerm == "truecolor" || colorTerm == "24bit" {
		return ProfileTrueCol
	}
	termVar := os.Getenv("TERM")
	switch {
	case strings.Contains(termVar, "256color"):
		return ProfileANSI256
	case termVar == "dumb" || termVar == "":
		return ProfileNone
	case strings.Contains(termVar, "kitty"), strings.Contains(termVar, "ghostty"),
		strings.Contains(termVar, "alacritty"), strings.Contains(termVar, "wezterm"):
		return ProfileTrueCol
	}
	return ProfileANSI
}

// Colour describes one logical UI colour at three render tiers. The reporter
// picks the highest tier the active profile supports. Callers should not build
// Colour values directly — use the predefined palette members below.
type Colour struct {
	True string // 24-bit hex without leading '#', e.g. "7d9eff"
	X256 uint8  // 256-colour palette index
	Ansi uint8  // 30–37 / 90–97 base
}

// FG returns the SGR foreground escape for c at profile p, or "" if p is None.
func (c Colour) FG(p Profile) string {
	switch p {
	case ProfileNone:
		return ""
	case ProfileANSI:
		return fmt.Sprintf("%s%dm", ansi.CSI, c.Ansi)
	case ProfileANSI256:
		return fmt.Sprintf("%s38;5;%dm", ansi.CSI, c.X256)
	case ProfileTrueCol:
		r, g, b := hexRGB(c.True)
		return fmt.Sprintf("%s38;2;%d;%d;%dm", ansi.CSI, r, g, b)
	}
	return ""
}

// BG returns the SGR background escape for c at profile p.
func (c Colour) BG(p Profile) string {
	switch p {
	case ProfileNone:
		return ""
	case ProfileANSI:
		return fmt.Sprintf("%s%dm", ansi.CSI, c.Ansi+10)
	case ProfileANSI256:
		return fmt.Sprintf("%s48;5;%dm", ansi.CSI, c.X256)
	case ProfileTrueCol:
		r, g, b := hexRGB(c.True)
		return fmt.Sprintf("%s48;2;%d;%d;%dm", ansi.CSI, r, g, b)
	}
	return ""
}

// Palette — the wake brand. Picked for legibility on both dark and light
// terminals at all three tiers. The Accent is the only brand-saturated colour;
// everything else stays in muted territory so the UI reads as calm.
var (
	// Accent — electric blue, used for the verb, the spinner, and the
	// progress-bar fill. The TUI's only deliberately attention-grabbing colour.
	Accent = Colour{True: "7d9eff", X256: 111, Ansi: 94}
	// Success — muted green, used for ✓ glyph and the success summary.
	Success = Colour{True: "7ec682", X256: 114, Ansi: 92}
	// Failure — desaturated red, used for ✗ and panel borders. Bold-adjacent
	// without being alarming; we let typography (bold) carry the urgency.
	Failure = Colour{True: "e06c75", X256: 167, Ansi: 91}
	// Warning — straw yellow, used for partial-success notes.
	Warning = Colour{True: "e5c07b", X256: 179, Ansi: 93}
	// Cached — soft cyan, used for cache hits to differentiate them from
	// "ran but quick" successes.
	Cached = Colour{True: "8ec5d8", X256: 109, Ansi: 96}
	// Dim — neutral grey for chrome (counters, separators, "→").
	Dim = Colour{True: "6e7785", X256: 244, Ansi: 90}
	// Subtle — slightly brighter than Dim; for the second-most-important text
	// in a row (a "1.2s" beside a step name).
	Subtle = Colour{True: "9aa5b8", X256: 248, Ansi: 37}
)

func hexRGB(h string) (r, g, b uint8) {
	if len(h) != 6 {
		return 0, 0, 0
	}
	var v uint32
	for i := 0; i < 6; i++ {
		c := h[i]
		var d uint32
		switch {
		case c >= '0' && c <= '9':
			d = uint32(c - '0')
		case c >= 'a' && c <= 'f':
			d = uint32(c-'a') + 10
		case c >= 'A' && c <= 'F':
			d = uint32(c-'A') + 10
		}
		v = v<<4 | d
	}
	return uint8(v >> 16), uint8(v >> 8), uint8(v)
}
