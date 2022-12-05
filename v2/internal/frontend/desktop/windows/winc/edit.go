//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package winc

import "github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"

type Edit struct {
	ControlBase
	onChange EventManager
}

const passwordChar = '*'
const nopasswordChar = ' '

func NewEdit(parent Controller) *Edit {
	edt := new(Edit)

	edt.InitControl("EDIT", parent, w32.WS_EX_CLIENTEDGE, w32.WS_CHILD|w32.WS_VISIBLE|w32.WS_TABSTOP|w32.ES_LEFT|
		w32.ES_AUTOHSCROLL)
	RegMsgHandler(edt)

	edt.SetFont(DefaultFont)
	edt.SetSize(200, 22)
	return edt
}

// Events.
func (ed *Edit) OnChange() *EventManager {
	return &ed.onChange
}

// Public methods.
func (ed *Edit) SetReadOnly(isReadOnly bool) {
	w32.SendMessage(ed.hwnd, w32.EM_SETREADONLY, uintptr(w32.BoolToBOOL(isReadOnly)), 0)
}

// Public methods
func (ed *Edit) SetPassword(isPassword bool) {
	if isPassword {
		w32.SendMessage(ed.hwnd, w32.EM_SETPASSWORDCHAR, uintptr(passwordChar), 0)
	} else {
		w32.SendMessage(ed.hwnd, w32.EM_SETPASSWORDCHAR, 0, 0)
	}
}

func (ed *Edit) WndProc(msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case w32.WM_COMMAND:
		switch w32.HIWORD(uint32(wparam)) {
		case w32.EN_CHANGE:
			ed.onChange.Fire(NewEvent(ed, nil))
		}
		/*case w32.WM_GETDLGCODE:
		println("Edit")
		if wparam == w32.VK_RETURN {
			return w32.DLGC_WANTALLKEYS
		}*/
	}
	return w32.DefWindowProc(ed.hwnd, msg, wparam, lparam)
}

// MultiEdit is multiline text edit.
type MultiEdit struct {
	ControlBase
	onChange EventManager
}

func NewMultiEdit(parent Controller) *MultiEdit {
	med := new(MultiEdit)

	med.InitControl("EDIT", parent, w32.WS_EX_CLIENTEDGE, w32.WS_CHILD|w32.WS_VISIBLE|w32.WS_TABSTOP|w32.ES_LEFT|
		w32.WS_VSCROLL|w32.WS_HSCROLL|w32.ES_MULTILINE|w32.ES_WANTRETURN|w32.ES_AUTOHSCROLL|w32.ES_AUTOVSCROLL)
	RegMsgHandler(med)

	med.SetFont(DefaultFont)
	med.SetSize(200, 400)
	return med
}

// Events
func (med *MultiEdit) OnChange() *EventManager {
	return &med.onChange
}

// Public methods
func (med *MultiEdit) SetReadOnly(isReadOnly bool) {
	w32.SendMessage(med.hwnd, w32.EM_SETREADONLY, uintptr(w32.BoolToBOOL(isReadOnly)), 0)
}

func (med *MultiEdit) AddLine(text string) {
	if len(med.Text()) == 0 {
		med.SetText(text)
	} else {
		med.SetText(med.Text() + "\r\n" + text)
	}
}

func (med *MultiEdit) WndProc(msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {

	case w32.WM_COMMAND:
		switch w32.HIWORD(uint32(wparam)) {
		case w32.EN_CHANGE:
			med.onChange.Fire(NewEvent(med, nil))
		}
	}
	return w32.DefWindowProc(med.hwnd, msg, wparam, lparam)
}
