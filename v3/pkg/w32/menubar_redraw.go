package w32

// sizeRedrawGate decides when the dark menubar needs a forced non-client
// redraw after a size change (#4462).
//
// Snapping a window to a screen edge (Win+Left/Right, or dragging to an edge)
// delivers WM_SIZE with SIZE_RESTORED when the window was already in the
// normal state, so the previous SIZE_RESTORED -> SIZE_RESTORED "no state
// change" heuristic never fires for it. Instead:
//
//   - WM_SIZE outside an interactive move/size loop (keyboard snap, caption
//     buttons, programmatic changes) redraws once.
//   - WM_SIZE inside the loop (WM_ENTERSIZEMOVE..WM_EXITSIZEMOVE) is
//     suppressed, because live drag-resizing sends one per frame and
//     redrawing each would flicker.
//   - WM_EXITSIZEMOVE redraws once, covering the drag-to-edge snap, whose
//     final size is applied while the loop is still active.
//
// Only touched from the window's UI thread (window procedure), so it needs no
// locking.
type sizeRedrawGate struct {
	inSizeMove bool
}

// enterSizeMove records WM_ENTERSIZEMOVE.
func (g *sizeRedrawGate) enterSizeMove() { g.inSizeMove = true }

// exitSizeMove records WM_EXITSIZEMOVE and reports whether to redraw.
func (g *sizeRedrawGate) exitSizeMove() bool {
	was := g.inSizeMove
	g.inSizeMove = false
	return was
}

// onSize reports whether a WM_SIZE should force a redraw. eligible is true
// for SIZE_RESTORED and SIZE_MAXIMIZED (not SIZE_MINIMIZED).
func (g *sizeRedrawGate) onSize(eligible bool) bool {
	return eligible && !g.inSizeMove
}
