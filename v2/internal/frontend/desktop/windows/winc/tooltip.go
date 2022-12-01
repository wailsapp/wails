//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package winc

import (
	"syscall"
	"unsafe"

	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
)

type ToolTip struct {
	ControlBase
}

func NewToolTip(parent Controller) *ToolTip {
	tp := new(ToolTip)

	tp.InitControl("tooltips_class32", parent, w32.WS_EX_TOPMOST, w32.WS_POPUP|w32.TTS_NOPREFIX|w32.TTS_ALWAYSTIP)
	w32.SetWindowPos(tp.Handle(), w32.HWND_TOPMOST, 0, 0, 0, 0, w32.SWP_NOMOVE|w32.SWP_NOSIZE|w32.SWP_NOACTIVATE)

	return tp
}

func (tp *ToolTip) SetTip(tool Controller, tip string) bool {
	var ti w32.TOOLINFO
	ti.CbSize = uint32(unsafe.Sizeof(ti))
	if tool.Parent() != nil {
		ti.Hwnd = tool.Parent().Handle()
	}
	ti.UFlags = w32.TTF_IDISHWND | w32.TTF_SUBCLASS /* | TTF_ABSOLUTE */
	ti.UId = uintptr(tool.Handle())
	ti.LpszText = syscall.StringToUTF16Ptr(tip)

	return w32.SendMessage(tp.Handle(), w32.TTM_ADDTOOL, 0, uintptr(unsafe.Pointer(&ti))) != w32.FALSE
}

func (tp *ToolTip) WndProc(msg uint, wparam, lparam uintptr) uintptr {
	return w32.DefWindowProc(tp.hwnd, uint32(msg), wparam, lparam)
}
