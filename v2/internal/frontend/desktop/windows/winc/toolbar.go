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

type Toolbar struct {
	ControlBase
	iml *ImageList

	buttons []*ToolButton
}

type ToolButton struct {
	tb *Toolbar

	text      string
	enabled   bool
	checkable bool
	checked   bool
	image     int

	onClick EventManager
}

func (bt *ToolButton) OnClick() *EventManager {
	return &bt.onClick
}

func (bt *ToolButton) update() { bt.tb.update(bt) }

func (bt *ToolButton) IsSeparator() bool { return bt.text == "-" }
func (bt *ToolButton) SetSeparator()     { bt.text = "-" }

func (bt *ToolButton) Enabled() bool     { return bt.enabled }
func (bt *ToolButton) SetEnabled(b bool) { bt.enabled = b; bt.update() }

func (bt *ToolButton) Checkable() bool     { return bt.checkable }
func (bt *ToolButton) SetCheckable(b bool) { bt.checkable = b; bt.update() }

func (bt *ToolButton) Checked() bool     { return bt.checked }
func (bt *ToolButton) SetChecked(b bool) { bt.checked = b; bt.update() }

func (bt *ToolButton) Text() string     { return bt.text }
func (bt *ToolButton) SetText(s string) { bt.text = s; bt.update() }

func (bt *ToolButton) Image() int     { return bt.image }
func (bt *ToolButton) SetImage(i int) { bt.image = i; bt.update() }

// NewHToolbar creates horizontal toolbar with text on same line as image.
func NewHToolbar(parent Controller) *Toolbar {
	return newToolbar(parent, w32.CCS_NODIVIDER|w32.TBSTYLE_FLAT|w32.TBSTYLE_TOOLTIPS|w32.TBSTYLE_WRAPABLE|
		w32.WS_CHILD|w32.TBSTYLE_LIST)
}

// NewToolbar creates toolbar with text below the image.
func NewToolbar(parent Controller) *Toolbar {
	return newToolbar(parent, w32.CCS_NODIVIDER|w32.TBSTYLE_FLAT|w32.TBSTYLE_TOOLTIPS|w32.TBSTYLE_WRAPABLE|
		w32.WS_CHILD /*|w32.TBSTYLE_TRANSPARENT*/)
}

func newToolbar(parent Controller, style uint) *Toolbar {
	tb := new(Toolbar)

	tb.InitControl("ToolbarWindow32", parent, 0, style)

	exStyle := w32.SendMessage(tb.hwnd, w32.TB_GETEXTENDEDSTYLE, 0, 0)
	exStyle |= w32.TBSTYLE_EX_DRAWDDARROWS | w32.TBSTYLE_EX_MIXEDBUTTONS
	w32.SendMessage(tb.hwnd, w32.TB_SETEXTENDEDSTYLE, 0, exStyle)
	RegMsgHandler(tb)

	tb.SetFont(DefaultFont)
	tb.SetPos(0, 0)
	tb.SetSize(200, 40)

	return tb
}

func (tb *Toolbar) SetImageList(imageList *ImageList) {
	w32.SendMessage(tb.hwnd, w32.TB_SETIMAGELIST, 0, uintptr(imageList.Handle()))
	tb.iml = imageList
}

func (tb *Toolbar) initButton(btn *ToolButton, state, style *byte, image *int32, text *uintptr) {
	*style |= w32.BTNS_AUTOSIZE

	if btn.checked {
		*state |= w32.TBSTATE_CHECKED
	}

	if btn.enabled {
		*state |= w32.TBSTATE_ENABLED
	}

	if btn.checkable {
		*style |= w32.BTNS_CHECK
	}

	if len(btn.Text()) > 0 {
		*style |= w32.BTNS_SHOWTEXT
	}

	if btn.IsSeparator() {
		*style = w32.BTNS_SEP
	}

	*image = int32(btn.Image())
	*text = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(btn.Text())))
}

func (tb *Toolbar) update(btn *ToolButton) {
	tbbi := w32.TBBUTTONINFO{
		DwMask: w32.TBIF_IMAGE | w32.TBIF_STATE | w32.TBIF_STYLE | w32.TBIF_TEXT,
	}

	tbbi.CbSize = uint32(unsafe.Sizeof(tbbi))

	var i int
	for i = range tb.buttons {
		if tb.buttons[i] == btn {
			break
		}
	}

	tb.initButton(btn, &tbbi.FsState, &tbbi.FsStyle, &tbbi.IImage, &tbbi.PszText)
	if w32.SendMessage(tb.hwnd, w32.TB_SETBUTTONINFO, uintptr(i), uintptr(unsafe.Pointer(&tbbi))) == 0 {
		panic("SendMessage(TB_SETBUTTONINFO) failed")
	}
}

func (tb *Toolbar) AddSeparator() {
	tb.AddButton("-", 0)
}

// AddButton creates and adds button to the toolbar. Use returned toolbutton to setup OnClick event.
func (tb *Toolbar) AddButton(text string, image int) *ToolButton {
	bt := &ToolButton{
		tb:      tb, // points to parent
		text:    text,
		image:   image,
		enabled: true,
	}
	tb.buttons = append(tb.buttons, bt)
	index := len(tb.buttons) - 1

	tbb := w32.TBBUTTON{
		IdCommand: int32(index),
	}

	tb.initButton(bt, &tbb.FsState, &tbb.FsStyle, &tbb.IBitmap, &tbb.IString)
	w32.SendMessage(tb.hwnd, w32.TB_BUTTONSTRUCTSIZE, uintptr(unsafe.Sizeof(tbb)), 0)

	if w32.SendMessage(tb.hwnd, w32.TB_INSERTBUTTON, uintptr(index), uintptr(unsafe.Pointer(&tbb))) == w32.FALSE {
		panic("SendMessage(TB_ADDBUTTONS)")
	}

	w32.SendMessage(tb.hwnd, w32.TB_AUTOSIZE, 0, 0)
	return bt
}

func (tb *Toolbar) WndProc(msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case w32.WM_COMMAND:
		switch w32.HIWORD(uint32(wparam)) {
		case w32.BN_CLICKED:
			id := uint16(w32.LOWORD(uint32(wparam)))
			btn := tb.buttons[id]
			btn.onClick.Fire(NewEvent(tb, nil))
		}
	}
	return w32.DefWindowProc(tb.hwnd, msg, wparam, lparam)
}
