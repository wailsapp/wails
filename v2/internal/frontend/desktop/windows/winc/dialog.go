//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package winc

import "github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"

// Dialog displayed as z-order top window until closed.
// It also disables parent window so it can not be clicked.
type Dialog struct {
	Form
	isModal bool

	btnOk     *PushButton
	btnCancel *PushButton

	onLoad   EventManager
	onOk     EventManager
	onCancel EventManager
}

func NewDialog(parent Controller) *Dialog {
	dlg := new(Dialog)

	dlg.isForm = true
	dlg.isModal = true
	RegClassOnlyOnce("winc_Dialog")

	dlg.hwnd = CreateWindow("winc_Dialog", parent, w32.WS_EX_CONTROLPARENT, /* IMPORTANT */
		w32.WS_SYSMENU|w32.WS_CAPTION|w32.WS_THICKFRAME /*|w32.WS_BORDER|w32.WS_POPUP*/)
	dlg.parent = parent

	// dlg might fail if icon resource is not embedded in the binary
	if ico, err := NewIconFromResource(GetAppInstance(), uint16(AppIconID)); err == nil {
		dlg.SetIcon(0, ico)
	}

	// Dlg forces display of focus rectangles, as soon as the user starts to type.
	w32.SendMessage(dlg.hwnd, w32.WM_CHANGEUISTATE, w32.UIS_INITIALIZE, 0)
	RegMsgHandler(dlg)

	dlg.SetFont(DefaultFont)
	dlg.SetText("Form")
	dlg.SetSize(200, 100)
	return dlg
}

func (dlg *Dialog) SetModal(modal bool) {
	dlg.isModal = modal
}

// SetButtons wires up dialog events to buttons. btnCancel can be nil.
func (dlg *Dialog) SetButtons(btnOk *PushButton, btnCancel *PushButton) {
	dlg.btnOk = btnOk
	dlg.btnOk.SetDefault()
	dlg.btnCancel = btnCancel
}

// Events
func (dlg *Dialog) OnLoad() *EventManager {
	return &dlg.onLoad
}

func (dlg *Dialog) OnOk() *EventManager {
	return &dlg.onOk
}

func (dlg *Dialog) OnCancel() *EventManager {
	return &dlg.onCancel
}

// PreTranslateMessage handles dialog specific messages. IMPORTANT.
func (dlg *Dialog) PreTranslateMessage(msg *w32.MSG) bool {
	if msg.Message >= w32.WM_KEYFIRST && msg.Message <= w32.WM_KEYLAST {
		if w32.IsDialogMessage(dlg.hwnd, msg) {
			return true
		}
	}
	return false
}

// Show dialog performs special setup for dialog windows.
func (dlg *Dialog) Show() {
	if dlg.isModal {
		dlg.Parent().SetEnabled(false)
	}
	dlg.onLoad.Fire(NewEvent(dlg, nil))
	dlg.Form.Show()
}

// Close dialog when you done with it.
func (dlg *Dialog) Close() {
	if dlg.isModal {
		dlg.Parent().SetEnabled(true)
	}
	dlg.ControlBase.Close()
}

func (dlg *Dialog) cancel() {
	if dlg.btnCancel != nil {
		dlg.btnCancel.onClick.Fire(NewEvent(dlg.btnCancel, nil))
	}
	dlg.onCancel.Fire(NewEvent(dlg, nil))
}

func (dlg *Dialog) WndProc(msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case w32.WM_COMMAND:
		switch w32.LOWORD(uint32(wparam)) {
		case w32.IDOK:
			if dlg.btnOk != nil {
				dlg.btnOk.onClick.Fire(NewEvent(dlg.btnOk, nil))
			}
			dlg.onOk.Fire(NewEvent(dlg, nil))
			return w32.TRUE

		case w32.IDCANCEL:
			dlg.cancel()
			return w32.TRUE
		}

	case w32.WM_CLOSE:
		dlg.cancel() // use onCancel or dlg.btnCancel.OnClick to close
		return 0

	case w32.WM_DESTROY:
		if dlg.isModal {
			dlg.Parent().SetEnabled(true)
		}
	}
	return w32.DefWindowProc(dlg.hwnd, msg, wparam, lparam)
}
