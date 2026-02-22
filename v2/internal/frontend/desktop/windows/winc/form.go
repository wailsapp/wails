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

type LayoutManager interface {
	Update()
}

// A Form is main window of the application.
type Form struct {
	ControlBase

	layoutMng LayoutManager

	// Fullscreen / Unfullscreen
	isFullscreen            bool
	previousWindowStyle     uint32
	previousWindowExStyle   uint32
	previousWindowPlacement w32.WINDOWPLACEMENT
}

func NewCustomForm(parent Controller, exStyle int, dwStyle uint) *Form {
	fm := new(Form)

	RegClassOnlyOnce("winc_Form")

	fm.isForm = true

	if exStyle == 0 {
		exStyle = w32.WS_EX_CONTROLPARENT | w32.WS_EX_APPWINDOW
	}

	if dwStyle == 0 {
		dwStyle = w32.WS_OVERLAPPEDWINDOW
	}

	fm.hwnd = CreateWindow("winc_Form", parent, uint(exStyle), dwStyle)
	fm.parent = parent

	// this might fail if icon resource is not embedded in the binary
	if ico, err := NewIconFromResource(GetAppInstance(), uint16(AppIconID)); err == nil {
		fm.SetIcon(0, ico)
	}

	// This forces display of focus rectangles, as soon as the user starts to type.
	w32.SendMessage(fm.hwnd, w32.WM_CHANGEUISTATE, w32.UIS_INITIALIZE, 0)

	RegMsgHandler(fm)

	fm.SetFont(DefaultFont)
	fm.SetText("Form")
	return fm
}

func NewForm(parent Controller) *Form {
	fm := new(Form)

	RegClassOnlyOnce("winc_Form")

	fm.isForm = true
	fm.hwnd = CreateWindow("winc_Form", parent, w32.WS_EX_CONTROLPARENT|w32.WS_EX_APPWINDOW, w32.WS_OVERLAPPEDWINDOW)
	fm.parent = parent

	// this might fail if icon resource is not embedded in the binary
	if ico, err := NewIconFromResource(GetAppInstance(), uint16(AppIconID)); err == nil {
		fm.SetIcon(0, ico)
	}

	// This forces display of focus rectangles, as soon as the user starts to type.
	w32.SendMessage(fm.hwnd, w32.WM_CHANGEUISTATE, w32.UIS_INITIALIZE, 0)

	RegMsgHandler(fm)

	fm.SetFont(DefaultFont)
	fm.SetText("Form")
	return fm
}

func (fm *Form) SetLayout(mng LayoutManager) {
	fm.layoutMng = mng
}

// UpdateLayout refresh layout.
func (fm *Form) UpdateLayout() {
	if fm.layoutMng != nil {
		fm.layoutMng.Update()
	}
}

func (fm *Form) NewMenu() *Menu {
	hMenu := w32.CreateMenu()
	if hMenu == 0 {
		panic("failed CreateMenu")
	}
	m := &Menu{hMenu: hMenu, hwnd: fm.hwnd}
	if !w32.SetMenu(fm.hwnd, hMenu) {
		panic("failed SetMenu")
	}
	return m
}

func (fm *Form) DisableIcon() {
	windowInfo := getWindowInfo(fm.hwnd)
	frameless := windowInfo.IsPopup()
	if frameless {
		return
	}
	exStyle := w32.GetWindowLong(fm.hwnd, w32.GWL_EXSTYLE)
	w32.SetWindowLong(fm.hwnd, w32.GWL_EXSTYLE, uint32(exStyle|w32.WS_EX_DLGMODALFRAME))
	w32.SetWindowPos(fm.hwnd, 0, 0, 0, 0, 0,
		uint(
			w32.SWP_FRAMECHANGED|
				w32.SWP_NOMOVE|
				w32.SWP_NOSIZE|
				w32.SWP_NOZORDER),
	)
}

func (fm *Form) Maximise() {
	w32.ShowWindow(fm.hwnd, w32.SW_MAXIMIZE)
}

func (fm *Form) Minimise() {
	w32.ShowWindow(fm.hwnd, w32.SW_MINIMIZE)
}

func (fm *Form) Restore() {
	// SC_RESTORE param for WM_SYSCOMMAND to restore app if it is minimized
	const SC_RESTORE = 0xF120
	// restore the minimized window, if it is
	w32.SendMessage(
		fm.hwnd,
		w32.WM_SYSCOMMAND,
		SC_RESTORE,
		0,
	)
	w32.ShowWindow(fm.hwnd, w32.SW_RESTORE)
}

// Public methods
func (fm *Form) Center() {

	windowInfo := getWindowInfo(fm.hwnd)
	frameless := windowInfo.IsPopup()

	info := getMonitorInfo(fm.hwnd)
	workRect := info.RcWork
	screenMiddleW := workRect.Left + (workRect.Right-workRect.Left)/2
	screenMiddleH := workRect.Top + (workRect.Bottom-workRect.Top)/2
	var winRect *w32.RECT
	if !frameless {
		winRect = w32.GetWindowRect(fm.hwnd)
	} else {
		winRect = w32.GetClientRect(fm.hwnd)
	}
	winWidth := winRect.Right - winRect.Left
	winHeight := winRect.Bottom - winRect.Top
	windowX := screenMiddleW - (winWidth / 2)
	windowY := screenMiddleH - (winHeight / 2)
	w32.SetWindowPos(fm.hwnd, w32.HWND_TOP, int(windowX), int(windowY), int(winWidth), int(winHeight), w32.SWP_NOSIZE)
}

func (fm *Form) Fullscreen() {
	if fm.isFullscreen {
		return
	}

	fm.previousWindowStyle = uint32(w32.GetWindowLongPtr(fm.hwnd, w32.GWL_STYLE))
	fm.previousWindowExStyle = uint32(w32.GetWindowLong(fm.hwnd, w32.GWL_EXSTYLE))

	monitor := w32.MonitorFromWindow(fm.hwnd, w32.MONITOR_DEFAULTTOPRIMARY)
	var monitorInfo w32.MONITORINFO
	monitorInfo.CbSize = uint32(unsafe.Sizeof(monitorInfo))
	if !w32.GetMonitorInfo(monitor, &monitorInfo) {
		return
	}
	if !w32.GetWindowPlacement(fm.hwnd, &fm.previousWindowPlacement) {
		return
	}
	// According to https://devblogs.microsoft.com/oldnewthing/20050505-04/?p=35703 one should use w32.WS_POPUP | w32.WS_VISIBLE
	w32.SetWindowLong(fm.hwnd, w32.GWL_STYLE, fm.previousWindowStyle & ^uint32(w32.WS_OVERLAPPEDWINDOW) | (w32.WS_POPUP|w32.WS_VISIBLE))
	w32.SetWindowLong(fm.hwnd, w32.GWL_EXSTYLE, fm.previousWindowExStyle & ^uint32(w32.WS_EX_DLGMODALFRAME))
	fm.isFullscreen = true
	w32.SetWindowPos(fm.hwnd, w32.HWND_TOP,
		int(monitorInfo.RcMonitor.Left),
		int(monitorInfo.RcMonitor.Top),
		int(monitorInfo.RcMonitor.Right-monitorInfo.RcMonitor.Left),
		int(monitorInfo.RcMonitor.Bottom-monitorInfo.RcMonitor.Top),
		w32.SWP_NOOWNERZORDER|w32.SWP_FRAMECHANGED)
}

func (fm *Form) UnFullscreen() {
	if !fm.isFullscreen {
		return
	}
	w32.SetWindowLong(fm.hwnd, w32.GWL_STYLE, fm.previousWindowStyle)
	w32.SetWindowLong(fm.hwnd, w32.GWL_EXSTYLE, fm.previousWindowExStyle)
	w32.SetWindowPlacement(fm.hwnd, &fm.previousWindowPlacement)
	fm.isFullscreen = false
	w32.SetWindowPos(fm.hwnd, 0, 0, 0, 0, 0,
		w32.SWP_NOMOVE|w32.SWP_NOSIZE|w32.SWP_NOZORDER|w32.SWP_NOOWNERZORDER|w32.SWP_FRAMECHANGED)
}

func (fm *Form) IsFullScreen() bool {
	return fm.isFullscreen
}

// IconType: 1 - ICON_BIG; 0 - ICON_SMALL
func (fm *Form) SetIcon(iconType int, icon *Icon) {
	if iconType > 1 {
		panic("IconType is invalid")
	}
	w32.SendMessage(fm.hwnd, w32.WM_SETICON, uintptr(iconType), uintptr(icon.Handle()))
}

func (fm *Form) EnableMaxButton(b bool) {
	SetStyle(fm.hwnd, b, w32.WS_MAXIMIZEBOX)
}

func (fm *Form) EnableMinButton(b bool) {
	SetStyle(fm.hwnd, b, w32.WS_MINIMIZEBOX)
}

func (fm *Form) EnableSizable(b bool) {
	SetStyle(fm.hwnd, b, w32.WS_THICKFRAME)
}

func (fm *Form) EnableDragMove(_ bool) {
	//fm.isDragMove = b
}

func (fm *Form) EnableTopMost(b bool) {
	tag := w32.HWND_NOTOPMOST
	if b {
		tag = w32.HWND_TOPMOST
	}
	w32.SetWindowPos(fm.hwnd, tag, 0, 0, 0, 0, w32.SWP_NOMOVE|w32.SWP_NOSIZE)
}

func (fm *Form) WndProc(msg uint32, wparam, lparam uintptr) uintptr {

	switch msg {
	case w32.WM_COMMAND:
		if lparam == 0 && w32.HIWORD(uint32(wparam)) == 0 {
			// Menu support.
			actionID := uint16(w32.LOWORD(uint32(wparam)))
			if action, ok := actionsByID[actionID]; ok {
				action.onClick.Fire(NewEvent(fm, nil))
			}
		}
	case w32.WM_KEYDOWN:
		// Accelerator support.
		key := Key(wparam)
		if uint32(lparam)>>30 == 0 {
			// Using TranslateAccelerators refused to work, so we handle them
			// ourselves, at least for now.
			shortcut := Shortcut{ModifiersDown(), key}
			if action, ok := shortcut2Action[shortcut]; ok {
				if action.Enabled() {
					action.onClick.Fire(NewEvent(fm, nil))
				}
			}
		}

	case w32.WM_CLOSE:
		return 0
	case w32.WM_DESTROY:
		w32.PostQuitMessage(0)
		return 0

	case w32.WM_SIZE, w32.WM_PAINT:
		if fm.layoutMng != nil {
			fm.layoutMng.Update()
		}
	case w32.WM_GETMINMAXINFO:
		mmi := (*w32.MINMAXINFO)(unsafe.Pointer(lparam))
		hasConstraints := false
		if fm.minWidth > 0 || fm.minHeight > 0 {
			hasConstraints = true

			width, height := fm.scaleWithWindowDPI(fm.minWidth, fm.minHeight)
			if width > 0 {
				mmi.PtMinTrackSize.X = int32(width)
			}
			if height > 0 {
				mmi.PtMinTrackSize.Y = int32(height)
			}
		}
		if fm.maxWidth > 0 || fm.maxHeight > 0 {
			hasConstraints = true

			width, height := fm.scaleWithWindowDPI(fm.maxWidth, fm.maxHeight)
			if width > 0 {
				mmi.PtMaxTrackSize.X = int32(width)
			}
			if height > 0 {
				mmi.PtMaxTrackSize.Y = int32(height)
			}
		}
		if hasConstraints {
			return 0
		}
	}

	return w32.DefWindowProc(fm.hwnd, msg, wparam, lparam)
}
