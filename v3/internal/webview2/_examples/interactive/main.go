//go:build windows

// Command interactive is a deliberately broad exercise of the *generated*
// WebView2 v2 bindings. Where the minimal example just opens a URL, this one
// touches one of each major kind of generated code so a single run proves the
// generator handles them all against the real runtime.
//
// How you know it worked: each exercised feature is recorded as a named check,
// and the program prints a PASS/FAIL checklist plus a single RESULT line when
// the run completes (on the Go<->JS round trip, or an 8s fallback, or when you
// close the window). RESULT: PASS means every generated-code kind worked.
//
// Run it on an interactive desktop (console/RDP), not over SSH.
//
//	go run .                      # opens a window; drives the Go<->JS bridge
//	go run . -screenshot out.png  # run the round trip, snap a PNG, exit
package main

import (
	"flag"
	"log"
	"runtime"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"

	webview2 "github.com/wailsapp/wails/v3/internal/webview2/pkg/webview2"
	"github.com/wailsapp/wails/v3/internal/webview2/webviewloader"
)

var (
	screenshotPath string
	mainHWND       uintptr
	captured       bool

	keepEnv        *webview2.ICoreWebView2Environment
	keepController *webview2.ICoreWebView2Controller
	keepWebView    *webview2.ICoreWebView2
	captureStream  *webview2.IStream
)

// ---- result checklist --------------------------------------------------------

// The named checks, in report order. Each maps to one kind of generated code.
const (
	cEnv      = "env created + reinterpret to generated type"
	cVer      = "GetBrowserVersionString  (sync retval string)"
	cQI       = "GetICoreWebView2Controller2  (versioned-interface QI)"
	cBg       = "PutDefaultBackgroundColor  (by-value struct arg)"
	cZoom     = "PutZoomFactor  (float64 setter)"
	cFocus    = "MoveFocus  (enum-by-value arg)"
	cSettings = "GetSettings + bool setters  (interface-returning getter)"
	cScript   = "AddScriptToExecuteOnDocumentCreated  (in-string + async handler)"
	cNav      = "NavigationCompleted  (Go-implemented event handler)"
	cExec     = "ExecuteScript document.title  (async retval string)"
	cJs2Go    = "WebMessageReceived JS->Go  (event + retval string)"
	cRound    = "Go->JS->title round-trip  (PostWebMessageAsString)"
)

var checkOrder = []string{cEnv, cVer, cQI, cBg, cZoom, cFocus, cSettings, cScript, cNav, cExec, cJs2Go, cRound}

type checkResult struct {
	seen   bool
	ok     bool
	detail string
}

var (
	results    = map[string]*checkResult{}
	summarized bool
)

func rec(name string, ok bool, detail string) {
	results[name] = &checkResult{seen: true, ok: ok, detail: detail}
	mark := "PASS"
	if !ok {
		mark = "FAIL"
	}
	if detail != "" {
		log.Printf("  [%s] %s — %s", mark, name, detail)
	} else {
		log.Printf("  [%s] %s", mark, name)
	}
	maybeFinish()
}

// maybeFinish prints the verdict once every check has been observed. Events do
// not fire in a fixed order (the injected-script message can round-trip before
// NavigationCompleted fires), so completion is "all checks seen", not any single
// event.
func maybeFinish() {
	for _, name := range checkOrder {
		if c := results[name]; c == nil || !c.seen {
			return
		}
	}
	finish("all checks observed")
}

func recErr(name string, err error, okDetail string) {
	if err != nil {
		rec(name, false, err.Error())
		return
	}
	rec(name, true, okDetail)
}

// printSummary prints the checklist and a single RESULT verdict, once.
func printSummary(reason string) {
	if summarized {
		return
	}
	summarized = true
	log.Printf("================ RESULT SUMMARY (%s) ================", reason)
	passed := 0
	for _, name := range checkOrder {
		c := results[name]
		switch {
		case c == nil || !c.seen:
			log.Printf("  [----] %s  (not observed)", name)
		case c.ok:
			passed++
			log.Printf("  [PASS] %s", name)
		default:
			log.Printf("  [FAIL] %s — %s", name, c.detail)
		}
	}
	total := len(checkOrder)
	log.Printf("====================================================")
	if passed == total {
		log.Printf("RESULT: PASS — %d/%d generated-code kinds exercised successfully", passed, total)
	} else {
		log.Printf("RESULT: FAIL — %d/%d passed (%d missing/failed)", passed, total, total-passed)
	}
	if screenshotPath == "" {
		log.Printf("(window stays open — close it to exit)")
	}
}

// finish records the terminal reason, prints the verdict, and — in screenshot
// mode — snaps the page then quits.
func finish(reason string) {
	printSummary(reason)
	killTimer()
	if screenshotPath != "" && !captured && keepWebView != nil {
		captured = true
		capturePreview(keepWebView, screenshotPath)
	}
}

// ---- page + injected bridge script ------------------------------------------

const page = `<!doctype html>
<html><head><meta charset="utf-8"><title>WebView2 Bridge Demo</title>
<style>
  body{font-family:Segoe UI,system-ui,sans-serif;margin:0;background:#0f172a;color:#e2e8f0}
  .wrap{max-width:760px;margin:0 auto;padding:48px 32px}
  h1{font-size:28px;margin:0 0 4px}
  .sub{color:#94a3b8;margin:0 0 28px}
  ul{line-height:1.9;color:#cbd5e1}
  code{background:#1e293b;padding:2px 6px;border-radius:4px;color:#f8fafc}
  #reply{margin-top:28px;padding:18px 20px;border-radius:10px;background:#1e293b;
    border-left:4px solid #ef4444;font-size:17px}
  .ok{color:#34d399;font-weight:600}
</style></head>
<body><div class="wrap">
  <h1>WebView2 v2 — Generated Bindings</h1>
  <p class="sub">Rendered by the real runtime, driven entirely by generated Go code.</p>
  <ul>
    <li>versioned-interface QI: <code>GetICoreWebView2Controller2</code></li>
    <li>by-value struct: <code>PutDefaultBackgroundColor</code></li>
    <li>interface getter + bool setters: <code>GetSettings</code></li>
    <li>injected script: <code>AddScriptToExecuteOnDocumentCreated</code></li>
    <li>events: NavigationCompleted, DocumentTitleChanged, WebMessageReceived</li>
    <li>async retval string: <code>ExecuteScript</code></li>
  </ul>
  <div id="reply">Waiting for a message from Go…</div>
</div></body></html>`

const injected = `
window.chrome.webview.addEventListener('message', e => {
  var box = document.getElementById('reply');
  if (box) { box.innerHTML = '<span class="ok">' + e.data + '</span>'; }
  document.title = 'Go↔JS round-trip complete';
});
window.chrome.webview.postMessage('hello from JS (injected script ran)');
`

func main() {
	flag.StringVar(&screenshotPath, "screenshot", "", "capture a PNG after the round trip and exit")
	flag.Parse()

	runtime.LockOSThread()
	const coinitApartmentThreaded = 0x2
	if err := windows.CoInitializeEx(0, coinitApartmentThreaded); err != nil {
		log.Printf("CoInitializeEx: %v (continuing)", err)
	}
	defer windows.CoUninitialize()

	mainHWND = createWindow("WebView2 v2 generated bindings — interactive", 900, 760)
	if screenshotPath == "" {
		showWindow(mainHWND)
	}
	log.Printf("exercising generated bindings — checklist follows, RESULT prints at the end")
	if err := webviewloader.CreateCoreWebView2Environment(&envHandler{}); err != nil {
		log.Fatalf("CreateCoreWebView2Environment: %v", err)
	}
	messageLoop()
	// If the window was closed before the round trip, still print a verdict.
	printSummary("window closed")
	log.Printf("done")
}

// ---- handlers ----------------------------------------------------------------

type unknownStub struct{}

func (unknownStub) QueryInterface(refiid, object uintptr) uintptr { return 0x80004002 } // E_NOINTERFACE
func (unknownStub) AddRef() uint32                                { return 1 }
func (unknownStub) Release() uint32                               { return 1 }

type envHandler struct{}

func (h *envHandler) EnvironmentCompleted(errorCode webviewloader.HRESULT, env *webviewloader.ICoreWebView2Environment) webviewloader.HRESULT {
	if errorCode != 0 || env == nil {
		rec(cEnv, false, "HRESULT "+hexHR(uint32(errorCode)))
		finish("environment creation failed")
		return 0
	}
	keepEnv = (*webview2.ICoreWebView2Environment)(unsafe.Pointer(env))
	keepEnv.AddRef()
	rec(cEnv, true, "")

	ver, err := keepEnv.GetBrowserVersionString()
	recErr(cVer, err, ver)

	if err := keepEnv.CreateCoreWebView2Controller(webview2.HWND(mainHWND), webview2.NewICoreWebView2CreateCoreWebView2ControllerCompletedHandler(&controllerHandler{})); err != nil {
		log.Fatalf("CreateCoreWebView2Controller: %v", err)
	}
	return 0
}

type controllerHandler struct{ unknownStub }

func (h *controllerHandler) CreateCoreWebView2ControllerCompleted(errorCode uintptr, controller *webview2.ICoreWebView2Controller) uintptr {
	if errorCode != 0 || controller == nil {
		log.Fatalf("controller creation failed: %s", hexHR(uint32(errorCode)))
	}
	keepController = controller
	keepController.AddRef()

	// Versioned-interface QueryInterface helper + by-value struct parameter.
	if c2, err := keepController.GetICoreWebView2Controller2(); err != nil {
		rec(cQI, false, err.Error())
		rec(cBg, false, "skipped: no ICoreWebView2Controller2")
	} else {
		rec(cQI, true, "")
		recErr(cBg, c2.PutDefaultBackgroundColor(webview2.COREWEBVIEW2_COLOR{A: 255, R: 15, G: 23, B: 42}), "color #0f172a")
	}

	// float64 setter + enum-by-value method.
	recErr(cZoom, keepController.PutZoomFactor(1.0), "zoom=1.0")
	recErr(cFocus, keepController.MoveFocus(webview2.COREWEBVIEW2_MOVE_FOCUS_REASON_PROGRAMMATIC), "reason=PROGRAMMATIC")

	wv, err := keepController.GetCoreWebView2()
	if err != nil {
		log.Fatalf("GetCoreWebView2: %v", err)
	}
	keepWebView = wv
	keepWebView.AddRef()

	// Interface-returning getter + several bool setters.
	if s, err := keepWebView.GetSettings(); err != nil {
		rec(cSettings, false, err.Error())
	} else {
		var firstErr error
		for _, set := range []func() error{
			func() error { return s.PutIsScriptEnabled(true) },
			func() error { return s.PutIsWebMessageEnabled(true) },
			func() error { return s.PutAreDevToolsEnabled(true) },
			func() error { return s.PutAreDefaultContextMenusEnabled(false) },
		} {
			if e := set(); e != nil && firstErr == nil {
				firstErr = e
			}
		}
		recErr(cSettings, firstErr, "script/webmessage/devtools on, context-menus off")
	}

	var rc webview2.RECT
	getClientRect(mainHWND, &rc)
	_ = keepController.PutBounds(rc)
	_ = keepController.PutIsVisible(true)

	// [in] string + async completion handler.
	recErr(cScript, keepWebView.AddScriptToExecuteOnDocumentCreated(injected, webview2.NewICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler(&addScriptHandler{})), "queued injected bridge script")

	// Go-implemented event handlers.
	_, _ = keepWebView.AddNavigationCompleted(webview2.NewICoreWebView2NavigationCompletedEventHandler(&navHandler{}))
	_, _ = keepWebView.AddDocumentTitleChanged(webview2.NewICoreWebView2DocumentTitleChangedEventHandler(&titleHandler{}))
	_, _ = keepWebView.AddWebMessageReceived(webview2.NewICoreWebView2WebMessageReceivedEventHandler(&webMsgHandler{}))

	if err := keepWebView.NavigateToString(page); err != nil {
		log.Fatalf("NavigateToString: %v", err)
	}
	return 0
}

type addScriptHandler struct{ unknownStub }

// The async id confirms the queued script call round-tripped through the DLL;
// the check itself was already recorded when the call was issued.
func (h *addScriptHandler) AddScriptToExecuteOnDocumentCreatedCompleted(errorCode uintptr, id string) uintptr {
	log.Printf("       (injected script accepted, id=%q)", id)
	return 0
}

type navHandler struct{ unknownStub }

func (h *navHandler) NavigationCompleted(sender *webview2.ICoreWebView2, args *webview2.ICoreWebView2NavigationCompletedEventArgs) uintptr {
	ok, _ := args.GetIsSuccess()
	rec(cNav, ok, "success="+boolStr(ok))

	// Start a fallback timer so a verdict always prints even if the bridge stalls.
	setTimer(8000)

	// async [out,retval] string: run JS, get the result back in the handler.
	_ = sender.ExecuteScript("document.title", webview2.NewICoreWebView2ExecuteScriptCompletedHandler(&execScriptHandler{}))

	// Kick the Go->JS side of the bridge.
	if err := sender.PostWebMessageAsString("Hello from Go — the generated bindings are talking to JS"); err != nil {
		log.Printf("[bridge] PostWebMessageAsString: %v", err)
	}
	return 0
}

type execScriptHandler struct{ unknownStub }

func (h *execScriptHandler) ExecuteScriptCompleted(errorCode uintptr, result string) uintptr {
	recErr(cExec, hrErr(uint32(errorCode)), "result="+result)
	return 0
}

type webMsgHandler struct{ unknownStub }

func (h *webMsgHandler) WebMessageReceived(sender *webview2.ICoreWebView2, args *webview2.ICoreWebView2WebMessageReceivedEventArgs) uintptr {
	msg, err := args.TryGetWebMessageAsString()
	recErr(cJs2Go, err, "JS said: "+msg)
	if err == nil {
		// Reply so JS can update the DOM + title (drives DocumentTitleChanged).
		_ = sender.PostWebMessageAsString("Go received: " + msg)
	}
	return 0
}

type titleHandler struct{ unknownStub }

func (h *titleHandler) DocumentTitleChanged(sender *webview2.ICoreWebView2, args *webview2.IUnknown) uintptr {
	title, _ := sender.GetDocumentTitle() // sync [out,retval] string
	// JS sets this title only after receiving Go's reply, so it marks a
	// completed round trip and is the natural point to finish.
	if strings.Contains(title, "round-trip") {
		rec(cRound, true, "JS set title to "+quote(title)+" after Go's reply")
	}
	return 0
}

type captureHandler struct{ unknownStub }

func (h *captureHandler) CapturePreviewCompleted(errorCode uintptr) uintptr {
	if captureStream != nil {
		_ = captureStream.Release()
	}
	if errorCode != 0 {
		log.Printf("[capture] failed: %s", hexHR(uint32(errorCode)))
	} else {
		log.Printf("[capture] wrote %s", screenshotPath)
	}
	postQuitMessage(0)
	return 0
}

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

// ---- small helpers -----------------------------------------------------------

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func quote(s string) string { return "\"" + s + "\"" }

func hexHR(hr uint32) string {
	const digits = "0123456789abcdef"
	b := []byte("0x00000000")
	for i := 0; i < 8; i++ {
		b[2+7-i] = digits[hr&0xf]
		hr >>= 4
	}
	return string(b)
}

func hrErr(hr uint32) error {
	if hr == 0 {
		return nil
	}
	return syscall.Errno(hr)
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
	pSetTimer           = user32.NewProc("SetTimer")
	pKillTimer          = user32.NewProc("KillTimer")
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
	wmTimer            = 0x0113
	wsOverlappedWindow = 0x00CF0000
	cwUseDefault       = ^uintptr(0x7FFFFFFF) // 0x80000000
	swShow             = 5
	timerID            = 1
)

func wndProc(hwnd, msg, wParam, lParam uintptr) uintptr {
	switch msg {
	case wmDestroy:
		postQuitMessage(0)
		return 0
	case wmTimer:
		finish("timeout — round-trip not observed within 8s")
		return 0
	}
	r, _, _ := pDefWindowProcW.Call(hwnd, msg, wParam, lParam)
	return r
}

func setTimer(ms int)  { pSetTimer.Call(mainHWND, timerID, uintptr(ms), 0) }
func killTimer()       { pKillTimer.Call(mainHWND, timerID) }
func postQuitMessage(code int) { pPostQuitMessage.Call(uintptr(code)) }

func createWindow(title string, width, height int) uintptr {
	hInstance, _, _ := pGetModuleHandleW.Call(0)
	className := windows.StringToUTF16Ptr("WebView2InteractiveExample")

	wc := wndClassExW{
		lpfnWndProc:   syscall.NewCallback(wndProc),
		hInstance:     hInstance,
		lpszClassName: className,
		hbrBackground: 6,
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

func messageLoop() {
	var msg msgW
	for {
		r, _, _ := pGetMessageW.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if int32(r) <= 0 {
			return
		}
		pTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		pDispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
	}
}

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
		1,
		0,
		uintptr(unsafe.Pointer(&stream)),
	)
	if hr != 0 {
		return nil, syscall.Errno(hr)
	}
	return (*webview2.IStream)(unsafe.Pointer(stream)), nil
}
