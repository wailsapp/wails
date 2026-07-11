//go:build windows

package application

import (
	"slices"
	"sync"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/w32"
	"github.com/wailsapp/wails/v3/internal/webview2/pkg/edge"
)

type nonClientHitTestState struct {
	mu      sync.RWMutex
	regions []nonClientHitTestRegion
}

func (s *nonClientHitTestState) set(regions []nonClientHitTestRegion) {
	s.mu.Lock()
	if len(regions) == 0 {
		s.regions = nil
	} else {
		s.regions = slices.Clone(regions)
	}
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

func (w *windowsWebviewWindow) applyCompositionCursor(cursor edge.HCURSOR, systemCursorID uint32) {
	hcursor := w32.HCURSOR(cursor)
	if hcursor == 0 && systemCursorID != 0 {
		hcursor = w32.LoadCursorWithResourceID(0, uint16(systemCursorID))
	}
	if hcursor == 0 {
		return
	}

	w.compositionCursor = hcursor
	w32.SetCursor(hcursor)
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

		w.trackCompositionMouseLeave(true)

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
		// Windows can emit a spurious NCMOUSELEAVE right after a forwarded
		// non-client button press. Suppress it until the captured mouse move
		// path below can decide whether the pointer really left the active button.
		if w32.GetCapture() == w.hwnd && w.activeNonClientButton != 0 {
			return 0, true
		}

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
	case w32.WM_NCRBUTTONUP:
		if wparam != w32.HTCAPTION && wparam != w32.HTSYSMENU {
			return 0, false
		}

		screenX := int(w32.GET_X_LPARAM(lparam))
		screenY := int(w32.GET_Y_LPARAM(lparam))

		if hitTest, ok := w.nonClientHitTestFromScreen(screenX, screenY); !ok || hitTest != w32.HTCAPTION {
			return 0, false
		}

		if w.showCaptionSystemMenu(screenX, screenY) {
			return 0, true
		}
	}

	return 0, false
}

func (w *windowsWebviewWindow) nonClientHitTestFromScreen(screenX, screenY int) (uintptr, bool) {
	regions := w.nonClientHitTest.snapshot()
	// Later frontend regions should win when rectangles overlap, matching the
	// DOM/CSS order used to collect them (for example caption buttons over a
	// caption drag area).
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

	dpi, _ := w.DPI()
	borderX := systemMetricForDPI(w32.SM_CXSIZEFRAME, dpi) + systemMetricForDPI(w32.SM_CXPADDEDBORDER, dpi)
	if borderX < 1 {
		borderX = 1
	}
	borderY := systemMetricForDPI(w32.SM_CYSIZEFRAME, dpi) + systemMetricForDPI(w32.SM_CXPADDEDBORDER, dpi)
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

func systemMetricForDPI(index int, dpi w32.UINT) int {
	if dpi != 0 && w32.HasGetSystemMetricsForDpiFunc() {
		return w32.GetSystemMetricsForDpi(index, dpi)
	}
	return w32.GetSystemMetrics(index)
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

	// Remember which caption button started the press. During capture, Windows
	// sends movement as WM_MOUSEMOVE, so later moves must be compared against
	// this original hit-test rather than any caption button under the cursor.
	if msg == w32.WM_NCLBUTTONDOWN || msg == w32.WM_NCLBUTTONDBLCLK {
		w.activeNonClientButton = wparam
		w.activeNonClientButtonHovered = true
	}

	w.updateCompositionMouseCapture(msg)

	return true
}

func (w *windowsWebviewWindow) showCaptionSystemMenu(screenX, screenY int) bool {
	menu := w32.GetSystemMenu(w.hwnd, false)
	if menu == 0 {
		return false
	}

	w32.SetForegroundWindow(w.hwnd)
	command := w32.TrackPopupMenuCommand(
		menu,
		w32.TPM_LEFTALIGN|w32.TPM_TOPALIGN|w32.TPM_RIGHTBUTTON,
		int32(screenX),
		int32(screenY),
		w.hwnd,
		nil,
	)
	w32.PostMessage(w.hwnd, w32.WM_NULL, 0, 0)
	if command == 0 {
		return true
	}

	w32.SendMessage(w.hwnd, w32.WM_SYSCOMMAND, command, 0)
	return true
}

func (w *windowsWebviewWindow) routeCompositionMouseInput(msg uint32, wparam, lparam uintptr) bool {
	eventKind, ok := webviewCompositionMouseEventKind(msg)
	if !ok {
		return false
	}

	if msg == w32.WM_MOUSEMOVE {
		w.trackCompositionMouseLeave(false)
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

	// While a forwarded caption button owns capture, WebView still receives
	// normal mouse moves. Only forward those moves while the cursor remains over
	// the same button; otherwise synthesize leave so pressed styling is cleared.
	if msg == w32.WM_MOUSEMOVE && w.activeNonClientButton != 0 && w32.GetCapture() == w.hwnd {
		if !w.updateActiveNonClientButtonHover() {
			return true
		}
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

	w.updateCompositionMouseCapture(msg)

	return true
}

func (w *windowsWebviewWindow) trackCompositionMouseLeave(nonClient bool) {
	flags := uint32(w32.TME_LEAVE)
	if nonClient {
		flags |= w32.TME_NONCLIENT
	}

	w32.TrackMouseEvent(&w32.TRACKMOUSEEVENT{
		CbSize:    uint32(unsafe.Sizeof(w32.TRACKMOUSEEVENT{})),
		DwFlags:   flags,
		HwndTrack: w.hwnd,
	})
}

func (w *windowsWebviewWindow) compositionMouseClientPoint(msg uint32, lparam uintptr) (int, int, bool) {
	switch msg {
	case w32.WM_MOUSELEAVE:
		return 0, 0, true
	case w32.WM_MOUSEWHEEL, w32.WM_MOUSEHWHEEL:
		return w32.ScreenToClient(w.hwnd, int(w32.GET_X_LPARAM(lparam)), int(w32.GET_Y_LPARAM(lparam)))
	default:
		return int(w32.GET_X_LPARAM(lparam)), int(w32.GET_Y_LPARAM(lparam)), true
	}
}

func (w *windowsWebviewWindow) updateCompositionMouseCapture(msg uint32) {
	if isCompositionMouseButtonDown(msg) {
		if w32.GetCapture() != w.hwnd {
			w32.SetCapture(w.hwnd)
		}
		return
	}

	if isCompositionMouseButtonUp(msg) && w32.GetCapture() == w.hwnd {
		w32.ReleaseCapture()

		// Button release ends the synthetic caption-button interaction.
		w.activeNonClientButton = 0
		w.activeNonClientButtonHovered = false
	}
}

// updateActiveNonClientButtonHover returns whether the current move should be
// forwarded to WebView. It sends a single LEAVE when capture moves off the
// originally pressed caption button, matching native caption-button behavior.
func (w *windowsWebviewWindow) updateActiveNonClientButtonHover() bool {
	screenX, screenY, ok := w32.GetCursorPos()
	if !ok {
		return true
	}

	hitTest, ok := w.nonClientHitTestFromScreen(screenX, screenY)
	hovered := ok && hitTest == w.activeNonClientButton
	if hovered {
		w.activeNonClientButtonHovered = true
		return true
	}

	if w.activeNonClientButtonHovered {
		_ = w.chromium.SendMouseInput(
			edge.COREWEBVIEW2_MOUSE_EVENT_KIND_LEAVE,
			edge.COREWEBVIEW2_MOUSE_EVENT_VIRTUAL_KEYS_NONE,
			0,
			0,
			0,
		)
		w.activeNonClientButtonHovered = false
	}

	return false
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

func isCompositionMouseButtonDown(msg uint32) bool {
	switch msg {
	case w32.WM_LBUTTONDOWN,
		w32.WM_NCLBUTTONDOWN,
		w32.WM_NCLBUTTONDBLCLK,
		w32.WM_RBUTTONDOWN,
		w32.WM_MBUTTONDOWN,
		w32.WM_XBUTTONDOWN:
		return true
	default:
		return false
	}
}

func isCompositionMouseButtonUp(msg uint32) bool {
	switch msg {
	case w32.WM_LBUTTONUP,
		w32.WM_NCLBUTTONUP,
		w32.WM_RBUTTONUP,
		w32.WM_MBUTTONUP,
		w32.WM_XBUTTONUP:
		return true
	default:
		return false
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
	default:
		return 0, false
	}
}
