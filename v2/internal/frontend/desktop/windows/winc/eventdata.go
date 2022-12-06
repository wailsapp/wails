//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package winc

import (
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
)

type RawMsg struct {
	Hwnd           w32.HWND
	Msg            uint32
	WParam, LParam uintptr
}

type MouseEventData struct {
	X, Y   int
	Button int
	Wheel  int
}

type DropFilesEventData struct {
	X, Y  int
	Files []string
}

type PaintEventData struct {
	Canvas *Canvas
}

type LabelEditEventData struct {
	Item ListItem
	Text string
	//PszText *uint16
}

/*type LVDBLClickEventData struct {
	NmItem *w32.NMITEMACTIVATE
}*/

type KeyUpEventData struct {
	VKey, Code int
}

type SizeEventData struct {
	Type uint
	X, Y int
}
