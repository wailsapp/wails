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

// TabView creates MultiPanel internally and manages tabs as panels.
type TabView struct {
	ControlBase

	panels           *MultiPanel
	onSelectedChange EventManager
}

func NewTabView(parent Controller) *TabView {
	tv := new(TabView)

	tv.InitControl("SysTabControl32", parent, 0,
		w32.WS_CHILD|w32.WS_VISIBLE|w32.WS_TABSTOP|w32.WS_CLIPSIBLINGS)
	RegMsgHandler(tv)

	tv.panels = NewMultiPanel(parent)

	tv.SetFont(DefaultFont)
	tv.SetSize(200, 24)
	return tv
}

func (tv *TabView) Panels() *MultiPanel {
	return tv.panels
}

func (tv *TabView) tcitemFromPage(panel *Panel) *w32.TCITEM {
	text := syscall.StringToUTF16(panel.Text())
	item := &w32.TCITEM{
		Mask:       w32.TCIF_TEXT,
		PszText:    &text[0],
		CchTextMax: int32(len(text)),
	}
	return item
}

func (tv *TabView) AddPanel(text string) *Panel {
	panel := NewPanel(tv.panels)
	panel.SetText(text)

	item := tv.tcitemFromPage(panel)
	index := tv.panels.Count()
	idx := int(w32.SendMessage(tv.hwnd, w32.TCM_INSERTITEM, uintptr(index), uintptr(unsafe.Pointer(item))))
	if idx == -1 {
		panic("SendMessage(TCM_INSERTITEM) failed")
	}

	tv.panels.AddPanel(panel)
	tv.SetCurrent(idx)
	return panel
}

func (tv *TabView) DeletePanel(index int) {
	w32.SendMessage(tv.hwnd, w32.TCM_DELETEITEM, uintptr(index), 0)
	tv.panels.DeletePanel(index)
	switch {
	case tv.panels.Count() > index:
		tv.SetCurrent(index)
	case tv.panels.Count() == 0:
		tv.SetCurrent(0)
	}
}

func (tv *TabView) Current() int {
	return tv.panels.Current()
}

func (tv *TabView) SetCurrent(index int) {
	if index < 0 || index >= tv.panels.Count() {
		panic("invalid index")
	}
	if ret := int(w32.SendMessage(tv.hwnd, w32.TCM_SETCURSEL, uintptr(index), 0)); ret == -1 {
		panic("SendMessage(TCM_SETCURSEL) failed")
	}
	tv.panels.SetCurrent(index)
}

func (tv *TabView) WndProc(msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case w32.WM_NOTIFY:
		nmhdr := (*w32.NMHDR)(unsafe.Pointer(lparam))

		switch int32(nmhdr.Code) {
		case w32.TCN_SELCHANGE:
			cur := int(w32.SendMessage(tv.hwnd, w32.TCM_GETCURSEL, 0, 0))
			tv.SetCurrent(cur)

			tv.onSelectedChange.Fire(NewEvent(tv, nil))
		}
	}
	return w32.DefWindowProc(tv.hwnd, msg, wparam, lparam)
}
