//go:build windows

package application

import (
	"sync"

	"github.com/wailsapp/wails/v3/pkg/w32"
	"github.com/wailsapp/wails/webview2/pkg/edge"
)

type nonClientHitTestState struct {
	mu      sync.RWMutex
	regions []nonClientHitTestRegion
}

func (s *nonClientHitTestState) set(regions []nonClientHitTestRegion) {
	s.mu.Lock()
	s.regions = regions
	s.mu.Unlock()
}

func (s *nonClientHitTestState) snapshot() []nonClientHitTestRegion {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.regions) == 0 {
		return nil
	}
	return append([]nonClientHitTestRegion(nil), s.regions...)
}

func (w *windowsWebviewWindow) setNonClientHitTestRegions(regions []nonClientHitTestRegion) {
	w.nonClientHitTest.set(regions)
}

func (w *windowsWebviewWindow) routeNonClientInput(msg uint32, wparam, lparam uintptr) (uintptr, bool) {
	switch msg {
	case w32.WM_NCHITTEST:
		if hitTest, handled := w32.DwmDefWindowProc(w.hwnd, msg, wparam, lparam); handled {
			if hitTest != w32.HTCLIENT && hitTest != w32.HTNOWHERE {
				return hitTest, true
			}
		}

		screenX := int(w32.GET_X_LPARAM(lparam))
		screenY := int(w32.GET_Y_LPARAM(lparam))

		if hitTest, ok := w.resizeBorderHitTest(screenX, screenY); ok {
			return hitTest, true
		}

		return w.nonClientHitTestFromScreen(screenX, screenY)
	case w32.WM_NCMOUSEMOVE:
		screenX := int(w32.GET_X_LPARAM(lparam))
		screenY := int(w32.GET_Y_LPARAM(lparam))

		hitTest, ok := w.nonClientHitTestFromScreen(screenX, screenY)
		if !ok || hitTest != wparam {
			return 0, false
		}

		clientX, clientY, ok := w32.ScreenToClient(
			w.hwnd,
			int(w32.GET_X_LPARAM(lparam)),
			int(w32.GET_Y_LPARAM(lparam)),
		)
		if !ok {
			return 0, false
		}

		_ = w.chromium.SendMouseInput(
			edge.COREWEBVIEW2_MOUSE_EVENT_KIND_MOVE,
			edge.COREWEBVIEW2_MOUSE_EVENT_VIRTUAL_KEYS_NONE,
			0,
			clientX,
			clientY,
		)

		return w32.DwmDefWindowProc(w.hwnd, msg, wparam, lparam)
	case w32.WM_NCMOUSELEAVE:
		_ = w.chromium.SendMouseInput(
			edge.COREWEBVIEW2_MOUSE_EVENT_KIND_LEAVE,
			edge.COREWEBVIEW2_MOUSE_EVENT_VIRTUAL_KEYS_NONE,
			0,
			0,
			0,
		)
	case w32.WM_NCLBUTTONDOWN, w32.WM_NCLBUTTONUP, w32.WM_NCLBUTTONDBLCLK:
		switch wparam {
		case w32.HTMINBUTTON, w32.HTMAXBUTTON, w32.HTCLOSE:
		default:
			return 0, false
		}

		if w.forwardFrontendNonClientButtonInput(msg, wparam, lparam) {
			return 0, true
		}
	}

	return 0, false
}

func (w *windowsWebviewWindow) nonClientHitTestFromScreen(screenX, screenY int) (uintptr, bool) {
	regions := w.nonClientHitTest.snapshot()
	for i := len(regions) - 1; i >= 0; i-- {
		r := regions[i]

		left, top := w32.ClientToScreen(w.hwnd, r.Left, r.Top)
		right, bottom := w32.ClientToScreen(w.hwnd, r.Right, r.Bottom)

		// outside app
		if screenX < left || screenX >= right || screenY < top || screenY >= bottom {
			continue
		}

		switch r.Kind {
		case nonClientHitTestKindMaximize:
			return w32.HTMAXBUTTON, true
		case nonClientHitTestKindCaption:
			return w32.HTCAPTION, true
		case nonClientHitTestKindMinimize:
			return w32.HTMINBUTTON, true
		case nonClientHitTestKindClose:
			return w32.HTCLOSE, true
		}
	}

	// fallback to default non-client region hit test (app-region)
	if !w.chromium.NonClientRegionSupportEnabled {
		return 0, false
	}

	clientX, clientY, ok := w32.ScreenToClient(w.hwnd, screenX, screenY)
	if !ok {
		return 0, false
	}

	region, handled, err := w.chromium.GetNonClientRegionAtPoint(int32(clientX), int32(clientY))
	if err != nil || !handled {
		return 0, false
	}

	switch region {
	case edge.COREWEBVIEW2_NON_CLIENT_REGION_KIND_CAPTION:
		return w32.HTCAPTION, true
	case edge.COREWEBVIEW2_NON_CLIENT_REGION_KIND_MINIMIZE:
		return w32.HTMINBUTTON, true
	case edge.COREWEBVIEW2_NON_CLIENT_REGION_KIND_MAXIMIZE:
		return w32.HTMAXBUTTON, true
	case edge.COREWEBVIEW2_NON_CLIENT_REGION_KIND_CLOSE:
		return w32.HTCLOSE, true
	case edge.COREWEBVIEW2_NON_CLIENT_REGION_KIND_CLIENT,
		edge.COREWEBVIEW2_NON_CLIENT_REGION_KIND_NOWHERE:
		return 0, false
	default:
		return 0, false
	}
}

func (w *windowsWebviewWindow) resizeBorderHitTest(screenX, screenY int) (uintptr, bool) {
	if w.parent.options.DisableResize || w.isMaximised() || w.isFullscreen() {
		return 0, false
	}

	rect := w32.GetWindowRect(w.hwnd)
	width := int(rect.Right - rect.Left)
	height := int(rect.Bottom - rect.Top)
	if width <= 0 || height <= 0 {
		return 0, false
	}

	x := screenX - int(rect.Left)
	y := screenY - int(rect.Top)
	if x < 0 || y < 0 || x >= width || y >= height {
		return 0, false
	}

	borderX := int(w.resizeBorderWidth) + w32.GetSystemMetrics(w32.SM_CXPADDEDBORDER)
	if borderX < 1 {
		borderX = 1
	}
	borderY := int(w.resizeBorderHeight) + w32.GetSystemMetrics(w32.SM_CXPADDEDBORDER)
	if borderY < 1 {
		borderY = 1
	}

	left := x < borderX
	right := x >= width-borderX
	top := y < borderY
	bottom := y >= height-borderY

	switch {
	case top && left:
		return w32.HTTOPLEFT, true
	case top && right:
		return w32.HTTOPRIGHT, true
	case bottom && left:
		return w32.HTBOTTOMLEFT, true
	case bottom && right:
		return w32.HTBOTTOMRIGHT, true
	case top:
		return w32.HTTOP, true
	case bottom:
		return w32.HTBOTTOM, true
	case left:
		return w32.HTLEFT, true
	case right:
		return w32.HTRIGHT, true
	default:
		return 0, false
	}
}

func (w *windowsWebviewWindow) forwardFrontendNonClientButtonInput(msg uint32, wparam, lparam uintptr) bool {
	screenX := int(w32.GET_X_LPARAM(lparam))
	screenY := int(w32.GET_Y_LPARAM(lparam))

	hitTest, ok := w.nonClientHitTestFromScreen(screenX, screenY)
	if !ok || hitTest != wparam {
		return false
	}

	clientX, clientY, ok := w32.ScreenToClient(w.hwnd, screenX, screenY)
	if !ok {
		return false
	}

	eventKind, virtualKeys, ok := nonClientLeftButtonMouseEvent(msg)
	if !ok {
		return false
	}

	if err := w.chromium.SendMouseInput(
		eventKind,
		virtualKeys,
		0,
		clientX,
		clientY,
	); err != nil {
		return false
	}

	return true
}

func (w *windowsWebviewWindow) routeCompositionMouseInput(msg uint32, wparam, lparam uintptr) bool {
	eventKind, ok := webviewCompositionMouseEventKind(msg)
	if !ok {
		return false
	}

	clientX, clientY, ok := w.compositionMouseClientPoint(msg, lparam)
	if !ok {
		return false
	}

	var mouseData uint32
	switch msg {
	case w32.WM_MOUSEWHEEL, w32.WM_MOUSEHWHEEL:
		mouseData = uint32(int16(wparam >> 16))
	case w32.WM_XBUTTONDOWN, w32.WM_XBUTTONUP, w32.WM_XBUTTONDBLCLK:
		mouseData = uint32(wparam >> 16)
	}

	if err := w.chromium.SendMouseInput(
		eventKind,
		edge.COREWEBVIEW2_MOUSE_EVENT_VIRTUAL_KEYS(uint32(wparam)&0xffff),
		mouseData,
		clientX,
		clientY,
	); err != nil {
		return false
	}

	return true
}

func (w *windowsWebviewWindow) compositionMouseClientPoint(msg uint32, lparam uintptr) (int, int, bool) {
	switch msg {
	case w32.WM_MOUSELEAVE:
		return 0, 0, true
	case w32.WM_MOUSEWHEEL, w32.WM_MOUSEHWHEEL,
		w32.WM_NCRBUTTONDOWN, w32.WM_NCRBUTTONUP:
		return w32.ScreenToClient(w.hwnd, int(w32.GET_X_LPARAM(lparam)), int(w32.GET_Y_LPARAM(lparam)))
	default:
		return int(w32.GET_X_LPARAM(lparam)), int(w32.GET_Y_LPARAM(lparam)), true
	}
}

func nonClientLeftButtonMouseEvent(msg uint32) (edge.COREWEBVIEW2_MOUSE_EVENT_KIND, edge.COREWEBVIEW2_MOUSE_EVENT_VIRTUAL_KEYS, bool) {
	switch msg {
	case w32.WM_NCLBUTTONDOWN:
		return edge.COREWEBVIEW2_MOUSE_EVENT_KIND_LEFT_BUTTON_DOWN,
			edge.COREWEBVIEW2_MOUSE_EVENT_VIRTUAL_KEYS_LEFT_BUTTON,
			true
	case w32.WM_NCLBUTTONUP:
		return edge.COREWEBVIEW2_MOUSE_EVENT_KIND_LEFT_BUTTON_UP,
			edge.COREWEBVIEW2_MOUSE_EVENT_VIRTUAL_KEYS_NONE,
			true
	case w32.WM_NCLBUTTONDBLCLK:
		return edge.COREWEBVIEW2_MOUSE_EVENT_KIND_LEFT_BUTTON_DOUBLE_CLICK,
			edge.COREWEBVIEW2_MOUSE_EVENT_VIRTUAL_KEYS_LEFT_BUTTON,
			true
	default:
		return 0, 0, false
	}
}

func webviewCompositionMouseEventKind(msg uint32) (edge.COREWEBVIEW2_MOUSE_EVENT_KIND, bool) {
	switch msg {
	case w32.WM_MOUSELEAVE:
		return edge.COREWEBVIEW2_MOUSE_EVENT_KIND_LEAVE, true
	case w32.WM_MOUSEMOVE:
		return edge.COREWEBVIEW2_MOUSE_EVENT_KIND_MOVE, true
	case w32.WM_LBUTTONDOWN:
		return edge.COREWEBVIEW2_MOUSE_EVENT_KIND_LEFT_BUTTON_DOWN, true
	case w32.WM_LBUTTONUP:
		return edge.COREWEBVIEW2_MOUSE_EVENT_KIND_LEFT_BUTTON_UP, true
	case w32.WM_LBUTTONDBLCLK:
		return edge.COREWEBVIEW2_MOUSE_EVENT_KIND_LEFT_BUTTON_DOUBLE_CLICK, true
	case w32.WM_RBUTTONDOWN:
		return edge.COREWEBVIEW2_MOUSE_EVENT_KIND_RIGHT_BUTTON_DOWN, true
	case w32.WM_RBUTTONUP:
		return edge.COREWEBVIEW2_MOUSE_EVENT_KIND_RIGHT_BUTTON_UP, true
	case w32.WM_RBUTTONDBLCLK:
		return edge.COREWEBVIEW2_MOUSE_EVENT_KIND_RIGHT_BUTTON_DOUBLE_CLICK, true
	case w32.WM_MBUTTONDOWN:
		return edge.COREWEBVIEW2_MOUSE_EVENT_KIND_MIDDLE_BUTTON_DOWN, true
	case w32.WM_MBUTTONUP:
		return edge.COREWEBVIEW2_MOUSE_EVENT_KIND_MIDDLE_BUTTON_UP, true
	case w32.WM_MBUTTONDBLCLK:
		return edge.COREWEBVIEW2_MOUSE_EVENT_KIND_MIDDLE_BUTTON_DOUBLE_CLICK, true
	case w32.WM_XBUTTONDOWN:
		return edge.COREWEBVIEW2_MOUSE_EVENT_KIND_X_BUTTON_DOWN, true
	case w32.WM_XBUTTONUP:
		return edge.COREWEBVIEW2_MOUSE_EVENT_KIND_X_BUTTON_UP, true
	case w32.WM_XBUTTONDBLCLK:
		return edge.COREWEBVIEW2_MOUSE_EVENT_KIND_X_BUTTON_DOUBLE_CLICK, true
	case w32.WM_MOUSEWHEEL:
		return edge.COREWEBVIEW2_MOUSE_EVENT_KIND_WHEEL, true
	case w32.WM_MOUSEHWHEEL:
		return edge.COREWEBVIEW2_MOUSE_EVENT_KIND_HORIZONTAL_WHEEL, true
	case w32.WM_NCRBUTTONDOWN:
		return edge.COREWEBVIEW2_MOUSE_EVENT_KIND_NON_CLIENT_RIGHT_BUTTON_DOWN, true
	case w32.WM_NCRBUTTONUP:
		return edge.COREWEBVIEW2_MOUSE_EVENT_KIND_NON_CLIENT_RIGHT_BUTTON_UP, true
	default:
		return 0, false
	}
}
