//go:build windows

package application

import (
	"strings"
	"sync"
	"unsafe"

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

type snapLayoutHoverState struct {
	tracking bool
}

type compositionMouseInputState struct {
	tracking  bool
	capturing bool
}

func (w *windowsWebviewWindow) setNonClientHitTestRegions(regions []nonClientHitTestRegion) {
	w.nonClientHitTest.set(regions)
}

func (w *windowsWebviewWindow) routeNonClientInput(msg uint32, wparam, lparam uintptr) (bool, uintptr) {
	switch msg {
	case w32.WM_NCHITTEST:
		if hitTest, ok := w.handleNonClientHitTest(wparam, lparam); ok {
			return true, hitTest
		}
	case w32.WM_NCMOUSEMOVE, w32.WM_NCMOUSEHOVER, w32.WM_NCMOUSELEAVE:
		if result, handled := w.handleSnapLayoutHoverMessage(msg, wparam, lparam); handled {
			return true, result
		}
	case w32.WM_NCLBUTTONDOWN, w32.WM_NCLBUTTONUP, w32.WM_NCLBUTTONDBLCLK:
		if w.forwardFrontendNonClientButtonInput(msg, wparam, lparam) {
			return true, 0
		}
	}

	return false, 0
}

func (w *windowsWebviewWindow) handleNonClientHitTest(wparam, lparam uintptr) (uintptr, bool) {
	if hitTest, handled := w32.DwmDefWindowProc(w.hwnd, w32.WM_NCHITTEST, wparam, lparam); handled {
		if authoritativeDwmHitTest(hitTest) {
			return hitTest, true
		}
	}

	screenX := int(w32.GET_X_LPARAM(lparam))
	screenY := int(w32.GET_Y_LPARAM(lparam))
	return w.nonClientHitTestFromScreen(screenX, screenY)
}

func authoritativeDwmHitTest(hitTest uintptr) bool {
	return hitTest != w32.HTCLIENT && hitTest != w32.HTNOWHERE
}

func (w *windowsWebviewWindow) nonClientHitTestFromScreen(screenX, screenY int) (uintptr, bool) {
	if hitTest, ok := w.resizeBorderHitTest(screenX, screenY); ok {
		return hitTest, true
	}
	if hitTest, ok := w.manualNonClientHitTestFromScreen(screenX, screenY); ok {
		return hitTest, true
	}
	return w.webviewNonClientRegionHitTestFromScreen(screenX, screenY)
}

func (w *windowsWebviewWindow) frontendNonClientHitTestFromScreen(screenX, screenY int) (uintptr, bool) {
	if hitTest, ok := w.manualNonClientHitTestFromScreen(screenX, screenY); ok {
		return hitTest, true
	}
	return w.webviewNonClientRegionHitTestFromScreen(screenX, screenY)
}

func (w *windowsWebviewWindow) manualNonClientHitTestFromScreen(screenX, screenY int) (uintptr, bool) {
	regions := w.nonClientHitTest.snapshot()
	for i := len(regions) - 1; i >= 0; i-- {
		region := regions[i]
		if !w.screenPointInClientRegion(screenX, screenY, region) {
			continue
		}
		if hitTest := nonClientRegionKindToHitTest(region.Kind); hitTest != 0 {
			return hitTest, true
		}
	}
	return 0, false
}

func (w *windowsWebviewWindow) screenPointInClientRegion(screenX, screenY int, region nonClientHitTestRegion) bool {
	left, top := w32.ClientToScreen(w.hwnd, region.Left, region.Top)
	right, bottom := w32.ClientToScreen(w.hwnd, region.Right, region.Bottom)
	return screenX >= left && screenX < right && screenY >= top && screenY < bottom
}

func nonClientRegionKindToHitTest(kind nonClientHitTestKind) uintptr {
	switch nonClientHitTestKind(strings.ToLower(string(kind))) {
	case "", nonClientHitTestKindMaximize:
		return w32.HTMAXBUTTON
	case nonClientHitTestKindCaption:
		return w32.HTCAPTION
	case nonClientHitTestKindMinimize:
		return w32.HTMINBUTTON
	case nonClientHitTestKindClose:
		return w32.HTCLOSE
	default:
		return 0
	}
}

func (w *windowsWebviewWindow) webviewNonClientRegionHitTestFromScreen(screenX, screenY int) (uintptr, bool) {
	if !w.parent.options.Windows.NonClientRegionSupport || w.chromium == nil {
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
	return webviewNonClientRegionKindToHitTest(region)
}

func webviewNonClientRegionKindToHitTest(region edge.COREWEBVIEW2_NON_CLIENT_REGION_KIND) (uintptr, bool) {
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

	if w.pointInDwmCaptionButtonBounds(screenX, screenY) {
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
	borderY := int(w.resizeBorderHeight) + w32.GetSystemMetrics(w32.SM_CXPADDEDBORDER)
	if borderX < 1 {
		borderX = 1
	}
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

func (w *windowsWebviewWindow) handleSnapLayoutHoverMessage(msg uint32, wparam, lparam uintptr) (uintptr, bool) {
	switch msg {
	case w32.WM_NCMOUSEMOVE:
		if w.isFrontendCaptionButtonHit(wparam, lparam) {
			if wparam == w32.HTMAXBUTTON {
				w.ensureSnapHoverTracking()
			} else {
				w.cancelSnapHoverTracking()
			}
			w.forwardFrontendNonClientHover(msg, wparam, lparam)
		} else {
			w.cancelSnapHoverTracking()
		}
	case w32.WM_NCMOUSEHOVER:
		w.snapHover.tracking = false
	case w32.WM_NCMOUSELEAVE:
		w.snapHover.tracking = false
		w.forwardFrontendNonClientHover(msg, wparam, lparam)
	default:
		return 0, false
	}

	result, handled := w32.DwmDefWindowProc(w.hwnd, msg, wparam, lparam)
	if handled {
		return result, true
	}
	return 0, false
}

func (w *windowsWebviewWindow) isFrontendCaptionButtonHit(expectedHitTest uintptr, lparam uintptr) bool {
	if !frontendCaptionButtonHitTest(expectedHitTest) {
		return false
	}

	screenX := int(w32.GET_X_LPARAM(lparam))
	screenY := int(w32.GET_Y_LPARAM(lparam))
	hitTest, ok := w.frontendNonClientHitTestFromScreen(screenX, screenY)
	return ok && hitTest == expectedHitTest
}

func (w *windowsWebviewWindow) ensureSnapHoverTracking() {
	if w.snapHover.tracking {
		return
	}
	w.snapHover.tracking = true
	w32.TrackMouseEvent(&w32.TRACKMOUSEEVENT{
		CbSize:      uint32(unsafe.Sizeof(w32.TRACKMOUSEEVENT{})),
		DwFlags:     w32.TME_HOVER | w32.TME_LEAVE | w32.TME_NONCLIENT,
		HwndTrack:   w.hwnd,
		DwHoverTime: w32.HOVER_DEFAULT,
	})
}

func (w *windowsWebviewWindow) cancelSnapHoverTracking() {
	if !w.snapHover.tracking {
		return
	}
	w.snapHover.tracking = false
	w32.TrackMouseEvent(&w32.TRACKMOUSEEVENT{
		CbSize:    uint32(unsafe.Sizeof(w32.TRACKMOUSEEVENT{})),
		DwFlags:   w32.TME_HOVER | w32.TME_LEAVE | w32.TME_NONCLIENT | w32.TME_CANCEL,
		HwndTrack: w.hwnd,
	})
}

func (w *windowsWebviewWindow) forwardFrontendNonClientHover(msg uint32, wparam, lparam uintptr) {
	if !w.canSendCompositionMouseInput() {
		return
	}

	switch msg {
	case w32.WM_NCMOUSEMOVE:
		if !w.isFrontendCaptionButtonHit(wparam, lparam) {
			return
		}
		w.sendCompositionMouseMoveFromNonClientLParam(lparam)
	case w32.WM_NCMOUSELEAVE:
		_ = w.sendCompositionMouseLeave()
	}
}

func (w *windowsWebviewWindow) sendCompositionMouseMoveFromNonClientLParam(lparam uintptr) bool {
	clientX, clientY, ok := w32.ScreenToClient(
		w.hwnd,
		int(w32.GET_X_LPARAM(lparam)),
		int(w32.GET_Y_LPARAM(lparam)),
	)
	if !ok {
		return false
	}

	return w.sendCompositionMouseInputAtClientPoint(
		edge.COREWEBVIEW2_MOUSE_EVENT_KIND_MOVE,
		edge.COREWEBVIEW2_MOUSE_EVENT_VIRTUAL_KEYS_NONE,
		0,
		clientX,
		clientY,
	) == nil
}

func (w *windowsWebviewWindow) forwardFrontendNonClientButtonInput(msg uint32, wparam, lparam uintptr) bool {
	if !w.canSendCompositionMouseInput() || !frontendCaptionButtonHitTest(wparam) {
		return false
	}

	screenX := int(w32.GET_X_LPARAM(lparam))
	screenY := int(w32.GET_Y_LPARAM(lparam))
	hitTest, ok := w.frontendNonClientHitTestFromScreen(screenX, screenY)
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

	if msg == w32.WM_NCLBUTTONDOWN {
		w.focus()
		if !w.compositionInput.capturing {
			w.compositionInput.capturing = true
			w32.SetCapture(w.hwnd)
		}
	} else if msg == w32.WM_NCLBUTTONUP {
		if w.compositionInput.capturing {
			w.compositionInput.capturing = false
			w32.ReleaseCapture()
		}
	}

	if err := w.sendCompositionMouseInputAtClientPoint(eventKind, virtualKeys, 0, clientX, clientY); err != nil {
		if msg == w32.WM_NCLBUTTONDOWN && w.compositionInput.capturing {
			w.compositionInput.capturing = false
			w32.ReleaseCapture()
		}
		return false
	}
	return true
}

func frontendCaptionButtonHitTest(hitTest uintptr) bool {
	switch hitTest {
	case w32.HTMINBUTTON, w32.HTMAXBUTTON, w32.HTCLOSE:
		return true
	default:
		return false
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

func (w *windowsWebviewWindow) dwmCaptionButtonBounds() (w32.RECT, bool) {
	var rect w32.RECT
	hr := w32.DwmGetWindowAttribute(
		w.hwnd,
		w32.DWMWA_CAPTION_BUTTON_BOUNDS,
		unsafe.Pointer(&rect),
		unsafe.Sizeof(rect),
	)
	return rect, hr == 0
}

func (w *windowsWebviewWindow) pointInDwmCaptionButtonBounds(screenX, screenY int) bool {
	rect, ok := w.dwmCaptionButtonBounds()
	if !ok || rect.Right <= rect.Left || rect.Bottom <= rect.Top {
		return false
	}

	clientX, clientY, clientOK := w32.ScreenToClient(w.hwnd, screenX, screenY)
	if clientOK &&
		clientX >= int(rect.Left) && clientX < int(rect.Right) &&
		clientY >= int(rect.Top) && clientY < int(rect.Bottom) {
		return true
	}

	windowRect := w32.GetWindowRect(w.hwnd)
	windowX := screenX - int(windowRect.Left)
	windowY := screenY - int(windowRect.Top)
	return windowX >= int(rect.Left) && windowX < int(rect.Right) &&
		windowY >= int(rect.Top) && windowY < int(rect.Bottom)
}

func (w *windowsWebviewWindow) routeCompositionMouseInput(msg uint32, wparam, lparam uintptr) bool {
	if !w.canSendCompositionMouseInput() {
		return false
	}

	// Caption-button non-client hover must remain visible to DefWindowProc/DWM.
	// Custom frontend buttons mirror it through handleSnapLayoutHoverMessage().
	if msg == w32.WM_NCMOUSEMOVE && frontendCaptionButtonHitTest(wparam) {
		return false
	}

	eventKind, ok := webviewCompositionMouseEventKind(msg)
	if !ok {
		return false
	}

	clientX, clientY, ok := w.compositionMouseClientPoint(msg, lparam)
	if !ok {
		return false
	}

	bounds := w.chromium.Bounds()
	mouseInWebView := msg == w32.WM_MOUSELEAVE ||
		w.compositionInput.capturing ||
		pointInEdgeRect(bounds, int32(clientX), int32(clientY))
	if !mouseInWebView {
		if msg == w32.WM_MOUSEMOVE && w.compositionInput.tracking {
			w.compositionInput.tracking = false
			w.cancelCompositionMouseLeaveTracking()
			_ = w.sendCompositionMouseLeave()
		}
		return false
	}

	w.updateCompositionMouseTracking(msg)
	w.updateCompositionMouseCapture(msg)

	mouseData := compositionMouseData(msg, wparam)
	if msg == w32.WM_MOUSELEAVE {
		wparam = 0
	}

	pointX, pointY := clientPointToWebViewPoint(bounds, msg, clientX, clientY)
	if err := w.chromium.SendMouseInput(
		eventKind,
		edge.COREWEBVIEW2_MOUSE_EVENT_VIRTUAL_KEYS(uint32(wparam)&0xffff),
		mouseData,
		pointX,
		pointY,
	); err != nil {
		return false
	}
	return true
}

func (w *windowsWebviewWindow) canSendCompositionMouseInput() bool {
	return w.chromium != nil &&
		w.parent.options.Windows.WebView2CompositionHosting &&
		w.chromium.CompositionControllerReady()
}

func (w *windowsWebviewWindow) compositionMouseClientPoint(msg uint32, lparam uintptr) (int, int, bool) {
	if msg == w32.WM_MOUSELEAVE {
		return 0, 0, true
	}

	if mouseMessageUsesScreenCoordinates(msg) {
		return w32.ScreenToClient(w.hwnd, int(w32.GET_X_LPARAM(lparam)), int(w32.GET_Y_LPARAM(lparam)))
	}

	return int(w32.GET_X_LPARAM(lparam)), int(w32.GET_Y_LPARAM(lparam)), true
}

func mouseMessageUsesScreenCoordinates(msg uint32) bool {
	switch msg {
	case w32.WM_MOUSEWHEEL, w32.WM_MOUSEHWHEEL,
		w32.WM_NCMOUSEMOVE,
		w32.WM_NCLBUTTONDOWN, w32.WM_NCLBUTTONUP, w32.WM_NCLBUTTONDBLCLK,
		w32.WM_NCRBUTTONDOWN, w32.WM_NCRBUTTONUP, w32.WM_NCRBUTTONDBLCLK,
		w32.WM_NCMBUTTONDOWN, w32.WM_NCMBUTTONUP, w32.WM_NCMBUTTONDBLCLK,
		w32.WM_NCXBUTTONDOWN, w32.WM_NCXBUTTONUP, w32.WM_NCXBUTTONDBLCLK:
		return true
	default:
		return false
	}
}

func (w *windowsWebviewWindow) updateCompositionMouseTracking(msg uint32) {
	if msg == w32.WM_MOUSEMOVE && !w.compositionInput.tracking {
		w.compositionInput.tracking = true
		w32.TrackMouseEvent(&w32.TRACKMOUSEEVENT{
			CbSize:    uint32(unsafe.Sizeof(w32.TRACKMOUSEEVENT{})),
			DwFlags:   w32.TME_LEAVE,
			HwndTrack: w.hwnd,
		})
	}
	if msg == w32.WM_MOUSELEAVE {
		w.compositionInput.tracking = false
	}
}

func (w *windowsWebviewWindow) cancelCompositionMouseLeaveTracking() {
	w32.TrackMouseEvent(&w32.TRACKMOUSEEVENT{
		CbSize:    uint32(unsafe.Sizeof(w32.TRACKMOUSEEVENT{})),
		DwFlags:   w32.TME_LEAVE | w32.TME_CANCEL,
		HwndTrack: w.hwnd,
	})
}

func (w *windowsWebviewWindow) updateCompositionMouseCapture(msg uint32) {
	switch {
	case compositionMouseButtonDownMessage(msg):
		w.focus()
		if !w.compositionInput.capturing {
			w.compositionInput.capturing = true
			w32.SetCapture(w.hwnd)
		}
	case compositionMouseButtonUpMessage(msg):
		if w.compositionInput.capturing {
			w.compositionInput.capturing = false
			w32.ReleaseCapture()
		}
	}
}

func compositionMouseButtonDownMessage(msg uint32) bool {
	switch msg {
	case w32.WM_LBUTTONDOWN, w32.WM_RBUTTONDOWN, w32.WM_MBUTTONDOWN, w32.WM_XBUTTONDOWN:
		return true
	default:
		return false
	}
}

func compositionMouseButtonUpMessage(msg uint32) bool {
	switch msg {
	case w32.WM_LBUTTONUP, w32.WM_RBUTTONUP, w32.WM_MBUTTONUP, w32.WM_XBUTTONUP:
		return true
	default:
		return false
	}
}

func compositionMouseData(msg uint32, wparam uintptr) uint32 {
	switch msg {
	case w32.WM_MOUSEWHEEL, w32.WM_MOUSEHWHEEL:
		return uint32(int32(int16(wparam >> 16)))
	case w32.WM_XBUTTONDOWN, w32.WM_XBUTTONUP, w32.WM_XBUTTONDBLCLK:
		return uint32(wparam >> 16)
	default:
		return 0
	}
}

func clientPointToWebViewPoint(bounds *edge.Rect, msg uint32, clientX, clientY int) (int32, int32) {
	if msg == w32.WM_MOUSELEAVE {
		return 0, 0
	}

	pointX := int32(clientX)
	pointY := int32(clientY)
	if bounds != nil {
		pointX -= bounds.Left
		pointY -= bounds.Top
	}
	return pointX, pointY
}

func (w *windowsWebviewWindow) sendCompositionMouseInputAtClientPoint(eventKind edge.COREWEBVIEW2_MOUSE_EVENT_KIND, virtualKeys edge.COREWEBVIEW2_MOUSE_EVENT_VIRTUAL_KEYS, mouseData uint32, clientX, clientY int) error {
	pointX, pointY := clientPointToWebViewPoint(w.chromium.Bounds(), 0, clientX, clientY)
	return w.chromium.SendMouseInput(eventKind, virtualKeys, mouseData, pointX, pointY)
}

func (w *windowsWebviewWindow) sendCompositionMouseLeave() error {
	return w.chromium.SendMouseInput(
		edge.COREWEBVIEW2_MOUSE_EVENT_KIND_LEAVE,
		edge.COREWEBVIEW2_MOUSE_EVENT_VIRTUAL_KEYS_NONE,
		0,
		0,
		0,
	)
}

func pointInEdgeRect(rect *edge.Rect, x, y int32) bool {
	if rect == nil {
		return true
	}
	return x >= rect.Left && x < rect.Right && y >= rect.Top && y < rect.Bottom
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
