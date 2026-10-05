//go:build windows

package edge

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
	"unsafe"
	"weak"

	"github.com/wailsapp/wails/v3/internal/webview2/internal/w32"
	"golang.org/x/sys/windows"
)

// These tests drive the installed WebView2 runtime, so they must run on a
// Windows desktop. Each locks its goroutine to the thread whose apartment
// creates the controller.
func lockWebView2Thread(t *testing.T) {
	runtime.LockOSThread()
	if err := windows.CoInitializeEx(0, windows.COINIT_APARTMENTTHREADED); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		windows.CoUninitialize()
		runtime.UnlockOSThread()
	})
}

func newLiveTestWindow(t *testing.T) uintptr {
	var instance windows.Handle
	if err := windows.GetModuleHandleEx(0, nil, &instance); err != nil {
		t.Fatal(err)
	}
	className, _ := windows.UTF16PtrFromString("WebView2LifetimeTest")
	wc := w32.WndClassExW{
		CbSize:        uint32(unsafe.Sizeof(w32.WndClassExW{})),
		HInstance:     instance,
		LpszClassName: className,
		LpfnWndProc:   windows.NewCallback(w32.DefWindowProc),
	}
	_, _, _ = w32.User32RegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))
	hwnd, _, err := w32.User32CreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(className)), 0,
		0xCF0000, // WS_OVERLAPPEDWINDOW
		uintptr(w32.CW_USEDEFAULT), uintptr(w32.CW_USEDEFAULT), 640, 480,
		0, 0, uintptr(instance), 0,
	)
	if hwnd == 0 {
		t.Fatal(err)
	}
	t.Cleanup(func() { w32.DestroyWindow(hwnd) })
	return hwnd
}

// pumpUntil dispatches this thread's messages until done reports true.
func pumpUntil(timeout time.Duration, done func() bool) bool {
	var msg w32.Msg
	deadline := time.Now().Add(timeout)
	for !done() {
		if time.Now().After(deadline) {
			return false
		}
		for {
			got, _, _ := w32.User32PeekMessageW.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0, w32.PM_REMOVE)
			if got == 0 {
				break
			}
			w32.User32TranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
			w32.User32DispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
		}
		time.Sleep(5 * time.Millisecond)
	}
	return true
}

func embedLive(t *testing.T, hwnd uintptr) *Chromium {
	e := NewChromium()
	// The browser process can outlive the test and lock its user data folder,
	// so share one outside t.TempDir.
	e.DataPath = filepath.Join(os.TempDir(), "wails-webview2-lifetime-test")
	e.SetErrorCallback(func(err error) { t.Errorf("fatal callback: %v", err) })
	if err := e.EmbedWithError(hwnd); err != nil {
		t.Fatalf("EmbedWithError: %v", err)
	}
	return e
}

func eventHandlers(e *Chromium) map[string]unsafe.Pointer {
	return map[string]unsafe.Pointer{
		"controllerCompleted":              unsafe.Pointer(e.controllerCompleted),
		"webMessageReceived":               unsafe.Pointer(e.webMessageReceived),
		"permissionRequested":              unsafe.Pointer(e.permissionRequested),
		"webResourceRequested":             unsafe.Pointer(e.webResourceRequested),
		"acceleratorKeyPressed":            unsafe.Pointer(e.acceleratorKeyPressed),
		"navigationStarting":               unsafe.Pointer(e.navigationStarting),
		"navigationCompleted":              unsafe.Pointer(e.navigationCompleted),
		"processFailed":                    unsafe.Pointer(e.processFailed),
		"containsFullScreenElementChanged": unsafe.Pointer(e.containsFullScreenElementChanged),
	}
}

func handlerAddrs(e *Chromium) map[string]comAddr {
	addrs := map[string]comAddr{}
	for name, p := range eventHandlers(e) {
		addrs[name] = comAddr(uintptr(p))
	}
	return addrs
}

func liveRefs(addrs map[string]comAddr) (total uint32, held map[string]uint32) {
	held = map[string]uint32{}
	for name, addr := range addrs {
		if n := handlerRefCount(unsafe.Pointer(uintptr(addr))); n != 0 {
			held[name] = n
			total += n
		}
	}
	return total, held
}

func checkNoStrayReferences(t *testing.T, queries, releases uint64) {
	t.Helper()
	if n := handlerQueryInterfaces.Load() - queries; n != 0 {
		t.Errorf("WebView2 queried Go handlers %d times; QueryInterface must be implemented", n)
	}
	if n := unbalancedHandlerReleases.Load() - releases; n != 0 {
		t.Errorf("WebView2 made %d Release calls without a matching AddRef", n)
	}
}

// startAbandonedController requests a controller for a Chromium, then abandons
// it before the completion is delivered, as the embed timeout does.
//
//go:noinline
func startAbandonedController(t *testing.T, env *ICoreWebView2Environment, hwnd uintptr) (comAddr, weak.Pointer[Chromium]) {
	e := NewChromium()
	e.hwnd = hwnd
	if err := env.CreateCoreWebView2Controller(hwnd, e.controllerCompleted); err != nil {
		t.Fatalf("CreateCoreWebView2Controller: %v", err)
	}
	handler := unsafe.Pointer(e.controllerCompleted)
	if handlerRefCount(handler) == 0 {
		t.Fatal("WebView2 did not AddRef the pending completion handler")
	}
	e.ShuttingDown()
	return comAddr(uintptr(handler)), weak.Make(e)
}

func TestLiveDelayedControllerCallbackAfterAbandon(t *testing.T) {
	lockWebView2Thread(t)
	queries, releases := handlerQueryInterfaces.Load(), unbalancedHandlerReleases.Load()
	host := embedLive(t, newLiveTestWindow(t))
	defer host.Close()

	handler, instance := startAbandonedController(t, host.environment, newLiveTestWindow(t))
	collectGarbage() // No messages are pumped, so the completion is still pending.
	if instance.Value() == nil {
		t.Fatal("abandoned Chromium was collected while its completion was pending")
	}
	if !pumpUntil(30*time.Second, func() bool { return handlerRefCount(unsafe.Pointer(uintptr(handler))) == 0 }) {
		t.Fatal("WebView2 did not deliver and release the late completion")
	}
	if !waitCollected(instance) {
		t.Fatal("abandoned Chromium was retained after WebView2 released its completion")
	}
	checkNoStrayReferences(t, queries, releases)
}

// embedAndClose embeds a Chromium whose callbacks reference a window, then
// closes it the way a failed or torn down instance is closed.
//
//go:noinline
func embedAndClose(t *testing.T, hwnd uintptr) (map[string]comAddr, weak.Pointer[Chromium], weak.Pointer[hostWindow]) {
	w := &hostWindow{}
	e := embedLive(t, hwnd)
	e.NavigationCompletedCallback = func(*ICoreWebView2, *ICoreWebView2NavigationCompletedEventArgs) { w.events++ }
	e.ProcessFailedCallback = func(*ICoreWebView2, *ICoreWebView2ProcessFailedEventArgs) { w.events++ }
	addrs := handlerAddrs(e)
	if total, held := liveRefs(addrs); total == 0 {
		t.Fatal("WebView2 holds no references to the registered event handlers")
	} else {
		t.Logf("native references while embedded: %v", held)
	}
	e.Close()
	return addrs, weak.Make(e), weak.Make(w)
}

func TestLiveClosedChromiumReleasesHandlers(t *testing.T) {
	lockWebView2Thread(t)
	queries, releases := handlerQueryInterfaces.Load(), unbalancedHandlerReleases.Load()
	addrs, instance, window := embedAndClose(t, newLiveTestWindow(t))
	if !waitCollected(window) {
		t.Fatal("closed Chromium retained its window")
	}
	if !pumpUntil(30*time.Second, func() bool { total, _ := liveRefs(addrs); return total == 0 }) {
		_, held := liveRefs(addrs)
		t.Fatalf("WebView2 kept handler references after Close: %v", held)
	}
	if !waitCollected(instance) {
		t.Fatal("closed Chromium was retained after WebView2 released its handlers")
	}
	checkNoStrayReferences(t, queries, releases)
}

// replaceChromium mirrors rebuildWebView: close the abandoned controller
// before embedding its replacement into the same window.
//
//go:noinline
func replaceChromium(t *testing.T, hwnd uintptr) (*Chromium, map[string]comAddr, weak.Pointer[Chromium], weak.Pointer[hostWindow]) {
	w := &hostWindow{}
	previous := embedLive(t, hwnd)
	previous.NavigationCompletedCallback = func(*ICoreWebView2, *ICoreWebView2NavigationCompletedEventArgs) { w.events++ }
	previous.ProcessFailedCallback = func(*ICoreWebView2, *ICoreWebView2ProcessFailedEventArgs) { w.events++ }
	addrs := handlerAddrs(previous)
	previous.Close()
	return embedLive(t, hwnd), addrs, weak.Make(previous), weak.Make(w)
}

func TestLiveReplacementReleasesAbandonedHandlers(t *testing.T) {
	lockWebView2Thread(t)
	queries, releases := handlerQueryInterfaces.Load(), unbalancedHandlerReleases.Load()
	replacement, addrs, previous, window := replaceChromium(t, newLiveTestWindow(t))
	defer replacement.Close()

	navigated := false
	replacement.NavigationCompletedCallback = func(*ICoreWebView2, *ICoreWebView2NavigationCompletedEventArgs) { navigated = true }
	replacement.NavigateToString("<p>replacement</p>")
	if !pumpUntil(30*time.Second, func() bool { return navigated }) {
		t.Fatal("replacement did not complete navigation")
	}
	if !waitCollected(window) {
		t.Fatal("abandoned Chromium retained its window")
	}
	if !pumpUntil(30*time.Second, func() bool { total, _ := liveRefs(addrs); return total == 0 }) {
		_, held := liveRefs(addrs)
		t.Fatalf("replacement retained abandoned controller handlers: %v", held)
	}
	if !waitCollected(previous) {
		t.Fatal("abandoned Chromium was retained after replacement")
	}
	checkNoStrayReferences(t, queries, releases)
}

// embedAndDestroyWindow mirrors the application's WM_CLOSE path, which marks
// the Chromium as shutting down and destroys its window without closing the
// controller.
//
//go:noinline
func embedAndDestroyWindow(t *testing.T, hwnd uintptr) (map[string]comAddr, weak.Pointer[Chromium]) {
	e := embedLive(t, hwnd)
	addrs := handlerAddrs(e)
	e.ShuttingDown()
	w32.DestroyWindow(hwnd)
	return addrs, weak.Make(e)
}

func TestLiveDestroyedWindowReleasesHandlers(t *testing.T) {
	lockWebView2Thread(t)
	queries, releases := handlerQueryInterfaces.Load(), unbalancedHandlerReleases.Load()
	addrs, instance := embedAndDestroyWindow(t, newLiveTestWindow(t))
	if !pumpUntil(30*time.Second, func() bool { total, _ := liveRefs(addrs); return total == 0 }) {
		_, held := liveRefs(addrs)
		t.Fatalf("WebView2 kept handler references after its window was destroyed: %v", held)
	}
	if !waitCollected(instance) {
		t.Fatal("Chromium was retained after its window was destroyed")
	}
	checkNoStrayReferences(t, queries, releases)
}
