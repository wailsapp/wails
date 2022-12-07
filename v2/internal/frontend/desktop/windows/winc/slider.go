//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 */

package winc

import "github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"

type Slider struct {
	ControlBase
	prevPos int

	onScroll EventManager
}

func NewSlider(parent Controller) *Slider {
	tb := new(Slider)

	tb.InitControl("msctls_trackbar32", parent, 0, w32.WS_TABSTOP|w32.WS_VISIBLE|w32.WS_CHILD /*|w32.TBS_AUTOTICKS*/)
	RegMsgHandler(tb)

	tb.SetFont(DefaultFont)
	tb.SetText("Slider")
	tb.SetSize(200, 32)

	tb.SetRange(0, 100)
	tb.SetPage(10)
	return tb
}

func (tb *Slider) OnScroll() *EventManager {
	return &tb.onScroll
}

func (tb *Slider) Value() int {
	ret := w32.SendMessage(tb.hwnd, w32.TBM_GETPOS, 0, 0)
	return int(ret)
}

func (tb *Slider) SetValue(v int) {
	tb.prevPos = v
	w32.SendMessage(tb.hwnd, w32.TBM_SETPOS, uintptr(w32.BoolToBOOL(true)), uintptr(v))
}

func (tb *Slider) Range() (min, max int) {
	min = int(w32.SendMessage(tb.hwnd, w32.TBM_GETRANGEMIN, 0, 0))
	max = int(w32.SendMessage(tb.hwnd, w32.TBM_GETRANGEMAX, 0, 0))
	return min, max
}

func (tb *Slider) SetRange(min, max int) {
	w32.SendMessage(tb.hwnd, w32.TBM_SETRANGE, uintptr(w32.BoolToBOOL(true)), uintptr(w32.MAKELONG(uint16(min), uint16(max))))
}

func (tb *Slider) SetPage(pagesize int) {
	w32.SendMessage(tb.hwnd, w32.TBM_SETPAGESIZE, 0, uintptr(pagesize))
}

func (tb *Slider) WndProc(msg uint32, wparam, lparam uintptr) uintptr {
	/*
		// REMOVE:
		// following code did not work, used workaround below
		code := w32.LOWORD(uint32(wparam))

		switch code {
		case w32.TB_ENDTRACK:
			tb.onScroll.Fire(NewEvent(tb, nil))
		}*/

	newPos := tb.Value()
	if newPos != tb.prevPos {
		tb.onScroll.Fire(NewEvent(tb, nil))
		tb.prevPos = newPos
	}

	return w32.DefWindowProc(tb.hwnd, msg, wparam, lparam)
}
