//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package winc

import (
	"unsafe"

	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
)

var wmInvokeCallback uint32

func init() {
	wmInvokeCallback = RegisterWindowMessage("WincV0.InvokeCallback")
}

func genPoint(p uintptr) (x, y int) {
	x = int(w32.LOWORD(uint32(p)))
	y = int(w32.HIWORD(uint32(p)))
	return
}

func genMouseEventArg(wparam, lparam uintptr) *MouseEventData {
	var data MouseEventData
	data.Button = int(wparam)
	data.X, data.Y = genPoint(lparam)

	return &data
}

func genDropFilesEventArg(wparam uintptr) *DropFilesEventData {
	hDrop := w32.HDROP(wparam)

	var data DropFilesEventData
	_, fileCount := w32.DragQueryFile(hDrop, 0xFFFFFFFF)
	data.Files = make([]string, fileCount)

	var i uint
	for i = 0; i < fileCount; i++ {
		data.Files[i], _ = w32.DragQueryFile(hDrop, i)
	}

	data.X, data.Y, _ = w32.DragQueryPoint(hDrop)
	w32.DragFinish(hDrop)
	return &data
}

func generalWndProc(hwnd w32.HWND, msg uint32, wparam, lparam uintptr) uintptr {

	switch msg {
	case w32.WM_HSCROLL:
		//println("case w32.WM_HSCROLL")

	case w32.WM_VSCROLL:
		//println("case w32.WM_VSCROLL")
	}

	if controller := GetMsgHandler(hwnd); controller != nil {
		ret := controller.WndProc(msg, wparam, lparam)

		switch msg {
		case w32.WM_NOTIFY: //Reflect notification to control
			nm := (*w32.NMHDR)(unsafe.Pointer(lparam))
			if controller := GetMsgHandler(nm.HwndFrom); controller != nil {
				ret := controller.WndProc(msg, wparam, lparam)
				if ret != 0 {
					w32.SetWindowLong(hwnd, w32.DWL_MSGRESULT, uint32(ret))
					return w32.TRUE
				}
			}
		case w32.WM_COMMAND:
			if lparam != 0 { //Reflect message to control
				h := w32.HWND(lparam)
				if controller := GetMsgHandler(h); controller != nil {
					ret := controller.WndProc(msg, wparam, lparam)
					if ret != 0 {
						w32.SetWindowLong(hwnd, w32.DWL_MSGRESULT, uint32(ret))
						return w32.TRUE
					}
				}
			}
		case w32.WM_CLOSE:
			controller.OnClose().Fire(NewEvent(controller, nil))
		case w32.WM_KILLFOCUS:
			controller.OnKillFocus().Fire(NewEvent(controller, nil))
		case w32.WM_SETFOCUS:
			controller.OnSetFocus().Fire(NewEvent(controller, nil))
		case w32.WM_DROPFILES:
			controller.OnDropFiles().Fire(NewEvent(controller, genDropFilesEventArg(wparam)))
		case w32.WM_CONTEXTMENU:
			if wparam != 0 { //Reflect message to control
				h := w32.HWND(wparam)
				if controller := GetMsgHandler(h); controller != nil {
					contextMenu := controller.ContextMenu()
					x, y := genPoint(lparam)

					if contextMenu != nil {
						id := w32.TrackPopupMenuEx(
							contextMenu.hMenu,
							w32.TPM_NOANIMATION|w32.TPM_RETURNCMD,
							int32(x),
							int32(y),
							controller.Handle(),
							nil)

						item := findMenuItemByID(int(id))
						if item != nil {
							item.OnClick().Fire(NewEvent(controller, genMouseEventArg(wparam, lparam)))
						}
						return 0
					}
				}
			}

		case w32.WM_LBUTTONDOWN:
			controller.OnLBDown().Fire(NewEvent(controller, genMouseEventArg(wparam, lparam)))
		case w32.WM_LBUTTONUP:
			controller.OnLBUp().Fire(NewEvent(controller, genMouseEventArg(wparam, lparam)))
		case w32.WM_LBUTTONDBLCLK:
			controller.OnLBDbl().Fire(NewEvent(controller, genMouseEventArg(wparam, lparam)))
		case w32.WM_MBUTTONDOWN:
			controller.OnMBDown().Fire(NewEvent(controller, genMouseEventArg(wparam, lparam)))
		case w32.WM_MBUTTONUP:
			controller.OnMBUp().Fire(NewEvent(controller, genMouseEventArg(wparam, lparam)))
		case w32.WM_RBUTTONDOWN:
			controller.OnRBDown().Fire(NewEvent(controller, genMouseEventArg(wparam, lparam)))
		case w32.WM_RBUTTONUP:
			controller.OnRBUp().Fire(NewEvent(controller, genMouseEventArg(wparam, lparam)))
		case w32.WM_RBUTTONDBLCLK:
			controller.OnRBDbl().Fire(NewEvent(controller, genMouseEventArg(wparam, lparam)))
		case w32.WM_MOUSEMOVE:
			controller.OnMouseMove().Fire(NewEvent(controller, genMouseEventArg(wparam, lparam)))
		case w32.WM_PAINT:
			canvas := NewCanvasFromHwnd(hwnd)
			defer canvas.Dispose()
			controller.OnPaint().Fire(NewEvent(controller, &PaintEventData{Canvas: canvas}))
		case w32.WM_KEYUP:
			controller.OnKeyUp().Fire(NewEvent(controller, &KeyUpEventData{int(wparam), int(lparam)}))
		case w32.WM_SIZE:
			x, y := genPoint(lparam)
			controller.OnSize().Fire(NewEvent(controller, &SizeEventData{uint(wparam), x, y}))
		case wmInvokeCallback:
			controller.invokeCallbacks()
		}
		return ret
	}

	return w32.DefWindowProc(hwnd, uint32(msg), wparam, lparam)
}
