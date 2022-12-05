//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 */

package winc

import (
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
)

type VResizer struct {
	ControlBase

	control1 Dockable
	control2 Dockable
	dir      Direction

	mouseLeft bool
	drag      bool
}

func NewVResizer(parent Controller) *VResizer {
	sp := new(VResizer)

	RegClassOnlyOnce("winc_VResizer")
	sp.hwnd = CreateWindow("winc_VResizer", parent, w32.WS_EX_CONTROLPARENT, w32.WS_CHILD|w32.WS_VISIBLE)
	sp.parent = parent
	sp.mouseLeft = true
	RegMsgHandler(sp)

	sp.SetFont(DefaultFont)
	sp.SetText("")
	sp.SetSize(20, 100)
	return sp
}

func (sp *VResizer) SetControl(control1, control2 Dockable, dir Direction, minSize int) {
	sp.control1 = control1
	sp.control2 = control2
	if dir != Left && dir != Right {
		panic("invalid direction")
	}
	sp.dir = dir

	// TODO(vi): ADDED
	/*internalTrackMouseEvent(control1.Handle())
	internalTrackMouseEvent(control2.Handle())

	control1.OnMouseMove().Bind(func(e *Event) {
		if sp.drag {
			x := e.Data.(*MouseEventData).X
			sp.update(x)
			w32.SetCursor(w32.LoadCursorWithResourceID(0, w32.IDC_SIZEWE))

		}
		fmt.Println("control1.OnMouseMove")
	})

	control2.OnMouseMove().Bind(func(e *Event) {
		if sp.drag {
			x := e.Data.(*MouseEventData).X
			sp.update(x)
			w32.SetCursor(w32.LoadCursorWithResourceID(0, w32.IDC_SIZEWE))

		}
		fmt.Println("control2.OnMouseMove")
	})

	control1.OnLBUp().Bind(func(e *Event) {
		sp.drag = false
		sp.mouseLeft = true
		fmt.Println("control1.OnLBUp")
	})

	control2.OnLBUp().Bind(func(e *Event) {
		sp.drag = false
		sp.mouseLeft = true
		fmt.Println("control2.OnLBUp")
	})*/

	// ---- finish ADDED

}

func (sp *VResizer) update(x int) {
	pos := x - 10

	w1, h1 := sp.control1.Width(), sp.control1.Height()
	if sp.dir == Left {
		w1 += pos
	} else {
		w1 -= pos
	}
	sp.control1.SetSize(w1, h1)
	fm := sp.parent.(*Form)
	fm.UpdateLayout()

	w32.SetCursor(w32.LoadCursorWithResourceID(0, w32.IDC_ARROW))
}

func (sp *VResizer) WndProc(msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case w32.WM_CREATE:
		internalTrackMouseEvent(sp.hwnd)

	case w32.WM_MOUSEMOVE:
		if sp.drag {
			x, _ := genPoint(lparam)
			sp.update(x)
		} else {
			w32.SetCursor(w32.LoadCursorWithResourceID(0, w32.IDC_SIZEWE))
		}

		if sp.mouseLeft {
			internalTrackMouseEvent(sp.hwnd)
			sp.mouseLeft = false
		}

	case w32.WM_MOUSELEAVE:
		sp.drag = false
		sp.mouseLeft = true

	case w32.WM_LBUTTONUP:
		sp.drag = false

	case w32.WM_LBUTTONDOWN:
		sp.drag = true
	}
	return w32.DefWindowProc(sp.hwnd, msg, wparam, lparam)
}

type HResizer struct {
	ControlBase

	control1  Dockable
	control2  Dockable
	dir       Direction
	mouseLeft bool
	drag      bool
}

func NewHResizer(parent Controller) *HResizer {
	sp := new(HResizer)

	RegClassOnlyOnce("winc_HResizer")
	sp.hwnd = CreateWindow("winc_HResizer", parent, w32.WS_EX_CONTROLPARENT, w32.WS_CHILD|w32.WS_VISIBLE)
	sp.parent = parent
	sp.mouseLeft = true
	RegMsgHandler(sp)

	sp.SetFont(DefaultFont)
	sp.SetText("")
	sp.SetSize(100, 20)

	return sp
}

func (sp *HResizer) SetControl(control1, control2 Dockable, dir Direction, minSize int) {
	sp.control1 = control1
	sp.control2 = control2
	if dir != Top && dir != Bottom {
		panic("invalid direction")
	}
	sp.dir = dir

}

func (sp *HResizer) update(y int) {
	pos := y - 10

	w1, h1 := sp.control1.Width(), sp.control1.Height()
	if sp.dir == Top {
		h1 += pos
	} else {
		h1 -= pos
	}
	sp.control1.SetSize(w1, h1)

	fm := sp.parent.(*Form)
	fm.UpdateLayout()

	w32.SetCursor(w32.LoadCursorWithResourceID(0, w32.IDC_ARROW))
}

func (sp *HResizer) WndProc(msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case w32.WM_CREATE:
		internalTrackMouseEvent(sp.hwnd)

	case w32.WM_MOUSEMOVE:
		if sp.drag {
			_, y := genPoint(lparam)
			sp.update(y)
		} else {
			w32.SetCursor(w32.LoadCursorWithResourceID(0, w32.IDC_SIZENS))
		}

		if sp.mouseLeft {
			internalTrackMouseEvent(sp.hwnd)
			sp.mouseLeft = false
		}

	case w32.WM_MOUSELEAVE:
		sp.drag = false
		sp.mouseLeft = true

	case w32.WM_LBUTTONUP:
		sp.drag = false

	case w32.WM_LBUTTONDOWN:
		sp.drag = true
	}
	return w32.DefWindowProc(sp.hwnd, msg, wparam, lparam)
}
