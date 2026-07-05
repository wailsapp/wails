//go:build windows

// Command minimal is the smallest real GUI application built on the generated
// WebView2 v2 bindings: it opens a top-level Win32 window, creates a CoreWebView2
// environment and controller, and navigates to a URL — all driven through the
// generated pkg/webview2 API.
//
// It exists to be *seen*. Run it on an interactive desktop (the physical console
// or an RDP session) and a window appears with the page rendered inside it.
//
//	go run .                       # opens https://wails.io in a window
//	go run . -url https://go.dev   # opens a different URL
//	go run . -screenshot out.png   # headless: load, write a PNG, exit
//
// IMPORTANT: a GUI window created inside an SSH session lives on a
// non-interactive window station and will not appear on anyone's desktop. Use
// -screenshot to get a visual artifact from a non-interactive session, or run
// it from an RDP/console session to see the window live.
package main

import (
	"flag"
	"log"
	"runtime"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"

	webview2 "github.com/wailsapp/wails/v3/internal/webview2/pkg/webview2"
	"github.com/wailsapp/wails/v3/internal/webview2/webviewloader"
)

// Application state. The WebView2 callbacks fire on the UI thread inside the
// message loop, so plain globals are safe here.
var (
	targetURL      string
	screenshotPath string
	mainHWND       uintptr
	captured       bool

	// Live COM objects kept alive for the lifetime of the window. The loader
	// and CreateXxx calls do not guarantee to hold these for us.
	keepEnv        *webview2.ICoreWebView2Environment
	keepController *webview2.ICoreWebView2Controller
	keepWebView    *webview2.ICoreWebView2
	captureStream  *webview2.IStream
)

func main() {
	flag.StringVar(&targetURL, "url", "https://wails.io", "URL to load")
	flag.StringVar(&screenshotPath, "screenshot", "", "if set, capture a PNG to this path after load and exit (no visible window needed)")
	flag.Parse()

	// WebView2 delivers all completion/event callbacks on this thread's COM
	// message queue, so pin the goroutine and drive a single-threaded apartment.
	runtime.LockOSThread()
	const coinitApartmentThreaded = 0x2
	if err := windows.CoInitializeEx(0, coinitApartmentThreaded); err != nil {
		log.Printf("CoInitializeEx: %v (continuing)", err)
	}
	defer windows.CoUninitialize()

	mainHWND = createWindow("WebView2 v2 generated bindings", 1024, 768)
	if screenshotPath == "" {
		showWindow(mainHWND)
	}

	// Create the environment via the loader (it handles locating + loading the
	// runtime DLL); everything after this is the generated bindings.
	if err := webviewloader.CreateCoreWebView2Environment(&envHandler{}); err != nil {
		log.Fatalf("CreateCoreWebView2Environment: %v", err)
	}

	log.Printf("loading %s ...", targetURL)
	messageLoop()
	log.Printf("done")
}

// ---- WebView2 callback handlers ----------------------------------------------

// unknownStub provides the IUnknown methods every callback handler needs. The
// AddRef/Release→1, QueryInterface→E_NOINTERFACE pattern is the standard
// approach for Go-implemented WebView2 callback objects.
type unknownStub struct{}

func (unknownStub) QueryInterface(refiid, object uintptr) uintptr { return 0x80004002 } // E_NOINTERFACE
func (unknownStub) AddRef() uint32                                { return 1 }
func (unknownStub) Release() uint32                               { return 1 }

// envHandler receives the created environment. It implements the loader's
// handler interface (combridge-based, no IUnknown needed).
type envHandler struct{}

func (h *envHandler) EnvironmentCompleted(errorCode webviewloader.HRESULT, env *webviewloader.ICoreWebView2Environment) webviewloader.HRESULT {
	if errorCode != 0 || env == nil {
		log.Fatalf("environment creation failed: HRESULT 0x%08x", uint32(errorCode))
	}
	// The loader hands back an already-QI'd ICoreWebView2Environment; it shares
	// the single-vtable-pointer layout of the generated type, so reinterpret it
	// to drive the generated bindings.
	keepEnv = (*webview2.ICoreWebView2Environment)(unsafe.Pointer(env))
	keepEnv.AddRef() // survive the loader's Release after this returns

	if err := keepEnv.CreateCoreWebView2Controller(webview2.HWND(mainHWND), webview2.NewICoreWebView2CreateCoreWebView2ControllerCompletedHandler(&controllerHandler{})); err != nil {
		log.Fatalf("CreateCoreWebView2Controller: %v", err)
	}
	return 0
}

// controllerHandler receives the created controller, attaches the WebView to the
// window, wires a navigation-completed handler, and starts navigation.
type controllerHandler struct{ unknownStub }

func (h *controllerHandler) CreateCoreWebView2ControllerCompleted(errorCode uintptr, controller *webview2.ICoreWebView2Controller) uintptr {
	if errorCode != 0 || controller == nil {
		log.Fatalf("controller creation failed: HRESULT 0x%08x", uint32(errorCode))
	}
	keepController = controller
	keepController.AddRef()

	wv, err := keepController.GetCoreWebView2()
	if err != nil {
		log.Fatalf("GetCoreWebView2: %v", err)
	}
	keepWebView = wv
	keepWebView.AddRef()

	// Fill the window with the WebView.
	var rc webview2.RECT
	getClientRect(mainHWND, &rc)
	if err := keepController.PutBounds(rc); err != nil {
		log.Fatalf("PutBounds: %v", err)
	}
	if err := keepController.PutIsVisible(true); err != nil {
		log.Fatalf("PutIsVisible: %v", err)
	}

	if _, err := keepWebView.AddNavigationCompleted(webview2.NewICoreWebView2NavigationCompletedEventHandler(&navHandler{})); err != nil {
		log.Fatalf("AddNavigationCompleted: %v", err)
	}
	if err := keepWebView.Navigate(targetURL); err != nil {
		log.Fatalf("Navigate: %v", err)
	}
	return 0
}

// navHandler logs the navigation result (a non-visual proof the page actually
// loaded) and, in screenshot mode, triggers a capture.
type navHandler struct{ unknownStub }

func (h *navHandler) NavigationCompleted(sender *webview2.ICoreWebView2, args *webview2.ICoreWebView2NavigationCompletedEventArgs) uintptr {
	ok, _ := args.GetIsSuccess()
	log.Printf("navigation completed: success=%v", ok)

	if screenshotPath == "" || captured {
		return 0
	}
	captured = true
	capturePreview(sender, screenshotPath)
	return 0
}

// captureHandler writes the capture to disk (by releasing the file stream) and
// quits the message loop.
type captureHandler struct{ unknownStub }

func (h *captureHandler) CapturePreviewCompleted(errorCode uintptr) uintptr {
	if captureStream != nil {
		_ = captureStream.Release() // flush + close the file
	}
	if errorCode != 0 {
		log.Printf("capture failed: HRESULT 0x%08x", uint32(errorCode))
	} else {
		log.Printf("wrote screenshot to %s", screenshotPath)
	}
	postQuitMessage(0)
	return 0
}

// capturePreview captures the rendered page to a PNG file via CapturePreview,
// backed by a file IStream so the bytes land on disk directly.
func capturePreview(wv *webview2.ICoreWebView2, path string) {
	stream, err := createFileStream(path)
	if err != nil {
		log.Fatalf("createFileStream: %v", err)
	}
	captureStream = stream
	if err := wv.CapturePreview(webview2.COREWEBVIEW2_CAPTURE_PREVIEW_IMAGE_FORMAT_PNG, stream, webview2.NewICoreWebView2CapturePreviewCompletedHandler(&captureHandler{})); err != nil {
		log.Fatalf("CapturePreview: %v", err)
	}
}

// ---- Win32 plumbing ----------------------------------------------------------

var (
	user32   = windows.NewLazySystemDLL("user32.dll")
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")
	shlwapi  = windows.NewLazySystemDLL("shlwapi.dll")

	pRegisterClassExW   = user32.NewProc("RegisterClassExW")
	pCreateWindowExW    = user32.NewProc("CreateWindowExW")
	pDefWindowProcW     = user32.NewProc("DefWindowProcW")
	pShowWindow         = user32.NewProc("ShowWindow")
	pUpdateWindow       = user32.NewProc("UpdateWindow")
	pGetMessageW        = user32.NewProc("GetMessageW")
	pTranslateMessage   = user32.NewProc("TranslateMessage")
	pDispatchMessageW   = user32.NewProc("DispatchMessageW")
	pPostQuitMessage    = user32.NewProc("PostQuitMessage")
	pGetClientRect      = user32.NewProc("GetClientRect")
	pGetModuleHandleW   = kernel32.NewProc("GetModuleHandleW")
	pSHCreateStreamFile = shlwapi.NewProc("SHCreateStreamOnFileEx")
)

type wndClassExW struct {
	cbSize        uint32
	style         uint32
	lpfnWndProc   uintptr
	cbClsExtra    int32
	cbWndExtra    int32
	hInstance     uintptr
	hIcon         uintptr
	hCursor       uintptr
	hbrBackground uintptr
	lpszMenuName  *uint16
	lpszClassName *uint16
	hIconSm       uintptr
}

type msgW struct {
	hwnd    uintptr
	message uint32
	wParam  uintptr
	lParam  uintptr
	time    uint32
	pt      struct{ x, y int32 }
}

const (
	wmDestroy          = 0x0002
	wsOverlappedWindow = 0x00CF0000
	cwUseDefault       = ^uintptr(0x7FFFFFFF) // 0x80000000
	swShow             = 5
)

func wndProc(hwnd, msg, wParam, lParam uintptr) uintptr {
	if msg == wmDestroy {
		postQuitMessage(0)
		return 0
	}
	r, _, _ := pDefWindowProcW.Call(hwnd, msg, wParam, lParam)
	return r
}

func createWindow(title string, width, height int) uintptr {
	hInstance, _, _ := pGetModuleHandleW.Call(0)
	className := windows.StringToUTF16Ptr("WebView2MinimalExample")

	wc := wndClassExW{
		lpfnWndProc:   syscall.NewCallback(wndProc),
		hInstance:     hInstance,
		lpszClassName: className,
		hbrBackground: 6, // COLOR_WINDOW+1
	}
	wc.cbSize = uint32(unsafe.Sizeof(wc))
	if atom, _, err := pRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc))); atom == 0 {
		log.Fatalf("RegisterClassExW: %v", err)
	}

	hwnd, _, err := pCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(title))),
		wsOverlappedWindow,
		cwUseDefault, cwUseDefault, uintptr(width), uintptr(height),
		0, 0, hInstance, 0,
	)
	if hwnd == 0 {
		log.Fatalf("CreateWindowExW: %v", err)
	}
	return hwnd
}

func showWindow(hwnd uintptr) {
	pShowWindow.Call(hwnd, swShow)
	pUpdateWindow.Call(hwnd)
}

func getClientRect(hwnd uintptr, rc *webview2.RECT) {
	pGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(rc)))
}

func postQuitMessage(code int) {
	pPostQuitMessage.Call(uintptr(code))
}

func messageLoop() {
	var msg msgW
	for {
		r, _, _ := pGetMessageW.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if int32(r) <= 0 { // 0 = WM_QUIT, -1 = error
			return
		}
		pTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		pDispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
	}
}

// createFileStream returns a writable file-backed IStream for CapturePreview.
func createFileStream(path string) (*webview2.IStream, error) {
	const (
		stgmReadWrite = 0x00000002
		stgmCreate    = 0x00001000
		fileAttrNorm  = 0x00000080
	)
	var stream uintptr
	hr, _, _ := pSHCreateStreamFile.Call(
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(path))),
		stgmReadWrite|stgmCreate,
		fileAttrNorm,
		1, // fCreate
		0, // pstmTemplate
		uintptr(unsafe.Pointer(&stream)),
	)
	if hr != 0 {
		return nil, syscall.Errno(hr)
	}
	return (*webview2.IStream)(unsafe.Pointer(stream)), nil
}
