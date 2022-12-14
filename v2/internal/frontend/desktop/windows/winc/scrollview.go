//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 */

package winc

import (
	"unsafe"

	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
)

type ScrollView struct {
	ControlBase
	child Dockable
}

func NewScrollView(parent Controller) *ScrollView {
	sv := new(ScrollView)

	RegClassOnlyOnce("winc_ScrollView")
	sv.hwnd = CreateWindow("winc_ScrollView", parent, w32.WS_EX_CONTROLPARENT,
		w32.WS_CHILD|w32.WS_HSCROLL|w32.WS_VISIBLE|w32.WS_VSCROLL)
	sv.parent = parent
	RegMsgHandler(sv)

	sv.SetFont(DefaultFont)
	sv.SetText("")
	sv.SetSize(200, 200)
	return sv
}

func (sv *ScrollView) SetChild(child Dockable) {
	sv.child = child
}

func (sv *ScrollView) UpdateScrollBars() {
	w, h := sv.child.Width(), sv.child.Height()
	sw, sh := sv.Size()

	var si w32.SCROLLINFO
	si.CbSize = uint32(unsafe.Sizeof(si))
	si.FMask = w32.SIF_PAGE | w32.SIF_RANGE

	si.NMax = int32(w - 1)
	si.NPage = uint32(sw)
	w32.SetScrollInfo(sv.hwnd, w32.SB_HORZ, &si, true)
	x := sv.scroll(w32.SB_HORZ, w32.SB_THUMBPOSITION)

	si.NMax = int32(h)
	si.NPage = uint32(sh)
	w32.SetScrollInfo(sv.hwnd, w32.SB_VERT, &si, true)
	y := sv.scroll(w32.SB_VERT, w32.SB_THUMBPOSITION)

	sv.child.SetPos(x, y)
}

func (sv *ScrollView) scroll(sb int32, cmd uint16) int {
	var pos int32
	var si w32.SCROLLINFO
	si.CbSize = uint32(unsafe.Sizeof(si))
	si.FMask = w32.SIF_PAGE | w32.SIF_POS | w32.SIF_RANGE | w32.SIF_TRACKPOS

	w32.GetScrollInfo(sv.hwnd, sb, &si)
	pos = si.NPos

	switch cmd {
	case w32.SB_LINELEFT: // == win.SB_LINEUP
		pos -= 20

	case w32.SB_LINERIGHT: // == win.SB_LINEDOWN
		pos += 20

	case w32.SB_PAGELEFT: // == win.SB_PAGEUP
		pos -= int32(si.NPage)

	case w32.SB_PAGERIGHT: // == win.SB_PAGEDOWN
		pos += int32(si.NPage)

	case w32.SB_THUMBTRACK:
		pos = si.NTrackPos
	}

	if pos < 0 {
		pos = 0
	}
	if pos > si.NMax+1-int32(si.NPage) {
		pos = si.NMax + 1 - int32(si.NPage)
	}

	si.FMask = w32.SIF_POS
	si.NPos = pos
	w32.SetScrollInfo(sv.hwnd, sb, &si, true)

	return -int(pos)
}

func (sv *ScrollView) WndProc(msg uint32, wparam, lparam uintptr) uintptr {
	if sv.child != nil {
		switch msg {
		case w32.WM_PAINT:
			sv.UpdateScrollBars()

		case w32.WM_HSCROLL:
			x, y := sv.child.Pos()
			x = sv.scroll(w32.SB_HORZ, w32.LOWORD(uint32(wparam)))
			sv.child.SetPos(x, y)

		case w32.WM_VSCROLL:
			x, y := sv.child.Pos()
			y = sv.scroll(w32.SB_VERT, w32.LOWORD(uint32(wparam)))
			sv.child.SetPos(x, y)

		case w32.WM_SIZE, w32.WM_SIZING:
			sv.UpdateScrollBars()
		}
	}
	return w32.DefWindowProc(sv.hwnd, msg, wparam, lparam)
}
