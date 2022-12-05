//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package winc

import (
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
)

type Controller interface {
	Text() string

	Enabled() bool
	SetFocus()

	Handle() w32.HWND
	Invalidate(erase bool)
	Parent() Controller

	Pos() (x, y int)
	Size() (w, h int)
	Height() int
	Width() int
	Visible() bool
	Bounds() *Rect
	ClientRect() *Rect

	SetText(s string)
	SetEnabled(b bool)
	SetPos(x, y int)
	SetSize(w, h int)
	EnableDragAcceptFiles(b bool)
	Show()
	Hide()

	ContextMenu() *MenuItem
	SetContextMenu(menu *MenuItem)

	Font() *Font
	SetFont(font *Font)
	InvokeRequired() bool
	Invoke(func())
	PreTranslateMessage(msg *w32.MSG) bool
	WndProc(msg uint32, wparam, lparam uintptr) uintptr

	//General events
	OnCreate() *EventManager
	OnClose() *EventManager

	// Focus events
	OnKillFocus() *EventManager
	OnSetFocus() *EventManager

	//Drag and drop events
	OnDropFiles() *EventManager

	//Mouse events
	OnLBDown() *EventManager
	OnLBUp() *EventManager
	OnLBDbl() *EventManager
	OnMBDown() *EventManager
	OnMBUp() *EventManager
	OnRBDown() *EventManager
	OnRBUp() *EventManager
	OnRBDbl() *EventManager
	OnMouseMove() *EventManager

	// OnMouseLeave and OnMouseHover does not fire unless control called internalTrackMouseEvent.
	// Use MouseControl for a how to example.
	OnMouseHover() *EventManager
	OnMouseLeave() *EventManager

	//Keyboard events
	OnKeyUp() *EventManager

	//Paint events
	OnPaint() *EventManager
	OnSize() *EventManager

	invokeCallbacks()
}
