package pulse

import (
	"io"
	"strings"

	"github.com/wailsapp/wails/v3/internal/report/pulse/ansi"
)

// liveRegion is a re-paintable bottom-of-terminal pane that sits *below* the
// scrollback. It is the core mechanism that lets the wake TUI keep a pinned
// "active steps + progress" panel visible while completed steps print above it
// and roll off into history.
//
// Model: at any moment we have written some number of pinned lines to the
// terminal. To redraw, we move the cursor up that many lines, clear them, and
// write the new content. To "promote" a completed line into permanent scrollback
// we first erase the pinned region, then write the permanent line, then re-emit
// the pinned region beneath it.
//
// Every paint cycle is wrapped in DEC mode 2026 (synchronized output). Modern
// terminals (Ghostty, Kitty, WezTerm, iTerm2) buffer the in-between writes and
// flip them atomically, eliminating tearing. Terminals that don't recognise
// the mode ignore the sequences — the worst case is identical to today.
type liveRegion struct {
	w     io.Writer
	lines []string // currently-painted lines (final state on screen)
	open  bool     // cursor is currently parked at the bottom of a paint cycle
}

func newLiveRegion(w io.Writer) *liveRegion { return &liveRegion{w: w} }

// paint replaces the pinned region with newLines. Empty newLines clears the
// region entirely. The whole erase-and-rewrite happens inside one DEC 2026
// synchronized-output block, so terminals that support it swap atomically.
func (r *liveRegion) paint(newLines []string) {
	r.update(nil, newLines)
}

// promote prints permanent lines above the pinned region AND swaps the pinned
// region to newPinned, in one synchronized write. Doing both as one atomic
// update prevents the flash where the just-finished step is briefly visible
// above its own "completed" line.
func (r *liveRegion) promote(permanent []string, newPinned []string) {
	r.update(permanent, newPinned)
}

// close erases the pinned region and forgets it. Used at BuildEnd in
// scrollback mode so the next permanent writes don't try to step on a
// region that's no longer there.
func (r *liveRegion) close() {
	r.paint(nil)
	r.lines = nil
	r.open = false
}

// detach forgets the pinned region without erasing it. Used at BuildEnd in
// skeleton mode: the settled skeleton stays exactly where it was drawn and
// flows into scrollback as ordinary lines. Subsequent writes (the build
// summary) flow beneath it without re-painting the area.
func (r *liveRegion) detach() {
	r.lines = nil
	r.open = false
}

// update is the single point through which all changes to the pinned region
// flow: erase the current region (if any), write any permanent lines, write
// the new region. The cycle is wrapped in DEC 2026 sync brackets only when a
// region transition actually happens — a pure promote with no region (e.g.
// non-TTY output) writes plain lines so the sync codes don't pollute logs.
func (r *liveRegion) update(permanent, newLines []string) {
	if !r.open && len(permanent) == 0 && len(newLines) == 0 {
		return
	}
	syncing := r.open || len(newLines) > 0
	var b strings.Builder
	if syncing {
		b.WriteString(ansi.SyncBegin)
	}
	if r.open {
		b.WriteString(ansi.EraseLines(len(r.lines)))
	}
	for _, line := range permanent {
		b.WriteString(line)
		b.WriteRune('\n')
	}
	for i, line := range newLines {
		b.WriteString(line)
		if i < len(newLines)-1 {
			b.WriteRune('\n')
		}
	}
	if syncing {
		b.WriteString(ansi.SyncEnd)
	}
	_, _ = io.WriteString(r.w, b.String())
	r.lines = append(r.lines[:0], newLines...)
	r.open = len(newLines) > 0
}
