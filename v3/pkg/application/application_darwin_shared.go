//go:build darwin && !ios && !server

package application

import (
	"time"
	"unsafe"

	"github.com/wailsapp/wails/v3/internal/assetserver/webview"
	"github.com/wailsapp/wails/v3/internal/operatingsystem"
)

type macosApp struct {
	applicationMenu unsafe.Pointer
	parent          *App
}

func (a *App) logPlatformInfo() {
	info, err := operatingsystem.Info()
	if err != nil {
		a.error("error getting OS info: %w", err)
		return
	}

	a.info("Platform Info:", info.AsLogSlice()...)

}

func (m *macosApp) GetFlags(options Options) map[string]any {
	if options.Flags == nil {
		options.Flags = make(map[string]any)
	}
	return options.Flags
}

// clearWindowDragCache removes cached references for a window
func clearWindowDragCache(windowID uint) {
	windowImplCache.Delete(windowID)

	// Cancel any pending timer
	if throttleVal, ok := dragThrottle.Load(windowID); ok {
		if throttle, ok := throttleVal.(*dragThrottleState); ok {
			throttle.mu.Lock()
			if throttle.timer != nil {
				throttle.timer.Stop()
			}
			throttle.mu.Unlock()
		}
	}
	dragThrottle.Delete(windowID)
}

func getNativeApplication() *macosApp {
	return globalApplication.impl.(*macosApp)
}

// sendDragUpdate sends the actual drag update to JavaScript
func sendDragUpdate(winID uint, x, y int) {
	// Try cached implementation first
	var darwinImpl *macosWebviewWindow
	var needsExecJS bool

	if cached, found := windowImplCache.Load(winID); found {
		darwinImpl = cached.(*macosWebviewWindow)
		if darwinImpl != nil && darwinImpl.nsWindow != nil {
			needsExecJS = true
		} else {
			// Invalid cache entry, remove it
			windowImplCache.Delete(winID)
		}
	}

	if !needsExecJS {
		// Fallback to full lookup
		window, ok := globalApplication.Window.GetByID(winID)
		if !ok || window == nil {
			return
		}

		// Type assert to WebviewWindow
		webviewWindow, ok := window.(*WebviewWindow)
		if !ok || webviewWindow == nil {
			return
		}

		// Get implementation
		darwinImpl, ok = webviewWindow.impl.(*macosWebviewWindow)
		if !ok {
			return
		}

		// Cache for next time
		windowImplCache.Store(winID, darwinImpl)
		needsExecJS = true
	}

	if !needsExecJS || darwinImpl == nil {
		return
	}

	// Protect shared buffer access
	dragOverJSMutex.Lock()

	// Build JS string with zero allocations
	// Format: "window._wails.handleDragOver(X,Y)"
	// Max length with int32 coords: 30 + 11 + 1 + 11 + 1 + 1 = 55 bytes
	n := copy(dragOverJSBuffer[:], dragOverJSPrefix)
	n += writeInt(dragOverJSBuffer[n:], x)
	if n < len(dragOverJSBuffer) {
		dragOverJSBuffer[n] = ','
		n++
	}
	n += writeInt(dragOverJSBuffer[n:], y)
	if n < len(dragOverJSBuffer) {
		dragOverJSBuffer[n] = ')'
		n++
	}
	if n < len(dragOverJSBuffer) {
		dragOverJSBuffer[n] = 0 // null terminate for C
	} else {
		// Buffer overflow - this should not happen with 128 byte buffer
		dragOverJSMutex.Unlock()
		return
	}

	// Call JavaScript with zero allocations
	darwinImpl.execJSDragOver(dragOverJSBuffer[:n+1]) // Include null terminator
	dragOverJSMutex.Unlock()
}

// writeInt writes an integer to a byte slice and returns the number of bytes written
func writeInt(buf []byte, n int) int {
	if n < 0 {
		if len(buf) == 0 {
			return 0
		}
		buf[0] = '-'
		return 1 + writeInt(buf[1:], -n)
	}
	if n == 0 {
		if len(buf) == 0 {
			return 0
		}
		buf[0] = '0'
		return 1
	}

	// Count digits
	tmp := n
	digits := 0
	for tmp > 0 {
		digits++
		tmp /= 10
	}

	// Bounds check
	if digits > len(buf) {
		return 0
	}

	// Write digits in reverse
	for i := digits - 1; i >= 0; i-- {
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return digits
}

// throttleDragOver coalesces drag-over positions per window: it sends at most one
// update immediately (on a 5px move) and trails the rest with a 50ms timer.
func throttleDragOver(winID uint, x int, y int) {
	intX, intY := x, y

	// Get or create throttle state
	throttleKey := winID
	throttleVal, _ := dragThrottle.LoadOrStore(throttleKey, &dragThrottleState{
		lastX: intX,
		lastY: intY,
	})
	throttle := throttleVal.(*dragThrottleState)

	throttle.mu.Lock()

	// Update pending position
	throttle.pendingX = intX
	throttle.pendingY = intY
	throttle.hasPending = true

	// If timer is already running, just update the pending position
	if throttle.timer != nil {
		throttle.mu.Unlock()
		return
	}

	// Apply 5-pixel threshold for immediate update
	dx := intX - throttle.lastX
	dy := intY - throttle.lastY
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}

	// Check if we should send an immediate update
	shouldSendNow := dx >= 5 || dy >= 5

	if shouldSendNow {
		// Update last position
		throttle.lastX = intX
		throttle.lastY = intY
		throttle.hasPending = false

		// Send this update immediately (unlock before JS call to avoid deadlock)
		throttle.mu.Unlock()
		sendDragUpdate(winID, intX, intY)
		throttle.mu.Lock()
	}

	// Start 50ms timer for next update (whether we sent now or not)
	throttle.timer = time.AfterFunc(50*time.Millisecond, func() {
		// Execute on main thread to ensure UI updates
		InvokeSync(func() {
			throttle.mu.Lock()
			// Clear timer reference
			throttle.timer = nil

			// Send pending update if any
			if throttle.hasPending {
				pendingX, pendingY := throttle.pendingX, throttle.pendingY
				throttle.lastX = pendingX
				throttle.lastY = pendingY
				throttle.hasPending = false
				throttle.mu.Unlock()
				sendDragUpdate(winID, pendingX, pendingY)
			} else {
				throttle.mu.Unlock()
			}
		})
	})
	throttle.mu.Unlock()
}

// handleURLRequest queues a WKURLSchemeTask for the window's asset server.
func handleURLRequest(windowID uint, wkURLSchemeTask unsafe.Pointer) {
	window, ok := globalApplication.Window.GetByID(windowID)
	if !ok || window == nil {
		globalApplication.debug("could not find window with id", "windowID", windowID)
		return
	}

	webviewRequests <- &webViewAssetRequest{
		Request:    webview.NewRequest(wkURLSchemeTask),
		windowId:   windowID,
		windowName: window.Name(),
	}
}

// handleDragEnter tells the window's JavaScript that a drag entered it.
func handleDragEnter(windowID uint) {
	window, ok := globalApplication.Window.GetByID(windowID)
	if !ok || window == nil {
		return
	}

	// Call JavaScript to show drag entered state
	window.ExecJS("window._wails.handleDragEnter();")
}

// handleDragExit tells the window's JavaScript that a drag left it.
func handleDragExit(windowID uint) {
	window, ok := globalApplication.Window.GetByID(windowID)
	if !ok || window == nil {
		return
	}

	// Call JavaScript to clean up drag state
	window.ExecJS("window._wails.handleDragLeave();")
}
