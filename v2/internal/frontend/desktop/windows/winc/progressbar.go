//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 */

package winc

import (
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
)

type ProgressBar struct {
	ControlBase
}

func NewProgressBar(parent Controller) *ProgressBar {
	pb := new(ProgressBar)

	pb.InitControl(w32.PROGRESS_CLASS, parent, 0, w32.WS_CHILD|w32.WS_VISIBLE)
	RegMsgHandler(pb)

	pb.SetSize(200, 22)
	return pb
}

func (pr *ProgressBar) Value() int {
	ret := w32.SendMessage(pr.hwnd, w32.PBM_GETPOS, 0, 0)
	return int(ret)
}

func (pr *ProgressBar) SetValue(v int) {
	w32.SendMessage(pr.hwnd, w32.PBM_SETPOS, uintptr(v), 0)
}

func (pr *ProgressBar) Range() (min, max uint) {
	min = uint(w32.SendMessage(pr.hwnd, w32.PBM_GETRANGE, uintptr(w32.BoolToBOOL(true)), 0))
	max = uint(w32.SendMessage(pr.hwnd, w32.PBM_GETRANGE, uintptr(w32.BoolToBOOL(false)), 0))
	return
}

func (pr *ProgressBar) SetRange(min, max int) {
	w32.SendMessage(pr.hwnd, w32.PBM_SETRANGE32, uintptr(min), uintptr(max))
}

func (pr *ProgressBar) WndProc(msg uint32, wparam, lparam uintptr) uintptr {
	return w32.DefWindowProc(pr.hwnd, msg, wparam, lparam)
}
