package tui

import (
	"github.com/gdamore/tcell/v2"

	"github.com/atterpac/dado/components"
)

// vimTableNav bridges vim movement (j/k/g/G) to a dado Table's row selection.
//
// dado's core Table only navigates with the arrow keys and never interprets
// j/k runes, so without this the cursor sits on the header row and vim keys do
// nothing. It returns true when it consumed a nav key (so the caller stops),
// false otherwise (arrow keys fall through to the table's native handling).
//
// Selecting a row calls Table.Select, which fires the table's
// SelectionChangedFunc, so any detail pane wired to selection updates for free.
func vimTableNav(t *components.Table, ev *tcell.EventKey) bool {
	if ev.Key() != tcell.KeyRune {
		return false
	}
	n := t.GetDataRowCount()
	if n == 0 {
		return false
	}
	target := t.SelectedRow()
	switch ev.Rune() {
	case 'j':
		target++
	case 'k':
		target--
	case 'g':
		target = 0
	case 'G':
		target = n - 1
	default:
		return false
	}
	if target < 0 {
		target = 0
	}
	if target >= n {
		target = n - 1
	}
	t.SelectRow(target)
	return true
}
