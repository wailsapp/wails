//go:build windows

/*
 * Copyright (C) 2019 The Winc Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package winc

import (
	"fmt"
	"runtime"
	"sync"
	"syscall"
	"unsafe"

	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
)

type ControlBase struct {
	hwnd        w32.HWND
	font        *Font
	parent      Controller
	contextMenu *MenuItem

	isForm bool

	minWidth, minHeight int
	maxWidth, maxHeight int

	// General events
	onCreate EventManager
	onClose  EventManager

	// Focus events
	onKillFocus EventManager
	onSetFocus  EventManager

	// Drag and drop events
	onDropFiles EventManager

	// Mouse events
	onLBDown    EventManager
	onLBUp      EventManager
	onLBDbl     EventManager
	onMBDown    EventManager
	onMBUp      EventManager
	onRBDown    EventManager
	onRBUp      EventManager
	onRBDbl     EventManager
	onMouseMove EventManager

	// use MouseControl to capture onMouseHover and onMouseLeave events.
	onMouseHover EventManager
	onMouseLeave EventManager

	// Keyboard events
	onKeyUp EventManager

	// Paint events
	onPaint EventManager
	onSize  EventManager

	m         sync.Mutex
	dispatchq []func()
}

// initControl is called by controls: edit, button, treeview, listview, and so on.
func (cba *ControlBase) InitControl(className string, parent Controller, exstyle, style uint) {
	cba.hwnd = CreateWindow(className, parent, exstyle, style)
	if cba.hwnd == 0 {
		panic("cannot create window for " + className)
	}
	cba.parent = parent
}

// InitWindow is called by custom window based controls such as split, panel, etc.
func (cba *ControlBase) InitWindow(className string, parent Controller, exstyle, style uint) {
	RegClassOnlyOnce(className)
	cba.hwnd = CreateWindow(className, parent, exstyle, style)
	if cba.hwnd == 0 {
		panic("cannot create window for " + className)
	}
	cba.parent = parent
}

// SetTheme for TreeView and ListView controls.
func (cba *ControlBase) SetTheme(appName string) error {
	if hr := w32.SetWindowTheme(cba.hwnd, syscall.StringToUTF16Ptr(appName), nil); w32.FAILED(hr) {
		return fmt.Errorf("SetWindowTheme %d", hr)
	}
	return nil
}

func (cba *ControlBase) Handle() w32.HWND {
	return cba.hwnd
}

func (cba *ControlBase) SetHandle(hwnd w32.HWND) {
	cba.hwnd = hwnd
}

func (cba *ControlBase) GetWindowDPI() (w32.UINT, w32.UINT) {
	if w32.HasGetDpiForWindowFunc() {
		// GetDpiForWindow is supported beginning with Windows 10, 1607 and is the most accureate
		// one, especially it is consistent with the WM_DPICHANGED event.
		dpi := w32.GetDpiForWindow(cba.hwnd)
		return dpi, dpi
	}

	if w32.HasGetDPIForMonitorFunc() {
		// GetDpiForWindow is supported beginning with Windows 8.1
		monitor := w32.MonitorFromWindow(cba.hwnd, w32.MONITOR_DEFAULTTONEAREST)
		if monitor == 0 {
			return 0, 0
		}
		var dpiX, dpiY w32.UINT
		w32.GetDPIForMonitor(monitor, w32.MDT_EFFECTIVE_DPI, &dpiX, &dpiY)
		return dpiX, dpiY
	}

	// If none of the above is supported fallback to the System DPI.
	screen := w32.GetDC(0)
	x := w32.GetDeviceCaps(screen, w32.LOGPIXELSX)
	y := w32.GetDeviceCaps(screen, w32.LOGPIXELSY)
	w32.ReleaseDC(0, screen)
	return w32.UINT(x), w32.UINT(y)
}

func (cba *ControlBase) SetAndClearStyleBits(set, clear uint32) error {
	style := uint32(w32.GetWindowLong(cba.hwnd, w32.GWL_STYLE))
	if style == 0 {
		return fmt.Errorf("GetWindowLong")
	}

	if newStyle := style&^clear | set; newStyle != style {
		if w32.SetWindowLong(cba.hwnd, w32.GWL_STYLE, newStyle) == 0 {
			return fmt.Errorf("SetWindowLong")
		}
	}
	return nil
}

func (cba *ControlBase) SetIsForm(isform bool) {
	cba.isForm = isform
}

func (cba *ControlBase) SetText(caption string) {
	w32.SetWindowText(cba.hwnd, caption)
}

func (cba *ControlBase) Text() string {
	return w32.GetWindowText(cba.hwnd)
}

func (cba *ControlBase) Close() {
	UnRegMsgHandler(cba.hwnd)
	w32.DestroyWindow(cba.hwnd)
}

func (cba *ControlBase) SetTranslucentBackground() {
	var accent = w32.ACCENT_POLICY{
		AccentState: w32.ACCENT_ENABLE_BLURBEHIND,
	}
	var data w32.WINDOWCOMPOSITIONATTRIBDATA
	data.Attrib = w32.WCA_ACCENT_POLICY
	data.PvData = unsafe.Pointer(&accent)
	data.CbData = unsafe.Sizeof(accent)

	w32.SetWindowCompositionAttribute(cba.hwnd, &data)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (cba *ControlBase) clampSize(width, height int) (int, int) {
	if cba.minWidth != 0 {
		width = max(width, cba.minWidth)
	}
	if cba.maxWidth != 0 {
		width = min(width, cba.maxWidth)
	}
	if cba.minHeight != 0 {
		height = max(height, cba.minHeight)
	}
	if cba.maxHeight != 0 {
		height = min(height, cba.maxHeight)
	}
	return width, height
}

func (cba *ControlBase) SetSize(width, height int) {
	x, y := cba.Pos()
	width, height = cba.clampSize(width, height)
	width, height = cba.scaleWithWindowDPI(width, height)
	w32.MoveWindow(cba.hwnd, x, y, width, height, true)
}

func (cba *ControlBase) SetMinSize(width, height int) {
	cba.minWidth = width
	cba.minHeight = height

	// Ensure we set max if min > max
	if cba.maxWidth > 0 {
		cba.maxWidth = max(cba.minWidth, cba.maxWidth)
	}
	if cba.maxHeight > 0 {
		cba.maxHeight = max(cba.minHeight, cba.maxHeight)
	}

	x, y := cba.Pos()
	currentWidth, currentHeight := cba.Size()
	clampedWidth, clampedHeight := cba.clampSize(currentWidth, currentHeight)
	if clampedWidth != currentWidth || clampedHeight != currentHeight {
		w32.MoveWindow(cba.hwnd, x, y, clampedWidth, clampedHeight, true)
	}
}
func (cba *ControlBase) SetMaxSize(width, height int) {
	cba.maxWidth = width
	cba.maxHeight = height

	// Ensure we set min if max > min
	if cba.maxWidth > 0 {
		cba.minWidth = min(cba.maxWidth, cba.minWidth)
	}
	if cba.maxHeight > 0 {
		cba.minHeight = min(cba.maxHeight, cba.minHeight)
	}

	x, y := cba.Pos()
	currentWidth, currentHeight := cba.Size()
	clampedWidth, clampedHeight := cba.clampSize(currentWidth, currentHeight)
	if clampedWidth != currentWidth || clampedHeight != currentHeight {
		w32.MoveWindow(cba.hwnd, x, y, clampedWidth, clampedHeight, true)
	}
}

func (cba *ControlBase) Size() (width, height int) {
	rect := w32.GetWindowRect(cba.hwnd)
	width = int(rect.Right - rect.Left)
	height = int(rect.Bottom - rect.Top)
	width, height = cba.scaleToDefaultDPI(width, height)
	return
}

func (cba *ControlBase) Width() int {
	rect := w32.GetWindowRect(cba.hwnd)
	return int(rect.Right - rect.Left)
}

func (cba *ControlBase) Height() int {
	rect := w32.GetWindowRect(cba.hwnd)
	return int(rect.Bottom - rect.Top)
}

func (cba *ControlBase) SetPos(x, y int) {
	info := getMonitorInfo(cba.hwnd)
	workRect := info.RcWork

	w32.SetWindowPos(cba.hwnd, w32.HWND_TOP, int(workRect.Left)+x, int(workRect.Top)+y, 0, 0, w32.SWP_NOSIZE)
}
func (cba *ControlBase) SetAlwaysOnTop(b bool) {
	if b {
		w32.SetWindowPos(cba.hwnd, w32.HWND_TOPMOST, 0, 0, 0, 0, w32.SWP_NOSIZE|w32.SWP_NOMOVE)
	} else {
		w32.SetWindowPos(cba.hwnd, w32.HWND_NOTOPMOST, 0, 0, 0, 0, w32.SWP_NOSIZE|w32.SWP_NOMOVE)
	}
}

func (cba *ControlBase) Pos() (x, y int) {
	rect := w32.GetWindowRect(cba.hwnd)
	x = int(rect.Left)
	y = int(rect.Top)
	if !cba.isForm && cba.parent != nil {
		x, y, _ = w32.ScreenToClient(cba.parent.Handle(), x, y)
	}
	return
}

func (cba *ControlBase) Visible() bool {
	return w32.IsWindowVisible(cba.hwnd)
}

func (cba *ControlBase) ToggleVisible() bool {
	visible := w32.IsWindowVisible(cba.hwnd)
	if visible {
		cba.Hide()
	} else {
		cba.Show()
	}
	return !visible
}

func (cba *ControlBase) ContextMenu() *MenuItem {
	return cba.contextMenu
}

func (cba *ControlBase) SetContextMenu(menu *MenuItem) {
	cba.contextMenu = menu
}

func (cba *ControlBase) Bounds() *Rect {
	rect := w32.GetWindowRect(cba.hwnd)
	if cba.isForm {
		return &Rect{*rect}
	}

	return ScreenToClientRect(cba.hwnd, rect)
}

func (cba *ControlBase) ClientRect() *Rect {
	rect := w32.GetClientRect(cba.hwnd)
	return ScreenToClientRect(cba.hwnd, rect)
}
func (cba *ControlBase) ClientWidth() int {
	rect := w32.GetClientRect(cba.hwnd)
	return int(rect.Right - rect.Left)
}

func (cba *ControlBase) ClientHeight() int {
	rect := w32.GetClientRect(cba.hwnd)
	return int(rect.Bottom - rect.Top)
}

func (cba *ControlBase) Show() {
	// WindowPos is used with HWND_TOPMOST to guarantee bring our app on top
	// force set our main window on top
	w32.SetWindowPos(
		cba.hwnd,
		w32.HWND_TOPMOST,
		0, 0, 0, 0,
		w32.SWP_SHOWWINDOW|w32.SWP_NOSIZE|w32.SWP_NOMOVE,
	)
	// remove topmost to allow normal windows manipulations
	w32.SetWindowPos(
		cba.hwnd,
		w32.HWND_NOTOPMOST,
		0, 0, 0, 0,
		w32.SWP_SHOWWINDOW|w32.SWP_NOSIZE|w32.SWP_NOMOVE,
	)
	// put main window on tops foreground
	w32.SetForegroundWindow(cba.hwnd)
}

func (cba *ControlBase) Hide() {
	w32.ShowWindow(cba.hwnd, w32.SW_HIDE)
}

func (cba *ControlBase) Enabled() bool {
	return w32.IsWindowEnabled(cba.hwnd)
}

func (cba *ControlBase) SetEnabled(b bool) {
	w32.EnableWindow(cba.hwnd, b)
}

func (cba *ControlBase) SetFocus() {
	w32.SetFocus(cba.hwnd)
}

func (cba *ControlBase) Invalidate(erase bool) {
	// pRect := w32.GetClientRect(cba.hwnd)
	// if cba.isForm {
	// 	w32.InvalidateRect(cba.hwnd, pRect, erase)
	// } else {
	// 	rc := ScreenToClientRect(cba.hwnd, pRect)
	// 	w32.InvalidateRect(cba.hwnd, rc.GetW32Rect(), erase)
	// }
	w32.InvalidateRect(cba.hwnd, nil, erase)
}

func (cba *ControlBase) Parent() Controller {
	return cba.parent
}

func (cba *ControlBase) SetParent(parent Controller) {
	cba.parent = parent
}

func (cba *ControlBase) Font() *Font {
	return cba.font
}

func (cba *ControlBase) SetFont(font *Font) {
	w32.SendMessage(cba.hwnd, w32.WM_SETFONT, uintptr(font.hfont), 1)
	cba.font = font
}

func (cba *ControlBase) EnableDragAcceptFiles(b bool) {
	w32.DragAcceptFiles(cba.hwnd, b)
}

func (cba *ControlBase) InvokeRequired() bool {
	if cba.hwnd == 0 {
		return false
	}

	windowThreadId, _ := w32.GetWindowThreadProcessId(cba.hwnd)
	currentThreadId := w32.GetCurrentThreadId()

	return windowThreadId != currentThreadId
}

func (cba *ControlBase) Invoke(f func()) {
	if cba.tryInvokeOnCurrentGoRoutine(f) {
		return
	}

	cba.m.Lock()
	cba.dispatchq = append(cba.dispatchq, f)
	cba.m.Unlock()
	w32.PostMessage(cba.hwnd, wmInvokeCallback, 0, 0)
}

func (cba *ControlBase) PreTranslateMessage(msg *w32.MSG) bool {
	if msg.Message == w32.WM_GETDLGCODE {
		println("pretranslate, WM_GETDLGCODE")
	}
	return false
}

// Events
func (cba *ControlBase) OnCreate() *EventManager {
	return &cba.onCreate
}

func (cba *ControlBase) OnClose() *EventManager {
	return &cba.onClose
}

func (cba *ControlBase) OnKillFocus() *EventManager {
	return &cba.onKillFocus
}

func (cba *ControlBase) OnSetFocus() *EventManager {
	return &cba.onSetFocus
}

func (cba *ControlBase) OnDropFiles() *EventManager {
	return &cba.onDropFiles
}

func (cba *ControlBase) OnLBDown() *EventManager {
	return &cba.onLBDown
}

func (cba *ControlBase) OnLBUp() *EventManager {
	return &cba.onLBUp
}

func (cba *ControlBase) OnLBDbl() *EventManager {
	return &cba.onLBDbl
}

func (cba *ControlBase) OnMBDown() *EventManager {
	return &cba.onMBDown
}

func (cba *ControlBase) OnMBUp() *EventManager {
	return &cba.onMBUp
}

func (cba *ControlBase) OnRBDown() *EventManager {
	return &cba.onRBDown
}

func (cba *ControlBase) OnRBUp() *EventManager {
	return &cba.onRBUp
}

func (cba *ControlBase) OnRBDbl() *EventManager {
	return &cba.onRBDbl
}

func (cba *ControlBase) OnMouseMove() *EventManager {
	return &cba.onMouseMove
}

func (cba *ControlBase) OnMouseHover() *EventManager {
	return &cba.onMouseHover
}

func (cba *ControlBase) OnMouseLeave() *EventManager {
	return &cba.onMouseLeave
}

func (cba *ControlBase) OnPaint() *EventManager {
	return &cba.onPaint
}

func (cba *ControlBase) OnSize() *EventManager {
	return &cba.onSize
}

func (cba *ControlBase) OnKeyUp() *EventManager {
	return &cba.onKeyUp
}

func (cba *ControlBase) scaleWithWindowDPI(width, height int) (int, int) {
	dpix, dpiy := cba.GetWindowDPI()
	scaledWidth := ScaleWithDPI(width, dpix)
	scaledHeight := ScaleWithDPI(height, dpiy)

	return scaledWidth, scaledHeight
}

func (cba *ControlBase) scaleToDefaultDPI(width, height int) (int, int) {
	dpix, dpiy := cba.GetWindowDPI()
	scaledWidth := ScaleToDefaultDPI(width, dpix)
	scaledHeight := ScaleToDefaultDPI(height, dpiy)

	return scaledWidth, scaledHeight
}

func (cba *ControlBase) tryInvokeOnCurrentGoRoutine(f func()) bool {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if cba.InvokeRequired() {
		return false
	}
	f()
	return true
}

func (cba *ControlBase) invokeCallbacks() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if cba.InvokeRequired() {
		panic("InvokeCallbacks must always be called on the window thread")
	}

	cba.m.Lock()
	q := append([]func(){}, cba.dispatchq...)
	cba.dispatchq = []func(){}
	cba.m.Unlock()
	for _, v := range q {
		v()
	}
}
