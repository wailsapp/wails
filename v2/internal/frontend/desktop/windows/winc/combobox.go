//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 */

package winc

import (
	"syscall"
	"unsafe"

	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
)

type ComboBox struct {
	ControlBase
	onSelectedChange EventManager
}

func NewComboBox(parent Controller) *ComboBox {
	cb := new(ComboBox)

	cb.InitControl("COMBOBOX", parent, 0, w32.WS_CHILD|w32.WS_VISIBLE|w32.WS_TABSTOP|w32.WS_VSCROLL|w32.CBS_DROPDOWNLIST)
	RegMsgHandler(cb)

	cb.SetFont(DefaultFont)
	cb.SetSize(200, 400)
	return cb
}

func (cb *ComboBox) DeleteAllItems() bool {
	return w32.SendMessage(cb.hwnd, w32.CB_RESETCONTENT, 0, 0) == w32.TRUE
}

func (cb *ComboBox) InsertItem(index int, str string) bool {
	lp := uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(str)))
	return w32.SendMessage(cb.hwnd, w32.CB_INSERTSTRING, uintptr(index), lp) != w32.CB_ERR
}

func (cb *ComboBox) DeleteItem(index int) bool {
	return w32.SendMessage(cb.hwnd, w32.CB_DELETESTRING, uintptr(index), 0) != w32.CB_ERR
}

func (cb *ComboBox) SelectedItem() int {
	return int(int32(w32.SendMessage(cb.hwnd, w32.CB_GETCURSEL, 0, 0)))
}

func (cb *ComboBox) SetSelectedItem(value int) bool {
	return int(int32(w32.SendMessage(cb.hwnd, w32.CB_SETCURSEL, uintptr(value), 0))) == value
}

func (cb *ComboBox) OnSelectedChange() *EventManager {
	return &cb.onSelectedChange
}

// Message processer
func (cb *ComboBox) WndProc(msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case w32.WM_COMMAND:
		code := w32.HIWORD(uint32(wparam))

		switch code {
		case w32.CBN_SELCHANGE:
			cb.onSelectedChange.Fire(NewEvent(cb, nil))
		}
	}
	return w32.DefWindowProc(cb.hwnd, msg, wparam, lparam)
	//return cb.W32Control.WndProc(msg, wparam, lparam)
}
