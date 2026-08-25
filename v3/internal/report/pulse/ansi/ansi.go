// Package ansi is a tiny pure-stdlib set of ANSI helpers, just enough to render
// the wake build TUI without pulling in a styling library.
//
// Three layers:
//   - SGR colour/attribute codes (Sgr, Reset, plus shortcuts)
//   - Cursor control (Up, Down, ClearLine, ClearBelow, Save, Restore)
//   - DEC private modes (HideCursor, ShowCursor, SyncBegin, SyncEnd)
//
// Nothing here speaks lipgloss's box-model; the wake reporter does its own
// padding/truncation because that work fits in a screenful of code and the
// reporter has very few unique layouts.
package ansi

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// CSI is the Control Sequence Introducer prefix.
const CSI = "\x1b["

// Common one-shots. Pre-formatted so renderers can write them as bytes.
const (
	Reset      = CSI + "0m"
	ClearLine  = CSI + "2K"
	ClearBelow = CSI + "0J"
	HideCursor = CSI + "?25l"
	ShowCursor = CSI + "?25h"
	// SyncBegin/SyncEnd wrap a paint cycle so tearing-free terminals
	// (Ghostty, Kitty, WezTerm, iTerm2) atomically swap the affected region.
	// Terminals that don't recognise mode 2026 ignore both sequences.
	SyncBegin = CSI + "?2026h"
	SyncEnd   = CSI + "?2026l"
)

// Up emits the "cursor up N rows, column 1" sequence. CUP at column 0 is what
// we always want before redrawing the pinned region.
func Up(n int) string {
	if n <= 0 {
		return "\r"
	}
	return fmt.Sprintf("%s%dF", CSI, n)
}

// Down emits the "cursor down N rows" sequence.
func Down(n int) string {
	if n <= 0 {
		return ""
	}
	return fmt.Sprintf("%s%dB", CSI, n)
}

// EraseLines erases the current line plus n-1 lines above it (i.e. n lines in
// total), and leaves the cursor at column 1 of the topmost erased line. Use
// when the caller knows it just wrote n lines without a trailing newline and
// the cursor is therefore parked at the end of the last one.
//
// n <= 0 is a no-op.
func EraseLines(n int) string {
	if n <= 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("\r")
	b.WriteString(ClearLine)
	for i := 1; i < n; i++ {
		b.WriteString(CSI + "1A")
		b.WriteString(ClearLine)
	}
	return b.String()
}

// Style wraps s in the given SGR sequence, with a reset at the end. The empty
// sgr (no styling requested) returns s unchanged so callers can compose without
// a special case.
func Style(sgr, s string) string {
	if sgr == "" {
		return s
	}
	return sgr + s + Reset
}

// Sgr builds a CSI ... m sequence from its component codes. Use the helpers
// below in preference; Sgr is the escape hatch for one-off combinations.
func Sgr(codes ...string) string {
	if len(codes) == 0 {
		return ""
	}
	return CSI + strings.Join(codes, ";") + "m"
}

// SGR attribute codes.
const (
	AttrBold  = "1"
	AttrFaint = "2"
)

// Bold returns the SGR for bold text.
func Bold() string { return CSI + AttrBold + "m" }

// Faint returns the SGR for "faint" (dim) text. Most terminals render this as
// reduced contrast against the background.
func Faint() string { return CSI + AttrFaint + "m" }

// IsTerminal reports whether f is a terminal. A nil file is reported as false.
func IsTerminal(f *os.File) bool {
	if f == nil {
		return false
	}
	return term.IsTerminal(int(f.Fd()))
}

// TermSize returns the terminal width and height for f, or (80, 24) if f isn't
// a terminal or the size can't be read.
func TermSize(f *os.File) (cols, rows int) {
	if f == nil {
		return 80, 24
	}
	w, h, err := term.GetSize(int(f.Fd()))
	if err != nil || w <= 0 {
		return 80, 24
	}
	return w, h
}
