package w32

import "testing"

func TestSizeRedrawGate(t *testing.T) {
	var g sizeRedrawGate

	// Keyboard snap: no size-move loop, one WM_SIZE -> redraw.
	if !g.onSize(true) {
		t.Fatal("WM_SIZE outside a size-move loop must redraw")
	}
	// Minimised is never eligible.
	if g.onSize(false) {
		t.Fatal("ineligible WM_SIZE must not redraw")
	}

	// Live drag: per-frame WM_SIZE suppressed, redraw once at exit.
	g.enterSizeMove()
	for i := 0; i < 5; i++ {
		if g.onSize(true) {
			t.Fatal("WM_SIZE during size-move loop must not redraw")
		}
	}
	if !g.exitSizeMove() {
		t.Fatal("exit of a size-move loop must redraw once")
	}
	if g.exitSizeMove() {
		t.Fatal("spurious exit without enter must not redraw")
	}
	if !g.onSize(true) {
		t.Fatal("gate must reopen after exit")
	}
}
